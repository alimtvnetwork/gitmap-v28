package cluster

import (
	"context"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
)

// Dispatch routes the subcommand to the correct executor for the node.
// This is currently a stub for executors.
func Dispatch(ctx context.Context, node ClusterNode, subCmd ClusterSubCommand) db.ClusterExecResult {
	start := time.Now()
	res := db.ClusterExecResult{
		NodeId:       node.ID,
		SubCommand:   subCmd.Kind.String(),
		ResultStatus: db.ResultStatusPending,
		StartedAt:    &start,
	}
	
	raw := subCmd.RawArg
	res.CommandText = &raw

	// Stub routing logic
	switch subCmd.Kind {
	case db.CommandKindPsCommand:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecPS stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindCmdCommand:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecCmd stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindInstall:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecInstall stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindGitPull, db.CommandKindGitPush, db.CommandKindGitCommit, db.CommandKindGitStatus:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecGit stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindProjRun, db.CommandKindProjCreateCICD:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecProj stubbed"
		res.ErrorMessage = &msg
	case db.CommandKindRestart, db.CommandKindShutdown, db.CommandKindLogoff:
		res.ResultStatus = db.ResultStatusSkipped
		msg := "ExecLifecycle stubbed"
		res.ErrorMessage = &msg
	default:
		res.ResultStatus = db.ResultStatusFailed
		msg := "unknown command kind"
		res.ErrorMessage = &msg
	}

	end := time.Now()
	res.FinishedAt = &end
	durMs := int(end.Sub(start).Milliseconds())
	res.DurationMs = &durMs

	return res
}
