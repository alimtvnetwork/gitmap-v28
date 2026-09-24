package cmdssh

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type sshCloneOptions struct {
	Repo     string
	DestPath string
	Target   string
	Except   string
	Exclude  string
	IsSSH    bool
	IsHTTPS  bool
	IsHelp   bool
}

func parseSSHCloneOptions(args []string) sshCloneOptions {
	if hasHelpFlag(args) || (len(args) > 0 && args[0] == "help") {
		return sshCloneOptions{IsHelp: true}
	}

	fs := flag.NewFlagSet("ssh-clone", flag.ContinueOnError)
	var opts sshCloneOptions
	fs.StringVar(&opts.Target, "target", "all", "Target machine alias or IP")
	fs.StringVar(&opts.Target, "t", "all", "Target machine alias or IP (shorthand)")
	fs.StringVar(&opts.Except, "except", "", "Exclude machines by alias, IP, or ID")
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.BoolVar(&opts.IsSSH, "ssh", false, "Force SSH protocol for clone")
	fs.BoolVar(&opts.IsHTTPS, "https", false, "Force HTTPS protocol for clone")

	filtered := extractLeadingCloneArgs(args)
	_ = fs.Parse(filtered.FlagArgs)

	assignPositionalCloneArgs(&opts, filtered.PosArgs)
	return opts
}

type cloneArgParts struct {
	PosArgs  []string
	FlagArgs []string
}

func extractLeadingCloneArgs(args []string) cloneArgParts {
	var pos []string
	var flags []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, args[i:]...)
			break
		}
		pos = append(pos, arg)
	}
	return cloneArgParts{PosArgs: pos, FlagArgs: flags}
}

func assignPositionalCloneArgs(opts *sshCloneOptions, pos []string) {
	if len(pos) > 0 {
		opts.Repo = pos[0]
	}
	if len(pos) > 1 {
		opts.DestPath = pos[1]
	}
}

func printSSHCloneUsage() {
	fmt.Printf("\n%sUsage:%s gitmap ssh clone <repo-name|url> [git|path] [flags]\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Aliases:")
	fmt.Println("  gitmap ssh-clone <repo-name|url> [git|path]")
	fmt.Println("  gitmap ssh-c <repo-name|url> [git|path]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -t, --target string   Target machine alias, IP, or 'all' (default: 'all')")
	fmt.Println("      --except string   Exclude machines by alias, IP, or ID")
	fmt.Println("      --ssh             Force clone over SSH (git@host:org/repo.git)")
	fmt.Println("      --https           Force clone over HTTPS")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ssh clone my-service")
	fmt.Println("  gitmap ssh clone alimtvnetwork/gitmap-v28 git")
	fmt.Println("  gitmap ssh clone https://github.com/org/repo.git /var/www/repo --target worker-1")
	fmt.Println()
}
