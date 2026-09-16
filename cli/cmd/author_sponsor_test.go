// Package cmd — author_sponsor_test.go tests author and sponsor commands.
package cmd

import (
	"testing"
)

func TestRunAuthor(t *testing.T) {
	err := runAuthor([]string{})
	if err != nil {
		t.Fatalf("runAuthor failed: %v", err)
	}
}

func TestRunSponsor(t *testing.T) {
	err := runSponsor([]string{})
	if err != nil {
		t.Fatalf("runSponsor failed: %v", err)
	}
}

func TestRunCredits(t *testing.T) {
	err := runCredits([]string{})
	if err != nil {
		t.Fatalf("runCredits failed: %v", err)
	}
}
