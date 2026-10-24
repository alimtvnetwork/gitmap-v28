package cmdnodes

import (
	"context"
	"fmt"
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

// RunNodesResources executes fleet resource monitoring across local and remote nodes.
func RunNodesResources(args []string) *apperror.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return printNodesResourcesHelp()
		}
	}
	opts := parseNodesResourcesOptions(args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	unifiedNodes, err := collectUnifiedNodes(ctx, true)
	if err != nil {
		return apperror.WrapSimple(err, "RunNodesResources.collectUnifiedNodes")
	}

	records := executeFleetResourcesProbeUnified(ctx, unifiedNodes, opts)
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

func parseBenchmarkArg(val string, fallback int) int {
	v, err := strconv.Atoi(val)
	if err == nil && v > 0 {
		return v
	}
	return fallback
}

func parseNodesResourcesOptions(args []string) (opts nodesResourcesOptions) {
	opts.benchmarkLimit = resolveDefaultBenchmarkLimit()
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" || a == "-j" {
			opts.isJSON = true
			continue
		}
		if a == "--open-only" || a == "--online-only" || a == "--running-only" || a == "-ro" {
			opts.isOpenOnly = true
			continue
		}
		hasBenchmarkFlag := (a == "--benchmark" || a == "-b" || a == "--threshold") && i+1 < len(args)
		if hasBenchmarkFlag {
			i++
			opts.benchmarkLimit = parseBenchmarkArg(args[i], opts.benchmarkLimit)
			continue
		}
		if isFilterCandidate(a, opts.targetFilter) {
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
	dbConn, err := store.OpenDefault()
	if err != nil {
		return constants.DefaultNodeResourceThreshold
	}
	defer dbConn.Close()

	val := dbConn.GetSetting(constants.SettingNodeResourceThreshold)
	n, errParse := strconv.Atoi(val)
	if errParse == nil && n > 0 {
		return n
	}
	return constants.DefaultNodeResourceThreshold
}

func resolveLocalHostname() string {
	h, err := os.Hostname()
	if err == nil && h != "" {
		return h
	}
	return "local"
}

func isNodeTargetIncluded(node UnifiedFleetNode, filter string) bool {
	isLocal := strings.EqualFold(node.Host, "127.0.0.1") || strings.EqualFold(node.Host, constants.Localhost) || strings.EqualFold(node.Alias, "local")
	if isLocal {
		return false
	}
	return isTargetMatching(node.Alias, node.Host, filter)
}

func filterOnlineRecords(records []NodeResourceRecord, isOpenOnly bool) []NodeResourceRecord {
	if isOpenOnly {
		return extractOnlineRecords(records)
	}
	return records
}

func extractOnlineRecords(records []NodeResourceRecord) []NodeResourceRecord {
	var online []NodeResourceRecord
	for _, r := range records {
		if r.IsOnline {
			online = append(online, r)
		}
	}
	return online
}

func executeFleetResourcesProbeUnified(ctx context.Context, nodes []UnifiedFleetNode, opts nodesResourcesOptions) []NodeResourceRecord {
	var records []NodeResourceRecord
	var mu sync.Mutex
	var wg sync.WaitGroup

	host := resolveLocalHostname()
	if isTargetMatching("local", "127.0.0.1", opts.targetFilter) {
		records = append(records, probeNodeExec(host, "127.0.0.1", runtime.GOOS == "windows", true, opts.benchmarkLimit, nil))
	}

	conns, errConns := cmdssh.FetchAllSSHConnections()
	connMap := make(map[string]db.SSHConnection)
	if errConns == nil {
		connMap = indexSSHConnections(conns)
	}

	for _, n := range nodes {
		if opts.isOpenOnly && !n.IsOnline {
			continue
		}
		hasTargetMatch := isNodeTargetIncluded(n, opts.targetFilter)
		if hasTargetMatch {
			wg.Add(1)
			go func(node UnifiedFleetNode) {
				defer wg.Done()
				alias, port, conn := resolveFleetNodeEndpoint(node, connMap)
				isWin := isWindowsNode(conn) || strings.Contains(strings.ToLower(node.Role), "win")
				rec := probeNodeExec(alias, fmt.Sprintf("%s:%d", node.Host, port), isWin, false, opts.benchmarkLimit, &conn)
				mu.Lock()
				records = append(records, rec)
				mu.Unlock()
			}(n)
		}
	}
	wg.Wait()
	return filterOnlineRecords(records, opts.isOpenOnly)
}

func isTargetMatching(alias, host, filter string) bool {
	if filter != "" {
		f := strings.ToLower(filter)
		return strings.Contains(strings.ToLower(alias), f) || strings.Contains(strings.ToLower(host), f)
	}
	return true
}

func resolveLocalShellCmd(isWin bool) (string, string) {
	if isWin {
		return "powershell", "-Command"
	}
	return "sh", "-c"
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
		bin, flag := resolveLocalShellCmd(isWin)
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
	return 0, 0, 0, 0, apperror.NewValidationError("invalid metrics")
}

func indexSSHConnections(conns []db.SSHConnection) map[string]db.SSHConnection {
	connMap := make(map[string]db.SSHConnection, len(conns)*2)
	for _, c := range conns {
		if c.IPAddress != "" {
			connMap[strings.ToLower(c.IPAddress)] = c
		}
		if c.Alias != "" {
			connMap[strings.ToLower(c.Alias)] = c
		}
	}
	return connMap
}
