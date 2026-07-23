# 16 — CI/CD Patterns

Pipeline structure, caching, deployment gates, and artifact management.

## Stage Ordering

Lint → Test → Build → Package → Deploy → Verify. Each must pass before the next.

## Pipeline-as-Code

Define in version-controlled YAML. Pin action versions to exact tags — never `@latest` or `@main`.

## Build Caching

Cache dependencies by lock file hash. Prefix keys with OS and tool version. Never cache secrets.

## CI Tool Versions

All pinned: `golangci-lint@v1.64.8`, `govulncheck@v1.1.4`. `@latest` is prohibited.

## Deployment Gates

dev → staging → production. Production requires approval for critical services.

## Build Once, Package Once

Binaries compiled exactly once per pipeline. All downstream steps operate on built artifacts — never rebuild.

## SHA Deduplication

Cache-based dedup with commit SHA key. Use passthrough gate pattern (step-level conditionals, not job-level `if`). Inline cache write in last validation job.

## Artifact Naming

`<binary>-<os>-<arch>.<ext>`. SHA-256 checksums for all release artifacts.

## Lessons Learned

- Never use `cd` in CI — use `working-directory`
- Validate build output directories before operating
- Never cancel release branch runs

---

Source: `spec/05-coding-guidelines/17-cicd-patterns.md`, `spec/05-coding-guidelines/29-ci-sha-deduplication.md`
