// sibling_parallel.go — intra-URL parallel sibling probing (#5).
//
// `gitmap clone-next` and the install scripts probe `<base>-vN`,
// `<base>-v{N+1}`, ... sequentially via tryLsRemote. For a base with
// 20+ siblings this dominates wall time. ProbeSiblingsParallel fans
// the probes out across a small worker pool (default 8) and returns
// the highest version that responded plus the per-candidate results.
//
// The implementation deliberately reuses tryLsRemote so failures and
// success criteria stay identical to the single-shot path.
package probe

import (
	"fmt"
	"sync"
)

// SiblingHit captures one probe outcome.
type SiblingHit struct {
	Version int
	URL     string
	Tag     string // highest tag at that URL (may be empty)
	OK      bool
}

// ProbeSiblingsParallel probes `baseURL`-v{start..start+span} in
// parallel using up to `workers` goroutines (clamped to [1,32]).
// urlFmt receives the version int and must return the full clone
// URL (e.g. fmt.Sprintf("%s-v%d", base, v)). The slice is sorted by
// Version ascending. Concurrency-safe; no shared mutable state
// beyond the result slice (protected by a mutex).
func ProbeSiblingsParallel(start, span, workers int, urlFmt func(v int) string) []SiblingHit {
	if span < 1 {
		return nil
	}
	if workers < 1 {
		workers = 1
	}
	if workers > 32 {
		workers = 32
	}
	jobs := make(chan int, span)
	for v := start; v < start+span; v++ {
		jobs <- v
	}
	close(jobs)

	var (
		mu  sync.Mutex
		out = make([]SiblingHit, 0, span)
		wg  sync.WaitGroup
	)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range jobs {
				url := urlFmt(v)
				tag, ok := tryLsRemote(url)
				mu.Lock()
				out = append(out, SiblingHit{Version: v, URL: url, Tag: tag, OK: ok})
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	// Stable ascending order by version for deterministic callers.
	sortHitsByVersion(out)
	return out
}

func sortHitsByVersion(hits []SiblingHit) {
	for i := 1; i < len(hits); i++ {
		for j := i; j > 0 && hits[j-1].Version > hits[j].Version; j-- {
			hits[j-1], hits[j] = hits[j], hits[j-1]
		}
	}
}

// FormatSiblingURL is the default URL builder: `<base>-v<N>`.
func FormatSiblingURL(base string, v int) string { return fmt.Sprintf("%s-v%d", base, v) }
