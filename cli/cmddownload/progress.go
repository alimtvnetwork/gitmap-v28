package cmddownload

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// ProgressBar manages symmetric terminal progress rendering.
type ProgressBar struct {
	filename   string
	totalBytes int64
	startTime  time.Time
	lastTime   time.Time
	lastBytes  int64
	speedBps   int64
	isQuiet    bool
	isJSON     bool
	isTTY      bool
	lastWidth  int
	mu         sync.Mutex
}

// NewProgressBar creates a new ProgressBar tracker.
func NewProgressBar(filename string, totalBytes int64, isQuiet, isJSON bool) *ProgressBar {
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	return &ProgressBar{
		filename:   filename,
		totalBytes: totalBytes,
		startTime:  time.Now(),
		lastTime:   time.Now(),
		isQuiet:    isQuiet,
		isJSON:     isJSON,
		isTTY:      isTTY,
	}
}

// Update redraws the centered progress bar based on downloaded bytes.
func (p *ProgressBar) Update(downloaded int64) {
	if p.isQuiet || p.isJSON || !p.isTTY {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	p.updateSpeed(downloaded, now)

	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth <= 0 {
		termWidth = 80
	}

	line := p.renderLine(downloaded, termWidth)
	padding := (termWidth - len(line)) / 2
	if padding < 0 {
		padding = 0
	}

	leftPad := strings.Repeat(" ", padding)
	fullLine := "\r" + leftPad + line

	// Clear trailing characters if line shortened
	if len(fullLine) < p.lastWidth {
		fullLine += strings.Repeat(" ", p.lastWidth-len(fullLine))
	}
	p.lastWidth = len(fullLine)

	fmt.Print(fullLine)
}

func (p *ProgressBar) updateSpeed(downloaded int64, now time.Time) {
	elapsed := now.Sub(p.lastTime)
	if elapsed < 200*time.Millisecond {
		return
	}

	deltaBytes := downloaded - p.lastBytes
	secs := elapsed.Seconds()
	if deltaBytes > 0 && secs > 0 {
		p.speedBps = int64(float64(deltaBytes) / secs)
	}

	p.lastBytes = downloaded
	p.lastTime = now
}

func clampBarWidth(termWidth int) int {
	barWidth := int(float64(termWidth) * 0.35)
	if barWidth > 30 {
		return 30
	}
	if barWidth < 10 {
		return 10
	}
	return barWidth
}

func computePercent(downloaded, total int64) float64 {
	if total <= 0 {
		return 0
	}
	pct := float64(downloaded) / float64(total)
	if pct > 1.0 {
		return 1.0
	}
	return pct
}

func renderBarString(pct float64, barWidth int, hasTotal bool) string {
	if !hasTotal {
		return "[......]"
	}
	filled := int(math.Floor(pct * float64(barWidth)))
	if filled > barWidth {
		filled = barWidth
	}
	if filled <= 0 {
		return "[" + strings.Repeat(" ", barWidth) + "]"
	}
	if filled >= barWidth {
		return "[" + strings.Repeat("=", barWidth) + "]"
	}
	return "[" + strings.Repeat("=", filled-1) + ">" + strings.Repeat(" ", barWidth-filled) + "]"
}

func calculateETAString(totalBytes, downloaded, speedBps int64) string {
	if totalBytes <= 0 || speedBps <= 0 {
		return ""
	}
	remaining := totalBytes - downloaded
	if remaining <= 0 {
		return ""
	}
	etaSec := remaining / speedBps
	return fmt.Sprintf(" · ETA %02d:%02d", etaSec/60, etaSec%60)
}

func truncateFilename(name string) string {
	if len(name) > 24 {
		return name[:21] + "..."
	}
	return name
}

func (p *ProgressBar) renderLine(downloaded int64, termWidth int) string {
	barWidth := clampBarWidth(termWidth)
	hasTotal := p.totalBytes > 0
	percent := computePercent(downloaded, p.totalBytes)
	barString := renderBarString(percent, barWidth, hasTotal)
	etaString := calculateETAString(p.totalBytes, downloaded, p.speedBps)

	pctText := fmt.Sprintf("%3.0f%%", percent*100)
	bytesText := FormatBytes(downloaded)
	if hasTotal {
		bytesText += "/" + FormatBytes(p.totalBytes)
	}

	speedText := ""
	if p.speedBps > 0 {
		speedText = fmt.Sprintf(" · %s/s", FormatBytes(p.speedBps))
	}

	name := truncateFilename(p.filename)
	return fmt.Sprintf("Downloading %s  %s %s (%s)%s%s", name, barString, pctText, bytesText, speedText, etaString)
}

// Finish prints the completion line.
func (p *ProgressBar) Finish(destPath string, totalBytes int64, duration time.Duration) {
	if p.isQuiet || p.isJSON {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth <= 0 {
		termWidth = 80
	}

	durSec := duration.Seconds()
	if durSec <= 0 {
		durSec = 0.001
	}

	cleanLine := fmt.Sprintf("✔ Download complete: %s (%s in %.1fs)", destPath, FormatBytes(totalBytes), durSec)
	padding := (termWidth - len(cleanLine)) / 2
	if padding < 0 {
		padding = 0
	}

	leftPad := strings.Repeat(" ", padding)
	fmt.Printf("\r%s%s%s\n", strings.Repeat(" ", p.lastWidth), "\r", leftPad+cleanLine)
}

// FormatBytes formats byte count into human-readable representation.
func FormatBytes(bytes int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ProgressReader wraps an io.Reader to update ProgressBar.
type ProgressReader struct {
	Reader     io.Reader
	Bar        *ProgressBar
	Downloaded int64
}

// Read implements io.Reader.
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n <= 0 {
		return n, err
	}

	pr.Downloaded += int64(n)
	if pr.Bar != nil {
		pr.Bar.Update(pr.Downloaded)
	}

	return n, err
}
