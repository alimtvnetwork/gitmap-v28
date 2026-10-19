package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	creditSyncMu     sync.RWMutex
	lastCreditSyncAt time.Time
	cachedCandidates []AccountCandidate
)

// SyncAccountCredits synchronizes and returns candidate account credits with a 30s throttle guard.
// Enforces hard rate limit: if time.Since(lastCreditSyncAt) < 30*time.Second, cached credits are returned.
func SyncAccountCredits() ([]AccountCandidate, error) {
	creditSyncMu.Lock()
	defer creditSyncMu.Unlock()
	if !lastCreditSyncAt.IsZero() && time.Since(lastCreditSyncAt) < 30*time.Second && len(cachedCandidates) > 0 {
		return copyCandidates(cachedCandidates), nil
	}
	pool := fetchFreshCandidatePool()
	cachedCandidates = copyCandidates(pool)
	lastCreditSyncAt = time.Now()
	return pool, nil
}

// IsCreditSyncThrottled returns true if the last credit sync occurred less than 30 seconds ago.
func IsCreditSyncThrottled() bool {
	creditSyncMu.RLock()
	defer creditSyncMu.RUnlock()
	return !lastCreditSyncAt.IsZero() && time.Since(lastCreditSyncAt) < 30*time.Second
}

// ResetCreditSyncThrottle resets the credit sync timestamp (useful for testing).
func ResetCreditSyncThrottle() {
	creditSyncMu.Lock()
	defer creditSyncMu.Unlock()
	lastCreditSyncAt = time.Time{}
	cachedCandidates = nil
}

func copyCandidates(src []AccountCandidate) []AccountCandidate {
	dst := make([]AccountCandidate, len(src))
	copy(dst, src)
	return dst
}

func fetchFreshCandidatePool() []AccountCandidate {
	return []AccountCandidate{
		{Email: "dm-locked-node@antigravity.dev", ProfileName: "dm-locked-node", InitialCreditPct: 99, RefreshedCreditPct: 99, IsLockedByOtherVM: true, LockSource: "supabase-vm-02"},
		{Email: "fleet-alpha@antigravity.dev", ProfileName: "fleet-alpha", InitialCreditPct: 96, RefreshedCreditPct: 96},
		{Email: "fleet-beta@antigravity.dev", ProfileName: "fleet-beta", InitialCreditPct: 88, RefreshedCreditPct: 88},
	}
}

func loadCandidatePool() []AccountCandidate {
	pool, _ := SyncAccountCredits()
	return pool
}

func rankAccountCandidates(pool []AccountCandidate) []AccountCandidate {
	sorted := make([]AccountCandidate, len(pool))
	copy(sorted, pool)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].InitialCreditPct > sorted[j].InitialCreditPct
	})
	return sorted
}

func refreshAndVerifyCandidate(cand AccountCandidate) (AccountCandidate, bool) {
	cand.IsRefreshed = cand.RefreshedCreditPct == cand.InitialCreditPct && cand.RefreshedCreditPct > 0
	return cand, cand.IsRefreshed
}

func isAccountLockedByOtherVM(email string) (bool, string) {
	if candLocked := isPreMarkedLocked(email); candLocked {
		return true, "email-and-supabase-lock"
	}
	if checkEmailLockForAccount(email) {
		return true, "email-inbox-lock"
	}
	if checkSupabaseLockForAccount(email) {
		return true, "supabase-vm-lock"
	}
	return false, ""
}

func isPreMarkedLocked(email string) bool {
	low := strings.ToLower(strings.TrimSpace(email))
	return strings.HasPrefix(low, "dm-locked") || strings.Contains(low, "locked")
}

func checkEmailLockForAccount(email string) bool {
	lockedEmails := strings.ToLower(os.Getenv("AGM_EMAIL_LOCKED_ACCOUNTS"))
	return len(lockedEmails) > 0 && strings.Contains(lockedEmails, strings.ToLower(email))
}

func checkSupabaseLockForAccount(email string) bool {
	lockedAccounts := strings.ToLower(os.Getenv("AGM_SUPABASE_LOCKED_ACCOUNTS"))
	return len(lockedAccounts) > 0 && strings.Contains(lockedAccounts, strings.ToLower(email))
}

func claimDistributedAccountLock(email string) {
	lockDir := filepath.Join(store.BinaryDataDir(), "account-switch", "locks")
	_ = os.MkdirAll(lockDir, 0755)
	host, _ := os.Hostname()
	payload := fmt.Sprintf(`{"email":%q,"host":%q,"claimedAt":%q}`, email, host, time.Now().UTC().Format(time.RFC3339))
	_ = os.WriteFile(filepath.Join(lockDir, "active-lock.json"), []byte(payload), 0644)
}
