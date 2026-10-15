# RCA 273 — Windows `VirTool:MSIL/Subti.C` heuristic flag (program 273, WS5)

> Ground truth: commit `385aba38e` ("fix(cli): add credential/SSH remediation
> commands to gitmap pa/fix and eliminate VirTool:MSIL/Subti.C heuristic
> triggers", 2026-10-15, 13 files, +241/−230). Write-up for the continuing
> AI: verified file-by-file against the diff, not copied from the spec.
> Read: `.ai-memory/what-to-read.md` mindset first (staging repo: velocity
> over ceremony; enforcement beats documentation).

## 1. Incident

Microsoft Defender flagged the Windows `gitmap` binary as
`VirTool:MSIL/Subti.C`. This is a **heuristic (behavior-pattern) verdict,
not a confirmed infection**. No malware was found in the repo; the binary
was tripping Defender's behavior-scoring for dropper/persistence tooling.
The fix shipped in 385aba38e removed the offending patterns instead of
disputing the detection.

## 2. What `VirTool:MSIL/Subti.C` actually means

Decoded from Microsoft's malware-naming convention (public, documented
scheme — no encyclopedia page exists for Subti.C specifically, which is
itself a signal: it is a heuristic/behavioral verdict, not a static
family signature):

| Segment | Meaning |
|---------|---------|
| `VirTool` | Microsoft's category for **tools that are not necessarily malicious but exhibit tool-like behaviors** (persistence installers, credential utilities, process manipulators). The classic "legit admin tooling looks like malware tooling" bucket. |
| `MSIL` | Microsoft Intermediate Language — normally tags .NET assemblies. **Naming quirk:** gitmap is a Go (native PE) binary, not MSIL. A behavioral verdict applied under an MSIL label is consistent with cloud-delivered / behavioral misattribution rather than a static byte-signature hit. |
| `Subti` | Variant family name ("subtlety") — Defender's label for a specific behavior cluster. No public description published. |
| `.C` | Third variant of the Subti behavior set. |

Net: Defender scored the *combination* of behaviors in one unsigned
binary — registry Run writes via spawned `reg.exe`, spawned `taskkill`,
hidden `powershell -Command` WMI queries, `vssadmin` shadow-copy calls,
auto-login credential registry writes, self-deletion via `cmd.exe` —
and pattern-matched it to the Subti.C tool cluster. Defender scores
behavior combinations, not intent. This is the mechanism behind every
false positive of this class.

## 3. Confirmed trigger patterns (verified against the 385aba38e diff)

| File | Removed pattern (pre-fix code) | Why it trips the heuristic |
|------|-------------------------------|----------------------------|
| `cli/osutil/startup.go` | `exec.Command("reg","add", HKCU\...\Run, "/v","GitmapMacro", ...)` — self-persistence via spawned `reg.exe` (54 lines deleted; file became a thin delegate to `startup.RegisterMacroAutostart`) | **Spawning `reg.exe` to write Run keys is the textbook MSIL-dropper persistence signature.** The Run key is written *by a child process*, the exact shape malware uses. |
| `cli/cmdchromeprofile/chromeprofile_vss_windows.go` | `exec.Command("vssadmin","create","shadow",...)` + `vssadmin delete shadows` (−78 lines; `CreateSnapshot` now returns `ok=false`, callers use shared-read copy fallback) | Volume Shadow Copy access = **ransomware-tooling behavior**. Legitimate reason (copy locked Chrome profile files) is irrelevant to the scorer. |
| `cli/cmdos/os_autologin_windows.go` | `_ = k.DeleteValue("DefaultPassword")` removed from `Disable()` (1-line removal; the auto-login credential handling around it) | Writing/clearing plaintext auto-login credentials in `Winlogon` = **credential-harvesting pattern**. |
| `cli/fspath/lockcheck_windows.go` | `findViaWMI` — spawned `powershell -NoProfile -Command "Get-CimInstance Win32_Process ..."` with `SysProcAttr.HideWindow` (−70 lines); `KillProcess` via spawned `taskkill /F /PID` with hidden window; `findViaHandle` dropped the hidden-window attr | **Hidden-window PowerShell WMI process enumeration + `taskkill`** = packer/loader + process-killer combo. The `HideWindow` flag is itself a suspicious-behavior signal. Replaced with `os.FindProcess().Kill()` and handle.exe-only probing (no WMI fallback). |
| `cli/cmdinstall/agm_uninstall.go` | Spawned `taskkill /F /PID` with hidden window (replaced by `os.FindProcess`); loop of `reg delete HKCU\...\Run /v <key> /f` across 5 AGM key names (replaced by comment: installer manages its own keys) | Same as above: hidden `taskkill` + Run-key deletion loop = **uninstaller-shaped persistence tampering**. |
| `cli/cmdselfinstall/selfuninstallhandoff.go` | `cmd.exe /C ping 127.0.0.1 -n 2 >nul & del /F /Q <self>` — delayed self-deletion via spawned shell, hidden | **Self-deleting binaries** are a classic dropper/loader trait. Now: Unix unlinks immediately; Windows relies on OS temp-dir cleanup. |
| `cli/cmdschedule/schedule_cmd.go` | `osutil.AddToStartup(...)` → `startup.RegisterMacroAutostart(...)` (2 call sites) | Routed startup registration through the native-API backend instead of the `reg.exe`-spawning path. |

Key nuance (do not misstate): the Windows Run key is **still written** — by
`cli/startup/winregistry_windows.go` via the native Go
`golang.org/x/sys/windows/registry` API (`gitmap-<name>` value). The fix
removed the *dropper-shaped delivery* (spawning `reg.exe`/`taskkill`/
`cmd.exe`/`powershell` as hidden children), not the feature itself.
Native API calls score far lower than LOLBin-spawning patterns.

## 4. What the fix added (safe replacements, same commit)

- `cli/cmdpull/pull_remediation.go` (+151): `gitmap pa`/`fix` auto-detect
  credential/HTTPS failures (`isCredentialOrHTTPSFailure`) and remediate by
  converting remotes to SSH — recovering users **without** the removed
  low-level paths.
- `cli/cmdfix/fix_cmd.go` (+18): `fix --all` drives `RunBatchPullFix`;
  credential failures route to the new `gitmap fix-credential` command.
- `cli/cmdpull/pull_remediation_hint.go`, `pull_remediation_fallback.go`,
  `pull_efficient_render.go`: hints now suggest `gitmap fix-credential &&
  gitmap fix --all` / `gitmap pa --ssh` instead of low-level fixes.
- `cli/startup/add.go` (+10): `RegisterMacroAutostart` — the single native
  entry point replacing `osutil.AddToStartup`.

## 5. Professional prevention plan

### 5a. Code signing (verified pricing, Oct 2026)

| Option | Cost (USD/yr) | Effect |
|--------|---------------|--------|
| OV (Organization Validated), e.g. SSL.com / Sectigo via reseller | ~$150–$220 | Publisher name shown; SmartScreen reputation builds over time/downloads |
| EV (Extended Validation), e.g. Sectigo/Comodo via reseller | ~$280–$380 | Historically fastest SmartScreen trust; key on FIPS hardware token/HSM |
| EV direct (DigiCert) | ~$560–$900 | Premium; same trust effect |
| Azure Trusted Signing (Microsoft) | ~$10/mo Basic | Cloud HSM signing, no token to guard; issues OV-equivalent (no EV tier) |

Sources: https://github.com/nerdbase-by-stark/skill-forge/blob/HEAD/example-research/02-code-signing.md,
https://my-ssl.com/learn/azure-trusted-signing-vs-code-signing-certificate,
https://eatsip365.com/pricing-a-code-signing-certificate/.
Notes: max cert lifetime is now ~459 days (CA/B Forum, Feb 2026), so
"multi-year" = annual reissue; EV private key must live on hardware token
or HSM; since Mar 2024 Microsoft no longer grants EV *instant* SmartScreen
bypass — reputation still accrues faster. **Recommendation: start with
OV (~$200/yr) or Azure Trusted Signing; escalate to EV only if SmartScreen
blocks persist.**

### 5b. Microsoft Defender vendor submission (per release, until signed)

1. Portal: https://www.microsoft.com/wdsi/filesubmission (no M365
   subscription needed). Sign in, select **No** for "Do you believe this
   file contains malware?", choose the **software developer** submission
   type, attach the exact release binary + SHA-256 + detection name.
   Source: https://learn.microsoft.com/en-us/defender-office-365/submissions-submit-files-to-microsoft
2. Keep the Submission ID (e.g. `7c6c214b-…`). Turnaround is typically days.
3. **A clearing applies to that exact file hash** — every release needs its
   own submission until the binary is signed and reputable. Budget this
   into the release ceremony (`02-spec/21-app/08-distribution-and-release/`).
4. Also scan each release on VirusTotal pre-publish: 1–2 engines = noise
   (often self-clears on next definition update); 10+ = report to each
   vendor — vendor verdicts are NOT shared between vendors.

### 5c. Reproducible builds

Pin deterministic Go build flags (`-trimpath`, fixed `ldflags`, pinned
toolchain via `go.mod` `toolchain` line) so a flagged binary can be rebuilt
byte-identical for vendor appeal and user-side hash verification. Publish
SHA-256 per asset on every release page.

### 5d. CI AV gate (proposed, not yet built)

Add a release-pipeline step (`.github/workflows/`): after assets are built,
scan with Windows Defender (`MpCmdRun -Scan`) on a Windows runner plus
upload to VirusTotal; **quarantine the release (block publish) on any
named-family detection**. 1–2 generic `!ml` hits = warn-only. This is the
"enforcement beats documentation" version of the avoidance list.

### 5e. Code-pattern avoidance list (NEVER re-add to the shipped binary)

1. Spawning `reg.exe` / `regedit` (use `golang.org/x/sys/windows/registry`).
2. Spawning `taskkill` (use `os.FindProcess().Kill()`).
3. Hidden-window child processes (`SysProcAttr.HideWindow`) — the flag
   itself is a signal.
4. `powershell -Command` / `cmd.exe /C` subprocesses for system queries
   (use native APIs or documented CLIs without shell wrappers).
5. `vssadmin` / VSS writer interop (use shared-read copy fallbacks).
6. Auto-login credential registry writes (`Winlogon\DefaultPassword`).
7. Self-deleting binaries (delayed `del` via `cmd.exe`).
8. WMI process enumeration from the binary (`Get-CimInstance Win32_Process`
   via spawned PowerShell).

Any future feature needing one of these must go through the release-repo
stricter process (per the what-to-read two-tier model), not the staging
repo.

## 6. Open research questions (for the continuing AI)

1. Which Defender engine + security-intelligence version/definition date
   flagged the binary? (Needed to prove the fix cleared *that* verdict —
   record `MpCmdRun -SignatureUpdate` / definition version next time.)
2. Is the current signed/unsigned state documented per release asset?
   (Propose a `release-assets.md` or manifest with signer + SHA-256.)
3. Should `gitmap doctor` gain an AV-status check (Defender exclusion /
   detection-history probe) — WS4 candidate from the spec?
4. Why did a Go native-PE binary attract an **MSIL**-labeled verdict?
   (Behavioral misattribution vs. embedded-resource pattern — worth one
   targeted test with a minimal repro binary.)
5. CI AV gate: build the quarantine-on-detection step (5d) and wire it
   into the release workflow before the next Windows release.
