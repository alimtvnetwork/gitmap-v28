# Migration & Verification Spec: Flexible Path Parsing and Asset Image Reading

- **Task Sequence:** 328
- **Task Slug:** `flexible-path-parsing-and-asset-image-reading`
- **Spec Path:** `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/03-migration-and-verification-spec.md`
- **Status:** APPROVED
- **Created At:** 2026-10-25

---

## 1. Executive Summary & Verification Scope

This specification defines the formal verification protocol, edge-case coverage matrix, backward compatibility guarantees, and documentation synchronization workflows for Task 328.

### 1.1 Strict No-Build, No-Test Protocol
In accordance with the task instructions:
- **Zero Build Execution:** Subagents must not run `go build`, `cargo build`, or binary compilers during specification and implementation phases.
- **Zero Test Execution:** Subagents must not run `go test`, `pytest`, or automated test suites during the execution cycle.
- **Static & Cognitive Verification:** All verification is conducted through static AST inspections, schema audits, boundary checks, and formal equivalence proofs.

---

## 2. Path Parsing Verification Matrix

Each of the four primary formats, plus critical edge cases, is evaluated against standard parsing criteria:

| Case ID | Input String Sample | Target Category | Expected Clean Path | Expected Flags |
| :--- | :--- | :--- | :--- | :--- |
| **TC-F1-01** | `@[C:\Users\Admin\img.png]` | Format 1 (Windows Absolute) | `C:/Users/Admin/img.png` | `hasBracketWrapper: true`, `isAbsolute: true` |
| **TC-F1-02** | `@[D:/work/repo/assets/logo.png]` | Format 1 (Forward Slash Absolute) | `D:/work/repo/assets/logo.png` | `hasBracketWrapper: true`, `isAbsolute: true` |
| **TC-F2-01** | `![Alt](file:///C:/Users/Admin/img.png)` | Format 2 (Markdown File URI) | `C:/Users/Admin/img.png` | `hasMarkdownWrapper: true`, `isAbsolute: true` |
| **TC-F2-02** | `![Screenshot](file:///D:/work/my%20assets/test.jpg)` | Format 2 (URI with `%20` space) | `D:/work/my assets/test.jpg` | `hasMarkdownWrapper: true`, `isAbsolute: true` |
| **TC-F3-01** | `1. @[assets/media_01.png] # relative from repo` | Format 3 (Numbered List + Repo Rel) | `assets/media_01.png` | `hasListPrefix: true`, `hasTrailingComment: true`, `isInsideRepo: true` |
| **TC-F3-02** | `04) @[src/components/button.tsx] // react component` | Format 3 (Paren list + `//` comment) | `src/components/button.tsx` | `hasListPrefix: true`, `hasTrailingComment: true`, `isInsideRepo: true` |
| **TC-F4-01** | `1. @[../repo-cache/assets/media_01.png] # outside repo` | Format 4 (Numbered List + Parent Rel) | `../repo-cache/assets/media_01.png` | `hasListPrefix: true`, `isParentRelative: true`, `isInsideRepo: false` |
| **TC-F4-02** | `- @[../../shared/vault/key.pem]` | Format 4 (Bullet List + Multi-Parent) | `../../shared/vault/key.pem` | `hasListPrefix: true`, `isParentRelative: true`, `isInsideRepo: false` |
| **TC-EDGE-01**| `   "@[assets/img.png]"   ` | Edge: Quoted with whitespace | `assets/img.png` | `hasBracketWrapper: true`, `isInsideRepo: true` |
| **TC-EDGE-02**| `file:///assets/icon.svg` | Edge: Relative File URI | `assets/icon.svg` | `isAbsolute: false`, `isInsideRepo: true` |

---

## 3. Safe Image Cat Verification Matrix

Evaluates behavior when `gitmap cat` is invoked on varying file types and output targets:

| Test Case | Target File Type | Flags / Output Target | Expected Behavior | Safety Notice Displayed |
| :--- | :--- | :--- | :--- | :--- |
| **TC-CAT-01**| PNG (`image/png`) | None (Interactive TTY) | Intercept binary, render Image Metadata Card | Yes (`[SAFE IMAGE CAT: ...]`) |
| **TC-CAT-02**| JPEG (`image/jpeg`)| None (Interactive TTY) | Intercept binary, render Image Metadata Card | Yes (`[SAFE IMAGE CAT: ...]`) |
| **TC-CAT-03**| GIF (`image/gif`) | None (Interactive TTY) | Intercept binary, render Image Metadata Card | Yes (`[SAFE IMAGE CAT: ...]`) |
| **TC-CAT-04**| WebP (`image/webp`)| None (Interactive TTY) | Intercept binary, render Image Metadata Card | Yes (`[SAFE IMAGE CAT: ...]`) |
| **TC-CAT-05**| SVG (`image/svg+xml`)| None (Interactive TTY) | Stream XML source text directly | No (Vector text is terminal-safe) |
| **TC-CAT-06**| Go / MD (`text/plain`)| None (Interactive TTY) | Stream plain text directly via stdout | No (Text file) |
| **TC-CAT-07**| PNG (`image/png`) | `--raw` flag supplied | Stream verbatim binary bytes to stdout | No (Explicit override) |
| **TC-CAT-08**| PNG (`image/png`) | Piped (`gitmap cat img.png > out.png`) | Stream verbatim binary bytes to pipe | No (Non-interactive terminal) |

---

## 4. Skills & Documentation Synchronization

GitMap skills in `.agents/skills/` guide autonomous agents on how to stream, search, and view files. When agents need to inspect image assets or handle paths from user messages:

1. **`cg-execute-in-below-steps` Skill:**
   - Update `gitmap cat <filepath>` documentation to note that `gitmap cat` natively accepts conversational paths (`@[...]`, `![...](...)`, `1. @[...] # ...`).
   - Highlight that `gitmap cat` is safe for image files and will not corrupt terminal contexts during agent tool execution.

2. **`gitmap` Skill (`.agents/skills/gitmap/SKILL.md`):**
   - Document the flexible path parsing behavior and safe image metadata inspection mode.
   - Explain `--raw` usage for extracting image bytes into pipelines.

3. **Fleet Synchronization Protocol:**
   - Use `03-ai-scripts/38-sync-prompts-skills-scripts.py` to synchronize updated skills across connected repositories.

---

## 5. Backward Compatibility & Release Governance

- **Zero Regression on Clean Paths:** Plain file paths (e.g., `gitmap cat readme.md` or `gitmap cat cli/main.go`) experience zero overhead or behavior alteration.
- **Zero Dependency Growth:** The path parsing and image inspection modules rely strictly on the Go standard library (`image`, `image/png`, `image/jpeg`, `image/gif`, `net/url`, `path/filepath`, `crypto/sha256`), introducing zero third-party dependencies.
- **Minor Release Alignment:** These additions constitute a backward-compatible feature expansion, designated for the next minor release ceremony.
