package usercontext

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func formatGhAuthLine(gh GhAuthInfo) string {
	if gh.IsAuthorized() {
		return fmt.Sprintf("%s● GitHub CLI:%s    %s✔ Authorized as %s%s",
			constants.ColorPastelCyan, constants.ColorReset,
			constants.ColorPastelGreen, gh.Username, constants.ColorReset)
	}
	if gh.Status == StatusToolMissing {
		return fmt.Sprintf("%s● GitHub CLI:%s    %s✖ gh tool not installed%s",
			constants.ColorPastelCyan, constants.ColorReset,
			constants.ColorDim, constants.ColorReset)
	}

	return fmt.Sprintf("%s● GitHub CLI:%s    %s✖ Not logged in%s",
		constants.ColorPastelCyan, constants.ColorReset,
		constants.ColorYellow, constants.ColorReset)
}

func formatActiveGitLine(summary UserSummary) string {
	user, scope := resolveActiveIdentity(summary)
	if !user.IsConfigured {
		return fmt.Sprintf("%s● Active Git:%s     %s(not configured)%s",
			constants.ColorPastelCyan, constants.ColorReset,
			constants.ColorDim, constants.ColorReset)
	}

	return fmt.Sprintf("%s● Active Git:%s     %s%s <%s>%s %s(%s)%s",
		constants.ColorPastelCyan, constants.ColorReset,
		constants.ColorWhite, user.Name, user.Email, constants.ColorReset,
		constants.ColorDim, scope, constants.ColorReset)
}

func resolveActiveIdentity(summary UserSummary) (GitIdentity, string) {
	if summary.IsInsideRepo && summary.LocalGitUser.IsConfigured {
		return summary.LocalGitUser, "local"
	}

	return summary.GlobalGitUser, "global"
}

func formatGlobalGitLine(global GitIdentity) string {
	if !global.IsConfigured {
		return fmt.Sprintf("%s● Global Git:%s     %s(not configured)%s",
			constants.ColorPastelCyan, constants.ColorReset,
			constants.ColorDim, constants.ColorReset)
	}

	return fmt.Sprintf("%s● Global Git:%s     %s%s <%s>%s",
		constants.ColorPastelCyan, constants.ColorReset,
		constants.ColorWhite, global.Name, global.Email, constants.ColorReset)
}

func formatProfileLine(profile string) string {
	val := profile
	if len(val) == 0 {
		val = "default"
	}

	return fmt.Sprintf("%s● GitMap Profile:%s %s%s%s",
		constants.ColorPastelCyan, constants.ColorReset,
		constants.ColorPastelYellow, val, constants.ColorReset)
}

func formatBindingLine(bound string) string {
	val := bound
	if len(val) == 0 {
		val = "none"
	}

	return fmt.Sprintf("%s● Project Bound:%s  %s%s%s",
		constants.ColorPastelCyan, constants.ColorReset,
		constants.ColorPastelMagenta, val, constants.ColorReset)
}
