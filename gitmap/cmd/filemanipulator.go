package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type lowercaseOptions struct {
	Source        string
	Target        string
	Except        []string
	IgnoreDefault bool
}

func runLowercase(args []string) error {
	opts, dirs := parseLowercaseArgs(args)
	if opts.Source == "" {
		printLowercaseUsage()
		return nil
	}
	opts.Except = appendDefaultIgnores(opts)
	dirs = ensureDefaultDir(dirs)
	for _, d := range dirs {
		processLowercaseDir(d, opts)
	}
	return nil
}

func parseLowercaseArgs(args []string) (lowercaseOptions, []string) {
	opts := lowercaseOptions{IgnoreDefault: true}
	var dirs []string
	for i := 0; i < len(args); i++ {
		i = handleLowercaseArg(args, i, &opts, &dirs)
	}
	return opts, dirs
}

func handleLowercaseArg(args []string, i int, opts *lowercaseOptions, dirs *[]string) int {
	arg := args[i]
	if arg == "-except" && i+1 < len(args) {
		opts.Except = parseExceptList(args[i+1])
		return i + 1
	}
	if arg == "-ignore" && i+1 < len(args) {
		opts.IgnoreDefault = (args[i+1] == "default")
		return i + 1
	}
	if !strings.HasPrefix(arg, "-") {
		assignPositional(arg, opts, dirs)
	}
	return i
}

func parseExceptList(val string) []string {
	parts := strings.Split(val, ",")
	for j := range parts {
		parts[j] = strings.TrimSpace(parts[j])
	}
	return parts
}

func assignPositional(arg string, opts *lowercaseOptions, dirs *[]string) {
	switch {
	case opts.Source == "":
		opts.Source = arg
	case opts.Target == "":
		opts.Target = arg
	default:
		*dirs = append(*dirs, arg)
	}
}

func printLowercaseUsage() {
	fmt.Println("Usage: gitmap lowercase <source_pattern> <target_pattern> [dirs...] -except \"paths\"")
	fmt.Println("Example: gitmap lowercase \"OLD.md\" \"old.md\" . -except \"node_modules/*,.git/*\"")
}

func appendDefaultIgnores(opts lowercaseOptions) []string {
	if opts.IgnoreDefault {
		return append(opts.Except, "node_modules/*", ".git/*")
	}
	return opts.Except
}

func ensureDefaultDir(dirs []string) []string {
	if len(dirs) == 0 {
		return []string{"."}
	}
	return dirs
}

func isLowercaseIgnored(path string, except []string) bool {
	for _, pattern := range except {
		if hasMatch(pattern, path) {
			return true
		}
	}
	return false
}

func hasMatch(pattern, path string) bool {
	matched, _ := filepath.Match(pattern, filepath.Base(path))
	if matched {
		return true
	}
	matched, _ = filepath.Match(pattern, path)
	return matched
}

func runGitMv(src, dst string) bool {
	cmd := exec.Command("git", "mv", src, dst)
	cmd.Dir = filepath.Dir(src)
	err := cmd.Run()
	return err == nil
}

func processLowercaseDir(dir string, opts lowercaseOptions) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return
	}
	absDir = normalizeWindowsPath(absDir)
	filepath.Walk(absDir, buildWalkFunc(opts))
}

func buildWalkFunc(opts lowercaseOptions) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if isLowercaseIgnored(path, opts.Except) {
			return skipIfDir(info)
		}
		if !info.IsDir() {
			renameIfMatch(path, opts)
		}
		return nil
	}
}

func skipIfDir(info os.FileInfo) error {
	if info.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

func renameIfMatch(path string, opts lowercaseOptions) {
	base := filepath.Base(path)
	if strings.Contains(base, opts.Source) {
		newBase := strings.ReplaceAll(base, opts.Source, opts.Target)
		newPath := filepath.Join(filepath.Dir(path), newBase)
		executeRename(path, newPath, base, newBase)
	}
}

func executeRename(path, newPath, base, newBase string) {
	if runGitMv(path, newPath) {
		fmt.Printf("✅ git mv: %s -> %s\n", base, newBase)
		return
	}
	if err := os.Rename(path, newPath); err == nil {
		fmt.Printf("✅ os.rename: %s -> %s\n", base, newBase)
	} else {
		fmt.Printf("❌ rename failed: %v\n", err)
	}
}

type seqFile struct {
	OriginalPath string
	Dir          string
	BaseName     string
	Seq          int
	HasSeq       bool
	Rest         string
	Time         int64
	NewSeq       int
}

type fixSeqOpts struct {
	IsOrderByTime  bool
	IsOrderByAZ    bool
	IsKeepOldOrder bool
	PinMap         map[string]int
}

func runFixSeqFiles(args []string) error {
	opts, dirs := parseFixSeqArgs(args)
	if len(dirs) == 0 {
		printFixSeqUsage()
		return nil
	}
	for _, d := range dirs {
		processFixSeqDir(d, opts)
	}
	return nil
}

func parseFixSeqArgs(args []string) (fixSeqOpts, []string) {
	opts := fixSeqOpts{PinMap: make(map[string]int)}
	var dirs []string
	for i := 0; i < len(args); i++ {
		i = handleFixSeqArg(args, i, &opts, &dirs)
	}
	return opts, dirs
}

func handleFixSeqArg(args []string, i int, opts *fixSeqOpts, dirs *[]string) int {
	arg := args[i]
	switch {
	case arg == "-orderbytime":
		opts.IsOrderByTime = true
	case arg == "-orderbyaz":
		opts.IsOrderByAZ = true
	case arg == "-keep-old-order":
		opts.IsKeepOldOrder = true
	case arg == "-pin" && i+1 < len(args):
		parsePinMap(args[i+1], opts.PinMap)
		return i + 1
	case !strings.HasPrefix(arg, "-"):
		*dirs = append(*dirs, arg)
	}
	return i
}

func parsePinMap(pinStr string, pinMap map[string]int) {
	pairs := strings.Split(pinStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			pinMap[strings.ToLower(strings.TrimSpace(parts[0]))] = val
		}
	}
}

func printFixSeqUsage() {
	fmt.Println("Usage: gitmap fix-seq-files <folder1> [folder2] [-orderbytime] [-orderbyaz] [-keep-old-order] [-pin \"draft=01\"]")
}

func processFixSeqDir(dir string, opts fixSeqOpts) {
	absDir, err := resolveNormPath(dir)
	if err != nil {
		return
	}
	entries, err := os.ReadDir(absDir)
	if err != nil {
		fmt.Printf("❌ Failed to read dir %s: %v\n", absDir, err)
		return
	}
	parsedFiles := parseSeqFiles(entries, absDir)
	sortParsedFiles(parsedFiles, opts)
	applyNewSeqs(parsedFiles, opts)
	renameSeqFiles(parsedFiles)
}

func resolveNormPath(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return normalizeWindowsPath(absDir), nil
}

func parseSeqFiles(entries []os.DirEntry, absDir string) []*seqFile {
	var parsedFiles []*seqFile
	re := regexp.MustCompile(`^(\d+)[-_](.*)$`)
	for _, entry := range entries {
		if sf := buildSeqFile(entry, absDir, re); sf != nil {
			parsedFiles = append(parsedFiles, sf)
		}
	}
	return parsedFiles
}

func buildSeqFile(entry os.DirEntry, absDir string, re *regexp.Regexp) *seqFile {
	if entry.IsDir() {
		return nil
	}
	info, err := entry.Info()
	if err != nil {
		return nil
	}
	sf := &seqFile{
		OriginalPath: filepath.Join(absDir, entry.Name()),
		Dir:          absDir,
		BaseName:     entry.Name(),
		Time:         info.ModTime().UnixNano(),
	}
	applyRegexExtract(sf, entry.Name(), re)
	return sf
}

func applyRegexExtract(sf *seqFile, name string, re *regexp.Regexp) {
	match := re.FindStringSubmatch(name)
	if len(match) == 3 {
		sf.HasSeq = true
		sf.Seq, _ = strconv.Atoi(match[1])
		sf.Rest = match[2]
	} else {
		sf.Rest = name
	}
}

func sortParsedFiles(files []*seqFile, opts fixSeqOpts) {
	if opts.IsOrderByTime {
		sort.Slice(files, func(i, j int) bool { return files[i].Time < files[j].Time })
	} else if opts.IsOrderByAZ {
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Rest) < strings.ToLower(files[j].Rest)
		})
	}
}

func applyNewSeqs(files []*seqFile, opts fixSeqOpts) {
	pinned, unpinned := partitionFiles(files, opts.PinMap)
	if opts.IsKeepOldOrder {
		sortKeepOldOrder(unpinned)
	}
	usedSeqs := buildUsedSeqs(pinned)
	assignUnpinnedSeqs(unpinned, usedSeqs)
}

func partitionFiles(files []*seqFile, pinMap map[string]int) ([]*seqFile, []*seqFile) {
	var pinned, unpinned []*seqFile
	for _, pf := range files {
		baseLower := strings.ToLower(strings.TrimSuffix(pf.Rest, filepath.Ext(pf.Rest)))
		if val, exists := pinMap[baseLower]; exists {
			pf.NewSeq = val
			pinned = append(pinned, pf)
		} else {
			unpinned = append(unpinned, pf)
		}
	}
	return pinned, unpinned
}

func sortKeepOldOrder(unpinned []*seqFile) {
	sort.Slice(unpinned, func(i, j int) bool {
		if unpinned[i].HasSeq && unpinned[j].HasSeq {
			return unpinned[i].Seq < unpinned[j].Seq
		}
		return unpinned[i].HasSeq
	})
}

func buildUsedSeqs(pinned []*seqFile) map[int]bool {
	used := make(map[int]bool)
	for _, pf := range pinned {
		used[pf.NewSeq] = true
	}
	return used
}

func assignUnpinnedSeqs(unpinned []*seqFile, usedSeqs map[int]bool) {
	currentSeq := 0
	for _, pf := range unpinned {
		for usedSeqs[currentSeq] {
			currentSeq++
		}
		pf.NewSeq = currentSeq
		usedSeqs[currentSeq] = true
		currentSeq++
	}
}

func renameSeqFiles(files []*seqFile) {
	maxSeq := findMaxSeq(files)
	digits := 2
	if maxSeq > 99 {
		digits = len(strconv.Itoa(maxSeq))
	}
	for _, pf := range files {
		format := fmt.Sprintf("%%0%dd-%%s", digits)
		newName := fmt.Sprintf(format, pf.NewSeq, pf.Rest)
		if newName != pf.BaseName {
			newPath := filepath.Join(pf.Dir, newName)
			executeRename(pf.OriginalPath, newPath, pf.BaseName, newName)
		}
	}
}

func findMaxSeq(files []*seqFile) int {
	maxSeq := 0
	for _, pf := range files {
		if pf.NewSeq > maxSeq {
			maxSeq = pf.NewSeq
		}
	}
	return maxSeq
}

func normalizeWindowsPath(path string) string {
	if len(path) > 248 && !strings.HasPrefix(path, `\\?\`) {
		return `\\?\` + path
	}
	return path
}
