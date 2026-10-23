package cmdpull

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func executePush(records []model.ScanRecord, prog *cloner.BatchProgress, opts pushOptions) {
	workers, ok := cloneconcurrency.Resolve(opts.parallel)
	if !ok {
		cliexit.HandleError(apperror.NewSimple("invalid concurrency", "E9000"), 1)
	}

	opts.parallel = workers
	if opts.parallel > 1 {
		runPushParallel(records, prog, opts.parallel, opts.stopOnFail)

		return
	}

	runPushSequential(records, prog)
}

func runPushSequential(records []model.ScanRecord, prog *cloner.BatchProgress) {
	for _, rec := range records {
		if prog.Stopped() {
			break
		}

		prog.BeginItem(rec.RepoName)
		pushOneRepoTracked(rec, prog)
	}
}

func pushOneRepoTracked(rec model.ScanRecord, prog *cloner.BatchProgress) {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		prog.Skip(rec.RepoName)

		return
	}

	result := cloner.SafePushOne(rec, rec.AbsolutePath)
	if result.IsFailed() {
		prog.FailWithError(rec.RepoName, result.Error)

		return
	}

	if result.Notes == "up-to-date" {
		prog.UpToDate(rec.RepoName)

		return
	}

	prog.Succeed(rec.RepoName)
}

// pushWithAutoRebase runs `git push`, and on non-fast-forward
// rejection auto-runs `git pull --rebase` then retries once.
func pushWithAutoRebase(cwd string, rest []string) {
	gitArgs := append([]string{"push"}, rest...)
	fmt.Printf("→ Running: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	runErr, stderr := runGitCapturingStderr(gitArgs)
	if runErr == nil {
		return
	}

	if !isNonFastForwardRejection(stderr) {
		handleGitExit("git push", runErr)

		return
	}

	retryPushWithRebase(cwd, gitArgs)
}

func retryPushWithRebase(cwd string, gitArgs []string) {
	fmt.Fprintln(os.Stderr, "↻ push rejected (non-fast-forward) — auto-running `git pull --rebase` and retrying")
	if pullErr := runGitInherit([]string{"pull", "--rebase"}); pullErr != nil {
		fmt.Fprintln(os.Stderr, "✗ auto pull --rebase failed — resolve conflicts then re-run `gitmap push`")
		handleGitExit("git pull --rebase", pullErr)

		return
	}

	fmt.Printf("→ Retrying: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	if retryErr := runGitInherit(gitArgs); retryErr != nil {
		handleGitExit("git push (retry)", retryErr)
	}
}

func runGitCapturingStderr(gitArgs []string) (error, string) {
	var buf bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)

	return cmd.Run(), buf.String()
}

func runGitInherit(gitArgs []string) error {
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isNonFastForwardRejection(stderr string) bool {
	lower := strings.ToLower(stderr)
	if !strings.Contains(lower, "[rejected]") && !strings.Contains(lower, "failed to push some refs") {
		return false
	}

	return strings.Contains(lower, "fetch first") || strings.Contains(lower, "non-fast-forward")
}

func handleGitExit(_ string, runErr error) {
	fmt.Fprintf(os.Stderr, "\n  %s[tip]%s To diagnose and auto-repair SSH authentication & push issues, run:\n    %sgitmap push-fix%s  (or: %sgitmap push fix%s)\n\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorGreen, constants.ColorReset,
		constants.ColorGreen, constants.ColorReset,
	)

	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		cliexit.Exit(exitErr.ExitCode())
		return
	}

	cliexit.Exit(1)
}

func joinForLog(args []string) string {
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}

		out += a
	}

	return out
}
