package cmdscan

// scan export / scan merge — portable repo-set JSON.
//
// `scan export` dumps the cached repo list (no rescan) to
// <out>/<machine-slug>/repos.json, a portable JSON array in the exact shape
// `gitmap clone-from` consumes. `scan merge` combines several export folders
// into one deduped JSON. Clone side needs no new code:
// `gitmap clone-from <file> --execute`.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// exportFileName is the JSON file written inside the machine-slug folder.
const exportFileName = "repos.json"

// isScanSubcommand reports whether arg is a scan verb (export|merge).
func isScanSubcommand(arg string) string {
	low := strings.TrimPrefix(strings.ToLower(arg), "--")
	switch low {
	case "export":
		return "export"
	case "merge":
		return "merge"
	default:
		return ""
	}
}

// resolveScanSubcommand splits a leading scan verb from its args.
func resolveScanSubcommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	sub := isScanSubcommand(args[0])
	if sub == "" {
		return "", nil
	}
	return sub, args[1:]
}

// dispatchScanSubcommand routes a scan verb to its handler.
func dispatchScanSubcommand(subcmd string, subArgs []string) error {
	switch subcmd {
	case "export":
		return RunScanExport(subArgs)
	case "merge":
		return RunScanMerge(subArgs)
	default:
		return nil
	}
}

// scanExportHelp prints usage for `gitmap scan export`.
func scanExportHelp() {
	fmt.Println(`Usage: gitmap scan export [--machine <name>] [--out <dir>]

Export the cached repo list to a portable JSON file inside a folder named
by machine slug: <out>/<machine-slug>/repos.json

The JSON array is the exact scan-record shape gitmap clone-from consumes,
so on any machine run:

    gitmap clone-from <out>/<machine-slug>/repos.json --execute

Flags:
    --machine <name>   Machine name for the folder slug (default: hostname)
    --out <dir>        Output directory (default: current directory)`)
}

// scanMergeHelp prints usage for `gitmap scan merge`.
func scanMergeHelp() {
	fmt.Println(`Usage: gitmap scan merge <dir>... [--out <file>]

Merge several scan-export folders into one deduped JSON file.
Each <dir> is an export folder containing repos.json (a direct .json
file path is also accepted). Repos are deduped by URL
(httpsUrl, else sshUrl, else discoveredUrl); first occurrence wins.

Then clone everything in one place:

    gitmap clone-from <file> --execute

Flags:
    --out <file>       Merged output file (default: ./merged-repos.json)`)
}

// exportRecord is a scan record plus the `url` key that `gitmap clone-from`
// requires (its JSON parser reads lowercase `url`; unknown fields are
// tolerated, so the scan-record shape stays intact for other consumers).
type exportRecord struct {
	model.ScanRecord
	URL string `json:"url"`
}

// preferredCloneURL picks the URL clone-from should use: https first (works
// on a fresh machine without SSH keys), then ssh, then discovered.
func preferredCloneURL(r model.ScanRecord) string {
	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}
	if r.SSHUrl != "" {
		return r.SSHUrl
	}
	return r.DiscoveredURL
}

// toExportRecord converts a scan record to the export shape.
func toExportRecord(r model.ScanRecord) exportRecord {
	return exportRecord{ScanRecord: r, URL: preferredCloneURL(r)}
}
func resolveExportMachine(flagMachine string) string {
	if strings.TrimSpace(flagMachine) != "" {
		return strings.TrimSpace(flagMachine)
	}
	name, err := os.Hostname()
	if err != nil || strings.TrimSpace(name) == "" {
		return "localhost"
	}
	return name
}

// hasHelpArg reports whether args contain -h or --help.
func hasHelpArg(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

// splitFlagArgs extracts --name value / --name=value pairs for the known
// string flags, from anywhere in args (stdlib flag stops at the first
// positional, so this pre-pass makes `merge <dirs> --out <file>` work).
// It returns the remaining positional args.
func splitFlagArgs(args []string, defs map[string]*string) []string {
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		matched := false
		for name, ptr := range defs {
			if a == "--"+name && i+1 < len(args) {
				*ptr = args[i+1]
				i++
				matched = true
				break
			}
			if strings.HasPrefix(a, "--"+name+"=") {
				*ptr = strings.TrimPrefix(a, "--"+name+"=")
				matched = true
				break
			}
		}
		if !matched {
			rest = append(rest, a)
		}
	}
	return rest
}

// RunScanExport implements `gitmap scan export`.
func RunScanExport(args []string) error {
	if hasHelpArg(args) {
		scanExportHelp()
		return nil
	}
	var machineFlag, outFlag string
	outFlag = "."
	rest := splitFlagArgs(args, map[string]*string{"machine": &machineFlag, "out": &outFlag})
	if len(rest) > 0 {
		return fmt.Errorf("scan export takes no positional args (got %q) — see: gitmap scan export --help", rest[0])
	}

	machine := resolveExportMachine(machineFlag)
	slug := store.SanitizeSlug(machine)

	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDBOpen)
	}
	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDBQuery)
	}

	targetDir := filepath.Join(outFlag, slug)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("could not create export folder %s: %w", targetDir, err)
	}
	path := filepath.Join(targetDir, exportFileName)
	if err := writeRecordsJSON(path, records); err != nil {
		return err
	}

	fmt.Printf("  exported %d repos — machine %q (slug %q) → %s\n", len(records), machine, slug, path)
	fmt.Printf("  clone on any machine: gitmap clone-from %s --execute\n", path)
	return nil
}

// RunScanMerge implements `gitmap scan merge`.
func RunScanMerge(args []string) error {
	if hasHelpArg(args) {
		scanMergeHelp()
		return nil
	}
	var outFlag string
	outFlag = "merged-repos.json"
	dirs := splitFlagArgs(args, map[string]*string{"out": &outFlag})

	if len(dirs) == 0 {
		scanMergeHelp()
		return fmt.Errorf("scan merge needs at least one export folder")
	}

	var merged []exportRecord
	seen := make(map[string]bool)
	total := 0
	for _, d := range dirs {
		recs, path, err := readExportRecords(d)
		if err != nil {
			return err
		}
		total += len(recs)
		for _, r := range recs {
			key := mergeDedupeKey(r)
			if seen[key] {
				continue
			}
			seen[key] = true
			merged = append(merged, r)
		}
		fmt.Printf("  %s: %d repos\n", path, len(recs))
	}

	if err := writeExportRecordsJSON(outFlag, merged); err != nil {
		return err
	}

	fmt.Printf("  merged %d repos (%d unique) → %s\n", total, len(merged), outFlag)
	fmt.Printf("  clone in one place: gitmap clone-from %s --execute\n", outFlag)
	return nil
}

// readExportRecords reads an export JSON array from an export folder
// (<dir>/repos.json) or a direct .json file path. Files written by older
// versions (scan-record shape without `url`) get the URL backfilled.
func readExportRecords(dirOrFile string) ([]exportRecord, string, error) {
	path := dirOrFile
	if info, err := os.Stat(dirOrFile); err == nil && info.IsDir() {
		path = filepath.Join(dirOrFile, exportFileName)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, path, fmt.Errorf("could not read %s: %w", path, err)
	}
	var recs []exportRecord
	if err := json.Unmarshal(data, &recs); err != nil {
		return nil, path, fmt.Errorf("could not parse %s as repo JSON: %w", path, err)
	}
	for i := range recs {
		if recs[i].URL == "" {
			recs[i].URL = preferredCloneURL(recs[i].ScanRecord)
		}
	}
	return recs, path, nil
}

// mergeDedupeKey returns the identity key for merge dedup: the clone URL,
// else the scan slug.
func mergeDedupeKey(r exportRecord) string {
	if r.URL != "" {
		return "url:" + r.URL
	}
	if r.HTTPSUrl != "" {
		return "https:" + r.HTTPSUrl
	}
	if r.SSHUrl != "" {
		return "ssh:" + r.SSHUrl
	}
	if r.DiscoveredURL != "" {
		return "discovered:" + r.DiscoveredURL
	}
	return "slug:" + r.Slug
}

// writeRecordsJSON writes records as an indented JSON array in the export
// shape (scan record + `url`), directly consumable by `gitmap clone-from`.
func writeRecordsJSON(path string, records []model.ScanRecord) error {
	exported := make([]exportRecord, 0, len(records))
	for _, r := range records {
		exported = append(exported, toExportRecord(r))
	}
	return writeExportRecordsJSON(path, exported)
}

// writeExportRecordsJSON writes export-shaped records as indented JSON,
// mirroring formatter.WriteJSON's indentation.
func writeExportRecordsJSON(path string, records []exportRecord) error {
	if records == nil {
		records = []exportRecord{}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", constants.JSONIndent)
	if err := enc.Encode(records); err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}
	return nil
}
