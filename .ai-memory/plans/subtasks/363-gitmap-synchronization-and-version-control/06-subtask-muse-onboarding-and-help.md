# Subtask 06: Muse Onboarding & Muse Help Terminal Command Engine

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-06`
- **Slug:** `gitmap-synchronization-and-version-control - Task 06`
- **Target Package:** `cli/cmdmuse/`, `cli/helpdoc/`
- **Status:** `pending`
- **Assigned Worker:** Muse Integration Subagent

---

## 1. Goal

Implement the native Muse AI agent onboarding engine and terminal help command suite in GitMap. Provide streamlined bootstrapping via `gitmap muse onboard` and curated, interactive terminal help via `gitmap muse help` (aliased to `gitmap thelp muse`) to guide operators through Muse memory IO, training vault synchronization, and LLM digest ingestion.

---

## 2. In Scope

1. **Muse Onboarding Engine (`cli/cmdmuse/muse_onboard.go`):**
   - Command: `gitmap muse onboard` (aliases: `gitmap muse init`, `gitmap muse setup`).
   - Automated environment verification:
     - Check local AI Auto login status; if unset, auto-register `gitmap login set ai-auto muse`.
     - Inspect `repo-secrets` special repository for `Muse/` structure (`Muse/README.md`, `Muse/01-*.md`, `Muse/chat-seeds/`).
     - Bootstrap local `.ai-memory/muse/` cache layout if missing.
     - Validate training vault access and digest freshness.
   - Flags:
     - `--check-only`: verify environment readiness without creating directories.
     - `--force`: re-initialize cache directories and reset stale seed pointers.
     - `--json`: emit JSON readiness report.

2. **Muse Rich Terminal Help Engine (`cli/cmdmuse/muse_help.go`):**
   - Command: `gitmap muse help` (and `gitmap thelp muse`, `gitmap muse --help`).
   - Dual-column Tokyo Night card interface:
     - Command catalog: `muse onboard`, `muse train`, `muse status`, `muse verify`.
     - Memory IO digest workflows (`gitmap llm train muse [--days N]`).
     - Multi-day diff inspection and AI traversal checklists.
     - Curated executable examples and troubleshooting tips.

3. **Subcommand Dispatcher Registration (`cli/root.go`):**
   - Register `muse` top-level verb in root dispatcher.
   - Wire `muse help` into `cli/richhelp/rich_help_dispatcher.go`.

---

## 3. Out of Scope

- Direct fine-tuning or neural network inference (delegated to LLM engines).
- Running external model training loops during onboarding.
- Full builds or test suite execution.

---

## 4. Architectural Invariants

- **Leaf Error Sovereignty:** All file I/O and store queries return `*apperror.AppError`. Do not call `os.Exit` inside `cmdmuse`.
- **Positive Booleans:** `isReady`, `isCheckOnly`, `hasMemoryVault`, `hasDigest`, `hasSeeds`.
- **Function Line Limit:** Strict 8–15 lines max. Decompose verification steps into discrete check functions.
- **Relative Paths:** All created directories and pointers must use 100% relative paths.

---

## 5. Verification Checklist

- [ ] `gitmap muse onboard` verifies prerequisites, registers identity, and bootstraps memory layout.
- [ ] `gitmap muse onboard --check-only` reports environment status without modifying filesystem.
- [ ] `gitmap muse help` and `gitmap thelp muse` render complete Tokyo cards with executable examples.
- [ ] Root dispatcher routes `muse` and all subcommands cleanly.
- [ ] All check functions adhere to positive boolean naming and error sovereignty.
