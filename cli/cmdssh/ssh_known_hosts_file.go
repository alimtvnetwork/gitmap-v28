package cmdssh

import (
	"bufio"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DefaultKnownHostsPath returns the standard ~/.ssh/known_hosts path.
func DefaultKnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "DefaultKnownHostsPath")
	}
	return filepath.Join(home, ".ssh", "known_hosts"), nil
}

func computeKeyFingerprint(keyType string, rawB64 string) string {
	rawBytes, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		h := sha256.Sum256([]byte(rawB64))
		return "SHA256:" + base64.RawStdEncoding.EncodeToString(h[:])
	}
	pubKey, parseErr := ssh.ParsePublicKey(rawBytes)
	if parseErr != nil {
		h := sha256.Sum256([]byte(rawB64))
		return "SHA256:" + base64.RawStdEncoding.EncodeToString(h[:])
	}
	return ssh.FingerprintSHA256(pubKey)
}

func parseKnownHostLine(line string) (store.SSHKnownHost, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return store.SSHKnownHost{}, false
	}
	fields := strings.Fields(trimmed)
	if len(fields) < 3 {
		return store.SSHKnownHost{}, false
	}
	host, keyType, pubKey := fields[0], fields[1], fields[2]
	comment := ""
	if len(fields) > 3 {
		comment = strings.Join(fields[3:], " ")
	}
	fp := computeKeyFingerprint(keyType, pubKey)
	id := fmt.Sprintf("kh-%s-%s", strings.ReplaceAll(host, ",", "_"), keyType)
	return store.SSHKnownHost{
		ID:          id,
		Host:        host,
		KeyType:     keyType,
		PublicKey:   pubKey,
		Fingerprint: fp,
		Comment:     comment,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, true
}

// ParseKnownHostsFile reads all known host entries from a file.
func ParseKnownHostsFile(path string) ([]store.SSHKnownHost, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ParseKnownHostsFile.Open")
	}
	defer file.Close()

	var result []store.SSHKnownHost
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, isOk := parseKnownHostLine(scanner.Text())
		if isOk {
			result = append(result, entry)
		}
	}
	return result, scanner.Err()
}

// AppendKnownHostFile appends a new host key entry to the known_hosts file.
func AppendKnownHostFile(path string, host string, keyType string, pubKey string) error {
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return apperror.WrapSimple(err, "AppendKnownHostFile.Open")
	}
	defer f.Close()

	line := fmt.Sprintf("%s %s %s\n", host, keyType, pubKey)
	_, err = f.WriteString(line)
	return err
}

func isMatchingHostLine(line, target string) bool {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return false
	}
	hosts := strings.Split(fields[0], ",")
	for _, h := range hosts {
		clean := strings.Trim(h, "[]")
		cleanTarget := strings.Trim(target, "[]")
		if strings.EqualFold(clean, cleanTarget) {
			return true
		}
	}
	return false
}

// RemoveFromKnownHostsFile removes all entries matching target from known_hosts file.
func RemoveFromKnownHostsFile(path string, target string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, apperror.WrapSimple(err, "RemoveFromKnownHostsFile.Read")
	}
	lines := strings.Split(string(data), "\n")
	var kept []string
	removed := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if isMatchingHostLine(line, target) {
			removed++
			continue
		}
		kept = append(kept, line)
	}
	out := strings.Join(kept, "\n")
	if len(kept) > 0 {
		out += "\n"
	}
	return removed, os.WriteFile(path, []byte(out), 0600)
}

// SyncKnownHostsFileWithDB synchronizes the local known_hosts file with SQLite DB.
func SyncKnownHostsFileWithDB(ctx context.Context, db *sql.DB) (int, error) {
	path, err := DefaultKnownHostsPath()
	if err != nil {
		return 0, err
	}
	entries, err := ParseKnownHostsFile(path)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		if err := store.UpsertSSHKnownHost(ctx, e, db); err == nil {
			count++
		}
	}
	return count, nil
}
