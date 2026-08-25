package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/osuser"
)

// dispatchUser handles the "user" command routing.
func dispatchUser(command string) bool {
	if command != "user" {
		return false
	}
	
	args := os.Args[2:]
	if len(args) == 0 {
		printUserUsage()
		return true
	}

	sub := args[0]
	switch sub {
	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: missing username\n\nUsage: gitmap user add <username> [--password <pwd>]\n")
			os.Exit(1)
		}
		
		username := args[1]
		password := ""
		for i, a := range args[2:] {
			if a == "--password" && i+1 < len(args[2:]) {
				password = args[2:][i+1]
			}
		}

		if err := osuser.AddUser(username, password); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✔ Successfully created user %q\n", username)

	case "rm", "delete", "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: missing username\n\nUsage: gitmap user rm <username>\n")
			os.Exit(1)
		}
		
		username := args[1]
		if err := osuser.RemoveUser(username); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✔ Successfully removed user %q\n", username)
		
	case "help", "--help", "-h":
		printUserUsage()
		
	default:
		fmt.Fprintf(os.Stderr, "Unknown user command: %s\n", sub)
		printUserUsage()
		os.Exit(1)
	}

	return true
}

func printUserUsage() {
	fmt.Println("Usage: gitmap user <command> [arguments]")
	fmt.Println()
	fmt.Println("The user command manages cross-platform OS-level user accounts (Windows, Ubuntu, Debian, Fedora).")
	fmt.Println("It seamlessly runs the appropriate underlying native commands (e.g. net user, useradd).")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <username> [--password <pwd>]    Create a new OS user")
	fmt.Println("  rm <username>                        Remove an OS user and their profile data")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap user add johndoe")
	fmt.Println("  gitmap user add janedoe --password secret123")
	fmt.Println("  gitmap user rm johndoe")
}
