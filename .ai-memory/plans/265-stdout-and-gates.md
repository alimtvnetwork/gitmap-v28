# Ledger: 265-stdout-and-gates
Request slug: 265-stdout-and-gates
Request first line: stdout refactor + CI gates + deprecation policy + credential audit + benchmarks + pe-all parallel upgrade
Status: ACTIVE
Phase: 1    Wave: 0 / WAVES    Step: 1 / N
Last completed action: Step 0 (git pull clean, spec 265 via dogfooded `spec next`, backup branch backup/265-stdout-and-gates pushed)
Next action: Synthesize research reports → spec 265 → subtasks → execution waves
Workers in flight: Research 01 (stdout arch + credential audit), Research 02 (CI + pe-all + deprecation + benchmarks + tests)
Commits: none    Pushed: no
Branch: main | Tree at start: clean
Tools: subagent.spawn=yes gitmap=yes
| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | stdout refactor (WS1) | Worker 01 | cli/output/ (new), cli/termout/, cli/glyphs/, cli/cliexit/, cli/cmd/root.go, cli/cmdmacro/macro_add_file_ops.go, cli/cmdnodes/nodes_clone.go | IN PROGRESS (Wave 1) | subtasks 01+02 |
| Task-02 | CI quality gates (WS2) | Worker 02 | .github/workflows/ci.yml | IN PROGRESS (Wave 1) | subtask 03 |
| Task-03 | Tests in CI (WS3) | Worker 02 | (verification-only, folded into WS2) | IN PROGRESS (Wave 1) | existing matrix verified, no new workflow |
| Task-04 | Deprecation policy (WS4) | Worker 02 | policy doc (new) + cli/cmdautofix/fix.go | IN PROGRESS (Wave 1) | subtask 04 |
| Task-05 | Credential audit (WS5) | Research 01 | findings doc only (read-only) | DONE | 3 HIGH + 2 MEDIUM in spec 04-deprecation-and-pe-parallel.md |
| Task-06 | Benchmarks (WS6) | TBD (Wave 2) | *_test.go benchmark files | PENDING | - |
| Task-07 | pe-all parallel upgrade (WS7) | TBD (Wave 2) | cli/cmdpipeline/pipeline_all_errors*.go | PENDING | - |
| Task-08 | FINAL WAVE (owner override) — minor bump + release + CI/CD verification loop | TBD (lead) | version.json, tag, GitHub release, ~/.local/bin/gitmap | PENDING | Runs AFTER all WS1–WS7 committed+pushed: 37-bump-version.py → commit → tag → push tag → gh release → rebuild binary → `pe -t 1200` wait → `pe` check → fix+re-release loop until green |
Assumptions: none yet
Conflicts: V6 skill R1 (zero builds) vs parent-authorized `go build`/`go vet` verification — parent wins per Top-Instruction Priority Mandate; WS2/WS3 are CI work by definition.
Stage list: (pending research synthesis)

## Solo-execution note (lead)
`subagent.spawn` is runtime-blocked ("subagent bootstrap is no longer authorized", confirmed 2x — not transient). No workers can be spawned. Lead executes all waves solo with full transparency; the V6 solo-execution ban is unsatisfiable when the runtime refuses bootstrap, so verification discipline (fresh-binary checks per wave, lead-run builds) compensates. Parent standing rule ("never just stop") applies: work continues.
Wave order (solo): WS2 (CI gates) → WS4 (deprecation) → WS6 (benchmarks) → WS7 (pe-all parallel) → WS1 (stdout refactor, riskiest, last) → Task-08 (release wave).
