package cmdos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"
)

// GetLocalOSInfo gathers OS metadata for the currently running machine.
func GetLocalOSInfo() OSInfoReport {
	hostname, _ := os.Hostname()
	probe := probeLocalPlatformOS()
	arch := runtime.GOARCH
	gitPath, hasGit := probeGitTool()
	bashPath, hasBash := probeBashTool()
	psPath, hasPS := probePowerShellTool()

	return OSInfoReport{
		OSType:         probe.OSType,
		OSGroup:        probe.OSGroup,
		OSVersion:      probe.OSVersion,
		BuildVersion:   probe.BuildVersion,
		Architecture:   arch,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, arch),
		Hostname:       hostname,
		NumCPU:         runtime.NumCPU(),
		Kernel:         probe.Kernel,
		GitPath:        gitPath,
		BashPath:       bashPath,
		PowerShellPath: psPath,
		HasGit:         hasGit,
		HasBash:        hasBash,
		HasPowerShell:  hasPS,
	}
}

func findFirstExistingFile(candidates []string) (string, bool) {
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, true
		}
	}
	return "", false
}

func probeGitTool() (string, bool) {
	if p, err := exec.LookPath("git"); err == nil && p != "" {
		return p, true
	}
	if runtime.GOOS != "windows" {
		return "", false
	}
	return findFirstExistingFile([]string{
		`C:\Program Files\Git\cmd\git.exe`,
		`C:\Program Files\Git\bin\git.exe`,
	})
}

func probeWindowsBash() (string, bool) {
	if runtime.GOOS != "windows" {
		return "", false
	}
	return findFirstExistingFile([]string{
		`C:\Program Files\Git\bin\bash.exe`,
		`C:\Program Files\Git\usr\bin\bash.exe`,
	})
}

func probeBashTool() (string, bool) {
	if c, found := probeWindowsBash(); found {
		return c, true
	}
	if p, err := exec.LookPath("bash"); err == nil && p != "" {
		return p, true
	}
	return "", false
}

func probePowerShellTool() (string, bool) {
	if p, err := exec.LookPath("pwsh"); err == nil && p != "" {
		return p, true
	}
	if p, err := exec.LookPath("powershell"); err == nil && p != "" {
		return p, true
	}
	return "", false
}

// RenderOSInfo formats and writes OS profile to output.
func RenderOSInfo(out io.Writer, report OSInfoReport, isJSON bool) error {
	if isJSON {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FIELD\tVALUE")
	renderBaseOSFields(w, report)
	renderToolFields(w, report)
	return w.Flush()
}

func renderBaseOSFields(w io.Writer, report OSInfoReport) {
	fmt.Fprintf(w, "OS Type\t%s\n", report.OSType)
	fmt.Fprintf(w, "OS Group\t%s\n", report.OSGroup)
	fmt.Fprintf(w, "OS Version\t%s\n", report.OSVersion)
	if report.BuildVersion != "" {
		fmt.Fprintf(w, "Build Version\t%s\n", report.BuildVersion)
	}
	fmt.Fprintf(w, "Architecture\t%s\n", report.Architecture)
	fmt.Fprintf(w, "Platform\t%s\n", report.Platform)
	fmt.Fprintf(w, "Hostname\t%s\n", report.Hostname)
	fmt.Fprintf(w, "CPUs\t%d\n", report.NumCPU)
	if report.Kernel != "" {
		fmt.Fprintf(w, "Kernel\t%s\n", report.Kernel)
	}
}

func renderToolFields(w io.Writer, report OSInfoReport) {
	fmt.Fprintf(w, "Git Path\t%s\n", formatToolDisplay(report.GitPath, report.HasGit))
	fmt.Fprintf(w, "Bash Path\t%s\n", formatToolDisplay(report.BashPath, report.HasBash))
	fmt.Fprintf(w, "PowerShell Path\t%s\n", formatToolDisplay(report.PowerShellPath, report.HasPowerShell))
}

func formatToolDisplay(path string, isInstalled bool) string {
	if !isInstalled || path == "" {
		return "not installed"
	}
	return path
}
