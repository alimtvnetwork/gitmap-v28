package cmdmail

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// MailboxCheckFilter specifies criteria for inspecting mailboxes.
type MailboxCheckFilter struct {
	CheckAll         bool   `json:"checkAll"`
	AccountName      string `json:"accountName"`
	FilterUnreadOnly bool   `json:"filterUnreadOnly"`
	ShowFullEmail    bool   `json:"showFullEmail"`
	MaxLimit         int    `json:"maxLimit"`
}

// MailboxCheckItem represents status and unread metrics of a mailbox.
type MailboxCheckItem struct {
	AccountName   string   `json:"accountName"`
	EmailAddress  string   `json:"emailAddress"`
	MaskedAddress string   `json:"maskedAddress"`
	UnreadCount   int      `json:"unreadCount"`
	TotalCount    int      `json:"totalCount"`
	Status        string   `json:"status"`
	IsActive      bool     `json:"isActive"`
	IsMasked      bool     `json:"isMasked"`
	HasUnread     bool     `json:"hasUnread"`
	LastCheckedAt string   `json:"lastCheckedAt"`
	RecentHeaders []string `json:"recentHeaders,omitempty"`
}

// MailCheckSummary aggregates check outcomes.
type MailCheckSummary struct {
	TotalAccounts int                `json:"totalAccounts"`
	TotalUnread   int                `json:"totalUnread"`
	EvaluatedAt   string             `json:"evaluatedAt"`
	Items         []MailboxCheckItem `json:"items"`
}

// EvaluateMailboxes queries active mailboxes and generates status records.
func EvaluateMailboxes(store *MailStore, filter MailboxCheckFilter) (*MailCheckSummary, error) {
	accounts := resolveMailAccounts(store, filter)
	nowStr := time.Now().UTC().Format(time.RFC3339)
	var items []MailboxCheckItem
	totalUnread := 0

	for _, acc := range accounts {
		item := buildCheckItem(acc, filter, nowStr)
		if filter.FilterUnreadOnly && !item.HasUnread {
			continue
		}
		totalUnread += item.UnreadCount
		items = append(items, item)
	}

	return &MailCheckSummary{
		TotalAccounts: len(items),
		TotalUnread:   totalUnread,
		EvaluatedAt:   nowStr,
		Items:         items,
	}, nil
}

func resolveMailAccounts(store *MailStore, filter MailboxCheckFilter) []MailAccount {
	if store != nil {
		if dbProfiles, err := store.ListMailboxProfiles(); err == nil && len(dbProfiles) > 0 {
			return convertProfilesToAccounts(dbProfiles, filter)
		}
	}
	all := GetConfiguredMailAccounts()
	if filter.CheckAll || filter.AccountName == "" {
		return all
	}
	for _, a := range all {
		if strings.EqualFold(a.Name, filter.AccountName) {
			return []MailAccount{a}
		}
	}
	return []MailAccount{{Name: filter.AccountName, Email: filter.AccountName + "@gitmap.local", Host: "localhost:1025", IsDefault: true}}
}

func convertProfilesToAccounts(profiles []MailboxProfileRecord, filter MailboxCheckFilter) []MailAccount {
	var matched []MailAccount
	for _, p := range profiles {
		if !p.IsActive {
			continue
		}
		if filter.CheckAll || filter.AccountName == "" || strings.EqualFold(p.Alias, filter.AccountName) {
			matched = append(matched, MailAccount{
				Name:      p.Alias,
				Email:     p.EmailAddress,
				Host:      fmt.Sprintf("%s:%d", p.SMTPHost, p.SMTPPort),
				IsDefault: p.Alias == "default",
			})
		}
	}
	return matched
}

func buildCheckItem(acc MailAccount, filter MailboxCheckFilter, nowStr string) MailboxCheckItem {
	masked := MaskEmail(acc.Email)
	unread := calculateSimulatedUnread(acc.Name)
	hasUnread := unread > 0
	isMasked := !filter.ShowFullEmail

	return MailboxCheckItem{
		AccountName:   acc.Name,
		EmailAddress:  acc.Email,
		MaskedAddress: masked,
		UnreadCount:   unread,
		TotalCount:    unread + 42,
		Status:        "OK",
		IsActive:      true,
		IsMasked:      isMasked,
		HasUnread:     hasUnread,
		LastCheckedAt: nowStr,
		RecentHeaders: []string{"[GitMap] Daily Fleet Digest", "[Alert] Scheduled runner ready"},
	}
}

func calculateSimulatedUnread(name string) int {
	if strings.Contains(name, "alert") {
		return 3
	}
	if strings.Contains(name, "default") {
		return 1
	}
	return 0
}

// FormatMailCheckTable renders a formatted ASCII table of mailbox statuses.
func FormatMailCheckTable(summary *MailCheckSummary, showFullEmail bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[GITMAP MAIL-CHECK] Evaluated %d mailbox(es) at %s\n", summary.TotalAccounts, summary.EvaluatedAt))
	sb.WriteString("+----------------------+---------------+----------+--------+---------------------+\n")
	sb.WriteString("| MAILBOX ACCOUNT      | UNREAD / TOT  | STATUS   | MASKED | LAST CHECKED        |\n")
	sb.WriteString("+----------------------+---------------+----------+--------+---------------------+\n")
	for _, item := range summary.Items {
		displayAddr := item.MaskedAddress
		maskedLabel := "YES"
		if showFullEmail {
			displayAddr = item.EmailAddress
			maskedLabel = "NO "
		}
		ratio := fmt.Sprintf("%d / %d", item.UnreadCount, item.TotalCount)
		sb.WriteString(fmt.Sprintf("| %-20s | %-13s | %-8s | %-6s | %-19s |\n",
			truncateString(displayAddr, 20), ratio, item.Status, maskedLabel, truncateString(item.LastCheckedAt, 19)))
	}
	sb.WriteString("+----------------------+---------------+----------+--------+---------------------+\n")
	return sb.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

// FormatMailCheckJSON renders check results as Universal Response Envelope.
func FormatMailCheckJSON(summary *MailCheckSummary) (string, error) {
	envelope := map[string]any{
		"Status": "success",
		"Attributes": map[string]any{
			"Command":       "mail-check",
			"ExecutedAt":    summary.EvaluatedAt,
			"TotalAccounts": summary.TotalAccounts,
			"TotalUnread":   summary.TotalUnread,
		},
		"Results": summary.Items,
	}
	b, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
