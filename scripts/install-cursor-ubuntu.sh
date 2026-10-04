#!/usr/bin/env bash
# ==============================================================================
# GitMap Ubuntu Cursor Provisioning & Project Synchronization Script
# ==============================================================================
set -euo pipefail

echo "=================================================================="
echo "  GitMap — Ubuntu Cursor IDE Automated Provisioning & Sync"
echo "=================================================================="

# 1. Platform & Architecture Verification
ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    echo "[!] Error: Cursor Linux binary requires x86_64 architecture (detected: $ARCH)."
    exit 1
fi

# 2. Determine target directories
if [ "$EUID" -eq 0 ]; then
    INSTALL_DIR="/opt/cursor"
    BIN_LINK="/usr/local/bin/cursor"
    DESKTOP_DIR="/usr/share/applications"
    ICON_DIR="/usr/share/icons/hicolor/scalable/apps"
else
    INSTALL_DIR="$HOME/.local/share/cursor"
    BIN_LINK="$HOME/.local/bin/cursor"
    DESKTOP_DIR="$HOME/.local/share/applications"
    ICON_DIR="$HOME/.local/share/icons/hicolor/scalable/apps"
fi

mkdir -p "$INSTALL_DIR" "$(dirname "$BIN_LINK")" "$DESKTOP_DIR" "$ICON_DIR"

# 3. Download Cursor AppImage
APPIMAGE_PATH="$INSTALL_DIR/cursor.appimage"
DOWNLOAD_URL="https://downloader.cursor.sh/linux/appImage/x64"

echo "[*] Downloading Cursor AppImage from official endpoint..."
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$DOWNLOAD_URL" -o "$APPIMAGE_PATH"
elif command -v wget >/dev/null 2>&1; then
    wget -qO "$APPIMAGE_PATH" "$DOWNLOAD_URL"
else
    echo "[!] Error: Neither curl nor wget is available."
    exit 1
fi

chmod +x "$APPIMAGE_PATH"

# 4. Create launcher wrapper script
WRAPPER_SCRIPT="$INSTALL_DIR/cursor"
cat << EOF > "$WRAPPER_SCRIPT"
#!/usr/bin/env bash
exec "$APPIMAGE_PATH" --no-sandbox "\$@"
EOF
chmod +x "$WRAPPER_SCRIPT"

ln -sf "$WRAPPER_SCRIPT" "$BIN_LINK"
echo "[+] Installed Cursor executable wrapper to: $BIN_LINK"

# 5. Write embedded SVG Icon
ICON_PATH="$ICON_DIR/cursor.svg"
cat << "EOF" > "$ICON_PATH"
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="cursor-grad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#000000" />
      <stop offset="100%" stop-color="#1e1e1e" />
    </linearGradient>
  </defs>
  <rect width="512" height="512" rx="100" fill="url(#cursor-grad)"/>
  <path d="M120 100 L380 256 L240 280 L320 400 L260 430 L180 310 L120 370 Z" fill="#58a6ff"/>
</svg>
EOF

# 6. Create Desktop Application Launcher
DESKTOP_FILE="$DESKTOP_DIR/cursor.desktop"
cat << EOF > "$DESKTOP_FILE"
[Desktop Entry]
Version=1.0
Type=Application
Name=Cursor
GenericName=AI Code Editor
Comment=Cursor AI-first Code Editor
Exec=$BIN_LINK %F
Icon=$ICON_PATH
Terminal=false
Categories=Development;IDE;TextEditor;
StartupWMClass=Cursor
MimeType=text/plain;inode/directory;
EOF
chmod +x "$DESKTOP_FILE"
echo "[+] Desktop application launcher registered: $DESKTOP_FILE"

# 7. Project Synchronization into Cursor Project Manager (projects.json)
CURSOR_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/Cursor/User/globalStorage/alefragnani.project-manager"
mkdir -p "$CURSOR_CONFIG_DIR"
PROJECTS_FILE="$CURSOR_CONFIG_DIR/projects.json"

echo "[*] Synchronizing Cursor Project Manager ($PROJECTS_FILE)..."
if command -v gitmap >/dev/null 2>&1; then
    echo "[*] Running gitmap cursor sync..."
    gitmap cursor sync || true
else
    if [ ! -f "$PROJECTS_FILE" ]; then
        echo "[]" > "$PROJECTS_FILE"
    fi
    echo "[+] Initialized Cursor Project Manager database: $PROJECTS_FILE"
fi

echo "=================================================================="
echo "  [SUCCESS] Cursor IDE provisioned and synchronized successfully!"
echo "  Launch anytime by running: cursor [target-directory]"
echo "=================================================================="
