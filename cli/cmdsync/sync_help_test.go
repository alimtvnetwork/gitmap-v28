package cmdsync

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
)

func TestIsSyncHelp(t *testing.T) {
	if !isSyncHelp(constants.FlagHelp) {
		t.Errorf("expected isSyncHelp(\"help\") to be true")
	}

	if !isSyncHelp(constants.FlagShortHelp) {
		t.Errorf("expected isSyncHelp(\"-h\") to be true")
	}

	if !isSyncHelp(constants.FlagLongHelp) {
		t.Errorf("expected isSyncHelp(\"--help\") to be true")
	}

	if isSyncHelp(constants.TargetIgnore) {
		t.Errorf("expected isSyncHelp(\"ignore\") to be false")
	}
}
