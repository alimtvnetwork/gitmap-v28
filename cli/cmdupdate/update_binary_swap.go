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
func runPureBinaryUpdate(opts UpdateOptions) error {
	return runBinaryReleaseUpdate(opts)
}

// runBinaryReleaseUpdate executes pure pre-compiled binary update from GitHub Releases.
func runBinaryReleaseUpdate(opts UpdateOptions) error {
	ui := NewUpdateUIRenderer(opts, constants.Version)
	defer ui.Close()

	cand, errCand := resolveCandidateWithUI(ui, opts)
	if errCand != nil {
		return errCand
	}

	if isUpdateSkipped(ui, opts, cand) {
		return nil
	}

	return runBinaryDownloadAndSwap(ui, opts, cand, time.Now())
}

func isUpdateSkipped(ui UpdateUIRenderer, opts UpdateOptions, cand *ReleaseCandidate) bool {
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
	stageDir, errDir := createUpdateStageDir(ui)
	if errDir != nil {
		return errDir
	}
	defer os.RemoveAll(stageDir)

	engine, errExec := executeUpdateStaging(ui, cand, stageDir)
	if errExec != nil {
		return errExec
	}

	return finalizeUpdate(ui, opts, cand, start, engine)
}

func createUpdateStageDir(ui UpdateUIRenderer) (string, error) {
	stageDir, errDir := os.MkdirTemp("", fmt.Sprintf("gitmap-update-%d-*", time.Now().UnixNano()))
	if errDir != nil {
		ui.OnStageFail(StageDownloading, errDir)
		return "", apperror.WrapWithDetails(errDir, "cmd.update.stageDir", "E1136", "failed staging dir", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return stageDir, nil
}

func executeUpdateStaging(ui UpdateUIRenderer, cand *ReleaseCandidate, stageDir string) (string, error) {
	archive, engine, errDown := downloadAndVerifyRelease(ui, cand, stageDir)
	if errDown != nil {
		return "", errDown
	}

	if errSwap := executeBinarySwap(ui, archive, stageDir); errSwap != nil {
		return "", errSwap
	}

	return engine, nil
}

func downloadAndVerifyRelease(ui UpdateUIRenderer, cand *ReleaseCandidate, stageDir string) (string, string, error) {
	archive, checksum, engine, errDown := downloadReleaseStaging(cand, stageDir, ui)
	if errDown != nil {
		ui.OnStageFail(StageDownloading, errDown)
		return "", "", errDown
	}

	if errVer := verifyReleaseArchive(ui, cand, archive, checksum); errVer != nil {
		return "", "", errVer
	}

	return archive, engine, nil
}

func executeBinarySwap(ui UpdateUIRenderer, archive, stageDir string) error {
	if isStandaloneBinary(archive) {
		return stageAndSwapDirectBinary(ui, archive, stageDir)
	}

	return swapReleaseBinary(ui, archive, stageDir)
}

func isStandaloneBinary(path string) bool {
	lower := strings.ToLower(path)
	if strings.Contains(lower, "updater") || strings.Contains(lower, "cli-updater") {
		return false
	}
	if isArchiveExtension(lower) {
		return false
	}
	if !strings.Contains(lower, "gitmap") {
		return false
	}

	return strings.HasSuffix(lower, ".exe") || strings.HasSuffix(lower, ".bin")
}

func isArchiveExtension(lower string) bool {
	isTar := strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz")
	isZip := strings.HasSuffix(lower, ".zip")

	return isTar || isZip
}

func resolveTargetBinaryName(assetPath string) string {
	lower := strings.ToLower(assetPath)
	if strings.HasSuffix(lower, ".exe") || runtime.GOOS == "windows" {
		return "gitmap.exe.new"
	}

	return "gitmap.new"
}

func stageAndSwapDirectBinary(ui UpdateUIRenderer, archive, stageDir string) error {
	targetName := resolveTargetBinaryName(archive)
	stagedBinary := filepath.Join(stageDir, targetName)
	if errStage := stageDirectBinary(archive, stagedBinary); errStage != nil {
		ui.OnStageFail(StageSwapping, errStage)
		return errStage
	}

	return swapDirectBinary(ui, stagedBinary)
}

func stageDirectBinary(src, dst string) error {
	if errCopy := copyFile(src, dst); errCopy != nil {
		return apperror.WrapWithDetails(errCopy, "cmd.update.stageDirect", "E1152", "failed staging direct binary", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return os.Chmod(dst, 0o755)
}

func swapDirectBinary(ui UpdateUIRenderer, stagedBinary string) error {
	ui.OnStageStart(StageSwapping, "Applying binary handoff and filesystem swap...")
	if errHandoff := applyBinaryHandoff(stagedBinary); errHandoff != nil {
		ui.OnStageFail(StageSwapping, errHandoff)
		return errHandoff
	}

	ui.OnStageSuccess(StageSwapping, "Binary handoff complete; cleanup scheduled")

	return nil
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
	if err := checkArchiveHash(archive, checksum, name); err != nil {
		ui.OnStageFail(StageVerifying, err)
		return err
	}

	ui.OnStageSuccess(StageVerifying, "Cryptographic checksum verified against checksums.txt")

	return nil
}

func checkArchiveHash(archive, checksum, name string) error {
	isMatch, err := verifyBinaryChecksum(archive, checksum, name)
	if err != nil {
		return err
	}

	if !isMatch {
		return apperror.NewWithDetails("cmd.update.checksumMismatch", "E1140", "checksum mismatch for "+name, "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}

	return nil
}

func swapReleaseBinary(ui UpdateUIRenderer, archive, stageDir string) error {
	ui.OnStageStart(StageSwapping, "Applying binary handoff and filesystem swap...")
	if err := performArchiveHandoff(archive, stageDir); err != nil {
		ui.OnStageFail(StageSwapping, err)
		return err
	}

	ui.OnStageSuccess(StageSwapping, "Binary handoff complete; cleanup scheduled")

	return nil
}

func performArchiveHandoff(archive, stageDir string) error {
	extracted, errExtract := extractBinaryFromArchive(archive, stageDir)
	if errExtract != nil {
		return errExtract
	}

	return applyBinaryHandoff(extracted)
}

func finalizeUpdate(ui UpdateUIRenderer, opts UpdateOptions, cand *ReleaseCandidate, start time.Time, engine string) error {
	if cand == nil {
		return apperror.NewValidationError("release candidate cannot be nil")
	}

	res := buildUpdateResult(cand, time.Since(start), engine)
	ui.RenderSummaryCard(res)
	ensurePostUpdateCompletions()

	return nil
}

func buildUpdateResult(cand *ReleaseCandidate, dur time.Duration, engine string) UpdateResult {
	res := createInitialUpdateResult(cand, dur, engine)
	if cand.ReleaseNotes != "" {
		res.ReleaseNotes = ExtractReleaseNotesSummary(cand.ReleaseNotes, 2)
	}

	return res
}

func createInitialUpdateResult(cand *ReleaseCandidate, dur time.Duration, engine string) UpdateResult {
	dir := resolveCurrentInstallDir()
	bin := resolveInstalledBinaryName()

	return UpdateResult{
		Success: true, Status: "updated", PreviousVersion: FormatVersionTag(constants.Version),
		CurrentVersion: FormatVersionTag(cand.Tag), TargetTag: FormatVersionTag(cand.Tag),
		InstallDir: dir, BinaryPath: filepath.Join(dir, bin), DownloaderEngine: resolveEngineName(engine),
		DownloadBytes: cand.AssetSize, DurationMs: dur.Milliseconds(), FallbackDepth: cand.FallbackDepth,
		IsCached: cand.IsCached, ScriptsUpdated: syncLocalRunnerScripts(dir, cand.Tag),
	}
}

func resolveInstalledBinaryName() string {
	if runtime.GOOS == "windows" {
		return constants.GitMapBin + ".exe"
	}

	return constants.GitMapBin
}

func resolveEngineName(engine string) string {
	if engine == "" {
		return downloadEngineSingle
	}

	return engine
}

func resolveBinaryReleaseCandidate(opts UpdateOptions) (*ReleaseCandidate, error) {
	slug := resolveReleaseSlug()
	cand, errCand := ResolveUpdateTargetWithFallback(slug, opts.TargetVersion, opts.MaxFallbackTags, opts.IsForce)
	if errCand != nil {
		return nil, errCand
	}

	return validateReleaseCandidate(cand, slug)
}

func resolveReleaseSlug() string {
	slug, _, errSlug := resolveTargetSlug()
	if errSlug != nil {
		return constants.UpdateCurrentRepoSlug
	}

	return slug
}

func validateReleaseCandidate(cand *ReleaseCandidate, slug string) (*ReleaseCandidate, error) {
	if len(cand.AssetURL) == 0 {
		msg := "no asset for " + runtime.GOOS + "/" + runtime.GOARCH
		return nil, apperror.NewWithDetails("cmd.update.noAsset", "E1137", msg, "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	if len(cand.ChecksumURL) == 0 {
		cand.ChecksumURL = fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/checksums.txt", constants.UpdateRepoOwner, slug, cand.Tag)
	}

	return cand, nil
}

func downloadReleaseStaging(cand *ReleaseCandidate, stageDir string, ui UpdateUIRenderer) (string, string, string, error) {
	name := filepath.Base(cand.AssetURL)
	ui.OnStageStart(StageDownloading, fmt.Sprintf("Downloading %s...", name))
	archPath := filepath.Join(stageDir, name)
	chkPath := filepath.Join(stageDir, "checksums.txt")

	engine, err := fetchReleaseAssets(cand, archPath, chkPath, ui)
	if err != nil {
		return "", "", "", err
	}

	ui.OnStageSuccess(StageDownloading, fmt.Sprintf("Downloaded %s (%.1f MB) via %s", name, float64(cand.AssetSize)/(1024*1024), engine))

	return archPath, chkPath, engine, nil
}

func fetchReleaseAssets(cand *ReleaseCandidate, archPath, chkPath string, ui UpdateUIRenderer) (string, error) {
	chkErr := startChecksumDownload(cand.ChecksumURL, chkPath)
	engine, errArch := downloadArchiveFast(cand.AssetURL, archPath, cand.AssetSize, MakeProgressCallback(ui))
	if errArch != nil {
		return "", apperror.WrapWithDetails(errArch, "cmd.update.downloadArchive", "E1138", "failed download", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"url": cand.AssetURL})
	}

	if errChk := awaitChecksumDownload(chkErr, cand.ChecksumURL); errChk != nil {
		return "", errChk
	}

	return engine, nil
}

func startChecksumDownload(checksumURL, chkPath string) <-chan error {
	chkErr := make(chan error, 1)
	go func() {
		chkErr <- downloadHTTP(checksumURL, chkPath, nil)
	}()

	return chkErr
}

func awaitChecksumDownload(chkErr <-chan error, checksumURL string) error {
	if errChk := <-chkErr; errChk != nil {
		return apperror.WrapWithDetails(errChk, "cmd.update.downloadChecksum", "E1139", "failed checksums.txt", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"url": checksumURL})
	}

	return nil
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

	return matchChecksumInFile(f, hash, expectedAssetName)
}

func matchChecksumInFile(r io.Reader, hash, expectedAssetName string) (bool, error) {
	clean := filepath.Base(expectedAssetName)
	scanner := bufio.NewScanner(r)
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
	if isStandaloneBinary(lower) {
		return extractDirectExecutable(archivePath, extractDir)
	}

	return extractCompressedArchive(lower, archivePath, extractDir)
}

func extractCompressedArchive(lower, archivePath, extractDir string) (string, error) {
	if strings.HasSuffix(lower, ".zip") {
		return extractZipBinary(archivePath, extractDir)
	}

	if strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz") {
		return extractTarGzBinary(archivePath, extractDir)
	}

	return "", apperror.NewWithDetails("cmd.update.unsupportedArchive", "E1146", "unsupported archive format", "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func extractDirectExecutable(archivePath, extractDir string) (string, error) {
	target := filepath.Join(extractDir, resolveTargetBinaryName(archivePath))
	if errCopy := copyFile(archivePath, target); errCopy != nil {
		return "", apperror.WrapWithDetails(errCopy, "cmd.update.copyDirect", "E1152", "failed direct copy", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	_ = os.Chmod(target, 0o755)

	return target, nil
}

func extractZipBinary(archivePath, extractDir string) (string, error) {
	zr, errZip := zip.OpenReader(archivePath)
	if errZip != nil {
		return "", apperror.WrapWithDetails(errZip, "cmd.update.zipOpen", "E1147", "failed to open zip", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	defer zr.Close()

	for _, f := range zr.File {
		if path, hasFound, errExt := extractZipFileEntry(f, extractDir); hasFound || errExt != nil {
			return path, errExt
		}
	}

	return "", apperror.NewSimple("gitmap binary not found in zip", "E1148")
}

func extractZipFileEntry(f *zip.File, extractDir string) (string, bool, error) {
	name := filepath.Base(filepath.Clean(f.Name))
	if !isZipBinaryTarget(name) {
		return "", false, nil
	}

	path, err := writeZipEntryFile(f, filepath.Join(extractDir, name+".new"))

	return path, true, err
}

func isZipBinaryTarget(name string) bool {
	return strings.EqualFold(name, "gitmap.exe") || strings.EqualFold(name, "gitmap")
}

func writeZipEntryFile(f *zip.File, target string) (string, error) {
	rc, errOpen := f.Open()
	if errOpen != nil {
		return "", errOpen
	}
	defer rc.Close()

	return writeExtractedBinary(rc, target)
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

	return iterateTarEntries(tar.NewReader(gzr), extractDir)
}

func iterateTarEntries(tr *tar.Reader, extractDir string) (string, error) {
	for {
		h, err := tr.Next()
		if err != nil {
			break
		}

		name := filepath.Base(filepath.Clean(h.Name))
		if isReleaseBinaryName(name) {
			return writeExtractedBinary(tr, filepath.Join(extractDir, name+".new"))
		}
	}

	return "", apperror.NewSimple("gitmap binary not found in tar", "E1148")
}

func isReleaseBinaryName(name string) bool {
	return name == "gitmap" || name == "gitmap.exe" || strings.HasPrefix(name, "gitmap-")
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
	activeExe, err := resolveActiveExecutable()
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return applyWindowsHandoff(activeExe, newBinaryPath)
	}

	return applyUnixHandoff(activeExe, newBinaryPath)
}

func resolveActiveExecutable() (string, error) {
	activeExe, err := os.Executable()
	if err != nil {
		return "", apperror.WrapWithDetails(err, "cmd.update.activeExe", "E1149", "failed active exe", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	if resolved, evalErr := filepath.EvalSymlinks(activeExe); evalErr == nil {
		return resolved, nil
	}

	return activeExe, nil
}

func applyUnixHandoff(activeExe, newBinaryPath string) error {
	if errMove := os.Rename(newBinaryPath, activeExe); errMove != nil && copyFile(newBinaryPath, activeExe) != nil {
		return apperror.NewSimple("failed to place new executable", "E1151")
	}

	_ = os.Chmod(activeExe, 0o755)

	return nil
}

func applyWindowsHandoff(activeExe, newBinaryPath string) error {
	backupPath := activeExe + ".old"
	if errBackup := backupActiveBinary(activeExe, backupPath); errBackup != nil {
		return errBackup
	}

	if errPlace := placeNewBinary(newBinaryPath, activeExe, backupPath); errPlace != nil {
		return errPlace
	}

	launchCleanupProcess(activeExe)

	return nil
}

func backupActiveBinary(activeExe, backupPath string) error {
	_ = os.Remove(backupPath)
	if errRename := os.Rename(activeExe, backupPath); errRename != nil {
		return apperror.WrapWithDetails(errRename, "cmd.update.renameBackup", "E1150", "failed backup rename", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"active": activeExe, "backup": backupPath})
	}

	return nil
}

func placeNewBinary(newBinaryPath, activeExe, backupPath string) error {
	if os.Rename(newBinaryPath, activeExe) != nil && copyFile(newBinaryPath, activeExe) != nil {
		_ = os.Rename(backupPath, activeExe)
		return apperror.NewSimple("failed placing new exe, restored backup", "E1151")
	}

	return nil
}

func launchCleanupProcess(activeExe string) {
	cmd := exec.Command(activeExe, constants.CmdUpdateCleanup)
	cmd.Env = append(os.Environ(), constants.EnvUpdateCleanupDelayMS+"=1500")
	setHiddenProcessAttr(cmd)
	_ = cmd.Start()
}
