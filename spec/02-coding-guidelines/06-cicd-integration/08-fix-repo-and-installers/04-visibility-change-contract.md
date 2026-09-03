# `visibility-change.sh` / `visibility-change.ps1` — Normative Contract

> **Version:** 1.0.0 · **Updated:** 2026-04-28
> Implements: `visibility-change.sh`, `visibility-change.ps1`, `scripts/visibility-change/*`

---

## 1. Goal

Toggle (or explicitly set) the visibility of the current repo on
**GitHub** or **GitLab**, with an interactive safety prompt for
`private → public` transitions.

---

## 2. Flag surface

| Bash | PowerShell | Meaning |
|---|---|---|
| `--visible pub` | `-Visible pub` | Force visibility to **public** |
| `--visible pri` | `-Visible pri` | Force visibility to **private** |
| (none) | (none) | **Toggle** current visibility |
| `--yes` / `-y` | `-Yes` | Skip the `private → public` confirmation prompt |
| `--dry-run` | `-DryRun` | Print what would change; no API calls |
| `--help` / `-h` | `-Help` / `-h` | Print help; exit `0` |

`--visible` accepts (case-insensitive): `pub`, `public`, `pri`, `private`.
Anything else → exit `6`.

---

## 3. Provider detection

1. `git remote get-url origin` → URL string.
2. Match host:
   - `github.com` or `ssh.github.com` → **GitHub**
   - `gitlab.com` or any host listed in `$VISIBILITY_GITLAB_HOSTS`
     (comma-separated env var) → **GitLab**
3. No match → exit `4` with a one-line message naming the detected host.

---

## 4. Auth backend

| Provider | Apply | Read | Missing CLI |
|---|---|---|---|
| GitHub | `gh repo edit <o>/<r> --visibility <v> --accept-visibility-change-consequences` | `gh repo view --json visibility -q .visibility` | exit `5` + install URL |
| GitLab | `glab repo edit <o>/<r> --visibility <v>` | `glab repo view -F json` | exit `5` + install URL |

Auth errors from the CLI are re-printed verbatim and surface as exit
`5`. The user is expected to run `gh auth login` / `glab auth login`
themselves; this script does NOT wrap the browser flow.

---

## 5. Toggle logic

```
current = read_current_visibility()
target  = args.Visible
       || (current == "public" ? "private" : "public")

if current == target:           print "already <target>"; exit 0
if target == "public" and !Yes: confirm_or_exit_7()
if DryRun:                      print "[dry-run] would change <c> → <t>"; exit 0

apply(target)
verify(target)                  # re-read; fail with exit 8 on mismatch
print "visibility: <c> → <t> (<provider>)"
exit 0
```

---

## 6. Confirmation prompt

Triggered ONLY when `current=private`, `target=public`, and `--yes` not set:

```
⚠  About to make <owner>/<repo> PUBLIC on <provider>.
   URL: <html_url>
   Type 'yes' to continue, anything else aborts:
```

Non-interactive stdin (piped) → exit `7` with hint to pass `--yes`.

---

## 7. Runner integration

`./run.sh visibility [flags]` and `run.ps1 visibility [flags]` forward
all flags verbatim. The standalone scripts are the source of truth; the
runner adds ≤ 5 lines per language.

---

## 8. Exit codes

| Code | Meaning |
|---|---|
| `0` | Success (changed, already-target, or dry-run) |
| `2` | Not inside a Git work tree |
| `3` | No `origin` remote configured |
| `4` | Unsupported provider |
| `5` | Auth/CLI failure (missing CLI, not logged in, API error) |
| `6` | Bad flag value |
| `7` | Confirmation required but stdin not interactive |
| `8` | Verification failed (apply returned 0 but visibility unchanged) |

---

## 9. Output contract

- Success → single line: `visibility: <old> → <new> (<provider>)`.
- Errors → stderr, one line per error, no stack traces.
- `--dry-run` → prefix every line with `[dry-run]`.
- Help text → plain ASCII, ≤ 80 columns.

---

## 10. CODE RED compliance

- Each script ≤ 300 lines, functions 8–15 effective lines.
- Zero nested conditionals — guard-clause only.
- Booleans positively named (`HasOrigin`, `IsGitHub`, `IsDryRun`).
- Errors never swallowed: every `try` has a logging `catch` + exit.
- Max 2 boolean operands per expression.
