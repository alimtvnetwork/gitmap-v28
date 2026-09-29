# Subtask 193.1: Typed JSON Envelope Architecture & Core Package (`cli/jsonenvelope`)
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/183-typed-json-envelope-and-format-inspection.md](../../../02-spec/21-app/183-typed-json-envelope-and-format-inspection.md)
Target Files: cli/jsonenvelope/envelope.go, cli/jsonenvelope/registry.go, cli/jsonenvelope/envelope_test.go
Action: Implement Envelope[T] and RawEnvelope structs with attributes (type, source, how, version, timestamp) and data. Add type detection, envelope wrapping, and transparent unmarshaling for both envelope and legacy JSONs.
Acceptance Criteria:
1. `Envelope` struct serializes to top-level `attributes` and `data` keys.
2. `DetectType` accurately detects typed envelope types as well as inferred legacy schemas.
3. `ExtractPayload` extracts raw data whether wrapped in an envelope or raw legacy JSON.
Targeted Verification: `go test -v ./jsonenvelope/...`
