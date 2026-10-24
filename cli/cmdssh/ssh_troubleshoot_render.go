package cmdssh

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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

// PrintSevenStepTroubleshootGuide renders the 7-step remediation workflow when SSH connection fails.
func PrintSevenStepTroubleshootGuide(target string, err error) {
	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("%s║           SSH CONNECTION TROUBLESHOOTING & REMEDIATION           ║%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorYellow, constants.ColorReset)
	if err != nil {
		fmt.Printf("  %sFailure Cause:%s %v\n\n", constants.ColorRed, constants.ColorReset, err)
	}

	fmt.Printf("  %s7-Step Remediation Checklist for '%s':%s\n", constants.ColorBold, target, constants.ColorReset)
	fmt.Printf("    1. %sCheck Host Power & Ping:%s Confirm the machine is online.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → ping %s\n", target)
	fmt.Printf("    2. %sEnable OpenSSH Server:%s Install Gitmap or enable sshd on the target.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → Run 'gitmap ssh enable' (or 'sudo systemctl enable --now ssh') on target\n")
	fmt.Printf("    3. %sCheck Target SSH Port:%s Verify what port sshd is listening on.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → Run 'gitmap ssh port ls' on target machine (default: 22)\n")
	fmt.Printf("    4. %sVerify Firewall Rules:%s Ensure inbound TCP port is allowed through the firewall.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → Run 'gitmap ssh port add <port>' or 'gitmap ssh enable-public <port>'\n")
	fmt.Printf("    5. %sDeploy Authentication Key:%s Push this machine's public key to the target.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → gitmap ssh copy-id %s (or gitmap ssh view to inspect key)\n", target)
	fmt.Printf("    6. %sRun Deep Diagnostics:%s Execute automated reachability and port probe.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → gitmap ssh troubleshoot %s\n", target)
	fmt.Printf("    7. %sManual Verbose Connection:%s Test with OpenSSH client directly.\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("       → ssh -vvv %s\n\n", target)
}
