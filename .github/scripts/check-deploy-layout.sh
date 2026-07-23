#!/usr/bin/env bash
# check-deploy-layout.sh
#
# Hard CI guarantee: the deploy-target folder MUST be `gitmap-cli/`
# (NOT `gitmap/`). Fails the build if any source file hardcodes a
# DEPLOY path using the legacy `gitmap` subfolder instead of reading
# from gitmap/constants/deploy-manifest.json.
#
# Scope: this guard targets *deploy-target* path construction only.
# It does NOT flag:
#   - Go import paths (github.com/.../gitmap/...)
#   - The source-repo subdirectory <RepoRoot>/gitmap/ (where code lives)
#   - URLs like raw.githubusercontent.com/.../gitmap/scripts/install.sh
#   - tmpDir-based source extracts (gitmap/scripts/install.ps1)
#   - Anything containing the marker `deploy-layout-allow`
#   - Anything on a line that also says `legacy` or `LegacyAppSubdirs`
#
# What it DOES flag — known-bad deploy-target idioms:
#
#   PowerShell:
#     Join-Path $deployPath  "gitmap"
#     Join-Path $cfg.deployPath "gitmap"
#     Join-Path $target "gitmap"
#     "$deployPath\gitmap\gitmap.exe"
#     \gitmap\gitmap.exe       (any literal Windows deploy path)
#
#   Go:
#     filepath.Join(deployPath, "gitmap", binaryName)
#     filepath.Join(cfg.DeployPath, "gitmap", ...)
#
#   Shell:
#     "$DEPLOY_DIR/gitmap/$BIN"
#
# Allowlist (file paths exempt from scanning):
EXEMPT_FILES=(
  "gitmap/constants/deploy-manifest.json"
  "gitmap/constants/deploy_manifest.go"
  ".github/scripts/check-deploy-layout.sh"
  ".github/scripts/check-legacy-refs.sh"
  ".github/scripts/smoke-installer.sh"
)
#
# Exit codes:
#   0 — clean
#   1 — at least one violation
#   2 — internal error

set -uo pipefail

ROOT="${1:-.}"

EXCLUDE_DIRS=(
  ".git" "node_modules" "dist" "build" "bin" ".next"
  ".gitmap" "vendor" "coverage" ".lovable" "spec"
)

# Tight, intent-focused regex set. Each anchors on a deploy-context
# variable name OR an absolute Windows path that is unmistakably a
# deploy target.
PATTERNS=(
  # PowerShell:  Join-Path $<deployVar> "gitmap"   (NOT followed by -cli)
  'Join-Path[[:space:]]+\$[A-Za-z_.]*[Dd]eploy[A-Za-z_.]*[[:space:]]+"gitmap"([^-]|$)'
  'Join-Path[[:space:]]+\$target[[:space:]]+"gitmap"([^-]|$)'
  # PowerShell/shell: "$<deployVar>\gitmap\..." or "/gitmap/..." after deploy var
  '\$[A-Za-z_.]*[Dd]eploy[A-Za-z_.]*[\\/]gitmap[\\/]'
  # Hardcoded Windows deploy path:  ...\gitmap\gitmap.exe  (binary inside same-named folder)
  '[\\/]gitmap[\\/]gitmap\.exe'
  # Go: filepath.Join(<deployVar>, "gitmap",   — deployVar contains 'eploy' (Deploy/deploy)
  'filepath\.Join\([^)]*[Dd]eploy[A-Za-z_]*,[[:space:]]*"gitmap"[[:space:]]*,'
  # Shell: "$DEPLOY.../gitmap/"
  '\$\{?[A-Z_]*DEPLOY[A-Z_]*\}?[/\\]gitmap[/\\]'
)

GREP_ARGS=(-RHInE)
for d in "${EXCLUDE_DIRS[@]}"; do
  GREP_ARGS+=(--exclude-dir="$d")
done
GREP_ARGS+=(--exclude="*.png" --exclude="*.jpg" --exclude="*.zip"
            --exclude="*.exe" --exclude="*.bin" --exclude="*.db"
            --exclude="*.sqlite" --exclude="*.woff*" --exclude="*.ttf"
            --exclude="*.json" --exclude="*.md" --exclude="*.html")

EXEMPT_FILES=(
  "gitmap/constants/deploy-manifest.json"
  "gitmap/constants/deploy_manifest.go"
  ".github/scripts/check-deploy-layout.sh"
  ".github/scripts/check-legacy-refs.sh"
  ".github/scripts/smoke-installer.sh"
)

violations_total=0
all_matches=""

for pat in "${PATTERNS[@]}"; do
  raw="$(grep "${GREP_ARGS[@]}" "$pat" "$ROOT" 2>/dev/null || true)"
  [ -z "$raw" ] && continue

  filtered="$raw"
  for f in "${EXEMPT_FILES[@]}"; do
    filtered="$(printf '%s\n' "$filtered" | grep -v "^\./$f:" | grep -v "^$f:" || true)"
  done

  # Strip explicit allow-marker, legacy-migration context, and false positives
  filtered="$(printf '%s\n' "$filtered" \
    | grep -v 'deploy-layout-allow' \
    | grep -viE '\blegacy\b|LegacyAppSubdirs' \
    | grep -v 'gitmap-cli' \
    | grep -v 'gitmap-updater' \
    || true)"

  if [ -n "$filtered" ]; then
    count="$(printf '%s\n' "$filtered" | wc -l | tr -d ' ')"
    violations_total=$((violations_total + count))
    all_matches="${all_matches}
::error::Pattern: $pat
${filtered}
"
  fi
done

if [ "$violations_total" -eq 0 ]; then
  echo "  [deploy-layout] OK — no hardcoded 'gitmap/' deploy paths found."
  echo "  [deploy-layout] Single source of truth: gitmap/constants/deploy-manifest.json"
  exit 0
fi

echo "::error::Found $violations_total hardcoded reference(s) to legacy 'gitmap/' deploy folder."
echo "::error::Deploy folder MUST be 'gitmap-cli/' (sourced from deploy-manifest.json)."
echo ""
printf '%s\n' "$all_matches"
echo ""
echo "  Fix options:"
echo "    Go:         use constants.GitMapCliSubdir (loaded from manifest)"
echo "    PowerShell: load via Get-DeployManifest in run.ps1"
echo "    Shell:      load via load_deploy_manifest in smoke-installer.sh"
echo "    Legitimate legacy-migration code: add 'legacy' to the line, or"
echo "    the marker '# deploy-layout-allow'."
exit 1
