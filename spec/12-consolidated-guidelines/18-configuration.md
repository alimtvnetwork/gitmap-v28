# 18 — Configuration Management

Environment variables, feature flags, config validation, and secrets rotation.

## Environment Variables

`UPPER_SNAKE_CASE`. Prefix with service name. Loading order: process env → `.env` file → hardcoded defaults. Process-level always wins. Defaults must be safe (restrictive).

Never log variable values. Validate at startup. Type-convert explicitly.

## Feature Flags

Decouple deployment from release. Types: release (days), ops (weeks), experiment (A/B), permission (permanent). Default to off. Name with `ENABLE_` prefix. Clean up after rollout. Test both paths.

## Config Validation

Validate at startup — fail fast on invalid config. Check types, ranges, enums, and paths. Log sanitized config (no secrets) for debugging.

## Three-Layer Merge

```
Hardcoded defaults → Config file → CLI flags / env vars
```

CLI flags have highest priority.

## Secrets Rotation

Support dual keys during rotation (`API_KEY` + `API_KEY_PREVIOUS`). Log rotation key usage. Never cache secrets indefinitely.

## Config File Standards

JSON format. camelCase field names. Arrays default to `[]`. Strings default to `""`. No `null` values.

## Constraints

No secrets in config files. No `null` values. Validate before use. Feature flags default to off. Config files are version-controlled (but `.env` with real values is not).

---

Source: `spec/05-coding-guidelines/21-configuration-management.md`
