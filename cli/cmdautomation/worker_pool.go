package cmdautomation

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// PartitionFiles divides the list of files evenly across the specified number of worker groups.
func PartitionFiles(files []string, numWorkers int) [][]string {
	if len(files) == 0 {
		return [][]string{}
	}
	workers := clampWorkerCount(numWorkers, len(files))
	chunks := make([][]string, workers)
	for i, f := range files {
		chunks[i%workers] = append(chunks[i%workers], f)
	}
	return chunks
}

func clampWorkerCount(requested, total int) int {
	if requested <= 0 {
		return 1
	}
	if requested > total {
		return total
	}
	return requested
}

// BuildFileContext populates relative/absolute paths, sizes, and timestamps for a file.
func BuildFileContext(relPath, repoRoot string) FileContext {
	normRel, absPath := filepath.ToSlash(relPath), filepath.ToSlash(filepath.Join(repoRoot, relPath))
	var size, modTime int64
	if info, err := os.Stat(absPath); err == nil {
		size, modTime = info.Size(), info.ModTime().Unix()
	}
	return FileContext{
		FilePath:                 normRel,
		FileName:                 filepath.Base(normRel),
		FileExtension:            filepath.Ext(normRel),
		ParentFolderPath:         filepath.ToSlash(filepath.Dir(normRel)),
		AbsoluteFilePath:         absPath,
		AbsoluteParentFolderPath: filepath.ToSlash(filepath.Dir(absPath)),
		FileSize:                 size,
		ModifiedTimestamp:        modTime,
		Encoding:                 EncodingUTF8,
	}
}

// ExecuteWorkerGroup runs a single worker group process on a chunk of files.
func ExecuteWorkerGroup(opts WorkerRunOptions, chunk []string, rec RuntimeRecord) WorkerRunResultMonad {
	if len(chunk) == 0 {
		return result.Ok(WorkerRunResult{ThreadsUsed: resolveThreads(opts.Threads)})
	}
	stdinData, encErr := buildWorkerStdin(chunk, opts)
	if encErr != nil {
		return result.Fail[WorkerRunResult](encErr)
	}
	return dispatchWorkerProcess(opts, chunk, rec, stdinData)
}

func dispatchWorkerProcess(opts WorkerRunOptions, chunk []string, rec RuntimeRecord, stdin []byte) WorkerRunResultMonad {
	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout(opts.TimeoutSec))
	defer cancel()
	cmd := buildWorkerCmd(ctx, opts, rec)
	cmd.Stdin, cmd.Dir = bytes.NewReader(stdin), opts.Dir
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outBuf, &errBuf
	start := time.Now()
	runErr := cmd.Run()
	res := formatWorkerResult(chunk, opts, &outBuf, &errBuf, cmd.ProcessState, runErr, time.Since(start))
	return result.Ok(res)
}

func formatWorkerResult(c []string, o WorkerRunOptions, out, errBuf *bytes.Buffer, st *os.ProcessState, err error, d time.Duration) WorkerRunResult {
	return WorkerRunResult{
		FilesMatched: len(c), FilesProcessed: len(c), WorkersUsed: 1,
		ThreadsUsed: resolveThreads(o.Threads), Duration: d, ExitCode: resolveExitCode(st, err),
		Stdout: DecodeFromStream(out.Bytes(), o.Encoding), Stderr: DecodeFromStream(errBuf.Bytes(), o.Encoding), HasErrors: err != nil,
	}
}

func buildWorkerStdin(chunk []string, opts WorkerRunOptions) ([]byte, *apperror.AppError) {
	var buf bytes.Buffer
	for _, file := range chunk {
		encoded, err := EncodeToStream(buildContextForFile(file, opts), opts.Encoding)
		if err != nil {
			return nil, err
		}
		buf.Write(encoded)
	}
	return buf.Bytes(), nil
}

func buildContextForFile(file string, opts WorkerRunOptions) FileContext {
	ctx := BuildFileContext(file, opts.Dir)
	if opts.IsPreRead {
		if data, err := os.ReadFile(ctx.AbsoluteFilePath); err == nil {
			ctx.Content = string(data)
		}
	}
	return ctx
}

func buildWorkerCmd(ctx context.Context, opts WorkerRunOptions, rec RuntimeRecord) *exec.Cmd {
	bin := rec.BinaryPath
	if bin == "" {
		bin = rec.BinaryName
	}
	return exec.CommandContext(ctx, bin, resolveWorkerArgs(opts, rec)...)
}

func resolveWorkerArgs(opts WorkerRunOptions, rec RuntimeRecord) []string {
	name, script := normalizeRuntimeKey(rec.Name), opts.Script
	if isScriptFile(opts.CommandType, script) {
		if name == "pwsh" {
			return []string{"-NoProfile", "-File", script}
		}
		if name == "go" || name == "rust" {
			return []string{"run", script}
		}
		return []string{script}
	}
	return resolveInlineArgs(name, script, opts.Encoding)
}

func resolveInlineArgs(name, script, encoding string) []string {
	if name == "node" {
		return []string{"-e", script}
	}
	if name == "go" || name == "rust" {
		return []string{"run", script}
	}
	if name == "pwsh" && !IsUTF16Encoding(encoding) {
		script = "[Console]::OutputEncoding = [System.Text.Encoding]::UTF8; $OutputEncoding = [System.Text.Encoding]::UTF8;\n" + script
	}
	if name == "pwsh" {
		return []string{"-NoProfile", "-Command", script}
	}
	return []string{"-c", script}
}

func isScriptFile(commandType, script string) bool {
	if commandType == "file" || strings.HasSuffix(commandType, "-file") {
		return true
	}
	if commandType == "inline" || commandType == "code" {
		return false
	}
	ext := filepath.Ext(script)
	return ext == ".py" || ext == ".js" || ext == ".ts" || ext == ".ps1" || ext == ".sh" || ext == ".go"
}

func resolveTimeout(sec int) time.Duration {
	if sec <= 0 {
		return 30 * time.Second
	}
	return time.Duration(sec) * time.Second
}

func resolveThreads(threads int) int {
	if threads <= 0 {
		return 1
	}
	return threads
}

func resolveExitCode(state *os.ProcessState, err error) int {
	if state != nil {
		return state.ExitCode()
	}
	if err != nil {
		return 1
	}
	return 0
}

// RunWorkerPool discovers the runtime, gathers files, and executes the worker pool.
func RunWorkerPool(opts WorkerRunOptions) WorkerRunResultMonad {
	probeRes := ProbeRuntime(opts.Runtime)
	if probeRes.IsFailure() {
		return result.Fail[WorkerRunResult](probeRes.AppError())
	}
	dir := opts.Dir
	if dir == "" {
		dir = "."
	}
	files := collectSearchFiles(SearchOptions{Dir: dir, IncludeBinaries: opts.IncludeBinaries, IncludeLargeJson: opts.IncludeLargeJson})
	return ExecuteWorkerGroup(opts, files, probeRes.Value)
}
