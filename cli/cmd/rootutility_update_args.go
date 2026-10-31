package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

const updateUsageHelpText = `Available update commands:
  gitmap update                          - Self-update local gitmap binary
  gitmap update all                      - Update all fleet cluster nodes in parallel (alias: gitmap ua)
  gitmap update all --include-others     - Update all cluster nodes and discovered hosts
  gitmap update all zip --include-others - Distribute zip packages via SCP across all nodes
  gitmap update-all-zip --include-others - Fast zip distribution across all nodes (alias: uaz)
  gitmap update ls                       - List installed software inventory across fleet nodes
  gitmap update cursor [node|all]        - Update Cursor IDE across fleet nodes
  gitmap update vscode [node|all]        - Update Visual Studio Code across fleet nodes
  gitmap update agm [node|all]           - Update Antigravity Manager across fleet nodes
  gitmap update agy [node|all]           - Update Antigravity Tools / IDE across fleet nodes
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

var updateTargetSuggestions = map[string]string{
	"al": "all", "allnodes": "all", "all-node": "all", "nodes": "all", "hosts": "all", "cluster": "all",
	"l": "ls", "lss": "ls", "inventory": "ls", "apps": "ls", "soft": "ls",
	"ag": "agm", "ang": "agm", "manager": "agm",
	"git": "gitmap", "gitm": "gitmap", "gm": "gitmap",
	"rem": "ssh", "ss": "ssh",
	"uaz": "update-all-zip", "update-all-zip": "update-all-zip", "updateallzip": "update-all-zip",
	"zp": "zip", "zi": "zip",
	"antigravity": "antigravity", "antigravity-ide": "antigravity",
	"ide": "ide", "agy-ide": "ide",
	"cur": "cursor", "vsc": "vscode", "vs": "vscode", "cod": "code",
}

func suggestUpdateTarget(token string) string {
	return updateTargetSuggestions[strings.ToLower(token)]
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
	return strutil.EqualFoldAny(arg, "antigravity", "antigravity-ide", "ide", "agy-ide")
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
	return strutil.EqualFoldAny(arg, "agm", "agy", "ag-manager", "antigravity-manager")
}

func hasPinArg(args []string) bool {
	for _, a := range args {
		if strutil.EqualFoldAny(a, "--pin", "-pin") {
			return true
		}
	}
	return false
}

func isSSHFlagToken(a string) bool {
	return strutil.EqualFoldAny(a, "--ssh", "-ssh", "--remote")
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
	return strutil.EqualFoldAny(token,
		"all", "all-nodes", "allnodes", "ls", "list", "ssh", "remote", "gitmap", "agm", "agy", "ag-manager", "antigravity-manager",
		"zip", "uaz", "update-all-zip", "updateallzip", "--include-others", "--include-other",
		"antigravity", "antigravity-ide", "ide", "agy-ide",
		"cursor", "cursor-ide", "vscode", "code", "vs-code", "agy-tools")
}

func isFleetSoftwareUpdateArg(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return isSoftwareUpdateToken(args[0])
}

func isSoftwareUpdateToken(token string) bool {
	return strutil.EqualFoldAny(token,
		"cursor", "cursor-ide",
		"vscode", "code", "vs-code",
		"agm", "ag-manager", "antigravity-manager",
		"agy", "agy-tools", "antigravity", "antigravity-ide", "ide",
	)
}

func isVersionListRequest(args []string) bool {
	for _, a := range args {
		if strutil.EqualFoldAnyTrim(a, "ls", "list", "versions", "tags", "ls-versions") {
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
