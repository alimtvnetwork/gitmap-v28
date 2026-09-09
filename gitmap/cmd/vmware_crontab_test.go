package cmd

import (
	"strings"
	"testing"
)

func TestCleanCrontabOutput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"no crontab for user", "no crontab for a", ""},
		{"no crontab for root with newline", "no crontab for root\n", ""},
		{"empty string", "", ""},
		{"whitespace only", "   \n  ", ""},
		{"valid cron job", "0 5 * * * /backup.sh", "0 5 * * * /backup.sh"},
		{"valid cron job with newlines", "\n0 5 * * * /backup.sh\n\n", "0 5 * * * /backup.sh"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanCrontabOutput(tc.input)
			if got != tc.expected {
				t.Errorf("cleanCrontabOutput(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestBuildUpdatedCrontabEmpty(t *testing.T) {
	rebootLine := "@reboot /usr/bin/vmhgfs-fuse -o allow_other -o auto_unmount .host:/ /mnt/hgfs"
	result := buildUpdatedCrontab("", rebootLine)

	expected := rebootLine + "\n"
	if result != expected {
		t.Errorf("buildUpdatedCrontab(\"\", line) = %q, want %q", result, expected)
	}

	if strings.Contains(result, "no crontab") {
		t.Errorf("buildUpdatedCrontab should not contain error string, got %q", result)
	}
}

func TestBuildUpdatedCrontabExisting(t *testing.T) {
	existing := "30 2 * * * /opt/scripts/nightly-backup.sh"
	rebootLine := "@reboot /usr/bin/vmhgfs-fuse -o allow_other -o auto_unmount .host:/ /mnt/hgfs"
	result := buildUpdatedCrontab(existing, rebootLine)

	expected := existing + "\n" + rebootLine + "\n"
	if result != expected {
		t.Errorf("buildUpdatedCrontab(existing, line) = %q, want %q", result, expected)
	}

	if !strings.HasPrefix(result, existing) {
		t.Errorf("expected original cron job preserved at beginning")
	}

	if !strings.Contains(result, rebootLine) {
		t.Errorf("expected rebootLine appended")
	}
}
