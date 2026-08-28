# Implementation Spec: Fix Release Tag Commit Ordering

## Context
When performing a release with Gitmap, the release automation currently generates the release branch and the tag from the *current* commit, and *then* writes the `.gitmap/release/latest.json` metadata file and performs an auto-commit on the original branch.
This causes the tag (e.g. `v1.28.0`) to point to the commit *before* the `.gitmap/release/*` files are tracked, creating a misaligned Git history.

## Solution Architecture
1. In `gitmap/release/workflow.go`, inside `performRelease`, swap the execution order.
2. Execute `writeMetadataIfRequired` and `AutoCommit` FIRST. This advances `HEAD` to the auto-commit containing the `.gitmap/release/*` files.
3. Execute `executeSteps` SECOND. Since `executeSteps` creates the release branch off `sourceRef` (which defaults to the new `HEAD`), the branch and the tag will correctly attach to the new auto-commit.
4. Correct variable shadowing (`err :=` vs `err =`) resulting from this swap to ensure successful compilation.

## Task-Specific Custom Rules
1. **Rule 1: Re-ordering Integrity:** Ensure `executeSteps` receives the updated `HEAD` if `sourceRef` is empty. The `executeSteps` function inherently uses the current `HEAD` if no `sourceRef` is passed to `git checkout -b`, which perfectly aligns with our new ordering.
2. **Rule 2: Error Scoping:** Go's strict variable declaration rules require `err := writeMetadataIfRequired(...)` and `err = executeSteps(...)` due to the block scope swap. Do not shadow `err` inside conditional blocks.
3. **Rule 3: Asset Preservation:** Moving metadata generation before `executeSteps` leaves the `assets` list technically empty in `latest.json`, but this matches previous behavior and prevents scope creep. Do not attempt to refactor `pushAndFinalize` asset collection in this PR.

## Subtasks
Subtasks generated in: `.lovable/plans/subtasks/03-fix-release-tag-ordering/01-task.md`
