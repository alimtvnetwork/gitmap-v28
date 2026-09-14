package cmdos

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runOSGroup(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) {
		printOSGroupUsage()

		return nil
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	return dispatchOSGroupSubcommand(sub, rest)
}

func dispatchOSGroupSubcommand(sub string, rest []string) error {
	switch sub {
	case "ls", "list":
		return runOSGroupList()
	case "add", "create":
		return runOSGroupAdd(rest)
	case "rm", "delete", "del":
		return runOSGroupRemove(rest)
	default:
		return apperror.NewSimple("unknown os group subcommand: "+sub, "E_OS_GROUP_INVALID")
	}
}

func runOSGroupList() error {
	fmt.Printf("▶ System Groups (%s):\n\n", runtime.GOOS)
	cmd := buildGroupListCmd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func buildGroupListCmd() *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return exec.Command("net", "localgroup")
	}

	if runtime.GOOS == constants.OSDarwin {
		return exec.Command("dscl", ".", "-list", "/Groups")
	}

	return exec.Command("cut", "-d:", "-f1", "/etc/group")
}

func runOSGroupAdd(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("group name required", "E_OS_GROUP_NAME_REQUIRED")
	}

	groupName := args[0]
	cmd := buildGroupAddCmd(groupName)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "add group "+groupName)
	}

	fmt.Printf("✔ Group %q created successfully\n", groupName)

	return nil
}

func buildGroupAddCmd(group string) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return exec.Command("net", "localgroup", group, "/add")
	}

	if runtime.GOOS == constants.OSDarwin {
		return exec.Command("dscl", ".", "-create", "/Groups/"+group)
	}

	return exec.Command("groupadd", group)
}

func runOSGroupRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("group name required", "E_OS_GROUP_NAME_REQUIRED")
	}

	groupName := args[0]
	cmd := buildGroupRemoveCmd(groupName)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "remove group "+groupName)
	}

	fmt.Printf("✔ Group %q removed successfully\n", groupName)

	return nil
}

func buildGroupRemoveCmd(group string) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return exec.Command("net", "localgroup", group, "/delete")
	}

	if runtime.GOOS == constants.OSDarwin {
		return exec.Command("dscl", ".", "-delete", "/Groups/"+group)
	}

	return exec.Command("groupdel", group)
}

func printOSGroupUsage() {
	fmt.Println(`Usage: gitmap os group [subcommand] [args]

Commands:
  ls, list           List system user groups
  add <group>        Create a new system user group
  rm <group>         Delete an existing user group`)
}
