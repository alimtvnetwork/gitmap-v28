package cmddb

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var defaultCandidatePorts = []int{8080, 8081, 8082, 3000, 5000}

// ParseCandidatePorts parses comma-separated port strings or single integer strings.
// Defaults to []int{8080, 8081, 8082, 3000, 5000} if rawInput is empty.
func ParseCandidatePorts(rawInput string) ([]int, error) {
	trimmed := strings.TrimSpace(rawInput)
	if trimmed == "" {
		return append([]int(nil), defaultCandidatePorts...), nil
	}

	parts := strings.Split(trimmed, ",")
	ports := make([]int, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) == "" {
			continue
		}
		port, err := parsePortToken(p)
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}

	if len(ports) == 0 {
		return append([]int(nil), defaultCandidatePorts...), nil
	}
	return ports, nil
}

func parsePortToken(token string) (int, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(token), "$")
	port, err := strconv.Atoi(clean)
	if err != nil || port < 1 || port > 65535 {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid port '%s' (must be between 1 and 65535)", token))
	}
	return port, nil
}

// CheckPortAvailable verifies whether a specific port is currently open for binding.
func CheckPortAvailable(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// PickAvailablePort probes candidate ports in order using net.Listen.
// Immediately closes listener upon success and returns the first free port.
func PickAvailablePort(candidates []int) (int, error) {
	if len(candidates) == 0 {
		return 0, apperror.NewValidationError("no candidate ports specified")
	}

	for i, port := range candidates {
		if CheckPortAvailable(port) {
			fmt.Printf("[PORT] Successfully bound to port %d.\n", port)
			return port, nil
		}
		logPortOccupied(port, i, candidates)
	}

	return 0, apperror.NewExecutionError(fmt.Sprintf("all candidate ports [%s] are currently occupied or unavailable", formatPortCandidates(candidates)))
}

func logPortOccupied(port, idx int, candidates []int) {
	if idx+1 < len(candidates) {
		fmt.Printf("[PORT] Port %d is currently occupied. Testing next candidate: %d...\n", port, candidates[idx+1])
		return
	}
	fmt.Printf("[PORT] Port %d is currently occupied. No further candidate ports.\n", port)
}

func formatPortCandidates(ports []int) string {
	strList := make([]string, len(ports))
	for i, p := range ports {
		strList[i] = strconv.Itoa(p)
	}
	return strings.Join(strList, ", ")
}
