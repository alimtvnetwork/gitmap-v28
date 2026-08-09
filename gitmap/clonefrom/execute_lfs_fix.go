package clonefrom

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	rxSmudgeFatal = regexp.MustCompile(`(?m)^fatal: (.*?): smudge filter lfs failed`)
	rxSmudgeError = regexp.MustCompile(`(?m)^Error downloading object: ([^\s]+)`)
)

// detectLFSSmudgeError scans the git clone standard error for a known
// Git LFS 404 smudge signature and returns the offending file path if found.
func detectLFSSmudgeError(output string) (string, bool) {
	if !strings.Contains(output, "smudge filter lfs failed") && !strings.Contains(output, "Object does not exist on the server") {
		return "", false
	}

	if match := rxSmudgeFatal.FindStringSubmatch(output); len(match) == 2 {
		return match[1], true
	}

	if match := rxSmudgeError.FindStringSubmatch(output); len(match) == 2 {
		return match[1], true
	}

	return "", false
}

// executeLFSFix applies the destructive fallback sequence to bypass a broken LFS pointer.
// It skips the smudge filter, restores HEAD, removes the cached pointer, and pushes.
func executeLFSFix(dest string, file string) error {
	// 1. git restore --source=HEAD :/
	cmdRest := exec.Command(constants.GitBin, "restore", "--source=HEAD", ":/")
	cmdRest.Dir = dest
	cmdRest.Env = append(os.Environ(), "GIT_LFS_SKIP_SMUDGE=1")
	if out, err := cmdRest.CombinedOutput(); err != nil {
		return fmt.Errorf("git restore failed: %v\n%s", err, string(out))
	}

	// 2. git rm --cached "<file>" -q
	cmdRm := exec.Command(constants.GitBin, "rm", "--cached", file, "-q")
	cmdRm.Dir = dest
	cmdRm.Env = append(os.Environ(), "GIT_LFS_SKIP_SMUDGE=1")
	if out, err := cmdRm.CombinedOutput(); err != nil {
		return fmt.Errorf("git rm --cached failed: %v\n%s", err, string(out))
	}

	// 3. Remove physical file
	_ = os.Remove(filepath.Join(dest, file))

	// 4. git commit -m "..."
	cmdCommit := exec.Command(constants.GitBin, "commit", "-m", "chore(lfs): remove pointer for missing LFS object")
	cmdCommit.Dir = dest
	cmdCommit.Env = append(os.Environ(), "GIT_LFS_SKIP_SMUDGE=1")
	if out, err := cmdCommit.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %v\n%s", err, string(out))
	}

	// 5. git push
	cmdPush := exec.Command(constants.GitBin, "push")
	cmdPush.Dir = dest
	cmdPush.Env = append(os.Environ(), "GIT_LFS_SKIP_SMUDGE=1")
	if out, err := cmdPush.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %v\n%s", err, string(out))
	}

	return nil
}
