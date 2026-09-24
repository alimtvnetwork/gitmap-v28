package cmddaemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ExecViaREST attempts to execute a command on a remote node daemon via REST.
func ExecViaREST(host string, port int, token string, req DaemonExecReq) (*DaemonExecResp, error) {
	targetURL := fmt.Sprintf("http://%s:%d/api/v1/exec", host, port)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, apperror.WrapSimple(err, "daemon_marshal")
	}

	httpReq, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, apperror.WrapSimple(err, "daemon_new_request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if token != "" {
		httpReq.Header.Set("X-GitMap-Token", token)
	}

	client := &http.Client{Timeout: time.Duration(resolveTimeoutMs(req.TimeoutMs)) * time.Millisecond}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, apperror.WrapSimple(err, "daemon_http_do")
	}
	defer resp.Body.Close()

	var execResp DaemonExecResp
	if decodeErr := json.NewDecoder(resp.Body).Decode(&execResp); decodeErr != nil {
		return nil, apperror.WrapSimple(decodeErr, "daemon_decode")
	}

	return &execResp, nil
}

func resolveTimeoutMs(timeoutMs int64) int64 {
	isDefault := timeoutMs <= 0
	if isDefault {
		return 30000
	}

	return timeoutMs
}
