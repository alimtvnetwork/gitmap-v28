# Secret Management

## 1. Zero Hardcoded Secrets

- Code MUST NEVER contain hardcoded secrets, API keys, passwords, or tokens.
- `.env` files containing secrets MUST NOT be committed to version control. Ensure `.env` is in `.gitignore`.

## 2. Secret Vaults

- Production environments must retrieve secrets dynamically at runtime or during the deployment phase from a secure vault (e.g., HashiCorp Vault, AWS Secrets Manager, Azure Key Vault).

## 3. Secret Rotation

- Services must be designed to handle dynamic secret rotation without requiring a full application restart if possible, or gracefully restart when the orchestrator updates the injected secrets.

## 4. Special Repository for Secrets (`repo-secrets` / `gitmap rs` / `gitmap cd rs`)

- All secrets, private credentials, authentication email and password pairs, `.env` files, API credentials, and private tokens MUST NEVER be committed to standard source repositories.
- **Context & Location Mandate:** Never specify or provide repository URLs, git remote URLs, or absolute folder paths. State strictly that if the `repo-secrets` folder exists in the default work directory, that is the context where secrets MUST be stored.
- If the `repo-secrets` folder exists in the default work directory, secrets MUST be stored into it using GitMap CLI commands or dedicated companion storage:
  - `gitmap rs file <filepath> [--repo <name>]`: Copies the secret file into `repo-secrets/XX-<repo-name>/01-<slug>.ext`, automatically commits, and pushes.
  - `gitmap rs folder <folderpath> [--repo <name>]`: Recursively archives a secret directory into `repo-secrets/XX-<repo-name>/01-<slug>/`, commits, and pushes.
  - `gitmap rs text "<secret-token>" [--slug <slug>] [--ext <ext>]`: Writes inline secrets or tokens into `repo-secrets/XX-<repo-name>/01-<slug>.ext` with deterministic numbering, commits, and pushes.
  - Jump directly to the secrets repository: `gitmap cd rs`.
  - Configurable repository name via settings: `gitmap settings set special_repos.secrets_name repo-secrets`.

## 5. Special Repository for Reusable & Scratch Scripts (`repo-cache` / `gitmap rc` / `gitmap cd rc`)

- Reusable temporary scripts (e.g. PowerShell `.ps1` automation, diagnostic test harnesses, scratch utilities) MUST NOT pollute application source repositories.
- Archive temporary scripts and test fixtures into `repo-cache` (`repo-storage`) using GitMap CLI commands:
  - `gitmap rc file <script.ps1> [--repo <name>]`: Archives script file into `repo-cache/XX-<repo-name>/01-<slug>.ps1`, automatically commits, and pushes.
  - `gitmap rc folder <folderpath> [--repo <name>]`: Archives test fixture or tool directory into `repo-cache/XX-<repo-name>/01-<slug>/`, commits, and pushes.
  - `gitmap rc text "<powershell-script>" [--slug <slug>] [--ext .ps1]`: Writes inline PowerShell snippet into `repo-cache/XX-<repo-name>/01-<slug>.ps1` with auto-sequencing, commits, and pushes.
  - Jump directly to the cache repository: `gitmap cd rc`.
  - Configurable repository name via settings: `gitmap settings set special_repos.cache_name repo-cache`.
