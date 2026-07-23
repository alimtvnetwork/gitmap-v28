// Package cmd — sshexisting.go handles the case where an SSH key already
// exists on disk when `gitmap ssh` is invoked. Instead of forwarding the
// stdin "Overwrite (y/n)?" prompt to `ssh-keygen` (which fails non-interactively
// and confuses users), we detect the existing key UP FRONT, print the public
// key + fingerprint, and exit cleanly. Pass `--force` to regenerate.
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// keyExistsOnDisk reports whether the private key file at keyPath exists.
// We check the private key (not .pub) because a missing private + present
// public would still cause ssh-keygen to refuse generation.
func keyExistsOnDisk(keyPath string) bool {
	_, err := os.Stat(keyPath)

	return err == nil
}

// printExistingKeyOnDisk reads the public key from keyPath+".pub" and prints
// it along with the fingerprint, file paths, and a "copy to GitHub" hint.
// Also upserts the key into the gitmap database so subsequent `ssh-cat` /
// `ssh-list` calls find it.
func printExistingKeyOnDisk(db *store.DB, name, keyPath, host string) {
	pub, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHReadPub, keyPath+".pub", err)
		fmt.Fprint(os.Stderr, constants.MsgSSHForceHint)
		os.Exit(1)
	}
	fingerprint := readFingerprint(keyPath)

	fmt.Fprintf(os.Stdout, constants.MsgSSHExistsOnDisk, keyPath)
	fmt.Fprintf(os.Stdout, constants.MsgSSHPath, keyPath)
	fmt.Fprintf(os.Stdout, constants.MsgSSHFingerprint, fingerprint)
	if host != constants.DefaultSSHHost {
		fmt.Fprintf(os.Stdout, constants.MsgSSHHostUsed, host)
	}
	fmt.Fprint(os.Stdout, constants.MsgSSHPubLabel)
	fmt.Fprintf(os.Stdout, "  %s\n", strings.TrimSpace(string(pub)))
	fmt.Fprint(os.Stdout, constants.MsgSSHCopyHint)
	fmt.Fprint(os.Stdout, constants.MsgSSHForceHint)

	upsertExistingKeyToDB(db, name, keyPath, string(pub), fingerprint)
}

// upsertExistingKeyToDB stores the disk-discovered key in the gitmap database
// so `ssh-cat` and `ssh-list` can find it later. Failures are logged but do
// not exit — the user already got their public key on stdout.
func upsertExistingKeyToDB(db *store.DB, name, keyPath, pub, fingerprint string) {
	email := resolveGitEmail()
	if db.SSHKeyExists(name) {
		if err := db.UpdateSSHKey(name, keyPath, pub, fingerprint, email); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not update SSH key in DB: %v\n", err)
		}

		return
	}
	if _, err := db.InsertSSHKey(name, keyPath, pub, fingerprint, email); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not register existing SSH key in DB: %v\n", err)
	}
}

// backupKeyForRegenerate is reserved for the --force flow: rename the existing
// key + .pub to *.bak.<unix-timestamp> before regenerating. Currently unused
// by the default path but kept here for the planned --force regenerate flow.
func backupKeyForRegenerate(keyPath string) error {
	stamp := time.Now().Unix()
	suffix := fmt.Sprintf(".bak.%d", stamp)
	if err := os.Rename(keyPath, keyPath+suffix); err != nil {
		return fmt.Errorf("backup private key: %w", err)
	}
	if _, err := os.Stat(keyPath + ".pub"); err == nil {
		if err := os.Rename(keyPath+".pub", keyPath+".pub"+suffix); err != nil {
			return fmt.Errorf("backup public key: %w", err)
		}
	}

	return nil
}

// (backupKeyForRegenerate is consumed by sshgen.go in the --force branch.)
