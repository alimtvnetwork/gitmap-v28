# GitMap UI and Functionality Improvements — Component Specification

**Version:** 6.630.0
**Updated:** 2026-10-29

Binding component contract for the embedded settings UI in `cli/cmdui`. This file specifies behavior. It does not change Go, CSS, or JavaScript by itself.

Paths below are repository-relative. Wire JSON for these routes keeps the camelCase names already returned by `cli/cmdui` (`isSuccess`, `instanceId`). New booleans use an `is` prefix.

---

## 1. Current contracts this spec binds

| Surface | Current contract |
| --- | --- |
| Theme | `cli/cmdui/ui_assets_theme.go` |
| Shell header and fullscreen control | `cli/cmdui/ui_assets_markup_misc.go` (`uiAssetsMarkupShell`) |
| Installer heading | `cli/cmdui/ui_assets_markup_ops.go`, heading text `Installer & Workstation Profile Tree` |
| Installer list and dead handlers | `cli/cmdui/ui_assets_js_settings.go` |
| Installer HTTP | `GET /api/installers/profiles` in `cli/cmdui/ui_server_macro_installer.go` |
| Prompts HTTP | `GET /api/prompts/instances`, `handleAPIPromptsInstances` in `cli/cmdui/ui_server_instances.go` |
| Prompts rendering | `cli/cmdui/ui_assets_js_prompts_help.go` |
| Help HTTP | `GET /api/help/doc`, `handleAPIHelpDoc` in `cli/cmdui/ui_server_help.go` |
| Help catalog entry with no file | `account-switch` in `cli/helpdoc/catalog.go`; `cli/helpdoc` has no `account-switch.md` |
| Instances | `GET /api/instances`, cards built by `loadAgyInstances` in `cli/cmdui/ui_assets_js_fleet.go` |
| Error log | `store.LogInternalError` inserts into `InternalErrorLog` |

---

## 2. Hover, active, and input motion

### 2.1 Verified hover and active translation

In `cli/cmdui/ui_assets_theme.go`, pointer movement changes position only on these rules:

| Rule | Current declaration | Line |
| --- | --- | --- |
| `.nav-btn:hover` | `transform: translateX(2px)` plus background and color | 92 |
| `.nav-btn:active` | `transform: translateY(1px)` | 93 |
| `.card:hover` | `transform: translateY(-1px)` plus shadow and border | 110 |
| `.btn:hover` | `transform: translateY(-1px)` plus background, filter, and shadow | 122 |
| `.btn:active` | `transform: translateY(1px)` | 123 |
| `.btn-secondary:active` | `transform: translateY(1px)` | 126 |

`.nav-btn` transitions `all` for `0.2s` (line 91) and sets `will-change: transform`. `.card` transitions transform, shadow, and border for `0.25s` (line 109), sets `transform: translateZ(0)`, and sets `will-change: transform`. `.btn` transitions `all` for `0.2s` (line 121). `input, select, textarea` transition border and shadow for `0.2s` (line 116). Inputs and textareas have no opacity rule.

Other `transition` durations in the same file are `0.15s` (`body`, `tbody tr`, `.sub-tab-btn`, `.chip`, `.skip-link`), `0.2s` (off-canvas `#sidebar`), and `0.25s` (`.toast-out`).

Entrance keyframes (`appEntrance`, `toast-in`, `pulseGlow`), `.toast-out`, `.skip-link`, and the off-canvas sidebar `translateX(-105%)` are not hover or active rules. Those transforms stay. Their `transition` durations still move to `0.3s`.

### 2.2 Required motion

Delete `transform` from every `:hover` and `:active` rule in `cli/cmdui/ui_assets_theme.go`, including the six rules in the table above. Delete `will-change: transform` from `.nav-btn` and `.card`.

Hover and active may change color, border, background, box-shadow, and filter. They must not change position, translation, margin, or `top`.

The resting `.card` declaration `transform: translateZ(0)` stays. It is not a hover or active rule.

Every `transition` in `cli/cmdui/ui_assets_theme.go` uses `0.3s`. That covers `body`, `.nav-btn`, `.card`, `input` / `select` / `textarea`, `.btn`, `tbody tr`, `.sub-tab-btn`, `.chip`, `.toast-out`, `.skip-link`, and the off-canvas `#sidebar`. Keyframe `animation` durations stay as they are.

### 2.3 Required field opacity

Split opacity off the shared `input, select, textarea` rule so `select` does not gain it.

- `input` and `textarea` use `opacity: 0.9` at rest.
- `input:hover`, `input:focus`, `textarea:hover`, and `textarea:focus` use `opacity: 1`.
- On hover, `input` and `textarea` use a border lighter than `--border`.

Add token `--border-hover` in both `:root` and `:root[data-theme="light"]`. In each theme its lightness is higher than `--border`. Hover uses `border-color: var(--border-hover)`. Focus keeps `border-color: var(--primary)` and the existing focus shadow, including when the pointer is also hovering. Field transitions include `opacity` and `border-color` at `0.3s`.

---

## 3. Settings header mark and fullscreen control

### 3.1 Verified header

`uiAssetsMarkupShell` in `cli/cmdui/ui_assets_markup_misc.go` renders `<h2 id="header-title">Settings</h2>` with no GitMap mark. The sidebar logo in the same constant is an inline lightning SVG. It has no `src` and no file path.

`#btn-fullscreen` uses `class="btn"` and the label `⛶ Fullscreen`. `.btn:hover` and `.btn:active` translate, so this control moves on pointer press.

### 3.2 Required header mark

Place that same inline SVG in the header, in the left cluster, before `#header-title`. The mark is inline markup. It has no `src`, no `href`, and no filesystem path. Give the cluster an accessible name `GitMap settings`. The SVG stays `aria-hidden="true"`.

### 3.3 Required fullscreen control

`#btn-fullscreen` must not use `btn`, `btn-secondary`, or any other class whose hover or active rule sets `transform`, `margin`, or `top`.

Give it a dedicated class `btn-fullscreen`. That class sets padding, font size, background, border, and color only. Its `:hover` and `:active` rules may change color, border, and background. They set no `transform`. Keep the fullscreen glyph and the visible word `Fullscreen`. Keep `id="btn-fullscreen"` so the existing fullscreen listener still finds it.

---

## 4. Installer profile tree and create-profile

### 4.1 Verified installer behavior

The installer tab heading says tree: `Installer & Workstation Profile Tree`. The list container is `#profiles-tree-list`.

`loadInstallerProfiles` in `cli/cmdui/ui_assets_js_settings.go` fetches `GET /api/installers/profiles` and paints one block per profile. Tools render as a flat wrapping row of chips inside that block. The layout is a flat list with a chip row, not a parent profile with child tool nodes.

`handleAPIInstallerProfiles` returns `success`, `profiles`, and each profile's `name`, `title`, `description`, `tools`, and `isInstalled`. The tools come from `cmdinstall.AllInstallProfiles()`.

`saveInstaller` and `testInstaller` are assigned in the dead-handler list in `cli/cmdui/ui_assets_js_settings.go`. Each shows the toast `Not implemented yet`.

`cmdinstall.InstallProfile` has a reader, `AllInstallProfiles`, and no store writer. `store.(*DB).CreateInstaller` inserts a `model.InstallerScript` into `installer_scripts`. That record has name, slug, target OS, and instructions. It has no child tool list. It is not a writer for `InstallProfile`.

### 4.2 Required tree

`loadInstallerProfiles` renders `#profiles-tree-list` as a tree:

- One parent node per profile, labeled with `title` and `name`.
- Description stays on the parent.
- Each string in `tools` is a child node of that parent.
- Children sit in a nested list under the parent (`ul` inside the parent `li`, with `role="tree"` on the outer list and `role="treeitem"` on each node).
- Existing Install Profile, Uninstall, and per-tool install actions stay on the matching node.

A chip row that is not nested under the profile fails this contract. The empty and error copy already in `loadInstallerProfiles` stays.

The GET response shape stays (`success`, `profiles`, `tools`) so the tree reads the handler that exists today.

### 4.3 Required create-profile handler

Remove `saveInstaller` and `testInstaller` from the dead-handler list. The Save Installer button calls create-profile. Create-profile sends `POST /api/installers/profiles` with the form fields `inst-name`, `inst-os`, and `inst-cmd`.

`handleAPIInstallerProfiles` must branch on method. GET keeps the current list. POST is create-profile. A POST must not fall through into the GET list.

No `InstallProfile` store writer exists. Create-profile therefore does all of the following:

1. Build a typed `*apperror.AppError` whose code is `constants.ErrCodeInstallerProfileWriterMissing` (`E9007`).
2. Call `store.LogInternalError` with that code, the operation name `createInstallerProfile`, the error text, the submitted name and target OS, and source file `cli/cmdui/ui_server_macro_installer.go`.
3. Respond with HTTP 409 and JSON `isSuccess: false`, `isFail: true`, `errorCode: "E9007"`, and `error` set to the typed error text.

The handler must not call `CreateInstaller`. That would store an installer script and report a profile write that did not happen.

The page reads `error` and shows that text. The toast text `Not implemented yet` is absent for Save Installer, create-profile, and Test on Selected Node.

### 4.4 Required test action

Test on Selected Node calls `POST /api/installers/test`. No node-test runner is wired in `cli/cmdui`. The handler returns HTTP 409 with `isSuccess: false`, `isFail: true`, `errorCode` `constants.ErrCodeInstallerTestRunnerMissing` (`E9011`), and `error` set to the typed error text. It logs that error through `store.LogInternalError` before writing the response. The page shows `error`.

Add both codes to `cli/constants/constants_error_codes.go`.

---

## 5. Prompts grouped by instance

### 5.1 Verified prompt flow

`handleAPIPromptsInstances` returns `AgyMultiInstancePromptResponse`: `instances[]` of `AgyInstancePromptPayload`, each with `instanceId`, `instanceName`, `running`, `queued`, and `recentConvs`.

`loadPromptsList` in `cli/cmdui/ui_assets_js_prompts_help.go` walks those arrays and pushes every row into one `allPrompts` array. `renderPromptsList` then sorts that flat array and prints `Instance:` as a line on each card. A row click copies text into the scratchpad. Rows have no Copy, Rerun, or Enqueue controls. No cmdui route reruns or enqueues a prompt.

### 5.2 Required grouping

Keep the GET handler per instance. Stop flattening instances into one list before paint.

`renderPromptsList` paints one group per instance that still has rows after the search and category filter. The group heading shows `instanceName` and `instanceId`. Rows from `running`, `queued`, and `recentConvs` render only inside that group. Status rank may still order rows inside the group. A group with no visible rows is omitted. The existing empty-filter sentence stays when every group is empty.

### 5.3 Required row actions

Each row has three controls: Copy, Rerun, and Enqueue.

Copy writes that row's text with `navigator.clipboard.writeText`. Success uses the existing copied toast. A clipboard rejection shows the rejection text.

Rerun sends `POST /api/prompts/rerun`. Enqueue sends `POST /api/prompts/enqueue`. The JSON body carries `instanceId` and the row `text`. Mount both routes beside `/api/prompts/instances`.

No rerun sender and no enqueue sender exist in `cli/cmdui`. Each handler therefore:

1. Returns a typed `*apperror.AppError`.
2. Rerun uses `constants.ErrCodePromptRerunSenderMissing` (`E9008`). Enqueue uses `constants.ErrCodePromptEnqueueSenderMissing` (`E9009`).
3. Calls `store.LogInternalError` with the code, the instance id, and the row text.
4. Responds with HTTP 409, `isSuccess: false`, `isFail: true`, `errorCode`, and `error`.

The click shows `error` from that JSON. It does not show `Not implemented yet`. The handlers must not open a socket, call an editor command, or invent an IDE injection protocol. Add `E9008` and `E9009` to `cli/constants/constants_error_codes.go`.

---

## 6. Help document contrast and missing topics

### 6.1 Verified help behavior

`GET /api/help/doc` is `handleAPIHelpDoc`. A missing embedded file takes the `helpdoc.ReadRaw` error branch and calls `http.Error` with status 404 and the body `topic documentation not found: ` plus the topic name.

`account-switch` is a catalog topic. `resolveHelpAlias` does not map it to another file. `ReadRaw` therefore reads `account-switch.md`, then the lowercased and dashed names. `cli/helpdoc` has no `account-switch.md`, so the branch returns 404. The same miss applies to a topic the catalog lists with `hasMarkdown: false`.

`loadHelpDoc` treats any non-OK response as `throw new Error('HTTP ' + res.status)`. The viewer then shows `Failed to load manual for account-switch: HTTP 404`. The response sentence is discarded.

`#help-doc-viewer` in `cli/cmdui/ui_assets_markup_misc.go` sets `color: var(--text)` and `background: var(--code-bg)`. Light theme sets `--text` to `#0f172a` and `--code-bg` to `hsl(230, 25%, 8%)`. Body text and the error string are both dark on that dark code background. `--code-text` is `#e2e8f0` on `:root` and is the light foreground already used by `#editor-area`. Light theme does not override `--code-text`.

### 6.2 Required contrast

`#help-doc-viewer` uses `color: var(--code-text)` on `background: var(--code-bg)` for the loaded manual and for the error string. Do not paint either string with `var(--text)` or with the light-theme `--error` value while the background is `--code-bg`.

Set `--code-text: #e2e8f0` on `:root[data-theme="light"]` as well, so the light theme owns the token the viewer uses.

### 6.3 Required missing-topic body

When `ReadRaw` fails, `handleAPIHelpDoc` still uses status 404 and writes a text body the viewer can show unchanged:

`No manual is embedded for '{topic}'. The help catalog lists this topic, and cli/helpdoc has no matching markdown file.`

`loadHelpDoc` reads `response.text()` for a non-OK response and sets the viewer to that body. The viewer must not contain the string `HTTP 404`. JSON `format=json` uses the same sentence in `error` with `isSuccess: false` and `isFail: true`.

---

## 7. Per-instance message composer

### 7.1 Verified gap

`loadAgyInstances` paints one card per `GET /api/instances` record inside `#instances-list-container`. The card shows name, id, running badge, process id, language server, and config path. It has no text field and no send control. `cli/cmdui/ui_server.go` does not mount a message route. No chat-to-IDE UI exists.

### 7.2 Required composer

Each detected-instance card includes one composer: a label `Message`, a `textarea`, and a send button. The composer sits on `--card` or `--raised` with foreground `--text`. It does not use `--text` on `--code-bg`.

Send posts `POST /api/instances/message` with JSON `instanceId` and `text`. Mount the route next to `GET /api/instances`. Reject an empty `text` or empty `instanceId` with HTTP 400, `isFail: true`, and a sentence that names the missing field.

### 7.3 Required accept result

No sender is wired from this route to an IDE. The handler stores the submitted text by calling `store.LogInternalError` with `constants.ErrCodeInstanceMessageSenderMissing` (`E9010`), operation `postInstanceMessage`, source file `cli/cmdui/ui_server_instances.go`, and details that include `instanceId` and the full `text`. That insert is the store.

The response is HTTP 409:

- `isSuccess`: `false`
- `isFail`: `true`
- `isAccepted`: `false`
- `errorCode`: `"E9010"`
- `error`: `The message was stored and was not delivered to the IDE. No sender is wired.`
- `text`: the stored text

The card shows `error`. It must not show a success toast, and it must not say the message reached Cursor or any other IDE. The handler must not call an editor, a language server, or an injection protocol. Add `E9010` to `cli/constants/constants_error_codes.go`.

---

## 8. Text to add to `.ai-memory/strictly-avoid.md`

Add the following section to `.ai-memory/strictly-avoid.md`. This spec only records the text. It does not edit that file.

```markdown
## Hover and Active Position Shift — TOTAL BAN

🔴 **Hover and active states must not change position, translation, margin, or top. Color change is allowed.**

Forbidden:
- ❌ `transform: translate`, `translateX`, or `translateY` on `:hover` or `:active`.
- ❌ Changing `margin` or `top` on `:hover` or `:active`.
- ❌ Giving a control a class whose hover or active rule shifts position. `#btn-fullscreen` must not use `btn` while `.btn:hover` or `.btn:active` translates.

Allowed work:
- ✅ Color, border, and background changes on hover and active.
- ✅ Entrance keyframes, toast motion, the skip-link reveal, and the off-canvas sidebar transform. Those are not hover or active states.

**Why:** The embedded settings UI in `cli/cmdui/ui_assets_theme.go` shifts `.nav-btn`, `.card`, and `.btn` under the pointer. The hit target must stay still.
```

---

## Acceptance Criteria

1. **Hover and active translation removed.** Pass when `.nav-btn`, `.card`, `.btn`, and `.btn-secondary` hover and active rules in `cli/cmdui/ui_assets_theme.go` contain no `transform`, and `will-change: transform` is gone from `.nav-btn` and `.card`. Fail when any of those rules still contains `translateX` or `translateY`.

2. **Hover and active position properties fixed.** Pass when hover and active rules change only color, border, background, box-shadow, or filter. Fail when a hover or active rule sets `transform`, `margin`, or `top`.

3. **Transition duration.** Pass when every `transition` in `cli/cmdui/ui_assets_theme.go` is `0.3s`. Fail when any `transition` in that file remains `0.15s`, `0.2s`, or `0.25s`.

4. **Field opacity and hover border.** Pass when `input` and `textarea` rest at opacity `0.9`, move to opacity `1` on hover and focus, and use `--border-hover` on hover while focus keeps `--primary`. Fail when inputs have no opacity rule, when hover opacity stays `0.9`, or when the hover border is the same token as the rest border.

5. **Fullscreen class.** Pass when `#btn-fullscreen` omits `btn` and uses `btn-fullscreen`, and `btn-fullscreen` hover and active rules have no `transform`. Fail when the control still has `class="btn"` or when its class translates on hover or active.

6. **Header mark.** Pass when the settings header shows the sidebar lightning SVG as inline markup before `#header-title`, with no `src` and no filesystem path. Fail when the header has no mark, or when the mark uses a path.

7. **Profile tree.** Pass when `#profiles-tree-list` renders each profile as a parent `treeitem` and each `tools` entry as a nested child `treeitem`. Fail when tools render only as a flat chip row.

8. **Create-profile handler.** Pass when Save Installer POSTs to `/api/installers/profiles`, the handler returns HTTP 409 with `errorCode` `E9007` and `isFail: true`, and `store.LogInternalError` records that code. Fail when the click still toasts `Not implemented yet`, when POST returns the profile list, or when the handler calls `CreateInstaller` and reports success.

9. **Installer test handler.** Pass when Test on Selected Node POSTs to `/api/installers/test` and the response shows the logged `E9011` error text. Fail when that button still toasts `Not implemented yet`.

10. **Prompts grouped by instance.** Pass when each visible instance is its own group and rows from that instance render only inside it. Fail when `renderPromptsList` paints one flat list sorted across instances.

11. **Prompt row actions.** Pass when every prompt row has Copy, Rerun, and Enqueue. Copy uses the clipboard. Rerun and Enqueue POST to `/api/prompts/rerun` and `/api/prompts/enqueue`, the handlers log `E9008` or `E9009` through `store.LogInternalError`, and the click shows the returned `error` text. Fail when either action toasts `Not implemented yet`, or when either handler calls an IDE injection protocol.

12. **Help foreground on the code background.** Pass when `#help-doc-viewer` uses `var(--code-text)` for the manual body and for the error string, and light theme sets `--code-text` to `#e2e8f0`. Fail when the viewer uses `var(--text)` on `var(--code-bg)`.

13. **Missing help topic copy.** Pass when `GET /api/help/doc?topic=account-switch` returns status 404 and the body `No manual is embedded for 'account-switch'. The help catalog lists this topic, and cli/helpdoc has no matching markdown file.`, and the viewer shows that sentence. Fail when the viewer shows `HTTP 404`.

14. **Instance composer and stored undelivered message.** Pass when each detected-instance card has a composer that POSTs to `/api/instances/message`, the handler writes the text through `store.LogInternalError` with `E9010`, and the JSON has `isAccepted: false` plus the error sentence `The message was stored and was not delivered to the IDE. No sender is wired.` Fail when a card has no composer, when `isAccepted` is `true`, or when the UI says the message reached Cursor.

15. **Strictly-avoid section recorded.** Pass when the section in §8 is the text intended for `.ai-memory/strictly-avoid.md`, including the sentence `Hover and active states must not change position, translation, margin, or top. Color change is allowed.` Fail when that sentence is missing from §8. This criterion does not require that file to have been edited.
