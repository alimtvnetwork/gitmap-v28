package cmdssh

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type livenessEntry struct {
	isOnline  bool
	checkedAt time.Time
	reason    string
}

var (
	livenessCacheMu sync.RWMutex
	livenessCache   = make(map[string]livenessEntry)
	livenessTTL     = 45 * time.Second
)

func getCachedLiveness(key string) (bool, string, bool) {
	livenessCacheMu.RLock()
	defer livenessCacheMu.RUnlock()

	entry, hasEntry := livenessCache[key]
	if !hasEntry || time.Since(entry.checkedAt) > livenessTTL {
		return false, "", false
	}

	return entry.isOnline, entry.reason, true
}

func setCachedLiveness(key string, isOnline bool, reason string) {
	livenessCacheMu.Lock()
	defer livenessCacheMu.Unlock()

	livenessCache[key] = livenessEntry{
		isOnline:  isOnline,
		checkedAt: time.Now(),
		reason:    reason,
	}
}

// CheckNodeLiveness tests if a node is reachable with short TTL caching.
func CheckNodeLiveness(ctx context.Context, host store.SSHHost, port int, timeout time.Duration) (bool, string) {
	port = resolveHealthPort(port)
	timeout = resolveHealthTimeout(timeout)
	key := host.IP + ":" + strconv.Itoa(port)

	if isOnline, reason, isCached := getCachedLiveness(key); isCached {
		return isOnline, reason
	}

	res := probeHostHealth(ctx, host, port, timeout)
	setCachedLiveness(key, res.IsOnline, res.Details)

	return res.IsOnline, res.Details
}

// CheckConnLiveness checks liveness for an SSH connection record.
func CheckConnLiveness(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
	host := store.SSHHost{IP: ip}

	return CheckNodeLiveness(ctx, host, port, timeout)
}

// InvalidateLivenessCache clears the liveness cache for fresh scans.
func InvalidateLivenessCache() {
	livenessCacheMu.Lock()
	defer livenessCacheMu.Unlock()

	livenessCache = make(map[string]livenessEntry)
}

func probeTCPQuick(ip string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(ip, strconv.Itoa(resolveHealthPort(port)))
	conn, err := net.DialTimeout("tcp", addr, resolveHealthTimeout(timeout))
	if err != nil {
		return false
	}

	_ = conn.Close()

	return true
}
