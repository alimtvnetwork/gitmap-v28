// Package scanner walks directories and detects Git repositories.
//
// The walker uses a small bounded worker pool so independent subtrees are
// crawled in parallel. On large folder trees this is I/O bound and yields
// a meaningful speedup; on small trees the pool collapses to effectively
// sequential work because the dispatch loop short-circuits when only one
// directory is in flight.
//
// Concurrency contract:
//   - Bounded by ScanWorkers (default = runtime.NumCPU(), capped by
//     scanWorkersMax to avoid pathological fd exhaustion on huge trees).
//   - Symlinks are NOT followed (consistent with the previous serial
//     implementation; see spec/01-app/03-scanner.md).
//   - When a `.git` directory is found the parent is recorded as a repo
//     and the subtree is NOT descended further (same rule as before).
//   - The first I/O error from any worker wins and is returned; remaining
//     workers drain and exit. Partial results discovered before the error
//     are still returned so callers can render what was found.
//
// Live progress: callers may pass a Progress callback via ScanOptions to
// observe directory-walked / repo-found counts in near-real time. The
// callback is invoked from a single dedicated goroutine on a fixed cadence
// (see progress.go) so handlers do NOT need to be reentrant or fast — but
// they MUST not block indefinitely or the closing snapshot will be delayed.
package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// scanWorkersMax caps the worker pool regardless of CPU count. Filesystem
// scans are I/O bound but each open dir consumes a file descriptor; 16 is
// well below the default ulimit on every supported platform.
const scanWorkersMax = 16

// MaxScanWorkers exposes the upper bound for callers (e.g. CLI flag
// validators) that want to clamp user-provided values into the supported
// range.
const MaxScanWorkers = scanWorkersMax

// gitFileSniffBytes is the read budget for `.git` regular files when
// checking for the `gitdir:` prefix that marks worktree / absorbed
// submodule checkouts. Real `.git` files are tens of bytes; 256 is a
// generous upper bound that keeps detection cheap on huge trees.
const gitFileSniffBytes = 256

// gitdirPrefix is the literal token a worktree/submodule .git file
// starts with: `gitdir: <path>`. Required prefix-match — anything else
// is treated as a non-git file to avoid false positives.
const gitdirPrefix = "gitdir:"

// DefaultMaxDepth is the hard cap on directory descent below the scan
// root, applied even when no repo has been found on the path. The scan
// root itself is depth 0, its immediate children depth 1, and so on —
// so a value of 4 walks up to four levels of subdirectories below the
// root and refuses to enqueue anything deeper. Chosen to comfortably
// cover typical "code/<org>/<project>/<service>/" layouts while
// preventing runaway walks into dependency trees that slipped past the
// exclude list. Override per-scan via ScanOptions.MaxDepth.
const DefaultMaxDepth = 4

// dirJob pairs a queued directory with its depth below the scan root so
// the worker can decide whether children are still in budget without
// recomputing the depth from path arithmetic. Using a struct (vs a
// `chan string` plus a side map) keeps the depth check lock-free.
type dirJob struct {
	path  string
	depth int
}

// RepoInfo holds raw data extracted from a discovered Git repo.
//
// Depth records the directory level at which the repo was found,
// counted from the scan root (depth 0 = root itself, depth 1 = its
// immediate children, …). Surfaced through ScanRecord.Depth so users
// can audit which repos sit at the boundary of the configured
// MaxDepth cap and decide whether to widen it.
type RepoInfo struct {
	AbsolutePath string
	RelativePath string
	Depth        int
}

// ScanProgress is a snapshot of in-flight scan counters delivered to a
// caller-supplied callback. Snapshots are emitted on a fixed cadence
// while the walker is running and once more when the walker terminates
// (Final == true) — even if the totals didn't change since the last
// emission, so renderers can clear / finalize their line.
type ScanProgress struct {
	// DirsWalked is the number of directories fully read by os.ReadDir.
	DirsWalked int64
	// ReposFound is the number of Git repositories discovered so far.
	ReposFound int64
	// Final marks the terminating snapshot for this scan.
	Final bool
}

// ScanOptions bundles optional hooks and tunables for ScanDirWithOptions.
// Zero-value is valid and equivalent to the legacy ScanDir signature
// (with the depth cap defaulting to DefaultMaxDepth).
type ScanOptions struct {
	// ExcludeDirs is the list of directory base names to skip.
	ExcludeDirs []string
	// Workers is the worker-pool size; <=0 picks the platform default.
	Workers int
	// Progress, when non-nil, is invoked from a single goroutine with
	// throttled snapshots while the scan runs and once more at the end.
	Progress func(ScanProgress)
	// MaxDepth caps the directory levels descended below the scan root.
	// Zero (the field's zero value) means "use DefaultMaxDepth"; a
	// negative value disables the cap entirely (legacy unbounded
	// behavior). Repos discovered at any depth still stop their own
	// subtree as before — the cap only matters for paths that have NOT
	// hit a `.git` marker yet.
	MaxDepth int
	// OnDirError, when non-nil, is invoked once per directory whose
	// ReadDir fails. Receives the absolute directory path and the
	// underlying error. Called from worker goroutines — implementations
	// MUST be goroutine-safe. Independent of the legacy first-error
	// return value, which is preserved for backward compat: callers
	// that want PER-DIR attribution use this callback; callers that
	// just want "did anything go wrong" keep using the err return.
	OnDirError func(path string, err error)
}

// ScanDir walks root recursively and returns all Git repo paths found.
// Subtrees are crawled by a bounded worker pool sized via
// defaultWorkerCount(); results are returned sorted by RelativePath
// (with AbsolutePath as a tiebreaker) so terminal, CSV, and JSON
// renderers produce byte-identical output across runs. See sort.go.
func ScanDir(root string, excludeDirs []string) ([]RepoInfo, error) {
	return ScanDirWithOptions(root, ScanOptions{ExcludeDirs: excludeDirs})
}

// ScanDirWithWorkers walks root using exactly `workers` goroutines.
// A value of 0 (or any negative number) selects the platform default
// from defaultWorkerCount(). Values larger than MaxScanWorkers are
// clamped down to keep the pool under the per-process fd budget.
func ScanDirWithWorkers(root string, excludeDirs []string, workers int) ([]RepoInfo, error) {
	return ScanDirWithOptions(root, ScanOptions{
		ExcludeDirs: excludeDirs,
		Workers:     workers,
	})
}

// ScanDirWithOptions is the full-fat entry point. Use it when you want to
// observe scan progress via opts.Progress. The two thin wrappers above
// exist for backward compatibility with callers that don't care.
func ScanDirWithOptions(root string, opts ScanOptions) ([]RepoInfo, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	return walkParallel(
		absRoot,
		buildExcludeSet(opts.ExcludeDirs),
		resolveWorkerCount(opts.Workers),
		resolveMaxDepth(opts.MaxDepth),
		opts.Progress,
		opts.OnDirError,
	)
}

// resolveWorkerCount normalizes a caller-supplied worker count: 0 / <0
// means "auto", and any positive value is clamped to [1, MaxScanWorkers].
func resolveWorkerCount(requested int) int {
	if requested <= 0 {
		return defaultWorkerCount()
	}
	if requested > scanWorkersMax {
		return scanWorkersMax
	}

	return requested
}

// resolveMaxDepth normalizes a caller-supplied depth cap. 0 (zero-value)
// picks DefaultMaxDepth; negative disables the cap; positive is honored
// verbatim. Returning -1 for "unlimited" lets the hot-path check stay a
// simple `depth+1 > cap`-style comparison while accepting any signed int.
func resolveMaxDepth(requested int) int {
	if requested == 0 {
		return DefaultMaxDepth
	}

	return requested
}

// defaultWorkerCount picks a sensible pool size for the host CPU.
func defaultWorkerCount() int {
	n := runtime.NumCPU()
	if n < 1 {
		return 1
	}
	if n > scanWorkersMax {
		return scanWorkersMax
	}

	return n
}

// buildExcludeSet converts a slice to a set for O(1) lookups.
func buildExcludeSet(dirs []string) map[string]bool {
	set := make(map[string]bool, len(dirs))
	for _, d := range dirs {
		set[d] = true
	}

	return set
}

// scanState bundles the shared mutable state passed to every worker. It
// keeps the worker function tiny (well under the per-func line limit) and
// makes the synchronization rules obvious in one place.
type scanState struct {
	root     string
	exclude  map[string]bool
	maxDepth int // negative = unbounded; otherwise inclusive cap below root

	queue chan dirJob    // pending directories + their depth
	wg    sync.WaitGroup // tracks outstanding queued items, NOT workers

	mu       sync.Mutex
	repos    []RepoInfo
	firstErr error

	// onDirError, when non-nil, is invoked from recordErr with the
	// failing directory path and its error. Set from ScanOptions —
	// see the field doc there for the contract. Stored on the state
	// (rather than passed through every helper) because the callback
	// is invoked from the deepest leaf of the call graph.
	onDirError func(path string, err error)

	dirsWalked atomic.Int64
	reposFound atomic.Int64
}

// snapshot returns the current counters in a single struct. The two
// atomic loads are independent — a snapshot is a near-monotonic estimate,
// not a transactional read — which is fine for human-facing UI updates.
func (st *scanState) snapshot(final bool) ScanProgress {
	return ScanProgress{
		DirsWalked: st.dirsWalked.Load(),
		ReposFound: st.reposFound.Load(),
		Final:      final,
	}
}

// walkParallel runs a fixed-size worker pool that consumes directories
// from an unbounded-capacity FIFO and enqueues child directories back.
// The queue is closed when wg drops to zero — i.e. every dispatched
// directory has been fully processed and produced no new work.
func walkParallel(root string, exclude map[string]bool, workers, maxDepth int, progress func(ScanProgress), onDirError func(string, error)) ([]RepoInfo, error) {
	st := &scanState{
		root:       root,
		exclude:    exclude,
		maxDepth:   maxDepth,
		onDirError: onDirError,
		// Buffer sized generously so workers rarely block on enqueue.
		// A bounded buffer is fine — if it fills, workers backpressure
		// each other, which is acceptable; deadlock is impossible
		// because every send is paired with a wg.Add and the closer
		// only fires after wg.Done across all sends.
		queue: make(chan dirJob, 1024),
	}

	st.wg.Add(1)
	st.queue <- dirJob{path: root, depth: 0}

	stopProgress := startProgress(st, progress)

	var workerWG sync.WaitGroup
	for i := 0; i < workers; i++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			for job := range st.queue {
				st.processDir(job)
				st.wg.Done()
			}
		}()
	}

	// Closer goroutine: once every queued dir has been processed (and
	// thus had a chance to enqueue its children), close the queue so
	// workers exit their range loop.
	go func() {
		st.wg.Wait()
		close(st.queue)
	}()

	workerWG.Wait()
	stopProgress() // emits the final snapshot exactly once

	st.mu.Lock()
	defer st.mu.Unlock()
	// Pin a deterministic order before returning so terminal, CSV,
	// and JSON renderers all see the same row sequence regardless of
	// goroutine scheduling. See sort.go for the chosen key order.
	SortRepos(st.repos)

	return st.repos, st.firstErr
}

// processDir reads one directory and dispatches its child directories
// back onto the queue. Errors short-circuit further enqueues for THIS
// dir but do not stop other workers — the first error is captured and
// returned at the end.
//
// Repo-detection is two-pass on purpose: we MUST scan all entries for a
// `.git` child first, and only descend into siblings if none was found.
// Otherwise a single-pass loop would enqueue earlier-listed subdirs
// (e.g. `outer/submodule/`) before discovering `.git` later in the same
// readdir, violating the "do not descend into a discovered repo" rule.
func (st *scanState) processDir(job dirJob) {
	entries, err := os.ReadDir(job.path)
	st.dirsWalked.Add(1)
	if err != nil {
		st.recordDirErr(job.path, err)

		return
	}
	if st.containsGitMarker(job.path, entries) {
		st.recordRepo(job.path, job.depth)

		return
	}
	// Children sit one level deeper. Skip the descend pass entirely
	// when even the closest child would exceed the depth budget — no
	// allocation, no enqueue, no spurious wg traffic.
	if !st.depthAllows(job.depth + 1) {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		st.handleSubdir(job.path, job.depth+1, entry)
	}
}

// depthAllows reports whether a job at the given depth may still be
// enqueued. Negative maxDepth disables the cap (legacy behavior). The
// scan root is depth 0, its children depth 1, and so on, so a cap of 4
// permits depths 0..4 inclusive (four levels of subdirectories below
// the root).
func (st *scanState) depthAllows(depth int) bool {
	if st.maxDepth < 0 {
		return true
	}

	return depth <= st.maxDepth
}

// containsGitMarker reports whether `dir` is a git repo root. A directory
// counts as a repo when it contains either:
//
//   - a `.git` subdirectory (the standard layout), OR
//   - a `.git` regular file whose contents start with `gitdir:` — the
//     layout used by `git worktree add` linked checkouts and by
//     submodules whose .git was absorbed into the superproject.
//
// The file form is gated on the `gitdir:` prefix so a stray `.git` text
// file (e.g. from a misconfigured editor) does not yield a false repo.
// We read at most gitFileSniffBytes to keep the check cheap on large
// trees — a real `.git` file is ~tens of bytes.
func (st *scanState) containsGitMarker(dir string, entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.Name() != constants.ExtGit {
			continue
		}
		if entry.IsDir() {
			return true
		}
		if isGitdirFile(filepath.Join(dir, entry.Name())) {
			return true
		}
	}

	return false
}

// isGitdirFile returns true when path is a regular file beginning with
// the `gitdir:` prefix. Read errors are treated as "not a marker" so a
// transient permission glitch silently skips the candidate rather than
// failing the whole scan.
func isGitdirFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, gitFileSniffBytes)
	n, _ := f.Read(buf)

	return strings.HasPrefix(string(buf[:n]), gitdirPrefix)
}

// handleSubdir applies the exclude filter and enqueues the subdir at
// `childDepth` for further walking. `.git` is handled by the caller
// (processDir) so it is never seen here. Caller is responsible for
// ensuring childDepth is in budget — handleSubdir itself does not
// re-check, since processDir's outer guard already did.
func (st *scanState) handleSubdir(parent string, childDepth int, entry os.DirEntry) {
	name := entry.Name()
	if st.exclude[name] {
		return
	}
	st.enqueue(dirJob{path: filepath.Join(parent, name), depth: childDepth})
}

// enqueue dispatches a directory job for processing.
func (st *scanState) enqueue(job dirJob) {
	st.wg.Add(1)
	st.queue <- job
}

// recordRepo appends a discovered repo (parent of the .git dir) under
// the shared mutex. Repo recording is the only mutex contention point.
// `depth` is the directory level at which the repo was found relative
// to the scan root and is propagated into RepoInfo.Depth so users can
// audit boundary cases against the configured MaxDepth cap.
func (st *scanState) recordRepo(repoPath string, depth int) {
	rel, err := filepath.Rel(st.root, repoPath)
	if err != nil {
		st.recordDirErr(repoPath, err)

		return
	}
	st.mu.Lock()
	st.repos = append(st.repos, RepoInfo{
		AbsolutePath: repoPath,
		RelativePath: rel,
		Depth:        depth,
	})
	st.mu.Unlock()
	st.reposFound.Add(1)
}

// recordErr stores the FIRST error to occur. Later errors are dropped
// to keep the public signature single-error and avoid a noisy
// multi-error. Prefer recordDirErr where a path is available so the
// optional OnDirError callback gets per-dir attribution.
func (st *scanState) recordErr(err error) {
	st.mu.Lock()
	if st.firstErr == nil {
		st.firstErr = err
	}
	st.mu.Unlock()
}

// recordDirErr captures err under the same first-error policy AND
// fires the optional OnDirError callback so callers can build a
// per-directory failure list. The callback runs OUTSIDE the state
// mutex to keep contention low and to let user callbacks block on
// their own collectors without serializing the whole walker.
func (st *scanState) recordDirErr(path string, err error) {
	st.recordErr(err)
	if st.onDirError != nil {
		st.onDirError(path, err)
	}
}
