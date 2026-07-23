---
name: Clone multi-URL
description: gitmap clone accepts many URLs in one call — space-separated, comma-separated, or both. Each positional arg is split on commas and flattened into a single ordered list. Shipped in v3.80.0 (was planned for v3.38.0; see issue 06 for the PowerShell crash that forced delivery).
type: feature
---

# Feature: `gitmap clone <url1> <url2,url3> ...` (shipped v3.80.0)

**Spec:** `spec/01-app/104-clone-multi.md`
**Depends on:** existing direct-URL clone (`mem://features/clone-direct-url`)
**Triggered by:** Pending issue 06 — PowerShell silently splits unquoted commas, so `clone a,b,c` was mis-parsing `b` as a folder name and producing illegal Windows paths like `D:\...\https:\...`.

## Behaviour

- **Both syntaxes accepted, mixable.** `gitmap clone a b c`, `gitmap clone a,b,c`, and `gitmap clone a,b c d,e` all work. Parser: for each positional arg, split on `,`, strip whitespace, drop empties, append to ordered list.
- **Dedup case-insensitively** with trailing `.git` normalised, preserving first-seen order.
- **Defensive folder-name guard:** the second positional is ignored as a folder name when it looks like a URL — prevents `clone url1 url2` from ever being interpreted as `<url1> <folder=url2>`.
- **Existing flags unchanged:** `--target-dir`, `--github-desktop`, `--no-replace`, `--ssh-key/-K`, `--safe-pull`, `--verbose`.
- **`--github-desktop` registers each successful clone immediately** (inline message, not at end of batch).
- **Per-repo progress + final summary.** Continues on failure (no `--stop-on-fail` yet).

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | All cloned successfully |
| 1 | One or more clones failed |
| 3 | All URLs invalid — nothing attempted |

## Implementation map

| File | Role |
|------|------|
| `cmd/clonemulti.go` (new) | `flattenURLArgs`, `classifyURLs`, `executeDirectCloneOne`, `resolveCloneFolder`, `normaliseURLKey` |
| `cmd/clone.go` | `runClone` dispatch + `shouldUseMultiClone` detector + `runCloneMulti` loop |
| `cmd/rootflags.go` | `CloneFlags` struct exposing full positional slice; `isLikelyURL` guard for folder-name disambiguation |
| `constants/constants_clone.go` | `MsgCloneInvalidURLFmt`, `MsgCloneSummaryMultiFmt`, `MsgCloneRegisteredInline`, `MsgCloneMultiBegin`, `MsgCloneMultiItem`, `ErrCloneAllInvalid`, `ErrCloneMultiFailedFmt`, `ExitCloneMultiPartialFail`, `ExitCloneMultiAllInvalid` |

## Detection heuristic

`shouldUseMultiClone(cf)` returns true when **either**:
1. Any positional arg contains `,` (covers `a,b,c` even if PowerShell did NOT split it), or
2. There are 2+ positional args AND both `Arg(0)` and `Arg(1)` parse as direct URLs (covers PowerShell's silent comma-split into separate argv entries).

Single-URL invocations with an explicit folder name (`gitmap clone https://x/y my-folder`) keep the original single-clone path.
