# GitMap Restriction Management

Version: 6.630.0
Updated: 2026-10-30

## Purpose

Named restrictions tell GitMap when commits are not allowed, which timezone applies, and which email owns the rule. The terminal prints the restricted status and the current time. A commit after the cutoff on a weekday is stamped before 08:00 when permit-stay is on. `--commit-only` keeps the restriction and does not push.

## Records

Each restriction has `Name`, `BlockedAfter`, `Email`, `Timezone`, `IsWeekendOpen`, and `IsPermitStay`. The file is `restrictions.json` under the user `.gitmap` directory.

## Commands

- `gitmap restrict status` prints name, time, timezone, blocked-after, email, weekend, permit-stay, and window.
- `gitmap restrict add` creates a record.
- `gitmap restrict edit` updates email, time, timezone, weekend, or permit-stay by name.

## Window

Weekdays before the cutoff are open. After the cutoff, the window is restricted unless `--auto-window` stamps 07:30 in the target timezone. `--weekend-open` leaves Saturday and Sunday unrestricted. `--weekend-closed` applies the cutoff on weekends too.

## Release

The minor version bump and the GitHub release both wait until the weekday window is open. A release during the restricted window is postponed. The command itself is in the tree at 6.630.0.
