package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// IsLocalHostAddress tests if an IP address belongs to the current local machine.
func IsLocalHostAddress(ip string) bool {
	clean := strings.TrimSpace(ip)
	if clean == "" || clean == "127.0.0.1" || clean == "::1" || strings.EqualFold(clean, constants.Localhost) {
		return true
	}
	if host, err := os.Hostname(); err == nil && strings.EqualFold(host, clean) {
		return true
	}
	target := net.ParseIP(clean)
	if target == nil {
		return false
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ipAddr net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ipAddr = v.IP
			case *net.IPAddr:
				ipAddr = v.IP
			}
			if ipAddr != nil && ipAddr.Equal(target) {
				return true
			}
		}
	}
	return false
}

func probeFleetLiveness(conns []db.SSHConnection) []fleetNodeLiveness {
	results := make([]fleetNodeLiveness, len(conns))
	var wg sync.WaitGroup
	for i, c := range conns {
		wg.Add(1)
		idx := i
		targetConn := c
		go func() {
			defer wg.Done()
			if IsLocalHostAddress(targetConn.IPAddress) {
				results[idx] = fleetNodeLiveness{
					conn:        targetConn,
					isOnline:    true,
					isLocalHost: true,
					version:     constants.Version,
				}
				return
			}
			isOnline, _ := CheckConnLiveness(context.Background(), targetConn.IPAddress, 22, 500*time.Millisecond)
			ver := "unknown"
			if isOnline {
				ver = queryFleetNodeVersion(targetConn)
			}
			if ver == "" {
				ver = "unknown"
			}
			results[idx] = fleetNodeLiveness{
				conn:        targetConn,
				isOnline:    isOnline,
				isLocalHost: false,
				version:     ver,
			}
		}()
	}
	wg.Wait()

	return results
}

func printFleetEnqueueBanner(probes []fleetNodeLiveness) {
	fmt.Printf("\n%s  Enqueuing 'pull-all' across SSH fleet:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, p := range probes {
		c := p.conn
		if p.isLocalHost {
			fmt.Printf("    • Remote Node [%s] (%s): %sLocal machine (covered by local execution, skipped)%s\n",
				c.Alias, c.IPAddress, constants.ColorYellow, constants.ColorReset)
			continue
		}
		if !p.isOnline {
			fmt.Printf("    • Remote Node [%s] (%s): %sOffline (skipped, no task enqueued)%s\n",
				c.Alias, c.IPAddress, constants.ColorYellow, constants.ColorReset)
			continue
		}
		verStr := p.version
		if verStr != "" && verStr != "unknown" && !strings.HasPrefix(verStr, "v") {
			verStr = "v" + verStr
		}
		fmt.Printf("    • Remote Node [%s] (%s) [GitMap %s]: %sOnline → Enqueued (async)%s\n",
			c.Alias, c.IPAddress, verStr, constants.ColorGreen, constants.ColorReset)
	}
	host := resolveLocalHostname()
	localVer := constants.Version
	if !strings.HasPrefix(localVer, "v") {
		localVer = "v" + localVer
	}
	fmt.Printf("    • Current Machine [%s (127.0.0.1)] [GitMap %s]: %sRunning locally (direct execution, not enqueued)%s\n\n",
		host, localVer, constants.ColorGreen, constants.ColorReset)
}

func resolveLocalHostname() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return constants.Localhost
	}
	return host
}

func queryFleetNodeVersion(c db.SSHConnection) string {
	if c.BuildVersion != "" {
		return cleanFleetVersionString(c.BuildVersion)
	}
	if ver := tryRESTNodeVersion(c.IPAddress); ver != "" {
		return ver
	}
	if ver := tryQuickSSHProbe(c); ver != "" {
		return ver
	}

	return "unknown"
}

func tryRESTNodeVersion(ip string) string {
	endpoints := []string{
		fmt.Sprintf("http://%s:49152/api/v1/version", ip),
		fmt.Sprintf("http://%s:49152/api/v1/status", ip),
	}
	client := http.Client{Timeout: 300 * time.Millisecond}
	for _, u := range endpoints {
		resp, err := client.Get(u)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		if readErr != nil {
			continue
		}
		var statusResp struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(body, &statusResp); err == nil && statusResp.Version != "" {
			return cleanFleetVersionString(statusResp.Version)
		}
		trimmed := strings.TrimSpace(string(body))
		if trimmed != "" && !strings.HasPrefix(trimmed, "<") && !strings.HasPrefix(trimmed, "{") {
			return cleanFleetVersionString(trimmed)
		}
	}
	return ""
}

func tryQuickSSHProbe(c db.SSHConnection) string {
	type probeResult struct {
		ver string
	}
	ch := make(chan probeResult, 1)
	go func() {
		client, err := dialNodeWithFallback(c, "")
		if err != nil || client == nil {
			ch <- probeResult{ver: ""}
			return
		}
		defer client.Close()
		osType := resolveTargetNodeOS(client, c)
		ver, isInstalled := queryNodeVersionViaSSH(client, osType)
		if isInstalled && ver != "" {
			ch <- probeResult{ver: cleanFleetVersionString(ver)}
			return
		}
		ch <- probeResult{ver: ""}
	}()

	select {
	case res := <-ch:
		return res.ver
	case <-time.After(800 * time.Millisecond):
		return ""
	}
}

func cleanFleetVersionString(v string) string {
	s := strings.TrimSpace(v)
	if idx := strings.IndexAny(s, "\r\n"); idx != -1 {
		s = strings.TrimSpace(s[:idx])
	}
	s = strings.TrimPrefix(s, "gitmap")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "version")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	return strings.TrimSpace(s)
}
