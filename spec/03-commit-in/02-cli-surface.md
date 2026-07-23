# 02 — CLI Surface

## 2.1 Aliases (Enum: `CommitInAlias`)

| Member          | Token        |
|-----------------|--------------|
| `CommitInLong`  | `commit-in`  |
| `CommitInShort` | `cin`        |

Both resolve to the same handler. CLI ID lives in `constants_cli.go`
(per Core memory rule: all CLI IDs in that single file).

## 2.2 Argv grammar (formal)

```
gitmap-v28 (commit-in | cin) <source> <inputs...> [flags]

<source>     := PATH | URL
<inputs...>  := INPUT (SEP INPUT)*  |  KEYWORD
INPUT        := PATH | URL | QUOTED
SEP          := SPACE+ | "," | "," SPACE+
KEYWORD      := "all" | "-" DIGIT+
QUOTED       := '"' (PATH | URL) '"'
```

- `<inputs...>` MAY be a single argv token containing internal commas
  (e.g. `"repo-a, repo-b, repo-c"`) — split by the parser using `SEP`.
- `<inputs...>` MAY also be many argv tokens (`repo-a repo-b repo-c`).
- Mixed forms (`"repo-a, repo-b" repo-c`) are valid.
- Quoting is optional and applies per item OR around the whole list.
- A `KEYWORD` MUST appear alone — never mixed with explicit inputs.

## 2.3 `<source>` resolution rule (no prompt, no flag)

Order of checks (first match wins):

1. `<source>` matches `^(https?://|git@|ssh://|git://)` → `git clone`
   into `CWD/<basename-without-.git>`. The clone target becomes the
   resolved source path.
2. `<source>` is an existing directory containing `.git/` (or is a
   bare repo) → reuse as-is.
3. `<source>` is an existing directory WITHOUT `.git/` → `git init`
   inside it. Existing files remain on disk, untracked.
4. `<source>` does NOT exist → `mkdir -p <source>` then `git init`
   inside it.

No interactive prompt. No `--init` flag. The behavior is deterministic
from the path alone. (Maps to INV-05.)

## 2.4 Special input keywords (Enum: `InputKeyword`)

| Member             | Token   | Effect                                                                                  |
|--------------------|---------|-----------------------------------------------------------------------------------------|
| `InputKeywordAll`  | `all`   | Discover sibling versioned repos of `<source>`'s basename. Walk in ascending version.   |
| `InputKeywordTail` | `-N`    | Same as `all` but truncated to the last `N` versioned siblings.                         |

Discovery rules:

- Search the **parent directory of `<source>`**.
- A sibling matches when its basename equals `<base>` or `<base>-vN`
  where `<base>` is `<source>`'s basename with any trailing `-vN`
  stripped (e.g. `gitmap-v28` → base `gitmap-v28`).
- Walk order: plain `<base>` FIRST (treated as `v0`), then `v1`, `v2`,
  …, `vK`. Numeric ascending. (Resolves Ambiguity #1, #2.)
- `<source>` itself is excluded from the input list.

## 2.5 Flags (full table)

| Long                     | Short | Type                          | Default        | Notes                                                                       |
|--------------------------|-------|-------------------------------|----------------|-----------------------------------------------------------------------------|
| `--default`              | `-d`  | bool                          | `false`        | Load profile where `IsDefault = 1` and `SourceRepoPath = <source>`.         |
| `--profile <name>`       |       | string                        | —              | Load `.gitmap/commit-in/profiles/<name>.json`.                              |
| `--save-profile <name>`  |       | string                        | —              | Persist current run's resolved answers as a profile.                        |
| `--save-profile-overwrite`|      | bool                          | `false`        | Allow `--save-profile` to overwrite an existing same-named profile.         |
| `--set-default`          |       | bool                          | `false`        | When combined with `--save-profile`, mark the saved profile `IsDefault=1`.  |
| `--author-name <s>`      |       | string                        | —              | Override author identity (Name).                                            |
| `--author-email <s>`     |       | string                        | —              | Override author identity (Email). Both required together or neither.        |
| `--conflict <mode>`      |       | enum `ConflictMode`           | `ForceMerge`   | `ForceMerge` \| `Prompt`.                                                   |
| `--exclude <csv>`        |       | csv of relative paths         | empty          | Each item classified as `PathFolder` or `PathFile` by trailing `/`.         |
| `--message-exclude <csv>`|       | csv of `Kind:Value` rules     | empty          | `StartsWith:`, `EndsWith:`, `Contains:` prefixes.                           |
| `--message-prefix <csv>` |       | csv                           | empty          | Single value or pool (random pick per commit).                              |
| `--message-suffix <csv>` |       | csv                           | empty          | Same.                                                                       |
| `--title-prefix <s>`     |       | string                        | empty          | Applied to first line only.                                                 |
| `--title-suffix <s>`     |       | string                        | empty          | Applied to first line only.                                                 |
| `--override-messages <csv>`|     | csv pool                      | empty          | Random pick per commit when override fires.                                 |
| `--override-only-weak`   |       | bool                          | `false`        | Override only when title's first word ∈ `WeakWords`.                        |
| `--weak-words <csv>`     |       | csv                           | `change,update,updates` | Override the default weak-word list.                               |
| `--function-intel <on\|off>`|    | enum                          | `off`          | Toggle per-language new-function detection block.                           |
| `--languages <csv>`      |       | csv of `FunctionIntelLanguage`| `Go`           | Subset of supported languages to scan when intel is on.                     |
| `--no-prompt`            |       | bool                          | `false`        | Refuse to interactively prompt; missing values become exit `MissingAnswer`. |
| `--dry-run`              |       | bool                          | `false`        | Walk + plan + log; never `git commit`. Still writes DB rows with outcome `Skipped`.|
| `--keep-temp`            |       | bool                          | `false`        | Don't delete `<.gitmap>/temp/<runId>/` on exit (debugging).                 |

## 2.6 Interactive prompts (only fired when value unset and `--no-prompt` not given)

Order is fixed; the prompt runner SHALL ask in this exact sequence:

1. `ConflictMode`             (default `ForceMerge`)
2. `Exclusions`               (CSV)
3. `MessageRules`             (CSV of `Kind:Value`)
4. `Author override`          (Name + Email, or skip)
5. `MessagePrefix` / `MessageSuffix` / `TitlePrefix` / `TitleSuffix`
6. `OverrideMessages` + `OverrideOnlyWeak` + `WeakWords`
7. `FunctionIntel` toggle + `Languages`
8. Save as profile? (`No` | `Yes` | `Yes + set as default`)

Each answered value populates the in-memory `RunConfig`. After the
prompts, the `RunConfig` is identical-shape to the JSON profile (see §05).

## 2.7 Exit codes (Enum: `CommitInExit`)

| Member                          | Code | Trigger                                                  |
|---------------------------------|------|----------------------------------------------------------|
| `CommitInExitOk`                | 0    | All inputs processed; `Run.Status = Completed`.          |
| `CommitInExitPartiallyFailed`   | 1    | At least one commit failed; `Run.Status = PartiallyFailed`.|
| `CommitInExitBadArgs`           | 2    | Argv grammar violation, unknown flag, mixed `KEYWORD`.   |
| `CommitInExitSourceUnusable`    | 3    | `<source>` could not be resolved (clone fail, perm).     |
| `CommitInExitInputUnusable`     | 4    | At least one `<input>` could not be cloned/opened.       |
| `CommitInExitDbFailed`          | 5    | SQLite migration or write failure.                       |
| `CommitInExitProfileMissing`    | 6    | `--default` / `--profile` requested but not found.       |
| `CommitInExitMissingAnswer`     | 7    | `--no-prompt` set and a required answer is unset.        |
| `CommitInExitConflictAborted`   | 8    | `Prompt` mode: user aborted the diff resolution.         |
| `CommitInExitLockBusy`          | 9    | Another `commit-in` is running in this workspace.        |
| `CommitInExitFunctionIntel`     | 10   | A language parser failed AND `--function-intel on` set.  |

## 2.8 STDOUT / STDERR contract

- `STDOUT` — only the final summary line: `commit-in: run=<runId> created=<n> skipped=<n> failed=<n>`.
- `STDERR` — every `▸` phase banner, every `✓ / ✗` per-commit line,
  every WARN/ERROR. Per zero-swallow rule, every error path writes
  exactly one `commit-in: <stage>: <message>` line to `STDERR` before
  exiting.