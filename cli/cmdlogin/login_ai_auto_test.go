package cmdlogin

import (
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestParseLoginArgsAIAuto(t *testing.T) {
	args := []string{"--use-ai-auto", "--restriction-time"}
	opts, rest, err := parseLoginArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.hasAIAuto {
		t.Errorf("expected hasAIAuto=true")
	}
	if !opts.hasRestrictionTime {
		t.Errorf("expected hasRestrictionTime=true")
	}
	if len(rest) != 0 {
		t.Errorf("expected rest to be empty, got: %v", rest)
	}
}

func TestParseLoginArgsRevert(t *testing.T) {
	args := []string{"--revert-to-original"}
	opts, rest, err := parseLoginArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.hasRevertOriginal {
		t.Errorf("expected hasRevertOriginal=true")
	}
	if len(rest) != 0 {
		t.Errorf("expected rest to be empty, got: %v", rest)
	}
}

func TestAppendOrUpdateBotProfile(t *testing.T) {
	profiles := []model.GitProfile{
		{ID: "personal", Name: "Personal User", Email: "personal@example.com"},
	}

	botProf := model.GitProfile{
		ID:         aiAutoBotAlias,
		Name:       aiAutoBotName,
		Email:      aiAutoBotEmail,
		Provider:   "git",
		Type:       "user",
		AuthMethod: "config",
		IsDefault:  true,
		LastUsedAt: time.Now(),
	}

	updated := appendOrUpdateBotProfile(profiles, botProf)
	if len(updated) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(updated))
	}

	botProf.Email = "updated@rgitmap.com"
	reupdated := appendOrUpdateBotProfile(updated, botProf)
	if len(reupdated) != 2 {
		t.Fatalf("expected 2 profiles after update, got %d", len(reupdated))
	}
	if reupdated[1].Email != "updated@rgitmap.com" {
		t.Errorf("expected updated email, got %s", reupdated[1].Email)
	}
}

func TestResolveOriginalIdentity(t *testing.T) {
	cfg := model.GitProfileConfig{
		OriginalIdentity: &model.OriginalIdentity{
			Name:         "Original Dev",
			Email:        "dev@example.com",
			ProfileAlias: "dev-alias",
			SavedAt:      time.Now(),
		},
	}

	orig := resolveOriginalIdentity(&cfg)
	if orig == nil {
		t.Fatalf("expected non-nil original identity")
	}
	if orig.Name != "Original Dev" {
		t.Errorf("expected 'Original Dev', got %s", orig.Name)
	}
	if orig.Email != "dev@example.com" {
		t.Errorf("expected 'dev@example.com', got %s", orig.Email)
	}
}

func TestApplyRestrictionTimeDefaults(t *testing.T) {
	err := applyRestrictionTimeDefaults()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveAIBotName(t *testing.T) {
	verifyBotName(t, "", aiAutoBotName)
	verifyBotName(t, "   ", aiAutoBotName)
	verifyBotName(t, "muse", "gitmap automation bot - muse - runs autonomously")
	verifyBotName(t, "  gemini  ", "gitmap automation bot - gemini - runs autonomously")
	verifyBotName(t, "[claude]", "gitmap automation bot - claude - runs autonomously")
}

func verifyBotName(t *testing.T, input, expected string) {
	actual := resolveAIBotName(input)
	if actual != expected {
		t.Errorf("resolveAIBotName(%q) = %q; want %q", input, actual, expected)
	}
}

func TestIsAIBotIdentity(t *testing.T) {
	verifyIsAIBot(t, "", true)
	verifyIsAIBot(t, "   ", true)
	verifyIsAIBot(t, "gitmap automation bot - runs autonomously", true)
	verifyIsAIBot(t, "gitmap automation bot - runs autonmously", true)
	verifyIsAIBot(t, "gitmap automation bot - muse - runs autonomously", true)
	verifyIsAIBot(t, "John Doe", false)
	verifyIsAIBot(t, "developer@example.com", false)
}

func verifyIsAIBot(t *testing.T, name string, isExpected bool) {
	isBot := isAIBotIdentity(name)
	if isBot != isExpected {
		t.Errorf("isAIBotIdentity(%q) = %v; want %v", name, isBot, isExpected)
	}
}

func TestParseLoginSetArgs(t *testing.T) {
	testParseLoginSetBasic(t)
	testParseLoginSetGlobal(t)
	testParseLoginSetAgent(t)
	testParseLoginSetFlags(t)
	testParseLoginSetInvalid(t)
}

func testParseLoginSetBasic(t *testing.T) {
	opts, err := parseLoginSetArgs([]string{"ai-auto"})
	if err != nil || !opts.hasAIAuto || opts.isGlobal || len(opts.agentName) > 0 {
		t.Fatalf("unexpected result for ai-auto: %+v, err: %v", opts, err)
	}
}

func testParseLoginSetGlobal(t *testing.T) {
	opts1, err1 := parseLoginSetArgs([]string{"global", "ai-auto"})
	if err1 != nil || !opts1.hasAIAuto || !opts1.isGlobal {
		t.Fatalf("unexpected result for global ai-auto: %+v, err: %v", opts1, err1)
	}

	opts2, err2 := parseLoginSetArgs([]string{"ai-auto", "--global"})
	if err2 != nil || !opts2.hasAIAuto || !opts2.isGlobal {
		t.Fatalf("unexpected result for ai-auto --global: %+v, err: %v", opts2, err2)
	}
}

func testParseLoginSetAgent(t *testing.T) {
	opts, err := parseLoginSetArgs([]string{"global", "ai-auto", "muse"})
	if err != nil || !opts.hasAIAuto || !opts.isGlobal || opts.agentName != "muse" {
		t.Fatalf("unexpected result for global ai-auto muse: %+v, err: %v", opts, err)
	}
}

func testParseLoginSetFlags(t *testing.T) {
	opts, err := parseLoginSetArgs([]string{"ai-auto", "--restriction-time"})
	if err != nil || !opts.hasAIAuto || !opts.hasRestrictionTime {
		t.Fatalf("unexpected result for restriction time: %+v, err: %v", opts, err)
	}

	helpOpts, helpErr := parseLoginSetArgs([]string{"--help"})
	if helpErr != nil || !helpOpts.isHelp {
		t.Fatalf("unexpected result for help: %+v, err: %v", helpOpts, helpErr)
	}
}

func testParseLoginSetInvalid(t *testing.T) {
	_, err := parseLoginSetArgs([]string{"--invalid-flag"})
	if err == nil {
		t.Fatalf("expected error for invalid flag, got nil")
	}
}

func TestMandatoryAgentValidation(t *testing.T) {
	err := runUseAIAutoWithAgent("", false, false)
	if err == nil {
		t.Fatalf("expected error when setting local ai-auto without agent name")
	}
	if !strings.Contains(err.Error(), "AI agent name is mandatory") {
		t.Errorf("expected mandatory agent name error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "Examples to teach on the fly") {
		t.Errorf("expected teaching examples in error, got: %v", err)
	}
}
