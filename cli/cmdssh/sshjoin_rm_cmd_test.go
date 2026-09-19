package cmdssh

import (
	"testing"
)

func TestExecuteSJRm(t *testing.T) {
	_ = executeSJRm
}

func TestResolveRmTarget_AllVariations(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{[]string{"--all"}, "all"},
		{[]string{"-a"}, "all"},
		{[]string{"all"}, "all"},
		{[]string{"devbox"}, "devbox"},
		{[]string{"192.168.1.50"}, "192.168.1.50"},
		{[]string{}, ""},
	}

	for _, c := range cases {
		actual := resolveRmTarget(c.args)
		if actual != c.expected {
			t.Errorf("resolveRmTarget(%v) = %q, expected %q", c.args, actual, c.expected)
		}
	}
}

func TestValidateRmTarget_ErrorOnEmpty(t *testing.T) {
	_, err := validateRmTarget([]string{})
	hasErr := err != nil
	if !hasErr {
		t.Errorf("expected error on empty args")
	}
}
