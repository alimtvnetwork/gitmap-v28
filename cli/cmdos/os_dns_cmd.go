package cmdos

import (
	"fmt"
	"strings"
)

func runOSDNSCommand(args []string) error {
	if len(args) == 0 {
		return handleDNSStatus()
	}

	subCmd := strings.ToLower(args[0])
	switch subCmd {
	case "help", "-h", "--help":
		printDNSUsage()

		return nil
	case "status", "st", "info":
		return handleDNSStatus()
	case "set":
		return handleDNSSet(args[1:])
	case "dhcp", "auto", "reset":
		return handleDNSDHCP()
	case "bench", "benchmark", "ping":
		return handleDNSBenchmark()
	default:
		return handleDNSDirectSet(subCmd)
	}
}

func handleDNSStatus() error {
	engine := newPlatformDNSEngine()
	iface, _ := engine.GetDefaultInterface()
	servers, err := engine.GetDNS(iface)
	if err != nil {
		return err
	}

	fmt.Println("▶ Active DNS Configuration:")
	fmt.Printf("  • Interface:   %s\n", iface)
	if len(servers) == 0 {
		fmt.Println("  • DNS Servers: Automatic (DHCP)")

		return nil
	}

	fmt.Printf("  • DNS Servers: %s\n", strings.Join(servers, ", "))

	return nil
}

func handleDNSDHCP() error {
	engine := newPlatformDNSEngine()
	iface, err := engine.GetDefaultInterface()
	if err != nil {
		return err
	}

	if err := engine.SetDHCP(iface); err != nil {
		return err
	}

	fmt.Printf("✔ DNS reverted to automatic DHCP on interface %s\n", iface)

	return nil
}

func handleDNSBenchmark() error {
	results := runDNSBenchmark()
	printDNSBenchmark(results)

	return nil
}
