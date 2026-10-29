package linkexpander

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandURL_ServiceActivation(t *testing.T) {
	raw := "https://serviceactivation.in/subscription/new/QuTYNtSw8PTD"
	res := ExpandURL(raw)
	if !res.IsExpanded {
		t.Fatalf("expected isExpanded to be true, got false")
	}
	expected := "https://serviceactivation.google.com/subscription/new/QuTYNtSw8PTD"
	if res.Expanded != expected {
		t.Fatalf("expected expanded URL %q, got %q", expected, res.Expanded)
	}
}

func TestExpandText_MultipleLinks(t *testing.T) {
	text := "Activation:\nLink 1: https://serviceactivation.in/subscription/new/QuTYNtSw8PTD\nLink 2: https://serviceactivation.in/subscription/new/w38mQH2MZ398\nThanks!"
	expandedText, links := ExpandText(text)
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
	if !strings.Contains(expandedText, "https://serviceactivation.google.com/subscription/new/QuTYNtSw8PTD") {
		t.Fatalf("expanded text missing Google link 1: %s", expandedText)
	}
	if !strings.Contains(expandedText, "https://serviceactivation.google.com/subscription/new/w38mQH2MZ398") {
		t.Fatalf("expanded text missing Google link 2: %s", expandedText)
	}
}

func TestExpandText_NoLinks(t *testing.T) {
	text := "Plain text without URLs"
	out, links := ExpandText(text)
	if len(links) != 0 {
		t.Fatalf("expected 0 links, got %d", len(links))
	}
	if out != text {
		t.Fatalf("expected text %q, got %q", text, out)
	}
}

func TestExecuteInstruction_ServiceActivationAndInterpolation(t *testing.T) {
	cfg := DefaultInstructionConfig()
	cfg.Links = []InstructionLink{
		{
			URL:          "https://serviceactivation.in/subscription/new/ALPHA100",
			Recipients:   []string{"user@example.com"},
			Subject:      "Notice for {{original}}",
			BodyTemplate: "Expanded to {{expanded}} from {{url}}",
		},
	}
	cfg.IsAutoSend = true

	res, err := ExecuteInstruction(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalLinks != 1 || res.ExpandedCount != 1 || res.SentCount != 1 {
		t.Fatalf("unexpected counts: %+v", res)
	}
	item := res.Results[0]
	if !item.IsExpanded || !item.IsSent {
		t.Fatalf("expected item expanded and sent: %+v", item)
	}
	expectedURL := "https://serviceactivation.google.com/subscription/new/ALPHA100"
	if item.ExpandedURL != expectedURL {
		t.Fatalf("expected %s, got %s", expectedURL, item.ExpandedURL)
	}
	if !strings.Contains(item.InterpolatedSubject, "https://serviceactivation.in/subscription/new/ALPHA100") {
		t.Fatalf("subject interpolation failed: %s", item.InterpolatedSubject)
	}
	if !strings.Contains(item.InterpolatedBody, expectedURL) {
		t.Fatalf("body interpolation failed: %s", item.InterpolatedBody)
	}
}

func TestExecuteInstruction_ProtocolCheck(t *testing.T) {
	cfg := DefaultInstructionConfig()
	cfg.Links = []InstructionLink{
		{URL: "file:///etc/passwd"},
		{URL: "javascript:alert(1)"},
	}

	res, err := ExecuteInstruction(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasErrors {
		t.Fatalf("expected res.HasErrors to be true")
	}
	if len(res.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res.Results))
	}
	if !strings.Contains(res.Results[0].Error, "unsupported protocol scheme 'file'") {
		t.Fatalf("expected file scheme error, got %s", res.Results[0].Error)
	}
	if !strings.Contains(res.Results[1].Error, "unsupported protocol scheme 'javascript'") {
		t.Fatalf("expected javascript scheme error, got %s", res.Results[1].Error)
	}
}

func TestExecuteInstruction_RecipientLimit(t *testing.T) {
	cfg := DefaultInstructionConfig()
	cfg.Links = []InstructionLink{
		{
			URL:        "https://example.com",
			Recipients: []string{"1@x.com", "2@x.com", "3@x.com", "4@x.com", "5@x.com", "6@x.com"},
		},
	}

	_, err := ExecuteInstruction(cfg)
	if err == nil || !strings.Contains(err.Error(), "recipient count exceeds maximum allowance of 5") {
		t.Fatalf("expected recipient limit error, got %v", err)
	}
}

func TestExecuteInstruction_LoopDetection(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, ts.URL+"/loop", http.StatusFound)
	}))
	defer ts.Close()

	cfg := DefaultInstructionConfig()
	cfg.Links = []InstructionLink{{URL: ts.URL + "/loop"}}

	res, err := ExecuteInstruction(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasErrors {
		t.Fatalf("expected hasErrors = true")
	}
	item := res.Results[0]
	if !strings.Contains(item.Error, "circular redirect detected") {
		t.Fatalf("expected circular redirect error, got %s", item.Error)
	}
}

func TestExecuteInstruction_MaxLevelsAndInfinite(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/1":
			http.Redirect(w, r, ts.URL+"/2", http.StatusFound)
		case "/2":
			http.Redirect(w, r, ts.URL+"/3", http.StatusFound)
		case "/3":
			http.Redirect(w, r, ts.URL+"/final", http.StatusFound)
		case "/final":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}
	}))
	defer ts.Close()

	cfg := DefaultInstructionConfig()
	cfg.MaxLevels = 2
	cfg.Links = []InstructionLink{{URL: ts.URL + "/1"}}

	res, err := ExecuteInstruction(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item := res.Results[0]
	if item.RedirectHops != 2 {
		t.Fatalf("expected 2 hops, got %d", item.RedirectHops)
	}
	if !strings.HasSuffix(item.ExpandedURL, "/3") {
		t.Fatalf("expected stop at /3, got %s", item.ExpandedURL)
	}

	cfg.IsInfinite = true
	resInf, errInf := ExecuteInstruction(cfg)
	if errInf != nil {
		t.Fatalf("unexpected error: %v", errInf)
	}
	itemInf := resInf.Results[0]
	if !strings.HasSuffix(itemInf.ExpandedURL, "/final") {
		t.Fatalf("expected final URL, got %s", itemInf.ExpandedURL)
	}
}

func TestExecuteInstruction_BodyLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		bigData := make([]byte, 10000)
		_, _ = w.Write(bigData)
	}))
	defer ts.Close()

	cfg := DefaultInstructionConfig()
	cfg.Safety.MaxBodyBytes = 1024
	cfg.Links = []InstructionLink{{URL: ts.URL}}

	res, err := ExecuteInstruction(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasErrors {
		t.Fatalf("expected no errors for body limit, got: %+v", res.Results[0])
	}
	if res.Results[0].StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.Results[0].StatusCode)
	}
}

func TestExecuteInstructionFile_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "instruction.json")

	cfg := DefaultInstructionConfig()
	cfg.Links = []InstructionLink{
		{
			URL:        "https://serviceactivation.in/subscription/new/FILETEST",
			Recipients: []string{"test@example.com"},
		},
	}
	data, marshalErr := json.Marshal(cfg)
	if marshalErr != nil {
		t.Fatalf("marshal error: %v", marshalErr)
	}
	if writeErr := os.WriteFile(filePath, data, 0600); writeErr != nil {
		t.Fatalf("write error: %v", writeErr)
	}

	res, execErr := ExecuteInstructionFile(filePath)
	if execErr != nil {
		t.Fatalf("execute file error: %v", execErr)
	}
	if res.TotalLinks != 1 || res.ExpandedCount != 1 {
		t.Fatalf("unexpected result: %+v", res)
	}
	expected := "https://serviceactivation.google.com/subscription/new/FILETEST"
	if res.Results[0].ExpandedURL != expected {
		t.Fatalf("expected %s, got %s", expected, res.Results[0].ExpandedURL)
	}
}
