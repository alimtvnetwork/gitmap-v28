# Execution Plan 242: aria2c Download, Smart Update, Diffstat Isolation & Remediation Sub-Node Tree

## User Request (Verbatim)
```text
# High Priority Instruction

Can you please try to fix this? You can see in some places the numbers exactly matches, which is insane, which couldn't be possible, right? That is one of the bug. I think when you summarize, you made a mistake in the counting. Find that issue, how you're making a mistake, so that you can fix it. That's one thing. The next thing is that when you write the code for the `gitmap`, one of the greatest issue that you did not include the download method. `gitmap` should have a download method that it basically going to use the area 2C automatically to download it. And also, one more factor I wanted to focus in is the update or installation. In update, if we give a flag, let's update to be some version or give the number of the version after the update. It should exactly update to that version. That's one thing. Another is when it's updating, that means you have to update the installed PowerShell or shell script. When it updates, if the update fails, then it should fall back to the previous tag and then previous. Basically, it's going to do the five times at least to find the last updated one. And during the update, it should check the last update, what is available, and what we have as a version. If these two matches, then it does not run the update script. It says, "Already updated." So remember to make it very fast so that we know when you run the command trying to make a casing call, usually in every day, it will check and save it to its SQLite DB, that if it has a version, current version, and what the tag is available. Not only the tag, it will also check along with the tag the executable files are there or not. In some cases, what happens is the tag is there, but executable file is not there. So in the database, we should have a table which actually have all these flags. And basically, if the executable file is not there yet, then it would not say the update is there. Same would happen if the user runs the update command. If the executable file is not there, it would say, "Executable file is not there for that update. Do you still want it to update to any one of the fallbacks?" Something like this, so it would be smart. So from the update command, we can do this stuff, but also from the installation, we have to have some smartness. And also when it downloads using the AI2C, make sure the padding is correct and it is in the center. And it needs to be exactly matching with the other padding, so make sure of that. Another factor is that when we update, we don't need to show so many script. It will just say the numbers. So update can have some flag. You can have JSON, and that would actually give some JSON as an output, like what is the progress or the percentage, what it is doing. And at the end, it will just summarize like from this to this, it has been updated. So these are the things I wanted you to have in the update mode. So make sure that you follow through all of these. And if there is a, let's say, if we do the pull all and if there is an issue during the pull all, let's say fail, and then you try to resolve it. Now, when you resolve it, the resolve should be the sub-node. But here everything looks like in the same node, which is wrong because you are executing a task inside the project, so it should be sub-node. The title and the sub-node, like is it successful or not. So this is how it needs to be. But it is not happening right now, so make sure you respect that and fix this issue. At the end do a minor bump and release.

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Fix the bug where numbers exactly match in some places, which shouldn't be possible.
5. Include a download method in `gitmap` that uses area 2C automatically.
6. Implement update functionality with a version flag to update to a specific version.
7. Ensure update process can fall back to the previous tag if it fails, attempting up to five times.
8. Check if the current version matches the available version before running the update script.
9. Store version and executable file availability in SQLite DB to determine update necessity.
10. Ensure AI2C download padding is correct and centered.
11. Simplify update script output to show only numbers and provide JSON output for progress.
12. Ensure pull all command issues are resolved as sub-nodes, not in the same node.
13. Perform a minor version bump and release.
```

## Architectural Overview

This plan orchestrates the implementation of 5 core capabilities across GitMap:
1. **Diffstat & Commit Counting Isolation**: Guard empty repo path evaluations across `gitutil.GetLastCommitSHA`, `queryGitDiffStat`, and ensure `cmd.Dir = cwd` in `executeGitPullCommand`.
2. **Remediation Sub-Node Tree UI**: Restructure batch pull failure remediation in `cli/cmdpull/pull_remediation.go` and `cli/cmdpull/pull_missing_remediate.go` to present repository as parent node (`• <repo>`) and remediation steps/results as indented tree sub-nodes (`├──`, `└──`).
3. **`gitmap download` (aria2c Acceleration & Centered Visual Padding)**: Unified CLI download command (`gitmap download <url> [--out <file>] [--json]`) automatically utilizing `aria2c` with fallback, centering the progress bar and aligning padding.
4. **Smart `gitmap update` Engine**:
   - Explicit version flag support: `gitmap update [version]` or `--version <v>`.
   - Skip update if already on target/latest version (`Already updated`).
   - Up to 5-tag fallback loop if update fails or release lacks executable assets.
   - `ReleaseCache` table in SQLite (`gitmap.db`) caching probed tags, versions, and positive boolean `HasExecutables`.
   - Simplified numerical / percentage progress reporting and structured `--json` payload.
   - Updating installed PowerShell (`install.ps1`, `gitmap.ps1`) and shell runner scripts.
5. **Release Ceremony**: Minor version bump and changelog update adhering strictly to relative Git paths.

## Subtask Decomposition
- Subtask 01 (`01-diffstat-and-remediation-subnode.md`): Diffstat empty-path isolation and remediation parent/sub-node tree restructuring.
- Subtask 02 (`02-aria2c-and-smart-update-engine.md`): `gitmap download` method, smart update with version targeting, 5-tag fallback, SQLite `ReleaseCache` schema, and simplified `--json` progress.

## References
- Spec: `02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/01-architecture-spec.md`
- Component Spec: `02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/02-download-and-update-spec.md`
