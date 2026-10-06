package tempdir

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// BackupVaultManifest catalogs all backed up blobs for an operation.
type BackupVaultManifest struct {
	OperationId int64            `json:"operationId"`
	RepoSlug    string           `json:"repoSlug"`
	CreatedAt   int64            `json:"createdAt"`
	TotalBlobs  int              `json:"totalBlobs"`
	TotalBytes  int64            `json:"totalBytes"`
	Blobs       []VaultBlobEntry `json:"blobs"`
}

// VaultBlobEntry describes a single backed-up blob.
type VaultBlobEntry struct {
	CommitSha     string `json:"commitSha"`
	RelativePath  string `json:"relativePath"`
	BackupRelPath string `json:"backupRelPath"`
	Sha256        string `json:"sha256"`
	SizeBytes     int64  `json:"sizeBytes"`
	FileMode      uint32 `json:"fileMode"`
}

// PurgeBlobTarget specifies a targeted blob to stage in temporary backup.
type PurgeBlobTarget struct {
	CommitSha    string `json:"commitSha"`
	RelativePath string `json:"relativePath"`
	BlobSha      string `json:"blobSha"`
	FileMode     uint32 `json:"fileMode"`
	Data         []byte `json:"data,omitempty"`
}

// HistoryBackupVault provides access to the temporary staging filesystem.
type HistoryBackupVault struct {
	baseDir     string
	repoSlug    string
	operationId int64
}

// NewHistoryBackupVault initializes a new vault helper for a repo and operation.
func NewHistoryBackupVault(repoSlug string, operationId int64) *HistoryBackupVault {
	opIdStr := fmt.Sprintf("%d", operationId)
	baseDir := RepoTempDir("history-backup", repoSlug, opIdStr)

	return &HistoryBackupVault{
		baseDir:     baseDir,
		repoSlug:    repoSlug,
		operationId: operationId,
	}
}

// BaseDir returns the root directory for this operation's vault.
func (v *HistoryBackupVault) BaseDir() string {
	return v.baseDir
}

// WriteBlob copies or writes a blob's content to the vault.
func (v *HistoryBackupVault) WriteBlob(commitSha, relPath string, data []byte, mode uint32) (*VaultBlobEntry, error) {
	if err := validateRelativePath(relPath); err != nil {
		return nil, err
	}

	destPath := filepath.Join(v.baseDir, commitSha, filepath.Clean(relPath))
	if err := writeVaultFile(destPath, data, mode); err != nil {
		return nil, err
	}

	entry := buildBlobEntry(commitSha, relPath, data, mode)

	return &entry, nil
}

// ReadBlob reads a previously backed up blob by commit and path.
func (v *HistoryBackupVault) ReadBlob(commitSha, relPath string) ([]byte, uint32, error) {
	if err := validateRelativePath(relPath); err != nil {
		return nil, 0, err
	}

	srcPath := filepath.Join(v.baseDir, commitSha, filepath.Clean(relPath))
	info, err := os.Stat(srcPath)
	if err != nil {
		return nil, 0, apperror.WrapSimple(err, "stat vault blob")
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, 0, apperror.WrapSimple(err, "read vault blob")
	}

	return data, uint32(info.Mode().Perm()), nil
}

// SaveManifest writes manifest.json to the root of the operation vault.
func (v *HistoryBackupVault) SaveManifest(entries []VaultBlobEntry) error {
	manifest := &BackupVaultManifest{
		OperationId: v.operationId,
		RepoSlug:    v.repoSlug,
		CreatedAt:   time.Now().Unix(),
		TotalBlobs:  len(entries),
		TotalBytes:  calculateTotalBytes(entries),
		Blobs:       entries,
	}

	manifestPath := filepath.Join(v.baseDir, "manifest.json")

	return writeManifestJson(manifestPath, manifest)
}

// LoadManifest reads and verifies manifest.json.
func (v *HistoryBackupVault) LoadManifest() (*BackupVaultManifest, error) {
	manifestPath := filepath.Join(v.baseDir, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read vault manifest")
	}

	var manifest BackupVaultManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, apperror.WrapSimple(err, "unmarshal vault manifest")
	}

	return &manifest, nil
}

// PurgeVault safely cleans up the temp backup directory after expiration or manual purge.
func (v *HistoryBackupVault) PurgeVault() error {
	if err := os.RemoveAll(v.baseDir); err != nil {
		return apperror.WrapSimple(err, "purge backup vault")
	}

	return nil
}

// StageBlobsToTempBackup stages all targeted files to the temporary backup vault.
func StageBlobsToTempBackup(repoDir, repoSlug, opId string, targetFiles []PurgeBlobTarget) (string, error) {
	vaultDir := RepoTempDir("history-backup", repoSlug, opId)
	entries, err := processAllBlobTargets(repoDir, vaultDir, targetFiles)
	if err != nil {
		return "", err
	}

	manifest := &BackupVaultManifest{
		RepoSlug:   repoSlug,
		CreatedAt:  time.Now().Unix(),
		TotalBlobs: len(entries),
		TotalBytes: calculateTotalBytes(entries),
		Blobs:      entries,
	}

	manifestPath := filepath.Join(vaultDir, "manifest.json")
	if err := writeManifestJson(manifestPath, manifest); err != nil {
		return "", err
	}

	return vaultDir, nil
}

func processAllBlobTargets(repoDir, vaultDir string, targetFiles []PurgeBlobTarget) ([]VaultBlobEntry, error) {
	var entries []VaultBlobEntry
	for _, target := range targetFiles {
		entry, err := stageSingleTarget(repoDir, vaultDir, target)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func stageSingleTarget(repoDir, vaultDir string, target PurgeBlobTarget) (VaultBlobEntry, error) {
	if err := validateRelativePath(target.RelativePath); err != nil {
		return VaultBlobEntry{}, err
	}

	data, err := resolveBlobData(repoDir, target)
	if err != nil {
		return VaultBlobEntry{}, err
	}

	destPath := filepath.Join(vaultDir, target.CommitSha, filepath.Clean(target.RelativePath))
	if err := writeVaultFile(destPath, data, target.FileMode); err != nil {
		return VaultBlobEntry{}, err
	}

	return buildBlobEntry(target.CommitSha, target.RelativePath, data, target.FileMode), nil
}

func validateRelativePath(relPath string) error {
	cleaned := filepath.Clean(relPath)
	if strings.HasPrefix(cleaned, "..") || filepath.IsAbs(cleaned) {
		return apperror.NewValidation("validateRelativePath", "E_PATH_TRAVERSAL", "path traversal detected: "+relPath)
	}

	return nil
}

func resolveBlobData(repoDir string, target PurgeBlobTarget) ([]byte, error) {
	if len(target.Data) > 0 {
		return target.Data, nil
	}

	if target.BlobSha != "" && repoDir != "" {
		return fetchBlobBySha(repoDir, target.BlobSha)
	}

	if target.CommitSha != "" && repoDir != "" {
		return fetchBlobByCommitPath(repoDir, target.CommitSha, target.RelativePath)
	}

	return nil, apperror.NewValidation("resolveBlobData", "E_EMPTY_BLOB", "no data or git locator provided")
}

func fetchBlobBySha(repoDir, blobSha string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoDir, "cat-file", "-p", blobSha)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "git cat-file blob: "+blobSha)
	}

	return out, nil
}

func fetchBlobByCommitPath(repoDir, commitSha, relPath string) ([]byte, error) {
	spec := fmt.Sprintf("%s:%s", commitSha, filepath.ToSlash(relPath))
	cmd := exec.Command("git", "-C", repoDir, "show", spec)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "git show blob: "+spec)
	}

	return out, nil
}

func buildBlobEntry(commitSha, relPath string, data []byte, mode uint32) VaultBlobEntry {
	hash := sha256.Sum256(data)
	fileMode := mode
	if fileMode == 0 {
		fileMode = 0644
	}

	return VaultBlobEntry{
		CommitSha:     commitSha,
		RelativePath:  filepath.ToSlash(filepath.Clean(relPath)),
		BackupRelPath: filepath.ToSlash(filepath.Join(commitSha, filepath.Clean(relPath))),
		Sha256:        hex.EncodeToString(hash[:]),
		SizeBytes:     int64(len(data)),
		FileMode:      fileMode,
	}
}

func writeVaultFile(destPath string, data []byte, mode uint32) error {
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir vault dir: "+dir)
	}

	perm := os.FileMode(mode)
	if perm == 0 {
		perm = 0600
	}

	if err := os.WriteFile(destPath, data, perm); err != nil {
		return apperror.WrapSimple(err, "write vault file: "+destPath)
	}

	return nil
}

func writeManifestJson(manifestPath string, manifest *BackupVaultManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal vault manifest")
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return apperror.WrapSimple(err, "write manifest json: "+manifestPath)
	}

	return nil
}

func calculateTotalBytes(entries []VaultBlobEntry) int64 {
	var total int64
	for _, e := range entries {
		total += e.SizeBytes
	}

	return total
}
