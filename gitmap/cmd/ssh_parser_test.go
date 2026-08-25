package cmd

import "testing"

func TestSSHTarget(t *testing.T) {
	target := SSHTarget{
		Username: "root",
		IP:       "10.0.0.1",
		Port:     22,
	}
	expected := "root@10.0.0.1"
	if got := target.String(); got != expected {
		t.Errorf("SSHTarget.String() = %q, want %q", got, expected)
	}
}
