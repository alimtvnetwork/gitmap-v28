# 06-storage-and-split-db-schema: Split-DB SQLite Storage & Schema Specification

- **Spec ID:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`
- **Subsystem:** SQLite Split-DB Engine & Mailing/Scheduling Schemas
- **Status:** APPROVED
- **Related Specs:**
  - `02-spec/21-app/06-database-and-split-db/01-architecture-spec.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md`
- **Implementation Subtask:** `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/05-subtask-embedded-mailing-and-schedule-uis.md`

---

## 1. Storage Architecture & Split-DB Partitioning

In accordance with GitMap's Three-Tier SQLite Split-DB architecture, mailing profiles, scheduled email jobs, folder watcher rules, and notification telemetry are partitioned into machine-scoped storage within Tier 2 (`installation.db`) and dedicated Tier 3 logging storage (`mail.db` or `repodb/mail.db`).

### 1.1 Database Tier Placement
- **Tier 2: Installation Database (`~/.gitmap/installation.db`):**
  - Stores persistent machine-wide configurations, mailbox connection credentials, folder watcher definitions, and recurring schedule rules.
  - Shared across all local repository workspaces.
- **Tier 3: Mail & Notification Telemetry (`.gitmap/repodb/mail.db` or `~/.gitmap/mail.db`):**
  - Stores transactional execution history, sent notification logs, delivery receipts, and debounced event audit records.
  - Prevents log bloat in primary configuration databases.

---

## 2. Schema Conventions & Guardrails

1. **PascalCase Tables:** All table names MUST be strict PascalCase:
   - `MailSchedule`
   - `MailboxProfile`
   - `FolderWatcherRule`
   - `NotificationLog`
   - `MailTemplate`
2. **camelCase Columns:** All column names MUST be strict camelCase:
   - `id`, `scheduleId`, `profileId`, `ruleId`, `folderPath`, `debounceSeconds`, `createdAt`.
3. **Positive Boolean Fields:** All boolean fields MUST use positive polarity (stored as `INTEGER` `0` or `1`):
   - `IsEnabled` (NOT `is_disabled`)
   - `IsActive` (NOT `is_inactive`)
   - `IsMasked` (NOT `is_unmasked`)
   - `HasAttachment` (NOT `is_attachmentless`)
   - `IsHTML` (NOT `is_plain_text`)
   - `IsDelivered` (NOT `is_undelivered`)
4. **Foreign Key Integrity:** `PRAGMA foreign_keys = ON;` is enforced on all opened database connections.

---

## 3. SQLite DDL Schema Definitions

### 3.1 Table: `MailboxProfile`
Stores configured IMAP/SMTP mailbox connection settings and privacy metadata.

```sql
CREATE TABLE IF NOT EXISTS MailboxProfile (
    profileId           TEXT PRIMARY KEY,
    alias               TEXT NOT NULL UNIQUE,
    emailAddress        TEXT NOT NULL,
    maskedAddress       TEXT NOT NULL,
    smtpHost            TEXT NOT NULL,
    smtpPort            INTEGER NOT NULL DEFAULT 587,
    smtpTLSMode         TEXT NOT NULL DEFAULT 'STARTTLS',
    imapHost            TEXT NOT NULL DEFAULT '',
    imapPort            INTEGER NOT NULL DEFAULT 993,
    imapTLS             INTEGER NOT NULL DEFAULT 1,
    credentialRef       TEXT NOT NULL,
    IsActive            INTEGER NOT NULL DEFAULT 1,
    IsMasked            INTEGER NOT NULL DEFAULT 1,
    createdAt           TEXT NOT NULL,
    updatedAt           TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mailbox_profile_alias ON MailboxProfile(alias);
CREATE INDEX IF NOT EXISTS idx_mailbox_profile_active ON MailboxProfile(IsActive);
```

### 3.2 Table: `MailSchedule`
Stores scheduled email dispatch tasks and periodic mailbox check jobs.

```sql
CREATE TABLE IF NOT EXISTS MailSchedule (
    scheduleId          TEXT PRIMARY KEY,
    jobType             TEXT NOT NULL, -- 'EMAIL_SEND', 'MAILBOX_CHECK'
    profileId           TEXT NOT NULL,
    cronExpression      TEXT NOT NULL DEFAULT '',
    intervalSeconds     INTEGER NOT NULL DEFAULT 0,
    nextRunAt           TEXT NOT NULL,
    lastRunAt           TEXT NOT NULL DEFAULT '',
    recipientList       TEXT NOT NULL DEFAULT '', -- JSON array of masked/unmasked emails
    subjectTemplate     TEXT NOT NULL DEFAULT '',
    bodyTemplate        TEXT NOT NULL DEFAULT '',
    IsHTML              INTEGER NOT NULL DEFAULT 0,
    HasAttachment       INTEGER NOT NULL DEFAULT 0,
    attachmentPaths     TEXT NOT NULL DEFAULT '[]', -- JSON array of file paths
    IsEnabled           INTEGER NOT NULL DEFAULT 1,
    createdAt           TEXT NOT NULL,
    updatedAt           TEXT NOT NULL,
    FOREIGN KEY(profileId) REFERENCES MailboxProfile(profileId) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_mail_schedule_next_run ON MailSchedule(nextRunAt, IsEnabled);
CREATE INDEX IF NOT EXISTS idx_mail_schedule_type ON MailSchedule(jobType);
```

### 3.3 Table: `FolderWatcherRule`
Stores event-driven folder and file mutation watcher configurations.

```sql
CREATE TABLE IF NOT EXISTS FolderWatcherRule (
    ruleId              TEXT PRIMARY KEY,
    ruleName            TEXT NOT NULL,
    folderPath          TEXT NOT NULL,
    watchEvents         TEXT NOT NULL, -- Comma-separated: 'CREATE,MODIFY,DELETE'
    filePattern         TEXT NOT NULL DEFAULT '*', -- Glob filter e.g. '*.go;*.md'
    debounceSeconds     INTEGER NOT NULL DEFAULT 8,
    profileId           TEXT NOT NULL,
    recipientList       TEXT NOT NULL, -- JSON array of recipient emails
    subjectTemplate     TEXT NOT NULL,
    bodyTemplate        TEXT NOT NULL,
    IsEnabled           INTEGER NOT NULL DEFAULT 1,
    lastTriggeredAt     TEXT NOT NULL DEFAULT '',
    triggerCount        INTEGER NOT NULL DEFAULT 0,
    createdAt           TEXT NOT NULL,
    updatedAt           TEXT NOT NULL,
    FOREIGN KEY(profileId) REFERENCES MailboxProfile(profileId) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_folder_watcher_enabled ON FolderWatcherRule(IsEnabled);
CREATE INDEX IF NOT EXISTS idx_folder_watcher_path ON FolderWatcherRule(folderPath);
```

### 3.4 Table: `NotificationLog`
Maintains an immutable audit ledger of all notification dispatches and watcher triggers.

```sql
CREATE TABLE IF NOT EXISTS NotificationLog (
    logId               TEXT PRIMARY KEY,
    ruleId              TEXT NOT NULL DEFAULT '',
    scheduleId          TEXT NOT NULL DEFAULT '',
    triggerSource       TEXT NOT NULL, -- 'FOLDER_WATCHER', 'SCHEDULED_JOB', 'MANUAL_DISPATCH'
    recipientSummary    TEXT NOT NULL, -- Masked recipient list preview
    subject             TEXT NOT NULL,
    mutatedFilesCount   INTEGER NOT NULL DEFAULT 0,
    mutatedFilesList    TEXT NOT NULL DEFAULT '[]', -- JSON array of affected file relative paths
    IsDelivered         INTEGER NOT NULL DEFAULT 0,
    errorMessage        TEXT NOT NULL DEFAULT '',
    dispatchDurationMs  INTEGER NOT NULL DEFAULT 0,
    createdAt           TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_notification_log_created ON NotificationLog(createdAt);
CREATE INDEX IF NOT EXISTS idx_notification_log_rule ON NotificationLog(ruleId);
CREATE INDEX IF NOT EXISTS idx_notification_log_delivered ON NotificationLog(IsDelivered);
```

### 3.5 Table: `MailTemplate`
Stores reusable subject and body email templates.

```sql
CREATE TABLE IF NOT EXISTS MailTemplate (
    templateId          TEXT PRIMARY KEY,
    templateName        TEXT NOT NULL UNIQUE,
    subjectTemplate     TEXT NOT NULL,
    bodyTemplate        TEXT NOT NULL,
    IsHTML              INTEGER NOT NULL DEFAULT 0,
    description         TEXT NOT NULL DEFAULT '',
    createdAt           TEXT NOT NULL,
    updatedAt           TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mail_template_name ON MailTemplate(templateName);
```

---

## 4. Concurrency Controls & Connection Pragmas

To prevent SQLite lock contention (`SQLITE_BUSY`) during concurrent CLI commands, background watcher routines, and web UI requests:

1. **WAL Journaling:** `PRAGMA journal_mode = WAL;` enabled at connection initialization.
2. **Synchronous Mode:** `PRAGMA synchronous = NORMAL;` optimizes write performance while guaranteeing filesystem consistency across crashes.
3. **Busy Timeout:** `PRAGMA busy_timeout = 5000;` waits up to 5,000 ms before returning a lock error.
4. **Single-Writer Pool:** Connection handles for transactional writes configured with `SetMaxOpenConns(1)` to eliminate parallel write races within the process.

---

## 5. Migration & Schema Bootstrap Engine

1. **Idempotent Bootstrapping:** Schema initialization executes `CREATE TABLE IF NOT EXISTS` blocks inside a single transaction.
2. **Version Tracking:** A lightweight `SchemaVersion` table records applied migration increments.
3. **Safe Column Extension:** Future column additions leverage `ALTER TABLE ... ADD COLUMN` guards to guarantee zero data loss on existing databases.
