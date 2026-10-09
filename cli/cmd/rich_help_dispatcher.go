package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstorage"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcluster"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmduser"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdllm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautofix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemplates"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfind"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
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
	if tryRenderCoreRichTopic(topic) {
		return true
	}

	return tryRenderWorkflowRichTopic(topic)
}

func tryRenderCoreRichTopic(topic string) bool {
	switch topic {
	case "search", "grep":
		cmdsearch.RenderSearchHelp()

		return true
	case "find", "find-files", "find-files-any", "find-files-startswith", "find-files-endswith", "ff", "ffa", "ffs", "ffe", "f":
		cmdfind.RenderFindHelp()

		return true
	case "commit", "cm", "commit-all", "ca", "commit-both", "commit-left", "commit-right", "cmb", "cml", "cmr", "cpf", "cpb", "cpr", "pcp", "commit-pull", "cpull":
		cmdcommit.RenderCommitHelp()

		return true
	case "status", "st":
		cmdstatus.RenderStatusHelp()

		return true
	case "diff", "df", "diff-profiles":
		cmddiff.RenderDiffHelp()

		return true
	case "templates", "template", "tpl":
		cmdtemplates.RenderTemplatesHelp()

		return true
	case "ui", "gui", "web":
		cmdui.RenderUIHelp()

		return true
	case "aum", "automation", "auto":
		RenderAumHelp()

		return true
	case "llm", "llm-train", "train", "chain", "llm-docs", "ld":
		cmdllm.RenderLlmHelp()

		return true
	case "branch", "b", "latest-branch", "lb":
		cmdbranch.RenderBranchHelp()

		return true
	case "stash", "wip", "discard":
		RenderStashHelp()

		return true
	case "fix":
		cmdautofix.RenderFixHelp()

		return true
	case "user", "os-user":
		cmduser.RenderUserHelp()

		return true
	}

	return false
}

func tryRenderInfraRichTopic(topic string) bool {
	switch topic {
	case "ssh", "sj":
		cmdssh.RenderSSHHelp()

		return true
	case "cluster":
		cmdcluster.RenderClusterHelp()

		return true
	case "sc", "clients", "servers":
		cmdschedule.RenderSCHelp()

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
	case "update", "ua", "update-all", "updateall":
		cmdupdate.RenderUpdateRichHelp()

		return true
	case "deploy", "deploy-right", "deploy-left", "deploy-config", "deploy-config-ssh", "deploy-ssh-config":
		cmdssh.RenderDeployRichHelp()

		return true
	case "storage":
		cmdstorage.RenderStorageHelp()

		return true
	case "sync":
		cmdsync.RenderSyncHelp()

		return true
	case "vscode", "code":
		cmdvscode.RenderVSCodeHelp()

		return true
	case "cursor", "cur":
		cmdcursor.RenderCursorHelp()

		return true
	case "github-desktop", "gd", "github", "desktop-sync", "ds":
		cmdvscode.RenderGitHubDesktopHelp()

		return true
	}

	return false
}

func tryRenderWorkflowRichTopic(topic string) bool {
	switch topic {
	case "shutdown-until", "shutdown-until-green", "sug":
		cmdagy.RenderAgySugHelp()

		return true
	case "watch-prompts-running", "watch-running-prompts", "wpr":
		cmdagy.RenderWPRHelp()

		return true
	case "running-prompts", "running-prompt", "backup-running-prompts", "restore-running-prompts", "rp-prompts":
		cmdagy.RenderRunningPromptsHelp()

		return true
	case "macro":
		cmdmacro.RenderMacroHelp()

		return true
	case "schedule", "sched":
		cmdschedule.RenderScheduleHelp()

		return true
	case "pipeline":
		cmdpipeline.RenderPipelineHelp()

		return true
	case "rerun", "rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue":
		cmdagy.RenderAgyRerunHelp()

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
