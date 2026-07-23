// Package cmd — reclonetransport.go wires the cfr/cfrp/clone-now
// pipelines into the store-backed IdentifiedTransport column added
// by migration 008. Two halves:
//
//  1. coerceURLToStoredTransport(url) — runs PRE-clone. If the URL
//     pair (HttpsUrl, SshUrl) has a stored transport verdict from a
//     prior scan/reclone, rewrite `url` to that transport so an
//     SSH-origin repo never silently downgrades to HTTPS on reclone.
//  2. persistRecloneTransport(url) — runs POST-clone. Records the
//     transport that was actually used and emits a `gitmap history`
//     row (Command="reclone-transport") so users can audit
//     transport flips with `gitmap history`.
//
// Both functions are fail-open: any store error is warned to stderr
// and the caller continues. The reclone is the user's primary goal;
// transport bookkeeping must never break it.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// coerceURLToStoredTransport returns the URL the actual `git clone`
// should use after consulting the store. When the stored transport
// disagrees with the positional URL's transport, the URL is rewritten
// (and the swap is logged to stderr so the user can audit it).
//
// No-ops (returns input unchanged) when:
//   - the store cannot be opened (warned to stderr),
//   - no row matches HttpsUrl/SshUrl = url,
//   - the stored value matches the URL's current transport,
//   - the rewrite helper rejects the URL shape.
func coerceURLToStoredTransport(url string) string {
	if url == "" {
		return url
	}
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: reclone-transport: store open failed: %v\n", err)
		return url
	}
	defer func() { _ = db.Close() }()

	stored, err := db.LookupRepoIdentifiedTransport(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: reclone-transport: lookup failed for %s: %v\n", url, err)
		return url
	}
	if stored == "" {
		return url
	}
	current := store.ClassifyURLTransport(url)
	if current == stored {
		return url
	}
	if stored == store.RepoTransportSSH {
		if out, ok := ConvertURLToSSH(url); ok && out != url {
			fmt.Fprintf(os.Stderr, "↪ reclone-transport: stored=ssh, coercing %s → %s\n", url, out)
			return out
		}
	}
	if stored == store.RepoTransportHTTPS {
		if out, ok := ConvertURLToHTTPS(url); ok && out != url {
			fmt.Fprintf(os.Stderr, "↪ reclone-transport: stored=https, coercing %s → %s\n", url, out)
			return out
		}
	}
	return url
}

// persistRecloneTransport classifies `url`, writes the verdict to
// every matching Repo row, and inserts a `gitmap history` event
// describing the action. Called from cfr/cfrp after a successful
// clone. Fail-open per the package doc.
func persistRecloneTransport(url string) {
	if url == "" {
		return
	}
	transport := store.ClassifyURLTransport(url)
	if transport == "" {
		return
	}
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: reclone-transport: store open failed: %v\n", err)
		return
	}
	defer func() { _ = db.Close() }()

	rows, err := db.SetRepoIdentifiedTransport(url, transport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: reclone-transport: persist failed for %s: %v\n", url, err)
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	rec := model.CommandHistoryRecord{
		Command:    "reclone-transport",
		Args:       url,
		Flags:      "transport=" + transport,
		StartedAt:  now,
		FinishedAt: now,
		Summary:    fmt.Sprintf("persisted transport=%s rows=%d", transport, rows),
		RepoCount:  int(rows),
	}
	if _, herr := db.InsertHistory(rec); herr != nil {
		fmt.Fprintf(os.Stderr, "warning: reclone-transport: history insert failed: %v\n", herr)
	}
}
