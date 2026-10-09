# 265 — CI Gates: 300-Line File Rule + Test-Matrix Verification + Dormant Workflow Decision

Status: spec (draft) · Date: 2026-10-10 · Owner task: `265-stdout-ci-gates`

## Owner request / scope

Close the REAL gaps in CI for task 265 ("stdout" hygiene + CI gates). Narrow scope:

1. Add a 300-line file-size gate job to CI.
2. Verify and document the existing test matrix as the WS3 deliverable (no new workflow).
3. Make a decision recommendation on the two dormant workflows — flag as an owner decision, do NOT delete anything in this spec.

No new lint behavior. No golangci-lint config changes.

## Current state (verified by research — do not re-derive)

`.github/workflows/` has 15 files. `ci.yml` ALREADY runs:

- gofmt self-heal
- `go vet ./...`
- compile gate (`go test -run=^$ ./...`)
- golangci-lint (strict)
- test matrix: unit / store / integration / tui

The REAL gaps:

- **(a)** The 300-line file rule is NOT a CI gate. The only automated check is a 200-line shell shim in `race-detector.yml`: `bash .github/scripts/file-size-check.sh 200`.
- **(b)** `coverage-floor.yml` and `mutation-tests.yml` are DORMANT: stale `gitmap/**` path filters (the tree restructured to `cli/**`) and go 1.24 pins while the module requires go 1.26.0.

Precedent for running `03-ai-scripts/*.py` from CI: `purge-actions-artifacts.yml` runs `python3 03-ai-scripts/34-...py` with `setup-python`. Follow that precedent exactly.

The Python gate tool already exists: `03-ai-scripts/52-file-size-check.py` supports `--staged`, `--diff`, `--files`, `--limit` modes; default limit 300.

## Deliverable 1 — 300-line file-size gate job in CI

Add a `file-size-gate` job to `ci.yml` (keep it in `ci.yml` unless the lead finds it unwieldy, in which case a dedicated `file-size-gate.yml` is acceptable):

```yaml
file-size-gate:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0   # --diff mode needs history
    - uses: actions/setup-python@v5
      with:
        python-version: "3.12"
    - name: 300-line file-size gate (PR diff)
      run: python3 03-ai-scripts/52-file-size-check.py --diff --limit 300
```

Decisions:

- **D1:** `--diff` mode on `pull_request` events (only files the PR touches). On `push` to main, run the same job against the pushed commit range. The job fails when any checked file exceeds 300 lines (respecting the script's existing exclusion list).
- **D2:** The limit is 300 — do NOT pass 200. The 200-line shim in `race-detector.yml` is a separate, stricter local rule, not the CI gate. Aligning the shim is a follow-up, not this spec.
- **D3:** Mirror the `purge-actions-artifacts.yml` Python setup (same action versions) so CI stays consistent.

Acceptance: a PR adding a 301+ line file fails the gate; the job log names the violating file(s).

## Deliverable 2 — Test matrix as the WS3 deliverable (verification only)

No new workflow. The implementer VERIFIES on the branch:

1. The `ci.yml` test-matrix jobs (unit / store / integration / tui) trigger and pass on the PR.
2. The `go vet` and compile gates pass.
3. golangci-lint (strict) passes.

Record the verification evidence (run URLs + commit SHA) in the subtask file. If any matrix job is red for pre-existing reasons, list it explicitly as a follow-up — do not silently fix unrelated failures here.

## Deliverable 3 — Decision on dormant workflows (owner decision, NOT deletion)

`coverage-floor.yml` and `mutation-tests.yml` are dormant because their path filters reference the stale `gitmap/**` tree and they pin go 1.24 while the module requires go 1.26.0.

Recommendation for the owner:

- **Option A (recommended):** FIX the path filters (`cli/**`, `03-ai-scripts/**`, `.github/workflows/**`) and bump the Go pin to 1.26, then re-enable coverage-floor at the currently measured coverage.
- **Option B:** DELETE both workflows if coverage/mutation gating is no longer a goal.

This spec does NOT delete them. The implementer prepares the Option A fix as a ready-to-review diff; the owner decides whether it merges or the workflows are removed.

## Non-goals

- No new linters, no golangci-lint config changes.
- No changes to the 200-line shell shim in `race-detector.yml`.
- No deletion of any workflow file.

## Acceptance criteria

1. `file-size-gate` job present in CI, green on main, fails on a 301-line file in a test PR.
2. Test-matrix verification evidence recorded (run URLs + SHA).
3. Dormant-workflow decision documented with the owner's call; Option A diff prepared if the owner picks it.
