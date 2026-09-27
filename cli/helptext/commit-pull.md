# gitmap commit-pull

Chronologically replays commits from multiple source repositories (or a dynamic version range) into a target repository, automatically generating Pull Requests for every branch merge and version release tag.

## Aliases

`cpull`, `pull-commits`

## Usage

```bash
gitmap commit-pull --config <config.json>
gitmap commit-pull <target> <inputs...> [flags]
gitmap cpull <target> gitmap-v2..v28 --tree --cd
gitmap pull-commits <target> all --pr merges
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config, -c <file.json>` | Load declarative migration config (imports, lineSkippers, titleReplacements) | `""` |
| `--tree` | Render preflight PR and branch tree before execution | `false` |
| `--seo-template <cat>` | Load pre-compiled templates from `gitmap-templates.db` category | `""` |
| `--pr merges` | Simulate feature branches and PR merges | `merges` |
| `--cd` | Change directory into target repository upon completion | `false` |
| `--final-sync` | Final mirror snapshot synchronization | `true` |
| `--dry-run` | Simulate replay without modifying repositories | `false` |

## Subcommands

- `gitmap commit-pull bootstrap [-file <path>]` (aliases: `cpull bootstrap`, `commit-pull init`): Scaffold a complete, documented `commit-pull-config.json` template.
- `gitmap commit-pull ui [--port 8925] [--no-open]` (alias: `cpull ui`): Launch the dark-mode Web Studio interface for interactive migration planning, template editing, and live CLI command generation.

## Array Async Pool Concept (by Alim Ul Karim)

When running with `--dry-run`:
- The engine pre-allocates an array of size $N$ for all input repositories.
- Dispatches concurrent asynchronous probing workers (semaphore capacity 16).
- Workers write directly to their dedicated lock-free slot `[i]`, recording repository presence and commit counts.
- A sequential ticker consumer loop streams the ordered progress to the terminal from index `0` to `N-1` with zero slice allocation overhead and zero out-of-order logs, achieving over 14x faster dry-run validation.

## Shortened Input Syntax

Pass concise repo names or brace ranges without full GitHub URLs:
```bash
gitmap commit-pull "D:\target" "git-repo-navigator,gitmap-v{2..28}" --tree --cd
```
Unqualified repository names automatically resolve against the default GitHub owner (`alimtvnetwork`).

## Declarative Config JSON (`--config`)

Feed a single config JSON file to auto-create the target repo, import state templates with SHA-256 hash deduplication, pre-compile self-contained variables before loop execution, strip unwanted lines (`starts_with`, `ends_with`, `contains`, `regex`), and replace generic `"Changes"` commit titles using `$files.2.names`:

```json
{
  "target": "D:\\test-gitmap\\test-gitmap",
  "inputs": [
    "https://github.com/alimtvnetwork/git-repo-navigator",
    "https://github.com/alimtvnetwork/gitmap-v{2..28}"
  ],
  "prMode": "merges",
  "tree": true,
  "finalSync": true,
  "cd": true,
  "imports": [".ai-memory/temp/seo-templates.json"],
  "lineSkippers": [
    { "mode": "starts_with", "pattern": "Co-authored-by:" },
    { "mode": "contains", "pattern": "X-Lovable-Edit-ID" }
  ],
  "titleReplacements": [
    { "matchMode": "equals", "match": "Changes", "replacement": "$files.2.names: $seo.title" }
  ],
  "suffixTemplates": ["seo"]
}
```

## Examples

```bash
# Bootstrap a new config template:
gitmap commit-pull bootstrap

# Launch the visual Web Studio:
gitmap commit-pull ui

# Fast dry-run validation via Array Async Pool:
gitmap commit-pull --config commit-pull-config.json --dry-run

# 1-Line Declarative Replay with Config JSON:
gitmap commit-pull --config .ai-memory/temp/commit-pull-config.json

# Shortened syntax multi-repo range replay:
gitmap commit-pull "D:\target" "git-repo-navigator,gitmap-v{2..28}" --tree --cd
```
