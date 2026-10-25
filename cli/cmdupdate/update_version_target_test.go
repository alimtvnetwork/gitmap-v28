package cmdupdate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestCompareSemVer(t *testing.T) {
	runSemVerEqualityTests(t)
	runSemVerOrderingTests(t)
	runSemVerBoundaryTests(t)
}

func runSemVerEqualityTests(t *testing.T) {
	assertCompare(t, "v1.2.3", "v1.2.3", 0)
	assertCompare(t, "1.2.3", "v1.2.3", 0)
	assertCompare(t, "v1.2.3-beta", "v1.2.3", 0)
}

func runSemVerOrderingTests(t *testing.T) {
	assertCompare(t, "v6.613.0", "v6.610.0", 1)
	assertCompare(t, "v6.610.0", "v6.613.0", -1)
	assertCompare(t, "v2.0.0", "v1.99.99", 1)
	assertCompare(t, "v1.99.99", "v2.0.0", -1)
	assertCompare(t, "v1.0.2", "v1.0.1", 1)
	assertCompare(t, "v1.0.1", "v1.0.2", -1)
}

func runSemVerBoundaryTests(t *testing.T) {
	assertCompare(t, "", "", 0)
	assertCompare(t, "", "v1.0.0", -1)
	assertCompare(t, "v1.0.0", "", 1)
}

func assertCompare(t *testing.T, v1, v2 string, want int) {
	got := CompareSemVer(v1, v2)
	if got != want {
		t.Errorf("CompareSemVer(%q, %q) = %d; want %d", v1, v2, got, want)
	}
}

func TestIsNewerVersion(t *testing.T) {
	if !IsNewerVersion("v6.613.0", "v6.610.0") {
		t.Error("expected 6.613.0 to be newer than 6.610.0")
	}

	if IsNewerVersion("v6.610.0", "v6.613.0") {
		t.Error("expected 6.610.0 not to be newer than 6.613.0")
	}

	if IsNewerVersion("v6.610.0", "v6.610.0") {
		t.Error("expected same version not to be newer")
	}
}

func TestIsAlreadyUpdated_DowngradePrevention(t *testing.T) {
	if !IsAlreadyUpdated("v6.613.0", "v6.610.0", false) {
		t.Error("expected true when target is older and force is false")
	}

	if !IsAlreadyUpdated("v6.613.0", "v6.613.0", false) {
		t.Error("expected true when versions are equal")
	}

	if IsAlreadyUpdated("v6.613.0", "v6.610.0", true) {
		t.Error("expected false when force is true")
	}

	if IsAlreadyUpdated("v6.610.0", "v6.613.0", false) {
		t.Error("expected false when target is newer")
	}
}

func TestIsUpdateAllowed(t *testing.T) {
	if IsUpdateAllowed("v6.613.0", "v6.610.0", false, false) {
		t.Error("expected update to older version to be rejected")
	}

	if !IsUpdateAllowed("v6.613.0", "v6.610.0", false, true) {
		t.Error("expected update to older version with allow-downgrade to be allowed")
	}

	if !IsUpdateAllowed("v6.613.0", "v6.610.0", true, false) {
		t.Error("expected update with force to be allowed")
	}

	if !IsUpdateAllowed("v6.610.0", "v6.613.0", false, false) {
		t.Error("expected update to newer version to be allowed")
	}
}

func TestParseUpdateOptions_AllowDowngrade(t *testing.T) {
	optsLong := ParseUpdateOptions([]string{"--allow-downgrade", "--version", "v6.610.0"})
	if !optsLong.IsAllowDowngrade {
		t.Error("expected IsAllowDowngrade = true with --allow-downgrade")
	}

	optsShort := ParseUpdateOptions([]string{"-d", "-v", "v6.610.0"})
	if !optsShort.IsAllowDowngrade {
		t.Error("expected IsAllowDowngrade = true with -d")
	}

	optsDefault := ParseUpdateOptions([]string{"--version", "v6.610.0"})
	if optsDefault.IsAllowDowngrade {
		t.Error("expected IsAllowDowngrade = false by default")
	}
}

func createOlderSiblingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "gitmap-v29") && strings.HasSuffix(req.URL.Path, "/releases/latest") {
			_ = json.NewEncoder(w).Encode(map[string]string{"tag_name": "v1.0.0"})
			return
		}

		if strings.HasSuffix(req.URL.Path, "/gitmap-v29") || (strings.Contains(req.URL.Path, constants.UpdateCurrentRepoSlug) && strings.HasSuffix(req.URL.Path, "/releases/latest")) {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestResolveLatestRepoSlug_SiblingOlderFallsBack(t *testing.T) {
	srv := createOlderSiblingServer()
	defer srv.Close()

	slug, source, err := resolveLatestRepoSlug(newTestClient(srv))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if slug != constants.UpdateCurrentRepoSlug || source != constants.UpdateProbeSourceRelease {
		t.Fatalf("expected fallback to %s (%s), got %s (%s)", constants.UpdateCurrentRepoSlug, constants.UpdateProbeSourceRelease, slug, source)
	}
}

func createNewerSiblingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, "gitmap-v29") && strings.HasSuffix(req.URL.Path, "/releases/latest") {
			_ = json.NewEncoder(w).Encode(map[string]string{"tag_name": "v99.0.0"})
			return
		}

		if strings.HasSuffix(req.URL.Path, "/gitmap-v29") {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestResolveLatestRepoSlug_SiblingNewerAccepted(t *testing.T) {
	srv := createNewerSiblingServer()
	defer srv.Close()

	slug, source, err := resolveLatestRepoSlug(newTestClient(srv))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if slug != "gitmap-v29" || source != constants.UpdateProbeSourceSibling {
		t.Fatalf("expected sibling gitmap-v29, got %s (%s)", slug, source)
	}
}
