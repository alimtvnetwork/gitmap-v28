package cmderrors

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const errorsLastHelp = `errors last — surface the last tracked internal error.

Description:
  Reads the most recent record from the error module's tracking store
  (gitmap-errors.db): the real message, its fallback context, and where
  it came from. Prints a clean "no tracked errors" message when the
  store is empty.

Usage:
  gitmap errors last [--json]
  gitmap errl [--json]

Flags:
  --json        output the record as JSON
  --help, -h    show this help

Example:
  gitmap errors last
`

// RunErrorsLast implements `gitmap errors last` (alias: errl).
func RunErrorsLast(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(errorsLastHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("errors last", flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "output as JSON")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(errorsLastHelp)
		return apperror.WrapSimple(parseErr, "errors.last.flags")
	}

	splitDB, openErr := store.OpenErrorsSplitDB()
	if openErr != nil {
		return apperror.WrapSimple(openErr, "errors.last.open")
	}
	defer splitDB.Close()

	records, listErr := splitDB.FederatedListErrors(1, false)
	if listErr != nil {
		if strings.Contains(strings.ToLower(listErr.Error()), "no such table") {
			fmt.Println("errors last: no tracked errors recorded. the error module has nothing on file.")
			return nil
		}
		return apperror.WrapSimple(listErr, "errors.last.list")
	}
	if len(records) == 0 {
		fmt.Println("errors last: no tracked errors recorded. the error module has nothing on file.")
		return nil
	}

	rec := records[0]
	if *asJSON {
		raw, jsonErr := json.MarshalIndent(rec, "", "  ")
		if jsonErr != nil {
			return apperror.WrapSimple(jsonErr, "errors.last.json")
		}
		fmt.Println(string(raw))

		return nil
	}

	message := strings.ReplaceAll(strings.SplitN(rec.Message, "\n", 2)[0], "\r", "")
	if len(message) > 160 {
		message = message[:157] + "..."
	}
	fmt.Println("last tracked error:")
	fmt.Printf("  id:        %d\n", rec.ID)
	fmt.Printf("  time:      %s\n", rec.CreatedAt)
	fmt.Printf("  code:      %s\n", rec.ErrorCode)
	fmt.Printf("  type:      %s\n", rec.ErrorType)
	fmt.Printf("  command:   %s\n", rec.Command)
	fmt.Printf("  message:   %s\n", message)
	if strings.TrimSpace(rec.RepoPath) != "" {
		fmt.Printf("  repo:      %s\n", rec.RepoPath)
	}
	if strings.TrimSpace(rec.SourceFile) != "" {
		fmt.Printf("  source:    %s\n", rec.SourceFile)
	}
	resolved := "no"
	if rec.IsResolved {
		resolved = "yes"
	}
	fmt.Printf("  resolved:  %s\n", resolved)

	return nil
}
