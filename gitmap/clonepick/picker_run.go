package clonepick

// picker_run.go: bubbletea program lifecycle for the --ask picker.
// Split out of picker.go so the model + key-handler file stays
// under the strict 200-line cap.

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// RunPicker enumerates plan.RepoUrl, opens the bubbletea picker, and
// returns the user-picked subset. plan.Paths seeds the initial
// selection so re-running with the same args is a no-op confirmation.
func RunPicker(plan Plan) ([]string, error) {
	picked, tmp, err := RunPickerKeep(plan)
	if len(tmp) > 0 {
		os.RemoveAll(tmp)
	}

	return picked, err
}

// RunPickerKeep is the clone-once variant: returns the picked paths
// AND the temp metadata-clone directory so the executor can promote
// it instead of re-cloning. The caller owns tmp and must remove it
// (or pass it to the executor via Plan.PreClonedSrc, which moves it
// into place). On error or cancellation tmp is already cleaned up.
func RunPickerKeep(plan Plan) ([]string, string, error) {
	all, tmp, err := ListRepoPathsKeep(plan)
	if err != nil {
		return nil, "", err
	}
	model := newPickerModel(all, plan.Paths)
	prog := tea.NewProgram(model)
	final, runErr := prog.Run()
	if runErr != nil {
		os.RemoveAll(tmp)

		return nil, "", fmt.Errorf("clone-pick: picker run: %w", runErr)
	}
	finished, _ := final.(pickerModel)
	if finished.cancelled {
		os.RemoveAll(tmp)

		return nil, "", ErrPickerCancelled
	}

	return finished.selected(), tmp, nil
}
