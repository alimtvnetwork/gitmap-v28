package cmdui

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func handleAPIInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": false,
			"error":     "method not allowed, GET required",
		})
		return
	}

	instances, err := cmdagy.DiscoverAllAgyInstances()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess":      false,
			"totalInstances": 0,
			"instances":      []cmdagy.AgyInstanceInfo{},
			"error":          err.Error(),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"isSuccess":      true,
		"totalInstances": len(instances),
		"instances":      instances,
	})
}

func handleAPIPromptsInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": false,
			"error":     "method not allowed, GET required",
		})
		return
	}

	q := r.URL.Query()
	instID := q.Get("instance")
	status := q.Get("status")
	limitStr := q.Get("limit")
	maxWordsStr := q.Get("max_words")
	includeConvsStr := q.Get("include_convs")

	limit := parsePositiveInt(limitStr, 10)
	maxWords := parsePositiveInt(maxWordsStr, 100)

	includeConvs := true
	if includeConvsStr != "" {
		includeConvs = includeConvsStr == "true" || includeConvsStr == "1"
	}

	opts := cmdagy.AgyInstancePromptQueryOptions{
		InstanceID:   instID,
		IsAll:        instID == "" || strings.EqualFold(instID, "all"),
		Status:       status,
		Limit:        limit,
		MaxWords:     maxWords,
		IncludeConvs: includeConvs,
	}

	resp, err := cmdagy.QueryInstancePrompts(opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": false,
			"error":     err.Error(),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func parsePositiveInt(str string, fallback int) int {
	if str == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(str)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
