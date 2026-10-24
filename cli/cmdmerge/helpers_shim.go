package cmdmerge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type moveOpts struct {
	yes           bool
	dryRun        bool
	isSkipVSCode  bool
	isSkipDesktop bool
}

// ResolveRepo finds a single repo matching target string across DB and JSON records.
func ResolveRepo(db *store.DB, target string) (*model.ScanRecord, error) {
	all := loadUnifiedCandidates(db)
	t := strings.TrimSpace(target)
	if len(t) == 0 || t == "." {
		return resolveByPWD(all)
	}

	if rec := resolveByPath(t, all); rec != nil {
		return rec, nil
	}

	if rec := resolveByAlias(db, t, all); rec != nil {
		return rec, nil
	}

	if rec := resolveByTarget(t, all); rec != nil {
		return rec, nil
	}

	return nil, apperror.NewNotFoundError(fmt.Sprintf("no repository matched %q", target))
}

func calculateDestPath(srcPath, destTarget string) (string, error) {
	cleanSrc := filepath.Clean(filepath.FromSlash(fsutil.TrimTrailingSlashes(srcPath)))
	repoName := filepath.Base(cleanSrc)
	t := fsutil.TrimTrailingSlashes(destTarget)

	if t == ".." {
		parent := filepath.Dir(cleanSrc)
		grandParent := filepath.Dir(parent)

		return filepath.Join(grandParent, repoName), nil
	}

	absDest, err := filepath.Abs(t)
	if err != nil {
		return "", err
	}

	if fsutil.DirExists(absDest) {
		return filepath.Join(absDest, repoName), nil
	}

	return absDest, nil
}

func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return store.OpenGlobalDefault()
	}
	_ = db.Migrate()
	if countRegisteredSSHHosts(db) == 0 {
		return resolveGlobalSSHDBFallback(db), nil
	}
	return db, nil
}

func parseMoveFlags(args []string) (moveOpts, []string) {
	opts := moveOpts{}
	var positional []string
	for _, a := range args {
		switch a {
		case "-y", "--yes":
			opts.yes = true
		case "--dry-run":
			opts.dryRun = true
		case "--no-vscode":
			opts.isSkipVSCode = true
		case "--no-desktop":
			opts.isSkipDesktop = true
		default:
			positional = append(positional, a)
		}
	}

	return opts, positional
}

func preflightMove(srcPath, destPath string) error {
	if !fsutil.DirExists(srcPath) {
		return apperror.NewNotFoundError(fmt.Sprintf("source directory %q does not exist", srcPath))
	}

	if fsutil.EqualPaths(srcPath, destPath) {
		return apperror.NewValidationError("source and destination are identical")
	}

	if fsutil.FileOrDirExists(destPath) {
		return apperror.NewValidationError(fmt.Sprintf("destination %q already exists", destPath))
	}

	if fsutil.IsSubdirectory(srcPath, destPath) {
		return apperror.NewValidationError("cannot move repository into its own subdirectory")
	}

	return nil
}

func printMoveDryRun(src, dest string) {
	fmt.Println("  (dry-run) Move operation preview:")
	fmt.Printf("     Source:      %s\n", src)
	fmt.Printf("     Destination: %s\n", dest)
	fmt.Println("     Actions:     Relocate directory, update SQLite Repo path, sync VS Code & GitHub Desktop")
}

func resolveInspectionFiles(fileArgs []string, targetDir string) ([]string, error) {
	if len(fileArgs) > 0 {
		return expandFileGlobs(fileArgs), nil
	}
	return scanDirectoryJSONFiles(targetDir)
}

func countRegisteredSSHHosts(db *store.DB) int {
	var count int
	row := db.SQL().QueryRow("SELECT count(*) FROM ssh_hosts")
	if scanErr := row.Scan(&count); scanErr != nil {
		return 0
	}
	return count
}
func resolveGlobalSSHDBFallback(localDB *store.DB) *store.DB {
	globalDB, err := store.OpenGlobalDefault()
	if err != nil {
		return localDB
	}
	_ = globalDB.Migrate()
	if countRegisteredSSHHosts(globalDB) > 0 {
		localDB.Close()
		return globalDB
	}
	globalDB.Close()
	return localDB
}

func loadUnifiedCandidates(db *store.DB) []model.ScanRecord {
	dbRepos := loadDbCandidates(db)
	jsonRepos := loadJSONCandidates()

	return mergeCandidateRepos(dbRepos, jsonRepos)
}

func loadDbCandidates(db *store.DB) []model.ScanRecord {
	if db == nil {
		return nil
	}

	repos, err := db.ListRepos()
	if err != nil {
		return nil
	}

	return repos
}

func loadJSONCandidates() []model.ScanRecord {
	paths := []string{
		filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile),
		filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
	}
	for _, p := range paths {
		if records, isOk := tryLoadJSONPath(p); isOk {
			return records
		}
	}

	return nil
}

func tryLoadJSONPath(path string) ([]model.ScanRecord, bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, false
	}

	records, err := model.LoadStatusRecords(path)
	if err != nil || len(records) == 0 {
		return nil, false
	}

	return records, true
}

func mergeCandidateRepos(primary, secondary []model.ScanRecord) []model.ScanRecord {
	seen := make(map[string]bool)
	var out []model.ScanRecord

	for _, r := range primary {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	for _, r := range secondary {
		key := fsutil.NormalizeSlashes(r.AbsolutePath)
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}

	return out
}

func resolveByPWD(all []model.ScanRecord) (*model.ScanRecord, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, apperror.WrapSimple(err, "get current directory")
	}

	for _, r := range all {
		if fsutil.EqualPaths(r.AbsolutePath, pwd) {
			return &r, nil
		}
	}

	return nil, apperror.NewNotFoundError(fmt.Sprintf("current directory (%s) is not a tracked repository", pwd))
}
