#!/usr/bin/env bash
# Repo-identity helpers for fix-repo.sh

get_repo_root() {
  git rev-parse --show-toplevel 2>/dev/null
}

get_remote_url() {
  local url
  url="$(git remote get-url origin 2>/dev/null || true)"
  if [ -n "$url" ]; then echo "$url"; return 0; fi
  git remote -v 2>/dev/null | awk '$3=="(fetch)"{print $2; exit}'
}

# Tries each regex in order; on first match sets PARSED_* via given group indices.
# Args: trimmed_url regex host_idx owner_idx repo_idx
_try_match_url() {
  local url="$1" re="$2" hi="$3" oi="$4" ri="$5"
  [[ "$url" =~ $re ]] || return 1
  PARSED_HOST="${BASH_REMATCH[$hi]}"
  PARSED_OWNER="${BASH_REMATCH[$oi]}"
  PARSED_REPO="${BASH_REMATCH[$ri]}"
  return 0
}

# Sets globals: PARSED_HOST, PARSED_OWNER, PARSED_REPO
parse_remote_url() {
  local url="$1"
  local trimmed="${url%/}"; trimmed="${trimmed%.git}"
  _try_match_url "$trimmed" '^https?://([^/:]+)(:[0-9]+)?/([^/]+)/([^/]+)$' 1 3 4 && return 0
  _try_match_url "$trimmed" '^git@([^:]+):([^/]+)/([^/]+)$'                 1 2 3 && return 0
  _try_match_url "$trimmed" '^ssh://git@([^/:]+)(:[0-9]+)?/([^/]+)/([^/]+)$' 1 3 4 && return 0
  return 1
}

# Sets globals: SPLIT_BASE, SPLIT_VERSION
split_repo_version() {
  local repo="$1"
  local re='^(.+)-v([0-9]+)$'
  if [[ "$repo" =~ $re ]]; then
    SPLIT_BASE="${BASH_REMATCH[1]}"; SPLIT_VERSION="${BASH_REMATCH[2]}"; return 0
  fi
  return 1
}
