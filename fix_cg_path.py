import re

with open('cli/cmd/repo_create_init.go', 'r') as f:
    content = f.read()

# Replace cgSrc := filepath.Join("02-spec", "02-coding-guidelines")
replacement = """
	baseDir := "."
	if constants.RepoPath != "" {
		baseDir = constants.RepoPath
	} else {
		exe, err := os.Executable()
		if err == nil {
			baseDir = filepath.Dir(filepath.Dir(exe)) // ../bin/gitmap.exe -> ../
		}
	}
	cgSrc := filepath.Join(baseDir, "02-spec", "02-coding-guidelines")
"""

content = content.replace('\tcgSrc := filepath.Join("02-spec", "02-coding-guidelines")', replacement)

# Add os to imports if not there
if '"os"' not in content:
    content = content.replace('"os/exec"', '"os"\n\t"os/exec"')

if '"github.com/alimtvnetwork/gitmap-v28/cli/constants"' not in content:
    content = content.replace('"github.com/alimtvnetwork/gitmap-v28/cli/apperror"', '"github.com/alimtvnetwork/gitmap-v28/cli/apperror"\n\t"github.com/alimtvnetwork/gitmap-v28/cli/constants"')


with open('cli/cmd/repo_create_init.go', 'w') as f:
    f.write(content)
