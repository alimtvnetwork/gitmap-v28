# Muse execute in a step

Paste this after the task text. Onboarding stays in `01-prompts/27-muse-prompts/01-muse-master-prompt.md`. This file only starts the work.

**Version:** 1.0.0
**Updated:** 30-Oct-26

## Must follow

1. Pull first with `gitmap pull`. If that fails, use `git pull` and leave the tree clean.
2. Show the task list in chat as `Task-01`, `Task-02`, and so on. Each line ends with `Understood: [YES]`.
3. The next line after the list must be exactly one of these, plus a finish estimate:
   - `RUNNING. ETA: <duration>.`
   - `NOT RUNNING. Blocker: <one sentence>.`
4. `LISTING-WITHOUT-STARTING` is a failure. A list with no `RUNNING` line does not count as started.
5. While `RUNNING`, send a status line every five minutes: task id, what changed, what is next.
6. Search with `gitmap aum search`, `gitmap find`, and `gitmap cat`. Do not use `rg`, `grep`, or `Select-String`.
7. Write only relative paths (`02-spec/...`, `.ai-memory/...`, `01-prompts/...`).
8. Finish with `gitmap cpf "<module> - <summary>"`. The module is the subsystem. The word `commit` is not a module. No colon inside the message.

## Output shape

```text
Task-01: <name> — Understood: [YES]
Task-02: <name> — Understood: [YES]
RUNNING. ETA: <duration>.
```

If a step cannot start, stop after the list and write `NOT RUNNING. Blocker: <one sentence>.` Do not invent a running state.
