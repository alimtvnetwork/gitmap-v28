package cmdssh

import (
	"fmt"
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ExtractRawTokens splits comma or space delimited IP tokens.
func ExtractRawTokens(args []string) []string {
	var tokens []string
	for _, arg := range args {
		splits := strings.Split(arg, ",")
		for _, s := range splits {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				tokens = append(tokens, trimmed)
			}
		}
	}
	return tokens
}

func extractAliasFromToken(token string) (string, string) {
	idxOpen := strings.Index(token, "(")
	idxClose := strings.Index(token, ")")
	if idxOpen != -1 && idxClose > idxOpen {
		ipPart := strings.TrimSpace(token[:idxOpen])
		alias := strings.TrimSpace(token[idxOpen+1 : idxClose])
		return ipPart, alias
	}
	return strings.TrimSpace(token), ""
}

func resolveDefaultAlias(fullIP, customAlias string) string {
	if customAlias != "" {
		return customAlias
	}
	return "node-" + strings.ReplaceAll(fullIP, ".", "-")
}

func parsePrefixAndIP(ipPart string) (string, string, *apperror.AppError) {
	parsed := net.ParseIP(ipPart)
	if parsed == nil {
		return "", "", apperror.NewValidationError(fmt.Sprintf("invalid IPv4 address: %s", ipPart))
	}
	lastDot := strings.LastIndex(ipPart, ".")
	newPrefix := ipPart[:lastDot+1]
	return ipPart, newPrefix, nil
}

func expandSingleIP(ipPart, currentPrefix string) (string, string, *apperror.AppError) {
	if strings.Contains(ipPart, ".") {
		return parsePrefixAndIP(ipPart)
	}

	if currentPrefix == "" {
		return "", "", apperror.NewValidationError(fmt.Sprintf("first IP token must be a full IPv4 address: %s", ipPart))
	}

	fullIP := currentPrefix + ipPart
	parsed := net.ParseIP(fullIP)
	if parsed == nil {
		return "", "", apperror.NewValidationError(fmt.Sprintf("invalid expanded IPv4 address: %s (from octet %s)", fullIP, ipPart))
	}
	return fullIP, currentPrefix, nil
}

// ParseCommonIPTokens parses tokens like 192.168.1.3(w1),7(w2),12(w3) into SSHCommonTarget structs.
func ParseCommonIPTokens(username string, defaultPort int, tokens []string) ([]SSHCommonTarget, *apperror.AppError) {
	if len(tokens) == 0 {
		return nil, apperror.NewValidationError("no IP targets provided for ssh-join-common")
	}

	port := 22
	if defaultPort > 0 {
		port = defaultPort
	}

	var targets []SSHCommonTarget
	currentPrefix := ""

	for _, token := range tokens {
		ipPart, customAlias := extractAliasFromToken(token)
		fullIP, newPrefix, err := expandSingleIP(ipPart, currentPrefix)
		if err != nil {
			return nil, err
		}
		currentPrefix = newPrefix

		alias := resolveDefaultAlias(fullIP, customAlias)
		targets = append(targets, SSHCommonTarget{
			FullIP:   fullIP,
			Port:     port,
			Alias:    alias,
			Username: username,
		})
	}

	return targets, nil
}
