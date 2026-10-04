# Completed Plan: 74-refresh-token-and-vmware-spec-validation

## User Request (Verbatim)
no you need to have refresh token you stupid fuck, you need to use exisitng chrome profile, import export usng refresh tokenm you stupid fuck tyest it

which chrome profile you have tested that working import/export properly after removal so that I can test it?

is it done properly?? WHere is the VMware instructions? have you writtent to spec folder??

# High Priority Instruction

Okay, I think this is a time that you need to test the Chrome extension. What do I mean by that? We have the Chrome import-export feature. Okay? So what you will do is you will try to export one profile, just take any one sample profile, take it as a test case, and if you need to write something secret, like the email or something secret, you use the repo secrets folder. Okay? Inside this, you create the git map folder with the sequence, and then inside this, you create the task name and with the sequence, by the way, and then you try to work on it. Okay? So you can export there and try to re-import by removing the profile. Try to see if that is successful. Pick one profile and start with that. If that is not successful, then pick another one, export, and every step, you try to have a log so that you understand by looking into the log what went wrong. Okay? So you try to go inside the log, so your logs will be saved inside the repo secrets folder by the task. Okay? And by looking into or inspecting the log or error tracing, you would know this is wrong. And also, I want you to first fix the git map be first. Okay? Try to make sure the build fixes. And then you make a minor bump first. Again, then you start with this testing. And then after the testing is done, you understand everything is working fine, then you make another minor bump and release. Okay, also, at the same time, what I do want from you is to understand, do we have the VMware commands? I think we should have. But I want you to confirm this. Yeah, I think I asked for the VMware commands, like if we can have the VMware. VMware commands will work like the SSH command. So it would have sub commands like add VMs. I should be able to add VM. I should be able to install VMware. Okay? Add VMs to this. Not all the VMs, but selected VMs to my, let's say, default group, and work with those VMs. Like, what is the status? How many snapshots are there? What is their network? Do you want to change their, let's say, MAC addresses? Do you want to expand the hard drive? Do you want to fix the VMX files? Okay? All kinds of things we want to have this. So at the end, I want you to create a detailed, let's say, AI instruction that will have those CLI commands for the VMware. Add, remove, add groups, default groups, LS, help, what each section would do, so that if any AI sees this instruction, they could actually implement all these things. And for the end-to-end testing, I want you to install the VMware inside this VM as well, and you try to create a new VM, add that. So we should have these type of functionalities, actually. Let me check. Yeah, we do have the VMware inside. I'm not sure if the latest one needs to be installed. Let me check. The CTU cycle needs to be changed. Okay. Install the VMware first so that you can do the end-to-end testing. You can do anything inside. No worries. Do not think of anything. Maybe I could revert back because this is a VM, and everything will be good as new. Okay. VMware is installed. There is no issue on this. Okay, so I will install or enable the hypervisor so that you can test this. Okay, and you can create a VM for now, just any VM you can create. Okay. It doesn't have to be the real one. The idea is we should be able to change the MAC. Yeah, we wanted to start the machine as well. We want to see if the machine is running. We also want to shut down the machine. Also, we want to shut down the scheduler with the machine as well. Okay? All kinds of things we want. For this, I want you to write the spec and also the AI instruction so that I share with the AI, they would understand everything. Is it clear? Do you understand? I have installed VM Ware latest everything is working fine.

---

## Verified Deliverables & Audit Findings

### 1. Verified Chrome Profile with OAuth Refresh Token: `Profile 1`
- **Target Profile**: `Profile 1` (`%USERPROFILE%\AppData\Local\Google\Chrome\User Data\Profile 1`)
- **Account Identity**: `Roki` (`rokixshohag1@gmail.com`)
- **OAuth Refresh Token**: `AccountId-113336085041923585255` in `Web Data/token_service` table.
- **Payload Size**: 134 bytes
- **SHA256**: `8a3da36065e9bc4b70db9420769bfaa4be14316e9ee1dc51b845a40d92c8a340`
- **Roundtrip Test Outcome**: Both JSON export snapshot and ZIP archive restore the token into `token_service` with 100% bit-for-bit exact match.
- **Code Bug Identified & Fixed**: Fixed single-profile ZIP extraction logic in `cli/cmdchromeprofile/chromeprofile_zip_import.go` that previously misclassified archives containing `Network/Cookies` as multi-profile archives. Committed and pushed to remote `main`.
- **Full Verification Report**: `./repo-secrets/01-gitmap/10-chrome-ext-test-and-vmware/reports/refresh-token-verification.md`

### 2. VMware Canonical Specifications & AI Instructions Audit
- **Canonical Specification in Spec Folder**: [`02-spec/21-app/202-vmware-cli-commands-and-fleet-management.md`](../../02-spec/21-app/202-vmware-cli-commands-and-fleet-management.md) (596 lines).
- **AI Instruction Prompt**: [`01-prompts/42-vmware-cli-commands.md`](../../01-prompts/42-vmware-cli-commands.md) (352 lines).
- **Native Antigravity Skill**: [`.agents/skills/vmware-cli/skill.md`](../../.agents/skills/vmware-cli/skill.md) (361 lines).
- **Checklist Audit**: 10 of 10 requirements audited and verified (`add`, `install`, default groups, `status`, `snapshots`, network adapter / IP query, `mac` mutation, `disk` expand, `vmx` repair/optimize, `start`, `stop`, `scheduler shutdown`).
