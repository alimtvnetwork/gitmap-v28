package clonepick

// picker_run.go: bubbletea program lifecycle for the --ask picker.
// Split out of picker.go so the model + key-handler file stays
// under the strict 200-line cap.

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunPicker enumerates plan.RepoUrl, opens the bubbletea picker, and
// returns the user-picked subset. plan.Paths seeds the initial
// selection so re-running with the same args is a no-op confirmation.
func RunPicker(plan Plan) ([]string, error) {
	res := RunPickerWrap(plan)
	if res.HasError() {
		return nil, res.Err
	}
	return res.Items(), nil
}

func RunPickerWrap(plan Plan) StringSliceWrap {
	picked, tmp, err := RunPickerKeep(plan)
	if len(tmp) > 0 {
		os.RemoveAll(tmp)
	}
	if err != nil {
		return result.FailSlice[string](apperror.WrapSimple(err, "RunPickerWrap"))
	}
	return result.OkSlice(picked)
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

		return nil, "", apperror.WrapSimple(runErr, "clone-pick: picker run")
	}

	finished, _ := final.(pickerModel)
	if finished.isCanceled {
		os.RemoveAll(tmp)

		return nil, "", ErrPickerCancelled
	}

	return finished.selected(), tmp, nil
}
