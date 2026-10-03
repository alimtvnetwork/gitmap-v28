package cmdssh

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	DiagFirewallInboundDrop = "FIREWALL INBOUND DROP"
	DiagDaemonNotRunning    = "SSH DAEMON NOT RUNNING"
	DiagHostOffline         = "HOST UNREACHABLE / OFFLINE"
	DiagSSHPortOpen         = "SSH PORT OPEN & ACCESSIBLE"
	DiagUnknownError        = "NETWORK PROBE ERROR"
)

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
	if runtime.GOOS == "windows" {
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
	diagnosis    string
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

func classifyTroubleshootDiagnosis(isSSHOpen, isTimeout, isRefused, isPingOk bool) string {
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

func printReachabilityMatrix(rep troubleshootReport) {
	sshStatus := formatProbeStatus(rep.isSSHOpen, "OPEN / RESPONDING", "UNREACHABLE", constants.ColorGreen, constants.ColorRed)
	if rep.isRefused {
		sshStatus = fmt.Sprintf("%sREFUSED (Daemon Down)%s", constants.ColorRed, constants.ColorReset)
	}
	if rep.isTimeout {
		sshStatus = fmt.Sprintf("%sTIMED OUT (Dropped)%s", constants.ColorYellow, constants.ColorReset)
	}

	pingStatus := formatProbeStatus(rep.isPingOk, "RESPONDED (Host Up)", "TIMED OUT / NO REPLY", constants.ColorGreen, constants.ColorDim)
	winrmStatus := formatProbeStatus(rep.isWinRMOpen, "OPEN (Ports 5985/5986)", "CLOSED", constants.ColorGreen, constants.ColorDim)
	rdpStatus := formatProbeStatus(rep.isRDPOpen, "OPEN (Port 3389)", "CLOSED", constants.ColorGreen, constants.ColorDim)

	fmt.Printf("  %s● Reachability & Port Matrix:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("    • Port %d (SSH):       %s\n", rep.port, sshStatus)
	fmt.Printf("    • ICMP Ping:           %s\n", pingStatus)
	fmt.Printf("    • WinRM (5985/5986):   %s\n", winrmStatus)
	fmt.Printf("    • RDP (3389):          %s\n", rdpStatus)
	fmt.Println()
}

func printTargetCommands(rep troubleshootReport) {
	fmt.Printf("  %s● Exact Commands to Run on Target Machine:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("    1. Install & start SSH:     %sgitmap ssh enable --port %d%s\n", constants.ColorCyan, rep.port, constants.ColorReset)
	fmt.Printf("    2. Check port listeners:    %sgitmap ports%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    3. Native Windows (Admin):  Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0; Start-Service sshd")
	fmt.Println("    4. Native Linux (Ubuntu):   sudo apt-get install -y openssh-server && sudo systemctl enable --now ssh")
	fmt.Println()
}

func printAlternativeOptions(rep troubleshootReport) {
	fmt.Printf("  %s● Alternative Remote Connection Options:%s\n", constants.ColorBold, constants.ColorReset)
	if rep.isWinRMOpen {
		fmt.Printf("    • PowerShell Remoting:  %sEnter-PSSession -ComputerName %s%s\n", constants.ColorGreen, rep.ip, constants.ColorReset)
		fmt.Printf("      Remote enable SSH:    %sInvoke-Command -ComputerName %s -ScriptBlock { gitmap ssh enable }%s\n", constants.ColorGreen, rep.ip, constants.ColorReset)
	}

	if rep.isRDPOpen {
		fmt.Printf("    • Remote Desktop (RDP): %smstsc /v:%s%s\n", constants.ColorGreen, rep.ip, constants.ColorReset)
	}

	if !rep.isWinRMOpen && !rep.isRDPOpen {
		fmt.Printf("    • No alternate remote services detected. Use local hypervisor/console or VM session.\n")
	}

	fmt.Println()
}

func printNextStepsChecklist(rep troubleshootReport) {
	fmt.Printf("  %s● Next Steps Checklist:%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("    [ ] Confirm machine is powered ON and has assigned IP %s%s%s.\n", constants.ColorCyan, rep.ip, constants.ColorReset)
	fmt.Printf("    [ ] Run '%sgitmap ssh enable%s' on target machine with elevated privileges.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    [ ] Ensure firewall rule allows inbound TCP port %d.\n", rep.port)
	fmt.Printf("    [ ] Retry: '%sgitmap ssh troubleshoot %s%s' or '%sgitmap ssh %s%s'\n", constants.ColorYellow, rep.target, constants.ColorReset, constants.ColorYellow, rep.target, constants.ColorReset)
	fmt.Println()
}

func renderTroubleshootHUD(rep troubleshootReport) {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s║             GITMAP SSH TROUBLESHOOTING & RECOVERY HUD            ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Target Target: %s%s%s (IP: %s, Port: %d)\n\n", constants.ColorWhite, rep.target, constants.ColorReset, rep.ip, rep.port)

	printReachabilityMatrix(rep)
	printDiagnosisExplanation(rep)
	printTargetCommands(rep)
	printAlternativeOptions(rep)
	printNextStepsChecklist(rep)
}

func printTroubleshootHelp() {
	fmt.Printf("\n%sUsage:%s gitmap ssh troubleshoot <target> [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Diagnoses SSH connectivity failures, classifies firewall drops vs daemon states, and provides recovery steps.")
	fmt.Println("\nAliases: doctor, diagnose")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap ssh troubleshoot 192.168.1.50")
	fmt.Println("  gitmap ssh doctor admin@192.168.1.50:2222")
	fmt.Println("  gitmap ssh diagnose my-server")
	fmt.Println()
}

func runSSHTroubleshootCLI(ctx context.Context, args []string) error {
	if len(args) == 0 {
		printTroubleshootHelp()
		return apperror.NewValidationError("missing required target host or IP address")
	}

	first := args[0]
	if first == "--help" || first == "-h" || first == "help" {
		printTroubleshootHelp()
		return nil
	}

	rep := performTroubleshootProbe(ctx, first)
	renderTroubleshootHUD(rep)

	return nil
}
