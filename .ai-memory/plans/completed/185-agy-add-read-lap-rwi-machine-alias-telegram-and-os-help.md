# Plan 185: AGY Add/Add-Read, Running-Projects Tree (`rp prompts ls`), Last-Active-Projects (`lap`), Rerun-With-ID (`rwi`/`rwc`), Machine/Alias (`--ssh`), Two-Way Telegram Bot, Email Speed Settings, and Modernized OS/AGY Help

**Status:** Completed (`2026-09-28`)
**Canonical Spec:** [02-spec/21-app/90-agy-add-read-lap-rwi-machine-alias-telegram-and-os-help.md](../../../02-spec/21-app/90-agy-add-read-lap-rwi-machine-alias-telegram-and-os-help.md)

---

## 1. Consolidated Subtasks & Verified Outcomes

- [x] **Subtask 01 (`01-agy-add-read-and-24h-seq-cache-lap-rwi.md`)**:
  - Created `AgySequenceCache` table in `gitmap-running-prompts.db` ([agy_sequence_cache.go](../../../cli/store/agy_sequence_cache.go)) with a 24-hour TTL (`ExpiresAt = CreatedAt + 86400`), persisting deterministic project sequences (`1..M`) and conversation/prompt sequences (`P1..PK`).
  - Implemented `gitmap agy add` (`.` / `<path>`) and `gitmap agy add-read` (`.` / `<path>`, alias `ar`) in [agy_add_read.go](../../../cli/cmdagy/agy_add_read.go) to resolve `.` or `<path>`, register the project in Antigravity, and inject/run the Read Memory prompt (`01-prompts/read.md`).
  - Enhanced `gitmap agy rp ls` ([agy_running_projects.go](../../../cli/cmdagy/agy_running_projects.go)) and implemented `gitmap agy rp prompts ls` and `gitmap agy last-active-projects` (`lap`) `N [ls/help] [--limit/-l Y] [--offset/--skip Z] [--page/-p P] [--wordcount/--wc T] [--json] [--file/-f <filepath>]` ([agy_lap_and_rp_tree.go](../../../cli/cmdagy/agy_lap_and_rp_tree.go)) with default `N=24` hours, `Y=10` limit notice, `T=200` words tree view, and bracketed `[ProjectID | ConvID | SeqID]`.
  - Implemented `gitmap agy rerun-with-id` (`rwi`) and `gitmap agy rerun-with-convid` (`rwc` / `rwp`) ([agy_rerun_with_id.go](../../../cli/cmdagy/agy_rerun_with_id.go)) supporting short sequence IDs (`P1`, `1`) from the 24h SQLite cache, file paths, and `-p/--prompt` named templates.

- [x] **Subtask 02 (`02-machine-alias-ssh-telegram-email-and-os-agy-help.md`)**:
  - Implemented cross-platform `gitmap machine` and `gitmap alias` (`ls/change/set/revert/help [--ssh] [-y]`) at both root level and under `gitmap os` ([os_machine_alias.go](../../../cli/cmdos/os_machine_alias.go)), displaying IP, Alias (auto-defaulting to IP when unset), OS Hostname, and OS Platform across Windows, macOS, and Ubuntu/Linux.
  - Modernized `gitmap os help` ([os_help_modern.go](../../../cli/cmdos/os_help_modern.go)) and updated `gitmap agy help` ([agy_help.go](../../../cli/cmdagy/agy_help.go) & [agy.md](../../../cli/helptext/agy.md)).
  - Implemented two-way Telegram chatbot (`gitmap telegram` / `gitmap agy telegram`), Email speed setup (`gitmap email`), unified speed settings (`gitmap settings` / `gitmap agy settings`) ([agy_telegram_email_settings.go](../../../cli/cmdagy/agy_telegram_email_settings.go)), and canonical prompt [01-prompts/telegram-bot-setup.md](../../../01-prompts/telegram-bot-setup.md).
