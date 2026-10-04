# Subtask 217.04: Fleet Verification, Parity Scorecard & App Issue Registry Indexing

**Subtask Code:** 217.04  
**Parent Task:** `217-antigravity-fleet-parity-theme-preset-plugins-and-delegation`  
**Owner:** Worker 02  
**Target Files:**
- `02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/02-component-and-cli-spec.md`
- `02-spec/22-app-issues/69-antigravity-fleet-parity-theme-preset-plugins-rca.md`
- `02-spec/22-app-issues/readme.md`
**Status:** Ready for Implementation  
**Date:** 2026-10-05  

---

## 1. Context & Objective

Following the generation of the component & CLI specification, 4-part RCA document, and implementation of `sync-antigravity-full-profile.ps1` and `gitmap agy deploy`, rigorous end-to-end fleet verification must be conducted on target workstation `u1` (`192.168.1.22`).

The objective of this subtask is to execute and validate the complete verification scorecard across UI theme, permission presets, plugins, skills, path hygiene, and security hardening, and register App Issue 69 in the master contents catalog of `02-spec/22-app-issues/readme.md`.

---

## 2. Comprehensive Verification Scorecard

The following 6 verification gates MUST be verified against `u1`:

| Gate ID | Verification Domain | Expected Parity Target | Verification Command / Evidence |
| :--- | :--- | :--- | :--- |
| **VG-01** | **UI Theme Parity** | Dracula Dark: background `#19191C`, primary `#BD93F9`, foreground `#F8F8F2` | Remote `jq` check on `/home/a/.gemini/config/config.json` under `.userSettings.customThemeSeedsDark`; visual confirmation in IDE |
| **VG-02** | **Permission Preset** | Unattended Eager / Turbo mode (`CASCADE_COMMANDS_AUTO_EXECUTION_EAGER`) | Inspect `.userSettings.autoExecutionPolicy` and per-project `.settings.autoExecutionPolicy`; UI dropdown reflects Turbo |
| **VG-03** | **Official Plugins** | 4 Plugins installed and enabled (`chrome-devtools`, `data-agent-kit`, `google-antigravity-sdk`, `modern-web-guidance`) | Verify `/home/a/.gemini/config/plugins/` contains all 4 subdirectories; `.plugins` in `config.json` has `enabled: true` |
| **VG-04** | **Plugin Skills Count** | 43 nested skills with `SKILL.md` discovered | Count `find /home/a/.gemini/config/plugins -name "SKILL.md" \| wc -l` yields 42-43 skills |
| **VG-05** | **Filesystem Hygiene** | Zero stray Windows `C:*` folders in `/home/a/`; sanitized `instances.json` | `find /home/a -maxdepth 1 -name 'C:*'` returns 0 results; `instances.json` contains no backslashes |
| **VG-06** | **Chrome Sandbox SUID** | Ownership `root:root`, permissions `4755` (`-rwsr-xr-x`) | `ls -la /home/a/.local/share/antigravity-ide/chrome-sandbox` confirms SUID bit set |

---

## 3. App Issue Registry Indexing Procedure

In `02-spec/22-app-issues/readme.md`, append App Issue 69 to the master Contents table immediately following row 68:

```markdown
| 69 | [69-antigravity-fleet-parity-theme-preset-plugins-rca.md](69-antigravity-fleet-parity-theme-preset-plugins-rca.md) | Antigravity Fleet Parity, Theme & Preset Synchronization, and Delegation: RCA & Fix | Resolved |
```

### Verification Checks for Index Registration:
1. Confirm the row number is strictly sequential (`69`).
2. Confirm the relative Markdown link `[69-antigravity-fleet-parity-theme-preset-plugins-rca.md](69-antigravity-fleet-parity-theme-preset-plugins-rca.md)` resolves to the created RCA file.
3. Confirm table pipe alignment is clean and uncorrupted.
4. Verify acceptance criteria AC-AI-001 (App issues triage conformance: every issue contains Symptom, Cause, Fix, Prevention, and references a commit).

---

## 4. Execution Steps for Verification Gate

### Step 1: Remote Configuration Verification via SSH
```bash
ssh a@192.168.1.22 'python3 -c "
import json
with open(\"/home/a/.gemini/config/config.json\") as f:
    cfg = json.load(f)
assert cfg[\"userSettings\"][\"customThemeSeedsDark\"][\"primary\"] == \"#BD93F9\", \"Theme mismatch\"
assert cfg[\"userSettings\"][\"autoExecutionPolicy\"] == \"CASCADE_COMMANDS_AUTO_EXECUTION_EAGER\", \"Preset mismatch\"
assert len(cfg[\"plugins\"]) == 4, \"Plugin count mismatch\"
print(\"CONFIG_PARITY_VERIFIED\")
"'
```

### Step 2: Remote Skills Inventory Check
```bash
ssh a@192.168.1.22 'find /home/a/.gemini/config/plugins -name "SKILL.md" | wc -l'
# Expected output: >= 42
```

### Step 3: Filesystem Hygiene Check
```bash
ssh a@192.168.1.22 'find /home/a -maxdepth 1 -name "C:*" | wc -l'
# Expected output: 0
```

### Step 4: Chrome Sandbox Permissions Check
```bash
ssh a@192.168.1.22 'stat -c "%a %U:%G" /home/a/.local/share/antigravity-ide/chrome-sandbox'
# Expected output: 4755 root:root
```

### Step 5: Master Index Update
- Edit `02-spec/22-app-issues/readme.md` to register Issue 69.
- Verify Markdown table formatting and link resolution.

---

## 5. Acceptance Criteria

- [ ] All 6 verification gates in Section 2 pass with 100% compliance.
- [ ] Remote node `u1` runs Antigravity IDE with Dracula Dark theme and Turbo permission preset.
- [ ] Stray `C:\Users` folder is completely purged from `/home/a/`.
- [ ] Chrome sandbox has SUID root permissions (`4755`).
- [ ] `02-spec/22-app-issues/readme.md` is updated with Issue 69 properly linked.
- [ ] AC-AI-001 audit criteria are fully satisfied.
