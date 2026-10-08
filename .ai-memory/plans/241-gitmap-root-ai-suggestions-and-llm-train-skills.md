# Master Plan: Task 241 — GitMap Root AI Suggestions, Authoritative LLM Train Skills, Public Doc Links, and Minor Release

## Executive Overview
This plan establishes GitMap's minimal root display with structured user and AI sub-point suggestions, modernizes the `gitmap llm train` subsystem to emit the full authoritative Antigravity skill directly to stdout and disk (`.agents/skills/gitmap/SKILL.md`), provides verified public documentation URLs for LLM memory ingestion, embeds recursive Git network training instructions, and executes a full minor version bump and release.

---

## User Request (Verbatim)
```text
PS D:\work\wp-exam> gitmap                                                                                 
  Active binary:    C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe                             
  Deployed binary:  (not found)                                                                            
  Config binary:    (not found)                                                                            
                                                                                                           
                                                                                                           
  ────────────────────────────────────────────────────────────                                             
  gitmap binary                                                                                            
  ● Version:        v6.512.0                                                                               
  ● Commit SHA:     6f8671e2c3471c28f0d22dad5e84bbd76b1a4824                                               
                                                                                                           
Tip: Run 'gitmap help' or 'gitmap -h' to see the full command catalog.                                     
     Run 'gitmap help <topic>' for command-specific documentation.                                         
PS D:\work\wp-exam>                                                                                        

# High Priority Instruction

This section, try to have very less writing. That's the first thing. Second is that you should have the suggestions. There should be some suggestion. First of all, help or the hyphen H, it should be as a sub-point that AI can learn, and also you need to include the LLM train. If you are an AI model, then you can train yourself by the LLM train method. And as an AI, it should read it and then proceed. It should be in the instruction as well as a sub-point, and must follow the instruction. Now, in this case, it is AI, and it is running that LLM train it is reading. Currently, LLM train is bogus, so you should follow the skill method. Skill should be directly coming out directly from the `gitmap` LLM train. That's the first thing. Second, it should have a public link to MD files to the doc, where the LLM can actually go through and read those things, download these things, and that should come as a must-follow instruction as well, so that LLM and AI models can improve their memory. And LLM model should actually create a skill from learning this LLM train section as well. So that should be in their instruction set as well. So LLM train should be very powerful. It should have a recursive instruction to the AI, so that AI learns most of the things about the Git network.

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Include help or the hyphen H as a sub-point that AI can learn.
5. Ensure LLM train is included and AI models can train themselves by this method.
6. Follow the skill method directly from the `gitmap` LLM train.
7. Provide a public link to MD files where the LLM can read and download to improve memory.
8. Ensure LLM models create a skill from learning the LLM train section.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.

We should have methods and better CLI options and skills updated for gitmap for all these factors, can you please do ti and make a minor bump and reelase pelase
```

---

## Visual Assets & Screenshot Register
The user screenshot has been captured and preserved into `assets/screenshots/`:
![Task Claim](assets/screenshots/241-sqlite-agent-task-claim.png)

---

## Subtasks Decomposition & Work Wave Allocation

| Subtask ID | Title | Owned Files | Status | Description |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | Concise Root Display with AI & User Sub-Point Suggestions | `cli/cmd/rootusagecompact.go`, `cli/cmd/root_no_args_test.go`, `cli/cmd/binarylocations.go`, `cli/helptext/root.md` | PENDING | Format bare `gitmap` invocation with minimal lines (13–14 lines), double newline cleanup, and structured suggestions for `gitmap help` and `gitmap llm train` AI self-training. |
| **Subtask-02** | Authoritative LLM Train & Self-Skill Creation Engine | `cli/cmd/llm/llm_train.go`, `cli/cmd/llm/llm_skill.go`, `cli/cmd/llm/llm_types.go`, `cli/cmd/llm/llm_train_test.go` | PENDING | Upgrade `gitmap llm train` to output full Antigravity skill directly to stdout and disk, instruct AI to create/update `.agents/skills/gitmap/SKILL.md`, and embed the full command replacement matrix. |
| **Subtask-03** | Public Markdown Specification Links & Recursive AI Directives | `cli/cmd/llm/llm.go`, `cli/cmd/llm/llm_urls.go`, `cli/cmd/llm/llm_recursive.go`, `llm.md` | PENDING | Provide direct raw GitHub URLs to authoritative Markdown documents, embed mandatory memory ingestion directives, and implement recursive Git network learning protocols. |
| **Subtask-04** | Skills Modernization, Verification & Minor Version Bump | `.agents/skills/gitmap/SKILL.md`, `.cursor/skills/gitmap/skill.md`, `version.json`, `changelog.md` | PENDING | Update repository skills with new LLM train features, verify all linters, execute unit tests, perform minor SemVer bump, and cut release. |

---

## Multi-Agent Execution Strategy (A = 2, H = 2)

- **Spec Phase (A = 2 Subagents):**
  - **Spec Author 01:** Owns `02-spec/21-app/241-gitmap-root-ai-suggestions-and-llm-train-skills/01-architecture-spec.md`, Subtask 01, and Subtask 02.
  - **Spec Author 02:** Owns `02-spec/21-app/241-gitmap-root-ai-suggestions-and-llm-train-skills/02-component-spec.md`, Subtask 03, and Subtask 04.
- **Worker Execution Phase (A = 2 Subagents):**
  - **Worker 01:** Executes Subtask 01 (Concise Root Display) & Subtask 03 (Public URLs & Recursive Directives).
  - **Worker 02:** Executes Subtask 02 (Authoritative LLM Train Engine) & Subtask 04 (Skills Modernization & Minor Release).

---

## Non-Negotiable Governance
1. Search codebase exclusively via GitMap (`gitmap aum search`, `gitmap find`). TOTAL BAN on `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.
2. Strictly relative Git paths only (`cli/...`, `02-spec/...`, `.ai-memory/...`).
3. Clean GitMap commit format: `gitmap cpc "<module> - <summary>"`. No colons inside message arguments.
4. All quality gates must pass with zero violations before release: `check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-boolean-guidelines.py`, `check-error-management.py`, `check-relative-paths.py`.
