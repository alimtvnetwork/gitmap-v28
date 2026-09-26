# Specification: Create Repo from Existing Folder, Recreate Repo, and Common Baseline Command

> **Specification ID:** 147-recreate-repo-and-folder-creation-specification  
> **Status:** APPROVED  
> **Category:** 21-app  
> **Sponsor:** Rise Up Asia LLC (Wyoming, California & New York)  
> **Principal Architects:** Marek Flejszman & Alim Ul Karim  

---

## 1. Scope & Objectives

This specification covers:
1. `gitmap create-repo` (alias `cr`) support for existing folders containing pre-existing user files.
2. `--common` baseline append and commit behavior.
3. `--cg` / `--coding-guideline` synchronization with mandatory pre-flight backup branch.
4. `gitmap common` command architecture (with `sync` / `sy` backward compatibility).
5. `gitmap recreate-repo` (alias `recreate`) remote re-anchoring on GitHub.
6. `gitmap var` persistent variable engine.

---

## 2. Technical Specification

### 2.1 `gitmap create-repo [path] [--common] [--cg]`
- **Alias:** `cr`.
- **Existing Directory Support:**
  - If executed inside an existing folder (`.` or specified path) containing files:
    - Automatically uses current folder name as repository name and slug if not explicitly passed.
    - Preserves all pre-existing files without blind overwriting (`README.md`, `.gitignore`, etc.).
    - Initializes Git if `.git` is absent.
    - Stages all files and creates an initial commit (`Initial commit with existing workspace files`).
  - **`--common` flag:**
    - Appends curated standard baselines (.gitignore, .gitattributes, git-lfs setup, .prettierignore, .prettierrc).
    - Commits as a clean sub-commit: `chore: apply common repository baselines`.
  - **`--cg` / `--coding-guideline` flag:**
    - Creates an automatic pre-flight backup branch: `backup/cg-sync-<timestamp>` before modifying any files.
    - Preserves user on `main` branch (`git branch <backup-name>` without switching).
    - Copies latest repository coding guidelines into `02-spec/02-coding-guidelines/`.
    - Commits as: `docs(guidelines): synchronize standard coding guidelines`.

### 2.2 `gitmap common [target] [flags]`
- **Renamed from:** `gitmap sync`.
- **Backward Compatibility:** `sync` and `sy` remain fully functional aliases.
- **Bare Invocation:** Defaults to `all`, running the entire baseline suite in sequence.
- **Targets:**
  - `all`: Full baseline setup.
  - `ignore`: Union-merges curated `.gitignore`.
  - `attributes`: Union-merges `.gitattributes`.
  - `lfs-install`: Sets up Git LFS file tracking.
  - `prettier-ignore`: Union-merges `.prettierignore`.
  - `prettier-rc`: Merges `.prettierrc` JSON options.
- **Flags:** `-n` / `--dry-run`, `-f` / `--force`, `-h` / `--help`.

### 2.3 `gitmap recreate-repo <name>`
- **Alias:** `recreate`.
- **Pre-flight Guard:** Creates `backup/recreate-<timestamp>` before remote modification.
- **GitHub Re-anchoring:** Uses `gh repo create <owner>/<name>` to publish the repository to GitHub, configures upstream tracking, and executes workspace sync.

### 2.4 `gitmap var <subcommand> [flags]`
- **Subcommands:** `set <key> <value>`, `get <key>`, `ls`, `rm <key>`, `export`.
- **Scoping:** `--scope <name>` (e.g. `--scope agy`). Defaults to `global`.
- **Storage:** `.gitmap/variables.json` with thread-safe file locking.
- **Expansion:** `config.ExpandVariables(input, scope)` expands `$VAR` and `${VAR}`.
- **Environment Export:** `export` writes variables to OS environment variables.
