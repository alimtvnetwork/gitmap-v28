# Plan: 02-github-desktop-sync-fix

## User Request (Verbatim)
> In Ubuntu, if we do run the GitKraken, it does not add all the projects to GitHub Desktop. It adds some, but then again crashes or I don't know what happens. You can try to investigate the GitKraken errors to see if there is any error logs and make sure you do not swallow any error or just skip any error. Remember that, try to fix it.

## Specifications
- Architecture Spec: `02-spec/21-app/02-github-desktop-sync-fix/01-architecture-spec.md`
- Component Spec: `02-spec/21-app/02-github-desktop-sync-fix/02-component-spec.md`

## Subtasks
1. `01-investigate-and-resolve-cli.md`: Fix CLI candidate resolution and command construction (`cli/desktop/resolve.go`, `cli/desktop/desktop.go`).
2. `02-pacing-and-error-surfacing.md`: Implement pacing/throttling and surface non-swallowed errors across IDE sync and Desktop sync (`cli/cmdide/ide_sync.go`, `cli/cmdide/ide_add.go`, `cli/cmddesktopsync/desktopsync_ops.go`).

## Wave Execution Strategy
- Wave 1: Subagent 1 authors Architecture Spec and Subtask 01; Subagent 2 authors Component Spec and Subtask 02.
- Wave 2: Worker 01 implements Task-01; Worker 02 implements Task-02.
- Wave 3: Targeted verification and atomic commit via GitMap.
