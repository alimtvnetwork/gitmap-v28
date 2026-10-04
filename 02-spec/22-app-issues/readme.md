# App Issues

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 22 App Issues.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `22-app-issues/`.
- [ ] `/learn` Adhere strictly to `.ai-memory/folder-structure.md` and `.ai-memory/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.


. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


**Version:** 3.2.0
**Updated:** 2026-04-16
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Overview

App-specific issue analysis, root-cause analysis, bug documentation, and solution guidance — for whatever project this repo ships. Whether the app is a web app, Chrome extension, browser plugin, CLI tool, mobile app, WordPress plugin, or desktop app, **its bug reports and post-mortems live here.**

This folder tracks problems encountered during application development, their diagnosis, and their resolution.

---

## Placement Rule

Any content that analyzes bugs, failures, root causes, or fixes for application-level work belongs here, regardless of the app's runtime. General coding-principle violations or cross-cutting concerns belong in the core fundamentals range (`01–20`).

---

## Contents

| Number | Spec File | Summary | Status |
|---|---|---|---|
| 34 | [34-hd-hosted-docs-fallback.md](34-hd-hosted-docs-fallback.md) | Hosted-Docs Fallback (No More Hard-Exit on Missing docs-site) | Fixed |
| 35 | [35-reconcile-prompt-nested-if-ci-failure.md](35-reconcile-prompt-nested-if-ci-failure.md) | Reconcile Prompt Nested If CI Failure: Root Cause Analysis & Prevention | Resolved |
| 36 | [36-ubuntu-agy-agm-execution-errors.md](36-ubuntu-agy-agm-execution-errors.md) | Ubuntu AGY & Antigravity Manager Remote Execution Errors: RCA & Remediation | Resolved |
| 37 | [37-pipeline-cancelled-shown-as-pass.md](37-pipeline-cancelled-shown-as-pass.md) | Pipeline Canceled / Timed Out Runs Erroneously Displayed as PASS: Root Cause Analysis & Fix | Resolved |
| 38 | [38-paet-table-latency-and-gh-pr-bottleneck-rca.md](38-paet-table-latency-and-gh-pr-bottleneck-rca.md) | PAET Table Command Excessive Execution Latency: Root Cause Analysis & Mitigation | Resolved |
| 39 | [39-pull-batch-abort-and-search-untracked-directory-crash.md](39-pull-batch-abort-and-search-untracked-directory-crash.md) | Pull Batch Abort & Search Untracked Directory Crash: RCA & Fix | Resolved |
| 40 | [40-cfr-short-name-clone-failure-and-missing-gh-resolution.md](40-cfr-short-name-clone-failure-and-missing-gh-resolution.md) | CFR Short-Name Clone Failure & Missing GitHub CLI Resolution: RCA & Fix | Resolved |
| 41 | [41-vscode-startup-failure-and-search-latency-rca.md](41-vscode-startup-failure-and-search-latency-rca.md) | VS Code Startup Failure & Investigation Latency: Root Cause Analysis & Mitigation | Resolved |
| 42 | [42-ssh-install-exec-upload-failed-and-os-misclassification-rca.md](42-ssh-install-exec-upload-failed-and-os-misclassification-rca.md) | SSH Install-Exec 0ms Upload Failure & Remote OS Misclassification: RCA & Fix | Resolved |
| 43 | [43-gitmap-lowercase-unresolved-conflicts-and-missing-push-rca.md](43-gitmap-lowercase-unresolved-conflicts-and-missing-push-rca.md) | GitMap Lowercase Unresolved Merge Conflicts, Working Tree Dirt Contamination & Missing Push: RCA & Fix | Resolved |
| 44 | [44-ssh-install-exec-nsis-silent-hanging-and-empty-registry-resolution-rca.md](44-ssh-install-exec-nsis-silent-hanging-and-empty-registry-resolution-rca.md) | SSH Install-Exec NSIS Silent Hanging, Unattended Detection & Empty Registry Discovery Resolution: RCA & Fix | Resolved |
| 45 | [45-ssh-install-exec-stale-process-file-lock-and-test-db-wiping-rca.md](45-ssh-install-exec-stale-process-file-lock-and-test-db-wiping-rca.md) | SSH Install-Exec Stale Process File Lock & Unit Test DB Erasure: RCA & Fix | Resolved |
| 46 | [46-update-all-ssh-dial-failure-and-missing-command-suggestions-rca.md](46-update-all-ssh-dial-failure-and-missing-command-suggestions-rca.md) | SSH Fleet Update Authentication Dial Failure & Missing Command Suggestions: RCA & Fix | Resolved |
| 47 | [47-agy-rerun-killed-ide-and-selected-wrong-project-rca.md](47-agy-rerun-killed-ide-and-selected-wrong-project-rca.md) | Antigravity Rerun Process Termination and Inverted Project Selection: RCA & Fix | Resolved |
| 48 | [48-ssh-auth-handshake-failure-and-misleading-node-status-rca.md](48-ssh-auth-handshake-failure-and-misleading-node-status-rca.md) | SSH Fleet Authentication Handshake Failure and Misleading Node Status: RCA & Fix | Resolved |
| 49 | [49-ssh-fleet-parallel-pull-machine-hangs-and-concurrency-multiplication-rca.md](49-ssh-fleet-parallel-pull-machine-hangs-and-concurrency-multiplication-rca.md) | SSH Fleet Parallel Pull Machine Hangs and Concurrency Multiplication: RCA & Fix | Resolved |
| 50 | [50-agy-rp-recreate-alias-hijack-and-ad-hoc-project-creation-rca.md](50-agy-rp-recreate-alias-hijack-and-ad-hoc-project-creation-rca.md) | AGY RP Subcommand Alias Hijack, Unconfirmed Ad-Hoc Project Creation & CWD Pollution: RCA & Fix | Resolved |
| 51 | [51-mesh-ssh-deploy-keys-windows-chmod-failure-rca.md](51-mesh-ssh-deploy-keys-windows-chmod-failure-rca.md) | Mesh SSH Deploy Keys Windows Chmod Command Failure & Misleading Output: RCA & Fix | Resolved |
| 52 | [52-ssh-join-re-enrollment-auth-failure-and-double-shell-wrap-rca.md](52-ssh-join-re-enrollment-auth-failure-and-double-shell-wrap-rca.md) | SSH Join Re-Enrollment Auth Failure, Premature Interactive Prompt & Windows Double Shell Wrapping: RCA & Fix | Resolved |
| 53 | [53-macro-phantom-steps-and-removal-failure-rca.md](53-macro-phantom-steps-and-removal-failure-rca.md) | Macro Deletion Failure on Missing Target (`rm test`) and Interactive Edit Step Discard: RCA & Fix | Resolved |
| 54 | [54-windows-node-bash-missing-clone-failure-rca.md](54-windows-node-bash-missing-clone-failure-rca.md) | Remote Windows Node Clone Failure (`'bash' is not recognized`) and Fleet Output Degradation: RCA & Fix | Resolved |
| 55 | [55-test-auto-dest-repo-pollution-of-vscode-project-manager-and-agy-rca.md](55-test-auto-dest-repo-pollution-of-vscode-project-manager-and-agy-rca.md) | Test Destination Auto-Provisioning Polluting VS Code Project Manager and Antigravity Projects: RCA & Fix | Resolved |
| 56 | [56-ssh-fleet-pull-subprocess-bypass-and-wincredman-credential-failure-rca.md](56-ssh-fleet-pull-subprocess-bypass-and-wincredman-credential-failure-rca.md) | SSH Fleet Pull Subprocess Bypass & Wincredman Credential Store Failure: RCA & Fix | Resolved |
| 57 | [57-fleet-nodes-clone-misalignment-and-w3-liveness-rca.md](57-fleet-nodes-clone-misalignment-and-w3-liveness-rca.md) | Fleet Nodes Clone Table Misalignment, Stderr Escape Leaks & Remote Liveness False-Offline: RCA & Fix | Resolved |
| 58 | [58-windows-git-cache-credential-store-failure-rca.md](58-windows-git-cache-credential-store-failure-rca.md) | Windows Git Cache Credential Store Failure, Pull-First Workflow, and GitIgnore Duplication: RCA & Fix | Resolved |
| 59 | [59-fleet-nodes-clone-duration-alignment-and-w3-liveness-rca.md](59-fleet-nodes-clone-duration-alignment-and-w3-liveness-rca.md) | Fleet Nodes Clone Duration Alignment, Column Gutter Spacing, and W3 Reachability Diagnostics: RCA & Fix | Resolved |
| 60 | [60-rootcore-commit-dispatch-shadowing-and-stale-binary-rca.md](60-rootcore-commit-dispatch-shadowing-and-stale-binary-rca.md) | RootCore Commit Dispatch Shadowing and Stale Binary Macro Execution Failure: RCA & Fix | Resolved |
| 61 | [61-pull-all-unconstrained-ignore-deadlock-and-ssh-concurrency-rca.md](61-pull-all-unconstrained-ignore-deadlock-and-ssh-concurrency-rca.md) | Pull-All Unconstrained Ignore Subprocess Deadlock and SSH Concurrency Multiplication: RCA & Fix | Resolved |
| 62 | [62-pull-all-duplicate-repositories-and-case-sensitivity-rca.md](62-pull-all-duplicate-repositories-and-case-sensitivity-rca.md) | Pull-All Duplicate Repositories, Database Case-Sensitivity & Concurrent Fetch Failures: RCA & Fix | Resolved |
| 63 | [63-commit-push-clean-working-tree-exit-failure-rca.md](63-commit-push-clean-working-tree-exit-failure-rca.md) | Commit-Push Clean Working Tree Exit Status 1 Failure & Missing Staged Changes Guard: RCA & Fix | Resolved |
| 64 | [64-ssh-connection-timeout-and-target-firewall-sshd-rca.md](64-ssh-connection-timeout-and-target-firewall-sshd-rca.md) | SSH Connection Timeout, Target OpenSSH Server Absence & Firewall Block: RCA & Fix | Resolved |
| 65 | [65-commit-non-git-repo-crash-and-commit-all-delegation-rca.md](65-commit-non-git-repo-crash-and-commit-all-delegation-rca.md) | Commit Non-Git Repository Execution Crash, Child Batch Delegation & Terminal UI Output: RCA & Fix | Resolved |
| 66 | [66-ssh-enable-wuauserv-disabled-capability-crash-rca.md](66-ssh-enable-wuauserv-disabled-capability-crash-rca.md) | OpenSSH Server Capability Installation Crash When Windows Update (wuauserv) Service is Disabled: RCA & Fix | Resolved |
| 67 | [67-detectedproject-foreign-key-constraint-787-rca.md](67-detectedproject-foreign-key-constraint-787-rca.md) | DetectedProject Foreign Key Constraint (787) During Repository Scan: RCA & Fix | Resolved |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| App Specs | [../21-app/01-index.md](../21-app/01-index.md) |
| Spec Authoring Guide | [../01-spec-authoring-guide/01-index.md](../01-spec-authoring-guide/01-index.md) |

---

## Verification

_Auto-generated section — see `02-spec/22-app-issues/97-acceptance-criteria.md` for the full criteria index._

### AC-AI-001: App issues triage conformance: Index

**Given** Audit issue write-ups for the required Reproduction / Cause / Fix / Prevention sections.
**When** Run the verification command shown below.
**Then** Every issue file contains all four sections and references at least one commit or PR.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
