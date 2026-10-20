# 04-privacy-and-email-masking: Privacy-Preserving Email Masking Specification

- **Spec ID:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Privacy Protection, Output Sanitization, CLI Envelopes, PII Redaction
- **Dependencies:** `cli/cmdmail`, `cli/cmdschedule`, `cli/jsonenv`, `cli/apperror`
- **Associated Plan:** `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`
- **Subtask Plan:** `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/03-subtask-multi-mailbox-and-schedule-email.md`

---

## 1. Executive Summary & Privacy Threat Model

GitMap commands frequently interact with sensitive communication channels, recipient rosters, server accounts, and alert dispatch lists. Displaying raw email addresses in terminal output, CI/CD runners, captured console logs, and public JSON envelopes exposes personal identifiable information (PII) to:
1. **Shoulder Surfing & Screen Recording:** Accidental exposure during demos, streaming, or pair programming.
2. **Public CI/CD Log Leakage:** Automated build steps capturing raw command output in public or shared pipeline consoles.
3. **Envelope Serialization Overexposure:** Machine-readable payloads stored in external diagnostics or bug reports.

This specification mandates a **Zero-Plaintext PII Invariant** across all GitMap display surfaces by default, defining deterministic masking standards and authenticated unmasking controls.

---

## 2. Zero-Plaintext PII Invariant

Across all CLI output, terminal tables (`termtable`), log messages, diagnostic reports, and embedded UI views:
- **Default Behavior:** Every email address (sender, recipient, CC/BCC, mailbox username) MUST be rendered in masked format (e.g., `u***r@d***n.com`).
- **Unmask Authorization:** Plaintext email display is strictly forbidden unless the caller explicitly passes `--show-full-email` or sets the environment guard `GITMAP_ALLOW_UNMASKED_PII=1`.
- **Decoupled Persistence:** Masking is an egress presentation filter. It must never overwrite or mutate canonical email addresses stored in Split-DB or secret vaults.

---

## 3. Email Masking Standard & Deterministic Algorithm

### 3.1 Masking Representation Format
An email address is split into two components: `[LocalPart]@[DomainPart]`. Both components undergo deterministic masking with preserved boundaries for recognition.

```
Raw Email:      user.name@example.com
Masked Email:   u***e@e***e.com

Raw Email:      devops@company.org
Masked Email:   d***s@c***y.org
```

### 3.2 Component Transformation Rules

#### 3.2.1 Local Part Transformation
Let $L$ be the string of the local part with length $|L|$:
- If $|L| \ge 3$: Output `L[0]` + `***` + `L[|L|-1]` (e.g., `"john"` $\to$ `"j***n"`).
- If $|L| == 2$: Output `L[0]` + `***` (e.g., `"me"` $\to$ `"m***"`).
- If $|L| == 1$: Output `L[0]` + `***` (e.g., `"a"` $\to$ `"a***"`).
- If $|L| == 0$: Output `"***"`.

#### 3.2.2 Domain Part Transformation
The domain is split by period `.` into domain labels and the final Top-Level Domain (TLD):
- **Root Domain Label:** For the primary domain name label (e.g., `"example"` in `"example.com"`):
  - If length $\ge 3$: First char + `***` + last char (e.g., `"e***e"`).
  - If length $< 3$: First char + `***` (e.g., `"i***"`).
- **Subdomain Labels:** Each subdomain label is transformed identically (e.g., `"mail.internal.org"` $\to$ `"m***l.i***l.org"`).
- **TLD Preservation:** The final TLD (e.g., `.com`, `.org`, `.net`, `.io`) is preserved unmasked to retain context for network routing and classification.

#### 3.2.3 Edge Case Fallbacks
- **Invalid / Malformed Emails:** If the input lacks an `@` symbol or contains illegal characters, the algorithm masks all characters except the first and last with `***` (e.g., `"invalidstring"` $\to$ `"i***g"`).
- **Recipient Lists:** Comma-separated or semicolon-separated strings are parsed into individual addresses, each masked individually and rejoined (e.g., `"a@b.com, c@d.org"` $\to$ `"a***@b***.com, c***@d***.org"`).

---

## 4. Go Type Definitions & Masking Engine API

```go
package cmdschedule

import (
    "strings"
)

type PrivacyMasker struct {
    AllowUnmasked bool
}

func NewPrivacyMasker(allowUnmasked bool) *PrivacyMasker {
    return &PrivacyMasker{
        AllowUnmasked: allowUnmasked,
    }
}

func (m *PrivacyMasker) MaskAddress(rawEmail string) string {
    if m.AllowUnmasked {
        return rawEmail
    }
    return MaskEmailAddress(rawEmail)
}

func (m *PrivacyMasker) MaskList(emails []string) []string {
    if m.AllowUnmasked {
        return emails
    }
    masked := make([]string, len(emails))
    for i, e := range emails {
        masked[i] = MaskEmailAddress(e)
    }
    return masked
}

func MaskEmailAddress(raw string) string {
    trimmed := strings.TrimSpace(raw)
    if trimmed == "" {
        return ""
    }
    parts := strings.Split(trimmed, "@")
    if len(parts) != 2 {
        return maskSimpleToken(trimmed)
    }
    localPart := maskToken(parts[0])
    domainPart := maskDomain(parts[1])
    return localPart + "@" + domainPart
}

func maskToken(token string) string {
    n := len(token)
    if n <= 1 {
        return token + "***"
    }
    if n == 2 {
        return string(token[0]) + "***"
    }
    return string(token[0]) + "***" + string(token[n-1])
}

func maskDomain(domain string) string {
    segments := strings.Split(domain, ".")
    if len(segments) <= 1 {
        return maskToken(domain)
    }
    lastIdx := len(segments) - 1
    tld := segments[lastIdx]
    var maskedLabels []string
    for i := 0; i < lastIdx; i++ {
        maskedLabels = append(maskedLabels, maskToken(segments[i]))
    }
    maskedLabels = append(maskedLabels, tld)
    return strings.Join(maskedLabels, ".")
}

func maskSimpleToken(s string) string {
    if len(s) <= 2 {
        return "***"
    }
    return string(s[0]) + "***" + string(s[len(s)-1])
}
```

---

## 5. Unmask Authorization & Security Controls

### 5.1 CLI Flag: `--show-full-email`
All email and scheduler commands accept an optional flag:
```bash
gitmap schedule mail-check --all --show-full-email
gitmap mail check --account work --show-full-email
```

### 5.2 Environment Override
Automated diagnostic pipelines running in secured environments may activate:
```bash
export GITMAP_ALLOW_UNMASKED_PII=1
```
When this variable is set to `1` or `true`, the unmasked representation is enabled globally for machine consumption.

### 5.3 Audit Telemetry
Whenever unmasking is active, GitMap records an audit log entry in the session diagnostics:
```text
[SECURITY AUDIT] Full PII email address unmasking requested by user via --show-full-email flag.
```

---

## 6. Output Surface Interception & Envelope Integration

### 6.1 Terminal Output Box (`termtable`)
When displaying mailbox status or scheduled mail tasks:

```text
+---------------------+-------------------+------------+----------------------+
| TASK ID             | SENDER            | STATUS     | NOTIFIER RECIPIENTS  |
+---------------------+-------------------+------------+----------------------+
| sch-mail-001        | d***s@c***y.org   | SCHEDULED  | a***t@e***e.com      |
| sch-mail-002        | b***d@i***l.net   | ACTIVE     | t***m@c***y.org      |
+---------------------+-------------------+------------+----------------------+
```

### 6.2 JSON Envelope Sanitization
When `--json` is supplied without `--show-full-email`, JSON output fields are masked automatically:

```json
{
  "version": "2.0",
  "status": "success",
  "data": {
    "taskId": "sch-mail-001",
    "maskedSender": "d***s@c***y.org",
    "maskedRecipients": [
      "a***t@e***e.com",
      "o***s@c***y.org"
    ],
    "isUnmasked": false
  }
}
```

---

## 7. Storage & In-Memory Isolation

To protect system integrity while upholding privacy:
1. **Raw Storage in Split-DB:** Tables `EmailWatcherRule`, `MailHistoryLog`, and `ScheduledMailJob` retain canonical addresses in their respective database columns so delivery pipelines function accurately.
2. **Masking at Boundary:** Functions that format tables, construct error envelopes, or print stdout strings execute `MaskAddress()` immediately before emitting text.
3. **Error Redaction:** `apperror.AppError` string formatters sanitize any email pattern detected in the error message before printing to terminal or saving to standard logs.

---

## 8. Verification & Acceptance Criteria

- **AC-PM-01:** Standard email `"developer@example.com"` is masked as `"d***r@e***e.com"` in terminal output.
- **AC-PM-02:** Short username `"me@test.org"` is masked as `"m***@t***t.org"`.
- **AC-PM-03:** Subdomain address `"alert@mon.corp.net"` is masked as `"a***t@m***n.c***p.net"`.
- **AC-PM-04:** Providing `--show-full-email` renders `"developer@example.com"` verbatim and logs an audit record.
- **AC-PM-05:** JSON output with `--json` produces masked strings in `data` unless `--show-full-email` is explicitly supplied.
- **AC-PM-06:** No plaintext email address appears in terminal error strings when an invalid recipient fails delivery.
