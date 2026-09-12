package cmd

import (
	"testing"
)

func TestParseBrowseArgs(t *testing.T) {
	url, isChrome := parseBrowseArgs([]string{"https://example.com"})
	if url != "https://example.com" || isChrome {
		t.Fatalf("expected url='https://example.com', isChrome=false; got url=%q, isChrome=%v", url, isChrome)
	}

	url2, isChrome2 := parseBrowseArgs([]string{"github.com", "--chrome"})
	if url2 != "github.com" || !isChrome2 {
		t.Fatalf("expected url='github.com', isChrome=true; got url=%q, isChrome=%v", url2, isChrome2)
	}
}

func TestNormalizeBrowseURL(t *testing.T) {
	if got := normalizeBrowseURL("github.com"); got != "https://github.com" {
		t.Fatalf("expected https://github.com, got %s", got)
	}

	if got := normalizeBrowseURL("http://localhost:8080"); got != "http://localhost:8080" {
		t.Fatalf("expected http://localhost:8080, got %s", got)
	}

	if got := normalizeBrowseURL("chrome://settings"); got != "chrome://settings" {
		t.Fatalf("expected chrome://settings, got %s", got)
	}
}

func TestHasURLProtocol(t *testing.T) {
	if !hasURLProtocol("https://google.com") {
		t.Fatal("expected https://google.com to have protocol")
	}

	if !hasURLProtocol("http://localhost") {
		t.Fatal("expected http://localhost to have protocol")
	}

	if !hasURLProtocol("chrome://version") {
		t.Fatal("expected chrome://version to have protocol")
	}

	if hasURLProtocol("google.com") {
		t.Fatal("expected google.com not to have protocol")
	}
}
