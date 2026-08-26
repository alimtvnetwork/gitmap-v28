#!/usr/bin/env bash
# Apply + verify visibility change via gh / glab.

apply_visibility() {
  local provider="$1" slug="$2" target="$3"
  if [ "$provider" = "github" ]; then
    gh repo edit "$slug" --visibility "$target" --accept-visibility-change-consequences
    return $?
  fi
  glab repo edit "$slug" --visibility "$target"
}

visibility_matches() {
  local provider="$1" slug="$2" target="$3"
  local actual; actual="$(get_current_visibility "$provider" "$slug")"
  [ "$actual" = "$target" ]
}

confirm_public_change() {
  local slug="$1" provider="$2"
  if [ ! -t 0 ]; then return 1; fi
  echo
  echo "⚠  About to make $slug PUBLIC on $provider."
  echo "   Type 'yes' to continue, anything else aborts:"
  local answer
  read -r answer
  [ "$answer" = "yes" ]
}
