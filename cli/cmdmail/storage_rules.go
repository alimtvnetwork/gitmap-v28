package cmdmail

import (
	"time"
)

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
