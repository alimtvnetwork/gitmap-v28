package cmdupdate

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestEligibleForSegmentedDownload_SmallSize(t *testing.T) {
	os.Unsetenv(constants.EnvDownloadSegments)

	isEligible := isEligibleForSegmentedDownload(4<<20, "https://example.com/asset.zip")
	if isEligible {
		t.Error("expected 4 MiB asset not to activate segmented download without env flag")
	}
}

func TestEligibleForSegmentedDownload_DisabledEnv(t *testing.T) {
	os.Setenv(constants.EnvDownloadSegments, "0")
	defer os.Unsetenv(constants.EnvDownloadSegments)

	isEligible := isEligibleForSegmentedDownload(10<<20, "https://example.com/asset.zip")
	if isEligible {
		t.Error("expected segmented download to be disabled when env is set to 0")
	}
}
