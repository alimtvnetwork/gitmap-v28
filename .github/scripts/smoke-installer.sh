#!/usr/bin/env bash
# Smoke test: verify a freshly-installed gitmap reports the expected version.
#
# Modes:
#   source   Build gitmap from the current checkout into a tempdir, then run
#            `<tempdir>/gitmap version` and assert it contains v$EXPECTED.
#            Used by ci.yml on every PR — no network release dependency.
#
#   release  Run gitmap/scripts/install.sh against a published GitHub release
#            (--version "v$EXPECTED" --no-discovery), then run the installed
#            binary and assert. Used by release.yml after the release is cut.
#
# Reads $EXPECTED (e.g. "4.1.0") from env. Falls back to constants.Version.
#
# Exit 0 on success, non-zero with diagnostic on failure.
set -euo pipefail

MODE="${1:-source}"
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
EXPECTED="${EXPECTED:-$(awk -F'"' '/^const Version/ {print $2}' "$REPO_ROOT/gitmap/constants/constants.go")}"
EXPECTED="${EXPECTED#v}"

load_deploy_manifest() {
  local manifest_path="$REPO_ROOT/gitmap/constants/deploy-manifest.json"
  APP_SUBDIR="gitmap-cli"
  BINARY_NAME_UNIX="gitmap"
  LEGACY_APP_SUBDIRS=("gitmap")

  if [ ! -f "$manifest_path" ]; then
    return
  fi

  # Prefer jq (always available on ubuntu-latest GitHub runners), fall back to
  # python3 (always available too), then awk as last resort. The previous
  # awk-only path silently produced empty values when the runner's awk locale
  # mishandled the embedded double-quote field separator, which left the
  # post-install path resolution searching for $DEST/gitmap (legacy layout)
  # instead of $DEST/gitmap-cli/gitmap and tripping the "Binary not found"
  # exit-3 even though install.sh wrote the binary correctly.
  local parsed_app="" parsed_bin="" parsed_legacy=""
  if command -v jq >/dev/null 2>&1; then
    parsed_app="$(jq -r '.appSubdir // empty' "$manifest_path" 2>/dev/null || true)"
    parsed_bin="$(jq -r '.binaryName.unix // empty' "$manifest_path" 2>/dev/null || true)"
    parsed_legacy="$(jq -r '(.legacyAppSubdirs // []) | .[]' "$manifest_path" 2>/dev/null || true)"
  elif command -v python3 >/dev/null 2>&1; then
    parsed_app="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1])).get("appSubdir",""))' "$manifest_path" 2>/dev/null || true)"
    parsed_bin="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1])).get("binaryName",{}).get("unix",""))' "$manifest_path" 2>/dev/null || true)"
    parsed_legacy="$(python3 -c 'import json,sys
[print(x) for x in json.load(open(sys.argv[1])).get("legacyAppSubdirs",[])]' "$manifest_path" 2>/dev/null || true)"
  else
    parsed_app="$(awk -F'"' '/"appSubdir"/ {print $4; exit}' "$manifest_path")"
    parsed_bin="$(awk -F'"' '/"unix"/ {print $4; exit}' "$manifest_path")"
    parsed_legacy="$(awk '
      /"legacyAppSubdirs"/ { in_block=1 }
      in_block {
        while (match($0, /"[^"]+"/)) {
          value = substr($0, RSTART + 1, RLENGTH - 2)
          if (value != "legacyAppSubdirs") { print value }
          $0 = substr($0, RSTART + RLENGTH)
        }
      }
      in_block && /\]/ { exit }
    ' "$manifest_path")"
  fi

  [ -n "$parsed_app" ] && APP_SUBDIR="$parsed_app"
  [ -n "$parsed_bin" ] && BINARY_NAME_UNIX="$parsed_bin"
  if [ -n "$parsed_legacy" ]; then
    mapfile -t LEGACY_APP_SUBDIRS <<< "$parsed_legacy"
  fi

  if [ -z "$APP_SUBDIR" ]; then
    APP_SUBDIR="gitmap-cli"
  fi
  if [ -z "$BINARY_NAME_UNIX" ]; then
    BINARY_NAME_UNIX="gitmap"
  fi
  if [ "${#LEGACY_APP_SUBDIRS[@]}" -eq 0 ]; then
    LEGACY_APP_SUBDIRS=("gitmap")
  fi
}

if [ -z "$EXPECTED" ]; then
  echo "::error::Could not determine expected version" >&2
  exit 2
fi

load_deploy_manifest

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

echo "▶ Smoke mode:    $MODE"
echo "▶ Expected:      v$EXPECTED"
echo "▶ Workdir:       $WORK"

case "$MODE" in
  source)
    echo "▶ Building gitmap from source into $WORK"
    (cd "$REPO_ROOT/gitmap" && go build -o "$WORK/gitmap" .)
    BIN="$WORK/gitmap"
    ;;
  release)
    echo "▶ Running install.sh --version v$EXPECTED --no-discovery"
    DEST="$WORK/install"
    mkdir -p "$DEST"
    bash "$REPO_ROOT/gitmap/scripts/install.sh" \
      --version "v$EXPECTED" \
      --dir "$DEST" \
      --no-path \
      --no-discovery
    # Resolve from the deploy manifest so the smoke test tracks the
    # installer's canonical layout instead of stale hardcoded paths.
    echo "▶ Manifest resolved values:"
    echo "    APP_SUBDIR        = '$APP_SUBDIR'"
    echo "    BINARY_NAME_UNIX  = '$BINARY_NAME_UNIX'"
    echo "    LEGACY_APP_SUBDIRS = [${LEGACY_APP_SUBDIRS[*]}]"

    # ALWAYS dump the post-install directory listing — invaluable for
    # diagnosing layout drift between install.sh and this resolver. Printed
    # before path resolution so the tree is visible even when later steps
    # short-circuit on success.
    echo "▶ Post-install tree under $DEST (depth ≤ 4):"
    if [ -d "$DEST" ]; then
      find "$DEST" -maxdepth 4 | sed 's/^/    /'
    else
      echo "    (DEST does not exist!)"
    fi

    echo "▶ BIN candidate probe order:"
    BIN=""
    primary_candidates=(
      "$DEST/$APP_SUBDIR/$BINARY_NAME_UNIX"
      "$DEST/$BINARY_NAME_UNIX"
    )
    for candidate in "${primary_candidates[@]}"; do
      if [ -f "$candidate" ] && [ -x "$candidate" ]; then
        echo "    [hit ] $candidate"
        BIN="$candidate"
        break
      else
        reason="missing"
        [ -f "$candidate" ] && reason="not executable"
        echo "    [miss] $candidate ($reason)"
      fi
    done
    if [ -z "$BIN" ]; then
      for legacy in "${LEGACY_APP_SUBDIRS[@]}"; do
        candidate="$DEST/$legacy/$BINARY_NAME_UNIX"
        if [ -f "$candidate" ] && [ -x "$candidate" ]; then
          echo "    [hit ] $candidate (legacy)"
          BIN="$candidate"
          break
        else
          reason="missing"
          [ -f "$candidate" ] && reason="not executable"
          echo "    [miss] $candidate (legacy, $reason)"
        fi
      done
    fi
    # ALWAYS run the recursive find fallback if direct paths missed — guards
    # against any future layout drift between install.sh and the manifest
    # reader. Tier 1: exact name + executable. Tier 2: exact name, any perms
    # (chmod +x and warn — surfaces a missing-bit install bug). Tier 3: glob
    # gitmap* executables (catches renamed/suffixed binaries like gitmap.bin).
    if [ -z "$BIN" ]; then
      echo "▶ Direct paths missed; recursive search under $DEST"
      echo "    tier 1: -name $BINARY_NAME_UNIX -perm -u+x"
      BIN="$(find "$DEST" -type f -name "$BINARY_NAME_UNIX" -perm -u+x 2>/dev/null | head -n1 || true)"
      if [ -z "$BIN" ]; then
        echo "    tier 2: -name $BINARY_NAME_UNIX (any perms)"
        BIN="$(find "$DEST" -type f -name "$BINARY_NAME_UNIX" 2>/dev/null | head -n1 || true)"
        if [ -n "$BIN" ]; then
          echo "    ⚠  found $BIN without +x bit; chmod +x and continuing"
          chmod +x "$BIN" 2>/dev/null || true
        fi
      fi
      if [ -z "$BIN" ]; then
        echo "    tier 3: -name '${BINARY_NAME_UNIX}*' -perm -u+x"
        BIN="$(find "$DEST" -type f -name "${BINARY_NAME_UNIX}*" -perm -u+x 2>/dev/null | head -n1 || true)"
      fi
      if [ -n "$BIN" ]; then
        echo "    [hit ] $BIN (recursive fallback)"
      else
        echo "    [miss] no binary matching ${BINARY_NAME_UNIX}* found anywhere under $DEST"
      fi
    fi
    if [ -z "$BIN" ]; then
      echo "::error::Could not locate installed gitmap binary under $DEST" >&2
      exit 3
    fi
    echo "▶ Located binary: $BIN"
    ;;
  *)
    echo "::error::Unknown mode '$MODE' (expected 'source' or 'release')" >&2
    exit 2
    ;;
esac

if [ ! -f "$BIN" ] || [ ! -x "$BIN" ]; then
  echo "::error::Binary not found or not executable at $BIN" >&2
  exit 3
fi

VERSION_OUTPUT="$("$BIN" version 2>&1)"
ACTUAL="$(printf '%s\n' "$VERSION_OUTPUT" | awk '/^gitmap v[0-9]/{print; exit}')"
echo "▶ Actual output: $ACTUAL"

EXPECTED_LINE="gitmap v$EXPECTED"
if [ "$ACTUAL" != "$EXPECTED_LINE" ]; then
  echo "::error::Version mismatch" >&2
  echo "  expected: $EXPECTED_LINE" >&2
  echo "  actual:   $ACTUAL" >&2
  exit 4
fi

echo "✅ Installer smoke test passed: $ACTUAL"
