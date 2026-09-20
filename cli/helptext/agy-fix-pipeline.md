# gitmap pipeline fix errors agy (aef)

Extract failing CI/CD pipeline error logs, automatically capture recent git commit logs, embed the complete error logs and 4-part Root Cause Analysis (RCA) prompt into `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`, inject the fix prompt directly into Google Antigravity (`agy` CLI / IDE), and stage a secondary follow-up verification prompt ("Is it fixed?") into the verification queue. Supports multi-project parallel batching across repositories with persistent batch cursor tracking.

## Synopsis

```bash
gitmap pipeline fix errors agy [repo] [flags]
gitmap pipeline errors agy fix [repo] [flags]
gitmap pipeline aef [repo] [flags]
gitmap aef [repo] [flags]
gitmap agy fix-pipeline [repo] [flags]
gitmap agy fix [repo] [flags]
```

## Options

```text
  -f, --force           Bypass duplicate check and re-send errors previously dispatched
  -r, --repo string     Target repository slug (e.g. alimtvnetwork/gitmap-v28)
      --all             Scan all tracked repositories for failing pipelines
      --projects int    Number of projects to batch-fix (default: 3)
      --limit int       Limit number of projects to batch-fix (default: 3)
      --reset-batch     Reset multi-project batch cursor to start from the first repository
      --no-inject       Skip direct Antigravity CLI/IDE injection (clipboard and disk only)
  -v, --detailed        Include verbose passing test lines in error logs
      --no-release      Feed CI/CD fix prompt without automated release
  -p, --prompt string   Path to custom prompt markdown template
      --no-clipboard    Skip writing payload to OS system clipboard
      --file string     Optional destination file path for prompt payload
  -d, --dry-run         Preview payload statistics without writing or copying
  -h, --help            Help for pipeline fix errors agy
```

## Direct Antigravity Injection & Dynamic Queue Protocol

1. **IDE-First Injection & Queue Protocol**:
   - **Priority Inversion**: Automatically detects active Antigravity IDE sessions before CLI runners.
   - **Execution State Detection**: Inspects conversation execution status from `transcript.jsonl` / `transcript_full.jsonl`.
   - **Dynamic Queue Delivery**:
     - When IDE conversation is **RUNNING**: The prompt is automatically appended to `.ai-memory/temp/agy-prompt-queue.json` (`queued`), avoiding disruption of the active session.
     - When IDE conversation is **IDLE**: The prompt is staged directly to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` and copied to the OS system clipboard (`ide`).
   - **CLI Suppression**: Background `agy.exe` runners are never spawned when an IDE process is already active.

2. **Full File Path Display**:
   - In all terminal outputs, GitMap displays the **complete absolute path** on disk for every generated file:
     - `• Saved Payload: <FullAbsolutePath>/.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`
     - `• Follow-up Verification: <FullAbsolutePath>/.ai-memory/temp/queued-agy-followup-prompt.txt`
     - `• Queue Ledger: <FullAbsolutePath>/.ai-memory/temp/agy-prompt-queue.json`

3. **Complete Error Logs Embedding**:
   - `active-agy-pipeline-fix-prompt.txt` embeds the full failing CI/CD logs, commit history, and the 4-part RCA fix directives.
   - `queued-agy-followup-prompt.txt` contains the `# Verification Check: Is It Fixed?` follow-up verification check.

## Multi-Project Parallel Batching

When run on a root folder or with `--all` / `--projects <N>`:
1. GitMap concurrently scans all tracked repositories for failing pipeline runs.
2. Batches execution with a default limit of **3 projects** per run.
3. Automatically saves the batch progress in `.ai-memory/temp/pipeline-fix-batch-cursor.json`.
4. Running the command again immediately processes the next batch (e.g. next 2 or 3 projects) without duplicate dispatch.
5. Use `--reset-batch` to reset the cursor to the beginning.

## Examples

### Fix pipeline errors for current repository with direct AGY injection
```bash
gitmap pipeline errors agy fix
```

### Batch fix failing pipelines across all projects (first 3 projects)
```bash
gitmap pipeline errors agy fix --all
```

### Process next batch of failing projects
```bash
gitmap pipeline errors agy fix --all
```

### Specify custom batch limit of 5 projects
```bash
gitmap pipeline errors agy fix --all --limit 5
```

### Force resend of previously dispatched errors
```bash
gitmap pipeline errors agy fix --force
```

