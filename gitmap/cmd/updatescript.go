package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// executeUpdate writes a temp script and runs it.
// On Windows it uses PowerShell; on Linux/macOS it uses run.sh directly.
func executeUpdate(repoPath string, report reportErrorsConfig) {
	if runtime.GOOS == "windows" {
		executeUpdateWindows(repoPath, report)

		return
	}

	executeUpdateUnix(repoPath, report)
}

// executeUpdateWindows writes a temp PS1 script and runs it.
func executeUpdateWindows(repoPath string, report reportErrorsConfig) {
	scriptPath, err := writeUpdateScript(repoPath)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.updatescript.writeWindows",
			"E1137",
			"failed to write Windows update script",
			"cmd.updatescript",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"repoPath": repoPath},
		)
		cliexit.HandleError(appErr, 1)
		return
	}
	defer os.Remove(scriptPath)

	log := verbose.Get()
	if log != nil {
		log.Log(constants.UpdateScriptLogExec, scriptPath)
	}

	runUpdateScript(scriptPath, report)
}

// executeUpdateUnix runs run.sh --update with the install path as deploy target.
func executeUpdateUnix(repoPath string, report reportErrorsConfig) {
	runSH := filepath.Join(repoPath, "run.sh")

	if _, err := os.Stat(runSH); err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.updatescript.statUnix",
			"E1138",
			fmt.Sprintf("run.sh not found in repo path '%s'", repoPath),
			"cmd.updatescript",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"repoPath": repoPath},
		)
		cliexit.HandleError(appErr, 1)
		return
	}

	log := verbose.Get()
	if log != nil {
		log.Log(constants.UpdateScriptLogExec, runSH)
	}

	// Resolve the active binary's installed directory.
	installDir := resolveInstalledDir()
	fmt.Printf(constants.MsgUpdateInstallDir, installDir)

	args := []string{runSH, "--update"}

	cmd := exec.Command("bash", args...)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = report.applyToEnv(os.Environ())

	err := cmd.Run()

	logScriptResult(err)

	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.updatescript.runUnix",
			"E1139",
			"failed to execute Unix update script",
			"cmd.updatescript",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
	}
}

// resolveInstalledDir returns the directory where the active gitmap binary lives.
func resolveInstalledDir() string {
	// First try: which gitmap on PATH
	path, err := exec.LookPath("gitmap")
	resolved := ""
	evalErr := error(nil)
	if err == nil {
		resolved, evalErr = filepath.EvalSymlinks(path)
	}
	if err == nil && evalErr == nil {
		return filepath.Dir(resolved)
	}
	if err == nil {
		return filepath.Dir(path)
	}

	// Fallback: current executable's directory
	selfPath, errExec := os.Executable()
	if errExec != nil {
		return ""
	}

	resolvedSelf, errSym := filepath.EvalSymlinks(selfPath)
	if errSym != nil {
		return filepath.Dir(selfPath)
	}

	return filepath.Dir(resolvedSelf)
}

// writeUpdateScript creates a temporary PowerShell script for self-update.
// Writes with UTF-8 BOM so PowerShell correctly handles Unicode characters.
func writeUpdateScript(repoPath string) (string, error) {
	runPS1 := filepath.Join(repoPath, "run.ps1")
	script := buildUpdateScript(repoPath, runPS1)

	return writeScriptToTemp(script)
}

// writeScriptToTemp writes script content to a temp file with UTF-8 BOM.
func writeScriptToTemp(script string) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), constants.UpdateScriptGlob)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := tmpFile.Write(bom); err != nil {
		return "", err
	}
	if _, err := tmpFile.WriteString(script); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// buildUpdateScript generates the PowerShell script content.
func buildUpdateScript(repoPath, runPS1 string) string {
	// Build the PowerShell @("a","b") literal of all known app subdirs
	// (current + legacy) from the embedded deploy manifest.
	knownSubdirs := buildPSSubdirArray()

	return fmt.Sprintf(constants.UpdatePSHeader, repoPath) +
		fmt.Sprintf(constants.UpdatePSDeployDetect,
			repoPath,
			constants.GitMapSubdir,
			constants.GitMapCliSubdir,
			constants.Manifest.BinaryName.Windows,
			knownSubdirs,
		) +
		constants.UpdatePSVersionBefore +
		fmt.Sprintf(constants.UpdatePSRunUpdate, runPS1) +
		constants.UpdatePSSync +
		constants.UpdatePSVersionAfter +
		fmt.Sprintf(constants.UpdatePSVerify, repoPath, repoPath) +
		constants.UpdatePSPostActions
}

// buildPSSubdirArray returns a PowerShell array literal — e.g. @("gitmap-cli","gitmap")
// — containing the current AppSubdir plus every LegacyAppSubdirs entry
// from gitmap/constants/deploy-manifest.json.
func buildPSSubdirArray() string {
	all := append([]string{constants.GitMapCliSubdir}, constants.LegacyAppSubdirs...)
	quoted := make([]string, 0, len(all))
	seen := map[string]bool{}
	for _, name := range all {
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		quoted = append(quoted, fmt.Sprintf("%q", name))
	}

	return "@(" + joinComma(quoted) + ")"
}

// joinComma joins strings with a comma — small helper to avoid pulling strings.Join.
func joinComma(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ","
		}
		out += p
	}

	return out
}

// runUpdateScript executes the PowerShell script with output piped to terminal.
func runUpdateScript(scriptPath string, report reportErrorsConfig) error {
	cmd := exec.Command(constants.PSBin, constants.PSExecPolicy, constants.PSBypass,
		constants.PSNoProfile, constants.PSNoLogo, constants.PSFile, scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = report.applyToEnv(os.Environ())
	err := cmd.Run()

	logScriptResult(err)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.updatescript.runScript",
			"E1140",
			"failed to run update script",
			"cmd.updatescript",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"scriptPath": scriptPath},
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	return nil
}

// logScriptResult logs the update script exit status if verbose is active.
func logScriptResult(err error) {
	log := verbose.Get()
	if log != nil {
		log.Log(constants.UpdateScriptLogExit, err)
	}
}
