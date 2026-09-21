package cmdssh

import (
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func validateSSHPublicKey(raw string) (string, string, error) {
	trimmed := strings.TrimSpace(raw)
	isEmpty := len(trimmed) == 0
	if isEmpty {
		return "", "", apperror.NewValidationError("SSH public key cannot be empty")
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(trimmed))
	if err != nil {
		return "", "", apperror.WrapSimple(err, "validateSSHPublicKey.Parse")
	}

	b64Blob := base64.StdEncoding.EncodeToString(pubKey.Marshal())

	return trimmed, b64Blob, nil
}

func isKeyInAuthorizedKeys(content, keyBlob string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		isCommentOrEmpty := trimmed == "" || strings.HasPrefix(trimmed, "#")
		if isCommentOrEmpty {
			continue
		}

		tokens := strings.Fields(trimmed)
		hasTokens := len(tokens) >= 2
		if hasTokens && tokens[1] == keyBlob {
			return true
		}
	}

	return false
}

func formatKeyAppend(existing, key string) string {
	isEmpty := len(existing) == 0
	if isEmpty {
		return key + "\n"
	}

	hasTrailingNL := strings.HasSuffix(existing, "\n")
	if hasTrailingNL {
		return existing + key + "\n"
	}

	return existing + "\n" + key + "\n"
}
