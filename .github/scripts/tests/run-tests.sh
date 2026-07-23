#!/usr/bin/env bash
# Unit tests for .github/scripts/check-single-linter-diff.sh
#
# Locks in the "new issues only" baseline-diff contract that all five
# guarded linters (unused, gosec G115, misspell, gocritic, exhaustive)
# now share. Each test:
#   1. seeds a fake $CURRENT_OUT JSON report (the script's golangci-lint
#      invocation is stubbed via a PATH shim so it leaves our fixture
#      in place rather than overwriting it),
#   2. optionally seeds a baseline JSON,
#   3. runs check-single-linter-diff.sh with a controlled env,
#   4. asserts on exit code + selected stdout/stderr substrings.
#
# Why a hand-rolled harness instead of bats?
#   - Zero external deps (CI image has bash + jq, nothing else needed).
#   - The contract is small enough that 6 focused cases give full
#     coverage of the diff/seeding/text-filter branches.
#
# Exit codes:
#   0 — all cases passed
#   1 — at least one case failed (per-case diagnostics on stderr)
#   2 — harness setup error (missing tool, can't write tmpdir, etc.)

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
SCRIPT_UNDER_TEST="$REPO_ROOT/.github/scripts/check-single-linter-diff.sh"

if [ ! -x "$SCRIPT_UNDER_TEST" ] && [ ! -f "$SCRIPT_UNDER_TEST" ]; then
  echo "ERROR: script under test not found: $SCRIPT_UNDER_TEST" >&2
  exit 2
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: jq not on PATH (required by script under test)" >&2
  exit 2
fi

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

# ---------------------------------------------------------------------
# PATH shim: a fake `golangci-lint` whose stdout is sourced from the
# fixture file at $TEST_FIXTURE_CURRENT (set per-case). The script
# under test redirects `golangci-lint ... > $CURRENT_OUT`, so whatever
# the shim writes to stdout becomes the "current" report. This keeps
# the contract test-controlled without touching the script under test.
# `cd "$LINT_DIR"` inside the script also needs $LINT_DIR to exist —
# we create $WORK/repo as a no-op stand-in.
# ---------------------------------------------------------------------
SHIM_DIR="$WORK/shim"
mkdir -p "$SHIM_DIR" "$WORK/repo"
cat > "$SHIM_DIR/golangci-lint" <<'SHIM'
#!/usr/bin/env bash
# Test shim — emits the pre-seeded current-report fixture on stdout so
# the script-under-test's `> $CURRENT_OUT` redirect captures it. The
# fixture path is passed via env to keep the shim stateless.
if [ -n "${TEST_FIXTURE_CURRENT:-}" ] && [ -f "$TEST_FIXTURE_CURRENT" ]; then
  cat "$TEST_FIXTURE_CURRENT"
else
  echo '{"Issues": []}'
fi
exit 0
SHIM
chmod +x "$SHIM_DIR/golangci-lint"

PASS=0
FAIL=0

# Build a golangci-lint-shaped JSON report with N issues.
# Each arg is "linter|file|line|col|text".
make_report() {
  local out="$1"; shift
  local issues="[]"
  local entry
  for entry in "$@"; do
    IFS='|' read -r linter file line col text <<<"$entry"
    issues=$(jq --arg l "$linter" --arg f "$file" \
      --argjson ln "$line" --argjson c "$col" --arg t "$text" \
      '. + [{
        FromLinter: $l,
        Text: $t,
        Pos: { Filename: $f, Line: $ln, Column: $c }
      }]' <<<"$issues")
  done
  jq -n --argjson i "$issues" '{Issues: $i}' > "$out"
}

# Run the script under test with a controlled env. Captures stdout +
# stderr + exit code into per-case files for assertion.
run_case() {
  local case_dir="$1"; shift
  mkdir -p "$case_dir"
  set +e
  PATH="$SHIM_DIR:$PATH" \
    bash "$SCRIPT_UNDER_TEST" "$WORK/repo" \
    > "$case_dir/stdout" 2> "$case_dir/stderr"
  echo "$?" > "$case_dir/exit"
  set -e
}

# Assertion helpers. On failure, dump captured output for debugging.
assert_exit() {
  local case_name="$1" case_dir="$2" want="$3"
  local got
  got=$(cat "$case_dir/exit")
  if [ "$got" = "$want" ]; then
    echo "  ✓ $case_name: exit=$got"
    PASS=$((PASS + 1))
  else
    echo "  ✗ $case_name: exit want=$want got=$got" >&2
    echo "    --- stdout ---" >&2; sed 's/^/    /' "$case_dir/stdout" >&2
    echo "    --- stderr ---" >&2; sed 's/^/    /' "$case_dir/stderr" >&2
    FAIL=$((FAIL + 1))
  fi
}
assert_stdout_has() {
  local case_name="$1" case_dir="$2" needle="$3"
  if grep -qF -- "$needle" "$case_dir/stdout"; then
    echo "  ✓ $case_name: stdout contains '$needle'"
    PASS=$((PASS + 1))
  else
    echo "  ✗ $case_name: stdout missing '$needle'" >&2
    sed 's/^/    /' "$case_dir/stdout" >&2
    FAIL=$((FAIL + 1))
  fi
}
assert_stdout_lacks() {
  local case_name="$1" case_dir="$2" needle="$3"
  if grep -qF -- "$needle" "$case_dir/stdout"; then
    echo "  ✗ $case_name: stdout unexpectedly contains '$needle'" >&2
    sed 's/^/    /' "$case_dir/stdout" >&2
    FAIL=$((FAIL + 1))
  else
    echo "  ✓ $case_name: stdout does not contain '$needle'"
    PASS=$((PASS + 1))
  fi
}

echo "============================================================"
echo "  Unit tests for check-single-linter-diff.sh"
echo "============================================================"

# ---------------------------------------------------------------------
# Case 1: empty current report → exit 0, "no new" message.
# Locks in the trivial happy path — nothing to report, nothing to fail.
# ---------------------------------------------------------------------
CASE="$WORK/case01-empty-current"
mkdir -p "$CASE"
make_report "$CASE/fixture.json"
make_report "$CASE/baseline.json"
export LINTER="unused" \
       BASELINE="$CASE/baseline.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="" LABEL=""
run_case "$CASE"
assert_exit "case01 empty-current"   "$CASE" "0"
assert_stdout_has "case01 empty-current" "$CASE" "OK: no new"

# ---------------------------------------------------------------------
# Case 2: current == baseline (one finding in both) → exit 0.
# Locks in: pre-existing findings are TOLERATED. This is the exact
# behavior change vs. the old hard-floor model.
# ---------------------------------------------------------------------
CASE="$WORK/case02-unchanged"
mkdir -p "$CASE"
make_report "$CASE/fixture.json"  "unused|pkg/a.go|10|1|func unusedFn is unused"
make_report "$CASE/baseline.json" "unused|pkg/a.go|10|1|func unusedFn is unused"
export LINTER="unused" \
       BASELINE="$CASE/baseline.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="" LABEL=""
run_case "$CASE"
assert_exit "case02 unchanged"   "$CASE" "0"
assert_stdout_has "case02 unchanged" "$CASE" "+ NEW    : 0"

# ---------------------------------------------------------------------
# Case 3: NEW finding present in current but not baseline → exit 1.
# Locks in: only NEW issues fail. Annotation must include file:line.
# ---------------------------------------------------------------------
CASE="$WORK/case03-new-finding"
mkdir -p "$CASE"
make_report "$CASE/fixture.json" \
  "unused|pkg/old.go|10|1|func oldFn is unused" \
  "unused|pkg/new.go|42|2|func newFn is unused"
make_report "$CASE/baseline.json" \
  "unused|pkg/old.go|10|1|func oldFn is unused"
export LINTER="unused" \
       BASELINE="$CASE/baseline.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="" LABEL=""
run_case "$CASE"
assert_exit "case03 new-finding"   "$CASE" "1"
assert_stdout_has "case03 new-finding" "$CASE" "pkg/new.go"
assert_stdout_has "case03 new-finding" "$CASE" "(NEW vs baseline)"
assert_stdout_lacks "case03 new-finding" "$CASE" "pkg/old.go"

# ---------------------------------------------------------------------
# Case 4: missing baseline → seeding mode (warn, don't fail) even with
# findings present. Locks in the bootstrap behavior on first PR before
# any baseline cache exists on main.
# ---------------------------------------------------------------------
CASE="$WORK/case04-seeding"
mkdir -p "$CASE"
make_report "$CASE/fixture.json" \
  "unused|pkg/x.go|1|1|func x is unused"
export LINTER="unused" \
       BASELINE="$CASE/does-not-exist.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="" LABEL=""
run_case "$CASE"
assert_exit "case04 seeding"   "$CASE" "0"
assert_stdout_has "case04 seeding" "$CASE" "(seeding baseline)"

# ---------------------------------------------------------------------
# Case 5: TEXT_FILTER scopes gosec to G115 only.
# Current has one G115 (NEW) and one G304 (also NEW). Baseline is empty.
# With TEXT_FILTER=G115, only the G115 finding should be reported as new
# — the G304 is filtered out before the diff runs. Locks in the gosec
# G115 sub-step's contract: scope to one rule, fail only on new.
# ---------------------------------------------------------------------
CASE="$WORK/case05-text-filter"
mkdir -p "$CASE"
make_report "$CASE/fixture.json" \
  "gosec|pkg/a.go|10|1|G115: integer overflow conversion int -> uint32" \
  "gosec|pkg/b.go|20|1|G304: file path provided as taint input"
make_report "$CASE/baseline.json"
export LINTER="gosec" \
       BASELINE="$CASE/baseline.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="G115" LABEL="gosec-G115"
run_case "$CASE"
assert_exit "case05 text-filter"   "$CASE" "1"
assert_stdout_has "case05 text-filter" "$CASE" "G115"
assert_stdout_has "case05 text-filter" "$CASE" "GOSEC-G115 DIFF"
assert_stdout_lacks "case05 text-filter" "$CASE" "G304"

# ---------------------------------------------------------------------
# Case 6: TEXT_FILTER applied to BOTH current and baseline.
# Both reports contain a G115 with the same key. Diff must be 0 — the
# filter must not accidentally drop the baseline-side entry, which
# would falsely flag the current G115 as "new". Locks in the
# apples-to-apples symmetry promised in the script header.
# ---------------------------------------------------------------------
CASE="$WORK/case06-filter-symmetry"
mkdir -p "$CASE"
make_report "$CASE/fixture.json" \
  "gosec|pkg/a.go|10|1|G115: integer overflow conversion int -> uint32" \
  "gosec|pkg/b.go|20|1|G304: file path provided as taint input"
make_report "$CASE/baseline.json" \
  "gosec|pkg/a.go|10|1|G115: integer overflow conversion int -> uint32"
export LINTER="gosec" \
       BASELINE="$CASE/baseline.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="G115" LABEL="gosec-G115"
run_case "$CASE"
assert_exit "case06 filter-symmetry"   "$CASE" "0"
assert_stdout_has "case06 filter-symmetry" "$CASE" "+ NEW    : 0"

# ---------------------------------------------------------------------
# Case 7: findings without Pos.Filename are dropped. A path-less issue
# (e.g. typecheck `r.Mode undefined`) in current must NOT count as a
# new finding. Locks in the "full-path only" guarantee that prevents
# stale annotations from masking real failures in the PR UI.
# ---------------------------------------------------------------------
CASE="$WORK/case07-pathless-dropped"
mkdir -p "$CASE"
# Hand-craft a report with one path-less entry (Filename: "").
jq -n '{Issues: [
  {FromLinter: "unused", Text: "phantom", Pos: {Filename: "", Line: 0, Column: 0}}
]}' > "$CASE/fixture.json"
make_report "$CASE/baseline.json"
export LINTER="unused" \
       BASELINE="$CASE/baseline.json" \
       CURRENT_OUT="$CASE/out.json" \
       TEST_FIXTURE_CURRENT="$CASE/fixture.json" \
       TEXT_FILTER="" LABEL=""
run_case "$CASE"
assert_exit "case07 pathless-dropped"   "$CASE" "0"
assert_stdout_has "case07 pathless-dropped" "$CASE" "+ NEW    : 0"

echo "============================================================"
echo "  Results: $PASS passed, $FAIL failed"
echo "============================================================"

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
exit 0
