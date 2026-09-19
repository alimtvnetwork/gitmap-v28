package cmdssh

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	defaultHostRole      = "worker"
	defaultHostUser      = "root"
	defaultHostAlias     = "-"
	defaultHostStatus    = "ready"
	defaultSSHPort       = 22
	dividerLength        = 100
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

func padVisual(s string, width int) string {
	rCount := utf8.RuneCountInString(s)
	if rCount >= width {
		return s
	}
	return s + strings.Repeat(" ", width-rCount)
}

func isControlPlaneRole(role string) bool {
	return role == "control" || role == "control-plane" || role == "master" || role == "cp"
}

func formatRoleColored(role string, width int) string {
	plain := padVisual(role, width)
	if isControlPlaneRole(role) {
		return constants.ColorYellow + plain + constants.ColorReset
	}
	if role == defaultHostRole {
		return constants.ColorCyan + plain + constants.ColorReset
	}
	return constants.ColorDim + plain + constants.ColorReset
}

func formatStatusColored(status string, width int) string {
	if status == "ready" || status == "active" {
		text := "● " + status
		plain := padVisual(text, width)
		return constants.ColorGreen + plain + constants.ColorReset
	}
	return constants.ColorDim + padVisual(status, width) + constants.ColorReset
}

func formatAliasColored(alias string, width int) string {
	plain := padVisual(alias, width)
	return constants.ColorBold + constants.ColorWhite + plain + constants.ColorReset
}

func formatHostPortColored(hostPort string, width int) string {
	plain := padVisual(hostPort, width)
	return constants.ColorWhite + plain + constants.ColorReset
}

func formatUserColored(user string, width int) string {
	plain := padVisual(user, width)
	return constants.ColorDim + plain + constants.ColorReset
}

func formatEnrolledColored(enrolled string, width int) string {
	plain := padVisual(enrolled, width)
	return constants.ColorDim + plain + constants.ColorReset
}

func renderHostsTableHeader(out io.Writer) error {
	colAlias := padVisual("ALIAS", 16)
	colRole := padVisual("ROLE", 14)
	colHost := padVisual("HOST (IP:PORT)", 22)
	colUser := padVisual("USER", 14)
	colStatus := padVisual("STATUS", 10)
	colEnrolled := padVisual("ENROLLED", 19)

	headerLine := fmt.Sprintf("  %s %s %s %s %s %s\n",
		colAlias, colRole, colHost, colUser, colStatus, colEnrolled)
	if _, err := fmt.Fprintf(out, "\n%s%s%s", constants.ColorCyan, headerLine, constants.ColorReset); err != nil {
		return apperror.WrapSimple(err, "renderHostsTableHeader")
	}
	divider := strings.Repeat("-", dividerLength)
	if _, err := fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset); err != nil {
		return apperror.WrapSimple(err, "renderHostsTableDivider")
	}
	return nil
}

func renderHostTableRow(out io.Writer, h store.SSHHost) error {
	alias := formatAliasColored(resolveTableAlias(h.Alias), 16)
	role := formatRoleColored(resolveTableRole(h.ClusterRole), 14)
	hostPort := formatHostPortColored(formatTableHostPort(h.IP, h.Port), 22)
	user := formatUserColored(resolveTableUser(h.Username), 14)
	status := formatStatusColored(defaultHostStatus, 10)
	enrolled := formatEnrolledColored(formatEnrolledTime(h.CreatedAt), 19)
	_, err := fmt.Fprintf(out, "  %s %s %s %s %s %s\n",
		alias, role, hostPort, user, status, enrolled)
	if err != nil {
		return apperror.WrapSimple(err, "renderHostTableRow")
	}
	return nil
}

func renderHostsTableFooter(out io.Writer, count int) error {
	_, err := fmt.Fprintf(out, "\n  %sTotal: %d registered node(s)%s\n\n",
		constants.ColorDim, count, constants.ColorReset)
	if err != nil {
		return apperror.WrapSimple(err, "renderHostsTableFooter")
	}
	return nil
}

func renderEmptyHostsNotice(out io.Writer) error {
	_, err := fmt.Fprintf(out, "\n  %s● No nodes registered. Enroll with: gitmap sj add <user@ip|ip> [alias]%s\n\n",
		constants.ColorYellow, constants.ColorReset)
	if err != nil {
		return apperror.WrapSimple(err, "renderEmptyHostsNotice")
	}
	return nil
}

// RenderSSHHostsTable renders a fixed-width colored table of SSH hosts.
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
