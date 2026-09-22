package cmdos

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

// AICleanCategory represents a discovered collection of purgeable cache files.
type AICleanCategory struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Paths       []string `json:"paths"`
	FileCount   int      `json:"file_count"`
	TotalBytes  int64    `json:"total_bytes"`
}

var (
	// OverrideAntigravityDir overrides ~/.gemini/antigravity for tests.
	OverrideAntigravityDir string

	// OverrideWorkspaceDir overrides workspace root for tests.
	OverrideWorkspaceDir string

	// OverrideTempDir overrides temporary directory for tests.
	OverrideTempDir string

	// OverrideGitmapDir overrides ~/.gitmap for tests.
	OverrideGitmapDir string
)

func resolveAntigravityDir() string {
	if OverrideAntigravityDir != "" {
		return OverrideAntigravityDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".gemini", "antigravity")
}

func resolveWorkspaceDir() string {
	if OverrideWorkspaceDir != "" {
		return OverrideWorkspaceDir
	}

	return "."
}

func resolveTempDir() string {
	if OverrideTempDir != "" {
		return OverrideTempDir
	}

	return os.TempDir()
}

func resolveGitmapDir() string {
	if OverrideGitmapDir != "" {
		return OverrideGitmapDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".gitmap")
}

func isProtectedFile(filePath string) bool {
	return osclean.IsAntigravityProtected(filePath)
}

type cleanAccumulator struct {
	paths      []string
	fileCount  int
	totalBytes int64
}

func (a *cleanAccumulator) addFilePath(filePath string, size int64) {
	a.paths = append(a.paths, filePath)
	a.fileCount++
	a.totalBytes += size
}

func scanDirectoryTree(dirPath string, acc *cleanAccumulator) {
	if !pathExists(dirPath) || osclean.IsAntigravityProtected(dirPath) {
		return
	}

	_ = filepath.Walk(dirPath, func(p string, info os.FileInfo, err error) error {
		processWalkEntry(p, info, err, acc)

		return nil
	})
}

func processWalkEntry(p string, info os.FileInfo, err error, acc *cleanAccumulator) {
	if err != nil || info == nil || info.IsDir() {
		return
	}
	if isProtectedFile(p) {
		return
	}

	acc.addFilePath(p, info.Size())
}

// DiscoverAntigravityBrainCaches finds scratch, crash logs, and temp media in Antigravity.
func DiscoverAntigravityBrainCaches() AICleanCategory {
	acc := &cleanAccumulator{}
	baseDir := resolveAntigravityDir()
	if baseDir != "" {
		scanAntigravityBrainDirs(baseDir, acc)
	}

	return AICleanCategory{
		Name:        "Antigravity Brain Caches",
		Description: "Antigravity scratch files, crash logs, and temp media",
		Paths:       acc.paths,
		FileCount:   acc.fileCount,
		TotalBytes:  acc.totalBytes,
	}
}

func scanAntigravityBrainDirs(baseDir string, acc *cleanAccumulator) {
	scanDirectoryTree(filepath.Join(baseDir, "crashes"), acc)
	scanDirectoryTree(filepath.Join(baseDir, "scratch"), acc)
	scanDirectoryTree(filepath.Join(baseDir, "brain", "tempmediaStorage"), acc)
	scanBrainSubdirScratches(filepath.Join(baseDir, "brain"), acc)
}

func scanBrainSubdirScratches(brainDir string, acc *cleanAccumulator) {
	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		processBrainSubdir(brainDir, entry, acc)
	}
}

func processBrainSubdir(brainDir string, entry os.DirEntry, acc *cleanAccumulator) {
	isDir := entry.IsDir()
	if !isDir {
		return
	}

	subName := entry.Name()
	isSpecialDir := subName == "tempmediaStorage"
	if isSpecialDir {
		return
	}

	scanDirectoryTree(filepath.Join(brainDir, subName, "scratch"), acc)
}

// DiscoverSystemGeneratedTasks finds system-generated background task logs.
func DiscoverSystemGeneratedTasks() AICleanCategory {
	acc := &cleanAccumulator{}
	scanWorkspaceTasks(acc)
	scanBrainGeneratedTasks(acc)

	return AICleanCategory{
		Name:        "System Generated Tasks",
		Description: "System-generated background execution tasks and worker logs",
		Paths:       acc.paths,
		FileCount:   acc.fileCount,
		TotalBytes:  acc.totalBytes,
	}
}

func scanWorkspaceTasks(acc *cleanAccumulator) {
	wsDir := resolveWorkspaceDir()
	tasksDir := filepath.Join(wsDir, ".system_generated", "tasks")
	scanDirectoryTree(tasksDir, acc)
}

func scanBrainGeneratedTasks(acc *cleanAccumulator) {
	baseDir := resolveAntigravityDir()
	if baseDir == "" {
		return
	}

	brainDir := filepath.Join(baseDir, "brain")
	scanBrainDirEntries(brainDir, acc)
}

func scanBrainDirEntries(brainDir string, acc *cleanAccumulator) {
	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		processBrainTaskSubdir(brainDir, entry, acc)
	}
}

func processBrainTaskSubdir(brainDir string, entry os.DirEntry, acc *cleanAccumulator) {
	isDir := entry.IsDir()
	if !isDir {
		return
	}

	taskPath := filepath.Join(brainDir, entry.Name(), ".system_generated", "tasks")
	scanDirectoryTree(taskPath, acc)
}

// DiscoverOSTempAIDumps finds temporary AI dumps, agent caches, and updater caches.
func DiscoverOSTempAIDumps() AICleanCategory {
	acc := &cleanAccumulator{}
	tempDir := resolveTempDir()
	scanTempWildcards(tempDir, acc)
	scanAdditionalTempPaths(acc)

	return AICleanCategory{
		Name:        "OS Temp AI Dumps",
		Description: "Operating system temporary AI caches and agent crash dumps",
		Paths:       acc.paths,
		FileCount:   acc.fileCount,
		TotalBytes:  acc.totalBytes,
	}
}

func scanTempWildcards(tempDir string, acc *cleanAccumulator) {
	prefixes := []string{"antigravity", "gemini", "agent"}
	for _, prefix := range prefixes {
		pattern := filepath.Join(tempDir, prefix+"*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		collectGlobMatches(matches, acc)
	}
}

func collectGlobMatches(matches []string, acc *cleanAccumulator) {
	for _, match := range matches {
		processGlobMatch(match, acc)
	}
}

func processGlobMatch(match string, acc *cleanAccumulator) {
	fi, err := os.Stat(match)
	if err != nil {
		return
	}
	if fi.IsDir() {
		scanDirectoryTree(match, acc)
		return
	}

	appendSingleGlobFile(match, fi.Size(), acc)
}

func appendSingleGlobFile(match string, size int64, acc *cleanAccumulator) {
	if isProtectedFile(match) {
		return
	}

	acc.addFilePath(match, size)
}

func scanAdditionalTempPaths(acc *cleanAccumulator) {
	scanTmpAntigravity(acc)
	scanUserCacheAntigravity(acc)
	scanUpdaterCache(acc)
}

func scanTmpAntigravity(acc *cleanAccumulator) {
	if runtime.GOOS == "windows" && OverrideTempDir == "" {
		return
	}

	matches, err := filepath.Glob("/tmp/antigravity*")
	if err != nil {
		return
	}

	collectGlobMatches(matches, acc)
}

func scanUserCacheAntigravity(acc *cleanAccumulator) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	cacheDir := filepath.Join(homeDir, ".cache", "antigravity")
	scanDirectoryTree(cacheDir, acc)
}

func scanUpdaterCache(acc *cleanAccumulator) {
	updaterDir := resolveUpdaterDir()
	if updaterDir == "" {
		return
	}

	scanDirectoryTree(updaterDir, acc)
}

func resolveUpdaterDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		return filepath.Join(localAppData, "antigravity-updater")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".local", "share", "antigravity-updater")
}

// DiscoverGitmapInstallerCaches finds GitMap downloads and temporary build caches.
func DiscoverGitmapInstallerCaches() AICleanCategory {
	acc := &cleanAccumulator{}
	scanGitmapHomeCaches(acc)
	scanGitmapTempCaches(acc)

	return AICleanCategory{
		Name:        "GitMap Installer Temp Caches",
		Description: "GitMap installer downloads, build sandboxes, and purge staging",
		Paths:       acc.paths,
		FileCount:   acc.fileCount,
		TotalBytes:  acc.totalBytes,
	}
}

func scanGitmapHomeCaches(acc *cleanAccumulator) {
	homeGitmap := resolveGitmapDir()
	if homeGitmap == "" {
		return
	}

	subDirs := []string{"downloads", "build", "sandbox", "purge"}
	for _, sub := range subDirs {
		scanDirectoryTree(filepath.Join(homeGitmap, sub), acc)
	}
}

func scanGitmapTempCaches(acc *cleanAccumulator) {
	tempDir := resolveTempDir()
	tempGitmap := filepath.Join(tempDir, "gitmap")
	subDirs := []string{"downloads", "build", "sandbox", "purge"}
	for _, sub := range subDirs {
		scanDirectoryTree(filepath.Join(tempGitmap, sub), acc)
	}
}

// DiscoverAllAICleanTargets aggregates all 4 cache categories.
func DiscoverAllAICleanTargets() []AICleanCategory {
	return []AICleanCategory{
		DiscoverAntigravityBrainCaches(),
		DiscoverSystemGeneratedTasks(),
		DiscoverOSTempAIDumps(),
		DiscoverGitmapInstallerCaches(),
	}
}
