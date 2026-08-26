import re

path = "gitmap/cmd/code.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

repl = """	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "paths":
			runCodePaths(args[1:])
			return
		case "pap", "prompt-all-project":
			fmt.Println("Feature [vscode pap] is not yet implemented")
			return
		case "plugins", "plugin":
			fmt.Println("Feature [vscode plugins] is not yet implemented")
			return
		case "add-project", "ap":
			fmt.Println("Feature [vscode add-project] is not yet implemented")
			return
		case "ls", "list":
			fmt.Println("Feature [vscode ls] is not yet implemented")
			return
		case "rm", "remove", "delete", "del":
			fmt.Println("Feature [vscode rm] is not yet implemented")
			return
		}
	}"""

content = content.replace("""	if len(args) > 0 && args[0] == "paths" {
		runCodePaths(args[1:])

		return
	}""", repl)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
