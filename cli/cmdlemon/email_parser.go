package cmdlemon

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseEmailCommand extracts a CommandEnvelope from multi-format email body text.
func ParseEmailCommand(msg EmailMessage) (*CommandEnvelope, error) {
	trimmed := strings.TrimSpace(msg.Body)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty email body")
	}

	if env, ok := tryParseJSON(trimmed); ok {
		PopulateDefaults(env, msg, EmailFormatTypeJSON)
		return env, nil
	}
	if env, ok := tryParseYAML(trimmed); ok {
		PopulateDefaults(env, msg, EmailFormatTypeYAML)
		return env, nil
	}
	if env, ok := tryParseText(trimmed); ok {
		PopulateDefaults(env, msg, EmailFormatTypeText)
		return env, nil
	}

	return nil, fmt.Errorf("unrecognized command format")
}

func tryParseJSON(body string) (*CommandEnvelope, bool) {
	raw := ExtractFencedBlock(body, "```json")
	if !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		return nil, false
	}

	var env CommandEnvelope
	if err := json.Unmarshal([]byte(raw), &env); err == nil && len(env.Command) > 0 {
		env.Command = strings.ToLower(env.Command)
		return &env, true
	}

	return nil, false
}

func tryParseText(body string) (*CommandEnvelope, bool) {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}
		cmd, args, isMatch := parseTextCommand(trimmed)
		if isMatch {
			return &CommandEnvelope{Command: cmd, Args: args}, true
		}
	}

	return nil, false
}

func parseTextCommand(line string) (string, []string, bool) {
	low := strings.ToLower(line)
	for _, p := range []string{"cmd:", "command:", "gitmap ", "/"} {
		if !strings.HasPrefix(low, p) {
			continue
		}
		clean := strings.TrimPrefix(strings.TrimSpace(line[len(p):]), "gitmap ")
		parts := strings.Fields(clean)
		if len(parts) > 0 {
			return strings.ToLower(parts[0]), parts[1:], true
		}
	}
	if low == "gitmap" {
		return "gitmap", nil, true
	}

	return "", nil, false
}
