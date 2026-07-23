#!/usr/bin/env bash
# scripts/smoke-act.sh — run the installer-smoke CI jobs locally via
# nektos/act so version-mismatch regressions can be caught pre-push.
#
# Motivation: the v6.74.0 drain bug only surfaced in the Windows CI
# smoke job (installer-smoke-windows) because that job runs the
# freshly-built binary under a real Windows pwsh host. Running act
# locally with `-P windows-latest=-self-hosted` is not possible (act
# has no Windows runner image), but we CAN run the ubuntu companion
# (installer-smoke) locally as a fast pre-push sanity check that
# covers everything except the Windows-specific pipe drain path.
#
# Usage:
#   bash scripts/smoke-act.sh                # runs installer-smoke
#   bash scripts/smoke-act.sh <job-name>     # runs a specific ci.yml job
#
# Requires: `act` (https://github.com/nektos/act) and Docker.

set -euo pipefail

JOB="${1:-installer-smoke}"

if ! command -v act >/dev/null 2>&1; then
  echo "✗ act not found on PATH." >&2
  echo "  Install: https://github.com/nektos/act#installation" >&2
  echo "  macOS:   brew install act" >&2
  echo "  Linux:   curl -sSfL https://raw.githubusercontent.com/nektos/act/master/install.sh | bash" >&2
  exit 2
fi

if ! docker info >/dev/null 2>&1; then
  echo "✗ Docker daemon is not reachable. act requires Docker to run." >&2
  exit 2
fi

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

echo "→ Running CI job '$JOB' locally via act..."
echo "  (Windows-only jobs cannot run under act; use CI for those.)"
echo

# -j filters to a single job. --pull=false avoids re-downloading the
# runner image on every invocation once it's cached. Medium image
# has go + bash + git preinstalled which matches ubuntu-latest close
# enough for the source-build smoke.
exec act \
  -j "$JOB" \
  --pull=false \
  -P ubuntu-latest=catthehacker/ubuntu:act-latest
