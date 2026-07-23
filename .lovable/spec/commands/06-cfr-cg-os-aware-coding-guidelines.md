# Command: cfr / cfrp with `cg` (Coding Guidelines) modifier

Slug: cfr-cg-os-aware-coding-guidelines
Status: captured
Created: 2026-07-16

## Verbatim user command

> cg = Coding Guidelines
> cfr cg = Clone repository with Coding Guidelines
> cfr p cg = Clone a public repository with Coding Guidelines
>
> When someone runs `cfr cg`, it should:
> 1. Clone the repository.
> 2. Detect the operating system (Windows or Linux/macOS).
> 3. Execute the appropriate Coding Guidelines setup for that OS.
> 4. Either auto-commit + push the initial Coding Guidelines changes, OR simply copy the Coding Guidelines into the project folder.

## OS-specific installers (canonical URLs)

- Windows (PowerShell):
  `irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/error-manage-install.ps1 | iex`
- Linux / macOS (bash):
  `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh | bash`

## Scope

- Applies to `gitmap cfr`, `gitmap cfrp` (public variant), and any future `clonefixrepo` entry points.
- The `cg` token is a modifier flag position, not an alias: `cfr cg <url>`, `cfr p cg <url>`, `cfr cg p <url>` should all be accepted.
- Existing `install clean-code` / `install cg` (see `installcleancode.go`) already targets `coding-guidelines-v15`; the new v24 URLs must NOT overwrite that alias. `cfr cg` is a distinct integration.
- Default post-install behavior: copy files into working tree, stage + commit as `chore: install coding guidelines (v24)`, push to `origin <current-branch>` only when the repo has a tracked upstream. Skip push with `--no-push`; skip commit with `--no-commit`.

## When it applies

Every future coding task that references `cfr`, `cfrp`, or the `cg` modifier. Windows uses PowerShell (prefer `powershell`, fallback `pwsh`). Non-Windows uses bash via `curl | bash`. Never invoke the wrong installer for the host OS.
