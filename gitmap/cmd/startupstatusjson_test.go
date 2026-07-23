package cmd

// Contract tests for the shared JSON status emitter used by
// `startup-add` and `startup-remove`. Pins the field set, key
// order, and per-status owner/action labels so a future refactor
// that reorders or renames keys breaks here loudly rather than
// silently drifting downstream jq pipelines.

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// expectedStartupStatusKeys is the wire-order list every emitted
// object must follow. Pinned here (not duplicated in each test)
// so a future field add/remove is one diff to review.
var expectedStartupStatusKeys = []string{
	"command", "action", "name", "target", "owner", "force_used", "dry_run",
}

// jsonKeyOrder re-parses one rendered object and returns the keys
// in source-text order. encoding/json maps would lose the order, so
// we walk the bytes directly using json.Decoder's Token() stream.
func jsonKeyOrder(t *testing.T, body []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(body))
	// Skip leading '['
	if tok, err := dec.Token(); err != nil || tok != json.Delim('[') {
		t.Fatalf("expected leading [, got tok=%v err=%v", tok, err)
	}
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		t.Fatalf("expected leading {, got tok=%v err=%v", tok, err)
	}
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("token: %v", err)
		}
		key, ok := tok.(string)
		if !ok {
			t.Fatalf("expected key string, got %T %v", tok, tok)
		}
		keys = append(keys, key)
		// Skip the value (single token for primitives — every
		// startupStatus value is a string or bool today).
		if _, err := dec.Token(); err != nil {
			t.Fatalf("skip value for %q: %v", key, err)
		}
	}

	return keys
}

// TestStartupStatusJSON_AddCreatedKeyOrder pins the field order for
// the created+force_used path. If a reflection-based encoder is ever
// reintroduced, key order could drift across Go versions and break
// this test — exactly the regression we want to catch.
func TestStartupStatusJSON_AddCreatedKeyOrder(t *testing.T) {
	var buf bytes.Buffer
	s := addResultToStatus("watch", true, startup.AddResult{
		Status: startup.AddCreated, Path: "/etc/xdg/autostart/gitmap-watch.desktop",
	})
	if err := writeStartupStatusJSON(&buf, s, 2); err != nil {
		t.Fatalf("emit: %v", err)
	}
	got := jsonKeyOrder(t, buf.Bytes())
	if len(got) != len(expectedStartupStatusKeys) {
		t.Fatalf("got %d keys, want %d: %v", len(got), len(expectedStartupStatusKeys), got)
	}
	for i, want := range expectedStartupStatusKeys {
		if got[i] != want {
			t.Errorf("key %d = %q, want %q (full: %v)", i, got[i], want, got)
		}
	}
}

// TestStartupStatusJSON_AddActionsAndOwners table-tests every
// AddStatus → (action, owner, target?) mapping. Centralized so a
// future AddStatus value is one row to add and an obvious failure
// when forgotten.
func TestStartupStatusJSON_AddActionsAndOwners(t *testing.T) {
	cases := []struct {
		status     startup.AddStatus
		path       string
		wantAction string
		wantOwner  string
		wantTarget string
	}{
		{startup.AddCreated, "/p", "created", "gitmap", "/p"},
		{startup.AddOverwritten, "/p", "overwritten", "gitmap", "/p"},
		{startup.AddExists, "/p", "exists", "gitmap", "/p"},
		{startup.AddRefused, "/p", "refused", "third-party", "/p"},
		{startup.AddBadName, "/ignored", "bad_name", "unknown", ""},
	}
	for _, tc := range cases {
		t.Run(tc.wantAction, func(t *testing.T) {
			s := addResultToStatus("watch", false, startup.AddResult{
				Status: tc.status, Path: tc.path,
			})
			if s.action != tc.wantAction || s.owner != tc.wantOwner || s.target != tc.wantTarget {
				t.Errorf("got action=%q owner=%q target=%q; want %s/%s/%s",
					s.action, s.owner, s.target,
					tc.wantAction, tc.wantOwner, tc.wantTarget)
			}
		})
	}
}

// TestStartupStatusJSON_RemoveActionsAndOwners parallels the Add
// table for RemoveStatus. The "deleted" / "noop" / "refused" /
// "bad_name" matrix proves every owner inference is correct.
func TestStartupStatusJSON_RemoveActionsAndOwners(t *testing.T) {
	cases := []struct {
		status     startup.RemoveStatus
		path       string
		dryRun     bool
		wantAction string
		wantOwner  string
		wantTarget string
		wantDryRun bool
	}{
		{startup.RemoveDeleted, "/p", false, "deleted", "gitmap", "/p", false},
		{startup.RemoveDeleted, "/p", true, "deleted", "gitmap", "/p", true},
		{startup.RemoveNoOp, "", false, "noop", "none", "", false},
		{startup.RemoveRefused, "/p", false, "refused", "third-party", "/p", false},
		{startup.RemoveBadName, "", false, "bad_name", "unknown", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.wantAction, func(t *testing.T) {
			s := removeResultToStatus("watch", startup.RemoveResult{
				Status: tc.status, Path: tc.path, DryRun: tc.dryRun,
			})
			if s.action != tc.wantAction || s.owner != tc.wantOwner ||
				s.target != tc.wantTarget || s.dryRun != tc.wantDryRun {
				t.Errorf("got action=%q owner=%q target=%q dryRun=%v; want %s/%s/%s/%v",
					s.action, s.owner, s.target, s.dryRun,
					tc.wantAction, tc.wantOwner, tc.wantTarget, tc.wantDryRun)
			}
		})
	}
}

// TestStartupStatusJSON_MinifiedShape confirms jsonIndent=0
// produces the compact `[{...}]\n` shape expected by line-oriented
// log pipelines (one record per line when callers chain commands).
func TestStartupStatusJSON_MinifiedShape(t *testing.T) {
	var buf bytes.Buffer
	s := removeResultToStatus("watch", startup.RemoveResult{
		Status: startup.RemoveDeleted, Path: "/p",
	})
	if err := writeStartupStatusJSON(&buf, s, 0); err != nil {
		t.Fatalf("emit: %v", err)
	}
	got := buf.String()
	if !strings.HasPrefix(got, `[{"command":"startup-remove"`) {
		t.Errorf("missing minified array+object opener:\n%s", got)
	}
	if !strings.HasSuffix(got, "}]\n") {
		t.Errorf("missing minified suffix:\n%s", got)
	}
	if strings.Contains(got, "\n  ") {
		t.Errorf("minified output should not contain indented lines:\n%s", got)
	}
}

// TestValidateStartupOutput_RejectsUnknownAndOutOfRange covers the
// shared validator both runners depend on. Empty / typical / bad
// values prove the errors fire the same way for either command name.
func TestValidateStartupOutput_RejectsUnknownAndOutOfRange(t *testing.T) {
	cases := []struct {
		name    string
		output  string
		indent  int
		wantErr bool
	}{
		{"terminal-default", "terminal", 2, false},
		{"json-default", "json", 2, false},
		{"json-min", "json", 0, false},
		{"json-max", "json", 8, false},
		{"unknown-output", "yaml", 2, true},
		{"indent-too-large", "json", 9, true},
		{"indent-negative", "json", -1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateStartupOutput("startup-add", tc.output, tc.indent)
			if (err != nil) != tc.wantErr {
				t.Errorf("output=%q indent=%d: wantErr=%v got %v",
					tc.output, tc.indent, tc.wantErr, err)
			}
		})
	}
}
