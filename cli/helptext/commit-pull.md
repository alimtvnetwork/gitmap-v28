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

## Declarative Config JSON (`--config`)

Feed a single config JSON file to auto-create the target repo, import state templates with SHA-256 hash deduplication, pre-compile variables before loop execution, strip unwanted lines (`starts_with`, `ends_with`, `contains`, `regex`), and replace generic `"Changes"` commit titles using `$files.2.names`:

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
# 1-Line Declarative Replay with Config JSON:
gitmap commit-pull --config .ai-memory/temp/commit-pull-config.json

# Replay commits into target repository with tree display:
gitmap commit-pull D:\target repo1 repo2 --tree --cd
```
