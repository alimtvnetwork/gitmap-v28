package cmdpull

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParsePullFlags_NoFlags(t *testing.T) {
	opts := parsePullFlags([]string{"my-repo"})
	if opts.slug != "my-repo" {
		t.Errorf("expected slug=my-repo, got %q", opts.slug)
	}
	if len(opts.group) > 0 || opts.all || opts.verbose {
		t.Error("expected no group/all/verbose")
	}
	if opts.parallel != 0 || opts.onlyAvailable {
		t.Errorf("expected default parallel=0 and onlyAvailable=false, got %+v", opts)
	}
}

func TestParsePullFlags_GroupLong(t *testing.T) {
	opts := parsePullFlags([]string{"--group", "backend"})
	if len(opts.slug) > 0 {
		t.Errorf("expected empty slug, got %q", opts.slug)
	}

	if opts.group != "backend" {
		t.Errorf("expected group=backend, got %q", opts.group)
	}

	if opts.all {
		t.Error("expected all=false")
	}
}

func TestParsePullFlags_GroupShort(t *testing.T) {
	opts := parsePullFlags([]string{"-g", "infra"})
	if opts.group != "infra" {
		t.Errorf("expected group=infra, got %q", opts.group)
	}
}

func TestParsePullFlags_All(t *testing.T) {
	opts := parsePullFlags([]string{"--all"})
	if len(opts.slug) > 0 || len(opts.group) > 0 {
		t.Error("expected empty slug and group")
	}

	if !opts.all {
		t.Error("expected all=true")
	}
}

func TestParsePullFlags_AllWithVerbose(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--verbose"})
	if !opts.all {
		t.Error("expected all=true")
	}

	if !opts.verbose {
		t.Error("expected verbose=true")
	}
}

func TestParsePullFlags_Parallel(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--parallel", "4"})
	if opts.parallel != 4 {
		t.Errorf("expected parallel=4, got %d", opts.parallel)
	}
}

func TestParsePullFlags_OnlyAvailable(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--only-available"})
	if !opts.onlyAvailable {
		t.Error("expected onlyAvailable=true")
	}
}

func TestParsePullFlags_SSH(t *testing.T) {
	useSSH, useHTTPS, rest := ExtractTransportFlags([]string{"--all", "--ssh"})
	opts := resolveParsedPullOptions(rest, useSSH, useHTTPS)
	if !opts.all {
		t.Error("expected all=true")
	}

	if !opts.useSSH {
		t.Error("expected useSSH=true")
	}

	if opts.useHTTPS {
		t.Error("expected useHTTPS=false")
	}
}

func TestParsePullFlags_SSHPositionalAndCaseInsensitive(t *testing.T) {
	useSSH, useHTTPS, rest := ExtractTransportFlags([]string{"--all", "ssh"})
	opts := resolveParsedPullOptions(rest, useSSH, useHTTPS)
	if !opts.useSSH {
		t.Error("expected useSSH=true for positional 'ssh'")
	}

	useSSHUpper, _, _ := ExtractTransportFlags([]string{"--all", "SSH"})
	if !useSSHUpper {
		t.Error("expected useSSH=true for uppercase 'SSH'")
	}

	_, useHTTPSPos, _ := ExtractTransportFlags([]string{"--all", "https"})
	if !useHTTPSPos {
		t.Error("expected useHTTPS=true for positional 'https'")
	}
}

func TestNormalizePullArgs_FastAndTableModes(t *testing.T) {
	gotPa := NormalizePullArgs([]string{"pa"})
	if len(gotPa) != 1 || gotPa[0] != "--all" {
		t.Fatalf("expected [--all] for 'pa', got %v", gotPa)
	}
	gotAllTable := NormalizePullArgs([]string{"all", "table"})
	if len(gotAllTable) != 2 || gotAllTable[0] != "--all" || gotAllTable[1] != "--status" {
		t.Fatalf("expected [--all --status] for 'all table', got %v", gotAllTable)
	}
	gotPat := NormalizePullArgs([]string{"pat"})
	if len(gotPat) != 2 || gotPat[0] != "--all" || gotPat[1] != "--status" {
		t.Fatalf("expected [--all --status] for 'pat', got %v", gotPat)
	}
}

func TestNormalizeAndParsePullFlags_PaStatusAndJSON(t *testing.T) {
	fastOpts := parsePullFlags(NormalizePullArgs([]string{"pa"}))
	if !fastOpts.all || fastOpts.showStatus || fastOpts.isJSON {
		t.Fatalf("expected fast mode all=true showStatus=false isJSON=false, got %+v", fastOpts)
	}
	statusOpts := parsePullFlags(NormalizePullArgs([]string{"pa", "--status"}))
	if !statusOpts.all || !statusOpts.showStatus || statusOpts.isJSON {
		t.Fatalf("expected all=true showStatus=true isJSON=false, got %+v", statusOpts)
	}
	jsonOpts := parsePullFlags(NormalizePullArgs([]string{"pa", "--json"}))
	if !jsonOpts.all || !jsonOpts.isJSON || jsonOpts.showStatus {
		t.Fatalf("expected all=true isJSON=true showStatus=false, got %+v", jsonOpts)
	}
}

func TestBuildPullBatchSummary_JSONStructure(t *testing.T) {
	states := []*PullRepoState{
		{RepoName: "repo-a", Changes: "2 commits", Step: PullStepDone},
	}
	summary := buildPullBatchSummary(1, states, 250*time.Millisecond)
	raw, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if decoded["total"] == nil || decoded["pulledCount"] == nil || decoded["durationMs"] == nil || decoded["states"] == nil {
		t.Fatalf("missing required JSON keys in %s", string(raw))
	}
}
