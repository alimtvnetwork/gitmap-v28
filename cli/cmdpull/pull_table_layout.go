package cmdpull

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const (
	minPullTableTermWidth       = 60
	defaultPullTableTermWidth   = 80
	widePullTableThreshold      = 100
	defaultPullTableColGapWide  = 3
	defaultPullTableColGapTight = 2
)

type PullTableLayout struct {
	TermWidth   int
	IsWide      bool
	ColGap      int
	MaxRepo     int
	MaxBranch   int
	MaxLatestBr int
	MaxRelease  int
	MaxRange    int
	MaxChanges  int
	MaxSHA      int
	MaxPR       int
	MaxStatus   int
	DividerLen  int
	Rows        []model.PullTableRow
}

func parseColumnsEnv() int {
	colsEnv := strings.TrimSpace(os.Getenv("COLUMNS"))
	parsed, err := strconv.Atoi(colsEnv)
	if err == nil && parsed >= minPullTableTermWidth {
		return parsed
	}

	return 0
}

func detectTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && width >= minPullTableTermWidth {
		return width
	}

	envWidth := parseColumnsEnv()
	if envWidth >= minPullTableTermWidth {
		return envWidth
	}

	return defaultPullTableTermWidth
}

func NewPullTableLayout(rows []model.PullTableRow) *PullTableLayout {
	termWidth := detectTerminalWidth()

	return NewPullTableLayoutWithWidth(rows, termWidth)
}

func NewPullTableLayoutWithWidth(rows []model.PullTableRow, termWidth int) *PullTableLayout {
	if termWidth < minPullTableTermWidth {
		termWidth = defaultPullTableTermWidth
	}

	isWide := termWidth >= widePullTableThreshold
	if isWide {
		return buildWidePullTableLayout(rows, termWidth)
	}

	return buildCompactPullTableLayout(rows, termWidth)
}

func buildWidePullTableLayout(rows []model.PullTableRow, termWidth int) *PullTableLayout {
	maxRepo, maxBranch, maxLatest, maxRel, maxSha, totalWidth, colGap := calcWideColumnWidths(termWidth)
	layout := newBaseWideLayout(termWidth, totalWidth, colGap, rows)
	layout.MaxRepo = maxRepo
	layout.MaxBranch = maxBranch
	layout.MaxLatestBr = maxLatest
	layout.MaxRelease = maxRel
	layout.MaxSHA = maxSha

	return layout
}

func newBaseWideLayout(termWidth, totalWidth, colGap int, rows []model.PullTableRow) *PullTableLayout {
	return &PullTableLayout{
		TermWidth: termWidth, IsWide: true, ColGap: colGap,
		MaxRange: 0, MaxChanges: 0, MaxPR: 4, MaxStatus: 14,
		DividerLen: totalWidth - 2, Rows: rows,
	}
}

func calcWideColumnWidths(termWidth int) (int, int, int, int, int, int, int) {
	colGap := resolveWideColGap(termWidth)
	maxBranch := resolveWideBranchWidth(termWidth)
	maxLatest := resolveWideLatestWidth(termWidth)
	maxRepo := resolveWideRepoWidth(termWidth)
	maxRelease := resolveWideReleaseWidth(termWidth)
	maxSHA := resolveWideSHAWidth(termWidth)
	maxPR := resolveWidePRWidth(termWidth)
	maxStatus := 14
	totalWidth := 2 + maxRepo + colGap + maxBranch + colGap + maxLatest + colGap + maxRelease + colGap + maxSHA + colGap + maxPR + colGap + maxStatus

	return maxRepo, maxBranch, maxLatest, maxRelease, maxSHA, totalWidth, colGap
}

func resolveWideColGap(termWidth int) int {
	isCompactWide := termWidth < 120
	if isCompactWide {
		return defaultPullTableColGapTight
	}

	return defaultPullTableColGapWide
}

func resolveWideRepoWidth(termWidth int) int {
	isExtraWide := termWidth >= 140
	if isExtraWide {
		return 34
	}

	isWide := termWidth >= 120
	if isWide {
		return 30
	}

	return 24
}

func resolveWideBranchWidth(termWidth int) int {
	isExtraWide := termWidth >= 140
	if isExtraWide {
		return 16
	}

	isWide := termWidth >= 120
	if isWide {
		return 13
	}

	return 11
}

func resolveWideLatestWidth(termWidth int) int {
	isExtraWide := termWidth >= 140
	if isExtraWide {
		return 18
	}

	isWide := termWidth >= 120
	if isWide {
		return 15
	}

	return 13
}

func resolveWideReleaseWidth(termWidth int) int {
	isSmallWide := termWidth < 120
	if isSmallWide {
		return 9
	}

	return 10
}

func resolveWideSHAWidth(termWidth int) int {
	isSmallWide := termWidth < 120
	if isSmallWide {
		return 5
	}

	return 6
}

func resolveWidePRWidth(termWidth int) int {
	isSmallWide := termWidth < 120
	if isSmallWide {
		return 3
	}

	return 4
}

func buildCompactPullTableLayout(rows []model.PullTableRow, termWidth int) *PullTableLayout {
	maxRepo, maxBranch, maxRel, maxSha, totalWidth := calcCompactColumnWidths(termWidth)
	layout := newBaseCompactLayout(termWidth, totalWidth, rows)
	layout.MaxRepo = maxRepo
	layout.MaxBranch = maxBranch
	layout.MaxRelease = maxRel
	layout.MaxSHA = maxSha

	return layout
}

func newBaseCompactLayout(termWidth, totalWidth int, rows []model.PullTableRow) *PullTableLayout {
	return &PullTableLayout{
		TermWidth: termWidth, IsWide: false, ColGap: defaultPullTableColGapTight,
		MaxRange: 0, MaxChanges: 0, MaxPR: 3, MaxStatus: 10,
		DividerLen: totalWidth - 2, Rows: rows,
	}
}

func calcCompactColumnWidths(termWidth int) (int, int, int, int, int) {
	maxBranch := 10
	maxRepo := resolveCompactRepoWidth(termWidth)
	maxRel := 9
	maxSha := 5
	maxPR := 3
	maxStatus := 10
	gap := defaultPullTableColGapTight
	totalWidth := 2 + maxRepo + gap + maxBranch + gap + maxRel + gap + maxSha + gap + maxPR + gap + maxStatus

	return maxRepo, maxBranch, maxRel, maxSha, totalWidth
}

func resolveCompactRepoWidth(termWidth int) int {
	isNarrow := termWidth < 75
	if isNarrow {
		return 16
	}

	return 18
}

func (l *PullTableLayout) PrintHeader() {
	if l.IsWide {
		l.printWideHeader()

		return
	}

	l.printCompactHeader()
}

func (l *PullTableLayout) printWideHeader() {
	sep := strings.Repeat(" ", l.ColGap)
	line := "  " +
		PadVisual("REPO", l.MaxRepo) + sep +
		PadVisual("BRANCH", l.MaxBranch) + sep +
		PadVisual("LATEST BRANCH", l.MaxLatestBr) + sep +
		PadVisual("RELEASE", l.MaxRelease) + sep +
		PadVisual("SHA", l.MaxSHA) + sep +
		PadVisual("PR", l.MaxPR) + sep +
		"STATUS"

	fmt.Println(line)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}

func (l *PullTableLayout) printCompactHeader() {
	sep := strings.Repeat(" ", l.ColGap)
	line := "  " +
		PadVisual("REPO", l.MaxRepo) + sep +
		PadVisual("BRANCH", l.MaxBranch) + sep +
		PadVisual("RELEASE", l.MaxRelease) + sep +
		PadVisual("SHA", l.MaxSHA) + sep +
		PadVisual("PR", l.MaxPR) + sep +
		"STATUS"

	fmt.Println(line)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}
