package e2e

import (
	"fmt"
	"os/exec"
	"strings"
)

// finalize streams every recorded commit into `git fast-import`, then
// creates each lightweight tag against the resolved commit. Must be
// called once, after all commit/tag calls and before any read.
func (r *gitRepo) finalize() {
	r.t.Helper()

	cmd := exec.Command("git", "-C", r.dir, "fast-import", "--quiet")
	cmd.Stdin = strings.NewReader(r.stream.String())

	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("fast-import failed: %v\n%s", err, out)
	}

	for _, tg := range r.tags {
		r.tagAtMark(tg)
	}
}

// tagAtMark resolves the commit hash that the recorded mark refers to
// (every mark sets refs/heads/main, so we walk the first-parent chain
// from main backward by `mark - count` steps) and creates a
// lightweight tag at that commit.
func (r *gitRepo) tagAtMark(tg taggedCommit) {
	r.t.Helper()

	steps := r.mark - tg.mark
	rev := fmt.Sprintf("main~%d", steps)

	if steps == 0 {
		rev = "main"
	}

	hash := r.capture("rev-parse", rev)
	r.run("tag", tg.name, hash)
}

func (r *gitRepo) run(args ...string) {
	r.t.Helper()

	cmd := exec.Command("git", append([]string{"-C", r.dir}, args...)...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func (r *gitRepo) capture(args ...string) string {
	r.t.Helper()

	cmd := exec.Command("git", append([]string{"-C", r.dir}, args...)...)

	out, err := cmd.Output()
	if err != nil {
		r.t.Fatalf("git %v failed: %v", args, err)
	}

	return strings.TrimSpace(string(out))
}
