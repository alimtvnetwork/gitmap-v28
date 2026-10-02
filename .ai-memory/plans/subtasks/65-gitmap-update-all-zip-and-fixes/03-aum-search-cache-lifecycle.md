# Subtask 03: AUM Search Cache Lifecycle & Safe Test Fixtures

> **Parent Plan:** [.ai-memory/plans/pending/65-gitmap-update-all-zip-and-fixes.md](../../pending/65-gitmap-update-all-zip-and-fixes.md)  
> **Spec Reference:** [02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md](../../../../02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md)  
> **Status:** `PENDING`  
> **Target Files:**
> - `cli/cmd/search.go`
> - `cli/searcher/search_history_db.go`
> - `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go`

---

## 1. Technical Objective

1. Clarify and verify the complete lifecycle of AUM search caching:
   - In-memory `MemoryCache` (`aumHotCache`) in package `searcher`.
   - SQLite Split-DB `SearchHotCache` in `.gitmap/data/search/sql.db` keyed by deterministic `DH2D` hashes.
   - Dynamic promotion to `HOT_MEMORY_CACHE (<0.04ms)` on the 2nd search execution hit.
2. Implement explicit `gitmap search clean` cache purge and SQLite vacuuming:
   - Resets memory cache safely under mutex locks.
   - Cleans persistent database tables (`SearchHotCache`, `SearchQueryLog`) and runs `VACUUM`.
3. Create hermetic E2E test fixtures in `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go`:
   - Enforce a strict ban on proprietary brand names (no Amazon, no Prize of Asia).
   - Use purely synthetic identifiers (`sample_service_node`, `generic_alpha_repo`, `test_fixture_beta`).

---

## 2. Implementation Scope & File Edits

### Target 1: `cli/searcher/search_history_db.go`
- **Implement `CleanSearchHotCache() (int, error)`**:
  ```go
  // CleanSearchHotCache purges in-memory search caches, deletes Split-DB cache records, and vacuums SQLite.
  func CleanSearchHotCache() (int, error) {
      aumHotCacheMu.Lock()
      memCount := len(aumHotCache)
      aumHotCache = make(map[string]AUMSearchStatRecord)
      aumHotCacheMu.Unlock()

      db, err := store.OpenSearchSplitDB()
      if err != nil {
          return memCount, nil
      }
      defer db.Close()

      conn := db.Conn()
      res, err := conn.Exec(`DELETE FROM SearchHotCache`)
      if err != nil {
          return memCount, apperror.WrapSimple(err, "CleanSearchHotCache.delete")
      }
      _, _ = conn.Exec(`DELETE FROM SearchQueryLog`)
      _, _ = conn.Exec(`VACUUM`)

      deletedRows, _ := res.RowsAffected()
      totalCleared := int(deletedRows)
      if memCount > totalCleared {
          totalCleared = memCount
      }
      return totalCleared, nil
  }
  ```

### Target 2: `cli/cmd/search.go`
- **Wire `gitmap search clean` in `executeSearchCommand`**:
  ```go
  func executeSearchCommand(query string, limit int, isAiCaller bool) error {
      lowQ := strings.ToLower(strings.TrimSpace(query))
      if lowQ == "history" || lowQ == "top" || lowQ == "stats" || lowQ == "dh2d" {
          return searcher.RenderAUMSearchHistoryTable(limit)
      }
      if lowQ == "clean" || lowQ == "clear" || lowQ == "purge" {
          cleared, err := searcher.CleanSearchHotCache()
          if err != nil {
              return err
          }
          fmt.Printf("\n  %s✓ Search cache purged and database vacuumed (%d entries cleared).%s\n\n",
              constants.ColorGreen, cleared, constants.ColorReset)
          return nil
      }
      // ... Proceed with hot cache lookup and search execution ...
  }
  ```

### Target 3: `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go`
- **Add Hermetic Lifecycle & Clean Test**:
  - Implement `TestTempE2E_AUMSearchCacheLifecycleAndClean(t *testing.T)` under `//go:build tempe2e`.
  - Use strictly synthetic symbols:
    - Search Query: `sample_service_node`
    - File Paths: `generic_alpha_repo/service.go`, `test_fixture_beta/mock.go`
  - Validate 4-step sequence:
    1. First call to `RecordAUMSearchExecution`: warm tier, `LookupHotCachedSearch` returns `isHot = false`.
    2. Second call to `RecordAUMSearchExecution`: promoted, `LookupHotCachedSearch` returns `isHot = true` with matching `DH2D` ID.
    3. Call `searcher.CleanSearchHotCache()`: verifies return count $\ge 1$.
    4. Post-clean verification: `LookupHotCachedSearch` returns `isHot = false`, confirming complete invalidation.

---

## 3. Verification Protocol

- Run targeted temporary E2E test:
  ```bash
  $env:RUN_TEMP_E2E="1"
  go test -v -tags=tempe2e ./cli/tests/e2e -run TestTempE2E_AUMSearchCacheLifecycleAndClean
  ```
- Run standard unit tests:
  ```bash
  go test -v ./cli/searcher -run TestAUMSearch
  ```
- CLI Manual Check:
  - Run `gitmap search "test_query"`.
  - Run `gitmap search "test_query"` (observe `<0.04ms` hit).
  - Run `gitmap search clean` (observe clean confirmation).
