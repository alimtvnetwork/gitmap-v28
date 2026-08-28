package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runSSHGenerate generates a new SSH key pair.
func runSSHGenerate(args []string) error {
	name, keyPath, email, force, host, confirm := parseSSHGenFlags(args)

	if err := validateSSHKeygen(); err != nil {
		fmt.Fprint(os.Stderr, constants.ErrSSHKeygenMissing)
		panic("error")
	}

	if len(email) == 0 {
		email = resolveGitEmail()
	}
	if len(email) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrSSHEmailResolve)
		panic("error")
	}

	keyPath = expandHome(keyPath)

	if confirm && askConfirm(name, keyPath) == false {
		fmt.Fprint(os.Stdout, constants.MsgSSHCanceled)

		return nil
	}

	db, err := openDB()
	if err != nil {
		panic("error")
	}
	defer db.Close()

	// Disk check FIRST — covers keys created outside gitmap (e.g. via raw
	// `ssh-keygen` or another tool). Without this, we'd fall through to
	// `ssh-keygen -f <existing>` which prompts "Overwrite (y/n)?" on stdin
	// and exits non-zero when stdin doesn't supply an answer — exactly the
	// bug reported in v3.30.x: "Overwrite (y/n)? Error: SSH key generation
	// failed at C:\Users\...\.ssh\id_rsa: exit status 1".
	if keyExistsOnDisk(keyPath) && !force {
		printExistingKeyOnDisk(db, name, keyPath, host)

		return nil
	}
	if keyExistsOnDisk(keyPath) && force {
		err2 := backupKeyForRegenerate(keyPath)
		exitOnBackupError(err2)
		fmt.Fprintf(os.Stdout, constants.MsgSSHBackedUp, keyPath)
	}

	shouldHandle := db.SSHKeyExists(name) && !force
	if shouldHandle && !handleExistingKey(db, name, &keyPath) {
		return nil
	}

	generateAndStore(db, name, keyPath, email, host)
	return nil
}

// parseSSHGenFlags parses flags for SSH key generation.
func parseSSHGenFlags(args []string) (name, keyPath, email string, force bool, host string, confirm bool) {
	fs := flag.NewFlagSet(constants.CmdSSH, flag.ExitOnError)
	nameFlag := fs.String("name", constants.DefaultSSHKeyName, "Key label")
	fs.StringVar(nameFlag, "n", constants.DefaultSSHKeyName, "Key label (short)")
	pathFlag := fs.String("path", "", "Key file path")
	fs.StringVar(pathFlag, "p", "", "Key file path (short)")
	emailFlag := fs.String("email", "", "Email comment")
	fs.StringVar(emailFlag, "e", "", "Email comment (short)")
	forceFlag := fs.Bool("force", false, "Skip prompt if key exists")
	fs.BoolVar(forceFlag, "f", false, "Skip prompt (short)")
	hostFlag := fs.String("host", constants.DefaultSSHHost, "Git provider hostname")
	fs.StringVar(hostFlag, "H", constants.DefaultSSHHost, "Git provider hostname (short)")
	confirmFlag := fs.Bool("confirm", false, "Require explicit confirmation")
	fs.Parse(args)

	name = *nameFlag
	email = *emailFlag
	// Accept positional args: an "@"-bearing token is treated as the
	// email comment, anything else as the key name. Lets users write
	// `gitmap ssh create me@x.com` or `gitmap ssh create mykey me@x.com`.
	for _, a := range fs.Args() {
		if strings.Contains(a, "@") && len(email) == 0 {
			email = a

			continue
		}
		if name == constants.DefaultSSHKeyName {
			name = a
		}
	}
	// `--name me@x.com` (or `-n me@x.com`) — treat as email too.
	if strings.Contains(name, "@") && len(email) == 0 {
		email = name
		name = constants.DefaultSSHKeyName
	}

	path := *pathFlag
	if len(path) == 0 {
		path = defaultSSHKeyPath(name)
	}

	return name, path, email, *forceFlag, *hostFlag, *confirmFlag
}

// handleExistingKey prompts the user when a key already exists.
// Returns true if generation should proceed, false to cancel.
func handleExistingKey(db *store.DB, name string, keyPath *string) bool {
	existing, _ := db.FindSSHKeyByName(name)
	fmt.Fprintf(os.Stdout, constants.MsgSSHExists, name, existing.PrivatePath)
	fmt.Fprintf(os.Stdout, constants.MsgSSHExistsFP, existing.Fingerprint)
	fmt.Fprint(os.Stdout, constants.MsgSSHPromptAction)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToUpper(input))

	if input == "R" {
		removeKeyFiles(existing.PrivatePath)
		*keyPath = existing.PrivatePath

		return true
	}
	if input == "N" {
		fmt.Fprint(os.Stdout, constants.MsgSSHNewPathPrompt)
		newPath, _ := reader.ReadString('\n')
		*keyPath = expandHome(strings.TrimSpace(newPath))

		return true
	}

	return false
}

// generateAndStore runs ssh-keygen and stores the result in the database.
func generateAndStore(db *store.DB, name, keyPath, email, host string) {
	if err := ensureSSHDir(filepath.Dir(keyPath)); err != nil {
		panic("error")
	}

	cmd := exec.Command(constants.SSHKeygenBin,
		"-t", constants.SSHKeyType,
		"-b", constants.SSHKeyBits,
		"-C", email,
		"-f", keyPath,
		"-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic("error")
	}

	pubKey, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		panic("error")
	}

	fingerprint := readFingerprint(keyPath)

	exists := db.SSHKeyExists(name)
	if exists == true {
		errUpdate := db.UpdateSSHKey(name, keyPath, string(pubKey), fingerprint, email)
		printDBError(errUpdate, "update")
	}
	if exists == false {
		_, errInsert := db.InsertSSHKey(name, keyPath, string(pubKey), fingerprint, email)
		printDBError(errInsert, "save")
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHGenerated, name)
	fmt.Fprintf(os.Stdout, constants.MsgSSHPath, keyPath)
	fmt.Fprintf(os.Stdout, constants.MsgSSHFingerprint, fingerprint)
	if host != constants.DefaultSSHHost {
		fmt.Fprintf(os.Stdout, constants.MsgSSHHostUsed, host)
	}
	fmt.Fprint(os.Stdout, constants.MsgSSHPubLabel)
	fmt.Fprintf(os.Stdout, "  %s\n", strings.TrimSpace(string(pubKey)))
	fmt.Fprint(os.Stdout, constants.MsgSSHCopyHint)
	copyPubKeyAndAnnounce(strings.TrimSpace(string(pubKey)))

	updateSSHConfig(db)
}

func askConfirm(name, keyPath string) bool {
	fmt.Fprintf(os.Stdout, constants.MsgSSHConfirmPrompt, name, keyPath)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

func exitOnBackupError(err error) {
	if err != nil {
		panic("error")
	}
}

func printDBError(err error, action string) {
	if err == nil {
		return
	}
	if action == "update" {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not update SSH key in DB: %v\n", err)
	}
	if action == "save" {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not save SSH key to DB: %v\n", err)
	}
}
