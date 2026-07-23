#!/usr/bin/env bash
# File-traversal + binary-detection helpers for fix-repo.sh

MAX_FILE_BYTES=$((5 * 1024 * 1024))

BINARY_EXTS_RE='\.(png|jpg|jpeg|gif|webp|ico|pdf|zip|tar|gz|tgz|bz2|xz|7z|rar|woff|woff2|ttf|otf|eot|mp3|mp4|mov|wav|ogg|webm|class|jar|so|dylib|dll|exe|pyc)$'

is_binary_extension() {
  local path="$1"
  local lower
  lower="$(printf '%s' "$path" | tr '[:upper:]' '[:lower:]')"
  [[ "$lower" =~ $BINARY_EXTS_RE ]]
}

is_symlink_path() {
  [ -L "$1" ]
}

is_oversized_file() {
  local path="$1"
  local size
  size="$(wc -c <"$path" 2>/dev/null || echo 0)"
  [ "$size" -gt "$MAX_FILE_BYTES" ]
}

has_null_byte() {
  local path="$1"
  LC_ALL=C head -c 8192 "$path" 2>/dev/null | LC_ALL=C tr -d '\000' | LC_ALL=C wc -c | {
    read -r kept
    local total
    total="$(LC_ALL=C head -c 8192 "$path" 2>/dev/null | LC_ALL=C wc -c | tr -d ' ')"
    kept="$(echo "$kept" | tr -d ' ')"
    [ "$kept" -lt "$total" ]
  }
}

is_scannable_file() {
  local path="$1"
  if is_symlink_path "$path";       then return 1; fi
  if is_oversized_file "$path";     then return 1; fi
  if is_binary_extension "$path";   then return 1; fi
  if has_null_byte "$path";         then return 1; fi
  return 0
}
