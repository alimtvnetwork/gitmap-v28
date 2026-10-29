package cmdui

import (
	"encoding/json"
	"net/http"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type historyPurgeRequest struct {
	Repo        string `json:"repo"`
	File        string `json:"file"`
	Folder      string `json:"folder"`
	Commit      string `json:"commit"`
	Range       string `json:"range"`
	IsDryRun    bool   `json:"isDryRun"`
	ConfirmText string `json:"confirmText"`
}

const historyConfirmWord = "yes"

func handleAPIHistoryPurge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	req, errRead := readHistoryPurge(r)
	if errRead != nil {
		writeHistoryResult(w, false, errRead.Error())

		return
	}
	if errGate := gateHistoryRewrite(req); errGate != nil {
		writeHistoryResult(w, false, errGate.Error())

		return
	}
	errRun := cmdhistory.RunHistoryPurgeCLI(historyPurgeArgs(req))
	if errRun != nil {
		writeHistoryResult(w, false, errRun.Error())

		return
	}
	writeHistoryResult(w, true, "history purge finished")
}

func readHistoryPurge(r *http.Request) (historyPurgeRequest, error) {
	var req historyPurgeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	return req, err
}

func gateHistoryRewrite(req historyPurgeRequest) error {
	if req.IsDryRun {
		return nil
	}
	if req.ConfirmText == historyConfirmWord {
		return nil
	}

	return errHistoryConfirm()
}

func errHistoryConfirm() error {
	return errString("type yes to rewrite history; dry-run does not need it")
}

type errString string

func (e errString) Error() string {
	return string(e)
}

func historyPurgeArgs(req historyPurgeRequest) []string {
	args := []string{"--repo", historyRepo(req.Repo)}
	args = append(args, historyTargetArgs(req)...)
	if req.IsDryRun {
		args = append(args, "--dry-run")
	}
	if req.ConfirmText == historyConfirmWord {
		args = append(args, "--yes")
	}

	return args
}

func historyRepo(repo string) string {
	if repo == "" {
		return "."
	}

	return repo
}

func historyTargetArgs(req historyPurgeRequest) []string {
	if req.File != "" {
		return []string{"--file", req.File}
	}
	if req.Folder != "" {
		return []string{"--folder", req.Folder}
	}
	if req.Commit != "" {
		return []string{"--commit", req.Commit}
	}
	if req.Range != "" {
		return []string{"--range", req.Range}
	}

	return nil
}

func writeHistoryResult(w http.ResponseWriter, isSuccess bool, message string) {
	_ = json.NewEncoder(w).Encode(map[string]any{
		"isSuccess": isSuccess,
		"message":   message,
	})
}
