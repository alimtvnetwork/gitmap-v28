package cmdide

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/setup"
)

func TestParseSetupOptions_Defaults(t *testing.T) {
	opts, err := parseSetupOptions([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Target != setup.TargetAll {
		t.Errorf("expected target 'all', got %q", opts.Target)
	}

	if opts.Provider != setup.ProviderKimi {
		t.Errorf("expected provider 'kimi', got %q", opts.Provider)
	}

	if opts.IsDryRun {
		t.Errorf("expected isDryRun false, got true")
	}
}

func TestParseSetupOptions_TargetAndProvider(t *testing.T) {
	args := []string{"cline", "-p", "deepseek", "-m", "deepseek-coder", "-k", "sk-custom"}
	opts, err := parseSetupOptions(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Target != setup.TargetCline {
		t.Errorf("expected target 'cline', got %q", opts.Target)
	}

	if opts.Provider != setup.ProviderDeepSeek {
		t.Errorf("expected provider 'deepseek', got %q", opts.Provider)
	}

	if opts.Model != "deepseek-coder" {
		t.Errorf("expected model 'deepseek-coder', got %q", opts.Model)
	}

	if opts.APIKey != "sk-custom" {
		t.Errorf("expected apiKey 'sk-custom', got %q", opts.APIKey)
	}
}

func TestParseSetupOptions_AllTargets(t *testing.T) {
	targets := map[string]setup.IDETargetType{
		"cline":    setup.TargetCline,
		"composer": setup.TargetCline,
		"continue": setup.TargetContinue,
		"tab":      setup.TargetContinue,
		"zed":      setup.TargetZed,
		"all":      setup.TargetAll,
	}

	for arg, expected := range targets {
		opts, err := parseSetupOptions([]string{arg})
		if err != nil {
			t.Errorf("unexpected error for %q: %v", arg, err)
		}
		if opts.Target != expected {
			t.Errorf("expected %q for %q, got %q", expected, arg, opts.Target)
		}
	}
}

func TestParseSetupOptions_BoolFlags(t *testing.T) {
	args := []string{"-n", "-j", "--test", "-q", "-f"}
	opts, err := parseSetupOptions(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.IsDryRun || !opts.IsJSON || !opts.IsTestEnabled || !opts.IsQuiet || !opts.IsForceCreate {
		t.Errorf("expected all bool flags true, got: %+v", opts)
	}
}

func TestParseSetupOptions_InvalidProvider(t *testing.T) {
	args := []string{"--provider", "non-existent-provider"}
	_, err := parseSetupOptions(args)
	if err == nil {
		t.Fatal("expected error for invalid provider, got nil")
	}
}

func TestRunIDESetup_DryRunJSON(t *testing.T) {
	out := captureStdout(t, func() {
		err := runIDESetup([]string{"--dry-run", "--json"})
		if err != nil {
			t.Fatalf("runIDESetup failed: %v", err)
		}
	})

	if !strings.Contains(out, `"isSuccess"`) {
		t.Errorf("expected json output with isSuccess, got: %s", out)
	}
	if !strings.Contains(out, `"target"`) {
		t.Errorf("expected json output with target, got: %s", out)
	}
}

func TestRunIDESetup_DryRunTerminal(t *testing.T) {
	out := captureStdout(t, func() {
		err := runIDESetup([]string{"zed", "--dry-run", "-p", "ollama"})
		if err != nil {
			t.Fatalf("runIDESetup failed: %v", err)
		}
	})

	if !strings.Contains(out, "GitMap AI IDE Setup") {
		t.Errorf("expected terminal header, got: %s", out)
	}
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run note, got: %s", out)
	}
}

func TestDispatchSubcommand_SetupAliases(t *testing.T) {
	aliases := []string{"setup", "cursor-setup", "setup-ai", "configure"}
	for _, alias := range aliases {
		err := dispatchSubcommand(alias, []string{"--dry-run", "-q"})
		if err != nil {
			t.Errorf("expected alias %q to succeed, got: %v", alias, err)
		}
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}
