package cmdports

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

var commonPortsList = []int{22, 80, 443, 3389, 5985, 5986, 8080}

var serviceNameToPortMap = map[string]int{
	"ssh":         22,
	"sshd":        22,
	"rdp":         3389,
	"http":        80,
	"web":         80,
	"https":       443,
	"ssl":         443,
	"winrm":       5985,
	"winrm-https": 5986,
	"winrm-ssl":   5986,
	"dns":         53,
	"smb":         445,
	"mysql":       3306,
	"postgres":    5432,
	"postgresql":  5432,
	"redis":       6379,
}

// Run executes the ports inspection CLI command.
func Run(args []string) error {
	opts, isHelp, err := parsePortsFlags(args)

	if err != nil {
		return err
	}

	if isHelp {
		printPortsHelp()

		return nil
	}

	entries, err := inspectPorts(opts)

	if err != nil {
		return apperror.WrapSimple(err, "ports inspection failed:")
	}

	if opts.JSONOutput {
		return outputPortsJSON(entries)
	}

	renderPortsTable(entries, opts)

	return nil
}

func parsePortsFlags(args []string) (PortsOptions, bool, error) {
	opts := PortsOptions{}
	i := 0

	for i < len(args) {
		arg := args[i]

		if isHelpFlag(arg) {
			return opts, true, nil
		}

		if isCommonFlag(arg) {
			opts.CommonOnly = true
			i++
			continue
		}

		if isFirewallFlag(arg) {
			opts.CommonOnly = true
			opts.FirewallOnly = true
			i++
			continue
		}

		if isJSONFlag(arg) {
			opts.JSONOutput = true
			i++
			continue
		}

		consumed, err := handlePortFlag(arg, args, i, &opts)

		if err != nil {
			return opts, false, err
		}

		if consumed > 0 {
			i += consumed
			continue
		}

		portNum, isInt := parsePortNumber(arg)

		if isInt {
			opts.TargetPort = portNum
			i++
			continue
		}

		servicePort, isService := resolveServiceNamePort(arg)

		if isService {
			opts.TargetPort = servicePort
			i++
			continue
		}

		return opts, false, apperror.NewValidationError(fmt.Sprintf("unknown ports argument: %s", arg))
	}

	return opts, false, nil
}

func isHelpFlag(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func isCommonFlag(arg string) bool {
	return arg == "-c" || arg == "--common" || arg == "common"
}

func isFirewallFlag(arg string) bool {
	norm := strings.ToLower(strings.TrimSpace(arg))

	return norm == "firewall" || norm == "fw" || norm == "--firewall"
}

func isJSONFlag(arg string) bool {
	return arg == "--json" || arg == "json"
}

func resolveServiceNamePort(name string) (int, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	port, isFound := serviceNameToPortMap[key]

	return port, isFound
}

func parsePortOrService(arg string) (int, bool) {
	portNum, isInt := parsePortNumber(arg)

	if isInt {
		return portNum, true
	}

	servicePort, isService := resolveServiceNamePort(arg)

	if isService {
		return servicePort, true
	}

	return 0, false
}

func handlePortFlag(arg string, args []string, index int, opts *PortsOptions) (int, error) {
	if strings.HasPrefix(arg, "--port=") {
		return parseEqualPort(arg, "--port=", opts)
	}

	if strings.HasPrefix(arg, "-p=") {
		return parseEqualPort(arg, "-p=", opts)
	}

	if arg == "--port" || arg == "-p" {
		return parseNextArgPort(args, index, opts)
	}

	return 0, nil
}

func parseEqualPort(arg string, prefix string, opts *PortsOptions) (int, error) {
	val := strings.TrimPrefix(arg, prefix)
	port, isValid := parsePortOrService(val)

	if !isValid || port < 1 || port > 65535 {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid port number: %s", val))
	}

	opts.TargetPort = port

	return 1, nil
}

func parseNextArgPort(args []string, index int, opts *PortsOptions) (int, error) {
	nextIndex := index + 1

	if nextIndex >= len(args) {
		return 0, apperror.NewValidationError("missing port value for flag")
	}

	val := args[nextIndex]
	port, isValid := parsePortOrService(val)

	if !isValid || port < 1 || port > 65535 {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid port number: %s", val))
	}

	opts.TargetPort = port

	return 2, nil
}

func parsePortNumber(arg string) (int, bool) {
	p, err := strconv.Atoi(arg)

	if err != nil || p < 1 || p > 65535 {
		return 0, false
	}

	return p, true
}

func printPortsHelp() {
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Usage: gitmap ports [options] [port|service]" + constants.ColorReset)
	fmt.Println()
	fmt.Println("Inspect listening TCP network ports, process owners, and firewall status.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -p, --port <port>   Inspect a specific port number (e.g. -p 22) or service name")
	fmt.Println("  -c, --common        Inspect common administrative ports (22, 80, 443, 3389, etc.)")
	fmt.Println("      --firewall, fw  Display firewall statuses for common management ports")
	fmt.Println("      --json          Output results in JSON format")
	fmt.Println("  -h, --help          Show this help message")
	fmt.Println()
	fmt.Println("Named Services:")
	fmt.Println("  gitmap ports ssh            # Inspect SSH service (port 22)")
	fmt.Println("  gitmap ports rdp            # Inspect Remote Desktop (port 3389)")
	fmt.Println("  gitmap ports winrm          # Inspect Windows Remote Management (port 5985)")
	fmt.Println("  gitmap ports firewall       # Display firewall status for all common ports")
	fmt.Println()
	fmt.Println("Supported Services:")
	fmt.Println("  ssh, sshd, rdp, http, web, https, ssl, winrm, winrm-https, winrm-ssl, dns, smb, mysql, postgres, postgresql, redis")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ports                # List active listening TCP ports")
	fmt.Println("  gitmap ports --common       # Check standard management and web ports")
	fmt.Println("  gitmap ports firewall       # Check firewall status on common ports")
	fmt.Println("  gitmap ports ssh            # Check SSH server status and firewall")
	fmt.Println("  gitmap ports rdp            # Check Remote Desktop port")
	fmt.Println("  gitmap ports winrm          # Check WinRM HTTP port")
	fmt.Println("  gitmap ports -p 22          # Check port 22 directly")
	fmt.Println("  gitmap ports ssh --json     # Export SSH port diagnosis to JSON")
	fmt.Println()
}

func outputPortsJSON(entries []PortEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal ports JSON failed:")
	}

	fmt.Println(string(data))

	return nil
}

func renderPortsTable(entries []PortEntry, opts PortsOptions) {
	fmt.Println()
	fmt.Printf("  %s%s Network Port & Firewall Diagnostics (%d ports)%s\n\n",
		constants.ColorCyan, "ℹ", len(entries), constants.ColorReset)

	cfg := buildPortsTableConfig(entries)
	termtable.PrintTable(cfg)
	printPortsFooter(entries)
}

func buildPortsColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "Port", MinWidth: 6, Align: termtable.AlignRight},
		{Title: "Proto", MinWidth: 6, Align: termtable.AlignLeft},
		{Title: "State", MinWidth: 10, Align: termtable.AlignLeft},
		{Title: "Process", MinWidth: 14, MaxWidth: 20, Align: termtable.AlignLeft},
		{Title: "PID", MinWidth: 6, Align: termtable.AlignRight},
		{Title: "Firewall", MinWidth: 12, MaxWidth: 18, Align: termtable.AlignLeft},
		{Title: "Recommendation", MinWidth: 24, MaxWidth: 48, Align: termtable.AlignLeft},
	}
}

func buildPortsRows(entries []PortEntry) []termtable.Row {
	rows := make([]termtable.Row, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, buildSinglePortRow(entry))
	}

	return rows
}

func buildSinglePortRow(entry PortEntry) termtable.Row {
	pidStr := "-"
	if entry.PID > 0 {
		pidStr = strconv.Itoa(entry.PID)
	}

	procName := entry.ProcessName
	if procName == "" {
		procName = "-"
	}

	rowColor := resolveRowColor(entry.State, entry.FirewallStatus)

	return termtable.Row{
		Color: rowColor,
		Cells: []string{
			strconv.Itoa(entry.Port),
			entry.Protocol,
			entry.State,
			procName,
			pidStr,
			entry.FirewallStatus,
			entry.Recommendation,
		},
	}
}

func resolveRowColor(state string, firewallStatus string) string {
	isListening := strings.EqualFold(state, "LISTENING")
	isAllowed := strings.Contains(strings.ToLower(firewallStatus), "allow")

	if isListening && isAllowed {
		return constants.ColorGreen
	}

	if isListening {
		return constants.ColorYellow
	}

	return constants.ColorDim
}

func buildPortsTableConfig(entries []PortEntry) termtable.TableConfig {
	return termtable.TableConfig{
		HeaderColor:  constants.ColorCyan,
		BorderColor:  constants.ColorDim,
		EllipsisText: "...",
		Columns:      buildPortsColumns(),
		Rows:         buildPortsRows(entries),
	}
}

func printPortsFooter(entries []PortEntry) {
	fmt.Println()
	hasSSH := false
	isSSHListening := false

	for _, e := range entries {
		if e.Port != 22 {
			continue
		}

		hasSSH = true
		isSSHListening = isSSHListening || strings.EqualFold(e.State, "LISTENING")
	}

	if hasSSH && !isSSHListening {
		fmt.Printf("  %s💡 Tip: Run `gitmap ssh enable` to install and start OpenSSH Server on this machine.%s\n\n",
			constants.ColorYellow, constants.ColorReset)
	}
}

func resolvePortRecommendation(port int, state string, firewallStatus string) string {
	isListening := strings.EqualFold(state, "LISTENING")
	isAllowed := strings.Contains(strings.ToLower(firewallStatus), "allow")

	if isListening {
		return resolveListeningRecommendation(port, isAllowed)
	}

	return resolveClosedRecommendation(port)
}

func resolveListeningRecommendation(port int, isAllowed bool) string {
	if !isAllowed {
		return fmt.Sprintf("Allow inbound port %d in firewall", port)
	}

	switch port {
	case 22:
		return "SSH service active and accessible"
	case 53:
		return "DNS service active"
	case 80:
		return "HTTP web service active"
	case 443:
		return "HTTPS secure web service active"
	case 445:
		return "SMB file sharing service active"
	case 3306:
		return "MySQL database service active"
	case 3389:
		return "Remote Desktop active"
	case 5432:
		return "PostgreSQL database service active"
	case 5985:
		return "WinRM HTTP active"
	case 5986:
		return "WinRM HTTPS active"
	case 6379:
		return "Redis cache service active"
	case 8080:
		return "HTTP Alt service active"
	default:
		return "Service listening and accessible"
	}
}

func resolveClosedRecommendation(port int) string {
	switch port {
	case 22:
		return "Install/start OpenSSH Server: gitmap ssh enable"
	case 53:
		return "Start DNS server service"
	case 80:
		return "Start HTTP web service (IIS/Nginx/Apache)"
	case 443:
		return "Start HTTPS web service / TLS listener"
	case 445:
		return "Start SMB / Server service"
	case 3306:
		return "Start MySQL database service"
	case 3389:
		return "Enable Remote Desktop in System Settings"
	case 5432:
		return "Start PostgreSQL database service"
	case 5985:
		return "Enable WinRM service: winrm quickconfig"
	case 5986:
		return "Enable WinRM HTTPS listener"
	case 6379:
		return "Start Redis cache service"
	case 8080:
		return "Start secondary HTTP service / proxy"
	default:
		return "Port not listening"
	}
}
