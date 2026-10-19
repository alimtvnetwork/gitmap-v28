# 29 — Chat Management

Side-chat lifecycle prompts: inventory, re-read, prune, and recreate the user's per-repository chats. Built for VM-replacement recovery — run the prompt on any fresh machine and the chat workspace rebuilds itself.

## Prompts

| # | File | Purpose |
|---|---|---|
| 01 | `01-chat-histories-reread-remove-and-recreate.md` | Re-read every repo side chat's history and memory, archive/delete stale or duplicate chats (with per-chat confirmation), and create exactly one `<repo> repo` chat per repository in the working set. |

## Conventions

- Repo chats are named exactly `<repo> repo` — no prefixes, suffixes, or variations.
- Archiving is the default for stale chats (reversible); deletion requires the user's explicit word per chat.
- `chat.create` forking is main-chat-only; side chats always create with `context_mode: "fresh"`.
