#!/usr/bin/env bash
# Post-rewrite audit: catches paired-literal desync where a
# {base}-v{Current} token was rewritten but a sibling bare digit
# literal representing the previous version was left behind
# (e.g. `"gitmap-v27", "12"`). Test-files only.
# See .lovable/memory/issues/2026-05-02-fixrepo-paired-literal-desync.md.

# Look-ahead: number of lines after a {base}-v{Current} hit we scan
# for a stale sibling digit. 2 covers the common Go map/slice layout.
PAIRED_AUDIT_LOOKAHEAD=2

is_test_file() {
  case "$1" in *_test.go) return 0 ;; *) return 1 ;; esac
}

# Args: file base current
# Echoes "<lineno>:<line>" for each hit; exits 0 (no early return on misses).
find_paired_literal_hits() {
  local file="$1" base="$2" current="$3"
  local prev=$((current - 1))
  [ "$prev" -ge 1 ] || return 0
  awk -v needle="$base-v$current" -v prev="$prev" -v window="$PAIRED_AUDIT_LOOKAHEAD" '
    BEGIN {
      qrx = "\"" prev "\""
      brx = "(^|[^v0-9])" prev "($|[^0-9])"
    }
    {
      lines[NR] = $0
    }
    END {
      for (i = 1; i <= NR; i++) {
        if (index(lines[i], needle) == 0) continue
        end = i + window
        if (end > NR) end = NR
        for (j = i; j <= end; j++) {
          if (lines[j] ~ qrx || lines[j] ~ brx) {
            print j ":" lines[j]
            break
          }
        }
      }
    }
  ' "$file"
}

# Args: base current dry changed-file [changed-file ...]
# Returns 0 = clean (or dry-run), 1 = hits found.
run_paired_literal_audit() {
  local base="$1" current="$2" dry="$3"; shift 3
  if [ "$dry" = "1" ]; then echo "audit:   skipped (dry-run)"; return 0; fi
  local total=0 files_with_hits=0 hits f
  for f in "$@"; do
    is_test_file "$f" || continue
    hits="$(find_paired_literal_hits "$f" "$base" "$current")"
    [ -n "$hits" ] || continue
    files_with_hits=$((files_with_hits + 1))
    while IFS= read -r line; do
      [ -n "$line" ] || continue
      total=$((total + 1))
      local lno content
      lno="${line%%:*}"
      content="${line#*:}"
      printf "fix-repo: AUDIT paired-literal at %s:%s: '%s-v%s' on/near sibling literal '%s'\n  -> line: %s\n" \
        "$f" "$lno" "$base" "$current" "$((current-1))" "$content" >&2
    done <<< "$hits"
  done
  if [ "$total" -eq 0 ]; then echo "audit:   no paired-literal desync detected"; return 0; fi
  printf "fix-repo: ERROR paired-literal audit failed: %d hit(s) in %d file(s) (E_PAIRED_LITERAL)\n" "$total" "$files_with_hits" >&2
  echo "  see .lovable/memory/issues/2026-05-02-fixrepo-paired-literal-desync.md" >&2
  echo "  fix: derive sibling literals from the same int via fmt.Sprintf" >&2
  return 1
}
