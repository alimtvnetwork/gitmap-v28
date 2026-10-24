//go:build !windows

package cmdinstall

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall/enums"

func detectWindowsVersion() enums.OSType {
	return enums.OSUnknown
}
