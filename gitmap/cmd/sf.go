package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runSf dispatches `gitmap sf <add|list|rm>`.
func runSf(args []string) {
	if len(args) == 0 {
		printSfUsage()
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case constants.SFSubAdd:
		runSfAdd(rest)
	case constants.SFSubList, constants.SFSubListAlias:
		runSfList(rest)
	case constants.SFSubRm, constants.SFSubRmAlias:
		runSfRemove(rest)
	default:
		printSfUsage()
		os.Exit(1)
	}
}

// printSfUsage prints the gitmap sf subcommand help.
func printSfUsage() {
	fmt.Fprintln(os.Stderr, constants.MsgSFUsageHeader)
	fmt.Fprintln(os.Stderr, constants.MsgSFUsageAdd)
	fmt.Fprintln(os.Stderr, constants.MsgSFUsageList)
	fmt.Fprintln(os.Stderr, constants.MsgSFUsageRm)
}

// runSfAdd registers a new scan folder.
func runSfAdd(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrSFMissingArg+"\n", "<absolute-path>")
		printSfUsage()
		os.Exit(1)
	}

	pathArg := args[0]
	label, notes := extractSfFlags(args[1:])

	absPath, err := filepath.Abs(pathArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSFAbsResolve+"\n", pathArg, err)
		os.Exit(1)
	}

	db := openSfDB()
	defer db.Close()

	existing, _ := db.ListScanFolders()
	folder, err := db.EnsureScanFolder(absPath, label, notes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if isExistingScanFolder(existing, folder.ID) {
		fmt.Printf(constants.MsgSFAddedExistsFmt, folder.AbsolutePath, folder.ID, folder.LastScannedAt)

		return
	}
	fmt.Printf(constants.MsgSFAddedFmt, folder.AbsolutePath, folder.ID)
}

// runSfList prints every registered scan folder.
func runSfList(_ []string) {
	db := openSfDB()
	defer db.Close()

	folders, err := db.ListScanFolders()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if len(folders) == 0 {
		fmt.Print(constants.MsgSFListEmpty)

		return
	}

	fmt.Printf(constants.MsgSFListHeaderFmt, len(folders))
	for _, f := range folders {
		printSfRow(db, f)
	}
}

// runSfRemove removes a scan folder by absolute path or numeric id.
func runSfRemove(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrSFMissingArg+"\n", "<absolute-path|id>")
		printSfUsage()
		os.Exit(1)
	}

	target := args[0]
	db := openSfDB()
	defer db.Close()

	folder, detached, err := removeSfTarget(db, target)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf(constants.MsgSFRemovedFmt, folder.AbsolutePath, folder.ID, detached)
}

// removeSfTarget resolves whether target is an id or a path and removes it.
func removeSfTarget(db *store.DB, target string) (model.ScanFolder, int, error) {
	if id, err := strconv.ParseInt(target, 10, 64); err == nil && id > 0 {
		return db.RemoveScanFolderByID(id)
	}

	absPath, err := filepath.Abs(target)
	if err != nil {
		return model.ScanFolder{}, 0, fmt.Errorf(constants.ErrSFAbsResolve, target, err)
	}

	return db.RemoveScanFolderByPath(absPath)
}

// extractSfFlags pulls --label and --notes out of the remaining args.
func extractSfFlags(args []string) (string, string) {
	label, notes := "", ""
	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case constants.SFFlagLabel:
			label = args[i+1]
		case constants.SFFlagNotes:
			notes = args[i+1]
		}
	}

	return label, notes
}

// isExistingScanFolder reports whether folder.ID was already in the list
// before the upsert ran (so we can show "already registered" messaging).
func isExistingScanFolder(existing []model.ScanFolder, id int64) bool {
	for _, f := range existing {
		if f.ID == id {
			return true
		}
	}

	return false
}

// printSfRow prints a single list row including the live repo count.
func printSfRow(db *store.DB, f model.ScanFolder) {
	count, err := db.CountReposInScanFolder(f.ID)
	if err != nil {
		count = -1
	}
	label := f.Label
	if len(label) == 0 {
		label = "(none)"
	}
	fmt.Printf(constants.MsgSFListRowFmt, f.ID, f.AbsolutePath, label, count, f.LastScannedAt)
}

// openSfDB opens the default profile DB and runs migrations. Mirrors
// the upsertToDB / scan helpers so `gitmap sf` shares the exact same
// resolution rules used by `gitmap scan`.
func openSfDB() *store.DB {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := db.Migrate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		db.Close()
		os.Exit(1)
	}

	return db
}
