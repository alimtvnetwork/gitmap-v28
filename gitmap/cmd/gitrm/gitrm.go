package gitrm

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/pterm/pterm"
)

// Run executes the git-rm command logic.
func Run(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("GitRmRun", "E_GITRM_MISSING_INPUT")
	}
	input := args[0]
	paths, err := parseInput(input)
	if err != nil {
		return apperror.WrapSimple(err, "git-rm: failed to parse input")
	}
	if len(paths) == 0 {
		return apperror.NewSimple("GitRmRun", "E_GITRM_NO_PATHS")
	}

	// 1. Backup paths to global location
	backupDir, err := createBackupDir()
	if err != nil {
		return apperror.WrapSimple(err, "git-rm: failed to create backup directory")
	}
	for _, p := range paths {
		backupFile(p, backupDir)
	}

	// 2. Remove from git history
	return rewriteHistory(paths)
}

func parseInput(input string) ([]string, error) {
	st, err := os.Stat(input)
	if err == nil && !st.IsDir() {
		ext := strings.ToLower(filepath.Ext(input))
		b, err := os.ReadFile(input)
		if err != nil {
			return nil, err
		}
		if ext == ".json" {
			var paths []string
			if err := json.Unmarshal(b, &paths); err == nil {
				return paths, nil
			}
		} else if ext == ".csv" {
			r := csv.NewReader(strings.NewReader(string(b)))
			records, err := r.ReadAll()
			if err == nil {
				var paths []string
				for _, rec := range records {
					paths = append(paths, rec...)
				}
				return paths, nil
			}
		} else {
			// Plain text file (lines)
			lines := strings.Split(string(b), "\n")
			var paths []string
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if l != "" {
					paths = append(paths, l)
				}
			}
			return paths, nil
		}
	}
	// Direct path or folder
	// Or comma-separated
	if strings.Contains(input, ",") {
		return strings.Split(input, ","), nil
	}
	return []string{input}, nil
}

func createBackupDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cwd, _ := os.Getwd()
	repoName := filepath.Base(cwd)
	bd := filepath.Join(home, ".gitmap", "backups", "git-rm", repoName)
	err = os.MkdirAll(bd, 0755)
	return bd, err
}

func backupFile(path, backupRoot string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return // Ignore missing files in working tree
	}
	dest := filepath.Join(backupRoot, path)
	os.MkdirAll(filepath.Dir(dest), 0755)
	os.WriteFile(dest, b, 0644)
}

func rewriteHistory(paths []string) error {
	pterm.Info.Printf("Removing %d files from git history...\n", len(paths))
	var quoted []string
	for _, p := range paths {
		quoted = append(quoted, fmt.Sprintf("'%s'", p))
	}
	filesStr := strings.Join(quoted, " ")
	filterCmd := fmt.Sprintf("git rm --cached --ignore-unmatch %s", filesStr)

	cmd := exec.Command("git", "filter-branch", "--force", "--index-filter", filterCmd, "--prune-empty", "--tag-name-filter", "cat", "--", "--all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "git-rm: history rewrite failed")
	}

	// Clean up refs
	exec.Command("git", "for-each-ref", "--format=%(refname)", "refs/original/").Run()
	exec.Command("git", "reflog", "expire", "--expire=now", "--all").Run()
	exec.Command("git", "gc", "--prune=now", "--aggressive").Run()

	pterm.Success.Println("Git history successfully rewritten.")
	return nil
}
