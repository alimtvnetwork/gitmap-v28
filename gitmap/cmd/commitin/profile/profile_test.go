package profile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	p := &Profile{
		Name:           "Default",
		SchemaVersion:  1,
		SourceRepoPath: "/abs/path",
		IsDefault:      true,
		ConflictMode:   constants.CommitInConflictModeForceMerge,
		Author:         &Author{Name: "Jane", Email: "j@x.io"},
		Exclusions:     []Exclusion{{Kind: "PathFolder", Value: "node_modules"}},
		MessageRules:   []MessageRule{{Kind: "StartsWith", Value: "Signed-off-by:"}},
		MessagePrefix:  []string{"chore:"},
		WeakWords:      []string{"update"},
		FunctionIntel:  FunctionIntel{IsEnabled: true, Languages: []string{"Go"}},
	}
	out, err := Encode(p)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !strings.Contains(string(out), `"SchemaVersion": 1`) {
		t.Fatalf("encoded missing SchemaVersion: %s", out)
	}
	got, err := Decode(out)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Name != "Default" || got.Author.Email != "j@x.io" {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	raw := []byte(`{"Name":"x","SchemaVersion":1,"Bogus":"nope"}`)
	if _, err := Decode(raw); err == nil {
		t.Fatal("expected unknown-field error")
	}
}

func TestDecodeRejectsBadSchemaVersion(t *testing.T) {
	raw := []byte(`{"Name":"x","SchemaVersion":99}`)
	_, err := Decode(raw)
	if err == nil || !strings.Contains(err.Error(), "SchemaVersion") {
		t.Fatalf("want SchemaVersion error, got %v", err)
	}
}

func TestSaveToDiskRefusesOverwrite(t *testing.T) {
	root := t.TempDir()
	p := &Profile{Name: "P", SchemaVersion: 1}
	if err := SaveToDisk(root, p, false); err != nil {
		t.Fatalf("first save: %v", err)
	}
	if err := SaveToDisk(root, p, false); err == nil {
		t.Fatal("expected refusal on second save")
	}
	if err := SaveToDisk(root, p, true); err != nil {
		t.Fatalf("overwrite save: %v", err)
	}
}

func TestLoadFromDiskMissing(t *testing.T) {
	_, err := LoadFromDisk(t.TempDir(), "Nope")
	var le *LoadError
	if !errors.As(err, &le) || le.Reason != "not found" {
		t.Fatalf("want LoadError not-found, got %v", err)
	}
}

func TestProfilePathLayout(t *testing.T) {
	got := ProfilePath("/root", "Mine")
	want := filepath.Join(string(filepath.Separator)+"root", ".gitmap", "commit-in", "profiles", "Mine.json")
	if got != want {
		t.Fatalf("path mismatch: %s vs %s", got, want)
	}
}

func TestSaveCreatesProfilesDir(t *testing.T) {
	root := t.TempDir()
	p := &Profile{Name: "Auto", SchemaVersion: 1}
	if err := SaveToDisk(root, p, false); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".gitmap", "commit-in", "profiles", "Auto.json")); err != nil {
		t.Fatalf("profile not at expected path: %v", err)
	}
}

func TestResolvePrecedence(t *testing.T) {
	prof := &Profile{ConflictMode: "Prompt", TitlePrefix: "[p]"}
	cliCM := constants.CommitInConflictModeForceMerge
	cli := &CliOverrides{ConflictMode: &cliCM}
	r := Resolve(cli, prof)
	if r.ConflictMode != "ForceMerge" {
		t.Fatalf("CLI should win, got %s", r.ConflictMode)
	}
	if r.TitlePrefix != "[p]" {
		t.Fatalf("profile TitlePrefix should win, got %q", r.TitlePrefix)
	}
	if r.WeakWords[0] != "change" {
		t.Fatalf("default weak words missing, got %v", r.WeakWords)
	}
}
