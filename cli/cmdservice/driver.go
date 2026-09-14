package cmdservice

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ResolveServiceDriver returns platform-specific service driver.
func ResolveServiceDriver() ServiceDriver {
	switch runtime.GOOS {
	case constants.OSWindows:
		return &WindowsServiceDriver{}
	case constants.OSDarwin:
		return &DarwinServiceDriver{}
	default:
		return &LinuxServiceDriver{}
	}
}
