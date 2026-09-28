// Package cmdagy — agy_account_switch_notify.go dispatches Telegram and Email notifications for account switching.
package cmdagy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AccountSwitchNotice holds structured step and completion telemetry for Telegram and Email notifications.
type AccountSwitchNotice struct {
	Stage              string   `json:"stage"`
	FromAccount        string   `json:"fromAccount"`
	TargetAccount      string   `json:"targetAccount"`
	ThresholdPct       int      `json:"thresholdPct"`
	CandidateCreditPct int      `json:"candidateCreditPct"`
	IsRefreshed        bool     `json:"isRefreshed"`
	IsLockVerified     bool     `json:"isLockVerified"`
	ProjectCount       int      `json:"projectCount"`
	ProjectNames       []string `json:"projectNames"`
	TotalPrompts       int      `json:"totalPrompts"`
	RestoredPrompts    int      `json:"restoredPrompts"`
	IsLivenessVerified bool     `json:"isLivenessVerified"`
	Timestamp          string   `json:"timestamp"`
}

func sendAccountSwitchNotifications(notice AccountSwitchNotice) (bool, bool) {
	notice.ProjectNames = sanitizeNotificationProjectNames(notice.ProjectNames)
	notice.ProjectCount = len(notice.ProjectNames)
	if len(notice.Timestamp) == 0 {
		notice.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	msg := formatAccountSwitchMessage(notice)
	subject := fmt.Sprintf("[GitMap AGY] Account Switch (%s) — %d Project(s)", notice.Stage, notice.ProjectCount)
	isTelegramSent := sendTelegramSwitchNotice(msg)
	isEmailSent := sendEmailSwitchNotice(subject, msg)
	recordAccountSwitchNoticeLog(notice, msg)
	return isTelegramSent, isEmailSent
}

func sanitizeNotificationProjectNames(names []string) []string {
	seen := make(map[string]bool)
	var clean []string
	for _, raw := range names {
		trimmed := strings.TrimSpace(raw)
		if isHumanReadableProjectName(trimmed, seen) {
			seen[trimmed] = true
			clean = append(clean, trimmed)
		}
	}
	return clean
}

func isHumanReadableProjectName(name string, seen map[string]bool) bool {
	if len(name) == 0 || seen[name] {
		return false
	}
	return !reRawUUID.MatchString(name)
}

func formatAccountSwitchMessage(notice AccountSwitchNotice) string {
	projList := formatNoticeProjectsList(notice.ProjectNames)
	return fmt.Sprintf(
		"Antigravity Account Switch [%s]\nThreshold: %d%% | Target: %s (%d%% remaining)\nRefresh Verified: %v | VM/Email Lock Clear: %v\nProjects Backed Up (%d): %s\nPrompts Backed Up: %d | Restored: %d | Running Liveness: %v\nTime: %s",
		notice.Stage, notice.ThresholdPct, notice.TargetAccount, notice.CandidateCreditPct,
		notice.IsRefreshed, notice.IsLockVerified,
		notice.ProjectCount, projList,
		notice.TotalPrompts, notice.RestoredPrompts, notice.IsLivenessVerified,
		notice.Timestamp,
	)
}

func formatNoticeProjectsList(names []string) string {
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

func sendTelegramSwitchNotice(text string) bool {
	botToken := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	chatID := strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID"))
	if len(botToken) == 0 || len(chatID) == 0 {
		return true
	}
	return postTelegramMessage(botToken, chatID, text)
}

func postTelegramMessage(botToken, chatID, text string) bool {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload, _ := json.Marshal(map[string]string{"chat_id": chatID, "text": text})
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func sendEmailSwitchNotice(subject, body string) bool {
	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	toAddr := strings.TrimSpace(os.Getenv("AGM_NOTIFY_EMAIL"))
	if len(smtpHost) == 0 || len(toAddr) == 0 {
		return true
	}
	return dispatchSmtpEmail(smtpHost, toAddr, subject, body)
}

func dispatchSmtpEmail(smtpHost, toAddr, subject, body string) bool {
	fromAddr := resolveSmtpFromAddress(toAddr)
	port := resolveSmtpPort()
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", toAddr, subject, body))
	err := smtp.SendMail(smtpHost+":"+port, nil, fromAddr, []string{toAddr}, msg)
	return err == nil
}

func resolveSmtpFromAddress(fallback string) string {
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if len(from) > 0 {
		return from
	}
	return fallback
}

func resolveSmtpPort() string {
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if len(port) > 0 {
		return port
	}
	return "25"
}

func recordAccountSwitchNoticeLog(notice AccountSwitchNotice, rendered string) {
	logDir := filepath.Join(store.BinaryDataDir(), "account-switch")
	_ = os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, "notifications.log")
	entry, _ := json.Marshal(notice)
	line := fmt.Sprintf("%s\n%s\n---\n", string(entry), rendered)
	appendTextToFile(logPath, line)
}

func appendTextToFile(path, content string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(content)
}
