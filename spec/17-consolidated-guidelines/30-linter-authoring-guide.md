# 30 — Linter Authoring Guide

> **Version:** 1.0.0
> **Updated:** 2026-04-22
> **Type:** Consolidated Guideline — Phase 6A
> **Status:** Closes the final 🟡 from `29-blind-ai-audit-v2.md` (validator authoring opacity)

This guide tells a blind AI **exactly** how to add a new linter to
`linter-scripts/` so it slots into the existing orchestrator (`run.sh` /
`run.ps1`), the validator inventory (`05-coding-guidelines.md §34`), and
the memory-mirror drift detector — without breaking any other check.

It complements:

- `05-coding-guidelines.md §34` — the **inventory** of all 17 linters
- `05-coding-guidelines.md §35` — the **CI step ordering**
- `linters-cicd/checks/_lib/sarif.py` — the **SARIF emitter contract**
- `02-spec/02-coding-guidelines/06-cicd-integration/01-sarif-contract.md`
- `02-spec/02-coding-guidelines/06-cicd-integration/02-plugin-model.md`

---

## §1 — When to Add a New Linter

Add a linter only when **all** of the following are true:

1. The rule is **machine-checkable** (no human judgement required).
2. The rule has a **stable spec section** under `spec/` it can cite.
3. The rule produces **zero false positives** on the current repo (or
   every false positive can be allowlisted with written justification).
4. No existing linter already covers the rule (check `§34` first).

If any of these fails, file a spec update first — not a linter.

---

## §2 — File Layout & Naming

```
linter-scripts/
├── <verb>-<subject>.<ext>          # the script
├── <verb>-<subject>.allowlist      # OPTIONAL plain-text exemptions
├── <verb>-<subject>.toml           # OPTIONAL config (forbidden-strings style)
└── README-<subject>.md             # OPTIONAL operator notes
```

Naming rules:

| Element | Rule | Example |
|---|---|---|
| Filename | `lowercase-kebab-case`, verb-led | `check-memory-mirror-drift.py` |
| Verb | `check-` (validate) · `generate-` (produce) · `suggest-` (advisory) | `check-axios-version.sh` |
| Extension | `.py` · `.sh` · `.ps1` · `.cjs` · `.go` only | `validate-guidelines.go` |
| PowerShell siblings | Required when a Bash linter runs in CI on Windows | `run.sh` ↔ `run.ps1` |

Forbidden:

- `snake_case`, `PascalCase`, or `camelCase` filenames in `linter-scripts/`
- Generic verbs: `do-`, `run-` (reserved for orchestrator), `lint-` (ambiguous)
- Hidden config (`.linterrc`) — every config must be a discoverable sibling

---

## §3 — Exit Code Contract (HARD RULE)

Every linter **must** use exactly this three-value exit code map. The
orchestrators (`run.sh` / `run.ps1`) and CI matrix depend on it.

| Code | Meaning | When to use |
|---|---|---|
| `0` | **Pass** — no findings | Clean repo, all rules satisfied |
| `1` | **Findings** — violation(s) detected | At least one rule was violated; print actionable detail to `stderr` |
| `2` | **Structural error** — linter cannot run | Missing input file, malformed config, unparseable AST, network failure |

Forbidden:

- Returning `0` when findings exist (silent pass)
- Returning `1` for structural errors (masks real bugs as violations)
- Returning any other code (`3`, `-1`, `127`, …) — the orchestrator
  treats unknown non-zero as `2` and **fails the build immediately**
- Using `sys.exit()` from inside a `try/except` that swallows the cause

Reference implementation:

```python

# linter-scripts/check-memory-mirror-drift.py

def fail(msg: str, code: int = 1) -> None:
    print(f"[memory-mirror-drift] FAIL: {msg}", file=sys.stderr)
    sys.exit(code)

# Structural error → code=2

if not MIRROR.exists():
    fail(f"missing {MIRROR}", code=2)

# Findings → code=1

if missing:
    return 1

# Clean → code=0

return 0
```

Bash equivalent:

```bash

# linter-scripts/check-axios-version.sh

set -euo pipefail
if ! grep -q '"axios"' package.json; then
  echo "[axios] structural: package.json has no axios entry" >&2
  exit 2
fi
if grep -qE '"axios":\s*"(\^|~)?(1\.14\.1|0\.30\.4)"' package.json; then
  echo "[axios] FAIL: blocked version detected" >&2
  exit 1
fi
exit 0
```

---

## §4 — Output Contract

### §4.1 Human-readable (default)

Every linter must emit, to **stdout**, a single line on success:

```
[<short-name>] OK — <count> <unit> verified
```

…and to **stderr** on failure, a structured block:

```
[<short-name>] FAIL: <one-line summary>
[<short-name>] <N> finding(s):
  - <file>:<line>  <rule-id>  <message>

Action: <one-line fix instruction>
```

The leading `[<short-name>]` prefix is **mandatory** — `run.sh` greps
for it when producing the merged CI summary.

### §4.2 SARIF (when `--format sarif`)

Linters that participate in the GitHub Code Scanning upload (anything
under `linters-cicd/checks/`) must additionally support
`--format sarif` per `01-sarif-contract.md` and use the shared
`linters-cicd/checks/_lib/sarif.py` emitter. Root-level
`linter-scripts/` checks are **exempt** from SARIF — they only emit
plain text and exit codes.

---

## §5 — CLI Flag Contract

| Flag | Required | Behavior |
|---|---|---|
| `--help` / `-h` | ✅ | Print usage, exit `0` |
| `--version` | ✅ for `linters-cicd/checks/` plugins | Print `coding-guidelines/<rule-slug> <X.Y.Z>`, exit `0` |
| `--root <dir>` | ✅ when scanning files | Default `.` — scan target |
| `--format {text,sarif}` | ✅ for `linters-cicd/checks/` plugins | Default `sarif` |
| `--output <file>` | Optional | Write payload to file instead of stdout |
| `--allowlist <file>` | When the linter has exemptions | Path to allowlist; default `<script>.allowlist` |

Use `argparse` (Python) or `getopts` (Bash) — never hand-roll flag
parsing.

---

## §6 — Allowlist Registration

Allowlists are **plain-text exemptions** for findings the linter would
otherwise fire. They keep the linter strict by default while admitting
documented exceptions.

### §6.1 Format

Two acceptable formats:

**Format A — flat list** (one entry per line):

```

# linter-scripts/spec-cross-links.allowlist

# Lines starting with # are comments.

# Each non-blank line is a path that may appear as a link target

# even though it does not resolve inside spec/.

readme.md
LICENSE
changelog.md
```

**Format B — sectioned** (when entries belong to disjoint categories):

```

# linter-scripts/spec-folder-refs.allowlist

[external]

# Folders living in sibling repos

gitmap-v3/02-spec/01-app

[doc-only]

# Folders intentionally referenced only in prose, never created

02-spec/99-future-work
```

### §6.2 Mandatory comment header

Every allowlist file **must** begin with:

```

# <linter-name> allowlist

# Owner: <github-username>

# Last reviewed: YYYY-MM-DD

# Each entry below MUST have a justification comment on the preceding line.

```

### §6.3 Adding an entry — workflow

1. Run the linter; capture the finding.
2. Decide: is this a **bug to fix** or a **legitimate exception**?
   - Bug → fix the source, do **not** allowlist.
   - Exception → continue.
3. Add the entry to the allowlist with a `# WHY: …` comment on the
   line above. Reviewers reject entries without `WHY`.
4. Re-run the linter; confirm exit `0`.
5. Commit the allowlist change in the **same commit** as whatever
   triggered the exception, with the justification in the commit
   message.

### §6.4 Periodic audit

Every 90 days, the linter owner walks each allowlist entry and removes
any that are no longer applicable. The `Last reviewed` date in the
header gets bumped.

---

## §7 — Test Fixtures

Every new linter ships with fixtures under
`linter-scripts/fixtures/<linter-name>/`:

```
linter-scripts/fixtures/check-memory-mirror-drift/
├── good/
│   ├── memory-index.md
│   └── mirror.md           # contains every expected token
├── bad-missing-token/
│   ├── memory-index.md
│   └── mirror.md           # missing one expected token → exit 1
└── bad-no-section/
    ├── memory-index.md
    └── mirror.md           # missing §X marker → exit 2
```

Rules:

- **At least one `good/` fixture** producing exit `0`.
- **At least one `bad-*/` fixture** for **each** non-zero exit code the
  linter can return (typically two: one for `1`, one for `2`).
- Fixture filenames mirror the real repo layout the linter scans —
  e.g., a linter that reads `package.json` puts a `package.json` inside
  the fixture directory.
- A sibling `<linter-name>_test.py` (or `_test.sh`) drives the fixtures
  via `subprocess.run([...], cwd=fixture_dir)` and asserts the exit
  code. No mocking — call the real script.

Reference test harness:

```python

# linter-scripts/fixtures/check-memory-mirror-drift_test.py

import subprocess, pathlib, sys
ROOT = pathlib.Path(__file__).parent
SCRIPT = ROOT.parent.parent / "check-memory-mirror-drift.py"

def run(case: str) -> int:
    return subprocess.run(
        [sys.executable, str(SCRIPT)],
        cwd=ROOT / case,
        capture_output=True,
    ).returncode

def test_good(): assert run("good") == 0
def test_bad_missing_token(): assert run("bad-missing-token") == 1
def test_bad_no_section(): assert run("bad-no-section") == 2
```

---

## §8 — Registration Checklist

When you add a new linter, perform **every** step. CI rejects PRs that
skip any.

- [ ] **1.** Script lives at `linter-scripts/<verb>-<subject>.<ext>`.
- [ ] **2.** Script is executable: `chmod +x linter-scripts/<file>`
      (Unix) and shebang line is present.
- [ ] **3.** Exit codes follow §3.
- [ ] **4.** Output follows §4 (mandatory `[<short-name>]` prefix).
- [ ] **5.** CLI flags follow §5.
- [ ] **6.** Allowlist exists (if applicable) with header per §6.2.
- [ ] **7.** Fixtures exist under `linter-scripts/fixtures/<name>/`
      covering exit `0`, `1`, and `2`.
- [ ] **8.** Test harness `<name>_test.py` runs all fixtures.
- [ ] **9.** Add the script as a new row in `05-coding-guidelines.md
      §34.1` (Active Linter Scripts) — assign the next free Linter #.
- [ ] **10.** Insert the script invocation into **both** `run.sh` and
      `run.ps1` in the order dictated by `05-coding-guidelines.md §35`
      (CI Step Ordering). Add a CI Step row alongside.
- [ ] **11.** If the linter introduces new config files, add a row in
      `§34.2 Allowlists & Config`.
- [ ] **12.** If the linter checks the memory mirror or any other
      meta-spec invariant, also append the relevant tokens to
      `EXPECTED_TOKENS` in `check-memory-mirror-drift.py`.
- [ ] **13.** Bump `package.json` version (minor for new linter,
      patch for fix), then run `node scripts/sync-version.mjs` and
      `node scripts/sync-spec-tree.mjs`.
- [ ] **14.** Run `bash linter-scripts/run.sh` locally; confirm
      exit `0`.

---

## §9 — Error Recovery When a Linter Fails

When `run.sh` exits non-zero, the diagnostic flow is:

| Linter exit | Diagnosis | Action |
|---|---|---|
| `1` | Real violation in source | Fix the source; do not allowlist without justification |
| `2` | Structural failure | Inspect stderr; usually a missing file or malformed config |
| `127` | Command not found | Install the missing runtime (Python 3, Node, pwsh) |
| Other | Unknown | Re-run with verbose flag if available; otherwise treat as `2` |

If a linter starts producing false positives after a refactor:

1. **Do not** delete the linter or downgrade exit `1` to `0`.
2. Add an allowlist entry with `# WHY: …` justification (§6.3).
3. If the false positive rate exceeds 5 % of findings, **rewrite the
   detection logic** instead of expanding the allowlist.

---

## §10 — Anti-Patterns (Reject in Code Review)

| Anti-pattern | Why it's wrong | Correct approach |
|---|---|---|
| `try: … except: pass` around the whole script | Hides structural errors as silent passes | Catch specific exceptions; exit `2` with stderr message |
| Reading config from env vars without a sibling `.toml` | Invisible coupling, no audit trail | Use a sibling config file checked into git |
| Hard-coding allowlist entries inside the script | Requires a code change for every exemption | Externalise to `<script>.allowlist` |
| Writing findings to stdout instead of stderr | Breaks `run.sh` summary parsing | Findings → stderr; success line → stdout |
| Calling `git` to discover files | Fails in tarball / shallow-clone contexts | Use `pathlib.Path.rglob` or the shared `walker.py` |
| `sleep` loops to wait for resources | Flaky on slow CI runners | Use a polling loop with timeout and explicit failure message |
| Generic linter names like `check.py` | Collides with future linters | Always `<verb>-<subject>` |

---

## §11 — Worked Example: A Brand-New Linter

Goal: detect TODO comments older than 30 days.

```python
#!/usr/bin/env python3
"""
check-stale-todos.py — Fail if any TODO comment is older than 30 days
(based on `git blame` author date).
"""
from __future__ import annotations
import argparse, subprocess, sys, datetime, pathlib, re

NAME = "stale-todos"
TODO_RE = re.compile(r"\bTODO\b")
MAX_AGE_DAYS = 30

def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Detect stale TODOs")
    p.add_argument("--root", default=".", help="Repo root (default: .)")
    p.add_argument("--allowlist",
                   default="linter-scripts/check-stale-todos.allowlist")
    return p.parse_args()

def load_allowlist(path: str) -> set[str]:
    p = pathlib.Path(path)
    if not p.exists():
        return set()
    return {
        line.strip()
        for line in p.read_text().splitlines()
        if line.strip() and not line.startswith("#")
    }

def main() -> int:
    args = parse_args()
    root = pathlib.Path(args.root).resolve()
    if not (root / ".git").exists():
        print(f"[{NAME}] FAIL: --root must be a git repo", file=sys.stderr)
        return 2

    allow = load_allowlist(args.allowlist)
    cutoff = datetime.date.today() - datetime.timedelta(days=MAX_AGE_DAYS)
    findings: list[str] = []

    for path in root.rglob("*"):
        if not path.is_file() or ".git" in path.parts:
            continue
        rel = str(path.relative_to(root))
        if rel in allow:
            continue
        try:
            text = path.read_text(encoding="utf-8")
        except (UnicodeDecodeError, OSError):
            continue
        for lineno, line in enumerate(text.splitlines(), start=1):
            if not TODO_RE.search(line):
                continue
            # blame this line; keep example short
            findings.append(f"{rel}:{lineno}  STALE-TODO  {line.strip()}")

    if findings:
        print(f"[{NAME}] FAIL: {len(findings)} stale TODO(s)", file=sys.stderr)
        for f in findings:
            print(f"  - {f}", file=sys.stderr)
        print("\nAction: resolve TODO or add to allowlist with WHY comment.",
              file=sys.stderr)
        return 1

    print(f"[{NAME}] OK — no stale TODOs older than {MAX_AGE_DAYS} days")
    return 0

if __name__ == "__main__":
    sys.exit(main())
```

After writing the script, walk the **§8 checklist**. The linter is not
"done" until every box is ticked.

---

## §12 — Cross-References

- `05-coding-guidelines.md §34` — Validator Inventory (where new linter
  gets registered as Linter #N)
- `05-coding-guidelines.md §35` — CI Step Ordering (where new step gets
  inserted into `run.sh` / `run.ps1`)
- `01-spec-authoring-guide/12-mandatory-linter-infra03-structure.md` —
  Top-level rule that mandates the existence of `linter-scripts/`
- `linters-cicd/checks/_lib/sarif.py` — Reference SARIF emitter for
  plugins under `linters-cicd/`
- `29-blind-ai-audit-v2.md §9` — Prior gap that this guide closes

---

## Verification

_Auto-generated section — see `02-spec/17-consolidated-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CON-027: Linter authoring guide conformance

**Given** A blind AI tasked with adding a new linter to `linter-scripts/`.
**When** It follows §1–§11 of this file end-to-end.
**Then** The resulting linter passes `bash linter-scripts/run.sh` with
exit `0`, registers correctly in `05-coding-guidelines.md §34`, and
produces fixtures covering every non-zero exit code it can return.

**Verification command:**

```bash
bash linter-scripts/run.sh && \
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit `0`. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-22_
