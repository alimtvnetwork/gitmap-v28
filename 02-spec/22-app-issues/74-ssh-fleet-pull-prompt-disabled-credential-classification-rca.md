# Issue 74: SSH Fleet Pull "Cannot Prompt" Credential Failure Shown Raw & Unclassified RCA

> **Issue ID:** `ISSUE-74`
> **Target Subsystem:** `cli/cmdpull`
> **Affected Nodes:** `w2` (192.168.1.7, Windows) via `gitmap pa --ssh` orchestrator on `AI-Main`
> **Status:** Resolved
> **Date:** 2026-10-24
> **Related:** ISSUE-73 (wincredman persist failure), `.ai-memory/issues/10-git-push-hang-credential-manager-rca.md`

---

## 1. Problem Reproduction & Symptom Analysis

Running `gitmap pa --ssh` from the orchestrator (`AI-Main`), node `w2` pulled 68 repositories: 54 succeeded, 14 failed, every failure rendered identically:

```text
├── alim-cv-v8                             failed
│   ↳ Reason: fatal: Cannot prompt because user interactivity has been disabled.
│   ↳ Next Step: gitmap status alim-cv-v8 or gitmap fix alim-cv-v8
```

Affected repositories: alim-cv-v8, bright-buddy-block, bsrm-presentation-hiltrax-v4, digital-name-card, flat-slide-show, gitlogger-new-v2, hiltrax-v1, ki-health-ppt-v5, kita-social-media-content-calender-v2, lara-publishing-v1, maid-app-spec-presentation-v1, remix-of-presentation-riseup-asia-v7, seo-packages-v2, slides-spec.

The rendered reason is Git Credential Manager's raw refusal line, and the next step is the generic fallback hint — neither tells the operator that the node's stored HTTPS credential is missing or expired, nor that the purpose-built fixes (`gitmap fix` auth auto-fix / SSH conversion, `gitmap fix-credential`) apply.

## 2. Root Cause Analysis (4-Part RCA)

### 1. Symptom
Fourteen repositories on one node fail a fleet pull with the same raw GCM line and a generic remediation hint, repeated per repository with no classification.

### 2. Root Cause
1. **Node credential state:** The failing repositories use HTTPS remotes whose credential is missing or expired in `w2`'s credential store. Inside the non-interactive SSH session, GCM cannot prompt (`GCM_INTERACTIVE=never`, `GIT_TERMINAL_PROMPT=0` per the headless hardening in `cli/gitutil/git_env.go`), so Git fails fast with `fatal: Cannot prompt because user interactivity has been disabled.` — the fail-fast behavior ISSUE-10 deliberately introduced so missing credentials error in ~1.5s instead of hanging.
2. **Classifier gap:** `cli/cmdpull/pull_remediation_hint.go` recognized the wincredman *persist* failure (ISSUE-73) and basic auth failures (`Permission denied (publickey)`, `Authentication failed`), but had no predicate for the prompt-disabled signature. Unclassified errors fall through to `extractFirstMeaningfulErrorLine`, which surfaces the first raw `fatal:` line as the reason, and to the generic `gitmap status / gitmap fix` hint.

### 3. Resolution
In `cli/cmdpull/pull_remediation_hint.go`:
- Added `isPromptDisabledFailure`: matches GCM's `Cannot prompt because user interactivity has been disabled`, plus Git's own non-GCM variants of the same class (`terminal prompts disabled`, `could not read Username`, `could not read Password`).
- `classifyPullErrorString` now returns a human reason: `Git credentials unavailable - credential prompt is disabled in this non-interactive session`.
- `buildActionableRemediation` now returns the targeted fix: run `gitmap fix <repo>` (auto-fix auth & convert remote to SSH) or sign in on the node via `gitmap fix-credential` (alias: fc).
- Structured remediation now appends the existing `Switch Remote to SSH` option (`gitmap pull <repo> --ssh`) for this failure class.
- In `cli/cmdpull/pull_remediation_fallback.go`, dual hints now route this class to `resolveCredentialAndSSHDualHints` (`Auto-Fix Auth & Convert to SSH` / `Fix Credential Store`), same as wincredman and auth failures.
- Added tests in `cli/cmdpull/pull_efficient_render_test.go`: `TestResolveDualPullRemediationHints_PromptDisabled`, `TestClassifyPullErrorString_PromptDisabled`, `TestResolvePullRemediationHint_PromptDisabled`.

### 4. Prevention & Learnings
- Every GCM/Git failure string that headless hardening can produce must have a classifier entry; a fail-fast error that renders raw is only half a fix. When adding a new headless guard or credential store, add its failure signatures to `cli/cmdpull` classification in the same change.
- The remediation for credential-class failures is node-level (store a credential once, or convert remotes to SSH); per-repo repetition of a generic hint hides that one node-level action clears the whole group.

## 3. Verification & Acceptance Testing

1. **Unit tests:** `go test ./cli/cmdpull -run "PromptDisabled"` — the three new tests pass (verified via CI).
2. **Fleet behavior:** after deploying the new build to `w2`, the same fourteen failures render the classified reason and the `gitmap fix` / `gitmap fix-credential` next step instead of the raw `fatal:` line; running `gitmap fix-credential` on the node (interactive sign-in once) or `gitmap fix <repo>` (SSH conversion) clears them.
