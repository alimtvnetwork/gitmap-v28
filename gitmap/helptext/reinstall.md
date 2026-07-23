# gitmap reinstall

Reinstall the gitmap binary itself. Auto-detects the right method based
on whether a source repo is linked.

## Usage

    gitmap reinstall [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --mode {auto\|repo\|self} | auto | Force a specific reinstall path |
| --yes / -y | false | Auto-confirm without prompting |

## Modes

**auto** (default) — picks the right path:

- If a source repo is linked (RepoPath set at build time), runs
  `run.ps1 -reinstall` (Windows) or `run.sh --reinstall` (Unix). This
  re-pulls, re-builds, re-deploys, and re-runs `gitmap setup` for you.
- Otherwise, runs `gitmap self-uninstall --confirm` then
  `gitmap self-install --yes`. This downloads the install script and
  reinstalls the binary into the same install dir.

**repo** — force the run-script path. Fails if no source repo is linked.

**self** — force the download path. Useful when you want to ignore the
linked repo and reinstall from the latest published release.

## Prerequisites

- For `--mode=repo`: a linked source repo (set automatically by the
  initial `run.ps1` deploy, or via `gitmap set-source-repo`).
- For `--mode=self`: network access to fetch `install.ps1` / `install.sh`.

## Examples

### Default (auto-detect)

    gitmap reinstall

Output (with linked repo):

      ╔══════════════════════════════════════╗
      ║         gitmap reinstall             ║
      ╚══════════════════════════════════════╝

      → Mode: auto (detected: repo)
      → Source repo: D:\wp-work\riseup-asia\gitmap
      → Proceed with reinstall? (yes/N): yes
      → Running run.ps1 -reinstall ...
    ...
      ✓ Reinstall complete.

### Force download path

    gitmap reinstall --mode=self --yes

Output:

      → Mode: self (forced via --mode)
      [1/2] self-uninstall
      ...
      [2/2] self-install
      ...
      ✓ Reinstall complete.

### Skip prompt

    gitmap reinstall -y

## See Also

- `gitmap self-install` — install the binary from scratch
- `gitmap self-uninstall` — remove the binary
- `gitmap install <tool>` — install third-party tools (vscode, npp, …)
- `gitmap update` — update without reinstalling

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter reinstall
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
