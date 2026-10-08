# Subtask 02: Multi-Supabase Database Vault & AES-GCM Encryption

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md`  
> **Target Subsystems:** `cli/cmdsupabase/`, `cli/cmd/rootdata.go`, `cli/cmd/roottooling.go`  
> **Owned Files:**  
> - `cli/cmdsupabase/supabase_types.go`  
> - `cli/cmdsupabase/supabase_vault.go`  
> - `cli/cmdsupabase/supabase_db.go`  
> - `cli/cmdsupabase/supabase_cmd.go`  
> - `cli/cmdsupabase/supabase_test.go`  
> - `cli/cmd/rootdata.go`  
> - `cli/cmd/roottooling.go`  

---

## 1. Concrete Objectives

1. **Multi-Supabase CLI Registration & Management:**
   - Implement single-line project registration:
     `gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]`
   - Parse optional flags: `--desc "<description>"`, `--project-ref "<ref>"`.
   - Validate URL format (`https://*.supabase.co` or self-hosted HTTP endpoints).
2. **Authenticated AES-GCM Vault Encryption (Zero Cleartext Invariant):**
   - Implement encryption engine in `cli/cmdsupabase/supabase_vault.go`.
   - Derive a 32-byte master encryption key combining the host machine identity (CPU/motherboard UUID) and an installation-scoped salt (`.gitmap/data/installation/vault.salt`).
   - Encrypt all sensitive credentials (`anon_key`, `service_key`, and `db_url`) using `cli/crypto.Encrypt(plaintext, key)` (AES-256-GCM with standard 12-byte nonce).
   - Under no circumstances may plaintext service keys, passwords, or tokens be stored in SQLite or written to logs.
3. **SQLite Split Database Persistence:**
   - In `cli/cmdsupabase/supabase_db.go`, initialize and manage the database at `.gitmap/data/installation/supabase/sql.db` (utilizing `store.ResolveInstallationDbPath("supabase/sql.db")`).
   - Create table `supabase_databases` with indexed columns.
   - Implement transactional CRUD operations: `InsertProject`, `GetProject`, `ListProjects`, `UpdateProjectStatus`, `DeleteProject`.
4. **Credential Masking & Inspection Operations:**
   - Implement `gitmap supabase list` (alias: `gitmap supabase ls`): Render ANSI table displaying `ALIAS`, `API URL`, `STATUS`, `ANON KEY (MASKED)`, `HAS DB URL`, `UPDATED AT`. Mask sensitive keys (e.g. `eyJhbGci...****`).
   - Implement `gitmap supabase get <alias> [--decrypt]`: Output project details. Only reveal decrypted secrets when the explicit `--decrypt` flag is passed.
   - Implement `gitmap supabase test <alias>` (alias: `gitmap supabase ping`): Test connection by sending an authenticated HTTP GET request to `<api_url>/rest/v1/` using the decrypted anon or service role key.
   - Implement `gitmap supabase remove <alias>` (alias: `gitmap supabase rm`): Safely remove project from vault.
   - Implement `gitmap supabase export-env <alias>`: Generate shell-compatible environment exports (`SUPABASE_URL`, `SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_ROLE_KEY`, `DATABASE_URL`).
5. **Command Dispatcher Routing:**
   - Register `supabase` and `sb` aliases in `cli/cmd/rootdata.go` within `dataDatabaseEntries()` or `dataExecutionEntries()`.
   - Add help menu entry to `cli/cmd/roottooling.go` ensuring discoverability in `gitmap --help`.

---

## 2. Core Domain Types & Structs

```go
package cmdsupabase

import (
	"time"
)

// SupabaseProjectRecord represents an encrypted project entry in the SQLite vault.
type SupabaseProjectRecord struct {
	Alias         string    `json:"alias"`
	ProjectRef    string    `json:"project_ref"`
	ApiUrl        string    `json:"api_url"`
	AnonKeyEnc    string    `json:"anon_key_enc"`
	ServiceKeyEnc string    `json:"service_key_enc"`
	DbUrlEnc      string    `json:"db_url_enc"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SupabaseProjectDecrypted represents in-memory credentials after authorized decryption.
type SupabaseProjectDecrypted struct {
	Alias       string `json:"alias"`
	ProjectRef  string `json:"project_ref"`
	ApiUrl      string `json:"api_url"`
	AnonKey     string `json:"anon_key"`
	ServiceKey  string `json:"service_key"`
	DbUrl       string `json:"db_url"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// SupabaseProjectView represents safe, masked metadata for tabular CLI display.
type SupabaseProjectView struct {
	Alias            string `json:"alias"`
	ProjectRef       string `json:"project_ref"`
	ApiUrl           string `json:"api_url"`
	AnonKeyMasked    string `json:"anon_key_masked"`
	ServiceKeyMasked string `json:"service_key_masked"`
	HasDbUrl         bool   `json:"has_db_url"`
	Status           string `json:"status"`
	UpdatedAt        string `json:"updated_at"`
}

// SupabasePingResult captures connectivity probe telemetry.
type SupabasePingResult struct {
	Alias          string        `json:"alias"`
	ApiUrl         string        `json:"api_url"`
	HttpStatusCode int           `json:"http_status_code"`
	Latency        time.Duration `json:"latency"`
	IsSuccess      bool          `json:"is_success"`
	ErrorMessage   string        `json:"error_message,omitempty"`
}
```

---

## 3. Implementation Checklist

- [ ] **1. Vault Types & Models (`cli/cmdsupabase/supabase_types.go`):**
  - Define `SupabaseProjectRecord`, `SupabaseProjectDecrypted`, `SupabaseProjectView`, and `SupabasePingResult`.
  - Define status constants: `StatusActive = "active"`, `StatusDisabled = "disabled"`, `StatusError = "error"`.
- [ ] **2. AES-GCM Encryption Engine (`cli/cmdsupabase/supabase_vault.go`):**
  - Implement `DeriveMasterVaultKey() ([]byte, error)`.
  - Read or generate installation salt at `.gitmap/data/installation/vault.salt`.
  - Combine machine fingerprint (UUID / Hostname) with salt via HMAC-SHA256 to produce 32-byte AES key.
  - Implement `EncryptSecret(plaintext string) (string, error)` calling `crypto.Encrypt([]byte(plaintext), key)`.
  - Implement `DecryptSecret(ciphertext string) (string, error)` calling `crypto.Decrypt(ciphertext, key)`.
  - Implement `MaskSecret(secret string) string` (e.g. keeping first 6 and last 4 chars, replacing middle with `****`).
- [ ] **3. SQLite Storage Engine (`cli/cmdsupabase/supabase_db.go`):**
  - Implement `OpenSupabaseDB() (*sql.DB, error)` resolving `.gitmap/data/installation/supabase/sql.db`.
  - Implement schema initialization with `CREATE TABLE IF NOT EXISTS supabase_databases`.
  - Implement `InsertProject(record SupabaseProjectRecord) error`.
  - Implement `GetProjectByAlias(alias string) (*SupabaseProjectRecord, error)`.
  - Implement `ListAllProjects() ([]SupabaseProjectRecord, error)`.
  - Implement `DeleteProjectByAlias(alias string) error`.
- [ ] **4. CLI Subcommand Handlers (`cli/cmdsupabase/supabase_cmd.go`):**
  - Implement `Run(args []string) error` with command router.
  - Implement `handleAdd(args []string) error`: Validate parameters, encrypt keys, store in DB.
  - Implement `handleList(args []string) error`: Fetch projects, format masked views, render ANSI table.
  - Implement `handleGet(args []string) error`: Retrieve project, check `--decrypt` flag before decrypting.
  - Implement `handleTest(args []string) error`: Decrypt key, execute HTTP probe to `<api_url>/rest/v1/`, report latency and status.
  - Implement `handleRemove(args []string) error`: Delete project with confirmation.
  - Implement `handleExportEnv(args []string) error`: Print export commands for shell evaluation.
- [ ] **5. Command Routing & Help Integration (`cli/cmd/rootdata.go`, `cli/cmd/roottooling.go`):**
  - Wire `supabase` and `sb` into `dataDatabaseEntries()` in `cli/cmd/rootdata.go`.
  - Add `supabase` documentation into CLI help texts.
- [ ] **6. Unit & Integration Testing (`cli/cmdsupabase/supabase_test.go`):**
  - Test encryption and decryption round-trip.
  - Test secret masking logic.
  - Test SQLite CRUD operations in an isolated temporary database.
  - Test validation error handling for missing arguments.

---

## 4. Acceptance Criteria

1. Running `gitmap supabase add test-proj https://xyz.supabase.co eyJanon eyJservice "postgresql://postgres:secret@db:5432/postgres"` successfully stores the entry.
2. In SQLite, `AnonKeyEnc`, `ServiceKeyEnc`, and `DbUrlEnc` are stored strictly as encrypted base64 ciphertexts. Cleartext passwords or keys are NEVER present in SQLite.
3. Running `gitmap supabase list` outputs an ANSI table showing masked keys (`eyJano...****`).
4. Running `gitmap supabase get test-proj` without `--decrypt` outputs masked secrets; with `--decrypt` outputs decrypted secrets.
5. Running `gitmap supabase test test-proj` performs an authenticated HTTP probe and outputs latency and status code.
6. Running `gitmap supabase remove test-proj` deletes the record.
7. All paths referenced in output and storage are strictly relative Git paths.

---

## 5. Verification Commands

```bash
# Run unit test suite
go test -v ./cli/cmdsupabase/...

# Verify build compilation
go build -v ./cli/...

# Manual verification tests
gitmap supabase help
gitmap supabase list
```
