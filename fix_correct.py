import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

correct = """
	HelpInstaller         = "  installer (in) [cmd]       Manage installer scripts and history (tree view) (use --help to expand)"
	HelpMacro             = "  macro (m) [cmd]            Manage and execute macros (tree view) (use --help to expand)"
	HelpSchedule          = "  schedule (sc) [cmd]        Schedule tasks and run jobs asynchronously (tree view) (use --help to expand)"
	HelpVSCode            = "  vscode (vsc) [cmd]         Manage VS Code Project Manager integrations (use --help to expand)"
	HelpAgy               = "  antigravity (ag) [cmd]     Manage Google Antigravity workspaces and plugins (use --help to expand)"
"""

content = content.replace('CompactNoMatchFmt = "  No group matching \'%s\'. Showing all groups:\\n"', 'CompactNoMatchFmt = "  No group matching \'%s\'. Showing all groups:\\n"' + correct)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

path2 = "gitmap/constants/constants_install.go"
with open(path2, "r", encoding="utf-8") as f:
    content2 = f.read()

content2 = content2.replace('HelpInstall   = "  install (i) [tools...]     Run setup profiles or install specific tools (tree view)"', 'HelpInstall   = "  install (i) [tools...]     Run setup profiles or install specific tools (tree view) (use --help to expand)"')

with open(path2, "w", encoding="utf-8") as f:
    f.write(content2)

print("done")
