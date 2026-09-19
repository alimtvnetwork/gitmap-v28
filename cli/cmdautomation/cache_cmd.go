package cmdautomation

import (
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunCacheStatus displays current in-memory cache capacity and stats.
func RunCacheStatus() *apperror.AppError {
	stats := GlobalCache().Stats()
	mb := float64(stats.TotalBytes) / (1024 * 1024)

	fmt.Printf("\n%s[GitMap In-Memory Automation Cache]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Cached Files:   %d\n", stats.TotalFiles)
	fmt.Printf("  Memory Usage:   %.2f MB (0 disk bytes)\n", mb)
	fmt.Printf("  Cache Hits:     %d\n", stats.Hits)
	fmt.Printf("  Cache Misses:   %d\n\n", stats.Misses)
	return nil
}

// RunCacheRead retrieves a file from memory cache or warms and returns it.
func RunCacheRead(path string) *apperror.AppError {
	start := time.Now()
	data, isHit := GlobalCache().GetFile(path)
	if isHit {
		renderReadResult(path, len(data), true, time.Since(start))
		return nil
	}

	diskData, err := os.ReadFile(path)
	if err != nil {
		ctx := map[string]any{"path": path, "err": err.Error()}
		return apperror.New("cache_read", "E_FILE_NOT_FOUND", ctx)
	}

	GlobalCache().SetFile(path, diskData)
	renderReadResult(path, len(diskData), false, time.Since(start))
	return nil
}

func renderReadResult(path string, bytesLen int, isHit bool, dur time.Duration) {
	tag := constants.ColorGreen + "HIT (memory)" + constants.ColorReset
	if !isHit {
		tag = constants.ColorYellow + "MISS (disk load)" + constants.ColorReset
	}
	fmt.Printf("\nRead: %s | Size: %d B | Status: %s | Latency: %s\n\n", path, bytesLen, tag, dur)
}

// RunCacheWarm pre-populates the in-memory cache for a directory.
func RunCacheWarm(dir string) *apperror.AppError {
	start := time.Now()
	targetDir := dir
	if targetDir == "" {
		targetDir = "."
	}

	count := GlobalCache().Warm(targetDir)
	dur := time.Since(start)
	fmt.Printf("\n%s✔ Warmed in-memory cache:%s %d files in %s (0 bytes temp disk usage)\n\n",
		constants.ColorGreen, constants.ColorReset, count, dur)
	return nil
}

// RunCacheClear purges the in-memory cache.
func RunCacheClear() *apperror.AppError {
	GlobalCache().Clear()
	fmt.Printf("\n%s✔ Purged automation in-memory cache.%s All memory reclaimed.\n\n",
		constants.ColorGreen, constants.ColorReset)
	return nil
}
