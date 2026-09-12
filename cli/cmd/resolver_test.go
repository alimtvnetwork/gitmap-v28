package cmd

import (
	"testing"
)

func TestResolver(t *testing.T) {
	// testing resolveEndpointString would require a DB or mocks
	// for now, we just test the URL pass-through
	hit := resolveEndpointString("https://github.com/a/b")
	if hit != "https://github.com/a/b" {
		t.Fatalf("url failed: %s", hit)
	}
}
