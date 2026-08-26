# Split DB Architecture: Reset API Standard

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**Parent:** [00-overview.md](../00-overview.md)

---

## Overview

All CLI tools in the ecosystem must implement a standardized **2-Step Reset API** for safe data deletion and fresh starts. This document defines the common interface.

---

## 2-Step Confirmation Flow

### Why 2 Steps?

1. **Prevent accidental deletion** - User sees exactly what will be deleted
2. **Audit trail** - All reset requests are logged
3. **Time-bounded** - 5-minute confirmation window prevents stale requests
4. **Consistent UX** - Same pattern across all CLIs

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    UNIVERSAL 2-STEP RESET CONFIRMATION                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   STEP 1: REQUEST                                                            │
│   ─────────────────                                                          │
│   POST /api/v1/reset/request                                                │
│   {                                                                          │
│     "Scope": "<scope>"                                                       │
│   }                                                                          │
│                                                                              │
│   RESPONSE:                                                                  │
│   {                                                                          │
│     "ResetId": "rst_abc123def456",      ◀── Unique confirmation token       │
│     "Scope": "<scope>",                                                      │
│     "ExpiresAt": "<now + 5 minutes>",   ◀── TTL: 5 minutes                  │
│     "AffectedItems": { ... },           ◀── Preview of what will be deleted │
│     "Message": "Confirm within 5 minutes"                                   │
│   }                                                                          │
│                                                                              │
│   ════════════════════════════════════════════════════════════════════════  │
│                                                                              │
│   STEP 2: CONFIRM                                                            │
│   ───────────────                                                            │
│   POST /api/v1/reset/confirm                                                │
│   {                                                                          │
│     "ResetId": "rst_abc123def456"       ◀── Token from Step 1               │
│   }                                                                          │
│                                                                              │
│   RESPONSE:                                                                  │
│   {                                                                          │
│     "Status": "completed",                                                   │
│     "DeletedDatabases": <count>,                                            │
│     "FreedBytes": <bytes>                                                   │
│   }                                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Standard Endpoints

All CLIs must implement these three endpoints:

### 1. Request Reset

```
POST /api/v1/reset/request

Request:
{
  "Scope": "all" | "<cli-specific-scopes>"
}

Response:
{
  "ResetId": "rst_{uuid}",
  "Scope": "<scope>",
  "ExpiresAt": "<ISO8601 datetime>",
  "AffectedItems": {
    // CLI-specific breakdown
  },
  "Message": "Review affected items and confirm within 5 minutes"
}
```

### 2. Confirm Reset

```
POST /api/v1/reset/confirm

Request:
{
  "ResetId": "rst_{uuid}"
}

Response (Success):
{
  "Status": "completed",
  "DeletedDatabases": <count>,
  "FreedBytes": <bytes>,
  "Duration": "<duration>"
}

Response (Expired):
{
  "Status": "expired",
  "Error": "Reset confirmation expired. Please request again."
}

Response (Invalid):
{
  "Status": "invalid", 
  "Error": "Invalid or already used reset ID"
}
```

### 3. Cancel Reset

```
POST /api/v1/reset/cancel

Request:
{
  "ResetId": "rst_{uuid}"
}

Response:
{
  "Status": "cancelled"
}
```

---

## Database Schema (Required in Root DB)

Every CLI's root database must include this table:

```sql
-- linter-waive: MISSING-DESC-001 reason="Split-DB architecture example; focus on per-DB isolation, not free-text columns"
CREATE TABLE ResetRequest (
    ResetRequestId INTEGER PRIMARY KEY AUTOINCREMENT,
    ResetToken TEXT UNIQUE NOT NULL,               -- rst_{uuid}
    Scope TEXT NOT NULL,
    RequestedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ExpiresAt DATETIME NOT NULL,                   -- RequestedAt + 5 minutes
    AffectedItems TEXT,                            -- JSON: CLI-specific preview
    ConfirmedAt DATETIME,
    CancelledAt DATETIME,
    CompletedAt DATETIME,
    Status TEXT DEFAULT 'pending',                 -- pending, confirmed, expired, cancelled, completed
    DeletedCount INTEGER,
    FreedBytes INTEGER,
    ErrorMessage TEXT
);

CREATE INDEX IdxResetRequest_Status ON ResetRequest(Status, ExpiresAt);
```

---

## CLI-Specific Scopes

### AI Bridge CLI

| Scope | Description |
|-------|-------------|
| `all` | Full system reset |
| `app` | Application-level reset (requires `AppName`) |
| `chat` | Chat sessions only |
| `rag` | RAG memory only |
| `search` | Search cache only |
| `seo` | SEO data only |

### GSearch CLI

| Scope | Description |
|-------|-------------|
| `all` | Full system reset |
| `cache` | Search cache only |
| `history` | Search history only |

### BRun CLI

| Scope | Description |
|-------|-------------|
| `all` | Full system reset |
| `runs` | Run history only |
| `profile:{name}` | Specific profile runs |

### Nexus Flow CLI

| Scope | Description |
|-------|-------------|
| `all` | Full system reset |
| `executions` | All execution history |
| `pipeline:{id}` | Specific pipeline data |
| `checkpoints` | Checkpoints only |

---

## Configuration (config.seed.json)

All CLIs should include this in their seedable config:

```json
{
  "Reset": {
    "ConfirmationTtlMinutes": 5,
    "RequireConfirmation": true,
    "LogResets": true
  }
}
```

---

## Implementation Guidelines

### 1. Generate Reset Id

```go
func generateResetId() string {
    return fmt.Sprintf("rst_%s", uuid.New().String()[:12])
}
```

### 2. Calculate Expiration

```go
func calculateExpiration(settings *Settings) time.Time {
    ttl := settings.GetInt("Reset.ConfirmationTtlMinutes", 5)
    return time.Now().Add(time.Duration(ttl) * time.Minute)
}
```

### 3. Validate Reset Request

```go
func validateResetRequest(db *sql.DB, resetId string) apperror.Result[*ResetRequest] {
    var req ResetRequest
    err := db.QueryRow(`
        SELECT ResetRequestId, ResetToken, Scope, ExpiresAt, Status 
        FROM ResetRequest 
        WHERE ResetToken = ?
    `, resetId).Scan(&req.ResetRequestId, &req.ResetToken, &req.Scope, &req.ExpiresAt, &req.Status)
    
    if err == sql.ErrNoRows {
        return nil, ErrInvalidResetId
    }
    if req.Status != "pending" {
        return nil, ErrResetAlreadyProcessed
    }
    if time.Now().After(req.ExpiresAt) {
        // Update status to expired
        db.Exec("UPDATE ResetRequest SET Status = 'expired' WHERE ResetToken = ?", resetId)
        return nil, ErrResetExpired
    }
    
    return &req, nil
}
```

### 4. Execute Reset

```go
func executeReset(db *sql.DB, req *ResetRequest) apperror.Result[*ResetResult] {
    // Start transaction
    tx, _ := db.Begin()
    defer tx.Rollback()
    
    // Delete databases based on scope
    deleted, freed := deleteByScope(req.Scope)
    
    // Update reset request
    tx.Exec(`
        UPDATE ResetRequest 
        SET Status = 'completed', 
            ConfirmedAt = ?,
            CompletedAt = ?,
            DeletedCount = ?,
            FreedBytes = ?
        WHERE Id = ?
    `, time.Now(), time.Now(), deleted, freed, req.Id)
    
    tx.Commit()
    
    return &ResetResult{
        Status: "completed",
        DeletedDatabases: deleted,
        FreedBytes: freed,
    }, nil
}
```

---

## Error Codes

Standard error codes for reset operations (add to each CLI's error registry):

| Range | CLI | Error Codes |
|-------|-----|-------------|
| 9401-9409 | AI Bridge | Reset errors |
| 7401-7409 | GSearch | Reset errors |
| 6401-6409 | BRun | Reset errors |
| 8401-8409 | Nexus Flow | Reset errors |

Common errors:
- `x401`: Reset expired
- `x402`: Invalid reset ID
- `x403`: Reset already confirmed
- `x404`: Reset cancelled

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Split DB Overview | `../00-overview.md` |
| CLI Examples | `./01-cli-examples.md` |
| AI Bridge Reset | `spec/22-ai-bridge-cli/01-backend/14-reset-and-export-api.md` |
| GSearch Database | `spec/20-gsearch-cli/01-backend/22-database-architecture.md` |
| BRun Database | `spec/21-brun-cli/01-backend/16-database-architecture.md` |
| Nexus Flow Database | `spec/24-nexus-flow-cli/01-backend/05-database-architecture.md` |
