import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# I will just remove the original ones and keep the ones at the bottom, or vice-versa.
# Let's remove the ones at the bottom and keep the original ones, but append `(use --help to expand)`.
content = re.sub(r'\tHelpInstall\s*=.*?\n', '', content)
# wait, there's `HelpInstaller` at line 57, `HelpMacro` at line 37
# The bottom ones I added are:
bottom_block = """	HelpInstall           = "  install (i) [tools...]     Run setup profiles or install specific tools (tree view) (use --help to expand)"
	HelpInstaller         = "  installer (in) [cmd]       Manage installer scripts and history (tree view) (use --help to expand)"
	HelpMacro             = "  macro (m) [cmd]            Manage and execute macros (tree view) (use --help to expand)"
	HelpSchedule          = "  schedule (sc) [cmd]        Schedule tasks and run jobs asynchronously (tree view) (use --help to expand)"
	HelpVSCode            = "  vscode (vsc) [cmd]         Manage VS Code Project Manager integrations (use --help to expand)"
	HelpAgy               = "  antigravity (ag) [cmd]     Manage Google Antigravity workspaces and plugins (use --help to expand)"

"""
content = content.replace(bottom_block, "")

# Now I will just add HelpInstall, HelpSchedule, HelpVSCode, HelpAgy cleanly.
clean_add = """	HelpInstall         = "  install (i) [tools...]     Run setup profiles or install specific tools (tree view) (use --help to expand)"
	HelpSchedule        = "  schedule (sc) [cmd]        Schedule tasks and run jobs asynchronously (tree view) (use --help to expand)"
	HelpVSCode          = "  vscode (vsc) [cmd]         Manage VS Code Project Manager integrations (use --help to expand)"
	HelpAgy             = "  antigravity (ag) [cmd]     Manage Google Antigravity workspaces and plugins (use --help to expand)"
"""

content = content.replace("	CompactNoMatchFmt = \"  No group matching '%s'. Showing all groups:\\n\"", "	CompactNoMatchFmt = \"  No group matching '%s'. Showing all groups:\\n\"\n" + clean_add)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
