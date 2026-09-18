#!/usr/bin/env bash
# 1. Directories setup
ICON_DIR="$HOME/.local/share/icons/hicolor/scalable/apps"
ICON_PATH="$ICON_DIR/antigravity-ide.svg"
APP_DIR="$HOME/.local/share/applications"
DESKTOP_ENTRY="$APP_DIR/antigravity-ide.desktop"
mkdir -p "$ICON_DIR" "$APP_DIR" "$HOME/Desktop"

# 2. Write embedded crisp vector SVG icon matching the Antigravity ribbon
cat << "EOF" > "$ICON_PATH"
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#21232b" />
      <stop offset="100%" stop-color="#14151a" />
    </linearGradient>
    <linearGradient id="antigravity-arch" x1="0%" y1="100%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#2979FF" />
      <stop offset="30%" stop-color="#00E5FF" />
      <stop offset="55%" stop-color="#76FF03" />
      <stop offset="75%" stop-color="#FFD600" />
      <stop offset="90%" stop-color="#FF3D00" />
      <stop offset="100%" stop-color="#E040FB" />
    </linearGradient>
    <filter id="soft-glow" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="12" result="blur" />
      <feMerge>
        <feMergeNode in="blur" />
        <feMergeNode in="SourceGraphic" />
      </feMerge>
    </filter>
  </defs>
  <rect x="24" y="24" width="464" height="464" rx="112" fill="url(#bg)" />
  <path d="M 148 376 C 148 376 195 210 256 160 C 317 210 364 376 364 376 C 336 348 304 265 256 222 C 208 265 176 348 148 376 Z"
        fill="url(#antigravity-arch)" filter="url(#soft-glow)" />
</svg>
EOF

# 3. Clean up any obsolete/duplicate .desktop files
rm -f "$APP_DIR/antigravity.desktop" "$HOME/Desktop/antigravity.desktop"

# 4. Generate the clean Antigravity IDE launcher
cat << EOF > "$DESKTOP_ENTRY"
[Desktop Entry]
Version=1.0
Name=Antigravity IDE
GenericName=Text Editor
Comment=Antigravity IDE
Exec=$HOME/.local/share/antigravity/antigravity %U
Icon=$ICON_PATH
Terminal=false
Type=Application
Categories=Development;IDE;
StartupWMClass=antigravity
StartupNotify=true
MimeType=text/plain;inode/directory;
EOF

chmod +x "$DESKTOP_ENTRY"

# 5. Place on Desktop & enable launching
cp "$DESKTOP_ENTRY" "$HOME/Desktop/Antigravity IDE.desktop"
chmod +x "$HOME/Desktop/Antigravity IDE.desktop"
gio set "$HOME/Desktop/Antigravity IDE.desktop" metadata::trusted true 2>/dev/null

# 6. Update icon cache & desktop database
gtk-update-icon-cache -f -t "$HOME/.local/share/icons/hicolor" 2>/dev/null
update-desktop-database "$APP_DIR" 2>/dev/null

# 7. Deduplicate dock favorites and pin Antigravity IDE
python3 -c "
import subprocess, ast
raw = subprocess.check_output([\"gsettings\", \"get\", \"org.gnome.shell\", \"favorite-apps\"]).decode().replace(\"@as \", \"\").strip()
favs = [x for x in ast.literal_eval(raw) if x not in [\"antigravity.desktop\", \"antigravity-ide.desktop\"]]
favs.append(\"antigravity-ide.desktop\")
subprocess.run([\"gsettings\", \"set\", \"org.gnome.shell\", \"favorite-apps\", str(favs)])
"

echo "Done! High-resolution Antigravity IDE icon installed and pinned."
