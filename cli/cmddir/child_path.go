package cmddir

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautomation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func init() {
	cmdautomation.RegisterChildPathRunner(RunChildPath)
}

// ChildPathOptions encapsulates parsed command-line flags and parameters for child-path.
type ChildPathOptions struct {
	TargetDir  string
	Glob       string
	Depth      int
	Extensions []string
	TypeFilter string
	Limit      int
	IsJson     bool
	IsAbs      bool
}

// ChildPathItem represents a single discovered filesystem entry.
type ChildPathItem struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

// RunChildPath executes the child-path inventory and search command.
func RunChildPath(args []string) error {
	checkHelp("child-path", args)

	opts, err := parseChildPathArgs(args)
	if err != nil {
		return err
	}

	if err := validateChildPathOptions(opts); err != nil {
		return err
	}

	return executeChildPath(opts)
}

func parseChildPathArgs(args []string) (ChildPathOptions, error) {
	opts := defaultChildPathOptions()
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		consumed, err := parseSingleChildPathArg(args, i, &opts, &nonFlags)
		if err != nil {
			return opts, err
		}

		i += consumed
	}

	resolveChildPathTargets(nonFlags, &opts)

	return opts, nil
}

func defaultChildPathOptions() ChildPathOptions {
	return ChildPathOptions{
		TargetDir:  ".",
		Glob:       "*",
		Depth:      1,
		TypeFilter: "all",
		Limit:      0,
		IsJson:     false,
		IsAbs:      false,
	}
}

func parseSingleChildPathArg(args []string, idx int, opts *ChildPathOptions, nonFlags *[]string) (int, error) {
	arg := args[idx]
	hasNext := idx+1 < len(args)
	nextArg := ""

	if hasNext {
		nextArg = args[idx+1]
	}

	if arg == "--json" || arg == "-j" {
		opts.IsJson = true

		return 0, nil
	}

	if arg == "--abs" {
		opts.IsAbs = true

		return 0, nil
	}

	if arg == "--rel" {
		opts.IsAbs = false

		return 0, nil
	}

	return parseChildPathParameterizedFlags(arg, nextArg, hasNext, opts, nonFlags)
}

func parseChildPathParameterizedFlags(arg string, nextArg string, hasNext bool, opts *ChildPathOptions, nonFlags *[]string) (int, error) {
	if consumed, err := parseChildPathDepthFlag(arg, nextArg, hasNext, opts); consumed > 0 || err != nil || isChildPathDepthFlag(arg) {
		return consumed, err
	}

	if consumed := parseChildPathExtFlag(arg, nextArg, hasNext, opts); consumed > 0 || isChildPathExtFlag(arg) {
		return consumed, nil
	}

	if consumed := parseChildPathTypeFlag(arg, nextArg, hasNext, opts); consumed > 0 || isChildPathTypeFlag(arg) {
		return consumed, nil
	}

	if consumed, err := parseChildPathLimitFlag(arg, nextArg, hasNext, opts); consumed > 0 || err != nil || isChildPathLimitFlag(arg) {
		return consumed, err
	}

	if !strings.HasPrefix(arg, "-") {
		*nonFlags = append(*nonFlags, arg)
	}

	return 0, nil
}

func isChildPathDepthFlag(arg string) bool {
	return strings.HasPrefix(arg, "--depth=") || strings.HasPrefix(arg, "-d=")
}

func isChildPathExtFlag(arg string) bool {
	return strings.HasPrefix(arg, "--ext=") || strings.HasPrefix(arg, "-e=")
}

func isChildPathTypeFlag(arg string) bool {
	return strings.HasPrefix(arg, "--type=") || strings.HasPrefix(arg, "-t=")
}

func isChildPathLimitFlag(arg string) bool {
	return strings.HasPrefix(arg, "--limit=") || strings.HasPrefix(arg, "-l=") || strings.HasPrefix(arg, "-n=")
}

func extractChildPathDepthRaw(arg string, nextArg string, hasNext bool) (string, int, bool) {
	if (arg == "--depth" || arg == "-d") && hasNext {
		return nextArg, 1, true
	}

	if isChildPathDepthFlag(arg) {
		parts := strings.SplitN(arg, "=", 2)

		return parts[1], 0, true
	}

	return "", 0, false
}

func parseChildPathDepthFlag(arg string, nextArg string, hasNext bool, opts *ChildPathOptions) (int, error) {
	raw, consumed, matched := extractChildPathDepthRaw(arg, nextArg, hasNext)
	if !matched {
		return 0, nil
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, apperror.NewValidation("child-path", "E1084", "depth must be an integer")
	}

	opts.Depth = val

	return consumed, nil
}

func parseChildPathExtFlag(arg string, nextArg string, hasNext bool, opts *ChildPathOptions) int {
	if (arg == "--ext" || arg == "-e") && hasNext {
		opts.Extensions = append(opts.Extensions, splitChildPathExtensions(nextArg)...)

		return 1
	}

	if isChildPathExtFlag(arg) {
		parts := strings.SplitN(arg, "=", 2)
		opts.Extensions = append(opts.Extensions, splitChildPathExtensions(parts[1])...)

		return 0
	}

	return 0
}

func parseChildPathTypeFlag(arg string, nextArg string, hasNext bool, opts *ChildPathOptions) int {
	if (arg == "--type" || arg == "-t") && hasNext {
		opts.TypeFilter = strings.ToLower(nextArg)

		return 1
	}

	if isChildPathTypeFlag(arg) {
		parts := strings.SplitN(arg, "=", 2)
		opts.TypeFilter = strings.ToLower(parts[1])

		return 0
	}

	return 0
}

func extractChildPathLimitRaw(arg string, nextArg string, hasNext bool) (string, int, bool) {
	if (arg == "--limit" || arg == "-l" || arg == "-n") && hasNext {
		return nextArg, 1, true
	}

	if isChildPathLimitFlag(arg) {
		parts := strings.SplitN(arg, "=", 2)

		return parts[1], 0, true
	}

	return "", 0, false
}

func parseChildPathLimitFlag(arg string, nextArg string, hasNext bool, opts *ChildPathOptions) (int, error) {
	raw, consumed, matched := extractChildPathLimitRaw(arg, nextArg, hasNext)
	if !matched {
		return 0, nil
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, apperror.NewValidation("child-path", "E1000", "limit must be an integer")
	}

	opts.Limit = val

	return consumed, nil
}

func splitChildPathExtensions(val string) []string {
	parts := strings.Split(val, ",")
	res := make([]string, 0, len(parts))

	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if clean != "" {
			res = append(res, clean)
		}
	}

	return res
}

func resolveChildPathTargets(nonFlags []string, opts *ChildPathOptions) {
	if len(nonFlags) == 0 {
		return
	}

	if len(nonFlags) == 1 {
		resolveSingleChildPathTarget(nonFlags[0], opts)

		return
	}

	opts.TargetDir = nonFlags[0]
	opts.Glob = nonFlags[1]
}

func resolveSingleChildPathTarget(arg string, opts *ChildPathOptions) {
	if isChildPathGlobPattern(arg) {
		opts.Glob = arg

		return
	}

	if info, err := os.Stat(arg); err == nil && info.IsDir() {
		opts.TargetDir = arg

		return
	}

	if strings.Contains(arg, ".") {
		opts.Glob = arg

		return
	}

	opts.TargetDir = arg
}

func isChildPathGlobPattern(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[]")
}

func validateChildPathOptions(opts ChildPathOptions) error {
	if opts.Depth < 0 {
		return apperror.NewValidation("child-path", "E1084", fmt.Sprintf("depth %d is out of bounds (must be >= 0)", opts.Depth))
	}

	if opts.TypeFilter != "f" && opts.TypeFilter != "d" && opts.TypeFilter != "all" {
		return apperror.NewValidation("child-path", "E1082", fmt.Sprintf("invalid type filter %q (must be 'f', 'd', or 'all')", opts.TypeFilter))
	}

	if _, err := filepath.Match(opts.Glob, ""); errors.Is(err, filepath.ErrBadPattern) {
		return apperror.NewValidation("child-path", "E1083", fmt.Sprintf("invalid glob pattern %q", opts.Glob))
	}

	return validateChildPathTargetDir(opts.TargetDir)
}

func validateChildPathTargetDir(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return apperror.NewNotFound("child-path", "E1081", fmt.Sprintf("target directory %q does not exist", dir))
	}

	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to access directory %q", dir))
	}

	if !info.IsDir() {
		return apperror.NewValidation("child-path", "E1081", fmt.Sprintf("target path %q is not a directory", dir))
	}

	return nil
}

func executeChildPath(opts ChildPathOptions) error {
	var items []ChildPathItem
	var err error

	if opts.Depth == 1 {
		items, err = scanChildPathDepthOne(opts)
	} else {
		items, err = scanChildPathDepthWalk(opts)
	}

	if err != nil {
		return apperror.WrapSimple(err, "child-path scan failed")
	}

	return outputChildPathResults(items, opts.IsJson)
}

func scanChildPathDepthOne(opts ChildPathOptions) ([]ChildPathItem, error) {
	entries, err := os.ReadDir(opts.TargetDir)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read directory failed")
	}

	items := make([]ChildPathItem, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() && isChildPathDirSkipped(name) {
			continue
		}

		if !matchesChildPathEntry(entry.IsDir(), name, opts) {
			continue
		}

		item := buildChildPathItem(opts.TargetDir, name, entry.IsDir(), opts.IsAbs)
		items = append(items, item)

		if opts.Limit > 0 && len(items) >= opts.Limit {
			break
		}
	}

	return items, nil
}

func isChildPathDirSkipped(name string) bool {
	return handleFindDirSkip(name) != nil
}

func scanChildPathDepthWalk(opts ChildPathOptions) ([]ChildPathItem, error) {
	cleanRoot := filepath.Clean(opts.TargetDir)
	items := make([]ChildPathItem, 0)

	errWalk := filepath.WalkDir(cleanRoot, func(currentPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return handleFindWalkErr(walkErr)
		}

		if currentPath == cleanRoot {
			return nil
		}

		return handleChildPathWalkEntry(currentPath, d, cleanRoot, opts, &items)
	})

	return items, errWalk
}

func resolveBeyondDepthAction(isDir bool) error {
	if isDir {
		return fs.SkipDir
	}

	return nil
}

func handleChildPathWalkEntry(currentPath string, d fs.DirEntry, cleanRoot string, opts ChildPathOptions, items *[]ChildPathItem) error {
	if d.IsDir() && handleFindDirSkip(d.Name()) != nil {
		return fs.SkipDir
	}

	rel, err := filepath.Rel(cleanRoot, currentPath)
	if err != nil {
		return nil
	}

	depth := calculateChildPathDepth(rel)
	isBeyondDepth := opts.Depth > 0 && depth > opts.Depth
	if isBeyondDepth {
		return resolveBeyondDepthAction(d.IsDir())
	}

	return processChildPathWalkMatch(currentPath, d, opts, items)
}

func processChildPathWalkMatch(currentPath string, d fs.DirEntry, opts ChildPathOptions, items *[]ChildPathItem) error {
	if !matchesChildPathEntry(d.IsDir(), d.Name(), opts) {
		return nil
	}

	item := buildChildPathItem(filepath.Dir(currentPath), d.Name(), d.IsDir(), opts.IsAbs)
	*items = append(*items, item)

	if opts.Limit > 0 && len(*items) >= opts.Limit {
		return fs.SkipAll
	}

	return nil
}

func calculateChildPathDepth(relPath string) int {
	norm := filepath.ToSlash(relPath)
	if norm == "." || norm == "" {
		return 0
	}

	return strings.Count(norm, "/") + 1
}

func matchesChildPathEntry(isDir bool, name string, opts ChildPathOptions) bool {
	if !matchesChildPathType(isDir, opts.TypeFilter) {
		return false
	}

	if !matchesChildPathExtensions(isDir, name, opts.Extensions) {
		return false
	}

	return matchesChildPathGlob(name, opts.Glob)
}

func matchesChildPathType(isDir bool, filter string) bool {
	switch filter {
	case "f":
		return !isDir
	case "d":
		return isDir
	default:
		return true
	}
}

func matchesChildPathExtensions(isDir bool, name string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}

	if isDir {
		return false
	}

	fileExt := strings.ToLower(filepath.Ext(name))
	fileExtTrimmed := strings.TrimPrefix(fileExt, ".")

	for _, targetExt := range exts {
		clean := strings.ToLower(strings.TrimPrefix(targetExt, "."))
		if clean == fileExtTrimmed {
			return true
		}
	}

	return false
}

func matchesChildPathGlob(name string, pattern string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}

	matched, err := filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
	if err != nil {
		return false
	}

	return matched
}

func buildChildPathItem(parentDir string, name string, isDir bool, isAbs bool) ChildPathItem {
	fullPath := filepath.Join(parentDir, name)
	displayPath := resolveChildPathDisplayPath(fullPath, isAbs)

	var size int64
	modTime := time.Now().Format(time.RFC3339)

	if info, err := os.Stat(fullPath); err == nil {
		size = info.Size()
		modTime = info.ModTime().Format(time.RFC3339)
	}

	return ChildPathItem{
		Path:    displayPath,
		Name:    name,
		IsDir:   isDir,
		Size:    size,
		ModTime: modTime,
	}
}

func resolveChildPathAbsDisplayPath(fullPath string) string {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return filepath.ToSlash(fullPath)
	}

	return filepath.ToSlash(absPath)
}

func resolveChildPathDisplayPath(fullPath string, isAbs bool) string {
	if isAbs {
		return resolveChildPathAbsDisplayPath(fullPath)
	}

	relPath, err := filepath.Rel(".", fullPath)
	if err == nil {
		return filepath.ToSlash(relPath)
	}

	return filepath.ToSlash(fullPath)
}

func outputChildPathResults(items []ChildPathItem, isJson bool) error {
	if isJson {
		return renderChildPathJson(items)
	}

	renderChildPathTable(items)

	return nil
}

func renderChildPathJson(items []ChildPathItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "json serialization failed")
	}

	fmt.Println(string(data))

	return nil
}

func renderChildPathTable(items []ChildPathItem) {
	if len(items) == 0 {
		fmt.Printf("\n%s(no matching entries found)%s\n\n", constants.ColorDim, constants.ColorReset)

		return
	}

	columns := []termtable.Column{
		{Title: "MODE", MinWidth: 6, Align: termtable.AlignLeft},
		{Title: "SIZE", MinWidth: 10, Align: termtable.AlignRight},
		{Title: "MODIFIED", MinWidth: 19, Align: termtable.AlignLeft},
		{Title: "PATH", MinWidth: 28, Align: termtable.AlignLeft},
	}

	rows := make([]termtable.Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, formatChildPathRow(item))
	}

	termtable.PrintTable(termtable.TableConfig{
		Columns: columns,
		Rows:    rows,
	})

	printChildPathSummary(items)
}

func formatChildPathRow(item ChildPathItem) termtable.Row {
	modeStr := "[FILE]"
	sizeStr := formatChildPathByteSize(item.Size)

	if item.IsDir {
		modeStr = "[DIR]"
		sizeStr = "-"
	}

	modStr := formatChildPathModTime(item.ModTime)

	return termtable.Row{
		Cells: []string{modeStr, sizeStr, modStr, item.Path},
	}
}

func formatChildPathModTime(rfc3339Time string) string {
	parsed, err := time.Parse(time.RFC3339, rfc3339Time)
	if err != nil {
		return rfc3339Time
	}

	return parsed.Format("2006-01-02 15:04:05")
}

func printChildPathSummary(items []ChildPathItem) {
	fileCount := 0
	dirCount := 0

	for _, item := range items {
		if item.IsDir {
			dirCount++
		} else {
			fileCount++
		}
	}

	fmt.Printf("\n  %sTotal:%s %d item(s) (%d files, %d directories)\n\n",
		constants.ColorBold, constants.ColorReset, len(items), fileCount, dirCount)
}

func formatChildPathByteSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}

	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}

	return fmt.Sprintf("%.1f GB", float64(bytes)/(1024*1024*1024))
}
