# Release Architecture Map

Releases for Gitmap are triggered by bumping the version in ersion.json. 
When the version is bumped, you must ensure the corresponding versions in 
eadme.md and changelog.md are synchronized.
Once the changes are committed, the CI/CD pipeline (when correctly configured) handles the build and artifact generation.

- v6.127.0: git-rm and folder export commands added
