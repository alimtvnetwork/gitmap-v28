# Subtask 03: Design System Multi-Theme Palettes & Split-DB SQLite Sharding

## Status: READY
- **Subtask Slug:** `03-design-tokens-and-splitdb`
- **Parent Plan:** `.ai-memory/plans/pending/340-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/04-improvements-361-to-400-design-tokens-splitdb-agy.md`
- **Owner:** Spec Author Subagent 3 / Design Tokens & Storage Engine Specialist
- **Target Subsystem:** `pkg/web/tokens/`, `pkg/db/splitdb/`, `pkg/db/cache/`, `pkg/mail/suppression/`
- **Exclusive Owned Files:**
  - `pkg/web/tokens/color_ramp.go`
  - `pkg/web/tokens/theme_registry.go`
  - `pkg/web/tokens/accessibility_manager.go`
  - `pkg/web/tokens/motion_engine.go`
  - `pkg/db/splitdb/pragmas.go`
  - `pkg/db/splitdb/shard_partitioner.go`
  - `pkg/db/cache/query_cache.go`
  - `pkg/mail/suppression/bloom_filter.go`
  - `web/static/css/tokens.css`
  - `cli/cmddb/cmd_shard.go`

---

## 1. Context & Objectives

This subtask delivers multi-theme design tokens, accessible motion physics, and high-performance Split-DB storage sharding:
1. **Design System Multi-Theme Palettes & Accessible Motion (Cluster 37, Items 361–370):** Mathematical OKLCH 10-step color ramp, 5 production glassmorphic themes, WCAG 2.1 AAA high-contrast modes, fluid typography clamp tokens, and zero-CDN self-contained asset bundling.
2. **Split-DB SQLite Sharding & High-Speed Query Cache (Cluster 38, Items 371–380):** High-throughput WAL pragmas, monthly historical table sharding, LRU memory query cache with TTL, in-memory Bloom filter suppression checks, and background WAL checkpoint governor.
3. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isThemeActive`, `hasOptimalPragmas`, `isCacheHit`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Design System CSS Custom Properties & Palettes
- Implement OKLCH mathematical color ramp generator in `pkg/web/tokens/color_ramp.go`.
- Define 5 glassmorphic theme palettes and runtime hot-reloader in `pkg/web/tokens/theme_registry.go`.
- Author responsive CSS design tokens and fluid typography clamp rules in `web/static/css/tokens.css`.
- Implement WCAG 2.1 AAA high-contrast toggle and motion spring adapter in `pkg/web/tokens/accessibility_manager.go` and `motion_engine.go`.

### Step 2: Split-DB SQLite Storage Sharding & Tuning
- Configure high-performance SQLite WAL pragmas (`journal_mode=WAL`, `synchronous=NORMAL`) in `pkg/db/splitdb/pragmas.go`.
- Build monthly audit table sharding and dynamic `ATTACH DATABASE` router in `pkg/db/splitdb/shard_partitioner.go`.
- Construct thread-safe LRU in-memory query cache with TTL in `pkg/db/cache/query_cache.go`.

### Step 3: Probabilistic Suppression & Memory Filters
- Implement space-optimized in-memory Bloom filter for instant suppression checking in `pkg/mail/suppression/bloom_filter.go`.
- Wire up SQLite shard inspection CLI commands in `cli/cmddb/cmd_shard.go`.
