# GitMap Settings UI and Functionality Improvements Architecture Specification

**Version:** 6.630.0
**Updated:** 2026-10-29
**Status:** Binding

## Purpose

This specification binds the settings UI fixes for GitMap. The settings UI is embedded Go strings in `cli/cmdui`. `mountSPARoutes` in `cli/cmdui/ui_server.go` serves that UI on the UI `ServeMux`.

The change makes SSH online, ping, and keepalive truthful, gives the page a working command section and two working terminals, and persists settings failures in `gitmap-errors.db`. Debug and verbose become settings toggles. Layout, installer profiles, prompts, help coloring, fullscreen, the GitMap icon, text-box opacity, and chat to a Cursor or Antigravity instance are part of the same settings UI behavior.

Implementation uses relative git paths, positive `is`/`has` booleans, and the existing writers and routes named below.

## Current Failure

1. SSH Online is wrong. `buildNodeSummaries` in `cli/cmdui/ui_server_ssh.go` sets `IsOnline` true for every node.
2. The Ping button calls `testNode`. `cli/cmdui/ui_assets_js_settings.go` replaces `window.testNode` with a toast `Not implemented yet`.
3. `GET /api/ssh/test` is mounted on the default mux. It is not mounted on the UI `ServeMux`, so the settings page cannot call it.
4. `executeNodePingCheck` is a 2 second TCP dial to port 22. It is not ICMP. The settings UI still presents that result as a general online flag that is hardcoded true before the dial matters.
5. Keepalive is not shown. `clientaliveinterval` is stripped in `cli/cmdssh/sshconfig.go`.
6. Two terminals share the DOM ids `term-output` and `term-input`, so the first terminal receives input and output and the second terminal does not.
7. `POST /api/terminal/exec` in `cli/cmdui/ui_server_terminal.go` is a real handler. `POST /api/command/exec` is the same handler. Settings JS does not call `POST /api/command/exec`.
8. UI errors go to `console.error` and toasts. They are not written to `gitmap-errors.db`. The writer is `LogInternalErrorRecord` in `cli/store/errors_split_ops.go`. The database tables are `InternalErrorLog` and `FailedCommand`.
9. Debug and verbose flags exist on the CLI (`cli/verbose`, `FlagVerbose`). They are absent from the settings UI.

## Target Behavior

### SSH online, ping, and keepalive

`buildNodeSummaries` sets `IsOnline` from `executeNodePingCheck`. A successful 2 second TCP dial to port 22 sets `IsOnline` true. A failed dial sets `IsOnline` false. The check stays a TCP dial to port 22.

The Ping button calls `GET /api/ssh/test` on the UI origin. Settings JS does not replace `window.testNode` with the toast `Not implemented yet`.

`cli/cmdssh/sshconfig.go` keeps `clientaliveinterval`. Each node summary shows that keepalive value. `HasClientAliveInterval` is true when the interval is present.

```json
{
  "IsOnline": false,
  "HasClientAliveInterval": true,
  "ClientAliveInterval": "60"
}
```

### Command section and terminals

The settings page has a command section. That section calls `POST /api/command/exec`.

Each terminal pane has its own output element id and its own input element id. The ids `term-output` and `term-input` are not reused across panes. Both panes call `POST /api/terminal/exec` and both accept typing and show output.

A failed command execution is stored in `FailedCommand` in `gitmap-errors.db`.

### Hover, opacity, fullscreen, and icon

Hover styles do not change the layout box of a control. Margin, padding, border width, and position stay constant while the pointer is over the control.

Text boxes use opacity `0.9` at rest and opacity `1` when engaged. The opacity transition duration is `0.3s`.

The settings UI exposes fullscreen and shows the GitMap icon.

### Installer profiles, prompts, help, and chat

The settings UI shows the installer profile tree and can create a profile.

Prompts are grouped by instance. Each instance group provides copy, rerun, and enqueue. The instance list comes from `GET /api/prompts/instances`.

Help manual content from `GET /api/help/doc` uses distinct colors for headings, code, and body text.

The user can chat to a Cursor instance or an Antigravity instance. The instance list comes from `GET /api/instances`.

### Debug, verbose, and saved errors

The settings UI has debug and verbose toggles. They are stored in the UI settings JSON read and written by `/api/settings`.

```json
{
  "IsDebugEnabled": false,
  "IsVerboseEnabled": false
}
```

The CLI flags `cli/verbose` and `FlagVerbose` remain the CLI switches. The settings page reads the JSON fields above.

Every user-visible settings failure still may show a toast, and it is also passed to `LogInternalErrorRecord`. That call writes `InternalErrorLog` in `gitmap-errors.db`. `IsResolved` stays false until a later resolve path marks it.

```json
{
  "ErrorCode": "ui.settings",
  "ErrorType": "ui",
  "Command": "settings",
  "Message": "ssh test failed",
  "IsResolved": false
}
```

### Spec file removal

Remove the generic next-100 UI specs and today's time-restriction specs by deleting those files in a normal new commit. Git history is not rewritten in this change.

## API Contracts

All routes below are registered on the UI `ServeMux` created in `cli/cmdui/ui_server.go`. `mountAPIRoutes` already mounts the SSH, terminal, settings, prompt, help, and instance routes. `mountSPARoutes` continues to serve the embedded settings document.

### Mount `GET /api/ssh/test` on the UI mux

| Item | Contract |
| --- | --- |
| Method | `GET` |
| Path | `/api/ssh/test` |
| Mount | UI `ServeMux`, in `mountSSHAndEditorRoutes`, beside `/api/ssh/nodes` |
| Behavior | Runs `executeNodePingCheck`: a 2 second TCP dial to port 22 |
| Success body | `IsOnline` is true |
| Failure body | `IsOnline` is false, and the failure is passed to `LogInternalErrorRecord` |

The default mux registration may remain. The settings page calls the UI mux route.

### Persist UI failures via `LogInternalErrorRecord`

Settings JS reports each user-visible failure to the UI server. The handler builds an `InternalErrorRecord` and calls `LogInternalErrorRecord` in `cli/store/errors_split_ops.go`.

| Item | Contract |
| --- | --- |
| Writer | `LogInternalErrorRecord` |
| Database | `gitmap-errors.db` |
| Internal UI row | `InternalErrorLog` |
| Failed command row | `FailedCommand` for a failed `POST /api/terminal/exec` or `POST /api/command/exec` |
| Empty message | `LogInternalErrorRecord` returns without inserting when `Message` is empty |
| Boolean | `IsResolved` is false on insert |

`console.error` is not the persistence path.

### Settings debug and verbose toggles in the UI settings JSON

`/api/settings` reads and writes the UI settings JSON. The document includes:

| Field | Meaning |
| --- | --- |
| `IsDebugEnabled` | Debug toggle from the settings UI |
| `IsVerboseEnabled` | Verbose toggle from the settings UI |

Both fields are booleans. A missing field is stored as false on the next save. The settings page renders the toggles from these fields.

### Routes the settings page already has

These routes stay on the UI mux and gain the behavior in Target Behavior. This change does not replace their handlers with stubs.

| Route | Use |
| --- | --- |
| `POST /api/terminal/exec` | Both terminal panes. Handler: `handleAPITerminalExec` in `cli/cmdui/ui_server_terminal.go` |
| `POST /api/command/exec` | Command section. Same handler as `POST /api/terminal/exec` |
| `GET /api/prompts/instances` | Prompts grouped by instance, with copy, rerun, and enqueue |
| `GET /api/instances` | Chat target list for a Cursor instance or an Antigravity instance |
| `GET /api/help/doc` | Help manual coloring |
| `GET /api/settings` | UI settings JSON, including `IsDebugEnabled` and `IsVerboseEnabled` |

## Out of Scope

This change does not edit context-menu or diff files. Those paths are owned by another agent:

- `cli/cmdinstall/installctx.go`
- `cli/cmdinstall/installctxentries.go`
- `cli/cmddiff/single_diff.go`
- `cli/constants/constants_installctx.go`

The Tortoise-style context menu gap, and the rule that a diff with no changes must say so and must not `Set-Location` a file path, stay recorded in User Request. This change does not implement them by editing the paths above.

This change does not rewrite git history. Spec removal is a forward delete.

This change does not replace the port-22 TCP dial with ICMP.

This change does not move the settings UI out of the embedded Go strings in `cli/cmdui`.

## User Request (Verbatim)

settings SSH keepalive/ping/online are wrong; a command section and a working terminal; hover must not move controls; installer profile tree and create profile; prompts grouped by instance with copy, rerun, enqueue; debug and verbose; errors saved to the GitMap errors database; text boxes opacity 0.9 to 1 with 0.3s transition; fullscreen and GitMap icon; help manual coloring; chat to Cursor or Antigravity instance; Tortoise-style context menu gap; diff with no changes must say so and must not Set-Location a file path; remove generic next-100 UI specs and today's time-restriction specs without rewriting git history in this change.

## Acceptance Criteria

1. PASS: `buildNodeSummaries` sets `IsOnline` from the port-22 TCP dial, and a failed dial yields `IsOnline` false. FAIL: every node is returned with `IsOnline` true.
2. PASS: Ping on the settings page calls `GET /api/ssh/test` on the UI `ServeMux` and shows the dial result. FAIL: `window.testNode` still toasts `Not implemented yet`, or `/api/ssh/test` remains unmounted on the UI mux.
3. PASS: `executeNodePingCheck` remains a 2 second TCP dial to port 22. FAIL: the settings ping path is ICMP, or the dial timeout is changed.
4. PASS: `clientaliveinterval` is kept in `cli/cmdssh/sshconfig.go` and the node summary shows it with `HasClientAliveInterval`. FAIL: the key is stripped, or keepalive is hidden.
5. PASS: the two terminal panes use distinct DOM ids and both run through `POST /api/terminal/exec`. FAIL: both panes still use `term-output` and `term-input`, or the second pane does not accept input.
6. PASS: the command section calls `POST /api/command/exec`, and that route remains the `handleAPITerminalExec` handler. FAIL: settings JS leaves `POST /api/command/exec` unused, or the route is a stub.
7. PASS: a user-visible settings failure is written by `LogInternalErrorRecord` into `InternalErrorLog` in `gitmap-errors.db`, and a failed terminal or command execution is stored in `FailedCommand`. FAIL: the failure exists only as `console.error` or a toast.
8. PASS: `/api/settings` persists `IsDebugEnabled` and `IsVerboseEnabled` in the UI settings JSON, and the settings page shows both toggles. FAIL: debug and verbose exist only as `cli/verbose` and `FlagVerbose`.
9. PASS: hovering a control leaves its margin, padding, border width, and position unchanged. FAIL: hover moves a neighboring control.
10. PASS: the settings UI shows the installer profile tree and can create a profile. FAIL: either action is missing.
11. PASS: prompts from `GET /api/prompts/instances` are grouped by instance, and each group offers copy, rerun, and enqueue. FAIL: prompts are ungrouped, or any of the three actions is missing.
12. PASS: text boxes transition opacity from `0.9` to `1` over `0.3s`. FAIL: opacity or duration differs.
13. PASS: the settings UI can enter fullscreen and shows the GitMap icon. FAIL: either is missing.
14. PASS: help manual content from `GET /api/help/doc` colors headings, code, and body text differently. FAIL: the manual is a single uncolored text run.
15. PASS: the user can chat to a Cursor instance or an Antigravity instance chosen from `GET /api/instances`. FAIL: chat has no instance target.
16. PASS: `cli/cmdinstall/installctx.go`, `cli/cmdinstall/installctxentries.go`, `cli/cmddiff/single_diff.go`, and `cli/constants/constants_installctx.go` are unchanged by this change. FAIL: this change edits any of those files.
17. PASS: the generic next-100 UI specs and today's time-restriction specs are deleted in a new commit, and git history is not rewritten. FAIL: those specs remain, or history is rewritten to remove them.
18. PASS: the settings document is still the embedded `cli/cmdui` strings served by `mountSPARoutes`. FAIL: the settings UI is served from a path outside `cli/cmdui`.
