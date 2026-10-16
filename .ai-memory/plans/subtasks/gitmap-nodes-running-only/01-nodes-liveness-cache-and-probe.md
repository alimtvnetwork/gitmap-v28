# Subtask 01 — Node Liveness Cache & Optimized Incremental Probe

## File Box
- `cli/cmdnodes/nodes_liveness_cache.go` (new file; create it)
- `cli/cmdnodes/nodes_probe.go` (modify)
- `cli/cmdnodes/nodes_collect.go` (modify)

Do NOT touch any other file outside this file box during the execution of this subtask.

---

## Context & Objectives
In fleet clusters, inactive machines (e.g. `w3`, `w2`, `w4`, `u1`) trigger a 1.2s TCP dial timeout on each probe. Repeated calls to fleet commands repeatedly stall CLI execution by 1.2s.

Subtask 01 implements a cross-process **Node Liveness Cache** with a 10-second TTL backed by a JSON file in `os.TempDir()`. When nodes are collected, cached live/offline statuses within the TTL window are reused immediately (0ms). Only uncached or expired nodes undergo concurrent TCP probing. Cache invalidation is supported via `--no-cache`, `--refresh`, and `-r` flags.

---

## Detailed Component Specifications

### 1. `cli/cmdnodes/nodes_liveness_cache.go` (New File)
Create a standalone cache manager providing thread-safe, cross-process caching:

- **Types & Schema:**
  ```go
  package cmdnodes

  import (
  	"sync"
  	"time"
  )

  const (
  	NodeLivenessDefaultTTL = 10 * time.Second
  	livenessCacheFileName  = "gitmap_nodes_liveness_cache.json"
  )

  type NodeLivenessEntry struct {
  	Alias        string `json:"alias"`
  	Host         string `json:"host"`
  	Port         int    `json:"port"`
  	Status       string `json:"status"`
  	IsOnline     bool   `json:"is_online"`
  	ProbedAtUnix int64  `json:"probed_at_unix"`
  }

  type NodeLivenessCache struct {
  	Version int                          `json:"version"`
  	Entries map[string]NodeLivenessEntry `json:"entries"`
  }
  ```

- **File Management & Atomic Writes:**
  - Path: `resolveLivenessCacheFilePath() string` returning `filepath.Join(os.TempDir(), livenessCacheFileName)`.
  - Atomic File Persistence: `writeLivenessCacheToDisk(cache *NodeLivenessCache) error`:
    - Marshals cache with indentation.
    - Writes to a temporary file: `filepath.Join(os.TempDir(), fmt.Sprintf("%s.tmp.%d", livenessCacheFileName, os.Getpid()))`.
    - Atomically replaces the target file via `os.Rename`.
  - Disk Read: `readLivenessCacheFromDisk() (*NodeLivenessCache, error)`:
    - If file does not exist, return a new cache initialized with an empty map.
    - Unmarshals and validates version.

- **In-Process Concurrency:**
  - Package-level `sync.RWMutex` protecting concurrent read/write calls within the same process.

- **Lookup & Freshness Logic:**
  - Key Generator: `resolveLivenessCacheKey(host string, port int) string` returning `strings.ToLower(fmt.Sprintf("%s:%d", host, port))`.
  - Freshness Guard:
    ```go
    func isCacheEntryFresh(entry NodeLivenessEntry, nowUnix int64, ttl time.Duration) bool {
    	ageSeconds := nowUnix - entry.ProbedAtUnix
    	maxAgeSeconds := int64(ttl.Seconds())
    	return ageSeconds >= 0 && ageSeconds < maxAgeSeconds
    }
    ```
  - Fetch Operations:
    - `GetCachedNodeLiveness(host string, port int, ttl time.Duration) (NodeLivenessEntry, bool)`: Returns entry and `isFresh`.
    - `BatchUpdateNodeLiveness(entries []NodeLivenessEntry) error`: Updates entries in-memory and writes back to disk atomically.
    - `InvalidateNodeLivenessCache() error`: Removes or resets cache entries.

### 2. `cli/cmdnodes/nodes_probe.go` (Modification)
Refactor `probeUnifiedNodesLive` to integrate incremental cache probing:

- **Cache-Aware Signature:**
  - Update `probeUnifiedNodesLive(ctx context.Context, nodes []UnifiedFleetNode, hosts []store.SSHHost, canUseCache bool) []UnifiedFleetNode`.
- **Partitioning Flow:**
  - If `canUseCache` is true:
    - Load cache once.
    - Iterate through `nodes`: check `GetCachedNodeLiveness`.
    - If fresh, populate `node.Status` and `node.IsOnline` directly from cache.
    - If missing or expired, collect into `staleNodes` slice.
  - If `canUseCache` is false or missing, mark all nodes as needing live probe.
- **Incremental Probe:**
  - If `len(staleNodes) == 0`, immediately return the cached `nodes` (0ms delay).
  - Probe only `staleNodes` via `cmdssh.ProbeHostsLiveStatus`.
  - For each probed node, construct a `NodeLivenessEntry` with `ProbedAtUnix: time.Now().Unix()` and batch update the cache.
  - Update the corresponding entries in the primary `nodes` slice.
- **Preserve Error Hygiene:**
  - All errors wrapped via `apperror.WrapSimple(err, "...")`.
  - Positive boolean naming used across all helper predicates.

### 3. `cli/cmdnodes/nodes_collect.go` (Modification)
Propagate caching flags into filter options and node collection:

- **Filter Options Enhancement:**
  - In `nodesFilterOptions`:
    ```go
    type nodesFilterOptions struct {
    	isJSON       bool
    	isFast       bool
    	isRefresh    bool
    	noCache      bool
    	filterSSH    bool
    	filterClust  bool
    	filterSC     bool
    	targetFilter string
    }
    ```
  - In `parseNodesFilterOptions(args []string)`:
    - Recognize `--refresh` and `-r` -> set `opts.isRefresh = true`.
    - Recognize `--no-cache` -> set `opts.noCache = true`.
- **Collect Function Refactoring:**
  - Update `collectUnifiedNodes(ctx context.Context, isFast bool, canUseCache bool) ([]UnifiedFleetNode, error)`:
    - Passes `canUseCache` to `probeUnifiedNodesLive(ctx, nodes, hosts, canUseCache)`.
  - Maintain compatibility by providing default/convenience overloads or cleanly updating internal call sites.

---

## Implementation Checklist
- [ ] Create `cli/cmdnodes/nodes_liveness_cache.go` with `NodeLivenessEntry`, `NodeLivenessCache`, atomic writes, and TTL check.
- [ ] Add unit-tested helper methods for reading and updating cache entries safely.
- [ ] Modify `cli/cmdnodes/nodes_probe.go` to split node list into cached hits and stale misses.
- [ ] Probe only stale misses via TCP dialer, record results, and update cache.
- [ ] Modify `cli/cmdnodes/nodes_collect.go` to parse `-r`, `--refresh`, and `--no-cache`.
- [ ] Forward `canUseCache` flag into `collectUnifiedNodes`.
- [ ] Maintain positive boolean conventions: `isOnline`, `isFresh`, `canUseCache`, `hasSubsystem`.
- [ ] Verify zero absolute paths and zero deep-nested conditionals.
