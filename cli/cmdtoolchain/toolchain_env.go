package cmdtoolchain

import (
	"os"
	"path/filepath"
	"strings"
)

// EnrichToolchainEnv augments PATH and LIBCLANG_PATH for child toolchain commands.
func EnrichToolchainEnv(baseBin string) []string {
	return EnrichEnv(os.Environ(), baseBin, ResolveLibclangDir())
}

// EnrichEnv augments environment variables with base binary dir and libclang.
func EnrichEnv(baseEnv []string, baseBin, libclangDir string) []string {
	envMap := sliceToEnvMap(baseEnv)
	ensureSystemEnvVars(envMap)
	augmentPathVar(envMap, baseBin, libclangDir)
	injectLibclangPath(envMap, libclangDir)

	return envMapToSlice(envMap)
}

func sliceToEnvMap(env []string) map[string]string {
	m := make(map[string]string, len(env)+4)
	for _, kv := range env {
		idx := strings.IndexByte(kv, '=')
		if idx > 0 {
			m[kv[:idx]] = kv[idx+1:]
		}
	}

	return m
}

func envMapToSlice(m map[string]string) []string {
	res := make([]string, 0, len(m))
	for k, v := range m {
		res = append(res, k+"="+v)
	}

	return res
}

func ensureSystemEnvVars(m map[string]string) {
	required := []string{"SYSTEMROOT", "COMSPEC", "TMP", "TEMP", "USERPROFILE"}
	for _, key := range required {
		hasKey := hasEnvKey(m, key)
		if !hasKey {
			setEnvFromSystem(m, key)
		}
	}
}

func hasEnvKey(m map[string]string, target string) bool {
	for k := range m {
		if strings.EqualFold(k, target) {
			return true
		}
	}

	return false
}

func setEnvFromSystem(m map[string]string, key string) {
	val := os.Getenv(key)
	if val != "" {
		m[key] = val
	}
}

func augmentPathVar(m map[string]string, baseBin, libclangDir string) {
	currentPath := getPathFromMap(m)
	sep := string(os.PathListSeparator)
	prefix := buildPathPrefix(baseBin, libclangDir, sep)
	if prefix == "" {
		return
	}

	newPath := prefix + sep + currentPath
	setPathInMap(m, newPath)
}

func buildPathPrefix(baseBin, libclangDir, sep string) string {
	parts := make([]string, 0, 2)
	if baseBin != "" {
		parts = append(parts, filepath.Dir(baseBin))
	}

	if libclangDir != "" {
		parts = append(parts, libclangDir)
	}

	return strings.Join(parts, sep)
}

func getPathFromMap(m map[string]string) string {
	for k, v := range m {
		if strings.EqualFold(k, "PATH") {
			return v
		}
	}

	return ""
}

func setPathInMap(m map[string]string, newPath string) {
	for k := range m {
		if strings.EqualFold(k, "PATH") {
			m[k] = newPath

			return
		}
	}

	m["PATH"] = newPath
}

func injectLibclangPath(m map[string]string, libclangDir string) {
	if libclangDir == "" {
		return
	}

	for k := range m {
		if strings.EqualFold(k, "LIBCLANG_PATH") {
			m[k] = libclangDir

			return
		}
	}

	m["LIBCLANG_PATH"] = libclangDir
}
