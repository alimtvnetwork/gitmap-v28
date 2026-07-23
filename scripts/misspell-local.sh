#!/usr/bin/env bash
# misspell-local.sh — run the same misspell check as the CI `spell-check`
# job, locally, against the same file filters loaded from
# gitmap/data/config.json (misspell.exclude / misspell.include).
#
# Usage:
#   scripts/misspell-local.sh                 # diff vs origin/main
#   scripts/misspell-local.sh --base HEAD~1   # diff vs an arbitrary ref
#   scripts/misspell-local.sh --all           # scan every tracked file
#   scripts/misspell-local.sh --staged        # scan staged (index) files
#   scripts/misspell-local.sh --files a.md b.go   # scan an explicit list
#
# Exit codes:
#   0 = clean, 1 = misspellings found, 2 = bad usage / missing tool.
#
# Mirrors .github/workflows/ci.yml `spell-check` job 1:1: same pinned
# misspell version (v0.3.4), same -locale US, same exclude/include
# resolution (jq → python3 → hard-coded fallback), same path-glob
# semantics (shell `case` patterns).
set -euo pipefail

MisspellVersion="v0.3.4"
ConfigPath="gitmap/data/config.json"
Mode="diff"
BaseRef="origin/main"
ExplicitFiles=()

die() { printf '%s\n' "$*" >&2; exit 2; }

while [ "$#" -gt 0 ]; do
  case "$1" in
    --all)     Mode="all"; shift ;;
    --staged)  Mode="staged"; shift ;;
    --base)    [ "$#" -ge 2 ] || die "--base requires a ref"; BaseRef="$2"; Mode="diff"; shift 2 ;;
    --files)   shift; while [ "$#" -gt 0 ] && [[ "$1" != --* ]]; do ExplicitFiles+=("$1"); shift; done; Mode="files" ;;
    -h|--help) sed -n '2,18p' "$0"; exit 0 ;;
    *)         die "unknown arg: $1 (try --help)" ;;
  esac
done

if ! command -v misspell >/dev/null 2>&1; then
  echo "misspell not found — installing pinned ${MisspellVersion}..." >&2
  command -v go >/dev/null 2>&1 || die "go toolchain required to install misspell"
  GOBIN="$(go env GOPATH)/bin" go install "github.com/client9/misspell/cmd/misspell@${MisspellVersion}" \
    || die "failed to install misspell ${MisspellVersion}"
  export PATH="$(go env GOPATH)/bin:${PATH}"
fi

# ---- load include/exclude from config (jq → python3 → fallback) ----------
read_patterns() {
  local key="$1"
  if command -v jq >/dev/null 2>&1 && [ -f "$ConfigPath" ]; then
    jq -r --arg k "$key" '.misspell[$k] // [] | .[]' "$ConfigPath" 2>/dev/null || true
  elif command -v python3 >/dev/null 2>&1 && [ -f "$ConfigPath" ]; then
    python3 -c "import json,sys; d=json.load(open('$ConfigPath')).get('misspell',{}).get('$1',[]); print('\n'.join(d))" 2>/dev/null || true
  fi
}
mapfile -t Excludes < <(read_patterns exclude)
mapfile -t Includes < <(read_patterns include)
if [ "${#Excludes[@]}" -eq 0 ]; then
  Excludes=(
    "*.png" "*.jpg" "*.jpeg" "*.gif" "*.ico" "*.svg" "*.webp"
    "*.zip" "*.tar" "*.gz" "*.exe"
    "*/testdata/*" "*/golden/*"
    "*/.gitmap/release/*" "*/.gitmap/release-assets/*"
    "gitmap/completion/allcommands_generated.go"
    ".lovable/*"
  )
fi

# ---- gather candidate file list per mode ---------------------------------
case "$Mode" in
  all)     mapfile -t Candidates < <(git ls-files) ;;
  staged)  mapfile -t Candidates < <(git diff --name-only --diff-filter=AM --cached) ;;
  files)   Candidates=("${ExplicitFiles[@]}") ;;
  diff)
    git rev-parse --verify "$BaseRef" >/dev/null 2>&1 \
      || die "base ref '$BaseRef' not found (fetch it or pass --base)"
    mapfile -t Candidates < <(git diff --name-only --diff-filter=AM "${BaseRef}"...HEAD)
    ;;
esac

if [ "${#Candidates[@]}" -eq 0 ]; then
  echo "no candidate files — nothing to scan"
  exit 0
fi

matches_any() {
  local path="$1"; shift
  local pat
  for pat in "$@"; do
    # shellcheck disable=SC2254
    case "$path" in $pat) return 0 ;; esac
  done
  return 1
}

Targets=()
for f in "${Candidates[@]}"; do
  [ -f "$f" ] || continue
  if matches_any "$f" "${Excludes[@]}"; then continue; fi
  if [ "${#Includes[@]}" -gt 0 ] && ! matches_any "$f" "${Includes[@]}"; then continue; fi
  Targets+=("$f")
done

if [ "${#Targets[@]}" -eq 0 ]; then
  echo "no scannable files after filters — nothing to scan"
  exit 0
fi

echo "misspell-local: mode=${Mode} | excludes=${#Excludes[@]} includes=${#Includes[@]} | scanning ${#Targets[@]} file(s)"
report="$(mktemp)"
misspell -locale US "${Targets[@]}" | tee "$report" || true

findings=0
while IFS= read -r line; do
  [[ "$line" =~ ^[^:]+:[0-9]+:[0-9]+: ]] && findings=$((findings + 1))
done < "$report"

if [ "$findings" -gt 0 ]; then
  echo "FAIL: ${findings} misspelling(s) found." >&2
  exit 1
fi
echo "OK: no misspellings."