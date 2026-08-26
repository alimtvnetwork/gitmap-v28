import re

with open("gitmap/constants/constants_install.go", "r", encoding="utf-8") as f:
    c = f.read()

c = c.replace(
    'ToolMacroAhk         = "macro-ahk"\n)',
    'ToolMacroAhk         = "macro-ahk"\n\tToolAgManager        = "ag-manager"\n)'
)

c = c.replace(
    'ToolMacroAhk:         "Install Macro AHK scripts",\n}',
    'ToolMacroAhk:         "Install Macro AHK scripts",\n\tToolAgManager:        "Install Antigravity Manager GUI",\n}'
)

c = c.replace(
    'ToolScriptsFixer, ToolCodingGuidelines, ToolMacroAhk,\n\t},',
    'ToolScriptsFixer, ToolCodingGuidelines, ToolMacroAhk, ToolAgManager,\n\t},'
)

with open("gitmap/constants/constants_install.go", "w", encoding="utf-8") as f:
    f.write(c)

print("done")
