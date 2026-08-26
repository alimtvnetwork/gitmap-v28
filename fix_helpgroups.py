import re

path = "gitmap/constants/constants_helpgroups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# Add to Help group headers
help_group_add = """	HelpGroupTemplates    = "  Templates & Scaffolding (.gitignore / .gitattributes / LFS):"
	HelpGroupCluster      = "  Cluster & Delegation (multi-machine networks):"
	HelpGroupInstallers   = "  Installers & Macros:"
	HelpGroupIntegrations = "  Integrations (VS Code, Antigravity, Scheduler):"
"""
content = content.replace('	HelpGroupTemplates    = "  Templates & Scaffolding (.gitignore / .gitattributes / LFS):"\n\tHelpGroupCluster      = "  Cluster & Delegation (multi-machine networks):"\n', help_group_add)

commands_add = """	HelpInstall           = "  install (i) [tools...]     Run setup profiles or install specific tools (tree view)"
	HelpInstaller         = "  installer (in) [cmd]       Manage installer scripts and history (tree view)"
	HelpMacro             = "  macro (m) [cmd]            Manage and execute macros (tree view)"
	HelpSchedule          = "  schedule (sc) [cmd]        Schedule tasks and run jobs asynchronously (tree view)"
	HelpVSCode            = "  vscode (vsc) [cmd]         Manage VS Code Project Manager integrations"
	HelpAgy               = "  antigravity (ag) [cmd]     Manage Google Antigravity workspaces and plugins"
"""

# Insert these anywhere in the const block. Let's just put them near the end of the first const block.
# We'll just replace `	HelpCluster           = "  cluster (cl)               Multi-machine remote sync dashboard"`
# Actually, I'll just append it before `)` of the first const block.
content = re.sub(r'(\n\))', "\n" + commands_add + r'\1', content, count=1)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
