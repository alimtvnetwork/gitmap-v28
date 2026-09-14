package cmdssh

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	clusterKeypairDir        string
	ensureRSAKeypairFn       = EnsureClusterRSAKeypair
	promptPasswordFn         = PromptSSHPassword
	bootstrapInjectKeyFn     = executeInjectKey
	bootstrapInjectSudoersFn = executeInjectSudoers
	bootstrapVerifyKeyAuthFn = verifyKeyAuth
)

func resolveSSHDir() (string, error) {
	if clusterKeypairDir != "" {
		return clusterKeypairDir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "os.UserHomeDir")
	}
	return filepath.Join(home, ".ssh"), nil
}

func resolveRSAKeyPaths() (string, string, error) {
	dir, err := resolveSSHDir()
	if err != nil {
		return "", "", err
	}
	return filepath.Join(dir, "id_rsa"), filepath.Join(dir, "id_rsa.pub"), nil
}

func generateRSA4096Key() (*rsa.PrivateKey, []byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "rsa.GenerateKey")
	}
	sshPubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "ssh.NewPublicKey")
	}
	return privKey, ssh.MarshalAuthorizedKey(sshPubKey), nil
}

func encodeRSAPrivatePEM(privKey *rsa.PrivateKey) []byte {
	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	return pem.EncodeToMemory(block)
}

func writeKeyFiles(privPath, pubPath string, privPEM, pubBytes []byte) error {
	if err := os.MkdirAll(filepath.Dir(privPath), 0700); err != nil {
		return apperror.WrapSimple(err, "os.MkdirAll")
	}
	if err := os.WriteFile(privPath, privPEM, 0600); err != nil {
		return apperror.WrapSimple(err, "os.WriteFile_priv")
	}
	if err := os.WriteFile(pubPath, pubBytes, 0644); err != nil {
		return apperror.WrapSimple(err, "os.WriteFile_pub")
	}
	return nil
}

func hasExistingKeyFiles(privPath, pubPath string) bool {
	_, privErr := os.Stat(privPath)
	_, pubErr := os.Stat(pubPath)
	return privErr == nil && pubErr == nil
}

func readExistingPubKey(privPath, pubPath string) (string, string, bool, error) {
	if !hasExistingKeyFiles(privPath, pubPath) {
		return "", "", false, nil
	}
	pubData, err := os.ReadFile(pubPath)
	if err != nil {
		return "", "", false, apperror.WrapSimple(err, "os.ReadFile_pub")
	}
	return privPath, strings.TrimSpace(string(pubData)), true, nil
}

func createNewRSAKeypair(privPath, pubPath string) (string, string, error) {
	privKey, pubBytes, err := generateRSA4096Key()
	if err != nil {
		return "", "", err
	}
	privPEM := encodeRSAPrivatePEM(privKey)
	if err := writeKeyFiles(privPath, pubPath, privPEM, pubBytes); err != nil {
		return "", "", err
	}
	return privPath, strings.TrimSpace(string(pubBytes)), nil
}

// EnsureClusterRSAKeypair discovers or creates a 4096-bit RSA keypair at ~/.ssh/id_rsa.
func EnsureClusterRSAKeypair() (string, string, error) {
	privPath, pubPath, err := resolveRSAKeyPaths()
	if err != nil {
		return "", "", err
	}
	if p, pub, isFound, checkErr := readExistingPubKey(privPath, pubPath); isFound || checkErr != nil {
		return p, pub, checkErr
	}
	return createNewRSAKeypair(privPath, pubPath)
}

type clusterBootstrapOptions struct {
	target        string
	password      string
	isSudoEnabled bool
	user          string
	port          int
	isShowHelp    bool
}

type BootstrapResult struct {
	Node     string
	IP       string
	HasSudo  bool
	HasKey   bool
	Status   string
	Duration time.Duration
	Err      error
}

func isSudoToggleFlag(arg string) (bool, bool) {
	if arg == "--sudo" || arg == "-s" {
		return true, true
	}
	if arg == "--sudo=false" || arg == "--no-sudo" {
		return true, false
	}
	return false, true
}

func parseBootstrapStringFlag(args []string, idx int, long, short string) (string, int, bool) {
	arg := args[idx]
	prefix := long + "="
	if strings.HasPrefix(arg, prefix) {
		return strings.TrimPrefix(arg, prefix), 1, true
	}
	if (arg == long || arg == short) && idx+1 < len(args) {
		return args[idx+1], 2, true
	}
	return "", 0, false
}

func parseIntFlag(args []string, idx int, long, short string) (int, int, bool) {
	valStr, consumed, isMatched := parseBootstrapStringFlag(args, idx, long, short)
	if !isMatched {
		return 0, 0, false
	}
	port, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, 0, false
	}
	return port, consumed, true
}

func parseBoolFlag(arg string, opts *clusterBootstrapOptions) bool {
	if isClusterHelpFlag(arg) {
		opts.isShowHelp = true
		return true
	}
	if isMatched, isVal := isSudoToggleFlag(arg); isMatched {
		opts.isSudoEnabled = isVal
		return true
	}
	return false
}

func parseUserOrPortFlag(args []string, idx int, opts *clusterBootstrapOptions) (int, bool) {
	if val, c, isMatched := parseBootstrapStringFlag(args, idx, "--user", "-u"); isMatched {
		opts.user = val
		return c, true
	}
	if val, c, isMatched := parseIntFlag(args, idx, "--port", "-p"); isMatched {
		opts.port = val
		return c, true
	}
	return 0, false
}

func parsePassFlag(args []string, idx int, opts *clusterBootstrapOptions) (int, bool) {
	if val, c, isMatched := parseBootstrapStringFlag(args, idx, "--password", "-P"); isMatched {
		opts.password = val
		return c, true
	}
	return 0, false
}

func parseBootstrapFlagAtIndex(args []string, idx int, opts *clusterBootstrapOptions) (int, bool) {
	if parseBoolFlag(args[idx], opts) {
		return 1, true
	}
	if c, isMatched := parsePassFlag(args, idx, opts); isMatched {
		return c, true
	}
	return parseUserOrPortFlag(args, idx, opts)
}

func resolveBootstrapPositional(opts *clusterBootstrapOptions, pos []string) (*clusterBootstrapOptions, error) {
	if opts.isShowHelp {
		return opts, nil
	}
	if len(pos) == 0 {
		return nil, apperror.NewValidationError("usage: gitmap cluster bootstrap <target> [password] [flags]")
	}
	opts.target = pos[0]
	if len(pos) > 1 && opts.password == "" {
		opts.password = pos[1]
	}
	return opts, nil
}

func parseClusterBootstrapArgs(args []string) (*clusterBootstrapOptions, error) {
	opts := &clusterBootstrapOptions{isSudoEnabled: true}
	var positional []string
	idx := 0
	for idx < len(args) {
		if c, isFlag := parseBootstrapFlagAtIndex(args, idx, opts); isFlag {
			idx += c
			continue
		}
		positional = append(positional, args[idx])
		idx++
	}
	return resolveBootstrapPositional(opts, positional)
}

func resolveHostUser(user string) string {
	if user != "" {
		return user
	}
	return "root"
}

func buildSyntheticHost(t *SSHTarget) store.SSHHost {
	return store.SSHHost{
		ID:          fmt.Sprintf("host-%s", t.IP),
		Alias:       t.IP,
		IP:          t.IP,
		Username:    t.Username,
		Port:        t.Port,
		ClusterRole: "worker",
		CreatedAt:   time.Now().UTC(),
	}
}

func resolveSyntheticTarget(raw string, defaultUser string, defaultPort int) ([]store.SSHHost, error) {
	t, err := ParseSSHTarget(raw, defaultUser, defaultPort)
	if err != nil {
		return nil, apperror.Wrap(err, "ParseSSHTarget", map[string]any{"target": raw})
	}
	return []store.SSHHost{buildSyntheticHost(t)}, nil
}

func applyUserAndPortOverrides(hosts []store.SSHHost, defaultUser string, defaultPort int) []store.SSHHost {
	res := make([]store.SSHHost, len(hosts))
	for i, h := range hosts {
		if defaultUser != "" {
			h.Username = defaultUser
		}
		if defaultPort > 0 {
			h.Port = defaultPort
		}
		res[i] = h
	}
	return res
}

func resolveBootstrapTargets(ctx context.Context, target string, db *sql.DB, defaultUser string, defaultPort int) ([]store.SSHHost, error) {
	hosts, err := store.ListHostsByTarget(ctx, target, db)
	if err == nil && len(hosts) > 0 {
		return applyUserAndPortOverrides(hosts, defaultUser, defaultPort), nil
	}
	return resolveSyntheticTarget(target, defaultUser, defaultPort)
}

func wipePasswordString(p *string) {
	if p == nil || *p == "" {
		return
	}
	bytes := []byte(*p)
	for i := range bytes {
		bytes[i] = 0
	}
	*p = ""
}

func encryptHostPasswordSafely(password string) string {
	if password == "" {
		return ""
	}
	enc, err := EncryptSSHPassword(password)
	if err != nil {
		return ""
	}
	return enc
}

func resolveStoredPassword(host store.SSHHost) (string, error) {
	if host.EncryptedPassword == "" {
		return "", nil
	}
	return DecryptSSHPassword(host.EncryptedPassword)
}

func resolveBootstrapTargetPassword(ctx context.Context, host store.SSHHost, cliPass string) (string, error) {
	if cliPass != "" {
		return cliPass, nil
	}
	dbPass, err := resolveStoredPassword(host)
	if err == nil && dbPass != "" {
		return dbPass, nil
	}
	prompt := fmt.Sprintf("Enter SSH password for %s@%s: ", host.Username, host.IP)
	return promptPasswordFn(ctx, prompt, int(os.Stdin.Fd()))
}

func buildInjectPublicKeyScript(pubKey string) string {
	escapedKey := escapeSingleQuotes(pubKey)
	return fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && (grep -qF '%s' ~/.ssh/authorized_keys || echo '%s' >> ~/.ssh/authorized_keys)", escapedKey, escapedKey)
}

func buildSudoersScript(user string) string {
	return fmt.Sprintf("echo '%s ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/%s && chmod 0440 /etc/sudoers.d/%s", user, user, user)
}

func buildBatchVerifyArgs(keyPath string, host store.SSHHost) []string {
	user := resolveHostUser(host.Username)
	target := fmt.Sprintf("%s@%s", user, host.IP)
	args := []string{"-i", keyPath, "-o", "BatchMode=yes", "-o", "StrictHostKeyChecking=no"}
	if isCustomSSHPort(host.Port) {
		args = append(args, "-p", strconv.Itoa(host.Port))
	}
	args = append(args, target, "echo ssh_ok")
	return args
}

func executeInjectKey(ctx context.Context, host store.SSHHost, pubKey, password string) error {
	script := buildInjectPublicKeyScript(pubKey)
	formattedCmd := FormatClusterCommand(script, "", false)
	sshArgs := BuildNodeSSHArgs(host, formattedCmd)
	cmd := SSHExecutor(ctx, "ssh", sshArgs...)
	cleanup := attachAskPass(cmd, password)
	defer cleanup()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "executeInjectKey", map[string]any{"out": string(out)})
	}
	return nil
}

func executeInjectSudoers(ctx context.Context, host store.SSHHost, password string) error {
	user := resolveHostUser(host.Username)
	script := buildSudoersScript(user)
	formattedCmd := FormatClusterCommand(script, password, user != "root")
	cmd := SSHExecutor(ctx, "ssh", BuildNodeSSHArgs(host, formattedCmd)...)
	cleanup := attachAskPass(cmd, password)
	defer cleanup()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "executeInjectSudoers", map[string]any{"out": string(out)})
	}
	return nil
}

func verifyKeyAuth(ctx context.Context, keyPath string, host store.SSHHost) error {
	args := buildBatchVerifyArgs(keyPath, host)
	cmd := SSHExecutor(ctx, "ssh", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "verifyKeyAuth", map[string]any{"out": string(out)})
	}
	if !strings.Contains(string(out), "ssh_ok") {
		return apperror.NewExecutionError("key auth verification output did not contain ssh_ok")
	}
	return nil
}

func buildBootstrappedHost(host store.SSHHost, encPass string) store.SSHHost {
	h := host
	if h.ID == "" {
		h.ID = fmt.Sprintf("host-%s", host.IP)
	}
	if h.Alias == "" {
		h.Alias = host.IP
	}
	if encPass != "" {
		h.EncryptedPassword = encPass
	}
	return h
}

func enrollBootstrappedHost(ctx context.Context, host store.SSHHost, encPass string, db *sql.DB) error {
	now := time.Now().UTC()
	h := buildBootstrappedHost(host, encPass)
	hist := store.SSHHistory{
		ID:       fmt.Sprintf("hist-%d", now.UnixNano()),
		HostIP:   host.IP,
		JoinedAt: now,
		User:     resolveHostUser(host.Username),
	}
	return store.EnrollSSHHost(ctx, h, hist, db)
}

func finishBootstrapFailure(res BootstrapResult, err error, start time.Time) BootstrapResult {
	res.Status = "FAILED"
	res.Duration = time.Since(start)
	res.Err = err
	return res
}

func runInjectionAndSudo(ctx context.Context, host store.SSHHost, pubKey, password string, hasSudo bool) (bool, error) {
	if err := bootstrapInjectKeyFn(ctx, host, pubKey, password); err != nil {
		return false, err
	}
	if !hasSudo {
		return false, nil
	}
	if err := bootstrapInjectSudoersFn(ctx, host, password); err != nil {
		return false, err
	}
	return true, nil
}

func verifyAndEnroll(ctx context.Context, host store.SSHHost, privKeyPath, encPass string, db *sql.DB, res *BootstrapResult) error {
	if err := bootstrapVerifyKeyAuthFn(ctx, privKeyPath, host); err != nil {
		return err
	}
	res.HasKey = true
	return enrollBootstrappedHost(ctx, host, encPass, db)
}

func finishBootstrapSuccess(res BootstrapResult, start time.Time) BootstrapResult {
	res.Status = "SUCCESS"
	res.Duration = time.Since(start)
	return res
}

func executeBootstrapSteps(ctx context.Context, host store.SSHHost, privKeyPath, pubKey, password string, opts *clusterBootstrapOptions, db *sql.DB, start time.Time) BootstrapResult {
	encPass := encryptHostPasswordSafely(password)
	defer wipePasswordString(&password)
	res := BootstrapResult{Node: resolveHostAlias(host), IP: host.IP}
	hasSudo, err := runInjectionAndSudo(ctx, host, pubKey, password, opts.isSudoEnabled)
	if err != nil {
		return finishBootstrapFailure(res, err, start)
	}
	res.HasSudo = hasSudo
	if err := verifyAndEnroll(ctx, host, privKeyPath, encPass, db, &res); err != nil {
		return finishBootstrapFailure(res, err, start)
	}
	return finishBootstrapSuccess(res, start)
}

func bootstrapSingleTarget(ctx context.Context, host store.SSHHost, privKeyPath, pubKey string, opts *clusterBootstrapOptions, db *sql.DB) BootstrapResult {
	start := time.Now()
	res := BootstrapResult{
		Node: resolveHostAlias(host),
		IP:   host.IP,
	}
	password, err := resolveBootstrapTargetPassword(ctx, host, opts.password)
	if err != nil {
		res.Status = "FAILED"
		res.Duration = time.Since(start)
		res.Err = err
		return res
	}
	return executeBootstrapSteps(ctx, host, privKeyPath, pubKey, password, opts, db, start)
}

func formatBoolStatus(val bool) string {
	if val {
		return "yes"
	}
	return "no"
}

func formatBootstrapRow(res BootstrapResult) string {
	durStr := res.Duration.Round(time.Millisecond).String()
	sudoStr := formatBoolStatus(res.HasSudo)
	keyStr := formatBoolStatus(res.HasKey)
	return fmt.Sprintf("%-16s %-16s %-10s %-10s %-10s %-10s\n",
		res.Node, res.IP, sudoStr, keyStr, res.Status, durStr)
}

func RenderBootstrapSummaryTable(results []BootstrapResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-16s %-16s %-10s %-10s %-10s %-10s\n",
		"NODE", "IP", "SUDO", "KEY_AUTH", "STATUS", "DURATION"))
	sb.WriteString(strings.Repeat("-", 78))
	sb.WriteString("\n")
	for _, res := range results {
		sb.WriteString(formatBootstrapRow(res))
	}
	return sb.String()
}

func PrintBootstrapSummaryTable(results []BootstrapResult) {
	fmt.Print(RenderBootstrapSummaryTable(results))
}

func hasBootstrapFailures(results []BootstrapResult) bool {
	for _, res := range results {
		if res.Status != "SUCCESS" || res.Err != nil {
			return true
		}
	}
	return false
}

func showClusterBootstrapHelp() error {
	fmt.Println("Usage: gitmap cluster bootstrap <target> [password] [flags]")
	fmt.Println("       gitmap sj bootstrap <target> [password] [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -s, --sudo        Configure passwordless sudoers rules (default true)")
	fmt.Println("  -u, --user        Remote SSH username")
	fmt.Println("  -p, --port        Remote SSH port (default 22)")
	fmt.Println("  -P, --password    Remote user password")
	fmt.Println("  -h, --help        Show this help message")
	return nil
}

func runBootstrapOnDB(ctx context.Context, opts *clusterBootstrapOptions, privPath, pubKey string, db *sql.DB) error {
	hosts, err := resolveBootstrapTargets(ctx, opts.target, db, opts.user, opts.port)
	if err != nil {
		return err
	}
	results := make([]BootstrapResult, len(hosts))
	for i, h := range hosts {
		results[i] = bootstrapSingleTarget(ctx, h, privPath, pubKey, opts, db)
	}
	PrintBootstrapSummaryTable(results)
	if hasBootstrapFailures(results) {
		return apperror.NewExecutionError("one or more cluster nodes failed bootstrap")
	}
	return nil
}

func executeBootstrapCLI(ctx context.Context, opts *clusterBootstrapOptions) error {
	privPath, pubKey, err := ensureRSAKeypairFn()
	if err != nil {
		return err
	}
	dbConn, err := openClusterDBFunc()
	if err != nil {
		return apperror.WrapSimple(err, "openClusterDB")
	}
	defer dbConn.Close()
	defer wipePasswordString(&opts.password)
	return runBootstrapOnDB(ctx, opts, privPath, pubKey, dbConn.Conn())
}

// RunClusterBootstrapCLI bootstraps cluster nodes with SSH keys and passwordless sudoers.
func RunClusterBootstrapCLI(args []string) error {
	opts, err := parseClusterBootstrapArgs(args)
	if err != nil {
		return err
	}
	if opts.isShowHelp {
		return showClusterBootstrapHelp()
	}
	return executeBootstrapCLI(context.Background(), opts)
}

// ClusterBootstrapCmd represents the cluster bootstrap subcommand.
var ClusterBootstrapCmd = &cobra.Command{
	Use:     "bootstrap <target> [password]",
	Aliases: []string{"bs"},
	Short:   "Bootstrap cluster nodes with SSH keys and passwordless sudo",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterBootstrapCLI(args)
	},
}
