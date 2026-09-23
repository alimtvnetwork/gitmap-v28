package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

	osType, fullVersion, osArch := probeTargetOSDetails(client)
	target.DetectedOS = osType
	target.OSVersion = fullVersion
	target.OSArch = osArch

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

	persistErr := persistHostWithEncryptedPass(ctx, host, hist, osType)
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
		statusStr := constants.ColorGreen + "SUCCESS" + constants.ColorReset
		details := "enrolled"
		if !t.IsSuccess {
			statusStr = constants.ColorRed + "FAILED" + constants.ColorReset
			details = t.ErrorMsg
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			statusStr, t.FullIP, t.Alias, t.Username, t.DetectedOS, t.OSVersion, details)
	}
	_ = w.Flush()
	fmt.Fprintf(out, "\n  Batch join summary: %d/%d succeeded (%d failed)\n\n",
		res.SuccessCount, res.TotalCount, res.FailureCount)
	return nil
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
