# Commit In (chronological multi-source replay)

Walk one or more SOURCE git repos in author-date order and APPEND
each commit (preserving BOTH AuthorDate and CommitterDate) into a
TARGET git repo. Useful for stitching together project history that
lives across forks, archives, or versioned siblings into a single
canonical timeline — without ever rewriting an existing commit.

## Alias

cin

## Usage

    gitmap commit-in <source> <input1,input2,...> [flags]
    gitmap cin       <source> all                  [flags]
    gitmap cin       <source> -5                   [flags]

`<source>` is the TARGET repo (the one receiving the appended commits).
The auto-init rule is fixed: URL → `git clone`; existing repo →
reuse; existing non-repo folder → `git init` in place; missing path
→ `mkdir -p && git init`. No prompt, no flag.

`<inputs>` is the comma- (or space-) separated list of sources to walk.
Two keywords are accepted IN PLACE of an explicit list:
- `all`  — every versioned sibling (`<source>-v1`, `<source>-v2`, …)
- `-N`   — only the latest N versioned siblings

Inputs are walked OLDEST → NEWEST by author date, deduped against
`ShaMap` so re-running is idempotent.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--default` / `-d` | off | Load the default profile bound to `<source>` |
| `--profile <name>` | — | Load `.gitmap/commit-in/profiles/<name>.json` |
| `--save-profile <name>` | — | Persist this run's resolved settings as a profile |
| `--save-profile-overwrite` | off | Allow `--save-profile` to overwrite |
| `--set-default` | off | Mark the saved profile as default for `<source>` |
| `--author-name <s>` | — | Override author name (requires `--author-email`) |
| `--author-email <s>` | — | Override author email (requires `--author-name`) |
| `--conflict <mode>` | `ForceMerge` | `ForceMerge` or `Prompt` |
| `--exclude <csv>` | — | Per-commit exclude list (trailing `/` = folder) |
| `--message-exclude <csv>` | — | `Kind:Value` rules: `StartsWith:`/`EndsWith:`/`Contains:` |
| `--message-prefix <csv>` | — | Random-pick pool prepended to every body |
| `--message-suffix <csv>` | — | Random-pick pool appended to every body |
| `--title-prefix <s>` | — | Prepended to the FIRST line only |
| `--title-suffix <s>` | — | Appended to the FIRST line only |
| `--override-messages <csv>` | — | Replaces the entire message (random pick) |
| `--override-only-weak` | off | Override only when title's first word is weak |
| `--weak-words <csv>` | `change,update,updates` | First-word triggers for override |
| `--function-intel <on\|off>` | `off` | Append per-language new-function block |
| `--languages <csv>` | `Go` | Languages scanned when intel is on |
| `--no-prompt` | off | Refuse interactive prompts; exit MissingAnswer if unset |
| `--dry-run` | off | Plan only; never run `git commit` |
| `--keep-temp` | off | Keep `.gitmap/temp/<runId>/` after exit |
| `--no-release-branch` | off | Suppress auto `release/<tag>` branch on version tags |

## Examples

### 1. Convert a plain folder of files into a git repo + replay history

Folder exists, no `.git/` yet — `commit-in` runs `git init` for you,
then appends history from the URL on top:

    gitmap commit-in ./my-project https://github.com/me/my-project-archive.git

### 2. Mix a local folder + a remote URL as INPUTS

Target = `./canonical`. Inputs = local checkout + 2 remote forks,
walked oldest → newest into one canonical timeline:

    gitmap cin ./canonical \
        ./old-local-checkout,https://github.com/me/old-fork.git,git@github.com:me/new-fork.git

### 3. Brand-new target folder from scratch

Path does not exist → `mkdir -p` + `git init` + replay, in one shot:

    gitmap commit-in ./brand-new-canonical \
        https://github.com/me/legacy-v1.git,https://github.com/me/legacy-v2.git

### 4. Replay every versioned sibling

`all` = every `<source>-vN`; `-N` = latest N only.

    gitmap commit-in ./canonical all --save-profile Default --set-default
    gitmap cin       ./canonical -3 --dry-run --function-intel on --languages Go,TypeScript

### 5. Override author + scrub commit messages

    gitmap cin git@github.com:me/canonical.git \
        https://github.com/me/old-fork.git,https://github.com/me/new-fork.git \
        --author-name "Jane Doe" --author-email jane@example.com \
        --message-exclude "StartsWith:Signed-off-by:,Contains:[skip ci]" \
        --title-suffix " — via gitmap"

### 6. Reuse a saved profile + only rewrite weak titles

    gitmap cin ./canonical all --default \
        --override-messages "Refine implementation,Improve module" \
        --override-only-weak

### 7. Headless CI run

    gitmap cin ./canonical all --profile CI --no-prompt

## Sample profile JSON

Drop this at `.gitmap/commit-in/profiles/Default.json` (workspace
root = nearest ancestor containing `.gitmap/`). Load with
`--profile Default` or `--default`. Strict PascalCase decoding —
unknown keys are an error. Every field maps 1:1 to a CLI flag above.

    {
      "Name": "Default",
      "SchemaVersion": 1,
      "SourceRepoPath": "/abs/path/to/canonical",
      "IsDefault": true,
      "ConflictMode": "ForceMerge",
      "Author": { "Name": "Jane Doe", "Email": "jane@example.com" },
      "Exclusions": [
        { "Kind": "PathFolder", "Value": "node_modules" },
        { "Kind": "PathFile",   "Value": "secrets.env" }
      ],
      "MessageRules": [
        { "Kind": "StartsWith", "Value": "Signed-off-by:" },
        { "Kind": "Contains",   "Value": "[skip ci]" }
      ],
      "MessagePrefix":   ["chore:", "feat:", "fix:"],
      "MessageSuffix":   [],
      "TitlePrefix":     "",
      "TitleSuffix":     " — via gitmap",
      "OverrideMessages": ["Improve module", "Refine implementation"],
      "OverrideOnlyWeak": true,
      "WeakWords":        ["change", "update", "updates"],
      "FunctionIntel": { "IsEnabled": true, "Languages": ["Go", "TypeScript"] }
    }

Tip: let gitmap write it for you first —
`gitmap cin ./canonical all --save-profile Default --set-default` —
then edit the resulting JSON. Re-saving requires
`--save-profile-overwrite`. Profiles bind by absolute filesystem
path, NOT by remote URL.

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Ok — every walked commit was Created or Skipped |
| `1`  | PartiallyFailed — at least one commit failed but others succeeded |
| `2`  | BadArgs — flag / positional validation failed |
| `3`  | SourceUnusable — `<source>` could not be resolved or initialized |
| `4`  | InputUnusable — at least one input could not be cloned / opened |
| `5`  | DbFailed — SQLite migration or write failed |
| `6`  | ProfileMissing — `--profile` / `--default` lookup empty |
| `7`  | MissingAnswer — `--no-prompt` set but a required value was unset |
| `8`  | ConflictAborted — `Prompt` mode and the user aborted the merge |
| `9`  | LockBusy — another `commit-in` run holds the workspace lock |
| `10` | FunctionIntel — a per-language detector panicked |

## See also

`gitmap clone-pick` (sparse-checkout single-repo clone),
`gitmap clone-from` (one-shot mirror from a manifest),
spec/03-commit-in/ for the full normative contract.
## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter commit-in
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
