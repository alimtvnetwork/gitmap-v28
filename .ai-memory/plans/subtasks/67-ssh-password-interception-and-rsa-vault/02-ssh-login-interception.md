# Subtask 02: SSH Password Interception, Interactive RSA Consent Prompt & AskPass Dispatch

> **Parent Plan:** `67-ssh-password-interception-and-rsa-vault.md`  
> **Status:** `PENDING`  
> **Spec Document:** `02-spec/21-app/204-ssh-password-interception-and-rsa-credential-vault/02-component-spec.md`  
> **Target Files:**  
> - `cli/cmdssh/ssh_login_cmd.go`  
> - `cli/cmdssh/ssh_login_pass_prompt.go`  
> - `cli/cmdssh/ssh_login_cmd_test.go`  

---

## 1. Technical Context & Scope

When a user executes `gitmap ssh login <target>` (or `gitmap ssh <target>`) without pre-existing SSH keys or saved credentials:
1. Historically, GitMap directly spawned the OpenSSH client.
2. OpenSSH took over the terminal to prompt for the password.
3. GitMap was completely blind to the password entry, preventing credential reuse or secure vault retention.
4. Users had to repeatedly enter passwords on every session, without being asked for consent or informed that credentials can be safeguarded with local RSA encryption.

### Core Objectives for Worker 02:
1. **Pre-Flight Authentication Pipeline**: In `cli/cmdssh/ssh_login_cmd.go`, evaluate public key auth (`tryConnectDefaultKey`) and stored vault credentials (`lookupVaultPassword`) before launching OpenSSH.
2. **Interactive Password Wrapping**: In new file `cli/cmdssh/ssh_login_pass_prompt.go`, wrap password entry in a masked terminal prompt (`term.ReadPassword`), ensuring raw characters are never echoed.
3. **Remote Dial Verification**: Actively verify the entered password against the remote sshd via `dialNodeWithPassword`. Bounded retry loop (up to 3 attempts) rejects invalid credentials before asking for retention.
4. **Explicit RSA Consent Prompt**: When credentials are confirmed valid, prompt the user:
   `Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: `
5. **Conditional Persistence**:
   - If confirmed (`y` or `yes`), encrypt via RSA-OAEP (`EncryptSSHPassword`) and save in SQLite `ssh_hosts`.
   - If rejected, use the password solely for the active session without touching the SQLite database.
6. **Seamless OpenSSH Dispatch**: Supply the verified password to `spawnSSHFn` with `attachAskPass`, allowing OpenSSH to connect without asking the operator for their password a second time.
7. **Comprehensive Unit Testing**: Author isolated unit tests in `cli/cmdssh/ssh_login_cmd_test.go` testing affirmative parsing, consent grant/decline, key bypass, and non-interactive handling.

---

## 2. Technical Specification & Implementation Plan

### 2.1 File 1: `cli/cmdssh/ssh_login_pass_prompt.go` (NEW FILE)

This file encapsulates all interactive terminal prompting, verification dial probes, and consent handling.

#### Pluggable Test Hooks & Data Types
```go
package cmdssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	readPasswordHook   = term.ReadPassword
	readConsentHook    = bufio.NewReader(os.Stdin).ReadString
	dialTargetPassHook = dialNodeWithPassword
	tryConnectKeyHook  = tryConnectDefaultKey
	isTerminalHook     = isInteractiveTerminal
)
```

#### Function Decomposition & Logic
1. `isConsentAffirmative(input string) bool`:
   - Trims whitespace from `input`.
   - Returns true if `strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes")`.
2. `formatPasswordPrompt(user, host string) string`:
   - Returns formatted prompt: `fmt.Sprintf("Enter password for %s@%s: ", user, host)`.
3. `formatConsentPrompt() string`:
   - Returns: `"Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: "`.
4. `readMaskedPassword(prompt string, fd int) (string, error)`:
   - Prints prompt to stdout.
   - Reads masked password via `readPasswordHook(fd)`.
   - Prints newline and returns trimmed string.
5. `readInteractiveConsent(r *bufio.Reader) (bool, error)`:
   - Reads line from `r` using `'\n'`.
   - Evaluates with `isConsentAffirmative`.
6. `verifyTargetPassword(target *SSHTarget, pass string) bool`:
   - Calls `dialTargetPassHook(target, pass)`.
   - If error is nil, closes client immediately and returns true.
   - If error is present, returns false.
7. `promptPasswordAttempt(ctx context.Context, target *SSHTarget, attempt int) (string, bool)`:
   - Formats prompt based on `attempt` count (showing `(attempt X of 3)` on retries).
   - Reads password via `readMaskedPassword`.
   - Verifies against remote host via `verifyTargetPassword`.
   - Returns captured password and boolean validity flag.
8. `captureVerifiedPassword(ctx context.Context, target *SSHTarget) (string, error)`:
   - Runs a bounded loop (up to 3 attempts).
   - If verified, prints green success checkmark and returns password.
   - If exhausted, returns `apperror.NewValidationError("authentication failed: invalid password")`.
9. `promptAndSaveConsent(ctx context.Context, alias string, target *SSHTarget, pass string) error`:
   - Displays consent prompt using `formatConsentPrompt()`.
   - Reads user input and checks affirmative consent.
   - If affirmative, calls `saveExplicitPassword(ctx, alias, target, pass)` and prints `✓ Password saved in local RSA vault for future logins.`.
   - If not affirmative, prints `ℹ Password will not be saved. Using for this session only.`.
10. `interceptAndResolvePassword(ctx context.Context, alias string, target *SSHTarget) (string, error)`:
    - Checks `isTerminalHook()`. If false, returns `apperror.NewValidationError("password required for SSH target but terminal is non-interactive")`.
    - Calls `captureVerifiedPassword(ctx, target)`.
    - If valid, calls `promptAndSaveConsent(ctx, alias, target, pass)`.
    - Returns verified password.
11. `tryKeyOrVaultProbe(ctx context.Context, alias string, target *SSHTarget) (string, bool)`:
    - Checks `tryConnectKeyHook(target)`. If client != nil, client.Close(), return `("", true)` (key auth succeeds).
    - Checks `resolveTargetPassword(target)`. If non-empty:
      - Validates with `verifyTargetPassword(target, storedPass)`.
      - If valid, return `(storedPass, true)`.
    - Returns `("", false)` (password prompt required).

---

### 2.2 File 2: `cli/cmdssh/ssh_login_cmd.go` (UPDATED)

Update `executeSSHLoginWithPassword` to integrate the interception pipeline:

```go
func executeSSHLoginWithPassword(ctx context.Context, target string, explicitPass string, force bool) error {
	sshTarget, err := ParseSSHTarget(target, "root", 22)
	if err != nil {
		return err
	}

	checkAndResolveIP(ctx, target, sshTarget)
	if err := checkAndResolveAlias(ctx, target, sshTarget); err != nil {
		return err
	}

	autoTrustTargetHostFn(ctx, sshTarget)

	password, err := resolveOrInterceptLoginPassword(ctx, target, sshTarget, explicitPass)
	if err != nil {
		return err
	}

	probeAndEnsureNodeProfile(ctx, target, sshTarget, password)
	return spawnSSHFn(ctx, *sshTarget, nil, password)
}

func resolveOrInterceptLoginPassword(ctx context.Context, target string, sshTarget *SSHTarget, explicitPass string) (string, error) {
	hasExplicitPass := explicitPass != ""
	if hasExplicitPass {
		saveExplicitPassword(ctx, target, sshTarget, explicitPass)
		return explicitPass, nil
	}

	knownPass, hasKnownAuth := tryKeyOrVaultProbe(ctx, target, sshTarget)
	if hasKnownAuth {
		return knownPass, nil
	}

	return interceptAndResolvePassword(ctx, target, sshTarget)
}
```

---

### 2.3 File 3: `cli/cmdssh/ssh_login_cmd_test.go` (NEW UNIT TESTS)

Add the following targeted tests to `cli/cmdssh/ssh_login_cmd_test.go`:

1. `TestIsConsentAffirmative_Table(t *testing.T)`:
   - Tests `"y"`, `"Y"`, `"yes"`, `"Yes"`, `"YES"`, `"  y  "` -> all return `true`.
   - Tests `"n"`, `"N"`, `"no"`, `""`, `"cancel"`, `"maybe"`, `"1"` -> all return `false`.
2. `TestFormatPrompts(t *testing.T)`:
   - Verifies `formatPasswordPrompt("admin", "192.168.1.10")` produces `"Enter password for admin@192.168.1.10: "`.
   - Verifies `formatConsentPrompt()` contains `[Encrypted locally with RSA algorithm]`.
3. `TestInterceptPassword_KeyAuthBypassesPrompt(t *testing.T)`:
   - Hooks `tryConnectKeyHook` to return a dummy `*ssh.Client`.
   - Tracks if `readPasswordHook` is invoked.
   - Verifies `resolveOrInterceptLoginPassword` returns empty password with zero prompts.
4. `TestInterceptPassword_ConsentGranted_SavesHost(t *testing.T)`:
   - Hooks `isTerminalHook` to return `true`.
   - Hooks `tryConnectKeyHook` to return `nil`.
   - Hooks `readPasswordHook` to return `"mySecretPass"`.
   - Hooks `dialTargetPassHook` to return `nil` (successful dial).
   - Hooks `readConsentHook` to return `"y\n"`.
   - Verifies password is returned and encrypted record is saved into SQLite `ssh_hosts`.
5. `TestInterceptPassword_ConsentDeclined_DoesNotSaveHost(t *testing.T)`:
   - Hooks `readConsentHook` to return `"n\n"`.
   - Verifies password is returned for session, but `ssh_hosts` table does NOT contain the host password.
6. `TestInterceptPassword_NonInteractiveAborts(t *testing.T)`:
   - Hooks `isTerminalHook` to return `false`.
   - Hooks `tryConnectKeyHook` to return `nil`.
   - Verifies structured `ValidationError` is returned without blocking.

---

## 3. Function Decomposition Plan (<= 15 Lines Per Function)

Every function implemented or refactored for this subtask must comply with the <= 15 line rule:

| Function Name | File | Target Lines | Purpose |
| :--- | :--- | :--- | :--- |
| `isConsentAffirmative` | `ssh_login_pass_prompt.go` | 6 lines | Case-insensitive affirmative check using `strings.EqualFold`. |
| `formatPasswordPrompt` | `ssh_login_pass_prompt.go` | 5 lines | Builds user@host password prompt string. |
| `formatConsentPrompt` | `ssh_login_pass_prompt.go` | 5 lines | Builds the RSA consent prompt string. |
| `readMaskedPassword` | `ssh_login_pass_prompt.go` | 10 lines | Terminal password reading with masked echo. |
| `readInteractiveConsent` | `ssh_login_pass_prompt.go` | 8 lines | Buffered stdin reader for consent line. |
| `verifyTargetPassword` | `ssh_login_pass_prompt.go` | 10 lines | Probes remote sshd via lightweight dialer. |
| `promptPasswordAttempt` | `ssh_login_pass_prompt.go` | 12 lines | Single prompt and verification attempt. |
| `captureVerifiedPassword` | `ssh_login_pass_prompt.go` | 14 lines | Retry loop (up to 3 attempts) for password entry. |
| `promptAndSaveConsent` | `ssh_login_pass_prompt.go` | 14 lines | Prompts for consent and saves on affirmative response. |
| `interceptAndResolvePassword` | `ssh_login_pass_prompt.go` | 12 lines | Orchestrates terminal check, capture, and consent. |
| `tryKeyOrVaultProbe` | `ssh_login_pass_prompt.go` | 14 lines | Evaluates key and vault auth pre-probes. |
| `resolveOrInterceptLoginPassword` | `ssh_login_cmd.go` | 12 lines | Decides whether explicit, stored, or intercepted credentials apply. |
| `executeSSHLoginWithPassword` | `ssh_login_cmd.go` | 14 lines | Refactored main login execution pipeline. |

---

## 4. Coding Guidelines & Invariant Rules

- **Positive Booleans Only**: Use `isConfirmed`, `hasKeyAuth`, `isValidPassword`, `isInteractive`, `hasExplicitPass`. Never use double-negative flags (`isNotSaved`, `hasNoKey`).
- **Zero-Allocation String Equality**: Always use `strings.EqualFold(a, b)` for case-insensitive checks. Never convert strings via `strings.ToLower()` for comparisons.
- **Function Sizing**: Strict <= 15 lines limit per function.
- **Unix LF Line Endings**: All newly created files and modifications must use Unix LF (`\n`).
- **Relative Path Hygiene**: All package references and file operations must use relative workspace paths.
- **AskPass Ephemerality**: Temporary scripts created by `attachAskPass` must always be cleaned up in a `defer` block.
- **Zero Test/Build Executions**: The worker must never run `go test` or `go build` directly. Verification is delegated to the CI/CD pipeline and designated QA workers.

---

## 5. Subtask Execution & Verification Checklist

- [ ] `cli/cmdssh/ssh_login_pass_prompt.go` created with all prompt, dial verification, and consent helpers.
- [ ] `cli/cmdssh/ssh_login_cmd.go` updated to integrate `resolveOrInterceptLoginPassword`.
- [ ] `cli/cmdssh/ssh_login_cmd_test.go` extended with comprehensive unit tests for consent, bypass, and persistence.
- [ ] Line count for all modified and created functions verified <= 15 lines.
- [ ] Affirmative consent strings verified using `strings.EqualFold`.
- [ ] RSA algorithm explicitly mentioned in consent prompt string.
- [ ] Ephemeral askpass lifecycle properly cleaned up.
