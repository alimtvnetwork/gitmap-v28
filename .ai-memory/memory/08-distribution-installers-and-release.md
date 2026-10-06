# 08 — Distribution, Installers & Release Ceremony Architecture

- **Domain:** Software Distribution, Cross-Platform Installers & Release Management
- **Authoritative Specification:** [08-distribution-and-release](../../02-spec/21-app/08-distribution-and-release/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. Unified Cross-Platform Runners (`run.ps1` & `run.sh`)

GitMap provides standardized entry point scripts for developer environments across all platforms:

- **Runners:** `run.ps1` (PowerShell for Windows) and `run.sh` (POSIX shell for Linux/macOS).
- **Configuration Driver (`run.config.json`):** Defines commands, default ports, environment variables, and pre-flight dependency prerequisites in a centralized configuration.
- **Dependency Self-Healing:** If required runtimes (Go, Node, Python) are missing, the runner delegates to local installer scripts (`local-install.ps1` or `local-install.sh`) and re-executes automatically.

---

## 2. Multi-OS Installer Framework

GitMap supports native packaging and automated setup across Windows, macOS, and Linux:

- **Windows NSIS Installer:**
  - Auto-detects architecture (`x64`, `arm64`) and installs GitMap CLI to `C:\Program Files\GitMap\`.
  - Appends installation directory to system `PATH` and registers uninstallation keys in the Windows Registry.
- **Linux Archive Packaging (`.tar.gz` / `.zip`):**
  - Packages binaries with proper POSIX file permissions (`0755`).
  - Generates XDG desktop launchers (`gitmap.desktop`) and installs icons to `/usr/share/icons/hicolor/`.
  - Distributes standalone tarballs (`gitmap-linux-amd64.tar.gz`) for zero-dependency server installs.
- **macOS Tarball & Homebrew:**
  - Distributes signed macOS binaries supporting both Intel (`amd64`) and Apple Silicon (`arm64`).

---

## 3. Centralized Version.json & SemVer Architecture

All version management in GitMap is driven by a single authoritative file:

```json
{
  "version": "6.499.0",
  "buildDate": "2026-10-06",
  "commitHash": "git",
  "releaseNotes": "Consolidated application specifications and AI memory"
}
```

- **Invariant:** Go code, PowerShell scripts, Python tools, and installer scripts must read version metadata from `version.json` (or embed it via Go `-ldflags` during compilation). Never hardcode version strings in source code.
- **Semantic Versioning Bumps:**
  - `patch-bump`: Bug fixes and minor tweaks (`x.y.Z+1`).
  - `minor-bump`: Standard features, performance improvements, and operational enhancements (`x.Y+1.0`).
  - `major-bump`: Breaking API changes and fundamental architectural redesigns (`X+1.0.0`).

---

## 4. Release Ceremony & Untouchable Pipeline Constraint

- **Release Ceremony Flow:**
  1. Verify all static linters and quality gates pass.
  2. Bump `version.json` and synchronize `package.json` / documentation headers.
  3. Generate release notes and update changelog.
  4. Create git tag `vX.Y.Z` and push to upstream repository.
- **Untouchable CI Release Pipeline Invariant:**
  - `.github/workflows/release.yml` and related release scripts (`.github/scripts/smoke-installer.*`) are strictly off-limits.
  - Subagents and refactoring tasks are forbidden from modifying or "improving" release workflow definitions.
