package cmdssh

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

type clusterInitFlags struct {
	controlIP  string
	workersRaw string
	userName   string
	userPass   string
	outFile    string
	isForce    bool
}

func registerClusterInitFlags(fs *flag.FlagSet, f *clusterInitFlags) {
	fs.StringVar(&f.controlIP, "control", "192.168.1.10", "Control-plane master node IP")
	fs.StringVar(&f.workersRaw, "workers", "192.168.1.11,192.168.1.12", "Comma-separated worker node IPs")
	fs.StringVar(&f.userName, "user", "ubuntu", "SSH user for node management")
	fs.StringVar(&f.userPass, "password", "ChangeMePassword123", "Initial node password")
	fs.StringVar(&f.outFile, "out", "01-config.json", "Output JSON configuration path")
	fs.BoolVar(&f.isForce, "force", false, "Overwrite existing config file")
	fs.BoolVar(&f.isForce, "f", false, "Short alias for --force")
}

func parseClusterInitFlags(args []string) (clusterInitFlags, error) {
	fs := flag.NewFlagSet("cluster init", flag.ContinueOnError)
	var f clusterInitFlags
	registerClusterInitFlags(fs, &f)

	if err := fs.Parse(args); err != nil {
		return f, apperror.WrapSimple(err, "parse cluster init flags")
	}

	if fs.NArg() > 0 && f.outFile == "01-config.json" {
		f.outFile = fs.Arg(0)
	}

	return f, nil
}

func buildWorkerMap(workersRaw string) map[string]string {
	nodes := make(map[string]string)
	parts := strings.Split(workersRaw, ",")
	for idx, p := range parts {
		ip := strings.TrimSpace(p)
		if ip != "" {
			alias := fmt.Sprintf("worker-%d", idx+1)
			nodes[alias] = ip
		}
	}
	if len(nodes) == 0 {
		nodes["worker-1"] = "192.168.1.11"
	}
	return nodes
}

func generateClusterConfigJSON(f clusterInitFlags) ClusterConfigJSON {
	control := map[string]string{
		"k8s-m1": f.controlIP,
	}
	nodes := buildWorkerMap(f.workersRaw)
	user := ClusterUserJSON{
		Name:     f.userName,
		Password: f.userPass,
	}
	return ClusterConfigJSON{
		Control: control,
		Nodes:   nodes,
		User:    user,
	}
}

func writeClusterConfigFile(path string, cfg ClusterConfigJSON, isForce bool) *apperror.AppError {
	if _, err := os.Stat(path); err == nil && !isForce {
		return apperror.NewValidationError(fmt.Sprintf("target config file %s already exists; use --force to overwrite", path))
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal cluster config")
	}

	if writeErr := os.WriteFile(path, append(data, '\n'), 0o644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write cluster config file")
	}

	return nil
}

func printClusterInitSuccess(path string, cfg ClusterConfigJSON) {
	fmt.Printf("✓ Generated cluster topology configuration: %s\n", path)
	fmt.Println("  Control Plane:")
	for alias, ip := range cfg.Control {
		fmt.Printf("    • %s: %s\n", alias, ip)
	}
	fmt.Println("  Worker Nodes:")
	for alias, ip := range cfg.Nodes {
		fmt.Printf("    • %s: %s\n", alias, ip)
	}
	fmt.Printf("  SSH User: %s\n\n", cfg.User.Name)
	fmt.Printf("  Next: Edit credentials, then import into cluster:\n")
	fmt.Printf("    gitmap cluster import %s\n", filepath.ToSlash(path))
}

func runClusterInit(args []string) error {
	if hasHelpFlag(args) {
		helptext.Print("cluster")
		return nil
	}

	flags, err := parseClusterInitFlags(args)
	if err != nil {
		return err
	}

	cfg := generateClusterConfigJSON(flags)
	if writeErr := writeClusterConfigFile(flags.outFile, cfg, flags.isForce); writeErr != nil {
		return writeErr
	}

	printClusterInitSuccess(flags.outFile, cfg)
	return nil
}

// RunClusterInitCLI executes cluster init command to scaffold JSON topology.
func RunClusterInitCLI(args []string) error {
	return runClusterInit(args)
}
