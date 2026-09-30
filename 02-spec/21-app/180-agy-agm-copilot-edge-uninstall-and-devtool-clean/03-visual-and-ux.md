# Spec 180: CLI Interaction, Visual Output & Safety UX

Spec Reference: [01-overview.md](01-overview.md)

---

## 1. Command Syntax & Routing Topology

### 1.1 Antigravity Uninstallation
```bash
# Standard uninstallation (keeps brains, logs, user configs)
gitmap uninstall agy
gitmap agy uninstall

# Complete purge with automatic project/conversation snapshot backup
gitmap uninstall agy-all
gitmap agy uninstall-all
gitmap agy uninstall --all

# Optional flags
gitmap agy uninstall-all --force           # Skip interactive confirmation
gitmap agy uninstall-all --dry-run         # Inspect what would be purged without deleting
gitmap agy uninstall-all --backup <path>   # Custom snapshot JSON output path
```

### 1.2 AGM (Antigravity Manager) Uninstallation
```bash
gitmap uninstall agm
gitmap agm uninstall
gitmap uninstall agm-all
gitmap agm uninstall-all
```

### 1.3 Windows Copilot & Edge Removal
```bash
gitmap uninstall copilot
gitmap winutil copilot uninstall

gitmap uninstall edge
gitmap winutil edge uninstall
gitmap winutil edge uninstall --keep-webview2 # Preserves WebView2 runtime
```

### 1.4 Enhanced DevTool Cache Cleaner
```bash
gitmap devtool clear
gitmap dt clear
gitmap clean-dev
gitmap os dev-clean --all
gitmap os dev-clean --dry-run
```

---

## 2. Terminal UI & Confirmation UX

### 2.1 AGY Complete Purge Banner & Interactive Confirmation
```text
  ⚠️ WARNING: Full Antigravity Purge Requested
  -------------------------------------------------------------
  This will remove:
    • .gemini configuration directory (~/.gemini)
    • Antigravity Brain database & transcripts
    • Local app cache & extension data
    • CLI wrapper scripts & PATH registrations

  ✓ Saved recovery snapshot to: ~/.gitmap/agy-snapshot-20260928-171500.json
    (Captured 3 project workspaces and 14 conversations)

  Are you sure you want to proceed? [y/N]:
```

### 2.2 DevTool Clean Aligned Table
```text
  ┌──────────────────────┬─────────────┬──────────────┬────────┐
  │ Category             │ Items Found │ Space Freed  │ Status │
  ├──────────────────────┼─────────────┼──────────────┼────────┤
  │ Go Build & Cache     │ 142 items   │ 450.2 MB     │ Done   │
  │ Node & npm Cache     │ 1,208 items │ 1.2 GB       │ Done   │
  │ Python / pip Cache   │ 85 items    │ 210.5 MB     │ Done   │
  │ Vite / Build Caches  │ 34 items    │ 95.0 MB      │ Done   │
  │ Antigravity Logs     │ 512 items   │ 680.1 MB     │ Done   │
  │ System Temp Files    │ 890 items   │ 2.4 GB       │ Done   │
  ├──────────────────────┼─────────────┼──────────────┼────────┤
  │ Total Savings        │ 2,871 items │ 5.03 GB      │ OK     │
  └──────────────────────┴─────────────┴──────────────┴────────┘
```
