package cmdinstall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func printInstallFailureDetails(
	tool, manager, version string,
	args []string, err error, logPath string,
) {
	versionLabel := resolveDisplayVersion(version)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailed, tool)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedVersion, versionLabel)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedManager, manager)
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedCmd, strings.Join(args, " "))
	fmt.Fprintf(os.Stderr, constants.ErrInstallFailedReason, err)

	if logPath != "" {
		fmt.Fprintf(os.Stderr, constants.ErrInstallFailedLog, logPath)
		fmt.Fprint(os.Stderr, constants.ErrInstallFailedHint)
	}
}

func resolveDisplayVersion(version string) string {
	if version != "" {
		return version
	}

	return "latest"
}

func writeInstallErrorLog(
	tool, manager, version string,
	args []string, output []byte, installErr error,
) string {
	logDir := constants.InstallLogDir

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not create log directory %s: %v\n", logDir, err)

		return ""
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("%s-error-%s.log", tool, time.Now().Format("2006-01-02_15-04-05")))
	content := buildInstallErrorLogContent(tool, manager, version, args, output, installErr)

	if err := os.WriteFile(logPath, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not write error log to %s: %v\n", logPath, err)

		return ""
	}

	return logPath
}

func buildInstallErrorLogContent(
	tool, manager, version string,
	args []string, output []byte, installErr error,
) string {
	versionLabel := resolveDisplayVersion(version)
	header := fmt.Sprintf("gitmap install error log\n========================\n\nTool:            %s\nVersion:         %s\nPackage Manager: %s\nCommand:         %s\nTimestamp:       %s\nError:           %v\n\n--- Installer Output ---\n\n",
		tool, versionLabel, manager, strings.Join(args, " "), time.Now().Format(time.RFC3339), installErr)

	if len(output) > 0 {
		return header + string(output)
	}

	return header + "(no output captured — verbose mode pipes directly to stdout/stderr)\n"
}
