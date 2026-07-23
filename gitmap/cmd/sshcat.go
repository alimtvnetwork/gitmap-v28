package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runSSHCat displays the public key for a named SSH key.
func runSSHCat(args []string) {
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
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	// Fallback: if the default name was requested and missing, and there
	// is exactly one stored key, use that one.
	if err != nil && name == constants.DefaultSSHKeyName {
		keys, lerr := db.ListSSHKeys()
		if lerr == nil && len(keys) == 1 {
			key = keys[0]
			err = nil
		}
	}
	// Fallback: key not in DB — check disk at the conventional path.
	// Covers keys created outside gitmap (raw `ssh-keygen`, OS imports).
	if err != nil {
		diskPath := defaultSSHKeyPath(name)
		if keyExistsOnDisk(diskPath) {
			pubBytes, rerr := os.ReadFile(diskPath + ".pub")
			if rerr == nil {
				pub := strings.TrimSpace(string(pubBytes))
				fp := readFingerprint(diskPath)
				upsertExistingKeyToDB(db, name, diskPath, string(pubBytes), fp)
				fmt.Println(pub)
				copyPubKeyAndAnnounce(pub)

				return
			}
		}
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
		printAvailableKeys(db)
		fmt.Fprint(os.Stderr, "\n  Hint: run `gitmap ssh` to generate a new SSH key for GitHub.\n")
		os.Exit(1)
	}

	pub := strings.TrimSpace(key.PublicKey)
	fmt.Println(pub)
	copyPubKeyAndAnnounce(pub)
}

// printAvailableKeys prints available SSH key names to stderr.
func printAvailableKeys(db *store.DB) {
	names, err := db.SSHKeyNames()
	if err != nil || len(names) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrSSHAvailable, strings.Join(names, ", "))
}
