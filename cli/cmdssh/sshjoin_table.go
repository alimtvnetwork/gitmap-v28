package cmdssh

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	defaultHostRole      = "worker"
	defaultHostUser      = "root"
	defaultHostAlias     = "-"
	defaultHostStatus    = "ready"
	defaultSSHPort       = 22
	dividerLength        = 95
	tableRowFormat       = "  %-16s %-10s %-22s %-14s %-10s %-19s\n"
	msgNoNodesRegistered = "No nodes registered. Enroll with: gitmap sj add <user@ip|ip> [alias]\n"
)

func resolveTableAlias(alias string) string {
	if alias == "" {
		return defaultHostAlias
	}
	return alias
}

func resolveTableRole(role string) string {
	if role == "" {
		return defaultHostRole
	}
	return role
}

func resolveTableUser(user string) string {
	if user == "" {
		return defaultHostUser
	}
	return user
}

func resolveTablePort(port int) int {
	if port <= 0 {
		return defaultSSHPort
	}
	return port
}

func formatTableHostPort(ip string, port int) string {
	p := resolveTablePort(port)
	return fmt.Sprintf("%s:%d", ip, p)
}

func formatEnrolledTime(t time.Time) string {
	if t.IsZero() {
		return defaultHostAlias
	}
	return t.Format("2006-01-02 15:04:05")
}

func renderHostsTableHeader(out io.Writer) error {
	_, err := fmt.Fprintf(out, tableRowFormat,
		"ALIAS", "ROLE", "HOST (IP:PORT)", "USER", "STATUS", "ENROLLED")
	if err != nil {
		return apperror.WrapSimple(err, "renderHostsTableHeader")
	}
	divider := strings.Repeat("-", dividerLength)
	if _, err := fmt.Fprintf(out, "  %s\n", divider); err != nil {
		return apperror.WrapSimple(err, "renderHostsTableDivider")
	}
	return nil
}

func renderHostTableRow(out io.Writer, h store.SSHHost) error {
	alias := resolveTableAlias(h.Alias)
	role := resolveTableRole(h.ClusterRole)
	user := resolveTableUser(h.Username)
	hostPort := formatTableHostPort(h.IP, h.Port)
	enrolled := formatEnrolledTime(h.CreatedAt)
	_, err := fmt.Fprintf(out, tableRowFormat,
		alias, role, hostPort, user, defaultHostStatus, enrolled)
	if err != nil {
		return apperror.WrapSimple(err, "renderHostTableRow")
	}
	return nil
}

func renderHostsTableFooter(out io.Writer, count int) error {
	_, err := fmt.Fprintf(out, "  Total: %d registered node(s)\n", count)
	if err != nil {
		return apperror.WrapSimple(err, "renderHostsTableFooter")
	}
	return nil
}

func renderEmptyHostsNotice(out io.Writer) error {
	_, err := fmt.Fprint(out, msgNoNodesRegistered)
	if err != nil {
		return apperror.WrapSimple(err, "renderEmptyHostsNotice")
	}
	return nil
}

// RenderSSHHostsTable renders a fixed-width ASCII table of SSH hosts.
func RenderSSHHostsTable(out io.Writer, hosts []store.SSHHost) error {
	if len(hosts) == 0 {
		return renderEmptyHostsNotice(out)
	}
	return renderHostsTableBody(out, hosts)
}

func renderHostsTableBody(out io.Writer, hosts []store.SSHHost) error {
	if err := renderHostsTableHeader(out); err != nil {
		return err
	}
	for _, h := range hosts {
		if err := renderHostTableRow(out, h); err != nil {
			return err
		}
	}
	return renderHostsTableFooter(out, len(hosts))
}
