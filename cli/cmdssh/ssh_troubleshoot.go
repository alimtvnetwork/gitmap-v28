package cmdssh

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// SSHDiagnosisType classifies the root-cause diagnosis of SSH connectivity probes.
type SSHDiagnosisType string

const (
	DiagFirewallInboundDrop SSHDiagnosisType = "FIREWALL INBOUND DROP"
	DiagDaemonNotRunning    SSHDiagnosisType = "SSH DAEMON NOT RUNNING"
	DiagHostOffline         SSHDiagnosisType = "HOST UNREACHABLE / OFFLINE"
	DiagSSHPortOpen         SSHDiagnosisType = "SSH PORT OPEN & ACCESSIBLE"
	DiagUnknownError        SSHDiagnosisType = "NETWORK PROBE ERROR"
)

// String returns the string representation.
func (d SSHDiagnosisType) String() string {
	return string(d)
}

// IsOpen reports whether SSH port is open and accessible.
func (d SSHDiagnosisType) IsOpen() bool {
	return d == DiagSSHPortOpen
}

type tcpProbeFunc func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, error)
type pingProbeFunc func(ctx context.Context, ip string, timeout time.Duration) bool

var defaultTCPProber tcpProbeFunc = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, error) {
	dialer := net.Dialer{Timeout: timeout}
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, err
	}

	_ = conn.Close()

	return true, nil
}

var defaultPingProber pingProbeFunc = func(ctx context.Context, ip string, timeout time.Duration) bool {
	cmd := buildPingCommand(ctx, ip, timeout)
	err := cmd.Run()

	return err == nil
}

func buildPingCommand(ctx context.Context, ip string, timeout time.Duration) *exec.Cmd {
	if constants.IsWindowsPlatform() {
		ms := strconv.Itoa(int(timeout / time.Millisecond))

		return exec.CommandContext(ctx, "ping", "-n", "1", "-w", ms, ip)
	}

	sec := resolvePingSeconds(timeout)

	return exec.CommandContext(ctx, "ping", "-c", "1", "-W", sec, ip)
}

func resolvePingSeconds(timeout time.Duration) string {
	secInt := int(timeout / time.Second)
	if secInt <= 0 {
		return "1"
	}

	return strconv.Itoa(secInt)
}

type troubleshootReport struct {
	target       string
	ip           string
	port         int
	isSSHOpen    bool
	isTimeout    bool
	isRefused    bool
	isPingOk     bool
	isWinRMOpen  bool
	isRDPOpen    bool
	diagnosis    SSHDiagnosisType
	probeErrText string
}

func splitHostPortPair(clean string) (string, int, bool) {
	if !strings.Contains(clean, ":") {
		return "", 0, false
	}

	host, portStr, splitErr := net.SplitHostPort(clean)
	if splitErr != nil {
		return "", 0, false
	}

	p, convErr := strconv.Atoi(portStr)
	if convErr != nil || p <= 0 {
		return "", 0, false
	}

	return host, p, true
}

func parseTroubleshootTarget(raw string) (string, int) {
	clean := strings.TrimSpace(raw)
	target, err := ParseSSHTarget(clean, "", 22)
	if err == nil && target != nil {
		return target.IP, target.Port
	}

	host, port, isSplit := splitHostPortPair(clean)
	if isSplit {
		return host, port
	}

	return clean, 22
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "timeout") || strings.Contains(msg, "timed out") || strings.Contains(msg, "i/o timeout")
}

func isRefusedError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "refused") || strings.Contains(msg, "wsaeconnrefused")
}

func classifyTroubleshootDiagnosis(isSSHOpen, isTimeout, isRefused, isPingOk bool) SSHDiagnosisType {
	if isSSHOpen {
		return DiagSSHPortOpen
	}

	if isRefused {
		return DiagDaemonNotRunning
	}

	if isTimeout && isPingOk {
		return DiagFirewallInboundDrop
	}

	if isTimeout && !isPingOk {
		return DiagHostOffline
	}

	return DiagUnknownError
}

func probeRemoteManagement(ctx context.Context, ip string, rep *troubleshootReport) {
	quickTimeout := 1200 * time.Millisecond
	winrmHttp, _ := defaultTCPProber(ctx, ip, 5985, quickTimeout)
	winrmHttps, _ := defaultTCPProber(ctx, ip, 5986, quickTimeout)
	rep.isWinRMOpen = winrmHttp || winrmHttps

	rdp, _ := defaultTCPProber(ctx, ip, 3389, quickTimeout)
	rep.isRDPOpen = rdp
}

func performTroubleshootProbe(ctx context.Context, rawTarget string) troubleshootReport {
	ip, port := parseTroubleshootTarget(rawTarget)
	rep := troubleshootReport{
		target: rawTarget,
		ip:     ip,
		port:   port,
	}

	sshOpen, sshErr := defaultTCPProber(ctx, ip, port, 2000*time.Millisecond)
	rep.isSSHOpen = sshOpen
	rep.isTimeout = isTimeoutError(sshErr)
	rep.isRefused = isRefusedError(sshErr)
	if sshErr != nil {
		rep.probeErrText = sshErr.Error()
	}

	if rep.isTimeout {
		rep.isPingOk = defaultPingProber(ctx, ip, 2000*time.Millisecond)
	}

	rep.diagnosis = classifyTroubleshootDiagnosis(rep.isSSHOpen, rep.isTimeout, rep.isRefused, rep.isPingOk)
	probeRemoteManagement(ctx, ip, &rep)

	return rep
}

func printDiagnosisExplanation(rep troubleshootReport) {
	fmt.Printf("  %s● Problem Diagnosis:%s\n", constants.ColorBold, constants.ColorReset)
	switch rep.diagnosis {
	case DiagFirewallInboundDrop:
		fmt.Printf("    Status: %s[%s]%s\n", constants.ColorYellow, rep.diagnosis, constants.ColorReset)
		fmt.Println("    Reason: Host responds to ICMP ping, but dropped TCP port 22 connection packets.")
		fmt.Println("    Cause:  Machine is ONLINE, but inbound TCP port 22 is blocked by Windows Firewall or iptables/ufw.")
	case DiagDaemonNotRunning:
		fmt.Printf("    Status: %s[%s]%s\n", constants.ColorRed, rep.diagnosis, constants.ColorReset)
		fmt.Println("    Reason: Target host actively refused TCP connection on port 22.")
		fmt.Println("    Cause:  Host and network are alive, but OpenSSH Server daemon (sshd) is NOT running or listening.")
	case DiagHostOffline:
		fmt.Printf("    Status: %s[%s]%s\n", constants.ColorRed, rep.diagnosis, constants.ColorReset)
		fmt.Println("    Reason: Target host did not respond to ICMP ping or TCP port 22 within 2 seconds.")
		fmt.Println("    Cause:  Machine may be powered off, suspended, disconnected, or on an unreachable subnet.")
	case DiagSSHPortOpen:
		fmt.Printf("    Status: %s[%s]%s\n", constants.ColorGreen, rep.diagnosis, constants.ColorReset)
		fmt.Println("    Reason: TCP port 22 is open and listening for incoming connections.")
		fmt.Println("    Cause:  OpenSSH Server is reachable. Check credentials, keys, or SSH config if login fails.")
	default:
		fmt.Printf("    Status: %s[%s]%s\n", constants.ColorYellow, rep.diagnosis, constants.ColorReset)
		fmt.Printf("    Reason: %s\n", rep.probeErrText)
	}
	fmt.Println()
}

func formatProbeStatus(isOpen bool, openText, closedText, openColor, closedColor string) string {
	if isOpen {
		return fmt.Sprintf("%s%s%s", openColor, openText, constants.ColorReset)
	}

	return fmt.Sprintf("%s%s%s", closedColor, closedText, constants.ColorReset)
}
