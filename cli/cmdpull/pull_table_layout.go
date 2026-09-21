package cmdpull

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const (
	minPullTableTermWidth       = 60
	defaultPullTableTermWidth   = 120
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
	isValid := err == nil && parsed >= minPullTableTermWidth
	if isValid {
		return parsed
	}

	return 0
}

func detectTerminalWidth() int {
	width := pterm.GetTerminalWidth()
	isValid := width >= minPullTableTermWidth
	if isValid {
		return width
	}

	return fallbackTerminalWidth()
}

func fallbackTerminalWidth() int {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	isValidTerm := err == nil && termWidth >= minPullTableTermWidth
	if isValidTerm {
		return termWidth
	}

	envWidth := parseColumnsEnv()
	isValidEnv := envWidth >= minPullTableTermWidth
	if isValidEnv {
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
	maxRelease := resolveWideReleaseWidth(termWidth)
	maxSHA := resolveWideSHAWidth(termWidth)
	maxPR := resolveWidePRWidth(termWidth)
	maxStatus := 14
	fixedWidth := 2 + (6 * colGap) + maxBranch + maxRelease + maxSHA + maxPR + maxStatus
	maxRepo, maxLatest := resolveWideRepoAndLatestWidths(termWidth, fixedWidth)
	totalWidth := fixedWidth + maxRepo + maxLatest

	return maxRepo, maxBranch, maxLatest, maxRelease, maxSHA, totalWidth, colGap
}

func resolveWideRepoAndLatestWidths(termWidth, fixedWidth int) (int, int) {
	avail := termWidth - fixedWidth
	isNarrow := avail < 36
	if isNarrow {
		return 24, 12
	}

	repoWidth := (avail * 58) / 100
	latestWidth := avail - repoWidth

	return repoWidth, latestWidth
}

func resolveWideColGap(termWidth int) int {
	isCompactWide := termWidth < 120
	if isCompactWide {
		return defaultPullTableColGapTight
	}

	return defaultPullTableColGapWide
}

func resolveWideBranchWidth(termWidth int) int {
	isUltraWide := termWidth >= 180
	if isUltraWide {
		return 18
	}

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
		return 7
	}

	return 8
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
	maxRel := 9
	maxSha := 7
	maxPR := 3
	maxStatus := 10
	gap := defaultPullTableColGapTight
	fixedWidth := 2 + (5 * gap) + maxBranch + maxRel + maxSha + maxPR + maxStatus
	maxRepo := resolveCompactRepoWidth(termWidth, fixedWidth)
	totalWidth := fixedWidth + maxRepo

	return maxRepo, maxBranch, maxRel, maxSha, totalWidth
}

func resolveCompactRepoWidth(termWidth, fixedWidth int) int {
	avail := termWidth - fixedWidth
	isNarrow := avail < 16
	if isNarrow {
		return 16
	}

	return avail
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
		PadVisual("COMMIT", l.MaxSHA) + sep +
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
		PadVisual("COMMIT", l.MaxSHA) + sep +
		PadVisual("PR", l.MaxPR) + sep +
		"STATUS"

	fmt.Println(line)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}
