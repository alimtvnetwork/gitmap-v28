package llm

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const llmMarkdownSpec = `# Gitmap LLM Specification

Gitmap is a powerful CLI designed for autonomous agents, LLMs, and developers to navigate, analyze, and modify codebases efficiently. 

## Capabilities

1. **Repository Discovery & Scanning**:
   - gitmap scan: Automatically locate and index Git repositories.
   - gitmap rescan: Refresh indexed local repositories.

2. **Repository Cloning & Synchronization**:
   - gitmap clone, gitmap clone-next: Manage bulk cloning.

3. **Intelligent File Search & Indexing (Split DB)**:
   - Gitmap utilizes a Split SQLite Database architecture for rapid indexing.
   - gitmap search <query>: Execute exact searches quickly on-the-fly.
   - gitmap repo-search <query>: Perform analytical, cache-backed searches across the repo.
   - gitmap repo-search-json <query>: Retrieve search results as structured JSON (ideal for LLM tool consumption).
   - *Note*: Gitmap natively skips .git, node_modules, and handles files >300KB using lazy regex to maintain high performance.

4. **Regex & Replacement**:
   - gitmap replace <query> <replacement>: Find and replace exact phrases.
   - gitmap replace-regex <regex> <replacement>: Regex replacement with automatic history tracking.
   - gitmap replace history: View replacement operations.

5. **Releases & Commits**:
   - gitmap release: Automate semantic version bumping and changelog generation.
   - gitmap commit-in: Perform robust, orchestrator-driven batch commits.

## Instructions for LLMs
- Always prefer JSON output commands (rsj, repo-search-json) when parsing data programmatically.
- Avoid modifying the root SQLite DB manually; rely on the CLI commands.
- Before running heavy regex operations across a large repo, use gitmap search to verify matches.

## Alternative Commands for AI (Instead of Raw Git)
When checking repository status or logs, DO NOT use raw 'git status && git log -1'. Instead, use:
- 'gitmap status': Checks dirty/clean status, ahead/behind, and stash for all tracked repos.
- 'gitmap history': View a rich audit log of command executions.
- 'gitmap changelog': See concise release notes and recent commits.
- Or use the standard aliases configured by gitmap setup: 'git st && git last'.

For committing and pushing, NEVER use raw 'git commit' and 'git push'. Use the semantic commit-push suite:
- 'gitmap commit-push-feature "<message>"' (alias 'gitmap cpf'): Commit and push a feature.
- 'gitmap commit-push-bug "<message>"' (alias 'gitmap cpb'): Commit and push a bug fix.
- 'gitmap commit-push-release "<message>"' (alias 'gitmap cpr'): Commit and push a release chore.
- 'gitmap pull-commit-push "<message>"' (alias 'gitmap pcp'): Pull latest, then commit and push.
- 'gitmap rm-git <last-4-digits>': Drop a recent commit safely using git reset --hard <sha>^.

## AI File Search Patterns
When searching codebases, LLMs can use native gitmap commands OR standard terminal tools.
Here are equivalent alternative command samples for LLM search operations:

- **Find a specific struct definition**:
  - 'gitmap file-search . "type SearchResult struct"'
  - 'Get-ChildItem -Path gitmap -Recurse -File | Select-String "type SearchResult struct"'

- **Search functions with Regex context**:
  - 'gitmap file-search cmd/ "func dispatch[A-Z]" 0 10'
  - 'Get-ChildItem -Path gitmap/cmd -Filter *.go | Select-String "func dispatch[A-Z]"'

- **Find specific function contexts**:
  - 'gitmap file-search cmd/root.go "func finishCommandAudit" 0 10'
  - 'cat gitmap/cmd/root.go | Select-String "func finishCommandAudit" -Context 0,10'

- **Global Text/Keyword Search (Instead of ripgrep/rg)**:
  - 'gitmap search "CHANGELOG"'
  - 'gitmap search "CHANGELOG\.md" --case-sensitive'
  *(Never use raw rg or grep when gitmap search utilizes the built-in Split DB index for much faster cross-repo lookups).*
`

// Run executes the llm command
func Run(args []string) *apperror.AppError {
	fs := flag.NewFlagSet("llm", flag.ExitOnError)
	isUrl := fs.Bool("url", false, "Output the URL to the LLM spec")

	if err := fs.Parse(args); err != nil {
		return apperror.WrapSimple(err, "parse flags")
	}

	if *isUrl {
		fmt.Println("https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md")
		return nil
	}

	fmt.Print(llmMarkdownSpec)
	return nil
}
