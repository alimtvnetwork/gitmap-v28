package uipref

import "testing"

func TestIsQuietHonorsEnv(t *testing.T) {
	t.Setenv(EnvQuiet, "1")
	if !IsQuiet() {
		t.Fatal("expected IsQuiet=true when GITMAP_QUIET=1")
	}
	t.Setenv(EnvQuiet, "0")
	if IsQuiet() {
		t.Fatal("expected IsQuiet=false when GITMAP_QUIET=0")
	}
}

func TestIsNoColorHonorsStdEnv(t *testing.T) {
	t.Setenv(EnvNoColorStd, "")
	if !IsNoColor() {
		t.Fatal("expected IsNoColor=true when NO_COLOR is set (even empty)")
	}
}

func TestIsNoColorHonorsProjectEnv(t *testing.T) {
	t.Setenv(EnvNoColor, "1")
	if !IsNoColor() {
		t.Fatal("expected IsNoColor=true when GITMAP_NO_COLOR=1")
	}
}

func TestIsMacroPwdVisibleDefaultAndOverride(t *testing.T) {
	ResetMacroPwdOverride()
	defer ResetMacroPwdOverride()

	if !IsMacroPwdVisible() {
		t.Fatal("expected IsMacroPwdVisible=true by default")
	}

	SetMacroPwdOverride(false)
	if IsMacroPwdVisible() {
		t.Fatal("expected IsMacroPwdVisible=false after override false")
	}

	SetMacroPwdOverride(true)
	if !IsMacroPwdVisible() {
		t.Fatal("expected IsMacroPwdVisible=true after override true")
	}
}

func TestIsMacroPwdVisibleHonorsEnv(t *testing.T) {
	ResetMacroPwdOverride()
	defer ResetMacroPwdOverride()

	t.Setenv(EnvMacroPwd, "0")
	if IsMacroPwdVisible() {
		t.Fatal("expected IsMacroPwdVisible=false when GITMAP_MACRO_PWD=0")
	}

	t.Setenv(EnvMacroPwd, "1")
	if !IsMacroPwdVisible() {
		t.Fatal("expected IsMacroPwdVisible=true when GITMAP_MACRO_PWD=1")
	}
}

func TestSetMacroShowPwdPersists(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)
	ResetMacroPwdOverride()
	defer ResetMacroPwdOverride()

	if err := SetMacroShowPwd(false); err != nil {
		t.Fatalf("SetMacroShowPwd(false) failed: %v", err)
	}

	ResetMacroPwdOverride()
	if IsMacroPwdVisible() {
		t.Fatal("expected IsMacroPwdVisible=false after persisting false")
	}
}

