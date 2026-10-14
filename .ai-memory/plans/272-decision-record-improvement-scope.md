# Decision Record — 2026-10-15: "What else can we improve" follow-up

## Owner decision (explicit, in chat 2026-10-15 ~07:42 +08)

From the 6-item improvement list, the owner approved **items 1–3 only**.
Items 4–6 were **explicitly NOT approved** — do not work on them, do not
propose them again without new evidence.

## Approved (program 272)

1. **Deprecation enforcement** — `warnIfDeprecated()` exists but zero commands
   use `markDeprecated`. Find the first real renamed-command candidate in the
   last 20 versions and mark it. If no genuine candidate exists, remove the
   policy doc rather than leave decoration.
2. **Coverage gate** — CI runs coverprofiles but has no floor. Add a coverage
   floor gate (initial: fail if coverage drops more than 2% vs baseline).
3. **`nc` codegen** — replace the hardcoded catalog with a `go:generate` step
   that emits the catalog from dispatcher registration at build time.

## NOT approved — do not touch

4. Pipeline DB retention policy — NOT approved.
5. "Staging repo" philosophy decision — NOT approved.
6. Actionable error messages overhaul — NOT approved.

## Standing note

An explicit "never approved" from the owner is a hard stop, not a backlog.
These three items stay parked until the owner reopens them.
