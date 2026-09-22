package cmdos

import (
	"fmt"
	"net"
	"time"
)

func runDNSBenchmark() []DNSBenchmarkResult {
	var results []DNSBenchmarkResult
	for _, p := range KnownDNSProviders {
		res := benchmarkProvider(p)
		results = append(results, res)
	}

	return results
}

func benchmarkProvider(p DNSProvider) DNSBenchmarkResult {
	start := time.Now()
	conn, err := net.DialTimeout("udp", p.Primary+":53", 2*time.Second)
	latency := time.Since(start)

	if err != nil {
		return DNSBenchmarkResult{
			Provider: p,
			Latency:  latency,
			Success:  false,
		}
	}
	defer conn.Close()

	return DNSBenchmarkResult{
		Provider: p,
		Latency:  latency,
		Success:  true,
	}
}

func printDNSBenchmark(results []DNSBenchmarkResult) {
	fmt.Println("▶ DNS Resolver Latency Benchmark:")
	for _, r := range results {
		status := "FAIL"
		if r.Success {
			status = fmt.Sprintf("%v", r.Latency.Round(time.Millisecond))
		}
		fmt.Printf("  • %-12s (%-15s): %s\n", r.Provider.Name, r.Provider.Primary, status)
	}
}
