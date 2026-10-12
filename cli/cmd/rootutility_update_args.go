package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const updateUsageHelpText = `Available update commands:
  gitmap update                          - Self-update local gitmap binary
  gitmap update all                      - Update all fleet cluster nodes in parallel (alias: gitmap ua)
  gitmap update all --include-others     - Update all cluster nodes and discovered hosts
  gitmap update all zip --include-others - Distribute zip packages via SCP across all nodes
  gitmap update-all-zip --include-others - Fast zip distribution across all nodes (alias: uaz)
  gitmap update ls                       - List installed software inventory across fleet nodes
  gitmap update agm                      - Update Antigravity Manager
  gitmap update antigravity              - Update Google Antigravity Desktop IDE (alias: gitmap update ide)
  gitmap update ssh <target>             - Update package on specified SSH target
`

func printUpdateUsageHelp() {
	fmt.Print(updateUsageHelpText)
}

func handleUnknownUpdateTarget(token string) {
	header := formatUnknownUpdateTargetHeader(token)
	sugg := suggestUpdateTarget(token)
	suggList := []string{"gitmap update", "gitmap update all", "gitmap update ls"}
	if sugg != "" {
		suggList = []string{"gitmap update " + sugg}
	}
	store.LogFailedCommand(token, strings.Join(os.Args[1:], " "), "update", "E1001", header, suggList)
	fmt.Printf("\n%s\n\n", header)
	printUpdateUsageHelp()
}

func formatUnknownUpdateTargetHeader(token string) string {
	suggestion := suggestUpdateTarget(token)
	if suggestion != "" {
		return fmt.Sprintf("gitmap update: Unknown update target '%s'.\n  Did you mean: gitmap update %s?", token, suggestion)
	}
	return fmt.Sprintf("gitmap update: Unknown update target '%s'.", token)
}

func suggestUpdateTarget(token string) string {
	switch strings.ToLower(token) {
	case "al", "allnodes", "all-node", "nodes", "hosts", "cluster":
		return "all"
	case "l", "lss", "inventory", "apps", "soft":
		return "ls"
	case "ag", "ang", "manager":
		return "agm"
	case "git", "gitm", "gm":
		return "gitmap"
	case "rem", "ss":
		return "ssh"
	case "uaz", "update-all-zip", "updateallzip":
		return "update-all-zip"
	case "zp", "zi":
		return "zip"
	case "antigravity", "antigravity-ide":
		return "antigravity"
	case "ide", "agy-ide":
		return "ide"
	default:
		return ""
	}
}

func resolveUpdatePackage(args []string) string {
	if isAgmUpdateTarget(args) {
		return "agm"
	}
	if isAntigravityIdeUpdateTarget(args) {
		return "antigravity"
	}
	return "gitmap"
}

func isAntigravityIdeUpdateTarget(args []string) bool {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") && isAntigravityIdeToken(arg) {
			return true
		}
	}
	return false
}

func isAntigravityIdeToken(arg string) bool {
	switch strings.ToLower(arg) {
	case "antigravity", "antigravity-ide", "ide", "agy-ide":
		return true
	default:
		return false
	}
}

func isAgmUpdateTarget(args []string) bool {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") && isAgmToken(arg) {
			return true
		}
	}
	return false
}

func isAgmToken(arg string) bool {
	switch strings.ToLower(arg) {
	case "agm", "agy", "ag-manager", "antigravity-manager":
		return true
	default:
		return false
	}
}

func hasPinArg(args []string) bool {
	for _, a := range args {
		if strings.EqualFold(a, "--pin") || strings.EqualFold(a, "-pin") {
			return true
		}
	}
	return false
}

func isSSHFlagToken(a string) bool {
	low := strings.ToLower(a)
	return low == "--ssh" || low == "-ssh" || low == "--remote"
}

func hasSSHArg(args []string) bool {
	for _, a := range args {
		if isSSHFlagToken(a) {
			return true
		}
	}
	return false
}

func stripSSHFlag(args []string) []string {
	var cleaned []string
	for _, a := range args {
		if !isSSHFlagToken(a) {
			cleaned = append(cleaned, a)
		}
	}
	return cleaned
}

func isKnownUpdateTargetOrFlag(token string) bool {
	if strings.HasPrefix(token, "-") || isSemverLike(token) {
		return true
	}
	switch strings.ToLower(token) {
	case "all", "all-nodes", "allnodes", "ls", "list", "ssh", "remote", "gitmap", "agm", "agy", "ag-manager", "antigravity-manager":
		return true
	case "zip", "uaz", "update-all-zip", "updateallzip", "--include-others", "--include-other":
		return true
	case "antigravity", "antigravity-ide", "ide", "agy-ide":
		return true
	default:
		return false
	}
}

func isVersionListRequest(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(strings.TrimSpace(a))
		if low == "ls" || low == "list" || low == "versions" || low == "tags" || low == "ls-versions" {
			return true
		}
	}
	return false
}

func isSemverLike(s string) bool {
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 || !isNumericString(p) {
			return false
		}
	}
	return true
}

func isNumericString(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func extractVersionFromArgs(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--version" && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			return args[i+1]
		}
		if strings.HasPrefix(a, "--version=") {
			return strings.TrimPrefix(a, "--version=")
		}
		if !strings.HasPrefix(a, "-") && isSemverLike(a) {
			return a
		}
	}
	return ""
}

func hasDryRunArg(args []string) bool {
	for _, a := range args {
		if a == "--dry-run" || a == "-n" {
			return true
		}
	}
	return false
}

func hasVerboseArg(args []string) bool {
	for _, a := range args {
		if a == "--verbose" || a == "-v" {
			return true
		}
	}
	return false
}


func hasYesArg(args []string) bool {
	for _, a := range args {
		if a == "--yes" || a == "-y" {
			return true
		}
	}
	return false
}

func hasDownloadMustArg(args []string) bool {
	for _, a := range args {
		if a == "--download-must" || a == "--redownload" || a == "--force-download" || a == "-fd" {
			return true
		}
	}
	return false
}


func extractPrefixFromArgs(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--prefix" || a == "--dir") && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(a, "--prefix=") {
			return strings.TrimPrefix(a, "--prefix=")
		}
		if strings.HasPrefix(a, "--dir=") {
			return strings.TrimPrefix(a, "--dir=")
		}
	}
	return ""
}
