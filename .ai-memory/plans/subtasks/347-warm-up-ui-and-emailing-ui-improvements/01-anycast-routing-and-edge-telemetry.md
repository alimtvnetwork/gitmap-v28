# Subtask Plan 01: Post-Anycast Routing & Edge Telemetry (Improvements 1001–1030)

- **Parent Task:** `347-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `347-01`
- **Spec Path:** `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/02-improvements-1001-to-1030-anycast-routing-and-edge-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 101: Post-Quantum TLS Handshake & Cryptographic Ciphers (1001–1010)
  - Cluster 102: Multi-Homed BGP Egress Binding & Dynamic Socket Balancing (1011–1020)
  - Cluster 103: RFC 8461 MTA-STS Autonomous Policy Discovery & DANE TLSA Sentinel (1021–1030)

## Implementation Steps

1. **Post-Quantum Cryptographic Ciphers & TLS Handshake Engine (1001–1010):**
   - Implement hybrid Kyber-1024 key encapsulation and Dilithium-5 digital signature verification in `pkg/mail/crypto/`.
   - Provide fallback cipher negotiation for legacy MTAs while prioritizing quantum-resistant key exchange.
   - Enforce positive booleans (`isQuantumSecure`, `hasValidCipher`, `isHybridHandshake`).

2. **Multi-Homed BGP Egress Binding & Socket Balancing (1011–1020):**
   - Architect dynamic interface binding across multihomed IP CIDR blocks in `pkg/mail/routing/`.
   - Implement low-latency round-robin and health-weighted socket dispatch pools.
   - Enforce strict error wrapping using `*apperror.AppError` and zero process exit calls in leaf functions.

3. **MTA-STS & DANE TLSA Autonomous Sentinel (1021–1030):**
   - Automate RFC 8461 MTA-STS policy discovery over HTTPS with cached policy persistence in `pkg/mail/mtasts/`.
   - Validate RFC 6698 DANE TLSA records via DNSSEC authenticated channels.
   - Persist egress routing and TLS telemetry into `repodb/pipeline.db` using positive schema flags (`isStsEnforced`, `hasDaneRecord`).

## Verification Checklist

- [ ] Post-quantum handshake benchmark negotiates hybrid Kyber-1024 ciphers without latency regression.
- [ ] Multi-homed socket manager routes across 4 virtual network interfaces with zero connection starvation.
- [ ] MTA-STS validator caches policy manifests and catches invalid certificates within 100ms.
