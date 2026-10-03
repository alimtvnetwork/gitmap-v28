# Rust Async Patterns (AI Execution Prompt)

> **/goal** Standardize async architecture, Tokio runtime execution, bounded channel communication, and cancellation-safe async patterns across Rust codebases.
> **/learn** Master async concurrency principles: multi-threaded Tokio runtime configuration, bounded MPSC and broadcast channels, cooperative cancellation tokens, and avoiding lost writes in `select!` branches.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Configure multi-threaded Tokio runtimes with bounded worker thread allocations.
- [ ] `/learn` Mandate bounded channel buffers (`mpsc::channel(capacity)`) and non-blocking `try_send` for high-throughput producers.
- [ ] `/goal` Implement graceful shutdown handlers using `broadcast` or `watch` channels to signal cooperative cancellation across tasks.
- [ ] `/learn` Ensure `tokio::select!` branches are strictly cancellation-safe, ensuring critical I/O operations finish before dropping futures.
- [ ] `/goal` Enforce `Send + Sync + 'static` trait bounds on all futures and shared state across spawned async tasks.
- [ ] `/learn` Verify zero guideline violations via `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

Async conventions for Tokio-based Rust applications. Covers task spawning, channel patterns, graceful shutdown, and cancellation safety.

---

## Runtime Configuration

```rust
// ✅ Correct — multi-threaded runtime with controlled thread count
#[tokio::main(flavor = "multi_thread", worker_threads = 2)]
async fn main() -> anyhow::Result<()> {
    // ...
}
```

- **Worker threads:** 2 (sufficient for I/O-bound activity tracking)
- **Blocking pool:** Default (for SQLite writes and screenshot encoding)

---

## Channel Patterns

### Event Bus (MPSC)

```rust
use tokio::sync::mpsc;

// ✅ Bounded channel with explicit capacity
const EVENT_BUFFER_SIZE: usize = 1_000;

pub type EventSender = mpsc::Sender<ActivityEvent>;
pub type EventReceiver = mpsc::Receiver<ActivityEvent>;

pub fn create_event_bus() -> (EventSender, EventReceiver) {
    mpsc::channel(EVENT_BUFFER_SIZE)
}
```

### Sending Events (Non-Blocking)

```rust
// ✅ REQUIRED — try_send to avoid blocking collectors
match sender.try_send(event) {
    Ok(()) => {},
    Err(mpsc::error::TrySendError::Full(_)) => {
        warn!(error_code = 15005, "Event bus full, dropping event");
    }
    Err(mpsc::error::TrySendError::Closed(_)) => {
        info!("Event bus closed, collector shutting down");
        return Ok(());
    }
}

// ❌ FORBIDDEN — blocking send in collector (could deadlock)
sender.send(event).await?;
```

### Shutdown Signal (Watch)

```rust
use tokio::sync::watch;

pub type ShutdownSender = watch::Sender<bool>;
pub type ShutdownReceiver = watch::Receiver<bool>;

pub fn create_shutdown_signal() -> (ShutdownSender, ShutdownReceiver) {
    watch::channel(false)
}

// In collector:
async fn run_collector(
    mut shutdown: ShutdownReceiver,
    sender: EventSender,
) {
    loop {
        tokio::select! {
            _ = shutdown.changed() => {
                if *shutdown.borrow() {
                    info!("Shutdown signal received");
                    break;
                }
            }
            _ = tokio::time::sleep(poll_interval) => {
                // Do collection work
            }
        }
    }
}
```

---

## Task Spawning Rules

### Rule 1: Always name spawned tasks

```rust
// ✅ Correct — named task for debugging
tokio::task::Builder::new()
    .name("browser-collector")
    .spawn(async move {
        browser_collector.run(shutdown_receiver, event_sender).await
    })?;

// ✅ Also acceptable — spawn with tracing span
tokio::spawn(
    async move {
        browser_collector.run(shutdown_receiver, event_sender).await
    }
    .instrument(tracing::info_span!("browser_collector"))
);
```

### Rule 2: Use `spawn_blocking` for CPU-bound or blocking I/O

```rust
// ✅ Correct — screenshot encoding is CPU-bound
let encoded = tokio::task::spawn_blocking(move || {
    encode_image(&raw_buffer, ImageFormat::WebP, quality)
}).await??;

// ✅ REQUIRED — SQLite is blocking I/O
let results = tokio::task::spawn_blocking(move || {
    storage.query_activities(from, to)
}).await??;

// ❌ FORBIDDEN — blocking call on async thread
let encoded = encode_image(&raw_buffer, ImageFormat::WebP, quality)?;
```

### Rule 3: Handle JoinHandle results

```rust
// ✅ Correct — handle both join error and task error
match collector_handle.await {
    Ok(Ok(())) => info!("Collector finished cleanly"),
    Ok(Err(error)) => error!(%error, "Collector failed"),
    Err(join_error) => error!(%join_error, "Collector task panicked"),
}
```

---

## Graceful Shutdown

```rust
pub async fn run_daemon(config: Config) -> Result<(), AppError> {
    let (shutdown_sender, shutdown_receiver) = create_shutdown_signal();
    let (event_sender, event_receiver) = create_event_bus();

    // Spawn collectors
    let mut handles = Vec::new();
    for collector in collectors {
        let shutdown = shutdown_receiver.clone();
        let sender = event_sender.clone();
        handles.push(tokio::spawn(async move {
            collector.run(shutdown, sender).await
        }));
    }

    // Spawn storage writer
    let storage_handle = tokio::spawn(async move {
        storage_writer.run(event_receiver, shutdown_receiver.clone()).await
    });

    // Wait for shutdown signal
    tokio::select! {
        _ = tokio::signal::ctrl_c() => {
            info!("SIGINT received, initiating shutdown");
        }
        _ = sigterm_future() => {
            info!("SIGTERM received, initiating shutdown");
        }
    }

    // Signal all tasks to stop
    shutdown_sender.send(true)?;

    // Wait for collectors with timeout
    let shutdown_timeout = Duration::from_secs(5);
    match tokio::time::timeout(shutdown_timeout, futures::future::join_all(handles)).await {
        Ok(results) => { /* check results */ }
        Err(_) => warn!(error_code = 15003, "Shutdown timed out after 5s"),
    }

    // Wait for storage to flush
    storage_handle.await??;

    Ok(())
}
```

---

## Cancellation Safety

### Safe Patterns

```rust
// ✅ Cancellation-safe — select on futures that can be safely dropped
tokio::select! {
    event = receiver.recv() => { /* process */ }
    _ = shutdown.changed() => { break; }
    _ = tokio::time::sleep(interval) => { /* poll */ }
}
```

### Unsafe Patterns to Avoid

```rust
// ❌ FORBIDDEN — partial write may be lost if canceled
tokio::select! {
    result = write_batch_to_database(&events) => { /* ... */ }
    _ = shutdown.changed() => { break; }  // Batch partially written!
}

// ✅ REQUIRED — complete the write before checking shutdown
write_batch_to_database(&events).await?;
if *shutdown.borrow() { break; }
```

---

## Send and Sync Bounds for Multi-Threaded Tokio

When spawning concurrent asynchronous tasks using `tokio::spawn`, all spawned futures and captured state must satisfy `Send + Sync + 'static`:

```rust
// ❌ FORBIDDEN — non-thread-safe reference counted pointers across tasks
use std::rc::Rc;
use std::cell::RefCell;

let shared_state = Rc::new(RefCell::new(Vec::new()));
tokio::spawn(async move {
    // Compile error: Rc is not Send
    shared_state.borrow_mut().push(1);
});

// ✅ REQUIRED — atomic reference counting with async-aware synchronization
use std::sync::Arc;
use tokio::sync::Mutex;

let shared_state = Arc::new(Mutex::new(Vec::new()));
let state_clone = Arc::clone(&shared_state);
tokio::spawn(async move {
    let mut lock = state_clone.lock().await;
    lock.push(1);
});
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Guidelines | `../01-cross-language/readme.md` |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/05-rust/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-RUST-004: Rust Async Tokio Cancellation and Send Bounds

**Given** Rust asynchronous tasks, channel endpoints, and Tokio event loops.
**When** Audited against async runtime discipline and cancellation safety standards.
**Then** All spawned futures satisfy `Send + Sync + 'static`, channels use bounded buffers, cancel-safe idioms protect I/O from lost writes, and zero violations are detected with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/05-rust --check-only
```
**Expected:** exit 0. Zero violations.
