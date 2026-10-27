# Improvements 901 to 930: Quantum Encryption, Multi-MTA Routing & Jitter Waveform TUI

This document details the complete specifications for Improvements 901 through 930 across Clusters 91, 92, and 93 of Sequence 346.

---

# Cluster 91: Autonomous Envelope Polyglot Encryption & Quantum-Resistant Handshakes (901–910)

---

### Improvement 901: NIST ML-KEM-768 / X25519 Hybrid Key Encapsulation Engine
- **ID**: `IMP-901`
- **Cluster**: 91 (Autonomous Envelope Polyglot Encryption & Quantum-Resistant Handshakes)
- **Goal**: Negotiates a hybrid post-quantum shared secret combining classical X25519 elliptic-curve Diffie-Hellman and NIST ML-KEM-768 (Kyber) for every outbound warming envelope handshake.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status card displaying hybrid key exchange metrics. Uses neon cyan (`#00F0FF`) for Kyber encapsulations, electric violet (`#7000FF`) for X25519, and emerald green (`#00E676`) for combined entropy verification.
  - *Web*: Custom Web Component `<hybrid-kem-visualizer>` displaying animated lattice key exchange vectors with native SVG (zero external CDN).
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum kem-handshake --target-mx=<host> [--kem-algorithm=<algo>] [--is-dry-run] [--has-pqc-enforced]`
  - Flag `--target-mx`: Hostname or IP of the destination mail exchanger.
  - Flag `--kem-algorithm`: KEM variant to encapsulate (default: `ML-KEM-768`).
  - Flag `--is-dry-run`: Performs simulated key encapsulation without transmitting wire packets.
  - Flag `--has-pqc-enforced`: Rejects connection if destination MX fails post-quantum negotiation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS KemHandshakes (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL UNIQUE,
      targetMxHost TEXT NOT NULL,
      kemAlgorithm TEXT NOT NULL DEFAULT 'ML-KEM-768',
      classicalAlgorithm TEXT NOT NULL DEFAULT 'X25519',
      isPqcEnforced INTEGER NOT NULL DEFAULT 1,
      hasHybridVerified INTEGER NOT NULL DEFAULT 1,
      handshakeDurationMs REAL NOT NULL DEFAULT 0.0,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type KemHandshake struct {
      SessionUUID         string    `json:"sessionUuid"`
      TargetMxHost        string    `json:"targetMxHost"`
      KemAlgorithm        string    `json:"kemAlgorithm"`
      ClassicalAlgorithm  string    `json:"classicalAlgorithm"`
      IsPqcEnforced       bool      `json:"isPqcEnforced"`
      HasHybridVerified   bool      `json:"hasHybridVerified"`
      HandshakeDurationMs float64   `json:"handshakeDurationMs"`
      CreatedAt           time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/crypto/kem.go` during initial TLS transport negotiation.

---

### Improvement 902: NIST ML-DSA-65 / Ed25519 Dual Envelope Digital Signer
- **ID**: `IMP-902`
- **Cluster**: 91
- **Goal**: Generates dual cryptographic digital signatures across warming message envelopes using classical Ed25519 and lattice-based ML-DSA-65 (Dilithium3) to resist quantum forgeability.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dual-signature badge displaying Ed25519 and ML-DSA signature digests side-by-side with border highlights in lavender (`#B388FF`).
  - *Web*: `<dual-signature-badge>` Web Component rendering interactive cryptographic proof cards with SVG lock glyphs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum sign-envelope --envelope-id=<uuid> [--signature-scheme=<name>] [--is-strict-eval] [--has-audit-hash]`
  - Flag `--envelope-id`: UUID of the warming envelope to sign.
  - Flag `--signature-scheme`: Signing mechanism (default: `ML-DSA-65+Ed25519`).
  - Flag `--is-strict-eval`: Requires both signatures to pass cryptographic verification before spooling.
  - Flag `--has-audit-hash`: Automatically writes signature hash into `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EnvelopeSignatures (
      id TEXT PRIMARY KEY,
      envelopeUuid TEXT NOT NULL UNIQUE,
      signatureScheme TEXT NOT NULL DEFAULT 'ML-DSA-65+Ed25519',
      ed25519DigestHex TEXT NOT NULL,
      mldsaDigestHex TEXT NOT NULL,
      isStrictEval INTEGER NOT NULL DEFAULT 1,
      hasAuditHash INTEGER NOT NULL DEFAULT 1,
      signedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EnvelopeSignature struct {
      EnvelopeUUID    string    `json:"envelopeUuid"`
      SignatureScheme string    `json:"signatureScheme"`
      Ed25519DigestHex string   `json:"ed25519DigestHex"`
      MldsaDigestHex  string    `json:"mldsaDigestHex"`
      IsStrictEval    bool      `json:"isStrictEval"`
      HasAuditHash    bool      `json:"hasAuditHash"`
      SignedAt        time.Time `json:"signedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/crypto/signer.go` before handoff to the MTA spooler.

---

### Improvement 903: Polymorphic Envelope Header Cipher & Nonce Mutator
- **ID**: `IMP-903`
- **Cluster**: 91
- **Goal**: Applies polymorphic AEAD ChaCha20-Poly1305 encryption to sensitive custom warming headers with per-message nonce mutation to evade deep packet pattern sniffers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss visual diff table showing plaintext custom headers versus dynamically mutated ciphertext and 96-bit nonce vectors in amber (`#FFB800`).
  - *Web*: `<header-cipher-table>` with interactive decrypt toggle for authorized seed mailboxes.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum cipher-headers --envelope-id=<uuid> [--nonce-entropy=<bits>] [--is-ephemeral] [--has-hmac-tag]`
  - Flag `--envelope-id`: Identifier of target warming envelope.
  - Flag `--nonce-entropy`: Random bit length for nonce generator (default: 96).
  - Flag `--is-ephemeral`: Ensures nonce and key are scrubbed from RAM immediately after encryption.
  - Flag `--has-hmac-tag`: Appends SHA-256 HMAC authentication tag to header payload.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CipheredHeaders (
      id TEXT PRIMARY KEY,
      envelopeUuid TEXT NOT NULL UNIQUE,
      nonceHex TEXT NOT NULL,
      headerPayloadCipherHex TEXT NOT NULL,
      hmacTagHex TEXT NOT NULL,
      isEphemeral INTEGER NOT NULL DEFAULT 1,
      hasHmacTag INTEGER NOT NULL DEFAULT 1,
      encryptedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CipheredHeader struct {
      EnvelopeUUID           string    `json:"envelopeUuid"`
      NonceHex               string    `json:"nonceHex"`
      HeaderPayloadCipherHex string    `json:"headerPayloadCipherHex"`
      HmacTagHex             string    `json:"hmacTagHex"`
      IsEphemeral            bool      `json:"isEphemeral"`
      HasHmacTag             bool      `json:"hasHmacTag"`
      EncryptedAt            time.Time `json:"encryptedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/crypto/poly_header.go`.

---

### Improvement 904: Ephemeral Session Key (ESK) Forward-Secrecy Double Ratchet
- **ID**: `IMP-904`
- **Cluster**: 91
- **Goal**: Derives per-message symmetric encryption keys using HKDF-SHA512 with ratchet advancement, guaranteeing forward secrecy across multi-turn conversation simulations.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss ratchet step counter showing root key entropy, ratchet iteration index, and key retirement alerts in electric blue (`#00B0FF`).
  - *Web*: `<double-ratchet-stepper>` SVG graph depicting root and chaining key derivations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum ratchet-step --session-id=<uuid> [--ratchet-rounds=<n>] [--is-forward-secure] [--has-zeroized-prev]`
  - Flag `--session-id`: Target warming session identifier.
  - Flag `--ratchet-rounds`: Number of ratchet advancements to simulate or advance.
  - Flag `--is-forward-secure`: Validates past keys cannot be derived from current ratchet state.
  - Flag `--has-zeroized-prev`: Confirms memory zeroization of retired ephemeral keys.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RatchetSessions (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL UNIQUE,
      currentStepIndex INTEGER NOT NULL DEFAULT 0,
      chainKeyHash TEXT NOT NULL,
      isForwardSecure INTEGER NOT NULL DEFAULT 1,
      hasZeroizedPrev INTEGER NOT NULL DEFAULT 1,
      lastRatchetAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RatchetSession struct {
      SessionUUID     string    `json:"sessionUuid"`
      CurrentStepIndex int      `json:"currentStepIndex"`
      ChainKeyHash    string    `json:"chainKeyHash"`
      IsForwardSecure bool      `json:"isForwardSecure"`
      HasZeroizedPrev bool      `json:"hasZeroizedPrev"`
      LastRatchetAt   time.Time `json:"lastRatchetAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/crypto/ratchet.go`.

---

### Improvement 905: Zero-Knowledge Proof (ZKP) Envelope Origin Attestation
- **ID**: `IMP-905`
- **Cluster**: 91
- **Goal**: Generates non-interactive zero-knowledge proofs (zk-SNARKs) verifying warming envelope legitimacy and DKIM domain ownership without disclosing sender private keys to intermediary MTAs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss compact attestation chip with green proof confirmation (`#00E676`) and proof generation latency in milliseconds.
  - *Web*: `<zkp-attestation-card>` showing cryptographic circuit constraint checks and verification status.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum zkp-attest --envelope-id=<uuid> [--proving-key=<path>] [--is-valid-proof] [--has-snark-trace]`
  - Flag `--envelope-id`: Envelope UUID to generate attestation proof for.
  - Flag `--proving-key`: Relative repository path to circuit proving key.
  - Flag `--is-valid-proof`: Boolean flag signaling successful circuit validation.
  - Flag `--has-snark-trace`: Writes detailed proof verification logs to diagnostic stdout.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ZkpAttestations (
      id TEXT PRIMARY KEY,
      envelopeUuid TEXT NOT NULL UNIQUE,
      proofHash TEXT NOT NULL,
      circuitIdentifier TEXT NOT NULL DEFAULT 'warmup_origin_v1',
      isValidProof INTEGER NOT NULL DEFAULT 1,
      hasSnarkTrace INTEGER NOT NULL DEFAULT 1,
      verifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ZkpAttestation struct {
      EnvelopeUUID      string    `json:"envelopeUuid"`
      ProofHash         string    `json:"proofHash"`
      CircuitIdentifier string    `json:"circuitIdentifier"`
      IsValidProof      bool      `json:"isValidProof"`
      HasSnarkTrace     bool      `json:"hasSnarkTrace"`
      VerifiedAt        time.Time `json:"verifiedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/crypto/zkp.go`.

---

### Improvement 906: Post-Quantum TLS 1.3 Client Hello Padding & ECH Obfuscator
- **ID**: `IMP-906`
- **Cluster**: 91
- **Goal**: Injects standardized random padding and Encrypted Client Hello (ECH) extensions into TLS 1.3 handshakes to obscure SNI and prevent ISP traffic analysis heuristics.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss packet inspector displaying raw Client Hello record size, padding variance (128–512 bytes), and ECH state in pastel orange (`#FF9100`).
  - *Web*: `<ech-padding-inspector>` rendering live byte distribution diagrams.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum ech-pad --target-mx=<host> [--min-pad=<n>] [--max-pad=<n>] [--is-ech-enabled] [--has-random-padding]`
  - Flag `--min-pad`: Minimum random padding bytes to append (default: 128).
  - Flag `--max-pad`: Maximum random padding bytes to append (default: 512).
  - Flag `--is-ech-enabled`: Activates Encrypted Client Hello extension.
  - Flag `--has-random-padding`: Enforces randomized non-zero padding patterns.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EchConfigs (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL UNIQUE,
      paddingMinBytes INTEGER NOT NULL DEFAULT 128,
      paddingMaxBytes INTEGER NOT NULL DEFAULT 512,
      isEchEnabled INTEGER NOT NULL DEFAULT 1,
      hasRandomPadding INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EchConfig struct {
      TargetMxHost     string    `json:"targetMxHost"`
      PaddingMinBytes  int       `json:"paddingMinBytes"`
      PaddingMaxBytes  int       `json:"paddingMaxBytes"`
      IsEchEnabled     bool      `json:"isEchEnabled"`
      HasRandomPadding bool      `json:"hasRandomPadding"`
      UpdatedAt        time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/crypto/tls_ech.go`.

---

### Improvement 907: Post-Quantum Hardware Security Module (HSM) PKCS#11 Vault Bridge
- **ID**: `IMP-907`
- **Cluster**: 91
- **Goal**: Interfaces GitMap's quantum signing engine with hardware security modules (HSM) and cloud KMS via PKCS#11 for zero-leak private key governance.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss HSM slot matrix displaying active cryptographic token slots, session handles, and hardware PIN states with high-contrast yellow tags (`#FFD600`).
  - *Web*: `<hsm-token-status>` Web Component showing real-time HSM connection latency.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum hsm-status [--module-path=<path>] [--slot-id=<n>] [--is-hsm-connected] [--has-hardware-backing]`
  - Flag `--module-path`: Relative path to PKCS#11 driver module.
  - Flag `--slot-id`: Hardware token slot index.
  - Flag `--is-hsm-connected`: Current connection status flag.
  - Flag `--has-hardware-backing`: Verifies non-exportable hardware key attribute.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HsmTokens (
      id TEXT PRIMARY KEY,
      slotId INTEGER NOT NULL UNIQUE,
      tokenLabel TEXT NOT NULL,
      isHsmConnected INTEGER NOT NULL DEFAULT 1,
      hasHardwareBacking INTEGER NOT NULL DEFAULT 1,
      lastPingMs REAL NOT NULL DEFAULT 0.0,
      checkedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HsmToken struct {
      SlotID             int       `json:"slotId"`
      TokenLabel         string    `json:"tokenLabel"`
      IsHsmConnected     bool      `json:"isHsmConnected"`
      HasHardwareBacking bool      `json:"hasHardwareBacking"`
      LastPingMs         float64   `json:"lastPingMs"`
      CheckedAt          time.Time `json:"checkedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/crypto/hsm.go`.

---

### Improvement 908: Quantum Cryptographic Health & Entropy Drain Sentinel
- **ID**: `IMP-908`
- **Cluster**: 91
- **Goal**: Monitors system entropy pool health (`/dev/urandom` / Windows CryptGenRandom) to prevent cryptographic degradation during high-throughput warming bursts.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss real-time entropy gauge rendered in glowing magenta (`#FF0055`) with low-entropy warning thresholds.
  - *Web*: `<entropy-sentinel-gauge>` rendering live SVG radial meter updating every 500ms.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum entropy-check [--min-entropy-bits=<n>] [--is-pool-healthy] [--has-entropy-reserve]`
  - Flag `--min-entropy-bits`: Minimum available entropy threshold (default: 2048).
  - Flag `--is-pool-healthy`: Boolean health state indicator.
  - Flag `--has-entropy-reserve`: Indicates secondary hardware entropy source availability.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EntropyMetrics (
      id TEXT PRIMARY KEY,
      availableBits INTEGER NOT NULL,
      drainRateBitsPerSec REAL NOT NULL,
      isPoolHealthy INTEGER NOT NULL DEFAULT 1,
      hasEntropyReserve INTEGER NOT NULL DEFAULT 1,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EntropyMetric struct {
      AvailableBits       int       `json:"availableBits"`
      DrainRateBitsPerSec float64   `json:"drainRateBitsPerSec"`
      IsPoolHealthy       bool      `json:"isPoolHealthy"`
      HasEntropyReserve   bool      `json:"hasEntropyReserve"`
      MeasuredAt          time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed as a background daemon in `pkg/warmup/crypto/entropy.go`.

---

### Improvement 909: Polymorphic Header Watermarking & Anti-Steganography Scrub
- **ID**: `IMP-909`
- **Cluster**: 91
- **Goal**: Embeds zero-width steganographic watermarks in MIME headers for sender authenticity verification while actively stripping unwanted third-party tracking headers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss forensic inspection panel highlighting detected zero-width glyphs in teal (`#00BFA5`) and purged tracker keys in bright red (`#FF1744`).
  - *Web*: `<watermark-scrubber-view>` showing before-and-after header sanitization tables.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum watermark-scrub --envelope-id=<uuid> [--is-watermark-valid] [--has-stripped-trackers]`
  - Flag `--envelope-id`: UUID of envelope to process.
  - Flag `--is-watermark-valid`: Confirms legitimate internal watermark presence.
  - Flag `--has-stripped-trackers`: True if unwanted tracking headers were removed.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WatermarkAudits (
      id TEXT PRIMARY KEY,
      envelopeUuid TEXT NOT NULL UNIQUE,
      watermarkFingerprint TEXT NOT NULL,
      strippedHeaderCount INTEGER NOT NULL DEFAULT 0,
      isWatermarkValid INTEGER NOT NULL DEFAULT 1,
      hasStrippedTrackers INTEGER NOT NULL DEFAULT 0,
      auditedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WatermarkAudit struct {
      EnvelopeUUID        string    `json:"envelopeUuid"`
      WatermarkFingerprint string   `json:"watermarkFingerprint"`
      StrippedHeaderCount int       `json:"strippedHeaderCount"`
      IsWatermarkValid    bool      `json:"isWatermarkValid"`
      HasStrippedTrackers bool      `json:"hasStrippedTrackers"`
      AuditedAt           time.Time `json:"auditedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/crypto/watermark.go`.

---

### Improvement 910: Autonomous Post-Quantum Migration Replay Simulator
- **ID**: `IMP-910`
- **Cluster**: 91
- **Goal**: Replays historical warming traffic batches through simulated quantum adversary decoders to verify zero-leak survivability and post-quantum migration readiness.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss simulation summary card displaying total envelopes replayed, quantum attack resilience rate (100%), and cipher break time estimations in gold (`#FFD700`).
  - *Web*: `<quantum-replay-dashboard>` rendering interactive attack scenario charts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup quantum sim-replay --batch-id=<uuid> [--adversary-model=<model>] [--is-simulation-passed] [--has-zero-leak]`
  - Flag `--batch-id`: Identifier of historical warming batch to simulate.
  - Flag `--adversary-model`: Simulated quantum model (default: `Shor_Universal`).
  - Flag `--is-simulation-passed`: Boolean pass/fail indicator.
  - Flag `--has-zero-leak`: Asserts zero plaintext payload exposure across all replayed envelopes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS QuantumSimulations (
      id TEXT PRIMARY KEY,
      batchUuid TEXT NOT NULL UNIQUE,
      adversaryModel TEXT NOT NULL DEFAULT 'Shor_Universal',
      resilienceScore REAL NOT NULL DEFAULT 1.0,
      isSimulationPassed INTEGER NOT NULL DEFAULT 1,
      hasZeroLeak INTEGER NOT NULL DEFAULT 1,
      executedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type QuantumSimulation struct {
      BatchUUID          string    `json:"batchUuid"`
      AdversaryModel     string    `json:"adversaryModel"`
      ResilienceScore    float64   `json:"resilienceScore"`
      IsSimulationPassed bool      `json:"isSimulationPassed"`
      HasZeroLeak        bool      `json:"hasZeroLeak"`
      ExecutedAt         time.Time `json:"executedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed via `pkg/warmup/crypto/sim_replay.go`.

---

# Cluster 92: High-Frequency Multi-MTA Zero-Loss Message Routing & Resilient Failover (911–920)

---

### Improvement 911: Autonomous Multi-MTA Daemon Registry & Health Prober
- **ID**: `IMP-911`
- **Cluster**: 92 (High-Frequency Multi-MTA Zero-Loss Message Routing & Resilient Failover)
- **Goal**: Discovers, registers, and continuously probes local and remote Postfix, Haraka, OpenSMTPD, and Exim instances, measuring queue latency and egress health.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss server matrix table displaying registered MTA nodes, daemon types, active queue counts, and health states in bright green (`#00E676`) or crimson (`#D50000`).
  - *Web*: `<mta-registry-grid>` Web Component rendering real-time node availability tiles.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta register --node-id=<id> --mta-type=<type> --socket=<addr> [--is-healthy] [--has-tls-support]`
  - Flag `--node-id`: Unique alphanumeric identifier for the MTA daemon.
  - Flag `--mta-type`: Supported engine (`postfix`, `haraka`, `opensmtpd`, `exim`).
  - Flag `--socket`: Unix domain socket path or IP:port address.
  - Flag `--is-healthy`: Initial operational status flag.
  - Flag `--has-tls-support`: Signals outbound STARTTLS capability.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RegisteredMtas (
      id TEXT PRIMARY KEY,
      nodeIdentifier TEXT NOT NULL UNIQUE,
      mtaType TEXT NOT NULL,
      socketAddress TEXT NOT NULL,
      isHealthy INTEGER NOT NULL DEFAULT 1,
      hasTlsSupport INTEGER NOT NULL DEFAULT 1,
      consecutiveFails INTEGER NOT NULL DEFAULT 0,
      lastProbedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RegisteredMta struct {
      NodeIdentifier   string    `json:"nodeIdentifier"`
      MtaType          string    `json:"mtaType"`
      SocketAddress    string    `json:"socketAddress"`
      IsHealthy        bool      `json:"isHealthy"`
      HasTlsSupport    bool      `json:"hasTlsSupport"`
      ConsecutiveFails int       `json:"consecutiveFails"`
      LastProbedAt     time.Time `json:"lastProbedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Registered in `pkg/warmup/mta/registry.go`.

---

### Improvement 912: Predictive Multi-MTA Egress Route Weight Optimizer
- **ID**: `IMP-912`
- **Cluster**: 92
- **Goal**: Dynamically computes routing weights across registered MTAs using historical ISP bounce ratios, TLS handshake speed, and active queue depths to maximize delivery reputation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss bar chart displaying active routing weight percentages across nodes (e.g. Postfix-01: 45%, Haraka-02: 35%, OpenSMTPD-01: 20%) in cyan (`#00F0FF`).
  - *Web*: `<route-weight-bars>` SVG chart with smooth CSS transitions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta calc-weights [--target-domain=<domain>] [--is-dynamic-rebalance] [--has-reputation-boost]`
  - Flag `--target-domain`: Recipient domain filter for weight optimization.
  - Flag `--is-dynamic-rebalance`: Automatically updates routing tables based on computed weights.
  - Flag `--has-reputation-boost`: Biases traffic toward nodes with 99%+ delivery success rates.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MtaWeights (
      id TEXT PRIMARY KEY,
      nodeIdentifier TEXT NOT NULL,
      targetDomain TEXT NOT NULL,
      weightPercentage INTEGER NOT NULL DEFAULT 50,
      isDynamicRebalance INTEGER NOT NULL DEFAULT 1,
      hasReputationBoost INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      UNIQUE(nodeIdentifier, targetDomain)
  );
  ```
  ```go
  type MtaWeight struct {
      NodeIdentifier     string    `json:"nodeIdentifier"`
      TargetDomain       string    `json:"targetDomain"`
      WeightPercentage   int       `json:"weightPercentage"`
      IsDynamicRebalance bool      `json:"isDynamicRebalance"`
      HasReputationBoost bool      `json:"hasReputationBoost"`
      UpdatedAt          time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Controlled in `pkg/warmup/mta/weight_optimizer.go`.

---

### Improvement 913: Microsecond Circuit Breaker & Spike Trip Sentinel
- **ID**: `IMP-913`
- **Cluster**: 92
- **Goal**: Evaluates sliding 500ms windows of SMTP response codes and immediately trips when deferrals (421/451) or TCP resets exceed 8%, isolating the affected MTA node.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-urgency Lipgloss banner alerting trip events in bold red (`#FF0055`) with instantaneous trip reason and cooldown timer.
  - *Web*: `<circuit-breaker-sentinel>` with animated pulse indicators and manual reset triggers.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta circuit-status --node-id=<id> [--is-circuit-tripped] [--has-cooldown-active]`
  - Flag `--node-id`: Target MTA node to inspect.
  - Flag `--is-circuit-tripped`: Indicates if breaker is currently in OPEN/TRIPPED state.
  - Flag `--has-cooldown-active`: True if breaker is in HALF-OPEN state undergoing cautious probation.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CircuitBreakers (
      id TEXT PRIMARY KEY,
      nodeIdentifier TEXT NOT NULL UNIQUE,
      isCircuitTripped INTEGER NOT NULL DEFAULT 0,
      hasCooldownActive INTEGER NOT NULL DEFAULT 0,
      failureRatePercent REAL NOT NULL DEFAULT 0.0,
      trippedAt DATETIME,
      cooldownUntil DATETIME
  );
  ```
  ```go
  type CircuitBreaker struct {
      NodeIdentifier    string     `json:"nodeIdentifier"`
      IsCircuitTripped  bool       `json:"isCircuitTripped"`
      HasCooldownActive bool       `json:"hasCooldownActive"`
      FailureRatePercent float64   `json:"failureRatePercent"`
      TrippedAt         *time.Time `json:"trippedAt"`
      CooldownUntil     *time.Time `json:"cooldownUntil"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforced in `pkg/warmup/mta/circuit_breaker.go`.

---

### Improvement 914: Zero-Drop Spool Failover & Re-Routing Engine
- **ID**: `IMP-914`
- **Cluster**: 92
- **Goal**: Atomically extracts in-flight warming spools from a tripped or degraded MTA and re-injects them into the healthiest secondary MTA without message loss or duplicate dispatches.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss transactional migration progress bar showing re-routed message IDs and target node reassignment in vivid amber (`#FFB800`).
  - *Web*: `<zero-drop-failover-feed>` Web Component displaying live re-routing event logs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta failover --from-node=<src> --to-node=<dst> [--is-atomic-swap] [--has-zero-loss]`
  - Flag `--from-node`: Source degraded MTA node ID.
  - Flag `--to-node`: Destination failover MTA node ID.
  - Flag `--is-atomic-swap`: Enforces SQLite single-transaction queue migration.
  - Flag `--has-zero-loss`: Verifies hash-matching spool counts before and after migration.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SpoolFailovers (
      id TEXT PRIMARY KEY,
      failoverUuid TEXT NOT NULL UNIQUE,
      sourceNodeIdentifier TEXT NOT NULL,
      destinationNodeIdentifier TEXT NOT NULL,
      migratedCount INTEGER NOT NULL DEFAULT 0,
      isAtomicSwap INTEGER NOT NULL DEFAULT 1,
      hasZeroLoss INTEGER NOT NULL DEFAULT 1,
      completedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SpoolFailover struct {
      FailoverUUID              string    `json:"failoverUuid"`
      SourceNodeIdentifier      string    `json:"sourceNodeIdentifier"`
      DestinationNodeIdentifier string    `json:"destinationNodeIdentifier"`
      MigratedCount             int       `json:"migratedCount"`
      IsAtomicSwap              bool      `json:"isAtomicSwap"`
      HasZeroLoss               bool      `json:"hasZeroLoss"`
      CompletedAt               time.Time `json:"completedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/mta/failover.go`.

---

### Improvement 915: Heterogeneous MTA Command Adapter (Postfix/OpenSMTPD/Haraka/Exim)
- **ID**: `IMP-915`
- **Cluster**: 92
- **Goal**: Provides a unified Go interface abstracting daemon-specific CLI commands (`postqueue`, `smtpctl`, `haraka-ctl`, `exim -bp`) into standard queue inspection and flush primitives.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss normalized command runner table displaying cross-daemon execution times and exit statuses in muted lavender (`#D1C4E9`).
  - *Web*: `<mta-adapter-tester>` Web Component allowing interactive command invocation.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta exec --node-id=<id> --action=<action> [--is-normalized] [--has-timeout-guard]`
  - Flag `--node-id`: Target registered MTA node.
  - Flag `--action`: Standard operation (`flush`, `hold`, `delete`, `summary`).
  - Flag `--is-normalized`: Converts daemon output into standard GitMap JSON envelope format.
  - Flag `--has-timeout-guard`: Terminates hung external commands after 3000ms.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MtaAdapterExecutions (
      id TEXT PRIMARY KEY,
      nodeIdentifier TEXT NOT NULL,
      actionName TEXT NOT NULL,
      executionDurationMs REAL NOT NULL DEFAULT 0.0,
      isNormalized INTEGER NOT NULL DEFAULT 1,
      hasTimeoutGuard INTEGER NOT NULL DEFAULT 1,
      executedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MtaAdapterExecution struct {
      NodeIdentifier      string    `json:"nodeIdentifier"`
      ActionName          string    `json:"actionName"`
      ExecutionDurationMs float64   `json:"executionDurationMs"`
      IsNormalized        bool      `json:"isNormalized"`
      HasTimeoutGuard     bool      `json:"hasTimeoutGuard"`
      ExecutedAt          time.Time `json:"executedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/mta/adapter.go`.

---

### Improvement 916: Zero-Copy Shared Memory Spool Pipeline
- **ID**: `IMP-916`
- **Cluster**: 92
- **Goal**: Facilitates high-throughput envelope transmission between GitMap and local MTA daemons via POSIX shared memory ring buffers, bypassing disk I/O bottlenecks.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss buffer utilization meter showing shared ring buffer memory usage (e.g. 64MB / 128MB) and throughput (msgs/sec) in electric green (`#76FF03`).
  - *Web*: `<shm-spool-meter>` SVG ring diagram updating in real-time via WebSocket.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta shm-status [--buffer-name=<name>] [--is-shm-active] [--has-zero-copy-verified]`
  - Flag `--buffer-name`: Named shared memory segment identifier (default: `gitmap_warmup_shm0`).
  - Flag `--is-shm-active`: True if shared memory mapping is active.
  - Flag `--has-zero-copy-verified`: Validates zero intermediate memory duplication.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ShmSpoolMetrics (
      id TEXT PRIMARY KEY,
      bufferName TEXT NOT NULL UNIQUE,
      allocatedBytes INTEGER NOT NULL,
      usedBytes INTEGER NOT NULL,
      isShmActive INTEGER NOT NULL DEFAULT 1,
      hasZeroCopyVerified INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ShmSpoolMetric struct {
      BufferName          string    `json:"bufferName"`
      AllocatedBytes      int64     `json:"allocatedBytes"`
      UsedBytes           int64     `json:"usedBytes"`
      IsShmActive         bool      `json:"isShmActive"`
      HasZeroCopyVerified bool      `json:"hasZeroCopyVerified"`
      UpdatedAt           time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/mta/shm_spool.go`.

---

### Improvement 917: Egress IP Pool Virtual Interface Balancer
- **ID**: `IMP-917`
- **Cluster**: 92
- **Goal**: Binds outbound warming connections across virtual network interfaces and dedicated IP blocks, rotating source IPs with warm reputation metrics per recipient domain.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss IP matrix displaying assigned CIDR pools, interface bindings (`eth0:1`, `eth0:2`), and domain-to-IP affinity rules in bright sky blue (`#40C4FF`).
  - *Web*: `<ip-pool-balancer>` component rendering IP rotation flowcharts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta ip-pool --interface=<iface> [--source-ip=<ip>] [--is-ip-warm] [--has-reverse-dns]`
  - Flag `--interface`: System network interface name.
  - Flag `--source-ip`: Egress IPv4/IPv6 address.
  - Flag `--is-ip-warm`: Reputation warming completion indicator.
  - Flag `--has-reverse-dns`: Verifies valid PTR record configuration matching forward HELO.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EgressIpPools (
      id TEXT PRIMARY KEY,
      sourceIpAddress TEXT NOT NULL UNIQUE,
      interfaceName TEXT NOT NULL,
      ptrRecordHost TEXT NOT NULL,
      isIpWarm INTEGER NOT NULL DEFAULT 1,
      hasReverseDns INTEGER NOT NULL DEFAULT 1,
      dailySendQuota INTEGER NOT NULL DEFAULT 500,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EgressIpPool struct {
      SourceIPAddress string    `json:"sourceIpAddress"`
      InterfaceName   string    `json:"interfaceName"`
      PtrRecordHost   string    `json:"ptrRecordHost"`
      IsIPWarm        bool      `json:"isIpWarm"`
      HasReverseDNS   bool      `json:"hasReverseDns"`
      DailySendQuota  int       `json:"dailySendQuota"`
      CreatedAt       time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/mta/ip_pool.go`.

---

### Improvement 918: Synchronous Cross-MTA Queue Flush Coordinator
- **ID**: `IMP-918`
- **Cluster**: 92
- **Goal**: Orchestrates coordinated, synchronized queue flushes across all active MTAs during scheduled warming time slots, preventing staggered batch anomalies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss countdown timer and synchronization radar showing synchronized trigger events across nodes in neon teal (`#64FFDA`).
  - *Web*: `<mta-flush-radar>` rendering circular sweep animation during flush cycles.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta sync-flush [--window-ms=<n>] [--is-synchronized] [--has-ack-quorum]`
  - Flag `--window-ms`: Maximum allowed divergence window in milliseconds across nodes (default: 50).
  - Flag `--is-synchronized`: Validates all nodes fired within window.
  - Flag `--has-ack-quorum`: Requires 100% of healthy nodes to acknowledge before executing.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SyncFlushEvents (
      id TEXT PRIMARY KEY,
      flushBatchUuid TEXT NOT NULL UNIQUE,
      participatingNodeCount INTEGER NOT NULL,
      divergenceMs REAL NOT NULL DEFAULT 0.0,
      isSynchronized INTEGER NOT NULL DEFAULT 1,
      hasAckQuorum INTEGER NOT NULL DEFAULT 1,
      triggeredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SyncFlushEvent struct {
      FlushBatchUUID         string    `json:"flushBatchUuid"`
      ParticipatingNodeCount int       `json:"participatingNodeCount"`
      DivergenceMs           float64   `json:"divergenceMs"`
      IsSynchronized         bool      `json:"isSynchronized"`
      HasAckQuorum           bool      `json:"hasAckQuorum"`
      TriggeredAt            time.Time `json:"triggeredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/mta/flush_coordinator.go`.

---

### Improvement 919: Autonomous Daemon Process Supervisor & Crash Reviver
- **ID**: `IMP-919`
- **Cluster**: 92
- **Goal**: Monitors MTA daemon child processes, detects abnormal exits or hung sockets, and automatically restarts service daemons within 250ms while preserving local spools.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss supervisor alert box displaying daemon PID, uptime counter, crash count, and recovery latency in neon amber (`#FFD600`).
  - *Web*: `<daemon-supervisor-card>` Web Component with live status pills and process graphs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta supervise --node-id=<id> [--is-process-alive] [--has-auto-restart]`
  - Flag `--node-id`: Target MTA daemon node identifier.
  - Flag `--is-process-alive`: Current PID liveness indicator.
  - Flag `--has-auto-restart`: Enables autonomous restart daemon loops.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DaemonSupervisors (
      id TEXT PRIMARY KEY,
      nodeIdentifier TEXT NOT NULL UNIQUE,
      processId INTEGER NOT NULL,
      restartCount INTEGER NOT NULL DEFAULT 0,
      isProcessAlive INTEGER NOT NULL DEFAULT 1,
      hasAutoRestart INTEGER NOT NULL DEFAULT 1,
      lastRestartAt DATETIME
  );
  ```
  ```go
  type DaemonSupervisor struct {
      NodeIdentifier string     `json:"nodeIdentifier"`
      ProcessID      int        `json:"processId"`
      RestartCount   int        `json:"restartCount"`
      IsProcessAlive bool       `json:"isProcessAlive"`
      HasAutoRestart bool       `json:"hasAutoRestart"`
      LastRestartAt  *time.Time `json:"lastRestartAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/mta/supervisor.go`.

---

### Improvement 920: Multi-MTA Forensic Egress Audit Ledger Sync
- **ID**: `IMP-920`
- **Cluster**: 92
- **Goal**: Ingests raw queue and delivery logs from all registered MTA nodes, parses RFC 3464 delivery status notifications (DSN), and writes cryptographic audit records into `repodb/audit.db`.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss stream viewer showing normalized forensic log events ingested across multi-MTA nodes with hash-chain indicators in purple (`#E040FB`).
  - *Web*: `<forensic-audit-log-stream>` component rendering searchable, filterable log tables.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mta audit-sync [--node-id=<id>] [--is-ledger-synced] [--has-tamper-proof]`
  - Flag `--node-id`: Specific MTA node filter (omit for all nodes).
  - Flag `--is-ledger-synced`: Verifies complete sync without pending uncommitted logs.
  - Flag `--has-tamper-proof`: Asserts SHA-256 hash chaining validity across imported records.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MtaEgressAudits (
      id TEXT PRIMARY KEY,
      recordHash TEXT NOT NULL UNIQUE,
      nodeIdentifier TEXT NOT NULL,
      envelopeUuid TEXT NOT NULL,
      recipientAddress TEXT NOT NULL,
      smtpResponseCode INTEGER NOT NULL,
      enhancedStatusCode TEXT NOT NULL,
      isLedgerSynced INTEGER NOT NULL DEFAULT 1,
      hasTamperProof INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MtaEgressAudit struct {
      RecordHash         string    `json:"recordHash"`
      NodeIdentifier     string    `json:"nodeIdentifier"`
      EnvelopeUUID       string    `json:"envelopeUuid"`
      RecipientAddress   string    `json:"recipientAddress"`
      SmtpResponseCode   int       `json:"smtpResponseCode"`
      EnhancedStatusCode string    `json:"enhancedStatusCode"`
      IsLedgerSynced     bool      `json:"isLedgerSynced"`
      HasTamperProof     bool      `json:"hasTamperProof"`
      RecordedAt         time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/mta/audit_sync.go`.

---

# Cluster 93: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (921–930)

---

### Improvement 921: Lipgloss Real-Time Waveform ASCII/Braille Streamer
- **ID**: `IMP-921`
- **Cluster**: 93 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Renders high-frequency (60 FPS) terminal waveforms using Unicode Braille patterns (`⣾`, `⣽`, `⣻`, `⢿`) displaying live SMTP packet latency and command jitter.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss real-time Braille waveform canvas rendered in fluorescent green (`#00E676`) with upper and lower jitter tolerance bounds.
  - *Web*: `<braille-waveform-canvas>` Web Component emulating terminal font typography and waveform animations with zero external libraries.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui waveform [--refresh-hz=<n>] [--is-realtime] [--has-jitter-bounds]`
  - Flag `--refresh-hz`: Refresh frequency in Hertz (default: 60).
  - Flag `--is-realtime`: Binds directly to streaming IPC socket without buffering.
  - Flag `--has-jitter-bounds`: Draws upper and lower dynamic threshold lines on the terminal graph.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WaveformSnapshots (
      id TEXT PRIMARY KEY,
      snapshotUuid TEXT NOT NULL UNIQUE,
      averageJitterMs REAL NOT NULL,
      peakJitterMs REAL NOT NULL,
      isRealtime INTEGER NOT NULL DEFAULT 1,
      hasJitterBounds INTEGER NOT NULL DEFAULT 1,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WaveformSnapshot struct {
      SnapshotUUID    string    `json:"snapshotUuid"`
      AverageJitterMs float64   `json:"averageJitterMs"`
      PeakJitterMs    float64   `json:"peakJitterMs"`
      IsRealtime      bool      `json:"isRealtime"`
      HasJitterBounds bool      `json:"hasJitterBounds"`
      CapturedAt      time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Rendered via Bubble Tea v2 in `pkg/warmup/tui/waveform.go`.

---

### Improvement 922: Sub-Millisecond SMTP Phase Latency Flamegraph
- **ID**: `IMP-922`
- **Cluster**: 93
- **Goal**: Generates hierarchical execution flamegraphs decomposing total message delivery time into DNS, TCP Connect, TLS Handshake, PQC KEM, MAIL FROM, RCPT TO, and DATA phases.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss colored horizontal stacked bars showing proportional phase durations in hex colors: DNS (`#80D8FF`), TCP (`#82B1FF`), TLS (`#B388FF`), PQC (`#EA80FC`), SMTP (`#FF8A80`).
  - *Web*: `<smtp-flamegraph-view>` pure SVG flamegraph with click-to-zoom phase inspection.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui flamegraph --envelope-id=<uuid> [--is-flamegraph-rendered] [--has-sub-millisecond-precision]`
  - Flag `--envelope-id`: Envelope UUID to generate phase flamegraph for.
  - Flag `--is-flamegraph-rendered`: True if terminal graphic rendered successfully.
  - Flag `--has-sub-millisecond-precision`: Employs microsecond timers (`time.Now()`) for phase metrics.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PhaseLatencies (
      id TEXT PRIMARY KEY,
      envelopeUuid TEXT NOT NULL UNIQUE,
      dnsDurationMs REAL NOT NULL,
      tcpDurationMs REAL NOT NULL,
      tlsDurationMs REAL NOT NULL,
      pqcDurationMs REAL NOT NULL,
      smtpDurationMs REAL NOT NULL,
      isFlamegraphRendered INTEGER NOT NULL DEFAULT 1,
      hasSubMillisecondPrecision INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PhaseLatency struct {
      EnvelopeUUID               string    `json:"envelopeUuid"`
      DnsDurationMs              float64   `json:"dnsDurationMs"`
      TcpDurationMs              float64   `json:"tcpDurationMs"`
      TlsDurationMs              float64   `json:"tlsDurationMs"`
      PqcDurationMs              float64   `json:"pqcDurationMs"`
      SmtpDurationMs             float64   `json:"smtpDurationMs"`
      IsFlamegraphRendered       bool      `json:"isFlamegraphRendered"`
      HasSubMillisecondPrecision bool      `json:"hasSubMillisecondPrecision"`
      RecordedAt                 time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tui/flamegraph.go`.

---

### Improvement 923: Zero-CDN Reactive SVG Telemetry Web Component (`<envelope-telemetry-hud>`)
- **ID**: `IMP-923`
- **Cluster**: 93
- **Goal**: Delivers an air-gapped, zero-external-dependency Custom Element rendering dynamic vector charts, gauge needles, and dispatch matrices using pure native JavaScript and SVG.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: CLI command emitting standalone HTML/SVG bundle for local browser viewing.
  - *Web*: `<envelope-telemetry-hud>` Custom Element encapsulated in Shadow DOM with dark theme palette (`#0A0E17`, `#00F0FF`, `#7000FF`).
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export-hud --output=<path> [--is-standalone] [--has-zero-cdn-enforced]`
  - Flag `--output`: Relative file destination for exported HTML/SVG file.
  - Flag `--is-standalone`: Embeds all styles and scripts inline without external links.
  - Flag `--has-zero-cdn-enforced`: Linter check failing build if remote URLs are detected.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HudExports (
      id TEXT PRIMARY KEY,
      exportPath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL,
      isStandalone INTEGER NOT NULL DEFAULT 1,
      hasZeroCdnEnforced INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HudExport struct {
      ExportPath         string    `json:"exportPath"`
      FileSizeBytes      int64     `json:"fileSizeBytes"`
      IsStandalone       bool      `json:"isStandalone"`
      HasZeroCdnEnforced bool      `json:"hasZeroCdnEnforced"`
      ExportedAt         time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/hud_component.go`.

---

### Improvement 924: Bubble Tea v2 Multi-Pane Interactive Navigation Hub
- **ID**: `IMP-924`
- **Cluster**: 93
- **Goal**: Implements an interactive terminal user interface with keyboard-driven tab switching (1: Overview, 2: Cryptography, 3: MTA Mesh, 4: Queues, 5: Audit Ledger) using Bubble Tea v2.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Multi-pane layout with responsive terminal sizing, active tab indicator in neon yellow (`#FFEA00`), and footer keybinding hints (`[Tab] Next [q] Quit [r] Rebalance`).
  - *Web*: Matching desktop navigation bar mirroring the TUI layout.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui dashboard [--default-tab=<tab>] [--is-interactive] [--has-mouse-support]`
  - Flag `--default-tab`: Initial tab to display (default: `overview`).
  - Flag `--is-interactive`: Enables keystroke capture and interactive filtering.
  - Flag `--has-mouse-support`: Enables terminal mouse click and scroll wheel events.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiSessions (
      id TEXT PRIMARY KEY,
      sessionUuid TEXT NOT NULL UNIQUE,
      activeTabName TEXT NOT NULL DEFAULT 'overview',
      isInteractive INTEGER NOT NULL DEFAULT 1,
      hasMouseSupport INTEGER NOT NULL DEFAULT 1,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiSession struct {
      SessionUUID     string    `json:"sessionUuid"`
      ActiveTabName   string    `json:"activeTabName"`
      IsInteractive   bool      `json:"isInteractive"`
      HasMouseSupport bool      `json:"hasMouseSupport"`
      StartedAt       time.Time `json:"startedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Main TUI loop in `pkg/warmup/tui/dashboard.go`.

---

### Improvement 925: TrueColor ANSI Color Token & Theme Dynamic Switcher
- **ID**: `IMP-925`
- **Cluster**: 93
- **Goal**: Provides dynamic terminal theme switching between `Cyberpunk`, `Monokai`, `Nord`, and `High-Contrast-A11y` palettes using Lipgloss TrueColor tokens with automatic 16-color ANSI fallbacks.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Color swatches previewing palette selection live in the console with contrast ratio compliance indicators.
  - *Web*: CSS variable mapping adhering to chosen theme tokens.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui theme set --name=<theme> [--is-truecolor] [--has-contrast-guard]`
  - Flag `--name`: Theme identifier (`cyberpunk`, `monokai`, `nord`, `a11y`).
  - Flag `--is-truecolor`: Detects 24-bit TrueColor terminal support (`$COLORTERM`).
  - Flag `--has-contrast-guard`: Validates minimum 4.5:1 WCAG contrast on all text elements.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TuiThemes (
      id TEXT PRIMARY KEY,
      themeName TEXT NOT NULL UNIQUE,
      isTruecolor INTEGER NOT NULL DEFAULT 1,
      hasContrastGuard INTEGER NOT NULL DEFAULT 1,
      selectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TuiTheme struct {
      ThemeName        string    `json:"themeName"`
      IsTruecolor      bool      `json:"isTruecolor"`
      HasContrastGuard bool      `json:"hasContrastGuard"`
      SelectedAt       time.Time `json:"selectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tui/theme.go`.

---

### Improvement 926: Real-Time Bounded Jitter Alert HUD Overlay
- **ID**: `IMP-926`
- **Cluster**: 93
- **Goal**: Displays non-blocking floating HUD alerts whenever inter-command jitter breaches the statistical 95th percentile, alerting operators to potential timing detection risks.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss floating overlay toast in bright amber (`#FF9100`) anchored to the top-right corner of the terminal window.
  - *Web*: `<jitter-alert-toast>` Web Component with smooth slide-in animation.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui jitter-alert [--threshold-ms=<n>] [--is-alert-active] [--has-sound-bell]`
  - Flag `--threshold-ms`: Jitter variance threshold triggering alerts (default: 85.0).
  - Flag `--is-alert-active`: True if jitter threshold is actively breached.
  - Flag `--has-sound-bell`: Emits terminal ASCII bell (`\a`) on critical violations.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterAlerts (
      id TEXT PRIMARY KEY,
      alertUuid TEXT NOT NULL UNIQUE,
      breachVarianceMs REAL NOT NULL,
      isAlertActive INTEGER NOT NULL DEFAULT 1,
      hasSoundBell INTEGER NOT NULL DEFAULT 0,
      triggeredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterAlert struct {
      AlertUUID        string    `json:"alertUuid"`
      BreachVarianceMs float64   `json:"breachVarianceMs"`
      IsAlertActive    bool      `json:"isAlertActive"`
      HasSoundBell     bool      `json:"hasSoundBell"`
      TriggeredAt      time.Time `json:"triggeredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/jitter_alert.go`.

---

### Improvement 927: Terminal Queue Depth Sparkline Historian
- **ID**: `IMP-927`
- **Cluster**: 93
- **Goal**: Generates inline ASCII sparklines (` ▂▃▄▅▆▇█`) tracking queue depth progression across all MTA spools over the preceding 60 minutes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Compact sparkline trend rows embedded in MTA status cards with minimum, maximum, and current metrics in cyan (`#00E5FF`).
  - *Web*: `<sparkline-historian>` inline SVG path rendering sparklines without external canvas libraries.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui sparklines [--window-minutes=<n>] [--is-sparkline-rendered] [--has-peak-marker]`
  - Flag `--window-minutes`: History duration to plot (default: 60).
  - Flag `--is-sparkline-rendered`: Render confirmation indicator.
  - Flag `--has-peak-marker`: Emphasizes historical peak point with a contrasting glyph (`▲`).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS QueueSparklines (
      id TEXT PRIMARY KEY,
      nodeIdentifier TEXT NOT NULL,
      sparklineGlyphs TEXT NOT NULL,
      minDepth INTEGER NOT NULL,
      maxDepth INTEGER NOT NULL,
      isSparklineRendered INTEGER NOT NULL DEFAULT 1,
      hasPeakMarker INTEGER NOT NULL DEFAULT 1,
      computedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type QueueSparkline struct {
      NodeIdentifier      string    `json:"nodeIdentifier"`
      SparklineGlyphs     string    `json:"sparklineGlyphs"`
      MinDepth            int       `json:"minDepth"`
      MaxDepth            int       `json:"maxDepth"`
      IsSparklineRendered bool      `json:"isSparklineRendered"`
      HasPeakMarker       bool      `json:"hasPeakMarker"`
      ComputedAt          time.Time `json:"computedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/tui/sparkline.go`.

---

### Improvement 928: Interactive Forensic Transaction Drill-Down Inspector
- **ID**: `IMP-928`
- **Cluster**: 93
- **Goal**: Allows operators to select any in-flight or delivered warming envelope in the TUI and expand its full cryptographic handshake, RFC 5322 headers, and SMTP conversation transcript.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-view drawer displaying raw SMTP dialogue (lines starting with `>>>` in green, `<<<` in blue) and cryptographic verification badges.
  - *Web*: `<transaction-drilldown-modal>` Web Component with searchable transcript pane.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui inspect --envelope-id=<uuid> [--is-transcript-complete] [--has-crypto-verified]`
  - Flag `--envelope-id`: Envelope UUID to drill down into.
  - Flag `--is-transcript-complete`: Confirms full captured conversation transcript.
  - Flag `--has-crypto-verified`: Displays verified cryptographic validation badges.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TransactionInspects (
      id TEXT PRIMARY KEY,
      envelopeUuid TEXT NOT NULL UNIQUE,
      rawTranscriptText TEXT NOT NULL,
      isTranscriptComplete INTEGER NOT NULL DEFAULT 1,
      hasCryptoVerified INTEGER NOT NULL DEFAULT 1,
      inspectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TransactionInspect struct {
      EnvelopeUUID         string    `json:"envelopeUuid"`
      RawTranscriptText    string    `json:"rawTranscriptText"`
      IsTranscriptComplete bool      `json:"isTranscriptComplete"`
      HasCryptoVerified    bool      `json:"hasCryptoVerified"`
      InspectedAt          time.Time `json:"inspectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/tui/inspect.go`.

---

### Improvement 929: Real-Time Headless Telemetry Streaming Web Socket Server
- **ID**: `IMP-929`
- **Cluster**: 93
- **Goal**: Hosts an ultra-lightweight, zero-external-dependency HTTP/WebSocket server in GitMap streaming real-time telemetry events to local browser dashboards or remote monitoring agents.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss connection card showing listening address (e.g. `127.0.0.1:8765`), connected client count, and outbound message rate in soft purple (`#CE93D8`).
  - *Web*: Native WebSocket connection lifecycle indicators inside `<envelope-telemetry-hud>`.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui ws-stream --port=<port> [--is-server-listening] [--has-auth-token]`
  - Flag `--port`: Local TCP port for WebSocket listener (default: 8765).
  - Flag `--is-server-listening`: True if socket is actively bound and accepting connections.
  - Flag `--has-auth-token`: Enforces pre-shared Bearer token authentication on incoming connections.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WsServerConfigs (
      id TEXT PRIMARY KEY,
      listenPort INTEGER NOT NULL UNIQUE,
      connectedClients INTEGER NOT NULL DEFAULT 0,
      isServerListening INTEGER NOT NULL DEFAULT 1,
      hasAuthToken INTEGER NOT NULL DEFAULT 1,
      startedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type WsServerConfig struct {
      ListenPort        int       `json:"listenPort"`
      ConnectedClients  int       `json:"connectedClients"`
      IsServerListening bool      `json:"isServerListening"`
      HasAuthToken      bool      `json:"hasAuthToken"`
      StartedAt         time.Time `json:"startedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/tui/ws_server.go`.

---

### Improvement 930: Multi-Format Telemetry Snapshot Exporter (JSON/CSV/SVG)
- **ID**: `IMP-930`
- **Cluster**: 93
- **Goal**: Exports frozen snapshots of current warming telemetry, jitter statistics, and multi-MTA queue distribution into machine-readable JSON, CSV, or standalone vector SVG files.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss export summary panel listing generated artifact paths, file sizes, and cryptographic SHA-256 digests in bright cyan (`#18FFFF`).
  - *Web*: `<snapshot-export-dialog>` Web Component with one-click download buttons.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tui export --format=<fmt> --out=<path> [--is-snapshot-valid] [--has-hash-digest]`
  - Flag `--format`: Output format (`json`, `csv`, `svg`).
  - Flag `--out`: Relative repository destination path.
  - Flag `--is-snapshot-valid`: Confirms snapshot data integrity.
  - Flag `--has-hash-digest`: Appends `.sha256` verification file alongside export.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TelemetryExports (
      id TEXT PRIMARY KEY,
      exportUuid TEXT NOT NULL UNIQUE,
      outputFormat TEXT NOT NULL,
      destinationPath TEXT NOT NULL,
      sha256DigestHex TEXT NOT NULL,
      isSnapshotValid INTEGER NOT NULL DEFAULT 1,
      hasHashDigest INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TelemetryExport struct {
      ExportUUID      string    `json:"exportUuid"`
      OutputFormat    string    `json:"outputFormat"`
      DestinationPath string    `json:"destinationPath"`
      Sha256DigestHex string    `json:"sha256DigestHex"`
      IsSnapshotValid bool      `json:"isSnapshotValid"`
      HasHashDigest   bool      `json:"hasHashDigest"`
      ExportedAt      time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/tui/exporter.go`.
