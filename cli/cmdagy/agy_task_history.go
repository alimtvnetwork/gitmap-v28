// Package cmdagy — agy_task_history.go logs AGY mutations to TaskHistory and outputs undo guidance.
package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func recordAgyTask(action, target, forwardPayload, inversePayload string) error {
	db, err := store.OpenSectionTasksDB("agy", "")
	if err != nil {
		return err
	}
	defer db.Close()

	taskID := fmt.Sprintf("agy-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC3339)
	queryHist := `INSERT INTO TaskHistory (TaskId, Section, Action, Target, ForwardPayload, InversePayload, CreatedAt) VALUES (?, 'agy', ?, ?, ?, ?, ?)`
	resHist := store.ExecWrapper(db.Conn(), queryHist, taskID, action, target, forwardPayload, inversePayload, now)
	if resHist.IsFailure {
		return apperror.WrapSimple(resHist.Error, "recordAgyTask")
	}

	return nil
}

func printAgyUndoGuidance() {
	fmt.Printf("\n  %sℹ%s Files on disk preserved: project files and git repository are not deleted.\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %sℹ%s Undo this action anytime with: %sgitmap agy undo%s (or %sgitmap undo%s)\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorGreen, constants.ColorReset,
		constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  %sℹ%s Review task history with: %sgitmap agy history%s (or %sgitmap task history%s)\n\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorYellow, constants.ColorReset,
		constants.ColorYellow, constants.ColorReset)
}
