package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SJClusterImportCmd is the cobra command for cluster topology import.
var SJClusterImportCmd = &cobra.Command{
	Use:     "import-cluster [path/to/01-config.json]",
	Aliases: []string{"import", "import-config"},
	Short:   "Import cluster nodes topology from JSON file and enroll hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterImportCLI(args)
	},
}

type clusterEnrollItem struct {
	host store.SSHHost
	hist store.SSHHistory
}

var clusterImportOut io.Writer = os.Stdout

func resolveImportConfigPath(args []string) string {
	hasArgs := len(args) > 0
	if hasArgs {
		return args[0]
	}
	return "./01-config.json"
}

func sortMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func buildClusterHost(alias, ip, user, role, encPass string, now time.Time) store.SSHHost {
	return store.SSHHost{
		ID:                fmt.Sprintf("host-%s", ip),
		Alias:             alias,
		IP:                ip,
		Username:          user,
		Port:              22,
		EncryptedPassword: encPass,
		ClusterRole:       role,
		CreatedAt:         now,
	}
}

func buildClusterHistory(ip, user string, now time.Time) store.SSHHistory {
	return store.SSHHistory{
		ID:       fmt.Sprintf("hist-%s-%d", ip, now.UnixNano()),
		HostIP:   ip,
		JoinedAt: now,
		User:     user,
	}
}

func collectControlItems(cfg *ClusterConfigJSON, encPass string, now time.Time) []clusterEnrollItem {
	keys := sortMapKeys(cfg.Control)
	items := make([]clusterEnrollItem, 0, len(keys))
	for _, alias := range keys {
		ip := cfg.Control[alias]
		host := buildClusterHost(alias, ip, cfg.User.Name, "control", encPass, now)
		hist := buildClusterHistory(ip, cfg.User.Name, now)
		items = append(items, clusterEnrollItem{host: host, hist: hist})
	}
	return items
}

func collectWorkerItems(cfg *ClusterConfigJSON, encPass string, now time.Time) []clusterEnrollItem {
	keys := sortMapKeys(cfg.Nodes)
	items := make([]clusterEnrollItem, 0, len(keys))
	for _, alias := range keys {
		ip := cfg.Nodes[alias]
		host := buildClusterHost(alias, ip, cfg.User.Name, "worker", encPass, now)
		hist := buildClusterHistory(ip, cfg.User.Name, now)
		items = append(items, clusterEnrollItem{host: host, hist: hist})
	}
	return items
}

func collectClusterItems(cfg *ClusterConfigJSON, encPass string) []clusterEnrollItem {
	now := time.Now().UTC()
	controls := collectControlItems(cfg, encPass, now)
	workers := collectWorkerItems(cfg, encPass, now)
	return append(controls, workers...)
}

func enrollClusterItems(ctx context.Context, db *sql.DB, items []clusterEnrollItem) error {
	for _, item := range items {
		if err := store.EnrollSSHHost(ctx, item.host, item.hist, db); err != nil {
			return apperror.Wrap(err, "enrollClusterItems", map[string]any{"alias": item.host.Alias, "ip": item.host.IP})
		}
	}
	return nil
}

func persistClusterTopology(ctx context.Context, items []clusterEnrollItem) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("persistClusterTopology", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()
	return enrollClusterItems(ctx, dbConn.SQL(), items)
}

func extractHostsFromItems(items []clusterEnrollItem) []store.SSHHost {
	hosts := make([]store.SSHHost, 0, len(items))
	for _, item := range items {
		hosts = append(hosts, item.host)
	}
	return hosts
}

func formatEncryptionStatus(cipherText string) string {
	isRSA := isRSACiphertext(cipherText)
	if isRSA {
		return "rsa:sha256 [encrypted]"
	}
	return "rsa:sha256 [encrypted]"
}

func flushTabWriter(w *tabwriter.Writer) error {
	if err := w.Flush(); err != nil {
		return apperror.WrapSimple(err, "flushTabWriter")
	}
	return nil
}

func renderClusterImportTable(out io.Writer, hosts []store.SSHHost) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ROLE\tALIAS\tIP\tUSERNAME\tSTATUS")
	for _, host := range hosts {
		status := formatEncryptionStatus(host.EncryptedPassword)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", host.ClusterRole, host.Alias, host.IP, host.Username, status)
	}
	return flushTabWriter(w)
}

func executeClusterImport(ctx context.Context, cfg *ClusterConfigJSON) error {
	if err := ValidateClusterUser(cfg.User); err != nil {
		return err
	}
	encPass, err := EncryptSSHPassword(cfg.User.Password)
	if err != nil {
		return err
	}
	items := collectClusterItems(cfg, encPass)
	if err := persistClusterTopology(ctx, items); err != nil {
		return err
	}
	return renderClusterImportTable(clusterImportOut, extractHostsFromItems(items))
}

// RunClusterImportCLI executes the cluster import command.
func RunClusterImportCLI(args []string) error {
	path := resolveImportConfigPath(args)
	cfg, err := LoadClusterConfigFile(path)
	if err != nil {
		return err
	}
	return executeClusterImport(context.Background(), cfg)
}
