package cmdinstall

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	// agmUpdateEnvMockScript points the unix installer at a local script
	// instead of the curl|bash one-liner. Unset = default behavior, unchanged.
	// Test seam for the e2e worker: GITMAP_AGM_INSTALL_SCRIPT=/path/to/mock.sh
	agmUpdateEnvMockScript = "GITMAP_AGM_INSTALL_SCRIPT"

	agmUpdateMaxLogLines  = 400
	agmUpdateMaxLineRunes = 500
	agmUpdateTailLines    = 15
	agmUpdateSpinnerRunes = 90
)

// agmAnsiEscape strips ANSI escapes from captured installer output before
// re-rendering; re-printing raw escapes uncoordinated is what corrupted the screen.
var agmAnsiEscape = regexp.MustCompile("\x1b\\[[0-9;?]*[a-zA-Z]|\x1b[()][0-9A-Z]")

// agmUpdateRenderer is the single coordinated writer for AGM install/update runs.
// Child stdio is captured through pipes and re-rendered here.
type agmUpdateRenderer struct {
	mu       sync.Mutex
	action   string // "Updating" | "Installing"
	via      string // "curl | bash" | "PowerShell"
	version  string
	plain    bool // non-TTY fallback: plain lines, zero escape sequences
	multi    pterm.MultiPrinter
	spinner  *pterm.SpinnerPrinter
	stage    string
	stages   []string
	log      []string
	finished bool
}

func newAgmUpdateRenderer(action, via, version string) *agmUpdateRenderer {
	r := &agmUpdateRenderer{action: action, via: via, version: version}
	r.plain = !isAgmUpdateTTY()
	r.printHeader()
	if r.plain {
		return r
	}
	r.multi = pterm.DefaultMultiPrinter
	r.multi.Start()
	spinner, err := pterm.DefaultSpinner.WithWriter(r.multi.NewWriter()).Start("Preparing…")
	if err != nil {
		r.multi.Stop()
		r.plain = true
		return r
	}
	r.spinner = spinner
	return r
}

func (r *agmUpdateRenderer) title() string {
	if r.version == "" {
		return r.action + " Antigravity Manager"
	}
	return r.action + " Antigravity Manager (" + r.version + ")"
}

func (r *agmUpdateRenderer) printHeader() {
	msg := r.title() + " via " + r.via + "..."
	if r.plain {
		fmt.Println(msg)
		return
	}
	pterm.Info.Println(msg)
}

func isAgmUpdateTTY() bool {
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice != 0
}

// FeedLine ingests one captured installer-output line. Safe for concurrent use
// (stdout/stderr are scanned on separate goroutines).
func (r *agmUpdateRenderer) FeedLine(raw string) {
	line := strings.TrimSpace(agmAnsiEscape.ReplaceAllString(raw, ""))
	if line == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.finished {
		return
	}
	r.appendLog(line)
	if name, isStage := agmStageFromLine(line); isStage {
		r.onStageLocked(name)
		return
	}
	r.onProgressLocked(line)
}

// agmStageFromLine recognizes installer stage headers. Grounded in the real
// install.sh: step() prints "==> <name>".
func agmStageFromLine(line string) (string, bool) {
	if !strings.HasPrefix(line, "==>") {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(line, "==>"))
	if name == "" {
		return "", false
	}
	return name, true
}

func (r *agmUpdateRenderer) onStageLocked(name string) {
	r.stage = name
	r.stages = append(r.stages, name)
	if r.plain {
		fmt.Printf("==> %s\n", name)
		return
	}
	pterm.Success.WithWriter(r.multi.NewWriter()).Println(name)
	r.updateSpinnerLocked(name)
}

func (r *agmUpdateRenderer) onProgressLocked(line string) {
	if r.plain {
		if isAgmNotableLine(line) {
			fmt.Println("    " + line)
		}
		return
	}
	r.updateSpinnerLocked(line)
}

func isAgmNotableLine(line string) bool {
	return strings.Contains(line, "[ERROR]") || strings.Contains(line, "[WARN]")
}

func (r *agmUpdateRenderer) updateSpinnerLocked(line string) {
	if r.spinner == nil {
		return
	}
	r.spinner.UpdateText(truncateAgmRunes(line, agmUpdateSpinnerRunes))
}

// renderSuccess stops the UI and prints the final success line.
func (r *agmUpdateRenderer) renderSuccess() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.finished = true
	msg := r.title() + " completed successfully."
	if r.plain {
		for _, s := range r.stages {
			fmt.Printf("  ✓ %s\n", s)
		}
		fmt.Println(msg)
		return
	}
	if r.spinner != nil {
		r.spinner.Success(msg)
	}
	r.multi.Stop()
}

// renderFailure stops the UI and prints a structured error panel: what
// failed, the current stage, and the last N installer lines. No stack trace.
func (r *agmUpdateRenderer) renderFailure(runErr error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.finished = true
	errText := strings.TrimSpace(runErr.Error())
	if r.plain {
		r.renderFailurePlain(errText)
		return
	}
	if r.spinner != nil {
		r.spinner.Fail(r.action + " failed")
	}
	r.multi.Stop()
	pterm.Error.Println(r.title() + " failed: " + errText)
	pterm.DefaultBox.WithTitle("Update failed").Println(r.failureBody(errText))
}

func (r *agmUpdateRenderer) renderFailurePlain(errText string) {
	fmt.Printf("  ✗ %s failed: %s\n", r.title(), errText)
	fmt.Printf("  Stage: %s\n", r.failureStage())
	fmt.Println("  Last installer output:")
	for _, l := range r.tailLines() {
		fmt.Println("    " + truncateAgmRunes(l, 120))
	}
	fmt.Println("  Next step: re-run the command, or inspect the last error with `gitmap error export <file>`.")
}

func (r *agmUpdateRenderer) failureBody(errText string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Stage: %s\n", r.failureStage())
	fmt.Fprintf(&b, "Error: %s\n", truncateAgmRunes(errText, 200))
	b.WriteString("Last installer output:\n")
	for _, l := range r.tailLines() {
		b.WriteString("  " + truncateAgmRunes(l, 120) + "\n")
	}
	b.WriteString("Next step: re-run the command, or inspect the last error with `gitmap error export <file>`.")
	return b.String()
}

func (r *agmUpdateRenderer) failureStage() string {
	if r.stage == "" {
		return "startup"
	}
	return r.stage
}

func (r *agmUpdateRenderer) tailLines() []string {
	if len(r.log) <= agmUpdateTailLines {
		return r.log
	}
	return r.log[len(r.log)-agmUpdateTailLines:]
}

func (r *agmUpdateRenderer) appendLog(line string) {
	r.log = append(r.log, truncateAgmRunes(line, agmUpdateMaxLineRunes))
	if len(r.log) > agmUpdateMaxLogLines {
		r.log = r.log[len(r.log)-agmUpdateMaxLogLines:]
	}
}

func truncateAgmRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}

// runAgmInstallerStreamed runs the installer with piped stdio (never
// inherited), feeding every output line to the renderer. Stdin is not
// inherited: the installer must not fight gitmap for the terminal.
func runAgmInstallerStreamed(name string, args []string, env []string, r *agmUpdateRenderer) error {
	cmd := exec.Command(name, args...)
	if env != nil {
		cmd.Env = env
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, pipe := range []io.Reader{stdout, stderr} {
		wg.Add(1)
		go func(rd io.Reader) {
			defer wg.Done()
			agmFeedPipe(rd, r)
		}(pipe)
	}
	wg.Wait()
	return cmd.Wait()
}

func agmFeedPipe(rd io.Reader, r *agmUpdateRenderer) {
	scanner := bufio.NewScanner(rd)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		r.FeedLine(scanner.Text())
	}
}

// resolveAgmUnixInstallCmd builds the unix installer one-liner.
// GITMAP_AGM_INSTALL_SCRIPT overrides the source (test seam); default unchanged.
func resolveAgmUnixInstallCmd(version string) string {
	if mock := strings.TrimSpace(os.Getenv(agmUpdateEnvMockScript)); mock != "" {
		return fmt.Sprintf("bash %q --version %q", mock, strings.TrimPrefix(version, "v"))
	}
	if version != "" {
		clean := strings.TrimPrefix(version, "v")
		return fmt.Sprintf(`curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash -s -- --version '%s'`, clean)
	}
	return constants.AgManagerUnixInstallCmd
}
