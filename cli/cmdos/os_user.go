package cmdos

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/osuser"
)

func runOSUser(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) || args[0] == "help" {
		printOSUserUsage()
		return nil
	}
	return dispatchOSUserSubcommand(strings.ToLower(args[0]), args[1:])
}

func dispatchOSUserSubcommand(sub string, args []string) error {
	switch sub {
	case "ls", "list":
		return runOSUserList()
	case "add":
		return handleOSUserAdd(args)
	case "edit":
		return handleOSUserEdit(args)
	case "export":
		return handleOSUserExport(args)
	case "export-all":
		return handleOSUserExportAll(args)
	case "import":
		return handleOSUserImport(args)
	case "import-all":
		return handleOSUserImportAll(args)
	case "rm", "delete", "remove", "del":
		return handleOSUserRm(args)
	case "create-root", "root":
		return handleOSUserCreateRoot(args)
	case "kill-processes", "kill":
		return handleOSUserKill(args)
	case "add-ssh-key", "key", "ssh-key":
		return handleOSUserSSHKey(args)
	default:
		printOSUserUsage()
		return apperror.NewSimple("unknown os user subcommand: "+sub, "E_INVALID_USER_SUBCMD")
	}
}

func runOSUserList() error {
	res := osuser.ExportAllUsers()
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("%-20s %-25s %-15s\n", "USERNAME", "HOMEDIR", "SHELL")
	fmt.Println(strings.Repeat("-", 62))
	for _, u := range res.Value {
		fmt.Printf("%-20s %-25s %-15s\n", u.Username, u.HomeDir, u.Shell)
	}
	return nil
}

func handleOSUserAdd(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing username for os user add", "E_MISSING_ARG")
	}
	pwd := extractArgFlag(args[1:], "--password")
	return osuser.AddUser(args[0], pwd)
}

func handleOSUserEdit(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os user edit <username> [flags]", "E_MISSING_ARG")
	}
	sh := extractArgFlag(args[1:], "--shell")
	home := extractArgFlag(args[1:], "--homedir")
	res := osuser.EditUser(args[0], sh, home)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ User %s updated successfully\n", args[0])
	return nil
}

func handleOSUserExport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os user export <username> [file.json]", "E_MISSING_ARG")
	}
	res := osuser.ExportUser(args[0])
	if res.IsFailure() {
		return res.AppError()
	}
	data, _ := json.MarshalIndent(res.Value, "", "  ")
	if len(args) >= 2 {
		return writeUserExportFile(args[1], data, args[0])
	}
	fmt.Println(string(data))
	return nil
}

func writeUserExportFile(destPath string, data []byte, username string) error {
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return apperror.WrapSimple(err, "write exported user")
	}
	fmt.Printf("✔ Exported user %s to %s\n", username, destPath)
	return nil
}

func handleOSUserExportAll(args []string) error {
	res := osuser.ExportAllUsers()
	if res.IsFailure() {
		return res.AppError()
	}
	data, _ := json.MarshalIndent(res.Value, "", "  ")
	if len(args) >= 1 {
		return writeAllUsersExportFile(args[0], data, len(res.Value))
	}
	fmt.Println(string(data))
	return nil
}

func writeAllUsersExportFile(destPath string, data []byte, count int) error {
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return apperror.WrapSimple(err, "write exported users")
	}
	fmt.Printf("✔ Exported %d users to %s\n", count, destPath)
	return nil
}

func handleOSUserImport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os user import <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import user")
	}
	var pu osuser.PortableUser
	if unmarshalErr := json.Unmarshal(content, &pu); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "parse user json")
	}
	res := osuser.ImportUser(pu)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Successfully imported user: %s\n", pu.Username)
	return nil
}

func handleOSUserImportAll(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os user import-all <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import-all users")
	}
	var users []osuser.PortableUser
	if unmarshalErr := json.Unmarshal(content, &users); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "parse users json")
	}
	res := osuser.ImportAllUsers(users)
	if res.IsFailure() {
		return res.AppError()
	}
	fmt.Printf("✔ Successfully imported %d users from %s\n", res.Value, args[0])
	return nil
}

func handleOSUserRm(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing username for os user rm", "E_MISSING_ARG")
	}
	opts := osuser.UserRemoveOptions{
		Username:       args[0],
		IsRemoveHome:   hasFlag(args[1:], "--remove-home") || !hasFlag(args[1:], "--keep-home"),
		IsCleanSudoers: true,
		IsForceKill:    hasFlag(args[1:], "--kill") || hasFlag(args[1:], "-k"),
	}
	return osuser.Remove(opts)
}

func handleOSUserCreateRoot(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing username for os user create-root", "E_MISSING_ARG")
	}
	opts := osuser.UserCreateOptions{
		Username:       args[0],
		Password:       extractArgFlag(args[1:], "--password"),
		Theme:          extractArgFlagWithDefault(args[1:], "--theme", "fletcherm"),
		HomeDir:        extractArgFlag(args[1:], "--homedir"),
		IsSudoer:       true,
		IsConfigureZsh: !hasFlag(args[1:], "--no-zsh"),
	}
	return osuser.CreateRoot(opts)
}

func handleOSUserKill(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing username for os user kill", "E_MISSING_ARG")
	}
	opts := osuser.UserKillOptions{
		Username: args[0],
		IsForce:  hasFlag(args[1:], "--force") || hasFlag(args[1:], "-f"),
	}
	return osuser.Kill(opts)
}

func handleOSUserSSHKey(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("usage: gitmap os user add-ssh-key <user> <key-or-file>", "E_MISSING_ARG")
	}
	opts := osuser.SSHKeyOptions{
		Username:  args[0],
		PublicKey: args[1],
	}
	return osuser.InstallKey(opts)
}

func printOSUserUsage() {
	fmt.Println("Usage: gitmap os user <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  ls, list                                     List all system users")
	fmt.Println("  add <username> [--password <pwd>]            Create standard user")
	fmt.Println("  edit <username> [--shell <sh>] [--homedir]   Update existing user")
	fmt.Println("  export <username> [file.json]                Export user configuration")
	fmt.Println("  export-all [file.json]                       Export all users to JSON")
	fmt.Println("  import <file.json>                           Import user configuration")
	fmt.Println("  import-all <file.json>                       Import multiple users")
	fmt.Println("  create-root <username> [flags]               Create root/sudo user")
	fmt.Println("  rm <username> [--remove-home] [--kill]       Remove user account")
	fmt.Println("  kill <username> [--force]                    Terminate user processes")
	fmt.Println("  add-ssh-key <username> <key-or-file>         Install public SSH key")
}
