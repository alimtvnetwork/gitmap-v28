#!/usr/bin/env bash
# Rewrite engine for fix-repo.sh — pure POSIX (awk/sed), no perl dependency.

# Echoes space-separated target versions
get_target_versions() {
  local current="$1" span="$2"
  local start=$((current - span))
  if [ "$start" -lt 1 ]; then start=1; fi
  local end=$((current - 1))
  if [ "$end" -lt "$start" ]; then return 0; fi
  seq "$start" "$end" | tr '\n' ' '
}

# Counts literal-token occurrences NOT followed by a digit, in $1.
# Args: file base n
count_token_occurrences() {
  local file="$1" base="$2" n="$3"
  local token="$base-v$n"
  awk -v tok="$token" '
    BEGIN { tlen = length(tok); total = 0 }
    {
      line = $0
      while ((p = index(line, tok)) > 0) {
        next_char = substr(line, p + tlen, 1)
        if (next_char !~ /[0-9]/) { total++ }
        line = substr(line, p + tlen)
      }
    }
    END { print total }
  ' "$file"
}

# awk program: replace literal $tok with $rep, but only when the char
# immediately following $tok is not a digit (numeric overflow guard).
read -r -d '' _SUBSTITUTE_AWK <<'AWK' || true
BEGIN { tlen = length(tok) }
{
  out = ""; line = $0
  while ((p = index(line, tok)) > 0) {
    next_char = substr(line, p + tlen, 1)
    if (next_char !~ /[0-9]/) {
      out = out substr(line, 1, p - 1) rep
    } else {
      out = out substr(line, 1, p + tlen - 1)
    }
    line = substr(line, p + tlen)
  }
  print out line
}
AWK

# Rewrites $1 in place via a temp file. Args: file base n current
substitute_token_in_file() {
  local file="$1" base="$2" n="$3" current="$4"
  local token="$base-v$n" replacement="$base-v$current"
  local tmp; tmp="$(mktemp)"
  awk -v tok="$token" -v rep="$replacement" "$_SUBSTITUTE_AWK" "$file" > "$tmp"
  mv "$tmp" "$file"
}

# Rewrites all targets in $1, echoes total replacement count.
# Args: path base current dry n1 [n2 ...]
rewrite_file() {
  local path="$1" base="$2" current="$3" dry="$4"; shift 4
  local total=0 n count
  for n in "$@"; do
    count="$(count_token_occurrences "$path" "$base" "$n")"
    if [ "$count" -gt 0 ] && [ "$dry" != "1" ]; then
      substitute_token_in_file "$path" "$base" "$n" "$current"
    fi
    total=$((total + count))
  done
  printf '%s' "$total"
}
