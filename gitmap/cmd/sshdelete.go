package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runSSHDelete removes an SSH key record and optionally its files.
func runSSHDelete(args []string) {
	fs := flag.NewFlagSet("ssh-delete", flag.ExitOnError)
	nameFlag := fs.String("name", "", "Key name")
	fs.StringVar(nameFlag, "n", "", "Key name (short)")
	filesFlag := fs.Bool("files", false, "Also delete key files from disk")
	fs.Parse(args)

	name := *nameFlag
	if len(name) == 0 && fs.NArg() > 0 {
		name = fs.Arg(0)
	}
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrSSHNameEmpty)
		os.Exit(1)
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHDeleteConfirm, name)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		return
	}

	if err := db.DeleteSSHKey(name); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHDelete, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHDeleted, name)

	if *filesFlag {
		removeKeyFiles(key.PrivatePath)
		fmt.Fprint(os.Stdout, constants.MsgSSHDeletedFiles)
	}

	updateSSHConfig(db)
}
