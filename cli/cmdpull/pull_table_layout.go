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
	widePullTableThreshold      = 92
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
	maxRepo, maxBranch, maxLatest, totalWidth := calcWideColumnWidths(termWidth)
	layout := newBaseWideLayout(termWidth, totalWidth, rows)
	layout.MaxRepo = maxRepo
	layout.MaxBranch = maxBranch
	layout.MaxLatestBr = maxLatest

	return layout
}

func newBaseWideLayout(termWidth, totalWidth int, rows []model.PullTableRow) *PullTableLayout {
	return &PullTableLayout{
		TermWidth: termWidth, IsWide: true, ColGap: defaultPullTableColGapWide,
		MaxRange: 0, MaxChanges: 0, MaxSHA: 0, MaxPR: 3, MaxStatus: 14,
		DividerLen: totalWidth - 2, Rows: rows,
	}
}

func calcWideColumnWidths(termWidth int) (int, int, int, int) {
	maxBranch := 8
	maxLatest := 14
	maxRepo := resolveWideRepoWidth(termWidth)
	totalWidth := 2 + maxRepo + 3 + maxBranch + 3 + maxLatest + 3 + 3 + 3 + 14

	return maxRepo, maxBranch, maxLatest, totalWidth
}

func resolveWideRepoWidth(termWidth int) int {
	isExtraWide := termWidth >= 120
	if isExtraWide {
		return 24
	}

	return 22
}

func buildCompactPullTableLayout(rows []model.PullTableRow, termWidth int) *PullTableLayout {
	maxRepo, maxBranch, totalWidth := calcCompactColumnWidths(termWidth)
	layout := newBaseCompactLayout(termWidth, totalWidth, rows)
	layout.MaxRepo = maxRepo
	layout.MaxBranch = maxBranch

	return layout
}

func newBaseCompactLayout(termWidth, totalWidth int, rows []model.PullTableRow) *PullTableLayout {
	return &PullTableLayout{
		TermWidth: termWidth, IsWide: false, ColGap: defaultPullTableColGapTight,
		MaxRange: 0, MaxChanges: 0, MaxSHA: 0, MaxPR: 3, MaxStatus: 12,
		DividerLen: totalWidth - 2, Rows: rows,
	}
}

func calcCompactColumnWidths(termWidth int) (int, int, int) {
	maxBranch := 10
	maxRepo := resolveCompactRepoWidth(termWidth)
	totalWidth := 2 + maxRepo + 2 + maxBranch + 2 + 3 + 2 + 12

	return maxRepo, maxBranch, totalWidth
}

func resolveCompactRepoWidth(termWidth int) int {
	isNarrow := termWidth < 70
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
	sep := "   "
	line := "  " +
		PadVisual("REPO", l.MaxRepo) + sep +
		PadVisual("BRANCH", l.MaxBranch) + sep +
		PadVisual("LATEST BRANCH", l.MaxLatestBr) + sep +
		PadVisual("PR", l.MaxPR) + sep +
		"STATUS"

	fmt.Println(line)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}

func (l *PullTableLayout) printCompactHeader() {
	sep := "  "
	line := "  " +
		PadVisual("REPO", l.MaxRepo) + sep +
		PadVisual("BRANCH", l.MaxBranch) + sep +
		PadVisual("PR", l.MaxPR) + sep +
		"STATUS"

	fmt.Println(line)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}
