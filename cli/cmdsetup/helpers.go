// Package cmdsetup — helpers.go provides shared utilities for setup subcommands.
package cmdsetup

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

type VHostConfig = cmdvhost.VHostConfig

const (
	VHostSiteTypeLaravel   = cmdvhost.VHostSiteTypeLaravel
	VHostSiteTypeWordpress = cmdvhost.VHostSiteTypeWordpress
)

var RenderVHostConfig = cmdvhost.RenderVHostConfig

var versionRegex = lazyregex.New(`v?(\d+\.\d+(?:\.\d+)?)`)

func reorderFlagsBeforeArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}

	return append(flags, positional...)
}

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helpdoc.Print(command)
			cliexit.Exit(0)
		}
	}
}

func parseVersionFromOutput(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		m := versionRegex.FindString(line)
		if m != "" {
			return m
		}
	}

	return ""
}
