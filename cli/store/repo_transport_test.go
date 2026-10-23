package store

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func seedTestRepo(t *testing.T, db *DB, rec model.ScanRecord) {
	t.Helper()
	if err := db.UpsertRepos([]model.ScanRecord{rec}); err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}
}

func fetchFirstRepo(t *testing.T, db *DB, pathStr string) model.ScanRecord {
	t.Helper()
	recs, err := db.FindByPath(pathStr)
	if err != nil || len(recs) == 0 {
		t.Fatalf("FindByPath(%q) failed: err=%v len=%d", pathStr, err, len(recs))
	}

	return recs[0]
}

func assertRepoTransport(t *testing.T, rec model.ScanRecord, wantURL, wantTransport string) {
	t.Helper()
	if rec.SSHUrl != wantURL {
		t.Errorf("expected SSHUrl %q, got %q", wantURL, rec.SSHUrl)
	}

	if rec.IdentifiedTransport != wantTransport {
		t.Errorf("expected IdentifiedTransport %q, got %q", wantTransport, rec.IdentifiedTransport)
	}
}

func assertRepoMoved(t *testing.T, db *DB, oldPath, newPath string) {
	t.Helper()
	oldRecs, err := db.FindByPath(oldPath)
	if err != nil || len(oldRecs) != 0 {
		t.Errorf("expected 0 records at old path, got %d (err: %v)", len(oldRecs), err)
	}

	newRecs, err := db.FindByPath(newPath)
	if err != nil || len(newRecs) != 1 {
		t.Errorf("expected 1 record at new path, got %d (err: %v)", len(newRecs), err)
	}
}

func TestUpdateRepoTransport_Success(t *testing.T) {
	db := openTempDB(t)
	absPath := "/tmp/test/repo1"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "repo1", RepoName: "repo1", AbsolutePath: absPath, Branch: "main",
	})

	sshURL := "git@github.com:org/repo1.git"
	if err := db.UpdateRepoTransport(absPath, sshURL, "ssh"); err != nil {
		t.Fatalf("UpdateRepoTransport failed: %v", err)
	}

	rec := fetchFirstRepo(t, db, absPath)
	assertRepoTransport(t, rec, sshURL, "ssh")
}

func TestUpdateRepoTransport_DefaultSSH(t *testing.T) {
	db := openTempDB(t)
	absPath := "/tmp/test/repo-default"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "repo-default", RepoName: "repo-default", AbsolutePath: absPath, Branch: "main",
	})

	sshURL := "git@github.com:org/repo-default.git"
	if err := db.UpdateRepoTransport(absPath, sshURL, ""); err != nil {
		t.Fatalf("UpdateRepoTransport with empty transport failed: %v", err)
	}

	rec := fetchFirstRepo(t, db, absPath)
	assertRepoTransport(t, rec, sshURL, constants.TransportSSH)
}

func TestUpdateRepoTransport_PathNormalization(t *testing.T) {
	db := openTempDB(t)
	normPath := "/work/norm-repo"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "norm-repo", RepoName: "norm-repo", AbsolutePath: normPath, Branch: "main",
	})

	sshURL := "git@github.com:org/norm-repo.git"
	backPath := strings.ReplaceAll(normPath, "/", "\\")
	if err := db.UpdateRepoTransport(backPath, sshURL, "ssh"); err != nil {
		t.Fatalf("UpdateRepoTransport with backslashes failed: %v", err)
	}

	rec := fetchFirstRepo(t, db, normPath)
	assertRepoTransport(t, rec, sshURL, "ssh")
}

func TestUpdateRepoAbsolutePath_Success(t *testing.T) {
	db := openTempDB(t)
	oldPath, newPath := "/tmp/test/old-path", "/tmp/test/new-path"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "move-repo", RepoName: "move-repo", AbsolutePath: oldPath, Branch: "main",
	})

	if err := db.UpdateRepoAbsolutePath(oldPath, newPath); err != nil {
		t.Fatalf("UpdateRepoAbsolutePath failed: %v", err)
	}

	assertRepoMoved(t, db, oldPath, newPath)
}

func TestUpdateRepoAbsolutePath_NoOp(t *testing.T) {
	db := openTempDB(t)
	if err := db.UpdateRepoAbsolutePath("", ""); err != nil {
		t.Errorf("expected nil for empty paths, got %v", err)
	}

	if err := db.UpdateRepoAbsolutePath("/same/path", "/same/path"); err != nil {
		t.Errorf("expected nil for identical paths, got %v", err)
	}
}

func TestGetRepoPrimaryURLByPath_EmptyAndMissing(t *testing.T) {
	db := openTempDB(t)
	id, url, err := db.GetRepoPrimaryURLByPath("")
	if err != nil || id != 0 || url != "" {
		t.Errorf("empty path: expected (0, '', nil), got (%d, %q, %v)", id, url, err)
	}

	id, url, err = db.GetRepoPrimaryURLByPath("/nonexistent/repo")
	if err == nil || id != 0 || url != "" {
		t.Errorf("missing path: expected non-nil error, got (%d, %q, %v)", id, url, err)
	}
}

func TestGetRepoPrimaryURLByPath_HttpsOnly(t *testing.T) {
	db := openTempDB(t)
	absPath := "/tmp/test/https-only"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "https-only", RepoName: "https-only", HTTPSUrl: "https://github.com/org/https-only",
		AbsolutePath: absPath, Branch: "main",
	})

	id, url, err := db.GetRepoPrimaryURLByPath(absPath)
	if err != nil || id == 0 || url != "https://github.com/org/https-only" {
		t.Errorf("expected https url, got (%d, %q, %v)", id, url, err)
	}
}

func TestGetRepoPrimaryURLByPath_SSHOnly(t *testing.T) {
	db := openTempDB(t)
	absPath := "/tmp/test/ssh-only"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "ssh-only", RepoName: "ssh-only", SSHUrl: "git@github.com:org/ssh-only.git",
		AbsolutePath: absPath, Branch: "main",
	})

	id, url, err := db.GetRepoPrimaryURLByPath(absPath)
	if err != nil || id == 0 || url != "git@github.com:org/ssh-only.git" {
		t.Errorf("expected ssh url, got (%d, %q, %v)", id, url, err)
	}
}

func TestGetRepoPrimaryURLByPath_BothPrioritizeSSH(t *testing.T) {
	db := openTempDB(t)
	absPath := "/tmp/test/both-ssh"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "both-ssh", RepoName: "both-ssh", HTTPSUrl: "https://github.com/org/both-ssh",
		SSHUrl: "git@github.com:org/both-ssh.git", IdentifiedTransport: "ssh",
		AbsolutePath: absPath, Branch: "main",
	})

	id, url, err := db.GetRepoPrimaryURLByPath(absPath)
	if err != nil || id == 0 || url != "git@github.com:org/both-ssh.git" {
		t.Errorf("expected ssh url prioritized, got (%d, %q, %v)", id, url, err)
	}
}

func TestGetRepoPrimaryURLByPath_BothPrioritizeHTTPS(t *testing.T) {
	db := openTempDB(t)
	absPath := "/tmp/test/both-https"
	seedTestRepo(t, db, model.ScanRecord{
		Slug: "both-https", RepoName: "both-https", HTTPSUrl: "https://github.com/org/both-https",
		SSHUrl: "git@github.com:org/both-https.git", IdentifiedTransport: "https",
		AbsolutePath: absPath, Branch: "main",
	})

	id, url, err := db.GetRepoPrimaryURLByPath(absPath)
	if err != nil || id == 0 || url != "https://github.com/org/both-https" {
		t.Errorf("expected https url prioritized, got (%d, %q, %v)", id, url, err)
	}
}
