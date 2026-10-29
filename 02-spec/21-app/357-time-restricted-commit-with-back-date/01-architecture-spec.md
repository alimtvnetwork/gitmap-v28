# Time-Restricted Commit with Back Date

Version: 6.630.0
Updated: 2026-10-30

## Purpose

When the weekday window is closed, GitMap stamps the commit at 07:30 in the restriction timezone and still pushes. The author date and the committer date both carry that timezone offset, so the back date is what Git shows.

## Timezone

The stamp uses the timezone already resolved for the commit: `--tz` when present, otherwise the machine restriction timezone. `Asia/Singapore` at 10:31 becomes `07:30:00+08:00` on the same morning. A missing timezone uses the host local zone.

## Push

The back date is applied before `git commit`. The push runs after that commit unless `--commit-only` is set. A push does not rewrite the stored author date.

## Icons

Every context-menu entry uses the GitMap executable icon (`{exe},0`) instead of `shell32.dll` placeholders. Explorer shows the new icons after the context menu is installed again.
