# Subtask 3: Final Release

1. Execute the release script: `python .lovable/release/bump_versions.py --type minor`. If it doesn't exist, manually bump the version in `constants/version.go` or equivalent.
2. Update the root `readme.md` (lowercase) and pin the new version in the Install section using the dynamic fallback `Invoke-WebRequest` snippet.
3. Add the release notes to `changelog.md` or `.lovable/memory/release-architecture-map.md`.
4. Run `git add .` and `go run main.go release` (or standard `git commit -m "chore: release vX.Y.Z" && git push` as a fallback if `gitmap release` is unavailable).
5. Mark task as `STATUS: DONE` in `.lovable/temp-agents/task-10-03.md`.
