package termgeneric

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// LineClearSequence is the standard ANSI escape code to return carriage and erase the entire line.
const LineClearSequence = "\r\033[2K"

// DefaultSpinnerFrames contains standard braille spinner animation characters.
var DefaultSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner provides a thread-safe terminal loading spinner.
type Spinner struct {
	mu         sync.Mutex
	frames     []string
	frameIndex int
	label      string
	out        io.Writer
	stopChan   chan struct{}
	isSpinning bool
}

// NewSpinner creates a new terminal spinner with default frames.
func NewSpinner(label string, out io.Writer) *Spinner {
	return &Spinner{
		frames:   DefaultSpinnerFrames,
		label:    label,
		out:      out,
		stopChan: make(chan struct{}),
	}
}

// SetLabel updates the label displayed next to the spinner.
func (s *Spinner) SetLabel(label string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.label = label
}

// SetFrames overrides the spinner animation frames.
func (s *Spinner) SetFrames(frames []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(frames) > 0 {
		s.frames = frames
		s.frameIndex = 0
	}
}

// Tick advances the spinner animation frame and writes the rendered line.
func (s *Spinner) Tick() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	frame := s.frames[s.frameIndex%len(s.frames)]
	s.frameIndex++
	rendered := fmt.Sprintf("%s %s", frame, s.label)

	hasWriter := s.out != nil
	if hasWriter {
		fmt.Fprintf(s.out, "%s%s", LineClearSequence, rendered)
	}

	return rendered
}

// Start begins a background goroutine ticking the spinner at interval.
func (s *Spinner) Start(interval time.Duration) {
	s.mu.Lock()
	if s.isSpinning {
		s.mu.Unlock()
		return
	}

	s.isSpinning = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go s.spinLoop(interval)
}

func (s *Spinner) spinLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.Tick()
		}
	}
}

// Stop terminates the background spinning goroutine.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isSpinning {
		return
	}

	s.isSpinning = false
	close(s.stopChan)
}

// Clear erases the current terminal line.
func (s *Spinner) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	hasWriter := s.out != nil
	if hasWriter {
		fmt.Fprint(s.out, LineClearSequence)
	}
}

// Success stops the spinner and writes a checkmark success line.
func (s *Spinner) Success(msg string) {
	s.Stop()
	s.Clear()

	hasWriter := s.out != nil
	if hasWriter {
		fmt.Fprintf(s.out, "✔ %s\n", msg)
	}
}

// Fail stops the spinner and writes an error cross failure line.
func (s *Spinner) Fail(msg string) {
	s.Stop()
	s.Clear()

	hasWriter := s.out != nil
	if hasWriter {
		fmt.Fprintf(s.out, "✖ %s\n", msg)
	}
}

// ProgressBarConfig parameterizes a terminal progress bar.
type ProgressBarConfig struct {
	Total      int64
	Width      int
	FilledChar string
	EmptyChar  string
	Label      string
	Writer     io.Writer
	HasPercent bool
	HasCount   bool
}

// ProgressBar renders a thread-safe terminal progress bar.
type ProgressBar struct {
	mu      sync.Mutex
	cfg     ProgressBarConfig
	current int64
}

// NewProgressBar constructs a progress bar with default styling.
func NewProgressBar(total int64, label string, out io.Writer) *ProgressBar {
	cfg := ProgressBarConfig{
		Total:      total,
		Width:      20,
		FilledChar: "█",
		EmptyChar:  "░",
		Label:      label,
		Writer:     out,
		HasPercent: true,
		HasCount:   true,
	}

	return NewProgressBarWithConfig(cfg)
}

// NewProgressBarWithConfig constructs a progress bar with specific settings.
func NewProgressBarWithConfig(cfg ProgressBarConfig) *ProgressBar {
	if cfg.Width <= 0 {
		cfg.Width = 20
	}
	if len(cfg.FilledChar) == 0 {
		cfg.FilledChar = "█"
	}
	if len(cfg.EmptyChar) == 0 {
		cfg.EmptyChar = "░"
	}

	return &ProgressBar{
		cfg: cfg,
	}
}

// SetCurrent sets the current progress value.
func (p *ProgressBar) SetCurrent(val int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = val
}

// SetLabel updates the label rendered alongside the progress bar.
func (p *ProgressBar) SetLabel(label string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cfg.Label = label
}

// Increment adds 1 to the current progress value.
func (p *ProgressBar) Increment() {
	p.Add(1)
}

// Add increases the current progress value by delta.
func (p *ProgressBar) Add(delta int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current += delta
}

// Current returns the current progress count.
func (p *ProgressBar) Current() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.current
}

// Total returns the target total value.
func (p *ProgressBar) Total() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.cfg.Total
}

// Percent returns the progress percentage as a float between 0 and 100.
func (p *ProgressBar) Percent() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.calcPercent()
}

func (p *ProgressBar) calcPercent() float64 {
	if p.cfg.Total <= 0 {
		return 100.0
	}

	ratio := float64(p.current) / float64(p.cfg.Total)
	if ratio > 1.0 {
		ratio = 1.0
	}
	if ratio < 0.0 {
		ratio = 0.0
	}

	return ratio * 100.0
}

// Render returns the formatted progress bar string.
func (p *ProgressBar) Render() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.renderLocked()
}

func (p *ProgressBar) renderLocked() string {
	pct := p.calcPercent()
	bar := p.formatBar(pct)

	var sb strings.Builder
	sb.WriteString("[" + bar + "]")
	p.appendProgressDetails(&sb, pct)

	return sb.String()
}

func (p *ProgressBar) formatBar(pct float64) string {
	filledCount := int(float64(p.cfg.Width) * (pct / 100.0))
	if filledCount > p.cfg.Width {
		filledCount = p.cfg.Width
	}
	emptyCount := p.cfg.Width - filledCount

	return strings.Repeat(p.cfg.FilledChar, filledCount) + strings.Repeat(p.cfg.EmptyChar, emptyCount)
}

func (p *ProgressBar) appendProgressDetails(sb *strings.Builder, pct float64) {
	if p.cfg.HasPercent {
		sb.WriteString(fmt.Sprintf(" %3.0f%%", pct))
	}
	if p.cfg.HasCount {
		sb.WriteString(fmt.Sprintf(" (%d/%d)", p.current, p.cfg.Total))
	}
	if len(p.cfg.Label) > 0 {
		sb.WriteString(" " + p.cfg.Label)
	}
}

// Draw thread-safely erases the current line and writes the progress bar.
func (p *ProgressBar) Draw() {
	p.mu.Lock()
	defer p.mu.Unlock()

	hasWriter := p.cfg.Writer != nil
	if !hasWriter {
		return
	}

	rendered := p.renderLocked()
	fmt.Fprintf(p.cfg.Writer, "%s%s", LineClearSequence, rendered)
}

// Clear erases the current line.
func (p *ProgressBar) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	hasWriter := p.cfg.Writer != nil
	if hasWriter {
		fmt.Fprint(p.cfg.Writer, LineClearSequence)
	}
}

// Finish erases the line and prints the final message followed by a newline.
func (p *ProgressBar) Finish(finalMsg string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	hasWriter := p.cfg.Writer != nil
	if !hasWriter {
		return
	}

	hasMsg := len(finalMsg) > 0
	if hasMsg {
		fmt.Fprintf(p.cfg.Writer, "%s%s\n", LineClearSequence, finalMsg)
		return
	}

	fmt.Fprint(p.cfg.Writer, "\n")
}
