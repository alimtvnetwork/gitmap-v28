package cmdnodes

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/charmbracelet/lipgloss"
)

var (
	resStyleTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	resStyleHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	resStyleDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	resStyleIdle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	resStyleNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c"))
	resStyleHigh   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
	resStyleExceed = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
)

func printResourcesJSON(records []NodeResourceRecord) *apperror.AppError {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "printResourcesJSON")
	}
	fmt.Println(string(data))
	return nil
}

func renderResourcesTable(out io.Writer, records []NodeResourceRecord, opts nodesResourcesOptions) *apperror.AppError {
	border, div := strings.Repeat("═", 116), resStyleDim.Render(strings.Repeat("─", 118))
	hdr := resStyleHeader.Render(fmt.Sprintf("%s %s %s %s %s %s %s", padCell("NODE (ALIAS)", 16), padCell("HOST (IP:PORT)", 20), padCell("CPU%", 8), padCell("RAM% (USED/TOTAL)", 26), padCell("STATUS", 14), padCell("BENCHMARK", 20), padCell("LATENCY", 10)))
	fmt.Fprintf(out, "\n  %s╔%s╗%s\n", constants.ColorDim, border, constants.ColorReset)
	fmt.Fprintf(out, "  %s║%s%s%s║%s\n", constants.ColorDim, constants.ColorReset, resStyleTitle.Render(padCell(" GITMAP FLEET NODES RESOURCE METRICS & BENCHMARK", 116)), constants.ColorDim, constants.ColorReset)
	fmt.Fprintf(out, "  %s╚%s╝%s\n", constants.ColorDim, border, constants.ColorReset)
	fmt.Fprintf(out, "  Probing %d node(s) (Benchmark: %d%%)...\n\n  %s\n  %s\n", len(records), opts.benchmarkLimit, hdr, div)

	var onlineCount, idleCount int
	var maxLat time.Duration
	for _, r := range records {
		if r.Latency > maxLat {
			maxLat = r.Latency
		}
		if r.IsOnline {
			onlineCount++
			if r.IsUnderLoaded {
				idleCount++
			}
		}
		renderResourceRow(out, r)
	}
	fmt.Fprintf(out, "\n  %s\n  Fleet Resources: %d/%d online (%d idle/boost-ready) | Limit: %d%% | Max Latency: %s\n\n",
		div, onlineCount, len(records), idleCount, opts.benchmarkLimit, formatDurationShort(maxLat))
	return nil
}

func renderResourceRow(out io.Writer, r NodeResourceRecord) {
	cpu, ram, status, bench := resStyleDim.Render(padCell("-", 8)), resStyleDim.Render(padCell("-", 26)), resStyleDim.Render(padCell("○ Offline", 14)), resStyleDim.Render(padCell("N/A", 20))
	if r.IsOnline {
		cpu = resPercentStyle(r.CPUPercent, r.BenchmarkLimit).Render(padCell(fmt.Sprintf("%d%%", r.CPUPercent), 8))
		ram = resPercentStyle(r.RAMPercent, r.BenchmarkLimit).Render(padCell(fmt.Sprintf("%d%% (%d/%d MB)", r.RAMPercent, r.UsedRAMMB, r.TotalRAMMB), 26))
		status, bench = resStyleNormal.Render(padCell("● Normal", 14)), resStyleExceed.Render(padCell("Exceeded", 20))
		if r.IsUnderLoaded {
			status, bench = resStyleIdle.Render(padCell("● Idle", 14)), resStyleIdle.Render(padCell(fmt.Sprintf("Passed (<%d%%)", r.BenchmarkLimit), 20))
		} else if r.LoadStatus == "High" {
			status = resStyleHigh.Render(padCell("▲ High", 14))
		}
	}
	fmt.Fprintf(out, "  %s %s %s %s %s %s %s\n", padCell(r.NodeAlias, 16), resStyleHeader.Render(padCell(r.Host, 20)), cpu, ram, status, bench, resStyleDim.Render(padCell(formatDurationShort(r.Latency), 10)))
}

func resPercentStyle(val, bench int) lipgloss.Style {
	if val < bench {
		return resStyleIdle
	}
	if val > 70 {
		return resStyleHigh
	}
	return resStyleNormal
}

func printNodesResourcesHelp() *apperror.AppError {
	fmt.Print("\nUsage: gitmap nodes resources [flags] [target]\n\nQuery CPU & RAM utilization across fleet nodes.\n\nAliases: gitmap nodes resources, top, metrics, res, usage\n\nFlags:\n  -b, --benchmark <N>  Utilization benchmark limit (default: 30%)\n      --open-only      Display only online nodes\n  -j, --json           Output JSON telemetry\n  -h, --help           Show help message\n\n")
	return nil
}
