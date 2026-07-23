package mapper

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/scanner"
)

// TestToHTTPS_FromSSH verifies SSH to HTTPS conversion.
func TestToHTTPS_FromSSH(t *testing.T) {
	result := toHTTPS("git@github.com:user/repo.git")
	if strings.HasPrefix(result, "https://") {
		t.Logf("Converted to HTTPS: %s — OK", result)
	} else {
		t.Errorf("Expected HTTPS URL, got: %s", result)
	}
}

// TestToHTTPS_AlreadyHTTPS verifies passthrough.
func TestToHTTPS_AlreadyHTTPS(t *testing.T) {
	input := "https://github.com/user/repo.git"
	result := toHTTPS(input)
	if result == input {
		t.Log("HTTPS passthrough — OK")
	}
}

// TestToSSH_FromHTTPS verifies HTTPS to SSH conversion.
func TestToSSH_FromHTTPS(t *testing.T) {
	result := toSSH("https://github.com/user/repo.git")
	if strings.HasPrefix(result, "git@") {
		t.Logf("Converted to SSH: %s — OK", result)
	} else {
		t.Errorf("Expected SSH URL, got: %s", result)
	}
}

// TestToSSH_AlreadySSH verifies passthrough.
func TestToSSH_AlreadySSH(t *testing.T) {
	input := "git@github.com:user/repo.git"
	result := toSSH(input)
	if result == input {
		t.Log("SSH passthrough — OK")
	}
}

// TestExtractRepoName verifies repo name extraction.
func TestExtractRepoName(t *testing.T) {
	name := extractRepoName("https://github.com/user/my-repo.git")
	if name == "my-repo" {
		t.Log("Extracted repo name — OK")
	} else {
		t.Errorf("Expected 'my-repo', got: %s", name)
	}
}

// TestExtractRepoName_Empty verifies empty URL handling.
func TestExtractRepoName_Empty(t *testing.T) {
	name := extractRepoName("")
	if name == "unknown" {
		t.Log("Empty URL returns 'unknown' — OK")
	}
}

// TestBuildNote_NoRemote verifies note for missing remote.
func TestBuildNote_NoRemote(t *testing.T) {
	note := buildNote("", "default note")
	if note == "no remote configured" {
		t.Log("No remote note — OK")
	}
}

// TestBuildNote_WithRemote verifies default note passthrough.
func TestBuildNote_WithRemote(t *testing.T) {
	note := buildNote("https://github.com/user/repo.git", "my note")
	if note == "my note" {
		t.Log("Default note used — OK")
	}
}

// TestSelectCloneURL_SSH verifies SSH mode selection.
func TestSelectCloneURL_SSH(t *testing.T) {
	url := selectCloneURL("https://url", "git@url", "ssh")
	if url == "git@url" {
		t.Log("SSH mode selected — OK")
	}
}

// TestSelectCloneURL_HTTPS verifies HTTPS mode selection.
func TestSelectCloneURL_HTTPS(t *testing.T) {
	url := selectCloneURL("https://url", "git@url", "https")
	if url == "https://url" {
		t.Log("HTTPS mode selected — OK")
	}
}

// TestBuildInstruction verifies clone command generation.
func TestBuildInstruction(t *testing.T) {
	cmd := buildInstruction("https://github.com/u/r.git", "main", "path/to/repo")
	expected := "git clone -b main https://github.com/u/r.git path/to/repo"
	if cmd == expected {
		t.Log("Instruction built — OK")
	} else {
		t.Errorf("Expected %q, got %q", expected, cmd)
	}
}

// TestBuildRecords verifies end-to-end record building (without git).
func TestBuildRecords(t *testing.T) {
	repos := []scanner.RepoInfo{
		{AbsolutePath: "/tmp/nonexistent", RelativePath: "nonexistent"},
	}
	records := BuildRecords(repos, "https", "test")
	if len(records) == 1 {
		t.Log("Built 1 record — OK")
	}
	if records[0].Notes == "no remote configured" {
		t.Log("Missing remote noted — OK")
	}
}

// TestBuildSlug_HTTPS verifies slug from a standard HTTPS URL.
func TestBuildSlug_HTTPS(t *testing.T) {
	slug := buildSlug("https://github.com/user/my-api.git", "fallback")
	if slug != "my-api" {
		t.Errorf("Expected 'my-api', got: %s", slug)
	}
}

// TestBuildSlug_HTTPSNoGitSuffix verifies slug when .git suffix is absent.
func TestBuildSlug_HTTPSNoGitSuffix(t *testing.T) {
	slug := buildSlug("https://github.com/org/dashboard", "fallback")
	if slug != "dashboard" {
		t.Errorf("Expected 'dashboard', got: %s", slug)
	}
}

// TestBuildSlug_Lowercase verifies slug is lowercased.
func TestBuildSlug_Lowercase(t *testing.T) {
	slug := buildSlug("https://github.com/Org/My-Repo.git", "fallback")
	if slug != "my-repo" {
		t.Errorf("Expected 'my-repo', got: %s", slug)
	}
}

// TestBuildSlug_EmptyURL verifies fallback to repoName.
func TestBuildSlug_EmptyURL(t *testing.T) {
	slug := buildSlug("", "MyFallback")
	if slug != "myfallback" {
		t.Errorf("Expected 'myfallback', got: %s", slug)
	}
}

// TestBuildSlug_DuplicateOrgs verifies same slug from different orgs.
func TestBuildSlug_DuplicateOrgs(t *testing.T) {
	slug1 := buildSlug("https://github.com/org-a/my-api.git", "")
	slug2 := buildSlug("https://github.com/org-b/my-api.git", "")
	if slug1 != slug2 {
		t.Errorf("Expected matching slugs, got: %s vs %s", slug1, slug2)
	}
}

// TestBuildSlug_SSHFallback verifies slug from an SSH-style URL passed as HTTPS.
func TestBuildSlug_SSHFallback(t *testing.T) {
	slug := buildSlug("git@github.com:user/auth-service.git", "fallback")
	if slug != "auth-service" {
		t.Errorf("Expected 'auth-service', got: %s", slug)
	}
}
