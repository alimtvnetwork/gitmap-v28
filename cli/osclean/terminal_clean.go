package osclean

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// CleanTerminalHistory clears terminal command history and re-seeds suggestions.
func CleanTerminalHistory(opts TerminalCleanOptions) result.Result[TerminalCleanSummary] {
	start := time.Now()
	summary := TerminalCleanSummary{
		IsDryRun: opts.IsDryRun,
	}

	cleaners := getShellCleaners()
	isGitmapPresent := isGitmapInstalledBefore()

	for _, c := range cleaners {
		if isShellSelected(c.Shell, opts.OnlyShells) {
			stats := c.CleanFunc(opts, isGitmapPresent)
			summary.TotalFilesCleared += stats.FilesCleared
			summary.TotalBytesFreed += stats.BytesFreed
			summary.TotalReseedCount += stats.ReseedCount
			summary.Shells = append(summary.Shells, stats)
		}
	}

	summary.DurationMs = time.Since(start).Milliseconds()
	return result.Ok(summary)
}

type shellCleanerDef struct {
	Shell     string
	Aliases   []string
	CleanFunc func(TerminalCleanOptions, bool) ShellCleanStats
}

func getShellCleaners() []shellCleanerDef {
	return []shellCleanerDef{
		{
			Shell:     "powershell",
			Aliases:   []string{"pwsh", "ps", "posh"},
			CleanFunc: cleanPowerShellTerminal,
		},
		{
			Shell:     "bash",
			Aliases:   []string{"git-bash", "sh-bash"},
			CleanFunc: cleanBashTerminal,
		},
		{
			Shell:     "zsh",
			Aliases:   []string{"zshell", "oh-my-zsh"},
			CleanFunc: cleanZshTerminal,
		},
		{
			Shell:     "fish",
			Aliases:   []string{"fish-shell"},
			CleanFunc: cleanFishTerminal,
		},
		{
			Shell:     "sh",
			Aliases:   []string{"shell", "dash"},
			CleanFunc: cleanShTerminal,
		},
	}
}

func isShellSelected(shell string, only []string) bool {
	if len(only) == 0 {
		return true
	}
	target := strings.ToLower(strings.TrimSpace(shell))
	for _, o := range only {
		norm := strings.ToLower(strings.TrimSpace(o))
		if norm == target || isShellAliasMatch(shell, norm) {
			return true
		}
	}
	return false
}

func isShellAliasMatch(shell, norm string) bool {
	for _, c := range getShellCleaners() {
		if c.Shell == shell && containsAlias(c.Aliases, norm) {
			return true
		}
	}
	return false
}

func cleanPowerShellTerminal(opts TerminalCleanOptions, isGitmapPresent bool) ShellCleanStats {
	stats := ShellCleanStats{
		Shell: "powershell",
		Label: "PowerShell & PSReadLine History",
	}

	targets := resolvePowerShellHistoryTargets()
	canonical := getCanonicalGitmapSuggestions()

	for _, p := range targets {
		processTerminalHistoryFile(&stats, p, canonical, opts.IsDryRun, opts.ShouldReseed && isGitmapPresent, false)
	}

	if opts.ShouldReseed && isGitmapPresent && !opts.IsDryRun {
		_ = completion.Install(constants.ShellPowerShell)
		stats.Notes = append(stats.Notes, "PowerShell argument completer reinstalled")
	}

	return stats
}

func cleanBashTerminal(opts TerminalCleanOptions, isGitmapPresent bool) ShellCleanStats {
	stats := ShellCleanStats{
		Shell: "bash",
		Label: "Bash History & Sessions",
	}

	targets := resolveBashHistoryTargets()
	canonical := getCanonicalGitmapSuggestions()

	for _, p := range targets {
		processTerminalHistoryFile(&stats, p, canonical, opts.IsDryRun, opts.ShouldReseed && isGitmapPresent, false)
	}

	cleanDirectorySessions(&stats, resolveBashSessionsDir(), opts.IsDryRun)

	if opts.ShouldReseed && isGitmapPresent && !opts.IsDryRun {
		_ = completion.Install(constants.ShellBash)
		stats.Notes = append(stats.Notes, "Bash tab completion script verified")
	}

	return stats
}

func cleanZshTerminal(opts TerminalCleanOptions, isGitmapPresent bool) ShellCleanStats {
	stats := ShellCleanStats{
		Shell: "zsh",
		Label: "Zsh History & Compdump",
	}

	targets := resolveZshHistoryTargets()
	canonical := getCanonicalGitmapSuggestions()

	for _, p := range targets {
		processTerminalHistoryFile(&stats, p, canonical, opts.IsDryRun, opts.ShouldReseed && isGitmapPresent, true)
	}

	cleanDirectorySessions(&stats, resolveZshSessionsDir(), opts.IsDryRun)
	cleanZshCompdumps(&stats, opts.IsDryRun)

	if opts.ShouldReseed && isGitmapPresent && !opts.IsDryRun {
		_ = completion.Install(constants.ShellZsh)
		stats.Notes = append(stats.Notes, "Zsh completion module verified")
	}

	return stats
}

func cleanFishTerminal(opts TerminalCleanOptions, isGitmapPresent bool) ShellCleanStats {
	stats := ShellCleanStats{
		Shell: "fish",
		Label: "Fish Shell History",
	}

	targets := resolveFishHistoryTargets()
	for _, p := range targets {
		processTerminalHistoryFile(&stats, p, nil, opts.IsDryRun, false, false)
	}

	if opts.ShouldReseed && isGitmapPresent && !opts.IsDryRun {
		_ = completion.Install(constants.ShellFish)
	}

	return stats
}

func cleanShTerminal(opts TerminalCleanOptions, isGitmapPresent bool) ShellCleanStats {
	stats := ShellCleanStats{
		Shell: "sh",
		Label: "POSIX Shell History",
	}

	targets := resolveShHistoryTargets()
	for _, p := range targets {
		processTerminalHistoryFile(&stats, p, nil, opts.IsDryRun, false, false)
	}

	return stats
}

func processTerminalHistoryFile(stats *ShellCleanStats, filePath string, suggestions []string, isDryRun, shouldReseed, isZsh bool) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return
	}

	bytes := fi.Size()
	stats.FilesCleared++
	stats.BytesFreed += bytes
	stats.ClearedPaths = append(stats.ClearedPaths, filePath)

	if isDryRun && shouldReseed {
		stats.IsReseeded = true
		stats.ReseedCount += len(suggestions)
		return
	}
	if isDryRun {
		return
	}

	if shouldReseed {
		writeReseededHistory(stats, filePath, suggestions, isZsh)
		return
	}

	truncateFile(stats, filePath)
}

func writeReseededHistory(stats *ShellCleanStats, filePath string, suggestions []string, isZsh bool) {
	var sb strings.Builder
	for _, cmd := range suggestions {
		if isZsh {
			sb.WriteString(fmt.Sprintf(": %d:0;%s\n", time.Now().Unix(), cmd))
		} else {
			sb.WriteString(cmd + "\n")
		}
	}

	if err := os.WriteFile(filePath, []byte(sb.String()), 0600); err != nil {
		stats.Errors = append(stats.Errors, fmt.Sprintf("failed writing %s: %v", filePath, err))
		return
	}

	stats.IsReseeded = true
	stats.ReseedCount += len(suggestions)
}

func truncateFile(stats *ShellCleanStats, filePath string) {
	if err := os.Truncate(filePath, 0); err != nil {
		_ = os.WriteFile(filePath, []byte{}, 0600)
	}
}

func cleanDirectorySessions(stats *ShellCleanStats, dirPath string, isDryRun bool) {
	if dirPath == "" {
		return
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, e := range entries {
		p := filepath.Join(dirPath, e.Name())
		fi, statErr := os.Stat(p)
		if statErr != nil {
			continue
		}
		stats.FilesCleared++
		stats.BytesFreed += fi.Size()
		stats.ClearedPaths = append(stats.ClearedPaths, p)
		if !isDryRun {
			_ = os.Remove(p)
		}
	}
}

func cleanZshCompdumps(stats *ShellCleanStats, isDryRun bool) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}

	matches, globErr := filepath.Glob(filepath.Join(home, ".zcompdump*"))
	if globErr != nil {
		return
	}

	for _, m := range matches {
		fi, statErr := os.Stat(m)
		if statErr != nil {
			continue
		}
		stats.FilesCleared++
		stats.BytesFreed += fi.Size()
		stats.ClearedPaths = append(stats.ClearedPaths, m)
		if !isDryRun {
			_ = os.Remove(m)
		}
	}
}

func resolvePowerShellHistoryTargets() []string {
	var targets []string
	seen := make(map[string]bool)

	addPath := func(p string) {
		clean := filepath.Clean(strings.TrimSpace(p))
		if clean != "" && !seen[clean] {
			seen[clean] = true
			targets = append(targets, clean)
		}
	}

	home, _ := os.UserHomeDir()
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")

	if appData != "" {
		addPath(filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt"))
		addPath(filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "Visual Studio Code Host_history.txt"))
	}
	if localAppData != "" {
		addPath(filepath.Join(localAppData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt"))
	}
	if home != "" {
		addPath(filepath.Join(home, "AppData", "Roaming", "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt"))
		addPath(filepath.Join(home, ".local", "share", "powershell", "PSReadLine", "ConsoleHost_history.txt"))
		addPath(filepath.Join(home, ".config", "powershell", "PSReadLine", "ConsoleHost_history.txt"))
	}

	return filterExistingPaths(targets)
}

func resolveBashHistoryTargets() []string {
	var targets []string
	home, _ := os.UserHomeDir()
	userProfile := os.Getenv("USERPROFILE")

	if home != "" {
		targets = append(targets, filepath.Join(home, ".bash_history"))
	}
	if userProfile != "" && userProfile != home {
		targets = append(targets, filepath.Join(userProfile, ".bash_history"))
	}

	return filterExistingPaths(targets)
}

func resolveBashSessionsDir() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}
	dir := filepath.Join(home, ".bash_sessions")
	if isDir(dir) {
		return dir
	}
	return ""
}

func resolveZshHistoryTargets() []string {
	var targets []string
	home, _ := os.UserHomeDir()
	userProfile := os.Getenv("USERPROFILE")

	if home != "" {
		targets = append(targets, filepath.Join(home, ".zsh_history"), filepath.Join(home, ".zhistory"))
	}
	if userProfile != "" && userProfile != home {
		targets = append(targets, filepath.Join(userProfile, ".zsh_history"))
	}

	return filterExistingPaths(targets)
}

func resolveZshSessionsDir() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}
	dir := filepath.Join(home, ".zsh_sessions")
	if isDir(dir) {
		return dir
	}
	return ""
}

func resolveFishHistoryTargets() []string {
	var targets []string
	home, _ := os.UserHomeDir()
	if home != "" {
		targets = append(targets, filepath.Join(home, ".local", "share", "fish", "fish_history"))
	}
	return filterExistingPaths(targets)
}

func resolveShHistoryTargets() []string {
	var targets []string
	home, _ := os.UserHomeDir()
	if home != "" {
		targets = append(targets, filepath.Join(home, ".sh_history"), filepath.Join(home, ".history"))
	}
	return filterExistingPaths(targets)
}

func filterExistingPaths(paths []string) []string {
	var existing []string
	seen := make(map[string]bool)
	for _, p := range paths {
		clean := filepath.Clean(p)
		if seen[clean] {
			continue
		}
		seen[clean] = true
		if _, err := os.Stat(clean); err == nil {
			existing = append(existing, clean)
		}
	}
	return existing
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

func isGitmapInstalledBefore() bool {
	home, _ := os.UserHomeDir()
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")

	checkFiles := []string{
		filepath.Join(appData, constants.CompDirName, constants.CompFilePS),
		filepath.Join(home, ".config", constants.CompDirName, constants.CompFilePS),
		filepath.Join(home, ".config", constants.CompDirName, constants.CompFileBash),
		filepath.Join(home, ".config", constants.CompDirName, constants.CompFileZsh),
		filepath.Join(home, ".gitmap-completion.bash"),
		filepath.Join(localAppData, "gitmap-cli", "gitmap.exe"),
		filepath.Join(home, ".local", "share", "gitmap-cli", "gitmap"),
	}

	for _, f := range checkFiles {
		if f != "" && fileExists(f) {
			return true
		}
	}

	if _, err := exec.LookPath("gitmap"); err == nil {
		return true
	}

	exe, err := os.Executable()
	if err == nil && strings.Contains(strings.ToLower(exe), "gitmap") {
		return true
	}

	return false
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func getCanonicalGitmapSuggestions() []string {
	return []string{
		"gitmap",
		"gitmap --help",
		"gitmap status",
		"gitmap scan",
		"gitmap cd",
		"gitmap pull",
		"gitmap deploy",
		"gitmap deploy-keys-all",
		"gitmap clear devtools",
		"gitmap clear terminal",
		"gitmap os",
		"gitmap machine",
		"gitmap ip",
		"gitmap ssh",
	}
}
