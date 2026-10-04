# Component Specification: SSH Password Interception, Interactive RSA Consent Prompt & AskPass Session Dispatch

- **Spec Document:** `02-spec/21-app/204-ssh-password-interception-and-rsa-credential-vault/02-component-spec.md`
- **Spec ID:** 204 Component Spec
- **Status:** Approved
- **Subsystem:** `cli/cmdssh`
- **Author:** Spec Author 02 (Component Specification & Subtask 02 Planning)
- **Parent Plan:** [.ai-memory/plans/pending/67-ssh-password-interception-and-rsa-vault.md](../../../.ai-memory/plans/pending/67-ssh-password-interception-and-rsa-vault.md)
- **Target Subtask:** [.ai-memory/plans/subtasks/67-ssh-password-interception-and-rsa-vault/02-ssh-login-interception.md](../../../.ai-memory/plans/subtasks/67-ssh-password-interception-and-rsa-vault/02-ssh-login-interception.md)

---

## 1. Executive Summary & Problem Analysis

### 1.1 Context & Problem Statement
GitMap provides high-speed SSH terminal access and fleet orchestration via commands such as `gitmap ssh login [target]`, `gitmap ssh [target]`, and cluster execution utilities. Historically, when connecting to a target machine for the first time without pre-configured SSH keys or a stored password in the SQLite credential registry (`ssh_hosts`), GitMap bypassed pre-flight interception and directly executed the OpenSSH binary:

```
gitmap ssh login target-node
  -> spawnSSHFn (ssh target-node)
  -> OpenSSH standard TTY prompt: "user@target-node's password: "
```

This legacy execution model suffered from four critical operational flaws:
1. **Zero Password Interception & Retention**: Because OpenSSH prompted the user directly in raw TTY mode, GitMap had no hook to capture or remember credentials. On every subsequent invocation, the engineer was forced to re-type their password.
2. **Missing User Agency & Consent**: Operating systems and secure CLI tools must never persist sensitive credentials without explicit, informed user consent. GitMap lacked any interactive prompt to ask users if they wanted their password retained.
3. **Absence of Cryptographic Disclosure**: Users were not informed about how credentials would be safeguarded at rest. Transparent security requires explicitly stating that passwords are encrypted locally using the RSA algorithm.
4. **Interactive AskPass Disconnect**: Even when credentials were known, spawning OpenSSH without proper `SSH_ASKPASS` orchestration resulted in redundant prompts or broken handshakes across heterogeneous operating systems (Windows, Linux, macOS).

### 1.2 Architectural Objective
This component specification establishes an automated **SSH Password Interception & Pre-Flight Probe Engine** inside `cli/cmdssh`. When an SSH login is executed:
- GitMap autonomously evaluates authentication capabilities before launching OpenSSH.
- If public key authentication succeeds, connection proceeds immediately with zero user friction.
- If stored vault credentials exist, GitMap verifies them against the target and injects them seamlessly via `SSH_ASKPASS`.
- If no credentials exist, GitMap wraps password entry in a masked terminal prompt (`term.ReadPassword`), validates the password against the target sshd in real time, requests user consent with explicit mention of RSA encryption, persists confirmed credentials in SQLite `ssh_hosts`, and launches OpenSSH with `SSH_ASKPASS` for a zero-reprompt interactive session.

---

## 2. Sequence Flow of SSH Connection Interception

The complete end-to-end interception and authentication lifecycle is structured into four deterministic phases:

```
[User: gitmap ssh login <target>]
               │
               ▼
   [Phase 1: Target Resolution]
   Resolve Alias / IP / User / Port
               │
               ▼
   [Phase 2: Pre-Probe 1 - Public Key Auth]
   tryConnectDefaultKey(sshTarget)
         │                       │
      (Success)               (Failed)
         │                       │
         ▼                       ▼
 [Spawn OpenSSH]    [Phase 3: Pre-Probe 2 - Vault Credential Check]
  (Zero prompts)     lookupVaultPassword(alias, ip)
                           │                        │
                        (Found)                 (Not Found)
                           │                        │
                           ▼                        ▼
                [dialNodeWithPassword]      [Phase 4: Pre-Probe 3 - Interactive Interception]
                    │            │          isInteractiveTerminal()
                 (Valid)      (Invalid)          │                   │
                    │            │            (True)              (False)
                    │            ▼               │                   │
                    │     [Purge Stale]          ▼                   ▼
                    │            │      [Masked Prompt:      [Fail Fast / Pass-through]
                    │            └────►  Enter password]
                    ▼                            │
             [Spawn OpenSSH]                     ▼
              with AskPass              [dialNodeWithPassword]
                                         │                   │
                                      (Valid)             (Invalid)
                                         │                   │
                                         ▼                   ▼
                                [Consent Prompt:       [Display Error &
                                 Save via RSA? (y/n)]   Reprompt / Exit]
                                  │                │
                               (Yes)              (No)
                                  │                │
                                  ▼                ▼
                           [Encrypt RSA-OAEP]  [Session Only]
                           [Persist in DB]         │
                                  │                │
                                  └───────┬────────┘
                                          │
                                          ▼
                                   [Spawn OpenSSH]
                                    with AskPass
```

### 2.1 Phase 1: Target Host Resolution
1. Target input argument (e.g. `t1`, `node-alias`, `admin@10.0.0.12:2222`) is parsed via `ParseSSHTarget`.
2. If target is an alias, `lookupSSHHostOrReport` queries SQLite `ssh_hosts` (and `ssh_connections` fallback) to resolve IP address, port, username, and any pre-existing `EncryptedPassword`.
3. If target is an IP address, `checkAndResolveIP` maps IP to known host records.
4. Auto-trust verification ensures known_hosts records are checked or initialized (`autoTrustTargetHostFn`).

### 2.2 Phase 2: Pre-Probe 1 - Public Key Authentication Test
1. Before attempting password operations, GitMap invokes `tryConnectDefaultKey(sshTarget)`.
2. `tryConnectDefaultKey` locates local user keys (`~/.ssh/id_rsa`, `~/.ssh/id_ed25519`, `~/.ssh/id_ecdsa`) via `findAllUserSSHKeys()`.
3. It performs a lightweight non-interactive SSH handshake with remote sshd using `newAutoAcceptHostKeyConfig`.
4. **Decision Branch:**
   - **Key Accepted**: Close probe client immediately. Invoke `spawnSSHFn(ctx, *sshTarget, nil, "")`. OpenSSH connects using public key authentication with zero prompts. Pre-flight complete.
   - **Key Rejected / No Keys**: Proceed to Phase 3.

### 2.3 Phase 3: Pre-Probe 2 - Vault Password Check & Verification
1. Inspect `sshTarget.EncryptedPassword` or invoke `lookupVaultPassword(alias, ip)`.
2. If an encrypted password is discovered:
   - Decrypt the secret in memory via `DecryptSSHPassword(encPass)`. Supports `rsa:` (RSA-OAEP with SHA-256), `aes:` (AES-GCM fallback), `salt:`, and legacy hashes.
   - Validate decrypted password against remote sshd via `dialNodeWithPassword(sshTarget, plainPass)`.
   - `dialNodeWithPassword` executes RFC 4252 password auth followed by PAM keyboard-interactive fallback.
3. **Decision Branch:**
   - **Vault Password Valid**: Close probe client. Spawn OpenSSH via `spawnSSHFn(ctx, *sshTarget, nil, plainPass)`. The password is supplied via `attachAskPass` using a temporary secure script. Terminal session opens seamlessly.
   - **Vault Password Invalid (Stale/Rotated)**: Log warning in terminal (`"Stored vault password rejected by remote host. Re-authenticating..."`). Invalidate stale in-memory password and fall through to Phase 4.

### 2.4 Phase 4: Pre-Probe 3 - Interactive Password Prompt, Dial Verification & RSA Consent
When neither public keys nor valid vault passwords exist:

#### Step 4.1: Terminal Environment Detection
- Evaluate `isInteractiveTerminal()` via `term.IsTerminal(int(os.Stdin.Fd()))`.
- If `!isInteractiveTerminal()`:
  - In automated CI/CD pipelines or headless scripts, non-interactive prompting is impossible.
  - Emit structured `apperror.NewValidationError("password required for SSH target '%s' but terminal is non-interactive", target)` and exit without hanging.

#### Step 4.2: Masked Password Entry
- Display clean prompt to standard output:
  ```
  Enter password for <user>@<host>: 
  ```
- Capture password using `term.ReadPassword(int(os.Stdin.Fd()))`.
- Input is masked: keystrokes do not echo to terminal display.
- Trim trailing carriage returns and newlines (`\r`, `\n`).
- If an empty password is submitted, report validation error and reprompt (bounded up to 3 attempts).

#### Step 4.3: Real-Time Remote Verification Probe
- Before asking the user about credential retention, GitMap verifies that the entered password is actually correct by dialing the remote host:
  ```go
  probeClient, dialErr := dialNodeWithPassword(sshTarget, enteredPassword)
  ```
- **Probe Failure**:
  - If remote sshd rejects the password, display clean error message:
    ```
    ✗ Authentication failed: Invalid password for admin@node-alias. Please try again.
    ```
  - Reprompt user (up to 3 total attempts). If all attempts fail, abort cleanly with exit code 1.
  - Crucially: **Invalid passwords are NEVER saved to the database.**
- **Probe Success**:
  - Close the probe client connection immediately (`_ = probeClient.Close()`).
  - Proceed to Step 4.4 with 100% verified credentials.

#### Step 4.4: User Consent Prompt with Explicit RSA Disclosure
- Prompt the operator for persistence consent:
  ```
  Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: 
  ```
- Read response from standard input using buffered line reader.
- Parse confirmation using zero-allocation case-insensitive evaluation (`strings.EqualFold`):
  ```go
  func isConsentAffirmative(input string) bool {
      trimmed := strings.TrimSpace(input)
      return strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes")
  }
  ```
- **User Confirms (`y` or `yes`)**:
  - Encrypt password using RSA-OAEP with SHA-256 via `EncryptSSHPassword(enteredPassword)`.
  - Store encrypted secret in SQLite `ssh_hosts` table (and sync with `ssh_connections`).
  - Emit positive confirmation:
    ```
      ✓ Password saved in local RSA vault for future logins.
    ```
- **User Declines (`n`, `no`, or any other input)**:
  - Do NOT modify SQLite database.
  - Emit informational notice:
    ```
      ℹ Password will not be saved. Using for this session only.
    ```

#### Step 4.5: OpenSSH Session Dispatch via AskPass
- Dispatch OpenSSH session via `spawnSSHFn(ctx, *sshTarget, nil, enteredPassword)`.
- `attachAskPass` generates an ephemeral script (`.bat` on Windows, `.sh` on Linux/macOS) with `0700` permissions.
- Injects environment variables:
  - `SSH_ASKPASS=<path-to-script>`
  - `SSH_ASKPASS_REQUIRE=force`
  - `DISPLAY=1`
  - `GITMAP_SSH_PASS=<entered-password>`
- Spawns interactive OpenSSH client connected to terminal `os.Stdin`, `os.Stdout`, and `os.Stderr`.
- Registers cleanup closure removing the temporary askpass script immediately upon process termination.
- Result: The user experiences an instant, seamless interactive SSH shell without being prompted for the password a second time by OpenSSH.

---

## 3. UI/UX Specifications

### 3.1 Terminal Output Presentation
All terminal messages must adhere strictly to GitMap's terminal styling guidelines, using ANSI colors from `cli/constants` and affirmative phrasing:

```text
Enter password for admin@node-alias: **********
  • Verifying credentials with remote host...
  ✓ Authentication verified.
Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: y
  ✓ Password saved in local RSA vault for future logins.

[OpenSSH interactive session begins below]
admin@ubuntu:~$ 
```

### 3.2 Declined Consent Terminal Presentation
When the operator declines credential persistence:

```text
Enter password for root@10.0.0.5: **********
  • Verifying credentials with remote host...
  ✓ Authentication verified.
Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: n
  ℹ Password will not be saved. Using for this session only.

[OpenSSH interactive session begins below]
root@server:~$ 
```

### 3.3 Authentication Error Terminal Presentation
When incorrect credentials are supplied:

```text
Enter password for ubuntu@node-main0: **********
  • Verifying credentials with remote host...
  ✗ Authentication failed: invalid password or remote server rejected credentials.
Enter password for ubuntu@node-main0 (attempt 2 of 3): 
```

### 3.4 Key UX Requirements & Invariants
1. **Masked Password Input**: Must never echo raw characters to the terminal.
2. **Explicit RSA Mention**: The consent prompt must explicitly state `[Encrypted locally with RSA algorithm]` to give operators confidence in local cryptographic security.
3. **Zero Double-Prompting**: Once the password is verified by GitMap, `attachAskPass` must feed it directly to OpenSSH so the user is never asked for the password again during that connection.
4. **Graceful Cancellation**: Pressing `Ctrl+C` or sending SIGINT during password entry or consent prompt must cleanly restore terminal echo mode and exit without panics or raw stack traces.

---

## 4. Edge Cases & Resilience Engineering

| Edge Case | Detection Condition | Remediation & Handling |
| :--- | :--- | :--- |
| **Non-Interactive Environment** | `isInteractiveTerminal() == false` (e.g. CI/CD runners, background daemons, pipes). | Do not block on `term.ReadPassword`. Return structured `apperror.NewValidationError` explaining that an interactive terminal or explicit credential flag is required. |
| **Invalid Password Entry** | Remote sshd returns `SSH_MSG_USERAUTH_FAILURE` or handshake error in `dialNodeWithPassword`. | Discard password immediately. Do not prompt for consent. Display clear error and allow up to 3 attempts before exiting. |
| **Consent Rejection** | User types `n`, `no`, or presses Enter without `y`/`yes`. | Proceed with `spawnSSHFn` for the current session using `attachAskPass`. Do NOT write to `ssh_hosts` or modify SQLite records. |
| **Stale / Rotated Vault Password** | `ssh_hosts` contains an encrypted password, but `dialNodeWithPassword` fails with auth error. | Invalidate in-memory vault password. Warn user that stored password is no longer accepted. Fall through to interactive password prompt to capture and update the new password. |
| **Host Unreachable / Network Timeout** | `dialNodeWithPassword` fails due to network timeout or connection refused. | Differentiate network failure from authentication failure via `isNetworkDialFailure`. Print diagnostic footer (`gitmap ssh troubleshoot <target>`) and exit without prompting for consent. |
| **Explicit Password CLI Flag** | User provided explicit password argument (`gitmap ssh login target pass123`). | Explicit command-line arguments indicate direct operator intent. Existing behavior preserves automatic persistence without interactive consent prompt. |
| **Missing RSA Keypair** | Local `~/.ssh/id_rsa` or `~/.gitmap/keys/vault_rsa` does not exist. | Cryptographic module automatically falls back to authenticated AES-256-GCM or triggers local RSA keypair generation. |

---

## 5. Component Architecture & Code Contracts

### 5.1 Package & File Boundaries
```
cli/cmdssh/
├── ssh_login_cmd.go            # Command routing & updated executeSSHLoginWithPassword pipeline
├── ssh_login_pass_prompt.go    # NEW: Interactive prompt, remote dial verification & RSA consent
├── ssh_login_cmd_test.go       # Comprehensive unit tests for interception, consent & decision tree
├── ssh_crypto_pass.go          # RSA-OAEP encryption/decryption engine
├── ssh_askpass.go              # Ephemeral askpass script generator & environment injector
└── sshjoin_dial.go             # dialNodeWithPassword & tryConnectDefaultKey probe dialers
```

### 5.2 Core Types & Function Signatures (`cli/cmdssh/ssh_login_pass_prompt.go`)

```go
package cmdssh

import (
	"context"
	"io"
)

// PasswordConsentResult captures user consent and captured credentials.
type PasswordConsentResult struct {
	Password    string
	IsConsented bool
	IsValid     bool
}

// InteractivePromptHooks enables isolated unit testing of terminal I/O.
type InteractivePromptHooks struct {
	ReadPasswordFn   func(prompt string, fd int) (string, error)
	ReadConsentFn    func(r io.Reader) (bool, error)
	DialPasswordFn   func(target *SSHTarget, pass string) error
	TryConnectKeyFn  func(target *SSHTarget) bool
	IsTerminalFn     func() bool
}

// isConsentAffirmative tests case-insensitive confirmation.
func isConsentAffirmative(input string) bool

// formatPasswordPrompt generates target-specific prompt string.
func formatPasswordPrompt(user, host string) string

// formatConsentPrompt returns the standard RSA consent question.
func formatConsentPrompt() string

// interceptLoginPassword manages the interactive prompt, verification, and consent workflow.
func interceptLoginPassword(ctx context.Context, alias string, target *SSHTarget) (string, error)

// tryKeyOrVaultAuth executes Pre-Probe 1 and Pre-Probe 2.
func tryKeyOrVaultAuth(ctx context.Context, alias string, target *SSHTarget) (string, bool)
```

---

## 6. Verification & Quality Gates

1. **Coding Guideline Conformance**:
   - Every function strictly <= 15 lines.
   - 100% positive booleans (`is`/`has` prefixes: `isConfirmed`, `hasKeyAuth`, `isValidPassword`, `isInteractive`).
   - Zero `strings.ToLower` comparisons (all case-insensitive logic uses `strings.EqualFold`).
   - Unix LF (`\n`) line endings preserved.
2. **Automated Unit Testing**:
   - Affirmative consent parsing table test (`"y"`, `"yes"`, `"Y"`, `"YES"` -> true; `"n"`, `"no"`, `""` -> false).
   - Pre-flight bypass test: key auth success immediately bypasses password prompting.
   - Consent granted test: verifies password is encrypted and persisted in SQLite `ssh_hosts`.
   - Consent declined test: verifies password is dispatched to `spawnSSHFn` but NOT written to SQLite.
   - Non-interactive test: verifies error is returned when `isInteractiveTerminal() == false`.
3. **Security Gate**:
   - Zero secrets committed to git.
   - Ephemeral askpass scripts created with `0700` permissions and immediately deleted on process exit.
   - Passwords stored exclusively as RSA-OAEP / AES encrypted ciphertext.
