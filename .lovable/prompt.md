# Lovable Prompts Reference

## Available Prompts

| Keyword | File | Purpose |
|---------|------|---------|
| `read memory` | [01-read-prompt.md](prompts/01-read-prompt.md) | Full AI onboarding protocol — reads all memory, guidelines, and specs |
| `write memory` / `end memory` | [02-write-prompt.md](prompts/02-write-prompt.md) | End-of-session persistence protocol — writes everything learned, done, and pending |
| `consolidate-plans` / `cleanup-plans` | [29-lovable-memory-consolidation-and-cleanup.md](prompts/29-lovable-memory-consolidation-and-cleanup.md) | Plan consolidation, audit folder removal & milestone compaction |


## Usage

- When the user says **"read memory"**, follow the onboarding protocol in `prompts/01-read-prompt.md`.
- When the user says **"write memory"** or **"end memory"**, follow the persistence protocol in `prompts/02-write-prompt.md`.
