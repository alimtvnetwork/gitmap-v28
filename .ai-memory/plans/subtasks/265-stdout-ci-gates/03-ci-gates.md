# Subtask 265-stdout-ci-gates · 03 — CI Gates

Scope: implementer of the CI changes for task 265.
Spec home: `02-spec/21-app/265-stdout-ci-gates/03-ci-gates.md` (read it first — it holds decisions D1–D3).
Rules: lowercase filenames, relative paths only. YAML edits keep the existing action versions/style in `ci.yml`. Do NOT run git commands — the lead owns commits/releases. Do NOT delete any workflow file.

## Steps

- [ ] Add the `file-size-gate` job to `.github/workflows/ci.yml`:
  - `actions/checkout@v4` with `fetch-depth: 0` (the `--diff` mode needs history)
  - `actions/setup-python@v5` with `python-version: "3.12"` (mirror the `purge-actions-artifacts.yml` precedent)
  - run step: `python3 03-ai-scripts/52-file-size-check.py --diff --limit 300`
- [ ] Sanity-check the script flags first via `gitmap py 03-ai-scripts/52-file-size-check.py --help` (Python only via `gitmap py`; never bare `python3` for repo scripts) and confirm `--diff --limit 300` passes on the current tree.
- [ ] Verify the test matrix as the WS3 deliverable: push the branch / open the PR and confirm the `ci.yml` jobs — gofmt self-heal, `go vet`, compile gate, golangci-lint (strict), and the test matrix (unit/store/integration/tui) — all trigger and pass; record run URLs + commit SHA in Evidence below.
- [ ] Negative test: on a scratch branch add a 301+ line file, confirm the `file-size-gate` job fails and names the file, then remove the scratch file.
- [ ] Dormant workflows: prepare the Option A fix diff for `coverage-floor.yml` + `mutation-tests.yml` (path filters → `cli/**`, `03-ai-scripts/**`, `.github/workflows/**`; Go pin 1.24 → 1.26) but DO NOT merge or delete — hand the diff to the owner for the decision call.
- [ ] If any matrix job is red for pre-existing reasons, list it explicitly as a follow-up — do not silently fix unrelated failures here.

## Evidence

(gate job run URL, matrix run URLs, commit SHA, Option A diff location, owner decision on dormant workflows)
