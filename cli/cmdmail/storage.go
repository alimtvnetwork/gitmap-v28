package cmdmail

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// MailStore wraps SQLite connection for mailing and watcher records.
type MailStore struct {
	db *sql.DB
}

// MailTemplateRecord stores reusable email template definitions.
type MailTemplateRecord struct {
	TemplateID      string `json:"templateId"`
	TemplateName    string `json:"templateName"`
	SubjectTemplate string `json:"subjectTemplate"`
	BodyTemplate    string `json:"bodyTemplate"`
	IsHTML          bool   `json:"isHTML"`
	Description     string `json:"description"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// FolderWatcherRuleRecord stores folder monitoring configurations.
type FolderWatcherRuleRecord struct {
	RuleID          string `json:"ruleId"`
	RuleName        string `json:"ruleName"`
	FolderPath      string `json:"folderPath"`
	WatchEvents     string `json:"watchEvents"`
	FilePattern     string `json:"filePattern"`
	DebounceSeconds int    `json:"debounceSeconds"`
	ProfileID       string `json:"profileId"`
	RecipientList   string `json:"recipientList"`
	SubjectTemplate string `json:"subjectTemplate"`
	BodyTemplate    string `json:"bodyTemplate"`
	IsEnabled       bool   `json:"isEnabled"`
	LastTriggeredAt string `json:"lastTriggeredAt"`
	TriggerCount    int    `json:"triggerCount"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// NotificationLogRecord tracks notification events and delivery status.
type NotificationLogRecord struct {
	LogID              string `json:"logId"`
	RuleID             string `json:"ruleId"`
	ScheduleID         string `json:"scheduleId"`
	TriggerSource      string `json:"triggerSource"`
	RecipientSummary   string `json:"recipientSummary"`
	Subject            string `json:"subject"`
	MutatedFilesCount  int    `json:"mutatedFilesCount"`
	MutatedFilesList   string `json:"mutatedFilesList"`
	IsDelivered        bool   `json:"isDelivered"`
	ErrorMessage       string `json:"errorMessage"`
	DispatchDurationMs int64  `json:"dispatchDurationMs"`
	CreatedAt          string `json:"createdAt"`
}

// MailboxProfileRecord holds connection configuration for a mailbox.
type MailboxProfileRecord struct {
	ProfileID     string `json:"profileId"`
	Alias         string `json:"alias"`
	EmailAddress  string `json:"emailAddress"`
	MaskedAddress string `json:"maskedAddress"`
	SMTPHost      string `json:"smtpHost"`
	SMTPPort      int    `json:"smtpPort"`
	SMTPTLSMode   string `json:"smtpTLSMode"`
	CredentialRef string `json:"credentialRef"`
	IsActive      bool   `json:"isActive"`
	IsMasked      bool   `json:"isMasked"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// OpenMailStore initializes a SQLite database connection for mail storage.
func OpenMailStore(customPath string) (*MailStore, error) {
	targetPath := resolveMailStorePath(customPath)
	db, err := sql.Open("sqlite", targetPath)
	if err != nil {
		return nil, err
	}
	if err := store.ConfigureSQLiteConn(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	s := &MailStore{db: db}
	if err := s.InitSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func resolveMailStorePath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	return store.ResolveSplitDbPath("mail", "default", "")
}

// Close closes the underlying SQLite database.
func (s *MailStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// InitSchema ensures all required tables exist in PascalCase with positive booleans.
func (s *MailStore) InitSchema() error {
	ddls := []string{
		ddlMailTemplate,
		ddlFolderWatcherRule,
		ddlNotificationLog,
		ddlMailboxProfile,
	}
	for _, ddl := range ddls {
		if _, err := s.db.Exec(ddl); err != nil {
			return err
		}
	}
	return nil
}

const ddlMailTemplate = `
CREATE TABLE IF NOT EXISTS MailTemplate (
    templateId TEXT PRIMARY KEY,
    templateName TEXT NOT NULL UNIQUE,
    subjectTemplate TEXT NOT NULL,
    bodyTemplate TEXT NOT NULL,
    IsHTML INTEGER NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    createdAt TEXT NOT NULL,
    updatedAt TEXT NOT NULL
);`

const ddlFolderWatcherRule = `
CREATE TABLE IF NOT EXISTS FolderWatcherRule (
    ruleId TEXT PRIMARY KEY,
    ruleName TEXT NOT NULL UNIQUE,
    folderPath TEXT NOT NULL,
    watchEvents TEXT NOT NULL,
    filePattern TEXT NOT NULL DEFAULT '*',
    debounceSeconds INTEGER NOT NULL DEFAULT 8,
    profileId TEXT NOT NULL DEFAULT 'default',
    recipientList TEXT NOT NULL,
    subjectTemplate TEXT NOT NULL,
    bodyTemplate TEXT NOT NULL,
    IsEnabled INTEGER NOT NULL DEFAULT 1,
    lastTriggeredAt TEXT NOT NULL DEFAULT '',
    triggerCount INTEGER NOT NULL DEFAULT 0,
    createdAt TEXT NOT NULL,
    updatedAt TEXT NOT NULL
);`

const ddlNotificationLog = `
CREATE TABLE IF NOT EXISTS NotificationLog (
    logId TEXT PRIMARY KEY,
    ruleId TEXT NOT NULL DEFAULT '',
    scheduleId TEXT NOT NULL DEFAULT '',
    triggerSource TEXT NOT NULL,
    recipientSummary TEXT NOT NULL,
    subject TEXT NOT NULL,
    mutatedFilesCount INTEGER NOT NULL DEFAULT 0,
    mutatedFilesList TEXT NOT NULL DEFAULT '[]',
    IsDelivered INTEGER NOT NULL DEFAULT 0,
    errorMessage TEXT NOT NULL DEFAULT '',
    dispatchDurationMs INTEGER NOT NULL DEFAULT 0,
    createdAt TEXT NOT NULL
);`

const ddlMailboxProfile = `
CREATE TABLE IF NOT EXISTS MailboxProfile (
    profileId TEXT PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    emailAddress TEXT NOT NULL,
    maskedAddress TEXT NOT NULL,
    smtpHost TEXT NOT NULL,
    smtpPort INTEGER NOT NULL DEFAULT 587,
    smtpTLSMode TEXT NOT NULL DEFAULT 'STARTTLS',
    credentialRef TEXT NOT NULL DEFAULT '',
    IsActive INTEGER NOT NULL DEFAULT 1,
    IsMasked INTEGER NOT NULL DEFAULT 1,
    createdAt TEXT NOT NULL,
    updatedAt TEXT NOT NULL
);`

// SaveTemplate persists or updates an email template.
func (s *MailStore) SaveTemplate(tmpl MailTemplateRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	isHTMLInt := boolToInt(tmpl.IsHTML)
	q := `INSERT INTO MailTemplate (templateId, templateName, subjectTemplate, bodyTemplate, IsHTML, description, createdAt, updatedAt)
          VALUES (?, ?, ?, ?, ?, ?, ?, ?)
          ON CONFLICT(templateName) DO UPDATE SET
          subjectTemplate=excluded.subjectTemplate, bodyTemplate=excluded.bodyTemplate,
          IsHTML=excluded.IsHTML, description=excluded.description, updatedAt=excluded.updatedAt;`
	_, err := s.db.Exec(q, tmpl.TemplateID, tmpl.TemplateName, tmpl.SubjectTemplate, tmpl.BodyTemplate, isHTMLInt, tmpl.Description, now, now)
	return err
}

// GetTemplate fetches a template by name.
func (s *MailStore) GetTemplate(name string) (*MailTemplateRecord, error) {
	row := s.db.QueryRow(`SELECT templateId, templateName, subjectTemplate, bodyTemplate, IsHTML, description, createdAt, updatedAt FROM MailTemplate WHERE templateName = ?`, name)
	var t MailTemplateRecord
	var isHTMLInt int
	err := row.Scan(&t.TemplateID, &t.TemplateName, &t.SubjectTemplate, &t.BodyTemplate, &isHTMLInt, &t.Description, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	t.IsHTML = isHTMLInt == 1
	return &t, nil
}

// ListTemplates retrieves all email templates.
func (s *MailStore) ListTemplates() ([]MailTemplateRecord, error) {
	rows, err := s.db.Query(`SELECT templateId, templateName, subjectTemplate, bodyTemplate, IsHTML, description, createdAt, updatedAt FROM MailTemplate ORDER BY templateName ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []MailTemplateRecord
	for rows.Next() {
		var t MailTemplateRecord
		var isHTMLInt int
		if err := rows.Scan(&t.TemplateID, &t.TemplateName, &t.SubjectTemplate, &t.BodyTemplate, &isHTMLInt, &t.Description, &t.CreatedAt, &t.UpdatedAt); err == nil {
			t.IsHTML = isHTMLInt == 1
			results = append(results, t)
		}
	}
	return results, nil
}

// DeleteTemplate removes a template by name.
func (s *MailStore) DeleteTemplate(name string) error {
	_, err := s.db.Exec(`DELETE FROM MailTemplate WHERE templateName = ?`, name)
	return err
}

// SaveWatcherRule stores or updates a folder watcher rule.
func (s *MailStore) SaveWatcherRule(rule FolderWatcherRuleRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	isEnabledInt := boolToInt(rule.IsEnabled)
	q := `INSERT INTO FolderWatcherRule (ruleId, ruleName, folderPath, watchEvents, filePattern, debounceSeconds, profileId, recipientList, subjectTemplate, bodyTemplate, IsEnabled, lastTriggeredAt, triggerCount, createdAt, updatedAt)
          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
          ON CONFLICT(ruleName) DO UPDATE SET
          folderPath=excluded.folderPath, watchEvents=excluded.watchEvents, filePattern=excluded.filePattern,
          debounceSeconds=excluded.debounceSeconds, recipientList=excluded.recipientList, subjectTemplate=excluded.subjectTemplate,
          bodyTemplate=excluded.bodyTemplate, IsEnabled=excluded.IsEnabled, updatedAt=excluded.updatedAt;`
	_, err := s.db.Exec(q, rule.RuleID, rule.RuleName, rule.FolderPath, rule.WatchEvents, rule.FilePattern, rule.DebounceSeconds, rule.ProfileID, rule.RecipientList, rule.SubjectTemplate, rule.BodyTemplate, isEnabledInt, rule.LastTriggeredAt, rule.TriggerCount, now, now)
	return err
}

// ListWatcherRules retrieves all configured watcher rules.
func (s *MailStore) ListWatcherRules() ([]FolderWatcherRuleRecord, error) {
	rows, err := s.db.Query(`SELECT ruleId, ruleName, folderPath, watchEvents, filePattern, debounceSeconds, profileId, recipientList, subjectTemplate, bodyTemplate, IsEnabled, lastTriggeredAt, triggerCount, createdAt, updatedAt FROM FolderWatcherRule ORDER BY createdAt DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []FolderWatcherRuleRecord
	for rows.Next() {
		var r FolderWatcherRuleRecord
		var isEnabledInt int
		if err := rows.Scan(&r.RuleID, &r.RuleName, &r.FolderPath, &r.WatchEvents, &r.FilePattern, &r.DebounceSeconds, &r.ProfileID, &r.RecipientList, &r.SubjectTemplate, &r.BodyTemplate, &isEnabledInt, &r.LastTriggeredAt, &r.TriggerCount, &r.CreatedAt, &r.UpdatedAt); err == nil {
			r.IsEnabled = isEnabledInt == 1
			results = append(results, r)
		}
	}
	return results, nil
}

// RecordNotificationLog saves an immutable notification dispatch entry.
func (s *MailStore) RecordNotificationLog(log NotificationLogRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	isDeliveredInt := boolToInt(log.IsDelivered)
	q := `INSERT INTO NotificationLog (logId, ruleId, scheduleId, triggerSource, recipientSummary, subject, mutatedFilesCount, mutatedFilesList, IsDelivered, errorMessage, dispatchDurationMs, createdAt)
          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := s.db.Exec(q, log.LogID, log.RuleID, log.ScheduleID, log.TriggerSource, log.RecipientSummary, log.Subject, log.MutatedFilesCount, log.MutatedFilesList, isDeliveredInt, log.ErrorMessage, log.DispatchDurationMs, now)
	return err
}

// ListNotificationLogs retrieves recent notification logs.
func (s *MailStore) ListNotificationLogs(limit int) ([]NotificationLogRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(`SELECT logId, ruleId, scheduleId, triggerSource, recipientSummary, subject, mutatedFilesCount, mutatedFilesList, IsDelivered, errorMessage, dispatchDurationMs, createdAt FROM NotificationLog ORDER BY createdAt DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []NotificationLogRecord
	for rows.Next() {
		var l NotificationLogRecord
		var isDeliveredInt int
		if err := rows.Scan(&l.LogID, &l.RuleID, &l.ScheduleID, &l.TriggerSource, &l.RecipientSummary, &l.Subject, &l.MutatedFilesCount, &l.MutatedFilesList, &isDeliveredInt, &l.ErrorMessage, &l.DispatchDurationMs, &l.CreatedAt); err == nil {
			l.IsDelivered = isDeliveredInt == 1
			results = append(results, l)
		}
	}
	return results, nil
}

// SaveMailboxProfile persists a mailbox profile.
func (s *MailStore) SaveMailboxProfile(prof MailboxProfileRecord) error {
	now := time.Now().UTC().Format(time.RFC3339)
	isActiveInt := boolToInt(prof.IsActive)
	isMaskedInt := boolToInt(prof.IsMasked)
	q := `INSERT INTO MailboxProfile (profileId, alias, emailAddress, maskedAddress, smtpHost, smtpPort, smtpTLSMode, credentialRef, IsActive, IsMasked, createdAt, updatedAt)
          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
          ON CONFLICT(alias) DO UPDATE SET
          emailAddress=excluded.emailAddress, maskedAddress=excluded.maskedAddress,
          smtpHost=excluded.smtpHost, smtpPort=excluded.smtpPort, smtpTLSMode=excluded.smtpTLSMode,
          credentialRef=excluded.credentialRef, IsActive=excluded.IsActive, IsMasked=excluded.IsMasked, updatedAt=excluded.updatedAt;`
	_, err := s.db.Exec(q, prof.ProfileID, prof.Alias, prof.EmailAddress, prof.MaskedAddress, prof.SMTPHost, prof.SMTPPort, prof.SMTPTLSMode, prof.CredentialRef, isActiveInt, isMaskedInt, now, now)
	return err
}

// ListMailboxProfiles retrieves all configured mailbox profiles.
func (s *MailStore) ListMailboxProfiles() ([]MailboxProfileRecord, error) {
	rows, err := s.db.Query(`SELECT profileId, alias, emailAddress, maskedAddress, smtpHost, smtpPort, smtpTLSMode, credentialRef, IsActive, IsMasked, createdAt, updatedAt FROM MailboxProfile ORDER BY alias ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []MailboxProfileRecord
	for rows.Next() {
		var p MailboxProfileRecord
		var isActiveInt, isMaskedInt int
		if err := rows.Scan(&p.ProfileID, &p.Alias, &p.EmailAddress, &p.MaskedAddress, &p.SMTPHost, &p.SMTPPort, &p.SMTPTLSMode, &p.CredentialRef, &isActiveInt, &isMaskedInt, &p.CreatedAt, &p.UpdatedAt); err == nil {
			p.IsActive = isActiveInt == 1
			p.IsMasked = isMaskedInt == 1
			results = append(results, p)
		}
	}
	return results, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
