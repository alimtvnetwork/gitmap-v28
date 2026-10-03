# Execution Plan: 67-ssh-password-interception-and-rsa-vault

**Status:** `Completed`  
**Completion Date:** 2026-10-03  
**Verified In:** `v6.471.0`  

## User Request (Verbatim)
```text
For the first time running, if the `Gitmap` knows that this login requires a password and it does not have a password, then it would be wrapped by `Gitmap` as a password, and when user types the password, it would be saved for the future use. You did not follow through this. Okay? Make sure you do that. And also, it will be prompted to the user that we can keep the password for the future use. Do you want to save it or not? Okay, yes or no? Answer. And also mention that we save the password with the RSA algorithm. So make sure you do that. Okay? When you save the password, save the password with the RSA algorithm with salt so that it cannot be retrieved directly. Or we need to retrieve it, so we cannot add the salt, I believe. Yeah, we cannot probably add the salt. Okay. Try this out, and let me know that you understood. So make a release after fixing this, and do a minor bump. Is it clear?
```

## Actionable Items Must Follow Non-Negotiable
1. Ensure `Gitmap` wraps the password when required and prompts the user to save it for future use. -> [COMPLETED]
2. Implement password saving using the RSA algorithm, considering the use of salt. -> [COMPLETED]
3. Make a release after implementing the changes and perform a minor version bump. -> [COMPLETED]

---

## 1. Architecture & Design

### 1.1 Password Interception Flow
1. Target Host Resolution: In `executeSSHLoginWithPassword` (`cli/cmdssh/ssh_login_cmd.go`), resolve target host alias (e.g. `t1`) and credentials.
2. Pre-Flight Authentication Probe:
   - Check if public key authentication succeeds with remote host (`tryConnectDefaultKey(sshTarget)`).
   - If public key succeeds: spawn OpenSSH directly without prompt.
   - If public key fails: check local vault for stored password (`lookupVaultPassword(alias)`).
   - If stored password found: verify and spawn OpenSSH with `SSH_ASKPASS` (`attachAskPass`).
   - If NO stored password in vault:
     - Check if interactive terminal (`isInteractiveTerminal()`).
     - Prompt user via GitMap with masked password input: `Enter password for <user>@<host>: `.
     - Validate password with remote sshd via lightweight dialer probe (`dialNodeWithPassword`).
     - If invalid: display error and reprompt / abort cleanly.
     - If valid: prompt user with explicit RSA mention:
       `Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: `
     - If user answers `y` / `yes`:
       - Encrypt password using RSA-OAEP with SHA-256 and store in vault / `ssh_hosts`.
     - Pass verified password to `attachAskPass` and spawn OpenSSH (`SpawnSSHWithPassword`).
     - Interactive session proceeds seamlessly with zero further prompts.

### 1.2 RSA Cryptographic Vault Design & Salt Considerations
1. Asymmetric Keypair Management:
   - Store or generate dedicated GitMap RSA keypair at `~/.gitmap/keys/vault_rsa` and `~/.gitmap/keys/vault_rsa.pub` (2048-bit RSA, permissions 0600).
   - If `~/.ssh/id_rsa` already exists, allow it as fallback keypair.
2. Encryption with Reversibility & Salt:
   - To retrieve and supply passwords to OpenSSH / automated SSH sessions, the ciphertext must be reversibly decryptable using the local private key.
   - Standard password hashing (like bcrypt / argon2) cannot be decrypted and therefore cannot be used for SSH authentication.
   - RSA-OAEP (Optimal Asymmetric Encryption Padding) with SHA-256 natively introduces a randomized salt/seed (`rand.Reader`) during OAEP encoding.
   - This ensures identical passwords encrypt to completely distinct, non-deterministic ciphertexts every time, preventing rainbow table attacks and plaintext inspection while remaining reversibly decryptable by GitMap.

---

## 2. Work Breakdown & Subtask Ownership

| Task ID | Subtask Name | Owner | Target Files | Status |
| :--- | :--- | :--- | :--- | :--- |
| `Task-01` | RSA Vault Cryptography & Key Management | Worker 01 | `cli/cmdssh/ssh_crypto_pass.go`, `cli/cmdssh/ssh_vault_rsa.go`, `cli/cmdssh/ssh_crypto_pass_test.go` | PASS |
| `Task-02` | SSH Password Interception, Consent & Flow | Worker 02 | `cli/cmdssh/ssh_login_cmd.go`, `cli/cmdssh/ssh_login_pass_prompt.go`, `cli/cmdssh/ssh_login_cmd_test.go` | PASS |

---

## 3. Targeted Quality Checks & Gates
1. `python .github/scripts/check-legacy-refs.py .` -> Exit 0.
2. `python .github/scripts/go-format-check.py --check-only` -> Exit 0.
3. `python .github/scripts/misspell-changed.py` -> Exit 0.
4. `python -m pytest .github/scripts/tests/test_ci_scripts.py` -> Exit 0 (18/18 passed).
5. Targeted Go unit tests (`cmdssh`) -> Exit 0 (11 passed).
6. Secrets Gate -> Exit 0 (0 hits).
