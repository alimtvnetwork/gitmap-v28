Windows Terminal Settings
========================

Place your settings.json file here.

gitmap install wt-settings handles sync automatically:
1. Finds %LOCALAPPDATA%\Packages\Microsoft.WindowsTerminal_*\LocalState\
2. Copies settings.json to that directory
3. Copies any additional files (themes, fragments) alongside it

Windows Terminal reads settings.json from LocalState on startup.

Usage:
  gitmap install wt-settings     # Sync settings only
