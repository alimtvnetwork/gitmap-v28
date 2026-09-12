package cmd

import (
	"strings"
	"testing"
)

func TestBuildUpdateScriptUsesPathAwareDeployVerification(t *testing.T) {
	script := buildUpdateScript(`D:\src\gitmap`, `D:\src\gitmap\run.ps1`)

	checks := []string{
		`$configDeployedBinary = $null`,
		`$activeCmdForDeploy = Get-Command gitmap -ErrorAction SilentlyContinue | Select-Object -First 1`,
		`$deployedBinary = Join-Path $effectiveDeployTarget "gitmap-cli\gitmap.exe"`,
		`Move-Item -Path $resolvedActive -Destination $backupPath -Force`,
		`Stop-Process -Id $proc.ProcessId -Force -ErrorAction SilentlyContinue`,
		`Write-Host "  [OK] Synced after killing stale processes." -ForegroundColor Green`,
		`powershell.json points to an older deploy location; using PATH-derived target for verification.`,
	}

	for _, check := range checks {
		if !strings.Contains(script, check) {
			t.Fatalf("expected script to contain %q", check)
		}
	}
}
