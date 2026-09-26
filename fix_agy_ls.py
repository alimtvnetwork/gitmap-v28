import re

with open('cli/cmdagy/agy_ls.go', 'r') as f:
    content = f.read()

content = content.replace('''
		if len(args) > 0 {
			if val, err := strconv.Atoi(args[0]); err == nil && val > 0 {
				n = val
			}
		}
''', '''
		n = parseArgCount(args, n)
''')

content = content.replace('''
func renderAgyLsResult(filtered []AgyProject, dirPath string) error {
	if agyLsJSON || agyLsFile != "" {
		if agyLsFile != "" {
			return outputAgyProjectsJSONFile(filtered, agyLsFile)
		}
		return outputAgyProjectsJSON(filtered)
	}

	renderAgyProjectsTable(filtered, dirPath)
	return nil
}
''', '''
func renderAgyLsResult(filtered []AgyProject, dirPath string) error {
	if agyLsFile != "" {
		return outputAgyProjectsJSONFile(filtered, agyLsFile)
	}
	if agyLsJSON {
		return outputAgyProjectsJSON(filtered)
	}

	renderAgyProjectsTable(filtered, dirPath)
	return nil
}

func parseArgCount(args []string, defaultN int) int {
	if len(args) == 0 {
		return defaultN
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val > 0 {
		return val
	}
	return defaultN
}
''')

with open('cli/cmdagy/agy_ls.go', 'w') as f:
    f.write(content)
