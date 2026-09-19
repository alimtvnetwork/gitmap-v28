package cmdai

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	cachedRuntime PythonRuntime
	errCached     *apperror.AppError
	detectorOnce  sync.Once
)

// DetectPythonRuntime discovers a local Python 3 runtime with thread-safe memoization.
func DetectPythonRuntime() (PythonRuntime, *apperror.AppError) {
	detectorOnce.Do(func() {
		cachedRuntime, errCached = detectRuntimeInternal()
	})

	hasErr := errCached != nil
	if hasErr {
		return PythonRuntime{}, errCached
	}

	return cachedRuntime, nil
}

// ResetPythonDetector clears memoized detection state for testing or path updates.
func ResetPythonDetector() {
	detectorOnce = sync.Once{}
	cachedRuntime = PythonRuntime{}
	errCached = nil
}

func detectRuntimeInternal() (PythonRuntime, *apperror.AppError) {
	candidates := candidateNames()
	for _, name := range candidates {
		rt, err := probePythonExecutable(name)
		isFound := (err == nil && rt.IsPython3)
		if isFound {
			return rt, nil
		}
	}

	ctx := map[string]any{"candidates": candidates}

	return PythonRuntime{}, apperror.New("detect", "E_PYTHON_NOT_FOUND", ctx)
}

func candidateNames() []string {
	envExe := os.Getenv("PYTHON")
	hasEnv := len(envExe) > 0
	if hasEnv {
		return []string{envExe, "python3", "python", "py"}
	}

	return []string{"python3", "python", "py"}
}

func probePythonExecutable(name string) (PythonRuntime, *apperror.AppError) {
	fullPath, pathErr := resolvePythonPath(name)
	hasPathErr := pathErr != nil
	if hasPathErr {
		return PythonRuntime{}, pathErr
	}

	return probePythonAtLocation(fullPath)
}

func resolvePythonPath(name string) (string, *apperror.AppError) {
	fullPath, lookErr := exec.LookPath(name)
	hasLookErr := lookErr != nil
	if hasLookErr {
		return "", apperror.WrapSimple(lookErr, "look_path")
	}

	return fullPath, nil
}

func probePythonAtLocation(fullPath string) (PythonRuntime, *apperror.AppError) {
	rawOut, probeErr := executeVersionProbe(fullPath)
	hasProbeErr := probeErr != nil
	if hasProbeErr {
		return PythonRuntime{}, probeErr
	}

	versionStr, major, minor, parseErr := parsePythonVersion(rawOut)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return PythonRuntime{}, parseErr
	}

	return buildRuntimeRecord(fullPath, versionStr, major, minor), nil
}

func executeVersionProbe(exePath string) (string, *apperror.AppError) {
	cmd := exec.Command(exePath, "--version")
	out, err := cmd.CombinedOutput()
	hasErr := err != nil
	if hasErr {
		return "", apperror.WrapSimple(err, "probe_version")
	}

	return string(out), nil
}

func extractPythonVersionString(raw string) (string, *apperror.AppError) {
	clean := strings.TrimSpace(raw)
	parts := strings.Fields(clean)
	isMalformed := (len(parts) < 2 || parts[0] != "Python")
	if isMalformed {
		return "", apperror.NewSimple("parse_version", "E_VERSION_PARSE_FAILED")
	}

	return parts[1], nil
}

func parsePythonVersion(raw string) (string, int, int, *apperror.AppError) {
	versionStr, err := extractPythonVersionString(raw)
	hasErr := err != nil
	if hasErr {
		return "", 0, 0, err
	}

	subparts := strings.Split(versionStr, ".")
	isShort := len(subparts) < 2
	if isShort {
		return versionStr, 0, 0, nil
	}

	return parseVersionNumbers(versionStr, subparts)
}

func parseVersionNumbers(versionStr string, subparts []string) (string, int, int, *apperror.AppError) {
	major, majorErr := strconv.Atoi(subparts[0])
	minor, minorErr := strconv.Atoi(subparts[1])
	hasParseErr := (majorErr != nil || minorErr != nil)
	if hasParseErr {
		return "", 0, 0, apperror.NewSimple("parse_digits", "E_VERSION_DIGIT_FAILED")
	}

	return versionStr, major, minor, nil
}

func buildRuntimeRecord(exePath, version string, major, minor int) PythonRuntime {
	isPython3 := major >= 3
	isAvailable := isPython3

	return PythonRuntime{
		ExecutablePath: exePath,
		Version:        version,
		MajorVersion:   major,
		MinorVersion:   minor,
		IsAvailable:    isAvailable,
		IsPython3:      isPython3,
	}
}
