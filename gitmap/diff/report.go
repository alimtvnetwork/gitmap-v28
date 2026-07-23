package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Summary is the per-kind tally returned alongside the entries.
type Summary struct {
	MissingLeft  int `json:"missing_left"`
	MissingRight int `json:"missing_right"`
	Conflicts    int `json:"conflicts"`
	Identical    int `json:"identical"`
}

// SummaryFor counts kinds across an Entry slice.
func SummaryFor(entries []Entry) Summary {
	var s Summary
	for _, e := range entries {
		switch e.Kind {
		case MissingLeft:
			s.MissingLeft++
		case MissingRight:
			s.MissingRight++
		case Conflict:
			s.Conflicts++
		case Identical:
			s.Identical++
		}
	}

	return s
}

// PrintOptions controls the human-readable report layout.
type PrintOptions struct {
	OnlyConflicts    bool
	OnlyMissing      bool
	IncludeIdentical bool
	JSON             bool
}

// Report renders the entries to out per opts (text or JSON).
func Report(out io.Writer, entries []Entry, opts PrintOptions) error {
	if opts.JSON {
		return reportJSON(out, entries)
	}

	return reportText(out, entries, opts)
}

// reportJSON emits a single JSON object: { summary, entries }.
func reportJSON(out io.Writer, entries []Entry) error {
	payload := struct {
		Summary Summary `json:"summary"`
		Entries []Entry `json:"entries"`
	}{Summary: SummaryFor(entries), Entries: entries}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")

	return enc.Encode(payload)
}

// reportText prints the four sections honoring the filter flags.
func reportText(out io.Writer, entries []Entry, opts PrintOptions) error {
	summary := SummaryFor(entries)
	if isAllZero(summary) {
		fmt.Fprintf(out, constants.DiffNothingFmt, constants.LogPrefixDiff)

		return nil
	}
	printSection(out, constants.DiffSectionConflicts, entries, Conflict, true)
	if !opts.OnlyConflicts {
		printSection(out, constants.DiffSectionMissingRight, entries, MissingRight, !opts.OnlyMissing || true)
		printSection(out, constants.DiffSectionMissingLeft, entries, MissingLeft, !opts.OnlyMissing || true)
	}
	if opts.IncludeIdentical {
		printSection(out, constants.DiffSectionIdentical, entries, Identical, true)
	}
	fmt.Fprintf(out, constants.DiffSummaryFmt, constants.LogPrefixDiff,
		summary.MissingLeft, summary.MissingRight, summary.Conflicts, summary.Identical)

	return nil
}

// isAllZero returns true when the summary has zero of every kind.
func isAllZero(s Summary) bool {
	return s.MissingLeft == 0 && s.MissingRight == 0 && s.Conflicts == 0 && s.Identical == 0
}

// printSection prints one labeled block of entries of a given kind.
func printSection(out io.Writer, header string, entries []Entry, kind EntryKind, enabled bool) {
	if !enabled {
		return
	}
	matched := filterByKind(entries, kind)
	if len(matched) == 0 {
		return
	}
	fmt.Fprintf(out, "  %s\n", header)
	for _, e := range matched {
		fmt.Fprintf(out, "    %s%s\n", e.RelPath, metaSuffix(e, kind))
	}
	fmt.Fprintln(out)
}

// filterByKind returns the subset of entries matching kind.
func filterByKind(entries []Entry, kind EntryKind) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if e.Kind == kind {
			out = append(out, e)
		}
	}

	return out
}

// metaSuffix renders compact size/mtime info next to a path.
func metaSuffix(e Entry, kind EntryKind) string {
	if kind == Conflict {
		return fmt.Sprintf("  (L: %s @ %s | R: %s @ %s)",
			humanSize(e.LeftSize), humanTime(e.LeftMTime),
			humanSize(e.RightSize), humanTime(e.RightMTime))
	}
	if kind == MissingRight {
		return fmt.Sprintf("  (L: %s @ %s)", humanSize(e.LeftSize), humanTime(e.LeftMTime))
	}
	if kind == MissingLeft {
		return fmt.Sprintf("  (R: %s @ %s)", humanSize(e.RightSize), humanTime(e.RightMTime))
	}

	return ""
}

// humanSize formats bytes as a short string.
func humanSize(b int64) string {
	const k = 1024
	if b < k {
		return fmt.Sprintf("%d B", b)
	}
	if b < k*k {
		return fmt.Sprintf("%.1f KB", float64(b)/float64(k))
	}

	return fmt.Sprintf("%.1f MB", float64(b)/float64(k*k))
}

// humanTime renders a Unix second as "YYYY-MM-DD HH:MM".
func humanTime(unix int64) string {
	if unix == 0 {
		return "-"
	}

	return time.Unix(unix, 0).Format("2006-01-02 15:04")
}
