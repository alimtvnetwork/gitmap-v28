# Spec 180: Verification Gates & Quality Acceptance

Spec Reference: [01-overview.md](01-overview.md)

---

## 1. Acceptance Criteria Index

| ID | Title | Verification Condition |
|---|---|---|
| **AC-SPEC180-01** | AGY Snapshot Generation | `gitmap uninstall agy-all` writes valid JSON snapshot containing project paths and conversation titles before any deletion occurs. |
| **AC-SPEC180-02** | Work Directory Invariant | AGY uninstallation NEVER purges or modifies directories within `d:\work` or the active repository workspace. |
| **AC-SPEC180-03** | AGY Full Purge | All `.gemini` traces, brains, and temporary caches are removed when `--force` or confirmed. |
| **AC-SPEC180-04** | AGM Uninstallation | `gitmap uninstall agm-all` terminates running AGM processes and cleans up associated binaries and registry entries. |
| **AC-SPEC180-05** | Windows Copilot Removal | `gitmap uninstall copilot` removes Copilot Appx package and configures Group Policy registry disabling Copilot. |
| **AC-SPEC180-06** | Edge Browser Removal | `gitmap uninstall edge` removes Edge application using Chris Titus WinUtil removal sequence and applies reinstall block policy. |
| **AC-SPEC180-07** | DevTool Cache Enhancement | `gitmap devtool clear` calculates and reports space savings across Go, Node, Python, Vite, and Temp caches. |
| **AC-SPEC180-08** | Hermetic Test Isolation | Unit tests use injectable executors to verify removal logic without deleting live files on the developer host. |
| **AC-SPEC180-09** | Coding Guidelines Compliance | Code strictly complies with zero nesting (depth <= 2), positive boolean naming, and Unix LF line endings. |

---

## 2. Targeted Verification Commands

```bash
# Boolean Guidelines Linter
python linter-scripts/check-boolean-guidelines.py

# Nested-If Linter (Changed Files)
python linter-scripts/check-nested-ifs.py --changed-only

# Newline Styling Linter (Unix LF)
python linter-scripts/check-newline-styling.py
```
