package cmdssh

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runSSHCat displays the public key for a named SSH key.
func runSSHCat(args []string) error {
	fs := flag.NewFlagSet("ssh-view", flag.ContinueOnError)
	nameFlag := fs.String("name", constants.DefaultSSHKeyName, "Key name")
	fs.StringVar(nameFlag, "n", constants.DefaultSSHKeyName, "Key name (short)")
	rawFlag := fs.Bool("raw", false, "Output only raw public key")
	fs.BoolVar(rawFlag, "r", false, "Output only raw public key (short)")
	_ = fs.Parse(args)

	name := *nameFlag
	isRaw := *rawFlag
	for _, a := range args {
		if a == "-r" || a == "--raw" {
			isRaw = true
		}
	}
	// Allow positional: `gitmap ssh view mykey`.
	for _, a := range fs.Args() {
		if !strings.HasPrefix(a, "-") {
			name = a

			break
		}
	}

	db, err := openDB()
	if err != nil {
		return apperror.WrapWithDetails(
			err,
			"cmd.sshcat.openDB",
			"E1150",
			"failed to open database for ssh-cat",
			"cmd.sshcat",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
	}

	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	// Fallback: if the default name was requested and missing, and there
	// is exactly one stored key, use that one.
	hasDefaultAndErr := err != nil && name == constants.DefaultSSHKeyName
	if hasDefaultAndErr {
		key, err = fallbackAndAssign(db, key, err)
	}

	// If key was found
	if err == nil {
		pub := strings.TrimSpace(key.PublicKey)
		printSSHKeyViewCard(key.Name, key.PrivatePath, pub, key.Fingerprint, isRaw)

		return nil
	}

	diskPath := defaultSSHKeyPath(name)
	exists := keyExistsOnDisk(diskPath)
	if !exists {
		printSSHNotFound(db, name)
	}

	pubBytes, rerr := os.ReadFile(diskPath + ".pub")
	if rerr != nil {
		printSSHNotFound(db, name)
	}

	pub := strings.TrimSpace(string(pubBytes))
	fp := readFingerprint(diskPath)
	upsertExistingKeyToDB(db, name, diskPath, string(pubBytes), fp)
	printSSHKeyViewCard(name, diskPath, pub, fp, isRaw)

	return nil
}

func printSSHKeyViewCard(name, keyPath, pub, fp string, isRaw bool) {
	if isRaw {
		fmt.Println(pub)

		return
	}

	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s║                     SSH PUBLIC KEY & IDENTITY                    ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Key Label:    %s%s%s\n", constants.ColorWhite, name, constants.ColorReset)
	if keyPath != "" {
		fmt.Printf("  Key Path:     %s%s%s\n", constants.ColorDim, keyPath, constants.ColorReset)
	}
	if fp != "" {
		fmt.Printf("  Fingerprint:  %s%s%s\n", constants.ColorCyan, fp, constants.ColorReset)
	}
	fmt.Println()
	displayPub := formatDisplayPublicKey(pub, false)
	fmt.Printf("  Public Key:\n    %s\n\n", displayPub)
	copyPubKeyAndAnnounce(pub)
	fmt.Printf("\n  %s[tip]%s To regenerate this key: %sgitmap ssh create -y%s (or pass --force)\n\n",
		constants.ColorYellow, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
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

func fallbackToSingleKey(
	db *store.DB,
	fallbackKey *model.SSHKey,
	fallbackErr error,
) (*model.SSHKey, error) {
	keys, lerr := db.ListSSHKeys()
	hasOneKey := lerr == nil && len(keys) == 1
	if hasOneKey {
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
