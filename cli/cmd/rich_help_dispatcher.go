package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
)

// tryRenderRichTopic renders a rich two-column help menu if the topic has one.
// Returns true if rendered, false otherwise.
func tryRenderRichTopic(topic string) bool {
	if tryRenderInfraRichTopic(topic) {
		return true
	}
	if tryRenderToolRichTopic(topic) {
		return true
	}

	return tryRenderWorkflowRichTopic(topic)
}

func tryRenderInfraRichTopic(topic string) bool {
	switch topic {
	case "ssh", "sj":
		cmdssh.RenderSSHHelp()

		return true
	case "cluster":
		RenderClusterHelp()

		return true
	case "sc", "clients", "servers":
		RenderSCHelp()

		return true
	case "vhost":
		cmdvhost.RenderVHostHelp()

		return true
	}

	return false
}

func tryRenderToolRichTopic(topic string) bool {
	switch topic {
	case "install", "in":
		cmdinstall.RenderInstallHelp()

		return true
	case "storage":
		RenderStorageHelp()

		return true
	case "sync":
		RenderSyncHelp()

		return true
	case "vscode", "code":
		cmdvscode.RenderVSCodeHelp()

		return true
	case "github-desktop", "gd", "github", "desktop-sync", "ds":
		RenderGitHubDesktopHelp()

		return true
	}

	return false
}

func tryRenderWorkflowRichTopic(topic string) bool {
	switch topic {
	case "macro":
		cmdmacro.RenderMacroHelp()

		return true
	case "schedule", "sched":
		cmdschedule.RenderScheduleHelp()

		return true
	case "pipeline":
		cmdpipeline.RenderPipelineHelp()

		return true
	case "agy", "antigravity":
		cmdagy.RenderAgyHelp()

		return true
	case "clone", "c":
		cmdclone.RenderCloneHelp()

		return true
	case "multiclone", "mutliclone", "mc":
		cmdclone.RenderMultiCloneHelp()

		return true
	case "scan", "s":
		cmdscan.RenderScanHelp()

		return true
	case "fix-git", "fg":
		cmdfixgit.RenderFixGitHelp()

		return true
	case "ai", "scripts":
		cmdai.RenderAiHelp()

		return true
	}

	return false
}
