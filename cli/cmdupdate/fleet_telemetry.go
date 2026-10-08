package cmdupdate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseFleetUpdateTelemetry parses JSON summary telemetry from remote execution.
func ParseFleetUpdateTelemetry(raw string, target FleetTarget, execErr error) FleetUpdateTelemetry {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return buildFallbackTelemetry(target, execErr, "Empty response")
	}

	jsonStr := extractJSONSubstring(trimmed)
	var parsed FleetUpdateTelemetry
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
		populateTelemetryDefaults(&parsed, target, execErr)
		return parsed
	}

	var items []map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &items); err == nil {
		return buildArrayTelemetry(items, target, execErr)
	}

	return buildFallbackTelemetry(target, execErr, cleanTelemetryDetails(trimmed))
}

func extractJSONSubstring(s string) string {
	startObj, startArr := strings.Index(s, "{"), strings.Index(s, "[")
	if startArr >= 0 && (startObj < 0 || startArr < startObj) {
		return extractDelimitedRange(s, startArr, "]")
	}
	if startObj >= 0 {
		return extractDelimitedRange(s, startObj, "}")
	}
	return s
}

func extractDelimitedRange(s string, start int, closeDelim string) string {
	end := strings.LastIndex(s, closeDelim)
	if end > start {
		return s[start : end+1]
	}
	return s
}

func cleanTelemetryDetails(raw string) string {
	lines := strings.Split(raw, "\n")
	var candidates []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if isIgnoredTelemetryLine(trimmed) {
			continue
		}
		candidates = append(candidates, trimmed)
	}
	if len(candidates) == 0 {
		return "OK"
	}
	first := candidates[0]
	if len(first) > 50 {
		return first[:47] + "..."
	}
	return first
}

func isIgnoredTelemetryLine(trimmed string) bool {
	if trimmed == "" {
		return true
	}
	low := strings.ToLower(trimmed)
	if strings.HasPrefix(low, "warning:") || strings.HasPrefix(low, "error:") {
		return true
	}
	if strings.Contains(low, "404") || strings.Contains(low, "[test-repoexists]") || strings.Contains(low, "[discovery]") {
		return true
	}
	if strings.Contains(low, "fetching") || strings.Contains(low, "url:") || strings.Contains(low, "gitmap shell wrapper") {
		return true
	}
	return strings.HasPrefix(low, "---") || strings.HasPrefix(low, "===") || strings.Contains(low, "done! run")
}

func populateTelemetryDefaults(t *FleetUpdateTelemetry, target FleetTarget, execErr error) {
	if t.NodeID == "" {
		t.NodeID = target.ID
	}
	if t.Alias == "" {
		t.Alias = target.Alias
	}
	if t.IP == "" {
		t.IP = target.IP
	}
	if execErr != nil {
		t.Success = false
		t.Details = execErr.Error()
	}
}

func buildArrayTelemetry(items []map[string]any, target FleetTarget, execErr error) FleetUpdateTelemetry {
	var updated []string
	var failed []string
	for _, item := range items {
		name, _ := item["name"].(string)
		status, _ := item["status"].(string)
		if status == "updated" || status == "success" || status == "ok" {
			updated = append(updated, name)
			continue
		}
		failed = append(failed, name)
	}
	isSuccess := execErr == nil && len(failed) == 0
	return FleetUpdateTelemetry{
		NodeID:  target.ID,
		Alias:   target.Alias,
		IP:      target.IP,
		Success: isSuccess,
		Updated: updated,
		Failed:  failed,
		Details: fmt.Sprintf("%d updated, %d failed", len(updated), len(failed)),
	}
}

func buildFallbackTelemetry(target FleetTarget, execErr error, raw string) FleetUpdateTelemetry {
	isSuccess := execErr == nil
	details := raw
	if execErr != nil {
		details = execErr.Error()
	}
	return FleetUpdateTelemetry{
		NodeID:  target.ID,
		Alias:   target.Alias,
		IP:      target.IP,
		Success: isSuccess,
		Details: details,
	}
}
