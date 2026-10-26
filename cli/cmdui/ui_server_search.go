package cmdui

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautomation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func handleAPISearchRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)

		return
	}

	opts := parseSearchRepoQueryParams(r)
	res, appErr := cmdautomation.RunSearch(opts)
	if appErr != nil {
		encodeSearchRepoError(w, appErr.Error())

		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func parseSearchRepoQueryParams(r *http.Request) cmdautomation.SearchOptions {
	q := r.URL.Query()
	rawExt := strings.TrimSpace(q.Get("ext"))
	var exts []string
	if len(rawExt) > 0 {
		exts = splitSearchExts(rawExt)
	}

	return cmdautomation.SearchOptions{
		Pattern:           strings.TrimSpace(q.Get("q")),
		Dir:               strings.TrimSpace(q.Get("dir")),
		IsRegex:           q.Get("isRegex") == "true",
		IsCaseInsensitive: q.Get("isCaseInsensitive") == "true",
		Extensions:        exts,
	}
}

func splitSearchExts(raw string) []string {
	parts := strings.Split(raw, ",")
	clean := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if len(trimmed) > 0 {
			clean = append(clean, trimmed)
		}
	}

	return clean
}

func encodeSearchRepoError(w http.ResponseWriter, errMsg string) {
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"matches":    []any{},
		"totalFiles": 0,
		"totalHits":  0,
		"error":      errMsg,
	})
}

func handleAPIOpenAntigravity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	targetPath := strings.TrimSpace(r.URL.Query().Get("path"))
	if len(targetPath) == 0 {
		http.Error(w, "path required", http.StatusBadRequest)

		return
	}

	err := cmdagy.RunAgyOpen(targetPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"isSuccess": false, "error": err.Error()})

		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"isSuccess": true, "path": targetPath})
}
