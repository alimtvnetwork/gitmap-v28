import re

path = "gitmap/cmd/vscode_cmd.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# I will add cases to dispatchVSCodeAction
new_cases = """	case "pap", "prompt-all-project":
		fmt.Println("Feature [vscode pap] is not yet implemented")
	case "plugins", "plugin":
		fmt.Println("Feature [vscode plugins] is not yet implemented")
"""
content = content.replace('	case "add":', new_cases + '	case "add", "add-project", "ap":')
content = content.replace('	case "rm", "remove", "delete":', '	case "rm", "remove", "delete", "del":')

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
