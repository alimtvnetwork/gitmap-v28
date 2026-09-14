package cmdservice

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"gopkg.in/yaml.v3"
)

// RunService handles gitmap service subcommands.
func RunService(args []string) error {
	if len(args) == 0 || isServiceHelp(args[0]) {
		printServiceHelp()

		return nil
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	return dispatchServiceSubcommand(sub, rest)
}

func isServiceHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func dispatchServiceSubcommand(sub string, rest []string) error {
	switch sub {
	case "ls", "list":
		return runServiceList(rest)
	case "status", "info":
		return runServiceStatus(rest)
	case "on", "start", "enable":
		return runServiceStart(rest)
	case "off", "stop", "disable":
		return runServiceStop(rest)
	case "create", "new", "add":
		return runServiceCreate(rest)
	case "rm", "delete", "del":
		return runServiceRemove(rest)
	case "export", "export-all":
		return runServiceExport(rest)
	case "import", "import-all":
		return runServiceImport(rest)
	default:
		return apperror.NewSimple("unknown service subcommand: "+sub, "E_SERVICE_INVALID_CMD")
	}
}

func runServiceList(args []string) error {
	driver := ResolveServiceDriver()
	services, err := driver.ListServices()
	if err != nil {
		return err
	}

	if isJson := hasServiceFlag(args, "--json"); isJson {
		return printServicesJSON(services)
	}

	renderServiceTable(services)

	return nil
}

func renderServiceTable(services []ServiceInfo) {
	fmt.Printf("\n  %s%s   %s   %s%s\n",
		constants.ColorWhite,
		cmdpull.PadVisual("SERVICE NAME", 35),
		cmdpull.PadVisual("STATUS", 15),
		"ENABLED",
		constants.ColorReset)
	fmt.Printf("  %s\n", strings.Repeat("─", 65))

	limit := 50
	if len(services) < limit {
		limit = len(services)
	}

	for i := 0; i < limit; i++ {
		renderSingleServiceRow(services[i])
	}
	fmt.Println()
}

func renderSingleServiceRow(s ServiceInfo) {
	enabledStr := "\033[31mno\033[0m"
	if s.IsEnabled {
		enabledStr = "\033[32myes\033[0m"
	}

	fmt.Printf("  %s   %s   %s\n",
		cmdpull.PadVisual(s.Name, 35),
		cmdpull.PadVisual(s.Status, 15),
		enabledStr)
}

func runServiceStatus(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("service name required", "E_SERVICE_NAME_REQUIRED")
	}

	driver := ResolveServiceDriver()
	svc, err := driver.GetService(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("\n  ▶ Service: %s\n", svc.Name)
	fmt.Printf("  • Status:  %s\n", svc.Status)
	fmt.Printf("  • Enabled: %t\n\n", svc.IsEnabled)

	return nil
}

func runServiceStart(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("service name required", "E_SERVICE_NAME_REQUIRED")
	}

	driver := ResolveServiceDriver()
	if err := driver.StartService(args[0]); err != nil {
		return err
	}

	fmt.Printf("✔ Service %q started successfully\n", args[0])

	return nil
}

func runServiceStop(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("service name required", "E_SERVICE_NAME_REQUIRED")
	}

	driver := ResolveServiceDriver()
	if err := driver.StopService(args[0]); err != nil {
		return err
	}

	fmt.Printf("✔ Service %q stopped successfully\n", args[0])

	return nil
}

func runServiceCreate(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("service name required", "E_SERVICE_NAME_REQUIRED")
	}

	name := args[0]
	execPath := extractServiceFlagValue(args, "--exec", name)
	desc := extractServiceFlagValue(args, "--desc", "GitMap managed service: "+name)
	driver := ResolveServiceDriver()
	if err := driver.CreateService(name, execPath, desc); err != nil {
		return err
	}

	fmt.Printf("✔ Service %q created successfully (exec: %s)\n", name, execPath)

	return nil
}

func runServiceRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("service name required", "E_SERVICE_NAME_REQUIRED")
	}

	driver := ResolveServiceDriver()
	if err := driver.RemoveService(args[0]); err != nil {
		return err
	}

	fmt.Printf("✔ Service %q removed successfully\n", args[0])

	return nil
}

func runServiceExport(args []string) error {
	driver := ResolveServiceDriver()
	services, err := driver.ListServices()
	if err != nil {
		return err
	}

	schema := ServiceExportSchema{
		ExportedAt: time.Now().Format(time.RFC3339),
		Platform:   runtime.GOOS,
		Services:   services,
	}

	return writeServiceExportOutput(schema, args)
}

func writeServiceExportOutput(schema ServiceExportSchema, args []string) error {
	outPath := extractServiceFlagValue(args, "--file", "")
	isYAML := hasServiceFlag(args, "--yaml")
	data, err := formatServiceExportBytes(schema, isYAML)
	if err != nil {
		return err
	}

	if outPath != "" {
		return os.WriteFile(outPath, data, 0644)
	}

	fmt.Println(string(data))

	return nil
}

func formatServiceExportBytes(schema ServiceExportSchema, isYAML bool) ([]byte, error) {
	if isYAML {
		return yaml.Marshal(schema)
	}

	return json.MarshalIndent(schema, "", "  ")
}

func runServiceImport(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("import file path required", "E_SERVICE_FILE_REQUIRED")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import file")
	}

	var schema ServiceExportSchema
	if err := unmarshalServiceExport(data, args[0], &schema); err != nil {
		return err
	}

	return applyImportedServices(schema.Services)
}

func unmarshalServiceExport(data []byte, filename string, schema *ServiceExportSchema) error {
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		return yaml.Unmarshal(data, schema)
	}

	return json.Unmarshal(data, schema)
}

func applyImportedServices(services []ServiceInfo) error {
	driver := ResolveServiceDriver()
	count := 0
	for _, s := range services {
		if s.ExecPath != "" {
			_ = driver.CreateService(s.Name, s.ExecPath, s.Description)
			count++
		}
	}

	fmt.Printf("✔ Imported and configured %d service(s)\n", count)

	return nil
}

func hasServiceFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}

	return false
}

func extractServiceFlagValue(args []string, flag, fallback string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == flag && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(args[i], flag+"=") {
			return strings.TrimPrefix(args[i], flag+"=")
		}
	}

	return fallback
}

func printServicesJSON(services []ServiceInfo) error {
	data, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal services")
	}

	fmt.Println(string(data))

	return nil
}

func printServiceHelp() {
	fmt.Println(`Usage: gitmap service [subcommand] [flags]

Commands:
  ls, list           List installed OS services
  status <name>      Inspect service status and state
  on, start <name>   Start and enable service
  off, stop <name>   Stop and disable service
  create <name>      Create and register a new service unit (--exec <cmd>)
  rm, delete <name>  Stop and remove a registered service
  export [--file]    Export service definitions to JSON or YAML
  import <file>      Import and configure service definitions`)
}
