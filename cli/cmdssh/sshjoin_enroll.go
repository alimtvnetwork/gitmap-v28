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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type enrollSession struct {
	client    *ssh.Client
	err       error
	osType    string
	osVersion string
	hasClient bool
}

type hostHistoryPair struct {
	host         store.SSHHost
	hist         store.SSHHistory
	osType       string
	osGroup      string
	osVersion    string
	buildVersion string
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

func promptUserPassword(ctx context.Context, target *SSHTarget) string {
	prompt := fmt.Sprintf("Enter SSH password for %s: ", target.String())
	pass, err := PromptSSHPassword(ctx, prompt, int(os.Stdin.Fd()))
	if err != nil {
		return ""
	}
	return pass
}

func connectWithGivenPass(target *SSHTarget, pass string) enrollSession {
	client, err := dialNodeWithPassword(target, pass)
	hasClient := client != nil
	osType, osVersion, _ := probeTargetWithWhichOS(client, target.IP, target.IP)

	return enrollSession{client: client, hasClient: hasClient, osType: osType, osVersion: osVersion, err: err}
}

func promptAndConnectTarget(ctx context.Context, opts *SSHJoinOptions) enrollSession {
	pass := promptUserPassword(ctx, opts.Target)
	hasPass := pass != ""
	if hasPass == false {
		return enrollSession{hasClient: false, osType: "linux"}
	}
	opts.Password = pass

	return connectWithGivenPass(opts.Target, pass)
}

func notifyPasswordEncryptedStorage(alias string) {
	fmt.Println("  ℹ Saving password as encrypted representation (RSA-OAEP/AES) in local vault.")
	hasAlias := alias != ""
	if hasAlias {
		fmt.Printf("  ℹ Review saved password anytime with: gitmap ssh pass show %s\n", alias)

		return
	}
	fmt.Println("  ℹ Review saved password anytime with: gitmap ssh pass ls")
}

func detectAndConnectAuth(ctx context.Context, opts *SSHJoinOptions) enrollSession {
	client := tryConnectDefaultKey(opts.Target)
	hasClient := client != nil
	if hasClient {
		osType, osVersion, _ := probeTargetWithWhichOS(client, opts.Alias, opts.Target.IP)

		return enrollSession{client: client, hasClient: true, osType: osType, osVersion: osVersion}
	}
	isInteractive := isInteractiveTerminal()
	if isInteractive == false {
		return enrollSession{hasClient: false, osType: "linux"}
	}

	return promptAndConnectTarget(ctx, opts)
}

func resolveTargetAddr(target *SSHTarget) string {
	hasCustomPort := target.Port > 0 && target.Port != 22
	if hasCustomPort {
		return net.JoinHostPort(target.IP, strconv.Itoa(target.Port))
	}

	return target.IP
}

func autoTrustTargetHost(ctx context.Context, target *SSHTarget) {
	isNil := target == nil || target.IP == ""
	if isNil {
		return
	}

	_ = PruneHostFromKnownHosts(target.IP, target.Port)

	addr := resolveTargetAddr(target)
	dbConn, err := openSSHDBFunc()
	if err != nil {
		_, _ = TrustRemoteTarget(ctx, addr, "", nil)

		return
	}

	defer dbConn.Close()
	_, _ = TrustRemoteTarget(ctx, addr, "", dbConn.SQL())
}

func resolveTargetClient(ctx context.Context, opts *SSHJoinOptions) enrollSession {
	trace := GetActiveSSHTrace()
	isReachable := probeTCPQuick(opts.Target.IP, opts.Target.Port, 800*time.Millisecond)
	if isReachable == false {
		err := fmt.Errorf("network dial unreachable: %s:%d (port 22 closed or timed out)", opts.Target.IP, opts.Target.Port)
		trace.AddStep("TCP Reachability", fmt.Sprintf("%s:%d unreachable", opts.Target.IP, opts.Target.Port), "FAILED", err)

		return enrollSession{hasClient: false, osType: "linux", err: err}
	}

	trace.AddStep("TCP Reachability", fmt.Sprintf("%s:%d reachable (port 22 open)", opts.Target.IP, opts.Target.Port), "SUCCESS", nil)
	autoTrustTargetHost(ctx, opts.Target)
	trace.AddStep("Host Key Trust", fmt.Sprintf("auto-trusted %s in known_hosts", resolveTargetAddr(opts.Target)), "SUCCESS", nil)

	hasPass := opts.Password != ""
	if hasPass {
		return connectWithGivenPass(opts.Target, opts.Password)
	}

	return detectAndConnectAuth(ctx, opts)
}

func probeRemoteOSType(client *ssh.Client) string {
	if client == nil {
		return "linux"
	}
	outVer, errVer := crypto.RunCommand(client, "cmd.exe /c ver", "")
	if errVer == nil && strings.Contains(strings.ToLower(outVer), "windows") {
		return "windows"
	}
	out, err := crypto.RunCommand(client, "uname -s", "")
	if err == nil {
		return resolveDetectedOSType(out)
	}
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	combined := strings.ToLower(outVer + " " + out + " " + errStr)
	if strings.Contains(combined, "windows") || strings.Contains(combined, "categoryinfo") || strings.Contains(combined, "powershell") {
		return "windows"
	}
	return "linux"
}

func probeRemoteOSVersion(client *ssh.Client, osType string) string {
	isNil := client == nil
	if isNil {
		return ""
	}
	isWin := isWindowsOS(osType)
	if isWin {
		out, _ := crypto.RunCommand(client, "powershell -NoProfile -Command \"(Get-CimInstance Win32_OperatingSystem).Caption\" 2>nul || ver", "")

		return strings.TrimSpace(out)
	}
	out, _ := crypto.RunCommand(client, "grep PRETTY_NAME /etc/os-release 2>/dev/null | cut -d= -f2 | tr -d '\"' || uname -srm", "")

	return strings.TrimSpace(out)
}

func probeRemoteOSArch(client *ssh.Client, osType string) string {
	if client == nil {
		return "amd64"
	}
	if isWindowsOS(osType) {
		return probeWindowsArch(client)
	}
	return probeUnixArch(client)
}

func probeWindowsArch(client *ssh.Client) string {
	out, _ := crypto.RunCommand(client, "powershell -NoProfile -Command \"$env:PROCESSOR_ARCHITECTURE\" 2>nul || echo %PROCESSOR_ARCHITECTURE%", "")
	trimmed := strings.ToLower(strings.TrimSpace(out))
	if trimmed == "" {
		return "amd64"
	}
	return trimmed
}

func probeUnixArch(client *ssh.Client) string {
	out, _ := crypto.RunCommand(client, "uname -m 2>/dev/null || echo amd64", "")
	trimmed := strings.ToLower(strings.TrimSpace(out))
	if trimmed == "x86_64" || trimmed == "" {
		return "amd64"
	}
	return trimmed
}

func formatOSVersionWithArch(osVer, osArch string) string {
	cleanVer := strings.TrimSpace(osVer)
	cleanArch := strings.TrimSpace(osArch)
	if cleanVer == "" && cleanArch != "" {
		return cleanArch
	}
	if cleanVer == "" {
		return ""
	}
	if cleanArch == "" || strings.Contains(strings.ToLower(cleanVer), cleanArch) {
		return cleanVer
	}
	return fmt.Sprintf("%s (%s)", cleanVer, cleanArch)
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
		OSGroup:           pair.osGroup,
		OSVersion:         pair.osVersion,
		BuildVersion:      pair.buildVersion,
		FirstRunAt:        pair.host.CreatedAt,
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

func buildEnrollPair(opts *SSHJoinOptions, osType, osVersion string) hostHistoryPair {
	now := time.Now().UTC()
	host := buildHostRecord(opts, now)
	hist := buildHistRecord(opts, now)
	host.Port = opts.Target.Port
	host.EncryptedPassword = resolveEncryptedPassword(opts.Password)

	return hostHistoryPair{host: host, hist: hist, osType: osType, osVersion: osVersion}
}

func persistEnrollmentDual(ctx context.Context, opts *SSHJoinOptions, osType, osVersion string) *apperror.AppError {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("persistEnrollmentDual", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()
	pair := buildEnrollPair(opts, osType, osVersion)

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
	trace := GetActiveSSHTrace()

	if err := ensureRemoteGitmapInstalled(session.client, opts.Alias, session.osType); err != nil {
		trace.AddStep("Remote GitMap Check", "bootstrap notice: "+err.Error(), "FAILED", err)
	} else {
		trace.AddStep("Remote GitMap Check", "gitmap available on remote", "SUCCESS", nil)
		reProfileTargetPostInstall(session.client, opts)
	}

	if err := deployHostPublicKey(session.client, opts.Alias, session.osType); err != nil {
		trace.AddStep("Deploy Host Key", "public key notice: "+err.Error(), "FAILED", err)
	} else {
		trace.AddStep("Deploy Host Key", "public key deployed", "SUCCESS", nil)
	}

	if ok := confirmBidirectionalComm(session.client, opts.Alias, session.osType); ok {
		trace.AddStep("Bidirectional Comm", "communication confirmed", "SUCCESS", nil)
	} else {
		trace.AddStep("Bidirectional Comm", "communication check notice", "FAILED", nil)
	}

	return nil
}

func reProfileTargetPostInstall(client *ssh.Client, opts *SSHJoinOptions) {
	rep, hasGitmap := probeRemoteGitmapWhichOS(client)
	if !hasGitmap || rep == nil {
		return
	}
	fmt.Printf("✔ Node %s: Profiled via GitMap which-os: %s (%s, %s, %s)\n",
		opts.Alias, rep.OSType, rep.OSGroup, rep.OSVersion, rep.Architecture)
	savePostInstallProfile(opts, rep)
}

func savePostInstallProfile(opts *SSHJoinOptions, rep *cmdos.OSInfoReport) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return
	}
	defer dbConn.Close()
	now := time.Now().UTC()
	conn := db.SSHConnection{
		Alias:             opts.Alias,
		IPAddress:         opts.Target.IP,
		Username:          opts.Target.Username,
		EncryptedPassword: resolveEncryptedPassword(opts.Password),
		KeyPath:           findDefaultUserSSHKey(),
		OS:                rep.OSType,
		OSGroup:           rep.OSGroup,
		OSVersion:         rep.OSVersion,
		BuildVersion:      rep.BuildVersion,
		FirstRunAt:        now,
		CreatedAt:         now,
	}
	_ = db.InsertOrUpdateSSHConnection(context.Background(), dbConn.SQL(), conn)
}

func checkEnrollAuth(opts *SSHJoinOptions, session enrollSession) error {
	isFailedAuth := opts.Password != "" && !session.hasClient
	if !isFailedAuth {
		return nil
	}

	trace := GetActiveSSHTrace()
	err := session.err
	if err == nil {
		err = fmt.Errorf("ssh: authentication failed: invalid password or remote server rejected credentials")
	}

	hint := buildAuthFailureHint(opts, err)
	trace.SetInternalError(err, hint)
	logPath := trace.PersistLog()

	printAuthFailureReport(err, trace, logPath)

	ctx := map[string]any{
		"target":    opts.Target.String(),
		"raw_error": err.Error(),
		"log_path":  logPath,
		"user":      opts.Target.Username,
		"host":      opts.Target.IP,
		"port":      opts.Target.Port,
	}

	return apperror.Wrap(err, "ExecuteSSHJoinEnrollment.auth", ctx)
}

func buildAuthFailureHint(opts *SSHJoinOptions, err error) string {
	if isNetworkDialFailure(err) {
		return fmt.Sprintf("Verify remote host %s is online, port %d open, and firewall allows SSH.", opts.Target.IP, opts.Target.Port)
	}

	return fmt.Sprintf("Remote host %s rejected credentials for '%s'. Check password or sshd_config on remote host.", opts.Target.IP, opts.Target.Username)
}

func printAuthFailureReport(err error, trace *SSHExecutionTrace, logPath string) {
	fmt.Fprintf(os.Stderr, "  ⚠ Failed to authenticate with remote machine: %v\n", err)
	if trace != nil {
		fmt.Fprint(os.Stderr, trace.FormatTerminalReport())
	}
	fmt.Fprintf(os.Stderr, "  ℹ To share or inspect full details, view: %s\n", logPath)
	fmt.Fprintf(os.Stderr, "    Or run: gitmap ssh error-logs\n")
}

func finalizeEnrollment(ctx context.Context, opts *SSHJoinOptions, session enrollSession) error {
	if err := persistEnrollmentDual(ctx, opts, session.osType, session.osVersion); err != nil {
		return err
	}

	if opts.Password != "" {
		notifyPasswordEncryptedStorage(opts.Alias)
	}

	_, _ = RecordSSHAddNodeTask(ctx, buildHostRecord(opts, time.Now().UTC()))
	_ = performConnectedBootstrap(session, opts)
	if err := pushAuthIfRequested(ctx, opts); err != nil {
		return apperror.WrapSimple(err, "pushAuthIfRequested")
	}

	printEnrollSuccess(opts.Alias, opts.Target.String())

	return nil
}

func ExecuteSSHJoinEnrollment(ctx context.Context, opts *SSHJoinOptions) error {
	if opts.Target == nil {
		return apperror.NewValidationError(msgMissingJoinTarget)
	}

	_ = BeginSSHTrace("ssh join", opts.Target.String(), opts.Target.Username, opts.Target.IP, opts.Target.Port)

	session := resolveTargetClient(ctx, opts)
	if err := checkEnrollAuth(opts, session); err != nil {
		return err
	}

	return finalizeEnrollment(ctx, opts, session)
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
