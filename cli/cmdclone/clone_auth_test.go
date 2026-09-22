// Package cmdclone — clone_auth_test.go tests auth store, token injection, and probe logic.
package cmdclone

import (
	"testing"
)

func TestInjectTokenIntoHTTPS(t *testing.T) {
	cases := []struct {
		url      string
		token    string
		expected string
	}{
		{"https://github.com/org/repo.git", "token123", "https://token123@github.com/org/repo.git"},
		{"https://old@github.com/org/repo.git", "token456", "https://token456@github.com/org/repo.git"},
		{"git@github.com:org/repo.git", "token789", "git@github.com:org/repo.git"},
		{"", "token", ""},
	}

	for _, tc := range cases {
		actual := InjectTokenIntoHTTPS(tc.url, tc.token)
		if actual != tc.expected {
			t.Errorf("InjectTokenIntoHTTPS(%q, %q) = %q; want %q", tc.url, tc.token, actual, tc.expected)
		}
	}
}

func TestGlobalAccessTokenStore(t *testing.T) {
	ClearGlobalAccessToken()
	_, isFound := GetGlobalAccessToken()
	if isFound {
		t.Errorf("expected empty token initially")
	}

	SetGlobalAccessToken("ghp_test_token")
	tok, isFound := GetGlobalAccessToken()
	if !isFound || tok != "ghp_test_token" {
		t.Errorf("GetGlobalAccessToken() = (%q, %v); want (ghp_test_token, true)", tok, isFound)
	}

	ClearGlobalAccessToken()
	_, isFoundAfterClear := GetGlobalAccessToken()
	if isFoundAfterClear {
		t.Errorf("expected empty token after clear")
	}
}
