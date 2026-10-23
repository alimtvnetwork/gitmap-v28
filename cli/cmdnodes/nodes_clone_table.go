// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PreFlightNodeInfo captures probing telemetry for an individual fleet node.
type PreFlightNodeInfo struct {
	Alias           string        `json:"alias"`
	Host            string        `json:"host"`
	Role            string        `json:"role"`
	OS              string        `json:"os"`
	Arch            string        `json:"arch"`
	Version         string        `json:"version"`
	Commit          string        `json:"commit"`
	IsOnline        bool          `json:"isOnline"`
	StatusBadge     string        `json:"statusBadge"`
	RemoteTargetDir string        `json:"remoteTargetDir"`
	ProbeDuration   time.Duration `json:"probeDuration"`
	Error           string        `json:"error,omitempty"`
}

// NodePreFlightInfo aliases PreFlightNodeInfo for cross-spec compatibility.
type NodePreFlightInfo = PreFlightNodeInfo

// FleetPreFlightReport aggregates pre-flight probes across the fleet.
type FleetPreFlightReport struct {
	TotalCount   int                 `json:"totalCount"`
	OnlineCount  int                 `json:"onlineCount"`
	OfflineCount int                 `json:"offlineCount"`
	Nodes        []PreFlightNodeInfo `json:"nodes"`
}

var (
	nodeVersionMu    sync.RWMutex
	nodeVersionCache = make(map[string]string)
)

func recordNodeVersion(alias, version string) {
	nodeVersionMu.Lock()
	defer nodeVersionMu.Unlock()

	nodeVersionCache[alias] = version
}

func getNodeVersion(alias string) string {
	nodeVersionMu.RLock()
	defer nodeVersionMu.RUnlock()

	return nodeVersionCache[alias]
}

func resolveDisplayValue(val string) string {
	if val == "" {
		return "-"
	}

	return val
}

func isANSIEscapeTerminator(r rune) bool {
	return r == 'm' || r == 'K' || r == 'J' || r == 'H'
}

func processANSIRune(r rune, isInEsc bool, b *strings.Builder) bool {
	if r == 0x1b {
		return true
	}

	if isInEsc {
		return !isANSIEscapeTerminator(r)
	}

	b.WriteRune(r)
	return false
}

func stripANSI(s string) string {
	var b strings.Builder
	isInEsc := false
	for _, r := range s {
		isInEsc = processANSIRune(r, isInEsc, &b)
	}

	return b.String()
}

func visualLen(s string) int {
	clean := stripANSI(s)
	return len([]rune(clean))
}

func padVisual(s string, width int) string {
	vl := visualLen(s)
	if vl >= width {
		return s
	}

	return s + strings.Repeat(" ", width-vl)
}

func resolvePreFlightBadge(isOnline bool, errStr string) string {
	if isOnline {
		return constants.ColorGreen + "● online" + constants.ColorReset
	}

	low := strings.ToLower(errStr)
	if strings.Contains(low, "auth") {
		return constants.ColorMagenta + "▲ auth_failed" + constants.ColorReset
	}

	if strings.Contains(low, "unreachable") {
		return constants.ColorRed + "✗ unreachable" + constants.ColorReset
	}

	return constants.ColorYellow + "○ offline" + constants.ColorReset
}

func renderPreFlightCounters(out io.Writer, total, online, offline int) {
	fmt.Fprintf(out, "  ▸ Fleet Readiness: %d registered node(s) | %d online | %d offline\n\n",
		total, online, offline)
}

func renderPreFlightHeader(out io.Writer) {
	fmt.Fprintf(out, "    %-16s %-18s %-10s %-14s %-32s %-14s\n",
		"NODE (ALIAS)", "HOST", "OS", "VERSION", "DESTINATION", "STATUS")
	fmt.Fprintln(out, "    --------------------------------------------------------------------------------------------------------")
}

func renderPreFlightRow(out io.Writer, node PreFlightNodeInfo) {
	ver := resolveDisplayValue(node.Version)
	dest := resolveDisplayValue(node.RemoteTargetDir)
	badge := resolvePreFlightBadge(node.IsOnline, node.Error)
	fmt.Fprintf(out, "    %-16s %-18s %-10s %-14s %-32s %s\n",
		node.Alias, node.Host, node.OS, ver, dest, padVisual(badge, 14))
}

func renderFleetPreFlightTable(out io.Writer, report FleetPreFlightReport) {
	fmt.Fprintln(out)
	renderPreFlightCounters(out, report.TotalCount, report.OnlineCount, report.OfflineCount)
	renderPreFlightHeader(out)
	for _, n := range report.Nodes {
		renderPreFlightRow(out, n)
	}

	fmt.Fprintln(out, "    --------------------------------------------------------------------------------------------------------")
	fmt.Fprintln(out)
}

func renderBannerHeaderBox(out io.Writer, kind NodesCloneKindType) {
	title := fmt.Sprintf("GITMAP FLEET NODES %s DISPATCH", strings.ToUpper(string(kind)))
	fmt.Fprintln(out, "    ┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Fprintf(out, "    │ %-100s │\n", title)
	fmt.Fprintln(out, "    └──────────────────────────────────────────────────────────────────────────────────────────────────────┘")
}

func formatBannerMode(kind NodesCloneKindType) string {
	switch kind {
	case CloneKindCFR:
		return "cfr (clone-fix-repo)"
	case CloneKindCFRP:
		return "cfrp (clone-fix-repo-pub)"
	case CloneKindClone:
		return "clone (multi-node clone)"
	default:
		return "clone (multi-node clone)"
	}
}

func formatBannerTarget(opts NodesCloneOptions) string {
	if opts.HasFile {
		return opts.DetectedFile + " (staged manifest)"
	}

	if len(opts.PassArgs) > 0 {
		return strings.Join(opts.PassArgs, ", ")
	}

	return "(auto-detected)"
}

func formatBannerWorkdir(opts NodesCloneOptions) string {
	if opts.HasCustomTargetDir && opts.TargetDir != "" {
		return fmt.Sprintf("%s (custom destination)", opts.TargetDir)
	}

	if opts.RelativeSubdir != "" {
		return fmt.Sprintf("D:\\work (preserved relative: %s)", opts.RelativeSubdir)
	}

	return "D:\\work"
}

func formatBannerScope(isSkipLocal bool, onlineCount int) string {
	if isSkipLocal {
		return fmt.Sprintf("%d active remote worker(s) (remote-only)", onlineCount)
	}

	return fmt.Sprintf("Local host + %d active remote worker(s)", onlineCount)
}

func renderBannerMetadata(out io.Writer, opts NodesCloneOptions, onlineCount int) {
	fmt.Fprintf(out, "    • Mode:        %s\n", formatBannerMode(opts.Kind))
	fmt.Fprintf(out, "    • Target:      %s\n", formatBannerTarget(opts))
	fmt.Fprintf(out, "    • Workdir:     %s\n", formatBannerWorkdir(opts))
	fmt.Fprintf(out, "    • Scope:       %s\n", formatBannerScope(opts.IsSkipLocal, onlineCount))
}

func formatDispatchPrefix(hasRoute bool) string {
	if hasRoute {
		return "                   "
	}

	return "    • Dispatch:    "
}

func renderSingleDispatchRoute(out io.Writer, n PreFlightNodeInfo, hasRoute bool) {
	dest := n.RemoteTargetDir
	if dest == "" {
		dest = "D:\\work"
	}

	prefix := formatDispatchPrefix(hasRoute)
	fmt.Fprintf(out, "%s%s -> %s:%s\n", prefix, n.Alias, n.Host, dest)
}

func renderBannerDispatchRoutes(out io.Writer, nodes []PreFlightNodeInfo) {
	hasRoute := false
	for _, n := range nodes {
		if !n.IsOnline {
			continue
		}

		renderSingleDispatchRoute(out, n, hasRoute)
		hasRoute = true
	}
}

func renderFleetStartBanner(out io.Writer, opts NodesCloneOptions, report FleetPreFlightReport) {
	fmt.Fprintln(out)
	renderBannerHeaderBox(out, opts.Kind)
	renderBannerMetadata(out, opts, report.OnlineCount)
	renderBannerDispatchRoutes(out, report.Nodes)
	fmt.Fprintln(out)
}
