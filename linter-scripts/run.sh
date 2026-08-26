#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────────────
# run.sh — Pull latest changes, run guidelines validator, run all linters.
#
# Usage:
#   ./linter-scripts/run.sh                       # full pipeline
#   ./linter-scripts/run.sh -d                    # git pull only
#   ./linter-scripts/run.sh --path cmd --max-lines 20
#   ./linter-scripts/run.sh --json
#   ./linter-scripts/run.sh --skip-linters        # skip Step 3
#   ./linter-scripts/run.sh --linters-only        # only Step 3
# ──────────────────────────────────────────────────────────────────────

set -uo pipefail

SCAN_PATH="src"
MAX_LINES=15
JSON_FLAG=""
SKIP_VALIDATION=false
SKIP_LINTERS=false
LINTERS_ONLY=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    -d)              SKIP_VALIDATION=true; shift ;;
    --path)          SCAN_PATH="$2"; shift 2 ;;
    --max-lines)     MAX_LINES="$2"; shift 2 ;;
    --json)          JSON_FLAG="--json"; shift ;;
    --skip-linters)  SKIP_LINTERS=true; shift ;;
    --linters-only)  LINTERS_ONLY=true; SKIP_VALIDATION=true; shift ;;
    -h|--help)
      sed -n '2,15p' "$0"
      exit 0
      ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_FILE="$SCRIPT_DIR/validate-guidelines.go"

# ── Step 1: Git Pull ──────────────────────────────────────────────
if [ "$LINTERS_ONLY" = false ]; then
  echo ""
  echo "═══ Step 1 — git pull ═══"
  if git pull; then
    echo "✅ Repository up to date."
  else
    echo "⚠️  git pull failed — continuing with local files..."
  fi
fi

if [ "$SKIP_VALIDATION" = true ] && [ "$LINTERS_ONLY" = false ]; then
  echo ""
  echo "⏭️  Skipping validation (-d flag)."
  exit 0
fi

VALIDATOR_EXIT=0

# ── Step 2: Go Validator ──────────────────────────────────────────
if [ "$LINTERS_ONLY" = false ]; then
  echo ""
  echo "═══ Step 2 — Coding-guidelines validator ═══"
  if [ ! -f "$GO_FILE" ]; then
    echo "❌ Cannot find $GO_FILE"
    exit 1
  fi
  if ! command -v go &>/dev/null; then
    echo "❌ Go is not installed or not in PATH (https://go.dev/dl/)."
    exit 1
  fi
  echo "Using $(go version)"
  echo "Scanning: $SCAN_PATH (max $MAX_LINES lines/function)"
  echo ""
  if go run "$GO_FILE" --path "$SCAN_PATH" --max-lines "$MAX_LINES" $JSON_FLAG; then
    echo "✅ Step 2 passed."
  else
    VALIDATOR_EXIT=$?
    echo "❌ Step 2 failed with CODE RED violations (exit=$VALIDATOR_EXIT)."
  fi
fi

# ── Step 3: Spec / docs linters ───────────────────────────────────
LINTERS_EXIT=0
PASSED=()
FAILED=()

run_linter() {
  local label="$1"; shift
  echo ""
  echo "── $label ──"
  if "$@"; then
    PASSED+=("$label")
  else
    local rc=$?
    FAILED+=("$label (exit=$rc)")
    LINTERS_EXIT=1
  fi
}

if [ "$SKIP_LINTERS" = false ]; then
  echo ""
  echo "═══ Step 3 — Spec / docs linters ═══"

  run_linter "tunable-constants"     python3 "$SCRIPT_DIR/check-tunable-constants.py"
  run_linter "mws-error-codes"       python3 "$SCRIPT_DIR/check-mws-error-codes.py"
  run_linter "function-lengths"      python3 "$SCRIPT_DIR/check-function-lengths.py"
  run_linter "forbidden-strings"     python3 "$SCRIPT_DIR/check-forbidden-strings.py"
  run_linter "placeholder-comments"  python3 "$SCRIPT_DIR/check-placeholder-comments.py"
  run_linter "memory-mirror-drift"   python3 "$SCRIPT_DIR/check-memory-mirror-drift.py"
  run_linter "prompts-loaded"        python3 "$SCRIPT_DIR/check-prompts-loaded.py"
  run_linter "readme-canonicals"     python3 "$SCRIPT_DIR/check-readme-canonicals.py"
  run_linter "readme-install"        python3 "$SCRIPT_DIR/check-readme-install-section.py"
  run_linter "root-readme"           python3 "$SCRIPT_DIR/check-root-readme.py"
  run_linter "spec-cross-links"      python3 "$SCRIPT_DIR/check-spec-cross-links.py"
  run_linter "spec-folder-refs"      python3 "$SCRIPT_DIR/check-spec-folder-refs.py"
  run_linter "axios-version"         bash    "$SCRIPT_DIR/check-axios-version.sh"
  run_linter "forbidden-spec-paths"  bash    "$SCRIPT_DIR/check-forbidden-spec-paths.sh"
  run_linter "runner-dispatch"       bash    "$SCRIPT_DIR/check-runner-dispatch-antipatterns.sh"
fi

# ── Summary ───────────────────────────────────────────────────────
echo ""
echo "═══ Summary ═══"
echo "Step 2 (validator): exit=$VALIDATOR_EXIT"
if [ "$SKIP_LINTERS" = false ]; then
  echo "Step 3 (linters): ${#PASSED[@]} passed, ${#FAILED[@]} failed"
  for item in "${FAILED[@]}"; do
    echo "  ❌ $item"
  done
fi

if [ "$VALIDATOR_EXIT" -ne 0 ] || [ "$LINTERS_EXIT" -ne 0 ]; then
  exit 1
fi
echo "✅ All checks passed."
exit 0
