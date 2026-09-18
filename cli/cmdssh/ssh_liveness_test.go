package cmdssh

import (
	"context"
	"testing"
	"time"
)

func TestLivenessCache(t *testing.T) {
	InvalidateLivenessCache()

	key := "127.0.0.1:2222"
	_, _, isCached := getCachedLiveness(key)
	if isCached {
		t.Fatalf("expected key %s to not be cached initially", key)
	}

	setCachedLiveness(key, true, "reachable")
	isOnline, reason, isCached := getCachedLiveness(key)
	if !isCached || !isOnline || reason != "reachable" {
		t.Fatalf("expected key to be cached as online, got cached=%v online=%v reason=%s", isCached, isOnline, reason)
	}

	InvalidateLivenessCache()
	_, _, isCachedAfter := getCachedLiveness(key)
	if isCachedAfter {
		t.Fatalf("expected cache to be invalidated")
	}
}

func TestProbeTCPQuickUnreachable(t *testing.T) {
	// 192.0.2.1 is reserved documentation TEST-NET-1, guaranteed unreachable
	isReachable := probeTCPQuick("192.0.2.1", 2222, 20*time.Millisecond)
	if isReachable {
		t.Fatalf("expected 192.0.2.1:2222 to be unreachable")
	}
}

func TestCheckConnLivenessOffline(t *testing.T) {
	ctx := context.Background()
	isOnline, reason := CheckConnLiveness(ctx, "192.0.2.1", 2222, 20*time.Millisecond)
	if isOnline {
		t.Fatalf("expected offline status for test net IP, got online with reason: %s", reason)
	}
}
