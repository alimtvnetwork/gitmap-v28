# gitmap setup

Interactive first-time configuration wizard that applies global Git settings and installs shell tab-completion.

## Alias

None

## Usage

    gitmap setup [--config <path>] [--dry-run]
    gitmap setup print-path-snippet --shell <bash|zsh|fish|pwsh> --dir <path> [--manager <label>]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --config \<path\> | data/git-setup.json beside gitmap | Path to git-setup.json config file |
| --dry-run | false | Preview changes without applying them |

## Subcommand: `print-path-snippet`

Emits the canonical marker-block PATH snippet to stdout. `run.sh` and
`gitmap/scripts/install.sh` shell out to this command so all three
drivers produce **byte-identical** rc-file output. Single source of
truth lives in `gitmap/constants/constants_pathsnippet.go`.

| Flag | Default | Description |
|------|---------|-------------|
| --shell \<sh\> | bash | Target shell: bash, zsh, fish, or pwsh |
| --dir \<path\> | _(required)_ | Directory to inject into the snippet's PATH line |
| --manager \<label\> | gitmap setup | Header label; different managers can coexist in one rc-file |

Example:

    gitmap setup print-path-snippet --shell zsh --dir ~/.local/bin/gitmap >> ~/.zshrc

## Prerequisites

- Git must be installed

## Examples

### Example 1: Run the setup wizard

    gitmap setup

**Output:**

    ■ Applying global Git configuration...
      ✓ core.autocrlf = true
      ✓ push.default = current
      ✓ pull.rebase = false
    ✓ 3 Git settings applied

    ■ Shell Completion
      Detected shell: powershell
      Installing completion to $PROFILE...
    ✓ Shell completion installed for PowerShell

    ■ Setup complete!
    → Run 'gitmap scan <directory>' to start tracking repos

### Example 2: Dry-run mode (preview only)

    gitmap setup --dry-run

**Output:**

    [DRY RUN] No changes will be made
    [DRY RUN] Would set core.autocrlf = true
    [DRY RUN] Would set push.default = current
    [DRY RUN] Would set pull.rebase = false
    [DRY RUN] Would install powershell completion to $PROFILE
    No changes made.

### Example 3: Setup with custom config file

    gitmap setup --config ./my-config/git-setup.json

**Output:**

    ■ Loading config from ./my-config/git-setup.json...
    ■ Applying global Git configuration...
      ✓ core.autocrlf = true
      ✓ init.defaultBranch = main
    ✓ 2 Git settings applied
    ✓ Setup complete!

## See Also

- [completion](completion.md) — Generate completion scripts manually
- [scan](scan.md) — Scan directories after setup
- [doctor](doctor.md) — Diagnose installation issues

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter setup
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
