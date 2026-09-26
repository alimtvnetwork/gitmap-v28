# gitmap migrate

Visualizes repository remapping, history consolidation, stacked PR pipeline,
and SemVer release ceremony directly in the terminal.

## Usage

    gitmap migrate
    gitmap migrate graph
    gitmap migrate preflight
    gitmap migrate plan

## Output

Renders the preflight execution graph merging `git-repo-navigator` and
`gitmap-v2` into `gitmap-v28`, showing approval gates and release tagging.

## Examples

```bash
# Render migration preflight graph in terminal
gitmap migrate graph

# Inspect preflight status
gitmap migrate preflight
```

