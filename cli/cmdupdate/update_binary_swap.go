package cmdupdate

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/subtle"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runPureBinaryUpdate aliases runBinaryReleaseUpdate for spec compliance.
func runPureBinaryUpdate(opts UpdateOptions) error { return runBinaryReleaseUpdate(opts) }

// runBinaryReleaseUpdate executes pure pre-compiled binary update from GitHub Releases.
func runBinaryReleaseUpdate(opts UpdateOptions) error {
	ui := NewUpdateUIRenderer(opts, constants.Version)
	defer ui.Close()

	cand, errCand := resolveCandidateWithUI(ui, opts)
	if errCand != nil {
		return errCand
	}
	if shouldSkipUpdate(ui, opts, cand) {
		return nil
	}
	return runBinaryDownloadAndSwap(ui, opts, cand, time.Now())
}

func shouldSkipUpdate(ui UpdateUIRenderer, opts UpdateOptions, cand *ReleaseCandidate) bool {
	cur := constants.Version
	if IsAlreadyUpdated(cur, cand.Version, opts.IsForce) {
		ui.RenderAlreadyUpdated(cur, cand.Version)
		return true
	}
	if opts.IsDryRun {
		outputDryRun(cand, cur, "go_http", resolveCurrentInstallDir(), opts.IsJSON)
		return true
	}
	return false
}

func runBinaryDownloadAndSwap(ui UpdateUIRenderer, opts UpdateOptions, cand *ReleaseCandidate, start time.Time) error {
	stageDir, errDir := os.MkdirTemp("", fmt.Sprintf("gitmap-update-%d-*", time.Now().UnixNano()))
	if errDir != nil {
		ui.OnStageFail(StageDownloading, errDir)
		return apperror.WrapWithDetails(errDir, "cmd.update.stageDir", "E1136", "failed staging dir", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	defer os.RemoveAll(stageDir)

	archive, checksum, errDown := downloadReleaseStaging(cand, stageDir, ui)
	if errDown != nil {
		ui.OnStageFail(StageDownloading, errDown)
		return errDown
	}
	if errVer := verifyReleaseArchive(ui, cand, archive, checksum); errVer != nil {
		return errVer
	}
	if errSwap := swapReleaseBinary(ui, archive, stageDir); errSwap != nil {
		return errSwap
	}
	return finalizeUpdate(ui, opts, cand, start)
}

func resolveCandidateWithUI(ui UpdateUIRenderer, opts UpdateOptions) (*ReleaseCandidate, error) {
	ui.OnStageStart(StageResolving, "Resolving latest release candidate...")
	cand, err := resolveBinaryReleaseCandidate(opts)
	if err != nil {
		ui.OnStageFail(StageResolving, err)
		return nil, err
	}
	ui.OnStageSuccess(StageResolving, fmt.Sprintf("Target release resolved: %s", cand.Tag))
	return cand, nil
}

func verifyReleaseArchive(ui UpdateUIRenderer, cand *ReleaseCandidate, archive, checksum string) error {
	ui.OnStageStart(StageVerifying, "Verifying SHA256 integrity against checksums.txt...")
	name := filepath.Base(cand.AssetURL)
	isMatch, err := verifyBinaryChecksum(archive, checksum, name)
	if err != nil {
		ui.OnStageFail(StageVerifying, err)
		return err
	}
	if !isMatch {
		errMismatch := apperror.NewWithDetails("cmd.update.checksumMismatch", "E1140", "checksum mismatch for "+name, "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
		ui.OnStageFail(StageVerifying, errMismatch)
		return errMismatch
	}
	ui.OnStageSuccess(StageVerifying, "Cryptographic checksum verified against checksums.txt")
	return nil
}

func swapReleaseBinary(ui UpdateUIRenderer, archive, stageDir string) error {
	ui.OnStageStart(StageSwapping, "Applying binary handoff and filesystem swap...")
	extracted, errExtract := extractBinaryFromArchive(archive, stageDir)
	if errExtract != nil {
		ui.OnStageFail(StageSwapping, errExtract)
		return errExtract
	}
	if errHandoff := applyBinaryHandoff(extracted); errHandoff != nil {
		ui.OnStageFail(StageSwapping, errHandoff)
		return errHandoff
	}
	ui.OnStageSuccess(StageSwapping, "Binary handoff complete; cleanup scheduled")
	return nil
}

func finalizeUpdate(ui UpdateUIRenderer, opts UpdateOptions, cand *ReleaseCandidate, start time.Time) error {
	dur := time.Since(start)
	dir := resolveCurrentInstallDir()
	bin := constants.GitMapBin
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	res := UpdateResult{
		Success: true, Status: "updated", PreviousVersion: FormatVersionTag(constants.Version),
		CurrentVersion: FormatVersionTag(cand.Tag), TargetTag: FormatVersionTag(cand.Tag),
		InstallDir: dir, BinaryPath: filepath.Join(dir, bin), DownloaderEngine: "go_http",
		DownloadBytes: cand.AssetSize, DurationMs: dur.Milliseconds(), FallbackDepth: cand.FallbackDepth,
		IsCached: cand.IsCached, ScriptsUpdated: syncLocalRunnerScripts(dir, cand.Tag),
	}
	ui.RenderSummaryCard(res)
	ensurePostUpdateCompletions()
	return nil
}

func resolveBinaryReleaseCandidate(opts UpdateOptions) (*ReleaseCandidate, error) {
	slug, _, errSlug := resolveTargetSlug()
	if errSlug != nil {
		slug = constants.UpdateCurrentRepoSlug
	}
	cand, errCand := ResolveUpdateTargetWithFallback(slug, opts.TargetVersion, opts.MaxFallbackTags, opts.IsForce)
	if errCand != nil {
		return nil, errCand
	}
	if len(cand.AssetURL) == 0 {
		return nil, apperror.NewWithDetails("cmd.update.noAsset", "E1137", "no asset for "+runtime.GOOS+"/"+runtime.GOARCH, "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	if len(cand.ChecksumURL) == 0 {
		cand.ChecksumURL = fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/checksums.txt", constants.UpdateRepoOwner, slug, cand.Tag)
	}
	return cand, nil
}

func downloadReleaseStaging(cand *ReleaseCandidate, stageDir string, ui UpdateUIRenderer) (string, string, error) {
	name := filepath.Base(cand.AssetURL)
	ui.OnStageStart(StageDownloading, fmt.Sprintf("Downloading %s via go_http...", name))
	archPath := filepath.Join(stageDir, name)
	if err := downloadHTTP(cand.AssetURL, archPath, MakeProgressCallback(ui)); err != nil {
		return "", "", apperror.WrapWithDetails(err, "cmd.update.downloadArchive", "E1138", "failed download", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"url": cand.AssetURL})
	}
	chkPath := filepath.Join(stageDir, "checksums.txt")
	if err := downloadHTTP(cand.ChecksumURL, chkPath, nil); err != nil {
		return "", "", apperror.WrapWithDetails(err, "cmd.update.downloadChecksum", "E1139", "failed checksums.txt", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"url": cand.ChecksumURL})
	}
	ui.OnStageSuccess(StageDownloading, fmt.Sprintf("Downloaded %s (%.1f MB)", name, float64(cand.AssetSize)/(1024*1024)))
	return archPath, chkPath, nil
}

func verifyBinaryChecksum(archivePath, checksumPath, expectedAssetName string) (bool, error) {
	hash, errHash := computeFileSHA256(archivePath)
	if errHash != nil {
		return false, apperror.WrapWithDetails(errHash, "cmd.update.hashCompute", "E1144", "hash compute error", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	f, errOpen := os.Open(checksumPath)
	if errOpen != nil {
		return false, apperror.WrapWithDetails(errOpen, "cmd.update.checksumOpen", "E1145", "checksum open error", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	clean := filepath.Base(expectedAssetName)
	for scanner.Scan() {
		fields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(fields) >= 2 && strings.EqualFold(filepath.Base(strings.TrimPrefix(fields[1], "*")), clean) {
			isMatch := subtle.ConstantTimeCompare([]byte(strings.ToLower(hash)), []byte(strings.ToLower(fields[0]))) == 1
			return isMatch, nil
		}
	}
	return false, apperror.NewWithDetails("cmd.update.checksumNotFound", "E1141", "asset not found in checksums", "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func extractBinaryFromArchive(archivePath, extractDir string) (string, error) {
	lower := strings.ToLower(archivePath)
	if strings.HasSuffix(lower, ".zip") {
		return extractZipBinary(archivePath, extractDir)
	}
	if strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz") {
		return extractTarGzBinary(archivePath, extractDir)
	}
	return "", apperror.NewWithDetails("cmd.update.unsupportedArchive", "E1146", "unsupported archive format", "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func extractZipBinary(archivePath, extractDir string) (string, error) {
	zr, errZip := zip.OpenReader(archivePath)
	if errZip != nil {
		return "", apperror.WrapWithDetails(errZip, "cmd.update.zipOpen", "E1147", "failed to open zip", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	defer zr.Close()
	for _, f := range zr.File {
		name := filepath.Base(filepath.Clean(f.Name))
		if !strings.EqualFold(name, "gitmap.exe") && !strings.EqualFold(name, "gitmap") {
			continue
		}
		rc, errOpen := f.Open()
		if errOpen != nil {
			return "", errOpen
		}
		defer rc.Close()
		return writeExtractedBinary(rc, filepath.Join(extractDir, name+".new"))
	}
	return "", apperror.NewSimple("gitmap binary not found in zip", "E1148")
}

func extractTarGzBinary(archivePath, extractDir string) (string, error) {
	f, errOpen := os.Open(archivePath)
	if errOpen != nil {
		return "", errOpen
	}
	defer f.Close()
	gzr, errGz := gzip.NewReader(f)
	if errGz != nil {
		return "", errGz
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		h, err := tr.Next()
		if err != nil {
			break
		}
		name := filepath.Base(filepath.Clean(h.Name))
		if name == "gitmap" || name == "gitmap.exe" {
			return writeExtractedBinary(tr, filepath.Join(extractDir, name+".new"))
		}
	}
	return "", apperror.NewSimple("gitmap binary not found in tar", "E1148")
}

func writeExtractedBinary(r io.Reader, target string) (string, error) {
	out, errCreate := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if errCreate != nil {
		return "", errCreate
	}
	defer out.Close()
	_, errCopy := io.Copy(out, r)
	return target, errCopy
}

func applyBinaryHandoff(newBinaryPath string) error {
	activeExe, err := os.Executable()
	if err != nil {
		return apperror.WrapWithDetails(err, "cmd.update.activeExe", "E1149", "failed active exe", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	if resolved, evalErr := filepath.EvalSymlinks(activeExe); evalErr == nil {
		activeExe = resolved
	}
	if runtime.GOOS == "windows" {
		return applyWindowsHandoff(activeExe, newBinaryPath)
	}
	if errMove := os.Rename(newBinaryPath, activeExe); errMove != nil && copyFile(newBinaryPath, activeExe) != nil {
		return apperror.NewSimple("failed to place new executable", "E1151")
	}
	_ = os.Chmod(activeExe, 0o755)
	return nil
}

func applyWindowsHandoff(activeExe, newBinaryPath string) error {
	backupPath := activeExe + ".old"
	_ = os.Remove(backupPath)
	if errRename := os.Rename(activeExe, backupPath); errRename != nil {
		return apperror.WrapWithDetails(errRename, "cmd.update.renameBackup", "E1150", "failed backup rename", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"active": activeExe, "backup": backupPath})
	}
	if os.Rename(newBinaryPath, activeExe) != nil && copyFile(newBinaryPath, activeExe) != nil {
		_ = os.Rename(backupPath, activeExe)
		return apperror.NewSimple("failed placing new exe, restored backup", "E1151")
	}
	cmd := exec.Command(activeExe, constants.CmdUpdateCleanup)
	cmd.Env = append(os.Environ(), constants.EnvUpdateCleanupDelayMS+"=1500")
	setHiddenProcessAttr(cmd)
	_ = cmd.Start()
	return nil
}
