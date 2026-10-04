# Completed Plan: 72-chrome-ext-test-and-vmware-spec

## User Request (Verbatim)
# High Priority Instruction

Okay, I think this is a time that you need to test the Chrome extension. What do I mean by that? We have the Chrome import-export feature. Okay? So what you will do is you will try to export one profile, just take any one sample profile, take it as a test case, and if you need to write something secret, like the email or something secret, you use the repo secrets folder. Okay? Inside this, you create the git map folder with the sequence, and then inside this, you create the task name and with the sequence, by the way, and then you try to work on it. Okay? So you can export there and try to re-import by removing the profile. Try to see if that is successful. Pick one profile and start with that. If that is not successful, then pick another one, export, and every step, you try to have a log so that you understand by looking into the log what went wrong. Okay? So you try to go inside the log, so your logs will be saved inside the repo secrets folder by the task. Okay? And by looking into or inspecting the log or error tracing, you would know this is wrong. And also, I want you to first fix the git map be first. Okay? Try to make sure the build fixes. And then you make a minor bump first. Again, then you start with this testing. And then after the testing is done, you understand everything is working fine, then you make another minor bump and release. Okay, also, at the same time, what I do want from you is to understand, do we have the VMware commands? I think we should have. But I want you to confirm this. Yeah, I think I asked for the VMware commands, like if we can have the VMware. VMware commands will work like the SSH command. So it would have sub commands like add VMs. I should be able to add VM. I should be able to install VMware. Okay? Add VMs to this. Not all the VMs, but selected VMs to my, let's say, default group, and work with those VMs. Like, what is the status? How many snapshots are there? What is their network? Do you want to change their, let's say, MAC addresses? Do you want to expand the hard drive? Do you want to fix the VMX files? Okay? All kinds of things we want to have this. So at the end, I want you to create a detailed, let's say, AI instruction that will have those CLI commands for the VMware. Add, remove, add groups, default groups, LS, help, what each section would do, so that if any AI sees this instruction, they could actually implement all these things. And for the end-to-end testing, I want you to install the VMware inside this VM as well, and you try to create a new VM, add that. So we should have these type of functionalities, actually. Let me check. Yeah, we do have the VMware inside. I'm not sure if the latest one needs to be installed. Let me check. The CTU cycle needs to be changed. Okay. Install the VMware first so that you can do the end-to-end testing. You can do anything inside. No worries. Do not think of anything. Maybe I could revert back because this is a VM, and everything will be good as new. Okay. VMware is installed. There is no issue on this. Okay, so I will install or enable the hypervisor so that you can test this. Okay, and you can create a VM for now, just any VM you can create. Okay. It doesn't have to be the real one. The idea is we should be able to change the MAC. Yeah, we wanted to start the machine as well. We want to see if the machine is running. We also want to shut down the machine. Also, we want to shut down the scheduler with the machine as well. Okay? All kinds of things we want. For this, I want you to write the spec and also the AI instruction so that I share with the AI, they would understand everything. Is it clear? Do you understand? I have installed VM Ware latest everything is working fine.

# Actionable Items Must Follow Non-Negotiable

1. Test the Chrome extension import-export feature with a sample profile.
2. Use the repo secrets folder for sensitive information.
3. Create a git map folder with sequence for task management.
4. Log each step to trace errors and save logs in the repo secrets folder.
5. Fix the git map build and make a minor bump before testing.
6. Confirm the availability of VMware commands and functionalities.
7. Create detailed AI instructions for VMware CLI commands.
8. Install VMware for end-to-end testing and create a new VM.
9. Ensure the ability to change MAC addresses, start, and shut down VMs.
10. Write specifications and AI instructions for VMware functionalities.

---

## Canonical Specifications
- Architecture Spec: [02-spec/21-app/72-chrome-ext-test-and-vmware-spec/01-architecture-spec.md](../../02-spec/21-app/72-chrome-ext-test-and-vmware-spec/01-architecture-spec.md)
- Component Spec: [02-spec/21-app/72-chrome-ext-test-and-vmware-spec/02-component-spec.md](../../02-spec/21-app/72-chrome-ext-test-and-vmware-spec/02-component-spec.md)
- AI Instruction Prompt: [01-prompts/42-vmware-cli-commands.md](../../01-prompts/42-vmware-cli-commands.md)

---

## Realized Subtask Outcomes & Evidence

### Subtask 01: Verify GitMap Build & Initial Minor Bump (v6.459.0)
- **Status**: Completed (`PASS exit 0`)
- **Outcome**: Inspected build integrity, ran SemVer minor bumper `03-ai-scripts/37-bump-version.py` from `6.458.0` to `6.459.0`.
- **Modified Manifests**: `version.json`, `package.json`, `cli/constants/constants.go`, `changelog.md`, `readme.md`, `what-to-read.md`.

### Subtask 02: Repo Secrets Sequence Provisioning & Discovery Logging
- **Status**: Completed (`PASS exit 0`)
- **Target Folder**: `./repo-secrets/01-gitmap/10-chrome-ext-test-and-vmware/`
- **Subdirectories**: `logs/`, `exports/`, `reports/`, `artifacts/`
- **Files Initialized**: `logs/execution.log`, `logs/audit-trace.json`, `logs/error.log`, `reports/verification-report.md`.
- **Host Discovery**: Discovered 47 Chrome profile folders, 46 registered profiles in `Local State`, cataloged `Profile 6` (isolated test sandbox) and `Profile 1` (rich extension candidate with 7 extensions).

### Subtask 03: Chrome Profile Export-Import E2E Test & Error Tracing
- **Status**: Completed (`PASS exit 0`)
- **Exports Validated**:
  - `Profile 6`: JSON, CSV, and ZIP export formats generated and checksummed.
  - `Profile 1`: JSON, CSV, and ZIP export with 7 active extensions preserved.
- **Destructive Deletion**: Executed `gitmap chrome-profile-delete "Profile 6" --yes`, confirmed directory removal and database row unlinking.
- **Restoration & Re-Import**: Executed `gitmap chrome import` from JSON snapshot and ZIP archive with explicit target. Restored Bookmarks, Preferences, and Local State registration with 100% data parity.
- **Error Diagnostics**: Traced and logged ZIP omission syntax error into `logs/error.log`.

### Subtask 04: VMware Workstation E2E VM Operations & Lifecycle
- **Status**: Completed (`PASS exit 0`)
- **Local Environment**: VMware Workstation Pro 25.0.1 build-25219725, `vmrun.exe` at `C:\Program Files (x86)\VMware\VMware Workstation\vmrun.exe`, `VMAuthdService` active.
- **Target VM**: `D:\testVM\Ubuntu 64-bit.vmx` (backed up to `Ubuntu 64-bit.vmx.bak`).
- **MAC Address Mutation**:
  - Tested static MAC configuration: `ethernet0.addressType = "static"`, `00:50:56:22:33:44`.
  - Tested dynamic MAC configuration: `ethernet0.addressType = "generated"`.
- **Lifecycle Verification**:
  - Headless power-on: `vmrun -T ws start "D:\testVM\Ubuntu 64-bit.vmx" nogui` -> exit 0.
  - Status query: `vmrun -T ws list` -> `Total running VMs: 1`.
  - Graceful ACPI shutdown: `vmrun -T ws stop "D:\testVM\Ubuntu 64-bit.vmx" soft` -> exit 0.
  - Final status query: `vmrun -T ws list` -> `Total running VMs: 0`, zero lingering `.lck` directories.
  - Scheduler / daemon check: `VMAuthdService` verified running.

### Subtask 05: VMware CLI AI Instruction Prompt Authoring & Cataloging
- **Status**: Completed (`PASS exit 0`)
- **Authored Prompt**: `01-prompts/42-vmware-cli-commands.md`
- **Catalog Registrations**: `01-prompts/readme.md`, `.ai-memory/prompts.md`.
- **Content**: Exhaustive CLI specification modelled after GitMap SSH fleet architecture (`add`, `install`, `group`, `default`, `ls`, `status`, `start`, `stop`, `mac`, `disk`, `snapshot`, `vmx`, `scheduler shutdown`), with target selector grammar, SQLite Split-DB schema, Go code standards, and micro-batch subagent execution roadmap.

### Subtask 06: Second Minor Bump (v6.460.0) & Final Release Sync
- **Status**: Completed (`PASS exit 0`)
- **Version Bump**: Bumped to `v6.460.0` via `python 03-ai-scripts/37-bump-version.py -t minor -s "vmware cli command suite specification and e2e test harness"`.
- **Synchronized Files**: `version.json`, `package.json`, `cli/constants/constants.go`, `changelog.md`, `.ai-memory/release/release-notes-v6.460.0.md`.
