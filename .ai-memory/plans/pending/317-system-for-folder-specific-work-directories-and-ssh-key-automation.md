# Master Plan: System for Folder-Specific Work Directories and SSH Key Automation

## Metadata
- **Task ID:** 317-system-for-folder-specific-work-directories-and-ssh-key-automation
- **Slug:** system-for-folder-specific-work-directories-and-ssh-key-automation
- **Spec Path:** `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/`
- **Status:** Completed (Verified, Tested, Released v6.604.0, Fleet Synced)
- **Concurrency Budget:** A = 3, H = 2, Concurrency Capacity = 6

---

## High-Level Objective
Architect, implement, and verify a complete folder-specific and work-group-specific work directory and automated SSH login engine backed by GitMap's root SQLite database, enhance `cpf`/`cbf` bracketed commit format syntax, provide intuitive grouping commands, and synchronize documentation, help menus, LLM train instructions, and Antigravity skills.

---

## Discrete Task Breakdown
- [x] **Task 317-01:** Workgroup & Folder Directory DB Schema & Repository Architecture (Split-DB SQLite).
- [x] **Task 317-02:** Autonomous Git Operation Interceptor (`pull`, `push`, `fetch`) with Folder/Group SSH Resolution.
- [x] **Task 317-03:** Workgroup CLI Commands (`gitmap group`, `gitmap workdir`, `gitmap user group`).
- [x] **Task 317-04:** Enhanced `cpf` & `cbf` Bracketed Module Format Syntax (`[module]: item`).
- [x] **Task 317-05:** End-to-End Test Suite for Folder-Specific Workdirs & SSH Automation.
- [x] **Task 317-06:** Help Documentation, HD Help, UI Menus, LLM Training & Skill Synchronization across Repositories.
