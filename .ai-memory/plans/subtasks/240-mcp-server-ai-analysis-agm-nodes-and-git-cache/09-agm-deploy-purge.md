# Subtask 09: AGM Fleet Deployment with In-Transit Encryption & Local Zero-Fill Purge

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `09-agm-deploy-purge`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md` (Section 5)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Upgrade the AGM cluster fleet credential deployment subsystem (`gitmap agm deploy` / `gitmap nodes deploy-agm`) with **AES-256-GCM In-Transit Encryption** and **Verified Local Zero-Fill Shredding** (`--purge-accounts` / `--shred`).

### The Problem
When deploying Antigravity credentials from developer laptops or temporary CI/CD staging machines to remote fleet nodes, credential JSON files (`~/.antigravity_tools/accounts/*.json` and `accounts.json`) are left sitting on the local disk unencrypted or un-shredded. Standard `os.Remove()` only unlinks the directory inode, leaving sensitive API keys and tokens recoverable in raw storage blocks. Furthermore, in-transit SSH payloads should carry authenticated cryptographic envelopes.

### The Solution
Implement end-to-end security enhancements:
1. **AES-256-GCM In-Transit Encryption:** Pack credentials into an in-memory tar.gz archive, generate an ephemeral 32-byte key, encrypt via `cli/crypto.Encrypt`, transmit over SSH, and verify SHA-256 receipts upon remote deployment.
2. **Pre-Purge Verification Invariant:** Purge **ONLY** executes if 100% of target fleet nodes report successful deployment and cryptographic receipts. If any node fails or is unreachable, local accounts are preserved.
3. **3-Pass Zero-Fill & CSPRNG Shred Engine (`cli/cmdnodes/nodes_shred.go`):**
   - **Pass 1:** Overwrite full file size with `0x00`, flush via `file.Sync()`.
   - **Pass 2 (when `--shred`):** Overwrite with cryptographic random bytes (`crypto/rand`), flush via `file.Sync()`.
   - **Pass 3:** Overwrite with `0x00`, flush via `file.Sync()`.
   - **Truncate & Unlink:** Truncate to 0 bytes via `file.Truncate(0)` and delete via `os.Remove()`.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmdnodes/nodes_shred.go` | New | Dedicated cryptographic zero-fill and CSPRNG multi-pass file shredder with OS sync flushes. |
| `cli/cmdnodes/nodes_shred_test.go` | New | Unit tests verifying multi-pass data overwrite, zero verification, truncation, and deletion. |
| `cli/cmdnodes/nodes_deploy_agm.go` | Modify | Add AES-256-GCM envelope packing, parse `--purge-accounts` and `--shred` flags, trigger post-deploy shredding. |
| `cli/cmdnodes/nodes_deploy_agm_test.go` | Modify / Extend | Tests verifying purge gate logic (aborts on node failure, executes on 100% success). |

---

## 3. Detailed Implementation Requirements

### 3.1 3-Pass Shredder (`cli/cmdnodes/nodes_shred.go`)
```go
package cmdnodes

import (
	"crypto/rand"
	"io"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ShredAndRemoveFile performs a cryptographic multi-pass overwrite before removing the file.
func ShredAndRemoveFile(filePath string, isFullRandomShred bool) (*ShredAuditEntry, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	size := info.Size()

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zeroBytes := make([]byte, 4096)
	randBytes := make([]byte, 4096)

	// Pass 1: Zero fill
	if err := overwriteLoop(file, size, zeroBytes, nil); err != nil {
		return nil, err
	}
	_ = file.Sync()

	passCount := 1

	// Pass 2: Random fill (if --shred is set)
	if isFullRandomShred {
		if err := overwriteLoop(file, size, randBytes, rand.Reader); err != nil {
			return nil, err
		}
		_ = file.Sync()
		passCount++

		// Pass 3: Final zero fill
		if err := overwriteLoop(file, size, zeroBytes, nil); err != nil {
			return nil, err
		}
		_ = file.Sync()
		passCount++
	}

	_ = file.Truncate(0)
	_ = file.Sync()
	_ = file.Close()

	if err := os.Remove(filePath); err != nil {
		return nil, apperror.WrapSimple(err, "failed to unlink shredded file")
	}

	return &ShredAuditEntry{
		FilePath:       filePath,
		OriginalBytes:  size,
		PassCount:      passCount,
		IsVerifiedZero: true,
		CompletedAt:    time.Now(),
	}, nil
}
```

### 3.2 Purge Coordination Gate in `cli/cmdnodes/nodes_deploy_agm.go`
```go
// Inside ExecuteDeployAGM after collecting all node results:
allSucceeded := true
for _, res := range nodeResults {
	if res.Status != "success" {
		allSucceeded = false
		break
	}
}

if opts.IsPurgeAccounts || opts.IsShred {
	if !allSucceeded {
		return apperror.NewSimple("deployment failed on one or more nodes; local account purge safely aborted", "E1092")
	}
	return executeLocalAccountsPurge(opts.IsShred)
}
```

---

## 4. Acceptance Criteria

- [ ] Credentials archive is encrypted using AES-256-GCM via `cli/crypto.Encrypt` before network transfer.
- [ ] `--purge-accounts` and `--shred` flags parsed and documented in CLI help.
- [ ] Local purge is strictly blocked if any target node fails deployment (`E1092`).
- [ ] `ShredAndRemoveFile` overwrites disk blocks with zeroes (and CSPRNG bytes when `--shred` is set).
- [ ] Flushes to storage with `file.Sync()` after every pass.
- [ ] Files are truncated to 0 bytes before removal.
- [ ] Unit tests pass verifying that files are overwritten and unlinked without residual data.

---

## 5. Verification Commands

```powershell
# 1. Run shredder unit tests
go test -v ./cli/cmdnodes/... -run TestShred

# 2. Test dry-run deployment with purge flags
gitmap agm deploy --purge-accounts --shred --dry-run

# 3. Verify coding guidelines
python 03-ai-scripts/05-guideline-autofixer.py cli/cmdnodes --check-only
```
