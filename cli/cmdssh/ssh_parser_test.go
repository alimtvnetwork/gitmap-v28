package cmdssh

import (
	"os"
	"testing"
)

func TestSSHTarget(t *testing.T) {
	target := SSHTarget{Username: "root", IP: "10.0.0.1", Port: 22}
	expected := "root@10.0.0.1"
	if got := target.String(); got != expected {
		t.Errorf("SSHTarget.String() = %q, want %q", got, expected)
	}
}

type targetTestCase struct {
	name     string
	raw      string
	wantUser string
	wantIP   string
	wantPort int
}

func getValidTargetCases() []targetTestCase {
	return []targetTestCase{
		{name: "ip@user", raw: "192.168.1.9@a", wantUser: "a", wantIP: "192.168.1.9", wantPort: 22},
		{name: "user@ip", raw: "a@192.168.1.9", wantUser: "a", wantIP: "192.168.1.9", wantPort: 22},
		{name: "no user", raw: "m1", wantUser: "root", wantIP: "m1", wantPort: 22},
		{name: "ip:port", raw: "10.0.0.1:2222", wantUser: "root", wantIP: "10.0.0.1", wantPort: 2222},
		{name: "user@ip:port", raw: "admin@10.0.0.1:2222", wantUser: "admin", wantIP: "10.0.0.1", wantPort: 2222},
		{name: "ipv6 bracketed", raw: "[2001:db8::1]", wantUser: "root", wantIP: "2001:db8::1", wantPort: 22},
		{name: "ipv6:port", raw: "[2001:db8::1]:2222", wantUser: "root", wantIP: "2001:db8::1", wantPort: 2222},
		{name: "user@ipv6:port", raw: "admin@[2001:db8::1]:2222", wantUser: "admin", wantIP: "2001:db8::1", wantPort: 2222},
		{name: "raw ipv6", raw: "2001:db8::1", wantUser: "root", wantIP: "2001:db8::1", wantPort: 22},
	}
}

func verifyTargetFields(t *testing.T, target *SSHTarget, tc targetTestCase) {
	if target.Username != tc.wantUser || target.IP != tc.wantIP || target.Port != tc.wantPort {
		t.Errorf("got %s@%s:%d, want %s@%s:%d",
			target.Username, target.IP, target.Port,
			tc.wantUser, tc.wantIP, tc.wantPort)
	}
}

func TestParseSSHTarget_Valid(t *testing.T) {
	for _, tc := range getValidTargetCases() {
		t.Run(tc.name, func(t *testing.T) {
			target, err := ParseSSHTarget(tc.raw, "root", 22)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			verifyTargetFields(t, target, tc)
		})
	}
}

func getInvalidTargetCases() []string {
	return []string{
		"",
		"10.0.0.1:70000",
		"10.0.0.1:0",
		"10.0.0.1:abc",
		"a@b@c",
		":22",
	}
}

func TestParseSSHTarget_Invalid(t *testing.T) {
	for _, raw := range getInvalidTargetCases() {
		t.Run("invalid_"+raw, func(t *testing.T) {
			_, err := ParseSSHTarget(raw, "root", 22)
			if err == nil {
				t.Errorf("expected error for raw %q, got nil", raw)
			}
		})
	}
}

func TestResolveDefaultUsername_User(t *testing.T) {
	oldUser := os.Getenv("USER")
	defer os.Setenv("USER", oldUser)

	os.Setenv("USER", "custom-user")
	if user := resolveDefaultUsername(); user != "custom-user" {
		t.Errorf("expected custom-user, got %s", user)
	}
}

func TestResolveDefaultUsername_Username(t *testing.T) {
	oldUser, oldUsername := os.Getenv("USER"), os.Getenv("USERNAME")
	defer func() {
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
	}()

	os.Unsetenv("USER")
	os.Setenv("USERNAME", "win-user")
	if user := resolveDefaultUsername(); user != "win-user" {
		t.Errorf("expected win-user, got %s", user)
	}
}

func TestResolveDefaultUsername_Root(t *testing.T) {
	oldUser, oldUsername := os.Getenv("USER"), os.Getenv("USERNAME")
	defer func() {
		os.Setenv("USER", oldUser)
		os.Setenv("USERNAME", oldUsername)
	}()

	os.Unsetenv("USER")
	os.Unsetenv("USERNAME")
	if user := resolveDefaultUsername(); user != "root" {
		t.Errorf("expected root fallback, got %s", user)
	}
}

func TestGenerateDefaultAlias(t *testing.T) {
	if alias := generateDefaultAlias("192.168.1.10", 22); alias != "host-192.168.1.10" {
		t.Errorf("expected host-192.168.1.10, got %s", alias)
	}

	if alias := generateDefaultAlias("192.168.1.10", 0); alias != "host-192.168.1.10" {
		t.Errorf("expected host-192.168.1.10, got %s", alias)
	}

	if alias := generateDefaultAlias("192.168.1.10", 2222); alias != "host-192.168.1.10-2222" {
		t.Errorf("expected host-192.168.1.10-2222, got %s", alias)
	}

	if alias := generateDefaultAlias("[2001:db8::1]", 22); alias != "host-2001:db8::1" {
		t.Errorf("expected host-2001:db8::1, got %s", alias)
	}
}

func TestParseSSHJoinOptions_PlainTarget(t *testing.T) {
	opts, err := ParseSSHJoinOptions([]string{"192.168.1.14"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Target.IP != "192.168.1.14" || opts.Alias != "host-192.168.1.14" {
		t.Errorf("unexpected options: %+v", opts)
	}
}

func TestParseSSHJoinOptions_PositionalAlias(t *testing.T) {
	opts, err := ParseSSHJoinOptions([]string{"192.168.1.14", "my-box"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Alias != "my-box" {
		t.Errorf("expected alias my-box, got %s", opts.Alias)
	}
}

func assertJoinAlias(t *testing.T, args []string, expectedAlias string) {
	opts, err := ParseSSHJoinOptions(args)
	if err != nil {
		t.Fatalf("unexpected error for %v: %v", args, err)
	}

	if opts.Alias != expectedAlias {
		t.Errorf("expected alias %s for %v, got %s", expectedAlias, args, opts.Alias)
	}
}

func TestParseSSHJoinOptions_NameFlags(t *testing.T) {
	assertJoinAlias(t, []string{"--name", "prod-box", "10.0.0.1"}, "prod-box")
	assertJoinAlias(t, []string{"--alias", "prod-box", "10.0.0.1"}, "prod-box")
	assertJoinAlias(t, []string{"-n", "prod-box", "10.0.0.1"}, "prod-box")
}

func TestParseSSHJoinOptions_UserAndPortFlags(t *testing.T) {
	args := []string{"-u", "admin", "-p", "2200", "10.0.0.1"}
	opts, err := ParseSSHJoinOptions(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Target.Username != "admin" || opts.Target.Port != 2200 {
		t.Errorf("expected admin and 2200, got %s and %d", opts.Target.Username, opts.Target.Port)
	}
}

func TestParseSSHJoinOptions_FlagOverrides(t *testing.T) {
	args := []string{"-u", "newuser", "-p", "9999", "origuser@10.0.0.1:22"}
	opts, err := ParseSSHJoinOptions(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Target.Username != "newuser" || opts.Target.Port != 9999 {
		t.Errorf("expected newuser and 9999, got %s and %d", opts.Target.Username, opts.Target.Port)
	}
}

func TestParseSSHJoinOptions_Auth(t *testing.T) {
	opts, err := ParseSSHJoinOptions([]string{"--auth", "10.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.IsPushAuth {
		t.Errorf("expected IsPushAuth true")
	}
}

func TestParseSSHJoinOptions_Help(t *testing.T) {
	opts, err := ParseSSHJoinOptions([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.IsShowHelp {
		t.Errorf("expected IsShowHelp true")
	}
}

func TestParseSSHJoinOptions_Invalid(t *testing.T) {
	_, err := ParseSSHJoinOptions([]string{})
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}

	_, err = ParseSSHJoinOptions([]string{"--port", "invalid", "10.0.0.1"})
	if err == nil {
		t.Fatal("expected error for invalid port flag, got nil")
	}
}
