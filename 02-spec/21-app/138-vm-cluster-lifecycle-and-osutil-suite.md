# Spec 138: VM Cluster Lifecycle, E2E Verification & OS Utility Suite

> **/goal** Provide robust cluster node lifecycle management (join, remove, list), multi-machine execution (`mm`), remote GitMap/AGM upgrades, Windows/Linux OS utility feature exposure (Auto-Login, Scheduled Tasks), and comprehensive local-only E2E tests without exposing secrets or regressing CI/CD.
> **/learn** VM cluster tests must be guarded under `//go:build e2e` and strictly use local `vmpass.json` (gitignored). Errors must be structured via `*appfault.AppError` with stack traces, never swallowed.

## User Request (Verbatim)

```text
is it done properly , please confirm the release and ci cd fix please??


https://prnt.sc/2wxb74N-bomk

AI Main	192.168.1.20	off
AI VM 01	192.168.1.3	w1
AI VM 02	192.168.1.7	w2
AI VM 03	192.168.1.12	w3
AI VM 04	192.168.1.13	w4
AI Ubuntu 01	192.168.1.22	u1

vmpass.json
D:\work\chris\winutil
D:\work\chris\linutil


Okay. So these are the VMs. I am going to open all the VM. All the VMs are running. The reason I'm sharing all this with you is that so that you could test locally, and you don't need to run this test on CI/CD. So you can mark this on local only and VM. If the user only requests, then and only then it will run. It's not like it will run every time, but only the test case, that's a special case, it will run again if we want it to. So remember, this is how you are going to write the end-to-end test, and I'm going to give you all the IPs, all the machine aliases. So W1, W2, these are the aliasing, so you could have the aliasing. So here, the AI main, the main machine, which is the IP 20, which is off. Okay? So you can mark this as off. So we have AI VM one, two, three. So these three are running. You can keep the IPs for this machine to somewhere as a JSON. That's all right. That would not be a problem. However, you need the password, right? You need the password to log in with SSH. So all of these has SSH enabled. I can log in using the password. I have also shared the password with you, but you should never commit the password or things to the Git, or let's say, expose the password. Okay? So for all the Windows machine, the username and the password is same. Username is administrator, the password is given in the JSON file. Okay? You can read, you can understand, but you should not commit. This file is marked as gitignore, so this should never be committed by me. So you just be aware of it. Okay, now the next steps. You have SSH join, you have clustering commands, lots of commands, which I tried several different ways after releasing. It does not work. Okay? Now I'm giving you all the access. You will write end-to-end local tests, okay, that might need the VM path if required to, or you can enhance that VM path using the IPs and other stuff to a formalized section like clustering we have seen before. You can do that, modify that, and you can use that schema. The way that this schema will work is that it's going to connect to SSH, first create the SSH join, all of these from the current machine, and then also it would send SSH IP command to see all of these machines' IP. It would also try to update Git maps from the current version to a latest version. It will try to execute some commands from Git map, try to execute commands in Ubuntu to install something and get the results here. And also we could do a JSON to get the result at the end as a JSON format to verify. So all of these, I want that you verify, you take your time, and you finally release the code as a final thing so that it works everywhere. Okay, I also want you to check the AGY commands. The Ubuntu does not work yet. It has the old version of the Git map. So I want you to update that automatically using our commands. And whatever you are doing at the end, I want a summary like, "This command works. This command does not work, and I fixed it. Now, this is the way you could prompt AGY, even from multiple machines, multiple VMs," things like that. Okay, I want this to be happening. And also the Ubuntu has Antigravity Manager, very old version. I want you to update that as well. During the process, you might encounter lots of error. I want to have all these errors as a markdown, so that we can fix it later on if, let's say, Antigravity Manager is having the issue. So I want that error so that I can pass it on and fix it. Try to understand that. And all these VMs, I want to run Git map pull all in the future, okay, with very easy functionalities. So these are the things I want things to be done. I can see the results in the terminal for the multiple machines and VMs. I should be able to remove, join node, add node. So these type of tests, I want to see that you perform using end-to-end tests, that it works, this final summary I want. And if there is any issues, you fix it immediately and make sure that the code actually has the app error and a stack trace so that you understand where it's coming from. And also the error needs to have full verbose thing. Let's say, Go gave you an error and you swallowed that error. It should never happen. You should fix that criteria. Okay? The error needs to be verbose. We need to see what is going on, and then we can decide, is it clear? Okay. So yeah, and make sure the logs does not have the password or any sort. Okay, be aware of it. Okay. Now, based on these factors, I want you to start working. And also at the end, I requested in the past to have a auto login feature. I'm not sure if you have integrated and added that. Please confirm that it is done. If not, please work on it and complete that. Auto scheduling feature, I asked you to learn from the WinUtil from Chris. I'm not sure if you have the repository. Let me just check one more time. Yes, we do have the Chris repositories and WinUtil. And in the past, I also asked you to improve some commands or add few commands for the Linux Util. I'm not sure if you have done it. So that is another part. So lots of things. So first plan all this stuff and then start working. Is it clear?

And finally make the release bump and test ci cd if it is working or not
```

## System Architecture & Deliverables

### 1. Cluster Node Lifecycle & Local E2E Test Suite
- **E2E Tag**: `//go:build e2e` strictly applied to `cli/tests/vm_cluster_e2e_test.go`.
- **Node Lifecycle Verification**:
  - Test node registration (`gitmap sj add-with-pass` or cluster join API).
  - Test node listing (`gitmap cluster ls`).
  - Test node removal (`gitmap cluster rm <node>` / `gitmap cluster remove <node>`).
  - Re-join and verify ready state.
- **Multi-Node Commands**:
  - `gitmap mm "hostname"`
  - `gitmap ssh <node> "<cmd>"` with `--json` output
  - Verify graceful offline skip of `w4` (192.168.1.13) and `AI Main` (192.168.1.20).

### 2. OS Utility Top-Level Command Aliases
Direct command exposure in `cli/cmd/rootutility.go`:
- `gitmap winutil <args...>` -> dispatches to `cmdos.RunOSCLI`
- `gitmap linutil <args...>` -> dispatches to `cmdos.RunOSCLI`
- `gitmap autologin <args...>` -> dispatches to `cmdos.RunOSCLI("autologin", ...)`
- `gitmap schedule <args...>` -> dispatches to `cmdos.RunOSCLI("schedule", ...)`

### 3. Remote Updating & Issue Root Cause Analysis
- Remote updates for GitMap across Linux and Windows nodes via `gitmap remote update --node <node> gitmap`.
- Document Ubuntu Antigravity IDE and Antigravity Manager issues in `02-spec/22-app-issues/14-ubuntu-agy-agm-execution-errors.md`.

## Quality Gates & Invariants
- Zero secrets committed: `vmpass.json` strictly gitignored and excluded from all git operations.
- Anti-nesting: Max conditional depth 1 across all Go code.
- Boolean rules: Only `is` and `has` prefixes.
- Zero swallowed errors: Wrap all errors with `*appfault.AppError` and preserve stack traces.
