# gitmap agy fix-pipeline

Extract failing CI/CD pipeline error logs and combine them with the repository's CI/CD fix prompt into the OS clipboard and active prompt file for direct feeding into the Google Antigravity IDE.

## Synopsis

```bash
gitmap agy fix-pipeline [repo] [flags]
gitmap agy fix [repo] [flags]
gitmap agy fp [repo] [flags]
```

## Aliases

- `gitmap agy fix-pipeline`
- `gitmap agy fix`
- `gitmap agy fp`
- `gitmap agy pipeline-fix`
- `gitmap agy fixpipeline`

## Options

```text
  -r, --repo string     Target repository slug (e.g. alimtvnetwork/gitmap-v28)
  -v, --detailed        Include verbose passing test lines in error logs
      --no-release      Feed CI/CD fix prompt without automated release (01-ci-cd-fix.md)
  -p, --prompt string   Path to custom prompt markdown template
      --no-clipboard    Skip writing payload to OS system clipboard
      --file string     Optional destination file path for prompt payload
  -d, --dry-run         Preview payload statistics without writing or copying
  -h, --help            Help for fix-pipeline
```

## How It Works

1. **Error Extraction**: Automatically queries the latest failing GitHub Actions workflow runs or cached errors from the CLI pipeline SQLite database (`data/pipeline/`).
2. **Two-Line Gap**: Appends two newlines (`\n\n`) as a distinct structural separator between error diagnostics and AI prompt instructions.
3. **Prompt Loading**: Ingests the canonical CI/CD fix prompt from `01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md` (or `01-ci-cd-fix.md` if `--no-release` is passed).
4. **Clipboard & Disk Persistence**: Copies the assembled payload to the system clipboard via `clipboard.WriteAll` and saves a copy to `.lovable/temp/active-agy-pipeline-fix-prompt.txt`.
5. **Direct Paste**: Open Antigravity IDE and press `Ctrl+V` (or `Cmd+V`) in the chat prompt to start the autonomous fix loop.

## Examples

### Feed failing pipeline logs with release prompt
```bash
gitmap agy fix-pipeline
```

### Use shorthand alias
```bash
gitmap agy fix
gitmap agy fp
```

### Target a specific sibling repository
```bash
gitmap agy fix alimtvnetwork/gitmap-v28
```

### Preview assembled payload without modifying clipboard
```bash
gitmap agy fix-pipeline --dry-run
```
