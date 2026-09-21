# `gitmap repo-create` (aliases: `repoc`, `create-repo`, `crepo`, `create-local-repo`, `clr`)

Create and initialize a local repository, automatically slugify human-readable names, and optionally provision and push to GitHub.

## Usage

```bash
# General creation (locally and on GitHub):
gitmap repo-create <name> [folder] [slug] [flags]
gitmap repoc <name> [folder] [slug] [flags]
gitmap create-repo <name> [folder] [slug] [flags]
gitmap crepo <name> [folder] [slug] [flags]

# Local-only creation:
gitmap create-local-repo <name> [folder] [slug] [flags]
gitmap clr <name> [folder] [slug] [flags]
gitmap repo-create --local <name> [flags]
```

## Features
- **Automatic Slugification**: Spaces, underscores, and special characters in `<name>` are automatically converted to a clean GitHub slug (e.g. `"My Cool Service"` → `"my-cool-service"`).
- **Custom Target Folder & Slug**: Specify an explicit target folder or custom slug.
- **Local Only Mode**: Using `create-local-repo`, `clr`, or `--local` / `--no-remote` creates the local git repository (`git init -b main`) without calling GitHub.
- **Auto-Provisioning in Replay Engines**: When using `commit-in`, `commit-left`, `commit-right`, or PR replay commands, any non-existent destination repository is automatically provisioned.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--local`, `--no-remote` | Local git repository initialization only (skip GitHub) | `false` |
| `--public` | Create repository as public on GitHub | `false` (private) |
| `--private` | Create repository as private on GitHub | `true` |
| `--dir <path>` | Destination local directory | `./<slug>` |
| `--slug <slug>` | Explicit repository slug | Auto-generated from `<name>` |
| `-d`, `--description <text>` | Repository description | `""` |
| `--profile <name>` | Target GitHub profile or organization | Active default |
| `--json` | Output creation metadata as structured JSON | `false` |

## Examples

```bash
# 1. Create repository with spaces (auto-slugified to my-new-project):
gitmap repo-create "My New Project"

# 2. Using short alias repoc:
gitmap repoc "Backend API Service"

# 3. Create with explicit folder and slug:
gitmap repo-create "Internal Tool" ./tools/internal internal-tool-v1

# 4. Create local-only repository (skip GitHub):
gitmap create-local-repo "Local Prototype"
gitmap clr "Local Sandbox"

# 5. Commit transfer with auto-provisioned destination:
gitmap commit-right ./source-repo "My New Destination"
gitmap commit-left "New Left Repo" ./source-repo
gitmap commit-in "Consolidated Suite" ./repo-a ./repo-b
```
