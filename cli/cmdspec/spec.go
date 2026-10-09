package cmdspec

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// specUsage is printed for `gitmap spec --help` or bare `gitmap spec`.
const specUsage = `Usage: gitmap spec <subcommand> [flags]

Spec-number issuance namespace (spec 252). Subcommands:

  next        Issue the next free spec number for the current
              repository and print it on stdout.

Flags (with next):
  --json          Print the issuance record as JSON instead of a
                  bare number
  -h, --help      Print usage, exit 0

Behavior:
  - The candidate is max(highest ^(\d+)- directory prefix under
    02-spec/21-app/, highest number recorded in the Tier-1 agent
    database) + 1.
  - The claim is atomic (INSERT OR IGNORE over a UNIQUE(number)
    column): parallel callers each get a distinct number; the
    loser retries to the next free one.
  - Plain mode prints ONLY the number on stdout (script-friendly).

Examples:
  gitmap spec next
  gitmap spec next --json
  gitmap spec --help
`

// specNextUsage is printed for `gitmap spec next --help`.
const specNextUsage = `Usage: gitmap spec next [--json]

Issues the next free spec number for the current repository and
prints it. Plain mode prints ONLY the number on stdout (e.g. 252).
With --json, stdout is exactly one pretty-printed object:

  {
    "repo": "<absolute repo path>",
    "number": 252,
    "issued_at": "<RFC3339>"
  }

Flags:
  --json          Print the issuance record as JSON
  -h, --help      Print this usage, exit 0

Exit codes:
  0   number issued (plain or --json), or usage printed
  2   usage error (unknown flag)
  1   issuance failure (DB open/scan/claim failed); the error goes
      to stderr and stdout contains no number
`

// specClaimAttempts bounds the claim retry loop (spec 252 §3.5).
const specClaimAttempts = 1000

// specDirPrefixPattern matches the leading-number prefix of spec
// directories, e.g. "252-spec-number-issuance".
var specDirPrefixPattern = regexp.MustCompile(`^(\d+)-`)

// DispatchSpec routes `gitmap spec <subcommand>`; returns (false, nil)
// for anything it does not own.
func DispatchSpec(command string) (bool, error) {
	if command != constants.CmdSpec {
		return false, nil
	}

	rest := os.Args[2:]
	if len(rest) == 0 {
		fmt.Print(specUsage)

		return true, nil
	}

	sub := rest[0]
	if sub == "--help" || sub == "-h" {
		fmt.Print(specUsage)

		return true, nil
	}

	if sub == constants.CmdSpecNext {
		return true, runSpecNext(rest[1:])
	}

	fmt.Print(specUsage)

	return true, fmt.Errorf("unknown spec subcommand: %s", sub)
}

// runSpecNext implements `gitmap spec next [--json]`.
func runSpecNext(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(specNextUsage)

			return nil
		}
	}

	fs := flag.NewFlagSet(constants.CmdSpecNext, flag.ExitOnError)
	asJSON := fs.Bool("json", false, "print the issuance record as JSON")
	if parseErr := fs.Parse(args); parseErr != nil {
		return parseErr
	}

	repoRoot := absSpecRepoRoot(findSpecRepoRoot())
	// Anchor the Tier-1 DB to the spec root (where 02-spec/ lives), NOT
	// to cwd: ResolveMasterAgentDbPath("") walks up for .gitmap/.git and
	// a nested .gitmap dir (e.g. cli/.gitmap) would resolve a SECOND
	// database, letting two agents issue the same number from two DBs.
	conn, openErr := store.InitMasterAgentDB(store.ResolveMasterAgentDbPath(repoRoot))
	if openErr != nil {
		return openErr
	}
	defer conn.Close()

	if schemaErr := store.EnsureSpecNumbersSchema(conn); schemaErr != nil {
		return schemaErr
	}

	candidate, seedErr := nextSpecCandidate(conn, repoRoot)
	if seedErr != nil {
		return seedErr
	}

	issued, claimErr := claimSpecNumber(conn, repoRoot, candidate)
	if claimErr != nil {
		return claimErr
	}

	return printIssued(repoRoot, issued, *asJSON)
}

// nextSpecCandidate returns max(disk scan, db max) + 1.
func nextSpecCandidate(conn *sql.DB, repoRoot string) (int, error) {
	candidate := scanSpecDirMax(repoRoot) + 1

	dbMax, dbErr := store.MaxIssuedSpecNumber(conn)
	if dbErr != nil {
		return 0, dbErr
	}
	if dbMax+1 > candidate {
		candidate = dbMax + 1
	}

	return candidate, nil
}

// scanSpecDirMax returns the highest leading-number directory prefix
// under <repoRoot>/02-spec/21-app/. Non-matching entries, files, and a
// missing or unreadable directory all contribute 0.
func scanSpecDirMax(repoRoot string) int {
	entries, readErr := os.ReadDir(filepath.Join(repoRoot, "02-spec", "21-app"))
	if readErr != nil {
		return 0
	}

	max := 0
	for _, entry := range entries {
		n := specDirNumber(entry)
		if n > max {
			max = n
		}
	}

	return max
}

// specDirNumber extracts the leading ^(\d+)- number from a directory
// entry, or 0 when the entry is a file or has no numeric prefix.
func specDirNumber(entry os.DirEntry) int {
	if !entry.IsDir() {
		return 0
	}
	match := specDirPrefixPattern.FindStringSubmatch(entry.Name())
	if len(match) < 2 {
		return 0
	}
	n, convErr := strconv.Atoi(match[1])
	if convErr != nil {
		return 0
	}

	return n
}

// findSpecRepoRoot walks up from the working directory to the nearest
// directory containing 02-spec/, or "" when none is found.
func findSpecRepoRoot() string {
	dir, wdErr := os.Getwd()
	if wdErr != nil {
		return ""
	}
	for {
		if isDirPresent(filepath.Join(dir, "02-spec")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// absSpecRepoRoot returns the absolute repo root, falling back to the
// absolute working directory when no 02-spec/ ancestor was found.
func absSpecRepoRoot(repoRoot string) string {
	if repoRoot == "" {
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			return "."
		}

		return wd
	}
	abs, absErr := filepath.Abs(repoRoot)
	if absErr != nil {
		return repoRoot
	}

	return abs
}

// isDirPresent reports whether path exists and is a directory.
func isDirPresent(path string) bool {
	info, statErr := os.Stat(path)
	if statErr != nil {
		return false
	}

	return info.IsDir()
}

// claimSpecNumber atomically claims the first free number at or above
// candidate. A lost race (INSERT OR IGNORE reports zero rows) moves to
// the next number; genuine DB failures abort immediately.
func claimSpecNumber(conn *sql.DB, repoRoot string, candidate int) (int, error) {
	for attempt := 0; attempt < specClaimAttempts; attempt++ {
		claimed, claimErr := store.ClaimSpecNumber(conn, candidate, repoRoot, "")
		if claimErr != nil {
			return 0, claimErr
		}
		if claimed {
			return candidate, nil
		}
		candidate++
	}

	return 0, fmt.Errorf("spec next: no free number after %d attempts", specClaimAttempts)
}

// specIssuance is the --json record emitted by `gitmap spec next`.
type specIssuance struct {
	Repo     string `json:"repo"`
	Number   int    `json:"number"`
	IssuedAt string `json:"issued_at"`
}

// printIssued prints the issued number: bare on stdout, or the JSON
// issuance record with --json. Nothing else reaches stdout.
func printIssued(repoRoot string, number int, asJSON bool) error {
	if !asJSON {
		fmt.Println(number)

		return nil
	}

	record := specIssuance{
		Repo:     repoRoot,
		Number:   number,
		IssuedAt: time.Now().UTC().Format(time.RFC3339),
	}
	encoded, encodeErr := json.MarshalIndent(record, "", "  ")
	if encodeErr != nil {
		return encodeErr
	}
	fmt.Println(string(encoded))

	return nil
}
