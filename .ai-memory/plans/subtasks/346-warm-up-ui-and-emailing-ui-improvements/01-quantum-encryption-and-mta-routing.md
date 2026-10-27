# Subtask Plan 01: Post-Quantum Encryption and Multi-Homed MTA Routing (Improvements 901–930)

- **Parent Task:** `346-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `346-01`
- **Spec Path:** `02-spec/21-app/346-warm-up-ui-and-emailing-ui-improvements/02-improvements-901-to-930-quantum-encryption-and-mta-routing.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 91: Post-Quantum TLS Handshake & Cryptographic Ciphers (901–910)
  - Cluster 92: Multi-Homed BGP Egress Binding & Dynamic Socket Balancing (911–920)
  - Cluster 93: RFC 8461 MTA-STS Autonomous Policy Discovery & DANE TLSA Sentinel (921–930)

## Implementation Steps

1. **Post-Quantum Cryptographic Ciphers & TLS Handshake Engine (901–910):**
   - Implement hybrid Kyber-1024 key encapsulation and Dilithium-5 digital signature verification in `pkg/mail/crypto/`.
   - Provide fallback cipher negotiation for legacy MTAs while prioritizing quantum-resistant key exchange.
   - Enforce positive booleans (`isQuantumSecure`, `hasValidCipher`, `isHybridHandshake`).

2. **Multi-Homed BGP Egress Binding & Socket Balancing (911–920):**
   - Architect dynamic interface binding across multihomed IP CIDR blocks in `pkg/mail/routing/`.
   - Implement low-latency round-robin and health-weighted socket dispatch pools.
   - Enforce strict error wrapping using `*apperror.AppError` and zero process exit calls in leaf functions.

3. **MTA-STS & DANE TLSA Autonomous Sentinel (921–930):**
   - Automate RFC 8461 MTA-STS policy discovery over HTTPS with cached policy persistence in `pkg/mail/mtasts/`.
   - Validate RFC 6698 DANE TLSA records via DNSSEC authenticated channels.
   - Persist egress routing and TLS telemetry into `repodb/pipeline.db` using positive schema flags (`isStsEnforced`, `hasDaneRecord`).

## Verification Checklist

- [ ] Post-quantum handshake benchmark negotiates hybrid Kyber-1024 ciphers without latency regression.
- [ ] Multi-homed socket manager routes across 4 virtual network interfaces with zero connection starvation.
- [ ] MTA-STS validator caches policy manifests and catches invalid certificates within 100ms.
