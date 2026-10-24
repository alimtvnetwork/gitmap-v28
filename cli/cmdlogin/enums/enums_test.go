package enums

import "testing"

func TestLoginAuthMethodType_ValidAndString(t *testing.T) {
	methods := AllLoginAuthMethods()
	if len(methods) != 4 {
		t.Fatalf("expected 4 auth methods, got %d", len(methods))
	}

	for _, m := range methods {
		if !m.IsValid() {
			t.Errorf("expected %s to be valid", m)
		}

		if m.String() != string(m) {
			t.Errorf("expected String() == %s, got %s", string(m), m.String())
		}
	}

	if LoginAuthMethodType("invalid").IsValid() {
		t.Errorf("expected invalid auth method to return false")
	}
}

func TestLoginAuthMethodType_Booleans(t *testing.T) {
	if !LoginAuthMethodSSH.IsSSH() || LoginAuthMethodToken.IsSSH() {
		t.Errorf("unexpected IsSSH result")
	}

	if !LoginAuthMethodToken.IsToken() || LoginAuthMethodBrowser.IsToken() {
		t.Errorf("unexpected IsToken result")
	}

	if !LoginAuthMethodBrowser.IsBrowser() || LoginAuthMethodStatus.IsBrowser() {
		t.Errorf("unexpected IsBrowser result")
	}

	if !LoginAuthMethodStatus.IsStatus() || LoginAuthMethodSSH.IsStatus() {
		t.Errorf("unexpected IsStatus result")
	}
}

func TestLoginAuthMethodType_ParseValid(t *testing.T) {
	parsedSSH, errSSH := ParseLoginAuthMethod("ssh")
	if errSSH != nil || parsedSSH != LoginAuthMethodSSH {
		t.Errorf("expected ssh, got %s, err: %v", parsedSSH, errSSH)
	}

	parsedPat, errPat := ParseLoginAuthMethodType(" PAT ")
	if errPat != nil || parsedPat != LoginAuthMethodToken {
		t.Errorf("expected token from pat, got %s, err: %v", parsedPat, errPat)
	}

	parsedBrowser, errBrowser := ParseLoginAuthMethod("browser-flow")
	if errBrowser != nil || parsedBrowser != LoginAuthMethodBrowser {
		t.Errorf("expected browser from browser-flow, got %s, err: %v", parsedBrowser, errBrowser)
	}

	parsedStatus, errStatus := ParseLoginAuthMethod("status")
	if errStatus != nil || parsedStatus != LoginAuthMethodStatus {
		t.Errorf("expected status, got %s, err: %v", parsedStatus, errStatus)
	}
}

func TestLoginAuthMethodType_ParseInvalidAndAliases(t *testing.T) {
	_, err := ParseLoginAuthMethod("unknown-method")
	if err == nil {
		t.Errorf("expected error for unknown auth method")
	}

	all := AllLoginAuthMethodTypes()
	if len(all) != 4 {
		t.Errorf("expected AllLoginAuthMethodTypes count 4, got %d", len(all))
	}
}

func TestLoginInteractiveChoiceType_ValidAndString(t *testing.T) {
	choices := AllLoginInteractiveChoices()
	if len(choices) != 4 {
		t.Fatalf("expected 4 choices, got %d", len(choices))
	}

	for _, c := range choices {
		if !c.IsValid() {
			t.Errorf("expected %s to be valid", c)
		}

		if c.String() != string(c) {
			t.Errorf("expected String() == %s, got %s", string(c), c.String())
		}
	}

	if LoginInteractiveChoiceType("9").IsValid() {
		t.Errorf("expected invalid choice to return false")
	}
}

func TestLoginInteractiveChoiceType_Booleans(t *testing.T) {
	if !LoginInteractiveChoiceSSH.IsSSH() || LoginInteractiveChoiceToken.IsSSH() {
		t.Errorf("unexpected IsSSH result")
	}

	if !LoginInteractiveChoiceToken.IsToken() || LoginInteractiveChoiceBrowser.IsToken() {
		t.Errorf("unexpected IsToken result")
	}

	if !LoginInteractiveChoiceBrowser.IsBrowser() || LoginInteractiveChoiceStatus.IsBrowser() {
		t.Errorf("unexpected IsBrowser result")
	}

	if !LoginInteractiveChoiceStatus.IsStatus() || LoginInteractiveChoiceSSH.IsStatus() {
		t.Errorf("unexpected IsStatus result")
	}
}

func TestLoginInteractiveChoiceType_ParseAndAliases(t *testing.T) {
	parsed, err := ParseLoginInteractiveChoice(" 1 ")
	if err != nil || parsed != LoginInteractiveChoiceSSH {
		t.Errorf("expected 1, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseLoginInteractiveChoiceType("2")
	if errAlias != nil || parsedAlias != LoginInteractiveChoiceToken {
		t.Errorf("expected 2, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseLoginInteractiveChoice("5")
	if invalidErr == nil {
		t.Errorf("expected error for choice 5")
	}

	if len(AllLoginInteractiveChoiceTypes()) != 4 {
		t.Errorf("expected AllLoginInteractiveChoiceTypes count 4")
	}
}

func TestGitAuthProviderType_ValidAndString(t *testing.T) {
	providers := AllGitAuthProviders()
	if len(providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(providers))
	}

	for _, p := range providers {
		if !p.IsValid() {
			t.Errorf("expected %s to be valid", p)
		}

		if p.String() != string(p) {
			t.Errorf("expected String() == %s, got %s", string(p), p.String())
		}
	}

	if GitAuthProviderType("bitbucket").IsValid() {
		t.Errorf("expected invalid provider to return false")
	}
}

func TestGitAuthProviderType_Booleans(t *testing.T) {
	if !GitAuthProviderGitHub.IsGitHub() || GitAuthProviderGitLab.IsGitHub() {
		t.Errorf("unexpected IsGitHub result")
	}

	if !GitAuthProviderGitLab.IsGitLab() || GitAuthProviderGitHub.IsGitLab() {
		t.Errorf("unexpected IsGitLab result")
	}
}

func TestGitAuthProviderType_ParseAndAliases(t *testing.T) {
	parsed, err := ParseGitAuthProvider(" GitHub ")
	if err != nil || parsed != GitAuthProviderGitHub {
		t.Errorf("expected github, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseGitAuthProviderType("gitlab")
	if errAlias != nil || parsedAlias != GitAuthProviderGitLab {
		t.Errorf("expected gitlab, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseGitAuthProvider("unknown-provider")
	if invalidErr == nil {
		t.Errorf("expected error for unknown provider")
	}

	if len(AllGitAuthProviderTypes()) != 2 {
		t.Errorf("expected AllGitAuthProviderTypes count 2")
	}
}

func TestClipboardToolType_ValidAndString(t *testing.T) {
	tools := AllClipboardTools()
	if len(tools) != 5 {
		t.Fatalf("expected 5 tools, got %d", len(tools))
	}

	for _, c := range tools {
		if !c.IsValid() {
			t.Errorf("expected %s to be valid", c)
		}

		if c.String() != string(c) {
			t.Errorf("expected String() == %s, got %s", string(c), c.String())
		}
	}

	if ClipboardToolType("unknown").IsValid() {
		t.Errorf("expected invalid tool to return false")
	}
}

func TestClipboardToolType_Booleans(t *testing.T) {
	if !ClipboardToolClip.IsClip() || ClipboardToolPbcopy.IsClip() {
		t.Errorf("unexpected IsClip result")
	}

	if !ClipboardToolPbcopy.IsPbcopy() || ClipboardToolWlCopy.IsPbcopy() {
		t.Errorf("unexpected IsPbcopy result")
	}

	if !ClipboardToolWlCopy.IsWlCopy() || ClipboardToolXclip.IsWlCopy() {
		t.Errorf("unexpected IsWlCopy result")
	}

	if !ClipboardToolXclip.IsXclip() || ClipboardToolXsel.IsXclip() {
		t.Errorf("unexpected IsXclip result")
	}

	if !ClipboardToolXsel.IsXsel() || ClipboardToolClip.IsXsel() {
		t.Errorf("unexpected IsXsel result")
	}
}

func TestClipboardToolType_ParseAndAliases(t *testing.T) {
	parsed, err := ParseClipboardTool(" CLIP ")
	if err != nil || parsed != ClipboardToolClip {
		t.Errorf("expected clip, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseClipboardToolType("wl-copy")
	if errAlias != nil || parsedAlias != ClipboardToolWlCopy {
		t.Errorf("expected wl-copy, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseClipboardTool("nonexistent")
	if invalidErr == nil {
		t.Errorf("expected error for nonexistent tool")
	}

	if len(AllClipboardToolTypes()) != 5 {
		t.Errorf("expected AllClipboardToolTypes count 5")
	}
}
