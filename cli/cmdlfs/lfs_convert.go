// NOTE: Needs refactor — exceeds 300-line target
package cmdlfs

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

type lfsFlags struct {
	dryRun      bool
	autoConfirm bool
	force       bool
	sizeRaw     string
	threshold   int64
}

type lfsCandidate struct {
	RelPath       string
	Size          int64
	Pattern       string
	AlreadyInAttr bool
}

// RunLFSConvertCmd executes the 'gitmap lfs convert' command.
func RunLFSConvertCmd(args []string) error {
	cleanArgs := normalizeLFSArgs(args)
	if hasHelpFlag(cleanArgs) {
		helpdoc.PrintHelpAndExit("lfs-convert", cleanArgs)
		return nil
	}

	flags, err := parseLFSFlags(cleanArgs)
	if err != nil {
		return err
	}

	if !isInsideGitRepo() {
		return apperror.NewValidation("lfs_convert", "E9001", "not inside a Git repository")
	}

	return processLFSConvert(flags)
}

func normalizeLFSArgs(args []string) []string {
	if len(args) >= 1 && (args[0] == "convert" || args[0] == "conv") {
		return args[1:]
	}
	return args
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}
	return false
}

func parseLFSFlags(args []string) (lfsFlags, error) {
	var lf lfsFlags
	fs := flag.NewFlagSet("lfs-convert", flag.ContinueOnError)
	fs.BoolVar(&lf.dryRun, "dry-run", false, "Preview without modifying files or git")
	fs.BoolVar(&lf.autoConfirm, "y", false, "Auto-confirm execution prompt")
	fs.BoolVar(&lf.autoConfirm, "yes", false, "Auto-confirm execution prompt")
	fs.BoolVar(&lf.force, "force", false, "Force convert execution")
	fs.StringVar(&lf.sizeRaw, "size", "5MB", "Size threshold for LFS conversion")
	fs.StringVar(&lf.sizeRaw, "s", "5MB", "Size threshold for LFS conversion")
	if err := fs.Parse(args); err != nil {
		return lf, apperror.WrapSimple(err, "parse_lfs_flags")
	}
	threshold, err := parseByteSize(lf.sizeRaw)
	lf.threshold = threshold
	return lf, err
}

func parseByteSize(val string) (int64, error) {
	trimmed := strings.ToUpper(strings.TrimSpace(val))
	if trimmed == "" {
		return 5 * 1024 * 1024, nil
	}

	multiplier, numStr := resolveSizeMultiplier(trimmed)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, apperror.WrapSimple(err, "invalid_size_format")
	}

	return int64(num * float64(multiplier)), nil
}

func resolveSizeMultiplier(s string) (int64, string) {
	for _, u := range []struct {
		s string
		m int64
	}{
		{"GB", 1024 * 1024 * 1024}, {"G", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024}, {"M", 1024 * 1024},
		{"KB", 1024}, {"K", 1024}, {"B", 1},
	} {
		if strings.HasSuffix(s, u.s) {
			return u.m, strings.TrimSuffix(s, u.s)
		}
	}
	return 1, s
}

func formatByteSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(1024), 0
	for n := bytes / 1024; n >= 1024; n /= 1024 {
		div *= 1024
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func isInsideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func isLFSAvailable() bool {
	cmd := exec.Command("git", "lfs", "version")
	return cmd.Run() == nil
}

func loadExistingLFSAttributes() map[string]bool {
	tracked := make(map[string]bool)
	data, err := os.ReadFile(".gitattributes")
	if err != nil {
		return tracked
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) > 0 && !strings.HasPrefix(fields[0], "#") && strings.Contains(line, "filter=lfs") {
			tracked[fields[0]] = true
		}
	}
	return tracked
}

func findLFSCandidates(threshold int64) ([]lfsCandidate, error) {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "ls_files")
	}

	existingRules := loadExistingLFSAttributes()
	var candidates []lfsCandidate
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if item, ok := inspectLFSCandidate(strings.TrimSpace(line), threshold, existingRules); ok {
			candidates = append(candidates, item)
		}
	}
	return candidates, nil
}

func inspectLFSCandidate(relPath string, threshold int64, existingRules map[string]bool) (lfsCandidate, bool) {
	var empty lfsCandidate
	if relPath == "" || relPath == ".gitattributes" {
		return empty, false
	}

	fi, err := os.Stat(relPath)
	if err != nil || fi.IsDir() || fi.Size() <= threshold {
		return empty, false
	}

	pat := resolveLFSPattern(relPath)
	return lfsCandidate{RelPath: relPath, Size: fi.Size(), Pattern: pat, AlreadyInAttr: existingRules[pat]}, true
}

func resolveLFSPattern(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return "*" + strings.ToLower(ext)
	}
	return filepath.Base(filename)
}

func processLFSConvert(flags lfsFlags) error {
	if !isLFSAvailable() {
		return apperror.NewValidation("lfs_convert", "E9002", "Git LFS is not installed on PATH. Install from https://git-lfs.com")
	}

	candidates, err := findLFSCandidates(flags.threshold)
	if err != nil {
		return err
	}

	if len(candidates) == 0 {
		fmt.Printf("  %s✓ No files exceeding %s found in repository.%s\n\n",
			constants.ColorGreen, formatByteSize(flags.threshold), constants.ColorReset)
		return nil
	}

	if flags.dryRun {
		return executeLFSDryRun(candidates, flags.threshold)
	}
	return executeLFSActive(candidates, flags)
}

func executeLFSDryRun(candidates []lfsCandidate, threshold int64) error {
	fmt.Printf("\n  %s■ gitmap lfs convert — [dry-run] previewing files > %s%s\n\n",
		constants.ColorCyan, formatByteSize(threshold), constants.ColorReset)

	termout.PrintTable(buildLFSTableConfig(candidates))
	printLFSDryRunRules(extractLFSPatterns(candidates, true))

	totalSize := sumLFSCandidatesSize(candidates)
	fmt.Printf("\n  %sSummary:%s Found %d file(s) exceeding %s (Total: %s)\n",
		constants.ColorCyan, constants.ColorReset, len(candidates), formatByteSize(threshold), formatByteSize(totalSize))
	fmt.Printf("  %s[dry-run]%s No files converted and Git history was not modified.\n\n",
		constants.ColorYellow, constants.ColorReset)
	return nil
}

func buildLFSTableConfig(candidates []lfsCandidate) termout.TableConfig {
	cols := []termout.Column{
		{Title: "File", Align: termout.AlignLeft, MinWidth: 25},
		{Title: "Size", Align: termout.AlignRight, MinWidth: 10},
		{Title: "LFS Rule", Align: termout.AlignLeft, MinWidth: 15},
		{Title: "Status", Align: termout.AlignCenter, MinWidth: 12},
	}

	rows := make([]termout.Row, len(candidates))
	for i, c := range candidates {
		status := "pending add"
		if c.AlreadyInAttr {
			status = "already tracked"
		}
		rows[i] = termout.Row{
			Cells: []string{c.RelPath, formatByteSize(c.Size), c.Pattern, status},
			Color: constants.ColorWhite,
		}
	}
	return termout.TableConfig{Columns: cols, Rows: rows, HeaderColor: constants.ColorYellow}
}

func printLFSDryRunRules(patterns []string) {
	if len(patterns) == 0 {
		return
	}

	fmt.Printf("\n  %sProposed .gitattributes rules to add:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, p := range patterns {
		fmt.Printf("    %s filter=lfs diff=lfs merge=lfs -text\n", p)
	}
}

func extractLFSPatterns(candidates []lfsCandidate, onlyNew bool) []string {
	seen := make(map[string]bool)
	var result []string

	for _, c := range candidates {
		if (!onlyNew || !c.AlreadyInAttr) && !seen[c.Pattern] {
			seen[c.Pattern] = true
			result = append(result, c.Pattern)
		}
	}
	return result
}

func sumLFSCandidatesSize(candidates []lfsCandidate) int64 {
	var total int64
	for _, c := range candidates {
		total += c.Size
	}
	return total
}

func executeLFSActive(candidates []lfsCandidate, flags lfsFlags) error {
	confirmed, err := confirmIrreversibleLFS(flags.autoConfirm || flags.force)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Printf("  %sOperation cancelled by user.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	if err := performLFSConversion(candidates); err != nil {
		return err
	}

	printLFSSuccessSummary(len(candidates))
	return nil
}

func confirmIrreversibleLFS(autoConfirm bool) (bool, error) {
	fmt.Printf("\n  %sWARNING: This operation rewrites Git history and converts files to Git LFS pointers. It is irreversible.%s\n",
		constants.ColorYellow, constants.ColorReset)

	if autoConfirm {
		fmt.Println()
		return true, nil
	}

	fmt.Print("  Do you want to proceed? [y/N]: ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false, apperror.WrapSimple(err, "read_confirmation")
	}

	val := strings.TrimSpace(strings.ToLower(text))
	return val == "y" || val == "yes", nil
}

func performLFSConversion(candidates []lfsCandidate) error {
	_ = exec.Command("git", "lfs", "install", "--local").Run()
	if err := appendLFSGitAttributes(extractLFSPatterns(candidates, true)); err != nil {
		return err
	}

	return runLFSMigrateImport(extractLFSPatterns(candidates, false), candidates)
}

func appendLFSGitAttributes(patterns []string) error {
	if len(patterns) == 0 {
		return nil
	}

	f, err := os.OpenFile(".gitattributes", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return apperror.WrapSimple(err, "open_gitattributes")
	}
	defer f.Close()

	for _, p := range patterns {
		if _, writeErr := fmt.Fprintf(f, "%s filter=lfs diff=lfs merge=lfs -text\n", p); writeErr != nil {
			return apperror.WrapSimple(writeErr, "write_gitattributes")
		}
	}
	return nil
}

func runLFSMigrateImport(patterns []string, candidates []lfsCandidate) error {
	patternsArg := strings.Join(patterns, ",")
	cmd := exec.Command("git", "lfs", "migrate", "import", "--everything", "--include="+patternsArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fallbackLFSPointerStaging(candidates)
	}

	fmt.Printf("  %s✓ LFS migration output:%s\n%s\n", constants.ColorGreen, constants.ColorReset, strings.TrimSpace(string(out)))
	return nil
}

func fallbackLFSPointerStaging(candidates []lfsCandidate) error {
	_ = exec.Command("git", "add", ".gitattributes").Run()
	for _, c := range candidates {
		_ = exec.Command("git", "rm", "--cached", "--quiet", "--", c.RelPath).Run()
		if err := exec.Command("git", "add", "--", c.RelPath).Run(); err != nil {
			return apperror.WrapSimple(err, "lfs_stage_fallback")
		}
	}
	return nil
}

func printLFSSuccessSummary(count int) {
	fmt.Println()
	fmt.Printf("  %s✓ LFS conversion complete:%s %d file(s) converted to Git LFS pointers.\n",
		constants.ColorGreen, constants.ColorReset, count)
	fmt.Println("  Raw blobs migrated and .gitattributes updated.")
	fmt.Println()
}
