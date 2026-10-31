package cmdcdp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func unmarshalTargets(body []byte) ([]CDPTarget, error) {
	var targets []CDPTarget
	if err := json.Unmarshal(body, &targets); err != nil {
		return nil, apperror.WrapSimple(err, "cdp.parse_targets")
	}

	return targets, nil
}

func decodeTargetsResponse(resp *http.Response) ([]CDPTarget, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cdp target discovery returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cdp.read_targets")
	}

	return unmarshalTargets(body)
}

// DiscoverTargets queries the given port for available CDP debugging targets.
func DiscoverTargets(port int) ([]CDPTarget, error) {
	urlStr := fmt.Sprintf("http://127.0.0.1:%d/json/list", port)
	client := http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(urlStr)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cdp.discover_targets")
	}

	defer resp.Body.Close()

	return decodeTargetsResponse(resp)
}

func locatePageTarget(targets []CDPTarget) *CDPTarget {
	for i := range targets {
		isPage := targets[i].Type == "page"
		hasWs := targets[i].WebSocketDebuggerURL != ""
		if isPage && hasWs {
			return &targets[i]
		}
	}

	return nil
}

func locateFallbackTarget(targets []CDPTarget, port int) (*CDPTarget, error) {
	for i := range targets {
		if targets[i].WebSocketDebuggerURL != "" {
			return &targets[i], nil
		}
	}

	return nil, fmt.Errorf("no usable target with WebSocket debugger URL on port %d", port)
}

func selectActiveTarget(targets []CDPTarget, port int) (*CDPTarget, error) {
	pageTarget := locatePageTarget(targets)
	if pageTarget != nil {
		return pageTarget, nil
	}

	return locateFallbackTarget(targets, port)
}

// FindActivePageTarget locates the best available page target with a WebSocket debugger URL.
func FindActivePageTarget(port int) (*CDPTarget, error) {
	targets, err := DiscoverTargets(port)
	if err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no CDP targets found on port %d", port)
	}

	return selectActiveTarget(targets, port)
}

func hasActiveCDPPort(port int) bool {
	client := http.Client{Timeout: 200 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/json/version", port))
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// DiscoverDefaultPort tests common CDP debugging ports and returns the first active port, falling back to 9222.
func DiscoverDefaultPort() int {
	candidatePorts := []int{9222, 49994, 9229, 8315}
	for _, port := range candidatePorts {
		if hasActiveCDPPort(port) {
			return port
		}
	}

	return 9222
}
