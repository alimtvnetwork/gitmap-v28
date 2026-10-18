package cmdlemon

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func tryParseYAML(body string) (*CommandEnvelope, bool) {
	low := strings.ToLower(body)
	if strings.HasPrefix(low, "# command:") {
		return parseYamlHeader(body[10:])
	}
	if strings.Contains(body, "```yaml") {
		return parseYamlFenced(body)
	}

	return nil, false
}

func parseYamlHeader(raw string) (*CommandEnvelope, bool) {
	parts := strings.Fields(strings.TrimSpace(raw))
	if len(parts) == 0 {
		return nil, false
	}

	return &CommandEnvelope{Command: strings.ToLower(parts[0]), Args: parts[1:]}, true
}

func parseYamlFenced(body string) (*CommandEnvelope, bool) {
	var env CommandEnvelope
	raw := ExtractFencedBlock(body, "```yaml")
	err := yaml.Unmarshal([]byte(raw), &env)
	if err != nil || len(env.Command) == 0 {
		return nil, false
	}

	env.Command = strings.ToLower(env.Command)
	return &env, true
}
