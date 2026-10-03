# Rust Error Handling (AI Execution Prompt)

> **/goal** Enforce structured, robust Rust error handling using the dual `thiserror` and `anyhow` pattern with numerical error codes and wire-format serialization.
> **/learn** Master the error handling taxonomy: `thiserror` for strongly-typed domain errors with project error codes, `anyhow` with context for application workflows, and serialized PascalCase `ErrorResponse` envelopes.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Define domain and library errors using `thiserror::Error` with strongly typed variants and numerical error codes.
- [ ] `/learn` Propagate application and CLI workflow errors using `anyhow::Result` with `.context()` annotations.
- [ ] `/goal` Map domain errors to serialized `ErrorResponse` and `ErrorDetail` structs using `#[serde(rename_all = "PascalCase")]`.
- [ ] `/learn` Prohibit bare `.unwrap()` or `.expect()` in production code paths; handle all `Result` and `Option` types explicitly.
- [ ] `/learn` Verify zero guideline violations via `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Error handling patterns for Rust projects. Uses the `thiserror` + `anyhow` dual-crate approach: `thiserror` for library/domain errors, `anyhow` for application-level error propagation.

---

## Error Type Hierarchy

```
AppError (thiserror)
├── ConfigError
│   ├── FileNotFound
│   ├── ParseError
│   └── InvalidValue
├── DatabaseError
│   ├── OpenFailed
│   ├── MigrationFailed
│   ├── WriteFailed
│   └── QueryFailed
├── CollectorError
│   ├── StartFailed
│   ├── StopFailed
│   └── HookRegistrationFailed
├── OsError
│   ├── PermissionDenied
│   ├── ApiUnavailable
│   └── WindowInfoFailed
├── CaptureError
│   ├── ScreenCaptureFailed
│   ├── EncodeFailed
│   └── WriteFailed
└── ApiError
    ├── BindFailed
    ├── Unauthorized
    └── InvalidParameter
```

---

## Defining Domain Errors with `thiserror`

```rust
use thiserror::Error;

// ✅ REQUIRED — each error variant maps to an error code
#[derive(Debug, Error)]
pub enum CollectorError {
    #[error("Collector '{name}' failed to start: {source}")]
    StartFailed {
        name: String,
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },

    #[error("Collector '{name}' failed to stop within timeout")]
    StopTimeout { name: String },

    #[error("OS hook registration failed: {0}")]
    HookRegistrationFailed(#[from] OsError),
}

impl CollectorError {
    /// Map to project error code (15000-range)
    pub fn error_code(&self) -> u32 {
        match self {
            Self::StartFailed { .. } => 15006,
            Self::StopTimeout { .. } => 15007,
            Self::HookRegistrationFailed(_) => 15150,
        }
    }
}
```

---

## Error Propagation Rules

### Rule 1: Use `Result<T, E>` everywhere — never panic in production code

```rust
// ✅ REQUIRED
pub fn load_config(path: &Path) -> Result<Config, ConfigError> {
    let content = std::fs::read_to_string(path)
        .map_err(|_| ConfigError::FileNotFound(path.to_path_buf()))?;
    toml::from_str(&content)
        .map_err(|error| ConfigError::ParseError { path: path.to_path_buf(), error })
}

// ❌ FORBIDDEN — unwrap in non-test code
pub fn load_config(path: &Path) -> Config {
    let content = std::fs::read_to_string(path).unwrap();
    toml::from_str(&content).unwrap()
}
```

### Rule 2: Use `anyhow` at application boundaries only

```rust
// ✅ REQUIRED — main.rs uses anyhow for top-level error handling
use anyhow::{Context, Result};

#[tokio::main]
async fn main() -> Result<()> {
    let config = Config::load(&config_path)
        .context("Failed to load configuration")?;

    let daemon = Daemon::new(config)
        .context("Failed to initialize daemon")?;

    daemon.run().await
        .context("Daemon exited with error")?;

    Ok(())
}

// ❌ FORBIDDEN — anyhow in library/domain code
pub fn start_collector(&mut self) -> anyhow::Result<()> {
    // Domain code should use specific error types
}
```

### Rule 3: Always add context when converting errors

```rust
// ✅ REQUIRED — context explains what operation failed
let database = SqliteStorage::open(&database_path)
    .map_err(|error| DatabaseError::OpenFailed {
        path: database_path.clone(),
        source: error,
    })?;

// ❌ FORBIDDEN — bare ? without context
let database = SqliteStorage::open(&database_path)?;
```

---

## Panic Policy

| Context | Panic Allowed? | Alternative |
|---------|---------------|-------------|
| Production code | ❌ Never | Return `Result` |
| Tests | ✅ Yes | `unwrap()`, `expect()`, `assert!()` |
| Const initialization | ✅ Yes | Compile-time guarantee |
| Unreachable branches | ✅ With `unreachable!()` | Only when logically provable |

### `expect()` vs `unwrap()`

```rust
// ✅ REQUIRED in tests — expect() with descriptive message
let config = Config::load(&path).expect("test config should be valid");

// ❌ FORBIDDEN in tests — bare unwrap gives no context
let config = Config::load(&path).unwrap();
```

---

## Error Logging

```rust
use tracing::{error, warn, info};

match collector.start(sender.clone()).await {
    Ok(()) => info!(collector = collector.name(), "Collector started"),
    Err(error) => {
        error!(
            collector = collector.name(),
            error_code = error.error_code(),
            error = %error,
            "Collector failed to start"
        );
        // Decide: skip this collector or abort daemon?
        if collector.is_required() {
            return Err(error.into());
        }
    }
}
```

---

## Error Response Serialization

For the HTTP API, errors serialize to the project-standard envelope:

```rust
#[derive(Debug, Serialize)]
#[serde(rename_all = "PascalCase")]
pub struct ErrorResponse {
    pub error: ErrorDetail,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "PascalCase")]
pub struct ErrorDetail {
    pub code: u32,
    pub message: String,
    pub detail: Option<String>,
}

impl From<&AppError> for ErrorResponse {
    fn from(error: &AppError) -> Self {
        Self {
            error: ErrorDetail {
                code: error.error_code(),
                message: error.to_string(),
                detail: error.source().map(|source| source.to_string()),
            },
        }
    }
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Error Code Registry | `../../03-error-manage/03-error-code-registry/01-registry.md` |
| Error Resolution Spec | `../../03-error-manage/readme.md` |
| Cross-Language Guidelines | `../01-cross-language/readme.md` |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/05-rust/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-RUST-003: Rust Result Handling, Error Trait and No Panics

**Given** Rust source code across library crates, domain modules, and application binaries.
**When** Audited against this specification using guideline linters and static checks.
**Then** All domain errors implement `std::error::Error` (via `thiserror`), application orchestration uses `anyhow::Result` with context, production code avoids `.unwrap()` and `.expect()` panics, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.
