package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isSSHViewSub(sub string) bool {
	return sub == constants.SubCmdSSHCat || sub == constants.SubCmdSSHView || sub == constants.SubCmdSSHViewS
}

func isSSHCopySub(sub string) bool {
	return sub == constants.SubCmdSSHCopy || sub == constants.SubCmdSSHCopyS
}

func isSSHListSub(sub string) bool {
	return sub == constants.SubCmdSSHList || sub == constants.SubCmdSSHListS
}

func isSSHDeleteSub(sub string) bool {
	return sub == constants.SubCmdSSHDelete || sub == constants.SubCmdSSHDeleteS
}

func isSSHStatusSub(sub string) bool {
	return sub == constants.SubCmdSSHStatus || sub == constants.SubCmdSSHStatusS
}

func dispatchConfigOrStatus(sub string, args []string) bool {
	if sub == constants.SubCmdSSHConfig {
		runSSHConfig(args)
		return true
	}
	if isSSHStatusSub(sub) {
		runSSHStatus(args)
		return true
	}
	return false
}

func dispatchOtherFallbackSSH(sub string, args []string) bool {
	if isSSHListSub(sub) {
		runSSHList(args...)
		return true
	}
	if isSSHDeleteSub(sub) {
		runSSHDelete(args)
		return true
	}
	return dispatchConfigOrStatus(sub, args)
}

func dispatchCreateSSH(sub string, args []string) bool {
	if sub == constants.SubCmdSSHCreate {
		runSSHGenerate(args)
		RenderSSHHelp()
		return true
	}
	return false
}

func dispatchFallbackSSH(sub string, args []string) bool {
	if isSSHViewSub(sub) {
		runSSHCat(args)
		return true
	}
	if isSSHCopySub(sub) {
		runSSHCopy(args)
		return true
	}
	if isCreated := dispatchCreateSSH(sub, args); isCreated {
		return true
	}
	return dispatchOtherFallbackSSH(sub, args)
}
