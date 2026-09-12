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
	widePullTableThreshold      = 105
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
	fixedWidth := 2 + 10 + 10 + 7 + 4 + (6 * defaultPullTableColGapWide)
	avail := termWidth - 2 - fixedWidth
	if avail < 36 {
		avail = 36
	}

	maxRepo := avail * 40 / 100
	maxBranch := avail * 30 / 100
	maxLatest := avail - maxRepo - maxBranch
	totalRowWidth := 2 + maxRepo + 3 + maxBranch + 3 + maxLatest + 3 + 10 + 3 + 10 + 3 + 7 + 3 + 4

	return &PullTableLayout{
		TermWidth:   termWidth,
		IsWide:      true,
		ColGap:      defaultPullTableColGapWide,
		MaxRepo:     maxRepo,
		MaxBranch:   maxBranch,
		MaxLatestBr: maxLatest,
		MaxPR:       10,
		MaxStatus:   10,
		MaxSHA:      7,
		DividerLen:  totalRowWidth - 2,
		Rows:        rows,
	}
}

func buildCompactPullTableLayout(rows []model.PullTableRow, termWidth int) *PullTableLayout {
	fixedWidth := 2 + 8 + 10 + 7 + 4 + (5 * defaultPullTableColGapTight)
	avail := termWidth - 2 - fixedWidth
	if avail < 16 {
		avail = 16
	}

	maxRepo := avail * 55 / 100
	maxBranch := avail - maxRepo
	if maxRepo < 10 {
		maxRepo = 10
	}

	if maxBranch < 8 {
		maxBranch = 8
	}

	totalRowWidth := 2 + maxRepo + 2 + maxBranch + 2 + 8 + 2 + 10 + 2 + 7 + 2 + 4

	return &PullTableLayout{
		TermWidth:   termWidth,
		IsWide:      false,
		ColGap:      defaultPullTableColGapTight,
		MaxRepo:     maxRepo,
		MaxBranch:   maxBranch,
		MaxLatestBr: 0,
		MaxPR:       8,
		MaxStatus:   10,
		MaxSHA:      7,
		DividerLen:  totalRowWidth - 2,
		Rows:        rows,
	}
}

func (l *PullTableLayout) PrintHeader() {
	if l.IsWide {
		l.printWideHeader()

		return
	}

	l.printCompactHeader()
}

func (l *PullTableLayout) printWideHeader() {
	fmt.Printf("  %-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %s\n",
		l.MaxRepo, "REPO",
		l.MaxBranch, "BRANCH",
		l.MaxLatestBr, "LATEST BRANCH",
		l.MaxPR, "PR/TRACK",
		l.MaxStatus, "STATUS",
		l.MaxSHA, "SHA",
		"TIME",
	)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}

func (l *PullTableLayout) printCompactHeader() {
	fmt.Printf("  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
		l.MaxRepo, "REPO",
		l.MaxBranch, "BRANCH",
		l.MaxPR, "PR/TRACK",
		l.MaxStatus, "STATUS",
		l.MaxSHA, "SHA",
		"TIME",
	)
	fmt.Printf("  %s\n", strings.Repeat("-", l.DividerLen))
}
