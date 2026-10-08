package cmdpull

import (
	"strings"
)

func ExtractTransportFlags(args []string) (bool, bool, []string) {
	var useSSH, useHTTPS bool
	rest := make([]string, 0, len(args))
	for _, a := range args {
		useSSH, useHTTPS, rest = classifyTransportToken(a, useSSH, useHTTPS, rest)
	}

	return useSSH, useHTTPS, rest
}

func classifyTransportToken(a string, useSSH, useHTTPS bool, rest []string) (bool, bool, []string) {
	switch strings.ToLower(a) {
	case "--ssh", "-ssh", "--sh", "-sh", "ssh":
		return true, useHTTPS, rest
	case "--https", "-https", "--ht", "-ht", "https":
		return useSSH, true, rest
	default:
		return useSSH, useHTTPS, append(rest, a)
	}
}
