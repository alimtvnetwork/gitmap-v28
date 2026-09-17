# gitmap pipeline fix errors agy (aef)

Extract failing CI/CD pipeline error logs, automatically capture recent git commit logs, embed the 4-part Root Cause Analysis (RCA) prompt, and dispatch the primary fix prompt to Google Antigravity while staging a secondary follow-up verification prompt ("Is it fixed?") into the verification queue.

## Synopsis

```bash
gitmap pipeline fix errors agy [repo] [flags]
gitmap pipeline-fix errors agy [repo] [flags]
gitmap pipeline agy errors fix [repo] [flags]
gitmap pipeline agy-errors-fix [repo] [flags]
gitmap pipeline aef [repo] [flags]
gitmap aef [repo] [flags]
gitmap fix-agy [repo] [flags]
gitmap agy fix-pipeline [repo] [flags]
```

## Aliases & Shorthands

- `gitmap pipeline fix errors agy`
- `gitmap pipeline-fix errors agy`
- `gitmap pipeline agy errors fix`
- `gitmap pipeline agy-errors-fix`
- `gitmap pipeline aef`
- `gitmap aef` (top-level shortcut)
- `gitmap fix-agy`
- `gitmap agy fix-pipeline`
- `gitmap agy fix`
- `gitmap agy fp`

## Options

```text
  -f, --force           Bypass duplicate check and re-send errors previously dispatched
  -r, --repo string     Target repository slug (e.g. alimtvnetwork/gitmap-v28)
  -v, --detailed        Include verbose passing test lines in error logs
      --no-release      Feed CI/CD fix prompt without automated release
  -p, --prompt string   Path to custom prompt markdown template
      --no-clipboard    Skip writing payload to OS system clipboard
      --file string     Optional destination file path for prompt payload
  -d, --dry-run         Preview payload statistics without writing or copying
  -h, --help            Help for pipeline fix errors agy
```

## Duplicate Detection & Re-sending (`--force` / `-f`)

GitMap computes a unique fingerprint for every failing workflow run (`repo:runId:commitSha:errorHash`).
If the same pipeline run and error logs have already been dispatched:
```text
  ⚠ Pipeline errors for run #12345 (commit abc1234) have already been sent to Antigravity!
    Previously sent at: 2026-09-17T15:30:00Z (dispatch count: 1)
    Do you want to send again? Use --force or -f to send again.
```
Pass `--force` or `-f` to bypass deduplication and send again.

## Two-Prompt Workflow Architecture

1. **Prompt 1 (Primary Fix with RCA)**:
   - **Directive Header**: Top-level `/goal` mandating autonomous resolution and strict coding guideline adherence.
   - **Automated Git Log**: Automatically runs `git log -n 5 --stat --no-merges` so the agent understands recent file and commit changes immediately.
   - **Pipeline Error Log**: Full diagnostic error output and step failures.
   - **Embedded RCA Prompt**: Ingests canonical 4-part RCA prompt (`01-prompts/07-bug-fix/01-fix-with-rca.md` or `01-prompts/16-ci-cd/01-ci-cd-fix.md`).
   - **Delivery**: Written to `.lovable/temp/active-agy-pipeline-fix-prompt.txt` and copied to OS clipboard for instant pasting (`Ctrl+V` / `Cmd+V`).

2. **Prompt 2 (Queued Follow-up Verification: "Is it fixed?")**:
   - Secondary verification prompt automatically generated and staged to `.lovable/temp/queued-agy-followup-prompt.txt`.
   - Recorded in prompt queue ledger `.lovable/temp/agy-prompt-queue.json`.
   - When the primary fix loop finishes, this queued prompt directs the agent to verify that all root causes have been fixed and all quality checks pass.

## Examples

### Dispatch pipeline fix prompt with automated git log
```bash
gitmap pipeline fix errors agy
```

### Use ultra-short alias
```bash
gitmap aef
```

### Force resend of previously sent errors
```bash
gitmap pipeline aef -f
gitmap aef --force
```

### Preview assembled payload without modifying clipboard
```bash
gitmap pipeline fix errors agy --dry-run
```
