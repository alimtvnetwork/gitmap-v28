# 31. Special Default Repositories (`repo-secrets` = `rs` & `repo-cache` = `rc`) (AI Execution Prompt)

> **/goal** Prevent credential leakage into public git repositories and eliminate throwaway script clutter by offloading secrets to `repo-secrets` and archiving reusable engineering scripts to `repo-cache`.
> **/learn** Master GitMap special repository shortcuts (`gitmap rs`, `gitmap rc`), enforce monotonic sequence prefixes (`XX-<repo-name>/NN-<slug>.ext`), prohibit hardcoded repository URLs or absolute paths for secrets, and leverage default work directory context.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Offload all passwords, API keys, tokens, and `.env` configurations to `repo-secrets` using `gitmap rs`.
- [ ] `/learn` Never commit credentials or secrets to standard or public repositories; never hardcode remote URLs or absolute paths in secret configurations.
- [ ] `/goal` Archive temporary verification scripts, diagnostic harnesses, and migration helpers into `repo-cache` using `gitmap rc`.
- [ ] `/learn` Enforce monotonic sequencing (`XX-<repo-name>/NN-<slug>.ext`) and automated commit/push behavior inside special repositories.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Status:** Active
**Scope:** Universal (All Repositories, Workflows & Languages)
**Strictness:** Non-Negotiable (Banned Anti-Pattern)

---

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

1. **Offloading Secrets (`gitmap rs` / Default Work Directory Context):**
   - Context & Location Mandate: Never specify or hardcode repository URLs, git remote URLs, or absolute folder paths for secrets. State strictly that if the `repo-secrets` folder exists in the default work directory, that is the context where secrets MUST be stored.
   - If anything is secret, private, or contains authentication credentials (e.g. email, password, API tokens, private keys, `.env` configurations), it MUST be stored inside the `repo-secrets` folder in the default work directory, and MUST NEVER be committed to standard source code repositories:
     - `gitmap rs file .env`
     - `gitmap rs folder ./secrets`
     - `gitmap rs text "PROD_DB_PASS=..." --slug prod-db-pass`
   - If the `repo-secrets` folder exists in the default work directory, offload and store credentials there immediately. If it does not exist, credentials must never be written in plain text or pushed to public/standard repositories.
2. **Saving Reusable Temporary Scripts (`gitmap rc`):**
   - When an AI agent writes a temporary `.ps1` script, migration helper, or E2E test harness that may be useful to future runs or other repositories:
     - `gitmap rc file ./verify-cluster.ps1`
     - `gitmap rc folder ./test-harnesses`
     - `gitmap rc text "Get-Process | Select-Object -First 5" --slug check-procs --ext .ps1`

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ARCH-031: Credential Offloading to `repo-secrets` and Script Archival to `repo-cache`

**Given** Application configurations, secret tokens, environment files, and temporary automation scripts across projects.
**When** Security linters, secret scanners, and repository hygiene auditors inspect the codebase.
**Then** Standard repositories contain zero committed secrets or plain-text credentials, secrets are segregated into `repo-secrets` with monotonic sequencing, and reusable temporary harnesses are preserved in `repo-cache`, achieving exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
