package probe

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// tryShallowClone is the heavyweight fallback. Clones into a temp dir with
// `--depth N --filter=blob:none --no-checkout` (treeless, no working copy),
// runs `git tag --sort=-v:refname`, and returns the top result. depth<1 is
// coerced to 1 so we always make a non-fatal request to the remote.
func tryShallowClone(url string, depth int) (string, error) {
	if depth < 1 {
		depth = 1
	}
	tmp, err := os.MkdirTemp("", "gitmap-probe-*")
	if err != nil {
		return "", fmt.Errorf("mkdtemp: %w", err)
	}
	defer os.RemoveAll(tmp)

	target := filepath.Join(tmp, "repo")
	clone := exec.Command("git", "clone",
		"--depth", fmt.Sprintf("%d", depth),
		"--filter=blob:none",
		"--no-checkout",
		url, target,
	)
	if out, err := clone.CombinedOutput(); err != nil {
		return "", fmt.Errorf(constants.ErrProbeCloneFail, summarize(out, err))
	}

	tags := exec.Command("git", "-C", target, "tag", "--sort=-v:refname")
	out, err := tags.Output()
	if err != nil {
		return "", fmt.Errorf("git tag: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		t := strings.TrimSpace(line)
		if t != "" {
			return t, nil
		}
	}

	return "", nil
}

// summarize folds combined-output bytes into a single-line error reason.
func summarize(out []byte, err error) string {
	tail := strings.TrimSpace(string(out))
	if tail == "" {
		return err.Error()
	}
	if idx := strings.LastIndex(tail, "\n"); idx >= 0 {
		tail = tail[idx+1:]
	}

	return tail
}
