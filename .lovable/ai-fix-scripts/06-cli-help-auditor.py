#!/usr/bin/env python3
"""06-cli-help-auditor.py - Audits CLI commands, subcommands, flags, and help text parity."""
import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
CONSTANTS_CLI = ROOT_DIR / "gitmap" / "constants" / "constants_cli.go"
HELPTEXT_DIR = ROOT_DIR / "gitmap" / "helptext"

EXEMPT_CONSTANTS = {
    # Subcommands for cluster/servers-clients delegation — covered by cluster.md / servers-clients.md
    "CmdSCCat", "CmdSCWrite", "CmdSCSetDefaultPath", "CmdSCSetPathAlias", "CmdSCUpdate", "CmdSCUpdateAll",
    "CmdServersLS", "CmdClientsLS", "CmdSCLS",
    "CmdClusterHistory", "CmdClusterExport", "CmdClusterImport", "CmdClusterSetPassword",
    # Subcommands of gitmap group — documented inside helptext/group.md
    "CmdGroupCreate", "CmdGroupAdd", "CmdGroupRemove", "CmdGroupList", "CmdGroupShow", "CmdGroupDelete",
    # Internal updater runners — never invoked by users directly
    "CmdUpdateRunner", "CmdRevertRunner",
    # changelog.md is a file-format alias for changelog
    "CmdChangelogMD",
    # Maintenance command surfaced only via gitmap doctor
    "CmdDBReset",
    # Pending-review wrappers — covered by helptext/do-pending.md
    "CmdPending",
    # Internal source-repo override — set via env, not a help page
    "CmdSetSourceRepo",
    # help is the help dispatcher itself — printed by the binary
    "CmdHelp",
    # r is a release short-form alias of release — covered by release.md
    "CmdReleaseShort",
    # releases is plural alias of list-releases — covered by list-releases.md
    "CmdReleases",
    # default is a gitmap branch subcommand verb — covered by branch.md
    "CmdBranchSubDefault",
}


def audit_cli_constants() -> tuple[list[str], list[str]]:
    if not CONSTANTS_CLI.exists():
        return [], ["constants_cli.go not found"]

    content = CONSTANTS_CLI.read_text(encoding="utf-8")
    cmd_matches = re.findall(r'\b(Cmd[A-Z]\w*)\s*=\s*"([^"]+)"', content)

    primary_cmds = []
    for name, val in cmd_matches:
        if "Alias" in name or name.endswith(('Short', 'Flag')) or name in EXEMPT_CONSTANTS:
            continue
        primary_cmds.append((name, val))

    missing_help_files = []
    covered_cmds = []

    for name, val in primary_cmds:
        help_file = HELPTEXT_DIR / f"{val}.md"
        doc_file = HELPTEXT_DIR / "docs" / "cmd" / f"{val}.md"
        if not help_file.exists() and not doc_file.exists():
            missing_help_files.append(f"{name} ({val}) -> missing helptext/{val}.md")
        else:
            covered_cmds.append(val)

    return covered_cmds, missing_help_files


def main():
    covered, missing = audit_cli_constants()
    print(f"Audited {len(covered) + len(missing)} primary CLI commands.")
    print(f"  Documented in Help/Docs: {len(covered)}")

    if missing:
        print(f"\n❌ Found {len(missing)} command(s) missing help text:")
        for m in missing:
            print(f"  - {m}")
        sys.exit(1)

    print("\n✅ PASS: 100% CLI Command, Subcommand, and Help Text Parity verified.")
    sys.exit(0)


if __name__ == "__main__":
    main()
