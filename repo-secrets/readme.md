# 🗝️ Centralized Repository Secrets & Operational Vault

Welcome to `repo-secrets` — the centralized external secrets vault and operational scratchpad companion for all ecosystem projects.

---

## 1. What This Repo Is About

- **Pure Secrets & Configuration Vault**: This repository is an external vault storing environment credentials, API keys, passwords, database connection strings, server node configurations, and operational scratchpads across all projects.
- **NOT an Application Codebase**: There is **no application source code**, no compiler build targets, no package managers (`npm`, `cargo`, `pip`), and no unit test suites.
- **Core Mission**: Decouple secrets from code. Its mission is keeping all project codebases completely clean of hardcoded credentials, sensitive tokens, or committed `.env` files.

---

## 2. How to Handle the `.ai-memory/` Folder

According to universal coding guidelines, operational memory and secrets must remain strictly separated:

> [!CAUTION]
> ### 🚨 CODE RED: ZERO Secrets in `.ai-memory/`
> Coding guidelines strictly ban putting secrets, tokens, or passwords into `.ai-memory/` in **ANY** repository. Placing credentials into `.ai-memory/` is a catastrophic security violation.

- **No Application `.ai-memory/` in This Repo**: Because `repo-secrets` is a pure secrets vault and configuration store, it intentionally does **not** have or need a standard application `.ai-memory/` tree (no `plans/pending/`, `memories/`, `memory/issues/`, etc.).
- **Where AI Scratchpads Go**: If an AI agent needs temporary working memory, execution logs, or scratch notes while working here, it must place them inside the target project's `temp/` folder (e.g., `01-gitmap/temp/`) or the operating system's temporary directory (`%TEMP%`).
  - **NEVER** dump scratch files at the repository root.
  - **NEVER** store unencrypted production secrets in plain markdown notes or scratchpads.
- **Strict Relative Paths**: All citations, links, and cross-references must use strictly relative paths from the repository root (e.g. `01-gitmap/04-ssh-nodes.json`). Absolute paths (such as local drive paths) and URI file schemes are strictly prohibited in committed documentation.

---

## 3. Repo Naming Standard: `xx-<repo-name>`

Every repository managed in this vault must adhere to a standardized sequence naming convention:

> [!IMPORTANT]
> **Mandatory Rule**: Any repo containing secrets **MUST** have `xx-` in the folder name.
> - **Format**: `xx-<repo-name>` (e.g. `01-gitmap`, `02-antigravity-manager`, `03-supabase`, `04-w1-machine`)
> - **Sequence Prefix**: Two-digit prefix (`01-`, `02-`, etc.) to establish deterministic ordering.
> - **Casing**: Strictly lowercase kebab-case. No uppercase letters, no underscores, and no spaces.

### Adding a New Repo
When onboarding a new repository, allocate the next available sequential two-digit prefix:
```text
08-<new-repo-name>/
├── vault/      # Encrypted secrets and credential configs
├── temp/       # In-flight agent scratchpads and queue state
└── readme.md   # Specific documentation for this repo's secrets
```

---

## 4. Guidelines for Placing Secrets vs. Temporary Items

Each `xx-<repo-name>/` directory maintains a strict boundary between long-lived credentials and ephemeral runtime state:

| Category | Target Location | Description & Allowed Formats |
|:---|:---|:---|
| **SECRETS** | `<namespace>/vault/` or credential files | Permanent secrets, encrypted DBs, JSON manifests, and Base64 configs (e.g., `email_passwords.db`, `accounts.json`, `supabase-credentials.json`). |
| **TEMPORARY ITEMS** | `<namespace>/temp/` | Ephemeral scratchpads, in-flight prompt queues, assignment inventories, migration blueprints, and delta audit logs (e.g., `file-assignments.json`, `agy-prompt-queue.json`). |

### Rules of Separation
- **NEVER** dump temporary files at the repository root.
- **NEVER** mix temporary files into `vault/`.
- **NEVER** place unencrypted plaintext secrets into `temp/`.

---

## 5. Current Repository Directory Tree

```text
repo-secrets/
├── setup-telegram.ps1                  # Telegram notification setup and test script (PowerShell)
├── setup-telegram.sh                   # Telegram notification setup script (Bash)
├── connect-supabase.ps1                # Supabase PostgREST connectivity verification script
├── vmpass.json                         # Root convenience: Shared VM credentials
├── gitmap-ssh-nodes.json               # Root convenience: Cluster SSH fleet nodes manifest
├── gitmap-final.json                   # Root convenience: Master repository scan inventory
├── readme.md                           # Master vault documentation (this file)
│
├── 01-gitmap/                          # GitMap CLI automation secrets & pipeline configs
│   ├── vault/ & config files           # Commit pull rules, git profiles, ssh-nodes, seo templates
│   ├── summaries/                      # Version metadata replay index (00-index.json .. 28-gitmap-v28.json)
│   └── temp/                           # Migration blueprints & template drafts
│
├── 02-antigravity-and-event-manager/   # Antigravity Manager full vault, broadcaster & deploy scripts
│   ├── scripts/                        # Telegram daemons & connectivity runners
│   ├── vault/                          # Google AI accounts, email DBs, security DB, Telegram configs
│   └── temp/                           # Prompt queues, agent file assignments, and runner ETA logs
│
├── 02-antigravity-manager/             # Lean companion AGM vault structure
│   ├── vault/                          # Core credentials, email DBs, and config files
│   └── temp/                           # Runtime prompt queues and delta logs
│
├── 03-supabase/                        # Supabase platform suites and synchronization settings
│   ├── 01-own/                         # Personal Supabase instance credentials & endpoints (Base64)
│   └── 02-lovable/                     # Lovable platform integration credentials & metadata (Base64)
│
├── 04-w1-machine/                      # Worker 1 (192.168.1.3) configuration snapshots & O&O ShutUp10 profile
├── 05-w2-machine/                      # Worker 2 (192.168.1.7) configuration snapshots & O&O ShutUp10 profile
├── 06-w3-machine/                      # Worker 3 (192.168.1.12) configuration snapshots & O&O ShutUp10 profile
├── 07-final-network-machine/           # Final Controller (192.168.1.20) master config, clone-gitmap.ps1 & manifests
├── 08-licenses/                        # Commercial software keys and license credentials
├── 09-antigravity-backup/              # Antigravity IDE project manifests, pinned configurations, cross-platform restoration automation, and migration runbook
│   ├── vault/                          # Portable JSON manifests (${WORKSPACE_ROOT} tokens for 78 repos, 23 pins, settings, plugins)
│   ├── scripts/                        # Zero-touch cross-platform restoration & verification suite (sh / ps1)
│   └── temp/                           # Ephemeral staging and migration scratchpads
└── vmware/                             # VMware Workstation & ESXi management automation (manage-vm.ps1)
```

---

## 6. How Any AI Agent Should Work in This Repo

Any AI agent operating within `repo-secrets` must follow this 6-point checklist:

1. **Locate the Assigned Repo Folder**: Identify the specific target folder using the `xx-<repo-name>` sequence prefix (e.g., `01-gitmap/`, `03-supabase/`).
2. **Access Secrets in `vault/`**: Read and write credentials only in designated `vault/` directories or defined credential files (e.g., `supabase-credentials.json`). Respect encoding standards (e.g., Base64 or AES-256 DBs).
3. **Stage Scratchpads Exclusively in `temp/`**: Place all temporary prompt queues, agent assignments, execution notes, and scratchpads in the target repo's `temp/` subfolder.
4. **Enforce CODE RED on `.ai-memory/`**: Never create standard `.ai-memory/` trees in this repository, and never write credentials or sensitive tokens into `.ai-memory/` files in any project.
5. **Use Strictly Relative Paths**: Always reference files using relative paths from the repo root (e.g., `02-antigravity-manager/vault/accounts.json`). Never output absolute paths or file URI schemes.
6. **Maintain Cleanliness & Zero Clutter**: Never drop scratch files or logs at the repository root. Keep root convenience manifests synchronized and leave the repo free of untracked debris.
