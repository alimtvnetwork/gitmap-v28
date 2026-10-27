# Subtask 03: Multi-Tenant Zero-Trust Mailbox Credential Vault & Forensic Header Trace

## Status: READY
- **Subtask Slug:** `03-credential-vault-and-header-trace`
- **Parent Plan:** `.ai-memory/plans/pending/343-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/04-improvements-661-to-700-sharding-audit-and-agy-fleet.md`
- **Target Range:** Improvements 661–680 (Clusters 67–68)
- **Target Subsystem:** `pkg/mail/vault/`, `pkg/mail/forensics/`, `pkg/mail/trace/`
- **Exclusive Owned Files:**
  - `pkg/mail/vault/enclave_manager.go`
  - `pkg/mail/vault/key_rotator.go`
  - `pkg/mail/vault/secret_cache.go`
  - `pkg/mail/vault/audit_logger.go`
  - `pkg/mail/vault/casbin_gate.go`
  - `pkg/mail/forensics/header_parser.go`
  - `pkg/mail/forensics/hop_timeline.go`
  - `pkg/mail/forensics/asn_lookup.go`
  - `pkg/mail/forensics/tls_cipher_inspector.go`
  - `pkg/mail/trace/visualizer.go`
  - `cli/cmdmail/vault_cmd.go`
  - `cli/cmdmail/trace_cmd.go`

---

## 1. Context & Objectives

This subtask delivers enterprise zero-trust mailbox credential protection and deep forensic email header hop tracking:
1. **Multi-Tenant Zero-Trust Mailbox Credential Vault & HSM Enclave HUD (Cluster 67, Items 661–670):** AES-256-GCM / ChaCha20-Poly1305 hardware-backed encryption, automated key rotation schedules, ephemeral memory caching, Casbin RBAC policy isolation, and tamper-evident audit logs.
2. **Forensic Header Trace & Cross-Gateway Hop Visualizer (Cluster 68, Items 671–680):** Deep RFC 5322 / RFC 5321 header decomposition, transit hop latency waterfall diagrams, intermediate gateway ASN/GeoIP enrichment, TLS cipher suite verification, and gateway tampering detection.
3. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isVaultUnlocked`, `hasValidKey`, `isEnclaveActive`, `hasValidSignature`, `isTlsEncrypted`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Mailbox Credential Vault & Enclave Engine
- Implement hardware-backed / envelope encryption manager in `pkg/mail/vault/enclave_manager.go`.
- Construct automated cryptographic key rotator in `pkg/mail/vault/key_rotator.go`.
- Implement zero-swap ephemeral memory credential cache in `pkg/mail/vault/secret_cache.go`.
- Build cryptographic audit logger with SHA-256 chain links in `pkg/mail/vault/audit_logger.go`.
- Implement fine-grained Casbin RBAC tenant gates in `pkg/mail/vault/casbin_gate.go`.

### Step 2: Forensic Header Parser & Hop Visualizer
- Implement RFC 5322 / RFC 5321 `Received:` header parser in `pkg/mail/forensics/header_parser.go`.
- Construct transit hop latency waterfall diagram builder in `pkg/mail/forensics/hop_timeline.go`.
- Build offline ASN and GeoIP lookup adapter in `pkg/mail/forensics/asn_lookup.go`.
- Implement intermediate TLS cipher suite and cipher grade inspector in `pkg/mail/forensics/tls_cipher_inspector.go`.
- Build interactive terminal and zero-CDN web hop visualizer in `pkg/mail/trace/visualizer.go`.

### Step 3: Command Integration
- Wire vault management subcommands into `cli/cmdmail/vault_cmd.go`.
- Wire forensic header trace commands into `cli/cmdmail/trace_cmd.go`.
