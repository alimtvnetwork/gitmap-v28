package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/osuser"
)

// dispatchUser handles the "user" command routing.
func dispatchUser(command string) (bool, error) {
	if command != "user" {
		return false, nil
	}

	args := os.Args[2:]
	if len(args) == 0 {
		printUserUsage()

		return true, nil
	}

	return true, dispatchUserSubcommand(args[0], args[1:])
}

func dispatchUserSubcommand(sub string, args []string) error {
	switch sub {
	case "add":
		return runUserAdd(args)
	case "rm", "delete", "remove", "del":
		return runUserRm(args)
	case "create-root", "root":
		return runUserCreateRoot(args)
	case "kill-processes", "kill":
		return runUserKill(args)
	case "add-ssh-key", "key", "ssh-key":
		return runUserSSHKey(args)
	case "help", "--help", "-h":
		printUserUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown user command: %s\n", sub)
		printUserUsage()
		cliexit.HandleError(nil, 1)
		return nil
	}
}

func runUserAdd(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing username\n\nUsage: gitmap user add <username> [--password <pwd>]\n")
		cliexit.HandleError(nil, 1)
		return nil
	}
	pwd := extractUserArgFlag(args[1:], "--password")
	if err := osuser.AddUser(args[0], pwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
		cliexit.HandleError(err, 1)
		return nil
	}
	fmt.Printf("✔ Successfully created user %q\n", args[0])
	return nil
}

func runUserRm(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing username\n\nUsage: gitmap user rm <username> [--remove-home] [--kill]\n")
		cliexit.HandleError(nil, 1)
		return nil
	}
	opts := osuser.UserRemoveOptions{
		Username:       args[0],
		IsRemoveHome:   !hasUserFlag(args[1:], "--keep-home"),
		IsCleanSudoers: true,
		IsForceKill:    hasUserFlag(args[1:], "--kill") || hasUserFlag(args[1:], "-k"),
	}
	if err := osuser.Remove(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing user: %v\n", err)
		cliexit.HandleError(err, 1)
		return nil
	}
	fmt.Printf("✔ Successfully removed user %q\n", args[0])
	return nil
}

func runUserCreateRoot(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing username\n\nUsage: gitmap user create-root <username> [flags]\n")
		cliexit.HandleError(nil, 1)
		return nil
	}
	opts := buildCreateRootOptions(args[0], args[1:])
	if err := osuser.CreateRoot(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating root user: %v\n", err)
		cliexit.HandleError(err, 1)
		return nil
	}
	fmt.Printf("✔ Successfully created root user %q\n", args[0])
	return nil
}

func buildCreateRootOptions(username string, args []string) osuser.UserCreateOptions {
	return osuser.UserCreateOptions{
		Username:       username,
		Password:       extractUserArgFlag(args, "--password"),
		Theme:          extractUserArgFlagDefault(args, "--theme", "fletcherm"),
		HomeDir:        extractUserArgFlag(args, "--homedir"),
		IsSudoer:       true,
		IsConfigureZsh: !hasUserFlag(args, "--no-zsh"),
	}
}

func runUserKill(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing username\n\nUsage: gitmap user kill <username> [--force]\n")
		cliexit.HandleError(nil, 1)
		return nil
	}
	opts := osuser.UserKillOptions{
		Username: args[0],
		IsForce:  hasUserFlag(args[1:], "--force") || hasUserFlag(args[1:], "-f"),
	}
	if err := osuser.Kill(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error terminating user processes: %v\n", err)
		cliexit.HandleError(err, 1)
		return nil
	}
	fmt.Printf("✔ Successfully terminated processes for user %q\n", args[0])
	return nil
}

func runUserSSHKey(args []string) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: missing key\n\nUsage: gitmap user add-ssh-key <username> <key-or-file>\n")
		cliexit.HandleError(nil, 1)
		return nil
	}
	opts := osuser.SSHKeyOptions{
		Username:  args[0],
		PublicKey: args[1],
	}
	if err := osuser.InstallKey(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing SSH key: %v\n", err)
		cliexit.HandleError(err, 1)
		return nil
	}
	fmt.Printf("✔ Successfully installed SSH key for user %q\n", args[0])
	return nil
}

func extractUserArgFlag(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func extractUserArgFlagDefault(args []string, flag, defVal string) string {
	val := extractUserArgFlag(args, flag)
	if val != "" {
		return val
	}
	return defVal
}

func hasUserFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func printUserUsage() {
	fmt.Println("Usage: gitmap user <command> [arguments]")
	fmt.Println()
	fmt.Println("The user command manages cross-platform OS-level user accounts (Windows, Ubuntu, Debian, Fedora).")
	fmt.Println("It seamlessly runs the appropriate underlying native commands (e.g. net user, useradd).")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <username> [--password <pwd>]            Create a new OS user")
	fmt.Println("  create-root <username> [flags]               Create root/sudo user with ZSH & SSH keys")
	fmt.Println("  rm <username> [--remove-home] [--kill]       Remove an OS user and cleanup sudoers")
	fmt.Println("  kill <username> [--force]                    Terminate all processes owned by user")
	fmt.Println("  add-ssh-key <username> <key-or-file>         Install public SSH key to authorized_keys")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap user add johndoe")
	fmt.Println("  gitmap user create-root deployer --password secret --ssh-key 'ssh-ed25519 AAA...'")
	fmt.Println("  gitmap user rm johndoe")
	fmt.Println("  gitmap user kill johndoe")
}
