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

func TestParseSSHTarget(t *testing.T) {
	t.Run("user@ip", func(t *testing.T) {
		target, err := ParseSSHTarget("admin@192.168.1.1", "root", 22)
		if err != nil {
			t.Fatal(err)
		}
		if target.Username != "admin" || target.IP != "192.168.1.1" {
			t.Errorf("unexpected parsing result: %+v", target)
		}
	})

	t.Run("ip@user", func(t *testing.T) {
		target, err := ParseSSHTarget("192.168.1.1@admin", "root", 22)
		if err != nil {
			t.Fatal(err)
		}
		if target.Username != "admin" || target.IP != "192.168.1.1" {
			t.Errorf("unexpected parsing result: %+v", target)
		}
	})

	t.Run("no user", func(t *testing.T) {
		target, err := ParseSSHTarget("192.168.1.1", "root", 22)
		if err != nil {
			t.Fatal(err)
		}
		if target.Username != "root" || target.IP != "192.168.1.1" {
			t.Errorf("unexpected parsing result: %+v", target)
		}
	})
}
