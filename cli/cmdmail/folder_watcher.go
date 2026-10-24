package cmdmail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdmail/enums"

const (
	MinPollingInterval     = 5 * time.Second
	DefaultPollingInterval = 8 * time.Second
	DefaultDebounceWindow  = 3 * time.Second
)

// IntervalGuard validates and clamps polling frequency to safe minimums.
type IntervalGuard struct {
	ConfiguredInterval time.Duration
	EffectiveInterval  time.Duration
	WasClamped         bool
}

// NewIntervalGuard constructs guard enforcing MinPollingInterval >= 5s.
func NewIntervalGuard(requested time.Duration) IntervalGuard {
	if requested <= 0 {
		return IntervalGuard{ConfiguredInterval: requested, EffectiveInterval: DefaultPollingInterval, WasClamped: false}
	}
	if requested < MinPollingInterval {
		return IntervalGuard{ConfiguredInterval: requested, EffectiveInterval: MinPollingInterval, WasClamped: true}
	}
	return IntervalGuard{ConfiguredInterval: requested, EffectiveInterval: requested, WasClamped: false}
}

// FileEventType identifies file mutation categories.
// FileEventType aliases the strongly-typed enum from the enums package.
type FileEventType = enums.FileEventType

const (
	EventCreate = enums.EventCreate
	EventModify = enums.EventModify
	EventDelete = enums.EventDelete
)

// FileItemSnapshot captures filesystem metadata for change detection.
type FileItemSnapshot struct {
	RelativePath string
	SizeBytes    int64
	ModTimeUnix  int64
	IsDir        bool
}

// FileEventRecord represents a detected file mutation.
type FileEventRecord struct {
	EventType FileEventType `json:"eventType"`
	Path      string        `json:"path"`
	SizeBytes int64         `json:"sizeBytes"`
	Timestamp time.Time     `json:"timestamp"`
	Details   string        `json:"details"`
}

// WatcherBatch bundles grouped file events dispatched to notifiers.
type WatcherBatch struct {
	RuleID       string            `json:"ruleId"`
	TargetFolder string            `json:"targetFolder"`
	TriggeredAt  time.Time         `json:"triggeredAt"`
	TotalEvents  int               `json:"totalEvents"`
	Events       []FileEventRecord `json:"events"`
	Recipients   []string          `json:"recipients"`
	HasEvents    bool              `json:"hasEvents"`
}

// FolderWatcherConfig specifies monitoring rules.
type FolderWatcherConfig struct {
	RuleID          string        `json:"ruleId"`
	TargetDirectory string        `json:"targetDirectory"`
	PollingInterval time.Duration `json:"pollingInterval"`
	DebounceWindow  time.Duration `json:"debounceWindow"`
	Recipients      []string      `json:"recipients"`
	SubjectTemplate string        `json:"subjectTemplate"`
	IsRecursive     bool          `json:"isRecursive"`
	IsActive        bool          `json:"isActive"`
	RunOnce         bool          `json:"runOnce"`
}

// ScanFolderTree walks folder and returns map of relative path to metadata.
func ScanFolderTree(root string, isRecursive bool) (map[string]FileItemSnapshot, error) {
	result := make(map[string]FileItemSnapshot)
	cleanRoot := filepath.Clean(root)
	err := filepath.Walk(cleanRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil || p == cleanRoot {
			return nil
		}
		rel, _ := filepath.Rel(cleanRoot, p)
		rel = filepath.ToSlash(rel)
		if isSkippedPath(rel) && info.IsDir() {
			return filepath.SkipDir
		}
		if isSkippedPath(rel) {
			return nil
		}
		if !isRecursive && info.IsDir() {
			return filepath.SkipDir
		}
		result[rel] = FileItemSnapshot{
			RelativePath: rel,
			SizeBytes:    info.Size(),
			ModTimeUnix:  info.ModTime().UTC().Unix(),
			IsDir:        info.IsDir(),
		}
		return nil
	})
	return result, err
}

func isSkippedPath(rel string) bool {
	return strings.HasPrefix(rel, ".git") || strings.HasPrefix(rel, "node_modules")
}

// DetectFileMutations compares current state against baseline snapshot.
func DetectFileMutations(prev, curr map[string]FileItemSnapshot) []FileEventRecord {
	var events []FileEventRecord
	now := time.Now().UTC()

	for path, item := range curr {
		old, exists := prev[path]
		if !exists {
			events = append(events, FileEventRecord{EventType: EventCreate, Path: path, SizeBytes: item.SizeBytes, Timestamp: now, Details: "File created"})
		} else if old.SizeBytes != item.SizeBytes || old.ModTimeUnix != item.ModTimeUnix {
			details := fmt.Sprintf("Modified (size: %d -> %d)", old.SizeBytes, item.SizeBytes)
			events = append(events, FileEventRecord{EventType: EventModify, Path: path, SizeBytes: item.SizeBytes, Timestamp: now, Details: details})
		}
	}
	for path, old := range prev {
		if _, exists := curr[path]; !exists {
			events = append(events, FileEventRecord{EventType: EventDelete, Path: path, SizeBytes: old.SizeBytes, Timestamp: now, Details: "File deleted"})
		}
	}
	return events
}

// ExecuteWatchCycle performs single pass comparison, notification, and logging.
func ExecuteWatchCycle(cfg FolderWatcherConfig, baseline map[string]FileItemSnapshot, store *MailStore, transmitter EmailTransmitter) (*WatcherBatch, map[string]FileItemSnapshot, error) {
	curr, err := ScanFolderTree(cfg.TargetDirectory, cfg.IsRecursive)
	if err != nil {
		return nil, baseline, err
	}
	events := DetectFileMutations(baseline, curr)
	hasEvents := len(events) > 0
	batch := &WatcherBatch{
		RuleID:       cfg.RuleID,
		TargetFolder: cfg.TargetDirectory,
		TriggeredAt:  time.Now().UTC(),
		TotalEvents:  len(events),
		Events:       events,
		Recipients:   cfg.Recipients,
		HasEvents:    hasEvents,
	}
	if hasEvents {
		_ = dispatchBatchNotification(cfg, batch, store, transmitter)
	}
	return batch, curr, nil
}

func dispatchBatchNotification(cfg FolderWatcherConfig, batch *WatcherBatch, store *MailStore, transmitter EmailTransmitter) error {
	subject := cfg.SubjectTemplate
	if subject == "" {
		subject = fmt.Sprintf("[GitMap Notifier] %d file change(s) detected in %s", batch.TotalEvents, cfg.TargetDirectory)
	}
	body := buildDigestBody(cfg.TargetDirectory, batch)
	sendOpts := MailSendOptions{
		AccountName: "default",
		From:        "notifier@gitmap.local",
		To:          cfg.Recipients,
		Subject:     subject,
		Body:        body,
		Host:        "localhost",
		Port:        1025,
	}
	rawMsg := BuildMIMEMessage(sendOpts)
	txErr := transmitter.SendMail(sendOpts, rawMsg)
	isDelivered := txErr == nil

	if store != nil {
		_ = recordWatcherLog(store, cfg.RuleID, subject, batch, isDelivered, txErr)
	}
	return txErr
}

func recordWatcherLog(store *MailStore, ruleID, subject string, batch *WatcherBatch, isDelivered bool, txErr error) error {
	errMsg := ""
	if txErr != nil {
		errMsg = txErr.Error()
	}
	return store.RecordNotificationLog(NotificationLogRecord{
		LogID:             fmt.Sprintf("log-%d", time.Now().UnixNano()),
		RuleID:            ruleID,
		TriggerSource:     "FOLDER_WATCHER",
		RecipientSummary:  MaskEmailCSV(strings.Join(batch.Recipients, ", ")),
		Subject:           subject,
		MutatedFilesCount: batch.TotalEvents,
		IsDelivered:       isDelivered,
		ErrorMessage:      errMsg,
	})
}

func buildDigestBody(target string, batch *WatcherBatch) string {
	var sb strings.Builder
	sb.WriteString("GitMap File Watcher Activity Report\n")
	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Target Directory: %s\n", target))
	sb.WriteString(fmt.Sprintf("Triggered At:     %s\n", batch.TriggeredAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Total Events:     %d\n\n", batch.TotalEvents))
	sb.WriteString("Changes:\n")
	for _, ev := range batch.Events {
		sb.WriteString(fmt.Sprintf("  [%s] %s (%s)\n", strings.ToUpper(string(ev.EventType)), ev.Path, ev.Details))
	}
	sb.WriteString("\n--\nAutomated notification generated by GitMap Scheduler.\n")
	return sb.String()
}
