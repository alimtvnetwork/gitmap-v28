# Repository Management and Testing Instructions

> **Directory:** `02-spec/21-app/316-repository-management-and-testing-instructions/`  
> **Status:** In Progress  
> **Spec Index:** `SPEC-316`

---

## Documents

| Document | Description |
|---|---|
| [01-architecture-spec.md](01-architecture-spec.md) | Scope, workstreams, verification gates, and release procedure |

## Plan

| Artifact | Path |
|---|---|
| Master plan | `.ai-memory/plans/316-repository-management-and-testing-instructions.md` |
| Subtasks | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/` |

## Workstreams

| Track | Repository | Goal |
|---|---|---|
| A | Antigravity-Manager | Verify auto-switch, fast-forward, and single-prompt rerun; CI green; minor-bump release with tags |
| B | gitmap-v28 | Review recent commands, verify fast-forward logic against spec, resolve CI failures; minor-bump release |
| C | img-pdf-v2 (local folder `img-pdf`) | Root-cause the right-click failure; explain the missing YouTube thumbnail download |
| D | yt-manage | Functional verification of the yt-manage CLI and any gitmap-side YT integration |
