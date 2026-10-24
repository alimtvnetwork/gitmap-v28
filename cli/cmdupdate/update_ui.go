package cmdupdate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var ansiPattern = lazyregex.New("\x1b\\[[0-9;?]*[a-zA-Z]|\x1b[()][0-9A-Z]")

func stripAnsiEscapes(raw string) string {
	return ansiPattern.ReplaceAllString(raw, "")
}

type updateUIRendererImpl struct {
	mu                     sync.Mutex
	opts                   UpdateOptions
	currentVer             string
	plain                  bool
	spinner                *pterm.SpinnerPrinter
	current                UpdateStage
	stageMsg               string
	stages                 []string
	hasPrintedFirstSuccess bool
	isClosed               bool
}

// NewUpdateUIRenderer constructs the coordinated terminal UI renderer.
func NewUpdateUIRenderer(opts UpdateOptions, currentVer string) UpdateUIRenderer {
	r := &updateUIRendererImpl{
		opts:       opts,
		currentVer: currentVer,
		plain:      !isUpdateTTY() || !isVTSupported(),
	}
	r.renderHeader()
	return r
}

func (r *updateUIRendererImpl) renderHeader() {
	if r.opts.IsJSON || r.opts.IsQuiet {
		return
	}
	if r.plain {
		fmt.Printf("  ==> Starting Gitmap CLI Update (current: %s)...\n", FormatVersionTag(r.currentVer))
		return
	}
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaCyan)).Render("Gitmap CLI Update")
	ver := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c")).Render(fmt.Sprintf("(current: %s)", FormatVersionTag(r.currentVer)))
	fmt.Printf("\n  %s %s\n\n", title, ver)
}

func (r *updateUIRendererImpl) startSpinner(text string) {
	if r.plain || r.isClosed {
		return
	}
	r.stopSpinner()
	spinner, err := pterm.DefaultSpinner.
		WithWriter(os.Stdout).
		WithSequence("  ▀", "  ▄").
		WithRemoveWhenDone(true).
		WithShowTimer(false).
		Start(text)
	if err != nil {
		r.plain = true
		return
	}
	r.spinner = spinner
}

func (r *updateUIRendererImpl) stopSpinner() {
	if r.spinner != nil {
		r.spinner.RemoveWhenDone = true
		_ = r.spinner.Stop()
		r.spinner = nil
	}
	if !r.plain {
		fmt.Print("\r\033[2K")
	}
}

func (r *updateUIRendererImpl) OnStageStart(stage UpdateStage, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.opts.IsJSON || r.opts.IsQuiet || r.isClosed {
		return
	}
	r.current = stage
	r.stageMsg = message
	cleanMsg := stripAnsiEscapes(message)
	stepTag := fmt.Sprintf("[%d/%d] %s", stage.StepNumber(), stage.TotalSteps(), cleanMsg)
	if r.plain {
		fmt.Printf("==> %s\n", stepTag)
		return
	}
	if r.spinner != nil && r.spinner.IsActive {
		r.spinner.UpdateText(stepTag)
		return
	}
	r.startSpinner(stepTag)
}

func (r *updateUIRendererImpl) OnStageProgress(detail string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.opts.IsJSON || r.opts.IsQuiet || r.isClosed || r.plain {
		return
	}
	cleanDetail := stripAnsiEscapes(detail)
	if r.spinner != nil && r.spinner.IsActive && len(cleanDetail) > 0 {
		text := fmt.Sprintf("[%d/%d] %s: %s", r.current.StepNumber(), r.current.TotalSteps(), r.stageMsg, cleanDetail)
		r.spinner.UpdateText(text)
	}
}

func (r *updateUIRendererImpl) OnStageSuccess(stage UpdateStage, summary string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.opts.IsJSON || r.opts.IsQuiet || r.isClosed {
		return
	}
	cleanSum := stripAnsiEscapes(summary)
	r.stages = append(r.stages, cleanSum)

	r.stopSpinner()
	r.padFirstSuccess()

	if r.plain {
		fmt.Printf("    ✓ %s\n", cleanSum)
		return
	}
	badge := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaGreen)).Render("SUCCESS")
	fmt.Printf("    %s %s\n", badge, cleanSum)
}

func (r *updateUIRendererImpl) padFirstSuccess() {
	if r.hasPrintedFirstSuccess {
		return
	}
	r.hasPrintedFirstSuccess = true
	fmt.Println()
}

func (r *updateUIRendererImpl) OnStageFail(stage UpdateStage, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.isClosed {
		return
	}
	r.isClosed = true
	r.stopSpinner()
	if r.opts.IsJSON {
		outputResolveError(err, r.currentVer, r.opts.TargetVersion, true)
		return
	}
	if r.opts.IsQuiet {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	r.renderFailureCard(stage, err)
}

func (r *updateUIRendererImpl) renderFailureCard(stage UpdateStage, err error) {
	if r.plain {
		fmt.Fprintf(os.Stderr, "  ✗ [%s] failed: %v\n", stage.String(), err)
		return
	}
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaRed))
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(DraculaRed)).Padding(0, 1)
	body := fmt.Sprintf("%s\n\nStage: %s\nError: %v\nHint : Run with --verbose or run gitmap error export",
		titleStyle.Render("✖ Gitmap Update Failed"), stage.String(), err)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, boxStyle.Render(body))
}

func (r *updateUIRendererImpl) RenderSummaryCard(result UpdateResult) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.isClosed {
		return
	}
	r.isClosed = true
	r.stopSpinner()
	if r.opts.IsJSON {
		enc, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(enc))
		return
	}
	if r.opts.IsQuiet {
		return
	}
	if r.plain {
		r.renderPlainSummary(result)
		return
	}
	r.renderLipglossSummary(result)
}

func (r *updateUIRendererImpl) renderPlainSummary(result UpdateResult) {
	sec := float64(result.DurationMs) / 1000.0
	fmt.Println()
	fmt.Println("  ============================================================")
	fmt.Printf("  ✔ Update Complete: %s → %s (took %.1fs)\n", result.PreviousVersion, result.CurrentVersion, sec)
	fmt.Printf("  Binary: %s\n", result.BinaryPath)
	printPlainReleaseNotes(result.ReleaseNotes)
	fmt.Println("  ============================================================")
}

func (r *updateUIRendererImpl) renderLipglossSummary(result UpdateResult) {
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(DraculaPurple)).
		Padding(0, 1).
		MarginLeft(2)

	lines := formatSummaryLines(result)
	fmt.Println()
	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
}

func (r *updateUIRendererImpl) RenderAlreadyUpdated(currentVer, targetVer string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.isClosed {
		return
	}
	r.isClosed = true
	r.stopSpinner()
	if r.opts.IsJSON {
		outputAlreadyUpdated(currentVer, targetVer, true)
		return
	}
	if r.opts.IsQuiet {
		return
	}
	if r.plain {
		fmt.Printf("  Gitmap CLI is already up to date (%s). Use --force to reinstall.\n", FormatVersionTag(currentVer))
		return
	}
	r.renderLipglossAlreadyUpdated(currentVer)
}

func (r *updateUIRendererImpl) renderLipglossAlreadyUpdated(currentVer string) {
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(DraculaCyan)).
		Padding(0, 1).
		MarginLeft(2)

	verTag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c")).Render(FormatVersionTag(currentVer))
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaGreen)).Render(fmt.Sprintf("Gitmap CLI is already up to date (%s)", verTag))
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color(DraculaGray)).Render("No update required. Pass --force to reinstall.")
	fmt.Println()
	fmt.Println(cardStyle.Render(fmt.Sprintf("%s\n%s", title, hint)))
}

func (r *updateUIRendererImpl) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.isClosed {
		return
	}
	r.isClosed = true
	r.stopSpinner()
}
