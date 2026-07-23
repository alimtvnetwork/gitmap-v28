#!/usr/bin/env bash
# pr-summary.sh — Build a concise PR comment summarizing full-suite-guard.
#
# Inputs:
#   $1  artifact_dir  — directory containing test-output.txt + lint-output.txt
#                       (downloaded full-suite-outputs artifact)
#   $2  out_path      — file to write the rendered Markdown comment to
#
# Required env:
#   FULL_SUITE_RESULT   — "success" | "failure" | "cancelled" | "skipped"
#   GITHUB_REPO         — "owner/repo"
#   GITHUB_RUN_ID       — Actions run ID for log deep links
#   GITHUB_SHA_SHORT    — full or short SHA for display
#
# Sentinel: <!-- gitmap-ci-summary --> is embedded so the sticky-comment
# action can find and replace this comment on subsequent pushes.

set -uo pipefail

readonly ARTIFACT_DIR="${1:?artifact_dir required}"
readonly OUT_PATH="${2:?out_path required}"

readonly TEST_OUT="$ARTIFACT_DIR/test-output.txt"
readonly LINT_OUT="$ARTIFACT_DIR/lint-output.txt"

readonly REPO="${GITHUB_REPO:-unknown/unknown}"
readonly RUN_ID="${GITHUB_RUN_ID:-0}"
readonly SHA="${GITHUB_SHA_SHORT:-HEAD}"
readonly RESULT="${FULL_SUITE_RESULT:-unknown}"

readonly RUN_URL="https://github.com/${REPO}/actions/runs/${RUN_ID}"
readonly ARTIFACTS_URL="${RUN_URL}#artifacts"

# Tally counts from the captured output files. Defensive defaults so a
# missing artifact (e.g. job cancelled before upload) still renders cleanly.
hasTestOutput=true
if [ -f "$TEST_OUT" ]; then
  hasTestOutput=true
else
  hasTestOutput=false
fi

hasLintOutput=true
if [ -f "$LINT_OUT" ]; then
  hasLintOutput=true
else
  hasLintOutput=false
fi

testsPassed=0
testsFailed=0
if [ "$hasTestOutput" = "true" ]; then
  testsPassed=$(grep -cE '^ok[[:space:]]'   "$TEST_OUT" || true)
  testsFailed=$(grep -cE '^FAIL[[:space:]]' "$TEST_OUT" || true)
fi

lintIssues=0
if [ "$hasLintOutput" = "true" ]; then
  lintIssues=$(grep -cE '^[^[:space:]].+:[0-9]+:[0-9]+:' "$LINT_OUT" || true)
fi

# Per-status emoji + headline so the comment is glanceable in the PR feed.
testsStatus="✅ pass"
if [ "$testsFailed" -gt 0 ]; then
  testsStatus="❌ fail"
fi
if [ "$hasTestOutput" = "false" ]; then
  testsStatus="⚠ no output"
fi

lintStatus="✅ clean"
if [ "$lintIssues" -gt 0 ]; then
  lintStatus="❌ $lintIssues issue(s)"
fi
if [ "$hasLintOutput" = "false" ]; then
  lintStatus="⚠ no output"
fi

overallEmoji="✅"
if [ "$RESULT" != "success" ]; then
  overallEmoji="❌"
fi

# First failing tests / first lint issues, capped to keep the comment short.
firstFailures=""
if [ "$testsFailed" -gt 0 ]; then
  firstFailures=$(grep -E '^--- FAIL:|^FAIL[[:space:]]' "$TEST_OUT" | head -n 10 || true)
fi

firstLintLines=""
if [ "$lintIssues" -gt 0 ]; then
  firstLintLines=$(grep -E '^[^[:space:]].+:[0-9]+:[0-9]+:' "$LINT_OUT" | head -n 10 || true)
fi

# Gate dashboard — turn each `needs.<job>.result` env var into a glanceable
# row (✅ pass / ❌ fail / ⏭ skipped / ♻ deduped / ⚠ unknown). We render a
# fixed list of gates so the table shape stays stable across pushes even
# when GitHub skips a job (e.g. SHA-dedup short-circuit). The env names
# are GATE_<UPPERCASE_NAME>; values come straight from `needs.*.result`.
gateRow() {
  local label="$1" envName="$2" detail="$3"
  local raw="${!envName:-}"
  local icon
  case "$raw" in
    success)   icon="✅ pass" ;;
    failure)   icon="❌ fail" ;;
    cancelled) icon="🛑 cancelled" ;;
    skipped)   icon="⏭ skipped" ;;
    "")        icon="⏭ not run" ;;
    *)         icon="⚠ ${raw}" ;;
  esac
  # SHA-dedup passthrough: a "success" with the dedup flag set means the
  # job intentionally short-circuited. Surface that as ♻ deduped so a
  # green row isn't mistaken for a real run.
  if [ "$raw" = "success" ] && [ "${SHA_DEDUPED:-false}" = "true" ]; then
    icon="♻ deduped"
  fi
  echo "| ${label} | ${icon} | ${detail} |"
}

{
  echo "<!-- gitmap-ci-summary -->"
  echo "## ${overallEmoji} CI Summary — \`${SHA:0:10}\`"
  echo ""
  echo "Full Suite Guard: **${RESULT}**"
  echo ""
  echo "### Gates"
  echo ""
  echo "| Gate | Status | Detail |"
  echo "| --- | --- | --- |"
  gateRow "Spell check (misspell)"        "GATE_SPELL_CHECK"        "US locale, whole repo"
  gateRow "Lint (golangci-lint + vet + gofmt + goimports)" "GATE_LINT" "v1.64.8 strict, --issues-exit-code=1"
  gateRow "Lint script unit tests"        "GATE_LINT_SCRIPT_TESTS"  "bash + jq"
  gateRow "Lint baseline diff"            "GATE_LINT_BASELINE_DIFF" "soft-gate, new-issues only"
  gateRow "Repo-policy checks"            "GATE_REPO_POLICY"        "naming, legacy-refs, generate-drift, layout"
  gateRow "Vulnerability scan (govulncheck)" "GATE_VULNCHECK"       "v1.1.4, third-party = fail"
  gateRow "Tests (go test ./... -count=1)" "GATE_TEST"              "${testsPassed} pkg(s) ok, ${testsFailed} failed"
  gateRow "JSON snapshot smoke"           "GATE_JSON_SNAPSHOT_FAST" "fast subset"
  gateRow "Installer smoke (Linux/macOS)" "GATE_INSTALLER_SMOKE"    "install.sh contract"
  gateRow "Installer smoke (Windows)"     "GATE_INSTALLER_SMOKE_WINDOWS" "install.ps1 contract"
  gateRow "Full-suite guard"              "GATE_FULL_SUITE_GUARD"   "test + lint replay, SARIF upload"
  gateRow "Cross-compile build"           "GATE_BUILD"              "6 GOOS/GOARCH targets"
  echo ""
  echo "### Detail"
  echo ""
  echo "| Stage | Result | Detail |"
  echo "| --- | --- | --- |"
  echo "| \`go test ./...\` | ${testsStatus} | ${testsPassed} package(s) ok, ${testsFailed} failed |"
  echo "| \`golangci-lint\` (strict) | ${lintStatus} | v1.64.8, --max-issues=0 |"
  echo ""

  if [ -n "$firstFailures" ]; then
    echo "### First failing tests"
    echo ""
    echo '```'
    echo "$firstFailures"
    echo '```'
    echo ""
  fi

  if [ -n "$firstLintLines" ]; then
    echo "### First lint findings"
    echo ""
    echo '```'
    echo "$firstLintLines"
    echo '```'
    echo ""
  fi

  echo "### Logs & artifacts"
  echo ""
  echo "- [Full run logs](${RUN_URL})"
  echo "- [Download artifacts](${ARTIFACTS_URL}) — \`full-suite-outputs\` contains the raw \`test-output.txt\` and \`lint-output.txt\`"
  echo ""
  echo "<sub>Reproduce locally: \`./scripts/preflight-ci.sh\`</sub>"
} > "$OUT_PATH"

echo "✓ pr-summary: wrote $OUT_PATH"
