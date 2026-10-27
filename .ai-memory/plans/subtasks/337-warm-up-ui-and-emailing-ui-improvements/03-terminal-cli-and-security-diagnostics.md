# Subtask 03: Terminal UI CLI Tools, SPF/DKIM/DMARC Diagnostics & Spam Score HUD

## Status: PLANNED
- **Subtask Slug:** `03-terminal-cli-and-security-diagnostics`
- **Parent Plan:** `.ai-memory/plans/pending/337-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/04-improvements-61-to-100-security-tokens-and-agy-integration.md`
- **Owner:** CLI & Security Subagent
- **Target Subsystem:** `cli/warmup/`, `cli/mail/`, `pkg/mail/auth/`, `pkg/mail/spam/`
- **Improvements Covered:** 51 to 70 (Clusters 6 and 7)
- **Exclusive Owned Files:**
  - `cli/warmup/warmup_cmd.go`
  - `cli/warmup/warmup_render.go`
  - `cli/mail/mail_cmd.go`
  - `cli/mail/mail_doctor.go`
  - `pkg/mail/auth/spf_validator.go`
  - `pkg/mail/auth/dkim_validator.go`
  - `pkg/mail/auth/dmarc_parser.go`
  - `pkg/mail/auth/dnsbl_checker.go`
  - `pkg/mail/spam/heuristic_scorer.go`

---

## 1. Executive Summary & Context

Subtask 03 couples developer-friendly command-line interfaces with deep cryptographic and deliverability diagnostic tools:
1. **Terminal UI CLI Tools (Improvements 51–60):** Rich Lipgloss/termpad terminal tools for inspecting warmup state, sending test messages, triggering health audits, and monitoring live log streams without opening a browser.
2. **Security & Authentication Diagnostics (Improvements 61–70):** SPF lookup depth matrix, DKIM public/private key bit-depth alignment, DMARC RUA parser, SpamAssassin/Rspamd heuristic scoring HUD, and parallel DNSBL blacklist scanner.

---

## 2. File Modification & Structural Ownership

```text
cli/
├── warmup/
│   ├── warmup_cmd.go           [NEW] - CLI dispatch for `gitmap warmup` subcommands
│   └── warmup_render.go        [NEW] - Lipgloss TokyoNight terminal layout & progress bars
├── mail/
│   ├── mail_cmd.go             [NEW] - CLI dispatch for `gitmap mail send/status`
│   └── mail_doctor.go          [NEW] - Interactive terminal diagnostics (`gitmap mail doctor`)
pkg/mail/
├── auth/
│   ├── spf_validator.go        [NEW] - Recursive SPF record lookup depth and mechanism tree
│   ├── dkim_validator.go       [NEW] - DNS TXT RSA key bit-depth and selector verifier
│   ├── dmarc_parser.go         [NEW] - DMARC policy inspector and RUA XML parser
│   └── dnsbl_checker.go        [NEW] - Concurrent DNSBL query runner (25 tier-1 blacklists)
└── spam/
    └── heuristic_scorer.go     [NEW] - Pre-send content, header, and MIME heuristic rule engine
```

---

## 3. Technical Requirements & Architecture

### 3.1 Terminal UI Commands
- `gitmap warmup status`: Renders visual ASCII progress bar, current daily quota, completed send count, and health grade.
- `gitmap warmup start [--curve=sigmoid]`: Launches background warmup runner for specified mailbox.
- `gitmap mail doctor [--domain=example.com]`: Executes full pre-flight authentication and reputation test suite.
- Supports `--json` flag across all commands for automated AI agent parsing.

### 3.2 DNS Authentication Diagnostics
- **SPF Validator:** Checks RFC 7208 compliance; recursively computes DNS query count and fails if > 10.
- **DKIM Validator:** Resolves `selector._domainkey.domain`; ensures key length is >= 2048 bits; validates base64 key structure.
- **DMARC Parser:** Parses `p=reject|quarantine|none`, `pct`, `sp`, `rua`, and `ruf` tags; alerts on permissive `p=none` configurations.
- **DNSBL Checker:** Dispatches parallel goroutines querying major blacklists with 800ms timeouts; returns list of active listings.

### 3.3 Spam Score Heuristic HUD
- Scans subject lines and message bodies for high-frequency promotional patterns, uppercase ratios, and invisible tracking pixels.
- Provides discrete score additions per matched rule; outputs aggregate score (0 to 10 scale).

---

## 4. Verification & Testing Checklist

- [ ] Unit tests in `pkg/mail/auth/spf_validator_test.go` checking simulated 5, 8, and 12-lookup records.
- [ ] Unit tests in `pkg/mail/auth/dkim_validator_test.go` checking 1024, 2048, and 4096-bit RSA keys.
- [ ] Heuristic scorer test validating detection of excessive capitalization and promotional trigger keywords.
- [ ] CLI rendering tests verifying zero ANSI color leaking and clean headless `--json` output.
- [ ] Coding guideline verification: function length <= 15 lines, error return sovereignty.
