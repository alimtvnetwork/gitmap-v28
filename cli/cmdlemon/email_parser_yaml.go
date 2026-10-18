package cmdlemon

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func tryParseYAML(body string) (*CommandEnvelope, bool) {
	low := strings.ToLower(body)
	hasCommandHeader := strings.HasPrefix(low, "# command:")
	if hasCommandHeader {
		return parseYamlHeader(body[10:])
	}

	hasFencedYaml := strings.Contains(body, "```yaml")
	if hasFencedYaml {
		return parseYamlFenced(body)
	}

	return nil, false
}

func parseYamlHeader(raw string) (*CommandEnvelope, bool) {
	parts := strings.Fields(strings.TrimSpace(raw))
	hasParts := len(parts) > 0
	if !hasParts {
		return nil, false
	}

	return &CommandEnvelope{Command: strings.ToLower(parts[0]), Args: parts[1:]}, true
}

func parseYamlFenced(body string) (*CommandEnvelope, bool) {
	var env CommandEnvelope
	raw := ExtractFencedBlock(body, "```yaml")
	err := yaml.Unmarshal([]byte(raw), &env)
	isUnmarshalSuccess := err == nil
	hasCommand := len(env.Command) > 0
	isValidYamlCommand := isUnmarshalSuccess && hasCommand
	if !isValidYamlCommand {
		return nil, false
	}

	env.Command = strings.ToLower(env.Command)

	return &env, true
}
