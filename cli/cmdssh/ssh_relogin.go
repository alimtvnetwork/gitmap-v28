package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AutoReloginAndRestoreKey locates, decrypts, and restores an SSH key to disk,
// re-binding git SSH command for transparent recovery.
func AutoReloginAndRestoreKey(keyName string) (*model.SSHKey, error) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "AutoReloginAndRestoreKey.openDB")
	}

	defer dbConn.Close()

	return restoreKeyFromConnection(dbConn, keyName)
}

func restoreKeyFromConnection(db *store.DB, keyName string) (*model.SSHKey, error) {
	key, err := findKeyForRelogin(db, keyName)
	if err != nil {
		return nil, err
	}

	if err := RestoreSSHKeyToDisk(key); err != nil {
		return nil, err
	}

	_ = RebindGitSSHCommand(key.PrivatePath)

	return &key, nil
}

func findKeyForRelogin(db *store.DB, keyName string) (model.SSHKey, error) {
	hasName := keyName != "" && keyName != "default"
	if hasName {
		return db.FindSSHKeyByName(keyName)
	}

	defaultKey, err := db.GetDefaultSSHKey()
	hasDefault := err == nil && defaultKey.Name != ""
	if hasDefault {
		return defaultKey, nil
	}

	return findFallbackKey(db)
}

func findFallbackKey(db *store.DB) (model.SSHKey, error) {
	keys, err := db.ListSSHKeys()
	if err != nil {
		return model.SSHKey{}, err
	}

	hasKeys := len(keys) > 0
	if hasKeys {
		return keys[0], nil
	}

	return model.SSHKey{}, apperror.NewNotFoundError("no SSH keys found in vault")
}

// RestoreSSHKeyToDisk writes decrypted private and public keys to disk with strict permissions.
func RestoreSSHKeyToDisk(key model.SSHKey) error {
	privBytes, err := resolveKeyPrivateBytes(key)
	if err != nil {
		return err
	}

	if err := writeKeyFiles(key, privBytes); err != nil {
		return err
	}

	injectKeyToAgent(key.PrivatePath)

	return nil
}

func resolveKeyPrivateBytes(key model.SSHKey) ([]byte, error) {
	hasEnc := key.EncryptedPrivateKey != ""
	if hasEnc {
		return DecryptPrivateKey(key.EncryptedPrivateKey)
	}

	hasDisk := isFileExisting(key.PrivatePath)
	if hasDisk {
		return os.ReadFile(key.PrivatePath)
	}

	return nil, apperror.NewNotFoundError("no private key material available for " + key.Name)
}

func writeKeyFiles(key model.SSHKey, privBytes []byte) error {
	dir := filepath.Dir(key.PrivatePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return apperror.WrapSimple(err, "writeKeyFiles.MkdirAll")
	}

	if err := os.WriteFile(key.PrivatePath, privBytes, 0600); err != nil {
		return apperror.WrapSimple(err, "writeKeyFiles.WritePrivate")
	}

	enforceFilePermissions(key.PrivatePath)
	writePublicKeyFile(key)

	return nil
}

func enforceFilePermissions(path string) {
	isWin := runtime.GOOS == "windows"
	if isWin {
		applyWindowsUserACL(path)

		return
	}

	_ = os.Chmod(path, 0600)
}

func writePublicKeyFile(key model.SSHKey) {
	hasPub := key.PublicKey != ""
	if hasPub {
		pubPath := key.PrivatePath + ".pub"
		_ = os.WriteFile(pubPath, []byte(key.PublicKey), 0644)
	}
}

func injectKeyToAgent(privatePath string) {
	hasSock := os.Getenv("SSH_AUTH_SOCK") != ""
	if hasSock {
		_ = exec.Command("ssh-add", privatePath).Run()
	}
}

func isFileExisting(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// CheckAndRecoverMissingKey checks if private key file exists; if missing, restores from vault.
func CheckAndRecoverMissingKey(targetPath string) (bool, error) {
	hasFile := isFileExisting(targetPath)
	if hasFile {
		return false, nil
	}

	_, err := AutoReloginAndRestoreKey(filepath.Base(targetPath))
	if err != nil {
		return false, err
	}

	return true, nil
}

// HandleSSHAuthFailure executes transparent auto-relogin upon authentication failure.
func HandleSSHAuthFailure(keyName string) error {
	_, err := AutoReloginAndRestoreKey(keyName)

	return err
}

// RebindGitSSHCommand pins core.sshCommand and sets GIT_SSH_COMMAND env for transparent recovery.
func RebindGitSSHCommand(privatePath string) error {
	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", privatePath)
	os.Setenv("GIT_SSH_COMMAND", sshCmd)

	isRepo := isGitRepoCWD()
	if isRepo {
		cfgCmd := fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", privatePath)
		_ = exec.Command("git", "config", "core.sshCommand", cfgCmd).Run()
	}

	return nil
}
