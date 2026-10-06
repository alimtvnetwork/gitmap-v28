# 10 — Cross-Platform Installers and Runner Parity

- **Subsystem:** Packaging, Deployment, and Execution Runners
- **Status:** Authoritative Reference

## 1. Installer Package Architecture
- Windows: NSIS auto-detection, Chocolatey and Winget package manifests.
- Linux: Tarball archives (`.tar.gz`), embedded standalone installer scripts (`local-install.sh`).
- Version Pinning: Supports freezing targets to specific versions (`gitmap pin <target> <version>`).

## 2. Execution Runner Parity
- Feature-identical execution across `run.ps1` (PowerShell) and `run.sh` (Bash).
- Central configuration in `run.config.json` drives dependency pre-flight checks and build execution.
