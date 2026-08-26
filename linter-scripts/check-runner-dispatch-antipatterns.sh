#!/usr/bin/env bash
# linter-scripts/check-runner-dispatch-antipatterns.sh
#
# Grep-based guard. Fails CI if run.sh or run.ps1 reintroduce forbidden
# dispatch anti-patterns in the `fix-repo` branch specifically.
#
# Spec: spec/15-distribution-and-runner/06-fix-repo-forwarding.md
#
# Strategy: extract ONLY the fix-repo dispatch region from each runner,
# then assert presence of required patterns and absence of forbidden
# patterns within that region. Findings are accumulated into a structured
# summary printed at the end, plus GitHub Actions `::error file=,line=::`
# annotations so violations surface inline on PRs.
#
# Exit codes:
#   0 — clean
#   1 — at least one anti-pattern found (or required pattern missing)
#   2 — a required runner file is missing, or region not extractable
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SH_PATH="$REPO_ROOT/run.sh"
PS_PATH="$REPO_ROOT/run.ps1"
SH_REL="run.sh"
PS_REL="run.ps1"

# Findings format (TSV): file<TAB>line<TAB>kind<TAB>reason<TAB>pattern<TAB>match<TAB>snippet
FINDINGS_FILE="$(mktemp)"
trap 'rm -f "$FINDINGS_FILE"' EXIT

# GitHub Actions annotation values must escape %, \r, \n, : and ,
# See: https://docs.github.com/actions/using-workflows/workflow-commands-for-github-actions
gha_escape() {
  printf '%s' "$1" \
    | sed -e 's/%/%25/g' -e 's/\r/%0D/g' -e 's/,/%2C/g' -e 's/:/%3A/g' \
    | tr '\n' ' '
}

record_finding() {
  local file="$1" line="$2" kind="$3" reason="$4" pattern="$5" match="$6" snippet="$7"
  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
    "$file" "$line" "$kind" "$reason" "$pattern" "$match" "$snippet" >> "$FINDINGS_FILE"
  # GitHub annotation: title carries the kind+reason; message body carries
  # the regex and matched substring so the developer can pinpoint the hit
  # without re-scanning the file by hand.
  local title body
  title="$(gha_escape "[$kind] $reason")"
  body="$(gha_escape "regex=$pattern  match=<<$match>>")"
  printf '::error file=%s,line=%s,title=%s::%s\n' "$file" "$line" "$title" "$body" >&2
}

extract_region_sh() { grep -nE '^[[:space:]]*fix-repo\)' "$SH_PATH" || true; }
extract_region_ps() { grep -nE '^[[:space:]]*"fix-repo"[[:space:]]*\{' "$PS_PATH" || true; }

region_line_no() { printf '%s' "$1" | head -n1 | cut -d: -f1; }
region_text()    { printf '%s' "$1" | head -n1 | cut -d: -f2-; }

# Print ±2 lines of context around $line in $file, prefixed with line numbers
# and a `>` marker on the offending line. Output is plain text, suitable
# for direct inclusion in the terminal summary and markdown fenced blocks.
CONTEXT_RADIUS="${RUNNER_GUARD_CONTEXT:-2}"
fetch_context() {
  local rel="$1" line="$2" abs
  case "$rel" in
    run.sh)  abs="$SH_PATH" ;;
    run.ps1) abs="$PS_PATH" ;;
    *)       return 0 ;;
  esac
  [ -f "$abs" ] || return 0
  local start=$((line - CONTEXT_RADIUS))
  local end=$((line + CONTEXT_RADIUS))
  [ "$start" -lt 1 ] && start=1
  awk -v s="$start" -v e="$end" -v hit="$line" '
    NR < s { next }
    NR > e { exit }
    { marker = (NR == hit) ? ">" : " "
      printf "    %s %4d | %s\n", marker, NR, $0 }
  ' "$abs"
}

forbid_in_region() {
  local file="$1" region="$2" pattern="$3" reason="$4"
  local line text match
  line="$(region_line_no "$region")"
  text="$(region_text "$region")"
  match="$(printf '%s' "$text" | { grep -oE -e "$pattern" || true; } | head -n1)"
  [ -n "$match" ] || return 0
  record_finding "$file" "$line" "FORBIDDEN" "$reason" "$pattern" "$match" "$text"
}

require_in_region() {
  local file="$1" region="$2" pattern="$3" reason="$4"
  local line text
  line="$(region_line_no "$region")"
  text="$(region_text "$region")"
  printf '%s' "$text" | grep -qE -e "$pattern" && return 0
  record_finding "$file" "$line" "MISSING" "$reason" "$pattern" "(no match — required pattern absent)" "$text"
}

[ -f "$SH_PATH" ] || { echo "::error file=$SH_REL::run.sh missing" >&2; exit 2; }
[ -f "$PS_PATH" ] || { echo "::error file=$PS_REL::run.ps1 missing" >&2; exit 2; }

SH_REGION="$(extract_region_sh)"
PS_REGION="$(extract_region_ps)"

[ -n "$SH_REGION" ] || { echo "::error file=$SH_REL::no 'fix-repo)' case arm found" >&2; exit 2; }
[ -n "$PS_REGION" ] || { echo "::error file=$PS_REL::no '\"fix-repo\"' switch arm found" >&2; exit 2; }

echo "── scanning fix-repo dispatch arms ─────────────────────"
echo "   $SH_REL:$(region_line_no "$SH_REGION")"
echo "   $PS_REL:$(region_line_no "$PS_REGION")"

# run.sh — forbidden
forbid_in_region "$SH_REL" "$SH_REGION" '"\$\*"'                      '"$*" collapses argv into one string; use "$@"'
forbid_in_region "$SH_REL" "$SH_REGION" '\beval\b'                    'eval-based dispatch breaks argv quoting'
forbid_in_region "$SH_REL" "$SH_REGION" '\bbash[[:space:]]+-c\b'      '`bash -c "..."` wrapper loses argv boundaries'
forbid_in_region "$SH_REL" "$SH_REGION" '\bsh[[:space:]]+-c\b'        '`sh -c "..."` wrapper loses argv boundaries'
forbid_in_region "$SH_REL" "$SH_REGION" 'fix-repo\.sh[^|]*\|[^|&]'    'piping fix-repo.sh output masks its exit code'
forbid_in_region "$SH_REL" "$SH_REGION" '\$\{@:[0-9]+\}'              '${@:N} slicing drops original argv; forward "$@" verbatim'
forbid_in_region "$SH_REL" "$SH_REGION" "printf[[:space:]]+['\"]?%[qs]" 'printf %q/%s rebuilds argv from a joined string'
forbid_in_region "$SH_REL" "$SH_REGION" '\bIFS=[^[:space:]]'          'mutating IFS in dispatch alters argv splitting'
forbid_in_region "$SH_REL" "$SH_REGION" '\$\([^)]*"\$@"[^)]*\)'       'command substitution on "$@" stringifies argv'
forbid_in_region "$SH_REL" "$SH_REGION" '\bxargs\b'                   'xargs reformats argv via stdin and loses quoting'
forbid_in_region "$SH_REL" "$SH_REGION" 'fix-repo\.sh[^&]*&[[:space:]]*$' \
                                                                       'background dispatch (&) detaches child; loses exit code'
# run.sh — required
require_in_region "$SH_REL" "$SH_REGION" '\bexec\b[^#]*fix-repo\.sh[^#]*"\$@"' \
  'must use `exec ... fix-repo.sh "$@"` to forward argv and exit code'

# run.ps1 — forbidden
forbid_in_region "$PS_REL" "$PS_REGION" '\$args[[:space:]]*-join'     '$args -join collapses argv into a single string'
forbid_in_region "$PS_REL" "$PS_REGION" '-join[[:space:]]+\$args'     '-join $args collapses argv into a single string'
forbid_in_region "$PS_REL" "$PS_REGION" '"\$args"'                    '"$args" interpolation flattens argv; use @args splatting'
forbid_in_region "$PS_REL" "$PS_REGION" '\[string\]::Join'            '[string]::Join on $args collapses argv into one string'
forbid_in_region "$PS_REL" "$PS_REGION" '\$args[[:space:]]*\.ToString\(' \
                                                                       '$args.ToString() flattens argv to a single string'
forbid_in_region "$PS_REL" "$PS_REGION" '\$args[[:space:]]+-as[[:space:]]+\[string\]' \
                                                                       '$args -as [string] flattens argv to a single string'
forbid_in_region "$PS_REL" "$PS_REGION" '\[string\]\$args'            '[string]$args casts argv array to a single joined string'
forbid_in_region "$PS_REL" "$PS_REGION" '\bInvoke-Expression\b'       'Invoke-Expression on argv breaks quoting and is unsafe'
forbid_in_region "$PS_REL" "$PS_REGION" '[[:space:]]iex[[:space:]]'   '`iex` (Invoke-Expression alias) on argv is unsafe'
forbid_in_region "$PS_REL" "$PS_REGION" 'Start-Process\b'             'Start-Process detaches the child; cannot propagate $LASTEXITCODE'
forbid_in_region "$PS_REL" "$PS_REGION" 'Start-Job\b'                 'Start-Job detaches the child; cannot propagate $LASTEXITCODE'
forbid_in_region "$PS_REL" "$PS_REGION" '\bcmd[[:space:]]+/c\b'       '`cmd /c` wrapper re-parses argv and loses quoting'
forbid_in_region "$PS_REL" "$PS_REGION" '\$args\[[0-9]+\.\.'          '$args[N..M] slicing drops original argv; forward @args verbatim'
# run.ps1 — required
require_in_region "$PS_REL" "$PS_REGION" '@args'                       'must invoke inner with @args splatting (preserves argv)'
require_in_region "$PS_REL" "$PS_REGION" 'exit[[:space:]]+\$LASTEXITCODE' \
  'must end with `exit $LASTEXITCODE` to propagate exit code'

print_summary_table() {
  local count="$1"
  echo
  echo "════════════════════ violation summary ════════════════════"
  printf '  %-3s  %-13s  %-10s  %s\n' "#" "FILE:LINE" "KIND" "REASON"
  printf '  %-3s  %-13s  %-10s  %s\n' "---" "-------------" "----------" "------"
  local n=0
  while IFS=$'\t' read -r file line kind reason pattern match snippet; do
    n=$((n + 1))
    printf '  %-3s  %-13s  %-10s  %s\n' "$n" "$file:$line" "$kind" "$reason"
    printf '       │ regex   : %s\n' "$pattern"
    printf '       │ match   : <<%s>>\n' "$match"
    printf '       │ context (±%s lines):\n' "$CONTEXT_RADIUS"
    fetch_context "$file" "$line" | sed 's/^/       │/'
    printf '       └ snippet : %s\n' "$snippet"
  done < "$FINDINGS_FILE"
  echo "═══════════════════════════════════════════════════════════"
  echo "❌ runner dispatch guard: $count violation(s)"
}

md_escape() { printf '%s' "$1" | sed -e 's/|/\\|/g' -e 's/`/\\`/g'; }

write_markdown_report() {
  local out="$1" status="$2" count="$3"
  {
    printf '# Runner Dispatch Guard Report\n\n'
    printf '%s\n' "- **Status:** $status"
    printf '%s\n' "- **Violations:** $count"
    printf '%s\n' "- **Generated:** $(date -u +'%Y-%m-%dT%H:%M:%SZ')"
    printf '%s\n' "- **Guard:** \`linter-scripts/check-runner-dispatch-antipatterns.sh\`"
    printf '%s\n\n' "- **Spec:** \`spec/15-distribution-and-runner/06-fix-repo-forwarding.md\` §7"
    if [ "$count" -eq 0 ]; then
      printf '✅ No anti-patterns found in `run.sh` or `run.ps1` `fix-repo` dispatch arms.\n'
      return 0
    fi
    printf '## Violations\n\n'
    printf '| # | File:Line | Kind | Reason | Regex | Match |\n'
    printf '|---|-----------|------|--------|-------|-------|\n'
    local n=0
    while IFS=$'\t' read -r file line kind reason pattern match _snippet; do
      n=$((n + 1))
      printf '| %d | `%s:%s` | %s | %s | `%s` | `%s` |\n' \
        "$n" "$file" "$line" "$kind" "$(md_escape "$reason")" \
        "$(md_escape "$pattern")" "$(md_escape "$match")"
    done < "$FINDINGS_FILE"
    printf '\n## Context (±%s lines)\n\n' "$CONTEXT_RADIUS"
    n=0
    while IFS=$'\t' read -r file line _kind _reason _pattern _match snippet; do
      n=$((n + 1))
      printf '### %d. `%s:%s`\n\n' "$n" "$file" "$line"
      printf '```\n'
      fetch_context "$file" "$line"
      printf '```\n\n'
      printf '_Matched line snippet:_ `%s`\n\n' "$(md_escape "$snippet")"
    done < "$FINDINGS_FILE"
  } > "$out"
  echo "📝 markdown report written: $out" >&2
}

REPORT_PATH="${RUNNER_GUARD_REPORT:-$REPO_ROOT/reports/runner-guard-report.md}"
mkdir -p "$(dirname "$REPORT_PATH")"

count="$(wc -l < "$FINDINGS_FILE" | tr -d ' ')"
if [ "$count" -eq 0 ]; then
  echo
  echo "✅ runner dispatch guard: no anti-patterns found"
  write_markdown_report "$REPORT_PATH" "PASS" 0
  exit 0
fi
print_summary_table "$count"
write_markdown_report "$REPORT_PATH" "FAIL" "$count"
exit 1
