package cmd

import (
	"errors"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// checkVSCodeProjectManager verifies the alefragnani.project-manager sync
// target is reachable. Three outcomes:
//
//   - Both user-data root AND extension storage dir present  -> OK + path.
//   - User-data root missing                                 -> Warn (VS Code
//     not detected — sync will be skipped on next scan).
//   - Root present but extension dir missing                 -> Warn (extension
//     not installed — sync will be skipped).
//
// Doctor checks must NEVER hard-fail on missing optional integrations, so
// this returns 0 on warn paths and only counts the explicit "broken"
// scenario (e.g. corrupt projects.json) as an issue.
func checkVSCodeProjectManager() int {
	path, err := vscodepm.ProjectsJSONPath()

	if err == nil {
		printOK(constants.DoctorVSCodePMOKFmt, path)

		return 0
	}

	if errors.Is(err, vscodepm.ErrUserDataMissing) {
		printWarn(constants.DoctorVSCodePMNoVSCode)

		return 0
	}

	if errors.Is(err, vscodepm.ErrExtensionMissing) {
		printWarn(constants.DoctorVSCodePMNoExtension)

		return 0
	}

	printIssue(constants.DoctorVSCodePMUnknownTitle, err.Error())

	return 1
}
