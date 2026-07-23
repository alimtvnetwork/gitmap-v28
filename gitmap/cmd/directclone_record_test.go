package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

func TestPopulateDirectCloneURLsPreservesSSHTransport(t *testing.T) {
	rec := model.ScanRecord{}
	populateDirectCloneURLs(&rec, "git@github.com:owner/repo.git")
	if rec.Transport != constants.ScanTransportSSH {
		t.Fatalf("Transport = %q, want ssh", rec.Transport)
	}
	if rec.SSHUrl == "" || rec.HTTPSUrl == "" {
		t.Fatalf("urls not both populated: ssh=%q https=%q", rec.SSHUrl, rec.HTTPSUrl)
	}
}

func TestPopulateDirectCloneURLsPreservesHTTPSTransport(t *testing.T) {
	rec := model.ScanRecord{}
	populateDirectCloneURLs(&rec, "https://github.com/owner/repo.git")
	if rec.Transport != constants.ScanTransportHTTPS {
		t.Fatalf("Transport = %q, want https", rec.Transport)
	}
	if rec.SSHUrl == "" || rec.HTTPSUrl == "" {
		t.Fatalf("urls not both populated: ssh=%q https=%q", rec.SSHUrl, rec.HTTPSUrl)
	}
}
