# Subtask Plan 06: yt-manage Functional Verification

- **Subtask Slug:** `06-yt-manage-functional-verification`
- **Parent Task:** `316-repository-management-and-testing-instructions`
- **Target Repository:** yt-manage

---

## 1. Overview & Objectives

Answer the user's question — does yt-manage work? — with evidence:

1. Map the CLI: entry points, command surface, and what each command needs
   (OAuth connect, local mirror, batch edit, scheduling, comments, API/web UI).
2. Determine what can run headlessly in this VM (no interactive OAuth) and verify
   everything verifiable: build health signals from CI if present, static wiring
   of commands to handlers, documented flows vs implementation.
3. Check for any gitmap-side YT integration ("gitmap YT manage") in gitmap-v28.

## 2. Method

- Search via GitMap (`gitmap aum search`, `gitmap find`, `gitmap cat`) inside
  yt-manage (and gitmap-v28 for the integration check). No rg/grep, no git
  commands, read-only.

## 3. Acceptance Criteria

- Verdict in one paragraph: works / partially works / broken, with the evidence
  split into "verified working", "verified broken", and "cannot verify headlessly
  (needs user OAuth)".
- Any defect found is cited with relative paths.
