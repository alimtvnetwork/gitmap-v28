# Update Terminal Visualization

Slug: update-terminal-visualization
Steps: 30
Status: pending
Created: 2026-08-20

## Context
The user requested a terminal visualization enhancement for the update process. It must clearly show the current version and target version before the update begins, and output a summary detailing the previous version, the updated version, and the source of the update. Affected files include `constants_messages.go`, `constants_update.go`, `update.go`, and `updateremoteinstall.go`.
Links:
- Spec: [.lovable/spec/tasks/07-update-terminal-visualization.md](file:///d:/work/gitmap/.lovable/spec/tasks/07-update-terminal-visualization.md)
- Adheres to: `.lovable/coding-guidelines.md` and `spec/12-consolidated-guidelines/02-go-code-style.md`

## Steps
1. Group 1: Constants Update (Execute via standalone agent "Constants Editor"). Read `gitmap/constants/constants_messages.go` to locate existing update messages.
2. In `gitmap/constants/constants_messages.go`, add `MsgUpdateVersionCompare` constant with the value `"\n  → Upgrading gitmap from v%s to v%s...\n"`.
3. In `gitmap/constants/constants_messages.go`, add `MsgUpdateSummaryDetail` constant with the value `"\n  ✓ Successfully updated from v%s to v%s\n  → Source: %s\n\n"`.
4. In `gitmap/constants/constants_messages.go`, locate and safely remove or replace the outdated `MsgUpdateVersion` constant if it conflicts with the new summary.
5. In `gitmap/constants/constants_update.go`, add `ErrUpdateVersionRead` constant with the value `"  ⚠ Could not read target version from %s: %v\n"`.
6. Commit the Constants Update group: verify absolutely NO test results, artifacts, or compiled binaries are staged before making the commit. Ensure `.gitignore` explicitly excludes them. Fix any git issues and push the commit to the repository immediately.
7. Group 2: Local Update Logic (Execute via standalone agent "Local Update Enhancer"). Read `gitmap/cmd/update.go` to analyze the `runUpdateRunner` execution flow.
8. Create a new helper function `readTargetVersion(repoPath string) string` in `gitmap/cmd/update.go` to parse the target version.
9. Implement `readTargetVersion` to read `version.json` from the provided `repoPath` and unmarshal the `"version"` string field.
10. Implement error handling in `readTargetVersion` to return `"unknown"` and log `ErrUpdateVersionRead` to `os.Stderr` if `version.json` is missing or invalid.
11. Modify `runUpdateRunner` in `gitmap/cmd/update.go` to extract the current system version using `constants.Version`.
12. Modify `runUpdateRunner` to call `readTargetVersion(repoPath)` and store the target version in a local variable.
13. Modify `runUpdateRunner` to print `MsgUpdateVersionCompare` using the current and target versions *before* calling `executeUpdate`.
14. Modify `runUpdateRunner` to print `MsgUpdateSummaryDetail` using the current version, target version, and `repoPath` *after* `runPostUpdateMigrate` and `report.summarize()` complete.
15. Commit the Local Update Logic group: verify absolutely NO artifacts are staged before making the commit. Ensure `.gitignore` excludes them. Fix any git issues and push the commit to the repository immediately.
16. Group 3: Remote Update Logic (Execute via standalone agent "Remote Update Enhancer"). Read `gitmap/cmd/updateremoteinstall.go` to analyze the `runUpdateRemoteInstall` execution flow.
17. Create a new helper function `fetchRemoteTargetVersion(slug string) string` in `gitmap/cmd/updateremoteinstall.go`.
18. Implement `fetchRemoteTargetVersion` to construct the remote raw URL for `version.json` based on the provided slug and `constants.UpdateRepoOwner`.
19. Implement an HTTP GET request within `fetchRemoteTargetVersion` with a standard 10-second timeout to fetch the remote `version.json`.
20. Implement JSON unmarshaling in `fetchRemoteTargetVersion` to extract and return the `"version"` string field, falling back to `"unknown"` and logging the error on network or parsing failures.
21. Modify `runUpdateRemoteInstall` in `gitmap/cmd/updateremoteinstall.go` to extract the current system version using `constants.Version`.
22. Modify `runUpdateRemoteInstall` to call `fetchRemoteTargetVersion(slug)` and store the remote target version in a local variable.
23. Modify `runUpdateRemoteInstall` to print `MsgUpdateVersionCompare` using the current and remote target versions *before* executing the installer script.
24. Modify `runUpdateRemoteInstall` to print `MsgUpdateSummaryDetail` using the current version, remote target version, and remote source URL upon successful completion.
25. Commit the Remote Update Logic group: verify absolutely NO artifacts are staged before making the commit. Ensure `.gitignore` excludes them. Fix any git issues and push the commit to the repository immediately.
26. Group 4: Verification and Release (Execute via standalone agent "Release Manager"). Run `go fmt ./gitmap/cmd/... ./gitmap/constants/...` to ensure formatting complies with the project style guide.
27. Run `goimports -w` on the modified files to ensure correct imports (`encoding/json`, `net/http`) are added and unused ones are removed.
28. Run `go test ./gitmap/cmd/...` to ensure the new helper functions and update flows do not break existing unit tests.
29. Run `golangci-lint run --path-prefix=gitmap --timeout=5m` to verify no new linting issues were introduced by the additions.
30. Execute the release ceremony per `11-release.md` (bump MINOR version, update changelog with the whole plan summary, pin new version in root README), and immediately commit and push the final release.

## Verification
- Run `go run ./gitmap update` in a valid source directory and observe the pre-update version announcement and post-update summary visualization in the terminal.
- Run `go run ./gitmap update` utilizing the remote fallback and verify the HTTP request fetches the version and outputs the same visualization structure.

## Appended from prior pending tasks
none
