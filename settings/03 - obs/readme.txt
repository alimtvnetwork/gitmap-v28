OBS Studio Settings Package
==========================

Contains scene collections (.json) and profile folders.

gitmap install obs-settings handles sync automatically:
1. Extracts the .zip to a temp directory
2. Copies .json files to %APPDATA%\obs-studio\basic\scenes\
3. Copies profile folders to %APPDATA%\obs-studio\basic\profiles\
4. Cleans up temp

OBS discovers scenes and profiles from these directories on startup.

Usage:
  gitmap install obs-settings   # Sync settings only
