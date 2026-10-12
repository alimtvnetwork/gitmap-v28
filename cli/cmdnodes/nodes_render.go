package cmdnodes

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func padCell(s string, width int) string {
	count := utf8.RuneCountInString(s)
	if count >= width {
		return s
	}
	return s + strings.Repeat(" ", width-count)
}

func renderUnifiedNodesTable(out io.Writer, nodes []UnifiedFleetNode) error {
	if len(nodes) == 0 {
		return renderEmptyUnifiedNodesNotice(out)
	}
	if err := renderCommandsMatrixTable(out, nodes); err != nil {
		return err
	}
	if err := renderPrimaryNodesTable(out, nodes); err != nil {
		return err
	}
	renderUnifiedNodesFooter(out)
	return nil
}

func renderEmptyUnifiedNodesNotice(out io.Writer) error {
	fmt.Fprintln(out)
	fmt.Fprintf(out, "  %s● No nodes registered across SSH, Cluster, or Server-Client networks.%s\n",
		constants.ColorYellow, constants.ColorReset)
	fmt.Fprintln(out, "    Enroll nodes with: gitmap ssh join <user@ip|ip> [alias]")
	fmt.Fprintln(out)
	renderUnifiedNodesFooter(out)
	return nil
}

func renderPrimaryNodesTable(out io.Writer, nodes []UnifiedFleetNode) error {
	if err := renderHostsHeader(out); err != nil {
		return err
	}
	for _, n := range nodes {
		renderPrimaryNodeRow(out, n)
	}
	fmt.Fprintf(out, "\n  %sTotal: %d registered node(s)%s\n\n", constants.ColorDim, len(nodes), constants.ColorReset)
	return nil
}

func renderHostsHeader(out io.Writer) error {
	colAlias := padCell("ALIAS", 16)
	colRole := padCell("ROLE", 14)
	colHost := padCell("HOST (IP:PORT)", 22)
	colUser := padCell("USER", 14)
	colStatus := padCell("STATUS", 21)
	colEnrolled := padCell("ENROLLED", 19)

	headerLine := fmt.Sprintf("  %s %s %s %s %s %s\n",
		colAlias, colRole, colHost, colUser, colStatus, colEnrolled)
	fmt.Fprintf(out, "\n%s%s%s", constants.ColorCyan, headerLine, constants.ColorReset)
	divider := strings.Repeat("-", 110)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)
	return nil
}

func renderPrimaryNodeRow(out io.Writer, n UnifiedFleetNode) {
	alias := formatCellAlias(n.Alias, 16)
	role := formatCellRole(n.Role, 14)
	hostPort := formatCellHost(fmt.Sprintf("%s:%d", n.Host, n.Port), 22)
	user := formatCellDim(n.User, 14)
	status := formatCellStatus(n.Status, 21)
	enrolled := formatCellDim(n.EnrolledAt, 19)

	fmt.Fprintf(out, "  %s %s %s %s %s %s\n", alias, role, hostPort, user, status, enrolled)
}

func renderCommandsMatrixTable(out io.Writer, nodes []UnifiedFleetNode) error {
	renderCommandsMatrixHeader(out)
	for _, n := range nodes {
		renderCommandsMatrixRow(out, n)
	}
	divider := strings.Repeat("-", 110)
	fmt.Fprintf(out, "  %s%s%s\n\n", constants.ColorDim, divider, constants.ColorReset)
	return nil
}

func renderCommandsMatrixHeader(out io.Writer) {
	colAlias := padCell("NODE (ALIAS)", 16)
	colSys := padCell("SUBSYSTEMS", 22)
	colCmds := padCell("SUPPORTED COMMANDS & CLUSTERS", 70)

	headerLine := fmt.Sprintf("  %s %s %s\n", colAlias, colSys, colCmds)
	fmt.Fprintf(out, "\n%s%s%s", constants.ColorCyan, headerLine, constants.ColorReset)
	divider := strings.Repeat("-", 110)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)
}

func renderCommandsMatrixRow(out io.Writer, n UnifiedFleetNode) {
	alias := formatCellAlias(n.Alias, 16)
	subsys := formatCellSubsystems(strings.Join(n.Subsystems, ", "), 22)
	cmds := formatSupportedCommands(n)

	fmt.Fprintf(out, "  %s %s %s\n", alias, subsys, cmds)
}

func formatSupportedCommands(n UnifiedFleetNode) string {
	var parts []string
	if hasSubsystem(n, "SSH") {
		parts = append(parts, fmt.Sprintf("gitmap ssh %s", n.Alias))
		parts = append(parts, fmt.Sprintf("gitmap exec %s", n.Alias))
	}
	if hasSubsystem(n, "Cluster") {
		parts = append(parts, fmt.Sprintf("gitmap cluster exec %s", n.Alias))
	}
	if hasSubsystem(n, "SC") {
		parts = append(parts, "gitmap sc exec")
	}
	return constants.ColorDim + strings.Join(parts, ", ") + constants.ColorReset
}

func formatCellAlias(alias string, width int) string {
	return constants.ColorBold + constants.ColorWhite + padCell(alias, width) + constants.ColorReset
}

func formatCellRole(role string, width int) string {
	plain := padCell(role, width)
	if role == "control" || role == "master" || role == "server" {
		return constants.ColorYellow + plain + constants.ColorReset
	}
	return constants.ColorCyan + plain + constants.ColorReset
}

func formatCellHost(host string, width int) string {
	return constants.ColorWhite + padCell(host, width) + constants.ColorReset
}

func formatCellDim(val string, width int) string {
	return constants.ColorDim + padCell(val, width) + constants.ColorReset
}

func formatCellSubsystems(sys string, width int) string {
	return constants.ColorCyan + padCell(sys, width) + constants.ColorReset
}

func formatCellStatus(status string, width int) string {
	if status == "ready" || strings.HasPrefix(status, "●") {
		text := ensurePrefix(status, "●")
		return constants.ColorGreen + padCell(text, width) + constants.ColorReset
	}
	if strings.Contains(status, "offline") || strings.HasPrefix(status, "○") {
		text := ensurePrefix(status, "○")
		return constants.ColorDim + padCell(text, width) + constants.ColorReset
	}
	if strings.Contains(status, "auth failed") || strings.HasPrefix(status, "▲") {
		text := ensurePrefix(status, "▲")
		return constants.ColorYellow + padCell(text, width) + constants.ColorReset
	}
	return constants.ColorDim + padCell(status, width) + constants.ColorReset
}

func ensurePrefix(s, pfx string) string {
	if strings.HasPrefix(s, pfx) {
		return s
	}
	return pfx + " " + s
}

func renderUnifiedNodesFooter(out io.Writer) {
	fmt.Fprintf(out, "  %s💡 Fleet Operations & Supported Command Suggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Fprintln(out, "    • Shell Access:            gitmap ssh <alias>                (or: gitmap ssh login <alias>)")
	fmt.Fprintln(out, "    • Command Execution:       gitmap exec <alias> \"<cmd>\"       (or: gitmap ssh exec <alias> \"<cmd>\")")
	fmt.Fprintln(out, "    • Cluster Remote Exec:     gitmap cluster exec <alias> \"<cmd>\"")
	fmt.Fprintln(out, "    • Broadcast to All Nodes:  gitmap sc exec \"<cmd>\"            (or: gitmap ssh exec all \"<cmd>\")")
	fmt.Fprintln(out, "    • Sync Keys & Ping Nodes:  gitmap ssh deploy-keys            | gitmap nodes ping")
	fmt.Fprintln(out, "    • Fleet Clones & CFR:      gitmap nodes clone <repo|file>    | gitmap nodes cfr <repo|file>")
	fmt.Fprintln(out)
}
