import re

path = "gitmap/cmd/rootusage_groups.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

funcs_add = """
func printGroupInstallers() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupInstallers))
	fmt.Println()
	fmt.Println(constants.HelpInstall)
	fmt.Println(constants.HelpInstaller)
	fmt.Println(constants.HelpMacro)
	fmt.Println(constants.HelpSetup)
}

func printGroupIntegrations() {
	fmt.Println()
	fmt.Println(colorGroupHeader(constants.HelpGroupIntegrations))
	fmt.Println()
	fmt.Println(constants.HelpVSCode)
	fmt.Println(constants.HelpAgy)
	fmt.Println(constants.HelpSchedule)
}
"""

content += funcs_add

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
