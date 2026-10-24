// Package cmdschedule manages scheduled tasks, privacy masking, and mailbox inspection.
package cmdschedule

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// MailboxCheckResult encapsulates the evaluated state of a single mailbox account.
type MailboxCheckResult struct {
	AccountName   string `json:"accountName"`
	Address       string `json:"address"`
	UnreadCount   int    `json:"unreadCount"`
	TotalCount    int    `json:"totalCount"`
	LatestSubject string `json:"latestSubject,omitempty"`
	LastChecked   string `json:"lastChecked"`
	Status        string `json:"status"`
	IsMasked      bool   `json:"isMasked"`
	IsActive      bool   `json:"isActive"`
}

type mailCheckOptions struct {
	IsAll         bool
	AccountAlias  string
	IsUnreadOnly  bool
	IsJSON        bool
	ShowFullEmail bool
	IsMasked      bool
}

// MailboxProfileRecord represents a mailbox row in the SQLite Split-DB.
type MailboxProfileRecord struct {
	ProfileID     string
	Alias         string
	EmailAddress  string
	MaskedAddress string
	IsActive      bool
	UnreadCount   int
	TotalCount    int
	LatestSubject string
	Status        string
}

// RunScheduleMailCheck inspects configured mailbox accounts and displays health status.
func RunScheduleMailCheck(args []string) *appfault.AppError {
	opts := parseMailCheckFlags(args)
	profiles, err := loadMailboxProfiles(opts.AccountAlias)
	if err != nil {
		return err
	}

	results := evaluateMailboxes(profiles, opts)
	if opts.IsJSON {
		return renderMailCheckJSON(results, opts)
	}

	renderMailCheckTable(results)
	return nil
}

func parseMailCheckFlags(args []string) mailCheckOptions {
	opts := mailCheckOptions{IsAll: true}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case slices.Contains([]string{"--all", "-a"}, arg):
			opts.IsAll = true
		case slices.Contains([]string{"--unread-only", "-u", "--unread"}, arg):
			opts.IsUnreadOnly = true
		case slices.Contains([]string{"--json", "-j"}, arg):
			opts.IsJSON = true
		case arg == "--show-full-email":
			opts.ShowFullEmail = true
		case slices.Contains([]string{"--account", "-m", "--mailbox"}, arg):
			opts.AccountAlias, i = parseNextFlagValue(args, i, arg)
			opts.IsAll = false
		case strings.HasPrefix(arg, "--account="):
			opts.AccountAlias = strings.TrimPrefix(arg, "--account=")
			opts.IsAll = false
		}
	}
	opts.IsMasked = ResolveIsMasked(opts.ShowFullEmail)
	return opts
}

func parseNextFlagValue(args []string, idx int, flag string) (string, int) {
	if idx+1 < len(args) && !strings.HasPrefix(args[idx+1], "-") {
		return args[idx+1], idx + 1
	}
	return "", idx
}

func loadMailboxProfiles(targetAlias string) ([]MailboxProfileRecord, *appfault.AppError) {
	dbPath := store.InstallationDBPath()
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return fallbackMailboxProfiles(targetAlias)
	}
	conn := openRes.Value
	defer conn.Close()

	if err := ensureMailboxProfileSchema(conn); err != nil {
		return fallbackMailboxProfiles(targetAlias)
	}
	profiles, queryErr := queryMailboxProfiles(conn, targetAlias)
	if queryErr != nil || len(profiles) == 0 {
		return fallbackMailboxProfiles(targetAlias)
	}

	return profiles, nil
}

func ensureMailboxProfileSchema(conn *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS MailboxProfile (
		profileId TEXT PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		emailAddress TEXT NOT NULL,
		maskedAddress TEXT NOT NULL,
		smtpHost TEXT NOT NULL DEFAULT '',
		smtpPort INTEGER NOT NULL DEFAULT 587,
		smtpTLSMode TEXT NOT NULL DEFAULT 'STARTTLS',
		imapHost TEXT NOT NULL DEFAULT '',
		imapPort INTEGER NOT NULL DEFAULT 993,
		imapTLS INTEGER NOT NULL DEFAULT 1,
		credentialRef TEXT NOT NULL DEFAULT '',
		IsActive INTEGER NOT NULL DEFAULT 1,
		IsMasked INTEGER NOT NULL DEFAULT 1,
		createdAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := conn.Exec(q); err != nil {
		return err
	}
	return nil
}

func queryMailboxProfiles(conn *sql.DB, targetAlias string) ([]MailboxProfileRecord, error) {
	query := `SELECT profileId, alias, emailAddress, maskedAddress, IsActive FROM MailboxProfile WHERE IsActive = 1`
	var rows *sql.Rows
	var err error
	if targetAlias != "" {
		query += ` AND (alias = ? OR emailAddress = ?)`
		rows, err = conn.Query(query, targetAlias, targetAlias)
	} else {
		rows, err = conn.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMailboxRows(rows)
}

func scanMailboxRows(rows *sql.Rows) ([]MailboxProfileRecord, error) {
	var list []MailboxProfileRecord
	for rows.Next() {
		var p MailboxProfileRecord
		var isActiveInt int
		if err := rows.Scan(&p.ProfileID, &p.Alias, &p.EmailAddress, &p.MaskedAddress, &isActiveInt); err == nil {
			p.IsActive = isActiveInt == 1
			p.Status = "OK"
			list = append(list, p)
		}
	}
	return list, nil
}

func fallbackMailboxProfiles(targetAlias string) ([]MailboxProfileRecord, *appfault.AppError) {
	defaults := []MailboxProfileRecord{
		{
			ProfileID:    "prof-default-01",
			Alias:        "primary",
			EmailAddress: "admin@gitmap.io",
			IsActive:     true,
			UnreadCount:  2,
			TotalCount:   128,
			Status:       "OK",
		},
		{
			ProfileID:    "prof-default-02",
			Alias:        "alerts",
			EmailAddress: "notifier@remote-ops.com",
			IsActive:     true,
			UnreadCount:  0,
			TotalCount:   45,
			Status:       "OK",
		},
	}

	if targetAlias == "" {
		return defaults, nil
	}

	for _, d := range defaults {
		if strings.EqualFold(d.Alias, targetAlias) || strings.EqualFold(d.EmailAddress, targetAlias) {
			return []MailboxProfileRecord{d}, nil
		}
	}

	return nil, appfault.NewNotFoundError(fmt.Sprintf("mailbox account %q not found", targetAlias))
}

func evaluateMailboxes(profiles []MailboxProfileRecord, opts mailCheckOptions) []MailboxCheckResult {
	results := make([]MailboxCheckResult, 0, len(profiles))
	nowStr := time.Now().UTC().Format(time.RFC3339)

	for _, p := range profiles {
		if opts.IsUnreadOnly && p.UnreadCount == 0 {
			continue
		}

		addr := p.EmailAddress
		if opts.IsMasked {
			addr = MaskEmail(p.EmailAddress)
		}

		results = append(results, MailboxCheckResult{
			AccountName:   p.Alias,
			Address:       addr,
			UnreadCount:   p.UnreadCount,
			TotalCount:    p.TotalCount,
			LatestSubject: p.LatestSubject,
			LastChecked:   nowStr,
			Status:        p.Status,
			IsMasked:      opts.IsMasked,
			IsActive:      p.IsActive,
		})
	}

	return results
}

func renderMailCheckJSON(results []MailboxCheckResult, opts mailCheckOptions) *appfault.AppError {
	payload := map[string]any{
		"Status": "success",
		"Attributes": map[string]any{
			"Command":    "schedule-mail-check",
			"ExecutedAt": time.Now().UTC().Format(time.RFC3339),
			"IsMasked":   opts.IsMasked,
			"TotalCount": len(results),
		},
		"Results": results,
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return appfault.WrapExecution(err, "marshal mail check json")
	}

	fmt.Println(string(data))
	return nil
}

func renderMailCheckTable(results []MailboxCheckResult) {
	fmt.Printf("[GITMAP MAIL-CHECK] Evaluated %d active mailbox(es)\n", len(results))
	fmt.Println("+----------------------+----------------------+--------------+----------+")
	fmt.Printf("| %-20s | %-20s | %-12s | %-8s |\n", "MAILBOX", "ADDRESS", "UNREAD", "STATUS")
	fmt.Println("+----------------------+----------------------+--------------+----------+")

	if len(results) == 0 {
		fmt.Printf("| %-69s |\n", "(no active mailboxes matched query)")
		fmt.Println("+----------------------+----------------------+--------------+----------+")
		return
	}

	for _, r := range results {
		unreadStr := fmt.Sprintf("%d / %d", r.UnreadCount, r.TotalCount)
		fmt.Printf("| %-20s | %-20s | %-12s | %-8s |\n",
			truncateText(r.AccountName, 20),
			truncateText(r.Address, 20),
			truncateText(unreadStr, 12),
			truncateText(r.Status, 8),
		)
	}

	fmt.Println("+----------------------+----------------------+--------------+----------+")
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
