# Component Spec: 311-cicd-issues-and-end-to-end-testing

## 1. CI/CD Failure Map (baseline at f723012c5 / 1f96b906)

| Failing check | Suspected root cause | Fix route |
|---|---|---|
| Lint Baseline Guard / Diff | Unused funcs (cli/cmddb/db_vacuum.go x4, cli/cmd/helpcheck.go, cli/cloner/safe_pull.go, formatFolderFileURI, hasSkipIDEsFlag, parseGitHubSlugFromURL); misspell DOUB | Delete dead code after caller verification via gitmap aum search; fix spelling at source; refresh baseline only for entries whose defect was the baseline itself |
| Legacy Refs Check | /gitmap-v[567]\b/ hits in 02-spec/01-spec-authoring-guide/13-root-readme-conventions.md lines 20, 257 | Reword the guide to the current naming without legacy version refs |
| Relative Path Check | Absolute paths / file URIs in tracked files | Coordinate with the in-flight spec-310 stream (its HEAD commit abb35f01d targets exactly this); fix only files outside its dirty set |
| Test (unit), go build + test x3, go test -race | Real failing tests introduced by the Oct 21-23 change stream; annotations are generic, so failing packages must be identified from the CI job structure and reproduced through CI iterations | Fix the underlying code/test defects; no test deletion to reach green |
| Full Suite Guard | Derivative of the unit/build failures | Clears when the suite itself is green |
| Lint Script Unit Tests | A python linter-script test failing | Identify via the script test suite layout (linter-scripts tests) and fix script or test at the source |

Fresh check-runs must be pulled for the current HEAD before fixing, because
the spec-310 stream is actively changing the failure set.

## 2. End-to-End Testing of All Methods

- Source of truth for "all methods": the gitmap command registry
  (cli/constants command registry + helpdoc index). Every registered command
  must be reachable and exercised.
- Harnesses already in CI and in scope: Installer Smoke (source-build) on
  linux + windows, history-rewrite smoke, command invocation tests
  (cli/cmd/cmd_invocation_test.go — note: currently dirty in the spec-310
  stream's box; coordinate before touching).
- Any method found untested gets the smallest e2e case that proves it runs
  end-to-end (help surface + one real invocation against a scratch repo),
  executed in CI, not locally.
- Deliverable evidence: CI run for the release commit with every smoke/e2e
  job green, plus a method-coverage list in this spec's evidence section.

## 3. Installation Download Speed

Current path (to be profiled before changing):

- Installer scripts at repo root (install.sh family) and the self-update path
  in cli/cmdupdate/ download the release tarball + checksums from GitHub
  releases in a single sequential connection, then swap the binary
  (cmdupdate/update_binary_swap.go — also the site of the known E1148 defect:
  the swap looks for a binary literally named `gitmap` while the tarball ships
  `gitmap-v<version>-<os>-<arch>`; the speed work must not regress the fix
  direction already reported for that defect).
- Levers, in priority order (consistent with spec 242 aria2c design and
  spec 286 speedup):
  1. Multi-connection segmented download (aria2c when available, with a
     built-in parallel range-download fallback) for the release asset.
  2. Fetch checksums.txt concurrently with the asset instead of after it.
  3. Resume support so a stalled download does not restart from zero.
  4. Skip-download fast path when the target version is already cached.
- Measurement: time a full install download before and after (same network,
  same asset), record both numbers here as evidence. The Installer Smoke CI
  jobs are the regression gate.

## 4. Release

- Bump: minor increment of version.json and every synced version file
  (current version to be read from version.json at execution time;
  v6.590.0 -> v6.591.0 at spec time).
- Tag: annotated tag for the new version, pushed (tag-based release
  convention on this machine; no GitHub token available locally).
- Release notes: a release-notes file for the new version listing the
  complete change set of the release (CI fixes, e2e coverage, download-speed
  improvement, and the pe false-green fix from db665c939/f723012c5 that this
  release carries). Relative paths only — no absolute paths, no filesystem
  URLs.

## 5. Evidence (to fill during execution)

- Baseline check-runs (release parent):
- Post-fix check-runs (release commit):
- Download timing before:
- Download timing after:
- Method coverage list:
