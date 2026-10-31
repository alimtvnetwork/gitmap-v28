# Plan 369 - CI/CD fix and release for multiple packages

Date: 31-Oct-26
Repo: gitmap-v28, base main 58292e98 (Release v6.643.0)
Spec: 02-spec/21-app/369-cicd-fix-and-release-for-multiple-packages/
Subtasks: .ai-memory/plans/subtasks/369-cicd-fix-and-release-for-multiple-packages/

## User request (verbatim, condensed)

High-priority autopilot task: pull first, spec + plan first, fix CI/CD
until fully green including every linter, minor version bump via the
repo's own bump script, create the release tag and release, verify
tag on remote and release workflow green with assets, final
verification that tree is clean and local == origin. Solo worker -
no subagents, no questions, fix immediately, push, report only when
done or genuinely blocked.

## Subtasks

1. 01-diagnose-ci-failures - reproduce the red jobs locally and list
   the exact failing files, tests, and linter rules.
2. 02-fix-linters-and-tests - smallest fixes that turn the red jobs
   green; workflow-file fixes, if any, committed locally only and
   reported BLOCKED.
3. 03-bump-and-release - minor bump, commit + push, release tag via
   gitmap release, verify tag, release run, assets, and final tree.

## Status

- Step 1 (pull): done, 58292e98, clean tree.
- Step 2 (spec + plan): this file, spec folder, and subtasks.
- Remaining steps follow the spec method in order.
