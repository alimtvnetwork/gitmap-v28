package cmdinstall

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// MusePlatformType represents the target operating system platform.
type MusePlatformType string

const (
	MusePlatformWindows MusePlatformType = "windows"
	MusePlatformLinux   MusePlatformType = "linux"
	MusePlatformDarwin  MusePlatformType = "darwin"
)

// InstallMuseOptions configures the installation process for Meta Muse.
type InstallMuseOptions struct {
	Platform   MusePlatformType `json:"platform,omitempty"`
	Version    string           `json:"version,omitempty"`
	Force      bool             `json:"force,omitempty"`
	DryRun     bool             `json:"dryRun,omitempty"`
	VerifyOnly bool             `json:"verifyOnly,omitempty"`
	Verify     bool             `json:"verify,omitempty"`
	IsJSON     bool             `json:"isJson,omitempty"`
	InstallDir string           `json:"installDir,omitempty"`
	Timeout    time.Duration    `json:"timeout,omitempty"`
}

// MuseInstallOptions is an alias for InstallMuseOptions for backward compatibility.
type MuseInstallOptions = InstallMuseOptions

// MuseInstallResult captures execution telemetry and post-install status.
type MuseInstallResult struct {
	IsSuccess       bool             `json:"isSuccess"`
	Success         bool             `json:"success"`
	Platform        MusePlatformType `json:"platform"`
	Version         string           `json:"version,omitempty"`
	VersionDetected string           `json:"versionDetected,omitempty"`
	BinaryPath      string           `json:"binaryPath,omitempty"`
	InstalledPath   string           `json:"installedPath,omitempty"`
	CommandExecuted string           `json:"commandExecuted"`
	ExitCode        int              `json:"exitCode"`
	DurationMs      int64            `json:"durationMs"`
	ErrorMessage    string           `json:"errorMessage,omitempty"`
}

// ResolveMusePlatform normalizes platform string or defaults to host runtime.GOOS.
func ResolveMusePlatform(override string) MusePlatformType {
	clean := strings.ToLower(strings.TrimSpace(override))
	switch clean {
	case "windows", "win", "win32", "win64":
		return MusePlatformWindows
	case "linux", "ubuntu", "debian":
		return MusePlatformLinux
	case "darwin", "mac", "macos", "osx":
		return MusePlatformDarwin
	}

	switch runtime.GOOS {
	case "windows":
		return MusePlatformWindows
	case "darwin":
		return MusePlatformDarwin
	default:
		return MusePlatformLinux
	}
}

// ResolveMuseInstallCommand returns the command runner and arguments for the platform.
func ResolveMuseInstallCommand(platform MusePlatformType) (string, []string, error) {
	return ResolveMuseInstallCommandWithVersion(platform, "")
}

// ResolveMuseInstallCommandWithVersion returns command runner and arguments with version support.
func ResolveMuseInstallCommandWithVersion(platform MusePlatformType, version string) (string, []string, error) {
	switch platform {
	case MusePlatformWindows:
		return resolveWindowsInstallCommand(version)
	case MusePlatformLinux:
		return resolveLinuxInstallCommand(version)
	case MusePlatformDarwin:
		return resolveDarwinInstallCommand(version)
	default:
		return "", nil, apperror.NewValidationError("unsupported platform for Muse installer: " + string(platform))
	}
}

func resolveWindowsInstallCommand(version string) (string, []string, error) {
	cmdScript := "irm " + constants.MetaMuseInstallUrlWindows + " | iex"
	if version != "" {
		cmdScript = fmt.Sprintf("$env:MUSE_VERSION='%s'; irm %s | iex", version, constants.MetaMuseInstallUrlWindows)
	}

	args := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", cmdScript}

	return "powershell", args, nil
}

func resolveLinuxInstallCommand(version string) (string, []string, error) {
	cmdScript := "curl -fsSL " + constants.MetaMuseInstallUrlUnix + " | bash"
	if version != "" {
		cmdScript = fmt.Sprintf("curl -fsSL %s | bash -s -- --version %s", constants.MetaMuseInstallUrlUnix, version)
	}

	args := []string{"-c", cmdScript}

	return "bash", args, nil
}

func resolveDarwinInstallCommand(version string) (string, []string, error) {
	cmdScript := "curl -fsSL " + constants.MetaMuseInstallUrlUnix + " | bash"
	if version != "" {
		cmdScript = fmt.Sprintf("curl -fsSL %s | bash -s -- --version %s", constants.MetaMuseInstallUrlUnix, version)
	}

	args := []string{"-c", cmdScript}

	return "bash", args, nil
}

// ProbeMuseEndpoint tests connectivity to the official Meta Muse distribution server.
func ProbeMuseEndpoint(timeout time.Duration) error {
	reqTimeout := timeout
	if reqTimeout <= 0 {
		reqTimeout = 5 * time.Second
	}

	client := &http.Client{Timeout: reqTimeout}
	resp, err := client.Head("https://dev.meta.ai")
	if err != nil {
		return apperror.WrapSimple(err, "muse: connectivity check failed")
	}
	defer resp.Body.Close()

	return nil
}

// CheckExistingMuse inspects system PATH and local directories for an installed muse binary.
func CheckExistingMuse() (string, string, bool) {
	return CheckExistingMuseWithDir("")
}

// CheckExistingMuseWithDir inspects custom directory, PATH, and platform locations for muse.
func CheckExistingMuseWithDir(customDir string) (string, string, bool) {
	binName := constants.MetaMuseBinaryNameUnix
	if runtime.GOOS == "windows" {
		binName = constants.MetaMuseBinaryNameWindows
	}

	foundPath, isFound := locateMuseBinary(binName, customDir)
	if !isFound {
		return "", "", false
	}

	version := detectMuseVersion(foundPath)

	return foundPath, version, true
}

func locateMuseBinary(binName, customDir string) (string, bool) {
	if customDir != "" {
		p := filepath.Join(customDir, binName)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, true
		}
	}

	if p, err := exec.LookPath(binName); err == nil {
		return p, true
	}

	return searchStandardMusePaths(binName)
}

func searchStandardMusePaths(binName string) (string, bool) {
	home, err := os.UserHomeDir()
	if err == nil {
		localBin := filepath.Join(home, ".local", "bin", binName)
		if info, err := os.Stat(localBin); err == nil && !info.IsDir() {
			return localBin, true
		}
	}

	if runtime.GOOS != "windows" {
		usrBin := filepath.Join("/usr/local/bin", binName)
		if info, err := os.Stat(usrBin); err == nil && !info.IsDir() {
			return usrBin, true
		}
	}

	return searchWindowsMusePaths(binName, home)
}

func searchWindowsMusePaths(binName, home string) (string, bool) {
	if runtime.GOOS != "windows" {
		return "", false
	}

	if localApp := os.Getenv("LOCALAPPDATA"); localApp != "" {
		p := filepath.Join(localApp, "Programs", "Muse", binName)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, true
		}
	}

	if home != "" {
		p := filepath.Join(home, ".muse", "bin", binName)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, true
		}
	}

	return "", false
}

func detectMuseVersion(binPath string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return "unknown"
}

// InstallMuse coordinates the installation of the Meta Muse AI developer agent.
func InstallMuse(opts InstallMuseOptions) error {
	res, err := RunMuseInstaller(opts)
	if err != nil {
		return err
	}

	if !res.Success {
		return apperror.NewSimple("muse: installation failed", constants.ErrMuseDownloadFailed)
	}

	return nil
}

// RunMuseInstaller executes the unified cross-platform Meta Muse installer pipeline.
func RunMuseInstaller(opts InstallMuseOptions) (*MuseInstallResult, error) {
	platform := opts.Platform
	if platform == "" {
		platform = ResolveMusePlatform("")
	}

	opts.Version = resolveEffectiveVersion(opts.Version)
	opts.IsJSON = resolveEffectiveJSON(opts.IsJSON)

	if opts.VerifyOnly {
		return handleVerifyOnly(opts, platform)
	}

	existingPath, existingVer, isInstalled := CheckExistingMuseWithDir(opts.InstallDir)
	if isInstalled && !opts.Force && !opts.DryRun {
		return buildExistingMuseResult(opts, platform, existingPath, existingVer)
	}

	bin, args, err := ResolveMuseInstallCommandWithVersion(platform, opts.Version)
	if err != nil {
		return nil, err
	}

	cmdStr := bin + " " + strings.Join(args, " ")
	result := &MuseInstallResult{
		Platform:        platform,
		CommandExecuted: cmdStr,
	}

	if opts.DryRun {
		return buildDryRunMuseResult(opts, result)
	}

	return executeMuseInstallPipeline(opts, bin, args, result)
}

func handleVerifyOnly(opts InstallMuseOptions, platform MusePlatformType) (*MuseInstallResult, error) {
	existingPath, existingVer, isInstalled := CheckExistingMuseWithDir(opts.InstallDir)
	res := &MuseInstallResult{
		IsSuccess:       isInstalled,
		Success:         isInstalled,
		Platform:        platform,
		InstalledPath:   existingPath,
		BinaryPath:      existingPath,
		VersionDetected: existingVer,
		Version:         existingVer,
	}

	if !isInstalled {
		res.ErrorMessage = "meta muse is not installed"
		if opts.IsJSON {
			_ = emitJSONReport(res)
		} else {
			fmt.Printf("%s○ Meta Muse is not installed.%s\n", constants.ColorYellow, constants.ColorReset)
		}

		return res, apperror.NewSimple("meta muse is not installed", constants.ErrMuseProbeFailed)
	}

	if opts.IsJSON {
		_ = emitJSONReport(res)
	} else {
		fmt.Printf("%s✔ Meta Muse is installed at: %s (version: %s)%s\n", constants.ColorGreen, existingPath, existingVer, constants.ColorReset)
	}

	return res, nil
}

func buildExistingMuseResult(opts InstallMuseOptions, platform MusePlatformType, path, ver string) (*MuseInstallResult, error) {
	res := &MuseInstallResult{
		IsSuccess:       true,
		Success:         true,
		Platform:        platform,
		InstalledPath:   path,
		BinaryPath:      path,
		VersionDetected: ver,
		Version:         ver,
	}

	if opts.IsJSON {
		_ = emitJSONReport(res)

		return res, nil
	}

	fmt.Printf("%s✓ Meta Muse is already installed at: %s (version: %s)%s\n", constants.ColorGreen, path, ver, constants.ColorReset)
	fmt.Printf("  Pass --force to reinstall.\n")

	return res, nil
}

func buildDryRunMuseResult(opts InstallMuseOptions, result *MuseInstallResult) (*MuseInstallResult, error) {
	result.Success = true
	result.IsSuccess = true

	if opts.IsJSON {
		_ = emitJSONReport(result)

		return result, nil
	}

	fmt.Printf("%s[DRY RUN] Meta Muse Install Command for %s:%s\n", constants.ColorYellow, result.Platform, constants.ColorReset)
	fmt.Printf("  %s\n", result.CommandExecuted)

	return result, nil
}

func executeMuseInstallPipeline(opts InstallMuseOptions, bin string, args []string, result *MuseInstallResult) (*MuseInstallResult, error) {
	if err := ProbeMuseEndpoint(opts.Timeout); err != nil && !opts.IsJSON {
		fmt.Printf("%s▲ Warning: connectivity probe to https://dev.meta.ai reported: %v (continuing)%s\n", constants.ColorYellow, err, constants.ColorReset)
	}

	execBin := resolvePlatformExecutor(result.Platform, bin)
	startTime := time.Now()

	if !opts.IsJSON {
		fmt.Printf("%s● Installing Meta Muse AI Companion (%s)%s\n\n", constants.ColorCyan, result.Platform, constants.ColorReset)
	}

	cmd := exec.Command(execBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	runErr := cmd.Run()
	result.DurationMs = time.Since(startTime).Milliseconds()
	result.ExitCode = resolveCmdExitCode(runErr)

	if runErr != nil {
		return handleInstallFailure(opts, runErr, result)
	}

	return handleInstallSuccess(opts, result)
}

func resolvePlatformExecutor(platform MusePlatformType, defaultBin string) string {
	if platform == MusePlatformWindows {
		if pwsh := resolvePowerShellBinary(); pwsh != "" {
			return pwsh
		}
	}

	return defaultBin
}

func handleInstallFailure(opts InstallMuseOptions, runErr error, result *MuseInstallResult) (*MuseInstallResult, error) {
	if result.Platform == MusePlatformDarwin {
		if brewErr := executeDarwinFallback(result); brewErr == nil {
			return handleInstallSuccess(opts, result)
		}
	}

	result.ErrorMessage = runErr.Error()
	result.Success = false
	result.IsSuccess = false

	if opts.IsJSON {
		_ = emitJSONReport(result)
	}

	return result, apperror.WrapSimple(runErr, "muse: installer execution failed")
}

func executeDarwinFallback(result *MuseInstallResult) error {
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return apperror.NewSimple("brew not found for Darwin fallback", constants.ErrMissingDownloadTool)
	}

	cmd := exec.Command(brewPath, "install", "meta/muse/muse")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	runErr := cmd.Run()
	if runErr != nil {
		return apperror.WrapSimple(runErr, "muse: brew fallback installation failed")
	}

	result.CommandExecuted = brewPath + " install meta/muse/muse"
	result.Success = true
	result.IsSuccess = true

	return nil
}

func handleInstallSuccess(opts InstallMuseOptions, result *MuseInstallResult) (*MuseInstallResult, error) {
	result.Success = true
	result.IsSuccess = true
	applyPostInstallVerification(opts.Verify || !opts.IsJSON, opts.InstallDir, result)

	recordVer := result.VersionDetected
	if opts.Version != "" {
		recordVer = opts.Version
	}

	recordMuseInstallation(recordVer, result.Platform, opts.IsJSON)

	if !opts.IsJSON {
		fmt.Printf("%s✓ Meta Muse installed successfully%s\n", constants.ColorGreen, constants.ColorReset)
	}

	if opts.IsJSON {
		_ = emitJSONReport(result)
	}

	return result, nil
}

func recordMuseInstallation(version string, platform MusePlatformType, isJSON bool) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}
	defer splitDB.Close()

	source := string(platform)
	if source == "" {
		source = runtime.GOOS
	}

	effectiveVer := version
	if effectiveVer == "" || effectiveVer == "unknown" {
		effectiveVer = "latest"
	}

	_ = splitDB.SaveInstalledTool(constants.ToolMuse, effectiveVer, source)
	_ = splitDB.RecordLog(store.InstallationLogRecord{
		Tool:           constants.ToolMuse,
		Action:         "install",
		Version:        effectiveVer,
		PackageManager: source,
		IsSuccess:      true,
	})

	if !isJSON && effectiveVer != "" {
		fmt.Printf(constants.MsgInstallRecorded, constants.ToolMuse, effectiveVer)
	}
}

func emitJSONReport(result *MuseInstallResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "muse: json serialize failed")
	}

	fmt.Println(string(data))

	return nil
}

func resolveCmdExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	return 1
}

func applyPostInstallVerification(shouldVerify bool, customDir string, result *MuseInstallResult) {
	newPath, newVer, isFound := CheckExistingMuseWithDir(customDir)
	if isFound {
		result.InstalledPath = newPath
		result.BinaryPath = newPath
		result.VersionDetected = newVer
		result.Version = newVer
	}

	if shouldVerify && isFound {
		fmt.Printf("%s✔ Verified Meta Muse binary: %s (%s)%s\n", constants.ColorGreen, result.InstalledPath, result.VersionDetected, constants.ColorReset)
	}
}

func resolveEffectiveVersion(optsVersion string) string {
	if strings.TrimSpace(optsVersion) != "" {
		return strings.TrimSpace(optsVersion)
	}

	for i, arg := range os.Args {
		if (arg == "--version" || arg == "-v") && i+1 < len(os.Args) {
			val := strings.TrimSpace(os.Args[i+1])
			if !strings.HasPrefix(val, "-") {
				return val
			}
		}
		if strings.HasPrefix(arg, "--version=") {
			return strings.TrimSpace(strings.TrimPrefix(arg, "--version="))
		}
	}

	return ""
}

func resolveEffectiveJSON(optsJSON bool) bool {
	if optsJSON {
		return true
	}

	for _, arg := range os.Args {
		if arg == "--json" || arg == "-j" {
			return true
		}
	}

	return false
}
