// Package cmdstale provides `gitmap stale` (sta): list local repos with no commits
// in N days, with optional --archive to move them to .gitmap/archive/.
package cmdstale

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhygiene"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale/enums"
)

type staleRepo struct {
	path string
	last time.Time
}

type staleOptions struct {
	staleDays    int
	scanRoot     string
	isArchive    bool
	isDryRun     bool
	outputFormat enums.StaleOutputFormatType
}

type staleRawFlags struct {
	days    *int
	root    *string
	archive *bool
	dryRun  *bool
	format  *string
}

func bindStaleFlags(fs *flag.FlagSet) staleRawFlags {
	return staleRawFlags{
		days:    fs.Int(constants.FlagDays, constants.DefaultStaleDays, constants.FlagDescDays),
		root:    fs.String(constants.FlagRoot, constants.DefaultScanRoot, constants.FlagDescRoot),
		archive: fs.Bool(constants.FlagArchive, false, constants.FlagDescArchive),
		dryRun:  fs.Bool(constants.FlagDryRun, false, constants.FlagDescDryRun),
		format:  fs.String(constants.FlagFormat, constants.DefaultOutputFormat, constants.FlagDescFormat),
	}
}

func buildStaleOptions(raw staleRawFlags, fmtKind enums.StaleOutputFormatType) staleOptions {
	return staleOptions{
		staleDays:    *raw.days,
		scanRoot:     *raw.root,
		isArchive:    *raw.archive,
		isDryRun:     *raw.dryRun,
		outputFormat: fmtKind,
	}
}

func parseStaleFlags(args []string) (staleOptions, error) {
	fs := flag.NewFlagSet(constants.CmdStale, flag.ContinueOnError)
	raw := bindStaleFlags(fs)
	if err := fs.Parse(args); err != nil {
		return staleOptions{}, appfault.NewValidationError(err.Error())
	}

	fmtKind, err := parseHygieneFormat(*raw.format)
	if err != nil {
		return staleOptions{}, err
	}

	return buildStaleOptions(raw, fmtKind), nil
}

// RunStale executes `gitmap stale`.
func RunStale(args []string) error {
	checkHelp(constants.CmdStale, args)
	opts, err := parseStaleFlags(args)
	if err != nil {
		return err
	}

	probed := probeStaleRepos(opts.scanRoot, opts.staleDays)
	emitStale(probed, opts.staleDays, opts.outputFormat)
	if opts.isArchive {
		archiveStaleRepos(probed, opts.isDryRun)
	}

	return nil
}

func probeStaleRepos(root string, days int) []staleRepo {
	repos := scanForReposParallel(root)
	cutoff := time.Now().AddDate(0, 0, -days)
	probed := cmdhygiene.MapReposParallel(repos, func(r string) (staleRepo, bool) {
		return probeRepoStaleness(r, cutoff)
	})
	sort.Slice(probed, func(i, j int) bool { return probed[i].last.Before(probed[j].last) })

	return probed
}

func probeRepoStaleness(r string, cutoff time.Time) (staleRepo, bool) {
	t, hasCommit := lastCommitTime(r)
	if !hasCommit || !t.Before(cutoff) {
		return staleRepo{}, false
	}

	return staleRepo{path: r, last: t}, true
}

// emitStale dispatches to the requested output format.
func emitStale(stale []staleRepo, days int, f enums.StaleOutputFormatType) {
	if f.IsJSON() {
		emitStaleJSON(stale)

		return
	}
	if f.IsCSV() {
		emitStaleCSV(stale)

		return
	}

	printStaleTable(stale, days)
}

type staleJSONRow struct {
	Path    string `json:"path"`
	LastUTC string `json:"last_commit_utc"`
	AgeDays int    `json:"age_days"`
}

func buildStaleJSONRows(stale []staleRepo, now time.Time) []staleJSONRow {
	rows := make([]staleJSONRow, 0, len(stale))
	for _, s := range stale {
		rows = append(rows, staleJSONRow{
			Path:    s.path,
			LastUTC: s.last.UTC().Format(time.RFC3339),
			AgeDays: int(now.Sub(s.last).Hours() / constants.HoursPerDay),
		})
	}

	return rows
}

func emitStaleJSON(stale []staleRepo) {
	rows := buildStaleJSONRows(stale, time.Now())
	emitJSON(rows)
}

func buildStaleCSVRows(stale []staleRepo, now time.Time) [][]string {
	rows := make([][]string, 0, len(stale))
	for _, s := range stale {
		age := int(now.Sub(s.last).Hours() / constants.HoursPerDay)
		ageStr := fmt.Sprintf(constants.AgeDaysIntFmt, age)
		rows = append(rows, []string{s.path, s.last.UTC().Format(time.RFC3339), ageStr})
	}

	return rows
}

func emitStaleCSV(stale []staleRepo) {
	header := []string{constants.HeaderPath, constants.HeaderLastCommitUTC, constants.HeaderAgeDays}
	rows := buildStaleCSVRows(stale, time.Now())
	emitCSV(header, rows)
}

// printStaleTable renders the stale-repo list.
func printStaleTable(stale []staleRepo, days int) {
	if len(stale) == 0 {
		fmt.Fprintf(os.Stdout, constants.MsgNoStaleReposFmt, days)

		return
	}

	printStaleRows(stale, days)
}

func printStaleRows(stale []staleRepo, days int) {
	fmt.Fprintf(os.Stdout, constants.MsgStaleHeaderFmt, len(stale), days)
	now := time.Now()
	for _, s := range stale {
		age := int(now.Sub(s.last).Hours() / constants.HoursPerDay)
		fmt.Fprintf(os.Stdout, constants.MsgStaleRowFmt, age, s.path, s.last.Format(constants.DateFormatStale))
	}

	fmt.Fprintln(os.Stdout, constants.EmptyString)
}

// LastCommitTime returns the last commit time for repo at dir.
func LastCommitTime(dir string) (time.Time, bool) {
	return lastCommitTime(dir)
}

func lastCommitTime(dir string) (time.Time, bool) {
	cmd := exec.Command(constants.GitBin, constants.GitFlagDir, dir, constants.GitSubCmdLog, constants.GitArgLimitOne, constants.GitFormatCommitEpoch)
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, false
	}

	return parseEpochOutput(string(out))
}

func parseEpochOutput(raw string) (time.Time, bool) {
	s := strings.TrimSpace(raw)
	if s == constants.EmptyString {
		return time.Time{}, false
	}

	var sec int64
	if _, err := fmt.Sscanf(s, constants.SecScanFmt, &sec); err != nil {
		return time.Time{}, false
	}

	return time.Unix(sec, 0), true
}

// archiveStaleRepos moves repos into .gitmap/archive/<timestamp>/.
func archiveStaleRepos(stale []staleRepo, isDryRun bool) {
	if len(stale) == 0 {
		return
	}

	stamp := time.Now().UTC().Format(constants.TimestampFormatArchive)
	archiveDir := filepath.Join(constants.DirGitmap, constants.DirArchive, stamp)
	for _, s := range stale {
		archiveSingleRepo(s.path, archiveDir, isDryRun)
	}
}

func archiveSingleRepo(srcPath, archiveDir string, isDryRun bool) {
	dest := filepath.Join(archiveDir, filepath.Base(srcPath))
	if isDryRun {
		fmt.Fprintf(os.Stdout, constants.MsgWouldArchiveFmt, srcPath, dest)

		return
	}

	executeArchiveMove(srcPath, dest, archiveDir)
}

func executeArchiveMove(srcPath, dest, archiveDir string) {
	if err := os.MkdirAll(archiveDir, constants.DefaultArchiveDirPerm); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrArchiveMkdirFmt, err)

		return
	}
	if err := os.Rename(srcPath, dest); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgFailedFmt, srcPath, err)

		return
	}

	fmt.Fprintf(os.Stdout, constants.MsgArchivedFmt, srcPath, dest)
}
