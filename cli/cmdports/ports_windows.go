//go:build windows

package cmdports

import (
	"encoding/csv"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type rawPortRecord struct {
	Port int
	PID  int
}

func inspectPorts(opts PortsOptions) ([]PortEntry, error) {
	return inspectPortsWindows(opts)
}

// inspectPortsWindows inspects listening TCP network ports on Windows.
func inspectPortsWindows(opts PortsOptions) ([]PortEntry, error) {
	records, err := queryListeningTCPPortsWindows()
	if err != nil {
		return nil, apperror.WrapSimple(err, "query listening ports failed:")
	}

	procMap := queryProcessMapWindows()
	fwMap := queryFirewallRulesWindows()

	return assembleWindowsEntries(records, procMap, fwMap, opts), nil
}

func queryListeningTCPPortsWindows() ([]rawPortRecord, error) {
	cmd := exec.Command("netstat", "-ano", "-p", "tcp")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseNetstatListeningLines(string(output)), nil
}

func parseNetstatListeningLines(content string) []rawPortRecord {
	lines := strings.Split(content, "\n")
	seen := make(map[int]rawPortRecord)

	for _, line := range lines {
		record, isMatch := parseSingleNetstatLine(line)
		if !isMatch {
			continue
		}

		if _, exists := seen[record.Port]; !exists {
			seen[record.Port] = record
		}
	}

	return flattenRecordMap(seen)
}

func parseSingleNetstatLine(line string) (rawPortRecord, bool) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return rawPortRecord{}, false
	}

	if !strings.EqualFold(fields[0], "TCP") {
		return rawPortRecord{}, false
	}

	if !strings.EqualFold(fields[3], "LISTENING") {
		return rawPortRecord{}, false
	}

	port := extractPortFromAddress(fields[1])
	if port <= 0 {
		return rawPortRecord{}, false
	}

	pid, _ := strconv.Atoi(fields[4])

	return rawPortRecord{Port: port, PID: pid}, true
}

func extractPortFromAddress(addr string) int {
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

func flattenRecordMap(seen map[int]rawPortRecord) []rawPortRecord {
	results := make([]rawPortRecord, 0, len(seen))
	for _, rec := range seen {
		results = append(results, rec)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

func queryProcessMapWindows() map[int]string {
	procMap := make(map[int]string)
	procMap[0] = "System Idle"
	procMap[4] = "System"

	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return procMap
	}

	r := csv.NewReader(strings.NewReader(string(output)))
	records, err := r.ReadAll()
	if err != nil {
		return procMap
	}

	for _, row := range records {
		processRowIntoMap(row, procMap)
	}

	return procMap
}

func processRowIntoMap(row []string, procMap map[int]string) {
	if len(row) < 2 {
		return
	}

	pid, err := strconv.Atoi(row[1])
	if err != nil {
		return
	}

	procMap[pid] = row[0]
}

func queryFirewallRulesWindows() map[int]string {
	fwMap := make(map[int]string)
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name=all", "dir=in")
	output, err := cmd.Output()
	if err != nil {
		return fwMap
	}

	return parseNetshRules(string(output))
}

func parseNetshRules(content string) map[int]string {
	fwMap := make(map[int]string)
	blocks := strings.Split(content, "Rule Name:")

	for _, block := range blocks {
		parseSingleNetshBlock(block, fwMap)
	}

	return fwMap
}

func parseSingleNetshBlock(block string, fwMap map[int]string) {
	lines := strings.Split(block, "\n")
	props := parseNetshProperties(lines)

	if !isRuleActiveTCP(props) {
		return
	}

	action := props["action"]
	ports := parseNetshPorts(props["localport"])

	for _, p := range ports {
		fwMap[p] = action
	}
}

func parseNetshProperties(lines []string) map[string]string {
	props := make(map[string]string)
	for _, l := range lines {
		parts := strings.SplitN(l, ":", 2)
		if len(parts) == 2 {
			k := strings.ToLower(strings.TrimSpace(parts[0]))
			v := strings.TrimSpace(parts[1])
			props[k] = v
		}
	}

	return props
}

func isRuleActiveTCP(props map[string]string) bool {
	if !strings.EqualFold(props["enabled"], "Yes") {
		return false
	}

	proto := strings.ToLower(props["protocol"])
	isTCP := proto == "tcp" || proto == "any"

	return isTCP
}

func parseNetshPorts(val string) []int {
	var ports []int
	if val == "" || strings.EqualFold(val, "Any") {
		return ports
	}

	for _, item := range strings.Split(val, ",") {
		p, err := strconv.Atoi(strings.TrimSpace(item))
		if err == nil && p > 0 && p <= 65535 {
			ports = append(ports, p)
		}
	}

	return ports
}

func assembleWindowsEntries(
	records []rawPortRecord,
	procMap map[int]string,
	fwMap map[int]string,
	opts PortsOptions,
) []PortEntry {
	recordMap := make(map[int]rawPortRecord)
	for _, r := range records {
		recordMap[r.Port] = r
	}

	if opts.TargetPort > 0 {
		return []PortEntry{buildTargetPortEntry(opts.TargetPort, recordMap, procMap, fwMap)}
	}

	if opts.CommonOnly {
		return buildCommonPortEntries(recordMap, procMap, fwMap)
	}

	return buildAllListeningEntries(records, procMap, fwMap)
}

func buildTargetPortEntry(
	port int,
	recordMap map[int]rawPortRecord,
	procMap map[int]string,
	fwMap map[int]string,
) PortEntry {
	rec, isListening := recordMap[port]
	if isListening {
		return buildListeningEntry(rec, procMap, fwMap)
	}

	return buildClosedEntry(port, fwMap)
}

func buildCommonPortEntries(
	recordMap map[int]rawPortRecord,
	procMap map[int]string,
	fwMap map[int]string,
) []PortEntry {
	entries := make([]PortEntry, 0, len(commonPortsList))
	for _, p := range commonPortsList {
		entry := buildTargetPortEntry(p, recordMap, procMap, fwMap)
		entries = append(entries, entry)
	}

	return entries
}

func buildAllListeningEntries(
	records []rawPortRecord,
	procMap map[int]string,
	fwMap map[int]string,
) []PortEntry {
	entries := make([]PortEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, buildListeningEntry(r, procMap, fwMap))
	}

	return entries
}

func buildListeningEntry(
	rec rawPortRecord,
	procMap map[int]string,
	fwMap map[int]string,
) PortEntry {
	proc := procMap[rec.PID]
	if proc == "" {
		proc = "-"
	}

	fw := resolveFirewallStatus(rec.Port, fwMap)
	recText := resolvePortRecommendation(rec.Port, "LISTENING", fw)

	return PortEntry{
		Port:           rec.Port,
		Protocol:       "TCP",
		ProcessName:    proc,
		PID:            rec.PID,
		State:          "LISTENING",
		FirewallStatus: fw,
		Recommendation: recText,
	}
}

func buildClosedEntry(port int, fwMap map[int]string) PortEntry {
	fw := resolveFirewallStatus(port, fwMap)
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

func resolveFirewallStatus(port int, fwMap map[int]string) string {
	status, exists := fwMap[port]
	if exists {
		return status
	}

	return "No Rule"
}
