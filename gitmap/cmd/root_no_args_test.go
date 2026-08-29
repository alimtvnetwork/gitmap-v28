package cmd

import (
	"os"
	"testing"
)

func TestRunNoArgsDoesNotPanic(t *testing.T) {
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"gitmap"}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Run() panicked with no arguments: %v", r)
		}
	}()

	Run()
}
