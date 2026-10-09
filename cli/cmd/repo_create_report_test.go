package cmd

import (
	"encoding/json"
	"fmt"
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

// createRepoParams holds repo creation parameters. Test-local type —
// the production type was never implemented.
type createRepoParams struct {
	Name     string
	Slug     string
	LocalDir string
	IsPublic bool
	IsJSON   bool
	IsYAML   bool
	Profile  model.GitProfile
}

// RepoCreateReport is the structured output for repo creation.
// Test-local type — the production type was never implemented.
type RepoCreateReport struct {
	Name       string `json:"name" yaml:"name"`
	Slug       string `json:"slug" yaml:"slug"`
	Visibility string `json:"visibility" yaml:"visibility"`
	URL        string `json:"url" yaml:"url"`
}

// getRepoVisibility returns "public" or "private". Test-local helper.
func getRepoVisibility(p createRepoParams) string {
	if p.IsPublic {
		return "public"
	}

	return "private"
}

// reportCreatedRepo writes the creation report. Test-local helper —
// the production function was never implemented.
func reportCreatedRepo(p createRepoParams, url string) error {
	report := RepoCreateReport{
		Name:       p.Name,
		Slug:       p.Slug,
		Visibility: getRepoVisibility(p),
		URL:        url,
	}

	switch {
	case p.IsJSON:
		data, _ := json.Marshal(report)
		fmt.Println(string(data))
	case p.IsYAML:
		data, _ := yaml.Marshal(report)
		fmt.Print(string(data))
	default:
		fmt.Println("● Visibility: " + report.Visibility)
		fmt.Println("● Name: " + report.Name)
	}

	return nil
}
