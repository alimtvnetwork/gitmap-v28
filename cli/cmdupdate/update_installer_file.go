package cmdupdate

import (
	"io"
	"os"
	"runtime"
)

func hasWindowsRuntime() bool {
	return runtime.GOOS == "windows"
}

func writeInstallerTempFile(body io.Reader) (string, error) {
	tmp, err := createInstallerTempFile()
	if err != nil {
		return "", err
	}

	return populateInstallerTempFile(tmp, body)
}

func createInstallerTempFile() (*os.File, error) {
	ext := getInstallerExtension()
	tmp, err := os.CreateTemp("", "gitmap-update-*"+ext)
	if err != nil {
		return nil, err
	}

	writeWindowsInstallerBOM(tmp)

	return tmp, nil
}

func populateInstallerTempFile(tmp *os.File, body io.Reader) (string, error) {
	if _, copyErr := io.Copy(tmp, body); copyErr != nil {
		tmp.Close()
		os.Remove(tmp.Name())

		return "", copyErr
	}

	tmp.Close()
	chmodUnixInstaller(tmp.Name())

	return tmp.Name(), nil
}

func getInstallerExtension() string {
	if hasWindowsRuntime() {
		return ".ps1"
	}

	return ".sh"
}

func writeWindowsInstallerBOM(f *os.File) {
	if hasWindowsRuntime() {
		_, _ = f.Write([]byte{0xEF, 0xBB, 0xBF})
	}
}

func chmodUnixInstaller(path string) {
	if !hasWindowsRuntime() {
		_ = os.Chmod(path, 0o755)
	}
}
