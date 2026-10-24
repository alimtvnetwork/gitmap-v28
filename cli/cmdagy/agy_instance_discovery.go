package cmdagy

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type diskInstancesFile struct {
	ActiveInstanceID string              `json:"active_instance_id"`
	Instances        []diskInstanceEntry `json:"instances"`
}

type diskInstanceEntry struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	DataDir        string  `json:"data_dir"`
	ExecutablePath *string `json:"executable_path"`
	ExtensionsDir  *string `json:"extensions_dir"`
	BoundAccountID string  `json:"bound_account_id"`
	BoundEmail     string  `json:"bound_email"`
	CreatedAt      int64   `json:"created_at"`
	LastUsed       int64   `json:"last_used"`
	IsDefault      bool    `json:"is_default"`
	PID            *int    `json:"pid"`
	SeqNum         int     `json:"seq_num"`
}

// DiscoverAllAgyInstances enumerates primary, registered, and active Antigravity instances.
func DiscoverAllAgyInstances() ([]AgyInstanceInfo, error) {
	var instances []AgyInstanceInfo
	seenIDs := make(map[string]bool)

	primaryInst, err := discoverPrimaryInstance()
	if err == nil {
		instances = append(instances, primaryInst)
		seenIDs[primaryInst.InstanceID] = true
		seenIDs["default"] = true
	}

	diskInstances := loadInstancesFromDiskRegistries()
	procInstances := scanRunningProcessesForInstances()

	for _, disk := range diskInstances {
		if seenIDs[disk.InstanceID] {
			continue
		}
		seenIDs[disk.InstanceID] = true

		for _, p := range procInstances {
			isMatch := p.InstanceID == disk.InstanceID || (disk.ConfigDir != "" && strings.Contains(p.ConfigDir, disk.ConfigDir))
			if isMatch {
				mergeProcMatch(&disk, p)
				break
			}
		}

		if disk.SummariesDBPath != "" {
			disk.ActiveWorkspaces = queryDistinctWorkspacesFromDB(disk.SummariesDBPath)
		}
		if len(disk.ActiveWorkspaces) == 0 && primaryInst.SummariesDBPath != "" {
			disk.ActiveWorkspaces = primaryInst.ActiveWorkspaces
		}

		instances = append(instances, disk)
	}

	for _, proc := range procInstances {
		if seenIDs[proc.InstanceID] {
			continue
		}
		seenIDs[proc.InstanceID] = true
		if proc.SummariesDBPath != "" {
			proc.ActiveWorkspaces = queryDistinctWorkspacesFromDB(proc.SummariesDBPath)
		}
		instances = append(instances, proc)
	}

	return instances, nil
}

func discoverPrimaryInstance() (AgyInstanceInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return AgyInstanceInfo{}, apperror.WrapSimple(err, "discover_primary_instance")
	}

	configDir := filepath.Join(home, ".gemini", "antigravity")
	brainDir := filepath.Join(configDir, "brain")
	summariesDB := filepath.Join(configDir, "conversation_summaries.db")

	inst := AgyInstanceInfo{
		InstanceID:       "primary",
		InstanceName:     "Default Antigravity Workspace",
		ConfigDir:        configDir,
		BrainDir:         brainDir,
		SummariesDBPath:  summariesDB,
		ActiveWorkspaces: queryDistinctWorkspacesFromDB(summariesDB),
		IsPrimary:        true,
		IsRunning:        false,
		LastActiveAt:     time.Now().UTC().Format(time.RFC3339),
	}

	procs := scanRunningProcessesForInstances()
	for _, p := range procs {
		if p.IsPrimary || strings.Contains(p.ConfigDir, ".config/Antigravity") || p.ConfigDir == configDir {
			inst.IsRunning = true
			inst.ProcessID = p.ProcessID
			inst.LanguageServer = p.LanguageServer
			break
		}
	}

	if !inst.IsRunning && len(procs) > 0 {
		attachFirstRunningProcess(&inst, procs)
	}

	return inst, nil
}

func mergeProcMatch(disk *AgyInstanceInfo, p AgyInstanceInfo) {
	disk.IsRunning = true
	if p.ProcessID > 0 {
		disk.ProcessID = p.ProcessID
	}
	if p.LanguageServer != "" {
		disk.LanguageServer = p.LanguageServer
	}
}

func attachFirstRunningProcess(inst *AgyInstanceInfo, procs []AgyInstanceInfo) {
	for _, p := range procs {
		if p.ProcessID > 0 {
			inst.IsRunning = true
			inst.ProcessID = p.ProcessID
			inst.LanguageServer = p.LanguageServer
			break
		}
	}
}

func loadInstancesFromDiskRegistries() []AgyInstanceInfo {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	candidateFiles := []string{
		filepath.Join(home, ".antigravity_tools", "instances", "instances.json"),
		filepath.Join(home, ".gemini", "antigravity", "instances.json"),
	}

	var results []AgyInstanceInfo
	for _, candidate := range candidateFiles {
		data, readErr := os.ReadFile(candidate)
		if readErr != nil {
			continue
		}

		var filePayload diskInstancesFile
		if jsonErr := json.Unmarshal(data, &filePayload); jsonErr != nil {
			continue
		}

		for _, item := range filePayload.Instances {
			if item.ID == "default" || item.ID == "primary" {
				continue
			}

			brainDir := filepath.Join(item.DataDir, "brain")
			if _, statErr := os.Stat(brainDir); statErr != nil {
				brainDir = filepath.Join(home, ".gemini", "antigravity", "brain")
			}

			summariesDB := filepath.Join(item.DataDir, "conversation_summaries.db")
			if _, statErr := os.Stat(summariesDB); statErr != nil {
				summariesDB = filepath.Join(home, ".gemini", "antigravity", "conversation_summaries.db")
			}

			pid := 0
			if item.PID != nil {
				pid = *item.PID
			}

			lastActive := time.Now().UTC().Format(time.RFC3339)
			if item.LastUsed > 0 {
				lastActive = time.Unix(item.LastUsed, 0).UTC().Format(time.RFC3339)
			}

			info := AgyInstanceInfo{
				InstanceID:      item.ID,
				InstanceName:    item.Name,
				ProcessID:       pid,
				ConfigDir:       item.DataDir,
				BrainDir:        brainDir,
				SummariesDBPath: summariesDB,
				IsPrimary:       item.IsDefault,
				IsRunning:       isProcessAlive(pid),
				LastActiveAt:    lastActive,
			}
			results = append(results, info)
		}
	}

	return results
}

func scanRunningProcessesForInstances() []AgyInstanceInfo {
	if runtime.GOOS == "windows" {
		return scanRunningProcessesWindows()
	}
	return scanRunningProcessesUnix()
}

func scanRunningProcessesUnix() []AgyInstanceInfo {
	cmd := exec.Command("ps", "-eo", "pid,args")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var found []AgyInstanceInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, "antigravity") && !strings.Contains(trimmed, "language_server") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}

		pid, parseErr := strconv.Atoi(parts[0])
		if parseErr != nil {
			continue
		}

		info := parseProcessCommandLine(pid, trimmed)
		if info != nil {
			found = append(found, *info)
		}
	}

	return found
}

func scanRunningProcessesWindows() []AgyInstanceInfo {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-CimInstance Win32_Process | Where-Object { $_.Name -like '*antigravity*' -or $_.CommandLine -like '*language_server*' } | Select-Object ProcessId, CommandLine | ConvertTo-Json -Compress")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	items := unmarshalWinProcItems(output)

	var found []AgyInstanceInfo
	for _, it := range items {
		info := parseProcessCommandLine(it.ProcessID, it.CommandLine)
		if info != nil {
			found = append(found, *info)
		}
	}

	return found
}

type winProcItem struct {
	ProcessID   int    `json:"ProcessId"`
	CommandLine string `json:"CommandLine"`
}
