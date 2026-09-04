# Architectural Research & Explorations

> **Location:** `research/`  
> **Purpose:** In-depth architectural research, design explorations, contract proposals, and review blueprints prior to production implementation.  
> **Rule:** All research documents must be strictly lowercase, self-contained, provide full working code samples, and define clear extension points for external developers.

---

## Index of Research Documents

| # | Document | Topic | Status | Date |
|---|---|---|---|---|
| 02 | [`02-pluggable-writer-architecture-and-composition.md`](02-pluggable-writer-architecture-and-composition.md) | Composable Writer contracts, BaseWriter embedding, RestAPIWriter batching, JSON/Text configurable formatting, and custom extension guides | Pending Review | 2026-09-03 |
| 03 | [`03-swappable-writer-methods-and-functional-injection.md`](03-swappable-writer-methods-and-functional-injection.md) | 4 patterns for swappable write methods, functional options injection, log-agnostic payloads, and AUK Go design comparison | Pending Review | 2026-09-03 |
| 04 | [`04-locked-and-lockless-streamers-with-self-binding-interfacer.md`](04-locked-and-lockless-streamers-with-self-binding-interfacer.md) | Two types of streamers (Locked vs Lockless), swappable StreamFunc, and self-binding AsInterfacer/AsWriter contracts | Pending Review | 2026-09-03 |
| 05 | [`05-streamer-and-writer-full-flow.md`](05-streamer-and-writer-full-flow.md) | Full flow implementation & verification of Locked/Lockless streamers, PluggableWriter, and CompositeLogger | Implemented | 2026-09-04 |
| 06 | [`06-generic-payload-and-ordered-compilation.md`](06-generic-payload-and-ordered-compilation.md) | Generic payload T, Compilable interface, and recursive order-wise transpilation for primitives, maps, slices, and structs | Implemented | 2026-09-04 |
| 07 | [`07-bytes-wrapper-and-apperror-standard.md`](07-bytes-wrapper-and-apperror-standard.md) | Bytes[T] monadic wrapper replacing ([]byte, error), and mandatory *appfault.AppError return standard | Implemented | 2026-09-04 |
