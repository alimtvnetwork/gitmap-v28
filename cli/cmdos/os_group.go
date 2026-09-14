package cmdos

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/osuser"
)

func runOSGroup(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) || args[0] == "help" {
		printOSGroupUsage()
		return nil
	}

	return dispatchOSGroupSubcommand(strings.ToLower(args[0]), args[1:])
}

func dispatchOSGroupSubcommand(sub string, rest []string) error {
	switch sub {
	case "ls", "list":
		return runOSGroupList()
	case "add", "create":
		return runOSGroupAdd(rest)
	case "edit":
		return runOSGroupEdit(rest)
	case "export":
		return runOSGroupExport(rest)
	case "export-all":
		return runOSGroupExportAll(rest)
	case "import":
		return runOSGroupImport(rest)
	case "import-all":
		return runOSGroupImportAll(rest)
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

func runOSGroupEdit(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("usage: gitmap os group edit <old-name> <new-name>", "E_MISSING_ARG")
	}
	res := osuser.EditGroup(args[0], args[1])
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Group %s renamed to %s\n", args[0], args[1])
	return nil
}

func runOSGroupExport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os group export <name> [file.json]", "E_MISSING_ARG")
	}
	res := osuser.ExportGroup(args[0])
	if res.IsFailure() {
		return res.AppError()
	}
	data, _ := json.MarshalIndent(res.Value, "", "  ")
	if len(args) >= 2 {
		return writeGroupExportFile(args[1], data, args[0])
	}
	fmt.Println(string(data))
	return nil
}

func writeGroupExportFile(destPath string, data []byte, name string) error {
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return apperror.WrapSimple(err, "write group export")
	}
	fmt.Printf("✔ Exported group %s to %s\n", name, destPath)
	return nil
}

func runOSGroupExportAll(args []string) error {
	res := osuser.ExportAllGroups()
	if res.IsFailure() {
		return res.AppError()
	}
	data, _ := json.MarshalIndent(res.Value, "", "  ")
	if len(args) >= 1 {
		return writeAllGroupsExportFile(args[0], data, len(res.Value))
	}
	fmt.Println(string(data))
	return nil
}

func writeAllGroupsExportFile(destPath string, data []byte, count int) error {
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return apperror.WrapSimple(err, "write all groups export")
	}
	fmt.Printf("✔ Exported %d groups to %s\n", count, destPath)
	return nil
}

func runOSGroupImport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os group import <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read group import")
	}
	var pg osuser.PortableGroup
	if unmarshalErr := json.Unmarshal(content, &pg); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "parse group json")
	}
	res := osuser.ImportGroup(pg)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Group %s imported successfully\n", pg.Name)
	return nil
}

func runOSGroupImportAll(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os group import-all <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read all groups import")
	}
	var groups []osuser.PortableGroup
	if unmarshalErr := json.Unmarshal(content, &groups); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "parse groups json")
	}
	res := osuser.ImportAllGroups(groups)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Imported %d groups successfully\n", res.Value)
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
  ls, list                      List system user groups
  add <group>                   Create a new system user group
  edit <old-name> <new-name>    Rename an existing group
  export <name> [file.json]     Export a group definition to JSON
  export-all [file.json]        Export all system groups to JSON
  import <file.json>            Import a group definition from JSON
  import-all <file.json>        Import multiple groups from JSON
  rm <group>                    Delete an existing user group`)
}
