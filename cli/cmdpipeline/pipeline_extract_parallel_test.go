package cmdpipeline

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParallelPreParseLogLines_SmallSlice(t *testing.T) {
	lines := []string{
		"build\tCompile\t2026-03-23T10:00:00Z warning: unused variable",
		"build\tCompile\t2026-03-23T10:00:01Z   --> main.go:10:2",
		"build\tCompile\t2026-03-23T10:00:02Z ##[error]undefined symbol",
		"build\tCompile\t2026-03-23T10:00:03Z     at compile (main.go:10)",
		"",
	}

	parsed := ParallelPreParseLogLines(lines)
	if len(parsed) != len(lines) {
		t.Fatalf("expected %d parsed lines, got %d", len(lines), len(parsed))
	}

	if !parsed[0].IsWarning {
		t.Errorf("line 0 expected warning")
	}
	if !parsed[1].IsWarningCont {
		t.Errorf("line 1 expected warning continuation")
	}
	if !parsed[2].IsError {
		t.Errorf("line 2 expected error")
	}
	if parsed[4].IsBlank != true {
		t.Errorf("line 4 expected blank")
	}
}

func TestParallelPreParseLogLines_LargeSlice_MatchesSequential(t *testing.T) {
	var lines []string
	for i := 0; i < 250; i++ {
		switch i % 5 {
		case 0:
			lines = append(lines, fmt.Sprintf("job1\tstep1\t2026-03-23T10:00:%02dZ ok\tpassing line %d", i%60, i))
		case 1:
			lines = append(lines, fmt.Sprintf("job1\tstep1\t2026-03-23T10:00:%02dZ warning: deprecation in module %d", i%60, i))
		case 2:
			lines = append(lines, fmt.Sprintf("job1\tstep1\t2026-03-23T10:00:%02dZ   --> src/file%d.go:1:1", i%60, i))
		case 3:
			lines = append(lines, fmt.Sprintf("job1\tstep1\t2026-03-23T10:00:%02dZ ##[error]syntax error on line %d", i%60, i))
		case 4:
			lines = append(lines, "")
		}
	}

	parallelResult := ParallelPreParseLogLines(lines)
	if len(parallelResult) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(parallelResult))
	}

	for i, line := range lines {
		sequentialItem := preParseSingleLogLine(line)
		if !reflect.DeepEqual(parallelResult[i], sequentialItem) {
			t.Fatalf("mismatch at line %d:\nparallel: %+v\nsequential: %+v", i, parallelResult[i], sequentialItem)
		}
	}
}

func TestParseFailedLogLinesParallel_MatchesParseFailedLogLines(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 300; i++ {
		switch i % 6 {
		case 0:
			sb.WriteString(fmt.Sprintf("build\tCompile Job\t2026-03-23T10:00:%02dZ PASS: test_%d\n", i%60, i))
		case 1:
			sb.WriteString(fmt.Sprintf("build\tCompile Job\t2026-03-23T10:00:%02dZ warning: unused variable `x%d`\n", i%60, i))
		case 2:
			sb.WriteString(fmt.Sprintf("build\tCompile Job\t2026-03-23T10:00:%02dZ   --> src/main.rs:%d:5\n", i%60, i))
		case 3:
			sb.WriteString(fmt.Sprintf("build\tCompile Job\t2026-03-23T10:00:%02dZ ##[error]FAIL: test_case_%d failed with exit status 1\n", i%60, i))
		case 4:
			sb.WriteString(fmt.Sprintf("build\tCompile Job\t2026-03-23T10:00:%02dZ     context stack line for failure %d\n", i%60, i))
		case 5:
			sb.WriteString("\n")
		}
	}

	rawLogs := sb.String()
	sequentialJobs := ParseFailedLogLines(rawLogs)
	parallelJobs := ParseFailedLogLinesParallel(rawLogs)

	if len(sequentialJobs) != len(parallelJobs) {
		t.Fatalf("expected %d jobs, got %d", len(sequentialJobs), len(parallelJobs))
	}

	for j := range sequentialJobs {
		sJob := sequentialJobs[j]
		pJob := parallelJobs[j]

		if sJob.JobName != pJob.JobName {
			t.Errorf("job %d name mismatch: seq %s != par %s", j, sJob.JobName, pJob.JobName)
		}
		if sJob.StepName != pJob.StepName {
			t.Errorf("job %d step mismatch: seq %s != par %s", j, sJob.StepName, pJob.StepName)
		}
		if sJob.FailureSummary != pJob.FailureSummary {
			t.Errorf("job %d summary mismatch: seq %s != par %s", j, sJob.FailureSummary, pJob.FailureSummary)
		}
		if len(sJob.ErrorLines) != len(pJob.ErrorLines) {
			t.Errorf("job %d error lines count mismatch: seq %d != par %d", j, len(sJob.ErrorLines), len(pJob.ErrorLines))
		}
		for l := range sJob.ErrorLines {
			if sJob.ErrorLines[l] != pJob.ErrorLines[l] {
				t.Errorf("job %d error line %d mismatch: seq %q != par %q", j, l, sJob.ErrorLines[l], pJob.ErrorLines[l])
			}
		}
	}
}

func TestParseFailedLogLinesParallel_EmptyAndBlank(t *testing.T) {
	if res := ParseFailedLogLinesParallel(""); res != nil {
		t.Errorf("expected nil for empty string, got: %v", res)
	}
	if res := ParseFailedLogLinesParallel("   \n\t  \n"); res != nil {
		t.Errorf("expected nil for blank string, got: %v", res)
	}
}

func TestParseFailedLogLinesParallel_AnsiEscapeStripping(t *testing.T) {
	rawLogs := "\x1b[31mbuild\tCompile\t2026-03-23T10:00:00Z ##[error]\x1b[1mFatal error occurred\x1b[0m\n"
	jobs := ParseFailedLogLinesParallel(rawLogs)
	if len(jobs) == 0 {
		t.Fatalf("expected at least 1 job, got 0")
	}

	if strings.Contains(jobs[0].FailureSummary, "\x1b") {
		t.Errorf("failure summary contains ANSI escape characters: %q", jobs[0].FailureSummary)
	}
}
