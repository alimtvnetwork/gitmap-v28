package cmdpipeline

import (
	"fmt"
)

func printPipelineErrorLogsHelp() {
	printPipelineErrorLogsUsage()
	printPipelineErrorLogsFlags()
}

func printPipelineErrorLogsUsage() {
	fmt.Println("Usage: gitmap pipeline error-logs [repo] [commit|-N|-Nn|HEAD~N] [flags]")
	fmt.Println("       gitmap pipeline errors [repo] [commit|-N] [clear [-y]] [flags]")
	fmt.Println("       gitmap pe [repo] [commit|-N|-Nn|HEAD~N] [clear [-y]] [flags]")
	fmt.Println("       gitmap pe -f <format.json|alias> [flags]")
	fmt.Println("       gitmap pe -f <alias> -test <filepath>")
	fmt.Println("       gitmap pe -f <alias> -test-commit <commit-sha> [-repo <path>]")
	fmt.Println("       gitmap pe add-format <file.json> [alias]")
	fmt.Println("       gitmap pe remove-format (rm-format) <name|alias>")
	fmt.Println("       gitmap pe add-all <folder-path>")
	fmt.Println("       gitmap pe list-formats")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  clear [-y]                     Purge error logs, reports, and reset pipeline DB for current repo")
	fmt.Println("  add-format <file.json> [alias] Register custom JSON format profile for error/warning filtering")
	fmt.Println("  remove-format <name|alias>     Remove registered format profile (alias: rm-format)")
	fmt.Println("  add-all <folder-path>          Batch register all format profiles (*.json) from folder")
	fmt.Println("  list-formats                   List all registered error log format profiles")
	fmt.Println("  preview-format <name|file>     Preview format configuration rules in JSON")
	fmt.Println()
	fmt.Println("Targeting:")
	fmt.Println("  [repo]                         Target repository by name, prefix, alias, or path")
	fmt.Println("  <commit-sha>                   Filter errors for specific commit (e.g. gitmap pe ee4a694)")
	fmt.Println("  -1, -2, -3, -1n, HEAD~1        Inspect errors for previous commits by relative offset")
	fmt.Println()
}

func printPipelineErrorLogsFlags() {
	fmt.Println("Flags:")
	printPEExecutionFlags()
	printPEOutputFlags()
}

func printPEExecutionFlags() {
	fmt.Println("  -f <format|file.json>          Apply custom format profile to filter and format errors (e.g. -f tauri)")
	fmt.Println("  -f <format> -test <file>       Test format profile against a local log file")
	fmt.Println("  -f <format> -test-commit <sha> Test format profile against remote GitHub Actions run for commit")
	fmt.Println("  -1, -2, -3, HEAD~N             Target past workflow run by relative commit offset or SHA")
	fmt.Println("  -t, --timeline                 Watch pipeline dynamic timeline until completion")
	fmt.Println("  -f, --fix                      Execute internal CI/CD diagnostic & auto-repair suite (when standalone)")
	fmt.Println("  -c, --check                    Run internal CI/CD checks without modifying files")
	fmt.Println("  -v, --detailed, --verbose      Show full raw error logs including passing ok lines")
	fmt.Println("  -y, --yes                      Auto-confirm prompts non-interactively")
	fmt.Println("  --force, --no-cache            Bypass local SQLite DB cache and pull fresh from GitHub")
}

func printPEOutputFlags() {
	fmt.Println("  --json                         Output data in structured JSON format")
	fmt.Println("  --file <path>                  Write error logs to specified file path")
	fmt.Println("  --tempfile <filename>          Write error logs to .ai-memory/temp/<filename>")
	fmt.Println("  -l, --limit, -n, --lines <N>   Limit displayed error log lines and failure run records (e.g. -l 10)")
	fmt.Println("  -n, --no-output-log            Stage error logs to disk without displaying in terminal")
}

func printPipelineLogsHelp() {
	fmt.Println("Usage: gitmap pipeline logs [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --json                  Output workflow status and URL in JSON format")
}
