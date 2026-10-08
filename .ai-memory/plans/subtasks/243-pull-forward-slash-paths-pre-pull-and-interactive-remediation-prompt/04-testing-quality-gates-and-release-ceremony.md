# Subtask 04: Testing, Quality Gates & Release Ceremony

- **Target Files:**
  - `cli/constants/constants.go`
  - `.agents/skills/gitmap/SKILL.md`
  - `.ai-memory/release/release-notes-v6.516.0.md`

- **Checklist:**
  - [x] Run `python linter-scripts/check-nested-ifs.py --all` -> 0 violations.
  - [x] Run `python linter-scripts/check-enum-and-boolean.py` -> 0 violations.
  - [x] Run `python .github/scripts/check-legacy-refs.py` -> passed.
  - [x] Run `python .github/scripts/go-format-check.py --check-only` -> gofmt clean.
  - [x] Run `golangci-lint run ./cmd/... ./cmdpull/...` from `cli` -> 0 errors.
  - [x] Bump version to `v6.516.0`.
  - [x] Commit with atomic hyphen format: `pull - normalize forward slash paths, enforce pull before changes, and prompt fix all single skip`.
  - [x] Push commit and tag `v6.516.0`.
  - [x] Monitor CI/CD with `gitmap pe -t` until 100% green.
