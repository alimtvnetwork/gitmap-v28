# V2 Prompt — Special Repositories (`repo-secrets` / `rs` & `repo-cache` / `rc`)

## Core Directive for AI Agents

When working across any repository managed by GitMap, you MUST enforce strict separation between **application source code**, **sensitive credentials (`repo-secrets`)**, and **reusable temporary/diagnostic scripts (`repo-cache`)**:

1. **Zero Secrets in Standard Repositories (`repo-secrets` / `gitmap rs`)**:
   - **NEVER** commit `.env` files, API tokens, bot tokens, SMTP passwords, private keys, or credential JSON files into standard source repositories.
   - **ALWAYS** store sensitive files, folders, or inline secret values into `repo-secrets` using `gitmap rs`:
     ```bash
     # Store a secret file (.env, credentials.json, keyfile)
     gitmap rs file <path>

     # Store an inline secret token or configuration block
     gitmap rs text "<content>" --slug <slug>

     # Store a folder of certificates or secret profiles
     gitmap rs folder <path> --slug <slug>
     ```
   - `gitmap rs` automatically places the secret under `XX-<repo-name>/01-<slug>.ext` inside `repo-secrets` and auto-commits & pushes it via `commit-in`.

2. **Reusable Temporary & Diagnostic Scripts (`repo-cache` / `gitmap rc`)**:
   - Whenever you create reusable temporary scripts—such as PowerShell (`.ps1`) automation scripts, diagnostic harnesses, database inspection snippets, or scratch tools—**ALWAYS** archive them into `repo-cache` (`repo-storage`) using `gitmap rc` so any project can reuse them without polluting the active workspace:
     ```bash
     # Store a reusable PowerShell (.ps1) or diagnostic script file
     gitmap rc file <script.ps1>

     # Store an inline PowerShell (.ps1) script directly into repo-cache
     gitmap rc text "<script>" --slug <slug> --ext .ps1

     # Store a folder of scratch utilities or test fixtures
     gitmap rc folder <dir> --slug <slug>
     ```
   - `gitmap rc` automatically sequences items under `XX-<repo-name>/01-<slug>.ext` inside `repo-cache` and auto-commits & pushes them.

3. **Fast Navigation (`gitmap cd rs` & `gitmap cd rc`)**:
   - Navigate immediately to `repo-secrets` or `repo-cache` from any working directory:
     ```bash
     gitmap cd rs
     gitmap cd rc
     ```

4. **Customizing Special Repository Names (`gitmap settings`)**:
   - Inspect or customize the special repository names globally if the workspace uses custom vault/cache repo names:
     ```bash
     gitmap settings set special_repos.secrets_name repo-secrets
     gitmap settings set special_repos.cache_name repo-cache
     gitmap settings
     ```
