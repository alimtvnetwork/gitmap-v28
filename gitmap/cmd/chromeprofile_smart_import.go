package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"gopkg.in/yaml.v3"
)

type importDestination struct {
	Dir         string
	Path        string
	DisplayName string
	Email       string
	IsNew       bool
	Action      string
}

type snapshotMetadata struct {
	FilePath          string
	FileName          string
	FileSize          int64
	Export            *chromeExport
	BookmarksCount    int
	ExtensionsCount   int
	HasEmail          bool
	Email             string
	DisplayName       string
	ProfileName       string
	ExportedAt        string
	TargetDestination importDestination
}

type bookmarkRoot struct {
	Roots map[string]bookmarkNode `json:"roots"`
}

type bookmarkNode struct {
	Type     string         `json:"type"`
	Children []bookmarkNode `json:"children"`
}

func countBookmarks(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}
	var root bookmarkRoot
	if err := json.Unmarshal(raw, &root); err != nil {
		return 0
	}
	count := 0
	for _, node := range root.Roots {
		count += countNodeBookmarks(node)
	}
	return count
}

func countNodeBookmarks(node bookmarkNode) int {
	count := 0
	if node.Type == "url" {
		count++
	}
	for _, child := range node.Children {
		count += countNodeBookmarks(child)
	}
	return count
}

func hasProfileBookmarks(profileDir string) bool {
	bmPath := filepath.Join(profileDir, "Bookmarks")
	info, err := os.Stat(bmPath)
	if err != nil {
		return false
	}
	return info.Size() > 50
}

func findNextAvailableProfileDir() string {
	stateSet := chromeLocalStateProfileDirs()
	root := chromeUserDataDir()
	for n := 1; n < 1000; n++ {
		candidate := fmt.Sprintf("Profile %d", n)
		if stateSet[candidate] {
			continue
		}
		if chromeProfilePathExists(filepath.Join(root, candidate)) {
			continue
		}
		return candidate
	}
	return fmt.Sprintf("Profile %d", time.Now().Unix())
}

func resolveImportDestination(exp *chromeExport, explicitTarget string, logSteps bool) importDestination {
	if explicitTarget != "" {
		return resolveExplicitDestination(explicitTarget, exp, logSteps)
	}
	if exp.Email != "" {
		return resolveDestinationByEmail(exp, logSteps)
	}
	return resolveDestinationDefault(exp, logSteps)
}

func resolveExplicitDestination(explicitTarget string, exp *chromeExport, logSteps bool) importDestination {
	path := chromeProfilePath(explicitTarget)
	isNew := !chromeProfilePathExists(path)
	logExplicitTargetStep(explicitTarget, logSteps)
	return importDestination{
		Dir:         explicitTarget,
		Path:        path,
		DisplayName: exp.DisplayName,
		Email:       exp.Email,
		IsNew:       isNew,
		Action:      fmt.Sprintf("Explicit target %q", explicitTarget),
	}
}

func logExplicitTargetStep(target string, logSteps bool) {
	if !logSteps {
		return
	}
	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Using explicit profile target: %q\n", target)
}

func resolveDestinationByEmail(exp *chromeExport, logSteps bool) importDestination {
	existingDir, found := findChromeProfileByEmail(exp.Email)
	if found {
		logExistingProfileMatched(existingDir, exp.Email, logSteps)
		return importDestination{
			Dir:         existingDir,
			Path:        chromeProfilePath(existingDir),
			DisplayName: exp.DisplayName,
			Email:       exp.Email,
			IsNew:       false,
			Action:      fmt.Sprintf("Matched existing profile %q (<%s>)", existingDir, exp.Email),
		}
	}

	nextDir := findNextAvailableProfileDir()
	dispName := resolveProfileDisplayName(exp)
	logNewProfileCreated(exp.Email, nextDir, logSteps)
	return importDestination{
		Dir:         nextDir,
		Path:        chromeProfilePath(nextDir),
		DisplayName: dispName,
		Email:       exp.Email,
		IsNew:       true,
		Action:      fmt.Sprintf("New profile %q (<%s>)", nextDir, exp.Email),
	}
}

func logExistingProfileMatched(existingDir, email string, logSteps bool) {
	if !logSteps {
		return
	}
	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Matched existing Chrome profile %q by email <%s>\n", existingDir, email)
}

func logNewProfileCreated(email, nextDir string, logSteps bool) {
	if !logSteps {
		return
	}
	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Email <%s> not found in Chrome. Creating new profile %q (protecting existing profiles)\n", email, nextDir)
}

func resolveProfileDisplayName(exp *chromeExport) string {
	if exp.DisplayName != "" {
		return exp.DisplayName
	}
	if exp.Name != "" {
		return exp.Name
	}
	if exp.Email != "" {
		return strings.Split(exp.Email, "@")[0]
	}
	return "Profile"
}

func findChromeProfileByEmail(email string) (string, bool) {
	for _, dir := range availableChromeProfileNames() {
		_, existingEmail := resolveProfileNameAndEmail(dir, nil)
		if existingEmail != "" && strings.EqualFold(existingEmail, email) {
			return dir, true
		}
	}
	return "", false
}

func resolveDestinationDefault(exp *chromeExport, logSteps bool) importDestination {
	candidate := exp.Name
	if candidate == "" {
		candidate = "Default"
	}
	targetPath := chromeProfilePath(candidate)
	if !chromeProfilePathExists(targetPath) {
		logProfileAvailable(candidate, logSteps)
		return importDestination{
			Dir:         candidate,
			Path:        targetPath,
			DisplayName: exp.DisplayName,
			Email:       "",
			IsNew:       true,
			Action:      fmt.Sprintf("Available profile %q", candidate),
		}
	}

	_, existingEmail := resolveProfileNameAndEmail(candidate, nil)
	isOccupied := existingEmail != "" || hasProfileBookmarks(targetPath)
	if isOccupied {
		nextDir := findNextAvailableProfileDir()
		logProfileOccupied(candidate, nextDir, logSteps)
		return importDestination{
			Dir:         nextDir,
			Path:        chromeProfilePath(nextDir),
			DisplayName: exp.DisplayName,
			Email:       "",
			IsNew:       true,
			Action:      fmt.Sprintf("Allocating new profile %q (protecting occupied %q)", nextDir, candidate),
		}
	}

	logProfileUpdating(candidate, logSteps)
	return importDestination{
		Dir:         candidate,
		Path:        targetPath,
		DisplayName: exp.DisplayName,
		Email:       "",
		IsNew:       false,
		Action:      fmt.Sprintf("Updating existing profile %q", candidate),
	}
}

func logProfileAvailable(name string, logSteps bool) {
	if !logSteps {
		return
	}
	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Target profile directory %q is available\n", name)
}

func logProfileOccupied(candidate, nextDir string, logSteps bool) {
	if !logSteps {
		return
	}
	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Profile %q is occupied; allocating new profile %q to protect existing data\n", candidate, nextDir)
}

func logProfileUpdating(candidate string, logSteps bool) {
	if !logSteps {
		return
	}
	fmt.Printf("      \033[1;94m[Step 2/5]\033[0m Updating profile %q\n", candidate)
}

func isExcepted(exceptRules []string, fileName, profileName, displayName, email string) (bool, string) {
	for _, rule := range exceptRules {
		r := strings.TrimSpace(rule)
		if r == "" {
			continue
		}
		if matchExceptRule(r, fileName) {
			return true, r
		}
		if matchExceptRule(r, profileName) {
			return true, r
		}
		if matchExceptRule(r, displayName) {
			return true, r
		}
		if matchExceptRule(r, email) {
			return true, r
		}
	}
	return false, ""
}

func matchExceptRule(rule, val string) bool {
	if val == "" {
		return false
	}
	lowRule := strings.ToLower(rule)
	lowVal := strings.ToLower(val)
	if strings.EqualFold(lowRule, lowVal) {
		return true
	}
	if strings.HasSuffix(lowRule, "*") {
		prefix := strings.TrimSuffix(lowRule, "*")
		return strings.HasPrefix(lowVal, prefix)
	}
	return strings.HasPrefix(lowVal, lowRule)
}

func readSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	lower := strings.ToLower(srcFile)
	if strings.HasSuffix(lower, constants.ExtJSON) {
		return readJSONSnapshotMetadata(srcFile)
	}
	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return readYAMLSnapshotMetadata(srcFile)
	}
	return readGenericSnapshotMetadata(srcFile)
}

func readJSONSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	raw, err := os.ReadFile(srcFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", srcFile, err)
	}
	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", srcFile, err)
	}
	extractEmailIfMissing(&exp)
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}
	bmsCount := countBookmarks(exp.Bookmarks)
	extsCount := len(exp.ExtensionIDs)
	target := resolveImportDestination(&exp, "", false)
	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          filepath.Base(srcFile),
		FileSize:          size,
		Export:            &exp,
		BookmarksCount:    bmsCount,
		ExtensionsCount:   extsCount,
		HasEmail:          exp.Email != "",
		Email:             exp.Email,
		DisplayName:       exp.DisplayName,
		ProfileName:       exp.Name,
		ExportedAt:        exp.ExportedAt,
		TargetDestination: target,
	}, nil
}

func extractEmailIfMissing(exp *chromeExport) {
	if exp.Email == "" && len(exp.Preferences) > 0 {
		exp.Email = extractEmailFromPreferences(exp.Preferences)
	}
}

func readYAMLSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	raw, err := os.ReadFile(srcFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", srcFile, err)
	}
	var exp chromeExport
	if err := yaml.Unmarshal(raw, &exp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", srcFile, err)
	}
	extractEmailIfMissing(&exp)
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}
	target := resolveImportDestination(&exp, "", false)
	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          filepath.Base(srcFile),
		FileSize:          size,
		Export:            &exp,
		BookmarksCount:    countBookmarks(exp.Bookmarks),
		ExtensionsCount:   len(exp.ExtensionIDs),
		HasEmail:          exp.Email != "",
		Email:             exp.Email,
		DisplayName:       exp.DisplayName,
		ProfileName:       exp.Name,
		ExportedAt:        exp.ExportedAt,
		TargetDestination: target,
	}, nil
}

func readGenericSnapshotMetadata(srcFile string) (*snapshotMetadata, error) {
	base := filepath.Base(srcFile)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	info, _ := os.Stat(srcFile)
	var size int64
	if info != nil {
		size = info.Size()
	}
	exp := &chromeExport{Name: name, DisplayName: name}
	target := resolveImportDestination(exp, "", false)
	return &snapshotMetadata{
		FilePath:          srcFile,
		FileName:          base,
		FileSize:          size,
		Export:            exp,
		BookmarksCount:    0,
		ExtensionsCount:   0,
		HasEmail:          false,
		Email:             "",
		DisplayName:       name,
		ProfileName:       name,
		ExportedAt:        "",
		TargetDestination: target,
	}, nil
}

func scanSnapshotFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if isImportableSnapshot(path) && isValidChromeSnapshotFile(path) {
			files = append(files, path)
		}
	}
	sort.Strings(files)
	return files
}

func isValidChromeSnapshotFile(path string) bool {
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, constants.ExtJSON) {
		return isChromeSnapshotJSON(path)
	}
	if strings.HasSuffix(lower, constants.ExtZIP) {
		return true
	}
	if isSQLiteSnapshot(lower) {
		return true
	}
	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return isChromeSnapshotYAML(path)
	}
	return false
}

func isChromeSnapshotJSON(path string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err == nil && isExportPopulated(&exp) {
		return true
	}
	var all chromeAllProfilesExport
	if err := json.Unmarshal(raw, &all); err == nil && len(all.Profiles) > 0 {
		return true
	}
	return false
}

func isExportPopulated(exp *chromeExport) bool {
	return exp.SchemaVersion > 0 ||
		exp.ExportedAt != "" ||
		len(exp.Bookmarks) > 0 ||
		len(exp.Preferences) > 0 ||
		len(exp.ExtensionIDs) > 0
}

func isChromeSnapshotYAML(path string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var exp chromeExport
	if err := yaml.Unmarshal(raw, &exp); err == nil {
		return exp.SchemaVersion > 0 || exp.ExportedAt != ""
	}
	return false
}

func registerImportedProfileToLocalState(dstDir, displayName, email string) error {
	path := filepath.Join(chromeUserDataDir(), constants.ChromeLocalStateFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	profile := ensureChromeLocalStateProfile(root)
	infoCache := ensureChromeLocalStateInfoCache(profile)

	entry, ok := infoCache[dstDir].(map[string]any)
	if !ok {
		entry = map[string]any{}
	}
	assignEntryDisplayName(entry, dstDir, displayName)
	assignEntryEmail(entry, email)
	entry["is_using_default_name"] = false
	entry["is_ephemeral"] = false
	infoCache[dstDir] = entry

	appendChromeProfileToOrder(profile, dstDir)
	return writeChromeLocalState(path, root)
}

func assignEntryDisplayName(entry map[string]any, dstDir, displayName string) {
	if displayName != "" {
		entry["name"] = displayName
		return
	}
	if entry["name"] == nil || entry["name"] == "" {
		entry["name"] = dstDir
	}
}

func assignEntryEmail(entry map[string]any, email string) {
	if email != "" {
		entry["user_name"] = email
	}
}

func importSingleSnapshotWithStepLogging(srcFile, explicitTarget string) error {
	meta, err := readSnapshotMetadata(srcFile)
	if err != nil {
		return err
	}
	exp := meta.Export
	fmt.Printf("  \033[1;94m[Step 1/5]\033[0m Inspecting snapshot: %s\n", srcFile)
	fmt.Printf("        → Name: %q | Display: %q | Email: %q | Bookmarks: %d | Extensions: %d\n",
		exp.Name, exp.DisplayName, exp.Email, meta.BookmarksCount, meta.ExtensionsCount)

	dest := resolveImportDestination(exp, explicitTarget, true)
	if err := os.MkdirAll(dest.Path, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", dest.Path, err)
	}

	fmt.Printf("  \033[1;94m[Step 3/5]\033[0m Restoring Bookmarks & Preferences into %q...\n", dest.Dir)
	if err := writeOptional(filepath.Join(dest.Path, "Bookmarks"), exp.Bookmarks); err != nil {
		return err
	}
	if err := writeOptional(filepath.Join(dest.Path, "Preferences"), exp.Preferences); err != nil {
		return err
	}

	fmt.Printf("  \033[1;94m[Step 4/5]\033[0m Staging extensions (%d pending hints)...\n", meta.ExtensionsCount)
	if err := writePendingExtensions(dest.Path, exp.ExtensionIDs); err != nil {
		return err
	}

	fmt.Printf("  \033[1;94m[Step 5/5]\033[0m Registering %q in Chrome Local State...\n", dest.Dir)
	if err := registerImportedProfileToLocalState(dest.Dir, dest.DisplayName, dest.Email); err != nil {
		fmt.Fprintf(os.Stderr, "        \033[1;93m⚠\033[0m Warning: Local State registration notice: %v\n", err)
	}

	fmt.Printf("  \033[1;92m✓ Successfully imported\033[0m %s → %s (%q)\n\n", srcFile, dest.Dir, dest.DisplayName)
	return nil
}

func runChromeProfileImportCheck(args []string) error {
	target := "."
	if len(args) > 0 && args[0] != "ls" && args[0] != "--help" && args[0] != "-h" {
		target = args[0]
	}
	checkHelp("import-check", args)

	files, err := resolveSnapshotCheckFiles(target)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Printf("No snapshot files found to check in %q\n", target)
		return nil
	}

	fmt.Printf("\n\033[1;96mChrome Snapshot Inspection & Import Preview (%s)\033[0m\n\n", target)
	for i, f := range files {
		meta, err := readSnapshotMetadata(f)
		if err != nil {
			fmt.Printf("  [%d] %s \033[1;91m(Error reading: %v)\033[0m\n", i+1, f, err)
			continue
		}
		printSnapshotCheckItem(i+1, meta)
	}

	fmt.Printf("Total snapshots inspected: %d\n", len(files))
	fmt.Printf("Usage hints:\n")
	fmt.Printf("  Import all:            gitmap chrome profile import *.json\n")
	fmt.Printf("  Import single email:   gitmap chrome profile import <email>\n")
	fmt.Printf("  Import with limit:     gitmap chrome profile import *.json --limit 1\n")
	fmt.Printf("  Import with exclusion: gitmap chrome profile import *.json --except <id|name>\n\n")
	return nil
}

func resolveSnapshotCheckFiles(target string) ([]string, error) {
	info, err := os.Stat(target)
	if err == nil && !info.IsDir() {
		return []string{target}, nil
	}
	files := scanSnapshotFiles(target)
	files = fallbackCheckFiles(files, target)
	return files, nil
}

func fallbackCheckFiles(files []string, target string) []string {
	if len(files) > 0 || target != "." {
		return files
	}
	gitmapChromeDir := filepath.Join(constants.GitMapDir, "chrome")
	if !chromeProfilePathExists(gitmapChromeDir) {
		return files
	}
	return scanSnapshotFiles(gitmapChromeDir)
}

func printSnapshotCheckItem(idx int, meta *snapshotMetadata) {
	fmt.Printf("  [%d] \033[1;97m%s\033[0m (%.1f KB)\n", idx, meta.FileName, float64(meta.FileSize)/1024.0)
	emailStr := meta.Email
	if emailStr == "" {
		emailStr = "(none)"
	}
	dispStr := meta.DisplayName
	if dispStr == "" {
		dispStr = "(none)"
	}
	fmt.Printf("      • Profile: %-12s | Display: %-14s | Email: %s\n", meta.ProfileName, dispStr, emailStr)
	fmt.Printf("      • Bookmarks: %-10d | Extensions: %-11d | Exported: %s\n",
		meta.BookmarksCount, meta.ExtensionsCount, meta.ExportedAt)

	printSnapshotAction(meta)
	fmt.Println()
}

func printSnapshotAction(meta *snapshotMetadata) {
	if meta.TargetDestination.IsNew {
		fmt.Printf("      \033[1;92m+ Action: Will CREATE NEW profile %q\033[0m (%s)\n",
			meta.TargetDestination.Dir, meta.TargetDestination.Action)
		return
	}
	fmt.Printf("      \033[1;93m↷ Action: Will UPDATE profile %q\033[0m (%s)\n",
		meta.TargetDestination.Dir, meta.TargetDestination.Action)
}

func listDiscoveredSnapshotsInDir(dir string) {
	files := scanSnapshotFiles(dir)
	files, dir = resolveFallbackSnapshots(files, dir)
	if len(files) == 0 {
		return
	}

	fmt.Printf("\n\033[1;96mDiscovered Backup / Snapshot Files in %s:\033[0m\n", dir)
	for i, f := range files {
		meta, err := readSnapshotMetadata(f)
		if err != nil {
			fmt.Printf("  [%d] %s\n", i+1, filepath.Base(f))
			continue
		}
		emailStr := meta.Email
		if emailStr == "" {
			emailStr = "(none)"
		}
		fmt.Printf("  [%d] \033[1;97m%-16s\033[0m (Display: %q, Email: %s, Bookmarks: %d, Ext: %d)\n",
			i+1, meta.FileName, meta.DisplayName, emailStr, meta.BookmarksCount, meta.ExtensionsCount)
	}
	fmt.Println()
}

func resolveFallbackSnapshots(files []string, dir string) ([]string, string) {
	if len(files) > 0 {
		return files, dir
	}
	gitmapChromeDir := filepath.Join(constants.GitMapDir, "chrome")
	if dir == gitmapChromeDir || !chromeProfilePathExists(gitmapChromeDir) {
		return files, dir
	}
	return scanSnapshotFiles(gitmapChromeDir), gitmapChromeDir
}

func runSmartChromeImport(opts chromeTransferOptions) error {
	candidates, explicitTarget, err := collectImportCandidates(opts)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageImport)
		return fmt.Errorf("no matching snapshot files found to import")
	}

	filtered := filterImportCandidates(candidates, opts)
	if len(filtered) == 0 {
		fmt.Println("No snapshot files remained after applying filters (--except / --email).")
		return nil
	}

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	fmt.Printf("\n\033[1;96mStarting Chrome Profile Import (%d snapshot(s))\033[0m\n\n", len(filtered))
	successCount := 0
	for _, f := range filtered {
		if err := importSingleSnapshotWithStepLogging(f, explicitTarget); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[1;91m✗ Failed to import %s:\033[0m %v\n", f, err)
			continue
		}
		successCount++
	}

	fmt.Printf("\033[1;92m✓ Chrome Profile Import Complete:\033[0m %d of %d profile(s) imported successfully.\n\n",
		successCount, len(filtered))
	return nil
}

func collectImportCandidates(opts chromeTransferOptions) ([]string, string, error) {
	if opts.Email != "" {
		return findCandidatesByEmail(opts.Email)
	}
	if len(opts.Positional) == 0 {
		return findCandidatesInDir(".")
	}
	if len(opts.Positional) == 1 {
		return resolveSinglePositionalCandidate(opts.Positional[0])
	}
	return resolveMultiplePositionalCandidates(opts.Positional)
}

func findCandidatesByEmail(email string) ([]string, string, error) {
	files := collectSearchSnapshotFiles(".")
	for _, f := range files {
		meta, err := readSnapshotMetadata(f)
		if err == nil && meta.Email != "" && strings.EqualFold(meta.Email, email) {
			return []string{f}, "", nil
		}
	}
	return nil, "", fmt.Errorf("no snapshot file found containing email %q (searched %d files in current directory)", email, len(files))
}

func collectSearchSnapshotFiles(dir string) []string {
	files := scanSnapshotFiles(dir)
	gitmapChromeDir := filepath.Join(constants.GitMapDir, "chrome")
	if chromeProfilePathExists(gitmapChromeDir) {
		files = append(files, scanSnapshotFiles(gitmapChromeDir)...)
	}
	return files
}

func findCandidatesInDir(dir string) ([]string, string, error) {
	files := scanSnapshotFiles(dir)
	files = fallbackCheckFiles(files, dir)
	return files, "", nil
}

func resolveSinglePositionalCandidate(pos0 string) ([]string, string, error) {
	if pos0 == "." || isDirectoryPath(pos0) {
		return findCandidatesInDir(pos0)
	}
	if strings.ContainsAny(pos0, "*?[") {
		return resolveGlobCandidate(pos0)
	}
	if strings.Contains(pos0, "@") && !isSnapshotFileExtension(pos0) {
		return findCandidatesByEmail(pos0)
	}
	if chromeProfilePathExists(pos0) {
		return []string{pos0}, "", nil
	}
	gitmapFile := filepath.Join(constants.GitMapDir, "chrome", pos0+".json")
	if chromeProfilePathExists(gitmapFile) {
		return []string{gitmapFile}, "", nil
	}
	return nil, "", fmt.Errorf("snapshot file or directory %q not found", pos0)
}

func resolveGlobCandidate(pattern string) ([]string, string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, "", err
	}
	return matches, "", nil
}

func resolveMultiplePositionalCandidates(pos []string) ([]string, string, error) {
	if len(pos) == 2 && isImportableSnapshot(pos[0]) && !isImportableSnapshot(pos[1]) {
		return []string{pos[0]}, pos[1], nil
	}
	var out []string
	for _, p := range pos {
		appendCandidateMatches(p, &out)
	}
	return out, "", nil
}

func appendCandidateMatches(p string, out *[]string) {
	if strings.ContainsAny(p, "*?[") {
		matches, _ := filepath.Glob(p)
		*out = append(*out, matches...)
		return
	}
	if chromeProfilePathExists(p) {
		*out = append(*out, p)
	}
}

func filterImportCandidates(candidates []string, opts chromeTransferOptions) []string {
	var out []string
	for _, f := range candidates {
		meta, err := readSnapshotMetadata(f)
		if err != nil {
			out = append(out, f)
			continue
		}
		if opts.Email != "" && !strings.EqualFold(meta.Email, opts.Email) {
			continue
		}
		if isEx, rule := isExcepted(opts.Except, meta.FileName, meta.ProfileName, meta.DisplayName, meta.Email); isEx {
			fmt.Printf("  \033[1;93m↷ Skipping\033[0m %s (matched --except %q)\n", f, rule)
			continue
		}
		out = append(out, f)
	}
	return out
}
