package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type enrollSession struct {
	client    *ssh.Client
	osType    string
	hasClient bool
}

type hostHistoryPair struct {
	host   store.SSHHost
	hist   store.SSHHistory
	osType string
}

type keySignerResult struct {
	signer    ssh.Signer
	hasSigner bool
}

func isInteractiveTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func newAutoAcceptHostKeyConfig(user string, auth []ssh.AuthMethod) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}
}

func loadPrivateKeySigner(keyPath string) keySignerResult {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return keySignerResult{hasSigner: false}
	}
	signer, parseErr := ssh.ParsePrivateKey(keyBytes)
	if parseErr != nil {
		return keySignerResult{hasSigner: false}
	}
	return keySignerResult{signer: signer, hasSigner: true}
}

func dialNodeWithSigner(target *SSHTarget, signer ssh.Signer) *ssh.Client {
	config := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{ssh.PublicKeys(signer)})
	addr := net.JoinHostPort(target.IP, strconv.Itoa(resolveHealthPort(target.Port)))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil
	}
	return client
}

func dialNodeWithPassword(target *SSHTarget, pass string) *ssh.Client {
	config := newAutoAcceptHostKeyConfig(target.Username, []ssh.AuthMethod{ssh.Password(pass)})
	addr := net.JoinHostPort(target.IP, strconv.Itoa(resolveHealthPort(target.Port)))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil
	}
	return client
}

func tryConnectDefaultKey(target *SSHTarget) *ssh.Client {
	keyPath := findDefaultUserSSHKey()
	if keyPath == "" {
		return nil
	}
	res := loadPrivateKeySigner(keyPath)
	if !res.hasSigner {
		return nil
	}
	return dialNodeWithSigner(target, res.signer)
}

func promptUserPassword(ctx context.Context, target *SSHTarget) string {
	prompt := fmt.Sprintf("Enter SSH password for %s: ", target.String())
	pass, err := PromptSSHPassword(ctx, prompt, int(os.Stdin.Fd()))
	if err != nil {
		return ""
	}
	return pass
}

func connectWithGivenPass(target *SSHTarget, pass string) enrollSession {
	client := dialNodeWithPassword(target, pass)
	hasClient := client != nil
	return enrollSession{client: client, hasClient: hasClient, osType: probeRemoteOSType(client)}
}

func promptAndConnectTarget(ctx context.Context, opts *SSHJoinOptions) enrollSession {
	pass := promptUserPassword(ctx, opts.Target)
	if pass == "" {
		return enrollSession{hasClient: false, osType: "linux"}
	}
	opts.Password = pass
	return connectWithGivenPass(opts.Target, pass)
}

func detectAndConnectAuth(ctx context.Context, opts *SSHJoinOptions) enrollSession {
	client := tryConnectDefaultKey(opts.Target)
	if client != nil {
		return enrollSession{client: client, hasClient: true, osType: probeRemoteOSType(client)}
	}
	if !isInteractiveTerminal() {
		return enrollSession{hasClient: false, osType: "linux"}
	}
	return promptAndConnectTarget(ctx, opts)
}

func resolveTargetClient(ctx context.Context, opts *SSHJoinOptions) enrollSession {
	if !probeTCPQuick(opts.Target.IP, opts.Target.Port, 800*time.Millisecond) {
		return enrollSession{hasClient: false, osType: "linux"}
	}
	if opts.Password != "" {
		return connectWithGivenPass(opts.Target, opts.Password)
	}
	return detectAndConnectAuth(ctx, opts)
}

func probeRemoteOSType(client *ssh.Client) string {
	if client == nil {
		return "linux"
	}
	out, err := crypto.RunCommand(client, "uname -s 2>/dev/null || echo Windows", "")
	if err != nil {
		return "linux"
	}
	return resolveDetectedOSType(out)
}

func resolveDefaultOS(osType string) string {
	if isWindowsOS(osType) {
		return "windows"
	}
	return "linux"
}

func persistDualTables(ctx context.Context, dbConn *sql.DB, pair hostHistoryPair) *apperror.AppError {
	if err := store.EnrollSSHHost(ctx, pair.host, pair.hist, dbConn); err != nil {
		return apperror.WrapSimple(err, "persistDualTables.EnrollSSHHost")
	}
	conn := db.SSHConnection{
		Alias:             pair.host.Alias,
		IPAddress:         pair.host.IP,
		Username:          pair.host.Username,
		EncryptedPassword: pair.host.EncryptedPassword,
		KeyPath:           findDefaultUserSSHKey(),
		OS:                resolveDefaultOS(pair.osType),
		CreatedAt:         pair.host.CreatedAt,
	}
	return db.InsertOrUpdateSSHConnection(ctx, dbConn, conn)
}

func resolveEncryptedPassword(pass string) string {
	if pass == "" {
		return ""
	}
	enc, err := EncryptSSHPassword(pass)
	if err != nil {
		return ""
	}
	return enc
}

func buildEnrollPair(opts *SSHJoinOptions, osType string) hostHistoryPair {
	now := time.Now().UTC()
	host := buildHostRecord(opts, now)
	hist := buildHistRecord(opts, now)
	host.Port = opts.Target.Port
	host.EncryptedPassword = resolveEncryptedPassword(opts.Password)
	return hostHistoryPair{host: host, hist: hist, osType: osType}
}

func persistEnrollmentDual(ctx context.Context, opts *SSHJoinOptions, osType string) *apperror.AppError {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("persistEnrollmentDual", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()
	pair := buildEnrollPair(opts, osType)
	return persistDualTables(ctx, dbConn.SQL(), pair)
}

func isRemoteGitmapInstalled(client *ssh.Client, osType string) bool {
	shell := resolveRemoteShell(osType)
	out, err := crypto.RunCommand(client, "gitmap --version", shell)
	return err == nil && strings.Contains(out, "gitmap")
}

func bootstrapRemoteGitmap(client *ssh.Client, alias, osType string) *apperror.AppError {
	fmt.Printf("ℹ [%s] GitMap missing on remote machine, auto-bootstrapping...\n", alias)
	installCmd := BuildGitmapInstallOneLiner(osType, "latest")
	out, err := crypto.RunCommand(client, installCmd, resolveRemoteShell(osType))
	reportRemoteExecution(fmt.Sprintf("[%s]", alias), "GitMap bootstrap", out, err)
	if err != nil {
		return apperror.WrapSimple(err, "bootstrapRemoteGitmap")
	}
	return nil
}

func ensureRemoteGitmapInstalled(client *ssh.Client, alias, osType string) *apperror.AppError {
	if isRemoteGitmapInstalled(client, osType) {
		return nil
	}
	return bootstrapRemoteGitmap(client, alias, osType)
}

func deployHostPublicKey(client *ssh.Client, alias, osType string) *apperror.AppError {
	pubKey, keyPath, err := loadDefaultPublicKey()
	if err != nil {
		fmt.Printf("  [%s] Public key discovery notice: %v\n", alias, err)
		return nil
	}
	script := buildInjectAuthKeyScript(pubKey, osType)
	out, runErr := crypto.RunCommand(client, script, resolveRemoteShell(osType))
	reportRemoteExecution(fmt.Sprintf("[%s]", alias), fmt.Sprintf("Authorized key (%s)", filepath.Base(keyPath)), out, runErr)
	if runErr != nil {
		return apperror.WrapSimple(runErr, "deployHostPublicKey")
	}
	return nil
}

func confirmBidirectionalComm(client *ssh.Client, alias, osType string) bool {
	checkCmd := "gitmap version 2>/dev/null || gitmap --version 2>/dev/null"
	out, err := crypto.RunCommand(client, checkCmd, resolveRemoteShell(osType))
	hasGitmap := err == nil && strings.Contains(out, "gitmap")
	if hasGitmap {
		fmt.Printf("✓ [%s] Bidirectional communication confirmed (node GitMap active).\n", alias)
		return true
	}
	fmt.Printf("! [%s] Bidirectional check notice: %v\n", alias, err)
	return false
}

func performConnectedBootstrap(session enrollSession, opts *SSHJoinOptions) *apperror.AppError {
	if !session.hasClient || session.client == nil {
		return nil
	}
	defer session.client.Close()
	_ = ensureRemoteGitmapInstalled(session.client, opts.Alias, session.osType)
	_ = deployHostPublicKey(session.client, opts.Alias, session.osType)
	_ = confirmBidirectionalComm(session.client, opts.Alias, session.osType)
	return nil
}

func ExecuteSSHJoinEnrollment(ctx context.Context, opts *SSHJoinOptions) error {
	if opts.Target == nil {
		return apperror.NewValidationError(msgMissingJoinTarget)
	}
	session := resolveTargetClient(ctx, opts)
	if err := persistEnrollmentDual(ctx, opts, session.osType); err != nil {
		return err
	}
	_ = performConnectedBootstrap(session, opts)
	if err := pushAuthIfRequested(ctx, opts); err != nil {
		return apperror.WrapSimple(err, "pushAuthIfRequested")
	}
	printEnrollSuccess(opts.Alias, opts.Target.String())
	return nil
}

func ExecuteEnrollmentCompletion(ctx context.Context, opts *SSHJoinOptions) error {
	session := resolveTargetClient(ctx, opts)
	_ = performConnectedBootstrap(session, opts)
	if err := pushAuthIfRequested(ctx, opts); err != nil {
		return apperror.WrapSimple(err, "pushAuthIfRequested")
	}
	printEnrollSuccess(opts.Alias, opts.Target.String())
	return nil
}
