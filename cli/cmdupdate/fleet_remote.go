package cmdupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

func executeDefaultRemoteUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	if opts.IsDryRun {
		return `{"success": true, "details": "dry-run simulated"}`, nil
	}
	if opts.IsZip {
		return ExecuteFleetZipUpdateFn(target, opts)
	}
	if out, ok := tryRestFleetUpdate(target, opts); ok {
		return out, nil
	}
	return executeSSHFleetUpdate(target, opts)
}

func tryRestFleetUpdate(target FleetTarget, opts FleetUpdateOptions) (string, bool) {
	url := fmt.Sprintf("http://%s:49152/api/v1/update", target.IP)
	payload := map[string]any{
		"package": opts.Pkg,
		"force":   opts.IsForce,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}
	client := http.Client{Timeout: 3000 * time.Millisecond}
	resp, reqErr := client.Post(url, "application/json", strings.NewReader(string(b)))
	if reqErr != nil {
		return "", false
	}
	defer resp.Body.Close()
	isOk := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated
	if !isOk {
		return "", false
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", false
	}
	return string(body), true
}

func executeSSHFleetUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	osType := target.OS
	if probed := cmdssh.ProbeRemoteOSType(client); probed != "" {
		osType = probed
	}

	cmd := resolveFleetUpdateCommand(osType, opts.Pkg)
	shell := resolveFleetShell(osType)
	out, err := crypto.RunCommand(client, cmd, shell)
	if err == nil && isAgmPkg(opts.Pkg) {
		collectAgmUpdateZip(client, osType)
	}
	return out, err
}

func executeSSHFleetZipUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	osType := target.OS
	if probed := cmdssh.ProbeRemoteOSType(client); probed != "" {
		osType = probed
	}

	zipData, err := getCachedUpdateZip(opts.Pkg, osType)
	if err != nil {
		return "", apperror.WrapSimple(err, "getCachedUpdateZip")
	}

	destZipPath := resolveRemoteZipPath(osType, opts.Pkg)
	err = StreamFileToRemoteFn(client, destZipPath, zipData, osType)
	if err != nil {
		return "", apperror.WrapSimple(err, "StreamFileToRemote")
	}

	cmd := resolveFleetZipInstallCommand(osType, opts.Pkg, destZipPath)
	shell := resolveFleetShell(osType)
	out, err := crypto.RunCommand(client, cmd, shell)
	if err == nil && isAgmPkg(opts.Pkg) {
		collectAgmUpdateZip(client, osType)
	}
	return out, err
}

var (
	zipCacheMu sync.Mutex
	zipCache   = make(map[string][]byte)
)
