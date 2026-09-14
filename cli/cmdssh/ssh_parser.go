package cmdssh

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type SSHTarget struct {
	Username string
	IP       string
	Port     int
}

func (t *SSHTarget) String() string {
	return fmt.Sprintf("%s@%s", t.Username, t.IP)
}

type SSHJoinOptions struct {
	RawTarget  string
	Target     *SSHTarget
	Alias      string
	IsPushAuth bool
	IsShowHelp bool
}

func (o *SSHJoinOptions) PushAuth() bool {
	return o.IsPushAuth
}

func (o *SSHJoinOptions) ShowHelp() bool {
	return o.IsShowHelp
}

func resolveDefaultUsername() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}

	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}

	return "root"
}

func isStandardPort(port int) bool {
	return port <= 0 || port == 22
}

func generateDefaultAlias(ip string, port int) string {
	cleanIP := strings.Trim(ip, "[]")
	if !isStandardPort(port) {
		return fmt.Sprintf("host-%s-%d", cleanIP, port)
	}

	return fmt.Sprintf("host-%s", cleanIP)
}

func isPortValid(port int) bool {
	return port >= 1 && port <= 65535
}

func resolveDefaultPort(port int) int {
	if port <= 0 {
		return 22
	}

	return port
}

func parsePortString(rawPort string) (int, error) {
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", rawPort, err)
	}

	if !isPortValid(port) {
		return 0, fmt.Errorf("port %d out of range (1-65535)", port)
	}

	return port, nil
}

func parseBracketedHost(raw string, defaultPort int) (string, int, error) {
	hasPort := strings.Contains(raw, "]:")
	if hasPort {
		return splitBracketedHostPort(raw)
	}

	return strings.Trim(raw, "[]"), resolveDefaultPort(defaultPort), nil
}

func splitBracketedHostPort(raw string) (string, int, error) {
	h, p, err := net.SplitHostPort(raw)
	if err != nil {
		return "", 0, err
	}

	port, err := parsePortString(p)
	return h, port, err
}

func parseColonHost(raw string, defaultPort int) (string, int, error) {
	isMultiColon := strings.Count(raw, ":") > 1
	if isMultiColon {
		return raw, resolveDefaultPort(defaultPort), nil
	}

	return splitSingleColonHost(raw)
}

func splitSingleColonHost(raw string) (string, int, error) {
	h, p, err := net.SplitHostPort(raw)
	if err != nil {
		return "", 0, err
	}
	if h == "" {
		return "", 0, errors.New("empty host")
	}
	port, err := parsePortString(p)
	return h, port, err
}

func splitHostAndPort(raw string, defaultPort int) (string, int, error) {
	if strings.HasPrefix(raw, "[") {
		return parseBracketedHost(raw, defaultPort)
	}

	if strings.Contains(raw, ":") {
		return parseColonHost(raw, defaultPort)
	}

	return raw, resolveDefaultPort(defaultPort), nil
}

func stripPortOrBrackets(s string) string {
	clean := strings.Trim(s, "[]")
	h, _, err := net.SplitHostPort(s)
	if err == nil {
		return strings.Trim(h, "[]")
	}

	return clean
}

func isLegacyIPAtUser(p1 string, p2 string) bool {
	cleanP1 := stripPortOrBrackets(p1)
	isP1IP := net.ParseIP(cleanP1) != nil
	isP2IP := net.ParseIP(p2) != nil

	return isP1IP && !isP2IP
}

func resolveTwoPartTarget(p1 string, p2 string, defaultUser string) (string, string, error) {
	if p1 == "" || p2 == "" {
		return "", "", errors.New("invalid target format")
	}

	if isLegacyIPAtUser(p1, p2) {
		return p2, p1, nil
	}

	return p1, p2, nil
}

func resolveTargetUser(defaultUser string) string {
	if defaultUser != "" {
		return defaultUser
	}

	return resolveDefaultUsername()
}

func splitUserAndHost(raw string, defaultUser string) (string, string, error) {
	user := resolveTargetUser(defaultUser)
	parts := strings.Split(raw, "@")
	if len(parts) == 1 {
		return user, parts[0], nil
	}

	if len(parts) == 2 {
		return resolveTwoPartTarget(parts[0], parts[1], user)
	}

	return "", "", errors.New("invalid target format")
}

func buildSSHTarget(user string, host string, port int) (*SSHTarget, error) {
	if host == "" {
		return nil, errors.New("empty host")
	}

	cleanHost := strings.Trim(host, "[]")

	return &SSHTarget{
		Username: user,
		IP:       cleanHost,
		Port:     port,
	}, nil
}

func wrapTargetError(raw string, msg string) *apperror.AppError {
	return &apperror.AppError{
		Op:    "ParseSSHTarget",
		Code:  "E_INTERNAL_ERROR",
		Ctx:   map[string]any{"raw": raw},
		Cause: errors.New(msg),
	}
}

// ParseSSHTarget parses a raw SSH target string into an SSHTarget struct.
func ParseSSHTarget(raw string, defaultUser string, defaultPort int) (*SSHTarget, error) {
	if raw == "" {
		return nil, wrapTargetError(raw, "empty target")
	}

	user, hostPort, err := splitUserAndHost(raw, defaultUser)
	if err != nil {
		return nil, wrapTargetError(raw, err.Error())
	}

	return resolveSSHTargetFromHost(user, hostPort, defaultPort, raw)
}

func resolveSSHTargetFromHost(user, hostPort string, defaultPort int, raw string) (*SSHTarget, error) {
	host, port, err := splitHostAndPort(hostPort, defaultPort)
	if err != nil {
		return nil, wrapTargetError(raw, err.Error())
	}

	return buildSSHTarget(user, host, port)
}

type parsedFlags struct {
	flagAlias  string
	flagUser   string
	flagPort   int
	isPushAuth bool
	isShowHelp bool
	positional []string
}

func isHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

func extractNextArg(args []string, idx int, flagName string) (string, int, bool, error) {
	if idx+1 >= len(args) {
		return "", idx, true, fmt.Errorf("flag %s requires an argument", flagName)
	}

	return args[idx+1], idx + 2, true, nil
}

func extractFlagValue(args []string, idx int, flags ...string) (string, int, bool, error) {
	arg := args[idx]
	for _, f := range flags {
		if arg == f {
			return extractNextArg(args, idx, f)
		}

		prefix := f + "="
		if strings.HasPrefix(arg, prefix) {
			return strings.TrimPrefix(arg, prefix), idx + 1, true, nil
		}
	}

	return "", idx, false, nil
}

func parseAliasFlag(args []string, idx int, p *parsedFlags) (int, bool, error) {
	val, nextIdx, isMatch, err := extractFlagValue(args, idx, "--name", "--alias", "-n")
	if !isMatch || err != nil {
		return nextIdx, isMatch, err
	}

	p.flagAlias = val

	return nextIdx, true, nil
}

func parseUserFlag(args []string, idx int, p *parsedFlags) (int, bool, error) {
	val, nextIdx, isMatch, err := extractFlagValue(args, idx, "--user", "-u")
	if !isMatch || err != nil {
		return nextIdx, isMatch, err
	}

	p.flagUser = val

	return nextIdx, true, nil
}

func parsePortFlag(args []string, idx int, p *parsedFlags) (int, bool, error) {
	val, nextIdx, isMatch, err := extractFlagValue(args, idx, "--port", "-p")
	if !isMatch || err != nil {
		return nextIdx, isMatch, err
	}

	port, err := parsePortString(val)
	if err != nil {
		return nextIdx, true, err
	}

	p.flagPort = port

	return nextIdx, true, nil
}

func tryParseOptionFlags(args []string, idx int, p *parsedFlags) (int, bool, error) {
	if nextIdx, isAlias, err := parseAliasFlag(args, idx, p); isAlias || err != nil {
		return nextIdx, isAlias, err
	}

	if nextIdx, isUser, err := parseUserFlag(args, idx, p); isUser || err != nil {
		return nextIdx, isUser, err
	}

	return parsePortFlag(args, idx, p)
}

func parseValuedFlagOrPositional(args []string, idx int, p *parsedFlags) (int, error) {
	nextIdx, isHandled, err := tryParseOptionFlags(args, idx, p)
	if err != nil {
		return nextIdx, err
	}

	if isHandled {
		return nextIdx, nil
	}

	p.positional = append(p.positional, args[idx])

	return idx + 1, nil
}

func parseJoinArg(args []string, idx int, p *parsedFlags) (int, error) {
	arg := args[idx]
	if isHelpFlag(arg) {
		p.isShowHelp = true
		return idx + 1, nil
	}
	if arg == "--auth" {
		p.isPushAuth = true
		return idx + 1, nil
	}
	return parseValuedFlagOrPositional(args, idx, p)
}

func parseAllJoinArgs(args []string) (*parsedFlags, error) {
	p := &parsedFlags{}
	idx := 0
	for idx < len(args) {
		nextIdx, err := parseJoinArg(args, idx, p)
		if err != nil {
			return nil, err
		}

		idx = nextIdx
	}

	return p, nil
}

func resolveJoinAlias(flagAlias string, positional []string, target *SSHTarget) string {
	if flagAlias != "" {
		return flagAlias
	}

	if len(positional) > 1 {
		return positional[1]
	}

	return generateDefaultAlias(target.IP, target.Port)
}

func applyFlagOverrides(target *SSHTarget, p *parsedFlags) {
	if p.flagUser != "" {
		target.Username = p.flagUser
	}

	if p.flagPort > 0 {
		target.Port = p.flagPort
	}
}

func buildJoinOptions(raw string, target *SSHTarget, p *parsedFlags) *SSHJoinOptions {
	applyFlagOverrides(target, p)
	alias := resolveJoinAlias(p.flagAlias, p.positional, target)

	return &SSHJoinOptions{
		RawTarget:  raw,
		Target:     target,
		Alias:      alias,
		IsPushAuth: p.isPushAuth,
		IsShowHelp: false,
	}
}

func assembleJoinOptions(p *parsedFlags) (*SSHJoinOptions, error) {
	if len(p.positional) == 0 {
		return nil, wrapTargetError("", "missing required target")
	}

	rawTarget := p.positional[0]
	target, err := ParseSSHTarget(rawTarget, p.flagUser, p.flagPort)
	if err != nil {
		return nil, err
	}

	return buildJoinOptions(rawTarget, target, p), nil
}

func parseSSHJoinOptions(args []string) (*SSHJoinOptions, error) {
	p, err := parseAllJoinArgs(args)
	if err != nil {
		return nil, wrapTargetError("", err.Error())
	}

	if p.isShowHelp {
		return &SSHJoinOptions{IsShowHelp: true}, nil
	}

	return assembleJoinOptions(p)
}

// ParseSSHJoinOptions parses command-line arguments into SSHJoinOptions.
func ParseSSHJoinOptions(args []string) (*SSHJoinOptions, error) {
	return parseSSHJoinOptions(args)
}
