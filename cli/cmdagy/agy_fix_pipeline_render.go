package cmdagy

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func toAbsPath(p string) string {
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}

	return p
}

func formatDisplayRepo(repo string) string {
	if len(repo) == 0 {
		return "(auto-detected current repo)"
	}

	return repo
}

func formatFailureStatus(hasFailures bool) string {
	if hasFailures {
		return constants.ColorRed + "failing error logs" + constants.ColorReset
	}

	return constants.ColorGreen + "clean status" + constants.ColorReset
}

func renderPayloadMetrics(repo, promptSource, logs, prompt, payload string, hasFailures bool) {
	fmt.Printf("    • Target Repo:    %s\n", formatDisplayRepo(repo))
	fmt.Printf("    • Pipeline Logs:  %.1f KB (%s)\n", float64(len(logs))/1024.0, formatFailureStatus(hasFailures))
	fmt.Printf("    • Fix Prompt:     %.1f KB (%s)\n", float64(len(prompt))/1024.0, promptSource)
	fmt.Printf("    • Total Payload:  %.1f KB\n", float64(len(payload))/1024.0)
	fmt.Printf("    • Saved Payload:  %s\n", toAbsPath(resolveActiveAgyPromptPath()))
}

func renderClipboardNotice(isNoClipboard bool) {
	if !isNoClipboard {
		fmt.Printf("    • Clipboard:      %sCopied to OS clipboard%s ✅\n",
			constants.ColorGreen, constants.ColorReset)
	}
	fmt.Printf("\n  %sReady! Prompt injected into Antigravity or paste via Ctrl+V to start fix loop.%s\n\n",
		constants.ColorYellow, constants.ColorReset)
}

func renderAgyFixFeedback(repo, promptSource, logs, prompt, payload string, hasFailures, noClip bool) {
	statusBanner := constants.ColorGreen + "✔" + constants.ColorReset
	fmt.Printf("\n  %s %sPrepared CI/CD pipeline fix prompt for Antigravity IDE!%s\n\n",
		statusBanner, constants.ColorCyan, constants.ColorReset)
	renderPayloadMetrics(repo, promptSource, logs, prompt, payload, hasFailures)
	renderClipboardNotice(noClip)
}

func renderQueuedVerificationNotice() {
	queuedAbs := toAbsPath(resolveQueuedAgyPromptPath())
	ledgerAbs := toAbsPath(resolveAgyPromptQueuePath())
	fmt.Printf("  %s✓ Follow-up Verification Prompt Queued: %s (\"Is it fixed?\")%s\n",
		constants.ColorGreen, queuedAbs, constants.ColorReset)
	fmt.Printf("    • Queue Ledger:   %s\n", ledgerAbs)
	if agyPath, hasAgy := resolveAntigravityBinary(); hasAgy {
		fmt.Printf("    • Antigravity CLI: Detected at %s\n", agyPath)
	}
	fmt.Println()
}

func renderAgyFixDryRun(repo, promptSource, logs, prompt, payload string, hasFailures bool) {
	fmt.Printf("\n  %s[dry-run] CI/CD Pipeline Fix Feed Preview:%s\n\n",
		constants.ColorYellow, constants.ColorReset)
	renderPayloadMetrics(repo, promptSource, logs, prompt, payload, hasFailures)
	fmt.Println("\n  [dry-run] Skipping clipboard copy, disk persistence, and Antigravity injection.")
}
