package constants

// Constants for the `--output terminal` summary block emitted by
// `gitmap clone-from --execute`. Split from constants_clonefrom.go
// to keep both files under the project's 200-line cap. All printf
// format verbs are documented inline so future translators / log
// parsers can map %v back to call sites.

// CloneFromTermSummaryHeader is the banner that opens the terminal-
// mode summary. Plain ASCII so the byte width matches the terminal
// column count Windows cmd.exe assumes (mem://constraints/powershell-encoding).
const CloneFromTermSummaryHeader = "\ngitmap clone-from: summary\n"

// CloneFromTermSummaryFoundFmt prints the "found" total — the number
// of repos in the user's manifest after dedup, NOT the number that
// were successfully cloned (use the status tally for that). Wording
// matches the user's phrasing in the request ("count of found repos").
// %d = total row count.
const CloneFromTermSummaryFoundFmt = "  found:    %d repo(s)\n"

// CloneFromTermSummarySchemeHeader introduces the per-scheme block.
// Indent matches the rest of the block so the renderer never has
// to compute column widths at runtime.
const CloneFromTermSummarySchemeHeader = "  by mode:\n"

// CloneFromTermSummarySchemeRowFmt renders one scheme row inside the
// "by mode" block. Two-space indent INSIDE the block (so 4 from the
// left edge) keeps the visual hierarchy clear without ANSI escapes.
// %-7s = left-padded scheme label (https / ssh / git / file / scp /
// other), %d = count.
const CloneFromTermSummarySchemeRowFmt = "    %-7s %d\n"

// CloneFromTermSummaryStatusFmt is the one-line status tally inside
// the terminal summary. Mirrors the wording of the legacy
// MsgCloneFromSummaryHeader so users who already know the default
// summary recognize the numbers. %d ok, %d skipped, %d failed,
// %d total.
const CloneFromTermSummaryStatusFmt = "  status:   %d ok, %d skipped, %d failed (%d total)\n"

// CloneFromTermSummaryReportFmt is one report-path line. Used twice
// (once for CSV, once for JSON) so the same format string covers
// both. %s = label ("csv"/"json"), %s = absolute path.
const CloneFromTermSummaryReportFmt = "  report %s: %s\n"

// CloneFromTermSummaryReportNoneFmt is the placeholder when --no-report
// is in effect or report writing failed. Shown in place of the path
// rows so the summary shape stays predictable for log scrapers.
const CloneFromTermSummaryReportNone = "  report:   (skipped — --no-report or write failed)\n"

// CloneFromSummaryTransportFmt is the one-line transport split shared
// by both the legacy RenderSummary and the enriched terminal block.
// SSH = ssh:// + scp-style; HTTPS = https://; OTHER folds http://,
// git://, file://, and unrecognized forms so the line stays at three
// stable columns regardless of manifest contents. Derived from
// ClassifyScheme so this counter and the per-scheme tally can never
// disagree. %d ssh, %d https, %d other.
const CloneFromSummaryTransportFmt = "transport: %d ssh, %d https, %d other\n"

// URL scheme labels — surfaced in the per-mode tally. Stable strings:
// renaming would break any downstream tooling that greps the
// terminal summary. "scp" covers the `[user@]host:path` form that
// validate.go's looksLikeSCP accepts.
const (
	CloneFromSchemeHTTPS = "https"
	CloneFromSchemeHTTP  = "http"
	CloneFromSchemeSSH   = "ssh"
	CloneFromSchemeGit   = "git"
	CloneFromSchemeFile  = "file"
	CloneFromSchemeSCP   = "scp"
	CloneFromSchemeOther = "other"
)

// JSON report on-disk filename pattern. Mirrors the CSV pattern so
// both files sort next to each other in `ls .gitmap/`. %d = unix
// seconds at write time.
const CloneFromReportJSONNameFmt = "clone-from-report-%d.json"
