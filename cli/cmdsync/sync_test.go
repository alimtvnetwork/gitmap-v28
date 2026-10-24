package cmdsync

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/enums"
)

func TestSyncOptionsEnums(t *testing.T) {
	optsDry := SyncOptions{DryRun: true}
	if optsDry.Mode() != enums.SyncModeDryRun {
		t.Errorf("expected SyncModeDryRun, got %s", optsDry.Mode())
	}

	optsLive := SyncOptions{DryRun: false}
	if optsLive.Mode() != enums.SyncModeFastForward {
		t.Errorf("expected SyncModeFastForward, got %s", optsLive.Mode())
	}

	optsAll := SyncOptions{Repos: nil}
	if optsAll.Target() != enums.SyncTargetAll {
		t.Errorf("expected SyncTargetAll, got %s", optsAll.Target())
	}

	optsSpecific := SyncOptions{Repos: []string{"gitmap"}}
	if optsSpecific.Target() != enums.SyncTargetDirty {
		t.Errorf("expected SyncTargetDirty, got %s", optsSpecific.Target())
	}
}

func TestRepoSyncResultOutcomeStatus(t *testing.T) {
	resOk := RepoSyncResult{Status: constants.StatusOk}
	if resOk.OutcomeStatus() != enums.SyncRepoStatusSynced {
		t.Errorf("expected SyncRepoStatusSynced, got %s", resOk.OutcomeStatus())
	}

	resFail := RepoSyncResult{Status: constants.StatusFail}
	if resFail.OutcomeStatus() != enums.SyncRepoStatusError {
		t.Errorf("expected SyncRepoStatusError, got %s", resFail.OutcomeStatus())
	}

	resSkipped := RepoSyncResult{Status: constants.StatusSkipped}
	if resSkipped.OutcomeStatus() != enums.SyncRepoStatusSkipped {
		t.Errorf("expected SyncRepoStatusSkipped, got %s", resSkipped.OutcomeStatus())
	}
}

func TestBumpPatchTag(t *testing.T) {
	if got := bumpPatchTag("v1.2.3"); got != "v1.2.4" {
		t.Errorf("expected v1.2.4, got %s", got)
	}

	if got := bumpPatchTag("1.0.0"); got != "1.0.1" {
		t.Errorf("expected 1.0.1, got %s", got)
	}

	if got := bumpPatchTag("invalid"); got != constants.DefaultFallbackTag {
		t.Errorf("expected %s, got %s", constants.DefaultFallbackTag, got)
	}
}
