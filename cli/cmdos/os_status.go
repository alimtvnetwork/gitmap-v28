package cmdos

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvmware"
)

func runOSStatus(args []string) error {
	checkHelp("os", args)
	fmt.Println("▶ gitmap os status")
	fmt.Printf("  • Operating System: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  • User Home:        %s\n", expandHome("~"))

	desktopDir := cmdvmware.ResolveUserDesktopDir()
	fmt.Printf("  • Desktop Dir:      %s\n", desktopDir)

	return inspectStandardLinksStatus(desktopDir)
}

func runOSCleanup(subArgs []string) error {
	if len(subArgs) > 0 && isDevTarget(subArgs[0]) {
		return RunOSDevClean(subArgs[1:])
	}

	return runOSClean(subArgs)
}

func runOSDevSubcommand(subArgs []string) error {
	if len(subArgs) == 0 {
		return RunOSDevClean(nil)
	}

	sub := strings.ToLower(subArgs[0])
	if sub == "clean" || sub == "cleanup" || sub == "dev-clean" {
		return RunOSDevClean(subArgs[1:])
	}

	return RunOSDevClean(subArgs)
}

func isDevTarget(s string) bool {
	low := strings.ToLower(strings.TrimSpace(s))

	return low == "dev" || low == "clean-dev" || low == "dev-clean"
}
