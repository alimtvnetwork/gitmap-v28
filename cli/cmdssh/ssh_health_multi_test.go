package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestMatchMultipleTargets_RegisteredAndAdhoc(t *testing.T) {
	hosts := []store.SSHHost{
		{Alias: "prod-db", IP: "10.0.0.1", Username: "ubuntu"},
		{Alias: "web-srv", IP: "10.0.0.2", Username: "www"},
	}

	targets := []string{"prod-db", "10.0.0.2", "192.168.1.99"}
	matched := matchMultipleTargets(hosts, targets)

	if len(matched) != 3 {
		t.Fatalf("expected 3 hosts resolved, got %d", len(matched))
	}

	if matched[0].Alias != "prod-db" || matched[0].IP != "10.0.0.1" {
		t.Errorf("unexpected first host: %+v", matched[0])
	}
	if matched[1].Alias != "web-srv" || matched[1].IP != "10.0.0.2" {
		t.Errorf("unexpected second host: %+v", matched[1])
	}
	if matched[2].Alias != "-" || matched[2].IP != "192.168.1.99" {
		t.Errorf("unexpected third (adhoc) host: %+v", matched[2])
	}
}

func TestMatchMultipleTargets_EmptyTargets(t *testing.T) {
	hosts := []store.SSHHost{
		{Alias: "srv1", IP: "10.0.0.1"},
	}
	matched := matchMultipleTargets(hosts, []string{})
	if len(matched) != 0 {
		t.Errorf("expected 0 matches for empty targets list, got %d", len(matched))
	}
}

func TestExtractTargetArg_MultiArgs(t *testing.T) {
	cases := []struct {
		input []string
		want  string
	}{
		{[]string{}, ""},
		{[]string{"devbox"}, "devbox"},
		{[]string{"devbox", "192.168.1.50"}, "devbox,192.168.1.50"},
		{[]string{"srv1", "srv2", "srv3"}, "srv1,srv2,srv3"},
	}

	for _, tc := range cases {
		got := extractTargetArg(tc.input)
		if got != tc.want {
			t.Errorf("extractTargetArg(%v) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

