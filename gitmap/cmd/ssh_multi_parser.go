// Package cmd — ssh_multi_parser.go parses comma- or whitespace-separated IP targets.
package cmd

import (
	"strings"
)

// ParseMultiIPList splits a combined list of host addresses or IPs.
func ParseMultiIPList(raw string) []string {
	var results []string
	clean := strings.ReplaceAll(raw, ",", " ")
	for _, token := range strings.Fields(clean) {
		trimmed := strings.TrimSpace(token)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}
