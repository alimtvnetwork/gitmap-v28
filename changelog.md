# Changelog

## [v6.504.0] - 2026-10-06

### Added
- `gitmap login`: unified GitHub authentication — `--token <PAT>` (validated via api.github.com), `--web` browser login (gh device flow or guided token page), interactive chooser, `--status`
- `gitmap logout`: removes the stored GitHub credential
- Logged-in token is reused automatically by clone/pull/push via token resolution
- Secure no-echo token prompt; fail-fast guidance when stdin is not a TTY

---

## [v6.503.0] - 2026-10-06


### Added
- chrome profile import export e2e test verification and test suite fixes


---


## [v6.502.0] - 2026-10-06


### Added
- resolve CI/CD sends push assertion and ERD parity cluster path


---


## [v6.501.0] - 2026-10-06


### Added
- Deep spec consolidation, canonical 8-cluster reduction, and memory unification


---


## [v6.500.0] - 2026-10-06


### Added
- Pre-deep-consolidation baseline release


---


## [v6.499.0] - 2026-10-06


### Added
- App spec and completed plans consolidation and reduction


---


## [v6.498.0] - 2026-10-06


### Added
- Pre-consolidation baseline release


---


## [v6.497.0] - 2026-10-06


### Added
- **Ubuntu IDE and GitHub Desktop Scan Sync Suite**: Resolved omission where `gitmap scan` on Ubuntu failed to add repositories into VS Code and GitHub Desktop.
- **Root Cause Analysis (Issue 71)**: Grounded 4-part RCA in `02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md` identifying missing directory creation in VS Code Project Manager extension path and missing Linux CLI resolution candidates for GitHub Desktop.
