# 04-fleet-nodes-and-ssh: Fleet Nodes Components, Vault CLI & Remote Clone Specification

- **Spec ID:** `04-fleet-nodes-and-ssh/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Fleet Nodes CLI, RSA Vault Implementation, Remote Clone Dispatcher, Ping Probe
- **Dependencies:** `cli/cmdnodes`, `cli/crypto`, `cli/cluster`, `cli/netip`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Fleet Nodes and SSH cluster encompasses four primary operational components:

```
04-fleet-nodes-and-ssh/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Fleet Nodes Controller** | `cli/cmdnodes/nodes.go`, `cli/cmdnodes/cluster.go` | Inventory management, multi-node grouping, parallel command dispatch. |
| **RSA Credential Vault** | `cli/crypto/vault.go`, `cli/cmdssh/askpass.go` | Asymmetric RSA keypair management, OAEP encryption, ephemeral askpass generation. |
| **Reachability Prober** | `cli/cluster/ping.go`, `cli/netip/probe.go` | Dual-stack ICMP and TCP ping probing with latency measurement. |
| **CFR Remote Cloner** | `cli/cluster/cfr.go`, `cli/cmdnodes/clone.go` | Cluster Fleet Remote manifest staging, except-self directory resolution. |

---

## 2. RSA Credential Vault Implementation

### 2.1 Cryptographic Parameters
- **Algorithm:** RSA 2048-bit with PKCS#8 ASN.1 DER encoding.
- **Encryption Scheme:** RSA-OAEP with SHA-256 hash and empty label.
- **Key Storage:**
  - Private key: `~/.gitmap/keys/vault_rsa` (`0600` file permission).
  - Public key: `~/.gitmap/keys/vault_rsa.pub` (`0644` file permission).

### 2.2 Ephemeral Askpass Script Generator
When invoking OpenSSH subprocesses, the vault creates a temporary shell script that prints the decrypted password to stdout and self-deletes upon completion:

```go
package cmdssh

import (
    "os"
    "path/filepath"
)

func CreateAskPassScript(password string) (string, func(), error) {
    // Writes temporary script printing password
    // Returns script path and cleanup closure
}
```

---

## 3. Unified Fleet Nodes CLI Commands

### 3.1 Fleet Nodes Command Suite
- `gitmap nodes`: Lists all configured cluster nodes with IP, OS, reachability status, and last sync timestamp.
- `gitmap nodes ping [node]`: Executes dual-stack reachability test (ICMP + TCP) and prints latency table.
- `gitmap nodes add <name> <ip> [flags]`: Adds a new machine to the cluster inventory.
- `gitmap nodes clone <repo> [flags]`: Dispatches repository cloning across all live cluster nodes (except self).
- `gitmap cluster exec <command>`: Executes an arbitrary shell command across all live cluster nodes concurrently.

---

## 4. CFR/CFRP Remote Fleet Staging

### 4.1 Manifest Format
The Cluster Fleet Remote manifest captures required repositories and branches for node deployment:

```json
{
  "version": "1.0.0",
  "sourceNode": "control-01",
  "timestamp": "2026-10-06T12:00:00Z",
  "repositories": [
    {
      "name": "gitmap",
      "remoteUrl": "git@github.com:alimtvnetwork/gitmap.git",
      "targetBranch": "main",
      "targetPath": "~/dev/gitmap"
    }
  ]
}
```

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isNodesListFormattedCorrectly: true
  isAskPassCleanedUpPostExecution: true
  isCfrManifestValidated: true
  isPositiveBooleansUsed: true
```

- [x] Askpass scripts guaranteed deleted even on command failure.
- [x] Nodes ping reports accurate round-trip latency.
- [x] Remote clone obeys except-self exclusion rules.
