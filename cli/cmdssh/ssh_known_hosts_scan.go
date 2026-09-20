package cmdssh

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveHostAndPort(target string) (string, int) {
	hasColon := strings.Contains(target, ":")
	if hasColon == false {
		t, err := ParseSSHTarget(target, "root", 22)
		if err == nil && t != nil && t.IP != "" {
			return t.IP, t.Port
		}
		return target, 22
	}
	host, portStr, err := net.SplitHostPort(target)
	p, convErr := strconv.Atoi(portStr)
	if err == nil && convErr == nil {
		return host, p
	}
	t, err := ParseSSHTarget(target, "root", 22)
	if err == nil && t != nil && t.IP != "" {
		return t.IP, t.Port
	}
	return target, 22
}

func buildScanClientConfig(capturedKey *ssh.PublicKey) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "gitmap-probe",
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			*capturedKey = key
			return nil
		},
		Timeout: 4 * time.Second,
	}
}

// ScanRemoteHostKey connects to remote host:port to capture its public host key.
func ScanRemoteHostKey(host string, port int) (ssh.PublicKey, error) {
	var captured ssh.PublicKey
	config := buildScanClientConfig(&captured)
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, config.Timeout)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ScanRemoteHostKey.Dial")
	}
	defer conn.Close()

	sshConn, _, _, _ := ssh.NewClientConn(conn, addr, config)
	if sshConn != nil {
		_ = sshConn.Close()
	}
	if captured != nil {
		return captured, nil
	}
	return nil, apperror.NewNotFoundError("no host key received from " + addr)
}

func buildKnownHostRecord(target string, pubKey ssh.PublicKey) store.SSHKnownHost {
	keyType := pubKey.Type()
	b64Key := base64.StdEncoding.EncodeToString(pubKey.Marshal())
	fp := ssh.FingerprintSHA256(pubKey)
	id := fmt.Sprintf("kh-%s-%s", strings.ReplaceAll(target, ",", "_"), keyType)
	now := time.Now().UTC()
	return store.SSHKnownHost{
		ID:          id,
		Host:        target,
		KeyType:     keyType,
		PublicKey:   b64Key,
		Fingerprint: fp,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func parseExplicitHostKey(target, explicitKey string) (store.SSHKnownHost, bool) {
	if explicitKey == "" {
		return store.SSHKnownHost{}, false
	}
	return parseKnownHostLine(target + " " + explicitKey)
}

func resolveOrScanRemoteHostKey(target string) (store.SSHKnownHost, error) {
	host, port := resolveHostAndPort(target)
	pubKey, scanErr := ScanRemoteHostKey(host, port)
	if scanErr != nil {
		return store.SSHKnownHost{}, scanErr
	}
	return buildKnownHostRecord(target, pubKey), nil
}

// TrustRemoteTarget scans or records a host key into known_hosts and SQLite.
func TrustRemoteTarget(ctx context.Context, target string, explicitKey string, db *sql.DB) (*store.SSHKnownHost, error) {
	kh, isParsed := parseExplicitHostKey(target, explicitKey)
	if isParsed == false {
		rec, err := resolveOrScanRemoteHostKey(target)
		if err != nil {
			return nil, err
		}
		kh = rec
	}
	return persistTrustedHost(ctx, kh, db)
}

func persistTrustedHost(ctx context.Context, kh store.SSHKnownHost, db *sql.DB) (*store.SSHKnownHost, error) {
	path, _ := DefaultKnownHostsPath()
	if path != "" {
		_ = AppendKnownHostFile(path, kh.Host, kh.KeyType, kh.PublicKey)
	}
	if db != nil {
		_ = store.UpsertSSHKnownHost(ctx, kh, db)
	}
	return &kh, nil
}

// UntrustRemoteTarget removes a host key from known_hosts and SQLite.
func UntrustRemoteTarget(ctx context.Context, target string, db *sql.DB) (int64, error) {
	path, _ := DefaultKnownHostsPath()
	if path != "" {
		_, _ = RemoveFromKnownHostsFile(path, target)
	}
	if db != nil {
		return store.DeleteSSHKnownHostByTarget(ctx, target, db)
	}
	return 0, nil
}
