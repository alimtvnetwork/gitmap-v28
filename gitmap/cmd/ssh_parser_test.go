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
	tests := []struct {
		name     string
		raw      string
		expected *SSHTarget
	}{
		{
			name: "ip@user",
			raw:  "192.168.1.9@a",
			expected: &SSHTarget{
				Username: "a",
				IP:       "192.168.1.9",
				Port:     22,
			},
		},
		{
			name: "user@ip",
			raw:  "a@192.168.1.9",
			expected: &SSHTarget{
				Username: "a",
				IP:       "192.168.1.9",
				Port:     22,
			},
		},
		{
			name: "no user",
			raw:  "m1",
			expected: &SSHTarget{
				Username: "root",
				IP:       "m1",
				Port:     22,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := ParseSSHTarget(tt.raw, "root", 22)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if target.Username != tt.expected.Username || target.IP != tt.expected.IP {
				t.Errorf("expected Username=%q IP=%q, got Username=%q IP=%q",
					tt.expected.Username, tt.expected.IP, target.Username, target.IP)
			}
		})
	}
}
