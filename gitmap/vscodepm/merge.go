package vscodepm

// mergePairsWithMode is the mode-aware merge engine. It mutates a
// copy of existing in place and returns the new slice + summary.
// Tag reconciliation is delegated to mergeTags; everything else
// (Name overwrite, Paths union, Added/Updated/Unchanged tally) is
// strategy-independent and lives here.
func mergePairsWithMode(existing []Entry, pairs []Pair, mode MergeMode) ([]Entry, SyncSummary) {
	indexByPath := make(map[string]int, len(existing))
	for i, e := range existing {
		indexByPath[normalizePath(e.RootPath)] = i
	}

	summary := SyncSummary{}

	for _, p := range pairs {
		key := normalizePath(p.RootPath)

		idx, found := indexByPath[key]
		if !found {
			// New entry: there's no "existing" set to merge with, so
			// every mode reduces to "write what the detector gave us"
			// — which is exactly what newEntry already does today.
			existing = append(existing, newEntry(p.RootPath, p.Name, p.Paths, p.Tags))
			indexByPath[key] = len(existing) - 1
			summary.Added++

			continue
		}

		existing, summary = applyMerge(existing, idx, p, mode, summary)
	}

	return existing, summary
}

// applyMerge reconciles one already-present entry against its incoming
// Pair under the given MergeMode and bumps the appropriate counter.
func applyMerge(existing []Entry, idx int, p Pair, mode MergeMode, summary SyncSummary) ([]Entry, SyncSummary) {
	mergedPaths := unionPaths(existing[idx].Paths, p.Paths)
	mergedTags := mergeTags(mode, existing[idx].Tags, p.Tags)
	nameChanged := existing[idx].Name != p.Name
	pathsChanged := len(mergedPaths) != len(existing[idx].Paths)
	tagsChanged := !sameTagSet(existing[idx].Tags, mergedTags)

	if !nameChanged && !pathsChanged && !tagsChanged {
		summary.Unchanged++

		return existing, summary
	}

	existing[idx].Name = p.Name
	existing[idx].Paths = mergedPaths
	existing[idx].Tags = mergedTags
	summary.Updated++

	return existing, summary
}

// sameTagSet reports whether a and b contain exactly the same tags
// (order-insensitive, dedup-aware). Used to decide Unchanged vs
// Updated under non-union modes where the merged slice can be
// SHORTER than the original — a length compare alone would miss it.
func sameTagSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[string]struct{}, len(a))
	for _, t := range a {
		seen[t] = struct{}{}
	}
	for _, t := range b {
		if _, ok := seen[t]; !ok {
			return false
		}
	}

	return true
}
