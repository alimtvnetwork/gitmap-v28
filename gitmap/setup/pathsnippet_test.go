package setup

import (
	"strings"
	"testing"
)

func TestRenderPathSnippet_Bash_ContainsMarkerAndPath(t *testing.T) {
	got, err := RenderPathSnippet("bash", "/opt/gitmap", "run.sh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wants := []string{
		"# gitmap shell wrapper v2 - managed by run.sh. Do not edit manually.",
		"export GITMAP_WRAPPER=1",
		`export PATH="$PATH:/opt/gitmap"`,
		"# gitmap shell wrapper v2 end",
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("snippet missing %q\nfull output:\n%s", want, got)
		}
	}
}

func TestRenderPathSnippet_Fish_UsesFishAddPath(t *testing.T) {
	got, err := RenderPathSnippet("fish", "/opt/gitmap", "installer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "fish_add_path /opt/gitmap") {
		t.Errorf("fish snippet should call fish_add_path; got:\n%s", got)
	}
	if !strings.Contains(got, "set -gx GITMAP_WRAPPER 1") {
		t.Errorf("fish snippet should set wrapper var; got:\n%s", got)
	}
}

func TestRenderPathSnippet_Pwsh_InstallsCommandWrapper(t *testing.T) {
	got, err := RenderPathSnippet("pwsh", `C:\gitmap-cli`, "installer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"function global:gitmap", "function global:gcd", "GITMAP_COMMAND_WRAPPER"} {
		if !strings.Contains(got, want) {
			t.Errorf("pwsh snippet missing %q\nfull output:\n%s", want, got)
		}
	}
}

func TestRenderPathSnippet_UnknownShell_Errors(t *testing.T) {
	_, err := RenderPathSnippet("ksh", "/opt/gitmap", "run.sh")
	if err == nil {
		t.Fatal("expected error for unknown shell")
	}
}

func TestRenderPathSnippet_EmptyDir_Errors(t *testing.T) {
	_, err := RenderPathSnippet("bash", "", "run.sh")
	if err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestRenderPathSnippet_DefaultManager(t *testing.T) {
	got, err := RenderPathSnippet("bash", "/opt/gitmap", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "managed by gitmap setup") {
		t.Errorf("empty manager should default to 'gitmap setup'; got:\n%s", got)
	}
}

func TestMarkerOpenFor_Format(t *testing.T) {
	got := MarkerOpenFor("run.sh")
	want := "# gitmap shell wrapper v2 - managed by run.sh. Do not edit manually."
	if got != want {
		t.Errorf("MarkerOpenFor: got %q want %q", got, want)
	}
}
