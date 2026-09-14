package cmdpull

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PullArrayUIPart1 struct{}

func InitPullArrayUIPart1() error          { return nil }
func (x *PullArrayUIPart1) Process1() bool { return true }

type PullArrayUIPart2 struct{}

func InitPullArrayUIPart2() error          { return nil }
func (x *PullArrayUIPart2) Process2() bool { return true }

// FormatPullDuration formats duration into ms or seconds string.
func FormatPullDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}

	return fmt.Sprintf("%.1fs", d.Seconds())
}

// ShortenSHA trims a commit SHA to 7 characters.
func ShortenSHA(sha string) string {
	clean := strings.TrimSpace(sha)
	if len(clean) > 7 {
		return clean[:7]
	}

	return clean
}

// FormatCommitRange formats old and new SHAs into a range or single SHA.
func FormatCommitRange(oldSHA, newSHA string) string {
	sOld := ShortenSHA(oldSHA)
	sNew := ShortenSHA(newSHA)
	if sOld == "" && sNew == "" {
		return ""
	}

	return formatValidCommitRange(sOld, sNew)
}

func formatValidCommitRange(sOld, sNew string) string {
	if sOld == "" || sOld == sNew {
		return sNew
	}
	if sNew == "" {
		return sOld
	}

	return fmt.Sprintf("%s..%s", sOld, sNew)
}

// FormatDiffStat converts git diffstat string into compact metrics.
func FormatDiffStat(statLine string) string {
	insertions := extractStatNumber(statLine, "insertion")
	deletions := extractStatNumber(statLine, "deletion")
	files := extractStatNumber(statLine, "file")
	if insertions == 0 && deletions == 0 {
		return formatZeroDiffStat(files)
	}

	return fmt.Sprintf("+%d/-%d (%d)", insertions, deletions, files)
}

func formatZeroDiffStat(files int) string {
	if files > 0 {
		return fmt.Sprintf("%d files", files)
	}

	return "synced"
}

func extractStatNumber(statLine, keyword string) int {
	idx := strings.Index(statLine, keyword)
	if idx <= 0 {
		return 0
	}

	return parsePrecedingNumber(statLine[:idx])
}

func parsePrecedingNumber(sub string) int {
	fields := strings.Fields(strings.TrimSpace(sub))
	if len(fields) == 0 {
		return 0
	}

	val, err := strconv.Atoi(fields[len(fields)-1])
	if err != nil {
		return 0
	}

	return val
}
