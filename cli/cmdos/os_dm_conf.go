package cmdos

import "strings"

func resolveWaylandLine(isEnabled bool) string {
	if isEnabled {
		return "WaylandEnable=true"
	}
	return "WaylandEnable=false"
}

func updateWaylandInConfig(content string, isEnabled bool) string {
	lines := strings.Split(content, "\n")
	var result []string
	found := false
	targetLine := resolveWaylandLine(isEnabled)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "WaylandEnable=") || strings.HasPrefix(trimmed, "#WaylandEnable=") {
			found = true
			result = append(result, targetLine)
			continue
		}

		result = append(result, line)
	}

	if !found {
		return insertWaylandUnderDaemon(result, isEnabled)
	}

	return strings.Join(result, "\n")
}

func insertWaylandUnderDaemon(lines []string, isEnabled bool) string {
	val := resolveWaylandLine(isEnabled)
	var output []string
	for _, line := range lines {
		output = append(output, line)
		if strings.TrimSpace(line) == "[daemon]" {
			output = append(output, val)
		}
	}

	return strings.Join(output, "\n")
}
