# Plan 192: Version Pinning, Macro Fleet Deployment, UI Settings Layout, and Secret Flags Deduplication

## User Request (Verbatim)

```text
/goal Okay, make sure that we have version pinning ability. So for any installation, we can go there and say pin version, and then we can say pin, and then installer object, and then give a version, and that version would be pinned. And when we do installer LS, we can also see the summary, like which versions are pinned. And we could do pin as a flag as well. That means when we are updating, we can say updated and a space V, specific version. It can be for gitmap, it can be for AGM. Both cases, we could do a specific version, and that update would fall to that specific version and pin to that if we use a hyphen pin flag. Same way we could run it for all the SSH machine just using an SSH flag in both cases. So also you need to confirm if we have a macro deploy command. That means we should have a deploy space macro space all that would deploy into all machines, all macros, or name of a macro, and then name of a machine. So both would be correct. Also we can have except flag that would actually take ID, aliasing IP, et cetera, to find and sequence the IP to not apply that. So these are very important factors, and we can always export the macros as many as we want and can import as well. So make sure that these commands are there. Okay. Also, you can confirm the UI section for the graphics and opening and changing the UI from the UI settings and everything. Do we have all the options? Do we have the setting option for the commit in, commit pull in, pull left, right, et cetera, all of these. And also you go to repo secrets, see the gitmap folder. There are this commit in pull JSON file, okay? That has same type of flags repeated multiple times. I hope you fix that and keep only one flag. For example, ECD apply, not the ECD. Okay? Something like this. So I think these are repeated in many places. So try to fix it in your code, fix it in the JSON, and please do a git pull on the repos before you do that, change that, and then make a commit. At the end, please make a minor bump, and also make sure that you check the CICD that it is green. Do you understand the task? And also at the end, include the execute parent task in NS steps. Okay? Can you please do that?
```

## Actionable Deliverables & Task Breakdown

- **Task-01: Version Pinning & Status in Installer Subsystem**
  - Implement `gitmap installer pin <target> <version>` and `gitmap pin <target> <version>` (and `gitmap pin installer <target> <version>`).
  - Store pinned version in SQLite Split-DB (`installation.db` / installation records).
  - Update `gitmap installer ls` / `installer list` output to show pinned version in summary and details.

- **Task-02: Target Version Update with `--pin` and Fleet `--ssh` Flag**
  - Support `gitmap update -v <version> [--pin]` / `gitmap update --version <version> [--pin]` for both `gitmap` and `agy` / `agm`.
  - When `--pin` is provided, persist the pinned version upon update.
  - Support fleet execution across all SSH machines using `--ssh` (e.g. `gitmap update -v <version> --pin --ssh`).

- **Task-03: Macro Fleet Deployment Commands (`deploy macro`, `macro deploy`)**
  - Implement `gitmap deploy macro all` and `gitmap macro deploy all` to deploy all macros to all SSH nodes.
  - Implement `gitmap deploy macro <macro-name> <node>` and `gitmap macro deploy <macro-name> <node>`.
  - Support `--except` flag to exclude specific nodes by ID, alias, or IP address.
  - Verify macro export (`gitmap macro export`) and import (`gitmap macro import`) commands with multi-macro support.

- **Task-04: UI Section Verification & Layout Settings Configuration**
  - Verify UI configuration for graphics, opening, and modifying UI from UI settings.
  - Verify and ensure settings for commit in, commit pull in, pull left, right, layout orientation, etc.

- **Task-05: Repo Secrets GitMap JSON Cleanup & Deduplication**
  - Git pull latest changes in `D:\work\repo-secrets`.
  - Locate `01-gitmap/` commit/pull JSON configuration files.
  - Deduplicate repeating flags (e.g. keeping `ECD apply`, removing duplicate `ECD`).
  - Align Go code parsing if needed, commit, and push changes in repo-secrets.

- **Task-06: Verification, AI/CI Prompt, Minor Bump & Green CI/CD**
  - Author CI/AI verification prompt in `01-prompts/`.
  - Run workstation cache and temp cleaner (`03-ai-scripts/42-clean-test-and-build-caches.py`).
  - Bump minor version and orchestrate release.
  - Verify GitHub Actions CI/CD workflows are 100% green.
