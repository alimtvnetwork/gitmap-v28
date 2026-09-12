package cmdpull

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RemediationItem holds remediation data for failed pulls.
type RemediationItem struct {
	RepoPath      string                      `json:"repoPath"`
	RepoName      string                      `json:"repoName"`
	SummaryReason string                      `json:"summaryReason"`
	Recipes       []gitutil.RemediationRecipe `json:"recipes"`
	Files         []string                    `json:"files,omitempty"`
}

// Delegate hooks
var (
	LoadAllRecordsDBFn                func() []model.ScanRecord
	LoadRecordsByGroupFn              func(string) []model.ScanRecord
	CreatePendingTaskFn               func(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB)
	CompletePendingTaskFn             func(db *store.DB, taskID int64)
	FailPendingTaskFn                 func(db *store.DB, taskID int64, reason string)
	RequireOnlineFn                   func()
	CheckHelpFn                       func(string, []string)
	PrintRemediationSummaryNoPromptFn func(items []RemediationItem)
	PrintRemediationSummaryAutoFixFn  func(items []RemediationItem)
	PrintRemediationSummaryFn         func(items []RemediationItem)
	ApplyTransportFlagFn              func(dir string, useSSH, useHTTPS bool) (bool, string, string, error)
	HasAliasFn                        func() bool
	GetAliasSlugFn                    func() string
	GetAliasPathFn                    func() string
	RunStatusFn                       func(args []string) error
)

func runStatus(args []string) error {
	if RunStatusFn != nil {
		return RunStatusFn(args)
	}
	return nil
}

func HasAlias() bool {
	if HasAliasFn != nil {
		return HasAliasFn()
	}
	return false
}

func GetAliasSlug() string {
	if GetAliasSlugFn != nil {
		return GetAliasSlugFn()
	}
	return ""
}

func GetAliasPath() string {
	if GetAliasPathFn != nil {
		return GetAliasPathFn()
	}
	return ""
}

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func loadAllRecordsDB() []model.ScanRecord {
	if LoadAllRecordsDBFn != nil {
		return LoadAllRecordsDBFn()
	}
	return nil
}

func loadRecordsByGroup(group string) []model.ScanRecord {
	if LoadRecordsByGroupFn != nil {
		return LoadRecordsByGroupFn(group)
	}
	return nil
}

func createPendingTask(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB) {
	if CreatePendingTaskFn != nil {
		return CreatePendingTaskFn(typeName, targetPath, workDir, sourceCmd, cmdArgs)
	}
	return 0, nil
}

func completePendingTask(db *store.DB, taskID int64) {
	if CompletePendingTaskFn != nil {
		CompletePendingTaskFn(db, taskID)
	}
}

func failPendingTask(db *store.DB, taskID int64, reason string) {
	if FailPendingTaskFn != nil {
		FailPendingTaskFn(db, taskID, reason)
	}
}

func requireOnline() {
	if RequireOnlineFn != nil {
		RequireOnlineFn()
	}
}

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

func PrintRemediationSummaryNoPrompt(items []RemediationItem) {
	if PrintRemediationSummaryNoPromptFn != nil {
		PrintRemediationSummaryNoPromptFn(items)
	}
}

func PrintRemediationSummaryAutoFix(items []RemediationItem) {
	if PrintRemediationSummaryAutoFixFn != nil {
		PrintRemediationSummaryAutoFixFn(items)
	}
}

func PrintRemediationSummary(items []RemediationItem) {
	if PrintRemediationSummaryFn != nil {
		PrintRemediationSummaryFn(items)
	}
}

func ApplyTransportFlag(dir string, useSSH, useHTTPS bool) (bool, string, string, error) {
	if ApplyTransportFlagFn != nil {
		return ApplyTransportFlagFn(dir, useSSH, useHTTPS)
	}
	return false, "", "", nil
}

func buildCommandArgs(args []string) string {
	return strings.Join(args, " ")
}

func printCanonicalCmdBanner(canonical, alias string) {
	fmt.Fprintf(os.Stderr, "  → Running: gitmap %s  (alias: %s)\n\n", canonical, alias)
}

func stringsEqualAbs(p1, p2 string) bool {
	return strings.EqualFold(filepath.Clean(p1), filepath.Clean(p2))
}

func stripANSI(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			i = j
			continue
		}
		out.WriteByte(s[i])
	}
	return out.String()
}
