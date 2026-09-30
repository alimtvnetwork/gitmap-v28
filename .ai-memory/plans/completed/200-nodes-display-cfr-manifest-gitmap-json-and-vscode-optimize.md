# Master Execution Plan: 200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md

> **Status:** `IN PROGRESS`
> **Target Version:** `v6.425.0`
> **Date:** `2026-09-30`
> **Author:** Antigravity Autonomous Lead Agent

---

## 1. User Request (Verbatim)

```text
PS D:\work> gitmap cfr
clone-fix-repo: ERROR <url> is required
  usage: gitmap clone-fix-repo <url> [folder]
  usage: gitmap clone-fix-repo-pub <url> [folder]
PS D:\work> gitmap cfr gitmap.json
Repository "gitmap.json" not found. Did you mean:
  gitmap-v28
PS D:\work> gitmap clone gitmap.json
PS D:\work>

Yeah, a couple of issues. So first of all, when we do the git-map nodes and when we do the SSH nodes, both of those are different actually things, how it displays. But I think they both can display same, and currently the SSH nodes looks very nice. It looks very professional. Try to make that compact output in the same case for the git-map nodes. Okay? Same thing. Okay. But also I think that also tries to showcase other stuff like which cluster or which command it is available to. For this, it can create another table in the nodes. The first table should show as SSH node does, and then it would actually give which command supports on which nodes. It would mention that, and then it just displays as suggestions. Rest of the things would be as same as the git-map SSH space node. That's the first thing. Second is that if we have a gitmap.json file, and if we just run the git-map CFR enter, then it will try to already clone the repositories which are missing only. Okay? Same thing will go for the git map clone, the gitmap JSON file, if we do it, or no file at all. It will also try to pick the gitmap.json file and try to clone. But also at the same time, if we go into the format of this JSON file, what we do see in the attributes is that-- So first of all, in all these attributes and all these JSON files regarding git map and others, you have default directory, one time mention, default path, one time mention, and also the variables one time mention. Why? It's very confusing. You just repeat. If you have the work directory, you just repeat that everywhere. Right? As a variable. Default path, it would be just the path variable, something like this. So you don't need to provide the same thing again and again. You just refer it from the variable. Okay, I hope it makes sense. Right? The default path would also be like this. Same path would be reused from the variable, which you are not doing. Please do that, and also which format is fine, but also at the same time, you don't have the import command. For the git map to clone, I think there are several commands. Right? So clone only the missing ones, CFR. Which each one will do, it also needs the help. So what you could do is add a section for the help, and then you can actually provide all these examples inside the help. Help would be a object rather than a single text. Okay, can you please do that? That's one thing. Second is that... Yeah. So the data should start from one, not zero. I think we did not corrected that. I requested it several times. And also make sure the git map merge JSON like this would work professionally without any hassle. Yeah, and try to have hard-coded object so that if things are failing, we know the stack trace and everything. Think about that, and also at the same time, for example, we did a gitignore stuff. Right? So for Antigravity Manager. So this way, we should be able to add gitignore stuff automatically. We could remove a file and then ignore it automatically, and that could be automation we could add inside the ignore. So add some help and example commands, other related commands that could automatically apply during the scan. And we can have the option, like do we want to apply during the scan or where we can disable, enable it from the settings as well. Try to understand this. This should be expandable setting. Okay. There is that. Now you tell me if you have any question and confusion. And also the absolute path that you have, that could also come using the variable. That means you could just use the work dir, and then just put it there. Yeah, you can just use the work dir and put it there. And also, I think we do have the PowerShell to run this thing, right? This JSON file. I think we need to update a little bit into the PowerShell code so that it can also understand this JSON and can clone from it. Because let's say we don't have the git map, we can actually bring that PowerShell file and then this JSON, and it should be able to clone using that as well. Is it clear? Do you understand the segments? Is everything clear?

make sure at the end you release minor bump and check ci cd

D:\work\repo-secrets\07-final-network-machine\gitmap.json

Also add one more command to the root command section. That would be git map VS Code optimize projects. That would actually remove the duplicates and also tell us where to move to duplicates projects. And if in the project section, if we have missing projects, that would also remove. So add the VS Code project optimizations for the project manager tool. Okay, mention this and also add it. Very important
```

---

## 2. Extracted Actionable Task List

- [x] **Task-01:** Harmonize `gitmap nodes` to match `gitmap ssh nodes` clean compact table formatting, add secondary Supported Commands Matrix table, and suggestions footer.
- [x] **Task-02:** Update `gitmap cfr` and `gitmap clone` to auto-detect `gitmap.json` when run without args or with `.json` manifests, cloning missing repositories only and avoiding auth hangs on existing repos.
- [x] **Task-03:** Refactor `D:\work\repo-secrets\07-final-network-machine\gitmap.json` to eliminate redundant work directory declarations, reuse `${workDir}`, upgrade `attributes.help` to rich object with examples, use 1-based indexing, and use `${workDir}` in `absolutePath`.
- [x] **Task-04:** Provide standalone `clone-gitmap.ps1` script (in `repo-secrets` and `scripts/`) to clone missing repositories from `gitmap.json` without requiring `gitmap` CLI binary.
- [x] **Task-05:** Implement expandable `autoGitignoreAgm` setting in settings store, config, and `gitmap gitignore agm` subcommands to automate AGM remediation on scan.
- [x] **Task-06:** Implement `gitmap vscode optimize-projects` (and aliases `vsc optimize-projects`, `vpm optimize`) to remove duplicates, advise relocation/consolidation, prune missing projects, and write clean `projects.json`.
- [x] **Task-07:** Run local Python linters, bump version to `v6.425.0`, update changelogs, commit, tag, push, and monitor CI/CD pipelines until green.
