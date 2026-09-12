package heavy_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestResolveProviderAndSlug_LocalRemote_ExitsZero(t *testing.T) {
	if url := os.Getenv("GITMAP_TEST_RESOLVE_LOCAL"); url != "" {
		cmd.ResolveProviderAndSlugOrExit(url)
		os.Exit(99)
	}

	urls := []string{
		"file:///tmp/gitmap-fixture.git",
		"/var/tmp/local-bare.git",
		"C:/repos/local-bare.git",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			c := exec.Command(os.Args[0], "-test.run=TestResolveProviderAndSlug_LocalRemote_ExitsZero")
			c.Env = append(os.Environ(), "GITMAP_TEST_RESOLVE_LOCAL="+url)
			out, err := c.CombinedOutput()

			exitCode := 0
			isExitErr := false
			var ee *exec.ExitError
			if err != nil {
				ee, isExitErr = err.(*exec.ExitError)
			}

			if err != nil && isExitErr {
				exitCode = ee.ExitCode()
			}

			if err != nil && !isExitErr {
				t.Fatalf("exec failed: %v\noutput:\n%s", err, out)
			}

			if exitCode == constants.ExitVisBadProvider {
				t.Fatalf("local remote %q wrongly rejected with ExitVisBadProvider (4)\noutput:\n%s", url, out)
			}

			if exitCode != constants.ExitVisOK {
				t.Fatalf("local remote %q exited with %d, want ExitVisOK (0)\noutput:\n%s", url, exitCode, out)
			}

			if !strings.Contains(string(out), "skipping local remote") {
				t.Errorf("expected local-skip stderr message, got:\n%s", out)
			}
		})
	}
}
