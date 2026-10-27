package cmdswap

import (
	"fmt"
	"strings"
)

// GenerateWindowsPowerShellScript builds PowerShell commands to terminate targets, backup, and swap binary.
func GenerateWindowsPowerShellScript(opts SwapOptions, pids []int) string {
	var sb strings.Builder

	sb.WriteString("# PowerShell binary swap script\n")
	sb.WriteString("$ErrorActionPreference = 'Stop'\n\n")

	writePowerShellKillSection(&sb, pids)
	writePowerShellBackupSection(&sb, opts.ToPath)
	writePowerShellCopySection(&sb, opts.FromPath, opts.ToPath)
	writePowerShellRestartSection(&sb, opts.ToPath, opts.HasRestart)

	return sb.String()
}

// GenerateUnixShellScript builds POSIX shell script to terminate targets, backup, and swap binary.
func GenerateUnixShellScript(opts SwapOptions, pids []int) string {
	var sb strings.Builder

	sb.WriteString("#!/bin/sh\n")
	sb.WriteString("set -e\n\n")

	writeUnixKillSection(&sb, pids)
	writeUnixBackupSection(&sb, opts.ToPath)
	writeUnixCopySection(&sb, opts.FromPath, opts.ToPath)
	writeUnixRestartSection(&sb, opts.ToPath, opts.HasRestart)

	return sb.String()
}

func writePowerShellKillSection(sb *strings.Builder, pids []int) {
	for _, pid := range pids {
		if pid <= 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("Stop-Process -Id %d -Force -ErrorAction SilentlyContinue\n", pid))
	}

	sb.WriteString("Start-Sleep -Milliseconds 200\n\n")
}

func writePowerShellBackupSection(sb *strings.Builder, toPath string) {
	sb.WriteString(fmt.Sprintf("if (Test-Path -LiteralPath '%s') {\n", toPath))
	sb.WriteString(fmt.Sprintf("    Copy-Item -LiteralPath '%s' -Destination '%s.old' -Force -ErrorAction SilentlyContinue\n", toPath, toPath))
	sb.WriteString("}\n\n")
}

func writePowerShellCopySection(sb *strings.Builder, fromPath, toPath string) {
	sb.WriteString("$maxRetries = 10\n")
	sb.WriteString("$isCopied = $false\n")
	sb.WriteString("for ($i = 0; $i -lt $maxRetries; $i++) {\n")
	sb.WriteString("    try {\n")
	sb.WriteString(fmt.Sprintf("        Copy-Item -LiteralPath '%s' -Destination '%s' -Force\n", fromPath, toPath))
	sb.WriteString("        $isCopied = $true\n")
	sb.WriteString("        break\n")
	sb.WriteString("    } catch {\n")
	sb.WriteString("        Start-Sleep -Milliseconds 300\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	sb.WriteString("if (-not $isCopied) {\n")
	sb.WriteString("    Write-Error 'Failed to swap binary after retries'\n")
	sb.WriteString("    exit 1\n")
	sb.WriteString("}\n\n")
}

func writePowerShellRestartSection(sb *strings.Builder, toPath string, hasRestart bool) {
	if !hasRestart {
		return
	}

	sb.WriteString(fmt.Sprintf("& '%s' --version\n", toPath))
}

func writeUnixKillSection(sb *strings.Builder, pids []int) {
	for _, pid := range pids {
		if pid <= 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("kill -9 %d 2>/dev/null || true\n", pid))
	}

	sb.WriteString("sleep 0.2\n\n")
}

func writeUnixBackupSection(sb *strings.Builder, toPath string) {
	sb.WriteString(fmt.Sprintf("if [ -f \"%s\" ]; then\n", toPath))
	sb.WriteString(fmt.Sprintf("    cp -f \"%s\" \"%s.old\" 2>/dev/null || true\n", toPath, toPath))
	sb.WriteString("fi\n\n")
}

func writeUnixCopySection(sb *strings.Builder, fromPath, toPath string) {
	sb.WriteString(fmt.Sprintf("cp -f \"%s\" \"%s\"\n", fromPath, toPath))
	sb.WriteString(fmt.Sprintf("chmod +x \"%s\"\n\n", toPath))
}

func writeUnixRestartSection(sb *strings.Builder, toPath string, hasRestart bool) {
	if !hasRestart {
		return
	}

	sb.WriteString(fmt.Sprintf("\"%s\" --version\n", toPath))
}
