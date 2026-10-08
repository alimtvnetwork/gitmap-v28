package cmddownload

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ExecuteDownload runs the multi-tier accelerated download pipeline.
func ExecuteDownload(opts DownloadOptions) ResultDownload {
	destDir, destFile, errDest := resolveDestination(opts.URL, opts.OutputPath)
	if errDest != nil {
		return result.Fail[DownloadResult](wrapDestError(errDest))
	}

	if errMkdir := os.MkdirAll(destDir, 0755); errMkdir != nil {
		return result.Fail[DownloadResult](apperror.NewWithDetails(
			"download.dest_mkdir",
			"E1202",
			fmt.Sprintf("destination directory is unwritable: %v", errMkdir),
			"cmddownload",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		))
	}

	fullPath := filepath.Join(destDir, destFile)
	if errExists := checkTargetExists(fullPath, opts.IsForceOverwrite); errExists != nil {
		return result.Fail[DownloadResult](errExists)
	}

	engine := opts.Engine
	if engine == "" {
		engine = EngineAuto
	}

	switch engine {
	case EngineAria2c:
		res, err := downloadWithAria2c(opts, destDir, destFile)
		if err != nil {
			return handleDownloadError(opts, fullPath, destFile, EngineAria2c, 1, err)
		}

		return result.Ok(res)

	case EngineCurl:
		res, err := downloadWithCurl(opts, destDir, destFile)
		if err != nil {
			return handleDownloadError(opts, fullPath, destFile, EngineCurl, 2, err)
		}

		return result.Ok(res)

	case EngineGoHTTP:
		res, err := downloadWithGoHTTP(opts, destDir, destFile)
		if err != nil {
			return handleDownloadError(opts, fullPath, destFile, EngineGoHTTP, 3, err)
		}

		return result.Ok(res)

	case EngineAuto:
		return runAutoFallbackPipeline(opts, destDir, destFile)

	default:
		return runAutoFallbackPipeline(opts, destDir, destFile)
	}
}

func wrapDestError(err error) *apperror.AppError {
	if appErr, isApp := err.(*apperror.AppError); isApp {
		return appErr
	}

	return apperror.WrapSimple(err, "download.dest")
}

func checkTargetExists(fullPath string, isForceOverwrite bool) *apperror.AppError {
	if isForceOverwrite {
		return nil
	}

	if _, statErr := os.Stat(fullPath); statErr != nil {
		return nil
	}

	return apperror.NewWithDetails(
		"download.file_exists",
		"E1202",
		fmt.Sprintf("file %s already exists; use --force to overwrite", fullPath),
		"cmddownload",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
}

func tryAria2cTier(opts DownloadOptions, destDir, destFile string) (DownloadResult, bool) {
	if !HasAria2c() {
		return DownloadResult{}, false
	}

	res, err := downloadWithAria2c(opts, destDir, destFile)
	if err == nil {
		return res, true
	}

	warnFallback(opts, "aria2c", "curl", err)

	return DownloadResult{}, false
}

func tryCurlTier(opts DownloadOptions, destDir, destFile string) (DownloadResult, bool) {
	if !HasCurl() {
		return DownloadResult{}, false
	}

	res, err := downloadWithCurl(opts, destDir, destFile)
	if err == nil {
		res.HasFallbackOccurred = true

		return res, true
	}

	warnFallback(opts, "curl", "native Go HTTP", err)

	return DownloadResult{}, false
}

func warnFallback(opts DownloadOptions, from, to string, err error) {
	if opts.IsQuiet || opts.IsJSON {
		return
	}

	fmt.Fprintf(os.Stderr, "  [warn] %s download failed (%v), falling back to %s...\n", from, err, to)
}

func runAutoFallbackPipeline(opts DownloadOptions, destDir, destFile string) ResultDownload {
	fullPath := filepath.Join(destDir, destFile)
	if res, ok := tryAria2cTier(opts, destDir, destFile); ok {
		return result.Ok(res)
	}

	if res, ok := tryCurlTier(opts, destDir, destFile); ok {
		return result.Ok(res)
	}

	// Tier 3: Native Go net/http
	res, err := downloadWithGoHTTP(opts, destDir, destFile)
	if err == nil {
		res.HasFallbackOccurred = true

		return result.Ok(res)
	}

	return handleDownloadError(opts, fullPath, destFile, EngineGoHTTP, 3, err)
}

func handleDownloadError(opts DownloadOptions, destPath, fileName string, engine DownloadEngine, tier int, err error) ResultDownload {
	errStr := err.Error()
	res := DownloadResult{
		Status:              "error",
		URL:                 opts.URL,
		DestinationPath:     destPath,
		FileName:            fileName,
		EngineUsed:          engine,
		Tier:                tier,
		HasFallbackOccurred: tier > 1,
		ErrorMessage:        &errStr,
	}

	if appErr, isApp := err.(*apperror.AppError); isApp {
		return result.Fail[DownloadResult](appErr)
	}

	return result.Fail[DownloadResult](apperror.NewWithDetails(
		"download.exhausted",
		"E1203",
		fmt.Sprintf("download failed: %v", err),
		"cmddownload",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		map[string]any{"result": res},
	))
}

func resolveDestination(urlStr, outPath string) (string, string, error) {
	if strings.TrimSpace(urlStr) == "" {
		return "", "", apperror.NewSimple("download URL cannot be empty", "E1201")
	}

	parsed, err := url.Parse(urlStr)
	if err != nil || parsed.Scheme == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", "", apperror.NewSimple(fmt.Sprintf("invalid download URL '%s': must begin with http:// or https://", urlStr), "E1201")
	}

	defaultName := path.Base(parsed.Path)
	if defaultName == "" || defaultName == "/" || defaultName == "." {
		defaultName = "downloaded_file"
	}

	if outPath == "" {
		return resolveCurrentDir(), defaultName, nil
	}

	info, errStat := os.Stat(outPath)
	if errStat == nil && info.IsDir() {
		return outPath, defaultName, nil
	}

	if strings.HasSuffix(outPath, "/") || strings.HasSuffix(outPath, "\\") {
		return outPath, defaultName, nil
	}

	return filepath.Dir(outPath), filepath.Base(outPath), nil
}

func resolveCurrentDir() string {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		return "."
	}

	return cwd
}

func outputDownloadJSON(opts DownloadOptions, res ResultDownload) error {
	val := resolveDownloadResultValue(opts, res)
	encoded, errEnc := json.MarshalIndent(val, "", "  ")
	if errEnc != nil {
		fmt.Printf("{\"status\":\"error\",\"error\":\"%s\"}\n", errEnc.Error())

		return errEnc
	}

	fmt.Println(string(encoded))

	if res.IsFailure() {
		return res.Err
	}

	return nil
}

func resolveDownloadResultValue(opts DownloadOptions, res ResultDownload) DownloadResult {
	if !res.IsFailure() {
		return res.Value
	}

	val := DownloadResult{
		Status:          "error",
		URL:             opts.URL,
		DestinationPath: opts.OutputPath,
		EngineUsed:      opts.Engine,
	}
	if res.Err != nil {
		msg := res.Err.Error()
		val.ErrorMessage = &msg
	}

	return val
}

// RunDownloadCLI dispatches the CLI download subcommand.
func RunDownloadCLI(args []string) error {
	opts, errParse := parseDownloadArgs(args)
	if errParse != nil {
		return errParse
	}

	res := ExecuteDownload(opts)
	if opts.IsJSON {
		return outputDownloadJSON(opts, res)
	}

	if res.IsFailure() {
		return res.Err
	}

	return nil
}

func parseDownloadArgs(args []string) (DownloadOptions, error) {
	opts := DownloadOptions{
		Engine:       EngineAuto,
		Threads:      16,
		Splits:       80,
		MinSplitSize: "1M",
	}

	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			printDownloadHelp()

			return opts, nil

		case arg == "-o" || arg == "--out":
			if i+1 < len(args) {
				opts.OutputPath = args[i+1]
				i++
			}

		case strings.HasPrefix(arg, "--out="):
			opts.OutputPath = strings.TrimPrefix(arg, "--out=")

		case strings.HasPrefix(arg, "-o="):
			opts.OutputPath = strings.TrimPrefix(arg, "-o=")

		case arg == "-j" || arg == "--json":
			opts.IsJSON = true

		case arg == "-q" || arg == "--quiet":
			opts.IsQuiet = true

		case arg == "-f" || arg == "--force":
			opts.IsForceOverwrite = true

		case arg == "-t" || arg == "--threads":
			if n, ok := parseIntArg(args, i+1); ok {
				opts.Threads = n
				i++
			}

		case strings.HasPrefix(arg, "--threads="):
			if n, err := strconv.Atoi(strings.TrimPrefix(arg, "--threads=")); err == nil {
				opts.Threads = n
			}

		case arg == "-s" || arg == "--splits":
			if n, ok := parseIntArg(args, i+1); ok {
				opts.Splits = n
				i++
			}

		case strings.HasPrefix(arg, "--splits="):
			if n, err := strconv.Atoi(strings.TrimPrefix(arg, "--splits=")); err == nil {
				opts.Splits = n
			}

		case arg == "--min-split":
			if i+1 < len(args) {
				opts.MinSplitSize = args[i+1]
				i++
			}

		case strings.HasPrefix(arg, "--min-split="):
			opts.MinSplitSize = strings.TrimPrefix(arg, "--min-split=")

		case arg == "--engine":
			if i+1 < len(args) {
				opts.Engine = DownloadEngine(strings.ToLower(args[i+1]))
				i++
			}

		case strings.HasPrefix(arg, "--engine="):
			opts.Engine = DownloadEngine(strings.ToLower(strings.TrimPrefix(arg, "--engine=")))

		case !strings.HasPrefix(arg, "-"):
			positional = append(positional, arg)
		}
	}

	if len(positional) == 0 {
		printDownloadHelp()

		return opts, apperror.NewSimple("missing required <url> argument", "E1201")
	}

	opts.URL = positional[0]

	return opts, nil
}

func printDownloadHelp() {
	fmt.Println(`Usage: gitmap download <url> [flags]
Alias: gitmap dl <url> [flags]

High-speed accelerated multi-tier download engine with aria2c auto-detection,
curl/Go fallbacks, centered visual progress bar, and JSON telemetry.

Flags:
  -o, --out <path>       Target output file or destination directory
  -j, --json             Emit structured JSON telemetry envelope to stdout
  -t, --threads <n>      Maximum parallel network connections (default 16)
  -s, --splits <n>       Maximum segmentation chunks (default 80)
      --min-split <sz>   Minimum chunk size per split (default 1M)
  -f, --force            Overwrite existing output file if present
  -q, --quiet            Suppress terminal progress output
      --engine <name>    Force downloader engine: 'aria2c', 'curl', 'go_http', or 'auto' (default)
  -h, --help             Show this help message`)
}

func parseIntArg(args []string, idx int) (int, bool) {
	if idx >= len(args) {
		return 0, false
	}

	n, err := strconv.Atoi(args[idx])
	if err != nil {
		return 0, false
	}

	return n, true
}
