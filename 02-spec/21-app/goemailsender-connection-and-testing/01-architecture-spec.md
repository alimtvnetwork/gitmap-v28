# GoEmailSender Connection, Multi-Profile Fallback & Fleet Delegation Architecture Specification

## User Request (Verbatim)

```text
Can you please verify and check the connection with the GoEmailSender? You can check in my repository using the Gitmap. It is GoEmailSender. Install that, and also clone that repository in the default work directory using Gitmap. Have the default email set up for the Gitmap. Gitmap should have this email sending capability along with the scheduling. I want you to test most of the commands using or creating a repository, creating some flat commit one. Also, for the multi login purpose, how do you handle this? I wanted to know your answer. And we can always have a default login. Let's say, profile that would be by default be used, but also we can give any other profile name, deciding path to change to or set to the default, and also we can provide it to the cloning time. Or also, let's say if we are trying to clone, and if we do not provide a profile, by default, it will first try to clone from the main profile. If it does not find in the main profile, it will also try into the other profiles parallely after main profile fails. If nothing is found, then and only then it says that it fail. So remember that, that type of code needs to be there. It should be for the clone, reclone, CFR, CFRP, all of these. Also the reclone as well. So email sending, I want you to test it. You get the credential and everything from the repo secrets. I think the passwords will be only missing for the email to connect. I will give you the password for the Gitmap to connect. So you will ask for the password, I will give you the password to the Gitmap to connect as a email thing. And also Gitmap can provide as a UI for the email sending, scheduling an email, things like that. Confirm that you understood the process and everything. If you have any question and concern, let me know. Also, other commands that you have recently created, I want you to do the end-to-end testing with proper test cases. So you write these test cases, but these test cases will not run in the CI/CD or in the local in the future, but only if the request is made. Do you understand? Also check the recent commands like the heatmap and the agent. All of these should actually reveal the database path as well. Absolute database path and also a relative database path, depending on where the database is. If the database is in the AI memory, then it would share the relative path. If it is in the other location using RC, then it would try to have the absolute path. Oh, I think rather than having the absolute path, it could actually create a relative path by going back from the repository and where the repo secrets are. And also put the repo secrets URL into the JSON output so that it knows where the database is or coming from. Also, the recent commands, let's say recent commands that we have added, make sure that those are also available through the nodes. For SSH machines. And also the delegation can happen with another Gitmap, so another Gitmap can pick commands as a delegation and execute and revert back as a response. I hope you understand all this and try to do this end-to-end testing and make sure the CI/CD is green as well at the end. And do a minor bump if you made any changes and make a release.
```

## 1. Architectural Intent & Systems Overview

This subsystem enhances GitMap with four interconnected operational pillars:
1. **GoEmailSender (`go-email-reader`) Integration**:
   - Connection resolution and cross-platform binary invocation (`.exe` Windows support).
   - Default email profile setup in GitMap (`gitmap space send-email` and `gitmap space schedule-email`).
   - Terminal interactive UI (`gitmap space email ui`) for viewing account state, drafting, sending, and scheduling email dispatches.
   - Credentials offloading and zero-plaintext security: passwords are NEVER stored in git repositories; connection metadata is pulled from `repo-secrets/11-go-email-reader/` and passwords are provided dynamically or securely prompted.
2. **Multi-Login Profile Fallback Strategy (Default First, Parallel Secondary Probing)**:
   - Configurable default GitHub profile (`alimtvnetwork` or specified alias).
   - Dynamic profile switching and explicit profile override at clone time (`--profile <name>`).
   - Fallback hierarchy for `gitmap clone`, `reclone`, `cfr` (clone-fix-repo), and `cfrp`:
     * Step 1: Probe default profile first.
     * Step 2: On failure (HTTP 401/404, repository not found, unauthorized), trigger concurrent parallel probing across all secondary profiles in `02-git-profiles.json`.
     * Step 3: Only when all profiles report failure does the operation abort with a unified diagnostic error.
3. **Database Path Exposure & Relative Traversal Hygiene**:
   - `gitmap heatmap` and `gitmap agent task/db` outputs must avoid raw absolute paths where possible.
   - For databases residing in `.ai-memory/`, output relative paths (e.g. `.ai-memory/heatmap.db`).
   - For external or RC databases, calculate relative traversal back from repository root towards `repo-secrets` (e.g. `../repo-secrets/...`) and inject `repoSecretsUrl` into JSON envelopes.
4. **Fleet Nodes & Peer GitMap Delegation**:
   - Expose recent commands (`heatmap`, `agent task`, `space send-email`, `space schedule-email`) through `gitmap nodes run` across remote SSH nodes.
   - Enable peer delegation protocol where one GitMap node dispatches a command payload to another node, monitors remote completion, and collects structured JSON responses.
5. **Isolated End-to-End Tests**:
   - New tests are strictly isolated via environment guards (`GITMAP_RUN_E2E=1`) or build tags (`//go:build e2e`) so routine `go test ./...` and CI/CD pipelines never execute slow external commands or network operations.
