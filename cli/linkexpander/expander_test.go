package linkexpander

import (
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
