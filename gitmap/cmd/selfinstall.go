package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/lockfile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/scripts"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// selfInstallOpts holds parsed flags for self-install.
type selfInstallOpts struct {
	Dir       string
	Yes       bool
	Version   string
	ShellMode string
	ShowPath  bool
	ForceLock bool
}

type selfInstallFlagVars struct {
	shellMode string
	profile   string
	dualShell bool
}

// runSelfInstall is the entry point for `gitmap self-install`.
func runSelfInstall(args []string) error {
	checkHelp(constants.CmdSelfInstall, args)
	opts := parseSelfInstallFlags(args)
	release := acquireSelfInstallLock(opts)
	defer release()

	runSelfInstallWorkflow(opts)
	return nil
}

func runSelfInstallWorkflow(opts selfInstallOpts) error {
	dir := resolveSelfInstallDir(opts)
	printSelfInstallStart(dir)
	executeSelfInstallScript(dir, opts)
	fmt.Print(constants.MsgSelfInstallDone)
	autoRunSetupAfterInstall()
	fmt.Print(constants.MsgSelfInstallReminder)
	return nil
}

func printSelfInstallStart(dir string) {
	fmt.Print(constants.MsgSelfInstallHeader)
	fmt.Printf(constants.MsgSelfInstallUsing, dir)
}

func executeSelfInstallScript(dir string, opts selfInstallOpts) {
	scriptName, scriptBody := loadInstallScript()
	tmpPath := writeInstallScriptTemp(scriptName, scriptBody)
	defer os.Remove(tmpPath)
	executeInstallScript(scriptName, tmpPath, dir, opts)
}

// autoRunSetupAfterInstall invokes `gitmap setup` as a final best-effort step.
func autoRunSetupAfterInstall() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, constants.MsgSelfInstallSetupSkipped, r)
		}
	}()
	fmt.Print(constants.MsgSelfInstallRunningSetup)
	runSetup(nil)
}

// acquireSelfInstallLock takes the duplicate-install guard.
func acquireSelfInstallLock(opts selfInstallOpts) lockfile.Releaser {
	if opts.ForceLock {
		return forceAcquireOrExit()
	}
	release, err := lockfile.Acquire(constants.SelfInstallLockName)
	if err == nil {
		return release
	}
	handleLockError(err)
	return func() {}
}

func handleLockError(err error) {
	if errors.Is(err, lockfile.ErrAlreadyHeld) {
		holder := lockfile.HolderPID(constants.SelfInstallLockName)
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallLockHeld, holder)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
	fmt.Fprintf(os.Stderr, constants.ErrSelfInstallLock, err)
	cliexit.HandleError(nil, constants.ExitCodeError)
}

func forceAcquireOrExit() lockfile.Releaser {
	release, err := lockfile.ForceAcquire(constants.SelfInstallLockName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallLock, err)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
	return release
}

// parseSelfInstallFlags reads flags and validates shell mode.
func parseSelfInstallFlags(args []string) selfInstallOpts {
	fs := flag.NewFlagSet(constants.CmdSelfInstall, flag.ExitOnError)
	opts := selfInstallOpts{}
	vars := selfInstallFlagVars{}
	bindSelfInstallFlags(fs, &opts, &vars)
	fs.Parse(reorderFlagsBeforeArgs(args))
	opts.ShellMode = resolveShellMode(vars.shellMode, vars.profile, vars.dualShell)
	validateShellMode(opts.ShellMode)
	return opts
}

func bindSelfInstallFlags(fs *flag.FlagSet, opts *selfInstallOpts, vars *selfInstallFlagVars) {
	bindSelfInstallDirAndConfirm(fs, opts)
	bindSelfInstallShellFlags(fs, vars)
	fs.StringVar(&opts.Version, constants.FlagNameVersion, "", constants.FlagDescSelfFromVersion)
	fs.BoolVar(&opts.ShowPath, constants.FlagNameShowPath, false, constants.FlagDescSelfShowPath)
	fs.BoolVar(&opts.ForceLock, constants.FlagNameForceLock, false, constants.FlagDescSelfForceLock)
}

func bindSelfInstallDirAndConfirm(fs *flag.FlagSet, opts *selfInstallOpts) {
	fs.StringVar(&opts.Dir, constants.FlagNameDir, "", constants.FlagDescSelfDir)
	fs.BoolVar(&opts.Yes, constants.FlagNameYes, false, constants.FlagDescSelfYes)
	fs.BoolVar(&opts.Yes, constants.FlagNameY, false, constants.FlagDescSelfYes)
}

func bindSelfInstallShellFlags(fs *flag.FlagSet, vars *selfInstallFlagVars) {
	fs.StringVar(&vars.shellMode, constants.FlagNameShellMode, "", constants.FlagDescSelfShellMode)
	fs.StringVar(&vars.profile, constants.FlagNameProfile, "", constants.FlagDescSelfProfile)
	fs.BoolVar(&vars.dualShell, constants.FlagNameDualShell, false, constants.FlagDescSelfDualShell)
}

func resolveShellMode(shellMode, profile string, dualShell bool) string {
	if len(shellMode) > 0 {
		return shellMode
	}
	if len(profile) > 0 {
		return profile
	}
	if dualShell {
		return constants.ShellModeBoth
	}
	return constants.ShellModeAuto
}

func validateShellMode(mode string) {
	if isValidSingletonShellMode(mode) || isValidComboShellMode(mode) {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrSelfInstallShellModeInvalid,
		mode, strings.Join(constants.SelfInstallShellModes, constants.ShellPipeSep))
	cliexit.HandleError(nil, constants.ExitCodeError)
}

func isValidSingletonShellMode(mode string) bool {
	for _, valid := range constants.SelfInstallShellModes {
		if mode == valid {
			return true
		}
	}
	return false
}

func isValidComboShellMode(mode string) bool {
	if !strings.Contains(mode, constants.ShellModeComboSep) {
		return false
	}
	tokens := strings.Split(mode, constants.ShellModeComboSep)
	if len(tokens) < 2 {
		return false
	}
	return validateComboTokens(tokens)
}

func validateComboTokens(tokens []string) bool {
	seen := map[string]bool{}
	for _, tok := range tokens {
		if !isConcreteShellFamily(tok) || seen[tok] {
			return false
		}
		seen[tok] = true
	}
	return true
}

func isConcreteShellFamily(tok string) bool {
	switch tok {
	case constants.ShellModeZsh, constants.ShellModeBash,
		constants.ShellModePwsh, constants.ShellModeFish:
		return true
	}
	return false
}

func resolveSelfInstallDir(opts selfInstallOpts) string {
	if len(opts.Dir) > 0 {
		return opts.Dir
	}
	def := defaultSelfInstallDir()
	if opts.Yes {
		return def
	}
	return promptInstallDir(def)
}

func defaultSelfInstallDir() string {
	if runtime.GOOS == constants.PlatformWindows {
		return constants.SelfInstallDefaultWindows
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return constants.SelfInstallDefaultUnixFallback
	}
	return filepath.Join(home, constants.SelfInstallDefaultUnix)
}

func promptInstallDir(def string) string {
	fmt.Printf(constants.MsgSelfInstallPrompt, def)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallReadStdin, err)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
	answer := strings.TrimSpace(line)
	if len(answer) == 0 {
		return def
	}
	return answer
}

func loadInstallScript() (string, []byte) {
	name := pickInstallScriptName()
	body, err := scripts.ReadFile(name)
	if err == nil && len(body) > 0 {
		fmt.Printf(constants.MsgSelfInstallEmbedded, name)
		return name, body
	}
	return name, downloadFallbackInstallScript()
}

func downloadFallbackInstallScript() []byte {
	remote := pickInstallScriptURL()
	fmt.Printf(constants.MsgSelfInstallRemote, remote)
	body, dlErr := downloadInstallScript(remote)
	if dlErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallDownload, remote, dlErr)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
	return body
}

func pickInstallScriptName() string {
	if runtime.GOOS == constants.PlatformWindows {
		return constants.SelfInstallScriptPwsh
	}
	return constants.SelfInstallScriptBash
}

func pickInstallScriptURL() string {
	if runtime.GOOS == constants.PlatformWindows {
		return constants.SelfInstallRemotePwsh
	}
	return constants.SelfInstallRemoteBash
}

func downloadInstallScript(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec // G107: URL is a compile-time constant.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(constants.ErrHTTPStatusFmt, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func writeInstallScriptTemp(name string, body []byte) string {
	pattern := tempScriptPattern(name)
	f, err := os.CreateTemp(os.TempDir(), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallScriptWrite, err)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
	defer f.Close()
	writeScriptBody(f, name, body)
	setScriptExecutable(f.Name(), name)
	return f.Name()
}

func tempScriptPattern(name string) string {
	if strings.HasSuffix(name, constants.ScriptExtPs1) {
		return constants.SelfInstallTempPrefix + constants.ScriptExtPs1
	}
	return constants.SelfInstallTempPrefix + constants.ScriptExtSh
}

func writeScriptBody(f *os.File, name string, body []byte) {
	if strings.HasSuffix(name, constants.ScriptExtPs1) {
		writeBOM(f)
	}
	if _, err := f.Write(body); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallScriptWrite, err)
		_ = f.Close()
		exitWith(constants.ExitCodeError)
	}
}

func setScriptExecutable(filePath, name string) {
	if !strings.HasSuffix(name, constants.ScriptExtPs1) {
		_ = os.Chmod(filePath, constants.DirPermission)
	}
}

func writeBOM(f *os.File) {
	_, err := f.Write(constants.UTF8BOM)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallScriptWrite, err)
		_ = f.Close()
		exitWith(constants.ExitCodeError)
	}
}

func executeInstallScript(name, path, dir string, opts selfInstallOpts) {
	cmd := buildSelfInstallCmd(name, path, dir, opts)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfInstallScriptRun, err)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
}

func buildSelfInstallCmd(name, path, dir string, opts selfInstallOpts) *exec.Cmd {
	if strings.HasSuffix(name, constants.ScriptExtPs1) {
		return buildSelfInstallPwshCmd(path, dir, opts)
	}
	return buildSelfInstallBashCmd(path, dir, opts)
}

func buildSelfInstallPwshCmd(path, dir string, opts selfInstallOpts) *exec.Cmd {
	args := []string{
		constants.PwshArgExecPolicy, constants.PwshArgBypass,
		constants.PwshArgNoProfile, constants.PwshArgNoLogo,
		constants.PwshArgFile, path,
		constants.PwshArgInstallDir, dir,
	}
	if len(opts.Version) > 0 {
		args = append(args, constants.PwshArgVersion, opts.Version)
	}
	return exec.Command(constants.ShellModePwsh, args...)
}

func buildSelfInstallBashCmd(path, dir string, opts selfInstallOpts) *exec.Cmd {
	args := buildSelfInstallBashArgs(path, dir, opts)
	cmd := exec.Command(constants.ShellModeBash, args...)
	if shellModeRequiresPwsh(opts.ShellMode) {
		cmd.Env = append(os.Environ(), constants.EnvGitmapDualShell)
	}
	return cmd
}

func buildSelfInstallBashArgs(path, dir string, opts selfInstallOpts) []string {
	args := []string{path, constants.FlagSelfDir, dir}
	if len(opts.Version) > 0 {
		args = append(args, constants.FlagSelfFromVersion, opts.Version)
	}
	args = append(args, constants.FlagSelfShellMode, opts.ShellMode)
	if opts.ShowPath {
		args = append(args, constants.FlagSelfShowPath)
	}
	return args
}

func shellModeRequiresPwsh(mode string) bool {
	if mode == constants.ShellModeBoth || mode == constants.ShellModePwsh {
		return true
	}
	if !strings.Contains(mode, constants.ShellModeComboSep) {
		return false
	}
	for _, tok := range strings.Split(mode, constants.ShellModeComboSep) {
		if tok == constants.ShellModePwsh {
			return true
		}
	}
	return false
}
