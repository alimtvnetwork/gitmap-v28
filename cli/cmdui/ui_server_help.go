package cmdui

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type topicSummaryItem struct {
	Name        string `json:"name"`
	Topic       string `json:"topic"`
	Summary     string `json:"summary"`
	HasMarkdown bool   `json:"hasMarkdown"`
}

const (
	helpTopicMissMessagePrefix = "embedded manual file is absent for topic "
	helpTopicMissErrorCode     = "ERR_HELP_TOPIC_MISS"
	helpTopicMissErrorType     = "HelpDocMiss"
	helpTopicMissCommand       = "ui help doc"
	helpTopicMissSourceFile    = "cli/cmdui/ui_server_help.go"
)

func handleAPIHelpTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{"isSuccess": false, "error": "method not allowed, GET required"})
		return
	}

	topics := helpdoc.AllHelpTopics()
	items := make([]topicSummaryItem, 0, len(topics))
	for _, t := range topics {
		items = append(items, topicSummaryItem{
			Name:        t,
			Topic:       t,
			Summary:     helpdoc.GetTopicDetailedSummary(t),
			HasMarkdown: helpdoc.HasTopic(t),
		})
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"isSuccess":   true,
		"totalTopics": len(items),
		"topics":      items,
	})
}

func handleAPIHelpDoc(w http.ResponseWriter, r *http.Request) {
	topic := strings.TrimSpace(r.URL.Query().Get("topic"))
	if topic == "" {
		topic = strings.TrimSpace(r.URL.Query().Get("command"))
	}
	if topic == "" {
		http.Error(w, "topic query param is required", http.StatusBadRequest)
		return
	}

	content, err := helpdoc.ReadRaw(topic)
	if err != nil {
		logHelpTopicMiss(topic)
		http.Error(w, helpTopicMissMessagePrefix + topic, http.StatusNotFound)

		return
	}

	format := r.URL.Query().Get("format")
	if format == "json" {
		w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON+"; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": true,
			"topic":     topic,
			"content":   string(content),
		})
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func logHelpTopicMiss(topic string) {
	store.LogInternalErrorRecord(store.InternalErrorRecord{
		ErrorCode:  helpTopicMissErrorCode,
		ErrorType:  helpTopicMissErrorType,
		Command:    helpTopicMissCommand,
		Message:    helpTopicMissMessagePrefix + topic,
		Details:    topic,
		SourceFile: helpTopicMissSourceFile,
	})
}
