package clonenext

import (
	"errors"
	"testing"
)

// stubProbe returns a probe function that says "exists" for every
// version up to and including maxExisting, then 404 above.
func stubProbe(maxExisting int) repoExistsFn {
	return func(_, repo string) (bool, error) {
		parsed := ParseRepoName(repo)
		if !parsed.HasVersion {
			return false, nil
		}

		return parsed.CurrentVersion <= maxExisting, nil
	}
}

func TestCheckRemoteForUpdate_NoSuffix_NoUpdate(t *testing.T) {
	parsed := ParseRepoName("plain-repo")
	got, err := checkRemoteForUpdateWith("acme", parsed, 5, stubProbe(99))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got.UpdateNeeded {
		t.Errorf("UpdateNeeded = true, want false (no -v<N> baseline)")
	}
}

func TestCheckRemoteForUpdate_RemoteHigher_TriggersUpdate(t *testing.T) {
	parsed := ParseRepoName("alpha-v3")
	got, err := checkRemoteForUpdateWith("acme", parsed, 10, stubProbe(5))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !got.UpdateNeeded {
		t.Fatalf("UpdateNeeded = false, want true (remote v5 > local v3)")
	}
	if got.RemoteVersion != 5 {
		t.Errorf("RemoteVersion = %d, want 5", got.RemoteVersion)
	}
	if got.TargetRepo != "acme/alpha-v5" {
		t.Errorf("TargetRepo = %q, want acme/alpha-v5", got.TargetRepo)
	}
}

func TestCheckRemoteForUpdate_RemoteSame_NoUpdate(t *testing.T) {
	parsed := ParseRepoName("alpha-v5")
	got, err := checkRemoteForUpdateWith("acme", parsed, 10, stubProbe(5))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got.UpdateNeeded {
		t.Errorf("UpdateNeeded = true, want false (remote == local)")
	}
}

func TestCheckRemoteForUpdate_FailFastOnMiss(t *testing.T) {
	// Probe says v3 and v4 exist, v5 misses → must stop at v4 even
	// if v6 would have hit (we're testing fail-fast semantics).
	calls := 0
	probe := func(_, repo string) (bool, error) {
		calls++
		parsed := ParseRepoName(repo)
		switch parsed.CurrentVersion {
		case 4:
			return true, nil
		case 5:
			return false, nil
		}

		return true, nil
	}
	parsed := ParseRepoName("alpha-v3")
	got, err := checkRemoteForUpdateWith("acme", parsed, 20, probe)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got.RemoteVersion != 4 {
		t.Errorf("RemoteVersion = %d, want 4", got.RemoteVersion)
	}
	if calls != 2 {
		t.Errorf("probe called %d times, want 2 (fail-fast)", calls)
	}
}

func TestCheckRemoteForUpdate_ProbeError_Surfaced(t *testing.T) {
	wantErr := errors.New("network down")
	probe := func(_, _ string) (bool, error) { return false, wantErr }
	parsed := ParseRepoName("alpha-v1")
	_, err := checkRemoteForUpdateWith("acme", parsed, 5, probe)
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want wrapped %v", err, wantErr)
	}
}

func TestCheckRemoteForUpdate_CeilingClamps(t *testing.T) {
	// Probe says everything exists; ceiling=2 means we should stop at v3 (1+2).
	parsed := ParseRepoName("alpha-v1")
	got, err := checkRemoteForUpdateWith("acme", parsed, 2, stubProbe(99))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got.RemoteVersion != 2 {
		t.Errorf("RemoteVersion = %d, want 2 (ceiling cap)", got.RemoteVersion)
	}
}
