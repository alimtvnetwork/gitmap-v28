# Forbidden Strings — checker, summary, and waiver tools

Three companion tools sit on top of the same TOML config:

| Tool | Purpose | Exit code on findings |
|------|---------|-----------------------|
| `check-forbidden-strings.py` | CI gate — full per-line listing of every match | `1` |
| `forbidden-strings-summary.py` | Per-rule digest + **exact runnable fix command** | `1` |
| `allowlist-forbidden-string.py` | Add legitimate paths to a rule's `allowlist` (with mandatory `--reason`) | `0` on success |

Single source of truth: [`forbidden-strings.toml`](./forbidden-strings.toml).

## Quick start

```bash

# 1. CI failed? Get a focused report with a one-liner fix per rule.

npm run lint:strings:report

# 2. Apply the fix for one rule directly.

python3 linter-scripts/forbidden-strings-summary.py \
    --rule STALE-MODULE-PATH --emit-fix-command | bash

# 3. The match is legitimate (audit trail, changelog, fixture) — waive it.

python3 linter-scripts/allowlist-forbidden-string.py \
    --rule STALE-MODULE-PATH \
    --reason "Historical changelog entry; documents the rule itself." \
    changelog.md

# 4. Re-verify.

npm run lint:strings
```

## TOML schema

```toml
[[rule]]
id          = "STALE-MODULE-PATH"
description = "Stale module path reference (movie-cli-v1). Current canonical namespace is movie-cli-v2."
pattern     = 'movie-cli-v1\b'           # Python regex; literal TOML string
replacement = "movie-cli-v2"             # Canonical replacement (used by the summary tool)
fix_hint    = "python3 linter-scripts/forbidden-strings-summary.py --rule STALE-MODULE-PATH --emit-fix-command | bash"
exclude_dirs  = ["release-artifacts"]    # never scanned
exclude_files = ["forbidden-strings.toml", "check-forbidden-strings.py", "forbidden-strings-summary.py"]
allowlist = [
  # Repo-relative paths or globs that may legitimately contain the pattern.
  # Add via `allowlist-forbidden-string.py` to keep the justification comment.
  "spec/14-update/23-install-script-version-probe.md",
]
```

The `replacement` field is **strongly recommended** for any rename
guard — without it, the summary tool falls back to "no `replacement`
declared — fix manually" and operators have to invent the substitution
themselves.

## Summary report — output formats

- **Text** (default): grouped per rule, lists every affected file with
  hit count, plus a single concrete `sed` command.
- **Markdown** (`--markdown`): same content as a Markdown table; also
  auto-appended to `$GITHUB_STEP_SUMMARY` when running in CI so the
  report shows inline on the PR.
- **Just the command** (`--emit-fix-command --rule <ID>`): emits one
  line, safe to pipe to `bash`. If there are no findings, it emits
  `: nothing to fix` (a no-op shell command) so piping is always safe.

## Why a separate summary tool

`check-forbidden-strings.py` is intentionally verbose — it shows file +
line + content for every hit so reviewers can audit. The summary tool
distils the same data into the form an operator actually needs to act
on it: *which files, what's the exact replacement, give me the
command.* And because the command uses the **actual matched text from
the scan** (not a hardcoded template), it can fan out across multiple
stale variants in a single invocation — for example, the `STALE-REPO-SLUG`
rule may match `coding-guidelines-v24`, `coding-guidelines-v24`, and
`coding-guidelines-v24` in the same file, and the summary will emit
three `-e` substitutions in one `sed` call rather than three separate
templates the operator has to expand by hand.
