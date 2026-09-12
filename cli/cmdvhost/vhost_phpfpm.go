package cmdvhost

import (
	"net"
	"path/filepath"
	"sort"
	"time"
)

var defaultPHPSocketPatterns = []string{
	"/run/php/php*-fpm.sock",
	"/var/run/php/php*-fpm.sock",
}

func sortSocketsDesc(socks []string) []string {
	sorted := make([]string, len(socks))
	copy(sorted, socks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] > sorted[j]
	})

	return sorted
}

func findSocketInDir(pattern string) string {
	matches, err := filepath.Glob(pattern)
	hasError := err != nil
	if hasError {
		return ""
	}

	hasMatches := len(matches) > 0
	if !hasMatches {
		return ""
	}

	sorted := sortSocketsDesc(matches)

	return sorted[0]
}

func findSocketByPatterns(patterns []string) string {
	for _, p := range patterns {
		sock := findSocketInDir(p)
		hasSock := sock != ""
		if hasSock {
			return sock
		}
	}

	return ""
}

func probeTCPPort(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", address, timeout)
	hasError := err != nil
	if hasError {
		return false
	}

	if conn != nil {
		conn.Close()
	}

	return true
}

// DiscoverFastCGIPassWithPatterns finds PHP-FPM socket or TCP loopback with custom glob patterns.
func DiscoverFastCGIPassWithPatterns(patterns []string) string {
	sock := findSocketByPatterns(patterns)
	hasSock := sock != ""
	if hasSock {
		return "unix:" + sock
	}

	isListening := probeTCPPort("127.0.0.1:9000", 200*time.Millisecond)
	if isListening {
		return "127.0.0.1:9000"
	}

	return defaultVHostFastCGIPass
}

// DiscoverFastCGIPass discovers dynamic PHP-FPM unix socket or 127.0.0.1:9000 loopback.
func DiscoverFastCGIPass() string {
	return DiscoverFastCGIPassWithPatterns(defaultPHPSocketPatterns)
}

// ResolveFastCGIPass returns custom pass if provided, otherwise discovers it.
func ResolveFastCGIPass(customPass string) string {
	hasCustom := customPass != ""
	if hasCustom {
		return customPass
	}

	return DiscoverFastCGIPass()
}
