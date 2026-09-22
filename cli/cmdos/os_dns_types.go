package cmdos

import "time"

// DNSProvider holds nameserver IPs for a secure resolver.
type DNSProvider struct {
	Name      string `json:"name"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

// DNSBenchmarkResult models latency results for DNS servers.
type DNSBenchmarkResult struct {
	Provider DNSProvider   `json:"provider"`
	Latency  time.Duration `json:"latency"`
	Success  bool          `json:"success"`
}

var KnownDNSProviders = map[string]DNSProvider{
	"cloudflare": {
		Name:      "Cloudflare",
		Primary:   "1.1.1.1",
		Secondary: "1.0.0.1",
	},
	"google": {
		Name:      "Google",
		Primary:   "8.8.8.8",
		Secondary: "8.8.4.4",
	},
	"quad9": {
		Name:      "Quad9",
		Primary:   "9.9.9.9",
		Secondary: "149.112.112.112",
	},
	"adguard": {
		Name:      "AdGuard",
		Primary:   "94.140.14.14",
		Secondary: "94.140.15.15",
	},
}

// DNSOperator abstracts OS-level DNS configuration.
type DNSOperator interface {
	SetDNS(iface string, p DNSProvider) error
	SetDHCP(iface string) error
	GetDNS(iface string) ([]string, error)
	GetDefaultInterface() (string, error)
}
