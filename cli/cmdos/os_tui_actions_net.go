package cmdos

import "strings"

func isDNSAction(id string) bool {
	return strings.HasPrefix(id, "dns-")
}

func dispatchDNSAction(id, title string) OSTUIActionResult {
	switch id {
	case "dns-cf":
		return actionResultFromError(title, handleDNSDirectSet("cloudflare"), "switched to Cloudflare DNS")
	case "dns-google":
		return actionResultFromError(title, handleDNSDirectSet("google"), "switched to Google DNS")
	case "dns-quad9":
		return actionResultFromError(title, handleDNSDirectSet("quad9"), "switched to Quad9 DNS")
	case "dns-dhcp":
		return actionResultFromError(title, handleDNSDHCP(), "reverted to automatic DHCP")
	default:
		return OSTUIActionResult{Title: title, IsSuccess: false, Message: "unknown DNS action"}
	}
}

func isCleanAction(id string) bool {
	return strings.HasPrefix(id, "clean-")
}

func dispatchCleanAction(id, title string) OSTUIActionResult {
	switch id {
	case "clean-dev":
		return actionResultFromError(title, RunOSDevClean(nil), "dev caches cleaned")
	case "clean-ai":
		return actionResultFromError(title, RunOSAICleanCLI([]string{"--force"}), "AI logs cleaned")
	case "clean-term":
		return actionResultFromError(title, RunTerminalCleanCLI(nil), "terminal logs cleaned")
	case "clean-sys":
		_, err := GetSystemCleanEngine().CleanSystemPackages()
		return actionResultFromError(title, err, "system packages cleaned")
	default:
		return OSTUIActionResult{Title: title, IsSuccess: false, Message: "unknown clean action"}
	}
}

func isUpdateAction(id string) bool {
	return strings.HasPrefix(id, "update-")
}

func dispatchUpdateAction(id, title string) OSTUIActionResult {
	switch id {
	case "update-repos":
		return actionResultFromError(title, runOSUpdateCommand(false, nil), "repositories refreshed")
	case "update-upgrade":
		return actionResultFromError(title, runOSUpdateCommand(true, nil), "system packages upgraded")
	case "update-mirrors":
		return actionResultFromError(title, FixRegionalMirrors(""), "mirrors updated to canonical US")
	default:
		return OSTUIActionResult{Title: title, IsSuccess: false, Message: "unknown update action"}
	}
}
