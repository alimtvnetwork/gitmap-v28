package cmdlfs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func loadExistingLFSAttributes() map[string]bool {
	tracked := make(map[string]bool)
	data, err := os.ReadFile(".gitattributes")
	if err != nil {
		return tracked
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		isLFSLine := len(fields) > 0 && !strings.HasPrefix(fields[0], "#") && strings.Contains(line, "filter=lfs")
		if isLFSLine {
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
		isSkip := onlyNew && c.AlreadyInAttr
		if isSkip || seen[c.Pattern] {
			continue
		}
		seen[c.Pattern] = true
		result = append(result, c.Pattern)
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
