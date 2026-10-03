package cmdssh

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	readPasswordHook   = term.ReadPassword
	readConsentHook    = bufio.NewReader(os.Stdin).ReadString
	dialTargetPassHook = dialNodeWithPassword
	tryConnectKeyHook  = tryConnectDefaultKey
	isTerminalHook     = isInteractiveTerminal
)

func isConsentAffirmative(input string) bool {
	trimmed := strings.TrimSpace(input)

	return strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes")
}

func formatPasswordPrompt(user, host string) string {
	return "Enter password for " + user + "@" + host + ": "
}

func formatConsentPrompt() string {
	return "Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: "
}

func readMaskedPassword(prompt string, fd int) (string, error) {
	fmt.Print(prompt)
	passBytes, err := readPasswordHook(fd)
	fmt.Println()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(passBytes)), nil
}

func readInteractiveConsent(r *bufio.Reader) (bool, error) {
	if r == nil {
		r = bufio.NewReader(os.Stdin)
	}

	line, err := r.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false, err
	}

	return isConsentAffirmative(line), nil
}

func verifyTargetPassword(target *SSHTarget, pass string) bool {
	client, err := dialTargetPassHook(target, pass)
	if err != nil {
		return false
	}

	if client != nil && client.Conn != nil {
		_ = client.Close()
	}

	return true
}

func buildAttemptPrompt(target *SSHTarget, attempt int) string {
	if attempt > 1 {
		return fmt.Sprintf("Enter password for %s@%s (attempt %d of 3): ", target.Username, target.IP, attempt)
	}

	return formatPasswordPrompt(target.Username, target.IP)
}

func reportFailedAttempt(target *SSHTarget, attempt int) {
	if attempt < 3 {
		fmt.Printf("✗ Authentication failed: invalid password for %s@%s. Please try again.\n", target.Username, target.IP)
	}
}

func promptAttempt(target *SSHTarget, attempt int) (string, bool) {
	prompt := buildAttemptPrompt(target, attempt)
	pass, err := readMaskedPassword(prompt, int(os.Stdin.Fd()))
	if err != nil || pass == "" {
		return "", false
	}

	isValid := verifyTargetPassword(target, pass)

	return pass, isValid
}

func promptAndVerifyPassword(ctx context.Context, target *SSHTarget) (string, bool) {
	for attempt := 1; attempt <= 3; attempt++ {
		if ctx.Err() != nil {
			return "", false
		}

		pass, isValid := promptAttempt(target, attempt)
		if isValid {
			return pass, true
		}

		reportFailedAttempt(target, attempt)
	}

	return "", false
}

func promptConsentAndPersist(ctx context.Context, alias string, target *SSHTarget, pass string) {
	fmt.Print(formatConsentPrompt())
	line, _ := readConsentHook('\n')
	if isConsentAffirmative(line) {
		saveExplicitPassword(ctx, alias, target, pass)
		fmt.Println("✓ Password saved in local RSA vault for future logins.")

		return
	}

	fmt.Println("ℹ Password will not be saved. Using for this session only.")
}

func hasKeyAuth(sshTarget *SSHTarget) bool {
	client := tryConnectKeyHook(sshTarget)
	if client == nil {
		return false
	}

	if client.Conn != nil {
		_ = client.Close()
	}

	return true
}

func promptAndPersistPassword(ctx context.Context, target string, sshTarget *SSHTarget) (string, error) {
	pass, isValid := promptAndVerifyPassword(ctx, sshTarget)
	if isValid {
		promptConsentAndPersist(ctx, target, sshTarget, pass)

		return pass, nil
	}

	return "", apperror.NewValidationError("authentication failed: invalid password")
}

func interceptSSHPasswordIfNeeded(ctx context.Context, target string, sshTarget *SSHTarget, currentPass string) (string, error) {
	if currentPass != "" {
		return currentPass, nil
	}

	if hasKeyAuth(sshTarget) {
		return "", nil
	}

	if !isTerminalHook() {
		return "", nil
	}

	return promptAndPersistPassword(ctx, target, sshTarget)
}
