package osutil

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/startup"
)

// RegisterStartup registers the command for OS login startup via the startup package.
func RegisterStartup(binPath string) error {
	return startup.RegisterMacroAutostart(binPath)
}
