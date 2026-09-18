# Subtask 07: Network IP Manager Factory Decoupling & Hermetic Testing

## Objective
Decouple `cli/cmdos/os_ip.go` from real network interface manipulation by introducing an injectable `NetIPManagerFactory`.

## Assigned Files
- `cli/cmdos/os_ip.go`
- `cli/cmdos/os_ip_test.go`

## Implementation Steps
1. In `cli/cmdos/os_ip.go`:
   - Introduce `type NetIPManagerFactory func() *netip.Manager`.
   - Define `var defaultNetIPManagerFactory NetIPManagerFactory = createNetIPManager`.
   - In `runOSIP`: call `defaultNetIPManagerFactory()`.
2. In `cli/cmdos/os_ip_test.go`:
   - Provide `setupMockIPManager()` with `mockIPDriver` and `defer` restoration.
   - Add unit tests verifying `runOSIP` help, invalid subcommand, and mocked show.
