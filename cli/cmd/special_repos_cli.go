// Package cmd — special_repos_cli.go implements CLI dispatch for 'gitmap rs' (repo-secrets), 'gitmap rc' (repo-cache), and 'gitmap cd rs|rc'.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type specialCLIOptions struct {
	subcmd       string
	primaryArg   string
	repoName     string
	slug         string
	ext          string
	isNoPush     bool
	isJSON       bool
	isAutoAccept bool
}

func runCDSpecialRepo(sub string, args []string) error {
	targetPath, _, err := resolveSpecialRepoRoot(sub)
	if err != nil {
		return err
	}
	targetPath = maybeAppendCDSubdir(targetPath, args)
	fmt.Print(targetPath)
	WriteShellHandoff(targetPath)
	warnIfNoWrapper()
	return nil
}

func maybeAppendCDSubdir(targetPath string, args []string) string {
	if len(args) == 0 || len(strings.TrimSpace(args[0])) == 0 {
		return targetPath
	}
	subDir := resolveMatchingRepoSubdir(targetPath, args[0])
	if len(subDir) == 0 {
		return targetPath
	}
	return subDir
}

func resolveMatchingRepoSubdir(rootDir, repoName string) string {
	clean := strings.ToLower(strings.TrimSpace(repoName))
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(strings.ToLower(entry.Name()), clean) {
			return filepath.Join(rootDir, entry.Name())
		}
	}
	return ""
}

func runSpecialRepoCLI(shortKey string, args []string) error {
	_, normShort := store.NormalizeSpecialRepoKey(shortKey)
	if len(args) == 0 || isSpecialHelpArg(args[0]) {
		fmt.Print(RenderSpecialRepoHelp(normShort))
		return nil
	}
	opts := parseSpecialCLIOptions(args)
	return dispatchSpecialRepoSubcmd(normShort, opts)
}

func isSpecialHelpArg(arg string) bool {
	cleaned := strings.ToLower(strings.TrimSpace(arg))
	return cleaned == "help" || cleaned == "-h" || cleaned == "--help"
}

func dispatchSpecialRepoSubcmd(shortKey string, opts specialCLIOptions) error {
	switch opts.subcmd {
	case "put", "add":
		return runSpecialAutoPutAction(shortKey, opts)
	case "file", "f":
		return runSpecialPutAction(shortKey, opts, executeSpecialRepoFile)
	case "folder", "dir", "d":
		return runSpecialPutAction(shortKey, opts, executeSpecialRepoFolder)
	case "text", "note", "t":
		return runSpecialPutAction(shortKey, opts, executeSpecialRepoText)
	case "ls", "list":
		return runSpecialRepoList(shortKey, opts)
	case "init", "ensure":
		return runSpecialRepoInit(resolveSpecialWorkBaseDir())
	case "scan-check":
		return runSpecialRepoScanCheck(opts)
	default:
		fmt.Print(RenderSpecialRepoHelp(shortKey))
		return nil
	}
}

func runSpecialAutoPutAction(shortKey string, opts specialCLIOptions) error {
	if len(strings.TrimSpace(opts.primaryArg)) == 0 {
		return apperror.NewValidationError("path or text required: gitmap " + shortKey + " put <path|text>")
	}
	actionFn := resolveAutoPutActionFunc(opts.primaryArg)
	return runSpecialPutAction(shortKey, opts, actionFn)
}

func resolveAutoPutActionFunc(target string) specialPutFunc {
	info, err := os.Stat(target)
	if err != nil {
		return executeSpecialRepoText
	}
	if info.IsDir() {
		return executeSpecialRepoFolder
	}
	return executeSpecialRepoFile
}

func parseSpecialCLIOptions(args []string) specialCLIOptions {
	opts := specialCLIOptions{subcmd: strings.ToLower(strings.TrimSpace(args[0]))}
	rest := args[1:]
	for i := 0; i < len(rest); i++ {
		i = consumeSpecialOptionToken(&opts, rest, i)
	}
	return opts
}

func consumeSpecialOptionToken(opts *specialCLIOptions, tokens []string, idx int) int {
	tok := tokens[idx]
	switch tok {
	case "--no-push", "-np":
		opts.isNoPush = true
	case "--json", "-j":
		opts.isJSON = true
	case "--auto-accept", "-y", "--yes":
		opts.isAutoAccept = true
	default:
		return consumeKeyValOrPositional(opts, tokens, idx)
	}
	return idx
}

func consumeKeyValOrPositional(opts *specialCLIOptions, tokens []string, idx int) int {
	tok := tokens[idx]
	hasNext := idx+1 < len(tokens)
	if isRepoFlag(tok) && hasNext {
		opts.repoName = tokens[idx+1]
		return idx + 1
	}
	if isSlugFlag(tok) && hasNext {
		opts.slug = tokens[idx+1]
		return idx + 1
	}
	if isExtFlag(tok) && hasNext {
		opts.ext = tokens[idx+1]
		return idx + 1
	}
	return assignPrimaryPositional(opts, tok, idx)
}

func isRepoFlag(tok string) bool {
	return tok == "--repo" || tok == "-r" || tok == "-repo"
}

func isSlugFlag(tok string) bool {
	return tok == "--slug" || tok == "-s" || tok == "--name" || tok == "-n"
}

func isExtFlag(tok string) bool {
	return tok == "--ext" || tok == "-e"
}

func assignPrimaryPositional(opts *specialCLIOptions, tok string, idx int) int {
	if !strings.HasPrefix(tok, "-") && len(opts.primaryArg) == 0 {
		opts.primaryArg = tok
	}
	return idx
}

type specialPutFunc func(db *store.SpecialReposSplitDB, shortKey, specialRoot string, opts specialCLIOptions) (*SpecialPutResult, error)

func runSpecialPutAction(shortKey string, opts specialCLIOptions, fn specialPutFunc) error {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	rootDir, _, err := resolveSpecialRepoRootWithDB(db, shortKey, resolveSpecialWorkBaseDir())
	if err != nil {
		return err
	}
	res, err := fn(db, shortKey, rootDir, opts)
	if err != nil {
		return err
	}
	return printSpecialPutResult(res, opts.isJSON)
}

func executeSpecialRepoFile(db *store.SpecialReposSplitDB, shortKey, specialRoot string, opts specialCLIOptions) (*SpecialPutResult, error) {
	if len(strings.TrimSpace(opts.primaryArg)) == 0 {
		return nil, apperror.NewValidationError("file path is required: gitmap " + shortKey + " file <filepath>")
	}
	repoName := resolveCurrentRepoName(opts.repoName)
	repoFolder, err := db.ResolveOrCreateRepoFolder(shortKey, specialRoot, repoName)
	if err != nil {
		return nil, err
	}
	destPath := allocateSequencedItemPath(repoFolder, filepath.Base(opts.primaryArg))
	if err := copySingleFile(opts.primaryArg, destPath); err != nil {
		return nil, err
	}
	return finalizeSpecialPut(shortKey, specialRoot, repoName, repoFolder, destPath, !opts.isNoPush), nil
}

func executeSpecialRepoFolder(db *store.SpecialReposSplitDB, shortKey, specialRoot string, opts specialCLIOptions) (*SpecialPutResult, error) {
	if len(strings.TrimSpace(opts.primaryArg)) == 0 {
		return nil, apperror.NewValidationError("folder path is required: gitmap " + shortKey + " folder <folderpath>")
	}
	repoName := resolveCurrentRepoName(opts.repoName)
	repoFolder, err := db.ResolveOrCreateRepoFolder(shortKey, specialRoot, repoName)
	if err != nil {
		return nil, err
	}
	destPath := allocateSequencedItemPath(repoFolder, filepath.Base(opts.primaryArg))
	if err := copyDirectoryRecursive(opts.primaryArg, destPath); err != nil {
		return nil, err
	}
	return finalizeSpecialPut(shortKey, specialRoot, repoName, repoFolder, destPath, !opts.isNoPush), nil
}

func executeSpecialRepoText(db *store.SpecialReposSplitDB, shortKey, specialRoot string, opts specialCLIOptions) (*SpecialPutResult, error) {
	repoName := resolveCurrentRepoName(opts.repoName)
	repoFolder, err := db.ResolveOrCreateRepoFolder(shortKey, specialRoot, repoName)
	if err != nil {
		return nil, err
	}
	fileName := deriveTextSlugAndFilename(opts.primaryArg, opts.slug, opts.ext)
	destPath := allocateSequencedItemPath(repoFolder, fileName)
	if err := os.WriteFile(destPath, []byte(opts.primaryArg), 0644); err != nil {
		return nil, apperror.WrapSimple(err, "write text note to special repo")
	}
	return finalizeSpecialPut(shortKey, specialRoot, repoName, repoFolder, destPath, !opts.isNoPush), nil
}

func finalizeSpecialPut(shortKey, specialRoot, repoName, repoFolder, destPath string, shouldPush bool) *SpecialPutResult {
	relPath, err := filepath.Rel(specialRoot, destPath)
	if err != nil {
		relPath = filepath.Base(destPath)
	}
	isCommitted, isPushed := autoCommitAndPushSpecialRepo(specialRoot, relPath, shortKey, repoName, shouldPush)
	return &SpecialPutResult{
		SpecialKey:   shortKey,
		SpecialRoot:  specialRoot,
		RepoName:     repoName,
		RepoFolder:   repoFolder,
		TargetPath:   destPath,
		RelativePath: filepath.ToSlash(relPath),
		IsCommitted:  isCommitted,
		IsPushed:     isPushed,
	}
}

func printSpecialPutResult(res *SpecialPutResult, isJSON bool) error {
	if isJSON {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	fmt.Printf("[%s] Stored %s -> %s (committed=%v, pushed=%v)\n", res.SpecialKey, res.RepoName, res.TargetPath, res.IsCommitted, res.IsPushed)
	return nil
}

// SpecialRepoTreeGroup represents a sequenced project folder and its stored files.
type SpecialRepoTreeGroup struct {
	Folder string   `json:"folder"`
	Items  []string `json:"items"`
}

func runSpecialRepoList(shortKey string, opts specialCLIOptions) error {
	rootDir, rec, err := resolveSpecialRepoRoot(shortKey)
	if err != nil {
		return err
	}
	groups := collectSpecialRepoGroups(rootDir, opts.repoName)
	if opts.isJSON {
		data, _ := json.MarshalIndent(groups, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	renderSpecialRepoTree(shortKey, rec.ConfiguredName, rootDir, groups)
	return nil
}

func collectSpecialRepoGroups(rootDir, filterRepo string) []SpecialRepoTreeGroup {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil
	}
	var out []SpecialRepoTreeGroup
	for _, entry := range entries {
		if isIgnoredSpecialFolder(entry) {
			continue
		}
		if isFilteredSpecialFolder(entry.Name(), filterRepo) {
			continue
		}
		out = append(out, buildRepoGroup(rootDir, entry.Name()))
	}
	return out
}

func isIgnoredSpecialFolder(entry os.DirEntry) bool {
	return !entry.IsDir() || strings.HasPrefix(entry.Name(), ".")
}

func isFilteredSpecialFolder(folderName, filterRepo string) bool {
	if len(filterRepo) == 0 {
		return false
	}
	return !strings.Contains(strings.ToLower(folderName), strings.ToLower(filterRepo))
}

func buildRepoGroup(rootDir, folderName string) SpecialRepoTreeGroup {
	subEntries, _ := os.ReadDir(filepath.Join(rootDir, folderName))
	var items []string
	for _, sub := range subEntries {
		name := sub.Name()
		if sub.IsDir() {
			name += "/"
		}
		items = append(items, name)
	}
	return SpecialRepoTreeGroup{
		Folder: folderName,
		Items:  items,
	}
}

func renderSpecialRepoTree(shortKey, repoName, rootDir string, groups []SpecialRepoTreeGroup) {
	totalItems := countTotalGroupItems(groups)
	fmt.Printf("\n  %s╔══ SPECIAL REPOSITORY: %s (%s) ═══════════════════════════════╗%s\n", constants.ColorCyan, repoName, shortKey, constants.ColorReset)
	fmt.Printf("  │ Path:  %-58s │\n", rootDir)
	fmt.Printf("  │ Total: %d project folder(s), %d item(s)                          │\n", len(groups), totalItems)
	fmt.Printf("  %s╚═════════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	for _, g := range groups {
		renderSingleRepoGroup(g)
	}
}

func countTotalGroupItems(groups []SpecialRepoTreeGroup) int {
	total := 0
	for _, g := range groups {
		total += len(g.Items)
	}
	return total
}

func renderSingleRepoGroup(g SpecialRepoTreeGroup) {
	fmt.Printf("  📁 %s%s/%s (%d items)\n", constants.ColorBold, g.Folder, constants.ColorReset, len(g.Items))
	for idx, it := range g.Items {
		isLast := idx == len(g.Items)-1
		branch := "├─"
		if isLast {
			branch = "└─"
		}
		fmt.Printf("    %s %s\n", branch, it)
	}
	fmt.Println()
}

func runSpecialRepoInit(workBaseDir string) error {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	for _, key := range []string{"rs", "rc"} {
		if err := initSingleSpecialRepo(db, key, workBaseDir); err != nil {
			return err
		}
	}
	return nil
}

func initSingleSpecialRepo(db *store.SpecialReposSplitDB, shortKey, workBaseDir string) error {
	targetDir, rec, err := resolveSpecialRepoRootWithDB(db, shortKey, workBaseDir)
	if err != nil {
		return err
	}
	if err := db.MarkPromptAnswered(shortKey, "initialized", targetDir, rec.RemoteURL); err != nil {
		return err
	}
	fmt.Printf("[%s] Initialized %s at %s\n", shortKey, rec.RepoKey, targetDir)
	return nil
}

func runSpecialRepoScanCheck(opts specialCLIOptions) error {
	records, err := CheckSpecialReposOnScan(resolveSpecialWorkBaseDir(), opts.isAutoAccept)
	if err != nil {
		return err
	}
	if opts.isJSON {
		data, _ := json.MarshalIndent(records, "", "  ")
		fmt.Println(string(data))
	}
	return nil
}

// RenderSpecialRepoHelp returns a boxed terminal help card for 'gitmap rs' or 'gitmap rc'.
func RenderSpecialRepoHelp(shortKey string) string {
	repoKey, normShort := store.NormalizeSpecialRepoKey(shortKey)
	return fmt.Sprintf(`╭── GitMap Special Repository (%s / %s) ──────────────────────────────────────╮
│ Navigation : gitmap cd %s [<repo>]                                          │
│ Auto-Put   : gitmap %s put <path|text> [--repo <name>] [--no-push] [--json] │
│ Store File : gitmap %s file <filepath> [--repo <name>] [--no-push] [--json] │
│ Store Dir  : gitmap %s folder <folder> [--repo <name>] [--no-push] [--json] │
│ Store Note : gitmap %s text "<text>" [--slug <s>] [--ext .ps1] [--repo <n>] │
│ List Tree  : gitmap %s ls [--repo <name>] [--json]                          │
│ Initialize : gitmap %s init | scan-check [--auto-accept]                    │
╰────────────────────────────────────────────────────────────────────────────╯
`, normShort, repoKey, normShort, normShort, normShort, normShort, normShort, normShort, normShort)
}
