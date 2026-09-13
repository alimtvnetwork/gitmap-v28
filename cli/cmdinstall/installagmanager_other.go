//go:build !windows

package cmdinstall

func isAgManagerWindowsInstalled() bool {
	return false
}
