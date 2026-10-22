package cmdlemon

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ParseEmailCommand extracts a CommandEnvelope from multi-format email body text.
func ParseEmailCommand(msg EmailMessage) (*CommandEnvelope, error) {
	trimmed := strings.TrimSpace(msg.Body)
	isEmpty := len(trimmed) == 0
	if isEmpty {
		return nil, apperror.NewValidationError("empty email body")
	}

	if env, isParsed := tryParseJSON(trimmed); isParsed {
		PopulateDefaults(env, msg, EmailFormatTypeJSON)

		return env, nil
	}

	if env, isParsed := tryParseYAML(trimmed); isParsed {
		PopulateDefaults(env, msg, EmailFormatTypeYAML)

		return env, nil
	}

	if env, isParsed := tryParseText(trimmed); isParsed {
		PopulateDefaults(env, msg, EmailFormatTypeText)

		return env, nil
	}

	return nil, apperror.NewValidationError("unrecognized command format")
}

func tryParseJSON(body string) (*CommandEnvelope, bool) {
	raw := ExtractFencedBlock(body, "```json")
	trimmedRaw := strings.TrimSpace(raw)
	isObject := strings.HasPrefix(trimmedRaw, "{")
	if !isObject {
		return nil, false
	}

	var env CommandEnvelope
	err := json.Unmarshal([]byte(trimmedRaw), &env)
	isUnmarshalSuccess := err == nil
	hasCommand := len(env.Command) > 0
	isValidJSONCommand := isUnmarshalSuccess && hasCommand
	if isValidJSONCommand {
		env.Command = strings.ToLower(env.Command)

		return &env, true
	}

	return nil, false
}

func tryParseText(body string) (*CommandEnvelope, bool) {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		isEmpty := len(trimmed) == 0
		isComment := strings.HasPrefix(trimmed, "#")
		isIgnoredLine := isEmpty || isComment
		if isIgnoredLine {
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
		hasPrefix := strings.HasPrefix(low, p)
		if !hasPrefix {
			continue
		}

		clean := strings.TrimPrefix(strings.TrimSpace(line[len(p):]), "gitmap ")
		parts := strings.Fields(clean)
		hasTokens := len(parts) > 0
		if hasTokens {
			return strings.ToLower(parts[0]), parts[1:], true
		}
	}

	isGitmapSelf := low == "gitmap"
	if isGitmapSelf {
		return "gitmap", nil, true
	}

	return "", nil, false
}
