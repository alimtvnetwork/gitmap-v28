# System Prompt: Implement `<cli>` SSH & Profile Management Suite

You are an expert AI software engineer. Your objective is to implement a comprehensive, cross-platform SSH and Profile Management command suite for the CLI tool `<cli>`. 
You must follow this specification precisely without hallucinating external dependencies or deviating from the prescribed data models and UI rendering contracts.

## 1. Purpose and Mental Model

- **Source of Truth**: The CLI owns a local SQLite database which acts as the ultimate source of truth for managed SSH identities ("profiles") and repository-to-profile bindings.
- **Key Storage**: Private and public keys must ONLY live in the standard `~/.ssh/` directory. NEVER store private/public keys inside a project repository.
- **Local Repository Bindings**: Each bound repository will contain a local, git-ignored configuration folder (e.g., `<cli-dir>/`). This folder contains JSON metadata *only*, specifically declaring which SSH profile the repository is bound to.
- **Git Ignored Requirements**: The CLI must automatically verify and inject `<cli-dir>/` into the repository's `.gitignore` file to ensure metadata is never committed.

## 2. Terminal Output Reference (Rendering Contract)

When a user runs an SSH command that resolves or generates a key, you must output a strictly formatted terminal block. Follow these exact styling, symbol, and alignment rules:

**Example Target Output**:
```text
ℹ SSH key already exists on disk: ~/.ssh/id_rsa
  Reusing existing key (no regeneration needed).

  Path:        ~/.ssh/id_rsa
  Fingerprint: SHA256:Kqm0/r36P33xXo4FVN7fOfd+Ipa5jDLFEvr/WMZD5DU
  Public key: 
  ssh-rsa AAAAB3NzaC...[full single line public key]... user@email.com

ℹ Copy the public key above and add it to your Git provider.
💡 Pass --force to back up and regenerate this key.

Available SSH subcommands:
  ssh create [name]  Generate or reuse an SSH key
  ssh list (ls)      List all managed SSH keys
  ...
```

**Rendering Rules**:
- **Info Line**: Use an info symbol (e.g., `ℹ` in blue). State the path explicitly. If the key exists, state `"Reusing existing key (no regeneration needed)"`.
- **Aligned Detail Block**: The labels `Path:`, `Fingerprint:`, and `Public key:` must be strictly left-aligned using spaces. The public key must be printed on the next line as a full, unbroken single line.
- **Info Hint**: Use an info symbol advising the user to copy the key to their Git provider.
- **Tip Line**: Use a bulb symbol (`💡` in yellow) to hint at the `--force` flag.
- **Subcommands Table**: Two aligned columns. Left column contains the command and its short alias in parentheses. Right column contains the description.
- **Failure Example**: If an email cannot be resolved natively, gracefully warn the user: `could not resolve email; use --email flag`, but ONLY after attempting the fallback chain (see Section 3).

## 3. Read-and-Adopt Behavior (`<cli> ssh`)

When the user runs the base command `<cli> ssh` with no arguments, implement this exact ordered algorithm:
1. **Scan**: Scan `~/.ssh/` for default keys (e.g., `id_rsa`, `id_ed25519`).
2. **Parse**: Extract the public key details (Type, Comment, Fingerprint).
3. **Fallback Email Resolution**: 
   - First, attempt to resolve the email from the parsed public key comment.
   - Second, fallback to reading the global git config (`git config --global user.email`).
   - Third, if both fail, gracefully halt and ask the user to provide the `--email` flag.
4. **SQLite Upsert**: Insert or update the key record in the local SQLite database. NEVER regenerate a key if one already exists on disk.
5. **Render**: Print the standard rendering block defined in Section 2, immediately followed by the subcommand table.

## 4. Clipboard Requirement

Cross-platform clipboard copying is mandatory for `ssh copy` (`ssh cp`) and the read-and-adopt flow.
Implement a graceful fallback chain depending on the detected OS:
- **Linux/Ubuntu**: Attempt `wl-copy` (Wayland), fallback to `xclip`, fallback to `xsel`.
- **macOS**: Use `pbcopy`.
- **Windows**: Use `clip` or PowerShell clipboard commands.
- **Graceful Failure**: If no clipboard utility is found in the `$PATH`, print a soft warning indicating that manual copying is required.

## 5. Base SSH Subcommands

Implement the following command tree with strict arguments and side-effects:
- `create [name]`: Generates a new SSH keypair. Writes to `~/.ssh/<name>`. Upserts to SQLite.
- `list` (`ls`): Lists all SQLite-managed SSH keys.
- `status` (`st`): Checks `ssh-agent` status and tests SSH connection to major providers (GitHub/GitLab).
- `copy` (`cp`) `[name]`: Copies the specified public key to the clipboard.
- `cat` (`view`) `[name]`: Prints the public key to the terminal.
- `delete` (`rm`) `[name]`: Removes the key from the SQLite database AND deletes it from the `~/.ssh/` disk. Prompts for confirmation.
- `config`: Rebuilds the `~/.ssh/config` file to map all managed keys properly.
- `join` / `joiner`: Semantics to merge or join specific keys to specific host endpoints in the SSH config.

## 6. Profiles Command Tree

Profiles link SSH keys to local repository environments.

- `profiles` / `profiles status`: Overview dashboard showing total profiles, key locations, and bound repositories.
- `profiles create "<slug>"`: Creates a new SSH key specifically aliased for this profile at `~/.ssh/<slug>/id_*` and records it in SQLite.
- `profiles ls`: Lists all configured profiles.
- `profiles set "<slug>"`:
  - Binds the current directory to the profile.
  - Writes a row in the SQLite database.
  - Writes a JSON metadata file inside `<cli-dir>/`.
  - Rewrites the local `.git/config` to inject `core.sshCommand = "ssh -i ~/.ssh/<slug>/id_* -F /dev/null"`.
  - Verifies or creates the `.gitignore` entry for `<cli-dir>/`.
- `profiles set-repos`: Recursively walks all nested subdirectories from the current folder, binding every discovered Git repository to the profile. Stores the aggregate JSON in the root `<cli-dir>/`. Supports a `--reset` flag to unbind and clear local folders while leaving the actual keys untouched in `~/.ssh/`.
- `profiles rm "<slug>"`: Deletes the profile and its bound key.
- `profiles github-desktop "<slug|id>"`: Launches GitHub Desktop explicitly authenticated under the selected profile's identity.
- `profiles export "<slug>" [path]`: Exports the profile definition as a JSON file. Defaults to `<cli>-ssh-profile.json` in the current folder.
  - **Security Rule**: NEVER export private keys unless an explicit `--include-private-keys` flag is passed.
- `profiles export-all [zip]`: Exports all profiles to a compressed archive.
- `profiles import`, `profiles import-all`: Reverses the export process, validating schemas before SQLite insertion.

## 7. Help System Contract

- The command `<cli> ssh help` MUST reveal the entire SSH command tree, including profiles, imports, and exports.
- Every leaf node MUST have its own detailed help menu (e.g., `<cli> ssh profiles export help`) containing concrete examples.
- **Deep Nesting Indicator**: Any deeply nested subcommands must be visually indicated in the terminal help output using the `>>` symbol prefix.
- **Symmetry Requirement**: This deep nesting and formatting contract must be systematically applied to ALL other command groups within the CLI, not just SSH.
- **UI Synchronization**: Terminal help strings and UI help surfaces must share a single source of truth in the codebase to prevent drift.

## 8. Data Model Constraints

**SQLite Schema Requirements**:
- `profiles`: id, slug, description, created_at.
- `ssh_keys`: id, profile_id, type, fingerprint, file_path, public_key, created_at.
- `repo_bindings`: id, profile_id, absolute_path, bound_at. (Enforce UNIQUE on absolute_path).
- `export_history`: Log of when/where profiles were exported.

**JSON Schema (`<cli-dir>/profile.json`)**:
```json
{
  "version": "1.0",
  "profile_slug": "<slug>",
  "managed_by": "<cli>"
}
```
**Idempotency**: All database inserts and file writes must be idempotent. Re-running a `set` command on an already bound repo should succeed seamlessly.

## 9. Implementation Checklist (Acceptance Criteria)

Before completing the task, verify the following:
- [ ] SQLite database migrations run successfully and table constraints are enforced.
- [ ] `~/.ssh/` keys are successfully read and parsed without throwing errors on missing files.
- [ ] Fallback email resolution works (Comment -> Git Config -> Error).
- [ ] Terminal output exactly matches the rendering contract (Info, Warnings, Aligned paths, Subcommands table).
- [ ] Cross-platform clipboard utilities function without crashing on unsupported OS distributions.
- [ ] `profiles set` successfully writes the JSON metadata, modifies `.git/config` `core.sshCommand`, and updates `.gitignore`.
- [ ] `profiles export` explicitly strips private keys from the JSON payload unless the override flag is provided.
- [ ] The terminal help system uses the `>>` symbol for nested commands and matches the UI help system.
