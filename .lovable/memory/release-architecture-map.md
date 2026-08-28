# Release Architecture Map

## v6.139.0
- Rewrote gitmap open to use native OS openers (Windows undll32, macOS open, Linux xdg-open).
- Centralized application errors to fully return *apperror.AppError and propagate correctly to inishCommandAudit.
- Included strict LLM grouping guidelines for AI agents to prevent un-orchestrated commits.


Releases for Gitmap are triggered by bumping the version in ersion.json. 
When the version is bumped, you must ensure the corresponding versions in 
eadme.md and changelog.md are synchronized.
Once the changes are committed, the CI/CD pipeline (when correctly configured) handles the build and artifact generation.

- v6.127.0: git-rm and folder export commands added

- v6.128.0: ignore and add commands added

- v6.129.0: ag and vscode install commands added

- v6.130.0: github-desktop apt install fix

- v6.133.0: search and llm feature spec added

- v6.134.0: Wired CLI File Search / Regex Search Commands to SplitDB Indexer Engine
