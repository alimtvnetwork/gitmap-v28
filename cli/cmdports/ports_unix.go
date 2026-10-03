//go:build !windows

package cmdports

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type rawPortRecordUnix struct {
	Port        int
	PID         int
	ProcessName string
}

func inspectPorts(opts PortsOptions) ([]PortEntry, error) {
	return inspectPortsUnix(opts)
}

// inspectPortsUnix inspects listening TCP network ports on Unix-like operating systems.
func inspectPortsUnix(opts PortsOptions) ([]PortEntry, error) {
	records, err := queryListeningTCPPortsUnix()
	if err != nil {
		return nil, apperror.WrapSimple(err, "query listening ports failed:")
	}

	fwMap := queryFirewallRulesUnix()

	return assembleUnixEntries(records, fwMap, opts), nil
}

func queryListeningTCPPortsUnix() ([]rawPortRecordUnix, error) {
	cmd := exec.Command("ss", "-tulpn")
	output, err := cmd.Output()

	if err == nil {
		return parseSSListeningLines(string(output)), nil
	}

	netCmd := exec.Command("netstat", "-tlpn")
	netOutput, netErr := netCmd.Output()

	if netErr == nil {
		return parseUnixNetstatListeningLines(string(netOutput)), nil
	}

	return nil, err
}

func parseSSListeningLines(content string) []rawPortRecordUnix {
	lines := strings.Split(content, "\n")
	seen := make(map[int]rawPortRecordUnix)

	for _, line := range lines {
		record, isMatch := parseSingleSSLine(line)
		if !isMatch {
			continue
		}

		if _, exists := seen[record.Port]; !exists {
			seen[record.Port] = record
		}
	}

	return flattenUnixRecordMap(seen)
}

func parseSingleSSLine(line string) (rawPortRecordUnix, bool) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return rawPortRecordUnix{}, false
	}

	proto := strings.ToLower(fields[0])
	if !strings.HasPrefix(proto, "tcp") {
		return rawPortRecordUnix{}, false
	}

	state := strings.ToUpper(fields[1])
	if state != "LISTEN" && state != "LISTENING" {
		return rawPortRecordUnix{}, false
	}

	port := extractPortFromAddressUnix(fields[4])
	if port <= 0 {
		return rawPortRecordUnix{}, false
	}

	procName, pid := parseSSProcessInfo(line)

	return rawPortRecordUnix{
		Port:        port,
		PID:         pid,
		ProcessName: procName,
	}, true
}

func parseSSProcessInfo(line string) (string, int) {
	idx := strings.Index(line, `users:(("`)
	if idx < 0 {
		return "-", 0
	}

	sub := line[idx+9:]
	quoteIdx := strings.Index(sub, `"`)
	procName := "-"

	if quoteIdx > 0 {
		procName = sub[:quoteIdx]
	}

	pidIdx := strings.Index(sub, "pid=")
	pid := extractPIDFromSub(sub, pidIdx)

	return procName, pid
}

func extractPIDFromSub(sub string, pidIdx int) int {
	if pidIdx <= 0 {
		return 0
	}

	pidSub := sub[pidIdx+4:]
	commaIdx := strings.IndexAny(pidSub, ",)")
	if commaIdx <= 0 {
		return 0
	}

	val, _ := strconv.Atoi(pidSub[:commaIdx])

	return val
}

func parseUnixNetstatListeningLines(content string) []rawPortRecordUnix {
	lines := strings.Split(content, "\n")
	seen := make(map[int]rawPortRecordUnix)

	for _, line := range lines {
		record, isMatch := parseSingleUnixNetstatLine(line)
		if !isMatch {
			continue
		}

		if _, exists := seen[record.Port]; !exists {
			seen[record.Port] = record
		}
	}

	return flattenUnixRecordMap(seen)
}

func parseSingleUnixNetstatLine(line string) (rawPortRecordUnix, bool) {
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return rawPortRecordUnix{}, false
	}

	proto := strings.ToLower(fields[0])
	if !strings.HasPrefix(proto, "tcp") {
		return rawPortRecordUnix{}, false
	}

	state := strings.ToUpper(fields[5])
	if state != "LISTEN" {
		return rawPortRecordUnix{}, false
	}

	port := extractPortFromAddressUnix(fields[3])
	if port <= 0 {
		return rawPortRecordUnix{}, false
	}

	procName, pid := parseUnixNetstatPID(fields, len(fields))

	return rawPortRecordUnix{
		Port:        port,
		PID:         pid,
		ProcessName: procName,
	}, true
}

func parseUnixNetstatPID(fields []string, length int) (string, int) {
	if length < 7 {
		return "-", 0
	}

	progField := fields[6]
	slashIdx := strings.Index(progField, "/")

	if slashIdx <= 0 {
		return "-", 0
	}

	pid, _ := strconv.Atoi(progField[:slashIdx])
	procName := progField[slashIdx+1:]

	return procName, pid
}

func extractPortFromAddressUnix(addr string) int {
	idx := strings.LastIndex(addr, ":")
	if idx < 0 || idx >= len(addr)-1 {
		return 0
	}

	port, err := strconv.Atoi(addr[idx+1:])
	if err != nil {
		return 0
	}

	return port
}

func flattenUnixRecordMap(seen map[int]rawPortRecordUnix) []rawPortRecordUnix {
	results := make([]rawPortRecordUnix, 0, len(seen))
	for _, rec := range seen {
		results = append(results, rec)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

func queryFirewallRulesUnix() map[int]string {
	fwMap := make(map[int]string)
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()

	if err == nil {
		parseUFWStatusLines(string(output), fwMap)
		return fwMap
	}

	fwCmd := exec.Command("firewall-cmd", "--list-ports")
	fwOutput, fwErr := fwCmd.Output()

	if fwErr == nil {
		parseFirewalldPortsUnix(string(fwOutput), fwMap)
		return fwMap
	}

	return fwMap
}

func parseUFWStatusLines(content string, fwMap map[int]string) {
	if !strings.Contains(content, "Status: active") {
		return
	}

	for _, line := range strings.Split(content, "\n") {
		parseSingleUFWLine(line, fwMap)
	}
}

func parseSingleUFWLine(line string, fwMap map[int]string) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return
	}

	portToken := fields[0]
	slashIdx := strings.Index(portToken, "/")
	rawPort := portToken

	if slashIdx > 0 {
		rawPort = portToken[:slashIdx]
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil || port < 1 || port > 65535 {
		return
	}

	action := fields[1]
	fwMap[port] = action
}

func parseFirewalldPortsUnix(content string, fwMap map[int]string) {
	for _, token := range strings.Fields(content) {
		slashIdx := strings.Index(token, "/")
		rawPort := token

		if slashIdx > 0 {
			rawPort = token[:slashIdx]
		}

		port, err := strconv.Atoi(rawPort)
		if err == nil && port >= 1 && port <= 65535 {
			fwMap[port] = "Allow"
		}
	}
}

func assembleUnixEntries(
	records []rawPortRecordUnix,
	fwMap map[int]string,
	opts PortsOptions,
) []PortEntry {
	recordMap := make(map[int]rawPortRecordUnix)
	for _, r := range records {
		recordMap[r.Port] = r
	}

	if opts.TargetPort > 0 {
		return []PortEntry{buildTargetUnixEntry(opts.TargetPort, recordMap, fwMap)}
	}

	if opts.CommonOnly {
		return buildCommonUnixEntries(recordMap, fwMap)
	}

	return buildAllListeningUnixEntries(records, fwMap)
}

func buildTargetUnixEntry(
	port int,
	recordMap map[int]rawPortRecordUnix,
	fwMap map[int]string,
) PortEntry {
	rec, isListening := recordMap[port]
	if isListening {
		return buildListeningUnixEntry(rec, fwMap)
	}

	return buildClosedUnixEntry(port, fwMap)
}

func buildCommonUnixEntries(
	recordMap map[int]rawPortRecordUnix,
	fwMap map[int]string,
) []PortEntry {
	entries := make([]PortEntry, 0, len(commonPortsList))
	for _, p := range commonPortsList {
		entry := buildTargetUnixEntry(p, recordMap, fwMap)
		entries = append(entries, entry)
	}

	return entries
}

func buildAllListeningUnixEntries(
	records []rawPortRecordUnix,
	fwMap map[int]string,
) []PortEntry {
	entries := make([]PortEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, buildListeningUnixEntry(r, fwMap))
	}

	return entries
}

func buildListeningUnixEntry(rec rawPortRecordUnix, fwMap map[int]string) PortEntry {
	fw := resolveUnixFirewallStatus(rec.Port, fwMap)
	recText := resolvePortRecommendation(rec.Port, "LISTENING", fw)

	return PortEntry{
		Port:           rec.Port,
		Protocol:       "TCP",
		ProcessName:    rec.ProcessName,
		PID:            rec.PID,
		State:          "LISTENING",
		FirewallStatus: fw,
		Recommendation: recText,
	}
}

func buildClosedUnixEntry(port int, fwMap map[int]string) PortEntry {
	fw := resolveUnixFirewallStatus(port, fwMap)
	recText := resolvePortRecommendation(port, "CLOSED", fw)

	return PortEntry{
		Port:           port,
		Protocol:       "TCP",
		ProcessName:    "-",
		PID:            0,
		State:          "CLOSED",
		FirewallStatus: fw,
		Recommendation: recText,
	}
}

func resolveUnixFirewallStatus(port int, fwMap map[int]string) string {
	status, exists := fwMap[port]
	if exists {
		return status
	}

	return "Allowed (Default)"
}
