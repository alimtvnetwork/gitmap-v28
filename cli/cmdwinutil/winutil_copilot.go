package cmdwinutil

import (
	"fmt"
	"os/exec"
	"runtime"
)

const (
	hkcuCopilotPolicyKey = `HKCU\Software\Policies\Microsoft\Windows\WindowsCopilot`
	hklmCopilotPolicyKey = `HKLM\Software\Policies\Microsoft\Windows\WindowsCopilot`
	hkcuExplorerKey      = `HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`
	copilotPolicyVal     = "TurnOffWindowsCopilot"
	copilotButtonVal     = "ShowCopilotButton"
	copilotAppxFilter    = "*Microsoft.Windows.Copilot*"
)

// RunCopilotUninstall removes Copilot Appx packages and injects disable policies.
func RunCopilotUninstall(isDryRun bool) (WinRemovalResult, error) {
	if runtime.GOOS != "windows" {
		return makeNonWindowsResult("copilot"), nil
	}
	res := WinRemovalResult{Target: "copilot", IsDryRun: isDryRun, IsSuccess: true}
	removeCopilotAppx(isDryRun, &res)
	applyCopilotGroupPolicy(isDryRun, &res)
	applyCopilotTaskbarPolicy(isDryRun, &res)
	return res, nil
}

func makeNonWindowsResult(target string) WinRemovalResult {
	return WinRemovalResult{
		Target:    target,
		IsSuccess: true,
		Actions:   []string{fmt.Sprintf("%s removal is only supported on Windows systems", target)},
	}
}

func removeCopilotAppx(isDryRun bool, res *WinRemovalResult) {
	if isDryRun {
		recordCopilotAppx("[dry-run] Would remove Appx package: "+copilotAppxFilter, res)
		return
	}
	cmdStr := fmt.Sprintf("Get-AppxPackage -AllUsers %s | Remove-AppxPackage -AllUsers", copilotAppxFilter)
	if _, err := runPowerShellScript(cmdStr); err != nil {
		res.Warnings = append(res.Warnings, fmt.Sprintf("Appx removal note: %v", err))
		return
	}
	recordCopilotAppx("Removed Appx package: "+copilotAppxFilter, res)
}

func recordCopilotAppx(action string, res *WinRemovalResult) {
	res.Actions = append(res.Actions, action)
	res.ItemsRemoved++
}

func applyCopilotGroupPolicy(isDryRun bool, res *WinRemovalResult) {
	setRegistryDWord(hkcuCopilotPolicyKey, copilotPolicyVal, 1, isDryRun, res)
	setRegistryDWord(hklmCopilotPolicyKey, copilotPolicyVal, 1, isDryRun, res)
}

func applyCopilotTaskbarPolicy(isDryRun bool, res *WinRemovalResult) {
	setRegistryDWord(hkcuExplorerKey, copilotButtonVal, 0, isDryRun, res)
}

func setRegistryDWord(fullKey, valName string, val uint32, isDryRun bool, res *WinRemovalResult) {
	if isDryRun {
		res.Actions = append(res.Actions, fmt.Sprintf("[dry-run] Would set %s\\%s = %d", fullKey, valName, val))
		return
	}
	valStr := fmt.Sprintf("%d", val)
	cmd := exec.Command("reg", "add", fullKey, "/v", valName, "/t", "REG_DWORD", "/d", valStr, "/f")
	if out, err := cmd.CombinedOutput(); err != nil {
		res.Warnings = append(res.Warnings, fmt.Sprintf("Reg warning (%s): %s", valName, string(out)))
		return
	}
	res.Actions = append(res.Actions, fmt.Sprintf("Set registry %s\\%s = %d", fullKey, valName, val))
}

func runPowerShellScript(script string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
