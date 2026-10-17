package cmdrepo

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runRepoCommand is the entry point for the gitmap repo command suite.
func RunRepoCommand(args []string) error {
	if len(args) == 0 {
		return runRepoDB(nil)
	}

	subcmd := strings.ToLower(strings.TrimSpace(args[0]))
	switch subcmd {
	case "db":
		return runRepoDB(args[1:])
	case "create", "new", "c", "repo-create", "create-repo", "repoc", "crepo", "cr":
		return runCreate(args[1:])
	case "recreate", "recreate-repo", "repo-recreate":
		return runRecreate(args[1:])
	case "create-local", "local", "clr", "create-local-repo", "create-repo-local", "repo-create-local":
		return runCreateLocal(args[1:])
	case "export", "repox":
		return RunRepoExport(args[1:])
	case "import", "repoi":
		return RunRepoImport(args[1:])
	case "help", "-h", "--help":
		printRepoHelp()

		return nil
	default:
		printRepoHelp()

		return fmt.Errorf("unknown repo subcommand: %s", subcmd)
	}
}

func printRepoHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap repo [command] [flags]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Commands:" + constants.ColorReset)
	fmt.Printf("  %-22s %s\n", "create (c)", "Create a new repository locally and on GitHub (auto-slugifies spaces)")
	fmt.Printf("  %-22s %s\n", "create-local (clr)", "Create a local-only git repository")
	fmt.Printf("  %-22s %s\n", "recreate-repo", "Duplicate an existing local repo to a new remote")
	fmt.Printf("  %-22s %s\n", "db", "Manage repository-specific split database (status, log, clear, reset, optimize)")
	fmt.Printf("  %-22s %s\n", "export", "Export a repo manifest JSON (clone URL, remotes, branches) to a portable file")
	fmt.Printf("  %-22s %s\n", "import", "Import a repo manifest: clone, restore remotes, register in gitmap.db")
	fmt.Printf("  %-22s %s\n", "help", "Show repository command help")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap repo create \"My New Repo\"")
	fmt.Println("  gitmap repo create \"My New Repo\" ./folder custom-slug")
	fmt.Println("  gitmap repo create-local \"My Local Tool\"")
	fmt.Println("  gitmap repo export story-telling-v1 --out ./story-telling-v1.gitmap.json")
	fmt.Println("  gitmap repo import ./story-telling-v1.gitmap.json --into ~/workspace/repos")
	fmt.Println("  gitmap repo db")
	fmt.Println("  gitmap repo db status")
	fmt.Println("  gitmap repo db log")
	fmt.Println("  gitmap repo db error-logs")
	fmt.Println("  gitmap repo db optimize")
	fmt.Println("  gitmap repo db clear -y")
}
