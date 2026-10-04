#!/usr/bin/env python3
"""
49-verify-privacy-and-relative-paths.py

Audits specification documents and AI memory plans for privacy leaks (raw IPs)
and workstation-specific absolute filesystem paths.
Enforces zero IP leaks and relative/portable paths across the repository.
"""

import os
import re
import sys
from pathlib import Path

TARGET_31_FILES = [
    # Group A: Specification Documents (02-spec/21-app/)
    "02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/01-architecture-spec.md",
    "02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/02-component-spec.md",
    "02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/02-full-os-setup-blueprint.md",
    "02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md",
    "02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/02-full-os-setup-blueprint.md",
    "02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md",
    "02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/02-gitmap-macro-and-installer-spec.md",
    "02-spec/21-app/213-antigravity-ubuntu-update-and-macro-automation/03-in-app-update-button-rca.md",
    "02-spec/21-app/214-ubuntu-fleet-automation-and-workstation-governance/01-architecture-spec.md",
    "02-spec/21-app/214-ubuntu-fleet-automation-and-workstation-governance/02-component-and-cli-spec.md",
    "02-spec/21-app/215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager/01-architecture-spec.md",
    "02-spec/21-app/215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager/02-component-and-cli-spec.md",
    "02-spec/21-app/216-gitmap-prompting-freeze-and-suggestion-engine-fix/01-architecture-spec.md",
    "02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md",
    "02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/02-component-and-cli-spec.md",
    "02-spec/21-app/182-version-pinning-macro-deploy-ui-settings-secret-flags.md",
    # Group B: Root Plans & Completed Plans (.ai-memory/plans/)
    ".ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md",
    ".ai-memory/plans/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md",
    ".ai-memory/plans/completed/218-nodes-agy-ui-remote-settings-and-cursor-automation.md",
    ".ai-memory/plans/completed/215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager.md",
    ".ai-memory/plans/completed/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md",
    # Group C: Subtask Execution Plans (.ai-memory/plans/subtasks/)
    ".ai-memory/plans/subtasks/212-ubuntu-fleet-full-customization-and-embedded-runner/01-desktop-wallpaper-and-gui-app-launch.md",
    ".ai-memory/plans/subtasks/212-ubuntu-fleet-full-customization-and-embedded-runner/02-vmware-shared-folder-deep-verification.md",
    ".ai-memory/plans/subtasks/212-ubuntu-fleet-full-customization-and-embedded-runner/03-antigravity-deep-conversation-and-settings-sync.md",
    ".ai-memory/plans/subtasks/212-ubuntu-fleet-full-customization-and-embedded-runner/04-self-contained-embedded-powershell-runner.md",
    ".ai-memory/plans/subtasks/212-ubuntu-fleet-full-customization-and-embedded-runner/05-detailed-engineering-log-and-future-roadmap.md",
    ".ai-memory/plans/subtasks/214-ubuntu-fleet-automation-and-workstation-governance/01-git-workspaces-clone-and-validation.md",
    ".ai-memory/plans/subtasks/214-ubuntu-fleet-automation-and-workstation-governance/02-gnome-ergonomics-and-keybindings.md",
    ".ai-memory/plans/subtasks/214-ubuntu-fleet-automation-and-workstation-governance/03-vmware-automount-and-gui-launching.md",
    ".ai-memory/plans/subtasks/214-ubuntu-fleet-automation-and-workstation-governance/04-antigravity-upgrade-and-deep-brain-migration.md",
    ".ai-memory/plans/subtasks/214-ubuntu-fleet-automation-and-workstation-governance/05-master-embedded-runner-and-scorecard.md",
]

# Patterns that represent violations
LEAK_PATTERNS = [
    ("Raw IP Address (192.168.x.x)", re.compile(r"\b192\.168\.\d{1,3}\.\d{1,3}\b")),
    ("Absolute Windows Work Path (d:\\work / d:/work)", re.compile(r"(?i)\b[dD]:[\\/]work\b")),
    ("Absolute Windows User Path (C:\\Users)", re.compile(r"(?i)\b[cC]:[\\/]users[\\/]administrator\b")),
    ("Hardcoded Linux Home (/home/a/)", re.compile(r"/home/a/(?:git-work|gitmap|\.gemini|\.antigravity|\.ssh|\.local|[a-zA-Z0-9_\-\.]+)")),
]

# Whitelist patterns: legitimate regex documentation or markdown table definitions in specifications
EXCLUSION_PATTERNS = [
    re.compile(r"Search Pattern.*Mandatory Replacement"),
    re.compile(r"`192\.168\.\d{1,3}\.\d{1,3}`\s*\|\s*`"),
    re.compile(r"re\.compile\("),
    re.compile(r"LEAK_PATTERNS"),
]

def is_line_excluded(line: str) -> bool:
    for pat in EXCLUSION_PATTERNS:
        if pat.search(line):
            return True
    return False

def audit_file(file_path: Path) -> list:
    violations = []
    if not file_path.exists():
        violations.append((0, "File Not Found", str(file_path), ""))
        return violations

    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception as e:
        violations.append((0, f"Read Error: {e}", "", ""))
        return violations

    for line_num, line in enumerate(content.splitlines(), start=1):
        if is_line_excluded(line):
            continue
        for name, pattern in LEAK_PATTERNS:
            match = pattern.search(line)
            if match:
                violations.append((line_num, name, match.group(0), line.strip()[:100]))
    return violations

def sanitize_content(text: str) -> str:
    """Applies the deterministic sanitization matrix to text."""
    # 1. IP address sanitization
    # 192.168.1.22 -> ubuntu-fleet-01
    text = text.replace("192.168.1.22", "ubuntu-fleet-01")
    # 192.168.1.10 -> node-main
    text = text.replace("192.168.1.10", "node-main")
    # 192.168.1.20 -> worker-1
    text = text.replace("192.168.1.20", "worker-1")
    # 192.168.1.14 -> devbox
    text = text.replace("192.168.1.14", "devbox")
    # 192.168.1.3 -> node-w1
    text = text.replace("192.168.1.3", "node-w1")
    # 192.168.1.7 -> node-w2
    text = text.replace("192.168.1.7", "node-w2")
    # 192.168.1.12 -> node-w3
    text = text.replace("192.168.1.12", "node-w3")

    # Catch any remaining 192.168.1.x (except documentation tables/regex)
    def ip_replacer(match):
        ip = match.group(0)
        last_octet = ip.split(".")[-1]
        return f"node-{last_octet}"
    # only if not in table
    # We will do regex for any remaining 192.168.1.x
    text = re.sub(r"\b192\.168\.1\.\d+\b", ip_replacer, text)

    # 2. Secrets directory
    text = re.sub(r"(?i)[dD]:[\\/]work[\\/]repo-secrets[\\/]", "$SECRETS_DIR/", text)
    text = re.sub(r"(?i)[dD]:[\\/]work[\\/]repo-secrets\b", "$SECRETS_DIR", text)

    # 3. Markdown links file:///d:/work/gitmap/ and file:///d%3A/work/gitmap/
    text = re.sub(r"file:///[dD]:[\\/]work[\\/]gitmap[\\/]", "./", text)
    text = re.sub(r"file:///[dD]%3A[\\/]work[\\/]gitmap[\\/]", "./", text)
    text = re.sub(r"file:///[dD]:[\\/]work[\\/]gitmap\b", "./", text)
    text = re.sub(r"file:///[dD]%3A[\\/]work[\\/]gitmap\b", "./", text)

    # Markdown links file:///d:/work/ and file:///d%3A/work/
    text = re.sub(r"file:///[dD]:[\\/]work[\\/]", "$WORKSPACE_DIR/", text)
    text = re.sub(r"file:///[dD]%3A[\\/]work[\\/]", "$WORKSPACE_DIR/", text)

    # 4. Workspace root: d:\work\gitmap, d:/work/gitmap, etc.
    text = re.sub(r"(?i)\bd:[\\/]work[\\/]gitmap\b", "$WORKSPACE_ROOT", text)
    text = re.sub(r"(?i)\bd:\\\\work\\\\gitmap\b", "$WORKSPACE_ROOT", text)

    # 5. Workspace parent: d:\work\, d:/work/, etc.
    text = re.sub(r"(?i)\bd:[\\/]work[\\/]", "$WORKSPACE_DIR/", text)
    text = re.sub(r"(?i)\bd:[\\/]work\b", "$WORKSPACE_DIR", text)
    text = re.sub(r"(?i)\bd:\\\\work\\\\", "$WORKSPACE_DIR/", text)

    # 6. Windows Administrator path
    text = re.sub(r"(?i)\bc:[\\/]users[\\/]administrator[\\/]", "$USERPROFILE/", text)
    text = re.sub(r"(?i)\bc:[\\/]users[\\/]administrator\b", "$USERPROFILE", text)
    text = re.sub(r"(?i)\bc:\\\\users\\\\administrator\\\\", "$USERPROFILE/", text)

    # 7. Linux home /home/a/
    text = re.sub(r"/home/a/git-work", "$HOME/git-work", text)
    text = re.sub(r"/home/a/\.gemini", "$HOME/.gemini", text)
    text = re.sub(r"/home/a/\.antigravity_tools", "$HOME/.antigravity_tools", text)
    text = re.sub(r"/home/a/\.antigravity", "$HOME/.antigravity", text)
    text = re.sub(r"/home/a/\.ssh", "$HOME/.ssh", text)
    text = re.sub(r"/home/a/\.local", "$HOME/.local", text)
    text = re.sub(r"/home/a/\.config", "$HOME/.config", text)
    text = re.sub(r"/home/a/Pictures", "$HOME/Pictures", text)
    text = re.sub(r"/home/a/", "$HOME/", text)

    return text

def fix_all_files(root: Path):
    print("=" * 72)
    print(" Applying sanitization across all 31 target files...")
    print("=" * 72)
    for rel_path in TARGET_31_FILES:
        target = root / rel_path
        if not target.exists():
            print(f"[SKIP] Not found: {rel_path}")
            continue
        original = target.read_text(encoding="utf-8", errors="replace")
        sanitized = sanitize_content(original)
        if original != sanitized:
            target.write_text(sanitized, encoding="utf-8")
            print(f"[FIXED] {rel_path}")
        else:
            print(f"[UNCHANGED] {rel_path}")

def main():
    root = Path(__file__).resolve().parent.parent
    if "--fix" in sys.argv:
        fix_all_files(root)

    total_violations = 0
    checked_files = 0

    print("=" * 72)
    print(" 49-verify-privacy-and-relative-paths: Privacy & Absolute Path Audit")
    print("=" * 72)

    for rel_path in TARGET_31_FILES:
        target = root / rel_path
        checked_files += 1
        violations = audit_file(target)
        if violations:
            print(f"\n[FAIL] {rel_path} ({len(violations)} violations):")
            for line_num, v_type, matched, context in violations:
                print(f"  L{line_num:4d}: {v_type} -> '{matched}' | {context}")
                total_violations += 1
        else:
            print(f"[PASS] {rel_path}")

    print("\n" + "-" * 72)
    print(f"Summary: Checked {checked_files} files. Total violations found: {total_violations}")
    print("-" * 72)

    if total_violations > 0:
        print("[RESULT] FAILED: Privacy leaks or absolute paths detected.")
        sys.exit(1)
    else:
        print("[RESULT] SUCCESS: Zero privacy leaks and zero absolute paths detected.")
        sys.exit(0)

if __name__ == "__main__":
    main()
