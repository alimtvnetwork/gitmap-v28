# Specification: 171 — SSH Fleet GitMap Auto-Update Guard & Concise Pull Remediation

## Status
`active`

## Domain
`cli / ssh fleet / pull / auto-update / json protocol`

## User Request (Verbatim)
```text
The first of all, if you run the Git map PA in a local directory, it will only show at the end whichever that actually has issues, failed, and the updated ones. Other ones, just ignore from the last summary. That's the first thing. Okay. Now, coming to the point, when we do the SSH, I mentioned this before, the communication needs to be done using JSON. Okay? The first thing it'll do is that definitely it would check that the Git map is installed and up to date. If not, then it would do that. So it did nicely how it represents the Git map SSH execute IP, just like how we do it. The way that it shows the commands that has been executed, the similar way it would do it here. But first thing it would say, "Yes, here we have enqueued the task into this other worker," and the current VM, it is going to run itself to be here. And the current VM is the one that it is this. So we are not sending the request, we're running it here. It mentions this. So all these IPs would show up where the request has been sent first. Okay? And then it would show the current one is also running. So once those async requests send it to those IPs, and using SSH would come back as a response. When it comes back on the async process, then we combine it and finally give a nice output, just like the current one from the JSON that we have received to the terminal. Okay, so the current response that you have done, it's trashy, and I'm not sure why there is this GitHub origin URL that we have to add, why we have to add this and display it everywhere. I don't know, so correct me if I'm wrong. Okay, so make sure that you do not do this type of output in the SSH. Whenever the SSH is mentioned. So all of these like update all or Git map AGM update hyphen hyphen SSH, then it indirectly send the request through these machines as a SSH, but it would receive the response as JSON. Remember that. And when it receives the JSON, it will just update what is the status of it nicely. Do you understand? Can you please help me with this?

Release minor and check gitmap pe please
```

## Architectural Context & Blast Radius
1. **Remote Node Auto-Update Guard:**
   - In `cli/cmdssh/ssh_pull_fleet.go`, `executeRemoteNodePull` connects to each remote SSH node.
   - Nodes running outdated GitMap binaries (prior to `v6.349.0`) lack `--json` support on batch pull, resulting in error `flag provided but not defined: -json`.
   - `ensureRemoteNodeGitmap` probes the remote OS, queries the installed GitMap version using `queryNodeVersionViaSSH`, and triggers `updateTargetNodeGitmap` via `resolveRemoteUpdateCommand(osType, "gitmap")` if outdated (< 6.349.0) or uninstalled.
   - `runRemotePullJSON` inspects the execution output: if `flag provided but not defined: -json` is detected, it triggers remote update and retries once automatically.
2. **Concise Pull Remediation Gate:**
   - In `cli/cmdpull/pull.go`, `handlePullRemediation` is guarded: when `opts.all && !opts.showStatus` (standard fast mode), interactive remediation prompts (`Remediate dirty repository(ies) now? [y/N]:`) are bypassed entirely, preventing process blocking on unattended and automated runs.
3. **Transport Rewrite Elimination on Batch Operations:**
   - In `cli/cmdpull/pull_efficient.go`, `maybeApplyTransportToRecords` is removed from `executeActiveEfficientBatch`. Batch pull operations never rewrite `remote.origin.url` or emit 64 lines of transport URL spam.

## Verification Gates & Invariants
- `check-nested-ifs.py`: 0 nested if violations.
- `check-relative-paths.py`: 0 absolute paths.
- `test_ci_scripts.py`: All unit tests pass.
- Version bump: SemVer minor release `v6.350.0`.
- All GitHub Actions CI/CD workflows complete with 100% green status.
