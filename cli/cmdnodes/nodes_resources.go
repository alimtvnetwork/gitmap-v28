package cmdnodes

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/charmbracelet/lipgloss"
)

// NodeResourceRecord stores telemetry metrics for a local or remote fleet node.
type NodeResourceRecord struct {
	NodeAlias      string        `json:"nodeAlias"`
	Host           string        `json:"host"`
	OS             string        `json:"os"`
	IsLocal        bool          `json:"isLocal"`
	IsOnline       bool          `json:"isOnline"`
	CPUPercent     int           `json:"cpuPercent"`
	RAMPercent     int           `json:"ramPercent"`
	TotalRAMMB     int           `json:"totalRamMb"`
	UsedRAMMB      int           `json:"usedRamMb"`
	BenchmarkLimit int           `json:"benchmarkLimit"`
	IsUnderLoaded  bool          `json:"isUnderLoaded"`
	LoadStatus     string        `json:"loadStatus"`
	Benchmark      string        `json:"benchmark"`
	Latency        time.Duration `json:"latency"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

type nodesResourcesOptions struct {
	benchmarkLimit int
	isJSON         bool
	isOpenOnly     bool
	targetFilter   string
}

var (
	resStyleTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	resStyleHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	resStyleDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	resStyleIdle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	resStyleNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c"))
	resStyleHigh   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
	resStyleExceed = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
)

// RunNodesResources executes fleet resource monitoring across local and remote nodes.
func RunNodesResources(args []string) *apperror.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return printNodesResourcesHelp()
		}
	}
	opts := parseNodesResourcesOptions(args)
	conns, _ := cmdssh.FetchAllSSHConnections()
	records := executeFleetResourcesProbe(conns, opts)
	if opts.isJSON {
		return printResourcesJSON(records)
	}
	return renderResourcesTable(os.Stdout, records, opts)
}

// RunNodesResourcesCLI wraps RunNodesResources to conform to standard error return signatures.
func RunNodesResourcesCLI(args []string) error {
	if err := RunNodesResources(args); err != nil {
		return err
	}
	return nil
}

func parseNodesResourcesOptions(args []string) (opts nodesResourcesOptions) {
	opts.benchmarkLimit = resolveDefaultBenchmarkLimit()
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" || a == "-j" {
			opts.isJSON = true
		} else if a == "--open-only" || a == "--online-only" {
			opts.isOpenOnly = true
		} else if (a == "--benchmark" || a == "-b" || a == "--threshold") && i+1 < len(args) {
			i++
			if v, err := strconv.Atoi(args[i]); err == nil && v > 0 {
				opts.benchmarkLimit = v
			}
		} else if isFilterCandidate(a, opts.targetFilter) {
			opts.targetFilter = a
		}
	}
	return opts
}

func isFilterCandidate(arg, current string) bool {
	tok := strings.ToLower(arg)
	if strings.HasPrefix(arg, "-") || current != "" || tok == "resources" || tok == "res" || tok == "top" || tok == "metrics" || tok == "usage" || tok == "nodes" {
		return false
	}
	return true
}

func resolveDefaultBenchmarkLimit() int {
	if dbConn, err := store.OpenDefault(); err == nil {
		defer dbConn.Close()
		if val, err := dbConn.GetSetting(constants.SettingNodeResourceThreshold); err == nil && val != "" {
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				return n
			}
		}
	}
	return constants.DefaultNodeResourceThreshold
}

func executeFleetResourcesProbe(conns []db.SSHConnection, opts nodesResourcesOptions) []NodeResourceRecord {
	var records []NodeResourceRecord
	var mu sync.Mutex
	var wg sync.WaitGroup
	if isTargetMatching("local", "127.0.0.1", opts.targetFilter) {
		host := "local"
		if h, err := os.Hostname(); err == nil && h != "" {
			host = h
		}
		records = append(records, probeNodeExec(host, "127.0.0.1", runtime.GOOS == "windows", true, opts.benchmarkLimit, nil))
	}
	for _, c := range conns {
		if isTargetMatching(c.Alias, c.IPAddress, opts.targetFilter) {
			wg.Add(1)
			go func(conn db.SSHConnection) {
				defer wg.Done()
				alias, port := conn.Alias, conn.Port
				if alias == "" {
					alias = conn.IPAddress
				}
				if port <= 0 {
					port = 22
				}
				rec := probeNodeExec(alias, fmt.Sprintf("%s:%d", conn.IPAddress, port), isWindowsNode(conn), false, opts.benchmarkLimit, &conn)
				mu.Lock()
				records = append(records, rec)
				mu.Unlock()
			}(c)
		}
	}
	wg.Wait()
	if opts.isOpenOnly {
		n := 0
		for _, r := range records {
			if r.IsOnline {
				records[n] = r
				n++
			}
		}
		return records[:n]
	}
	return records
}

func isTargetMatching(alias, host, filter string) bool {
	if filter != "" {
		f := strings.ToLower(filter)
		return strings.Contains(strings.ToLower(alias), f) || strings.Contains(strings.ToLower(host), f)
	}
	return true
}

func probeNodeExec(alias, host string, isWin, isLocal bool, benchmark int, conn *db.SSHConnection) NodeResourceRecord {
	start := time.Now()
	nodeOS, cmd, shell := "linux", `sh -c 'cpu=$(grep "^cpu " /proc/stat 2>/dev/null | awk "{u=\$2+\$4; t=\$2+\$4+\$5; if(t>0) print int(u*100/t); else print 0}"); eval $(awk "/MemTotal/ {t=\$2} /MemAvailable/ {a=\$2} END {u=t-a; if(t>0) r=int(u*100/t); else r=0; printf \"ram=%d; tot=%d; used=%d\n\", r, int(t/1024), int(u/1024)}" /proc/meminfo 2>/dev/null); echo "$cpu,$ram,$tot,$used"'`, "bash"
	if isWin {
		nodeOS, cmd, shell = "windows", `$c=(Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average; $os=Get-CimInstance Win32_OperatingSystem; $tot=[math]::Round($os.TotalVisibleMemorySize/1024); $free=[math]::Round($os.FreePhysicalMemory/1024); $used=$tot-$free; $ram=[math]::Round(($used/$tot)*100); Write-Output ('{0},{1},{2},{3}' -f [int]$c, [int]$ram, [int]$tot, [int]$used)`, "powershell"
	}
	rec := NodeResourceRecord{NodeAlias: alias, Host: host, OS: nodeOS, IsLocal: isLocal, BenchmarkLimit: benchmark}
	var out string
	var err error
	if conn == nil {
		bin, flag := "sh", "-c"
		if isWin {
			bin, flag = "powershell", "-Command"
		}
		var b []byte
		b, err = exec.Command(bin, flag, cmd).Output()
		out = string(b)
	} else {
		out, err = currentSSHExecutor.Execute(*conn, cmd, shell)
	}
	rec.Latency = time.Since(start)
	if err != nil {
		rec.LoadStatus, rec.Benchmark, rec.ErrorMessage = "Offline", "N/A", err.Error()
		return rec
	}
	cpu, ram, tot, used, errParse := parseCSV(out)
	if errParse != nil {
		rec.IsOnline, rec.LoadStatus, rec.Benchmark, rec.ErrorMessage = true, "Unknown", "N/A", errParse.Error()
		return rec
	}
	rec.IsOnline, rec.CPUPercent, rec.RAMPercent, rec.TotalRAMMB, rec.UsedRAMMB = true, cpu, ram, tot, used
	rec.IsUnderLoaded = (cpu < benchmark) && (ram < benchmark)
	rec.LoadStatus, rec.Benchmark = "Normal", "Exceeded"
	if rec.IsUnderLoaded {
		rec.LoadStatus, rec.Benchmark = "Idle", "Passed"
	} else if cpu > 70 || ram > 70 {
		rec.LoadStatus = "High"
	}
	return rec
}

func parseCSV(raw string) (cpu, ram, tot, used int, err error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	parts := strings.Split(strings.TrimSpace(lines[len(lines)-1]), ",")
	if len(parts) >= 4 {
		c, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		r, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		t, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		u, _ := strconv.Atoi(strings.TrimSpace(parts[3]))
		return c, r, t, u, nil
	}
	return 0, 0, 0, 0, fmt.Errorf("invalid metrics")
}

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
	fmt.Fprintf(out, "\n  %s╔%s╗%s\n  %s║%s%s%s%s║%s\n  %s╚%s╝%s\n  Probing %d node(s) (Benchmark: %d%%)...\n\n  %s\n  %s\n",
		constants.ColorDim, border, constants.ColorReset, constants.ColorDim, constants.ColorReset, resStyleTitle.Render(padCell(" GITMAP FLEET NODES RESOURCE METRICS & BENCHMARK", 116)), constants.ColorDim, constants.ColorReset, constants.ColorDim, border, constants.ColorReset, len(records), opts.benchmarkLimit, hdr, div)

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
