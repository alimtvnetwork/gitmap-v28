//go:build !windows

package cmdinstall

func detectWindowsVersion() OSType {
	return OSUnknown
}
