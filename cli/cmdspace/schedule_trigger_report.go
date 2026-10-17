// schedule_trigger_report.go holds the machine-identity collection and
// the HTML notification report builder for email task triggers. Trigger
// storage, registration, and firing live in schedule_trigger.go.
package cmdspace

import (
	"html"
	"net"
	"os"
	"os/user"
	"strings"
)

const (
	// outboundProbeAddr is dialed over UDP to discover the outbound IP;
	// connect() on UDP sends no traffic.
	outboundProbeAddr = "8.8.8.8:80"
	// bannerOKColor and bannerFailColor color the email status banner.
	bannerOKColor   = "#1a7f37"
	bannerFailColor = "#d1242f"
)

// MachineInfo carries host identity rows for the trigger email body.
type MachineInfo struct {
	Host string
	User string
	IP   string
}

// CollectMachineInfo gathers host identity rows for trigger emails.
// It never errors: every field falls back to "unknown".
func CollectMachineInfo() MachineInfo {
	return MachineInfo{
		Host: machineHostname(),
		User: machineUsername(),
		IP:   machineOutboundIP(),
	}
}

// machineHostname returns os.Hostname() or "unknown".
func machineHostname() string {
	host, err := os.Hostname()
	if err == nil && len(strings.TrimSpace(host)) > 0 {
		return host
	}

	return unknownMachineValue
}

// machineUsername returns the OS login name or "unknown".
func machineUsername() string {
	current, err := user.Current()
	if err == nil && len(strings.TrimSpace(current.Username)) > 0 {
		return current.Username
	}

	envUser := strings.TrimSpace(os.Getenv("USER"))
	if len(envUser) > 0 {
		return envUser
	}

	return unknownMachineValue
}

// machineOutboundIP discovers the outbound IP via a UDP dial (no traffic
// sent), then the first non-loopback interface address, then "unknown".
func machineOutboundIP() string {
	if ip := outboundIPViaDial(); ip != "" {
		return ip
	}

	if ip := outboundIPViaInterfaces(); ip != "" {
		return ip
	}

	return unknownMachineValue
}

// outboundIPViaDial returns the local address of a UDP socket dialed at
// the probe address. Dialing UDP only resolves the route — no packet
// leaves the machine.
func outboundIPViaDial() string {
	conn, err := net.Dial("udp", outboundProbeAddr)
	if err != nil {
		return ""
	}
	defer conn.Close()

	udpAddr, isUDP := conn.LocalAddr().(*net.UDPAddr)
	if isUDP {
		return udpAddr.IP.String()
	}

	return ""
}

// outboundIPViaInterfaces returns the first non-loopback IPv4 address.
func outboundIPViaInterfaces() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if ip := firstUsableInterfaceIP(iface); ip != "" {
			return ip
		}
	}

	return ""
}

// firstUsableInterfaceIP returns the first non-loopback IPv4 on an
// interface, or "".
func firstUsableInterfaceIP(iface net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ip := usableIPFromAddr(addr); ip != "" {
			return ip
		}
	}

	return ""
}

// usableIPFromAddr extracts a non-loopback IPv4 from a net.Addr.
func usableIPFromAddr(addr net.Addr) string {
	ip := ipFromAddr(addr)
	if ip == nil || ip.IsLoopback() {
		return ""
	}

	flat := ip.To4()
	if flat == nil {
		return ""
	}

	return flat.String()
}

// ipFromAddr unwraps the IP from the two common net.Addr shapes.
func ipFromAddr(addr net.Addr) net.IP {
	switch typed := addr.(type) {
	case *net.IPNet:
		return typed.IP
	case *net.IPAddr:
		return typed.IP
	default:
		return nil
	}
}

// BuildTriggerEmailHTML renders the notification email: a green/red
// status banner plus a table of machine-info and outcome rows. Every
// interpolated value is HTML-escaped.
func BuildTriggerEmailHTML(task EmailTaskType, outcome TaskOutcome, info MachineInfo) string {
	color, status := triggerBannerStyle(outcome.IsSuccess)

	var sb strings.Builder
	sb.WriteString("<html><body style=\"font-family:sans-serif;max-width:640px;margin:0 auto;\">")
	sb.WriteString("<h2>gitmap " + html.EscapeString(string(task)) + " — " + html.EscapeString(status) + "</h2>")
	sb.WriteString("<div style=\"background:" + color + ";color:#ffffff;padding:12px 16px;border-radius:6px;\">")
	sb.WriteString("<strong>" + html.EscapeString(status) + "</strong> — gitmap " + html.EscapeString(outcome.Operation))
	sb.WriteString("</div>")
	sb.WriteString("<table style=\"margin-top:12px;border-collapse:collapse;width:100%;\">")
	writeTriggerRow(&sb, "Host", info.Host)
	writeTriggerRow(&sb, "User", info.User)
	writeTriggerRow(&sb, "Outbound IP", info.IP)
	writeTriggerRow(&sb, "Repo", outcome.Repo)
	writeTriggerRow(&sb, "Branch", outcome.Branch)
	writeTriggerRow(&sb, "Operation", outcome.Operation)
	writeTriggerRow(&sb, "Outcome", outcome.Summary)
	writeTriggerRow(&sb, "Started", outcome.StartedAt)
	writeTriggerRow(&sb, "Finished", outcome.FinishedAt)
	sb.WriteString("</table></body></html>")

	return sb.String()
}

// triggerBannerStyle picks the banner color and label per outcome.
func triggerBannerStyle(isSuccess bool) (color, label string) {
	if isSuccess {
		return bannerOKColor, statusOKText
	}

	return bannerFailColor, statusFailText
}

// writeTriggerRow appends one escaped label/value row to the email table.
func writeTriggerRow(sb *strings.Builder, label, value string) {
	sb.WriteString("<tr><td style=\"padding:6px 10px;border:1px solid #dddddd;font-weight:bold;\">")
	sb.WriteString(html.EscapeString(label))
	sb.WriteString("</td><td style=\"padding:6px 10px;border:1px solid #dddddd;\">")
	sb.WriteString(html.EscapeString(value))
	sb.WriteString("</td></tr>")
}
