# gitmap commit-pull

Chronologically replays commits from multiple source repositories (or a dynamic version range) into a target repository, automatically generating Pull Requests for every branch merge and version release tag.

## Aliases

`cpull`, `pull-commits`

## Usage

```bash
gitmap commit-pull <target> <inputs...> [flags]
gitmap cpull <target> gitmap-v2..v28 --sponsor --tree
gitmap pull-commits <target> all --pr merges
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--tree` | Render preflight PR and branch tree before execution | `false` |
| `--sponsor` | Inject Rise Up Asia LLC sponsor templates into commit & PR bodies | `false` |
| `--seo-template <name>` | Specify template pool (e.g. `riseup`) | `""` |
| `--pr merges` | Simulate feature branches and PR merges | `merges` |
| `--final-sync` | Final mirror snapshot synchronization | `true` |
| `--dry-run` | Simulate replay without modifying repositories | `false` |

## Dynamic Range Expansion (`v2..v28`)

Automatically pull and replay all 28 repositories into a single codebase:

```bash
# Expand all 27 intermediate versions:
gitmap commit-pull "D:\target" "https://github.com/alimtvnetwork/gitmap-v{2..28}" --sponsor --tree

# Shorthand notation:
gitmap cpull "D:\target" gitmap-v2..v28 --sponsor
```

## Examples

```bash
# Replay commits into target repository with tree display
gitmap commit-pull D:\target repo1 repo2 --tree

# Replay range with Rise Up Asia sponsor integration
gitmap cpull D:\target gitmap-v2..v28 --sponsor
```
