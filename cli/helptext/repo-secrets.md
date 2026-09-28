# `gitmap repo-secrets` (`rs`) — Special Secrets Repository Vault & Sequenced Archiver

## Synopsis

```bash
gitmap repo-secrets <subcommand> [arguments] [flags]
gitmap rs <file|folder|text|ls|init|help> [arguments] [flags]
gitmap cd rs
```

`gitmap repo-secrets` (alias `gitmap rs`) manages the dedicated **`repo-secrets`** special repository so secrets, `.env` files, API tokens, SSH keys, and environment credentials are never committed into application source repositories. Every item added via `gitmap rs` is automatically organized under a deterministic 2-digit repository folder (`XX-<repo-name>/`) and 2-digit file/folder sequence (`01-<slug>.ext`), then automatically committed and pushed via `commit-in`.

---

## Sequenced Hierarchy (`XX-<repo>/01-<slug>.ext`)

When storing secrets from any active repository (e.g., `gitmap-v28`), `gitmap rs` resolves or allocates a 2-digit project folder inside `repo-secrets` and increments the inner item sequence automatically:

```text
repo-secrets/
├── 01-gitmap-v28/
│   ├── 01-env-production.env
│   ├── 02-telegram-bot-token.txt
│   └── 03-certs-bundle/
└── 02-webapp-portal/
    └── 01-stripe-webhook-secret.txt
```

- **Project Folder Allocation**: Reuses existing `XX-<repo>` directory if present, or assigns the next 2-digit sequence (`01-`, `02-`, `03-`, ...). Override target project name with `--repo <name>`.
- **Item Sequence Allocation**: Assigns the next 2-digit prefix (`01-`, `02-`, ...) inside `XX-<repo>/` using `--slug <slug>` or the source filename.
- **Auto-Commit & Push**: Automatically runs `commit-in` inside `repo-secrets` unless `--no-push` is specified.

---

## Navigation Shortcut (`gitmap cd rs`)

Jump immediately into the local `repo-secrets` repository from anywhere:

```bash
gitmap cd rs
gitmap cd repo-secrets
```

---

## Examples

### 1. Store a Secret File (`file`)
Copy a `.env`, credential JSON, or certificate file into `repo-secrets/XX-<repo>/01-<slug>.ext` and auto-commit/push:

```bash
gitmap rs file .env.production
gitmap rs file ./configs/service-account.json --slug gcp-sa --repo payment-api
```

### 2. Store a Secret Folder (`folder`)
Recursively archive a directory of secrets or certificates into `repo-secrets/XX-<repo>/01-<slug>/`:

```bash
gitmap rs folder ./certs --slug tls-certs
gitmap rs folder ./secrets-bundle --repo auth-service
```

### 3. Store Inline Secret Text (`text`)
Write a raw token, password, or key directly into a sequenced secret file without creating a local scratch file:

```bash
gitmap rs text "TELEGRAM_BOT_TOKEN=123456:ABC-DEF" --slug telegram-token --ext .env
gitmap rs text "sk-live-987654321" --slug openai-api-key
```

### 4. List Stored Secrets (`ls`)
Inspect all sequenced project folders and secret entries in `repo-secrets`:

```bash
gitmap rs ls
gitmap rs ls --json
```

### 5. Initialize or Verify `repo-secrets` (`init`)
Ensure the local `repo-secrets` directory exists and is initialized as a Git repository:

```bash
gitmap rs init
```

---

## Customizing the Special Secrets Repository Name (`gitmap settings`)

By default, GitMap resolves the repository named `repo-secrets`. Customize the repository folder name globally via `gitmap settings`:

```bash
gitmap settings set special_repos.secrets_name repo-secrets
gitmap settings set secrets_repo my-private-secrets
gitmap settings
```
