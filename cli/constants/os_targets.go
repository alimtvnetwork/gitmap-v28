// Package constants — os_targets.go defines supported target operating system identifiers.
package constants

const (
	// Target OS identifiers for installers
	OSTargetWin    = "win"
	OSTargetUbuntu = "ubuntu"
	OSTargetCentOS = "centos"
	OSTargetDebian = "debian"
	OSTargetFedora = "fedora"
	OSTargetArch   = "arch-linux"
	OSTargetMac    = "macos"
	OSTargetUnix   = "unix"
	OSTargetAll    = "all"

	// Execution ordering modes for multi-OS installers
	OrderUnixFirst = "unix-first"
	OrderOSFirst   = "os-first"
	OrderOSOnly    = "os-only"
	OrderFallback  = "fallback"
)

// SupportedOSTargets returns the list of all supported installer operating systems.
var SupportedOSTargets = []string{
	OSTargetWin,
	OSTargetUbuntu,
	OSTargetDebian,
	OSTargetCentOS,
	OSTargetFedora,
	OSTargetArch,
	OSTargetMac,
	OSTargetUnix,
	OSTargetAll,
}
