# Plan 34: Generic AI Scripts Creator and Scaffolding Engine

> **Origin & Execution Context:** Started from user verification request inquiring about operational status of SSH update, SSH join/exec with live terminal UI, cluster/sc command parity, SSH authorized key deployment, ZSH automation, and Antigravity error/prompt injection, combined with a mandate to analyze `03-ai-scripts/` and implement a native, generic AI scripts creation and scaffolding engine in GitMap CLI (`gitmap ai create` / `gitmap ai new`).
> **Loops to Complete:** Completed in 3 continuous self-loops across 3 subtasks without a single test/build failure or execution stall.

## Consolidated Outcomes & Architectural Delivery

### 1. Verification of Prior Systems
- **SSH Fleet Update:** Confirmed `gitmap ssh update --all` / `-a` correctly parses arguments and updates GitMap across all registered cluster/SSH nodes.
- **SSH Exec & Join (`se` / `sj`):** Confirmed multi-node and single-node command dispatch with rich terminal output, status spinners, middle-ellipsized table formatting (`termtable`), and universal `*apperror.AppError` envelopes.
- **Cluster & Servers-Clients (`sc`) Parity:** Confirmed 100% command parity across `node`, `bootstrap`, `k8s`, `init`, `exec`, `install`, `import` between `gitmap cluster` and `gitmap sc`.
- **SSH Authorized Key Deployment:** Confirmed `gitmap ssh auth-key deploy <target>` discovering `-i <path>` or standard keys (`id_ed25519.pub`, `id_rsa.pub`, `id_ecdsa.pub`), enforcing `0700` and `0600` permissions with `grep -qF` deduplication.
- **ZSH Theme Management & Automation:** Confirmed `gitmap zsh themes` (cataloging 41 Oh-My-Zsh themes from Adem's Kubernetes automation suite), `gitmap zsh set-theme <theme>`, and node user scaffolding with Oh-My-Zsh and `zsh-autosuggestions`.
- **Antigravity (AGY) Integration:** Confirmed `gitmap agy fix-pipeline` / `gitmap pipeline fix-agy` error log and context injection into Antigravity IDE/CLI, `gitmap prompt-template` rerun suites, and prompt queue batching.

### 2. Code Analysis of `03-ai-scripts/` & Scaffolding Engine
Analyzed all 44 Python scripts in `03-ai-scripts/` and extracted 5 standard archetypes:
- **Linter (`TemplateLinter`):** Read-only pattern and AST validation with exit 0/1 contracts and optional multi-threaded workers.
- **Fixer (`TemplateFixer`):** Two-pass non-destructive in-place code/text transformation with `--fix` application mode.
- **Auditor (`TemplateAuditor`):** Cross-language parity and spec verification.
- **Generator (`TemplateGenerator`):** Single-source-of-truth code/schema generator.
- **Util (`TemplateUtil`):** General automation script using `02-shared-engine.py` utilities.

### 3. Native AI Script Creator & Dynamic Catalog Discovery
- Implemented `gitmap ai create <name>` (aliases: `new`, `scaffold`, `gen`) in `cli/cmdai/ai_create.go`, `cli/cmdai/ai_cmd_create.go`, and `cli/cmdai/ai_create_templates.go`.
- Added options: `--type` (`linter`, `fixer`, `auditor`, `generator`, `util`), `--desc`, `--parallel`, `--fix`, `--dry-run`, `--force`.
- Implemented dynamic script scanning in `cli/cmdai/ai_catalog_discover.go` and `cli/cmdai/ai_catalog_meta.go` so that newly created scripts in `03-ai-scripts/` are dynamically discovered and immediately runnable via `gitmap ai run <name>` and listed in `gitmap ai list` without hardcoding.
- Added comprehensive unit tests in `cli/cmdai/ai_create_test.go`.

### 4. Help Text & Documentation Parity
- Updated `cli/helptext/ai.md` with syntax, options, and examples for `gitmap ai create`.
- Updated `cli/helptext/catalog.go` and `cli/helptext/print.go` with subcommand summaries and aliases.
- Synchronized documentation across `llm.md`, `cli/llm.md`, and `cli/helptext/llm.md`.
