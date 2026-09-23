package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// HandlePEFormatCommands inspects args for format management and testing commands.
func HandlePEFormatCommands(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "add-format", "addformat", "save-format":
		return true, handleAddFormatCLI(args[1:])
	case "rm-format", "rmformat", "delete-format":
		return true, handleRmFormatCLI(args[1:])
	case "add-all", "addall", "import-formats":
		return true, handleAddAllFormatsCLI(args[1:])
	case "list-formats", "formats", "format-list":
		return true, handleListFormatsCLI()
	case "preview-format", "previewformat":
		return true, handlePreviewFormatCLI(args[1:])
	}
	if hasTestFlags(args) {
		return true, handleFormatTestCLI(args)
	}
	return false, nil
}

func hasTestFlags(args []string) bool {
	return hasArgFlag(args, "-test") || hasArgFlag(args, "--test") ||
		hasArgFlag(args, "-test-commit") || hasArgFlag(args, "--test-commit")
}

func handleAddFormatCLI(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap pe add-format <file.json> [alias]")
	}
	filePath := args[0]
	alias := ""
	if len(args) > 1 {
		alias = args[1]
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return apperror.WrapSimple(err, "read format json file "+filePath)
	}
	var p PEFormatProfile
	if jsonErr := json.Unmarshal(data, &p); jsonErr != nil {
		return apperror.WrapSimple(jsonErr, "parse format json file "+filePath)
	}
	if p.Name == "" {
		p.Name = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}
	if alias != "" {
		p.Alias = alias
	}
	if saveErr := SavePEFormatProfile(p); saveErr != nil {
		return saveErr
	}
	fmt.Printf("✔ Format profile '%s' (alias: '%s') registered successfully.\n", p.Name, p.Alias)
	return nil
}

func handleRmFormatCLI(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap pe rm-format <name|alias>")
	}
	name := args[0]
	if err := DeletePEFormatProfile(name); err != nil {
		return err
	}
	fmt.Printf("✔ Format profile '%s' removed.\n", name)
	return nil
}

func handleAddAllFormatsCLI(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap pe add-all <folder-path>")
	}
	folder := args[0]
	count, err := ImportPEFormatProfilesFromFolder(folder)
	if err != nil {
		return err
	}
	fmt.Printf("✔ Registered %d format profile(s) from '%s'.\n", count, folder)
	return nil
}

func handleListFormatsCLI() error {
	profiles, err := ListPEFormatProfiles()
	if err != nil {
		return err
	}
	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "NAME", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "STRIP PREFIXES", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "DESCRIPTION", Align: termtable.AlignLeft, MinWidth: 40},
		},
		Rows: buildFormatRows(profiles),
	}
	fmt.Println("\nRegistered Pipeline Error Format Profiles:")
	termtable.PrintTable(cfg)
	fmt.Println()

	return nil
}

func buildFormatRows(profiles []PEFormatProfile) []termtable.Row {
	var rows []termtable.Row
	for _, p := range profiles {
		stripCount := fmt.Sprintf("%d prefixes", len(p.StripPrefixes))
		rows = append(rows, termtable.Row{
			Cells: []string{p.Name, p.Alias, stripCount, p.Description},
		})
	}

	return rows
}

func handlePreviewFormatCLI(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap pe preview-format <name|file.json>")
	}
	p, err := LoadPEFormatProfile(args[0])
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(data))
	return nil
}

func handleFormatTestCLI(args []string) error {
	formatName := extractFlagVal(args, "-f")
	if formatName == "" {
		formatName = extractFlagVal(args, "--format")
	}
	if formatName == "" {
		formatName = "default"
	}
	profile, err := LoadPEFormatProfile(formatName)
	if err != nil {
		return err
	}
	testFile := extractFlagVal(args, "-test")
	if testFile == "" {
		testFile = extractFlagVal(args, "--test")
	}
	if testFile != "" {
		return executeFormatTestOnFile(testFile, profile)
	}
	return executeFormatTestOnCommit(args, profile)
}

func executeFormatTestOnFile(filePath string, p *PEFormatProfile) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return apperror.WrapSimple(err, "read test file "+filePath)
	}
	filtered := FilterLogWithProfile(string(data), p)
	fmt.Printf("\n▶ Test Filter Results using format '%s' on %s:\n", p.Name, filePath)
	fmt.Println(strings.Repeat("─", 60))
	if strings.TrimSpace(filtered) == "" {
		fmt.Println("(All lines filtered out or no errors/warnings detected)")
	} else {
		fmt.Println(filtered)
	}
	fmt.Println(strings.Repeat("─", 60))
	return nil
}

func executeFormatTestOnCommit(args []string, p *PEFormatProfile) error {
	sha := extractFlagVal(args, "-test-commit")
	if sha == "" {
		sha = extractFlagVal(args, "--test-commit")
	}
	repo := extractFlagVal(args, "-repo")
	if repo == "" {
		repo = extractFlagVal(args, "--repo")
	}
	if repo == "" {
		repo = resolveCurrentRepoSlug()
	}
	fmt.Printf("\n▶ Test Filter Results using format '%s' on commit %s (%s):\n", p.Name, sha, repo)
	runs := queryWorkflowRunsForTarget(repo, sha)
	if len(runs) == 0 {
		fmt.Println("No workflow runs found for commit " + sha)
		return nil
	}
	for _, r := range runs {
		rawLogs := queryAllRunLogs(repo, r.DatabaseId)
		filtered := FilterLogWithProfile(rawLogs, p)
		if strings.TrimSpace(filtered) != "" {
			fmt.Printf("─── Workflow: %s (#%d) ───\n%s\n", r.Name, r.DatabaseId, filtered)
		}
	}
	return nil
}
