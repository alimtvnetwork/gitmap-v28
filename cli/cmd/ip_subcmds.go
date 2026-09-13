package cmd

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/netip"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func dispatchIPSubcommand(ctx context.Context, args []string) error {
	sub := strings.ToLower(args[0])
	mgr := createCMDNetIPManager()

	switch sub {
	case "help", "-h", "--help":
		printIPUsage()
		return nil
	case "show", "get", "list", "ls", "status", "st":
		return runIPShow(ctx, mgr, args[1:])
	case "set", "apply":
		return runIPSet(ctx, mgr, args[1:])
	case "change", "edit", "modify":
		return runIPChange(ctx, mgr, args[1:])
	case "switch", "toggle":
		return runIPSwitch(ctx, mgr, args[1:])
	case "revert", "rollback", "undo":
		return runIPRevert(ctx, mgr, args[1:])
	default:
		if looksLikeIPv4Addr(args[0]) {
			return runIPSet(ctx, mgr, args)
		}
		printIPUsage()
		return apperror.NewSimple("unknown ip subcommand: "+sub, "E_INVALID_IP_SUBCMD")
	}
}

func createCMDNetIPManager() *netip.Manager {
	driver := netip.DetectDriver()
	var db *store.DB
	if rootDb, err := store.OpenRootDB(); err == nil {
		db = rootDb
	}
	return netip.NewManager(driver, db)
}

func runIPShow(ctx context.Context, mgr *netip.Manager, args []string) error {
	res := mgr.ListInterfaces(ctx)
	if res.IsFailure() {
		return res.AppError()
	}
	printCMDInterfaceList(res.Value)
	return nil
}

func printCMDInterfaceList(ifaces []netip.InterfaceInfo) {
	fmt.Printf("%-18s %-16s %-16s %-16s %-6s %-6s\n", "INTERFACE", "IP ADDRESS", "NETMASK", "GATEWAY", "DHCP", "STATUS")
	fmt.Println(strings.Repeat("-", 82))
	for _, iface := range ifaces {
		dhcpStr := "no"
		if iface.IsDHCP {
			dhcpStr = "yes"
		}
		statusStr := "down"
		if iface.IsUp {
			statusStr = "up"
		}
		fmt.Printf("%-18s %-16s %-16s %-16s %-6s %-6s\n",
			iface.Name, iface.IP, iface.Netmask, iface.Gateway, dhcpStr, statusStr)
	}
}

func runIPSet(ctx context.Context, mgr *netip.Manager, args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing IP address for set", "E_MISSING_ARG")
	}
	opts, valOpts := parseCMDChangeOptions(args, false)
	res := mgr.ChangeIP(ctx, opts, valOpts)
	return handleCMDRollbackResult(res)
}

func runIPChange(ctx context.Context, mgr *netip.Manager, args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing IP address for change", "E_MISSING_ARG")
	}
	opts, valOpts := parseCMDChangeOptions(args, false)
	res := mgr.ChangeIP(ctx, opts, valOpts)
	return handleCMDRollbackResult(res)
}

func runIPSwitch(ctx context.Context, mgr *netip.Manager, args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap ip switch <dhcp|static> [flags]", "E_MISSING_ARG")
	}
	mode := strings.ToLower(args[0])
	isDHCP := mode == "dhcp"
	opts, valOpts := parseCMDChangeOptions(args[1:], isDHCP)
	opts.IsDHCP = isDHCP
	res := mgr.ChangeIP(ctx, opts, valOpts)
	return handleCMDRollbackResult(res)
}

func runIPRevert(ctx context.Context, mgr *netip.Manager, args []string) error {
	iface := resolveCMDTargetInterface(args)
	res := mgr.RevertIP(ctx, iface)
	return handleCMDRollbackResult(res)
}

func handleCMDRollbackResult(res netip.RollbackResultWrap) error {
	if res.IsFailure() {
		return res.AppError()
	}
	val := res.Value
	if val.IsReverted {
		fmt.Printf("⚠️  Network change failed and was rolled back: %s\n", val.Message)
		return apperror.NewSimple(val.Message, "E_IP_REVERTED")
	}
	fmt.Printf("✔ %s\n", val.Message)
	return nil
}

func parseCMDChangeOptions(args []string, isDHCP bool) (netip.ChangeOptions, netip.ValidationOptions) {
	ipVal := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		ipVal = args[0]
	}
	gw := extractCMDArgFlag(args, "--gateway")
	if gw == "" && len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		gw = args[1]
	}
	mask := extractCMDArgFlagWithDefault(args, "--mask", "255.255.255.0")
	iface := resolveCMDTargetInterface(args)
	dnsStr := extractCMDArgFlagWithDefault(args, "--dns", "8.8.8.8,1.1.1.1")
	var dnsList []string
	if dnsStr != "" {
		dnsList = strings.Split(dnsStr, ",")
	}
	opts := netip.ChangeOptions{
		InterfaceName: iface,
		IP:            ipVal,
		Netmask:       mask,
		Gateway:       gw,
		DNS:           dnsList,
		IsDHCP:        isDHCP,
		IsDryRun:      hasCMDFlag(args, "--dry-run") || hasCMDFlag(args, "-n"),
	}
	valOpts := netip.ValidationOptions{
		IsValidationActive: !hasCMDFlag(args, "--no-validate"),
		TargetHost:         extractCMDArgFlagWithDefault(args, "--ping-target", "8.8.8.8"),
		Count:              3,
		TimeoutSec:         2,
	}
	return opts, valOpts
}

func resolveCMDTargetInterface(args []string) string {
	iface := extractCMDArgFlag(args, "--interface")
	if iface != "" {
		return iface
	}
	iface = extractCMDArgFlag(args, "-i")
	if iface != "" {
		return iface
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "eth0"
	}
	for _, i := range ifaces {
		isLoop := (i.Flags & net.FlagLoopback) != 0
		isUp := (i.Flags & net.FlagUp) != 0
		if !isLoop && isUp {
			return i.Name
		}
	}
	return "eth0"
}

func extractCMDArgFlag(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func extractCMDArgFlagWithDefault(args []string, flag, defVal string) string {
	val := extractCMDArgFlag(args, flag)
	if val != "" {
		return val
	}
	return defVal
}

func hasCMDFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func looksLikeIPv4Addr(str string) bool {
	parsed := net.ParseIP(str)
	return parsed != nil && parsed.To4() != nil
}

func printIPUsage() {
	fmt.Println("Usage: gitmap ip [command] [arguments] [flags]")
	fmt.Println()
	fmt.Println("When called without arguments, gitmap ip prints the local primary IPv4 address.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  show [iface]                         Display network interfaces, IPs, and status")
	fmt.Println("  set <ip> [gateway] [flags]           Configure static IP with ping validation & rollback")
	fmt.Println("  change <ip> [gateway] [flags]        Guided IP change workflow")
	fmt.Println("  switch <dhcp|static> [flags]         Switch configuration mode between DHCP and Static")
	fmt.Println("  revert [iface]                       Restore previous IP configuration from snapshot")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -i, --interface <name>               Target network interface (auto-detected if omitted)")
	fmt.Println("  -g, --gateway <ip>                   Default gateway IPv4 address")
	fmt.Println("  -m, --mask <mask>                    Subnet mask (default: 255.255.255.0)")
	fmt.Println("      --dns <ip,ip>                    DNS servers (default: 8.8.8.8,1.1.1.1)")
	fmt.Println("      --no-validate                    Bypass connectivity ping verification")
	fmt.Println("  -n, --dry-run                        Simulate changes without applying")
}
