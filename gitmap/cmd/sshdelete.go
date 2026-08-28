package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runSSHDelete removes an SSH key record and optionally its files.
func runSSHDelete(args []string) error {
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
		return apperror.NewSimple(constants.ErrSSHNameEmpty, "E9000")
	}

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrSSHQuery)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		return apperror.NewSimple(constants.ErrSSHNotFound, "E9000")
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHDeleteConfirm, name)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		return nil
	}

	if err := db.DeleteSSHKey(name); err != nil {
		return apperror.WrapSimple(err, constants.ErrSSHDelete)
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHDeleted, name)

	if *filesFlag {
		removeKeyFiles(key.PrivatePath)
		fmt.Fprint(os.Stdout, constants.MsgSSHDeletedFiles)
	}

	updateSSHConfig(db)
	return nil
}
