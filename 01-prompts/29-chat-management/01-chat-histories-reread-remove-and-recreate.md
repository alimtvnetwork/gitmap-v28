# Chat Histories: Re-read, Remove & Recreate (per-repo chats)

**Version:** 1.0.0
**Date:** 2026-10-10
**Purpose:** Recover the user's chat workspace after a VM replacement — or any drift. Re-read every repository side chat's history and memory, remove (archive) stale or duplicate chats, and guarantee exactly one side chat per repository named `<repo> repo`. Paste into any Muse chat and run.

---

## Operating contract (always in force)

1. **Inventory before action.** Phase 0 is strictly read-only. Never archive, delete, rename, or create a chat until the full inventory map is built and shown.
2. **Destructive actions need per-chat confirmation.** Archiving is the default for stale chats (reversible). `chat.delete` is irreversible — use it only when the user explicitly says "delete" for that specific chat. Never delete the main chat. Never delete anything on your own judgment.
3. **One chat per repo.** Every repository in the working set gets exactly one side chat named exactly `<repo> repo` (e.g. `gitmap-v28 repo`). No prefixes, no suffixes, no variations.
4. **Repo working set = source of truth.** Reconcile `~/workspace/repos/` (directories on disk) with the `Repos` registry in `~/MEMORY.md`. On conflict, the on-disk directories win; note the conflict in the report.
5. **Memory lives beside the chat.** Each side chat's durable memory is `~/side-chats/<chat-id>/MEMORY.md` (plus `memory/` underneath). A chat with no memory dir is "memoryless" — flag it, don't punish it.
6. **Fork rule.** `chat.create` with `context_mode: "fork"` is only allowed from the main chat. When this prompt runs inside a side chat, always use `context_mode: "fresh"`.
7. **Foolproof rule.** Every phase has a check and a fallback. A missing dir, an empty chat list, or a failed read is recorded as a blocker in plain language; everything independent continues. Never halt the whole run on one failure. Never invent chat IDs, repo names, or paths — copy them from tool output.

---

## Phase 0 — Inventory (read-only)

[/goal](slashCommand;goal) Build the complete map before touching anything.

1. `chat.list` → active chats (chat_id, name, is_main).
2. `chat.list` with `include_archived=true` → archived chats.
3. `ls ~/side-chats/` → local memory directories (one per chat_id).
4. `ls ~/workspace/repos/` → repository working set.
5. Read `~/MEMORY.md` → `Repos` registry (expected repo → side-chat mapping).

Produce the inventory table and show it:

| Repo | Expected chat name | chat_id (or MISSING) | Status | Memory dir |
|---|---|---|---|---|
| gitmap-v28 | `gitmap-v28 repo` | … | ok / missing / duplicate / orphan | yes / no |

Also list: chats with no matching repo (candidates for Phase 2), memory dirs with no matching chat (orphans), and repos with no chat (Phase 3 work).

**Done when:** every repo, every active chat, and every memory dir appears exactly once in the map.

## Phase 1 — Re-read (restore context)

[/learn](slashCommand;learn) For each repo chat in the map:

1. Read `~/side-chats/<chat-id>/MEMORY.md` if it exists; then any `memory/*.md` files under it (newest first, cap at 5 files).
2. `chat.read_messages` with `limit=5` on that chat_id to capture recent context (page further only if the summary needs it).
3. Write a 2–3 line brief per chat: what the repo is, what the chat was last working on, any open commitments.

**Done when:** you can state, without guessing, what each repo chat knows and where it left off. Record chats whose history is unreachable (list error, empty transcript) as blockers — do not fabricate their context.

## Phase 2 — Remove (prune stale chats)

[/goal](slashCommand;goal) Propose, then confirm, then act.

1. Build the removal candidate list from Phase 0: duplicate repo chats (keep the one with the richest memory/history, propose the rest), chats for repos no longer in the working set, and orphans.
2. Present the list to the user with one line of reasoning per chat and ask for explicit confirmation per chat: **archive** (default, reversible) or **delete** (only if the user says the word "delete").
3. Execute only what was confirmed: `chat.archive` for archives, `chat.delete` for confirmed deletions.
4. Never archive or delete the main chat. Never act on a chat the user did not confirm.

**Done when:** every confirmed removal is executed and the survivors are listed.

## Phase 3 — Create (one chat per repo)

[/goal](slashCommand;goal) Close the gaps found in Phase 0.

1. For each repo with no side chat: `chat.create` with `context_mode="fresh"`, `name="<repo> repo"`, and a `message` describing the chat's purpose (e.g. `"Dedicated working chat for the <repo> repository: deep repo knowledge, memory, and task execution."`).
2. For chats with the right repo but a wrong name: `chat.rename` to `<repo> repo` (confirm with the user first if the chat has significant history under its old name).
3. Seed each newly created chat's memory minimally: create `~/side-chats/<new-chat-id>/MEMORY.md` with the repo name, path `~/workspace/repos/<repo>`, and a one-line purpose. Full learning happens in that chat's own onboarding — do not bulk-copy other chats' memories.

**Done when:** the final map shows exactly one `<repo> repo` chat per repo in the working set.

---

## Completion report

Reply with:

```
Chat recovery complete.

- Chats inventoried: [N] active, [M] archived
- Repo chats re-read: [K] (memory + recent history)
- Removed: [A] archived, [D] deleted (user-confirmed)
- Created: [C] new repo chats
- Blockers: [list, or "none"]

Final map:
| Repo | Chat name | chat_id |
|---|---|---|
| … | … | … |
```

Then stop. No next-step suggestions.

---

## Metadata

- slug: chat-histories-reread-remove-and-recreate
- status: active
- category: 29-chat-management
