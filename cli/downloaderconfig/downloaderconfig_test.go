package downloaderconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestDefaults_ValidationAndMarshal(t *testing.T) {
	doc := Defaults()
	if err := Validate(doc); err != nil {
		t.Fatalf("expected Defaults() to pass validation, got: %v", err)
	}

	res := Marshal(doc)
	if res.IsFailure() {
		t.Fatalf("expected Marshal to succeed, got: %v", res.Err)
	}

	raw := res.Value

	str := string(raw)
	if !strings.Contains(str, `"PreferredDownloader":`) {
		t.Errorf("expected marshaled JSON to contain PreferredDownloader, got:\n%s", str)
	}

	if !strings.Contains(str, `"FallbackDownloader":`) {
		t.Errorf("expected marshaled JSON to contain FallbackDownloader, got:\n%s", str)
	}

	hash := SeedHash(doc)
	if len(hash) != 64 {
		t.Errorf("expected 64-char hex SHA-256 hash, got: %s", hash)
	}
}

func TestParse_SeedFileRoundTrip(t *testing.T) {
	seedPath := filepath.Join("..", "data", "downloader-config.json")
	data, err := os.ReadFile(seedPath)
	if err != nil {
		t.Skipf("seed file not found at %s: %v", seedPath, err)
	}

	res2 := Parse(data)
	if res2.IsFailure() {
		t.Fatalf("expected Parse to succeed on seed file, got: %v", res2.Err)
	}

	doc := res2.Value

	if doc.DownloaderConfig.PreferredDownloader != "Aria2C" {
		t.Errorf("expected PreferredDownloader Aria2C, got: %s", doc.DownloaderConfig.PreferredDownloader)
	}

	if doc.DownloaderConfig.ParallelDownloads != 15 {
		t.Errorf("expected ParallelDownloads 15, got: %d", doc.DownloaderConfig.ParallelDownloads)
	}

	if doc.DatabaseVersion.LastKnownVersion != constants.Version {
		t.Errorf("expected LastKnownVersion to resolve to %s, got: %s", constants.Version, doc.DatabaseVersion.LastKnownVersion)
	}
}

func TestValidate_Errors(t *testing.T) {
	badDoc := Defaults()
	badDoc.DownloaderConfig.PreferredDownloader = ""
	if err := Validate(badDoc); err == nil {
		t.Errorf("expected error for empty PreferredDownloader")
	}

	badDoc = Defaults()
	badDoc.DownloaderConfig.ParallelDownloads = 0
	if err := Validate(badDoc); err == nil {
		t.Errorf("expected error for ParallelDownloads = 0")
	}

	badDoc = Defaults()
	badDoc.DownloaderConfig.ParallelDownloads = 100
	if err := Validate(badDoc); err == nil {
		t.Errorf("expected error for ParallelDownloads = 100")
	}
}
