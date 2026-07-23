#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────────────
# fix-repo.sh — rewrite prior versioned-repo-name tokens to current
#
# Spec: spec-authoring/22-fix-repo/01-spec.md
#
# Usage:
#   ./fix-repo.sh                  # default: replace last 2 versions
#   ./fix-repo.sh --2              # explicit
#   ./fix-repo.sh --3              # last 3 versions
#   ./fix-repo.sh --5              # last 5 versions
#   ./fix-repo.sh --all            # every prior version (1..Current-1)
#   ./fix-repo.sh --dry-run        # report only
#   ./fix-repo.sh --verbose        # list each modified file
# ──────────────────────────────────────────────────────────────────────

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=scripts/fix-repo/repo-identity.sh
. "$SCRIPT_DIR/scripts/fix-repo/repo-identity.sh"
# shellcheck source=scripts/fix-repo/file-scan.sh
. "$SCRIPT_DIR/scripts/fix-repo/file-scan.sh"
# shellcheck source=scripts/fix-repo/rewrite.sh
. "$SCRIPT_DIR/scripts/fix-repo/rewrite.sh"
# shellcheck source=scripts/fix-repo/config.sh
. "$SCRIPT_DIR/scripts/fix-repo/config.sh"
# shellcheck source=scripts/fix-repo/paired-literal-audit.sh
. "$SCRIPT_DIR/scripts/fix-repo/paired-literal-audit.sh"

EXIT_OK=0
EXIT_NOT_A_REPO=2
EXIT_NO_REMOTE=3
EXIT_NO_VERSION_SUFFIX=4
EXIT_BAD_VERSION=5
EXIT_BAD_FLAG=6
EXIT_WRITE_FAILED=7
EXIT_BAD_CONFIG=8
# E_PAIRED_LITERAL: post-rewrite audit found a {base}-v{Current} token
# next to a stale sibling literal of the previous version (most often
# in *_test.go fixtures). See .lovable/memory/issues/2026-05-02-fixrepo-
# paired-literal-desync.md.
EXIT_PAIRED_LITERAL=10

MODE="--2"
DRY_RUN=0
VERBOSE_FLAG=0
CONFIG_PATH=""

is_mode_flag() {
  case "$1" in --2|--3|--5|--all) return 0 ;; *) return 1 ;; esac
}

parse_args() {
  local mode_count=0 a expect_config=0
  for a in "$@"; do
    if [ "$expect_config" = "1" ]; then CONFIG_PATH="$a"; expect_config=0; continue; fi
    if is_mode_flag "$a"; then
      MODE="$a"; mode_count=$((mode_count + 1)); continue
    fi
    case "$a" in
      --dry-run) DRY_RUN=1 ;;
      --verbose) VERBOSE_FLAG=1 ;;
      --config)  expect_config=1 ;;
      --config=*) CONFIG_PATH="${a#--config=}" ;;
      -h|--help) print_help; exit 0 ;;
      *) echo "fix-repo: ERROR unknown flag '$a' (E_BAD_FLAG)" >&2; exit $EXIT_BAD_FLAG ;;
    esac
  done
  if [ "$expect_config" = "1" ]; then
    echo "fix-repo: ERROR --config requires a path (E_BAD_FLAG)" >&2; exit $EXIT_BAD_FLAG
  fi
  if [ "$mode_count" -gt 1 ]; then
    echo "fix-repo: ERROR multiple mode flags (E_BAD_FLAG)" >&2; exit $EXIT_BAD_FLAG
  fi
}

print_help() {
  cat <<'EOF'
fix-repo.sh — rewrite prior versioned-repo-name tokens to current.

Usage: ./fix-repo.sh [--2|--3|--5|--all] [--dry-run] [--verbose] [--config <path>]

Config file (default: ./fix-repo.config.json) supports:
  ignoreDirs:     repo-relative directory prefixes to skip
  ignorePatterns: glob patterns (** = any depth, * = within segment)

Spec: spec-authoring/22-fix-repo/01-spec.md
EOF
}

get_span_from_mode() {
  local mode="$1" current="$2"
  case "$mode" in
    --2)   echo 2 ;;
    --3)   echo 3 ;;
    --5)   echo 5 ;;
    --all) echo $((current - 1)) ;;
  esac
}

_assert_repo_root_resolved() {
  REPO_ROOT="$(get_repo_root || true)"
  [ -n "$REPO_ROOT" ] && return 0
  echo "fix-repo: ERROR not a git repository (E_NOT_A_REPO)" >&2; exit $EXIT_NOT_A_REPO
}

_assert_remote_parsed() {
  local url; url="$(get_remote_url || true)"
  [ -n "$url" ] || { echo "fix-repo: ERROR no remote URL found (E_NO_REMOTE)" >&2; exit $EXIT_NO_REMOTE; }
  parse_remote_url "$url" \
    || { echo "fix-repo: ERROR cannot parse remote URL '$url'" >&2; exit $EXIT_NO_REMOTE; }
}

_assert_version_suffix() {
  split_repo_version "$PARSED_REPO" \
    || { echo "fix-repo: ERROR no -vN suffix on repo name '$PARSED_REPO' (E_NO_VERSION_SUFFIX)" >&2; exit $EXIT_NO_VERSION_SUFFIX; }
  [ "$SPLIT_VERSION" -ge 1 ] \
    || { echo "fix-repo: ERROR version <= 0 (E_BAD_VERSION)" >&2; exit $EXIT_BAD_VERSION; }
}

resolve_identity() {
  _assert_repo_root_resolved
  _assert_remote_parsed
  _assert_version_suffix
}

print_header() {
  local current="$1" mode="$2" targets_str="$3"
  echo "fix-repo  base=$SPLIT_BASE  current=v$current  mode=$mode"
  if [ -z "$targets_str" ]; then
    echo "targets:  (none)"
  else
    local pretty=""
    for n in $targets_str; do pretty="$pretty v$n,"; done
    echo "targets: $(echo "$pretty" | sed 's/,$//')"
  fi
  echo "host:     $PARSED_HOST  owner=$PARSED_OWNER"
  echo
}

print_summary() {
  local scanned="$1" changed="$2" reps="$3" dry="$4"
  local label="write"
  if [ "$dry" = "1" ]; then label="dry-run"; fi
  echo
  echo "scanned: $scanned files"
  echo "changed: $changed files ($reps replacements)"
  echo "mode:    $label"
}

# Process one file. Updates: SWEEP_SCANNED, SWEEP_CHANGED, SWEEP_REPS, SWEEP_FAILED.
# Args: rel current target_arr_var_name
_process_one_file() {
  local rel="$1" current="$2"
  local full="$REPO_ROOT/$rel"
  [ -f "$full" ] || return 0
  is_ignored_path "$rel" && return 0
  is_scannable_file "$full" || return 0
  SWEEP_SCANNED=$((SWEEP_SCANNED + 1))
  local reps
  reps="$(rewrite_file "$full" "$SPLIT_BASE" "$current" "$DRY_RUN" "${_TARGET_ARR[@]}")" \
    || { echo "fix-repo: ERROR write failed for $rel" >&2; SWEEP_FAILED=1; return 0; }
  [ "$reps" -gt 0 ] || return 0
  SWEEP_CHANGED=$((SWEEP_CHANGED + 1))
  SWEEP_REPS=$((SWEEP_REPS + reps))
  SWEEP_CHANGED_FILES+=("$full")
  [ "$VERBOSE_FLAG" = "1" ] && echo "modified: $rel ($reps replacements)"
  return 0
}

run_sweep() {
  local current="$1" targets_str="$2" rel n
  SWEEP_SCANNED=0; SWEEP_CHANGED=0; SWEEP_REPS=0; SWEEP_FAILED=0
  SWEEP_CHANGED_FILES=()
  _TARGET_ARR=()
  for n in $targets_str; do _TARGET_ARR+=("$n"); done
  while IFS= read -r -d '' rel; do
    _process_one_file "$rel" "$current"
  done < <(cd "$REPO_ROOT" && git ls-files -z)
}

main() {
  parse_args "$@"
  resolve_identity
  load_fixrepo_config "$CONFIG_PATH" "$REPO_ROOT" \
    || exit $EXIT_BAD_CONFIG
  local current="$SPLIT_VERSION"
  local span; span="$(get_span_from_mode "$MODE" "$current")"
  local targets_str; targets_str="$(get_target_versions "$current" "$span" | sed 's/ *$//')"
  print_header "$current" "$MODE" "$targets_str"
  if [ -z "$targets_str" ]; then
    print_summary 0 0 0 "$DRY_RUN"
    echo "fix-repo: nothing to replace"
    exit $EXIT_OK
  fi
  run_sweep "$current" "$targets_str"
  print_summary "$SWEEP_SCANNED" "$SWEEP_CHANGED" "$SWEEP_REPS" "$DRY_RUN"
  if [ "$SWEEP_FAILED" = "1" ]; then exit $EXIT_WRITE_FAILED; fi
  if ! run_paired_literal_audit "$SPLIT_BASE" "$current" "$DRY_RUN" "${SWEEP_CHANGED_FILES[@]:-}"; then
    exit $EXIT_PAIRED_LITERAL
  fi
  exit $EXIT_OK
}

main "$@"
