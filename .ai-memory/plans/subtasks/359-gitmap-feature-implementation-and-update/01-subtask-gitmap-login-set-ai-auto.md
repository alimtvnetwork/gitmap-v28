# Subtask 01: GitMap `login set ai-auto` Implementation

- **Parent Plan:** `359-gitmap-feature-implementation-and-update`
- **Slug:** `gitmap-feature-implementation-and-update - Task 01`
- **Status:** `pending`
- **Assigned Subagent:** Worker 01

---

## Deliverables
1. Update `cli/cmdlogin/login_cmd.go` to intercept `args[0] == "set"` and dispatch to `runLoginSet(args[1:])`.
2. Implement parsing in `cli/cmdlogin/login_set.go` or `cli/cmdlogin/login_cmd.go`:
   - `gitmap login set ai-auto` -> `gitmap automation bot - runs autonomously`
   - `gitmap login set ai-auto <agent-name>` -> `gitmap automation bot - <agent-name> - runs autonomously`
   - `gitmap login set global ai-auto` -> sets global git identity to `gitmap automation bot - runs autonomously`
   - `gitmap login set global ai-auto <agent-name>` -> sets global git identity to `gitmap automation bot - <agent-name> - runs autonomously`
   - Supports `--global` flag and `--restriction-time` flag.
3. Update `cli/cmdlogin/login_ai_auto.go`:
   - Correct typo in `aiAutoBotName`: `"gitmap automation bot - runs autonmously"` -> `"gitmap automation bot - runs autonomously"`.
   - Add `isAIBotIdentity(name string) bool` to guard against backing up bot names as original human user identities.
   - Update `updateAndActivateAutobotProfile` to store the active profile.
4. Add unit tests in `cli/cmdlogin/login_ai_auto_test.go` verifying all argument permutations and name formatting.
5. Coding Guidelines Compliance:
   - Functions <= 15 lines.
   - Zero nested ifs.
   - Positive boolean naming (`isGlobal`, `hasAIAuto`, `hasRestrictionTime`).
   - Leaf error sovereignty.
