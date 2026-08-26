# Rust Testing Standards

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Testing conventions for Rust projects. Follows the cross-language [Test Naming and Structure](../01-cross-language/14-test-naming-and-structure.md) guidelines adapted for Rust's built-in test framework.

---

## Test Naming Convention

Pattern: `test_{function}_{scenario}_{expected_result}`

```rust
#[cfg(test)]
mod tests {
    use super::*;

    // ✅ Correct — descriptive three-part name
    #[test]
    fn test_parse_browser_title_chrome_returns_page_title() {
        let title = "GitHub - Google Chrome";
        let result = parse_browser_title(title, BrowserType::Chrome);
        assert_eq!(result.unwrap().title, "GitHub");
    }

    #[test]
    fn test_parse_browser_title_empty_string_returns_none() {
        let result = parse_browser_title("", BrowserType::Chrome);
        assert!(result.is_none());
    }

    #[test]
    fn test_dwell_tracker_short_visit_discards_event() {
        let mut tracker = DwellTracker::new(3.0);
        // Visit for < 3 seconds
        let completed = tracker.on_tab_change(tab_info("Site A"));
        assert!(completed.is_none()); // First visit, nothing to complete
    }
}
```

---

## Test Organization

### Unit Tests — Same file, `#[cfg(test)]` module

```rust
// src/collectors/browser.rs

pub fn classify_url(url: &str, categories: &[UrlCategory]) -> &str {
    // ...
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_classify_url_github_returns_work() {
        let categories = default_categories();
        assert_eq!(classify_url("https://github.com/repo", &categories), "Work");
    }
}
```

### Integration Tests — `tests/` directory

```
tests/
├── daemon_lifecycle_test.rs    # Start/stop/pause/resume
├── storage_test.rs             # SQLite read/write round-trips
├── browser_collector_test.rs   # Browser detection integration
└── api_test.rs                 # HTTP API endpoint tests
```

```rust
// tests/storage_test.rs
use timelog::storage::SqliteStorage;
use tempfile::tempdir;

#[tokio::test]
async fn test_storage_write_and_read_round_trip() {
    let directory = tempdir().unwrap();
    let storage = SqliteStorage::open(directory.path().join("test.db"))
        .expect("should open test database");

    let session = create_test_session();
    storage.insert_session(&session).expect("should insert session");

    let retrieved = storage.get_session(&session.id)
        .expect("should retrieve session");
    assert_eq!(retrieved.id, session.id);
}
```

---

## Assertion Patterns

```rust
// ✅ Use assert_eq! with descriptive messages
assert_eq!(result.dwell_seconds, 45.0, "Dwell time should be 45 seconds");

// ✅ Use assert! for boolean conditions
assert!(filter.should_track("https://github.com", false), "GitHub should be tracked");

// ✅ Use assert_matches! for enum variants (nightly or matches! macro)
assert!(matches!(error, CollectorError::StartFailed { .. }));

// ✅ Test error cases explicitly
let result = Config::load(Path::new("/nonexistent"));
assert!(result.is_err());
assert!(matches!(result.unwrap_err(), ConfigError::FileNotFound(_)));
```

---

## Test Helpers

```rust
#[cfg(test)]
mod test_helpers {
    use super::*;

    pub fn create_test_config() -> Config {
        Config {
            general: GeneralConfig {
                data_directory: tempdir().unwrap().into_path(),
                log_level: "Debug".into(),
                auto_start: false,
            },
            // ... minimal config for tests
        }
    }

    pub fn create_test_session() -> Session {
        Session {
            id: Uuid::new_v4(),
            started_at: Utc::now(),
            ended_at: None,
            duration_seconds: None,
            end_reason: "Active".into(),
        }
    }
}
```

---

## Async Test Patterns

```rust
// ✅ Use #[tokio::test] for async tests
#[tokio::test]
async fn test_event_bus_sends_and_receives() {
    let (sender, mut receiver) = create_event_bus();

    sender.send(ActivityEvent::SessionStart {
        session_id: Uuid::new_v4(),
        timestamp: Utc::now(),
    }).await.unwrap();

    let event = receiver.recv().await.expect("should receive event");
    assert!(matches!(event, ActivityEvent::SessionStart { .. }));
}

// ✅ Use tokio::time::pause for time-dependent tests
#[tokio::test]
async fn test_idle_detection_triggers_after_threshold() {
    tokio::time::pause();
    let mut idle_detector = IdleDetector::new(Duration::from_secs(300));

    tokio::time::advance(Duration::from_secs(301)).await;

    assert!(idle_detector.is_idle());
}
```

---

## Mocking with Traits

```rust
// Define trait for mockable dependencies
#[async_trait]
pub trait ScreenCapture: Send + Sync {
    async fn capture(&self, mode: CaptureMode) -> Result<RawBuffer, CaptureError>;
}

// Production implementation
pub struct OsScreenCapture;

#[async_trait]
impl ScreenCapture for OsScreenCapture {
    async fn capture(&self, mode: CaptureMode) -> Result<RawBuffer, CaptureError> {
        platform::capture_screen(mode)
    }
}

// Test mock
#[cfg(test)]
pub struct MockScreenCapture {
    pub should_fail: bool,
}

#[cfg(test)]
#[async_trait]
impl ScreenCapture for MockScreenCapture {
    async fn capture(&self, _mode: CaptureMode) -> Result<RawBuffer, CaptureError> {
        if self.should_fail {
            Err(CaptureError::ScreenCaptureFailed("mock error".into()))
        } else {
            Ok(RawBuffer::test_buffer(100, 100))
        }
    }
}
```

---

## Coverage Targets

| Category | Target |
|----------|--------|
| Domain logic (parsing, classification, aggregation) | ≥ 90% |
| Storage layer (SQLite operations) | ≥ 80% |
| API handlers | ≥ 80% |
| Platform-specific code (FFI wrappers) | Manual testing only |
| Configuration loading | ≥ 90% |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Test Standards | `../01-cross-language/14-test-naming-and-structure.md` |
| Error Handling (testing error paths) | `./02-error-handling.md` |
