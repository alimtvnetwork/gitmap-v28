package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

func processCommonTarget(ctx context.Context, target *SSHCommonTarget, password, encPass string, isDryRun bool) {
	if isDryRun {
		target.IsSuccess = true
		target.DetectedOS = "preview"
		return
	}

	client, dialErr := crypto.ConnectWithPassword(target.FullIP, target.Username, password)
	if dialErr != nil || client == nil {
		target.IsSuccess = false
		target.ErrorMsg = resolveDialError(dialErr)
		return
	}
	defer client.Close()

	rep := probeTargetReport(client, target.Alias, target.FullIP)
	fullVersion := formatOSVersionWithArch(rep.OSVersion, rep.Architecture)
	target.DetectedOS = rep.OSType
	target.OSVersion = fullVersion
	target.OSArch = rep.Architecture

	host := store.SSHHost{
		Alias:             target.Alias,
		IP:                target.FullIP,
		Username:          target.Username,
		Port:              target.Port,
		EncryptedPassword: encPass,
		CreatedAt:         time.Now().UTC(),
	}
	hist := store.SSHHistory{
		HostIP:   target.FullIP,
		JoinedAt: time.Now().UTC(),
		User:     target.Username,
	}

	persistErr := persistHostWithDetails(
		ctx, host, hist, rep.OSType, rep.OSGroup, fullVersion, rep.BuildVersion,
	)
	if persistErr != nil {
		target.IsSuccess = false
		target.ErrorMsg = persistErr.Error()
		return
	}

	target.IsSuccess = true
}

func resolveDialError(err error) string {
	if err == nil {
		return "connection failed"
	}
	return err.Error()
}

func probeTargetReport(client *ssh.Client, alias, ip string) *cmdos.OSInfoReport {
	rep, hasGitmap := probeRemoteGitmapWhichOS(client)
	if hasGitmap && rep != nil {
		fmt.Printf("✔ Node %s: Identified via GitMap which-os: %s (%s, %s, %s)\n",
			alias, rep.OSType, rep.OSGroup, rep.OSVersion, rep.Architecture)
		return rep
	}
	notifyRemoteGitmapAdvice(alias, ip)
	osType, fullVer, arch := probeTargetOSDetails(client)
	return &cmdos.OSInfoReport{
		OSType:       osType,
		OSGroup:      resolveOSGroupFromType(osType),
		OSVersion:    fullVer,
		Architecture: arch,
	}
}

func probeTargetWithWhichOS(client *ssh.Client, alias, ip string) (string, string, string) {
	rep := probeTargetReport(client, alias, ip)
	fullVer := formatOSVersionWithArch(rep.OSVersion, rep.Architecture)
	return rep.OSType, fullVer, rep.Architecture
}

func probeRemoteGitmapWhichOS(client *ssh.Client) (*cmdos.OSInfoReport, bool) {
	if client == nil {
		return nil, false
	}
	cmd := "gitmap which-os --json 2>/dev/null || powershell -NoProfile -Command \"gitmap which-os --json\" 2>$null || \"$env:LOCALAPPDATA\\gitmap-cli\\gitmap.exe\" which-os --json 2>$null || ~/.local/bin/gitmap which-os --json 2>/dev/null"
	out, err := crypto.RunCommand(client, cmd, "")
	if err != nil {
		return nil, false
	}
	clean := strings.TrimSpace(out)
	idx := strings.Index(clean, "{")
	lastIdx := strings.LastIndex(clean, "}")
	if idx < 0 || lastIdx <= idx {
		return nil, false
	}
	var report cmdos.OSInfoReport
	if err := json.Unmarshal([]byte(clean[idx:lastIdx+1]), &report); err != nil {
		return nil, false
	}
	if report.OSType == "" {
		return nil, false
	}
	return &report, true
}

func notifyRemoteGitmapAdvice(alias, ip string) {
	fmt.Printf("ℹ Node %s (%s): GitMap not yet installed. Install via: gitmap ssh deploy %s or gitmap ssh install gitmap -t %s\n",
		alias, ip, alias, alias)
}

func probeTargetOSDetails(client *ssh.Client) (string, string, string) {
	osType := probeRemoteOSType(client)
	osVer := probeRemoteOSVersion(client, osType)
	osArch := probeRemoteOSArch(client, osType)
	fullVersion := formatOSVersionWithArch(osVer, osArch)
	return osType, fullVersion, osArch
}

func countJoinOutcomes(targets []SSHCommonTarget) (int, int) {
	success := 0
	failure := 0
	for _, t := range targets {
		if t.IsSuccess {
			success++
		} else {
			failure++
		}
	}
	return success, failure
}

func renderCommonJoinTable(out io.Writer, res *SSHCommonJoinResult) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "STATUS\tIP\tALIAS\tUSER\tOS\tVERSION\tDETAILS")
	for _, t := range res.Targets {
		renderCommonJoinTargetRow(w, t)
	}
	_ = w.Flush()
	fmt.Fprintf(out, "\n  Batch join summary: %d/%d succeeded (%d failed)\n\n",
		res.SuccessCount, res.TotalCount, res.FailureCount)
	return nil
}

func renderCommonJoinTargetRow(w io.Writer, t SSHCommonTarget) {
	statusStr := constants.ColorRed + "FAILED" + constants.ColorReset
	details := t.ErrorMsg
	if t.IsSuccess {
		statusStr = constants.ColorGreen + "SUCCESS" + constants.ColorReset
		details = "enrolled"
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		statusStr, t.FullIP, t.Alias, t.Username, t.DetectedOS, t.OSVersion, details)
}

func renderCommonJoinJSON(out io.Writer, res *SSHCommonJoinResult) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(res)
}

func resolveCommonPassword(ctx context.Context, pwd string, isDryRun bool) (string, *apperror.AppError) {
	if pwd != "" || isDryRun {
		return pwd, nil
	}
	prompted, err := promptPasswordIfMissing(ctx, "")
	if err != nil {
		return "", apperror.WrapSimple(err, "prompt password")
	}
	return prompted, nil
}

func resolveCommonEncryptedPassword(pwd string, isDryRun bool) (string, *apperror.AppError) {
	if isDryRun || pwd == "" {
		return "", nil
	}
	encrypted, err := EncryptSSHPassword(pwd)
	if err != nil {
		return "", apperror.WrapSimple(err, "encrypt password")
	}
	return encrypted, nil
}

// ExecuteCommonJoin runs batch SSH join across parsed targets using shared credentials.
func ExecuteCommonJoin(ctx context.Context, out io.Writer, opts SSHCommonJoinOptions) (*SSHCommonJoinResult, *apperror.AppError) {
	tokens := ExtractRawTokens(opts.RawIPs)
	targets, parseErr := ParseCommonIPTokens(opts.Username, opts.Port, tokens)
	if parseErr != nil {
		return nil, parseErr
	}

	password, pwdErr := resolveCommonPassword(ctx, opts.Password, opts.DryRun)
	if pwdErr != nil {
		return nil, pwdErr
	}

	encPass, encErr := resolveCommonEncryptedPassword(password, opts.DryRun)
	if encErr != nil {
		return nil, encErr
	}

	for i := range targets {
		processCommonTarget(ctx, &targets[i], password, encPass, opts.DryRun)
	}

	succ, fail := countJoinOutcomes(targets)
	result := &SSHCommonJoinResult{
		Username:     opts.Username,
		TotalCount:   len(targets),
		SuccessCount: succ,
		FailureCount: fail,
		Targets:      targets,
	}

	if opts.IsJSON {
		_ = renderCommonJoinJSON(out, result)
		return result, nil
	}

	_ = renderCommonJoinTable(out, result)
	return result, nil
}
