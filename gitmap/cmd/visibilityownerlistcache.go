// Package cmd — visibilityownerlistcache.go: TTL-bounded SQLite cache
// in front of listOwnerRepos. Cuts repeated `gh repo list` round-trips
// when users iterate on patterns within the cache window. The cache
// is bypassed entirely when the resolved TTL is 0.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §parallel.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// listOwnerReposCached wraps listOwnerRepos with a SQLite-backed TTL
// cache. On HIT it returns the cached slice; on MISS it refreshes
// from the provider and persists the new snapshot.
func listOwnerReposCached(provider, owner string, flags bulkFlags) ([]string, error) {
	ttl := resolveOwnerRepoListTTL(flags)
	db, dbErr := openDB()
	if dbErr == nil && ttl > 0 {
		if names, age, ok := readOwnerRepoListCache(db, provider, owner, ttl); ok {
			fmt.Fprintf(os.Stdout, constants.MsgBulkCacheHitFmt, len(names), age.Round(time.Second))

			return names, nil
		}
	}

	fmt.Fprintf(os.Stdout, constants.MsgBulkCacheMissFmt, providerCLI(provider))
	names, err := listOwnerRepos(provider, owner, flags.Verbose)
	if err != nil {
		return nil, err
	}

	if dbErr == nil && ttl > 0 {
		writeOwnerRepoListCache(db, provider, owner, names)
	}

	return names, nil
}

// resolveOwnerRepoListTTL picks the per-invocation flag, then the
// persisted Setting, then the compiled default.
func resolveOwnerRepoListTTL(flags bulkFlags) time.Duration {
	if flags.CacheTTLSet {
		return time.Duration(flags.CacheTTLSecs) * time.Second
	}
	if db, err := openDB(); err == nil {
		if raw := db.GetSetting(constants.SettingOwnerRepoListCacheTTL); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil && n >= 0 {
				return time.Duration(n) * time.Second
			}
		}
	}

	return time.Duration(constants.OwnerRepoListCacheTTLSeconds) * time.Second
}

// readOwnerRepoListCache returns (names, age, true) on HIT, or
// (nil, 0, false) on MISS / expired.
func readOwnerRepoListCache(db *store.DB, provider, owner string, ttl time.Duration) ([]string, time.Duration, bool) {
	raw, fetchedAt, ok := db.LookupOwnerRepoListCache(provider, owner)
	if !ok {
		return nil, 0, false
	}
	age := time.Since(fetchedAt)
	if age > ttl {
		return nil, 0, false
	}
	var names []string
	if err := json.Unmarshal([]byte(raw), &names); err != nil {
		return nil, 0, false
	}

	return names, age, true
}

// writeOwnerRepoListCache persists the freshly-fetched names list.
// Errors are non-fatal: the user already has the names in-memory.
func writeOwnerRepoListCache(db *store.DB, provider, owner string, names []string) {
	raw, err := json.Marshal(names)
	if err != nil {
		return
	}
	now := time.Now()
	if err := db.UpsertOwnerRepoListCache(provider, owner, string(raw), now); err != nil {
		fmt.Fprintf(os.Stderr, "make-all-*: cache write failed: %v\n", err)
	}
	if err := db.EnsureOwnerRepoNameIndex(); err == nil {
		if err := db.UpsertOwnerRepoNameIndex(provider, owner, names, now); err != nil {
			fmt.Fprintf(os.Stderr, "make-all-*: name-index write failed: %v\n", err)
		}
	}
}
