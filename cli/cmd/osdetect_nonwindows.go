//go:build !windows

package cmd

func detectWindowsVersion() OSType {
	return OSUnknown
}
