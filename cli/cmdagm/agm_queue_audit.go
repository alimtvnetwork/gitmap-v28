package cmdagm

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const agmQueueAuditHelp = `agm queue audit — cross-check the AGM prompt queue vs dispatch receipts.

Description:
  Reads the Antigravity Manager prompt queue state and compares entries
  marked dispatched against dispatch receipts. Flags PHANTOM entries
  (marked dispatched with no receipt). Gracefully degrades when AGM
  state is absent. Read-only.

Usage:
  gitmap agm queue audit [--repo <dir>]

Flags:
  --repo <dir>    repo root holding .ai-memory/temp/ (default: cwd)
  --help, -h      show this help

Example:
  gitmap agm queue audit
`

// agmQueueEntry mirrors the Antigravity prompt queue entry shape.
type agmQueueEntry struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// agmQueueFile mirrors the queue JSON file shape.
type agmQueueFile struct {
	Active    *agmQueueEntry  `json:"active,omitempty"`
	Queued    []agmQueueEntry `json:"queued,omitempty"`
	UpdatedAt string          `json:"updatedAt"`
}

// agmReceipt is one dispatch receipt record.
type agmReceipt struct {
	ID int    `json:"id"`
	At string `json:"at"`
}

// RunAgmQueueAudit implements `gitmap agm queue audit`.
func RunAgmQueueAudit(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(agmQueueAuditHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("agm queue audit", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "repo root holding .ai-memory/temp/")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(agmQueueAuditHelp)
		return apperror.WrapSimple(parseErr, "agm.queue.audit.flags")
	}
	repoRoot := strings.TrimSpace(*repoFlag)
	if repoRoot == "" {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return apperror.WrapSimple(cwdErr, "agm.queue.audit.cwd")
		}
		repoRoot = cwd
	}

	queuePath := filepath.Join(repoRoot, ".ai-memory", "temp", "agy-prompt-queue.json")
	receiptsPath := filepath.Join(repoRoot, ".ai-memory", "temp", "agy-dispatch-receipts.json")

	queue, isFound := loadAgmQueue(queuePath)
	if !isFound {
		fmt.Printf("agm queue audit: no AGM queue state at %s (AGM not initialized here).\n", queuePath)
		return nil
	}

	entries := queue.Queued
	if queue.Active != nil {
		entries = append([]agmQueueEntry{*queue.Active}, entries...)
	}

	receipts, hasReceipts := loadAgmReceipts(receiptsPath)
	pending := 0
	dispatched := 0
	verified := 0
	var phantoms []agmQueueEntry
	var unverified []agmQueueEntry
	for _, e := range entries {
		status := strings.ToLower(strings.TrimSpace(e.Status))
		isDispatched := status == "dispatched" || status == "sent" || status == "done"
		if !isDispatched {
			pending++
			continue
		}
		dispatched++
		if !hasReceipts {
			unverified = append(unverified, e)
			continue
		}
		if receipts[e.ID] {
			verified++
		} else {
			phantoms = append(phantoms, e)
		}
	}

	fmt.Printf("agm queue audit (%s):\n", repoRoot)
	fmt.Printf("  queue: %d pending · %d dispatched\n", pending, dispatched)
	if !hasReceipts {
		fmt.Printf("  receipts: no receipts log at %s\n", receiptsPath)
		fmt.Printf("  unverified dispatched: %d (cannot confirm sent)\n", len(unverified))
	} else {
		fmt.Printf("  receipts: %d verified ✓", verified)
		if len(phantoms) > 0 {
			fmt.Printf(" · %d PHANTOM", len(phantoms))
		}
		fmt.Println()
	}
	if len(phantoms) > 0 {
		sort.Slice(phantoms, func(i, j int) bool { return phantoms[i].ID < phantoms[j].ID })
		fmt.Println("  phantom entries (marked dispatched, no receipt):")
		for _, p := range phantoms {
			title := p.Title
			if title == "" {
				title = p.Type
			}
			fmt.Printf("    #%d %s\n", p.ID, title)
		}
	}

	return nil
}

// loadAgmQueue reads the queue file; isFound is false when absent/unparseable.
func loadAgmQueue(path string) (agmQueueFile, bool) {
	var q agmQueueFile
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return q, false
	}
	if jsonErr := json.Unmarshal(data, &q); jsonErr != nil {
		return q, false
	}

	return q, true
}

// loadAgmReceipts reads the receipts file into an id set.
func loadAgmReceipts(path string) (map[int]bool, bool) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, false
	}
	var list []agmReceipt
	if jsonErr := json.Unmarshal(data, &list); jsonErr != nil {
		return nil, false
	}
	set := make(map[int]bool, len(list))
	for _, r := range list {
		set[r.ID] = true
	}

	return set, true
}
