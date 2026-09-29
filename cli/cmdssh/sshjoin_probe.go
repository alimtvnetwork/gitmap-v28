package cmdssh

import (
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func probeHostsLiveStatus(ctx context.Context, hosts []store.SSHHost) []store.SSHHost {
	if len(hosts) == 0 || os.Getenv("GITMAP_SKIP_SSH_PROBE") != "" {
		return hosts
	}
	results := make([]store.SSHHost, len(hosts))
	var wg sync.WaitGroup
	for i, h := range hosts {
		wg.Add(1)
		go func(idx int, targetHost store.SSHHost) {
			defer wg.Done()
			results[idx] = probeSingleHostStatus(ctx, targetHost)
		}(i, h)
	}
	wg.Wait()
	return results
}

func probeSingleHostStatus(ctx context.Context, host store.SSHHost) store.SSHHost {
	if host.Status != "" {
		return host
	}
	port := resolveTablePort(host.Port)
	addr := net.JoinHostPort(host.IP, strconv.Itoa(port))
	dialer := net.Dialer{Timeout: 1200 * time.Millisecond}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		host.Status = classifyHostOfflineStatus(err)
		return host
	}
	_ = conn.Close()

	c := db.SSHConnection{
		Alias:             host.Alias,
		IPAddress:         host.IP,
		Username:          host.Username,
		EncryptedPassword: host.EncryptedPassword,
	}
	client, dialErr := dialNodeWithFallback(c, "")
	if dialErr != nil {
		host.Status = "▲ auth failed"
		return host
	}
	_ = client.Close()
	host.Status = "● ready"
	return host
}

func classifyHostOfflineStatus(err error) string {
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "timeout") {
		return "○ offline (timeout)"
	}
	if strings.Contains(msg, "refused") {
		return "○ offline (refused)"
	}
	return "○ offline"
}

// ProbeHostsLiveStatus tests network reachability and authentication for hosts.
func ProbeHostsLiveStatus(ctx context.Context, hosts []store.SSHHost) []store.SSHHost {
	return probeHostsLiveStatus(ctx, hosts)
}

