# Time-Restricted Execution Configuration

Version: 6.630.0
Updated: 2026-10-30

## Purpose

One machine setting decides whether the weekday cutoff applies to every repository. Auto is the default. Flag mode waits for `--schedule`. Settings and `gitmap restrict` read and write the same value, and `gitmap restrict status` prints it.

## Mode

`EnforcementMode` lives on the restriction file next to the named window.

- `auto`: every GitMap commit on this machine sets time enforcement, including commits that omit `--schedule`.
- `flag`: enforcement runs only when `--schedule` or `-sc` is passed.

An empty or missing mode is auto.

## Commands

- `gitmap restrict status` prints the clock, the window, and the enforcement line.
- `gitmap restrict mode` prints the active mode.
- `gitmap restrict mode auto` and `gitmap restrict mode flag` switch it.

## Settings

The settings select `restriction.enforcement` loads the machine mode and saves it back to the same file. Saving settings does not create a second store.

## Release

No release in this change. The weekday window is closed after 08:00.
