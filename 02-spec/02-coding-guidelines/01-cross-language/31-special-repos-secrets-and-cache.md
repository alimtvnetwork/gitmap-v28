# 31. Special Default Repositories (`repo-secrets` = `rs` & `repo-cache` = `rc`)

## 1. Purpose & Architecture

To prevent credential leakage in public/standard repositories and eliminate disposable script clutter while preserving reusable engineering scripts across projects, GitMap provisions two fixed default companion repositories in the work directory (customizable via `gitmap settings`):

| Special Repo | Shortcut | Navigation | Default Setting Key | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| **`repo-secrets`** | `rs` | `gitmap cd rs` | `special_repos.secrets_name` | Stores `.env` files, tokens, passwords, and config secrets under `XX-<repo-name>/01-<slug>.ext`. |
| **`repo-cache`** (`repo-storage`) | `rc` | `gitmap cd rc` | `special_repos.cache_name` | Stores reusable `.ps1` PowerShell scripts, test fixtures, and temporary verification harnesses under `XX-<repo-name>/01-<slug>.ps1`. |

## 2. Sequenced Folder & File Convention (`XX-<repo-name>/NN-<slug>.ext`)

Inside both `repo-secrets` and `repo-cache`:
- Each source repository receives a monotonically sequenced folder (`01-gitmap`, `02-coding-guidelines`, `03-prompt-architect`).
- Every stored file, folder, or text note inside `<XX-repo-name>/` receives a two-digit prefix (`01-`, `02-`, `03-`) and is **automatically committed and pushed** upon execution.

## 3. Mandatory AI Agent Workflow

1. **Offloading Secrets (`gitmap rs`):**
   - When an AI agent encounters or generates a secret file, credential JSON, or password:
     - `gitmap rs file .env`
     - `gitmap rs folder ./secrets`
     - `gitmap rs text "PROD_DB_PASS=..." --slug prod-db-pass`
2. **Saving Reusable Temporary Scripts (`gitmap rc`):**
   - When an AI agent writes a temporary `.ps1` script, migration helper, or E2E test harness that may be useful to future runs or other repositories:
     - `gitmap rc file ./verify-cluster.ps1`
     - `gitmap rc folder ./test-harnesses`
     - `gitmap rc text "Get-Process | Select-Object -First 5" --slug check-procs --ext .ps1`
