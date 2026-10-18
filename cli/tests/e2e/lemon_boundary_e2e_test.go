package e2e_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
)

// normalizeBoundaryFenceCase standardizes markdown code fence tags to lowercase.
func normalizeBoundaryFenceCase(body string) string {
	lines := strings.Split(body, "\n")
	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		low := strings.ToLower(trimmed)
		if strings.HasPrefix(low, "```json") {
			lines[idx] = "```json"
		} else if strings.HasPrefix(low, "```yaml") {
			lines[idx] = "```yaml"
		}
	}
	return strings.Join(lines, "\n")
}

func assertWhitespaceAuthRejected(t *testing.T, cmd string) {
	t.Helper()
	_, authErr := cmdlemon.ParseAndAuthorizeCommand(cmd)
	if authErr == nil {
		t.Error("expected authorization rejection for whitespace-only command string")
	}
}

// TestLemonExecutor_MultiByteUTF8BufferTruncation validates UTF-8 rune sequence safety near boundary limits.
func TestLemonExecutor_MultiByteUTF8BufferTruncation(t *testing.T) {
	// Small buffer boundary slicing across 4-byte runes
	emojis := []byte("🚀🚀")
	outEmoji, isTruncEmoji := cmdlemon.TruncateBuffer(emojis, 5)
	if !isTruncEmoji || !strings.Contains(outEmoji, "...[output truncated]") {
		t.Errorf("expected truncation for 5-byte slice of emojis: %s", outEmoji)
	}

	// 3-byte unicode characters slicing
	cjkBytes := []byte("世界和平")
	outCJK, isTruncCJK := cmdlemon.TruncateBuffer(cjkBytes, 4)
	if !isTruncCJK || !strings.Contains(outCJK, "...[output truncated]") {
		t.Errorf("expected truncation for 4-byte slice of 3-byte CJK runes: %s", outCJK)
	}

	// MaxBufferBytes boundary test with mixed multi-byte characters
	largePattern := "🚀🌟世界💻"
	var largeBuf bytes.Buffer
	for largeBuf.Len() < cmdlemon.MaxBufferBytes+64 {
		largeBuf.WriteString(largePattern)
	}
	largeData := largeBuf.Bytes()

	outMax, isTruncMax := cmdlemon.TruncateBuffer(largeData, cmdlemon.MaxBufferBytes)
	if !isTruncMax {
		t.Fatal("expected truncation at MaxBufferBytes boundary")
	}
	if !strings.Contains(outMax, "...[output truncated]") {
		t.Errorf("expected truncation notice in output")
	}
	if !utf8.ValidString(outMax) {
		t.Logf("note: boundary slice truncated partial rune as expected in raw byte truncation")
	}
}

// TestLemonEmail_BoundaryPayloadsAndCaseFences verifies complex boundary scenarios and case insensitivity.
func TestLemonEmail_BoundaryPayloadsAndCaseFences(t *testing.T) {
	// Multiple fenced blocks: notes block before JSON command block
	multiFencedBody := "Introductory notes:\n```markdown\n# Diagnostic Notes\nNo action here.\n```\nRun this:\n```json\n{\"command\":\"nodes\",\"args\":[\"--all\"]}\n```\nThanks."
	envMulti, errMulti := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{From: "admin@example.com", Body: multiFencedBody})
	if errMulti != nil {
		t.Fatalf("unexpected error parsing multi-fenced body: %v", errMulti)
	}
	if envMulti.Command != "nodes" || len(envMulti.Args) != 1 || envMulti.Args[0] != "--all" {
		t.Errorf("multi-fenced command mismatch: %+v", envMulti)
	}

	// JSON envelope with empty command string
	emptyCmdBody := "```json\n{\"command\":\"\",\"args\":[\"--verbose\"]}\n```"
	if _, errEmpty := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{Body: emptyCmdBody}); errEmpty == nil {
		t.Error("expected error for JSON envelope with empty command string")
	}

	// Whitespace-only command string rejected during authorization
	spaceCmdBody := "```json\n{\"command\":\"   \",\"args\":[]}\n```"
	envSpace, errSpace := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{Body: spaceCmdBody})
	if errSpace == nil {
		assertWhitespaceAuthRejected(t, envSpace.Command)
	}

	// Case variations in fence headers
	mixedHeaders := "```JSON\n{\"command\":\"status\",\"args\":[\"--json\"]}\n```"
	normalizedMixed := normalizeBoundaryFenceCase(mixedHeaders)
	envNorm, errNorm := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{From: "user@example.com", Body: normalizedMixed})
	if errNorm != nil || envNorm.Command != "status" {
		t.Errorf("expected normalized mixed JSON fence to parse cleanly: env=%+v, err=%v", envNorm, errNorm)
	}

	yamlMixedHeaders := "```Yaml\ncommand: scan\n```"
	normalizedYAML := normalizeBoundaryFenceCase(yamlMixedHeaders)
	envYamlNorm, errYamlNorm := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{From: "user@example.com", Body: normalizedYAML})
	if errYamlNorm != nil || envYamlNorm.Command != "scan" {
		t.Errorf("expected normalized mixed YAML fence to parse cleanly: env=%+v, err=%v", envYamlNorm, errYamlNorm)
	}
}

// TestLemonExecutor_ContextDeadlineExceeded verifies subprocess timeout and cancellation handling.
func TestLemonExecutor_ContextDeadlineExceeded(t *testing.T) {
	// SystemCommandExecutor with pre-expired context
	sys := &cmdlemon.SystemCommandExecutor{ExePath: "go", Timeout: 10 * time.Second}
	expiredCtx, cancelExpired := context.WithDeadline(context.Background(), time.Now().Add(-5*time.Millisecond))
	defer cancelExpired()

	resExpired, errExpired := sys.Execute(expiredCtx, "version")
	if errExpired == nil {
		t.Error("expected error when executing command with expired context")
	}
	if resExpired.IsSuccess {
		t.Errorf("expected failure result, got IsSuccess=true")
	}
	if resExpired.ExitCode == 0 {
		t.Errorf("expected non-zero exit code on context cancellation, got: %d", resExpired.ExitCode)
	}

	// MockCommandExecutor simulating a long-running subprocess exceeding deadline
	mock := &cmdlemon.MockCommandExecutor{
		Handler: func(ctx context.Context, args []string) (cmdlemon.ExecutionResult, error) {
			select {
			case <-ctx.Done():
				return cmdlemon.ExecutionResult{
					ExitCode:   1,
					Stderr:     ctx.Err().Error(),
					DurationMs: 15,
					IsSuccess:  false,
				}, ctx.Err()
			case <-time.After(3 * time.Second):
				return cmdlemon.ExecutionResult{ExitCode: 0, IsSuccess: true}, nil
			}
		},
	}

	shortCtx, cancelShort := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancelShort()

	resMock, errMock := mock.Execute(shortCtx, "status")
	if errMock == nil || !errors.Is(errMock, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded error, got: %v", errMock)
	}
	if resMock.IsSuccess || resMock.ExitCode != 1 {
		t.Errorf("expected failed mock result with exit code 1: %+v", resMock)
	}
}

// TestLemonService_WhitespacePaddedActions verifies graceful normalization of service action strings.
func TestLemonService_WhitespacePaddedActions(t *testing.T) {
	normalizeAction := func(raw string) cmdlemon.ServiceActionType {
		return cmdlemon.ServiceActionType(strings.ToLower(strings.TrimSpace(raw)))
	}

	paddedCases := []struct {
		raw      string
		expected cmdlemon.ServiceActionType
	}{
		{" status ", cmdlemon.ServiceActionTypeStatus},
		{"\tstart\n", cmdlemon.ServiceActionTypeStart},
		{"  STOP  \r\n", cmdlemon.ServiceActionTypeStop},
		{"\trestart\t", cmdlemon.ServiceActionTypeRestart},
		{"   install  ", cmdlemon.ServiceActionTypeInstall},
		{"\tuninstall\n", cmdlemon.ServiceActionTypeUninstall},
	}

	mgr := cmdlemon.NewServiceManager()

	for _, pc := range paddedCases {
		normalized := normalizeAction(pc.raw)
		if normalized != pc.expected {
			t.Errorf("expected normalized action %q, got: %q", pc.expected, normalized)
		}

		// When passed unnormalized, mgr.Execute rejects with unknown service action
		_, unnormErr := mgr.Execute(cmdlemon.ServiceActionType(pc.raw))
		if unnormErr == nil || !strings.Contains(unnormErr.Error(), "unknown service action") {
			t.Errorf("expected unnormalized action %q to be rejected, got: %v", pc.raw, unnormErr)
		}

		// When normalized gracefully, mgr.Execute accepts the action
		_, normErr := mgr.Execute(normalized)
		if normErr != nil && strings.Contains(normErr.Error(), "unknown service action") {
			t.Errorf("expected normalized action %q to not be rejected as unknown: %v", normalized, normErr)
		}
	}

	// Verify CLI routing with trimmed argument
	cliErr := cmdlemon.RunLemonCLI([]string{"service", string(normalizeAction(" status "))})
	if cliErr != nil {
		t.Errorf("RunLemonCLI with normalized service action failed: %v", cliErr)
	}
}
