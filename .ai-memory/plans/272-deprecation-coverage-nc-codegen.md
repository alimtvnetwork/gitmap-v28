# Program 272 — Deprecation + Coverage Gate + nc Codegen

## Status
- [x] Step 0: Sync, spec 272, backup branch `backup/272-improvements` pushed
- [ ] WS1: Deprecation — first real candidate (or remove policy doc)
- [ ] WS2: Coverage floor gate (CI job, >2% drop vs baseline fails)
- [ ] WS3: nc codegen (go:generate → nc_catalog_gen.go)
- [ ] Final: Build verify, commit, push

## Scope authority
.ai-memory/plans/272-decision-record-improvement-scope.md — items 1-3 approved ONLY.
Items 4-6 (pipeline DB retention, staging repo philosophy, error messages) NOT approved — do not touch.

## Waves
### WS1 — Deprecation candidate
- [ ] Search last 20 versions for renamed/repurposed/superseded commands
- [ ] If found: wire markDeprecated with warning + replacement + removal version
- [ ] If none: DELETE docs/deprecation-policy.md + hook (authorized)

### WS2 — Coverage floor
- [ ] New CI job: fail if coverage drops >2% vs main baseline
- [ ] Reuse coverprofile artifacts; baseline via CI cache; seed-and-pass if missing
- [ ] Validate YAML

### WS3 — nc codegen
- [ ] go:generate in cli/cmdpending/ parsing dispatcher registration
- [ ] Emits nc_catalog_gen.go (// Code generated DO NOT EDIT)
- [ ] Must include spec next, te all, all 105 entries; nc --json works
- [ ] All files ≤300 lines; NO runtime AST parsing
