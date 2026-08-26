# Error Code Registry - Master List

> **Last Updated:** 2026-03-09  
**Version:** 3.2.0  
> **Maintainer:** AI/Human collaboration  
> **Wave 1 Remediation:** Error collisions resolved, paths corrected, missing ranges registered

---

## Registered Project Prefixes

| Prefix | Project | Range | Spec Location | Status |
|--------|---------|-------|---------------|--------|
| `GEN` | General/Shared | 1000-1999 | (embedded) | ✅ Active |
| `SM` | Spec Management Software | 2000-2999 | `spec/02-spec-management-software/` | ✅ Active |
| `LM` | Link Manager | 3000-3999 | `spec/13-wp-plugin/04-link-manager/` | ⚠️ Deprecated — range reassigned to 15000-15999 |
| `CLI` | CLI Tools (legacy) | 4000-4999 | (deprecated) | ⚠️ Deprecated |
| `GS` | GSearch CLI Core | 7000-7099 | `spec/09-gsearch-cli/` | ✅ Active |
| `BR` | BRun CLI | 7100-7599 | `spec/10-brun-cli/` | ✅ Active |
| `GS` | GSearch Movie Search | 7600-7609 | `spec/09-gsearch-cli/01-backend/` | ✅ Active |
| `GS` | GSearch BI Suite | 7700-7839 | `spec/09-gsearch-cli/01-backend/openapi-bi-suite.yaml` | ✅ Active |
| `GS` | GSearch Multi-Source | 7840-7859 | `spec/09-gsearch-cli/01-backend/` | ✅ Active |
| `GS` | GSearch Scheduled | 7860-7879 | `spec/09-gsearch-cli/01-backend/` | ✅ Active |
| `GS` | GSearch Chrome Extension | 7880-7899 | `spec/09-gsearch-cli/01-backend/` | ✅ Active |
| `GS` | GSearch Enum Architecture | 7900-7919 | `spec/09-gsearch-cli/01-backend/` | ✅ Active |
| `GS` | GSearch Provider Integration | 7920-7949 | `spec/09-gsearch-cli/01-backend/` | ✅ Active |
| `NF` | Nexus Flow | 8000-8399 | `spec/12-nexus-flow-cli/` | ✅ Active |
| `AB` | AI Bridge Core | 9000-9499 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `PS` | PowerShell Integration | 9500-9599 | `spec/11-powershell-integration/` | ✅ Active |
| `AB` | AI Bridge SEO | 9500-9540 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Extended | 9600-9699 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Revisions/Suggestions | 9700-9749 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge RAG Session Memory | 9750-9809 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Adaptive Reasoning | 9810-9829 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge WebSocket Resilience | 9830-9839 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Context Integration | 9840-9847 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Research Mode | 9848-9849 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Code Pattern Learning | 9850-9869 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Plan Generation | 9870-9889 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Plan Synchronization | 9890-9909 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Plan Templates | 9910-9929 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Execution Monitoring | 9930-9949 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Retry Strategies | 9950-9969 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Long-Chain Commands | 9970-9989 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Vector DB | 9990-9999 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB` | AI Bridge Lovable Reasoning | 19000-19019 | `spec/11-ai-bridge-cli/` | ✅ Active |
| `AB-TR` | AI Bridge Non-Vector RAG | 20000-20999 | `spec/33-ai-bridge-non-vector-rag/` | ✅ Active |
| `WPB` | WP Plugin Builder | 10000-10499 | `spec/14-wp-plugin-builder/` | ✅ Active |
| `SRC` | Spec Reverse CLI | 11000-11999 | `spec/15-spec-reverse-cli/` | ✅ Active |
| `WSP` | WP SEO Publish | 12000-12599 | `spec/21-wp-seo-publish-cli/` | ✅ Active |
| `WPP` | WP Plugin Publish | 13000-13999 | `spec/13-wp-plugin/05-wp-plugin-publish/` | ✅ Active |
| `AIT` | AI Transcribe CLI | 14000-14499 | `spec/16-ai-transcribe-cli/` | ✅ Active |
| `EQM` | Exam Manager | 14500-14999 | `spec/13-wp-plugin/03-exam-manager/` | ✅ Active |
| `LM` | Link Manager | 15000-15999 | `spec/13-wp-plugin/04-link-manager/` | ✅ Active |
| `SM-CG` | SM Code Generation | 16000-16799 | `spec/02-spec-management-software/05-features/24-code-generation-system/` | ✅ Active |
| `SM-PE` | SM Project Editor | 17000-17999 | `spec/02-spec-management-software/05-features/28-project-editor/` | ✅ Active |
| `SM-GS` | SM GSearch CLI | 18000-18249 | `spec/02-spec-management-software/05-features/22-golang-search-cli/` | ✅ Active |
| `SM-RT` | SM Realtime | 2800-2849 | `spec/02-spec-management-software/05-features/18-realtime/` | ✅ Active |
| `SM-RV` | SM Registry Validator | 2850-2859 | `spec/03-error-manage/03-error-code-registry/08-overlap-validator.md` | ✅ Active |
| `CAST` | Type Casting (Cross-Cutting) | GEN-600-01 to GEN-600-10 | `spec/02-coding-guidelines/01-cross-language/03-casting-elimination-patterns.md` | ✅ Active |
| `MWS` | Main/Worker Service (Worker tier) | 21000-21099 | `spec/19-main-worker-service/13-error-codes.md` | ✅ Active |
| `MWS` | Main/Worker Service (Main tier) | 21100-21199 | `spec/19-main-worker-service/13-error-codes.md` | ✅ Active |

---

## Collision Resolution Log

> This section documents all error code collisions resolved in Wave 1 remediation (2026-02-07).
> 📋 **Summary report:** [03-collision-resolution-summary.md](./03-collision-resolution-summary.md) — consolidated before/after table of all 13 resolutions.

### Resolution 1: AI Bridge Lovable Reasoning vs WebSocket Resilience (Phase 14, CRIT-07)

**Collision:** `42-lovable-reasoning-defaults.md` defined codes 9830-9835 for Lovable Reasoning, but 9830-9839 was already assigned to WebSocket Connection Manager.

**Resolution:** Lovable Reasoning reassigned to **10500-10519** (new range outside existing AB allocation). WebSocket Resilience retains 9830-9839.

### Resolution 2: Nexus Flow Reset API vs State Errors (Phase 15, W-12)

**Collision:** `09-reset-api.md` defined NF-8301 through NF-8308, but `04-error-codes.md` defines 8301-8303 as "State Errors".

**Resolution:** Reset API errors reassigned to **8350-8369** (within NF's 8000-8399 range, previously unallocated). State Errors retain 8301-8303.

### Resolution 3: Link Manager vs AI Transcribe (Phase 16, Finding 1.1)

**Collision:** Link Manager spec claimed 14000-14999, directly overlapping AI Transcribe's canonical 14000-14499.

**Resolution:** Link Manager reassigned to **15000-15999**. Legacy LM range 3000-3999 deprecated. AI Transcribe retains 14000-14499.

### Resolution 4: AI Transcribe Internal — Voice Commands vs Model Download (Phase 17, Finding 1.1)

**Collision:** Both Voice Commands and Model Download used codes starting at 14200.

**Resolution:** Model Download errors reassigned to **14470-14489** (within AIT's 14000-14499 range). Voice Commands retain 14200-14249.

### Resolution 5: WP Plugin Publish Local Codes (Phase 16, Finding 1.3)

**Issue:** WP Plugin Publish used local `E{x}xxx` codes disconnected from ecosystem registry.

**Resolution:** Assigned new prefix `WPP` with range **13000-13999**. All `E{x}xxx` codes must be converted to `13xxx` integers.

### Resolution 6: WP Plugin Builder Range Compressed (Revised 2026-02-28)

**Issue:** WPB was originally registered at 10000-10999. Resolution 1 moved Lovable Reasoning to 10500-10519, and Resolution 6 narrowed WPB to 10000-10499. WPB's own specs used codes in 10500-10899 (Code Generation, Spec Processing, Server/API, Settings, Reset).

**Resolution:** All WPB 10500-10899 codes compressed into 10000-10499 using 20-slot sub-ranges (10420-10439 Code Gen, 10440-10459 Spec Processing, 10460-10479 Server/API, 10480-10489 Settings, 10490-10499 Reset). AB Lovable Reasoning reassigned to **19000-19019**. See Resolution 13.

### Resolution 7: SM Code Generation vs WP SEO Publish (2026-02-28)

**Collision:** SM Code Generation System used 12000-12799, directly overlapping WSP (WP SEO Publish CLI, 12000-12599).

**Resolution:** SM Code Generation reassigned to **16000-16799**. WSP retains 12000-12599. All 12xxx codes in `spec/02-spec-management-software/05-features/24-code-generation-system/16-error-codes.md` remapped to 16xxx equivalents.

### Resolution 8: SM Project Editor vs WP Plugin Publish (2026-02-28)

**Collision:** SM Project Editor used 13000-13999, directly overlapping WPP (WP Plugin Publish, 13000-13999).

**Resolution:** SM Project Editor reassigned to **17000-17999**. WPP retains 13000-13999. All 13xxx codes in `spec/02-spec-management-software/05-features/28-project-editor/05-error-codes.md` remapped to 17xxx equivalents.

### Resolution 9: PS/AB SEO 9500 Range Overlap (2026-02-28, Documented)

**Overlap:** PowerShell Integration (PS, 9500-9599) and AI Bridge SEO (AB, 9500-9540) share the 9500-9540 sub-range.

**Resolution:** This is a **known intentional overlap** that does not cause runtime collisions due to format separation:
- **PS codes** use prefixed format: `PS-9500-00`, `PS-9501-01`, etc. (3-segment, string-based)
- **AB SEO codes** use flat integer format: `9501`, `9502`, etc. (Go constants, integer-based)

The two formats are distinguishable at parse time by their encoding. No reassignment is required. Both modules document their codes independently:
- PS: `spec/11-powershell-integration/04-error-codes.md`
- AB SEO: `spec/11-ai-bridge-cli/01-backend/16-ai-seo-error-codes.md`

**Contingency:** If future ambiguity arises (e.g., a unified logging system that strips prefixes), AB SEO should migrate to **9541-9599** (currently reserved for SEO expansion).

### Resolution 10: AIT Voice Codes vs WPP (2026-02-28)

**Collision:** AI Transcribe CLI voice-related specs used 13200-13308 (voice commands, voice cloning, TTS providers), directly overlapping WPP (WP Plugin Publish, 13000-13999).

**Resolution:** All AIT voice codes reassigned to their canonical 14xxx sub-ranges per the AIT allocation map:
- Voice Commands: 13200-13206 → **14200-14206**
- Voice Cloning: 13250-13257 → **14300-14307**
- TTS Providers: 13300-13308 → **14150-14158**

WPP retains 13000-13999.

### Resolution 11: SM Realtime vs WSP (2026-02-28)

**Collision:** SM Realtime (feature 18) used 12001-12031, directly overlapping WSP (WP SEO Publish CLI, 12000-12599).

**Resolution:** SM Realtime reassigned to **SM-RT 2800-2849** (within SM's base 2000-2999 range). WSP retains 12000-12599.

### Resolution 12: SM GSearch CLI vs Multiple Ranges (2026-02-28)

**Collision:** SM GSearch CLI (feature 22) used local 1xxx-12xxx codes across 12 domains, colliding with GEN (1xxx), SM (2xxx), WSP (12xxx), and others.

**Resolution:** All 92 GSearch CLI codes reassigned to **SM-GS 18000-18249** with 20-slot sub-ranges per domain. Full migration table in `spec/02-spec-management-software/05-features/22-golang-search-cli/15-error-codes.md`.

### Resolution 13: AB Lovable Reasoning vs WPB (2026-02-28)

**Collision:** AB Lovable Reasoning occupied 10500-10519 (per Resolution 1), within WPB's original 10000-10999 range. Resolution 6 narrowed WPB to 10000-10499 but WPB specs used 10500-10899.

**Resolution:** WPB compressed into **10000-10499** (all 10500-10899 codes remapped to 10420-10499). AB Lovable Reasoning reassigned from 10500-10519 to **19000-19019** (new dedicated range). Full migration table in WPB specs.

---

## Standalone Specification Error Ranges

| Module | Prefix | Range | Frontend Sub-range | Error Codes Doc |
|--------|--------|-------|-------------------|-----------------|
| GSearch CLI Core | `GS` | 7000-7099 | 7050-7069 | `spec/09-gsearch-cli/01-backend/15-error-codes.md` |
| BRun CLI | `BR` | 7100-7599 | 7150-7169 | `spec/10-brun-cli/01-backend/06-error-handling.md` |
| GSearch Movie Search | `GS` | 7600-7609 | N/A | `spec/09-gsearch-cli/01-backend/` |
| GSearch BI Suite | `GS` | 7700-7839 | 7800-7819 | `spec/09-gsearch-cli/01-backend/openapi-bi-suite.yaml` |
| GSearch Multi-Source | `GS` | 7840-7859 | N/A | `spec/09-gsearch-cli/01-backend/` |
| GSearch Scheduled | `GS` | 7860-7879 | N/A | `spec/09-gsearch-cli/01-backend/` |
| GSearch Chrome Extension | `GS` | 7880-7899 | N/A | `spec/09-gsearch-cli/01-backend/` |
| GSearch Enum Architecture | `GS` | 7900-7919 | N/A | `spec/09-gsearch-cli/01-backend/` |
| GSearch Provider Integration | `GS` | 7920-7949 | N/A | `spec/09-gsearch-cli/01-backend/` |
| Nexus Flow | `NF` | 8000-8399 | 8050-8069 | `spec/12-nexus-flow-cli/01-backend/04-error-codes.md` |
| AI Bridge Core | `AB` | 9000-9499 | 9050-9069 | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge SEO | `AB` | 9500-9540 | N/A | `spec/11-ai-bridge-cli/01-backend/16-ai-seo-error-codes.md` |
| AI Bridge Extended | `AB` | 9600-9699 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Revisions/Suggestions | `AB` | 9700-9749 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge RAG Session Memory | `AB` | 9750-9809 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Adaptive Reasoning | `AB` | 9810-9829 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge WebSocket Resilience | `AB` | 9830-9839 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Context Integration | `AB` | 9840-9847 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Research Mode | `AB` | 9848-9849 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Code Pattern Learning | `AB` | 9850-9869 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Plan Generation | `AB` | 9870-9889 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Plan Synchronization | `AB` | 9890-9909 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Plan Templates | `AB` | 9910-9929 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Execution Monitoring | `AB` | 9930-9949 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Retry Strategies | `AB` | 9950-9969 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Long-Chain Commands | `AB` | 9970-9989 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Vector DB | `AB` | 9990-9999 | N/A | `spec/11-ai-bridge-cli/01-backend/05-error-codes.md` |
| AI Bridge Lovable Reasoning | `AB` | 19000-19019 | N/A | `spec/11-ai-bridge-cli/01-backend/42-lovable-reasoning-defaults.md` |
| WP Plugin Builder | `WPB` | 10000-10499 | N/A | `spec/14-wp-plugin-builder/10-error-handling.md` |
| Spec Reverse CLI | `SRC` | 11000-11999 | N/A | `spec/15-spec-reverse-cli/01-backend/05-error-codes.md` |
| WP SEO Publish | `WSP` | 12000-12599 | N/A | `spec/21-wp-seo-publish-cli/01-backend/` |
| WP Plugin Publish | `WPP` | 13000-13999 | N/A | `spec/13-wp-plugin/05-wp-plugin-publish/66-shared-constants.md` |
| AI Transcribe CLI | `AIT` | 14000-14499 | N/A | `spec/16-ai-transcribe-cli/01-backend/10-error-codes.md` |
| Exam Manager | `EQM` | 14500-14999 | N/A | `spec/13-wp-plugin/03-exam-manager/66-shared-constants.md` |
| Link Manager | `LM` | 15000-15999 | N/A | `spec/13-wp-plugin/04-link-manager/` |
| AI Bridge Non-Vector RAG | `AB-TR` | 20000-20999 | N/A | `spec/33-ai-bridge-non-vector-rag/08-error-codes.md` |

### AI Bridge Error Code Details (9000-10519)

| Range | Category | Description |
|-------|----------|-------------|
| 9000-9049 | Backend Core | Go backend errors |
| 9050-9069 | Frontend | React frontend errors |
| 9100-9199 | Provider | AI provider connection errors |
| 9200-9299 | Streaming | SSE/WebSocket streaming errors |
| 9301-9310 | RAG Validation | RAG chunk config validation errors |
| 9311-9319 | Request Processing | Request validation errors |
| 9320-9339 | Model & Generation | Model loading, generation errors |
| 9400-9499 | Rate Limiting | Rate limit and quota errors |
| 9500-9540 | AI SEO Generate | SEO module errors |
| 9541-9599 | (Reserved) | Future SEO expansion |
| 9600-9699 | Extended Core | Advanced backend features |
| 9700-9749 | Revisions/Suggestions | Content revision system |
| 9750-9809 | RAG Session Memory | Session-scoped RAG errors |
| 9810-9829 | Adaptive Reasoning | Reasoning mode errors |
| 9830-9839 | WebSocket Resilience | Connection manager errors |
| 9840-9847 | Context Integration | GSearch context injection |
| 9848-9849 | Research Mode | Research delegation errors |
| 9850-9869 | Code Pattern Learning | Pattern detection/enforcement |
| 9870-9889 | Plan Generation | Plan creation errors |
| 9890-9909 | Plan Synchronization | Plan sync/conflict errors |
| 9910-9929 | Plan Templates | Template management errors |
| 9930-9949 | Execution Monitoring | Execution tracking errors |
| 9950-9969 | Retry Strategies | Retry/backoff errors |
| 9970-9989 | Long-Chain Commands | Multi-step command errors |
| 9990-9999 | Vector DB Integration | Vector search errors |
| 19000-19019 | Lovable Reasoning | Reasoning defaults/questions |

### AI Bridge Non-Vector RAG Error Code Details (20000-20999)

| Range | Category | Description |
|-------|----------|-------------|
| 20000-20003 | General | Startup, config, database, model availability |
| 20100-20108 | Code Parsing | AST/regex parsing, file read, encoding, timeout |
| 20200-20206 | Document Parsing | Markdown, HTML, CSV, frontmatter parsing |
| 20300-20306 | Tree Construction | LLM enrichment, batch processing, caching |
| 20400-20404 | Tree Storage | SQLite writes, FTS5, migrations, capacity |
| 20500-20507 | Retrieval | Query analysis, traversal, scoring, context assembly |
| 20600-20604 | API | Index/job lookup, validation, rate limiting |
| 20700-20702 | Configuration | Config file loading and validation |
| 20800-20802 | AI Bridge Integration | Bridge communication, router, session |
| 20803-20999 | (Reserved) | Future expansion |

### Nexus Flow Error Code Details (8000-8399)

| Range | Category | Description |
|-------|----------|-------------|
| 8000-8049 | Backend Core | Go backend errors |
| 8050-8069 | Frontend | React frontend errors |
| 8100-8199 | Pipeline | Pipeline execution errors |
| 8200-8299 | Block | Block/node errors |
| 8300-8303 | State | State transition errors |
| 8304-8349 | (Available) | Unallocated |
| 8350-8369 | Reset API | Reset/teardown errors |
| 8370-8399 | (Reserved) | Future expansion |

### AI Transcribe Error Code Details (14000-14499)

| Range | Category | Description |
|-------|----------|-------------|
| 14000-14049 | General | Initialization, config errors |
| 14050-14099 | Audio Pipeline | Audio processing errors |
| 14100-14149 | STT Providers | Speech-to-text provider errors |
| 14150-14199 | TTS Providers | Text-to-speech provider errors |
| 14200-14249 | Voice Commands | Voice command parsing errors |
| 14250-14299 | Real-time | WebSocket conversation errors |
| 14300-14349 | Voice Cloning | Voice cloning errors |
| 14350-14399 | (Reserved) | Future expansion |
| 14400-14449 | (Reserved) | Future expansion |
| 14450-14469 | Provider Registry | Provider management errors |
| 14470-14489 | Model Download | Model download/verification errors |
| 14490-14499 | (Reserved) | Future expansion |

---

## RAG Configuration Validation Errors (9301-9310)

| Code | Name | Message |
|------|------|---------|
| AB-9301 | `ErrRagChunkSizeInvalid` | Chunk size outside valid range (256-8192) |
| AB-9302 | `ErrRagChunkSizeNotMultiple` | Chunk size not multiple of 256 |
| AB-9303 | `ErrRagOverlapTooLarge` | Chunk overlap exceeds 25% of chunk size |
| AB-9304 | `ErrRagContextBudgetInvalid` | Context token budget outside range (512-16384) |
| AB-9305 | `ErrRagEmbeddingModelInvalid` | Embedding model not supported |
| AB-9306 | `ErrRagSimilarityThresholdInvalid` | Similarity threshold outside range (0.0-1.0) |
| AB-9307 | `ErrRagTopkInvalid` | TopK outside valid range (1-50) |
| AB-9308 | `ErrRagConfigLoadFailed` | Failed to load RAG configuration |
| AB-9309 | `ErrRagConfigSaveFailed` | Failed to save RAG configuration |
| AB-9310 | `ErrRagConfigSourceConflict` | Conflicting config from multiple sources |

## AI SEO Generate Error Codes (9500-9524)

| Code | Name | Message |
|------|------|---------|
| AB-9501 | `ErrSeoPresetNotFound` | Industry preset not found |
| AB-9502 | `ErrSeoTemplateNotFound` | Template file not found in preset |
| AB-9503 | `ErrSeoVariableMissing` | Required template variable not provided |
| AB-9504 | `ErrSeoGenerationFailed` | Page generation failed |
| AB-9505 | `ErrSeoContextInjectionFailed` | Failed to inject RAG context |
| AB-9506 | `ErrSeoLlmTimeout` | LLM response timeout during generation |
| AB-9507 | `ErrSeoTokenLimitExceeded` | Content exceeds token limit |
| AB-9508 | `ErrSeoOutputWriteFailed` | Failed to write generated output |
| AB-9509 | `ErrSeoBatchLimitExceeded` | Batch exceeds maximum page limit |
| AB-9510 | `ErrSeoUploadFailed` | ZIP upload failed |
| AB-9511 | `ErrSeoUploadTooLarge` | Upload exceeds size limit |
| AB-9512 | `ErrSeoUnsupportedFormat` | File format not supported |
| AB-9513 | `ErrSeoCircularDependency` | Circular dependency detected in files |
| AB-9514 | `ErrSeoParseMarkdownFailed` | Failed to parse Markdown file |
| AB-9515 | `ErrSeoParseHtmlFailed` | Failed to parse HTML template |
| AB-9516 | `ErrSeoParseCsvFailed` | Failed to parse CSV file |
| AB-9517 | `ErrSeoParseSqliteFailed` | Failed to read SQLite database |
| AB-9518 | `ErrSeoExtractionFailed` | ZIP extraction failed |
| AB-9519 | `ErrSeoIngestFailed` | Failed to ingest files into RAG |
| AB-9520 | `ErrSeoJobNotFound` | Generation job not found |
| AB-9521 | `ErrSeoJobAlreadyRunning` | Job already in progress |

## GSearch Movie Search Error Codes (7600-7609)

| Code | Name | Message |
|------|------|---------|
| GS-7600 | `ErrMovieSearchFailed` | Movie search request failed |
| GS-7601 | `ErrMovieNotFound` | Movie not found in database |
| GS-7602 | `ErrTmdbApiError` | TMDB API request failed |
| GS-7603 | `ErrOmdbApiError` | OMDB API request failed |
| GS-7604 | `ErrMovieCacheError` | Movie cache read/write failed |
| GS-7605 | `ErrTvEpisodeNotFound` | TV episode not found |
| GS-7606 | `ErrFolderScanError` | Folder batch scan failed |
| GS-7607 | `ErrFilenameParsError` | Failed to parse filename for metadata |
| GS-7608 | `ErrApiKeyRotationFailed` | API key rotation exhausted |
| GS-7609 | `ErrMetadataTtlExpired` | Cached metadata expired (30-day TTL) |

## GSearch Business Intelligence Suite Error Codes (7700-7839)

| Range | Category | Description |
|-------|----------|-------------|
| 7700-7709 | Multi-Engine Search | Search engine errors |
| 7710-7719 | FAQ Discovery | FAQ extraction errors |
| 7720-7729 | SERP Tracking | Position tracking errors |
| 7730-7739 | Contact Extraction | Contact parsing errors |
| 7740-7749 | Maps Search | Google Maps errors |
| 7750-7759 | Response Formatting | Cache and format errors |
| 7760-7769 | Webhook | Webhook notification errors |
| 7800-7819 | Frontend | BI Suite frontend errors |
| 7820-7839 | Reserved | Future expansion |

### Frontend Error Code Pattern

All CLI frontends use a consistent error code pattern at offset +50 from their base range:

| Offset | Error | Description |
|--------|-------|-------------|
| +50 | ErrWsConnectionFailed | WebSocket connection failure |
| +51 | ErrWsDisconnected | WebSocket unexpectedly closed |
| +52 | ErrSettingsLoadFailed | Failed to load settings |
| +53 | ErrSettingsSaveFailed | Failed to save settings |
| +54 | ErrApiTimeout | API request timeout |
| +55 | ErrApiError | API returned error response |
| +56 | ErrConfigParseError | Failed to parse config |
| +57 | ErrVersionMismatch | Frontend/backend version mismatch |
| +58 | ErrPortUnavailable | Configured port not available |
| +59 | ErrFirewallBlocked | Firewall blocking connection |

---

## GEN: General/Shared Errors (1000-1999)

### GEN-000: Initialization

| Code | Name | Message |
|------|------|---------|
| GEN-000-01 | ErrConfigMissing | Configuration file not found |
| GEN-000-02 | ErrConfigInvalid | Configuration file is malformed |
| GEN-000-03 | ErrEnvMissing | Required environment variable not set |

### GEN-100: Authentication

| Code | Name | Message |
|------|------|---------|
| GEN-100-01 | ErrAuthRequired | Authentication required |
| GEN-100-02 | ErrTokenExpired | Authentication token has expired |
| GEN-100-03 | ErrTokenInvalid | Authentication token is invalid |
| GEN-100-04 | ErrCredentialsInvalid | Invalid username or password |

### GEN-200: Authorization

| Code | Name | Message |
|------|------|---------|
| GEN-200-01 | ErrAccessDenied | Access denied to this resource |
| GEN-200-02 | ErrRoleRequired | Insufficient role privileges |
| GEN-200-03 | ErrPermissionDenied | Permission not granted |

### GEN-300: Validation

| Code | Name | Message |
|------|------|---------|
| GEN-300-01 | ErrFieldRequired | Required field is missing |
| GEN-300-02 | ErrFieldInvalid | Field value is invalid |
| GEN-300-03 | ErrFormatInvalid | Input format is invalid |
| GEN-300-04 | ErrLengthExceeded | Input exceeds maximum length |

### GEN-400: Business Logic

| Code | Name | Message |
|------|------|---------|
| GEN-400-01 | ErrOperationFailed | Business operation failed |
| GEN-400-02 | ErrStateInvalid | Invalid state for requested operation |
| GEN-400-03 | ErrConflict | Operation conflicts with current state |
| GEN-400-04 | ErrLimitExceeded | Operation limit exceeded |

### GEN-500: Database

| Code | Name | Message |
|------|------|---------|
| GEN-500-01 | ErrDbConnection | Database connection failed |
| GEN-500-02 | ErrDbQuery | Database query failed |
| GEN-500-03 | ErrDbTransaction | Transaction failed |
| GEN-500-04 | ErrRecordNotFound | Record not found |
| GEN-500-05 | ErrDuplicateRecord | Record already exists |

### GEN-600: Type Casting / Conversion

> Cross-cutting errors used by `pkg/typecast/`. See [Casting Elimination Patterns](../../02-coding-guidelines/01-cross-language/03-casting-elimination-patterns.md) for usage.

| Code | Name | Message |
|------|------|---------|
| GEN-600-01 | `ErrCastTypeAssertionFailed` | Type assertion failed: expected `T`, got `U` |
| GEN-600-02 | `ErrCastSliceElementFailed` | Slice element cast failed at index `N`: expected `T`, got `U` |

**Implementation mapping:**

| Spec Code | Go Constant | Used By |
|-----------|-------------|---------|
| `GEN-600-01` | `ErrCastTypeAssertionFailed` | `typecast.CastOrFail[T]()` |
| `GEN-600-02` | `ErrCastSliceElementFailed` | `typecast.CastSliceOrFail[T]()` |

**Rules:**
- These codes are emitted exclusively by `pkg/typecast/` — never constructed manually
- The `AppError` includes `.WithSkip(1)` so stack traces point to the caller
- Cast errors must **never** be swallowed — see §7.2 and §10 in [03-casting-elimination-patterns.md](../../02-coding-guidelines/01-cross-language/03-casting-elimination-patterns.md)


### GEN-700: File System

| Code | Name | Message |
|------|------|---------|
| GEN-700-01 | ErrFileNotFound | File not found at specified path |
| GEN-700-02 | ErrFileReadFailed | Failed to read file |
| GEN-700-03 | ErrFileWriteFailed | Failed to write file |
| GEN-700-04 | ErrDirCreateFailed | Failed to create directory |
| GEN-700-05 | ErrFsPermissionDenied | Insufficient file system permissions |

### GEN-800: Network

| Code | Name | Message |
|------|------|---------|
| GEN-800-01 | ErrNetworkError | Network request failed |
| GEN-800-02 | ErrTimeout | Request timed out |
| GEN-800-03 | ErrServiceUnavailable | Service temporarily unavailable |

### GEN-900: Reserved

> Reserved for future cross-cutting error categories. No codes should be registered in this range until a new category is formally defined and documented.

| Code | Name | Message |
|------|------|---------|
| *(none allocated)* | — | Reserved for future use |

---

## SM: Spec Management Software (2000-2999)

### SM-000: Initialization

| Code | Name | Message |
|------|------|---------|
| SM-000-01 | ErrSpecRootMissing | Spec root directory not found |
| SM-000-02 | ErrIndexMissing | Master index file not found |

### SM-400: Business Logic

| Code | Name | Message |
|------|------|---------|
| SM-400-01 | ErrSpecParseError | Failed to parse specification file |
| SM-400-02 | ErrCircularDependency | Circular dependency detected |
| SM-400-03 | ErrVersionConflict | Version conflict detected |
| SM-400-04 | ErrTemplateError | Template rendering failed |

### SM-500: Database

| Code | Name | Message |
|------|------|---------|
| SM-500-01 | ErrMigrationFailed | Database migration failed |
| SM-500-02 | ErrSeedFailed | Database seeding failed |

### SM-600: External Services

| Code | Name | Message |
|------|------|---------|
| SM-600-01 | ErrAiApiError | AI service request failed |
| SM-600-02 | ErrEmbeddingFailed | Vector embedding generation failed |
| SM-600-03 | ErrRagSearchFailed | RAG search query failed |

---

## LM: Link Manager (15000-15999)

> ⚠️ **REASSIGNED:** Previously at 3000-3999 and 14000-14999 (both deprecated). Moved to 15000-15999 to resolve collision with AI Transcribe CLI (14000-14499).

### LM-15000: Initialization

| Code | Name | Message |
|------|------|---------|
| LM-15000-01 | ErrWpNotDetected | WordPress environment not detected |
| LM-15000-02 | ErrPluginConflict | Plugin conflict detected |

### LM-15400: Business Logic

| Code | Name | Message |
|------|------|---------|
| LM-15400-01 | ErrLinkInvalid | Invalid link format |
| LM-15400-02 | ErrRedirectLoop | Redirect loop detected |
| LM-15400-03 | ErrDomainBlocked | Domain is blocked |

---

## CLI: CLI Tools (4000-4999) — DEPRECATED

> ⚠️ **DEPRECATED:** Legacy range. Do not allocate new codes. GSearch and BRun now use 7000+ ranges.

### CLI-000: Initialization

| Code | Name | Message |
|------|------|---------|
| CLI-000-01 | ErrBinaryNotFound | Required binary not found |
| CLI-000-02 | ErrPathNotSet | PATH environment not configured |

### CLI-400: gsearch (legacy)

| Code | Name | Message |
|------|------|---------|
| CLI-400-01 | ErrPatternInvalid | Invalid search pattern |
| CLI-400-02 | ErrNoResults | No results found |
| CLI-400-03 | ErrIndexStale | Search index is stale |

### CLI-500: brun (legacy)

| Code | Name | Message |
|------|------|---------|
| CLI-500-01 | ErrBuildFailed | Build step failed |
| CLI-500-02 | ErrRunFailed | Run step failed |
| CLI-500-03 | ErrDepsMissing | Dependencies not installed |

---

## PS: PowerShell Integration (9500-9599)

> See `spec/11-powershell-integration/04-error-codes.md` for full list.

| Code | Name | Message |
|------|------|---------|
| PS-9500-00 | PsSuccess | Operation completed successfully |
| PS-9501-01 | ErrPsGoMissing | Go runtime not found |
| PS-9502-01 | ErrPsNodeMissing | Node.js not found |
| PS-9510-01 | ErrPsConfigMissing | powershell.json not found |
| PS-9520-01 | ErrPsBuildFailed | Frontend build failed |
| PS-9530-01 | ErrPsBackendFailed | Backend failed to start |

---

## Error Code Format Reference

The ecosystem uses **two** error code formats:

| Format | Used By | Example | Pattern |
|--------|---------|---------|---------|
| `XX-NNN-NN` | General specs, PHP plugins | `SM-400-01` | `^[A-Z]{2,4}-[0-9]{3}-[0-9]{2}$` |
| Integer | Go CLI tools | `7001`, `9301` | `^[0-9]{4,5}$` |

Both formats are valid. Go CLI tools use flat integers in API responses, constants, and logs. PHP plugins and general specs use the prefixed format. The JSON schema supports both.

---

## Range Allocation Map

```
1000-1999  GEN (General/Shared)
2000-2999  SM  (Spec Management)
3000-3999  LM  [DEPRECATED - moved to 15000]
4000-4999  CLI [DEPRECATED]
5000-6999  --- [UNALLOCATED]
7000-7099  GS  (GSearch Core)
7100-7599  BR  (BRun)
7600-7609  GS  (Movie Search)
7610-7699  --- [UNALLOCATED]
7700-7839  GS  (BI Suite)
7840-7859  GS  (Multi-Source)
7860-7879  GS  (Scheduled)
7880-7899  GS  (Chrome Extension)
7900-7919  GS  (Enum Architecture)
7920-7949  GS  (Provider Integration)
7950-7999  --- [UNALLOCATED]
8000-8399  NF  (Nexus Flow)
8400-8999  --- [UNALLOCATED]
9000-9499  AB  (AI Bridge Core)
9500-9599  PS  (PowerShell) / AB (AI SEO 9500-9540)
9600-9999  AB  (AI Bridge Extended)
10000-10499 WPB (WP Plugin Builder) ✅ Compressed from 10000-10999
10500-10999 --- [UNALLOCATED]
11000-11999 SRC (Spec Reverse)
12000-12599 WSP (WP SEO Publish)
12600-12999 --- [UNALLOCATED]
13000-13999 WPP (WP Plugin Publish)
14000-14499 AIT (AI Transcribe)
14500-14999 EQM (Exam Manager)
15000-15999 LM  (Link Manager)
16000-16799 SM-CG (SM Code Generation) ✅ Reassigned from 12xxx
16800-16999 --- [UNALLOCATED]
17000-17999 SM-PE (SM Project Editor) ✅ Reassigned from 13xxx
18000-18249 SM-GS (SM GSearch CLI) ✅ Reassigned from 1xxx-12xxx
18250-18999 --- [UNALLOCATED]
19000-19019 AB  (Lovable Reasoning) ✅ Reassigned from 10500-10519
19020-19999 --- [UNALLOCATED]
20000-20999 AB-TR (Non-Vector RAG) ✅ Reassigned from 12000-12999
20999+      --- [UNALLOCATED]
```

---

## Adding New Codes

1. Check the Range Allocation Map above for available ranges
2. Claim a project prefix in the Registered Project Prefixes table
3. Use the next available number in the category
4. Follow naming convention: `ErrPascalCase` (PascalCase with `Err` prefix)
5. Provide clear, actionable message
6. Update this registry in the same commit
7. Never reuse deprecated or retired codes
8. If the code is a Go domain error, also register it as an `apperrtype` enum — see [05-apperrtype-enums.md](../02-error-architecture/06-apperror-package/01-apperror-reference/05-apperrtype-enums.md)

---

## Relationship to `apperrtype` Enums

The `apperror` package uses a separate `E{x}xxx` string code format for Go-internal domain errors (e.g., `E2010` = site not found). These do **not** collide with this registry's prefixed format (`GEN-000-01`, `AB-9301`, etc.) because the formats are distinguishable at parse time.

For the full `apperrtype` enum reference and collision analysis, see:
→ [apperrtype Domain Error Enums](../02-error-architecture/06-apperror-package/01-apperror-reference/05-apperrtype-enums.md)
