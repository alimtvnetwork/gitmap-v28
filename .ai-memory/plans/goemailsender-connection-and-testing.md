# Execution Plan: goemailsender-connection-and-testing

## User Request (Verbatim)

```text
Can you please verify and check the connection with the GoEmailSender? You can check in my repository using the Gitmap. It is GoEmailSender. Install that, and also clone that repository in the default work directory using Gitmap. Have the default email set up for the Gitmap. Gitmap should have this email sending capability along with the scheduling. I want you to test most of the commands using or creating a repository, creating some flat commit one. Also, for the multi login purpose, how do you handle this? I wanted to know your answer. And we can always have a default login. Let's say, profile that would be by default be used, but also we can give any other profile name, deciding path to change to or set to the default, and also we can provide it to the cloning time. Or also, let's say if we are trying to clone, and if we do not provide a profile, by default, it will first try to clone from the main profile. If it does not find in the main profile, it will also try into the other profiles parallelly after main profile fails. If nothing is found, then and only then it says that it fail. So remember that, that type of code needs to be there. It should be for the clone, reclone, CFR, CFRP, all of these. Also the reclone as well. So email sending, I want you to test it. You get the credential and everything from the repo secrets. I think the passwords will be only missing for the email to connect. I will give you the password for the Gitmap to connect. So you will ask for the password, I will give you the password to the Gitmap to connect as a email thing. And also Gitmap can provide as a UI for the email sending, scheduling an email, things like that. Confirm that you understood the process and everything. If you have any question and concern, let me know. Also, other commands that you have recently created, I want you to do the end-to-end testing with proper test cases. So you write these test cases, but these test cases will not run in the CI/CD or in the local in the future, but only if the request is made. Do you understand? Also check the recent commands like the heatmap and the agent. All of these should actually reveal the database path as well. Absolute database path and also a relative database path, depending on where the database is. If the database is in the AI memory, then it would share the relative path. If it is in the other location using RC, then it would try to have the absolute path. Oh, I think rather than having the absolute path, it could actually create a relative path by going back from the repository and where the repo secrets are. And also put the repo secrets URL into the JSON output so that it knows where the database is or coming from. Also, the recent commands, let's say recent commands that we have added, make sure that those are also available through the nodes. For SSH machines. And also the delegation can happen with another Gitmap, so another Gitmap can pick commands as a delegation and execute and revert back as a response. I hope you understand all this and try to do this end-to-end testing and make sure the CI/CD is green as well at the end. And do a minor bump if you made any changes and make a release.
```

## System Concurrency & Wave Budget
- Concurrency Capacity: A = 3 agents × H = 2 hands = 6 concurrent operations
- Total Steps Budget: N = 400
- Phase 1: Planning, Discovery & Spec Authoring (Steps 1–100)
- Phase 2: Execution Waves (Steps 101–250)
- Phase 3: Targeted Verification, E2E Isolation, Release Ceremony (Steps 251–400)

## Subtasks Allocation

- **Task-01: GoEmailSender Discovery, Windows Path Resolution & Installation Validation**
  - Files: `cli/cmdspace/send_email.go`, `cli/cmdspace/send_email_windows.go`
  - Owner: Worker 01

- **Task-02: Multi-Login Profile Fallback (Default First, Parallel Secondary)**
  - Files: `cli/cmdclone/clone_auth_resolve.go`, `cli/cmdclone/gh_resolver.go`, `cli/cmdclone/clone.go`, `cli/usercontext/project_binder.go`
  - Owner: Worker 02

- **Task-03: Email Sending, Scheduling & Terminal UI**
  - Files: `cli/cmdspace/send_email_ui.go`, `cli/cmdspace/schedule_send_email.go`
  - Owner: Worker 03

- **Task-04: Relative DB Paths & Repo-Secrets URL Exposure in Heatmap and Agent**
  - Files: `cli/heatmap/heatmap_view.go`, `cli/cmdagent/agent_task_status.go`, `cli/heatmap/heatmap_rc.go`
  - Owner: Worker 01

- **Task-05: Fleet Nodes SSH Parity & Peer GitMap Delegation**
  - Files: `cli/cmdnodes/nodes_remote.go`, `cli/cluster/exec_delegate.go`
  - Owner: Worker 02

- **Task-06: Isolated E2E Test Suite & Minor Version Bump Release**
  - Files: `cli/tests/e2e/goemailsender_tempe2e_test.go`, `version.json`
  - Owner: Lead Orchestrator

- **Task-07: SSH Fleet Pull UX & Local Node Detection**
  - Files: `cli/cmdssh/ssh_pull_fleet.go`, `cli/cmdssh/ssh_pull_fleet_heartbeat.go`, `cli/cmdssh/ssh_pull_fleet_liveness.go`, `cli/cmdssh/ssh_pull_fleet_dispatch.go`, `cli/cmdssh/ssh_pull_fleet_outcome.go`, `.github/workflows/ci.yml`
  - Owner: Lead Orchestrator
