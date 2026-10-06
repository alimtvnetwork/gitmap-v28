# 04 — Fleet Nodes, SSH Delegation & RSA Credential Vault

- **Domain:** Distributed Nodes, Remote SSH Execution & Credential Management
- **Authoritative Specification:** [04-fleet-nodes-and-ssh](../../02-spec/21-app/04-fleet-nodes-and-ssh/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. Fleet Nodes Topology & Node Membership

GitMap manages distributed developer nodes and remote servers through a unified node registry:

- **Node Types:** Workstations (Ubuntu, macOS, Windows), VM runners, and remote deployment servers.
- **Node Join Protocol (`gitmap nodes join`):** Adds a new machine into the local cluster database (`SshHost` table), recording hostname, IP addresses, SSH user, port, and key fingerprint.
- **Node Discovery & Recall:** Subcommands support host lookup by friendly alias or IP prefix; autocompletion enumerates registered nodes.
- **Except-Self Filter:** When executing cluster broadcast commands or remote cloning, GitMap automatically filters out the local executing host (`except-self`) to prevent duplicate operations or loops.

---

## 2. Port 22 Concurrent TCP Liveness Probing

Before executing remote commands across nodes, GitMap conducts pre-flight network checks:

- **Dual-Stack Probe:** Probes both ICMP echo and TCP port 22 handshake concurrently.
- **Dial Timeout Cap:** Enforces a strict 1500ms TCP dial timeout per endpoint.
- **Fail-Fast Status:** Unreachable or slow nodes are marked offline in the status table without blocking subsequent execution on healthy cluster nodes.
- **Terminal UI Badging:** Displays latency in milliseconds alongside connection health (`[HEALTHY] 12ms`, `[TIMEOUT] >1500ms`).

---

## 3. RSA-OAEP Credential Vault & Password Capture

- **Encrypted Local Vault:** Credentials (SSH passphrases, sudo tokens, API secrets) are secured in `credentials.vault` using 2048-bit RSA-OAEP encryption combined with AES-256-GCM authenticated payloads.
- **Zero Plaintext Storage:** Passwords entered in CLI prompts are masked during input (`******`) and wiped from memory buffers immediately after encryption.
- **Token Vault Storage:** Public keys are managed in `~/.ssh/id_rsa.pub` or custom key paths, while private keys remain locked behind OS keychain or vault encryption.
- **SSH Config Sync:** Safely injects and maintains SSH host stanzas in `~/.ssh/config` without overwriting existing user-defined hosts.

---

## 4. Multi-Target Remote Execution Engine (`gitmap se`)

- **Space-Delimited Target Parsing:** Commands support multiple target nodes (e.g., `gitmap se "worker1 worker2 vm3" "uptime"`).
- **Parallel Output Multiplexing:** Stdout and stderr from remote nodes are multiplexed to the terminal with node-prefixed labels (`[worker1]`, `[worker2]`), preserving execution order.
- **Remote Clone Engine:** Dispatches repository clone directives to remote nodes, verifying disk space and Git credentials on the remote end before invoking clone.
