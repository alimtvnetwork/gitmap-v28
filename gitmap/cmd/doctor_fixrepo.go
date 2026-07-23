// Package cmd — `gitmap doctor fix-repo` subcommand.
//
// Runs targeted probes for the fix-repo → gofmt pipeline: gofmt on
// PATH, gofmt actually executes, argv budget sanity, and a self-test
// of the chunker. Introduced in v6.80.1 after the Windows argv
// overflow reported at
// .lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md to give users
// a one-shot health check before or after tuning
// --gofmt-max-cmd-len.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDoctorFixRepo executes the fix-repo probe suite. Called by
// runDoctor when args[0] == "fix-repo" (or "fr"). Honors --json and
// --budget <N>.
func runDoctorFixRepo(args []string) {
	wantJSON := false
	budget := constants.FixRepoGofmtMaxCmdLen
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" {
			wantJSON = true

			continue
		}
		if a == "--budget" && i+1 < len(args) {
			if n, err := strconv.Atoi(args[i+1]); err == nil && n >= constants.FixRepoGofmtMinCmdLen {
				budget = n
			}
			i++
		} else if strings.HasPrefix(a, "--budget=") {
			if n, err := strconv.Atoi(a[len("--budget="):]); err == nil && n >= constants.FixRepoGofmtMinCmdLen {
				budget = n
			}
		}
	}
	results := doctorFixRepoProbes(budget)
	failed := 0
	for _, r := range results {
		if !r.OK {
			failed++
		}
	}
	if wantJSON {
		emitDoctorFixRepoJSON(results, failed, budget)
	} else {
		emitDoctorFixRepoText(results, failed, budget)
	}
	if failed > 0 {
		os.Exit(1)
	}
}

// doctorFixRepoProbes runs each probe in order and returns results.
func doctorFixRepoProbes(budget int) []DoctorResult {
	return []DoctorResult{
		probeGofmtPresent(),
		probeGofmtRuns(),
		probeArgvBudget(budget),
		probeChunkerSelfTest(budget),
	}
}

// probeGofmtPresent checks gofmt is on PATH.
func probeGofmtPresent() DoctorResult {
	path, err := exec.LookPath("gofmt")
	if err != nil {
		return DoctorResult{
			Name:    "gofmt-present",
			OK:      false,
			Detail:  "gofmt not on PATH",
			FixHint: "Install the Go toolchain from https://go.dev/dl/ and reopen the shell",
		}
	}

	return DoctorResult{Name: "gofmt-present", OK: true, Detail: path}
}

// probeGofmtRuns writes a scratch .go file to a tempdir and shells
// out to `gofmt -l` against it. Catches PATH-injected shims and
// unexecutable binaries — cases where LookPath succeeds but exec
// fails with the same "filename or extension is too long" class of
// error we hit in production.
func probeGofmtRuns() DoctorResult {
	dir, err := os.MkdirTemp("", "gitmap-doctor-gofmt-")
	if err != nil {
		return DoctorResult{Name: "gofmt-runs", OK: false, Detail: "mktemp: " + err.Error()}
	}
	defer os.RemoveAll(dir)
	sample := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(sample, []byte("package sample\n"), 0o644); err != nil {
		return DoctorResult{Name: "gofmt-runs", OK: false, Detail: "write sample: " + err.Error()}
	}
	cmd := exec.Command("gofmt", "-l", sample)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return DoctorResult{
			Name:    "gofmt-runs",
			OK:      false,
			Detail:  fmt.Sprintf("exec failed: %v (%s)", err, strings.TrimSpace(string(out))),
			FixHint: "Reinstall Go; on Windows also check that no AV shim proxies gofmt.exe",
		}
	}

	return DoctorResult{Name: "gofmt-runs", OK: true, Detail: "gofmt -l executed cleanly"}
}

// probeArgvBudget reports the configured budget. On Windows it also
// measures the effective argv cap by attempting exec.Command with a
// doubling filler until failure, so users can see if their real cap
// is below the documented 32,767. On non-Windows the ARG_MAX is
// ~2 MB and never the bottleneck, so measurement is skipped.
func probeArgvBudget(budget int) DoctorResult {
	if runtime.GOOS != "windows" {
		return DoctorResult{
			Name:   "argv-budget",
			OK:     true,
			Detail: fmt.Sprintf("configured=%d (measurement skipped on %s; ARG_MAX ~2MB)", budget, runtime.GOOS),
		}
	}
	measured := measureWindowsArgvCap()
	ok := measured == 0 || measured >= budget
	detail := fmt.Sprintf("configured=%d, measured=%d", budget, measured)
	hint := ""
	if !ok {
		hint = fmt.Sprintf("Re-run with --gofmt-max-cmd-len %d (or lower)", measured/2)
	}

	return DoctorResult{Name: "argv-budget", OK: ok, Detail: detail, FixHint: hint}
}

// measureWindowsArgvCap probes the real CreateProcess limit by
// running `gofmt -l` with progressively larger filler argv until it
// fails. Returns the largest joined length that succeeded, or 0 when
// gofmt itself is missing.
func measureWindowsArgvCap() int {
	if _, err := exec.LookPath("gofmt"); err != nil {
		return 0
	}
	// Try 8k, 16k, 24k, 32k joined arg bytes. Anything above 32k
	// exceeds the documented cap.
	sizes := []int{8000, 16000, 24000, 32000}
	dir, err := os.MkdirTemp("", "gitmap-doctor-argv-")
	if err != nil {
		return 0
	}
	defer os.RemoveAll(dir)
	sample := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(sample, []byte("package sample\n"), 0o644); err != nil {
		return 0
	}
	last := 0
	for _, target := range sizes {
		copies := target / (len(sample) + 1)
		if copies < 1 {
			copies = 1
		}
		args := make([]string, 0, copies+1)
		args = append(args, "-l")
		for i := 0; i < copies; i++ {
			args = append(args, sample)
		}
		if err := exec.Command("gofmt", args...).Run(); err != nil {
			return last
		}
		last = target
	}

	return last
}

// probeChunkerSelfTest exercises chunkPathsForGofmt against three
// synthetic inputs and asserts invariants (batch count, per-batch
// budget respected).
func probeChunkerSelfTest(budget int) DoctorResult {
	if got := chunkPathsForGofmt(nil, budget); got != nil {
		return DoctorResult{Name: "chunker-selftest", OK: false, Detail: "empty input did not return nil"}
	}
	small := []string{"a.go", "b.go", "c.go"}
	if got := chunkPathsForGofmt(small, budget); len(got) != 1 {
		return DoctorResult{
			Name: "chunker-selftest", OK: false,
			Detail: fmt.Sprintf("expected 1 chunk for %d small paths, got %d", len(small), len(got)),
		}
	}
	long := strings.Repeat("x", 200)
	overflow := make([]string, 500)
	for i := range overflow {
		overflow[i] = long
	}
	batches := chunkPathsForGofmt(overflow, budget)
	if len(batches) < 2 {
		return DoctorResult{
			Name: "chunker-selftest", OK: false,
			Detail: fmt.Sprintf("overflow input yielded only %d batch(es)", len(batches)),
		}
	}
	for i, b := range batches {
		if len(b) > 1 && batchCmdLen(b)-gofmtArgvOverhead > budget {
			return DoctorResult{
				Name: "chunker-selftest", OK: false,
				Detail: fmt.Sprintf("batch %d exceeds budget %d", i+1, budget),
			}
		}
	}

	return DoctorResult{
		Name:   "chunker-selftest",
		OK:     true,
		Detail: fmt.Sprintf("empty ok; small ok (1 chunk); overflow ok (%d chunks)", len(batches)),
	}
}

func emitDoctorFixRepoText(results []DoctorResult, failed, budget int) {
	fmt.Printf("gitmap doctor fix-repo  budget=%d\n\n", budget)
	for _, r := range results {
		mark := "[ok]  "
		if !r.OK {
			mark = "[fail]"
		}
		fmt.Printf("%s %-18s %s\n", mark, r.Name, r.Detail)
		if !r.OK && r.FixHint != "" {
			fmt.Printf("                    fix: %s\n", r.FixHint)
		}
	}
	if failed > 0 {
		fmt.Printf("\n%d probe(s) failed.\n", failed)

		return
	}
	fmt.Println("\nfix-repo → gofmt pipeline nominal.")
}

func emitDoctorFixRepoJSON(results []DoctorResult, failed, budget int) {
	payload := struct {
		Command string         `json:"command"`
		Budget  int            `json:"budget"`
		Failed  int            `json:"failed"`
		Results []DoctorResult `json:"results"`
	}{
		Command: "doctor fix-repo",
		Budget:  budget,
		Failed:  failed,
		Results: results,
	}
	buf, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(buf))
}
