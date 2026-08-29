import os
import re

with open('gitmap/cmd/update.go', 'r', encoding='utf-8') as f:
    t = f.read()

# I will just write a regex to replace the function body properly.
func_body = """func resolveRepoPath() (string, error) {
	for _, path := range []string{
		resolveRepoPathFromFlag(),
		resolveRepoPathFromEmbedded(),
		resolveRepoPathFromDB(),
	} {
		if len(path) > 0 {
			saveRepoPathToDB(path)
			return path, nil
		}
	}

	if prompted := promptRepoPath(); len(prompted) > 0 {
		saveRepoPathToDB(prompted)
		return prompted, nil
	}

	// Try to fall back to gitmap-updater for release-based update
	if tryUpdaterFallback() {
		os.Exit(0)
	}

	return "", apperror.NewSimple("no repo path resolved", "E9024")
}"""

t = re.sub(r'func resolveRepoPath\(\).*?^}', func_body, t, flags=re.MULTILINE|re.DOTALL)

with open('gitmap/cmd/update.go', 'w', encoding='utf-8') as f:
    f.write(t)
