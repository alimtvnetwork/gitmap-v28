package cmdos

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// HardwareMetricsTTL defines the duration hardware telemetry is cached across processes.
const HardwareMetricsTTL = 5 * time.Second

var (
	hwCacheLock   sync.RWMutex
	hwMemoryCache *HardwareMetrics
)

func isMemoryCacheValid() bool {
	hasCache := hwMemoryCache != nil
	if hasCache {
		isFresh := time.Since(hwMemoryCache.Timestamp) < HardwareMetricsTTL
		return isFresh
	}
	return false
}

// GetHardwareMetrics returns cached hardware metrics or triggers a fresh platform probe.
func GetHardwareMetrics() (HardwareMetrics, error) {
	hwCacheLock.RLock()
	isMemValid := isMemoryCacheValid()
	if isMemValid {
		cached := *hwMemoryCache
		cached.IsCached = true
		hwCacheLock.RUnlock()
		return cached, nil
	}
	hwCacheLock.RUnlock()

	hwCacheLock.Lock()
	defer hwCacheLock.Unlock()

	isMemValidSecond := isMemoryCacheValid()
	if isMemValidSecond {
		cached := *hwMemoryCache
		cached.IsCached = true
		return cached, nil
	}

	diskMetrics, hasDiskCache := readDiskHWCache()
	if hasDiskCache {
		hwMemoryCache = &diskMetrics
		cached := diskMetrics
		cached.IsCached = true
		return cached, nil
	}

	start := time.Now()
	fresh, err := executePlatformHWProbe()
	hasProbeErr := err != nil
	if hasProbeErr {
		return HardwareMetrics{}, apperror.WrapSimple(err, "GetHardwareMetrics.executePlatformHWProbe")
	}
	fresh.Timestamp = time.Now()
	fresh.DurationMs = time.Since(start).Milliseconds()
	fresh.IsCached = false

	hwMemoryCache = &fresh
	_ = writeDiskHWCache(fresh)

	return fresh, nil
}

func resolveHWCacheFilePath() string {
	return filepath.Join(os.TempDir(), "gitmap_hw_metrics_cache.json")
}

func readDiskHWCache() (HardwareMetrics, bool) {
	cachePath := resolveHWCacheFilePath()
	data, err := os.ReadFile(cachePath)
	hasReadErr := err != nil
	if hasReadErr {
		return HardwareMetrics{}, false
	}

	var env HardwareCacheEnvelope
	unmarshalErr := json.Unmarshal(data, &env)
	hasUnmarshalErr := unmarshalErr != nil
	if hasUnmarshalErr {
		return HardwareMetrics{}, false
	}

	isExpired := time.Since(env.CachedAt) >= HardwareMetricsTTL
	if isExpired {
		return HardwareMetrics{}, false
	}

	metrics := env.Metrics
	metrics.Timestamp = env.CachedAt
	return metrics, true
}

func writeDiskHWCache(m HardwareMetrics) error {
	cachePath := resolveHWCacheFilePath()
	env := HardwareCacheEnvelope{
		CachedAt:   m.Timestamp,
		TTLSeconds: int(HardwareMetricsTTL.Seconds()),
		Metrics:    m,
	}
	data, err := json.Marshal(env)
	hasMarshalErr := err != nil
	if hasMarshalErr {
		return apperror.WrapSimple(err, "writeDiskHWCache.marshal")
	}

	tmpPath := fmt.Sprintf("%s.%d.tmp", cachePath, time.Now().UnixNano())
	writeErr := os.WriteFile(tmpPath, data, 0600)
	hasWriteErr := writeErr != nil
	if hasWriteErr {
		return apperror.WrapSimple(writeErr, "writeDiskHWCache.writeFile")
	}

	renameErr := os.Rename(tmpPath, cachePath)
	hasRenameErr := renameErr != nil
	if hasRenameErr {
		_ = os.Remove(tmpPath)
		return apperror.WrapSimple(renameErr, "writeDiskHWCache.rename")
	}
	return nil
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	hasErr := err != nil
	if hasErr {
		return apperror.WrapSimple(err, "printJSON")
	}
	fmt.Println(string(data))
	return nil
}

func runCPUCLI(args []string) error {
	checkHelp("cpu", args)
	isJSON := hasOSJSONArg(args)
	metrics, err := GetHardwareMetrics()
	hasErr := err != nil
	if hasErr {
		return err
	}
	if isJSON {
		return printJSON(metrics.CPU)
	}
	renderCPUView(metrics.CPU, metrics.IsCached, metrics.DurationMs)
	return nil
}

func runRAMCLI(args []string) error {
	checkHelp("ram", args)
	isJSON := hasOSJSONArg(args)
	metrics, err := GetHardwareMetrics()
	hasErr := err != nil
	if hasErr {
		return err
	}
	if isJSON {
		return printJSON(metrics.RAM)
	}
	renderRAMView(metrics.RAM, metrics.IsCached, metrics.DurationMs)
	return nil
}

func runDiskCLI(args []string) error {
	checkHelp("disk", args)
	isJSON := hasOSJSONArg(args)
	metrics, err := GetHardwareMetrics()
	hasErr := err != nil
	if hasErr {
		return err
	}
	if isJSON {
		return printJSON(metrics.Disks)
	}
	renderDiskView(metrics.Disks, metrics.IsCached, metrics.DurationMs)
	return nil
}
