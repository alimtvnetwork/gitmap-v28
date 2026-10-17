# Spec 276 — Acceptance criteria and evidence

How to read this file: each AC is numbered (AC-01…). Every AC names the
exact evidence the implementer must produce. "Evidence" means commands run
and their real output, not claims. Per the repo's standing rule, never
trust a self-reported "it works" — verify personally and paste the
receipts.

Global gates for ALL ACs: `go build ./...` exit 0, `go vet` clean on
touched packages, every new Go file ≤ 300 lines, no nested `if`s, boolean
names `is*`/`has*` only, no magic strings, errors wrapped with context.
Unit tests colocated as `cli/cmdspaceemail/*_test.go`.

---

## A. `space send-email` — argument handling

**AC-01** — `--to` is required.
`gitmap space send-email --subject S --body B` (no `--to`) exits non-zero
with a message naming the missing `--to` flag.
Evidence: terminal transcript.

**AC-02** — `--file` is mutually exclusive with `--subject`/`--body`.
`gitmap space send-email --to X --file F --subject S` exits non-zero on
BOTH paths (delegation and built-in) with a uniform gitmap-styled message,
before any send is attempted.
Evidence: transcript on a machine WITH go-email-reader on PATH and on a
machine WITHOUT it.

**AC-03** — `--subject`/`--body` required when `--file` is absent.
`gitmap space send-email --to X` (neither `--file` nor `--subject`) exits
non-zero naming what is missing.
Evidence: transcript.

**AC-04** — `--help` renders.
`gitmap space send-email --help` and `gitmap space --help` list the new
subcommand with its flags and one example each.
Evidence: transcript.

## B. Binary resolution order (AC-05 – AC-08)

**AC-05** — PATH wins.
With `go-email-reader` on PATH, output contains
`sent via go-email-reader` and the resolved PATH location; the child
receives exactly `--to/--subject/--body/[--file]/[--html]` and nothing
else (verify by stubbing `go-email-reader` with a script that echoes its
argv into a file).
Evidence: stub-argv capture file + transcript.

**AC-06** — `~/go/bin` fallback.
With PATH cleared of go-email-reader but an executable at
`~/go/bin/go-email-reader`, the `~/go/bin` binary is used and its path is
printed.
Evidence: transcript with `PATH` manipulated in the test shell.

**AC-07** — `email.go_email_reader_bin` setting fallback.
With no PATH/`~/go/bin` binary, setting
`email.go_email_reader_bin=/custom/path/go-email-reader` (via
`config.SetVariable`-equivalent user command) routes delegation to that
exact path. A non-existent or non-executable configured path is treated as
NOT found (falls through to built-in) AND prints a warning naming the bad
path.
Evidence: transcript for both the good-path and bad-path cases.

**AC-08** — No credential ever on the child command line.
The stub-argv capture from AC-05/AC-07 contains no password, token, or
`--password` in any form, even when the built-in path would have prompted.
Evidence: captured argv file inspected and pasted.

## C. Built-in SMTP sender (AC-09 – AC-14)

**AC-09** — Authenticated, never nil-auth.
Code inspection: the new sender constructs `smtp.PlainAuth` (STARTTLS) or
uses an implicit-TLS client; there is NO `smtp.SendMail(…, nil, …)` call
in any new file. The legacy `dispatchSmtpEmail` in
`cli/cmdagy/agy_account_switch_notify.go` is untouched.
Evidence: `gitmap aum search "SendMail"` output showing the only nil-auth
call site remains the legacy one.

**AC-10** — TLS modes behave.
Against a local mock SMTP server: `starttls` upgrades and authenticates;
`implicit` dials TLS first; `none` against a non-loopback host fails closed
with an explicit error; `none` against `127.0.0.1` sends. Unknown
`email.smtp_tls_mode` value fails closed with the valid values listed.
Evidence: mock-server logs + transcripts.

**AC-11** — MIME correctness.
A sent message captured by the mock server has `From`, `To`, `Subject`,
`MIME-Version: 1.0`, `Content-Type: text/plain; charset=utf-8` (default)
or `text/html; charset=utf-8` (with `--html`), `\r\n` line endings.
Evidence: raw captured message pasted.

**AC-12** — Password never in plaintext settings.
After a successful built-in send, `gitmap aum search "email.password"`
shows NO new write site; the settings store contains no new
password-bearing key; the only acceptable secret key is
`email.smtp_password_ref` holding a prefixed ciphertext (`rsa:` etc.).
Evidence: search output + settings dump (secret value redacted, prefix
visible).

**AC-13** — Masked prompt, fail closed on non-TTY.
Interactive: password entry echoes nothing (document the manual check).
Piped stdin (`echo secret | gitmap space send-email …` with no vault ref):
exits non-zero with the message naming `email.smtp_password_ref`, sends
nothing, and does not hang.
Evidence: transcript of the piped case with `timeout 10` proving no hang.

**AC-14** — Vault ref path.
`email.smtp_password_ref` holding an `rsa:` ciphertext decrypts via
`secrets.DecryptStoredPassword` and sends without prompting (test with a
locally generated keypair per `cli/secrets/ssh_password.go` expectations).
A ref whose value decrypts to itself unchanged is treated as absent →
falls back to the prompt (AC-13).
Evidence: transcript + unit test.

## D. `--file` parity with go-email-reader (AC-15 – AC-18)

Parity target: `~/workspace/go-email-reader/internal/mail/parse_email_file.go`
(`EmailFile{To, CC, BCC, Subject, Body, HTML}`).

**AC-15** — Table-driven unit tests cover, at minimum:
- YAML: all six keys; unknown keys ignored; malformed YAML → error.
- JSON: same as YAML.
- MD: first `# ` heading → subject, rest → body; lines before the first
  heading dropped; file with no heading → empty subject, whole body.
- Text: first non-empty line → subject; `Subject:` prefix stripped
  case-insensitively (`SUBJECT:`, `subject:`); blank lines inside body
  preserved; leading blank lines of body skipped; trailing newlines trimmed.
- CRLF input normalized in all four formats.
Evidence: `go test ./cli/cmdspaceemail/ -run TestParseEmailFile -v` output.

**AC-16** — Flag precedence on the built-in path.
A file with `to/subject/html:false` + explicit `--to`/`--subject`/`--html`
flags → flags win for to/subject; `--html` forces HTML. `cc`/`bcc` come
from the file only.
Evidence: mock-server captured message.

**AC-17** — Delegation path passes the FILE, not parsed content.
With `--file`, the child argv contains `--file <exact path given>` and no
`--subject`/`--body` (go-email-reader stays the parser of record).
Evidence: stub-argv capture.

**AC-18** — Missing/unreadable `--file` fails before any send.
Evidence: transcript (no mock-server connection attempted — assert via
mock log showing zero connections).

## E. `space schedule-send-email` (AC-19 – AC-25)

**AC-19** — Enum validation.
`--task bogus` exits non-zero listing exactly
`release | pull | clone | scan | push | backup`. Each valid value registers.
Evidence: transcript.

**AC-20** — Register / list / remove.
`--to X --task pull` writes `email.schedule.pull`; bare invocation lists a
table with the entry; `--task pull --remove` deletes it (second remove is a
no-op success); invalid address rejected via `net/mail.ParseAddress`.
Evidence: transcript + settings dump.

**AC-21** — Trigger fires after a real operation.
In a scratch clone: register `--task pull`, run `gitmap pull`, assert
exactly one email was captured by the mock SMTP server (or delegated
child), with subject `[gitmap] pull succeeded — <repo> (<branch>) on
<hostname>`.
Evidence: mock-server capture + transcript.

**AC-22** — Failure outcome reported.
Register `--task clone`, run a `gitmap clone` against an unreachable URL,
assert the email subject contains `failed` and the body contains a capped
error excerpt (≤ 20 lines / 4 KB), and the clone's own non-zero exit code
is unchanged.
Evidence: mock-server capture + `echo $?`.

**AC-23** — Trigger send never changes the operation's result.
Force the trigger send to fail (bad SMTP host in settings); the firing
command's exit code and stdout are identical to the no-trigger run; a
stderr warning names the trigger failure.
Evidence: two transcripts diffed.

**AC-24** — No recursion, no fire on unmapped commands.
`gitmap space send-email` and `gitmap space schedule-send-email` never
fire triggers; running `gitmap status` (unmapped) with triggers registered
sends nothing; `GITMAP_EMAIL_TRIGGER_DEPTH=1` suppresses firing.
Evidence: mock-server log showing zero sends for these cases.

**AC-25** — Report content completeness.
The HTML body contains: hostname, OS user, IP address (or `unavailable`),
operation + full argv + cwd + timestamps, repo root + branch + origin URL
(or `unknown`), outcome section. Fleet operations (`pull` over many repos)
produce ONE aggregate email (ok/failed counts, failed list capped at 20),
not one per repo.
Evidence: captured HTML pasted (redact nothing except secrets — there are
none in the report).

## F. Decision matrix (AC-26 – AC-27)

**AC-26** — Matrix row coverage.
Each row of the decision matrix in `01-scope-and-plan.md §5` has at least
one AC above (map: delegate→AC-05/AC-08; built-in→AC-09–AC-14;
broken-binary→AC-07; `--file`→AC-17/AC-18; non-TTY→AC-13;
trigger-failure→AC-23). No row is untested.
Evidence: this checklist, checked off in the implementation report.

**AC-27** — Path visibility.
Every send (manual or triggered) prints which path was used and, for
delegation, the resolved binary path; for built-in, `host:port` and TLS
mode. The user can always answer "how did this email go out?".
Evidence: transcripts.

## G. Coding-constraints audit (AC-28)

**AC-28** — Automated + manual constraint sweep on all new files:
`gitmap aum search` for `!` in `if` conditions (must be zero outside
positively named booleans), file line counts ≤ 300, `go vet` clean,
settings keys referenced only via named constants, error wraps present at
every boundary (`%w`), enum switches exhaustive.
Evidence: search outputs + line counts pasted into the implementation
report.

---

## Evidence packaging

The implementer delivers, in the implementation report (not in this spec):

1. Per-AC transcripts / test outputs as run (copy-paste, unedited).
2. The AC-26 matrix checklist with each row checked.
3. `go build ./...`, `go vet`, and full `go test ./cli/cmdspaceemail/`
   outputs.
4. Confirmation that `cli/cmdagy/agy_account_switch_notify.go` and
   `cli/cmdagy/agy_telegram_email_settings.go` are byte-identical to before
   (`git diff --stat` shows no changes under `cli/cmdagy/`).
