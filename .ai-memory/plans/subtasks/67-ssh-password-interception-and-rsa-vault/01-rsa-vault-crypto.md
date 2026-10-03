# Subtask 67-01: RSA Credential Vault Keypair & RSA-OAEP Cryptographic Engine

**Subtask ID:** 67-01  
**Status:** Pending  
**Spec Reference:** [02-spec/21-app/204-ssh-password-interception-and-rsa-credential-vault/01-architecture-spec.md](../../../../02-spec/21-app/204-ssh-password-interception-and-rsa-credential-vault/01-architecture-spec.md)  
**Assigned Worker:** Worker 01  
**Target Files:**  
- `cli/cmdssh/ssh_vault_rsa.go` [NEW]  
- `cli/cmdssh/ssh_crypto_pass.go` [MODIFY / ENHANCE]  
- `cli/cmdssh/ssh_crypto_pass_test.go` [MODIFY / ENHANCE]  

---

## 1. Objective & Scope

Implement the dedicated GitMap RSA Credential Vault and upgrade SSH password encryption to use RSA-OAEP with SHA-256 and intrinsic non-deterministic salt/padding. This ensures GitMap can securely encrypt SSH passwords at rest in SQLite `ssh_hosts` and reversibly decrypt them into memory for automated headless execution, while preventing rainbow table and frequency analysis attacks.

---

## 2. File Implementation Specifications

### 2.1 Component 1: `cli/cmdssh/ssh_vault_rsa.go` [NEW FILE]
Create a modular key management file responsible for dedicated vault keypair provisioning, filesystem permissions, and key serialization.

#### Functions to Implement (strictly <= 15 lines each):
1. **`resolveVaultKeysDir() (string, error)`**
   - Resolves user home directory via `os.UserHomeDir()`.
   - Returns absolute path to `~/.gitmap/keys`.
2. **`resolveVaultKeyPaths() (string, string, error)`**
   - Returns `(privKeyPath, pubKeyPath, error)` for `vault_rsa` and `vault_rsa.pub`.
3. **`hasVaultRSAKey() bool`**
   - Checks if `vault_rsa` exists on disk.
4. **`EnsureVaultRSAKeyPair() (*rsa.PrivateKey, error)`**
   - Checks `hasVaultRSAKey()`. If present, loads and returns key via `LoadVaultRSAPrivateKey()`.
   - If absent, provisions directory `~/.gitmap/keys` with `0700` permissions and calls `generateAndSaveVaultRSAKeyPair()`.
5. **`LoadVaultRSAPrivateKey() (*rsa.PrivateKey, error)`**
   - Reads `vault_rsa` from disk and parses PKCS#1 or PKCS#8 PEM blocks.
6. **`generateAndSaveVaultRSAKeyPair(privPath, pubPath string) (*rsa.PrivateKey, error)`**
   - Generates 2048-bit key via `rsa.GenerateKey(rand.Reader, 2048)`.
   - Saves private key PEM with `0600` permissions.
   - Saves public key PEM with `0644` permissions.
7. **`encodeRSAPrivateKeyPEM(priv *rsa.PrivateKey) []byte`**
   - Encodes private key using `x509.MarshalPKCS1PrivateKey` into a `pem.Block{Type: "RSA PRIVATE KEY"}`.
8. **`encodeRSAPublicKeyPEM(pub *rsa.PublicKey) ([]byte, error)`**
   - Encodes public key using `x509.MarshalPKIXPublicKey` into a `pem.Block{Type: "PUBLIC KEY"}`.
9. **`saveVaultFile(path string, data []byte, perm os.FileMode) error`**
   - Writes file atomically and wraps any filesystem error using `apperror.WrapSimple`.
10. **`parseRSAPrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error)`**
    - Decodes PEM block and parses using `x509.ParsePKCS1PrivateKey` (or fallback `x509.ParsePKCS8PrivateKey`).

---

### 2.2 Component 2: `cli/cmdssh/ssh_crypto_pass.go` [MODIFY / ENHANCE]
Upgrade password encryption and decryption to integrate with the dedicated vault keypair while preserving backward compatibility.

#### Key Enhancements:
1. **Vault-First Key Loading:**
   - Update `loadOrEnsureRSAPrivateKey()` to call `EnsureVaultRSAKeyPair()`.
   - Fall back to `loadLocalSSHRSAKey()` (`~/.ssh/id_rsa`) if vault key generation fails.
   - Fall back to `encryptWithFallbackAES(plain)` only if all RSA mechanisms are unavailable.
2. **Vault-First Key Decryption:**
   - In `decryptRSAPassword(body string)`:
     - Attempt decryption with `LoadVaultRSAPrivateKey()`.
     - If key not found or decryption fails, attempt `loadLocalSSHRSAKey()`.
     - Ensures smooth zero-downtime migration for existing credentials.
3. **Coding Guideline & Sizing Refactor:**
   - Audit all functions in `ssh_crypto_pass.go`.
   - Break down any function exceeding 15 lines (e.g. `DecryptSSHPassword` or `parseRawRSAPrivateKey`).
   - Use positive boolean variables (`hasKey`, `isRSA`, `hasVaultKey`, `isAES`).

---

### 2.3 Component 3: `cli/cmdssh/ssh_crypto_pass_test.go` [MODIFY / ENHANCE]
Expand unit test coverage with dedicated test cases validating RSA-OAEP properties and vault key management.

#### Test Cases to Add:
1. **`TestEnsureVaultRSAKeyPair_GenerationAndPermissions`**
   - Tests directory creation (`0700`) and private key file creation (`0600`).
   - Validates that repeated calls return the same key without overwriting.
2. **`TestRSAOAEP_NonDeterministicSaltCiphertext`**
   - Encrypts the identical password `MyP@ssw0rd!2026` twice.
   - Asserts `cipher1 != cipher2` (proving random OAEP seed/salt uniqueness).
   - Decrypts both ciphertexts and asserts both yield `MyP@ssw0rd!2026`.
3. **`TestVaultRSA_EncryptDecryptRoundtrip`**
   - Tests roundtrip encryption and decryption with special characters, symbols, and long passphrases.
4. **`TestVaultRSA_FallbackToSSHRSA`**
   - Simulates absent vault key and validates successful fallback to `~/.ssh/id_rsa`.
5. **`TestVaultRSA_FallbackToAES`**
   - Simulates absent RSA keys and validates fallback to AES encryption.
6. **`TestDecryptSSHPassword_LegacyCompatibility`**
   - Validates decoding of `rsa:`, `aes:`, `salt:`, and `caesar:` prefixed strings.

---

## 3. Strict Coding Guidelines & Constraints

- **Function Sizing:** No function may exceed 15 lines of code. Split complex logic into focused single-responsibility helpers.
- **Boolean Naming:** Use positive prefixes only (`hasKey`, `isRSA`, `isValid`, `isGenerated`). Never use inverted/negative booleans (`isNotValid`, `noKey`).
- **Error Wrapping:** Always use `apperror.WrapSimple(err, op)` or `apperror.NewValidationError(...)`.
- **Paths:** Strict relative paths from repository root or standard user directory utilities (`os.UserHomeDir`).
- **Line Endings:** Unix LF line endings across all files.
- **ZERO GIT COMMANDS:** Never execute any git commands (`git status`, `git add`, `git commit`, `git push`, `git checkout`).
- **ZERO BUILDS OR TESTS:** Never run `go build` or `go test`. Worker focuses purely on surgical code and test authoring.

---

## 4. Acceptance Criteria Checklist

- [ ] `cli/cmdssh/ssh_vault_rsa.go` is created with pure Go RSA generation (`rsa.GenerateKey(rand.Reader, 2048)`).
- [ ] Directory `~/.gitmap/keys` is created with `0700` permissions.
- [ ] Private key `vault_rsa` is saved with `0600` permissions.
- [ ] Public key `vault_rsa.pub` is saved with `0644` permissions.
- [ ] `EncryptSSHPassword` uses `EnsureVaultRSAKeyPair()` with RSA-OAEP SHA-256 by default.
- [ ] Encrypting the same password multiple times produces distinct non-deterministic ciphertexts.
- [ ] `DecryptSSHPassword` handles `rsa:`, `aes:`, `salt:`, and `caesar:` prefixes seamlessly.
- [ ] All functions in `ssh_vault_rsa.go`, `ssh_crypto_pass.go`, and `ssh_crypto_pass_test.go` are <= 15 lines.
- [ ] All boolean variables follow positive naming conventions (`has...`, `is...`).
