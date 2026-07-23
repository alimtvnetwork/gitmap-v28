#!/usr/bin/env bash
# ----------------------------------------------------------------------
# visibility-change.sh - toggle/set GitHub/GitLab repo visibility
#
# Mirrors the Go-native `gitmap make-public` / `gitmap make-private`
# commands so CI and shell users have a script-only fallback.
# ----------------------------------------------------------------------

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=visibility-change/provider.sh
. "$SCRIPT_DIR/visibility-change/provider.sh"
# shellcheck source=visibility-change/apply.sh
. "$SCRIPT_DIR/visibility-change/apply.sh"

EXIT_OK=0
EXIT_NOT_A_REPO=2
EXIT_NO_ORIGIN=3
EXIT_BAD_PROVIDER=4
EXIT_AUTH_FAILED=5
EXIT_BAD_FLAG=6
EXIT_CONFIRM_REQ=7
EXIT_VERIFY_FAILED=8

VISIBLE_RAW=""
YES_FLAG=0
DRY_RUN=0

print_help() {
  cat <<'EOF'
visibility-change.sh - toggle/set GitHub/GitLab repo visibility.

Usage:
  ./visibility-change.sh                  # toggle current visibility
  ./visibility-change.sh --visible pub    # force public
  ./visibility-change.sh --visible pri    # force private
  ./visibility-change.sh --yes            # skip private->public prompt
  ./visibility-change.sh --dry-run        # preview, no API call
  ./visibility-change.sh -h | --help

Env:
  VISIBILITY_GITLAB_HOSTS  comma-separated allow-list of self-hosted GitLab hosts
EOF
}

err() { echo "$@" >&2; }

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      --visible)   VISIBLE_RAW="${2:-}"; shift 2 ;;
      --yes|-y)    YES_FLAG=1; shift ;;
      --dry-run)   DRY_RUN=1; shift ;;
      -h|--help)   print_help; exit $EXIT_OK ;;
      *) err "visibility-change: ERROR unknown flag '$1'"; exit $EXIT_BAD_FLAG ;;
    esac
  done
}

# Echoes 'public'|'private' for valid input, '' for empty, 'INVALID' otherwise
resolve_target_value() {
  local raw="$1"
  if [ -z "$raw" ]; then echo ""; return 0; fi
  local v; v="$(echo "$raw" | tr '[:upper:]' '[:lower:]')"
  case "$v" in
    pub|public)  echo "public" ;;
    pri|private) echo "private" ;;
    *)           echo "INVALID" ;;
  esac
}

is_repo_root() {
  git rev-parse --show-toplevel >/dev/null 2>&1
}

required_cli_for() {
  if [ "$1" = "github" ]; then echo "gh"; else echo "glab"; fi
}

resolve_context() {
  if ! is_repo_root; then
    err "visibility-change: ERROR not a git repository"; exit $EXIT_NOT_A_REPO
  fi
  URL="$(get_origin_url)"
  if [ -z "$URL" ]; then
    err "visibility-change: ERROR no origin remote"; exit $EXIT_NO_ORIGIN
  fi
  resolve_provider "$URL"
  if [ -z "$PROVIDER" ]; then
    err "visibility-change: ERROR unsupported host in '$URL'"; exit $EXIT_BAD_PROVIDER
  fi
  resolve_owner_repo "$URL"
  if [ -z "$SLUG" ]; then
    err "visibility-change: ERROR cannot parse owner/repo from '$URL'"; exit $EXIT_BAD_PROVIDER
  fi
}

assert_cli_ready() {
  local cli="$1"
  if is_cli_available "$cli"; then return 0; fi
  err "visibility-change: ERROR '$cli' not found on PATH (install: https://cli.github.com / https://gitlab.com/gitlab-org/cli)"
  exit $EXIT_AUTH_FAILED
}

resolve_next_target() {
  local forced="$1" current="$2"
  if [ -n "$forced" ]; then echo "$forced"; return 0; fi
  if [ "$current" = "public" ]; then echo "private"; else echo "public"; fi
}

maybe_confirm() {
  local current="$1" target="$2"
  if [ "$target" != "public" ]; then return 0; fi
  if [ "$current" != "private" ]; then return 0; fi
  if [ "$YES_FLAG" = "1" ]; then return 0; fi
  if confirm_public_change "$SLUG" "$PROVIDER"; then return 0; fi
  err "visibility-change: ERROR confirmation required (pass --yes for non-interactive)"
  exit $EXIT_CONFIRM_REQ
}

_validate_forced() {
  local forced="$1"
  [ "$forced" = "INVALID" ] || return 0
  err "visibility-change: ERROR bad --visible value '$VISIBLE_RAW' (use pub|public|pri|private)"
  exit $EXIT_BAD_FLAG
}

_read_current_or_die() {
  local cur; cur="$(get_current_visibility "$PROVIDER" "$SLUG")"
  [ -n "$cur" ] || { err "visibility-change: ERROR cannot read current visibility (auth?)"; exit $EXIT_AUTH_FAILED; }
  printf '%s' "$cur"
}

_apply_and_verify() {
  local current="$1" target="$2"
  apply_visibility "$PROVIDER" "$SLUG" "$target" \
    || { err "visibility-change: ERROR apply failed"; exit $EXIT_AUTH_FAILED; }
  visibility_matches "$PROVIDER" "$SLUG" "$target" \
    || { err "visibility-change: ERROR verification failed (visibility did not change)"; exit $EXIT_VERIFY_FAILED; }
}

main() {
  parse_args "$@"
  local forced; forced="$(resolve_target_value "$VISIBLE_RAW")"
  _validate_forced "$forced"
  resolve_context
  assert_cli_ready "$(required_cli_for "$PROVIDER")"
  local current; current="$(_read_current_or_die)"
  local target; target="$(resolve_next_target "$forced" "$current")"
  [ "$current" != "$target" ] || { echo "visibility: already $current ($PROVIDER)"; exit $EXIT_OK; }
  maybe_confirm "$current" "$target"
  [ "$DRY_RUN" != "1" ] || { echo "[dry-run] visibility: $current -> $target ($PROVIDER)"; exit $EXIT_OK; }
  _apply_and_verify "$current" "$target"
  echo "visibility: $current -> $target ($PROVIDER)"
}

main "$@"
