package cmdlist

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/stablejson"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractFlagVal(args []string, flagName string) string {
	for i, arg := range args {
		if arg == flagName && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(arg, flagName+"=") {
			return strings.TrimPrefix(arg, flagName+"=")
		}
	}

	return ""
}

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}

// listGroupsHints returns hints shown after gitmap ls groups.
func listGroupsHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintGroupShow, constants.HintGroupShowDesc},
	}
}

// listHints returns hints shown after gitmap list.
func listHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintLsType, constants.HintLsTypeDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
	}
}

func printHints(hints []hintEntry) {
	fmt.Fprint(os.Stderr, constants.MsgHintHeader)
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, constants.MsgHintRowFmt, h.command, h.description)
	}
}

// runGroupList handles "group list".
func runGroupList() error {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	groups, err := db.ListGroups()
	if err != nil && isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)

		return apperror.WrapSimple(err, "fatal error")
	}

	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	printGroupList(db, groups)
	printHints(groupListHints())

	return nil
}

// runProjectRepos handles go-repos, node-repos, react-repos, cpp-repos, csharp-repos.
func runProjectRepos(typeKey string, args []string) error {
	checkHelp(typeKey+"-repos", args)
	jsonOut, countOnly := parseProjectReposFlags(args)
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprint(os.Stderr, constants.MsgProjectNoDB)
		cliexit.HandleError(apperror.WrapSimple(err, "open default store"), 1)
	}

	defer db.Close()

	if countOnly {
		printProjectCount(db, typeKey)

		return nil
	}

	printProjectList(db, typeKey, jsonOut)
	if !jsonOut {
		printHints(projectReposHints())
	}

	return nil
}

type hintEntry struct {
	command     string
	description string
}

// groupListHints returns hints shown after gitmap group list.
func groupListHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintGroupShow, constants.HintGroupShowDesc},
		{constants.HintGroupDelete, constants.HintGroupDeleteDesc},
	}
}

// projectReposHints returns hints shown after go-repos, node-repos, etc.
func projectReposHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupAdd, constants.HintGroupAddDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
		{constants.HintPullGroup, constants.HintPullGroupDesc},
	}
}

// parseProjectReposFlags parses --json and --count flags.
func parseProjectReposFlags(args []string) (bool, bool) {
	fs := flag.NewFlagSet("project-repos", flag.ExitOnError)
	jsonOut := fs.Bool(constants.FlagProjectJSON, false, "Output as JSON")
	countOnly := fs.Bool(constants.FlagProjectCount, false, "Print count only")
	_ = fs.Parse(args)

	return *jsonOut, *countOnly
}

// printProjectCount prints the count of projects for a type.
func printProjectCount(db *store.DB, typeKey string) {
	count, err := db.CountProjectsByTypeKey(typeKey)

	isLegacyErr := err != nil && isLegacyDataError(err)
	if isLegacyErr {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		cliexit.HandleError(apperror.WrapSimple(err, "legacy project data error"), 1)
	}

	if err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, constants.ErrProjectQuery), 1)
	}

	fmt.Printf(constants.MsgProjectCount, count)
}

// printProjectList queries and displays projects for a type.
func printProjectList(db *store.DB, typeKey string, jsonOut bool) {
	projects, err := db.SelectProjectsByTypeKey(typeKey)

	isLegacyErr := err != nil && isLegacyDataError(err)
	if isLegacyErr {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		cliexit.HandleError(apperror.WrapSimple(err, "legacy project data error"), 1)
	}

	if err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, constants.ErrProjectQuery), 1)
	}

	if len(projects) == 0 {
		fmt.Printf(constants.MsgProjectNoneFound, typeKey)

		return
	}

	if jsonOut {
		printProjectsJSON(projects)

		return
	}

	printProjectsTerminal(projects)
	printProjectsSummary(projects)
}

// printProjectsTerminal prints projects in labeled terminal format.
func printProjectsTerminal(projects []model.DetectedProject) {
	for _, p := range projects {
		slug := deriveProjectSlug(p)
		fmt.Printf("  %-6s %s\n", p.ProjectType, p.ProjectName)
		fmt.Printf("         Repo:      %s\n", slug)
		fmt.Printf("         Path:      %s\n", p.AbsolutePath)
		fmt.Printf("         Indicator: %s\n", p.PrimaryIndicator)
		fmt.Printf("         → gitmap cd %s\n\n", slug)
	}
}

// deriveProjectSlug returns the best cd-friendly name for a project.
// It prefers RepoName (the Git folder name) and falls back to the
// last segment of AbsolutePath.
func deriveProjectSlug(p model.DetectedProject) string {
	if p.RepoName != "" {
		return p.RepoName
	}

	return filepath.Base(p.AbsolutePath)
}

// printProjectsSummary prints count summary after the list.
func printProjectsSummary(projects []model.DetectedProject) {
	fmt.Fprintf(os.Stderr, constants.MsgProjectListCount, len(projects))
}

// printProjectsJSON prints projects as formatted JSON via stablejson.
func printProjectsJSON(projects []model.DetectedProject) {
	if err := encodeProjectReposJSON(os.Stdout, projects); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode projects JSON: %v\n", err)
	}
}

// project-repos wire keys. Names + order are the contract.
const (
	projectReposKeyID               = "id"
	projectReposKeyRepoID           = "repoId"
	projectReposKeyRepoName         = "repoName"
	projectReposKeyProjectTypeID    = "projectTypeId"
	projectReposKeyProjectType      = "projectType"
	projectReposKeyProjectName      = "projectName"
	projectReposKeyAbsolutePath     = "absolutePath"
	projectReposKeyRepoPath         = "repoPath"
	projectReposKeyRelativePath     = "relativePath"
	projectReposKeyPrimaryIndicator = "primaryIndicator"
	projectReposKeyDetectedAt       = "detectedAt"
)

// encodeProjectReposJSON writes projects as a stablejson 2-space-indented
// array. Empty input emits `[]\n`.
func encodeProjectReposJSON(w io.Writer, projects []model.DetectedProject) error {
	return stablejson.WriteArray(w, buildProjectReposJSONItems(projects))
}

// buildProjectReposJSONItems is the single source of (field name,
// field order, value) for project-repos.
func buildProjectReposJSONItems(projects []model.DetectedProject) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(projects))
	for _, p := range projects {
		items = append(items, []stablejson.Field{
			{Key: projectReposKeyID, Value: p.ID},
			{Key: projectReposKeyRepoID, Value: p.RepoID},
			{Key: projectReposKeyRepoName, Value: p.RepoName},
			{Key: projectReposKeyProjectTypeID, Value: p.ProjectTypeID},
			{Key: projectReposKeyProjectType, Value: p.ProjectType},
			{Key: projectReposKeyProjectName, Value: p.ProjectName},
			{Key: projectReposKeyAbsolutePath, Value: p.AbsolutePath},
			{Key: projectReposKeyRepoPath, Value: p.RepoPath},
			{Key: projectReposKeyRelativePath, Value: p.RelativePath},
			{Key: projectReposKeyPrimaryIndicator, Value: p.PrimaryIndicator},
			{Key: projectReposKeyDetectedAt, Value: p.DetectedAt},
		})
	}

	return items
}

// printGroupList renders the group table to stdout.
func printGroupList(db *store.DB, groups []model.Group) {
	if len(groups) == 0 {
		fmt.Println(constants.MsgGroupEmpty)

		return
	}

	fmt.Println(constants.MsgGroupHeader)
	fmt.Println(constants.MsgListSeparator)
	for _, g := range groups {
		count, countErr := db.CountGroupRepos(g.Name)
		if countErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not count repos for group %s: %v\n", g.Name, countErr)
		}

		fmt.Printf(constants.MsgGroupRowFmt, g.Name, count, g.Description)
	}
}

// resolveCurrentRepoID returns the RepoId for the cwd. Caller must handle
// the "repo not scanned" error path explicitly.
func resolveCurrentRepoID(db *store.DB) (int64, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	return db.ResolveCurrentRepoID(cwd)
}
