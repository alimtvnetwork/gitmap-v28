# Codebase Type Safety Remediation Plan

> **Created:** 2026-02-12  
**Version:** 3.2.0  
> **Priority:** CRITICAL  
> **Spec Reference:** `spec/02-coding-guidelines/02-typescript/08-typescript-standards-reference.md` v2.0.0  
> **Goal:** Eliminate all `any`, `unknown`, `Record<string, unknown>`, string union types, and magic strings/numbers — use proper enums with PascalCase values and typed generics

---

## Audit Summary

| Category | Count | Severity |
|----------|-------|----------|
| `catch (err: any)` | 11 | 🔴 Critical |
| `as any` type assertions | 5 | 🔴 Critical |
| `getQueryData<any>` | 1 | 🔴 Critical |
| `unknown` in API method signatures | 18 | 🔴 Critical |
| `Record<string, unknown>` in API params/returns | 14 | 🔴 Critical |
| `request<unknown>` (untyped endpoints) | 8 | 🟡 High |
| `unknown` in envelope/client internals | 6 | 🟢 Acceptable (parse boundaries) |
| Magic strings (status checks, action types) | ~50+ | 🟡 High |
| Magic numbers (staleTime, limits, intervals) | ~20+ | 🟡 High |

**Total violations: ~130+**

---

## Phase 1: API Type Definitions (Priority: 🔴 CRITICAL)

Create specific interfaces to replace `Record<string, unknown>` and `unknown` in `src/lib/api/`.

### P1.1 — Snapshot operation types

**File:** `src/lib/api/types.ts`

Create:
```typescript
enum SnapshotScopeType {
  All = "All",
  Wordpress = "Wordpress",
  Content = "Content",
  Custom = "Custom",
}

enum SnapshotVariantType {
  Full = "Full",
  Incremental = "Incremental",
}

interface CreateSnapshotOptions {
  readonly name?: string;
  readonly scope?: SnapshotScopeType;
  readonly snapshotType?: SnapshotVariantType;
  readonly parentId?: number;
  readonly tables?: readonly string[];
  readonly workerCount?: number;
}

interface SnapshotOperationResult {
  id: number;
  status: string;
  message?: string;
  snapshotId?: number;
  filename?: string;
}

enum RestoreModeType {
  Full = "Full",
  Selective = "Selective",
}

interface RestoreSnapshotOptions {
  readonly mode?: RestoreModeType;
  readonly tables?: readonly string[];
  readonly preBackup?: boolean;
}

interface CleanupSnapshotOptions {
  readonly dryRun?: boolean;
  readonly maxAge?: number;
  readonly maxCount?: number;
}

interface CleanupSnapshotResult {
  readonly deleted: number;
  readonly dryRun: boolean;
  readonly candidates?: readonly string[];
}

interface SnapshotImportResult {
  id: number;
  filename: string;
  tables: number;
  totalRows: number;
}
```

### P1.2 — Site Health types

**File:** `src/lib/api/types.ts`

Create:
```typescript
enum SiteHealthStatusType {
  Healthy = "Healthy",
  Warning = "Warning",
  Critical = "Critical",
}

interface SiteHealthCheck {
  readonly siteId: number;
  readonly status: SiteHealthStatusType;
  readonly checks: readonly SiteHealthCheckItem[];
  readonly score: number;
  readonly checkedAt: string;
}

interface SiteHealthSummary {
  readonly siteId: number;
  readonly siteName: string;
  readonly status: SiteHealthStatusType;
  readonly score: number;
  readonly lastCheckedAt: string;
}

interface SiteHealthStats {
  readonly totalSites: number;
  readonly healthy: number;
  readonly warning: number;
  readonly critical: number;
  readonly avgScore: number;
}
```

### P1.3 — E2E Test types

**File:** `src/lib/api/types.ts`

Create:
```typescript
enum E2eCaseStatusType {
  Pending = "Pending",
  Running = "Running",
  Passed = "Passed",
  Failed = "Failed",
  Skipped = "Skipped",
}

enum E2eRunStatusType {
  Pending = "Pending",
  Running = "Running",
  Completed = "Completed",
  Aborted = "Aborted",
  Failed = "Failed",
}

interface E2eSuite {
  readonly id: string;
  readonly name: string;
  readonly description?: string;
  readonly caseCount: number;
}

interface E2eCase {
  readonly id: string;
  readonly suiteId: string;
  readonly name: string;
  readonly status: E2eCaseStatusType;
}

interface E2eRun {
  readonly runId: string;
  readonly status: E2eRunStatusType;
  readonly totalTests: number;
  readonly passed: number;
  readonly failed: number;
  readonly skipped: number;
  readonly startedAt: string;
  readonly endedAt?: string;
  readonly results?: readonly E2eTestResult[];
}
```

### P1.4 — Update `methods.ts` signatures

Replace all `Record<string, unknown>` params and `request<unknown>` calls with the specific types from P1.1–P1.3.

---

## Phase 2: Eliminate `catch (err: any)` (Priority: 🔴 CRITICAL)

**Pattern:** Replace all 11 occurrences with bare `catch (err)` + `instanceof` narrowing.

| File | Line(s) | Current | Fix |
|------|---------|---------|-----|
| `SnapshotRetentionPolicy.tsx` | 80 | `catch (err: any)` | `catch (err) { const msg = err instanceof Error ? err.message : String(err); }` |
| `SnapshotRestoreDialog.tsx` | 67 | `catch (err: any)` | Same pattern |
| `Settings.tsx` | 896 | `catch (err: any)` | Same pattern |
| `RemoteSnapshotsPanel.tsx` | 217, 331 | `catch (err: any)` | Same pattern |
| `SyncProgressDialog.tsx` | 423 | `catch (err: any)` | Same pattern |
| `SnapshotSettingsTab.tsx` | 613, 627, 642 | `catch (err: any)` | Same pattern |

---

## Phase 3: Eliminate `as any` Assertions (Priority: 🔴 CRITICAL)

| File | Line | Current | Fix |
|------|------|---------|-----|
| `SnapshotRetentionPolicy.tsx` | 78 | `(result as any)?.deleted` | Type the `cleanupRemoteSnapshots` return as `CleanupSnapshotResult` |
| `ThemeSelector.tsx` | 149, 169 | `v as any` | Create `FontSize` and `BorderRadius` string literal types |
| `useDashboardStats.ts` | 65, 71 | `(data as any).entries` | Type the publish history response properly |
| `useDashboardStats.ts` | 77, 84, 88, 89 | `(e: any)`, `(s: any)`, `(p: any)` | Use typed array callbacks with `ErrorHistoryRecord`, `Site`, `Plugin` |
| `useTheme.ts` | 91 | `(appearance as any).sidebarTheme` | Add `sidebarTheme` to the `Settings.appearance` interface |
| `Dashboard.tsx` | 66 | `getQueryData<any>` | Use `getQueryData<DashboardStats>` with proper type |
| `BackendSection.tsx` | 368, 383, 402, 465–467 | Multiple `any` | Create typed `StackFrame` interface for Go/PHP frames |

---

## Phase 4: Make Envelope Generic (Priority: 🟡 HIGH)

**File:** `src/lib/api/envelope.ts`

```typescript
// Current — ❌ uses unknown
export interface RawEnvelope {
  Results: unknown[];
}

// Target — ✅ generic with typed parameter (never use unknown or any as T)
export interface RawEnvelope<T> {
  Status: EnvelopeStatusType;
  Attributes: EnvelopeAttributes;
  Results: T[];
  Navigation?: EnvelopeNavigation;
  Errors?: EnvelopeErrors;
  MethodsStack?: EnvelopeMethodsStack;
}

// Usage: RawEnvelope<SnapshotOperationResult>, RawEnvelope<SiteHealthCheck>
// ❌ NEVER: RawEnvelope<unknown>, RawEnvelope<any>, RawEnvelope<Record<string, unknown>>
```

---

## Phase 5: Constants & Enums for Magic Strings (Priority: 🟡 HIGH)

### P5.1 — Create `src/lib/constants.ts`

```typescript
// Connection
export enum ConnectionStatusType {
  Connected = "Connected",
  Disconnected = "Disconnected",
}

// Snapshots
export enum SnapshotStatusType {
  Complete = "Complete",
  InProgress = "InProgress",
  Running = "Running",
  Pending = "Pending",
  Failed = "Failed",
}

export enum SnapshotActionType {
  Create = "Create",
  Restore = "Restore",
  Delete = "Delete",
  Export = "Export",
  Import = "Import",
  Cleanup = "Cleanup",
}

// Activity
export enum ActivitySourceType {
  Go = "Go",
  Wordpress = "Wordpress",
}

// Publish
export enum PublishStatusType {
  Success = "Success",
  Failed = "Failed",
  Partial = "Partial",
}

// Timing constants
export const STALE_TIME_MS = 60_000 as const;
export const SNAPSHOT_POLL_INTERVAL_MS = 5_000 as const;
export const DEFAULT_PAGE_SIZE = 25 as const;
export const TOAST_DURATION_ERROR_MS = 10_000 as const;
```

### P5.2 — Migrate all consumers

Search-and-replace all inline string comparisons with enum references.

---

## Phase 6: ActivityEntry Metadata Typing (Priority: 🟡 HIGH)

Replace `metadata: Record<string, unknown>` with a discriminated union using proper enums and typed metadata interfaces:

```typescript
enum ActivityType {
  Publish = "Publish",
  Snapshot = "Snapshot",
  Plugin = "Plugin",
  Config = "Config",
  Connection = "Connection",
}

interface BaseActivityEntry<T extends ActivityType, M> {
  readonly id: string;
  readonly timestamp: string;
  readonly siteId: number;
  readonly siteName: string;
  readonly type: T;
  readonly action: string;
  readonly title: string;
  readonly metadata: M;
  readonly source: ActivitySourceType;
  readonly machineName?: string;
  readonly version?: string;
}

// ✅ Each activity type has its own typed metadata — no Record<string, unknown>
type ActivityEntry =
  | BaseActivityEntry<ActivityType.Publish, PublishMetadata>
  | BaseActivityEntry<ActivityType.Snapshot, SnapshotMetadata>
  | BaseActivityEntry<ActivityType.Plugin, PluginMetadata>
  | BaseActivityEntry<ActivityType.Config, ConfigMetadata>
  | BaseActivityEntry<ActivityType.Connection, ConnectionMetadata>;

// ❌ NEVER use: metadata: Record<string, unknown>
// ❌ NEVER use: metadata: unknown
// ✅ ALWAYS define a concrete metadata interface per activity type
```

---

## Execution Order

| Priority | Phase | Effort | Dependencies |
|----------|-------|--------|--------------|
| 1 | P1 — API type definitions | Medium | None |
| 2 | P2 — Catch block fixes | Small | None |
| 3 | P3 — `as any` elimination | Small | P1 (some fixes need new types) |
| 4 | P5.1 — Constants file | Small | None |
| 5 | P1.4 — Update methods.ts | Medium | P1.1–P1.3 |
| 6 | P4 — Generic envelope | Small | None |
| 7 | P5.2 — Magic string migration | Large | P5.1 |
| 8 | P6 — Activity metadata | Medium | P5.1 |
