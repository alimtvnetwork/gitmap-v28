# Split DB Architecture: Complete Database Flow Diagram

**Version:** 3.2.0  
**Status:** Active  
**Updated:** 2026-04-16  
**Parent:** [00-overview.md](../00-overview.md)

---

## Overview

This document provides visual architecture diagrams showing the complete database flow from API request to SQLite storage across all CLI tools.

---

## Master Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    SPLIT DB ARCHITECTURE - ALL CLIs                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                  │
│   ┌───────────────────────────────────────────────────────────────────────────────────────────────────────────┐ │
│   │                                        SHARED INFRASTRUCTURE                                               │ │
│   │                                                                                                            │ │
│   │   ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐           │ │
│   │   │  config.seed.json │    │   Seedable       │    │   PascalCase     │    │   2-Step Reset   │           │ │
│   │   │  (Defaults)       │───▶│   Config         │───▶│   Naming         │───▶│   API            │           │ │
│   │   └──────────────────┘    └──────────────────┘    └──────────────────┘    └──────────────────┘           │ │
│   │                                                                                                            │ │
│   └───────────────────────────────────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                                                  │
│   ┌─────────────────────────────────────┐   ┌─────────────────────────────────────┐                             │
│   │         AI BRIDGE CLI               │   │           GSearch CLI                │                             │
│   │            :8089                     │   │             :8087                    │                             │
│   ├─────────────────────────────────────┤   ├─────────────────────────────────────┤                             │
│   │                                      │   │                                      │                             │
│   │   data/aibridge.db ◀── ROOT DB      │   │   data/gsearch.db ◀── ROOT DB       │                             │
│   │         │                            │   │         │                            │                             │
│   │   data/{app}/                        │   │   data/searches/                     │                             │
│   │   ├── ai/chat/{seq}-{id}.db         │   │   ├── search.db ◀── APP DB          │                             │
│   │   ├── rag/documents/{id}.db         │   │   └── cache/                         │                             │
│   │   ├── rag/cache/search/{id}.db      │   │       └── {seq}-{slug}.db           │                             │
│   │   └── seo/jobs/{id}.db              │   │                                      │                             │
│   │                                      │   │                                      │                             │
│   └─────────────────────────────────────┘   └─────────────────────────────────────┘                             │
│                                                                                                                  │
│   ┌─────────────────────────────────────┐   ┌─────────────────────────────────────┐                             │
│   │           BRun CLI                   │   │        Nexus Flow CLI                │                             │
│   │            :8088                     │   │            :8085                     │                             │
│   ├─────────────────────────────────────┤   ├─────────────────────────────────────┤                             │
│   │                                      │   │                                      │                             │
│   │   data/brun.db ◀── ROOT DB          │   │   data/nexusflow.db ◀── ROOT DB     │                             │
│   │         │                            │   │         │                            │                             │
│   │   data/runs/                         │   │   data/workflows/                    │                             │
│   │   └── {seq}-{id}.db ◀── SESSION DB  │   │   └── pipeline-{id}/                 │                             │
│   │                                      │   │       ├── meta.db                    │                             │
│   │                                      │   │       ├── executions/{id}.db        │                             │
│   │                                      │   │       └── checkpoints/{id}.db       │                             │
│   │                                      │   │                                      │                             │
│   └─────────────────────────────────────┘   └─────────────────────────────────────┘                             │
│                                                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## AI Bridge: API to Database Flow

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                 AI BRIDGE: REQUEST → DATABASE FLOW                                               │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                  │
│   CLIENT REQUEST                         API HANDLER                           DATABASE OPERATION               │
│   ─────────────────                      ───────────                           ──────────────────               │
│                                                                                                                  │
│   POST /chat/sessions ──────────────────▶ CreateSession() ─────────────────────▶ Creates:                       │
│   {                                        │                                      data/{app}/ai/chat/           │
│     "AppName": "myapp",                    ├─▶ Get next sequence from            {seq}-{id}.db                  │
│     "Title": "Code Review"                 │   aibridge.db → Counters                                           │
│   }                                        │                                                                     │
│                                            └─▶ Register in aibridge.db                                          │
│                                                → DbRegistry                                                      │
│                                                                                                                  │
│   ─────────────────────────────────────────────────────────────────────────────────────────────────────────────  │
│                                                                                                                  │
│   POST /chat/sessions/:id/messages ─────▶ SendMessage() ───────────────────────▶ Writes to:                     │
│   {                                        │                                      {seq}-{id}.db                  │
│     "Content": "Explain main.go"           ├─▶ Load RAG context from               │                             │
│   }                                        │   rag/documents/*.db                  ├─▶ Messages table            │
│                                            │   (Vector search)                     └─▶ ToolCalls table           │
│                                            │                                                                     │
│                                            ├─▶ Call LLM with context                                            │
│                                            │                                                                     │
│                                            └─▶ Stream response via SSE                                          │
│                                                                                                                  │
│   ─────────────────────────────────────────────────────────────────────────────────────────────────────────────  │
│                                                                                                                  │
│   POST /rag/documents ──────────────────▶ IngestDocument() ────────────────────▶ Creates:                       │
│   (multipart: file)                        │                                      data/{app}/rag/documents/     │
│                                            ├─▶ Chunk content (2048 tokens)        {seq}-{id}.db                  │
│                                            │                                        │                             │
│                                            ├─▶ Generate embeddings                  ├─▶ DocumentMeta             │
│                                            │                                        └─▶ Chunks + Embeddings      │
│                                            └─▶ Register in aibridge.db                                          │
│                                                                                                                  │
│   ─────────────────────────────────────────────────────────────────────────────────────────────────────────────  │
│                                                                                                                  │
│   POST /reset/request ──────────────────▶ RequestReset() ──────────────────────▶ Writes to:                     │
│   { "Scope": "app" }                       │                                      aibridge.db                    │
│                                            └─▶ Create ResetRequest                  │                             │
│                                                ExpiresAt = NOW + 5min               └─▶ ResetRequests table      │
│                                                                                                                  │
│   POST /reset/confirm ──────────────────▶ ConfirmReset() ──────────────────────▶ Deletes:                       │
│   { "ResetId": "rst_abc" }                 │                                      All DBs in scope               │
│                                            ├─▶ Validate not expired                                              │
│                                            └─▶ Delete databases                                                  │
│                                                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## RAG Context Loading Flow

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                       RAG MEMORY LOADING FLOW                                                    │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                  │
│   1. USER MESSAGE RECEIVED                                                                                       │
│      "For this codebase, explain the authentication flow"                                                        │
│                                                                                                                  │
│   2. EMBEDDING GENERATION                                                                                        │
│      ┌──────────────────┐                                                                                        │
│      │ Message → Vector │  (NomicEmbedText)                                                                    │
│      │ [0.12, -0.45...] │                                                                                        │
│      └────────┬─────────┘                                                                                        │
│               │                                                                                                  │
│   3. VECTOR SEARCH ACROSS RAG DOCUMENTS                                                                          │
│               ▼                                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────────────────────────────────────┐   │
│   │   data/{app}/rag/documents/                                                                              │   │
│   │                                                                                                          │   │
│   │   ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐                │   │
│   │   │ 001-readme.db    │  │ 002-auth.db      │  │ 003-config.db    │  │ 004-main.db      │                │   │
│   │   │                  │  │                  │  │                  │  │                  │                │   │
│   │   │ Chunks:          │  │ Chunks:          │  │ Chunks:          │  │ Chunks:          │                │   │
│   │   │  ├ chunk_001     │  │  ├ chunk_001 ◀───┼──┼──similarity: 0.94│  │  ├ chunk_001     │                │   │
│   │   │  ├ chunk_002     │  │  ├ chunk_002 ◀───┼──┼──similarity: 0.89│  │  └ chunk_002     │                │   │
│   │   │  └ chunk_003     │  │  └ chunk_003 ◀───┼──┼──similarity: 0.82│  │                  │                │   │
│   │   └──────────────────┘  └──────────────────┘  └──────────────────┘  └──────────────────┘                │   │
│   │                                                                                                          │   │
│   └─────────────────────────────────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                                                  │
│   4. CONTEXT ASSEMBLY (Token Budget: 4096)                                                                       │
│      ┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│      │ TOP-K CHUNKS (K=10, threshold > 0.7)                                                                 │    │
│      │                                                                                                      │    │
│      │   Rank 1: 002-auth.db/chunk_001  (0.94) → "Authentication uses JWT tokens..."     [512 tokens]      │    │
│      │   Rank 2: 002-auth.db/chunk_002  (0.89) → "The login flow validates..."           [512 tokens]      │    │
│      │   Rank 3: 002-auth.db/chunk_003  (0.82) → "Session management handles..."         [512 tokens]      │    │
│      │   ...                                                                                                │    │
│      │                                                                                                      │    │
│      │   Total: 3,584 tokens (under 4,096 budget)                                                          │    │
│      └─────────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                                  │
│   5. SYSTEM PROMPT INJECTION                                                                                     │
│      ┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│      │ [System] You are an AI assistant. Use the following context:                                         │    │
│      │                                                                                                      │    │
│      │ --- CONTEXT FROM CODEBASE ---                                                                        │    │
│      │ [File: auth.go, Chunk 1] Authentication uses JWT tokens...                                           │    │
│      │ [File: auth.go, Chunk 2] The login flow validates...                                                 │    │
│      │ [File: auth.go, Chunk 3] Session management handles...                                               │    │
│      │ --- END CONTEXT ---                                                                                  │    │
│      │                                                                                                      │    │
│      │ [User] For this codebase, explain the authentication flow                                            │    │
│      └─────────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                                  │
│   6. LLM RESPONSE                                                                                                │
│      "Based on the codebase, authentication works as follows..."                                                 │
│                                                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Search Cache Flow (GSearch + AI Bridge)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    SEARCH CACHE FLOW (DELEGATION)                                                │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                  │
│   AI BRIDGE                                          GSearch CLI                                                 │
│   ─────────                                          ───────────                                                 │
│                                                                                                                  │
│   User: "Search for Go concurrency patterns"                                                                     │
│         │                                                                                                        │
│         ▼                                                                                                        │
│   ┌─────────────────┐                                                                                            │
│   │ AI Bridge       │                                                                                            │
│   │ Agentic Mode    │─────── gsearch search ────────▶┌──────────────────────────────────────────────────────┐   │
│   │ (Tool Call)     │        --query "..."           │                                                       │   │
│   └─────────────────┘        --output json           │   1. Check cache hash                                 │   │
│                                                       │      data/gsearch.db → DbRegistry                    │   │
│                              ◀───────────────────────│                                                       │   │
│                              JSON results             │   2a. CACHE HIT                                      │   │
│                                                       │       └─▶ Read from cache/{seq}-{slug}.db           │   │
│   ┌─────────────────┐                                │                                                       │   │
│   │ AI Bridge       │                                │   2b. CACHE MISS                                      │   │
│   │ Cache Result    │                                │       ├─▶ Execute search                             │   │
│   └────────┬────────┘                                │       ├─▶ Create cache/{seq}-{slug}.db               │   │
│            │                                          │       └─▶ Set ExpiresAt = NOW + 5 days               │   │
│            ▼                                          │                                                       │   │
│   data/{app}/rag/cache/search/                       │   3. Return results                                   │   │
│   {seq}-{slug}.db                                    │                                                       │   │
│    ├── CacheMeta                                     └──────────────────────────────────────────────────────┘   │
│    ├── Results                                                                                                   │
│    └── Embeddings (optional)                                                                                     │
│                                                                                                                  │
│   ─────────────────────────────────────────────────────────────────────────────────────────────────────────────  │
│                                                                                                                  │
│   TTL CONFIGURATION (Seedable)                                                                                   │
│   ┌────────────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│   │                                                                                                         │    │
│   │   Priority Order:                                                                                       │    │
│   │   1. App-level override:  data/{app}/search.db → CacheSettings.TtlDays                                 │    │
│   │   2. Root DB setting:     data/aibridge.db → Settings["Search.Cache.TtlDays"]                          │    │
│   │   3. Seed default:        config.seed.json → toolDelegations.webSearch.cacheTtl = 5                    │    │
│   │                                                                                                         │    │
│   └────────────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Reset API Flow (All CLIs)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                         2-STEP RESET FLOW                                                        │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                  │
│   STEP 1: REQUEST                                                                                                │
│   ───────────────                                                                                                │
│                                                                                                                  │
│   POST /api/v1/reset/request                                                                                    │
│   { "Scope": "app", "AppName": "myapp" }                                                                        │
│         │                                                                                                        │
│         ▼                                                                                                        │
│   ┌─────────────────────────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                                          │   │
│   │   1. Scan affected databases                                                                             │   │
│   │      └─▶ Query DbRegistry WHERE ApplicationId = "myapp"                                                 │   │
│   │                                                                                                          │   │
│   │   2. Calculate impact                                                                                    │   │
│   │      └─▶ { ChatSessions: 12, RagDocuments: 45, SearchCaches: 30, TotalBytes: 500MB }                   │   │
│   │                                                                                                          │   │
│   │   3. Create ResetRequest                                                                                 │   │
│   │      ┌─────────────────────────────────────────────────────────────────────────────────────────────┐    │   │
│   │      │ INSERT INTO ResetRequests (Id, Scope, AppName, ExpiresAt, AffectedItems, Status)             │    │   │
│   │      │ VALUES ('rst_abc123', 'app', 'myapp', datetime('now', '+5 minutes'), '{...}', 'pending')     │    │   │
│   │      └─────────────────────────────────────────────────────────────────────────────────────────────┘    │   │
│   │                                                                                                          │   │
│   └─────────────────────────────────────────────────────────────────────────────────────────────────────────┘   │
│         │                                                                                                        │
│         ▼                                                                                                        │
│   Response: { "ResetId": "rst_abc123", "ExpiresAt": "...", "AffectedItems": {...} }                             │
│                                                                                                                  │
│   ══════════════════════════════════════════════════════════════════════════════════════════════════════════════│
│                                                                                                                  │
│   STEP 2: CONFIRM (within 5 minutes)                                                                            │
│   ─────────────────────────────────                                                                              │
│                                                                                                                  │
│   POST /api/v1/reset/confirm                                                                                    │
│   { "ResetId": "rst_abc123" }                                                                                   │
│         │                                                                                                        │
│         ▼                                                                                                        │
│   ┌─────────────────────────────────────────────────────────────────────────────────────────────────────────┐   │
│   │                                                                                                          │   │
│   │   1. Validate request                                                                                    │   │
│   │      ├─▶ Check Id exists                                                                                │   │
│   │      ├─▶ Check Status = 'pending'                                                                       │   │
│   │      └─▶ Check ExpiresAt > NOW                                                                          │   │
│   │                                                                                                          │   │
│   │   2. Execute deletion                                                                                    │   │
│   │      ┌─────────────────────────────────────────────────────────────────────────────────────────────┐    │   │
│   │      │ FOR EACH db IN DbRegistry WHERE ApplicationId = "myapp":                                     │    │   │
│   │      │   a. Close database connection                                                               │    │   │
│   │      │   b. DELETE file: data/{app}/ai/chat/{seq}-{id}.db                                          │    │   │
│   │      │   c. DELETE from DbRegistry                                                                  │    │   │
│   │      │   d. Update Counters                                                                         │    │   │
│   │      └─────────────────────────────────────────────────────────────────────────────────────────────┘    │   │
│   │                                                                                                          │   │
│   │   3. Update ResetRequest                                                                                 │   │
│   │      └─▶ Status = 'completed', DeletedCount = 87, FreedBytes = 524288000                               │   │
│   │                                                                                                          │   │
│   └─────────────────────────────────────────────────────────────────────────────────────────────────────────┘   │
│         │                                                                                                        │
│         ▼                                                                                                        │
│   Response: { "Status": "completed", "DeletedDatabases": 87, "FreedBytes": 524288000 }                          │
│                                                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Database Terminology Reference

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                       DATABASE TERMINOLOGY                                                       │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                  │
│   TERM              │ MEANING                              │ EXAMPLE PATH                                        │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   ROOT DB           │ Global registry, settings, app list  │ data/aibridge.db                                   │
│   (Setting DB)      │                                      │ data/gsearch.db                                    │
│                     │                                      │ data/brun.db                                       │
│                     │                                      │ data/nexusflow.db                                  │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   APP DB            │ Application-scoped metadata          │ data/{appName}/search.db                           │
│                     │                                      │ data/{appName}/settings/config.db                  │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   SESSION DB        │ Per-session isolated storage         │ data/{app}/ai/chat/001-{id}.db                     │
│                     │                                      │ data/runs/001-{id}.db                              │
│                     │                                      │ data/workflows/{p}/executions/001-{id}.db          │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   CACHE DB          │ Cached results with TTL              │ data/{app}/rag/cache/search/001-{slug}.db          │
│                     │                                      │ data/searches/cache/001-{slug}.db                  │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   DOCUMENT DB       │ RAG chunks + embeddings              │ data/{app}/rag/documents/001-{id}.db               │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   META DB           │ Definition/configuration             │ data/workflows/pipeline-001/meta.db                │
│   ──────────────────┼──────────────────────────────────────┼─────────────────────────────────────────────────── │
│   CHECKPOINT DB     │ RES fault tolerance data             │ data/workflows/{p}/checkpoints/001-{id}.db         │
│                                                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Split DB Overview | `../00-overview.md` |
| CLI Examples | `./01-cli-examples.md` |
| Reset API Standard | `./02-reset-api-standard.md` |
| AI Bridge DB | `../22-ai-bridge-cli/01-backend/12-database-architecture.md` |
| GSearch DB | `../20-gsearch-cli/01-backend/22-database-architecture.md` |
| BRun DB | `../21-brun-cli/01-backend/16-database-architecture.md` |
| Nexus Flow DB | `../24-nexus-flow-cli/01-backend/05-database-architecture.md` |
