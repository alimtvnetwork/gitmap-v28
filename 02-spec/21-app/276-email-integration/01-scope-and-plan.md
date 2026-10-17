# Spec 276 — Email integration: scope and plan

**Version:** 1.0.0
**Updated:** 2026-10-17
**Spec:** 276 (auto-discovered next free number; highest existing is 275)
**Parent program:** email-notification-and-gitmap-integration, Part C
(go-email-reader Part B spec lives at `~/workspace/go-email-reader/02-spec/21-app/email-notification-and-gitmap-integration/README.md`)

---

## 1. Scope

IN SCOPE:

1. `gitmap space send-email --to X [--subject S] [--body B] [--file F] [--html]`
   - Resolve the go-email-reader binary (PATH → `~/go/bin` → gitmap setting
     `email.go_email_reader_bin`); when found, delegate by exec'ing
     `go-email-reader send --to … --subject … --body … [--file …] [--html]`.
     Never pass passwords on the command line (go-email-reader reads its own
     RSA-encrypted store).
   - When NOT found, use gitmap's NEW built-in SMTP sender: SMTP
     host/port/TLS-mode from gitmap settings, authenticated
     (STARTTLS/587 or implicit TLS/465). Password is ALWAYS prompted at
     runtime via a masked prompt (`golang.org/x/term` ReadPassword, fail
     closed on non-TTY) or read from the existing `cli/secrets` RSA-OAEP
     vault — never stored in plaintext settings.
   - `--file` accepts YAML/JSON/MD/plain-text with subject/body rules
     IDENTICAL to go-email-reader's `send --file`.
   - `--html` sends the body as `text/html`.
2. `gitmap space schedule-send-email --to X --task <enum>`
   - Task-type enum as a dedicated Go string type with validation:
     `release | pull | clone | scan | push | backup` (extendable).
   - Registers a trigger in gitmap settings; when gitmap executes that
     operation it sends a nicely formatted HTML email to the address with
     machine info: hostname, OS user, IP address, repo, branch, operation,
     outcome/summary.
3. A delegate-vs-builtin decision matrix.
4. Coding constraints per `.ai-memory/coding-guidelines.md`.

OUT OF SCOPE:

- Changing the existing AGM account-switch notify path
  (`cli/cmdagy/agy_account_switch_notify.go`, nil-auth `dispatchSmtpEmail`,
  port 25 default, env `SMTP_HOST`/`SMTP_FROM`/`SMTP_PORT`) — that path stays
  as-is; new code must NOT copy the nil-auth pattern.
- Changing the existing `gitmap email setup` (the plaintext
  `email.password` anti-pattern is called out in §3.3 but migrating it is a
  separate decision, not this spec).
- Attachments (go-email-reader supports `--attach`; the gitmap flag set here
  deliberately omits it — v2 feature).
- Time-based scheduling (`--when` cron-style): v1 triggers fire when gitmap
  executes the registered operation, not on a clock. `cron`/`hooks` namespaces
  are the venue for wall-clock scheduling; not designed here.
- No Go code, no commits — this spec is spec files only.

## 2. Existing-code context (verified by research — cite these paths)

| Path | What it gives this spec |
|---|---|
| `cli/cmdspace/space.go` | `DispatchSpace(command)` routes `gitmap space <subcommand>` off `os.Args[2:]`; new subcommands plug in here alongside `backup-branch` and `PR-merge` |
| `cli/cmd/cmdspace` registration | `cli/cmd/root.go:650` — `found, err = cmdspace.DispatchSpace(command)`; success exits via `handleDispatchResult` (`cli/cmd/root.go:544`) |
| `cli/constants/constants_commands_extra.go` | `CmdSpace`, `CmdSpaceBackupBranch`, `CmdSpacePRMerge` pattern — new constants `CmdSpaceSendEmail`, `CmdSpaceScheduleSendEmail` follow it |
| `cli/cmdagy/agy_account_switch_notify.go` | `dispatchSmtpEmail` — stdlib `net/smtp`, `smtp.SendMail(host+":"+port, nil, …)` with nil auth and port 25 default. THE ANTI-PATTERN: new code must ADD auth (STARTTLS/587 or implicit TLS/465) and must not copy the nil-auth call |
| `cli/cmdagy/agy_telegram_email_settings.go:256-260` | `saveEmailSpeedConfig` → `setGlobalIfNonEmpty("email.password", "SMTP_PASSWORD", pass)` → `config.SetVariable("global", "email.password", …)` — the plaintext-password anti-pattern to AVOID. New code must never write a password into settings |
| `cli/cmdagy/agy_telegram_email_settings.go:454` | `readSpeedVar("email.smtp_host", defaultIfEmpty(os.Getenv("SMTP_HOST"), "smtp.gmail.com:587"))` — existing settings/env fallback pattern new code may reuse for non-secret fields |
| `cli/secrets/ssh_password.go` | `DecryptStoredPassword(cipherText)` — RSA-OAEP vault reader; approved pattern when a stored secret must be used. Accepts `rsa:` / `aes:` / `salt:` / `caesar:` prefixed ciphertext, falls back to the original string |
| `cli/cmdlogin/login_browser.go:70` | `term.ReadPassword(int(os.Stdin.Fd()))` — masked prompt precedent |
| `cli/cmdos/os_autologin_prompt.go:53-65` | `promptPassword(label)` — masked prompt precedent (label + ReadPassword + newline, empty string on error) |
| `cli/cmdssh/ssh_client.go:215`, `cli/cmdssh/ssh_login_pass_prompt.go:16` | `readPasswordHook = term.ReadPassword` — hookable masked-prompt precedent (unit tests can inject a fake) |
| `cli/config/variables.go:117` | `config.GetVariable(scope, key)` — scope → global → OS env lookup; settings storage for all new `email.*` keys |
| `cli/cmdip/ip_resolver.go:26` | `GetLocalIP(ctx, skipLoopback, ifaceName)` — machine IP for trigger emails |
| `~/workspace/go-email-reader/cmd_mail.go:105-150` | `go-email-reader send` flag surface: `--to`, `--subject`, `--body`, `--body-file`, `--file`, `--html`, `--template`, `--var`, `--attach` — delegation targets only the subset in §3 |
| `~/workspace/go-email-reader/internal/mail/parse_email_file.go` | `ParseEmailFile` + `EmailFile{To, CC, BCC, Subject, Body, HTML}` — EXACT parsing semantics §3.5 must replicate |

## 3. Command 1 — `gitmap space send-email`

### 3.1 Invocation and flags

```
gitmap space send-email --to X [--subject S] [--body B] [--file F] [--html]
```

| Flag | Meaning |
|---|---|
| `--to X` (required) | Recipient address(es), comma-separated |
| `--subject S` | Subject; required unless `--file` supplies one |
| `--body B` | Plain-text body; required unless `--file` supplies one |
| `--file F` | Build the email from a YAML/JSON/Markdown/text file (§3.5); mutually exclusive with `--subject`/`--body` — refuse with a non-zero exit when combined |
| `--html` | Send the body as `text/html`; a file's `html: true` key has the same effect |

Routing: add `CmdSpaceSendEmail = "send-email"` to
`cli/constants/constants_commands_extra.go` (beside `CmdSpaceBackupBranch`),
and a branch in `DispatchSpace` in `cli/cmdspace/space.go` following the
`backup-branch`/`PR-merge` pattern (`--help`/`-h` prints usage, else run).

### 3.2 Binary resolution order (highest priority first)

1. **PATH**: `exec.LookPath("go-email-reader")`.
2. **`~/go/bin`**: `filepath.Join(os.UserHomeDir(), "go", "bin", "go-email-reader")`
   — the manual Go toolchain location per `TOOLS.md`; the go-email-reader
   binary is built there (`~/workspace/go-email-reader/go-email-reader`).
3. **gitmap setting**: `email.go_email_reader_bin` (absolute path). Read via
   `config.GetVariable("global", "email.go_email_reader_bin")`; when set,
   verify the path is an executable file before use.

The first candidate that resolves to an existing executable file wins. When
none resolves, the built-in SMTP sender (§3.3) is used. Resolution result is
surfaced in output so the user always knows which path sent the mail, e.g.
`✔ sent via go-email-reader (/home/user/go/bin/go-email-reader)` vs
`✔ sent via built-in SMTP (smtp.gmail.com:587)`.

### 3.3a Delegation path (binary found)

Exec `go-email-reader send` with the user's flags passed through verbatim:

```
go-email-reader send --to <to> [--subject <s>] [--body <b>] [--file <f>] [--html]
```

Rules:

- Build the argument list programmatically from the parsed gitmap flags;
  pass through ONLY `--to`, `--subject`, `--body`, `--file`, `--html`.
- NEVER pass a password (or any credential) on the command line.
  go-email-reader reads its own RSA-encrypted SQLite store
  (`~/workspace/go-email-reader` memory notes: 2048-bit keypair, 0600
  private key, passwords write-only). gitmap has nothing to supply.
- A non-zero exit from the child is propagated: wrap the exit error with
  context (`sending via go-email-reader`) and exit non-zero. Do not swallow
  stderr — relay the child's error text to the user.
- `--file` exclusivity is enforced by go-email-reader itself too
  (`cmd_mail.go`: `--file` cannot be combined with `--body`, `--body-file`,
  or `--template`); gitmap enforces the same rule before exec'ing so the
  error message is gitmap-styled and consistent across both paths.

### 3.3b Built-in SMTP path (binary NOT found) — NEW code

New authenticated SMTP sender. It must NOT reuse `dispatchSmtpEmail` as-is
(nil auth, port 25). A small, dedicated package file holds the sender.

**Settings (all non-secret, stored via `config.SetVariable("global", …)`):**

| Setting key | Meaning | Default |
|---|---|---|
| `email.smtp_host` | SMTP host (no port) | `smtp.gmail.com` |
| `email.smtp_port` | SMTP port | `587` for STARTTLS, `465` for implicit TLS |
| `email.smtp_tls_mode` | `starttls` \| `implicit` \| `none` | `starttls` |
| `email.from` | Envelope + From header address | required (no default) |
| `email.smtp_user` | SMTP auth username (usually = from) | falls back to `email.from` |

Note: the legacy AGM setup wrote `email.smtp_host` as `host:port`
(`smtp.gmail.com:587`); the new keys split host and port. On read, if
`email.smtp_host` contains a colon, split it once and prefer the explicit
`email.smtp_port` only when set — preserving backward compatibility with
existing AGM-configured machines without silently mis-parsing.

**Password (NEVER stored in plaintext settings):**

1. First, look for a vaulted password in `email.smtp_password_ref` — the
   setting holds only a `rsa:`-prefixed ciphertext (or any
   `DecryptStoredPassword`-compatible encoding). Decrypt via
   `secrets.DecryptStoredPassword`; if decryption yields the unchanged input
   and it does not carry a known prefix, treat the ref as absent rather than
   as a plaintext password.
2. Otherwise, prompt at runtime with a masked prompt:
   `term.ReadPassword(int(os.Stdin.Fd()))`, same shape as
   `promptPassword` in `cli/cmdos/os_autologin_prompt.go`. Precede it with a
   `term.IsTerminal(int(os.Stdin.Fd()))` check — when stdin is NOT a TTY,
   FAIL CLOSED with an explicit error
   (`email password required but stdin is not a terminal; set email.smtp_password_ref to a vaulted secret`),
   never reading a line from piped stdin (that would echo-allow secrets into
   logs and hang automation).
3. An optional `--password-stdin` escape hatch is explicitly REJECTED for v1
   (command-line/process-list exposure); document the vault ref as the
   automation answer.

**Sending rules:**

- `starttls`: dial `host:port`, `smtp.Client.Hello`, upgrade with
  `StartTLS` (verify server cert — no `InsecureSkipVerify`), then
  `smtp.PlainAuth("", user, password, host)`.
- `implicit`: `tls.Dial` first (port 465 default), then `smtp.NewClient`.
- `none`: allowed only for local/dev relays; requires the host to be
  `localhost`/`127.0.0.1`/`::1`, otherwise fail closed with an error.
  (This is the narrow, fenced carve-out for the legacy port-25 relay case —
  never a silent default.)
- MIME: `From`, `To`, `Subject` headers; `Content-Type: text/plain;
  charset=utf-8` or `text/html; charset=utf-8` per `--html`;
  `MIME-Version: 1.0`; `\r\n` line endings; body encoded as-is (UTF-8).
- All errors wrapped with context (`fmt.Errorf("smtp send via %s:%s: %w",
  …)`); a failed send exits non-zero with the provider's message relayed.

### 3.4 `--html`

Body sent with `Content-Type: text/html; charset=utf-8`. On the delegation
path this is the `--html` flag; on the built-in path it sets the MIME
content type. A file with `html: true` (§3.5) implies it.

### 3.5 `--file` parsing — IDENTICAL to go-email-reader

Replicate `~/workspace/go-email-reader/internal/mail/parse_email_file.go`
semantics exactly (extension selects the parser):

- **YAML** (`.yaml`/`.yml`) / **JSON** (`.json`): keys
  `to, cc, bcc, subject, body, html` (all strings except `html` bool).
  Unknown keys ignored. Malformed document → hard error, no partial send.
- **Markdown** (`.md`/`.markdown`): the first line whose trimmed form starts
  with `# ` becomes the subject (prefix stripped, trimmed); everything after
  it is the body. Lines before the first `# ` heading are dropped.
- **Plain text** (anything else): the first non-empty line is the subject —
  an optional `Subject:` prefix is stripped case-insensitively
  (`len > 8 && EqualFold(s[:8], "subject:")`); everything after it is the
  body. Blank lines inside the body are preserved; leading blank lines of
  the body are skipped; trailing newlines trimmed.
- CRLF is normalized to LF before parsing in all cases.
- Precedence when `--file` is combined with explicit flags (delegation
  path relies on go-email-reader's own rule — explicit flags win — but the
  built-in path applies the same): `--to` overrides file `to`;
  `--subject` overrides file `subject`; `--html` forces HTML even when the
  file says `html: false`; `cc`/`bcc` come from the file only.
- `--file` combined with `--subject` or `--body` is refused on BOTH paths
  (gitmap-side check first, so the error is uniform).

## 4. Command 2 — `gitmap space schedule-send-email`

### 4.1 Invocation

```
gitmap space schedule-send-email --to X --task <enum> [--remove]
```

- `--task` is required and validated against the enum below.
- `--remove` deletes the registered trigger for that task (idempotent).
- Bare `gitmap space schedule-send-email` (no flags) lists registered
  triggers as a table (`task | to | added-at`).

### 4.2 The task-type enum — `EmailTaskType`

Dedicated Go string type in its own file (guideline: definitions live in
dedicated files; enum name ends with `Type`):

```go
type EmailTaskType string
const (
    EmailTaskRelease EmailTaskType = "release"
    EmailTaskPull    EmailTaskType = "pull"
    EmailTaskClone   EmailTaskType = "clone"
    EmailTaskScan    EmailTaskType = "scan"
    EmailTaskPush    EmailTaskType = "push"
    EmailTaskBackup  EmailTaskType = "backup"
)
```

- `ParseEmailTaskType(s string) (EmailTaskType, error)` — trims,
  lowercases, validates; unknown values return a descriptive error listing
  the valid values.
- All switches on `EmailTaskType` are exhaustive (or carry a default branch
  that fails closed).
- Extendable: adding a value means adding a constant + the command-name
  mapping below; the parser and switch stay exhaustive by construction.

Mapping from the executed top-level gitmap command to `EmailTaskType`
(lives next to the enum; command names are the existing `cli/constants`
symbols, not raw strings):

| gitmap command | EmailTaskType |
|---|---|
| `release` | `EmailTaskRelease` |
| `pull` | `EmailTaskPull` |
| `clone` | `EmailTaskClone` |
| `scan` | `EmailTaskScan` |
| `push` | `EmailTaskPush` |
| `backup` | `EmailTaskBackup` |

### 4.3 Registration (settings)

One setting per task: `email.schedule.<task>` = recipient address.
Example: `email.schedule.pull` = `user@example.com`.

- `schedule-send-email --to X --task pull` writes
  `config.SetVariable("global", "email.schedule.pull", X)`.
- `--remove` deletes the key (settings delete must not leave empty values;
  use the config package's delete path, not an empty-string write).
- Multiple tasks may each have their own address; one task has at most one
  address (last write wins).
- Validation on write: `--to` must be a syntactically valid address
  (parse via `net/mail.ParseAddress`); `--task` must parse.

### 4.4 Trigger mechanism (the hook)

Single choke point: `handleDispatchResult` in `cli/cmd/root.go:544` — every
dispatched command funnels through it with `(command, found, err)`.

New function `maybeFireScheduleTrigger(command string, err error)` called
from `handleDispatchResult` AFTER the existing audit finish, BEFORE the
function returns, with these rules:

1. Only when `found` is true (an actual command ran).
2. Map `command` → `EmailTaskType`; unmapped commands return immediately
   (no settings read — zero overhead for the 500+ other commands).
3. Read `email.schedule.<task>`; absent → return silently.
4. Build the report (§4.5) and send it via the SAME send pipeline as
   `space send-email` (binary resolution → delegation or built-in SMTP),
   forced to HTML. The send is best-effort and non-blocking for the user's
   command result: a trigger-send failure is logged (stderr warning) but
   NEVER changes the exit code of the operation that fired it, and NEVER
   panics.
5. Recursion guard: the trigger must not fire for `space` itself
   (`space send-email` / `space schedule-send-email` never appear in the
   mapping table, but the guard is explicit — a named constant set of
   excluded commands), and must not fire when the send was itself triggered
   (process-env sentinel `GITMAP_EMAIL_TRIGGER_DEPTH=1` set around the
   send; when present, `maybeFireScheduleTrigger` returns immediately).

`err` from the dispatch gives the outcome: nil → success; non-nil → failure
with the message.

### 4.5 Report email content (HTML, nicely formatted)

Subject: `[gitmap] <Task> <succeeded|failed> — <repo> (<branch>) on <hostname>`

Body sections (HTML table / headings, inline CSS only, no external assets):

- **Machine**: hostname (`os.Hostname()`), OS user
  (`os/user.Current().Username`), IP address
  (`cmdip.GetLocalIP(ctx, true, "")` — skip loopback; on error show
  `unavailable` rather than failing the email).
- **Operation**: task, full argv (`gitmap pull --all` etc. from
  `os.Args`), working directory, started/finished timestamps
  (RFC3339, local time).
- **Repository**: repo root path, current branch (`git rev-parse
  --abbrev-ref HEAD`), remote origin URL (best-effort; `unknown` when
  not a repo).
- **Outcome**: `succeeded` or `failed`; on failure the first N lines of the
  error (cap at 20 lines / 4 KB — never dump unbounded logs into mail);
  on success a one-line summary passed by the caller where available.
- Footer: `Sent by gitmap schedule-send-email (spec 276)`.

For operations that span many repos (e.g. `gitmap pull` over a fleet),
the report covers the aggregate outcome (counts of ok/failed, failed repo
list capped at 20) rather than one email per repo — exactly one email per
fired trigger per invocation.

## 5. Decision matrix — delegate vs built-in

| Condition | Path | Failure behavior |
|---|---|---|
| go-email-reader executable found (PATH → `~/go/bin` → `email.go_email_reader_bin`) | DELEGATE: `go-email-reader send …` | Child non-zero → wrap + propagate non-zero; stderr relayed verbatim. go-email-reader owns account selection, auth, and its RSA store — gitmap adds nothing |
| No binary found | BUILT-IN: authenticated SMTP from §3.3b | Auth/send error → wrapped `apperror`/fmt error, non-zero exit, provider message relayed. Password never logged, never echoed |
| Binary found but not executable / broken | Treat as NOT found (fall through to built-in) — but print a warning naming the bad path so the user can fix `email.go_email_reader_bin` | Same as built-in row |
| `--file` given | Parsing happens gitmap-side for validation; on the delegation path the FILE PATH is passed through (`--file <f>`) — not the parsed content — so go-email-reader remains the parser of record and semantics cannot drift | Parse error → hard error before any send attempt, on both paths |
| Non-TTY stdin and password needed (built-in path) | Fail closed: error, no send, no fallback to reading piped stdin | Exit non-zero with the actionable message naming `email.smtp_password_ref` |
| Trigger send fails (`schedule-send-email` hook) | Best-effort: stderr warning only | NEVER alters the firing command's exit code |

## 6. Settings key table (all under `config` global scope)

| Key | Secret? | Set by |
|---|---|---|
| `email.go_email_reader_bin` | no | user (manual) |
| `email.smtp_host` | no | user / future setup command |
| `email.smtp_port` | no | user / future setup command |
| `email.smtp_tls_mode` | no | user / future setup command |
| `email.from` | no | user / future setup command |
| `email.smtp_user` | no | user / future setup command |
| `email.smtp_password_ref` | YES (ciphertext only, `rsa:`/`aes:`/`salt:`/`caesar:` prefix) | user via vault tooling — NEVER via a plaintext `--password` flag |
| `email.schedule.<task>` | no | `space schedule-send-email` |

Hard rule: no new setting key may ever hold a plaintext password. The
existing `email.password` key (written by `gitmap email setup` at
`cli/cmdagy/agy_telegram_email_settings.go:259`) is the documented
anti-pattern; the new code must not read it as a password source and must
not write it.

## 7. File plan (new package `cli/cmdspaceemail/`, each file ≤ 300 lines)

The `space` namespace owns the UX (`DispatchSpace` branches in
`cli/cmdspace/space.go` + `CmdSpaceSendEmail`/`CmdSpaceScheduleSendEmail`
constants + usage text), but the email logic lives in a dedicated package
so `cli/cmdspace` stays within its size caps and concerns stay grouped.

| File | Responsibility | ≤ lines |
|---|---|---|
| `cli/cmdspaceemail/email_task_type.go` | `EmailTaskType`, constants, `ParseEmailTaskType`, command→task mapping | 120 |
| `cli/cmdspaceemail/binary_resolve.go` | Resolution order PATH → `~/go/bin` → `email.go_email_reader_bin`; executable check; result struct `{Path, IsFound, Source}` | 150 |
| `cli/cmdspaceemail/send_email.go` | `space send-email` arg parsing, `--file`/`--subject`/`--body` exclusivity check, path selection (delegate vs built-in), user-facing output | 250 |
| `cli/cmdspaceemail/delegate_send.go` | Build `go-email-reader send …` argv, exec, stream output, wrap child errors | 150 |
| `cli/cmdspaceemail/parse_email_file.go` | YAML/JSON/MD/text parsing per §3.5 (mirror of go-email-reader semantics; cite the source file in a comment) | 200 |
| `cli/cmdspaceemail/smtp_builtin.go` | Authenticated SMTP: settings load, TLS modes, `smtp.PlainAuth`, MIME build, send | 280 |
| `cli/cmdspaceemail/smtp_password.go` | Vault-ref decrypt (`secrets.DecryptStoredPassword`) → masked prompt (`term.ReadPassword`, fail closed on non-TTY, hookable var like `cli/cmdssh/ssh_login_pass_prompt.go:16`) | 150 |
| `cli/cmdspaceemail/schedule_send_email.go` | `space schedule-send-email` arg parsing, `--task` validation, register/remove/list | 200 |
| `cli/cmdspaceemail/schedule_trigger.go` | `maybeFireScheduleTrigger`, recursion guard, report data collection (hostname/user/IP/repo/branch), HTML render, fire-and-log send | 280 |
| `cli/cmdspaceemail/help.go` | Usage strings for both subcommands | 120 |
| `cli/cmdspaceemail/*_test.go` | Unit tests colocated per file (parser table tests, enum parse tests, resolution-order tests with temp dirs, TLS-mode validation) | — |

`DispatchSpace` additions (in `cli/cmdspace/space.go`, following the existing
branch pattern):

```go
if sub == constants.CmdSpaceSendEmail {
    return true, cmdspaceemail.RunSpaceSendEmail(rest[1:])
}
if sub == constants.CmdSpaceScheduleSendEmail {
    return true, cmdspaceemail.RunSpaceScheduleSendEmail(rest[1:])
}
```

`handleDispatchResult` addition (in `cli/cmd/root.go:544`): one call —
`cmdspaceemail.MaybeFireScheduleTrigger(command, err)` — after
`cmdaudit.FinishCommandAudit`, guarded internally per §4.4. Its own errors
are swallowed into a stderr warning inside the function; it never returns
an error to the dispatcher.

## 8. Coding constraints (from `.ai-memory/coding-guidelines.md` — binding)

- **File size**: any Go file 300 lines max (file plan above is sized for it).
- **Function length**: 8 lines preferred, 15 lines hard cap (blank lines and
  comments excluded); extract helpers instead of growing functions.
- **No nested `if`**: flatten with early returns and guard clauses.
- **Positive booleans only**: every bool variable/function/field starts with
  `is` or `has` (`isFound`, `hasAuth`, `isTerminal`); no `!` on complex
  conditions — extract positively named booleans. No `can`/`should`/`was`.
- **No magic strings/numbers**: settings keys, flag names, TLS modes, ports
  (587/465), task names — all named constants or the `EmailTaskType` enum.
- **Enum rules**: `EmailTaskType` ends with `Type`; definitions in a
  dedicated file; switches exhaustive or fail-closed default.
- **Errors**: wrapped with context at every boundary
  (`fmt.Errorf("resolve go-email-reader binary: %w", err)`); user-facing
  failures carry the operation name and key inputs; nothing swallowed —
  except the deliberate, documented best-effort swallow inside
  `MaybeFireScheduleTrigger`, which logs to stderr.
- **No `any`/`interface{}`** in new signatures; narrow at trust boundaries
  (YAML/JSON decode targets are typed structs, same as go-email-reader's
  `emailFileDoc`).
- **DRY**: the send pipeline (resolve → delegate-or-builtin) is implemented
  ONCE and shared by `space send-email` and the schedule trigger.
- **Tests**: table-driven unit tests for the file parser (parity with
  go-email-reader edge cases), enum parsing, resolution order; behavioral
  tests against a mock SMTP server verifying MIME; E2E: real
  `gitmap pull` in a scratch clone with a registered trigger.
