# gitmap scan

> **Generated — do not hand-edit.** This topic is rendered from the
> `cmdscan.ScanHelpDisplay` HelpDisplay struct. Regenerate with:
>
> `gitmap py 03-ai-scripts/51-helptext-generator.py --topics scan`

```text
Repository Discovery Scanner (gitmap scan)

  Usage:
    gitmap scan [dir] [flags]
    gitmap s [dir] [flags]
    gitmap scan --fix
    gitmap scan ~ --force-include .oh-my-zsh
    gitmap scan ~ -fi .oh-my-zsh,node_modules
    gitmap scan /home/user --force-include all
    gitmap scan export [--machine <name>] [--out <dir>]
    gitmap scan merge <dir>... [--out <file>]

Output & Formatting:
  --output <mode>  Output format: terminal (default), csv, or json
  --output-path <dir>  Directory to write scan artifacts (.gitmap/output)
  --manifest <dir>  Alias for --output-path (aligned with reclone)
  --relative-root <dir>  Pin base directory for byte-stable relative paths
Portable Repo Sets:
  export [--machine <name>] [--out <dir>]  Export cached repos to <out>/<machine-slug>/repos.json (clone with: gitmap clone-from <file> --execute)
  merge <dir>... [--out <file>]  Merge export folders into one deduped JSON (dedupe by URL)
Scanner Walk & Performance:
  --workers <n>  Parallel directory walker pool size (1-16, auto: NumCPU)
  --max-depth <n>  Max folder depth to descend (default: 4, -1: unlimited)
  --default-branch <b>  Fallback branch name when HEAD detection returns empty
  --force-include <dir>, -fi  Force scan of default excluded directories (.oh-my-zsh, or 'all')
  --config <path>  Path to configuration file (default: ./data/config.json)
Integrations & Probing:
  --fix  Run database reconciliation (removes stale rows, adds missing)
  --no-vscode-sync  Skip syncing discovered repos into VS Code Project Manager
  --no-auto-tags  Skip auto-derived project tags (git, node, go, etc.)
  --github-desktop  Auto-register discovered repositories with GitHub Desktop
  --no-probe  Skip background remote branch version probe entirely
Flags:
  --quiet  Suppress interactive walker spinner and clone hints
  --open  Automatically open output directory after scan completes
  -v, --verbose  Emit verbose directory traversal and probe logs
  -h, --help  Show this scanner help menu

  Tip: Run 'gitmap scan . --output json' to generate a machine-readable repo catalog.
  Tip: Use 'gitmap scan --fix' to reconcile local disk with GitMap database.
  Tip: Pass '--max-depth 2' to quickly scan shallow project directories.
  Tip: Pass '--force-include .oh-my-zsh' (alias: -fi) to force scanning excluded directories or 'all' for everything.
```
