// Package constants — os_targets.go defines supported target operating system identifiers.
package constants

const (
	// Target OS identifiers for installers
	OSTargetWin    = "win"
	OSTargetUbuntu = "ubuntu"
	OSTargetCentOS = "centos"
	OSTargetDebian = "debian"
	OSTargetArch   = "arch-linux"
	OSTargetAll    = "all"
)

// SupportedOSTargets returns the list of all supported installer operating systems.
var SupportedOSTargets = []string{
	OSTargetWin,
	OSTargetUbuntu,
	OSTargetCentOS,
	OSTargetDebian,
	OSTargetArch,
	OSTargetAll,
}
