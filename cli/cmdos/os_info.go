package cmdos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"text/tabwriter"
)

// GetLocalOSInfo gathers OS metadata for the currently running machine.
func GetLocalOSInfo() OSInfoReport {
	hostname, _ := os.Hostname()
	osType, osVer, kernel := probeLocalPlatformOS()
	arch := runtime.GOARCH
	platform := fmt.Sprintf("%s/%s", runtime.GOOS, arch)

	return OSInfoReport{
		OSType:       osType,
		OSVersion:    osVer,
		Architecture: arch,
		Hostname:     hostname,
		NumCPU:       runtime.NumCPU(),
		Kernel:       kernel,
		Platform:     platform,
	}
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
	fmt.Fprintf(w, "OS Type\t%s\n", report.OSType)
	fmt.Fprintf(w, "OS Version\t%s\n", report.OSVersion)
	fmt.Fprintf(w, "Architecture\t%s\n", report.Architecture)
	fmt.Fprintf(w, "Platform\t%s\n", report.Platform)
	fmt.Fprintf(w, "Hostname\t%s\n", report.Hostname)
	fmt.Fprintf(w, "CPUs\t%d\n", report.NumCPU)
	if report.Kernel != "" {
		fmt.Fprintf(w, "Kernel\t%s\n", report.Kernel)
	}
	return w.Flush()
}
