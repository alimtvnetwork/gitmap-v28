package cmdos

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/netip"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runOSIP(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) {
		printOSIPUsage()
		return nil
	}

	sub := strings.ToLower(args[0])
	ctx := context.Background()
	mgr := createNetIPManager()

	switch sub {
	case "show", "get", "list", "ls", "status", "st":
		return handleOSIPShow(ctx, mgr, args[1:])
	case "set", "apply":
		return handleOSIPSet(ctx, mgr, args[1:])
	case "change", "edit", "modify":
		return handleOSIPChange(ctx, mgr, args[1:])
	case "switch", "toggle":
		return handleOSIPSwitch(ctx, mgr, args[1:])
	case "revert", "rollback", "undo":
		return handleOSIPRevert(ctx, mgr, args[1:])
	default:
		if looksLikeIPv4(args[0]) {
			return handleOSIPSet(ctx, mgr, args)
		}
		printOSIPUsage()
		return apperror.NewSimple("unknown os ip subcommand: "+sub, "E_INVALID_IP_SUBCMD")
	}
}

func createNetIPManager() *netip.Manager {
	driver := netip.DetectDriver()
	var db *store.DB
	if rootDb, err := store.OpenRootDB(); err == nil {
		db = rootDb
	}
	return netip.NewManager(driver, db)
}

func handleOSIPShow(ctx context.Context, mgr *netip.Manager, args []string) error {
	res := mgr.ListInterfaces(ctx)
	if res.IsFailure() {
		return res.AppError()
	}
	printInterfaceList(res.Value)
	return nil
}

func printInterfaceList(ifaces []netip.InterfaceInfo) {
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

func handleOSIPSet(ctx context.Context, mgr *netip.Manager, args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing IP address for set", "E_MISSING_ARG")
	}
	opts, valOpts := parseChangeOptions(args, false)
	res := mgr.ChangeIP(ctx, opts, valOpts)
	return handleRollbackResult(res)
}

func handleOSIPChange(ctx context.Context, mgr *netip.Manager, args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("missing IP address for change", "E_MISSING_ARG")
	}
	opts, valOpts := parseChangeOptions(args, false)
	res := mgr.ChangeIP(ctx, opts, valOpts)
	return handleRollbackResult(res)
}

func handleOSIPSwitch(ctx context.Context, mgr *netip.Manager, args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap os ip switch <dhcp|static> [flags]", "E_MISSING_ARG")
	}
	mode := strings.ToLower(args[0])
	isDHCP := mode == "dhcp"
	opts, valOpts := parseChangeOptions(args[1:], isDHCP)
	opts.IsDHCP = isDHCP
	res := mgr.ChangeIP(ctx, opts, valOpts)
	return handleRollbackResult(res)
}

func handleOSIPRevert(ctx context.Context, mgr *netip.Manager, args []string) error {
	iface := resolveTargetInterface(args)
	res := mgr.RevertIP(ctx, iface)
	return handleRollbackResult(res)
}

func handleRollbackResult(res netip.RollbackResultWrap) error {
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

func parseChangeOptions(args []string, isDHCP bool) (netip.ChangeOptions, netip.ValidationOptions) {
	ipVal := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		ipVal = args[0]
	}
	gw := extractArgFlag(args, "--gateway")
	if gw == "" && len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		gw = args[1]
	}
	mask := extractArgFlagWithDefault(args, "--mask", "255.255.255.0")
	iface := resolveTargetInterface(args)
	dnsStr := extractArgFlagWithDefault(args, "--dns", "8.8.8.8,1.1.1.1")
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
		IsDryRun:      hasFlag(args, "--dry-run") || hasFlag(args, "-n"),
	}
	valOpts := netip.ValidationOptions{
		IsValidationActive: !hasFlag(args, "--no-validate"),
		TargetHost:         extractArgFlagWithDefault(args, "--ping-target", "8.8.8.8"),
		Count:              3,
		TimeoutSec:         2,
	}
	return opts, valOpts
}

func resolveTargetInterface(args []string) string {
	iface := extractArgFlag(args, "--interface")
	if iface != "" {
		return iface
	}
	iface = extractArgFlag(args, "-i")
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

func looksLikeIPv4(str string) bool {
	parsed := net.ParseIP(str)
	return parsed != nil && parsed.To4() != nil
}

func printOSIPUsage() {
	fmt.Println("Usage: gitmap os ip <command> [arguments] [flags]")
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
