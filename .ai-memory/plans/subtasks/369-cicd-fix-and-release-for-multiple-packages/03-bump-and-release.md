# Subtask 03 - Bump and release

Date: 31-Oct-26
Owner: solo worker (parent task forbids subagents).

- Minor version bump via the repo's own bump mechanism.
- Commit + push via the gitmap commit suite.
- Release tag via gitmap release; verify tag with git ls-remote,
  verify the release workflow run is green, the release is published
  (not draft, not prerelease), and assets are present.
- Final check: tree clean, local == origin on main, latest CI on
  the release tip green. Report with run IDs/URLs, SHAs, version
  old -> new, tag, release URL, asset count, and BLOCKED section.
