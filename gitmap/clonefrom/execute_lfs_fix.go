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
	rxSmudgeFatal = regexp.MustCompile(constants.LFSSmudgeFatalPattern)
	rxSmudgeError = regexp.MustCompile(constants.LFSSmudgeErrorPattern)
)

// detectLFSSmudgeError scans the git clone standard error for a known
// Git LFS 404 smudge signature and returns the offending file path if found.
func detectLFSSmudgeError(output string) (string, bool) {
	hasSmudgeErr := strings.Contains(output, constants.LFSSmudgeFilterFailedSignature)
	hasServerErr := strings.Contains(output, constants.LFSServerObjectMissingSignature)
	if hasSmudgeErr || hasServerErr {
		if match := rxSmudgeFatal.FindStringSubmatch(output); len(match) == constants.LFSSubmatchLength {
			return match[1], true
		}

		if match := rxSmudgeError.FindStringSubmatch(output); len(match) == constants.LFSSubmatchLength {
			return match[1], true
		}
	}

	return "", false
}

// executeLFSFix applies the destructive fallback sequence to bypass a broken LFS pointer.
// It skips the smudge filter, restores HEAD, removes the cached pointer, and pushes.
func executeLFSFix(dest string, file string) error {
	// 1. git restore --source=HEAD :/
	cmdRest := exec.Command(constants.GitBin, constants.GitRestore, constants.GitRestoreSourceHEAD, constants.GitPathspecRoot)
	cmdRest.Dir = dest
	cmdRest.Env = append(os.Environ(), constants.EnvGitLFSSkipSmudge)
	if out, err := cmdRest.CombinedOutput(); err != nil {
		return fmt.Errorf(constants.ErrCloneFromLFSRestore, err, string(out))
	}

	// 2. git rm --cached "<file>" -q
	cmdRm := exec.Command(constants.GitBin, constants.GitRm, constants.GitCachedFlag, file, constants.GitQuietFlag)
	cmdRm.Dir = dest
	cmdRm.Env = append(os.Environ(), constants.EnvGitLFSSkipSmudge)
	if out, err := cmdRm.CombinedOutput(); err != nil {
		return fmt.Errorf(constants.ErrCloneFromLFSRmCached, err, string(out))
	}

	// 3. Remove physical file
	_ = os.Remove(filepath.Join(dest, file))

	// 4. git commit -m "..."
	cmdCommit := exec.Command(constants.GitBin, constants.GitCommitCmd, constants.GitTagMessageFlag, constants.LFSFixCommitMessage)
	cmdCommit.Dir = dest
	cmdCommit.Env = append(os.Environ(), constants.EnvGitLFSSkipSmudge)
	if out, err := cmdCommit.CombinedOutput(); err != nil {
		return fmt.Errorf(constants.ErrCloneFromLFSCommit, err, string(out))
	}

	// 5. git push
	cmdPush := exec.Command(constants.GitBin, constants.GitPush)
	cmdPush.Dir = dest
	cmdPush.Env = append(os.Environ(), constants.EnvGitLFSSkipSmudge)
	if out, err := cmdPush.CombinedOutput(); err != nil {
		return fmt.Errorf(constants.ErrCloneFromLFSPush, err, string(out))
	}

	return nil
}
