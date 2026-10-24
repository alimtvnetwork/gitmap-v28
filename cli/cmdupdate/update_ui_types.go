package cmdupdate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate/enums"
)

// UpdateStageType is the canonical update stage type, defined in enums.
type UpdateStageType = enums.UpdateStageType

// UpdateStage is an alias for UpdateStageType.
type UpdateStage = UpdateStageType

const (
	StageNone        = enums.StageNone
	StageResolving   = enums.StageResolving
	StageDownloading = enums.StageDownloading
	StageVerifying   = enums.StageVerifying
	StageSwapping    = enums.StageSwapping
	StageComplete    = enums.StageComplete
)

// UpdateUIRenderer coordinates progress display and summary cards.
type UpdateUIRenderer interface {
	OnStageStart(stage UpdateStageType, message string)
	OnStageProgress(detail string)
	OnStageSuccess(stage UpdateStageType, summary string)
	OnStageFail(stage UpdateStageType, err error)
	RenderSummaryCard(result UpdateResult)
	RenderAlreadyUpdated(currentVersion, targetVersion string)
	Close()
}

// Dracula / Modern Cyber 24-bit color palette tokens.
const (
	DraculaCyan   = "#8be9fd"
	DraculaPurple = "#bd93f9"
	DraculaGreen  = "#50fa7b"
	DraculaOrange = "#ffb86c"
	DraculaRed    = "#ff5555"
	DraculaGray   = "#6272a4"
	DraculaFg     = "#f8f8f2"
)

// ANSI 24-bit Truecolor escape sequences.
const (
	AnsiDraculaCyan   = "\033[38;2;139;233;253m"
	AnsiDraculaPurple = "\033[38;2;189;147;249m"
	AnsiDraculaGreen  = "\033[38;2;80;250;123m"
	AnsiDraculaOrange = "\033[38;2;255;184;108m"
	AnsiDraculaRed    = "\033[38;2;255;85;85m"
	AnsiDraculaGray   = "\033[38;2;98;114;164m"
	AnsiReset         = "\033[0m"
	AnsiBold          = "\033[1m"
)

// isUpdateTTY checks if stdout is connected to a character device terminal.
func isUpdateTTY() bool {
	if os.Getenv("TERM") == "dumb" {

		return false
	}
	stat, err := os.Stdout.Stat()
	if err != nil {

		return false
	}

	return stat.Mode()&os.ModeCharDevice != 0
}

// isColorSupported checks if colored terminal output is supported.
func isColorSupported() bool {
	if os.Getenv("NO_COLOR") != "" {

		return false
	}

	return isUpdateTTY()
}

// isVTSupported checks if virtual terminal processing is available.
func isVTSupported() bool {
	if !isColorSupported() {

		return false
	}
	if runtime.GOOS != "windows" {

		return true
	}

	return os.Getenv("WT_SESSION") != "" || os.Getenv("TERM_PROGRAM") != "" || isUpdateTTY()
}

// ProgressReader wraps an io.Reader and periodically reports read progress.
type ProgressReader struct {
	io.Reader
	total      int64
	read       int64
	lastUpdate time.Time
	onProgress func(written, total int64)
}

// NewProgressReader constructs a throttled ProgressReader.
func NewProgressReader(r io.Reader, total int64, onProgress func(written, total int64)) *ProgressReader {

	return &ProgressReader{Reader: r, total: total, onProgress: onProgress}
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.read += int64(n)
	now := time.Now()
	if pr.onProgress != nil && (now.Sub(pr.lastUpdate) > 100*time.Millisecond || err != nil) {
		pr.lastUpdate = now
		pr.onProgress(pr.read, pr.total)
	}

	return n, err
}

func copyHTTPBody(resp *http.Response, dstPath string, onProg func(w, t int64)) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	reader := io.Reader(resp.Body)
	if onProg != nil {
		reader = NewProgressReader(resp.Body, resp.ContentLength, onProg)
	}
	_, errCopy := io.Copy(outFile, reader)
	return errCopy
}

func downloadHTTP(url, dstPath string, onProg func(w, t int64)) error {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return apperror.NewExecutionError(fmt.Sprintf("HTTP %d from %s", resp.StatusCode, url))
	}
	return copyHTTPBody(resp, dstPath, onProg)
}

func computeFileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {

		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {

		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
