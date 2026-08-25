package cmd

import (
	"bytes"
	"testing"
)

func TestInteractiveSSHClient(t *testing.T) {
	client := InteractiveSSHClient{
		Stdin:  &bytes.Buffer{},
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}
	
	if client.Stdin == nil || client.Stdout == nil || client.Stderr == nil {
		t.Errorf("InteractiveSSHClient fields should not be nil")
	}
}
