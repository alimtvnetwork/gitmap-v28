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
