#!/usr/bin/env bash
# Config loader + ignore-matching for fix-repo.sh.
#
# Loads a JSON config (default: <repo-root>/fix-repo.config.json) with:
#   ignoreDirs:     array of repo-relative directory prefixes to skip
#   ignorePatterns: array of glob patterns (** = any depth, * = within segment)
#
# Populates these globals:
#   FIXREPO_IGNORE_DIRS      (newline-separated list)
#   FIXREPO_IGNORE_PATTERNS  (newline-separated list)
#   FIXREPO_CONFIG_PATH      (resolved path, or empty if none)

FIXREPO_IGNORE_DIRS=""
FIXREPO_IGNORE_PATTERNS=""
FIXREPO_CONFIG_PATH=""

_extract_json_array() {
  # Args: <file> <key>. Emits one element per line. Pure-bash, no jq.
  local file="$1" key="$2"
  python3 - "$file" "$key" <<'PY' 2>/dev/null || true
import json, sys
try:
    with open(sys.argv[1], 'r', encoding='utf-8') as f:
        data = json.load(f)
except Exception:
    sys.exit(0)
arr = data.get(sys.argv[2], [])
if not isinstance(arr, list):
    sys.exit(0)
for item in arr:
    if isinstance(item, str) and item:
        print(item)
PY
}

resolve_config_path() {
  local explicit="$1" repo_root="$2"
  if [ -n "$explicit" ]; then
    [ -f "$explicit" ] || { echo "fix-repo: ERROR config file not found: $explicit" >&2; return 1; }
    FIXREPO_CONFIG_PATH="$explicit"; return 0
  fi
  local default="$repo_root/fix-repo.config.json"
  [ -f "$default" ] && FIXREPO_CONFIG_PATH="$default"
  return 0
}

load_fixrepo_config() {
  local explicit="$1" repo_root="$2"
  resolve_config_path "$explicit" "$repo_root" || return 1
  [ -n "$FIXREPO_CONFIG_PATH" ] || return 0
  FIXREPO_IGNORE_DIRS="$(_extract_json_array "$FIXREPO_CONFIG_PATH" ignoreDirs)"
  FIXREPO_IGNORE_PATTERNS="$(_extract_json_array "$FIXREPO_CONFIG_PATH" ignorePatterns)"
  return 0
}

_path_starts_with_dir() {
  local rel="$1" dir="$2"
  dir="${dir%/}"
  [ -z "$dir" ] && return 1
  [ "$rel" = "$dir" ] && return 0
  case "$rel" in "$dir"/*) return 0 ;; esac
  return 1
}

is_ignored_dir() {
  local rel="$1" dir
  [ -n "$FIXREPO_IGNORE_DIRS" ] || return 1
  while IFS= read -r dir; do
    [ -n "$dir" ] || continue
    _path_starts_with_dir "$rel" "$dir" && return 0
  done <<EOF
$FIXREPO_IGNORE_DIRS
EOF
  return 1
}

_glob_to_regex() {
  # Translate **, *, ? into ERE. Escapes regex metacharacters.
  local pat="$1" out="" i ch
  for (( i=0; i<${#pat}; i++ )); do
    ch="${pat:i:1}"
    case "$ch" in
      '*')
        if [ "${pat:i+1:1}" = "*" ]; then out+=".*"; i=$((i+1)); else out+="[^/]*"; fi ;;
      '?') out+="[^/]" ;;
      '.'|'+'|'('|')'|'['|']'|'{'|'}'|'^'|'$'|'|'|'\\') out+="\\$ch" ;;
      *) out+="$ch" ;;
    esac
  done
  printf '^%s$' "$out"
}

is_ignored_pattern() {
  local rel="$1" pat re
  [ -n "$FIXREPO_IGNORE_PATTERNS" ] || return 1
  while IFS= read -r pat; do
    [ -n "$pat" ] || continue
    re="$(_glob_to_regex "$pat")"
    [[ "$rel" =~ $re ]] && return 0
  done <<EOF
$FIXREPO_IGNORE_PATTERNS
EOF
  return 1
}

is_ignored_path() {
  local rel="$1"
  is_ignored_dir "$rel" && return 0
  is_ignored_pattern "$rel" && return 0
  return 1
}
