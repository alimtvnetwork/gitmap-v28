# Rust Async Patterns

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
// ✅ Correct — try_send to avoid blocking collectors
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

// ❌ Forbidden — blocking send in collector (could deadlock)
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

// ✅ Correct — SQLite is blocking I/O
let results = tokio::task::spawn_blocking(move || {
    storage.query_activities(from, to)
}).await??;

// ❌ Forbidden — blocking call on async thread
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
// ❌ Dangerous — partial write may be lost if cancelled
tokio::select! {
    result = write_batch_to_database(&events) => { /* ... */ }
    _ = shutdown.changed() => { break; }  // Batch partially written!
}

// ✅ Fix — complete the write before checking shutdown
write_batch_to_database(&events).await?;
if *shutdown.borrow() { break; }
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Cross-Language Guidelines | `../01-cross-language/00-overview.md` |
