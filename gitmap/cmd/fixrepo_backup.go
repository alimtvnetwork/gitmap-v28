package cmd

// Backup helpers for `gitmap fix-repo` (v5.40.0+).
//
// Layout under the repo root:
//
//   .gitmap/backup/<repo-name>/v<current>/fix-repo/<UTC-timestamp>/
//     manifest.json            — list of backed-up rel paths + meta
//     files/<rel/path>         — verbatim pre-rewrite copies
//
// Each `gitmap fix-repo` invocation that writes >0 files creates one
// timestamped snapshot. `gitmap undo` restores the latest snapshot for
// the current repo + current version (see undo.go).

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// fixRepoBackupManifest is the per-snapshot index file. Read by undo.
type fixRepoBackupManifest struct {
	SchemaVersion int      `json:"schemaVersion"`
	Repo          string   `json:"repo"`
	CurrentV      int      `json:"currentVersion"`
	Timestamp     string   `json:"timestamp"`
	GitmapVersion string   `json:"gitmapVersion"`
	Files         []string `json:"files"`
}

// fixRepoBackupSession holds the resolved snapshot directory + the
// running manifest. Created lazily on first backup so dry-run / no-op
// sweeps never create empty dirs.
type fixRepoBackupSession struct {
	root      string // repo root
	repoName  string // base + version, e.g. gitmap-v27
	version   int
	timestamp string
	dir       string // <root>/.gitmap/backup/<repo>/v<N>/fix-repo/<ts>
	files     []string
	disabled  bool // set if creating the snapshot dir failed
}

// newFixRepoBackupSession constructs a session with the canonical
// UTC timestamp. Filesystem is only touched on the first BackupFile.
func newFixRepoBackupSession(identity fixRepoIdentity) *fixRepoBackupSession {
	repo := identity.base + "-v" + strconv.Itoa(identity.current)
	ts := time.Now().UTC().Format(constants.FixRepoBackupTimestampFmt)
	dir := filepath.Join(identity.root, constants.GitMapDir,
		constants.FixRepoBackupSubdir, repo,
		"v"+strconv.Itoa(identity.current),
		constants.CmdFixRepo, ts)

	return &fixRepoBackupSession{
		root: identity.root, repoName: repo, version: identity.current,
		timestamp: ts, dir: dir,
	}
}

// BackupFile copies `<root>/<rel>` to the snapshot at
// `<dir>/files/<rel>`. Idempotent per rel — if already backed up,
// returns nil without re-copying so the snapshot always reflects the
// pre-rewrite state (the FIRST observation wins).
func (s *fixRepoBackupSession) BackupFile(rel string) {
	if s.disabled {
		return
	}
	dst := filepath.Join(s.dir, constants.FixRepoBackupFilesSubdir, rel)
	if _, err := os.Stat(dst); err == nil {
		return
	}
	if err := copyFileForBackup(filepath.Join(s.root, rel), dst); err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoBackupErrFmt, rel, err)
		s.disabled = true

		return
	}
	s.files = append(s.files, rel)
}

// Finalize writes the manifest.json. Skipped when no files were
// backed up (so empty snapshots don't litter the backup directory).
func (s *fixRepoBackupSession) Finalize() {
	if s.disabled || len(s.files) == 0 {
		return
	}
	manifest := fixRepoBackupManifest{
		SchemaVersion: constants.FixRepoBackupSchemaVersion,
		Repo:          s.repoName, CurrentV: s.version,
		Timestamp: s.timestamp, GitmapVersion: constants.Version,
		Files: s.files,
	}
	writeFixRepoManifest(s.dir, manifest)
	fmt.Fprintf(os.Stderr, constants.FixRepoBackupMsgFmt,
		len(s.files), filepath.ToSlash(relOrAbs(s.root, s.dir)))
}

// writeFixRepoManifest persists the manifest as pretty JSON.
func writeFixRepoManifest(dir string, m fixRepoBackupManifest) {
	body, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoBackupManifestErrFmt, err)

		return
	}
	path := filepath.Join(dir, constants.FixRepoBackupManifestName)
	if err := os.WriteFile(path, body, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoBackupManifestErrFmt, err)
	}
}

// copyFileForBackup mkdirs the dest tree and streams src → dst.
func copyFileForBackup(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

// relOrAbs returns rel(base, target) when possible, else target.
func relOrAbs(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}

	return rel
}
