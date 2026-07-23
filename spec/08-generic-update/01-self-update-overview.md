# 01 — Self-Update Overview

## Purpose

Explain why CLI self-update is non-trivial, what platform-specific
constraints exist, and how the overall architecture addresses them.

---

## The Problem

When a user runs `<binary> update`, the tool must:

1. Fetch the latest source or download a pre-built binary.
2. Build (if from source) or extract the new binary.
3. Replace the currently running binary with the new one.
4. Verify the update succeeded.
5. Clean up temporary artifacts.

**Step 3 is the hard part.** On Windows, a running `.exe` file is
locked by the OS — it cannot be overwritten or deleted while the
process is alive. On Linux/macOS, the file can be replaced in-place
(the OS uses inode references, not file paths, for running processes).

---

## Platform Behavior

| Operation | Windows | Linux / macOS |
|-----------|---------|---------------|
| Overwrite running binary | ❌ Blocked (file lock) | ✅ Works |
| Rename running binary | ✅ Allowed | ✅ Works |
| Delete running binary | ❌ Blocked | ✅ Works |
| Replace after rename | ✅ Works | ✅ Works |

**Key insight**: Windows allows **renaming** a running executable but
not **overwriting** or **deleting** it. This is the foundation of the
rename-first deploy strategy.

---

## Two Update Strategies

### Strategy 1: Source-Based Update (Build from Repo)

Used when the binary was installed from a source repository:

```
1. Resolve the source repo location
2. Pull latest code
3. Resolve dependencies
4. Build new binary
5. Deploy to the installed location (rename-first)
6. Verify version
7. Clean up
```

**Requires**: Source repository accessible on the machine, Go toolchain
installed.

**Advantage**: The user always gets a binary built for their exact
platform with embedded repo path.

### Strategy 2: Binary-Based Update (Download Pre-Built)

Used when the binary was installed via an installer script or package
manager (no source repo available):

```
1. Fetch the latest version tag from GitHub API
2. Download the install script for the platform
3. Execute the install script (handles download + verification + install)
4. Verify version
5. Clean up
```

**Requires**: Internet access, GitHub releases with install scripts.

**Advantage**: No build tools needed on the user's machine.

---

## Update Command Flow

```
<binary> update
    │
    ├── Can we find the source repo?
    │   ├── YES → Source-Based Update (build from repo)
    │   └── NO  → Binary-Based Update (download pre-built)
    │
    ├── Platform?
    │   ├── Windows → Handoff mechanism (copy self, run worker)
    │   └── Unix    → Direct in-place update
    │
    ├── Deploy new binary (rename-first)
    ├── Verify: <binary> version == expected
    └── Cleanup temporary files
```

---

## Source Repository Resolution

When using source-based updates, the tool must locate the source repo.
Use a multi-tier resolution strategy:

| Priority | Source | Description |
|----------|--------|-------------|
| 1 | CLI flag | `--repo-path /path/to/repo` |
| 2 | Embedded constant | Built into the binary via `-ldflags` at compile time |
| 3 | Database/config | Stored path from a previous update |
| 4 | Interactive prompt | Ask the user for the path |

If the provided path does not exist, offer to clone the repository:

```
Source repo not found at /path/to/repo.
Clone from https://github.com/org/repo? [Y/n]
```

### Path Validation

After resolving a path, validate it contains the expected project:

1. Check for `go.mod` (or equivalent manifest).
2. Verify the module name matches the expected module.
3. Check for the build script (`run.ps1` / `run.sh`).

---

## Version Comparison

After deploying the new binary, compare versions:

```bash
old_version="1.2.0"
new_version=$(<binary> version)

if [[ "$new_version" == "$old_version" ]]; then
    echo "  !! Warning: version unchanged after update ($old_version)"
    echo "     The source may already be up to date."
fi
```

Always normalize versions before comparing (strip `v` prefix, trim
whitespace).

---

## Error Handling

| Error | Response |
|-------|----------|
| Source repo not found | Prompt user or fall back to binary-based update |
| Git pull fails | Print error, exit without building |
| Build fails | Print error, exit — previous binary is untouched |
| Deploy fails (all retries exhausted) | Rollback from `.old` backup |
| Version unchanged after update | Warn but don't fail |
| No internet for binary-based update | Print error, suggest source-based |

---

## Constraints

- The update command must be **synchronous** — the user sees all output
  in the same terminal session.
- Never leave the system in a state without a working binary.
- Always attempt rollback on deploy failure.
- The update must work from any working directory (not just the repo).
- No interactive prompts during the build/deploy phase (only during
  initial repo resolution if needed).

## Application-Specific References

The following app-level specs apply this generic overview:

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Full build pipeline, config, deploy, and run patterns |
| [03-self-update-mechanism.md](../03-general/03-self-update-mechanism.md) | Handoff, skip-if-current, version verification, rollback |
| [09-build-deploy.md](../01-app/09-build-deploy.md) | gitmap-specific build, deploy, and `-R` flag behavior |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
