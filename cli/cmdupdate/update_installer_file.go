package cmdupdate

import "runtime"

func hasWindowsRuntime() bool {
	return runtime.GOOS == "windows"
}
