package cmdos

import "fmt"

func printDNSUsage() {
	fmt.Println("Usage: gitmap os dns <subcommand>")
	fmt.Println("       gitmap os dns status                           (Inspect current DNS servers)")
	fmt.Println("       gitmap os dns set <cloudflare|google|quad9|adguard>")
	fmt.Println("       gitmap os dns <cloudflare|google|quad9|adguard> (Direct provider shortcut)")
	fmt.Println("       gitmap os dns dhcp                             (Revert to automatic DHCP DNS)")
	fmt.Println("       gitmap os dns benchmark                        (Test latency of DNS providers)")
}
