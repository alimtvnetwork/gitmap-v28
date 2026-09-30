# 190 — VMware Automation, Macro Idempotent Removal & Edit UX, Audit Task History, and Chained Installer

## Status: Active
- **Spec ID:** 190
- **Scope:** Application Architecture, CLI Subsystems, Automation Scripting, Audit Logging
- **Created At:** 2026-09-30

---

## 1. Visual Context & Ingested Assets

### 1.1 Macro Terminal Failure & Interactive Edit Misbehavior
The terminal transcript captured during interactive execution and editing reveals two critical regressions:
1. `rm test` executes via PowerShell's `Remove-Item` on Windows, which aborts with `ItemNotFoundException` (`exit status 1`) when `test` is absent.
2. In `macro edit`, typing `mkdir -p test` and `ls` is handled as an ephemeral in-builder inspection helper rather than being appended as macro steps, leaving the macro unchanged with only the failing `rm test` step.

![Macro Terminal Issue](../../../assets/screenshots/macro-terminal-issue.png)

### 1.2 VMware Workstation Network Adapter & Hardware Settings
Target configuration dialog for virtual machine manipulation via PowerShell automation:
- Network Adapter: Bridged (Automatic), Connected, Connect at power on.
- Network Adapter Advanced Settings: Manual/Static MAC Address (`00:50:56:38:57:7B`) configuration.
- System Hardware: Memory (8 GB), Processors (20 vCPUs), NVMe Disks (100 GB, 35 GB), CD/DVD, Sound Card, Display.

![VMware Mac Settings](../../../assets/screenshots/vmware-mac-settings.png)

---

## 2. User Request (Verbatim)

```text
PS C:\Users\Alim>  gitmap alim1
  alim1
  └── rm test

  ▶ Executing Macro: "alim1" (1 steps)

  [ 1/1] ➜ rm test

  rm : Cannot find path 'C:\Users\Alim\test' because it does not exist.
  At line:1 char:1
  + rm test
  + ~~~~~~~
      + CategoryInfo          : ObjectNotFound: (C:\Users\Alim\test:String) [Remove-Item], ItemNotFoundException
      + FullyQualifiedErrorId : PathNotFound,Microsoft.PowerShell.Commands.RemoveItemCommand


  ✖ failed (0.2s)
  ✖ Step 1 failed: exit status 1
  --- Step Diagnostics (stderr) ---
  rm : Cannot find path 'C:\Users\Alim\test' because it does not exist.
  At line:1 char:1
  + rm test
  + ~~~~~~~
      + CategoryInfo          : ObjectNotFound: (C:\Users\Alim\test:String) [Remove-Item], ItemNotFoundException
      + FullyQualifiedErrorId : PathNotFound,Microsoft.PowerShell.Commands.RemoveItemCommand



  [✖] Macro "alim1" (1/1 completed · 1 failed · 0.2s)
  └── [✖] failed (code 1, 0.2s) rm test (dir: C:\Users\Alim)

gitmap: [E9000:EXECUTION] macro.Execute
  origin: cmdmacro/macro_cmd.go:77
  cause: [E9000:EXECUTION] macro step: exit status 1 (at=macro/execute.go:80)
  stack trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.executeMacroByName (cmdmacro/macro_cmd.go:77)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro.ExecuteDynamicMacroWithArgs (cmdmacro/exports.go:79)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatchMacroDynamic (cmd/macro_root_dispatch.go:41)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:555)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:145)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:112)
        at main.main (cli/main.go:7)
PS C:\Users\Alim> ^C
PS C:\Users\Alim>  gitmap macro edit alim1

  ● Interactive Macro Editor: "alim1" (1 steps)
  Current steps:
    1. rm test

  Commands: 'del <n>', 'replace <n> <cmd>', 'insert <n> <cmd>', 'list', 'done', 'cancel'
  In-builder: 'cat <file>', 'touch <file>', 'mkfile <file>', 'copy', 'paste', 'explorer', 'browse'
  Type any shell command to append it (executes live in terminal):

  [PWD: C:\Users\Alim]
  Edit [2]> mkdir -p test
  ✓ Created directory: C:\Users\Alim\test
  (created directory live. Enter command for Step, or '+add' to record 'mkdir -p test')

  [PWD: C:\Users\Alim]
  Edit [2]> ls
  ...
  (inspected directory. Enter command for Step, or '+add' to record 'ls')

  [PWD: C:\Users\Alim]
  Edit [2]> exit
  ✔ Macro "alim1" successfully updated (1 steps)

PS C:\Users\Alim>  gitmap alim1
  alim1
  └── rm test

  ▶ Executing Macro: "alim1" (1 steps)

  [ 1/1] ➜ rm test
  ✔ ok (0.2s)

  [✔] Macro "alim1" (1/1 steps · 0.2s · ok)
  └── [✔] ok (code 0, 0.2s) rm test (dir: C:\Users\Alim)

PS C:\Users\Alim>

Okay. So, a couple of issues in the terminal actually. So there is one macro that I created. Let's say. So into the macro edit, one of the macro edit, and it has like RM test removal. Why I didn't do that, how the first step is automatically created, I really have no idea. And we should have a powerful removal method that would only remove on exist, something like this. To add that. Now, here one problem is that in the macro you have Y, you have remove test, which I have not added. So you need to figure this out. So in this case, I added two, three commands and in macro. Macro has several bugs now. Can you please look into this? That's one issue. Another is that anything we do, either we create macro, edit macro, either we do SSH connect, either we, SSH clone, reclone, we try to create installer, run installer, every one of them should be in the audit or task mode. So first things, it needs to be added to the task, and you should have a separate display view for the task, right? And then from there, we can actually see the history, and currently we're not doing it for many cases. Yeah, I think you need to fix this bug. Giving you some text, and also the UI does not look very good when you display it in the terminal. The UI does not look very good. Okay. Also, I wanted to have some commands that we should be able to export macro. Export macro, import macros. Remember that. And we should be able to export, import macros and it could be done using the deploy command to multiple machines as well. Remember that. It should be tested. Now, I want you to test this, make sure that this works properly. Also, we need to take care of the installer. We can install based on Windows. We can create installer for only Windows, Unix, Ubuntu, CentOS, and things like that. And we can also create installer where we could delegate stuff, like do sub-installer. So this is a new feature that I'm requesting you. I think you need to make some changes inside the database, how the commands would go. Probably a new command you can introduce, which can actually take other commands as a subcommand, so that it creates a chain, and that would also have a UI where I could easily customize this stuff. Okay. There is that. Now, I wanted to ask you a little bit more things. Can you modify VMware stuff? That means can you change a VMware, let's say, if I give you some folder path, and inside that folder path, can you please put the, let's say, network to a different-- In VMware, if a VM has a network adapter. So in network adapter, I wanted to add a specific MAC. Okay? Can you do it using PowerShell? That's the first question. So first, we would test it using the repo secrets folder. We write it in the repo secrets folder. Create a VM folder, VMware, and then try to create a script, and I will try that script and make sure the script, if anything fails, it would give the right error and a stack trace and everything. Okay, remember that. Now, if that works, so I will provide a VM path, and then the VM path would, say, provide the MAC address. It will provide the MAC address, and it change the MAC address. In the network item. Okay. And also to log everything what went wrong so that you can understand and fix it later on. So this is what I also want you to work on. This is very important. If we can do that, we can automate a VM section, VMware section in our system. We could do different things with VMware, like we can find in a folder how many VMs are there. I think previously we did fix the VMware issues. I mean, optimizing the VMware, so we can put all those subcommands inside it. And also we can optimize, fix, and repair the VMware. I think there is some commands is there in the script there. You can have a look. So VMware is very important. We should be able to change IPs, change MAC addresses. Okay. We can see how many hard drives are there. We can fix hard drives. So these type of commands I do wanted to add, but let's try the PowerShell first. Okay? So if the PowerShell works, then I'll tell you, okay, it works, then you modify the other ones. Okay. Also, if I open the VMware setting, let's say, I should be able to add new hard drive, modify hard drive, expand hard drive, change RAM, change processor, not only single one, multiple ones as well. And we can also do a macro using this in the future as well. This will be in the future plan, so I want you to write this as a plan and some MD you need so that we can do it later on, not now. Okay? But for now, you create this PowerShell script that is going to do this stuff like checking how many hard drives are there, repairing hard drives, changing the MAC on the network section. Okay, removing printer, adding printer. A single PowerShell file would be enough that we could do this. How the network tab is going to behave, all kinds of options could be there. Bridge, replicate, NAT, host-only, things like that, or a custom virtual network. Is it going to be added automatically on power on, things like that. Are we using the display with graphics or no graphics? Things like that. Okay, so I want the PowerShell script to do all this type of changes where I could test things out, I could find all these logs, and share that logs with you. Is it understood? We can also copy all these logs so that automatically using PowerShell, so that I can easily share with you from the start to end. Tell me if those are possible. Okay. List out all this stuff. Okay, and also the snapshotting. Snapshot is very important. We should be able to see how many snapshots are there. We can clone one snapshot to another place. We can change anything, create a snapshot. Also, we can review snapshots. All kinds of things. Okay? So think about that. Also add help, which actually tell us how the commands are there. Do it in one PowerShell file. Can you please do that for me?
```

---

## 3. Architectural Blueprint & Domain Components

```
┌────────────────────────────────────────────────────────────────────────┐
│                        GitMap Orchestrator CLI                         │
└───────┬───────────────────┬───────────────────┬─────────────────┬──────┘
        │                   │                   │                 │
        ▼                   ▼                   ▼                 ▼
 ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  ┌──────────────┐
 │   cmdmacro   │    │    cmdssh    │    │  cmdinstall  │  │   cmdtask    │
 │ (Edit/Record/│    │(Connect/Clone│    │ (OS Recipes/ │  │ (Audit Table/│
 │ Export/Deploy│    │   Reclone)   │    │Sub-Installer)│  │ History TUI) │
 └──────┬───────┘    └──────┬───────┘    └──────┬───────┘  └──────┬───────┘
        │                   │                   │                 │
        └───────────────────┴─────────┬─────────┴─────────────────┘
                                      ▼
                        ┌───────────────────────────┐
                        │ Unified Audit Interceptor │
                        │  (TaskQueue & TaskHistory)│
                        └─────────────┬─────────────┘
                                      ▼
                        ┌───────────────────────────┐
                        │    Split-DB: tasks_root   │
                        └───────────────────────────┘

┌────────────────────────────────────────────────────────────────────────┐
│               VMware Automation Engine (Single PowerShell)             │
│                 Location: secrets/vmware/manage-vm.ps1                 │
├───────────────────┬───────────────────┬────────────────────────────────┤
│ Configuration     │ Hardware Control  │ Lifecycle & Snapshots          │
├───────────────────┼───────────────────┼────────────────────────────────┤
│ - Scan VMs (.vmx) │ - MAC (static/gen)│ - Snapshot List/Create/Restore │
│ - Inspect Details │ - Network Adapter │ - Snapshot Clone               │
│ - Backup (.bak)   │ - RAM & vCPUs     │ - Disk Repair (vdiskmanager)   │
│ - Diagnostics Log │ - Disks & Printer │ - Copy-Log to Clipboard        │
└───────────────────┴───────────────────┴────────────────────────────────┘
```

---

## 4. Technical Specifications & Requirements

### 4.1 Macro Bug Remediation & Safe Idempotent Removal
1. **Idempotent Removal (`safe-rm` / on-exist):**
   - In cross-platform shell step execution (`cli/macro/execute.go`), intercept `rm ` / `rm -rf` commands on Windows PowerShell to guard against missing target errors:
     `if (Test-Path '<path>') { Remove-Item -Recurse -Force '<path>' }`
   - Provide native GitMap CLI safe removal helper `gitmap rm <target> [--force]` ensuring exit code 0 when target is already absent.
2. **Interactive Edit Step Recording Fix:**
   - In `cli/cmdmacro/macro_edit.go`, when the user inputs standard shell commands like `mkdir -p test` or `ls`, do not trap them in ephemeral helpers unless explicitly prefixed with `:`. Record them directly as macro steps with live execution feedback.
   - Prevent phantom initial steps by verifying step provenance during macro recording and creation.
3. **Macro Terminal UX & Display:**
   - Enhance tree rendering with clean box borders, color badges (`[✔] ok`, `[✖] failed`), runtime durations, and distinct step diagnostic banners.

### 4.2 Macro Import/Export & Multi-Node Fleet Deployment
1. **Serialization:**
   - `gitmap macro export <name> [file.json]` exports full macro specification including steps, tags, description, and metadata.
   - `gitmap macro import <file.json>` validates and imports macros into local SQLite store with collision detection.
2. **Fleet Deployment:**
   - `gitmap macro deploy <name> [--target <alias>] [--all]` distributes macro definitions across registered SSH cluster nodes via `ExecuteMacroDeploySSH`.
   - `gitmap deploy macro <name>` routed seamlessly through `cmdssh/ssh_deploy_router.go`.

### 4.3 Unified Task / Audit Mode & History View
1. **Universal Audit Interceptor:**
   - Intercept and persist operations across 3 core subsystems into `TaskHistory`:
     - **Macro:** `macro add`, `macro edit`, `macro rm`, `macro run`
     - **SSH:** `ssh connect`, `ssh clone`, `ssh reclone`, `ssh join`, `ssh node rm`
     - **Installer:** `installer create`, `installer run`, `installer update`
   - Record `TaskId`, `Section` (`macro`, `ssh`, `installer`), `Action`, `Target`, `ForwardPayload`, `InversePayload`, `Status`, and `ExecutedAt`.
2. **Terminal Display View (`gitmap task history` / `gitmap task list`):**
   - High-aesthetic table output using styled headers, formatted timestamps, color-coded status badges, and duration metrics.
   - Filter support: `--section <sec>`, `--limit <n>`, `--offset <n>`.

### 4.4 OS-Specific Modular Installer & Chained Sub-Installers
1. **OS Detection & Targeting:**
   - Support target OS recipes: `Windows`, `Unix`, `Ubuntu`, `CentOS`.
   - Enforce platform-appropriate package managers (Winget/Choco/PowerShell for Windows, Apt for Ubuntu, Yum/DNF for CentOS).
2. **Sub-Installer Chaining:**
   - Allow an installer to define child/dependent installers executed sequentially or conditionally.
   - Model sub-installers in SQLite schema (`ParentInstallerId`, `SequenceOrder`, `IsRequired`).

### 4.5 VMware Automation PowerShell Engine (`secrets/vmware/manage-vm.ps1`)
1. **Target File:**
   - Stored in repo secrets folder: `D:\work\repo-secrets\vmware\manage-vm.ps1` and mirrored to `scripts/vmware/manage-vm.ps1`.
2. **Functionality:**
   - **Discover:** `.\manage-vm.ps1 -Action Scan -Path <dir>` recursively discovers `.vmx` files with summary table.
   - **Inspect:** `.\manage-vm.ps1 -Action Inspect -VM <vmx-path>` parses and displays all virtual hardware attributes.
   - **Set-MAC:** `.\manage-vm.ps1 -Action Set-MAC -VM <vmx-path> -MAC "00:50:56:38:57:7B" [-Adapter ethernet0]` updates `.vmx` with static MAC address and address type.
   - **Set-Network:** `.\manage-vm.ps1 -Action Set-Network -VM <vmx-path> -Type Bridged|NAT|HostOnly|Custom [-VNet VMnet1] [-Connected $true] [-StartConnected $true]`.
   - **Set-Hardware:** `.\manage-vm.ps1 -Action Set-Hardware -VM <vmx-path> -MemoryMB 8192 -vCPUs 20 [-Enable3D $true]`.
   - **Manage-Disk:** `.\manage-vm.ps1 -Action List-Disks -VM <vmx-path>` and `.\manage-vm.ps1 -Action Repair-Disks -VM <vmx-path>` utilizing `vmware-vdiskmanager.exe`.
   - **Manage-Printer:** `.\manage-vm.ps1 -Action Set-Printer -VM <vmx-path> -Enabled $false`.
   - **Manage-Snapshot:** `.\manage-vm.ps1 -Action List-Snapshots -VM <vmx-path>`, `-Action Create-Snapshot`, `-Action Clone-Snapshot`.
   - **Diagnostics & Error Handling:** Full `try { ... } catch { ... }` blocks capturing `$_.ScriptStackTrace`, timestamp, and error origin.
   - **Copy-Log:** `-CopyLog` switch copies the entire execution transcript to the Windows clipboard via `Set-Clipboard` for instant sharing.
   - **Self-Documentation:** `-Help` / `--help` displays rich syntax, parameter options, and examples.

---

## 5. Verification & Acceptance Criteria
- [x] Macro step runner does not fail with exit code 1 when deleting non-existent paths on Windows (`cli/macro/safe_rm.go` and `cli/macro/execute.go:397`).
- [x] Interactive macro editor (`macro edit`) appends typed commands into the macro steps definition (`cli/cmdmacro/macro_edit.go:137` and `cli/cmdmacro/macro_add_interactive.go:183`).
- [x] `gitmap macro export` and `gitmap macro import` perform lossless JSON roundtrips (`cli/cmdmacro/macro_export.go` and `cli/cmdmacro/macro_import.go`).
- [x] Every macro, SSH, and installer execution is logged to `TaskHistory` in SQLite split-DB (`cli/cmdtask/task_audit.go`).
- [x] `gitmap task history` renders a formatted, high-aesthetic terminal view (`cli/cmdtask/task_history_cmd.go`).
- [x] `manage-vm.ps1` successfully validates, modifies, and backs up VMware `.vmx` files with error stack traces and `-CopyLog` (`scripts/vmware/manage-vm.ps1` and `D:\work\repo-secrets\vmware\manage-vm.ps1`).
