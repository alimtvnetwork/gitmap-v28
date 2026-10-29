# Execution Plan: Git Prime Machine OS Version and UI Help Docs Improvement

- **Slug:** `git-prime-machine-os-version-and-ui-help-docs-improvement`
- **Spec Directory:** `02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/`
- **Subtasks Directory:** `.ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/`
- **Database Path:** `.ai-memory/temp-agents/02-git-prime-machine-os-version-and/agent-task.db`
- **Status:** ACTIVE
- **Total Steps Budget:** 300
- **Operational Waves:** 4

---

## User Request (Verbatim)

```text
PS D:\work> gitmap machine                                                                                      
● Machine Identity:                                                                                             
  • Alias:           192.168.1.20                                                                               
  • Hostname:        AI-Main                                                                                    
  • IP Address:      192.168.1.20                                                                               
  • OS Platform:     Windows (windows) (Windows Server 2022 Standard 21H2 (Build 20348))                        
  • Architecture:    amd64 (20 CPUs)                                                                            
  • User:            Administrator                                                                              
  • GitMap Version:  6.533.0                                                                                    
PS D:\work> gitmap machine --json                                                                               
[                                                                                                               
  {                                                                                                             
    "sequence": 1,                                                                                              
    "nodeId": "local-01",                                                                                       
    "ipAddress": "192.168.1.20",                                                                                
    "alias": "192.168.1.20",                                                                                    
    "machineName": "AI-Main",                                                                                   
    "osHostname": "AI-Main",                                                                                    
    "osPlatform": "Windows (windows)",                                                                          
    "osVersion": "Windows Server 2022 Standard 21H2 (Build 20348)",                                             
    "architecture": "amd64",                                                                                    
    "cpuCores": 20,                                                                                             
    "currentUser": "Administrator",                                                                             
    "gitmapVersion": "6.533.0",                                                                                 
    "scope": "local"                                                                                            
  }                                                                                                             
]                                                                                                               
PS D:\work>                                                                                                     

# Git Prime Machine OS Version and UI Help Docs Improvement: high priority instruction, non-negotiable task

Okay, I think there is a serious issue when we say Git Prime machine. It should always find the actual OS version. It's not like this, what you are showing. So the OS version should be a block of object that would have the type of OS, that would be Windows in this case. Is it a Windows Server? It would say, "Is it Windows Server?" Yes. Is it a pro? No. Is it a Windows Server? Yes. Okay. So it will have the version information also, like what version it is, exact version, and also the versions, major, minor, all these things in a breakdown manner so that anyone can read and use this to understand, do many more stuff. Okay? The JSON version should have this detail, I think, which we are missing right now. And also, one more serious issue is in the UI of the HD docs, help docs. We don't have the new commands help. That is terrible. And also the merge AI you have created, you did not give a parameterized example like where to pass the destination, the other URLs, the text version example. Nothing is there. It's pretty terrible on the new commands that you have implemented. Can you just work on it and improve and fix and make a minor bump and release?

## slug: git-prime-machine-os-version-and-ui-help-docs-improvement

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/ and enqueue plan task in .ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md (subtasks in .ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Use `gitmap` AI agents to enter data
5. Task completion includes committing and pushing to Git
6. Ensure the OS version is accurately detected and represented as a JSON object with type, version, major, minor details
7. Update the UI of the HD docs to include help for new commands
8. Provide parameterized examples for the merge AI commands, including destination, URLs, and text version examples
9. Implement improvements, fix issues, and make a minor version bump for release

## Must follow and spawn agent using

@[\1](file;.agents/skills/\1)

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[\1](file;.agents/skills/\1) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## Actionable Deliverables & Wave Schedule

### Wave 1: Discovery, Architecture Specification & Subtask Decomposition (Phase 1)
- **Task-01:** Discovery across Go machine identity, HD docs UI, and merge-ai commands (Research Subagents 01 & 02).
- **Task-02:** Author canonical specs `01-architecture-spec.md`, `02-component-and-cli-spec.md`, and 4 subtask plans in `.ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/` (Spec Workers 01 & 02).

### Wave 2: Core Go Engine Implementation (Phase 2)
- **Task-03:** Implement `OSVersionInfo` struct in `cli/cmdos/os_info_types.go`, Windows registry probing in `cli/cmdos/os_info_windows.go`, and Unix/macOS parsing in `cli/cmdos/os_info_other.go`.
- **Task-04:** Update `MachineIdentity` and `renderSingleMachineIdentity` in `cli/cmdos/os_machine_alias.go` with formatted terminal block and structured JSON output under `osVersionDetails`.
- **Task-05:** Implement CLI help flag routing, named flag parsing (`--destination`, `-d`, `--source`, `-s`, `--url`, text manifest, JSON config), and parameterized examples in `cli/cmdmergeai/merge_ai_core.go` and `cli/cmdmergeai/merge_ai_help.go`.
- **Task-06:** Add CLI help markdown topics `cli/helpdoc/machine.md` and `cli/helpdoc/merge-ai.md`, and register in `cli/helpdoc/catalog.go`.

### Wave 3: Help Docs (HD) Web UI Enhancement (Phase 2)
- **Task-07:** Expand command database in `src/data/commands.ts` with new commands (`merge-ai`, `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, `cpar`) with full parameter examples.
- **Task-08:** Update `src/pages/HelpIndex.tsx`, `src/components/docs/DocsSidebar.tsx`, and `src/App.tsx` routing. Create dedicated documentation page `src/pages/MergeAIPage.tsx`.

### Wave 4: Quality Gate, Minor Version Bump & Release Ceremony (Phase 3)
- **Task-09:** Validate with `python 03-ai-scripts/50-fastgate.py` (<1.5s SLA) ensuring relative paths and zero guideline violations.
- **Task-10:** Execute minor version bump (`v6.533.0` &rarr; `v6.534.0`), update `version.json`, `package.json`, `changelog.md`, consolidate completed plan, atomic commit, and upstream push via GitMap (`gitmap cpf "machine - structured os version object, hd docs ui and merge-ai parameterized examples"`).
