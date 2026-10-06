package cmdide

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdantigravity"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func runIDEList(args []string) error {
	opts, _ := parseIDEOptions(args)
	repos, err := collectRepoRegistrations()
	if err != nil {
		return err
	}
	if opts.IsJSON {
		b, _ := json.MarshalIndent(repos, "", "  ")
		fmt.Println(string(b))
		return nil
	}
	printRegistrationTable(repos)
	return nil
}

func collectRepoRegistrations() ([]IDERepoRegistration, error) {
	paths, err := fetchStoreRepos()
	if err != nil {
		return nil, err
	}
	hasDesktop := desktop.ResolveCLI() != ""
	items := make([]IDERepoRegistration, len(paths))
	for i, p := range paths {
		items[i] = IDERepoRegistration{
			Name:                filepath.Base(p),
			Path:                p,
			IsVSCodeLinked:      isRepoInVSCode(p),
			IsCursorLinked:      isRepoInCursor(p),
			IsAntigravityLinked: cmdantigravity.IsRepoRegisteredInAgy(p),
			IsDesktopLinked:     hasDesktop,
		}
	}
	return items, nil
}

func fetchStoreRepos() ([]string, error) {
	s, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdide.fetchStoreRepos: open store")
	}
	defer s.Close()
	repos, rErr := s.ListRepos()
	if rErr != nil {
		return nil, apperror.WrapSimple(rErr, "cmdide.fetchStoreRepos: list repos")
	}
	paths := make([]string, 0, len(repos))
	for _, r := range repos {
		if r.AbsolutePath != "" {
			paths = append(paths, r.AbsolutePath)
		} else if r.RelativePath != "" {
			paths = append(paths, r.RelativePath)
		}
	}
	return paths, nil
}

func isPathInEntries(entries []vscodepm.Entry, p string) bool {
	clean := filepath.Clean(p)
	for _, e := range entries {
		if filepath.Clean(e.RootPath) == clean {
			return true
		}
	}
	return false
}

func isRepoInVSCode(p string) bool {
	entries, err := vscodepm.ListEntries()
	return err == nil && isPathInEntries(entries, p)
}

func isRepoInCursor(p string) bool {
	jsonPath, err := cmdcursor.GetCursorProjectsJSONPath()
	if err != nil {
		return false
	}
	entries, err := vscodepm.ReadEntries(jsonPath)
	return err == nil && isPathInEntries(entries, p)
}

func printRegistrationTable(repos []IDERepoRegistration) {
	fmt.Printf("\n%s● IDE Repository Registrations (%d repositories):%s\n", constants.ColorCyan, len(repos), constants.ColorReset)
	fmt.Printf("  %-24s %-6s %-6s %-6s %-6s\n", "REPOSITORY", "VSCODE", "CURSOR", "AGY", "DESKTOP")
	fmt.Printf("  %-24s %-6s %-6s %-6s %-6s\n", "----------", "------", "------", "---", "-------")
	for _, r := range repos {
		fmt.Printf("  %-24s %-6s %-6s %-6s %-6s\n",
			truncateName(r.Name, 24),
			formatCheck(r.IsVSCodeLinked),
			formatCheck(r.IsCursorLinked),
			formatCheck(r.IsAntigravityLinked),
			formatCheck(r.IsDesktopLinked),
		)
	}
	fmt.Println()
}

func formatCheck(val bool) string {
	if val {
		return constants.ColorGreen + "✔" + constants.ColorReset
	}
	return constants.ColorDim + "✖" + constants.ColorReset
}

func truncateName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen-3] + "..."
	}
	return name
}
