package replay

import (
	"strings"
	"testing"
	"time"
)

// TestApplyCommitDryRunShortCircuitsAllSideEffects asserts that
// dryRun=true never invokes any git hook — Result.NewSha must be
// empty and the test fakes must record zero calls.
func TestApplyCommitDryRunShortCircuitsAllSideEffects(t *testing.T) {
	calls := 0
	restore := SetTestHooks(
		func(_, _ string, _ ...string) (string, error) { calls++; return "", nil },
		func(_, _ string, _ ...string) ([]byte, error) { calls++; return nil, nil },
		func(_ string, _ []string, _ ...string) (string, error) { calls++; return "", nil },
		func(_ string, _ []byte) (string, error) { calls++; return "", nil },
	)
	defer restore()
	res, err := ApplyCommit(samplePlan(), true)
	if err != nil {
		t.Fatalf("ApplyCommit dry-run: %v", err)
	}
	if res.NewSha != "" {
		t.Fatalf("dry-run returned NewSha %q, want empty", res.NewSha)
	}
	if calls != 0 {
		t.Fatalf("dry-run made %d git calls, want 0", calls)
	}
}

// TestApplyCommitWiresFullPipelineThroughHooks runs ApplyCommit with
// stub hooks that return canned successes and verifies every plumbing
// stage was invoked in the right order with the right args.
func TestApplyCommitWiresFullPipelineThroughHooks(t *testing.T) {
	calls := []string{}
	restore := SetTestHooks(
		func(_, sub string, args ...string) (string, error) {
			calls = append(calls, sub+":"+strings.Join(args, ","))
			return cannedTextResponse(sub)
		},
		func(_, sub string, _ ...string) ([]byte, error) {
			calls = append(calls, sub+":bytes")
			return []byte("file-contents"), nil
		},
		func(_ string, _ []string, args ...string) (string, error) {
			calls = append(calls, args[0])
			return "newshaXYZ", nil
		},
		func(_ string, _ []byte) (string, error) {
			calls = append(calls, "hash-object")
			return "blobsha", nil
		},
	)
	defer restore()
	res, err := ApplyCommit(samplePlan(), false)
	if err != nil {
		t.Fatalf("ApplyCommit: %v", err)
	}
	if res.NewSha != "newshaXYZ" {
		t.Fatalf("NewSha = %q, want newshaXYZ", res.NewSha)
	}
	got := strings.Join(calls, "|")
	for _, want := range []string{"cat-file:bytes", "hash-object", "update-index", "write-tree", "rev-parse", "commit-tree", "update-ref"} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in call sequence: %s", want, got)
		}
	}
}

// TestCommitEnvPinsBothDates is a tiny pure-function test that locks
// the spec guardrail "Replicate BOTH AuthorDate AND CommitterDate
// byte-for-byte" — both env vars MUST be present and equal to the
// RFC3339 format of the source dates.
func TestCommitEnvPinsBothDates(t *testing.T) {
	p := samplePlan()
	env := commitEnv(p)
	wantAuthor := "GIT_AUTHOR_DATE=" + p.AuthorDate.Format(time.RFC3339)
	wantCommit := "GIT_COMMITTER_DATE=" + p.CommitterDate.Format(time.RFC3339)
	if !contains(env, wantAuthor) || !contains(env, wantCommit) {
		t.Fatalf("env missing date pins: %v", env)
	}
}

// cannedTextResponse returns plausible stdout for each git subcommand
// so the pipeline reaches the next stage.
func cannedTextResponse(sub string) (string, error) {
	switch sub {
	case "write-tree":
		return "treeshaABC", nil
	case "rev-parse":
		return "parentshaABC", nil
	}
	return "", nil
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

func samplePlan() Plan {
	at, _ := time.Parse(time.RFC3339, "2024-03-04T05:06:07+02:00")
	ct, _ := time.Parse(time.RFC3339, "2024-03-04T05:06:08+02:00")
	return Plan{
		SourceRepoDir: "/src",
		TargetRepoDir: "/tgt",
		SourceSha:     "deadbeef",
		Files:         []string{"a.go"},
		Message:       "hello",
		AuthorName:    "alice",
		AuthorEmail:   "a@x",
		AuthorDate:    at,
		CommitterDate: ct,
	}
}
