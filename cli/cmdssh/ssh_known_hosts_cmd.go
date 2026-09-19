package cmdssh

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunSSHKnownHostsCLI handles the 'gitmap ssh known-hosts' subcommand.
func RunSSHKnownHostsCLI(args []string) error {
	if len(args) == 0 {
		return executeKnownHostsList(context.Background())
	}
	sub := strings.ToLower(args[0])
	if sub == "--help" || sub == "-h" || sub == "help" {
		printKnownHostsHelp()
		return nil
	}
	return routeKnownHostsSubcommand(sub, args[1:])
}

func routeKnownHostsSubcommand(sub string, rest []string) error {
	ctx := context.Background()
	if sub == "ls" || sub == "list" || sub == "nodes" {
		return executeKnownHostsList(ctx)
	}
	if sub == "add" || sub == "trust" {
		return executeKnownHostsAdd(ctx, rest)
	}
	if sub == "rm" || sub == "remove" || sub == "delete" || sub == "untrust" {
		return executeKnownHostsRm(ctx, rest)
	}
	if sub == "sync" {
		return executeKnownHostsSync(ctx)
	}
	return executeKnownHostsAdd(ctx, append([]string{sub}, rest...))
}

// RunSSHTrustCLI handles the 'gitmap ssh trust <target>' command.
func RunSSHTrustCLI(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printTrustHelp()
		return nil
	}
	return executeKnownHostsAdd(context.Background(), args)
}

// RunSSHUntrustCLI handles the 'gitmap ssh untrust <target>' command.
func RunSSHUntrustCLI(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("Usage: gitmap ssh untrust <target>")
		return nil
	}
	return executeKnownHostsRm(context.Background(), args)
}

func printKnownHostsHelp() {
	fmt.Printf("\n%sManage SSH Known Hosts & Host Keys%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Usage:")
	fmt.Println("  gitmap ssh known-hosts [ls|list]          List all tracked known hosts")
	fmt.Println("  gitmap ssh known-hosts add <host> [key]   Trust and add a host to known_hosts & DB")
	fmt.Println("  gitmap ssh known-hosts rm <host>          Remove a host from known_hosts & DB")
	fmt.Println("  gitmap ssh known-hosts sync               Sync ~/.ssh/known_hosts to SQLite DB")
	fmt.Println("  gitmap ssh trust <host>                   Auto-scan & trust a remote host key")
	fmt.Println("  gitmap ssh untrust <host>                 Untrust & delete a remote host key")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap ssh known-hosts ls")
	fmt.Println("  gitmap ssh trust 192.168.1.5")
	fmt.Println("  gitmap ssh trust u2")
	fmt.Println("  gitmap ssh untrust 192.168.1.5")
	fmt.Println("  gitmap ssh known-hosts sync")
	fmt.Println()
}

func printTrustHelp() {
	fmt.Println("Usage: gitmap ssh trust <target> [public-key]")
	fmt.Println("Auto-scans and trusts the remote host key, saving it to known_hosts and SQLite DB.")
}

func executeKnownHostsList(ctx context.Context) error {
	dbConn, err := openSSHDB()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	_, _ = SyncKnownHostsFileWithDB(ctx, dbConn.SQL())
	hosts, err := store.ListSSHKnownHosts(ctx, dbConn.SQL())
	if err != nil {
		return err
	}
	renderKnownHostsTable(os.Stdout, hosts)
	return nil
}

func executeKnownHostsAdd(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("target host is required")
	}
	dbConn, err := openSSHDB()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	target := args[0]
	key := ""
	if len(args) > 1 {
		key = strings.Join(args[1:], " ")
	}
	kh, err := TrustRemoteTarget(ctx, target, key, dbConn.SQL())
	if err != nil {
		return err
	}
	fmt.Printf("%s✓ Host '%s' trusted successfully!%s\n", constants.ColorGreen, target, constants.ColorReset)
	fmt.Printf("  Key Type:    %s\n", kh.KeyType)
	fmt.Printf("  Fingerprint: %s\n\n", kh.Fingerprint)
	return nil
}

func executeKnownHostsRm(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("target host is required")
	}
	dbConn, err := openSSHDB()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	target := args[0]
	_, err = UntrustRemoteTarget(ctx, target, dbConn.SQL())
	if err != nil {
		return err
	}
	fmt.Printf("%s✓ Removed '%s' from known_hosts and database.%s\n\n", constants.ColorGreen, target, constants.ColorReset)
	return nil
}

func executeKnownHostsSync(ctx context.Context) error {
	dbConn, err := openSSHDB()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	count, err := SyncKnownHostsFileWithDB(ctx, dbConn.SQL())
	if err != nil {
		return err
	}
	fmt.Printf("%s✓ Synchronized %d known host(s) to database.%s\n\n", constants.ColorGreen, count, constants.ColorReset)
	return nil
}

func renderKnownHostsTable(w io.Writer, hosts []store.SSHKnownHost) {
	if len(hosts) == 0 {
		fmt.Fprintf(w, "\n%sNo known SSH hosts tracked yet.%s\n\n", constants.ColorMuted, constants.ColorReset)
		fmt.Fprintf(w, "Run 'gitmap ssh trust <host>' or 'gitmap ssh known-hosts sync' to import.\n\n")
		return
	}
	fmt.Fprintf(w, "\n%sKnown SSH Hosts (%d tracked):%s\n\n", constants.ColorCyan, len(hosts), constants.ColorReset)
	fmt.Fprintf(w, "%-22s %-14s %-50s %-16s\n", "HOST / IP", "KEY TYPE", "FINGERPRINT", "UPDATED")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 106))
	for _, h := range hosts {
		updated := h.UpdatedAt.Format("2006-01-02 15:04")
		fmt.Fprintf(w, "%-22s %-14s %-50s %-16s\n", h.Host, h.KeyType, h.Fingerprint, updated)
	}
	fmt.Fprintln(w)
}
