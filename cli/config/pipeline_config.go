package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// ResolveActiveConfig loads config.json from local repo, AppData, or user home, falling back to defaults.
func ResolveActiveConfig() model.Config {
	for _, candidate := range CandidateConfigPaths() {
		cfg, isLoaded := tryLoadCandidateConfig(candidate)
		if isLoaded {
			return cfg
		}
	}

	return model.DefaultConfig()
}

func tryLoadCandidateConfig(path string) (model.Config, bool) {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		return model.Config{}, false
	}
	cfg, loadErr := LoadFromFile(path)
	if loadErr != nil {
		return model.Config{}, false
	}

	return cfg, true
}

// CandidateConfigPaths returns discovery paths for config.json.
func CandidateConfigPaths() []string {
	var paths []string
	paths = append(paths, "./data/config.json", ".gitmap/config.json")
	if appData := os.Getenv("LOCALAPPDATA"); len(appData) > 0 {
		paths = append(paths, filepath.Join(appData, "gitmap-cli", "data", "config.json"))
	}
	if home, err := os.UserHomeDir(); err == nil && len(home) > 0 {
		paths = append(paths, filepath.Join(home, ".gitmap", "config.json"))
	}

	return paths
}

func parseEnvMaxDbBytes(envVal string) int64 {
	if len(envVal) == 0 {
		return 0
	}
	mb, err := strconv.ParseInt(envVal, 10, 64)
	if err != nil || mb <= 0 {
		return 0
	}

	return mb * 1024 * 1024
}

// ResolvePipelineMaxDbBytes returns the pipeline DB size ceiling in bytes.
// It checks GITMAP_PIPELINE_MAX_DB_MB, then config.json's pipelineMaxDbSizeMb, defaulting to 10 MB.
func ResolvePipelineMaxDbBytes() int64 {
	const defaultMaxBytes = int64(10 * 1024 * 1024)
	envBytes := parseEnvMaxDbBytes(os.Getenv("GITMAP_PIPELINE_MAX_DB_MB"))
	if envBytes > 0 {
		return envBytes
	}

	cfg := ResolveActiveConfig()
	if cfg.PipelineMaxDbSizeMB > 0 {
		return int64(cfg.PipelineMaxDbSizeMB) * 1024 * 1024
	}

	return defaultMaxBytes
}

func parseEnvCacheTTL(envVal string) time.Duration {
	if len(envVal) == 0 {
		return 0
	}
	sec, err := strconv.Atoi(envVal)
	if err != nil || sec <= 0 {
		return 0
	}

	return time.Duration(sec) * time.Second
}

// ResolvePipelineCacheTTL returns the cache TTL window duration.
// It checks GITMAP_PIPELINE_CACHE_TTL_SEC, then config.json's pipelineCacheTtlSec, defaulting to 5 seconds.
func ResolvePipelineCacheTTL() time.Duration {
	const defaultTTL = 5 * time.Second
	envTTL := parseEnvCacheTTL(os.Getenv("GITMAP_PIPELINE_CACHE_TTL_SEC"))
	if envTTL > 0 {
		return envTTL
	}

	cfg := ResolveActiveConfig()
	if cfg.PipelineCacheTTLSec > 0 {
		return time.Duration(cfg.PipelineCacheTTLSec) * time.Second
	}

	return defaultTTL
}
