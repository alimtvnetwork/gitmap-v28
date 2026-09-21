package cmdssh

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var reservedSSHVerbs = []string{
	"clear", "reset", "rm", "remove", "delete",
	"ls", "list", "nodes", "node", "help",
	"join", "sj", "add", "enroll", "keys", "key",
}

func isReservedSSHVerb(target string) bool {
	clean := strings.ToLower(strings.TrimSpace(target))
	for _, verb := range reservedSSHVerbs {
		isMatch := clean == verb
		if isMatch {
			return true
		}
	}

	return false
}

func isPrefixMatch(target, field string) bool {
	cleanTarget := strings.ToLower(strings.TrimSpace(target))
	cleanField := strings.ToLower(strings.TrimSpace(field))
	hasPrefix := strings.HasPrefix(cleanField, cleanTarget)

	return hasPrefix
}

func matchNodeCandidate(target string, host store.SSHHost) bool {
	clean := strings.TrimSpace(target)
	isAlias := strings.EqualFold(host.Alias, clean) || isPrefixMatch(clean, host.Alias)
	isIP := host.IP == clean || isPrefixMatch(clean, host.IP)
	isID := strings.EqualFold(host.ID, clean) || isPrefixMatch(clean, host.ID)
	userHost := host.Username + "@" + host.IP
	userAlias := host.Username + "@" + host.Alias
	isUserTarget := strings.EqualFold(userHost, clean) || strings.EqualFold(userAlias, clean)
	isCandidate := isAlias || isIP || isID || isUserTarget

	return isCandidate
}

func filterCandidateNodes(target string, hosts []store.SSHHost) []store.SSHHost {
	isAll := target == "all" || target == "--all"
	if isAll {
		return hosts
	}

	var candidates []store.SSHHost
	for _, host := range hosts {
		isMatch := matchNodeCandidate(target, host)
		if isMatch {
			candidates = append(candidates, host)
		}
	}

	return candidates
}

func parseRmTokens(args []string) (string, []string) {
	var target string
	var flags []string
	for _, a := range args {
		clean := strings.TrimSpace(a)
		isFlag := strings.HasPrefix(clean, "-")
		if isFlag {
			flags = append(flags, clean)
		} else if target == "" {
			target = clean
		}
	}

	return target, flags
}
