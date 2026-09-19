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
		MaxRange: 12, MaxChanges: 7, MaxSHA: 7, MaxPR: 3, MaxStatus: 10,
		DividerLen: totalWidth - 2, Rows: rows,
	}
}

func calcWideColumnWidths(termWidth int) (int, int, int, int) {
	avail := calcWideAvailableWidth(termWidth)
	maxBranch := avail * 20 / 100
	if maxBranch > 10 {
		maxBranch = 10
	}
	if maxBranch < 6 {
		maxBranch = 6
	}
	maxLatest := avail * 28 / 100
	if maxLatest > 14 {
		maxLatest = 14
	}
	if maxLatest < 9 {
		maxLatest = 9
	}
	maxRepo := avail - maxBranch - maxLatest
	if maxRepo < 14 {
		maxRepo = 14
	}
	totalWidth := 2 + maxRepo + 3 + maxBranch + 3 + maxLatest + 3 + 12 + 3 + 7 + 3 + 3 + 3 + 10

	return maxRepo, maxBranch, maxLatest, totalWidth
}

func calcWideAvailableWidth(termWidth int) int {
	fixedWidth := 2 + 12 + 7 + 3 + 10 + (6 * defaultPullTableColGapWide)
	avail := termWidth - fixedWidth
	if avail < 30 {
		return 30
	}

	return avail
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
		MaxRange: 8, MaxChanges: 7, MaxSHA: 7, MaxPR: 0, MaxStatus: 8,
		DividerLen: totalWidth - 2, Rows: rows,
	}
}

func calcCompactColumnWidths(termWidth int) (int, int, int) {
	avail := calcCompactAvailableWidth(termWidth)
	maxRepo := avail * 55 / 100
	if maxRepo < 10 {
		maxRepo = 10
	}
	maxBranch := avail - maxRepo
	if maxBranch < 8 {
		maxBranch = 8
	}
	totalWidth := 2 + maxRepo + 2 + maxBranch + 2 + 8 + 2 + 7 + 2 + 8

	return maxRepo, maxBranch, totalWidth
}

func calcCompactAvailableWidth(termWidth int) int {
	fixedWidth := 2 + 8 + 7 + 8 + (4 * defaultPullTableColGapTight)
	avail := termWidth - fixedWidth
	if avail < 16 {
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
	sep := "   "
	line := "  " +
		PadVisual("REPO", l.MaxRepo) + sep +
		PadVisual("BRANCH", l.MaxBranch) + sep +
		PadVisual("LATEST BRANCH", l.MaxLatestBr) + sep +
		PadVisual("COMMIT RANGE", l.MaxRange) + sep +
		PadVisual("CHANGES", l.MaxChanges) + sep +
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
		PadVisual("RANGE", l.MaxRange) + sep +
		PadVisual("CHANGES", l.MaxChanges) + sep +
		"STATUS"

	fmt.Println(line)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}
