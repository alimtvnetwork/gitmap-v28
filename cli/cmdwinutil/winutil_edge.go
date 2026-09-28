package cmdwinutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	edgeRegBlockerKey = `HKLM\SOFTWARE\Microsoft\EdgeUpdate`
	edgeRegBlockerVal = "DoNotUpdateToEdgeWithChromium"
	edgeAppxName      = "*Microsoft.MicrosoftEdge*"
	edgeDevToolsAppx  = "*MicrosoftEdgeDevToolsClient*"
)

// RunEdgeUninstall executes Chris Titus WinUtil Edge removal sequence.
func RunEdgeUninstall(hasKeepWebView2, isDryRun bool) (WinRemovalResult, error) {
	if runtime.GOOS != "windows" {
		return makeNonWindowsResult("edge"), nil
	}
	res := WinRemovalResult{Target: "edge", IsDryRun: isDryRun, IsSuccess: true}
	terminateEdgeProcesses(hasKeepWebView2, isDryRun, &res)
	setupPath := findEdgeSetupPath(hasKeepWebView2)
	executeEdgeSetup(setupPath, isDryRun, &res)
	removeEdgeAppxPackages(isDryRun, &res)
	disableEdgeServices(isDryRun, &res)
	applyEdgeRegistryBlocker(isDryRun, &res)
	return res, nil
}

func terminateEdgeProcesses(hasKeepWebView2, isDryRun bool, res *WinRemovalResult) {
	killEdgeProcess("msedge.exe", isDryRun, res)
	if !hasKeepWebView2 {
		killEdgeProcess("msedgewebview2.exe", isDryRun, res)
		return
	}
	res.Actions = append(res.Actions, "Preserved WebView2 runtime processes")
}

func killEdgeProcess(procName string, isDryRun bool, res *WinRemovalResult) {
	if isDryRun {
		res.Actions = append(res.Actions, "[dry-run] Would terminate process: "+procName)
		return
	}
	_ = exec.Command("taskkill", "/F", "/IM", procName).Run()
	res.Actions = append(res.Actions, "Terminated process: "+procName)
}

func findEdgeSetupPath(hasKeepWebView2 bool) string {
	for _, base := range getEdgeSetupDirs(hasKeepWebView2) {
		if found := findSetupInDir(base); len(found) > 0 {
			return found
		}
	}
	return ""
}

func getEdgeSetupDirs(hasKeepWebView2 bool) []string {
	dirs := []string{
		`C:\Program Files (x86)\Microsoft\Edge\Application`,
		`C:\Program Files (x86)\Microsoft\EdgeCore`,
		`C:\Program Files\Microsoft\Edge\Application`,
		`C:\Program Files\Microsoft\EdgeCore`,
	}
	if !hasKeepWebView2 {
		dirs = append(dirs, `C:\Program Files (x86)\Microsoft\EdgeWebView\Application`)
	}
	return dirs
}

func findSetupInDir(baseDir string) string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		sub := filepath.Join(baseDir, e.Name(), "Installer", "setup.exe")
		if info, sErr := os.Stat(sub); sErr == nil && !info.IsDir() {
			return sub
		}
	}
	return ""
}

func executeEdgeSetup(setupPath string, isDryRun bool, res *WinRemovalResult) {
	if len(setupPath) == 0 {
		res.Warnings = append(res.Warnings, "Edge installer setup.exe not found on system")
		return
	}
	if isDryRun {
		recordEdgeSetup("[dry-run] Would execute: "+setupPath, res)
		return
	}
	runEdgeSetupCommand(setupPath, res)
}

func recordEdgeSetup(action string, res *WinRemovalResult) {
	res.Actions = append(res.Actions, action)
	res.ItemsRemoved++
}

func runEdgeSetupCommand(setupPath string, res *WinRemovalResult) {
	cmd := exec.Command(setupPath, "--uninstall", "--system-level", "--verbose-logging", "--force-uninstall")
	if out, err := cmd.CombinedOutput(); err != nil {
		res.Warnings = append(res.Warnings, fmt.Sprintf("Setup note: %s", string(out)))
		return
	}
	recordEdgeSetup("Executed Edge force uninstall via "+setupPath, res)
}

func removeEdgeAppxPackages(isDryRun bool, res *WinRemovalResult) {
	removeSingleAppxPackage(edgeAppxName, isDryRun, res)
	removeSingleAppxPackage(edgeDevToolsAppx, isDryRun, res)
}

func removeSingleAppxPackage(pkgFilter string, isDryRun bool, res *WinRemovalResult) {
	if isDryRun {
		res.Actions = append(res.Actions, "[dry-run] Would remove Appx: "+pkgFilter)
		res.ItemsRemoved++
		return
	}
	cmdStr := fmt.Sprintf("Get-AppxPackage -AllUsers %s | Remove-AppxPackage -AllUsers", pkgFilter)
	_, _ = runPowerShellScript(cmdStr)
	res.Actions = append(res.Actions, "Removed Appx: "+pkgFilter)
	res.ItemsRemoved++
}

func disableEdgeServices(isDryRun bool, res *WinRemovalResult) {
	disableSingleService("edgeupdate", isDryRun, res)
	disableSingleService("edgeupdatem", isDryRun, res)
}

func disableSingleService(svcName string, isDryRun bool, res *WinRemovalResult) {
	if isDryRun {
		res.Actions = append(res.Actions, "[dry-run] Would disable service: "+svcName)
		return
	}
	_ = exec.Command("sc.exe", "stop", svcName).Run()
	_ = exec.Command("sc.exe", "config", svcName, "start=disabled").Run()
	res.Actions = append(res.Actions, "Disabled service: "+svcName)
}

func applyEdgeRegistryBlocker(isDryRun bool, res *WinRemovalResult) {
	setRegistryDWord(edgeRegBlockerKey, edgeRegBlockerVal, 1, isDryRun, res)
}
