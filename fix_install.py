import re

with open("gitmap/cmd/install.go", "r", encoding="utf-8") as f:
    c = f.read()

c = c.replace(
    'constants.ToolMacroAhk:         func(installOptions) { runInstallCustomTool("macro-ahk") },\n\t}[tool]',
    'constants.ToolMacroAhk:         func(installOptions) { runInstallCustomTool("macro-ahk") },\n\t\tconstants.ToolAgManager:        func(installOptions) { runInstallAgManager() },\n\t}[tool]'
)

c = c.replace(
    'if _, exists := constants.InstallToolDescriptions[tool]; exists {\n\t\treturn\n\t}',
    'if tool == "ag-m" {\n\t\treturn\n\t}\n\tif _, exists := constants.InstallToolDescriptions[tool]; exists {\n\t\treturn\n\t}'
)

c = c.replace(
    'func executeInstall(opts installOptions) {\n\tif handler := specialInstallHandler(opts.Tool); handler != nil {',
    'func executeInstall(opts installOptions) {\n\tif opts.Tool == "ag-m" { opts.Tool = constants.ToolAgManager }\n\tif handler := specialInstallHandler(opts.Tool); handler != nil {'
)

with open("gitmap/cmd/install.go", "w", encoding="utf-8") as f:
    f.write(c)
print("done")
