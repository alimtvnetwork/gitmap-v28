package heatmap

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TrackAgentAction records a discrete agent lifecycle event into the root heatmap database.
func TrackAgentAction(taskSlug, actionType, summary string) error {
	cleanSlug := strings.TrimSpace(taskSlug)
	cleanAction := strings.TrimSpace(actionType)
	hasMissingInput := cleanSlug == "" || cleanAction == ""
	if hasMissingInput {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.cwd")
	}

	return recordAgentActionAtRoot(cwd, cleanSlug, cleanAction, summary)
}

func recordAgentActionAtRoot(cwd, taskSlug, actionType, summary string) error {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		TaskSlug:       taskSlug,
		ActionType:     actionType,
		ReasoningNotes: summary,
		ReasoningState: "EXPLORATION",
	})

	return logErr
}
