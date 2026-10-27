package heatmap

import (
	"io"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunHeatmapFailures queries and renders recent expectation failures.
func RunHeatmapFailures(isRc bool, flags map[string]string, out io.Writer) error {
	query := "SELECT HeatmapActionLogId, TaskSlug, ActionType, FilePath, ExpectedOutcome, ActualOutcome FROM HeatmapActionLog WHERE IsExpectationMet = 0 ORDER BY HeatmapActionLogId DESC LIMIT 25"

	return executeShortcutQuery(isRc, flags, query, out)
}

// RunHeatmapHotfiles queries and renders top touched files and read/edit ratios.
func RunHeatmapHotfiles(isRc bool, flags map[string]string, out io.Writer) error {
	query := "SELECT FilePath, COUNT(*) AS Touches, SUM(CASE WHEN ActionType = 'READ' THEN 1 ELSE 0 END) AS Reads, SUM(CASE WHEN ActionType = 'EDIT' THEN 1 ELSE 0 END) AS Edits FROM HeatmapActionLog GROUP BY FilePath ORDER BY Touches DESC LIMIT 20"

	return executeShortcutQuery(isRc, flags, query, out)
}

// RunHeatmapRecent queries and renders recent action logs.
func RunHeatmapRecent(isRc bool, flags map[string]string, out io.Writer) error {
	query := "SELECT HeatmapActionLogId, TaskSlug, ActionType, FilePath, ReasoningState, ExpectedOutcome, ActualOutcome, CreatedAt FROM HeatmapActionLog ORDER BY HeatmapActionLogId DESC LIMIT 20"

	return executeShortcutQuery(isRc, flags, query, out)
}

// RunHeatmapTasks queries and renders registered parent tasks.
func RunHeatmapTasks(isRc bool, flags map[string]string, out io.Writer) error {
	query := "SELECT TaskSlug, TaskName, Status, CreatedAt, UpdatedAt FROM HeatmapParentTask ORDER BY UpdatedAt DESC LIMIT 20"

	return executeShortcutQuery(isRc, flags, query, out)
}

func executeShortcutQuery(isRc bool, flags map[string]string, query string, out io.Writer) error {
	dbPath := resolveQueryDbPath(isRc, flags, query)

	start := time.Now()
	cols, rows, err := store.ExecuteReadOnlyQuery(dbPath, query)
	durationMs := float64(time.Since(start).Microseconds()) / 1000.0
	if err != nil {
		return apperror.WrapSimple(err, "heatmap_shortcut.exec")
	}

	return renderQueryResult(dbPath, query, cols, rows, durationMs, resolveFormat(flags), out)
}
