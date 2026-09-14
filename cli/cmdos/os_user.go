package cmdos

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/osuser"
)

func runOSUser(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) {
		printOSUserUsage()
		return nil
	}

	sub := args[0]
	switch sub {
	case "add":
		return handleOSUserAdd(args[1:])
	case "rm", "delete", "remove", "del":
		return handleOSUserRm(args[1:])
	case "create-root", "root":
		return handleOSUserCreateRoot(args[1:])
	case "kill-processes", "kill":
		return handleOSUserKill(args[1:])
	case "add-ssh-key", "key", "ssh-key":
		return handleOSUserSSHKey(args[1:])
	default:
		printOSUserUsage()
		return apperror.NewSimple("unknown os user subcommand: "+sub, "E_INVALID_USER_SUBCMD")
	}
}

func handleOSUserAdd(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing username for os user add", "E_MISSING_ARG")
	}
	pwd := extractArgFlag(args[1:], "--password")
	return osuser.AddUser(args[0], pwd)
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

func extractArgFlag(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func extractArgFlagWithDefault(args []string, flag, defVal string) string {
	val := extractArgFlag(args, flag)
	if val != "" {
		return val
	}
	return defVal
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func printOSUserUsage() {
	fmt.Println("Usage: gitmap os user <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <username> [--password <pwd>]            Create standard user")
	fmt.Println("  create-root <username> [flags]               Create root/sudo user with ZSH & SSH keys")
	fmt.Println("  rm <username> [--remove-home] [--kill]       Remove user account and cleanup sudoers")
	fmt.Println("  kill <username> [--force]                    Terminate all processes owned by user")
	fmt.Println("  add-ssh-key <username> <key-or-file>         Install public SSH key to authorized_keys")
}
