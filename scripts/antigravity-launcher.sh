#!/bin/bash
if [ -x "/home/a/.local/share/antigravity-ide/antigravity" ]; then
    DIR="/home/a/.local/share/antigravity-ide"
elif [ -d "/opt/antigravity" ]; then
    DIR="/opt/antigravity"
else
    DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fi
export LD_LIBRARY_PATH="${DIR}:${DIR}/lib:${LD_LIBRARY_PATH:-}"
EXEC="${DIR}/antigravity"

export ELECTRON_OZONE_PLATFORM_HINT="auto"
export DONT_PROMPT_WSL_INSTALL=1

# Unconfined tarball Electron builds require --no-sandbox on modern Linux (AppArmor userns restriction)
exec "${EXEC}" --no-sandbox "$@"
