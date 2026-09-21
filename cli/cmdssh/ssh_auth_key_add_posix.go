package cmdssh

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolvePosixAuthKeysPath() (string, error) {
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return "", apperror.WrapSimple(err, "resolvePosixAuthKeysPath.Home")
	}

	return filepath.Join(home, ".ssh", "authorized_keys"), nil
}

func installPosixAuthKeys(key, keyBlob string) (int, error) {
	path, err := resolvePosixAuthKeysPath()
	hasErr := err != nil
	if hasErr {
		return 0, err
	}

	sshDir := filepath.Dir(path)
	_ = os.MkdirAll(sshDir, 0700)
	_ = os.Chmod(sshDir, 0700)

	data, _ := os.ReadFile(path)
	hasKey := isKeyInAuthorizedKeys(string(data), keyBlob)
	if hasKey {
		fmt.Printf("  ℹ Key already present in: %s\n", path)

		return 0, nil
	}

	updated := formatKeyAppend(string(data), key)
	writeErr := os.WriteFile(path, []byte(updated), 0600)
	if writeErr != nil {
		return 0, apperror.WrapSimple(writeErr, "installPosixAuthKeys.Write")
	}

	_ = os.Chmod(path, 0600)
	fmt.Printf("  ✓ Key added to: %s\n", path)

	return 1, nil
}
