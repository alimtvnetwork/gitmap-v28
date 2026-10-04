#!/bin/bash
set -euo pipefail

# ==============================================================================
# run-migration.sh: Automated AGM Settings Migration & Provisioning on Linux
# ==============================================================================

TAR_ARCHIVE="${1:-/tmp/agm_windows_export.tar.gz}"
AGM_TOOLS_DIR="${HOME}/.antigravity_tools"
GEMINI_DIR="${HOME}/.gemini"
CONFIG_DIR="${HOME}/.config/Antigravity"
LOCAL_BIN="${HOME}/.local/bin"

echo "[1/6] Validating target directories..."
mkdir -p "${AGM_TOOLS_DIR}/accounts"
mkdir -p "${GEMINI_DIR}"
mkdir -p "${CONFIG_DIR}/User/globalStorage"
mkdir -p "${LOCAL_BIN}"

if [ ! -f "${TAR_ARCHIVE}" ]; then
    echo "ERROR: Migration archive not found at: ${TAR_ARCHIVE}" >&2
    exit 1
fi

echo "[2/6] Extracting settings archive into ${AGM_TOOLS_DIR}..."
tar -xzf "${TAR_ARCHIVE}" -C "${AGM_TOOLS_DIR}"

echo "[3/6] Applying secure file permissions..."
chmod 700 "${AGM_TOOLS_DIR}"
chmod 700 "${AGM_TOOLS_DIR}/accounts"
chmod 600 "${AGM_TOOLS_DIR}"/*.json "${AGM_TOOLS_DIR}"/*.db 2>/dev/null || true
chmod 600 "${AGM_TOOLS_DIR}/accounts"/*.json 2>/dev/null || true

echo "[4/6] Installing fixed Antigravity IDE launcher to prevent recursion loops..."
cat << 'LAUNCHER_EOF' > /tmp/antigravity_launcher_fixed.sh
#!/bin/bash
if [ -x "${HOME}/.local/share/antigravity-ide/antigravity" ]; then
    DIR="${HOME}/.local/share/antigravity-ide"
elif [ -d "/opt/antigravity" ]; then
    DIR="/opt/antigravity"
else
    DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fi
export LD_LIBRARY_PATH="${DIR}:${DIR}/lib:${LD_LIBRARY_PATH:-}"
EXEC="${DIR}/antigravity"

export ELECTRON_OZONE_PLATFORM_HINT="auto"
export DONT_PROMPT_WSL_INSTALL=1

exec "${EXEC}" --no-sandbox "$@"
LAUNCHER_EOF

chmod +x /tmp/antigravity_launcher_fixed.sh
cp /tmp/antigravity_launcher_fixed.sh "${LOCAL_BIN}/antigravity"
cp /tmp/antigravity_launcher_fixed.sh "${LOCAL_BIN}/antigravity-ide"
if command -v sudo >/dev/null 2>&1; then
    sudo cp /tmp/antigravity_launcher_fixed.sh /usr/local/bin/antigravity 2>/dev/null || true
    sudo cp /tmp/antigravity_launcher_fixed.sh /usr/local/bin/antigravity-ide 2>/dev/null || true
fi

echo "[5/6] Verifying required system dependencies (libsecret-tools)..."
if ! command -v secret-tool >/dev/null 2>&1; then
    if command -v sudo >/dev/null 2>&1 && command -v apt-get >/dev/null 2>&1; then
        echo "Installing libsecret-tools..."
        sudo apt-get update -qq && sudo apt-get install -y -qq libsecret-tools
    fi
fi

echo "[6/6] Verifying AGM account index..."
if command -v agm >/dev/null 2>&1; then
    ACCOUNT_COUNT=$(find "${AGM_TOOLS_DIR}/accounts" -name "*.json" | wc -l)
    echo "SUCCESS: Migration complete. ${ACCOUNT_COUNT} accounts registered."
    echo "Run 'agm list' or 'agm switch <email>' to activate an account."
else
    echo "WARNING: 'agm' binary not found in PATH. Install AGM via gitmap install agm."
fi
