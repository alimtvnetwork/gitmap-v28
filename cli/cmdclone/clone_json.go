package cmdclone

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DirectCloneJSONResponse captures structured clone outcome for JSON delegates.
type DirectCloneJSONResponse struct {
	Success  bool   `json:"success"`
	RepoName string `json:"repoName"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Path     string `json:"path"`
}

func emitDirectCloneJSON(resp DirectCloneJSONResponse) {
	b, err := json.Marshal(resp)
	if err == nil {
		fmt.Println(string(b))
	}
}

func parseJSONLine(line string) (DirectCloneJSONResponse, bool) {
	var resp DirectCloneJSONResponse
	if err := json.Unmarshal([]byte(line), &resp); err != nil || resp.Status == "" {
		return resp, false
	}
	return resp, true
}

// ParseCloneJSONResponse extracts structured DirectCloneJSONResponse from raw output.
func ParseCloneJSONResponse(out string) (DirectCloneJSONResponse, bool) {
	lines := strings.Split(out, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
			continue
		}
		if resp, isOk := parseJSONLine(line); isOk {
			return resp, true
		}
	}
	return DirectCloneJSONResponse{}, false
}
