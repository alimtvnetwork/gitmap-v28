package cmdstatus

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	statusconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RunStatus handles the "status" subcommand.
func RunStatus(args []string) error {
	checkHelp(statusconstants.CmdStatus, args)
	groupName, isAll, isOnlyDirty := parseStatusFlags(args)
	records, err := loadStatusByScope(groupName, isAll)
	if err != nil {
		return err
	}

	return renderStatusFiltered(records, isOnlyDirty)
}

func renderStatusFiltered(records []model.ScanRecord, isOnlyDirty bool) error {
	if isOnlyDirty {
		records = filterDirtyRecords(records)
	}

	if isOnlyDirty && len(records) == 0 {
		fmt.Println(statusconstants.StatusCleanSuccessMsg)

		return nil
	}

	return renderStatusView(records)
}

func renderStatusView(records []model.ScanRecord) error {
	printStatusBanner(len(records))
	prog := cloner.NewBatchProgress(len(records), statusconstants.StatusBatchProgressLabel, true)
	summary := printStatusTableTracked(records, prog)
	printStatusSummary(summary)

	return nil
}

func filterDirtyRecords(records []model.ScanRecord) []model.ScanRecord {
	var dirty []model.ScanRecord
	for _, rec := range records {
		rs := gitutil.Status(rec.AbsolutePath)
		if rs.Dirty || rs.Ahead > 0 || rs.Behind > 0 || rs.StashCount > 0 {
			dirty = append(dirty, rec)
		}
	}

	return dirty
}

// parseStatusFlags parses --group, --all, and --dirty flags.
func parseStatusFlags(args []string) (string, bool, bool) {
	fs := flag.NewFlagSet(statusconstants.CmdStatus, flag.ExitOnError)
	gFlag := fs.String(statusconstants.FlagGroup, "", constants.FlagDescGroup)
	fs.StringVar(gFlag, statusconstants.FlagGroupShort, "", constants.FlagDescGroup)
	aFlag := fs.Bool(statusconstants.FlagAll, false, constants.FlagDescAll)
	dFlag := fs.Bool(statusconstants.FlagDirty, false, statusconstants.FlagDescDirty)
	fs.BoolVar(dFlag, statusconstants.FlagOnlyDirty, false, statusconstants.FlagDescDirty)
	fs.Parse(args)

	isAll := *aFlag
	isOnlyDirty := *dFlag

	return *gFlag, isAll, isOnlyDirty
}

// statusSummary aggregates counts across all repos.
type statusSummary struct {
	Total   int
	Clean   int
	Dirty   int
	Ahead   int
	Behind  int
	Stashed int
	Missing int
}
