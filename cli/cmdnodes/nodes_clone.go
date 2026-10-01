// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
)

// IsNodesCloneCommand identifies if an argument invokes fleet clone.
func IsNodesCloneCommand(cmd string) (NodesCloneKind, bool) {
	low := strings.ToLower(cmd)
	switch low {
	case "clone", "clone-except-self", "clone-noself", "clones":
		return CloneKindClone, true
	case "cfr", "clone-fix-repo", "cfr-except-self":
		return CloneKindCFR, true
	case "cfrp", "clone-fix-repo-pub", "cfr-pub", "cfrp-except-self":
		return CloneKindCFRP, true
	default:
		return "", false
	}
}

// isExceptSelfToken checks if a token requests excluding local host.
func isExceptSelfToken(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))
	return low == "except-self" || low == "exceptself" || low == "no-self" || low == "noself" ||
		low == "without-self" || low == "--except-self" || low == "--exceptself" ||
		low == "--no-self" || low == "--noself" || low == "--without-self" ||
		low == "--skip-local" || low == "--remote-only"
}

func extractTargetDirFromPassArgs(passArgs []string, hasFile bool) string {
	if hasFile || len(passArgs) < 2 {
		return ""
	}
	candidate := strings.TrimSpace(passArgs[1])
	if candidate == "" || strings.HasPrefix(candidate, "-") {
		return ""
	}
	if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") || strings.HasPrefix(candidate, "git@") {
		return ""
	}
	return candidate
}

// RunNodesClone orchestrates fleet clone across local host and remote nodes.
func RunNodesClone(args []string) error {
	if len(args) == 0 {
		return PrintNodesCloneHelp(CloneKindClone)
	}
	kind, isClone := IsNodesCloneCommand(args[0])
	if !isClone {
		return PrintNodesCloneHelp(CloneKindClone)
	}
	opts, isProceed := parseNodesCloneOptions(kind, args[1:])
	if !isProceed {
		return nil
	}
	SetFleetCloneActive(true)
	defer SetFleetCloneActive(false)
	return dispatchFleetExecution(opts)
}

func parseNodesCloneOptions(kind NodesCloneKind, raw []string) (NodesCloneOptions, bool) {
	opts := NodesCloneOptions{Kind: kind, RawArgs: raw}
	if strings.Contains(strings.ToLower(string(kind)), "except-self") || strings.Contains(strings.ToLower(string(kind)), "noself") {
		opts.IsSkipLocal = true
	}
	for i := 0; i < len(raw); i++ {
		a := raw[i]
		if a == "-h" || a == "--help" || a == "help" {
			PrintNodesCloneHelp(kind)
			return opts, false
		}
		if (a == "-t" || a == "--target") && i+1 < len(raw) {
			opts.TargetFilter = raw[i+1]
			i++
			continue
		}
		if a == "--exclude" && i+1 < len(raw) {
			opts.ExcludeFilter = raw[i+1]
			i++
			continue
		}
		if isExceptSelfToken(a) {
			opts.IsSkipLocal = true
			continue
		}
		if a == "--dry-run" {
			opts.IsDryRun = true
			continue
		}
		if a == "-j" || a == "--json" {
			opts.IsJSON = true
			continue
		}
		opts.PassArgs = append(opts.PassArgs, a)
	}
	file, hasFile := DetectCloneFile(opts.PassArgs)
	opts.DetectedFile = file
	opts.HasFile = hasFile
	opts.TargetDir = extractTargetDirFromPassArgs(opts.PassArgs, opts.HasFile)
	return opts, true
}

func loadCloneFileBytes(opts NodesCloneOptions) ([]byte, string) {
	if !opts.HasFile {
		return nil, ""
	}
	b, f, err := ReadCloneFileBytes(opts.DetectedFile)
	if err != nil {
		return nil, ""
	}
	return b, f
}

func dispatchFleetExecution(opts NodesCloneOptions) error {
	conns, _ := cmdssh.FetchAllSSHConnections()
	remoteConns := filterRemoteConnections(conns, opts)
	fileBytes, fileName := loadCloneFileBytes(opts)
	if !opts.IsJSON {
		renderFleetStartBanner(os.Stdout, opts, len(remoteConns))
	}
	results := executeFleetNodesParallel(remoteConns, opts, fileBytes, fileName)
	isLocalOk, localDetails, localDur := executeLocalClone(opts)
	recordFleetTaskAudit(opts, isLocalOk)
	if opts.IsJSON {
		return emitFleetJSON(results, isLocalOk)
	}
	renderFleetResultsTable(os.Stdout, results, isLocalOk, localDetails, localDur, opts)
	return nil
}

func recordFleetTaskAudit(opts NodesCloneOptions, isLocalOk bool) {
	status := "completed"
	if !isLocalOk {
		status = "partial"
	}
	target := strings.Join(opts.PassArgs, " ")
	if target == "" && opts.HasFile {
		target = opts.DetectedFile
	}
	cmdtask.RecordTaskAudit("nodes", string(opts.Kind), target, opts.TargetDir, status)
}

func captureOutput(fn func() error) (string, error) {
	origStdout := os.Stdout
	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		return "", fn()
	}
	os.Stdout = w
	os.Stderr = w
	runErr := fn()
	_ = w.Close()
	os.Stdout = origStdout
	os.Stderr = origStderr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String(), runErr
}

func dispatchLocalKind(kind NodesCloneKind, args []string) error {
	switch kind {
	case CloneKindClone:
		return cmdclone.RunClone(args)
	case CloneKindCFR:
		return cmdclone.RunCloneFixRepo(args)
	case CloneKindCFRP:
		return cmdclone.RunCloneFixRepoPub(args)
	default:
		return cmdclone.RunClone(args)
	}
}

func containsJSONArg(args []string) bool {
	for _, a := range args {
		if a == "--json" || a == "-j" {
			return true
		}
	}
	return false
}

func prepareLocalCloneArgs(opts NodesCloneOptions) []string {
	args := append([]string{}, opts.PassArgs...)
	if len(args) == 0 && opts.HasFile {
		args = []string{opts.DetectedFile}
	}
	if !containsJSONArg(args) {
		args = append(args, "--json")
	}
	return args
}

func executeLocalClone(opts NodesCloneOptions) (bool, string, time.Duration) {
	if opts.IsSkipLocal {
		return true, "skipped local execution (except-self)", 0
	}
	args := prepareLocalCloneArgs(opts)
	start := time.Now()
	out, err := captureOutput(func() error {
		return dispatchLocalKind(opts.Kind, args)
	})
	dur := time.Since(start)
	if payload, hasJSON := cmdclone.ParseCloneJSONResponse(out); hasJSON {
		return payload.Success, payload.Message, dur
	}
	if err != nil {
		return false, err.Error(), dur
	}
	return true, "completed successfully", dur
}

func emitFleetJSON(results []RemoteCloneNodeResult, isLocalOk bool) error {
	payload := map[string]any{
		"localSuccess": isLocalOk,
		"nodes":        results,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
