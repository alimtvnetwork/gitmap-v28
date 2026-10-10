// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"encoding/json"
	"testing"
)

func TestParseNewCommandsOptionsDefaults(t *testing.T) {
	opts, err := parseNewCommandsOptions([]string{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if opts.Limit != 100 {
		t.Errorf("expected default limit 100, got %d", opts.Limit)
	}
	if opts.IsJSON {
		t.Errorf("expected IsJSON to be false")
	}
	if opts.Category != "" {
		t.Errorf("expected empty category, got %s", opts.Category)
	}
}

func TestParseNewCommandsOptionsCustom(t *testing.T) {
	args := []string{"--limit", "25", "--category", "commits", "--filter", "cpar", "--json"}
	opts, err := parseNewCommandsOptions(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertCustomOptions(t, opts)
}

func assertCustomOptions(t *testing.T, opts NewCommandsOptions) {
	if opts.Limit != 25 {
		t.Errorf("expected limit 25, got %d", opts.Limit)
	}
	if opts.Category != "commits" {
		t.Errorf("expected category commits, got %s", opts.Category)
	}
	if opts.Filter != "cpar" {
		t.Errorf("expected filter cpar, got %s", opts.Filter)
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON to be true")
	}
}

func TestParseNewCommandsOptionsErrors(t *testing.T) {
	testInvalidLimit(t)
	testMissingFlagValues(t)
}

func testInvalidLimit(t *testing.T) {
	_, err := parseNewCommandsOptions([]string{"--limit", "invalid"})
	if err == nil {
		t.Errorf("expected validation error for invalid limit, got nil")
	}
}

func testMissingFlagValues(t *testing.T) {
	_, errLimit := parseNewCommandsOptions([]string{"--limit"})
	if errLimit == nil {
		t.Errorf("expected error for trailing --limit")
	}
	_, errCat := parseNewCommandsOptions([]string{"--category"})
	if errCat == nil {
		t.Errorf("expected error for trailing --category")
	}
}

func TestBuildNewCommandsCatalog(t *testing.T) {
	catalog := buildNewCommandsCatalog()
	if len(catalog) != 100 {
		t.Fatalf("expected exactly 100 catalog entries, got %d", len(catalog))
	}
	for i, entry := range catalog {
		assertCatalogEntryValid(t, i, entry)
	}
}

func assertCatalogEntryValid(t *testing.T, idx int, entry NewCommandEntry) {
	assertEntryCoreFields(t, idx, entry)
	assertEntryDetails(t, idx, entry)
}

func assertEntryCoreFields(t *testing.T, idx int, entry NewCommandEntry) {
	if entry.Name == "" {
		t.Errorf("entry [%d] has empty Name", idx)
	}
	if entry.Version == "" {
		t.Errorf("entry [%d] (%s) has empty Version", idx, entry.Name)
	}
	if entry.Category == "" {
		t.Errorf("entry [%d] (%s) has empty Category", idx, entry.Name)
	}
}

func assertEntryDetails(t *testing.T, idx int, entry NewCommandEntry) {
	if entry.Description == "" {
		t.Errorf("entry [%d] (%s) has empty Description", idx, entry.Name)
	}
	if entry.Example == "" {
		t.Errorf("entry [%d] (%s) has empty Example", idx, entry.Name)
	}
}

func TestFilterNewCommandsLimit(t *testing.T) {
	catalog := buildNewCommandsCatalog()
	opts := NewCommandsOptions{Limit: 5}
	filtered := filterNewCommands(catalog, opts)
	if len(filtered) != 5 {
		t.Fatalf("expected 5 filtered commands, got %d", len(filtered))
	}
}

func TestFilterNewCommandsCategory(t *testing.T) {
	catalog := buildNewCommandsCatalog()
	opts := NewCommandsOptions{Limit: 100, Category: "commits"}
	filtered := filterNewCommands(catalog, opts)
	if len(filtered) == 0 {
		t.Fatalf("expected matching commands for category 'commits'")
	}
	for _, cmd := range filtered {
		if cmd.Category != "commits" {
			t.Errorf("expected category 'commits', got '%s'", cmd.Category)
		}
	}
}

func TestFilterNewCommandsQuery(t *testing.T) {
	catalog := buildNewCommandsCatalog()
	opts := NewCommandsOptions{Limit: 100, Filter: "cpar"}
	filtered := filterNewCommands(catalog, opts)
	if len(filtered) == 0 {
		t.Fatalf("expected matches for query 'cpar'")
	}
	if filtered[0].Name != "cpar" {
		t.Errorf("expected first match to be cpar, got %s", filtered[0].Name)
	}
}

func TestNewCommandsJSONOutput(t *testing.T) {
	catalog := buildNewCommandsCatalog()
	opts := NewCommandsOptions{Limit: 10, Category: "commits"}
	payload := buildNewCommandsPayload(catalog, opts)
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}
	assertPayloadUnmarshal(t, data)
}

func assertPayloadUnmarshal(t *testing.T, data []byte) {
	var decoded NewCommandsPayload
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if decoded.TotalCommands != 100 {
		t.Errorf("expected TotalCommands 100, got %d", decoded.TotalCommands)
	}
	if decoded.FilteredCommands != len(decoded.Commands) {
		t.Errorf("mismatch in filtered commands count")
	}
}

func TestRunNewCommandsHelpAndExecution(t *testing.T) {
	testHelpExecution(t)
	testJSONExecution(t)
}

func testHelpExecution(t *testing.T) {
	if err := RunNewCommands([]string{"--help"}); err != nil {
		t.Errorf("expected nil error on --help, got %v", err)
	}
	if err := RunNewCommands([]string{"-h"}); err != nil {
		t.Errorf("expected nil error on -h, got %v", err)
	}
}

func testJSONExecution(t *testing.T) {
	if err := RunNewCommands([]string{"--json", "--limit", "2"}); err != nil {
		t.Errorf("expected nil error on valid JSON run, got %v", err)
	}
}
