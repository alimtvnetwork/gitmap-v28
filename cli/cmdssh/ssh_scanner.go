package cmdssh

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SSHScanOptions holds configuration parameters for network subnet scanning.
type SSHScanOptions struct {
	Subnet  string
	Port    int
	Timeout time.Duration
	Workers int
}

// SSHScanHost represents a discovered machine on the network.
type SSHScanHost struct {
	IP            string
	Port          int
	IsOnline      bool
	Latency       time.Duration
	Hostname      string
	IsEnrolled    bool
	EnrolledAlias string
	Suggestion    string
}

func resolveScanPort(port int) int {
	if port > 0 {
		return port
	}
	return 22
}

func resolveScanTimeout(d time.Duration) time.Duration {
	if d > 0 {
		return d
	}
	return 800 * time.Millisecond
}

func resolveScanWorkers(workers int) int {
	if workers > 0 {
		return workers
	}
	return 40
}

func isScanLoopback(ip net.IP) bool {
	return ip.IsLoopback()
}

func isIPv4Address(ip net.IP) bool {
	return ip.To4() != nil
}

func formatSubnetSlash24(ip net.IP) string {
	v4 := ip.To4()
	return fmt.Sprintf("%d.%d.%d.0/24", v4[0], v4[1], v4[2])
}

func extractIPv4FromAddr(addr net.Addr) (net.IP, bool) {
	ipNet, isIPNet := addr.(*net.IPNet)
	if isIPNet && isIPv4Address(ipNet.IP) && !isScanLoopback(ipNet.IP) {
		return ipNet.IP, true
	}
	return nil, false
}

func scanAddrsForIPv4(addrs []net.Addr) (net.IP, bool) {
	for _, addr := range addrs {
		if ip, isFound := extractIPv4FromAddr(addr); isFound {
			return ip, true
		}
	}
	return nil, false
}

func findInterfaceIPv4(iface net.Interface) (net.IP, bool) {
	if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
		return nil, false
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, false
	}
	return scanAddrsForIPv4(addrs)
}

func detectLocalSubnet() (string, *apperror.AppError) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "192.168.1.0/24", nil
	}
	for _, iface := range ifaces {
		if ip, isFound := findInterfaceIPv4(iface); isFound {
			return formatSubnetSlash24(ip), nil
		}
	}
	return "192.168.1.0/24", nil
}

func parseIPToSlash24(trimmed string) (string, *apperror.AppError) {
	ip := net.ParseIP(trimmed)
	if ip == nil {
		return "", apperror.NewValidationError(fmt.Sprintf("invalid subnet IP: %s", trimmed))
	}
	return formatSubnetSlash24(ip), nil
}

func normalizeSubnetCIDR(raw string) (string, *apperror.AppError) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return detectLocalSubnet()
	}
	if !strings.Contains(trimmed, "/") {
		return parseIPToSlash24(trimmed)
	}
	return trimmed, nil
}

func parseCIDRRange(cidr string) (*net.IPNet, *apperror.AppError) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, apperror.WrapValidation(err, fmt.Sprintf("invalid CIDR notation: %s", cidr))
	}
	return ipNet, nil
}

func calculateHostCount(mask net.IPMask) uint32 {
	maskVal := binary.BigEndian.Uint32(mask)
	return ^maskVal + 1
}

func intToIPv4String(val uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, val)
	return ip.String()
}

func generateSmallRange(start, total uint32) []string {
	ips := make([]string, 0, total)
	for i := uint32(0); i < total; i++ {
		ips = append(ips, intToIPv4String(start+i))
	}
	return ips
}

func generateIPRangeSlice(start, total uint32) []string {
	if total <= 2 {
		return generateSmallRange(start, total)
	}
	ips := make([]string, 0, total-2)
	for i := uint32(1); i < total-1; i++ {
		ips = append(ips, intToIPv4String(start+i))
	}
	return ips
}

func expandSubnetIPs(cidr string) ([]string, *apperror.AppError) {
	ipNet, err := parseCIDRRange(cidr)
	if err != nil {
		return nil, err
	}
	start := binary.BigEndian.Uint32(ipNet.IP.To4())
	total := calculateHostCount(ipNet.Mask)
	if total > 1024 {
		return nil, apperror.NewValidationError(fmt.Sprintf("subnet too large (%d hosts, max 1024): %s", total, cidr))
	}
	return generateIPRangeSlice(start, total), nil
}

func probeHostTCP(ctx context.Context, ip string, port int, timeout time.Duration) (bool, time.Duration) {
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	dialer := net.Dialer{Timeout: timeout}
	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, 0
	}
	_ = conn.Close()
	return true, time.Since(start).Round(time.Millisecond)
}

func resolveHostName(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "-"
	}
	return strings.TrimSuffix(names[0], ".")
}

func buildEnrolledMap(hosts []store.SSHHost) map[string]string {
	res := make(map[string]string, len(hosts))
	for _, h := range hosts {
		if h.IP != "" {
			res[h.IP] = h.Alias
		}
	}
	return res
}

func fetchEnrolledHostsMap(ctx context.Context) map[string]string {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return map[string]string{}
	}
	defer dbConn.Close()
	hosts, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		return map[string]string{}
	}
	return buildEnrolledMap(hosts)
}

func buildEnrolledSuggestion(alias string) string {
	return fmt.Sprintf("gitmap ssh %s", alias)
}

func buildNewHostSuggestion(ip string) string {
	return fmt.Sprintf("gitmap sj %s", ip)
}

func resolveScanSuggestion(ip, alias string, isEnrolled bool) string {
	if isEnrolled {
		return buildEnrolledSuggestion(alias)
	}
	return buildNewHostSuggestion(ip)
}

func createScanHostResult(ip string, port int, latency time.Duration, enrolledMap map[string]string) SSHScanHost {
	alias, isEnrolled := enrolledMap[ip]
	return SSHScanHost{
		IP:            ip,
		Port:          port,
		IsOnline:      true,
		Latency:       latency,
		Hostname:      resolveHostName(ip),
		IsEnrolled:    isEnrolled,
		EnrolledAlias: alias,
		Suggestion:    resolveScanSuggestion(ip, alias, isEnrolled),
	}
}

func dispatchScanWorker(ctx context.Context, ipChan <-chan string, resChan chan<- SSHScanHost, port int, timeout time.Duration, enrolledMap map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	for ip := range ipChan {
		isOnline, latency := probeHostTCP(ctx, ip, port, timeout)
		if isOnline {
			resChan <- createScanHostResult(ip, port, latency, enrolledMap)
		}
	}
}

func feedIPsToChan(ips []string, ipChan chan<- string) {
	for _, ip := range ips {
		ipChan <- ip
	}
	close(ipChan)
}

func launchScanWorkers(ctx context.Context, workers int, ipChan <-chan string, resChan chan<- SSHScanHost, port int, timeout time.Duration, enrolledMap map[string]string) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go dispatchScanWorker(ctx, ipChan, resChan, port, timeout, enrolledMap, &wg)
	}
	return &wg
}

func collectScanResults(resChan <-chan SSHScanHost) []SSHScanHost {
	var results []SSHScanHost
	for res := range resChan {
		results = append(results, res)
	}
	return results
}

func sortScanResultsByIP(results []SSHScanHost) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].IP < results[j].IP
	})
}

func runConcurrentScanner(ctx context.Context, ips []string, opts SSHScanOptions, enrolledMap map[string]string) []SSHScanHost {
	workers := resolveScanWorkers(opts.Workers)
	ipChan := make(chan string, len(ips))
	resChan := make(chan SSHScanHost, len(ips))
	go feedIPsToChan(ips, ipChan)
	wg := launchScanWorkers(ctx, workers, ipChan, resChan, opts.Port, opts.Timeout, enrolledMap)
	wg.Wait()
	close(resChan)
	results := collectScanResults(resChan)
	sortScanResultsByIP(results)
	return results
}

func formatEnrollmentStatus(host SSHScanHost) string {
	if host.IsEnrolled {
		return fmt.Sprintf("[ENROLLED: %s]", host.EnrolledAlias)
	}
	return "[NEW]"
}

func printScanTableHeader(w io.Writer) {
	fmt.Fprintln(w, "IP\tPORT\tSTATUS\tLATENCY\tHOSTNAME\tENROLLMENT\tSUGGESTION")
}

func printScanTableRow(w io.Writer, h SSHScanHost) {
	status := "ONLINE"
	enroll := formatEnrollmentStatus(h)
	fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n",
		h.IP, h.Port, status, h.Latency.String(), h.Hostname, enroll, h.Suggestion)
}

func printScanResultTable(out io.Writer, hosts []SSHScanHost, subnet string) error {
	if len(hosts) == 0 {
		fmt.Fprintf(out, "No active SSH machines discovered on subnet %s (port 22).\n", subnet)
		return nil
	}
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	printScanTableHeader(w)
	for _, h := range hosts {
		printScanTableRow(w, h)
	}
	flushErr := w.Flush()
	fmt.Fprintf(out, "\n  Total active SSH machines discovered: %d on subnet %s\n\n", len(hosts), subnet)
	return flushErr
}

func prepareScanOptions(opts *SSHScanOptions) {
	opts.Port = resolveScanPort(opts.Port)
	opts.Timeout = resolveScanTimeout(opts.Timeout)
}

// ExecuteSubnetScan executes a full subnet discovery scan and prints formatted table output.
func ExecuteSubnetScan(ctx context.Context, out io.Writer, opts SSHScanOptions) ([]SSHScanHost, *apperror.AppError) {
	cidr, err := normalizeSubnetCIDR(opts.Subnet)
	if err != nil {
		return nil, err
	}
	ips, err := expandSubnetIPs(cidr)
	if err != nil {
		return nil, err
	}
	prepareScanOptions(&opts)
	hosts := runConcurrentScanner(ctx, ips, opts, fetchEnrolledHostsMap(ctx))
	if flushErr := printScanResultTable(out, hosts, cidr); flushErr != nil {
		return nil, apperror.WrapSimple(flushErr, "ExecuteSubnetScan")
	}
	return hosts, nil
}
