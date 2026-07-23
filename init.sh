#!/usr/bin/env bash
# ----------------------------------------------------------------------
# init.sh - one-shot repo init: ensure repo is public, then rewrite
#           stale version tokens via fix-repo. Both steps always run.
#
# Order (per spec/03-general/11-init-pipeline.md):
#   1) visibility-change.sh --visible pub --yes  (no-op if already public)
#   2) fix-repo.sh --all
#
# Failure policy: best-effort. Both steps run regardless of the first
# step's exit code. Exits 0 only if both succeeded; otherwise exits
# with the first non-zero step exit code and prints a combined report.
#
# --yes is forwarded to visibility-change so the private->public
# confirmation never blocks. Pass --dry-run to preview both steps.
# ----------------------------------------------------------------------

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

DRY_RUN=0

print_help() {
  cat <<'EOF'
init.sh - run visibility-change (force public, auto-yes) then fix-repo --all.

Usage:
  ./init.sh             # ensure public, then rewrite stale version tokens
  ./init.sh --dry-run   # preview both steps
  ./init.sh -h | --help

Behavior:
  - Both steps always run (best-effort), even if the first fails.
  - Exit 0 only if both succeeded; otherwise exits with the first
    non-zero step exit code and prints a combined report.
EOF
}

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      --dry-run) DRY_RUN=1; shift ;;
      -h|--help) print_help; exit 0 ;;
      *) echo "init: ERROR unknown flag '$1'" >&2; exit 6 ;;
    esac
  done
}

run_step() {
  # run_step <label> <script> <gitmap-subcmd> <script-args-csv> <gitmap-args-csv>
  local label="$1" script="$2" sub="$3" sargs_csv="$4" gargs_csv="$5"
  local -a sargs=() gargs=()
  IFS='|' read -r -a sargs <<<"$sargs_csv"
  IFS='|' read -r -a gargs <<<"$gargs_csv"
  echo
  if [ -x "$SCRIPT_DIR/$script" ] || [ -f "$SCRIPT_DIR/$script" ]; then
    echo "==> [$label] $script ${sargs[*]}"
    "$SCRIPT_DIR/$script" "${sargs[@]}"
    return $?
  fi
  if command -v gitmap >/dev/null 2>&1; then
    echo "==> [$label] gitmap $sub ${gargs[*]}"
    gitmap "$sub" "${gargs[@]}"
    return $?
  fi
  echo "==> [$label] SKIP — neither $script nor 'gitmap' binary found on PATH" >&2
  return 127
}

write_summary() {
  local vis_rc="$1" fix_rc="$2"
  echo
  echo "==> init summary"
  echo "    visibility-change : exit $vis_rc"
  echo "    fix-repo          : exit $fix_rc"
}

main() {
  parse_args "$@"

  local vis_sargs="--visible|pub|--yes"
  local vis_gargs="--yes"
  if [ "$DRY_RUN" = "1" ]; then
    vis_sargs="$vis_sargs|--dry-run"
    vis_gargs="$vis_gargs|--dry-run"
  fi
  run_step visibility visibility-change.sh make-public "$vis_sargs" "$vis_gargs"
  local vis_rc=$?

  local fix_sargs="--all"
  local fix_gargs="--all"
  if [ "$DRY_RUN" = "1" ]; then
    fix_sargs="$fix_sargs|--dry-run"
    fix_gargs="$fix_gargs|--dry-run"
  fi
  run_step fix-repo fix-repo.sh fix-repo "$fix_sargs" "$fix_gargs"
  local fix_rc=$?

  write_summary "$vis_rc" "$fix_rc"

  [ "$vis_rc" = "0" ] || exit "$vis_rc"
  [ "$fix_rc" = "0" ] || exit "$fix_rc"
  exit 0
}

main "$@"