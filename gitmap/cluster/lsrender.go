package cluster

import (
	"fmt"
	"strconv"
	"strings"
)

func RenderNodeTable(nodes []ClusterNode, showRole bool) string {
	if len(nodes) == 0 {
		return ""
	}
	idLen := 2
	ipLen := 10
	nameLen := 12
	roleLen := 4
	osLen := 2
	statusLen := 6
	hbLen := 14

	for _, n := range nodes {
		if len(strconv.Itoa(n.DisplayId)) > idLen {
			idLen = len(strconv.Itoa(n.DisplayId))
		}
		if len(n.IP) > ipLen {
			ipLen = len(n.IP)
		}
		if len(n.Alias) > nameLen {
			nameLen = len(n.Alias)
		}
		roleStr := "client"
		if n.IsServer {
			roleStr = "server"
		}
		if len(roleStr) > roleLen {
			roleLen = len(roleStr)
		}
		if len(n.OS) > osLen {
			osLen = len(n.OS)
		}
	}

	var sb strings.Builder

	idH := padRight("ID", idLen)
	ipH := padRight("IP Address", ipLen)
	nameH := padRight("Machine Name", nameLen)
	roleH := padRight("Role", roleLen)
	osH := padRight("OS", osLen)
	statusH := padRight("Status", statusLen)
	hbH := padRight("Last Heartbeat", hbLen)

	if showRole {
		sb.WriteString(fmt.Sprintf("%s  %s  %s  %s  %s  %s  %s\n", idH, ipH, nameH, roleH, osH, statusH, hbH))
	} else {
		sb.WriteString(fmt.Sprintf("%s  %s  %s  %s  %s  %s\n", idH, ipH, nameH, osH, statusH, hbH))
	}

	for _, n := range nodes {
		roleStr := "client"
		if n.IsServer {
			roleStr = "server"
		}

		cId := padRight(strconv.Itoa(n.DisplayId), idLen)
		cIp := padRight(n.IP, ipLen)
		cName := padRight(n.Alias, nameLen)
		cRole := padRight(roleStr, roleLen)
		cOs := padRight(n.OS, osLen)
		cStatus := padRight("Online", statusLen)
		cHb := padRight("-", hbLen)

		if showRole {
			sb.WriteString(fmt.Sprintf("%s  %s  %s  %s  %s  %s  %s\n", cId, cIp, cName, cRole, cOs, cStatus, cHb))
		} else {
			sb.WriteString(fmt.Sprintf("%s  %s  %s  %s  %s  %s\n", cId, cIp, cName, cOs, cStatus, cHb))
		}
	}
	return sb.String()
}

func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}

