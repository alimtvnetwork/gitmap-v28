# 266 — `prompt show --copy`: clipboard copy for prompt templates

## 1. Problem statement

`gitmap prompt` manages AI prompt templates (list/show/add/rm/export/import/inject),
but there was no way to copy a template body to the clipboard directly — "see,
use, copy" required manual selection. Agents and users had to re-select long
template bodies by hand.

## 2. Design

- New `--copy` flag on `gitmap prompt show <slug>`: prints the template as
  before, then copies the body to the system clipboard via the existing
  `github.com/atotto/clipboard` dependency (already used by `cmdagy`).
- Headless/no-clipboard environments: warn on stderr
  ("clipboard unavailable ... — template shown above"), exit 0. Showing the
  template is the primary contract; copy is best-effort.
- No changes to the template store schema or other subcommands.

## 3. Implementation

- `cli/cmdprompt/prompt_cmd.go`: register `--copy` bool flag; pass to runner.
- `cli/cmdprompt/prompt_ops.go`: `runPromptShow(slug, copyToClipboard)`;
  clipboard write with graceful warning.

## 4. Acceptance

- `go build ./...` exit 0; gofmt clean.
- `gitmap prompt show <slug> --copy` prints the template and reports
  copied / warns when no clipboard utility exists.
- `gitmap prompt show <slug>` (no flag) behaves exactly as before.
