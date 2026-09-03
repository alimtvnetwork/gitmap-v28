# Repository Creation Engine (`gitmap create`)

Create and initialize local Git repositories and automatically provision, configure, and push them to remote Git hosting providers (GitHub or GitLab) under the active default account or specified organization.

---

## Command Overview

```bash
gitmap create [repo] <name> [flags]
```

Aliases: `gitmap cr <name>`, `gitmap create repo <name>`

### Key Features
1. **Zero-Configuration Defaults**: Automatically detects and uses the active default Git profile or organization configured in GitMap.
2. **End-to-End Automation**: Creates the local workspace folder, runs `git init -b main`, writes an initial `README.md` and `.gitignore`, creates the initial commit, and provisions the remote repository via `gh` CLI in a single step.
3. **Multi-Account & Organization Support**: Target specific GitHub users or organization accounts with `--profile <name|index>` or `--org <name>`.
4. **Visibility Control**: Default to private repositories for safety, or pass `--public` for open-source repositories.
5. **Local-Only Mode**: Pass `--no-remote` to initialize a pristine local Git repository without publishing to cloud providers.
6. **Machine-Readable**: Use `--json` for automation and CI/CD pipelines.

---

## Flags Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--profile <name\|1-N>` | String / Int | Active profile | Profile name or numeric sequence number from `gitmap profiles ls` |
| `--org <name>` | String | Active profile | Explicit organization name on GitHub / GitLab |
| `--private` | Boolean | `true` | Create repository with private visibility (default) |
| `--public` | Boolean | `false` | Create repository with public visibility |
| `--dir <path>` | String | `./<name>` | Local destination directory path |
| `-d`, `--description <text>` | String | `""` | Repository description included in README and remote metadata |
| `--no-remote` | Boolean | `false` | Skip remote creation (local repository only) |
| `--json` | Boolean | `false` | Output structured JSON metadata |

---

## Practical Examples

### 1. Create a Private Repository Under Default Account

```bash
gitmap create api-gateway --description "Production API Gateway service"
```

Output:
```text
  ✓ Repository created successfully!
  ● Name:      api-gateway
  ● Path:      D:\wp-work\api-gateway
  ● Profile:   alimtvnetwork (github)
  ● Remote:    https://github.com/alimtvnetwork/api-gateway
```

### 2. Create a Public Repository Under an Organization

```bash
gitmap create react-dashboard --org auktvgo --public -d "Next.js UI dashboard"
```

### 3. Select Account by Numeric Sequence Number

```bash
# Target sequence #2 from `gitmap profiles ls`
gitmap create worker-daemon --profile 2
```

### 4. Create Local-Only Repository (No Cloud Remote)

```bash
gitmap create sandbox-experiment --no-remote --dir ./experiments/sandbox
```

### 5. Automated CI/CD Scripting with Structured JSON

```bash
gitmap create ci-artifact-repo --json
```

Output:
```json
{
  "name": "ci-artifact-repo",
  "path": "D:\\wp-work\\ci-artifact-repo",
  "profile": "alimtvnetwork",
  "provider": "github",
  "remoteUrl": "https://github.com/alimtvnetwork/ci-artifact-repo"
}
```
