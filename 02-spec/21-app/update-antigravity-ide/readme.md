# Specification: Update Antigravity IDE Across Developer Toolchain

## Specification Index

- **`01-architecture-and-gitmap-command.md`**: Architectural breakdown and implementation specification for `gitmap update antigravity` / `gitmap update ide` and Antigravity tools updater integration.
- **`02-scripts-fixer-and-live-system-update.md`**: Specification for `scripts-fixer` update action addition and live Linux host update verification.

## Target Repositories

1. **`gitmap` (`git-work/gitmap`)**: Core developer CLI engine providing `gitmap update antigravity` and `gitmap update ide`.
2. **`scripts-fixer` (`git-work/scripts-fixer`)**: Scripts fixture containing Linux automation modules (`scripts-linux/69-install-antigravity/`).
3. **`antigravity-manager` (`git-work/antigravity-manager`)**: Antigravity tools desktop and CLI manager hosting maintenance automation scripts.

## Core Mandates

- **Strict Relative Git Paths**: Zero absolute paths (`/home/a/...` or `C:\...`) across all documents, specs, code, and release notes.
- **GitMap Acceleration Primacy**: All discovery and file queries handled via `gitmap aum search`, `gitmap find`, `gitmap cat`.
- **Zero Raw Git Commits**: Releases and pushes executed exclusively via `gitmap cpf "<module> - <summary>"`.
