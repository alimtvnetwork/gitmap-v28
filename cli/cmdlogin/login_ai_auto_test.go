package cmdlogin

import (
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
