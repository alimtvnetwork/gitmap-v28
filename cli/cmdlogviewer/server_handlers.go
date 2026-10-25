package cmdlogviewer

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func handleAPICommits(repoRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		limit := parseQueryLimit(r.URL.Query().Get("limit"))
		commits, err := QueryCommits(repoRoot, path, limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, commits)
	}
}

func parseQueryLimit(raw string) int {
	if raw == "" {
		return 100
	}

	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return 100
	}

	return val
}

func handleAPICommitDetail(repoRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Query().Get("hash")
		if hash == "" {
			writeJSONError(w, http.StatusBadRequest, "missing commit hash parameter")
			return
		}

		detail, err := QueryCommitDetail(repoRoot, hash)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, detail)
	}
}

func handleAPIDiff(repoRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DiffRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid diff request body")
			return
		}

		if req.Commit == "" || req.Path == "" {
			writeJSONError(w, http.StatusBadRequest, "commit and path are required")
			return
		}

		executeDiffRequest(w, repoRoot, req)
	}
}

func executeDiffRequest(w http.ResponseWriter, repoRoot string, req DiffRequest) {
	tool, err := TriggerExternalDiff(repoRoot, req.Commit, req.Path)
	if err != nil {
		resp := DiffResponse{IsSuccess: false, Error: err.Error()}
		writeJSON(w, http.StatusInternalServerError, resp)
		return
	}

	resp := DiffResponse{
		IsSuccess: true,
		Tool:      tool,
		Message:   "Compare tool launched: " + tool,
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleAPIBranches(repoRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		branches, err := QueryBranches(repoRoot)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, branches)
	}
}

func handleAPIBranchCompare(repoRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BranchCompareRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid branch compare body")
			return
		}

		if req.BranchA == "" || req.BranchB == "" {
			writeJSONError(w, http.StatusBadRequest, "branchA and branchB are required")
			return
		}

		executeBranchCompareRequest(w, repoRoot, req)
	}
}

func executeBranchCompareRequest(w http.ResponseWriter, repoRoot string, req BranchCompareRequest) {
	tool, err := TriggerBranchDiff(repoRoot, req.BranchA, req.BranchB, req.Path)
	if err != nil {
		resp := BranchCompareResponse{IsSuccess: false, Error: err.Error()}
		writeJSON(w, http.StatusInternalServerError, resp)
		return
	}

	resp := BranchCompareResponse{
		IsSuccess: true,
		Message:   "Branch comparison launched with " + tool,
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
