package completion

// generatePowerShell returns the PowerShell completion script.
func generatePowerShell() string {
	return `Register-ArgumentCompleter -CommandName gitmap -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $elems = $commandAst.CommandElements | Select-Object -Skip 1
    $cmd = if ($elems.Count -gt 0) { $elems[0].ToString() } else { "" }
    $prev = if ($elems.Count -gt 1) { $elems[$elems.Count - 1].ToString() } else { "" }
    $sub = if ($elems.Count -gt 1) { $elems[1].ToString() } else { "" }

    if (($cmd -eq "cd" -or $cmd -eq "go") -and ($prev -eq "--group" -or $prev -eq "-g")) {
        gitmap completion --list-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "cd" -or $cmd -eq "go") {
        $items = @(gitmap completion --list-repos) + @("repos", "set-default", "clear-default")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "pull") {
        gitmap completion --list-repos | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "exec" -and ($prev -eq "--group")) {
        gitmap completion --list-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "group" -or $cmd -eq "g") {
        $subs = @("create", "add", "remove", "list", "show", "delete", "pull", "status", "exec", "clear")
        $groups = @(gitmap completion --list-groups)
        $items = $subs + $groups
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "list" -or $cmd -eq "ls") {
        $items = @("go", "node", "nodejs", "react", "cpp", "csharp", "groups", "--group", "--verbose")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "multi-group" -or $cmd -eq "mg") {
        $subs = @("pull", "status", "exec", "clear")
        $groups = @(gitmap completion --list-groups)
        $items = $subs + $groups
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if (($cmd -eq "release" -or $cmd -eq "r") -and ($prev -eq "--zip-group")) {
        gitmap completion --list-zip-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "release" -or $cmd -eq "r") {
        $items = @("--assets", "--commit", "--branch", "--bump", "--draft", "--dry-run", "--compress", "--checksums", "--bin", "--targets", "--list-targets", "--verbose", "--zip-group", "-Z", "--bundle")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "release-branch" -or $cmd -eq "rb") {
        $items = @("--assets", "--draft", "--dry-run", "--compress", "--checksums", "--bin", "--targets")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "alias" -or $cmd -eq "a") {
        $subs = @("set", "remove", "list", "show", "suggest")
        $aliases = @(gitmap completion --list-aliases)
        $items = $subs + $aliases
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if (($cmd -eq "zip-group" -or $cmd -eq "z") -and ($sub -eq "add" -or $sub -eq "show" -or $sub -eq "delete" -or $sub -eq "remove" -or $sub -eq "rename")) {
        gitmap completion --list-zip-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "zip-group" -or $cmd -eq "z") {
        $subs = @("create", "add", "remove", "list", "show", "delete", "rename")
        $zgroups = @(gitmap completion --list-zip-groups)
        $items = $subs + $zgroups
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "dashboard" -or $cmd -eq "db") {
        $items = @("--limit", "--since", "--no-merges", "--out-dir", "--open")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "ssh" -and ($sub -eq "cat" -or $sub -eq "delete" -or $sub -eq "rm") -and ($prev -eq "--name" -or $prev -eq "-n")) {
        gitmap completion --list-ssh-keys | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "ssh") {
        $subs = @("cat", "list", "ls", "delete", "rm", "config", "--name", "--path", "--email", "--force")
        $subs | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "help" -and ($prev -eq "--compact")) {
        gitmap completion --list-help-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "help") {
        $items = @("--compact")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "clone-next" -or $cmd -eq "cn") {
        $items = @("v++", "--delete", "--keep", "--no-desktop", "--ssh-key", "--verbose", "--force", "-f")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "version-history" -or $cmd -eq "vh") {
        $items = @("--limit", "--json")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "mv" -or $cmd -eq "move") {
        $items = @("--prefer-newer", "--prefer-left", "--prefer-right", "--dry-run", "--verbose")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "merge-both" -or $cmd -eq "mb" -or $cmd -eq "merge-left" -or $cmd -eq "ml" -or $cmd -eq "merge-right" -or $cmd -eq "mr") {
        $items = @("--prefer-newer", "--prefer-left", "--prefer-right", "--prefer-larger", "--dry-run", "--verbose")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if (($cmd -eq "vscode-pm-sync" -or $cmd -eq "vpm") -and ($prev -eq "--mode")) {
        $items = @("union", "replace", "intersection")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if (($cmd -eq "vscode-pm-sync" -or $cmd -eq "vpm") -and ($prev -eq "--projects-json")) {
        return
    }

    if ($cmd -eq "vscode-pm-sync" -or $cmd -eq "vpm") {
        $items = @("--dry-run", "--projects-json", "--tag", "--mode")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if (($cmd -eq "llm-docs" -or $cmd -eq "ld") -and ($prev -eq "--format")) {
        $items = @("markdown", "json")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if (($cmd -eq "llm-docs" -or $cmd -eq "ld") -and ($prev -eq "--sections")) {
        $items = @("commands", "architecture", "flags", "conventions", "structure", "database", "installation", "patterns")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "llm-docs" -or $cmd -eq "ld") {
        $items = @("--stdout", "--format", "--sections")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($prev -eq "-A" -or $prev -eq "--alias") {
        gitmap completion --list-aliases | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($prev -eq "--zip-group") {
        gitmap completion --list-zip-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($prev -eq "--ssh-key" -or $prev -eq "-K") {
        gitmap completion --list-ssh-keys | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    gitmap completion --list-commands | Where-Object { $_ -like "$wordToComplete*" } |
        ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
}
`
}
