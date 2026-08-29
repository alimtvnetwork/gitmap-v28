package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runSSHCat displays the public key for a named SSH key.
func runSSHCat(args []string) error {
	fs := flag.NewFlagSet("ssh-cat", flag.ExitOnError)
	nameFlag := fs.String("name", constants.DefaultSSHKeyName, "Key name")
	fs.StringVar(nameFlag, "n", constants.DefaultSSHKeyName, "Key name (short)")
	fs.Parse(args)

	name := *nameFlag
	// Allow positional: `gitmap ssh view mykey`.
	for _, a := range fs.Args() {
		if !strings.HasPrefix(a, "-") {
			name = a

			break
		}
	}

	db, err := openDB()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.sshcat.openDB",
			"E1150",
			"failed to open database for ssh-cat",
			"cmd.sshcat",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	// Fallback: if the default name was requested and missing, and there
	// is exactly one stored key, use that one.
	hasDefaultAndErr := err != nil && name == constants.DefaultSSHKeyName
	if hasDefaultAndErr == true {
		key, err = fallbackAndAssign(db, key, err)
	}

	// If key was found
	if err == nil {
		pub := strings.TrimSpace(key.PublicKey)
		fmt.Println(pub)
		copyPubKeyAndAnnounce(pub)

		return nil
	}

	diskPath := defaultSSHKeyPath(name)
	exists := keyExistsOnDisk(diskPath)
	if exists == false {
		printSSHNotFound(db, name)
	}

	pubBytes, rerr := os.ReadFile(diskPath + ".pub")
	if rerr != nil {
		printSSHNotFound(db, name)
	}

	pub := strings.TrimSpace(string(pubBytes))
	fp := readFingerprint(diskPath)
	upsertExistingKeyToDB(db, name, diskPath, string(pubBytes), fp)
	fmt.Println(pub)
	copyPubKeyAndAnnounce(pub)
	return nil
}

func printSSHNotFound(db *store.DB, name string) {
	fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
	printAvailableKeys(db)
	appErr := apperror.NewWithDetails(
		"ssh.findKey",
		"E2022",
		fmt.Sprintf(constants.ErrSSHNotFound, name),
		"cmd.sshcat",
		apperror.ErrorTypeNotFound,
		apperror.SeverityError,
		map[string]any{"name": name},
	)
	cliexit.HandleError(appErr, 1)
}
func fallbackToSingleKey(db *store.DB, fallbackKey *model.SSHKey, fallbackErr error) (*model.SSHKey, error) {
	keys, lerr := db.ListSSHKeys()
	hasOneKey := lerr == nil && len(keys) == 1
	if hasOneKey == true {
		return &keys[0], nil
	}
	return fallbackKey, fallbackErr
}

// printAvailableKeys prints available SSH key names to stderr.
func printAvailableKeys(db *store.DB) {
	names, err := db.SSHKeyNames()
	if err != nil || len(names) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrSSHAvailable, strings.Join(names, ", "))
}

func fallbackAndAssign(db *store.DB, origKey model.SSHKey, origErr error) (model.SSHKey, error) {
	keyPtr, errFallback := fallbackToSingleKey(db, &origKey, origErr)
	if keyPtr != nil {
		return *keyPtr, errFallback
	}
	return origKey, errFallback
}
