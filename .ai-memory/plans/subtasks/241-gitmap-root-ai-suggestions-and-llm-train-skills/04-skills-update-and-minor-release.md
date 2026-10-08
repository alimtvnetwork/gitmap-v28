# Subtask Plan 04: Skills Modernization, Verification & Minor Version Bump

## Overview
Synchronize `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md` with the new LLM train skill output, public doc URLs, and recursive directives. Verify all coding guidelines and linters across modified files. Perform a SemVer minor bump from `v6.512.0` to `v6.513.0` across `version.json` and `changelog.md`. Recompile `gitmap.exe` binary and verify bare output and `gitmap llm train`.

---

## Owned Files
- `.agents/skills/gitmap/SKILL.md`
- `.cursor/skills/gitmap/skill.md`
- `version.json`
- `changelog.md`

---

## Step-by-Step Implementation

1. **Update Skills Files:**
   - Update `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md`:
     * Document bare `gitmap` concise output and AI suggestions.
     * Document `gitmap llm train` emitting the full skill to stdout and disk.
     * Document `--urls` flag and public GitHub Markdown URLs.
     * Document recursive Git network directives and self-skill creation mandate.

2. **Run Linter Quality Gates:**
   - `python linter-scripts/check-relative-paths.py` (0 violations)
   - `python .github/scripts/go-format-check.py --check-only` (pass)
   - `python linter-scripts/check-nested-ifs.py` (0 violations)
   - `python linter-scripts/check-enum-and-boolean.py` (0 violations)
   - `python linter-scripts/check-boolean-guidelines.py` (0 violations)
   - `python linter-scripts/check-error-management.py` (0 violations)

3. **Rebuild Binary & Test Live:**
   - Run tests: `go test -v ./cli/cmd/...`
   - Recompile: `go build -p 4 -o C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe .` in `cli/`
   - Verify `gitmap` bare output line count <= 15.
   - Verify `gitmap llm train` outputs skill.

4. **Minor Version Bump:**
   - Bump `version.json`: `6.512.0` -> `6.513.0`
   - Update `changelog.md` with release notes for `v6.513.0`.

---

## Acceptance Criteria
- [ ] Skills synchronized in both `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md`.
- [ ] All linters and unit tests pass with zero violations.
- [ ] `gitmap` bare output is under 15 lines.
- [ ] Version bumped to `6.513.0` in `version.json` and `changelog.md`.
- [ ] Local binary compiled and verified.
