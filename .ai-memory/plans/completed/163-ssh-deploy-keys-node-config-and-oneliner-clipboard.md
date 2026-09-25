# Consolidated Completed Plan 163: SSH Fleet Deploy Keys, Node-Config, and One-Liner Clipboard

Canonical Specification: [02-spec/21-app/163-ssh-deploy-keys-node-config-and-oneliner-clipboard.md](../../../02-spec/21-app/163-ssh-deploy-keys-node-config-and-oneliner-clipboard.md)
Execution Summary: Completed in 1 continuous loop (4 discrete subtasks consolidated).

## User Request (Verbatim)

```text
gitmap ssh deploy keys all [--except id, ip, aliasing]
gitmap ssh deploy node-config all [--except id, ip, aliasing]

gitmap ssh export-oneliner (eo) # should export a single line not too many seperate ones, please fix it, and you use the copy method to copy in memroy as well and inform the user about it? ckear????

There is a error from your one liner export SSH export one liner. Also, I do not know if you have a deploy key for SSH deploy keys or something like this, which is basically going to deploy the SSH keys, the public keys from the current machine to all these machines. And also, all these machines should know about the node information. So seemingly, we could have a specific deploy node config, something like this. Please tell me if you have any one of these. We could do accept if we like. This is also an optional section where we could do IP, IP, or aliasing. Usually mention that keys should first should deploy the node config. That means all the node aliasing, IP, and things. And then on those machines, the current nodes public keys automatically will be added as a trusted key. So it will happen for all these machines. The way it will work, let's say I have three nodes. Okay? So the SSH will try to log into each one of them and try to get the SSH public key. Once it gets the public key of these three machines, then it will check if those public keys are same. Try to see which one is unique. So it will do the unique, and then the only unique keys will be added to all those machines as a authorized key, SSH authorized key, so that the password is not needed. Do you understand? Can you please do that, and confirm the implementation, verify the implementation, and follow the coding quality, and also at the end, check git map space PE. Do you understand? Can you please validate this?
```

## Consolidated Subtasks & Outcomes

### Subtask 01 (Task-01): Refactor Export One-Liner and Clipboard Copy
- Extracted and modularized one-liner export into `cli/cmdssh/ssh_oneliner.go`.
- Emits a clean, cross-platform single-line command: `gitmap ssh nodes import-json --base64 "<b64>"`.
- Automatically copies the command to the system clipboard using `atotto/clipboard` in memory.
- Prints confirmation to the user: `📋 Copied single-line import command to clipboard!`.
- Registered `eo` alias in `ssh.go` and `ssh_ls_cmd.go`.

### Subtask 02 (Task-02): Enhance Deploy Node-Config and Ordering Guidance
- Implemented `gitmap ssh deploy node-config [all] [--except <id,ip,alias>]` in `cli/cmdssh/ssh_deploy_node_config.go` and `cli/cmdssh/ssh_deploy_node_config_render.go`.
- Added advice note recommending `deploy node-config` run prior to `deploy keys`.
- Filters target fleet nodes by ID, worker ID, IP, or alias.

### Subtask 03 (Task-03): Implement Mesh Public Key Gathering, Deduplication, and Deployment
- Implemented `gitmap ssh deploy keys [all] [--except <id,ip,alias>]` across:
  - `cli/cmdssh/ssh_deploy_keys.go`: Orchestration and execution workflow.
  - `cli/cmdssh/ssh_deploy_keys_gather.go`: Gathering local `~/.ssh/*.pub` and remote public keys, with key normalization and deduplication into unique keys.
  - `cli/cmdssh/ssh_deploy_keys_apply.go`: Applying missing unique public keys into `~/.ssh/authorized_keys` with strict directory/file permissions (`chmod 700` and `600`).
  - `cli/cmdssh/ssh_deploy_keys_render.go`: Clean terminal reporting table and summary.
  - `cli/cmdssh/ssh_deploy_keys_types.go`: Typed models for key sync results.
  - `cli/cmdssh/ssh_deploy_keys_test.go`: Unit tests for key deduplication, signature extraction, and flag parsing.

### Subtask 04 (Task-04): CLI Dispatch, Help Documentation, Linters, and Pipeline Verification
- Created `cli/cmdssh/ssh_deploy_router.go` routing `gitmap ssh deploy` to `keys` and `node-config`.
- Updated `cli/cmdssh/ssh.go` to handle `eo`, `deploy`, `deploy-keys` (`dk`), and `deploy-node-config` (`nc`).
- Updated `cli/cmdssh/ssh_help_sections.go` and `cli/helptext/ssh.md`.
- Verified all Go files <= 100 lines, zero nested ifs, and positive booleans.
- Verified remote CI pipeline status using `gitmap pe` (100% green).
