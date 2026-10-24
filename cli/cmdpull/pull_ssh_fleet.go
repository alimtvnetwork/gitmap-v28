package cmdpull

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func handleSSHFleetPullAll(args []string) error {
	if hasFleetRemediateFlag(args) && RunRemoteSSHPullFleetRemediateFn != nil {
		return RunRemoteSSHPullFleetRemediateFn(args)
	}
	cleanArgs := stripSSHFleetFlags(args)
	hasHook := RunRemoteSSHPullAllFleetFn != nil
	if !hasHook {
		return apperror.NewSimple("ssh fleet pull-all is not wired", "E_SSH_PULLALL")
	}

	return executeSSHFleetPullTask(cleanArgs)
}

func hasFleetRemediateFlag(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "--fix" || low == "-fix" || low == "--remediate" || low == "-remediate" || low == "--fleet-fix" {
			return true
		}
	}

	return false
}

func executeSSHFleetPullTask(cleanArgs []string) error {
	queueId, tDB := enqueueTaskQueue("pull-all-ssh", "fleet")
	defer closeTaskDB(tDB)
	updateTaskQueue(tDB, queueId, "running")

	errFleet := RunRemoteSSHPullAllFleetFn(cleanArgs)
	hasErr := errFleet != nil
	if hasErr {
		return handleSSHFleetFailure(tDB, queueId, errFleet)
	}

	updateTaskQueue(tDB, queueId, "completed")

	return nil
}

func handleSSHFleetFailure(tDB *store.TasksSplitDB, queueId string, errFleet error) error {
	store.LogInternalError("PULL_ALL_SSH", "SSH_DELEGATION_ERROR", errFleet.Error(), "", "")
	updateTaskQueue(tDB, queueId, "failed")

	return apperror.WrapSimple(errFleet, "SSH fleet pull-all failed")
}

func hasSSHFleetFlag(args []string) bool {
	for _, a := range args {
		if isSSHFleetToken(a) {
			return true
		}
	}

	return false
}

func isSSHFleetToken(token string) bool {
	low := strings.ToLower(token)

	return low == "--ssh" || low == "-ssh" || low == "--sh" || low == "-sh" || low == "ssh" || low == "pas" || low == "pull-all-ssh"
}

func stripSSHFleetFlags(args []string) []string {
	var clean []string
	for _, a := range args {
		if isSSHFleetToken(a) {
			continue
		}
		clean = append(clean, a)
	}

	return clean
}

func enqueueTaskQueue(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())
	tasksDB, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDB.Conn(), "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'pull', ?, ?, 'pending', ?, ?)", queueId, action, target, now, now)

	return queueId, tasksDB
}

func updateTaskQueue(db *store.TasksSplitDB, queueId, status string) {
	hasDB := db != nil
	if !hasDB {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(db.Conn(), "UPDATE TaskQueue SET Status = ?, UpdatedAt = ? WHERE QueueId = ?", status, now, queueId)
}

func closeTaskDB(db *store.TasksSplitDB) {
	if db != nil {
		_ = db.Close()
	}
}
