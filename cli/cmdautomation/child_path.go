package cmdautomation

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

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

var (
	aumChildPathRunner func(args []string) error

	childPathCmd = &cobra.Command{
		Use:                "child-path [dir] [glob]",
		Aliases:            []string{"childpath", "cpth"},
		Short:              "High-speed directory inventory and search replacement for Get-ChildItem and git grep",
		DisableFlagParsing: true,
		Example: `  gitmap aum child-path . "*.go" --depth 2 --type f
  gitmap aum child-path cli "*.go" --depth 1 --json --limit 5
  gitmap aum cpth src -e ts,tsx -d 3`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunAumChildPath(args)
		},
	}
)

func init() {
	AutomationCmd.AddCommand(childPathCmd)
}

// RegisterChildPathRunner connects package cmd's executor to automation.
func RegisterChildPathRunner(fn func(args []string) error) {
	aumChildPathRunner = fn
}

// RunAumChildPath exposes the child-path automation command.
func RunAumChildPath(args []string) error {
	if aumChildPathRunner != nil {
		return aumChildPathRunner(args)
	}

	return runChildPathDirect(args)
}

type aumChildOptions struct {
	targetDir  string
	glob       string
	depth      int
	extensions []string
	typeFilter string
	limit      int
	isJson     bool
	isAbs      bool
}

type aumChildItem struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

func runChildPathDirect(args []string) error {
	opts, err := parseAumChildArgs(args)
	if err != nil {
		return err
	}

	if err := validateAumChildOptions(opts); err != nil {
		return err
	}

	return executeAumChildPath(opts)
}

func parseAumChildArgs(args []string) (aumChildOptions, error) {
	opts := defaultAumChildOptions()
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		consumed, err := parseSingleAumArg(args, i, &opts, &nonFlags)
		if err != nil {
			return opts, err
		}

		i += consumed
	}

	resolveAumTargets(nonFlags, &opts)

	return opts, nil
}

func defaultAumChildOptions() aumChildOptions {
	return aumChildOptions{
		targetDir:  ".",
		glob:       "*",
		depth:      1,
		typeFilter: "all",
		limit:      0,
		isJson:     false,
		isAbs:      false,
	}
}

func parseSingleAumArg(args []string, idx int, opts *aumChildOptions, nonFlags *[]string) (int, error) {
	arg := args[idx]
	hasNext := idx+1 < len(args)
	nextArg := ""

	if hasNext {
		nextArg = args[idx+1]
	}

	if arg == "--json" || arg == "-j" {
		opts.isJson = true

		return 0, nil
	}

	if arg == "--abs" {
		opts.isAbs = true

		return 0, nil
	}

	if arg == "--rel" {
		opts.isAbs = false

		return 0, nil
	}

	return parseAumParameterizedFlags(arg, nextArg, hasNext, opts, nonFlags)
}

func parseAumParameterizedFlags(arg string, nextArg string, hasNext bool, opts *aumChildOptions, nonFlags *[]string) (int, error) {
	if consumed, err := parseAumDepth(arg, nextArg, hasNext, opts); consumed > 0 || err != nil || isAumDepth(arg) {
		return consumed, err
	}

	if consumed := parseAumExt(arg, nextArg, hasNext, opts); consumed > 0 || isAumExt(arg) {
		return consumed, nil
	}

	if consumed := parseAumType(arg, nextArg, hasNext, opts); consumed > 0 || isAumType(arg) {
		return consumed, nil
	}

	if consumed, err := parseAumLimit(arg, nextArg, hasNext, opts); consumed > 0 || err != nil || isAumLimit(arg) {
		return consumed, err
	}

	if !strings.HasPrefix(arg, "-") {
		*nonFlags = append(*nonFlags, arg)
	}

	return 0, nil
}

func isAumDepth(arg string) bool {
	return strings.HasPrefix(arg, "--depth=") || strings.HasPrefix(arg, "-d=")
}

func isAumExt(arg string) bool {
	return strings.HasPrefix(arg, "--ext=") || strings.HasPrefix(arg, "-e=")
}

func isAumType(arg string) bool {
	return strings.HasPrefix(arg, "--type=") || strings.HasPrefix(arg, "-t=")
}

func isAumLimit(arg string) bool {
	return strings.HasPrefix(arg, "--limit=") || strings.HasPrefix(arg, "-l=") || strings.HasPrefix(arg, "-n=")
}

func extractAumDepthRaw(arg string, nextArg string, hasNext bool) (string, int, bool) {
	if (arg == "--depth" || arg == "-d") && hasNext {
		return nextArg, 1, true
	}

	if isAumDepth(arg) {
		parts := strings.SplitN(arg, "=", 2)

		return parts[1], 0, true
	}

	return "", 0, false
}

func parseAumDepth(arg string, nextArg string, hasNext bool, opts *aumChildOptions) (int, error) {
	raw, consumed, matched := extractAumDepthRaw(arg, nextArg, hasNext)
	if !matched {
		return 0, nil
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, apperror.NewValidation("child-path", "E1084", "depth must be an integer")
	}

	opts.depth = val

	return consumed, nil
}

func parseAumExt(arg string, nextArg string, hasNext bool, opts *aumChildOptions) int {
	if (arg == "--ext" || arg == "-e") && hasNext {
		opts.extensions = append(opts.extensions, splitExtensionsAum(nextArg)...)

		return 1
	}

	if isAumExt(arg) {
		parts := strings.SplitN(arg, "=", 2)
		opts.extensions = append(opts.extensions, splitExtensionsAum(parts[1])...)

		return 0
	}

	return 0
}

func parseAumType(arg string, nextArg string, hasNext bool, opts *aumChildOptions) int {
	if (arg == "--type" || arg == "-t") && hasNext {
		opts.typeFilter = strings.ToLower(nextArg)

		return 1
	}

	if isAumType(arg) {
		parts := strings.SplitN(arg, "=", 2)
		opts.typeFilter = strings.ToLower(parts[1])

		return 0
	}

	return 0
}

func extractAumLimitRaw(arg string, nextArg string, hasNext bool) (string, int, bool) {
	if (arg == "--limit" || arg == "-l" || arg == "-n") && hasNext {
		return nextArg, 1, true
	}

	if isAumLimit(arg) {
		parts := strings.SplitN(arg, "=", 2)

		return parts[1], 0, true
	}

	return "", 0, false
}

func parseAumLimit(arg string, nextArg string, hasNext bool, opts *aumChildOptions) (int, error) {
	raw, consumed, matched := extractAumLimitRaw(arg, nextArg, hasNext)
	if !matched {
		return 0, nil
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, apperror.NewValidation("child-path", "E1000", "limit must be an integer")
	}

	opts.limit = val

	return consumed, nil
}

func splitExtensionsAum(val string) []string {
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

func resolveAumTargets(nonFlags []string, opts *aumChildOptions) {
	if len(nonFlags) == 0 {
		return
	}

	if len(nonFlags) == 1 {
		resolveSingleAumTarget(nonFlags[0], opts)

		return
	}

	opts.targetDir = nonFlags[0]
	opts.glob = nonFlags[1]
}

func resolveSingleAumTarget(arg string, opts *aumChildOptions) {
	if strings.ContainsAny(arg, "*?[]") {
		opts.glob = arg

		return
	}

	if info, err := os.Stat(arg); err == nil && info.IsDir() {
		opts.targetDir = arg

		return
	}

	if strings.Contains(arg, ".") {
		opts.glob = arg

		return
	}

	opts.targetDir = arg
}

func validateAumChildOptions(opts aumChildOptions) error {
	if opts.depth < 0 {
		return apperror.NewValidation("child-path", "E1084", fmt.Sprintf("depth %d is out of bounds (must be >= 0)", opts.depth))
	}

	if opts.typeFilter != "f" && opts.typeFilter != "d" && opts.typeFilter != "all" {
		return apperror.NewValidation("child-path", "E1082", fmt.Sprintf("invalid type filter %q (must be 'f', 'd', or 'all')", opts.typeFilter))
	}

	if _, err := filepath.Match(opts.glob, ""); errors.Is(err, filepath.ErrBadPattern) {
		return apperror.NewValidation("child-path", "E1083", fmt.Sprintf("invalid glob pattern %q", opts.glob))
	}

	return validateAumDir(opts.targetDir)
}

func validateAumDir(dir string) error {
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

func executeAumChildPath(opts aumChildOptions) error {
	var items []aumChildItem
	var err error

	if opts.depth == 1 {
		items, err = scanAumDepthOne(opts)
	} else {
		items, err = scanAumDepthWalk(opts)
	}

	if err != nil {
		return apperror.WrapSimple(err, "child-path scan failed")
	}

	return outputAumResults(items, opts.isJson)
}

func scanAumDepthOne(opts aumChildOptions) ([]aumChildItem, error) {
	entries, err := os.ReadDir(opts.targetDir)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read directory failed")
	}

	items := make([]aumChildItem, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() && isAumSkipDir(name) {
			continue
		}

		if !matchesAumEntry(entry.IsDir(), name, opts) {
			continue
		}

		item := buildAumItem(opts.targetDir, name, entry.IsDir(), opts.isAbs)
		items = append(items, item)

		if opts.limit > 0 && len(items) >= opts.limit {
			break
		}
	}

	return items, nil
}

func isAumSkipDir(name string) bool {
	return name == ".git" || name == ".tmp" || name == "node_modules" || name == ".gitmap"
}

func scanAumDepthWalk(opts aumChildOptions) ([]aumChildItem, error) {
	cleanRoot := filepath.Clean(opts.targetDir)
	items := make([]aumChildItem, 0)

	errWalk := filepath.WalkDir(cleanRoot, func(currentPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || currentPath == cleanRoot {
			return nil
		}

		return handleAumWalkStep(currentPath, d, cleanRoot, opts, &items)
	})

	return items, errWalk
}

func resolveBeyondAumDepthAction(isDir bool) error {
	if isDir {
		return fs.SkipDir
	}

	return nil
}

func handleAumWalkStep(currentPath string, d fs.DirEntry, cleanRoot string, opts aumChildOptions, items *[]aumChildItem) error {
	if d.IsDir() && isAumSkipDir(d.Name()) {
		return fs.SkipDir
	}

	rel, err := filepath.Rel(cleanRoot, currentPath)
	if err != nil {
		return nil
	}

	depth := calculateAumDepth(rel)
	isBeyondDepth := opts.depth > 0 && depth > opts.depth
	if isBeyondDepth {
		return resolveBeyondAumDepthAction(d.IsDir())
	}

	return processAumWalkMatch(currentPath, d, opts, items)
}

func processAumWalkMatch(currentPath string, d fs.DirEntry, opts aumChildOptions, items *[]aumChildItem) error {
	if !matchesAumEntry(d.IsDir(), d.Name(), opts) {
		return nil
	}

	item := buildAumItem(filepath.Dir(currentPath), d.Name(), d.IsDir(), opts.isAbs)
	*items = append(*items, item)

	if opts.limit > 0 && len(*items) >= opts.limit {
		return fs.SkipAll
	}

	return nil
}

func calculateAumDepth(relPath string) int {
	norm := filepath.ToSlash(relPath)
	if norm == "." || norm == "" {
		return 0
	}

	return strings.Count(norm, "/") + 1
}

func matchesAumEntry(isDir bool, name string, opts aumChildOptions) bool {
	if !matchesAumType(isDir, opts.typeFilter) {
		return false
	}

	if !matchesAumExtensions(isDir, name, opts.extensions) {
		return false
	}

	return matchesAumGlob(name, opts.glob)
}

func matchesAumType(isDir bool, filter string) bool {
	switch filter {
	case "f":
		return !isDir
	case "d":
		return isDir
	default:
		return true
	}
}

func matchesAumExtensions(isDir bool, name string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}

	if isDir {
		return false
	}

	fileExtTrimmed := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")

	for _, targetExt := range exts {
		clean := strings.TrimPrefix(strings.ToLower(targetExt), ".")
		if clean == fileExtTrimmed {
			return true
		}
	}

	return false
}

func matchesAumGlob(name string, pattern string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}

	matched, err := filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
	if err != nil {
		return false
	}

	return matched
}

func buildAumItem(parentDir string, name string, isDir bool, isAbs bool) aumChildItem {
	fullPath := filepath.Join(parentDir, name)
	displayPath := resolveAumDisplayPath(fullPath, isAbs)

	var size int64
	modTime := time.Now().Format(time.RFC3339)

	if info, err := os.Stat(fullPath); err == nil {
		size = info.Size()
		modTime = info.ModTime().Format(time.RFC3339)
	}

	return aumChildItem{
		Path:    displayPath,
		Name:    name,
		IsDir:   isDir,
		Size:    size,
		ModTime: modTime,
	}
}

func resolveAumAbsDisplayPath(fullPath string) string {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return filepath.ToSlash(fullPath)
	}

	return filepath.ToSlash(absPath)
}

func resolveAumDisplayPath(fullPath string, isAbs bool) string {
	if isAbs {
		return resolveAumAbsDisplayPath(fullPath)
	}

	relPath, err := filepath.Rel(".", fullPath)
	if err == nil {
		return filepath.ToSlash(relPath)
	}

	return filepath.ToSlash(fullPath)
}

func outputAumResults(items []aumChildItem, isJson bool) error {
	if isJson {
		return renderAumJson(items)
	}

	renderAumTable(items)

	return nil
}

func renderAumJson(items []aumChildItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "json serialization failed")
	}

	fmt.Println(string(data))

	return nil
}

func renderAumTable(items []aumChildItem) {
	if len(items) == 0 {
		fmt.Printf("\n%s(no matching entries found)%s\n\n", constants.ColorDim, constants.ColorReset)

		return
	}

	columns := []termout.Column{
		{Title: "MODE", MinWidth: 6, Align: termout.AlignLeft},
		{Title: "SIZE", MinWidth: 10, Align: termout.AlignRight},
		{Title: "MODIFIED", MinWidth: 19, Align: termout.AlignLeft},
		{Title: "PATH", MinWidth: 28, Align: termout.AlignLeft},
	}

	rows := make([]termout.Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, formatAumRow(item))
	}

	termout.PrintTable(termout.TableConfig{
		Columns: columns,
		Rows:    rows,
	})

	printAumSummary(items)
}

func formatAumRow(item aumChildItem) termout.Row {
	modeStr := "[FILE]"
	sizeStr := formatAumByteSize(item.Size)

	if item.IsDir {
		modeStr = "[DIR]"
		sizeStr = "-"
	}

	modStr := item.ModTime
	if parsed, err := time.Parse(time.RFC3339, item.ModTime); err == nil {
		modStr = parsed.Format("2006-01-02 15:04:05")
	}

	return termout.Row{
		Cells: []string{modeStr, sizeStr, modStr, item.Path},
	}
}

func printAumSummary(items []aumChildItem) {
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

func formatAumByteSize(bytes int64) string {
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
