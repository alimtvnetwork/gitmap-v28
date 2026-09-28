package completion

// generatePowerShell returns the modern dynamic PowerShell completion script with PSReadLine prediction.
func generatePowerShell() string {
	return `Register-ArgumentCompleter -Native -CommandName gitmap -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commandElements = $commandAst.CommandElements
    $command = @(
        for ($i = 1; $i -lt $commandElements.Count; $i++) {
            $element = $commandElements[$i]
            if ($element -is [System.Management.Automation.Language.StringConstantExpressionAst] -and
                $element.StringConstantType -ne [System.Management.Automation.Language.StringConstantType]::BareWord) {
                $element.Value
            } else {
                $element.Extent.Text
            }
        }
    )

    if ([string]::IsNullOrEmpty($wordToComplete)) {
        $command += ""
    }

    $completions = @(gitmap __complete @command 2>$null)
    if ($completions.Count -eq 0) {
        return
    }

    $directive = $completions[-1]
    if ($directive -match '^:[0-9]+$') {
        $completions = $completions[0..($completions.Count - 2)]
    }

    foreach ($line in $completions) {
        if ([string]::IsNullOrWhiteSpace($line)) { continue }
        $parts = $line -split [char]9, 2
        $val = $parts[0]
        $desc = if ($parts.Count -gt 1) { $parts[1] } else { $val }
        [System.Management.Automation.CompletionResult]::new($val, $val, 'ParameterValue', $desc)
    }
}

if ((Get-Module -ListAvailable -Name PSReadLine) -and -not [Console]::IsOutputRedirected) {
    try {
        Set-PSReadLineOption -PredictionSource History -ErrorAction SilentlyContinue
        Set-PSReadLineOption -PredictionViewStyle ListView -ErrorAction SilentlyContinue
    } catch {}
}
`
}
