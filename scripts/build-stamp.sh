#!/usr/bin/env bash
# build-stamp.sh — pre-build provenance stamp for stale-checkout detection.
#
# Why this exists
# ---------------
# CI users have hit `cmd/updaterepo.go:118:6: fileExists redeclared in this
# block` errors that only reproduce on stale checkouts predating the
# v3.92.0 rename + v3.113.0 fsutil migration. Without a build-time
# provenance line, the only signal is a cryptic Go error pointing at line
# numbers that don't match the current source. This script prints the
# exact commit SHA, branch, declared `constants.Version`, and a tiny
# fingerprint of the two files that historically caused the redeclaration
# — so a stale checkout is obvious in the very first lines of the build
# log, before `go build` runs.
#
# Output is a single fenced block on stdout. Failures are non-fatal: every
# probe falls back to "(unknown)" so the script never blocks a build that
# could otherwise succeed (e.g. shallow clone without git history,
# offline build from a tarball). The build-stamp itself is informational
# — it fails loud only when the user explicitly compares the printed SHA
# against what they expected.
#
# Usage
# -----
#   bash scripts/build-stamp.sh          # prints to stdout
#   bash scripts/build-stamp.sh --strict # exit 1 if git is missing
#
# Called from
# -----------
#   .github/workflows/ci.yml — pre-`go build` step in the build job
#   run.ps1 / run.sh         — pre-`go build` step in local builds
set -u

readonly STAMP_SCRIPT_VERSION="1.0.0"
readonly REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
readonly CONSTANTS_FILE="${REPO_ROOT}/gitmap/constants/constants.go"
readonly UPDATEREPO_FILE="${REPO_ROOT}/gitmap/cmd/updaterepo.go"
readonly UPDATEDEBUG_FILE="${REPO_ROOT}/gitmap/cmd/updatedebugwindows.go"

strict_mode="false"
if [ "${1:-}" = "--strict" ]; then
  strict_mode="true"
fi

# probe_git runs a git command and echoes its output, or "(unknown)" if
# git is unavailable or the command fails. Strict mode escalates the
# git-missing case to a hard error because CI explicitly opts in.
probe_git() {
  if ! command -v git >/dev/null 2>&1; then
    if [ "$strict_mode" = "true" ]; then
      echo "build-stamp: git not found in PATH (strict mode)" >&2
      exit 1
    fi
    echo "(unknown — git not in PATH)"
    return
  fi

  local out
  out="$(git -C "$REPO_ROOT" "$@" 2>/dev/null)" || out=""
  if [ -z "$out" ]; then
    echo "(unknown)"
    return
  fi
  echo "$out"
}

# probe_constants_version greps the declared version from constants.go
# without invoking Go. The grep is anchored so a renamed/moved constant
# (e.g. someone introducing `const VersionPrev`) doesn't false-match.
probe_constants_version() {
  if [ ! -f "$CONSTANTS_FILE" ]; then
    echo "(unknown — constants.go missing)"
    return
  fi
  grep -E '^const Version = ' "$CONSTANTS_FILE" \
    | head -1 \
    | sed -E 's/.*"([^"]+)".*/\1/' \
    || echo "(unknown — pattern miss)"
}

# fingerprint_file prints "<sha256-prefix> <line-count> <path>" for one
# file. The 12-char SHA prefix is enough to detect stale content while
# staying readable. Falls back to "(missing)" if the file isn't there —
# which is itself a useful stale-checkout signal (the file was renamed
# or removed in a later commit).
fingerprint_file() {
  local label="$1" path="$2"
  if [ ! -f "$path" ]; then
    printf '  %-22s (missing — %s)\n' "$label" "${path#$REPO_ROOT/}"
    return
  fi
  local sha lines
  if command -v sha256sum >/dev/null 2>&1; then
    sha="$(sha256sum "$path" | cut -c1-12)"
  elif command -v shasum >/dev/null 2>&1; then
    sha="$(shasum -a 256 "$path" | cut -c1-12)"
  else
    sha="(no-sha-tool)"
  fi
  lines="$(wc -l <"$path" | tr -d ' ')"
  printf '  %-22s sha256:%s  lines:%s  %s\n' \
    "$label" "$sha" "$lines" "${path#$REPO_ROOT/}"
}

# detect_redecl_risk grep-scans the two historically-problematic files
# for local `func fileExists` / `func fileExistsLoose` declarations. If
# either is found in BOTH files, the build is going to fail with the
# redeclaration error — we surface that prediction up front so the user
# doesn't have to wait for `go build` to confirm it.
detect_redecl_risk() {
  if [ ! -f "$UPDATEREPO_FILE" ] || [ ! -f "$UPDATEDEBUG_FILE" ]; then
    echo "  redecl-risk-check       skipped (one or both source files missing)"
    return
  fi
  local repo_has debug_has
  repo_has="$(grep -cE '^func (fileExists|fileExistsLoose)\(' "$UPDATEREPO_FILE" 2>/dev/null | head -1)"
  debug_has="$(grep -cE '^func (fileExists|fileExistsLoose)\(' "$UPDATEDEBUG_FILE" 2>/dev/null | head -1)"
  repo_has="${repo_has:-0}"
  debug_has="${debug_has:-0}"

  if [ "$repo_has" -gt 0 ] && [ "$debug_has" -gt 0 ]; then
    echo "  redecl-risk-check       ⚠ FAIL — fileExists/fileExistsLoose declared in both files"
    echo "                           (this checkout predates the v3.113.0 fsutil migration)"
    echo "                           expected fix: git pull origin main"
    if [ "$strict_mode" = "true" ]; then
      echo "build-stamp: redeclaration risk detected in strict mode" >&2
      exit 1
    fi
  else
    echo "  redecl-risk-check       ok (no local fileExists* in cmd/ — fsutil migration present)"
  fi
}

cat <<EOF
=== gitmap build-stamp v${STAMP_SCRIPT_VERSION} ====================================
Provenance for stale-checkout detection. Compare these against the SHA
and version you expected to build — if they don't match, run
'git pull origin main' before debugging the build error.

git
  commit                  $(probe_git rev-parse HEAD)
  short                   $(probe_git rev-parse --short=10 HEAD)
  branch                  $(probe_git rev-parse --abbrev-ref HEAD)
  describe                $(probe_git describe --tags --always --dirty)
  commit-date             $(probe_git log -1 --format=%cI)
  commit-subject          $(probe_git log -1 --format=%s)

source
  declared-version        $(probe_constants_version)
$(fingerprint_file "constants.go" "$CONSTANTS_FILE")
$(fingerprint_file "updaterepo.go" "$UPDATEREPO_FILE")
$(fingerprint_file "updatedebugwindows.go" "$UPDATEDEBUG_FILE")

guards
$(detect_redecl_risk)
=====================================================================
EOF
