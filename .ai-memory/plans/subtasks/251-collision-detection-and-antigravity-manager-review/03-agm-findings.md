# Worker 03 Report — Antigravity-Manager Investigation (Task 251, Subtask 03)

**Date:** 2026-10-09
**Scope:** READ-ONLY investigation. No AGM files modified. No merges performed.
**Repo:** `~/workspace/repos/Antigravity-Manager` (fork: `alimtvnetwork/Antigravity-Manager`, upstream: `lbjlaq/Antigravity-Manager`)

---

## 1. Attribution: Who is "Jeik"?

**Answer: Jeik is an UPSTREAM maintainer/contributor, NOT a fork contributor.**

### Evidence

**Git log on openai.rs (local main):**
| Author | Commits |
|--------|---------|
| lbjlaq <469054080@qq.com> | 38 |
| Jeik <jeikliu@outlook.com> | 22 |
| Jeikl <jeikliu@outlook.com> | 1 (same email = Jeik, total 23) |
| new-Beginner <3205732937@qq.com> | 15 |
| buluw <ocgbuluw@gmail.com> | 10 |
| 8 others | 1-3 each |

**Upstream/main contains Jeik's commits:**
```
$ git log upstream/main -- src-tauri/src/proxy/handlers/openai.rs
Jeik <jeikliu@outlook.com> d0fccd32 fix(proxy): 修复 3.1-flash-lite 误重定向...
Jeik <jeikliu@outlook.com> d00085dd feat(proxy): 实现纯思考空回复流式自愈门禁 (Pipeline First)
Jeik <jeikliu@outlook.com> b57872fc feat(proxy): 增加模型配置面板与协议无关的多模态保鲜滑动窗口...
Jeik <jeikliu@outlook.com> dbbfb073 fix(proxy): drop placeholder thinking blocks, keep anchor signatures
... (same commits as local main)
```

**GitHub API — upstream contributors (`lbjlaq/Antigravity-Manager`):**
| GitHub User | Contributions | Role |
|-------------|---------------|------|
| `lbjlaq` (id 22748003) | 671 | Repo owner, #1 |
| `jeikl` (id 100782431) | 251 | **#2 contributor** |
| `new-Beginner` | 31 | #3 |
| `buluw` | 21 | #4 |

**Conclusion:** "Jeik" (jeikliu@outlook.com) = GitHub user `jeikl`. With 251 upstream contributions (#2 after the repo owner), Jeik is a core upstream maintainer. The fork (`alimtvnetwork`) syncs FROM upstream, which is why Jeik's commits appear in both trees. The 7,583-line openai.rs was authored primarily UPSTREAM, not by the fork.

**Implication for the "leakage" question:** The pipeline/adapter boundary violation in openai.rs is an UPSTREAM design outcome, not something the fork introduced. Any rescue must either be contributed upstream or maintained as a fork divergence.

---

## 2. openai.rs: 7,583-Line Concern Breakdown

**Total:** 7,583 lines, 56 functions, 2 public (`get_cached_tool_call`, `insert_cached_tool_call`).

### By line range

| Lines | Size | Concern |
|-------|------|---------|
| 1–920 | ~920 (12%) | Image utilities + Responses API content building |
| 923–3,060 | ~2,140 (28%) | Codex transcript/ledger + Responses routing |
| 3,060–4,953 | ~1,890 (25%) | Main `handle_completions` request handler |
| 4,953–5,591 | ~640 (8%) | Image generations/edits endpoints |
| 5,591–6,065 | ~475 (6%) | Misc helpers |
| 6,065–7,583 | ~1,520 (20%) | Codex WebSocket handler + global tool-call cache |

### By keyword (rough `aum search` counts)

| Concern | Matching lines |
|---------|---------------|
| image/Image | 357 |
| stream/Stream/SSE/chunk | 183 |
| token/Token | 146 |
| tool_call/function_call | 111 |
| responses/Responses | 113 |
| thinking/Thinking | 75 |
| signature/Signature | 48 |
| audio/Audio | 14 |

### Key image functions (lines 46–750, ~700 lines — the extraction target)

**Parsing & validation:**
- `validate_input_image_limits` (46), `normalized_image_from_bytes` (72)
- `parse_image_data_url` (87), `parse_image_data_url_parts` (107)
- `parse_generation_input_images` (126), `generation_image_size_param` (159)
- `validate_responses_image_data_url` (505), `validate_responses_input_image_limits` (528)
- `decoded_base64_len` (492)

**Content building:**
- `build_image_contents` (208), `build_image_edit_body` (231)
- `image_inline_part` (199), `is_edit_image_field` (175), `edit_size_input` (183)
- `image_account_selection_target` (195)

**Media history management:**
- `historical_media_placeholder` (574), `history_without_inline_media` (588)
- `into_history_without_inline_media` (616), `omit_media_before_latest_user_turn` (647)
- `build_responses_tool_output_content` (680), `debug_value_without_inline_data` (693)

**Media detection:**
- `response_has_inline_image_data` (288), `text_has_nonempty_image_data_url` (312)
- `value_has_nonempty_image_data_url` (328), `stream_chunk_has_image_data` (337)

**What "7.5K lines means":** Only ~25% (1,890 lines) is the actual OpenAI protocol adapter (`handle_completions`). The remaining 75% is: image processing that should be pipeline-generic (20%), Codex-specific logic that should be its own module (28%), and WebSocket infrastructure that should be in `proxy/upstream/` (20%). The file violates the repo's own "Pipeline First" rule — it contains three separate subsystems masquerading as one adapter.

---

## 3. PR #6 "Sync v5" Review

**Status:** OPEN, never merged. Created 2026-09-30 by `alimtvnetwork`.
**Scope:** 140 files, +31,206 / −9,802. 100 commits from `lbjlaq:main` → `alimtvnetwork:main`.
**Critical finding:** Our fork does NOT have 3 new pipeline files from this PR (`auto_heal.rs`, `estimator.rs`, `official_alignment_tests.rs` — all MISSING). The fork has diverged; PR #6 would be a major sync.

### Top 10 Pipeline/Adapter Changes

| # | File | Change | What it does | In our fork? | Cherry-pick safety |
|---|------|--------|--------------|--------------|-------------------|
| 1 | `proxy/pipeline/inbound.rs` | +2,241/−372 | Core thinking-budget pipeline rewrite, official alignment | Diverged (different HEADs) | **RISKY** — histories diverged; expect conflicts |
| 2 | `proxy/pipeline/auto_heal.rs` | +587 NEW | Stream self-healing pipeline module | **MISSING** | **SAFE** — new file, no conflicts |
| 3 | `proxy/pipeline/estimator.rs` | +552 NEW | Token/cost estimation pipeline module | **MISSING** | **SAFE** — new file, no conflicts |
| 4 | `proxy/pipeline/official_alignment_tests.rs` | +845 NEW | Protocol alignment test suite | **MISSING** | **SAFE** — new file, no conflicts |
| 5 | `proxy/handlers/claude.rs` | +1,554/−352 | Claude adapter: catalog, 128k output, circuit breaker | Diverged | **RISKY** — large diff on diverged file |
| 6 | `proxy/token_manager.rs` | +1,314/−589 | Token quota/rate-limit management | Unknown | **CAUTION** — check for fork customizations first |
| 7 | `proxy/thinking_store.rs` | +1,065/−470 | Thinking-block persistence layer | Unknown | **CAUTION** — DB schema may differ |
| 8 | `proxy/handlers/openai.rs` | +417/−396 | OpenAI adapter fixes | Diverged | **RISKY** — our fork has 38 lbjlaq + 23 Jeik commits; 3-way merge needed |
| 9 | `proxy/handlers/gemini.rs` | +391/−141 | Gemini adapter updates | Diverged | **RISKY** — same as above |
| 10 | `proxy/common/model_mapping.rs` | +885/−145 | Model name normalization mappings | Unknown | **CAUTION** — may conflict with fork's model list |

### Cherry-pick safety summary

**SAFE (new files, zero conflict risk):**
- `proxy/pipeline/auto_heal.rs` — pure addition
- `proxy/pipeline/estimator.rs` — pure addition
- `proxy/pipeline/official_alignment_tests.rs` — pure addition (but requires the modules it tests)

**DO NOT cherry-pick blindly:**
- `pipeline/inbound.rs` — 2,241 added lines on a diverged file. Our fork's HEAD (`5321c496`) and upstream's HEAD (`04198598`) have different thinking-block logic. A blind cherry-pick will either conflict or silently regress our fork's placeholder-block fix.
- `handlers/openai.rs`, `handlers/claude.rs`, `handlers/gemini.rs` — all diverged. The fork has local fixes (e.g., our `dbbfb073` thinking-block fix) that may not exist upstream in the same form.

**Recommended approach:** Do NOT merge PR #6 as-is (140 files, 100 commits, 3 months stale). Instead:
1. Cherry-pick the 3 SAFE new files first.
2. For `inbound.rs`: do a targeted 3-way diff of the thinking-budget logic only, not the whole file.
3. Leave the handler files alone until the rescue (section 4) is done — merging upstream handler changes into our 7.5k-line file will only entrench the problem.

---

## 4. Rescue Plan: Extract Image Logic to `proxy/pipeline/`

**Goal:** Move ~700 lines of protocol-generic image logic from `openai.rs` (adapter) to `proxy/pipeline/media.rs` (pipeline), per the repo's "Pipeline First" rule. No behavior change.

### New file: `src-tauri/src/proxy/pipeline/media.rs`

**Move these functions verbatim (lines 46–271, ~225 lines):**
| Function | Lines | Notes |
|----------|-------|-------|
| `validate_input_image_limits` | 46–71 | Generic validation |
| `normalized_image_from_bytes` | 72–86 | Generic bytes → normalized |
| `parse_image_data_url` | 87–106 | Generic data-URL parsing |
| `parse_image_data_url_parts` | 107–125 | Generic |
| `parse_generation_input_images` | 126–158 | Generic input parsing |
| `generation_image_size_param` | 159–174 | Generic param extraction |
| `is_edit_image_field` | 175–182 | Generic field check |
| `edit_size_input` | 183–194 | Generic size logic |
| `image_account_selection_target` | 195–198 | Generic routing |
| `image_inline_part` | 199–207 | Generic Gemini part builder |
| `build_image_contents` | 208–230 | Generic contents builder |
| `build_image_edit_body` | 231–271 | Generic edit body builder |
| `decoded_base64_len` | 492–504 | Generic utility |

**Move these functions (lines 505–750, ~245 lines):**
| Function | Lines | Notes |
|----------|-------|-------|
| `validate_responses_image_data_url` | 505–527 | Rename → `validate_image_data_url` (drop "responses" prefix; it's generic) |
| `validate_responses_input_image_limits` | 528–573 | Rename → `validate_input_image_limits_v2` or merge with line 46 version |
| `historical_media_placeholder` | 574–587 | Generic history management |
| `history_without_inline_media` | 588–615 | Generic |
| `into_history_without_inline_media` | 616–646 | Generic |
| `omit_media_before_latest_user_turn` | 647–679 | Generic sliding-window logic |
| `build_responses_tool_output_content` | 680–692 | Rename → `build_tool_output_content` (generic) |
| `debug_value_without_inline_data` | 693–740 | Generic debug utility |

**Move these detection helpers (lines 288–347, ~60 lines):**
| Function | Lines | Notes |
|----------|-------|-------|
| `response_has_inline_image_data` | 288–311 | Generic detection |
| `text_has_nonempty_image_data_url` | 312–327 | Generic |
| `value_has_nonempty_image_data_url` | 328–336 | Generic |
| `stream_chunk_has_image_data` | 337–347 | Generic |

### Steps (no behavior change)

1. **Create** `src-tauri/src/proxy/pipeline/media.rs` with the moved functions. Change `fn` → `pub fn` for all moved functions. Keep signatures identical.
2. **Update** `src-tauri/src/proxy/pipeline/mod.rs`: add `pub mod media;` and re-export as needed.
3. **Update** `src-tauri/src/proxy/handlers/openai.rs`: delete the moved functions (lines 46–271, 288–347, 492–504, 505–750), add `use crate::proxy::pipeline::media::*;` at the top.
4. **Verify:** `cargo check` must pass with zero errors. The moved functions have no OpenAI-specific dependencies (verified: they operate on `serde_json::Value` and byte slices only).
5. **Do NOT move** (OpenAI-specific, stays in handler):
   - `handle_completions` (3,066–4,953) — the actual adapter entry point
   - `handle_images_generations` / `handle_images_edits` (4,991–5,591) — OpenAI Images API endpoints
   - `intercept_chat_to_image` (4,991) — OpenAI chat→image routing
   - `stream_chunk_has_error_event` (272–287) — OpenAI error format specific

### Expected outcome
- `openai.rs`: 7,583 → ~6,850 lines (−730, −10%)
- New `pipeline/media.rs`: ~730 lines, reusable by Claude/Gemini handlers
- Zero behavior change (pure move, identical signatures)

### Follow-up (not in this rescue)
- The Codex block (lines 923–3,060, ~2,140 lines) should become `proxy/handlers/codex.rs` — separate rescue.
- The WebSocket block (lines 6,065–7,583, ~1,520 lines) should move to `proxy/upstream/` — separate rescue.
- After all three extractions, `openai.rs` would be ~2,500 lines (the actual adapter), which is still large but defensible.

---

## Files Referenced (read-only)

- `~/workspace/repos/Antigravity-Manager/src-tauri/src/proxy/handlers/openai.rs` (7,583 lines)
- `~/workspace/repos/Antigravity-Manager/src-tauri/src/proxy/pipeline/` (inbound.rs, mod.rs, policy.rs, usage.rs)
- `~/workspace/repos/Antigravity-Manager/AGENTS.md`
- GitHub: `lbjlaq/Antigravity-Manager` contributors API, PR #6 files API
