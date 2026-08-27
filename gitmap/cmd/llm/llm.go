package llm

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
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
`

// Run executes the llm command
func Run(args []string) {
	fs := flag.NewFlagSet("llm", flag.ExitOnError)
	isUrl := fs.Bool("url", false, "Output the URL to the LLM spec")
	
	if err := fs.Parse(args); err != nil {
		cliexit.Fail("llm", "parse_flags", "flags", apperror.Wrap(err, "failed to parse llm flags", nil), 1)
	}

	if *isUrl {
		fmt.Println("https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/LLM.md")
		return
	}

	fmt.Print(llmMarkdownSpec)
}
