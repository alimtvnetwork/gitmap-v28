// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

func isExceptSelfToken(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))
	return low == "except-self" || low == "exceptself" || low == "no-self" || low == "noself" ||
		low == "without-self" || low == "--except-self" || low == "--exceptself" ||
		low == "--no-self" || low == "--noself" || low == "--without-self" ||
		low == "--skip-local" || low == "--remote-only"
}

func isExcludeAlias(arg string) bool {
	return arg == "--exclude" || arg == "--except" || arg == "--excep" || arg == "--accept"
}

func parseExcludeFlag(raw []string, idx int) (string, int, bool) {
	arg := raw[idx]
	for _, prefix := range []string{"--exclude=", "--except=", "--excep=", "--accept="} {
		if strings.HasPrefix(arg, prefix) {
			return strings.TrimPrefix(arg, prefix), idx, true
		}
	}
	if isExcludeAlias(arg) && idx+1 < len(raw) {
		return raw[idx+1], idx + 1, true
	}
	return "", idx, false
}

func parseTargetFlag(raw []string, idx int) (string, int, bool) {
	arg := raw[idx]
	if strings.HasPrefix(arg, "--target=") {
		return strings.TrimPrefix(arg, "--target="), idx, true
	}
	if (arg == "-t" || arg == "--target") && idx+1 < len(raw) {
		return raw[idx+1], idx + 1, true
	}
	return "", idx, false
}

func isDestFlag(arg string) bool {
	return arg == "-d" || arg == "--dest" || arg == "--dir" || arg == "--target-dir"
}

func cleanTargetDir(dir string) string {
	cleaned := strings.TrimRight(strings.TrimSpace(dir), `/\`)
	if cleaned == "" {
		return dir
	}
	return cleaned
}

func parseDestFlag(raw []string, idx int) (string, int, bool) {
	arg := raw[idx]
	for _, prefix := range []string{"--dest=", "--dir=", "--target-dir=", "-d="} {
		if strings.HasPrefix(arg, prefix) {
			return cleanTargetDir(strings.TrimPrefix(arg, prefix)), idx, true
		}
	}
	if isDestFlag(arg) && idx+1 < len(raw) {
		return cleanTargetDir(raw[idx+1]), idx + 1, true
	}
	return "", idx, false
}

func parseBoolFlags(opts *NodesCloneOptions, a string) bool {
	if isExceptSelfToken(a) {
		opts.IsSkipLocal = true
		return true
	}
	if a == "--dry-run" {
		opts.IsDryRun = true
		return true
	}
	if a == "-j" || a == "--json" {
		opts.IsJSON = true
		return true
	}
	return false
}

func parseTargetOrExclude(opts *NodesCloneOptions, raw []string, i int) (int, bool) {
	if val, nextI, hasVal := parseTargetFlag(raw, i); hasVal {
		opts.TargetFilter = val
		return nextI, true
	}
	if val, nextI, hasVal := parseExcludeFlag(raw, i); hasVal {
		opts.ExcludeFilter = val
		return nextI, true
	}
	return i, false
}

func parseParamFlags(opts *NodesCloneOptions, raw []string, i int) (int, bool) {
	if nextI, hasVal := parseTargetOrExclude(opts, raw, i); hasVal {
		return nextI, true
	}
	if val, nextI, hasVal := parseDestFlag(raw, i); hasVal {
		opts.TargetDir = val
		opts.HasCustomTargetDir = true
		return nextI, true
	}
	return i, false
}

func parseSingleToken(opts *NodesCloneOptions, raw []string, i int) (int, bool) {
	if nextI, hasParam := parseParamFlags(opts, raw, i); hasParam {
		return nextI, true
	}
	if parseBoolFlags(opts, raw[i]) {
		return i, true
	}
	opts.PassArgs = append(opts.PassArgs, raw[i])
	return i, true
}

func parseTokens(opts *NodesCloneOptions, raw []string) bool {
	for i := 0; i < len(raw); i++ {
		if raw[i] == "-h" || raw[i] == "--help" || raw[i] == "help" {
			return false
		}
		nextI, _ := parseSingleToken(opts, raw, i)
		i = nextI
	}
	return true
}

func extractTargetDirFromPassArgs(passArgs []string, hasFile bool) string {
	if hasFile || len(passArgs) < 2 {
		return ""
	}
	candidate := strings.TrimSpace(passArgs[1])
	if candidate == "" || strings.HasPrefix(candidate, "-") {
		return ""
	}
	if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") || strings.HasPrefix(candidate, "git@") || strings.Contains(candidate, ",") {
		return ""
	}
	return candidate
}

func defaultOSWorkBase() string {
	if runtime.GOOS == "windows" {
		return `D:\work`
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, "work")
	}
	return "~/work"
}

func resolveLocalWorkBase() string {
	db, err := store.OpenDefault()
	if err != nil {
		return defaultOSWorkBase()
	}
	defer db.Close()

	wd, errGet := db.GetDefaultWorkDir()
	if errGet == nil && wd != nil && wd.AbsolutePath != "" {
		return filepath.Clean(wd.AbsolutePath)
	}

	return defaultOSWorkBase()
}

func inspectCwdRelativeWorkDir(workBase string) (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	rel, errRel := filepath.Rel(filepath.Clean(workBase), filepath.Clean(cwd))
	if errRel != nil || strings.HasPrefix(rel, "..") {
		return "", false
	}
	if rel == "." {
		return "", true
	}
	return filepath.ToSlash(rel), true
}

func applyPositionalTargetDir(opts *NodesCloneOptions) {
	if opts.HasCustomTargetDir {
		return
	}
	if posDest := extractTargetDirFromPassArgs(opts.PassArgs, opts.HasFile); posDest != "" {
		opts.TargetDir = cleanTargetDir(posDest)
		opts.HasCustomTargetDir = true
	}
}

func resolveCloneDestination(opts *NodesCloneOptions) {
	applyPositionalTargetDir(opts)
	workBase := resolveLocalWorkBase()
	relSub, isInside := inspectCwdRelativeWorkDir(workBase)
	opts.IsInsideWorkDir = isInside
	opts.RelativeSubdir = relSub
	opts.Destination = NodesCloneDestination{
		TargetDir:          opts.TargetDir,
		RelativeSubdir:     opts.RelativeSubdir,
		HasCustomTargetDir: opts.HasCustomTargetDir,
		IsInsideWorkDir:    opts.IsInsideWorkDir,
	}
}

func finalizeCloneOptions(opts *NodesCloneOptions) {
	file, hasFile := DetectCloneFile(opts.PassArgs)
	opts.DetectedFile = file
	opts.HasFile = hasFile
	resolveCloneDestination(opts)
}

func initNodesCloneOptions(kind NodesCloneKind, raw []string) NodesCloneOptions {
	opts := NodesCloneOptions{Kind: kind, RawArgs: raw}
	lowKind := strings.ToLower(string(kind))
	if strings.Contains(lowKind, "except-self") || strings.Contains(lowKind, "noself") {
		opts.IsSkipLocal = true
	}
	return opts
}

func parseNodesCloneOptions(kind NodesCloneKind, raw []string) (NodesCloneOptions, bool) {
	opts := initNodesCloneOptions(kind, raw)
	if !parseTokens(&opts, raw) {
		PrintNodesCloneHelp(kind)
		return opts, false
	}
	finalizeCloneOptions(&opts)
	return opts, true
}

func validateCloneCommand(args []string) (NodesCloneKind, bool) {
	if len(args) == 0 {
		return "", false
	}
	return IsNodesCloneCommand(args[0])
}

// RunNodesClone orchestrates fleet clone across local host and remote nodes.
func RunNodesClone(args []string) error {
	kind, isClone := validateCloneCommand(args)
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

func renderPreFlightAndBanner(opts NodesCloneOptions, report FleetPreFlightReport) {
	if !opts.IsJSON {
		renderFleetPreFlightTable(os.Stdout, report)
		renderFleetStartBanner(os.Stdout, opts, report)
	}
}

func finishFleetDispatch(opts NodesCloneOptions, report FleetPreFlightReport, results []RemoteCloneNodeResult, isLocalOk bool, localDetails string, localDur time.Duration) error {
	if opts.IsJSON {
		return emitFleetJSON(report, results, isLocalOk)
	}
	renderFleetResultsTable(os.Stdout, results, isLocalOk, localDetails, localDur, opts)
	return nil
}

func dispatchFleetExecution(opts NodesCloneOptions) error {
	conns, _ := cmdssh.FetchAllSSHConnections()
	remoteConns := filterRemoteConnections(conns, opts)
	fileBytes, fileName := loadCloneFileBytes(opts)
	report := probeFleetNodesPreFlight(context.Background(), remoteConns, opts)
	renderPreFlightAndBanner(opts, report)
	results := executeFleetNodesParallel(remoteConns, opts, fileBytes, fileName)
	isLocalOk, localDetails, localDur := executeLocalClone(opts)
	recordFleetTaskAudit(opts, isLocalOk)
	return finishFleetDispatch(opts, report, results, isLocalOk, localDetails, localDur)
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

func readPipedOutput(r *os.File) string {
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func captureOutput(fn func() error) (string, error) {
	origStdout, origStderr := os.Stdout, os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		return "", fn()
	}
	os.Stdout, os.Stderr = w, w
	runErr := fn()
	_ = w.Close()
	os.Stdout, os.Stderr = origStdout, origStderr
	return readPipedOutput(r), runErr
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

func parseLocalResult(out string, err error, dur time.Duration) (bool, string, time.Duration) {
	if payload, hasJSON := cmdclone.ParseCloneJSONResponse(out); hasJSON {
		return payload.Success, payload.Message, dur
	}
	if err != nil {
		return false, err.Error(), dur
	}
	return true, "completed successfully", dur
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
	return parseLocalResult(out, err, time.Since(start))
}

func emitFleetJSON(report FleetPreFlightReport, results []RemoteCloneNodeResult, isLocalOk bool) error {
	payload := map[string]any{
		"localSuccess": isLocalOk,
		"preflight":    report,
		"nodes":        results,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
