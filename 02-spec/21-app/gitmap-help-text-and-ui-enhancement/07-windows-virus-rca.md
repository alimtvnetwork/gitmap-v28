# 07 — Windows VirTool:MSIL/Subti.C RCA

## 1. Ground truth

Commit `385aba38e` ("fix(cli): add credential/SSH remediation commands to
gitmap pa/fix and eliminate VirTool:MSIL/Subti.C heuristic triggers",
2026-10-15): 13 files, +241/−230. Microsoft Defender flagged the Windows
gitmap binary as `VirTool:MSIL/Subti.C` — a HEURISTIC detection (behavioral
pattern match), not a confirmed infection.

## 2. Heuristic triggers removed (the actual fix)

| File | Removed behavior | Why it trips the heuristic |
|------|-----------------|----------------------------|
| `cli/osutil/startup.go` | `reg add HKCU\...\Run` self-persistence (54 lines deleted; package reduced to a thin delegate) | Registry Run-key writes are the classic persistence signature for MSIL droppers |
| `cli/cmdos/os_autologin_windows.go` | Auto-login credential handling | Credential-harvesting pattern |
| `cli/cmdchromeprofile/chromeprofile_vss_windows.go` | VSS (Volume Shadow Copy) writer interop (−78 lines) | Shadow-copy access is ransomware-tooling behavior |
| `cli/fspath/lockcheck_windows.go` | Low-level file-lock probing (−70 lines) | File-system probing matches packer/loader behavior |
| `cli/cmdinstall/agm_uninstall.go`, `cli/cmdselfinstall/selfuninstallhandoff.go`, `cli/cmdschedule/schedule_cmd.go` | Startup/scheduler side paths trimmed | Scheduled-task + startup persistence combo |

Plus: `gitmap pa`/`fix` gained credential/SSH remediation commands so users
recover WITHOUT the removed low-level paths.

## 3. Root cause (stated plainly)

Legitimate admin features (startup persistence, auto-login, VSS snapshots,
lock probing) are byte-pattern-indistinguishable from the behaviors the
Subti.C heuristic hunts. Defender does not know intent — it scores behavior
combinations. Fix = remove the risky behaviors, provide safe alternatives.

## 4. Professional prevention (for the continuing AI)

1. **Code signing** — EV code-signing cert (~USD 300–500/yr) builds
   SmartScreen reputation; unsigned binaries start at zero trust.
2. **Microsoft Defender vendor submission** — submit each release binary to
   the Microsoft malware analysis portal; false-positive whitelisting
   typically clears within days.
3. **Reproducible builds** — deterministic build flags so a flagged binary
   can be rebuilt byte-identical for vendor appeal.
4. **CI AV gate** — scan release assets with Defender + 2–3 engines before
   publishing; quarantine on detection.
5. **Avoidance list (never re-add):** registry Run writes, VSS interop,
   auto-login credential code, raw file-lock probing in the shipped binary.

## 5. Open questions for the continuing AI

- Which Defender engine version + definition date flagged the binary?
- Is the current signed/unsigned state documented per release asset?
- Should a `gitmap doctor` AV-status check be proposed (WS4 candidate)?
