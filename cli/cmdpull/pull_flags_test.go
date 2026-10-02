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
	if opts.parallel <= 0 || opts.onlyAvailable {
		t.Errorf("expected auto parallel > 0 and onlyAvailable=false, got %+v", opts)
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
		{RepoName: "repo-a", Changes: "2 commits"},
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

func TestNormalizePullArgs_TrailingFlags(t *testing.T) {
	opts := parsePullFlags(NormalizePullArgs([]string{"my-repo", "--json"}))
	if opts.slug != "my-repo" || !opts.isJSON {
		t.Fatalf("expected slug=my-repo isJSON=true, got slug=%q isJSON=%v", opts.slug, opts.isJSON)
	}
}

func TestHasSSHFleetFlag(t *testing.T) {
	if !hasSSHFleetFlag([]string{"pa", "--ssh"}) {
		t.Fatal("expected hasSSHFleetFlag=true for --ssh")
	}
	if !hasSSHFleetFlag([]string{"pull-all", "--sh"}) {
		t.Fatal("expected hasSSHFleetFlag=true for --sh")
	}
	if hasSSHFleetFlag([]string{"pa", "--status"}) {
		t.Fatal("expected hasSSHFleetFlag=false without ssh flag")
	}
}

func TestStripSSHFleetFlags(t *testing.T) {
	clean := stripSSHFleetFlags([]string{"pa", "--ssh", "--status"})
	if len(clean) != 2 || clean[0] != "pa" || clean[1] != "--status" {
		t.Fatalf("expected [pa, --status], got %v", clean)
	}
}

func TestParsePullFlags_Probe(t *testing.T) {
	opts := parsePullFlags(NormalizePullArgs([]string{"pa", "--probe"}))
	if !opts.all {
		t.Fatal("expected all=true for pa")
	}
	if !opts.isProbe {
		t.Fatal("expected isProbe=true for --probe")
	}
}

func TestParsePullFlags_ProbeReposAlias(t *testing.T) {
	opts := parsePullFlags(NormalizePullArgs([]string{"pull", "all", "--probe-repos"}))
	if !opts.all {
		t.Fatal("expected all=true for pull all")
	}
	if !opts.isProbe {
		t.Fatal("expected isProbe=true for --probe-repos")
	}
}

func TestParsePullFlags_ProbeWithYes(t *testing.T) {
	opts := parsePullFlags(NormalizePullArgs([]string{"pa", "--probe", "-y"}))
	if !opts.all || !opts.isProbe || !opts.yes {
		t.Fatalf("expected all=true, isProbe=true, yes=true, got %+v", opts)
	}
}

func TestParsePullFlags_ProbeWithStatus(t *testing.T) {
	opts := parsePullFlags(NormalizePullArgs([]string{"pa", "--probe", "--status"}))
	if !opts.all || !opts.isProbe || !opts.showStatus {
		t.Fatalf("expected all=true, isProbe=true, showStatus=true, got %+v", opts)
	}
}

func TestParsePullFlags_ProbeWithJSON(t *testing.T) {
	opts := parsePullFlags(NormalizePullArgs([]string{"pa", "--probe", "--json"}))
	if !opts.all || !opts.isProbe || !opts.isJSON {
		t.Fatalf("expected all=true, isProbe=true, isJSON=true, got %+v", opts)
	}
}

func TestExtractEfficientFlags_Probe(t *testing.T) {
	if !hasProbeFlag([]string{"pae", "--probe"}) {
		t.Fatal("expected hasProbeFlag=true for --probe")
	}
	if !hasProbeFlag([]string{"paet", "--probe-repos"}) {
		t.Fatal("expected hasProbeFlag=true for --probe-repos")
	}
	_, _, _, _, _, rest := extractEfficientFlags([]string{"pae", "--probe", "-y"})
	for _, r := range rest {
		if r == "--probe" {
			t.Fatalf("expected --probe stripped from rest, got %v", rest)
		}
	}
}

func TestParsePullFlags_AutoScale(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--auto-scale"})
	if !opts.isAutoScale || opts.parallel <= 0 {
		t.Fatalf("expected isAutoScale=true parallel>0, got %+v", opts)
	}
}

func TestParsePullFlags_HighPerfAndTurbo(t *testing.T) {
	optsHigh := parsePullFlags([]string{"--all", "--high-perf"})
	if !optsHigh.isHighPerf || optsHigh.parallel <= 0 {
		t.Fatalf("expected isHighPerf=true parallel>0, got %+v", optsHigh)
	}
	optsTurbo := parsePullFlags([]string{"--all", "--turbo"})
	if !optsTurbo.isHighPerf || optsTurbo.parallel <= 0 {
		t.Fatalf("expected isHighPerf=true via turbo, got %+v", optsTurbo)
	}
}

func TestParsePullFlags_LowCPUAndConservative(t *testing.T) {
	optsLow := parsePullFlags([]string{"--all", "--low-cpu"})
	if !optsLow.isLowCPU || optsLow.parallel <= 0 || optsLow.parallel > 2 {
		t.Fatalf("expected isLowCPU=true parallel 1-2, got %+v", optsLow)
	}
	optsCons := parsePullFlags([]string{"--all", "--conservative"})
	if !optsCons.isLowCPU || optsCons.parallel <= 0 || optsCons.parallel > 2 {
		t.Fatalf("expected isLowCPU=true via conservative, got %+v", optsCons)
	}
}

func TestParsePullFlags_ParallelOverridesPreset(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--high-perf", "-p", "7"})
	if opts.parallel != 7 {
		t.Fatalf("expected parallel=7 override, got %d", opts.parallel)
	}
}

func TestParsePullFlags_ConcurrencyFlag(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--concurrency", "5"})
	if opts.parallel != 5 {
		t.Fatalf("expected parallel=5 via --concurrency, got %d", opts.parallel)
	}
}
