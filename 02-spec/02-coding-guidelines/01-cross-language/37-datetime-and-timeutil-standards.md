# Date-Time & Centralized Time Utility Standards — Coding Guidelines

## 1. Overview & Core Mission

Across multi-repository, polyglot software architectures, date and time manipulation represents one of the most prolific sources of bugs, drift, and fragmentation:
1. Developers invent ad-hoc time string formats (`"2006-01-02 15:04:05"`, `RFC3339`, `RFC3339Nano`, ISO 8601, millisecond timestamps) scattered inline throughout business logic.
2. Inconsistent timezone handling leads to mixed local and UTC timestamps in persistent databases.
3. Scattered direct calls to standard library time functions (such as `time.Now()` in Go, `new Date()` in TypeScript, `datetime.now()` in Python, `date()` in PHP) make deterministic unit testing and clock virtualization impossible.

This specification establishes the **Centralized Time Utility (`timeutil`) Standard** as a non-negotiable coding guideline. Direct ad-hoc date-time string formatting is strictly forbidden. All timestamp generation and parsing must route through centralized, mockable time utilities adhering strictly to UTC ISO 8601 (`YYYY-MM-DDTHH:MM:SSZ`).

---

## 2. Anti-Patterns: Ad-Hoc Time Formatting & Direct Clocks

### ❌ Anti-Pattern: Scattered Direct Format Calls in Go
```go
// ❌ FORBIDDEN: Direct ad-hoc timestamp formatting scattered in domain logic
now := time.Now().UTC().Format(time.RFC3339)
updated := time.Now().Format("2006-01-02T15:04:05Z07:00")
stamp := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
```

### ❌ Anti-Pattern: Non-UTC and Ad-Hoc Strings in TypeScript
```typescript
// ❌ FORBIDDEN: Inconsistent Date formatting and local timezone leakage
const timestamp = new Date().toISOString(); // scattered unmocked call
const localTime = new Date().toLocaleString(); // local timezone contamination
```

### ❌ Anti-Pattern: Naive Datetimes in Python
```python
# ❌ FORBIDDEN: Naive local time and scattered strftime calls
from datetime import datetime
now = datetime.now().strftime("%Y-%m-%d %H:%M:%S") # naive local time
```

### ❌ Anti-Pattern: Ad-Hoc Date Strings in PHP
```php
// ❌ FORBIDDEN: Direct date() calls with inconsistent formatting
$createdAt = date('Y-m-d H:i:s');
```

---

## 3. The Canonical Pattern: Centralized `timeutil` Engine

### 3.1 Standard Canonical Layouts
- **Primary ISO 8601 UTC:** `YYYY-MM-DDTHH:MM:SSZ` (RFC3339 without sub-second noise, standard for DB and API payloads).
- **Compact Timestamp (for IDs and slugs):** `HHMMSS` or `YYYYMMDDHHMMSS`.
- **Date Only:** `YYYY-MM-DD`.

### 3.2 Canonical Implementation in Go (`cli/timeutil/timeutil.go`)

```go
// Package timeutil provides centralized RFC3339 UTC time formatting helpers.
package timeutil

import (
	"time"
)

// TimeProvider defines an injectable clock interface for deterministic unit testing.
type TimeProvider interface {
	Now() time.Time
}

type realTimeProvider struct{}

func (realTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

var currentProvider TimeProvider = realTimeProvider{}

// SetTimeProvider overrides the global time provider (for testing).
func SetTimeProvider(p TimeProvider) {
	if p == nil {
		currentProvider = realTimeProvider{}
		return
	}
	currentProvider = p
}

// ResetTimeProvider resets the provider to standard UTC clock.
func ResetTimeProvider() {
	currentProvider = realTimeProvider{}
}

// NowUTC returns the current UTC time from the active provider.
func NowUTC() time.Time {
	return currentProvider.Now()
}

// NowUTCISO returns the current time formatted as UTC RFC3339 (ISO 8601).
func NowUTCISO() string {
	return currentProvider.Now().Format(time.RFC3339)
}

// FormatUTCISO formats the specified time as UTC RFC3339.
func FormatUTCISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// NowTimestampCompact returns a compact UTC timestamp (HHMMSS).
func NowTimestampCompact() string {
	return currentProvider.Now().Format("150405")
}
```

### 3.3 Canonical Implementation in TypeScript (`src/utils/timeutil.ts`)

```typescript
export interface TimeProvider {
  now(): Date;
}

class SystemTimeProvider implements TimeProvider {
  now(): Date {
    return new Date();
  }
}

let activeProvider: TimeProvider = new SystemTimeProvider();

export function setTimeProvider(provider: TimeProvider): void {
  activeProvider = provider;
}

export function resetTimeProvider(): void {
  activeProvider = new SystemTimeProvider();
}

/** Returns the current time formatted as ISO 8601 UTC string. */
export function nowUtcIso(): string {
  return activeProvider.now().toISOString();
}

/** Formats a given Date instance into ISO 8601 UTC string. */
export function formatUtcIso(d: Date): string {
  return d.toISOString();
}

/** Returns a compact HHMMSS timestamp in UTC. */
export function nowTimestampCompact(): string {
  const d = activeProvider.now();
  const pad = (n: number) => n.toString().padStart(2, '0');
  return `${pad(d.getUTCHours())}${pad(d.getUTCMinutes())}${pad(d.getUTCSeconds())}`;
}
```

### 3.4 Canonical Implementation in Python (`pkg/timeutil.py`)

```python
from datetime import datetime, timezone
from typing import Protocol

class TimeProvider(Protocol):
    def now(self) -> datetime:
        ...

class SystemTimeProvider:
    def now(self) -> datetime:
        return datetime.now(timezone.utc)

_active_provider: TimeProvider = SystemTimeProvider()

def set_time_provider(provider: TimeProvider) -> None:
    global _active_provider
    _active_provider = provider

def reset_time_provider() -> None:
    global _active_provider
    _active_provider = SystemTimeProvider()

def now_utc_iso() -> str:
    """Returns the current UTC time as an ISO 8601 string."""
    return _active_provider.now().strftime("%Y-%m-%dT%H:%M:%SZ")

def format_utc_iso(dt: datetime) -> str:
    """Formats any datetime into UTC ISO 8601."""
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    return dt.astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
```

### 3.5 Canonical Implementation in PHP (`src/Utils/DateUtil.php`)

```php
<?php

declare(strict_types=1);

namespace GitMap\Utils;

use DateTimeImmutable;
use DateTimeZone;

class DateUtil
{
    private static ?DateTimeImmutable $mockedNow = null;

    public static function setMockedNow(?DateTimeImmutable $mock): void
    {
        self::$mockedNow = $mock;
    }

    public static function nowUtc(): DateTimeImmutable
    {
        if (self::$mockedNow !== null) {
            return self::$mockedNow->setTimezone(new DateTimeZone('UTC'));
        }

        return new DateTimeImmutable('now', new DateTimeZone('UTC'));
    }

    public static function nowUtcIso(): string
    {
        return self::nowUtc()->format('Y-m-d\TH:i:s\Z');
    }

    public static function formatUtcIso(DateTimeImmutable $dt): string
    {
        return $dt->setTimezone(new DateTimeZone('UTC'))->format('Y-m-d\TH:i:s\Z');
    }
}
```

---

## 4. Deterministic Unit Testing with Virtual Clocks

When unit tests assert timestamp columns or temporal expiration windows, never call `time.Sleep()` or rely on wall clock advances:

```go
// ✅ TEST PATTERN: Deterministic time testing via MockTimeProvider
type mockTimeProvider struct {
	fixedTime time.Time
}

func (m mockTimeProvider) Now() time.Time {
	return m.fixedTime
}

func TestTaskExpiry(t *testing.T) {
	baseTime := time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
	timeutil.SetTimeProvider(mockTimeProvider{fixedTime: baseTime})
	defer timeutil.ResetTimeProvider()

	task := CreateTask("task-01")
	if task.CreatedAt != "2026-10-22T12:00:00Z" {
		t.Fatalf("expected fixed timestamp, got %s", task.CreatedAt)
	}
}
```

---

## 5. Summary Checklist

- [ ] All timestamps generated across the repository call `timeutil.NowUTCISO()` (or language counterpart).
- [ ] Zero inline string formatting `Format(time.RFC3339)` in business logic or command handlers.
- [ ] All database columns storing timestamps (`CreatedAt`, `UpdatedAt`) are populated via `timeutil`.
- [ ] Clock providers support deterministic override for zero-sleep unit testing.
