package cluster

import (
	"fmt"
	"sync"
)

const (
	// ClearScreenSeq is the ANSI escape code to clear the terminal screen.
	ClearScreenSeq = "\033[H\033[2J"

	// DefaultProgress is the starting progress.
	DefaultProgress = 0

	// FormatAggregated is the format string for aggregated progress.
	FormatAggregated = "Aggregated Cluster Progress: %d%%\n"

	// FormatNodeProgress is the format string for individual node progress.
	FormatNodeProgress = "Node %s: %d%%\n"

	// SeparatorLine is a visual separator.
	SeparatorLine = "--------------------------------\n"

	// MsgNoClients is shown when no clients are connected.
	MsgNoClients = "Cluster Progress: No clients connected yet.\n"
)

// TerminalUI displays aggregated cluster progress.
type TerminalUI struct {
	mu       sync.Mutex
	progress map[string]int
}

// NewTerminalUI creates a new TerminalUI.
func NewTerminalUI() *TerminalUI {
	return &TerminalUI{
		progress: make(map[string]int),
	}
}

// OnLogReceived updates the progress map and redraws the UI.
func (ui *TerminalUI) OnLogReceived(log ProgressLog) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.progress[log.ClientID] = log.Progress
	ui.render()
}

// render draws the current progress to the terminal.
func (ui *TerminalUI) render() {
	fmt.Print(ClearScreenSeq)

	totalProgress := DefaultProgress
	numClients := len(ui.progress)
	hasNoClients := numClients == DefaultProgress

	if isEmptyClients {
		fmt.Print(MsgNoClients)
		return
	}

	for _, p := range ui.progress {
		totalProgress += p
	}

	aggregated := totalProgress / numClients

	fmt.Printf(FormatAggregated, aggregated)
	fmt.Print(SeparatorLine)
	for id, p := range ui.progress {
		fmt.Printf(FormatNodeProgress, id, p)
	}
}
