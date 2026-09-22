package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"gopkg.in/yaml.v3"
)

func TestGetRepoVisibility(t *testing.T) {
	pPrivate := createRepoParams{IsPublic: false}
	if got := getRepoVisibility(pPrivate); got != "private" {
		t.Errorf("got %q, want 'private'", got)
	}

	pPublic := createRepoParams{IsPublic: true}
	if got := getRepoVisibility(pPublic); got != "public" {
		t.Errorf("got %q, want 'public'", got)
	}
}

func TestReportCreatedRepo_JSON(t *testing.T) {
	p := createRepoParams{
		Name: "pwp-mobile", Slug: "pwp-mobile", LocalDir: "/tmp/pwp-mobile",
		IsPublic: false, IsJSON: true,
		Profile: model.GitProfile{Name: "alimtvnetwork", Provider: "github"},
	}

	out, _ := captureStdout(t, func() int {
		_ = reportCreatedRepo(p, "https://github.com/alimtvnetwork/pwp-mobile")

		return 0
	})

	var report RepoCreateReport
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v, raw: %s", err, out)
	}

	if report.Visibility != "private" {
		t.Errorf("got visibility %q, want 'private'", report.Visibility)
	}
	if report.Name != "pwp-mobile" {
		t.Errorf("got name %q, want 'pwp-mobile'", report.Name)
	}
}

func TestReportCreatedRepo_YAML(t *testing.T) {
	p := createRepoParams{
		Name: "pwp-mobile", Slug: "pwp-mobile", LocalDir: "/tmp/pwp-mobile",
		IsPublic: false, IsYAML: true,
		Profile: model.GitProfile{Name: "alimtvnetwork", Provider: "github"},
	}

	out, _ := captureStdout(t, func() int {
		_ = reportCreatedRepo(p, "https://github.com/alimtvnetwork/pwp-mobile")

		return 0
	})

	var report RepoCreateReport
	if err := yaml.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("failed to unmarshal YAML: %v, raw: %s", err, out)
	}

	if report.Visibility != "private" {
		t.Errorf("got visibility %q, want 'private'", report.Visibility)
	}
	if report.Slug != "pwp-mobile" {
		t.Errorf("got slug %q, want 'pwp-mobile'", report.Slug)
	}
}

func TestReportCreatedRepo_Text(t *testing.T) {
	p := createRepoParams{
		Name: "pwp-mobile", Slug: "pwp-mobile", LocalDir: "/tmp/pwp-mobile",
		IsPublic: false,
		Profile:  model.GitProfile{Name: "alimtvnetwork", Provider: "github"},
	}

	out, _ := captureStdout(t, func() int {
		_ = reportCreatedRepo(p, "https://github.com/alimtvnetwork/pwp-mobile")

		return 0
	})

	if !strings.Contains(out, "● Visibility: private") {
		t.Errorf("expected text output to contain '● Visibility: private', got:\n%s", out)
	}
}
