# Parent Plan: gitmap-login-and-cloning-enhancements

## Description
Eliminate SSH cloning blocks and false PAT prompts, suppress Windows Credential Manager startup popups, enforce global Git identity reporting, and implement a unified `gitmap login` subsystem with multi-account SSH profile management.

## Status
- **Status:** COMPLETED
- **Created:** 2026-10-18
- **Phase:** 3 (Verification & Finalization)

## Architectural Decomposition
- **Subtask 01:** SSH Probe Fix, Direct SSH Clone Bypass, and `--ssh` Pipeline Routing
  - Fix `-h HEAD` bug in `cli/cmdclone/clone_auth_probe.go`.
  - Bypass auth probing for explicit SSH URLs in `clone.go`, `clonemulti.go`, and `clone_auth_resolve.go`.
  - Plumb `UseSSH` into `DirectCloneParams` and ensure `--ssh` routes directly to Git execution without terminal blocking.
- **Subtask 02:** GCM Modal Suppression, Global Git Identity Enforcement, and Unified `gitmap login` Subsystem
  - Suppress GCM Windows modal popup in `cli/secrets/ghtoken.go`.
  - Enforce strictly global `user.name` and `user.email` in `cli/cmd/rootusagecompact.go` and remove placeholder defaults.
  - Implement `gitmap login` interactive flow and `gitmap login ssh [alias] [-y]` in `cli/cmdlogin/`.
  - Add multi-account SSH binding in `cli/model/git_profile.go` and `cli/usercontext/project_binder.go`.

## Subtasks
- [.ai-memory/plans/subtasks/gitmap-login-and-cloning-enhancements/subtask-01.md](.ai-memory/plans/subtasks/gitmap-login-and-cloning-enhancements/subtask-01.md)
- [.ai-memory/plans/subtasks/gitmap-login-and-cloning-enhancements/subtask-02.md](.ai-memory/plans/subtasks/gitmap-login-and-cloning-enhancements/subtask-02.md)
