package cmdautomation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type checkFn func(opts PreflightOptions) PreflightCheck

// RunPreflight executes parallel multi-core checks across CPU cores.
func RunPreflight(opts PreflightOptions) PreflightResultMonad {
	start := time.Now()
	checks := resolvePreflightChecks(opts)
	outcomes := executePreflightPool(checks, opts)
	return result.Ok(aggregatePreflightResult(outcomes, start))
}

func resolvePreflightChecks(opts PreflightOptions) []checkFn {
	all := []checkFn{
		runFormatCheck,
		runRelPathsCheck,
		runNestedIfCheck,
		runNamingCheck,
		runGoModulesCheck,
	}
	return filterPreflightChecks(all, opts.Filter)
}

func filterPreflightChecks(checks []checkFn, filter string) []checkFn {
	if len(filter) == 0 {
		return checks
	}
	filt := strings.ToLower(filter)
	var filtered []checkFn
	for _, fn := range checks {
		probe := fn(PreflightOptions{Dir: "."})
		if strings.Contains(strings.ToLower(probe.Name), filt) {
			filtered = append(filtered, fn)
		}
	}
	return filtered
}

func executePreflightPool(checks []checkFn, opts PreflightOptions) []PreflightCheck {
	workers := resolveWorkerCount(opts.Workers, len(checks))
	jobs := make(chan checkFn, len(checks))
	results := make(chan PreflightCheck, len(checks))
	spawnPreflightWorkers(jobs, results, opts, workers)
	for _, fn := range checks {
		jobs <- fn
	}
	close(jobs)
	return collectPreflightResults(results, len(checks), opts.IsFailFast)
}

func resolveWorkerCount(requested, maxJobs int) int {
	if requested > 0 {
		return requested
	}
	cores := runtime.NumCPU()
	if cores > maxJobs {
		return maxJobs
	}
	if cores < 1 {
		return 1
	}
	return cores
}

func spawnPreflightWorkers(jobs <-chan checkFn, results chan<- PreflightCheck, opts PreflightOptions, count int) {
	for i := 0; i < count; i++ {
		go func() {
			for fn := range jobs {
				results <- fn(opts)
			}
		}()
	}
}

func collectPreflightResults(results <-chan PreflightCheck, total int, isFailFast bool) []PreflightCheck {
	var collected []PreflightCheck
	for i := 0; i < total; i++ {
		res := <-results
		collected = append(collected, res)
		if isFailFast && res.IsFail {
			break
		}
	}
	return collected
}

func runFormatCheck(opts PreflightOptions) PreflightCheck {
	start := time.Now()
	dir := resolveAuditDir(opts.Dir)
	cmd := exec.Command("gofmt", "-l", dir)
	out, err := cmd.Output()
	isFail := err != nil || len(strings.TrimSpace(string(out))) > 0
	msg := "All Go files are properly formatted"
	if isFail {
		msg = "Unformatted Go files detected"
	}
	return PreflightCheck{
		Name:     "Go Format Check",
		Category: "format",
		IsPass:   !isFail,
		IsFail:   isFail,
		Duration: time.Since(start),
		Message:  msg,
	}
}

func runRelPathsCheck(opts PreflightOptions) PreflightCheck {
	start := time.Now()
	monad := RunRelPathsAudit(RelPathOptions{Dir: opts.Dir})
	isFail := monad.IsFailure() || len(monad.Value.Violations) > 0
	msg := "Zero forbidden absolute paths found"
	if isFail {
		msg = fmt.Sprintf("Found %d forbidden absolute path(s)", len(monad.Value.Violations))
	}
	return PreflightCheck{
		Name:     "Relative Path Check",
		Category: "relpaths",
		IsPass:   !isFail,
		IsFail:   isFail,
		Duration: time.Since(start),
		Message:  msg,
	}
}

func runNestedIfCheck(opts PreflightOptions) PreflightCheck {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	vios := scanNestedIfs(root)
	isFail := len(vios) > 0
	msg := "Zero nested if statements found"
	if isFail {
		msg = fmt.Sprintf("Found %d nested if violation(s)", len(vios))
	}
	return PreflightCheck{
		Name:     "Nested If Linter",
		Category: "nested_if",
		IsPass:   !isFail,
		IsFail:   isFail,
		Duration: time.Since(start),
		Message:  msg,
		Details:  vios,
	}
}

func scanNestedIfs(root string) []string {
	files, _ := collectTextFiles(root, []string{".go"})
	var vios []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		vios = append(vios, checkFileNestedIfs(f, string(data))...)
	}
	return vios
}

func checkFileNestedIfs(path, content string) []string {
	var vios []string
	depth := 0
	for idx, line := range strings.Split(content, "\n") {
		depth = depth + strings.Count(line, "{") - strings.Count(line, "}")
		if depth < 0 {
			depth = 0
		}
		if strings.HasPrefix(strings.TrimSpace(line), "if ") && depth > 2 {
			vios = append(vios, fmt.Sprintf("%s:%d: nested if", path, idx+1))
		}
	}
	return vios
}

func runNamingCheck(opts PreflightOptions) PreflightCheck {
	start := time.Now()
	monad := RunNamingAudit(NamingOptions{Dir: opts.Dir})
	isFail := monad.IsFailure() || len(monad.Value.Violations) > 0
	msg := "All files adhere to affirmative boolean and naming conventions"
	if isFail {
		msg = fmt.Sprintf("Found %d boolean/naming violation(s)", len(monad.Value.Violations))
	}
	return PreflightCheck{
		Name:     "Boolean & Naming Check",
		Category: "naming",
		IsPass:   !isFail,
		IsFail:   isFail,
		Duration: time.Since(start),
		Message:  msg,
	}
}

func runGoModulesCheck(opts PreflightOptions) PreflightCheck {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	mods := findGoModDirs(root)
	msg := fmt.Sprintf("Found %d Go module(s) ready for preflight", len(mods))
	return PreflightCheck{
		Name:     "Go Modules Discovery",
		Category: "gotest",
		IsPass:   true,
		IsFail:   false,
		Duration: time.Since(start),
		Message:  msg,
		Details:  mods,
	}
}

func findGoModDirs(root string) []string {
	var dirs []string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info == nil || !info.IsDir() {
			return nil
		}
		if info.Name() == ".git" || info.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if _, statErr := os.Stat(filepath.Join(p, "go.mod")); statErr == nil {
			dirs = append(dirs, p)
		}
		return nil
	})
	return dirs
}

func aggregatePreflightResult(checks []PreflightCheck, start time.Time) PreflightResult {
	failed := 0
	for _, c := range checks {
		if c.IsFail {
			failed++
		}
	}
	return PreflightResult{
		TotalChecks:  len(checks),
		PassedChecks: len(checks) - failed,
		FailedChecks: failed,
		Checks:       checks,
		Duration:     time.Since(start),
		IsPass:       failed == 0,
	}
}
