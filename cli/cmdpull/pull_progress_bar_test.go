package cmdpull

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestFormatVisualBarSafe(t *testing.T) {
	bar := FormatVisualBar(5, 10, 10, true)
	isExpected := bar == "[=====-----]"
	if !isExpected {
		t.Fatalf("expected [=====-----], got %q", bar)
	}

	fullBar := FormatVisualBar(10, 10, 10, true)
	isFullExpected := fullBar == "[==========]"
	if !isFullExpected {
		t.Fatalf("expected [==========], got %q", fullBar)
	}
}

func TestFormatVisualBarRich(t *testing.T) {
	bar := FormatVisualBar(5, 10, 10, false)
	hasBlocks := strings.Contains(bar, "█")
	if !hasBlocks {
		t.Fatalf("expected filled blocks, got %q", bar)
	}

	hasShaded := strings.Contains(bar, "░")
	if !hasShaded {
		t.Fatalf("expected empty shaded blocks, got %q", bar)
	}
}

func TestPullProgressBarLifecycle(t *testing.T) {
	var buf bytes.Buffer
	bar := NewPullProgressBar(2, false, false)
	bar.SetOutput(&buf)
	bar.SetTTY(false)
	bar.Start()
	bar.UpdateRepoStep("repo-a", PullStepTypeFetching, "pulling")
	state := &PullRepoState{RepoName: "repo-a", Step: PullStepTypeUpToDate, Changes: "up-to-date"}
	bar.CompleteRepo(state)
	bar.Stop()
	isCountValid := bar.Succeeded() == 1
	if !isCountValid {
		t.Fatalf("expected succeeded 1, got %d", bar.Succeeded())
	}
}

func TestPullProgressBarStopOnFail(t *testing.T) {
	var buf bytes.Buffer
	bar := NewPullProgressBar(5, false, true)
	bar.SetOutput(&buf)
	bar.SetTTY(false)
	bar.CompleteRepo(&PullRepoState{
		RepoName: "bad-repo",
		Step:     PullStepTypeError,
		Changes:  "fatal: error",
	})
	isStopped := bar.IsStopped()
	if !isStopped {
		t.Fatalf("expected bar to be stopped on error with stopOnFail")
	}
	assertFailedCount(t, bar, 1)
}

func assertFailedCount(t *testing.T, bar *PullProgressBar, expected int) {
	isFailedExpected := bar.Failed() == expected
	if !isFailedExpected {
		t.Fatalf("expected failed %d, got %d", expected, bar.Failed())
	}
}

func TestPullProgressBarTickerStartStop(t *testing.T) {
	var buf bytes.Buffer
	bar := NewPullProgressBar(1, false, false)
	bar.SetOutput(&buf)
	bar.SetTTY(true)
	bar.Start()
	time.Sleep(100 * time.Millisecond)
	bar.Stop()
	bar.Stop()
	isStopped := bar.IsStopped()
	if !isStopped {
		t.Fatalf("expected bar to be marked stopped")
	}
}

func TestPullProgressBarWorkerSlotRegistry(t *testing.T) {
	bar := NewPullProgressBar(4, false, false)
	bar.RegisterWorker(0, "repo-alpha")
	bar.RegisterWorker(1, "repo-beta")
	assertActiveWorkerCount(t, bar, 2)

	bar.UpdateWorkerProgress(0, PullStepTypeFetching, "Receiving (50%)", 50)
	bar.UnregisterWorker(0)
	assertActiveWorkerCount(t, bar, 1)
}

func assertActiveWorkerCount(t *testing.T, bar *PullProgressBar, expected int) {
	count := len(bar.ActiveWorkers())
	isExpected := count == expected
	if !isExpected {
		t.Fatalf("expected %d active workers, got %d", expected, count)
	}
}

func TestPullProgressBarSingleRepoSubSteps(t *testing.T) {
	bar := NewPullProgressBar(1, false, false)
	bar.SetSubStep(2, 4, "fetching remote")
	idx, total := bar.SubStepMilestone()
	isSubStepValid := idx == 2 && total == 4
	if !isSubStepValid {
		t.Fatalf("expected step 2/4, got %d/%d", idx, total)
	}
	title := StepMilestoneTitle(2, PullStepTypeFetching)
	isTitleValid := title == "Fetching remote objects"
	if !isTitleValid {
		t.Fatalf("expected title 'Fetching remote objects', got %q", title)
	}
}

func TestParseGitProgressLineObjects(t *testing.T) {
	line := "Receiving objects:  45% (9/20), 1.20 MiB"
	event, hasEvent := ParseGitProgressLine(line)
	if !hasEvent {
		t.Fatalf("expected event for receiving objects")
	}
	isPhaseReceiving := event.Phase == "Receiving"
	if !isPhaseReceiving {
		t.Fatalf("expected phase Receiving, got %q", event.Phase)
	}
	isPercentValid := event.Percent == 45
	if !isPercentValid {
		t.Fatalf("expected percent 45, got %d", event.Percent)
	}
}

func TestParseGitProgressLineFastForward(t *testing.T) {
	line := "Updating 03be798..f1d94df"
	event, hasEvent := ParseGitProgressLine(line)
	if !hasEvent {
		t.Fatalf("expected event for updating line")
	}
	isRangeValid := event.CommitRange == "03be798..f1d94df"
	if !isRangeValid {
		t.Fatalf("expected commit range 03be798..f1d94df, got %q", event.CommitRange)
	}
	isPhaseFF := event.Phase == "Fast-forward"
	if !isPhaseFF {
		t.Fatalf("expected phase Fast-forward, got %q", event.Phase)
	}
}

func TestSortStatesAlphabetically(t *testing.T) {
	states := []*PullRepoState{
		{RepoName: "zebra"},
		{RepoName: "Alpha"},
		{RepoName: "beta"},
	}
	sorted := sortStatesAlphabetically(states)
	isFirstAlpha := sorted[0].RepoName == "Alpha"
	isSecondBeta := sorted[1].RepoName == "beta"
	isThirdZebra := sorted[2].RepoName == "zebra"
	if !isFirstAlpha || !isSecondBeta || !isThirdZebra {
		t.Fatalf("expected [Alpha, beta, zebra], got [%s, %s, %s]",
			sorted[0].RepoName, sorted[1].RepoName, sorted[2].RepoName)
	}
}

func TestPullProgressBarTTYRender(t *testing.T) {
	var buf bytes.Buffer
	bar := NewPullProgressBar(1, false, false)
	bar.SetOutput(&buf)
	bar.SetTTY(true)
	bar.UpdateRepoStep("repo-x", PullStepTypeFetching, "Receiving")
	out := buf.String()
	hasEscape := strings.Contains(out, "\r\033[K")
	if !hasEscape {
		t.Fatalf("expected escape sequence in TTY output, got %q", out)
	}
}

func TestClampLineToTermWidth(t *testing.T) {
	short := "hello world"
	clamped := clampLineToTermWidth(short)
	if clamped != short {
		t.Fatalf("expected %q, got %q", short, clamped)
	}

	longLine := strings.Repeat("A", 200)
	clampedLong := clampLineToTermWidth(longLine)
	if len(clampedLong) >= 200 {
		t.Fatalf("expected clamped line < 200, got len %d", len(clampedLong))
	}
	if !strings.HasSuffix(clampedLong, "...") {
		t.Fatalf("expected ellipsis suffix, got %q", clampedLong)
	}
}

func TestFormatActiveWorkersListBounded(t *testing.T) {
	bar := NewPullProgressBar(5, false, false)
	bar.RegisterWorker(0, "repo-1")
	bar.RegisterWorker(1, "repo-2")
	bar.RegisterWorker(2, "repo-3")

	summary := bar.formatActiveWorkersList()
	if !strings.Contains(summary, "(+1 more)") {
		t.Fatalf("expected summary to contain '(+1 more)', got %q", summary)
	}
}

func TestColorizeBarAndBadge(t *testing.T) {
	bar := NewPullProgressBar(1, false, false)
	cyanBar := bar.colorizeBar("[====]", 50)
	isCyan := strings.Contains(cyanBar, constants.ColorCyan)
	if !isCyan {
		t.Fatalf("expected cyan bar, got %q", cyanBar)
	}
	greenBar := bar.colorizeBar("[====]", 100)
	isGreen := strings.Contains(greenBar, constants.ColorGreen)
	if !isGreen {
		t.Fatalf("expected green bar, got %q", greenBar)
	}
	badge := bar.colorizeBadge("✔", PullStepTypeUpToDate)
	isBadgeGreen := strings.Contains(badge, constants.ColorGreen)
	if !isBadgeGreen {
		t.Fatalf("expected green badge, got %q", badge)
	}
}

func TestPrependAll(t *testing.T) {
	res1 := prependAll([]string{"--verbose"})
	isFirstAll := len(res1) == 2 && res1[0] == "--all"
	if !isFirstAll {
		t.Fatalf("expected --all prepended, got %v", res1)
	}
	res2 := prependAll([]string{"--all", "--verbose"})
	isUnchanged := len(res2) == 2 && res2[0] == "--all"
	if !isUnchanged {
		t.Fatalf("expected idempotent prepend, got %v", res2)
	}
}

func TestNormalizePullArgsPa(t *testing.T) {
	res := NormalizePullArgs([]string{"pa", "--verbose"})
	hasAll := len(res) == 2 && res[0] == "--all"
	if !hasAll {
		t.Fatalf("expected 'pa' normalized to '--all', got %v", res)
	}
	res2 := NormalizePullArgs([]string{"pull-all"})
	hasPullAll := len(res2) == 1 && res2[0] == "--all"
	if !hasPullAll {
		t.Fatalf("expected 'pull-all' normalized to '--all', got %v", res2)
	}
}
