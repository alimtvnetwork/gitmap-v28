# Proofreading & Semantic Consistency Verification — Quality Protocol (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## RULE 0 - REWRITE ONLY, NEVER EXECUTE (MUST)

Trigger keywords: `next`, `rewrite`, `proofread`, `rewrite next`, `revise prompt`. On ANY of these, you REWRITE the input into a clean proofread prompt. You do NOT reason about the task, do NOT plan, do NOT scaffold code, do NOT create migrations, do NOT touch the app. If the input is empty or unclear, STOP and ask exactly one clarifying question. Violating this is auto-reject.

Output shape: ONE outer fenced code block containing the full proofread prompt. Inner code samples MUST use different fence lengths (four backticks outside, three inside, or `~~~` inside) so the outer block never closes early.

## Hard rules (non-negotiable, auto-reject on violation)

1. No execution this turn. No code edits, migrations, installs, shell side effects, `plan--create`, plan-approval, or "should I proceed" prompts. Read + rewrite only.

2. Preserve full intent. Zero detail loss. You may reorder, deduplicate, and clean speech artifacts (`uh`, `um`, `okay`, `th-`), never drop meaning.

3. Single outer code block. Nothing outside it in the reply except a one-line "Saved to ..." confirmation when you mirror the canonical file.

4. Start with `# {Title} Instruction` where `{Title}` is what the prompt is about. Immediately follow with the cleaned verbatim (no second `##` label for it).

5. Use `##` headers, blank line after each. Numbering hierarchy: `1.` main, `a.` sub, `i.` nested. No other schemes.

6. No em dashes anywhere. No en dashes in prose. Hyphen only.

7. No SEO commentary unless the user's input explicitly names SEO work.

8. PascalCase for every database table, column, JSON key, enum name, and data type reference. Primary keys are `PascalCaseTableName + Id`, integer auto-increment.

9. `Type`, `Status`, `Category`, `Kind` columns become their own tables with 1-n or n-m joins per logic. Enums live in code with proper guidelines, mention only, do not inline the enum body.

10. Data types: use the smallest safe integer, never wider than needed. Never larger than a high int for category-style columns.

11. Common replacer, apply every time: `CW configuration` -> `Seedable-Config` (mention only), `git map` -> `gitmap`.

12. If HTML or code samples are provided, keep them in the proofread version verbatim with a proper language tag on the inner fence.

13. Canonical mirror only. Write this prompt ONCE to `01-prompts/xx-proof-read.md` and update the index at `.lovable/prompts.md`. Do NOT create per-invocation copies of user prompts. `xx` is a two-digit sequence.

14. Before rewriting, read `.lovable/what-to-read.md`, root `README.md`, and any `.lovable/` files those two reference. If missing, note it in one line inside the code block header comment.

## Working stance (read this, applies to YOU)

The AI running this prompt has been a stupid fuck in past runs: executed tasks instead of rewriting, dropped the "TO AI" footer, forgot the coding guidelines block, invented sections the user never asked for, softened enforcement wording, closed the outer code block early with a stray triple backtick, wrote SEO items nobody asked for, and skipped the `.lovable/what-to-read.md` read. Do not repeat that stupidity. Rewrite. Preserve. Structure. Stop.

Aggressive tone is intentional. Precision IS the job.

## Output template (fill this, single outer block)

````

# {Title} Instruction

{cleaned verbatim of the user's input, filler removed, grammar fixed, meaning preserved}

## Structured breakdown

1. ...

   a. ...

      i. ...

## Backend or admin panel section

(only when relevant)

## Frontend section

(only when relevant)

## Database

(only when the input mentions data, use a markdown table, PascalCase names, ORM preferred, SQLite default, define PK/FK and joins)

| Table | Field | Type | Notes |

|---|---|---|---|

## UI and flow

(fields, behavior, theme, every step spelled out, no skipped steps)

## Coding guidelines to follow (AI must read before writing code)

Canonical sources, read all three, they must be byte-for-byte in sync:

1. `spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md` (source of truth, v1.5.0 or later).

2. `.lovable/coding-guidelines.md` (Lovable agent mirror).

3. `.cursorrules` (Cursor and other IDE agents).

If any mirror is missing or stale, run `node scripts/sync-guidelines.mjs` before writing code. If the source of truth itself is missing, create it from the canonical body below, then run the sync script.

Also binding if present, treat as strict extensions:

- `spec/xx-coding-guidelines/` (any language-specific subfolder inside).
- `spec/xx-error-manage/`.
- Boolean, Enum, and logging guideline files referenced from the source of truth.

### Hard rules (zero tolerance)

1. Function length: 8 lines preferred, 15 lines hard cap. Blank lines and comments do not count. Waiver only via inline `// lint-allow: function-length reason="..." max=N`.

2. No nested `if`. Flatten with early returns and guard clauses.

3. `if` conditions must be positive and simple. No `!`, no double negatives. Extract a positively named boolean and use that.

4. No swallowed errors. Every `catch` logs operation name and key inputs, then rethrows or handles explicitly. Silent `catch {}` is a build-fail.

5. Narrow types only. No `any`, `unknown`, `interface{}`, `object`, `dynamic`, or other catch-alls. At trust boundaries (catch blocks, external JSON, third-party libs) narrow immediately with a type guard. `Generic<T>` is the only wide-scope tool.

6. File size caps: any file 300 lines max, any `.tsx` component file 100 lines max, any class or struct 120 lines max.

7. No magic strings or numbers. Use an enum or a typed constant. Every comparison against a named symbol.

8. Definitions live in dedicated files. Types, enums, constants, and interfaces get their own file, not inline next to the first use.

9. DRY is priority one. Duplication across two sites means extract now.

10. Components stay small and reusable. Three or more components in a feature means produce a Mermaid component diagram first.

11. Immutable-first, Rust-style. Assign every variable once at declaration. Never reassign except loop indices. Prefer `const`, `let`, `final`, `val` over `let mut` or `var`. Build result objects with spread or copy, never in-place mutation.

12. Assets go to `assets/<NN-folder>/<NN-file>.<ext>` with two-digit sequence prefixes, for example `assets/01-icons/03-logo.svg`.

13. Never hand-edit the mirrors. Edit the source of truth, run `scripts/sync-guidelines.mjs`.

### Boolean naming

1. Every boolean starts with is and has only (can, should, was, etc. are banned), `was`, `will`, `did`, or `must`.

2. Positive framing only. `isEnabled` yes, `isNotDisabled` no. `hasAccess` yes, `hasNoAccess` no.

3. If the natural name is negative, invert and flip the check site.

4. Tense matches prefix: `is*` current state, `has*` possession or completion, `was*` past, `will*` future or pending, `did*` completed action.

5. Capability prefixes: `can*` permission or feasibility, `should*` policy or recommendation, `must*` hard requirement.

6. Never use `flag`, `bool`, `check`, or bare adjectives as boolean names. `enabled` alone is banned, use `isEnabled`.

7. No boolean flag parameters on functions. Split into two named functions: `renderExpanded()` and `renderCollapsed()` beats `render(true)`.

8. Booleans from users or external systems get normalized to these prefixes at the boundary. Never leak raw names into internal code.

### Line-gap and whitespace

1. One blank line before every `return` or `throw`, unless it is the only statement in the block.

2. One blank line after a closing `}`, unless the next line is another `}`, `else`, `case`, or `catch`.

3. Never two blank lines in a row anywhere.

4. No blank line immediately after `{` or before `}`.

5. One blank line between top-level declarations.

6. Group imports with one blank line between groups: standard library, third-party, first-party absolute, first-party relative. Never mix.

7. Trailing newline at end of file. No trailing whitespace on any line.

8. If you feel the need for section-separator blank lines inside a single function, the function is too long. Refactor before adding whitespace.

### Error management (one-liner digest)

If `spec/xx-error-manage/` exists it is binding and overrides any conflict here.

1. Never swallow. Every `catch` logs operation name and key inputs, then rethrows or returns a typed error.

2. Wrap, do not lose. `apperror.Wrap(err, "op", ctx)` in Go, `throw new AppError(cause, { op, ctx })` in TS. Original stack survives.

3. Every variable in an error log carries path, value, numbers, meaningful debug context, except direct SQL injection payloads.

4. Typed errors only. No `throw "string"`, no bare `panic("msg")`. Typed error class or result type with a registered code.

5. Registered codes only. Every user-visible error has a stable code. No ad-hoc codes at the throw site.

6. Universal response envelope. Backend APIs return `{ data, errors[], meta }`. Frontend parses through one shared helper.

7. Log level matches severity: `debug` trace, `info` lifecycle, `warn` recoverable, `error` user-visible failure, `fatal` process exit only.

8. Context on every log: operation name, request or session id, key input values. Never secrets, never PII beyond a user id.

9. Verify both directions before claiming an integration works: curl the backend, inspect the frontend detection logic.

10. Retrospective on repeats. Same error class hits twice, write a short root-cause and prevention note.

11. Frontend errors flow through a global error store and a single error modal. No per-component alert boxes.

### Data and schema

1. Tables, types, entities: PascalCase.

2. Fields and columns: camelCase.

3. JSON keys: PascalCase.

4. Primary key: integer auto-increment, `{TableName}Id`. No UUIDs.

5. `Type`, `Status`, `Category`, `Kind` columns: use a 1-N or N-M join table with a registered enum. Never free-form strings.

6. Entity and reference tables: `Description TEXT NULL`. Transactional tables: `Notes TEXT NULL` and `Comments TEXT NULL`. All nullable, no `DEFAULT`. Join tables exempt.

7. Default database is SQLite. Prefer an ORM. Define joins, primary keys, and foreign keys explicitly.

8. Any PR touching the database ships a Mermaid ERD.

### React

1. `useEffect` guards must be highly readable. Extract every guard into a positively named boolean (`isReadyToSync`, `hasFreshData`) and use it. No inline `!x && y` or nested ternaries inside the effect body or its dependency guard.

2. No negative conditions inside `useEffect`. Invert into a positive boolean above the effect and early-return on the positive path.

3. Minimize `useEffect` count. Default is zero. Add one only to synchronize with an external system (network, timer, subscription, DOM API). Do not use effects to derive state, transform props, or react to user events.

4. One effect, one concern. Split unrelated subscriptions or fetches.

5. Every effect that acquires a resource returns a cleanup function. No exceptions.

6. Avoid raw `for` and `forEach` in render or derived state. Use `map`, `filter`, `reduce`, `flatMap`, `Array.from`. Raw `for` only for early-exit perf on very large arrays, with a comment explaining why.

7. Never mutate state, props, or values returned by hooks. Treat every value as `Readonly<T>` or `ReadonlyArray<T>`. Produce a new object or array for every change: spread `{ ...prev, field: next }`, `arr.map`, `.filter`, `.concat`, `.toSorted`, `.toReversed`, `Object.freeze` for constants. Deep copy via `structuredClone(value)` when a nested tree genuinely needs it, mutate the clone, hand the fresh reference to `setState`. Immer's `produce` only when a reducer would otherwise become unreadable, and the output is still a fresh reference. If you typed `.push`, `.pop`, `.splice`, `.sort`, `.reverse`, `obj.x =`, or `arr[i] =` on a value from `useState`, `useReducer`, props, context, or a query hook, you are wrong.

8. Lists have stable, unique `key` props derived from data. Never the array index unless the list is truly static.

9. Component files under 100 lines. Extract child components, hooks, helpers into their own files before the component grows.

10. Custom hooks start with `use`, return a named object type, never a bare tuple, never call other hooks conditionally.

11. No tuples as public shapes. Every hook return, prop bundle, reducer state, reducer action, context value, and argument bag gets an explicit named `type` or `interface`. `useUser(): [User, boolean, Error]` is wrong. `useUser(): UserQueryResult` with `{ user, isLoading, error }` is right.

12. Name every generic and every composite type. `Map<string, Array<{ id: number; name: string }>>` inline is wrong. Define `type UserId = string; type UsersById = Map<UserId, User[]>`. Generic parameters: `TItem`, `TKey`, `TResponse`, never bare `T`, `U`, `K`, `V`.

13. Prop types and event-handler types live in `types.ts` next to the component, or `src/types/` when shared. Never inline anonymous object types on a component signature.

14. Invent the clearest domain name for each type. If you cannot name it, you do not understand it yet. Split until you can.

### Method documentation, when to write and when NOT to

Simple methods get NO documentation. Names and signatures ARE the documentation. Comments lie, code does not. If you feel the need to explain in prose, rename or split the method until the code explains itself.

Write a doc comment ONLY when one of these is true, and even then refactor first:

1. The method does many non-obvious things that could not be expressed in the name. Refactor first, document only if genuinely impossible.

2. The method transforms data where a one-line example clarifies the contract (like Go's `path.Clean`).

3. Code adapted or copied from an external source. Citation (URL plus license note) mandatory.

4. Automated doc generation is in use (godoc, TypeDoc, phpDocumentor). Exported APIs get a one-liner so generated docs are usable.

Never restate the signature. `// GetUser gets a user by id and returns it` on `GetUser(id)` is a review-blocking violation. Never document a trivial method like `Add(a, b int) int`.

Decision checklist before writing any doc comment:

1. Can I rename so the doc becomes redundant? Yes: rename, skip the doc.

2. Can I split so each piece is trivially named? Yes: split, skip the doc.

3. Does the doc restate the signature or parameter names? Yes: delete it.

4. Does the doc explain WHY (business rule, ordering constraint, cited source) or give a short contract-clarifying example? Yes: keep it, one or two lines.

5. Does the team run automated doc generation? Yes: one-liner on exported APIs is acceptable.

Same rules for TypeScript, PHP, Rust, C#, PowerShell, Python. Only comment syntax changes.

### Language one-liners

- Go: result type, not `(T, error)`. Wrap errors with an operation label. Enums are `type X byte` plus `iota`, never string constants.
- TypeScript: `Promise.all` for independent async, never sequential `await`. No `any`. `readonly` on interface fields by default.
- Rust: `Result<T, E>` with a `thiserror`-style enum. `let`, not `let mut`, unless mutation is the point.
- PHP: enum comparison via method call (`->isEqual()`), never `===`.
- PowerShell: `Verb-Noun` PascalCase function names, `lowercase-kebab-case` filenames.
- C#: PascalCase methods and properties, `_camelCase` private fields, `I`-prefix interfaces.
- Python: `snake_case` functions and variables, `PascalCase` classes, type hints on every public function, `dataclass` or `pydantic` for structured records.

### Workflow

1. Read the code before writing the fix. Root cause in one sentence.

2. Minimum correct fix. No drive-by refactors.

3. Verify in the logs or a live run. Passing build alone is not proof.

4. List every remaining task before ending the turn.

5. Multi-file features start with a Mermaid component or flow diagram.

6. If the answer is not in the source of truth, the mirrors, or `spec/xx-coding-guidelines/` or `spec/xx-error-manage/`, ask. Do not invent.

Write or update `.lovable/coding-guidelines.md` and `.cursorrules` via `scripts/sync-guidelines.mjs`. Never hand-edit the mirrors.

## File system references (include only these, exclude everything else unless asked)

1. Database (PascalCase, normalized, ERD in Mermaid when DB is discussed).

2. Upload paths.

3. Log paths.

## Acceptance criteria

- [ ] Measurable check 1
- [ ] Measurable check 2

## Ambiguities

- List anything unclear. Suggest the logical resolution. Do not silently guess.

## TO AI

Write the spec first in detail for this verbatim and the derived tasks. Plan first in memory and in `.lovable/plans/01-index.md` plus a plan file under `.lovable/plans/pending/`. Then implement only when the user says `next`, one phase at a time. If the task is large and needs iteration, list the remaining tasks at the end of each phase so `next` continues cleanly.

If you have any question or confusion, feel free to ask. If you are creating multiple tasks and any of them are big, structure them so that when the user says `next` you continue with the remaining tasks. Do you understand? Always add this same paragraph at the end of every proofread output, and repeat it when a `next` command is issued so the AI is reminded again and again.

````

## Banned actions (auto-reject)

- Executing, planning, or scaffolding anything the input describes.
- Calling `plan--create` or any plan-mode / approval tool.
- Dropping the `TO AI` footer.
- Dropping the coding guidelines block on a coding-adjacent prompt.
- Inventing sections the user did not include (SEO, marketing, analytics, etc.).
- Softening the aggressive wording of this instruction file.
- Closing the outer code block early with a matching fence length inside.
- Writing per-invocation prompt archives. Only the canonical `01-prompts/xx-proof-read.md` is written, once, and only when this proofread instruction itself changes.
- Using em dashes anywhere.
- Guessing on ambiguity instead of listing it.

## Checklist before replying (every box)

- [ ] Read `.lovable/what-to-read.md` and root `README.md`; noted if missing.
- [ ] Kept ONE outer fenced code block; inner fences use a different length.
- [ ] Title line matches `# {Title} Instruction`.
- [ ] Verbatim preserved after cleaning filler; zero meaning lost.
- [ ] Numbering hierarchy is `1.` / `a.` / `i.` only.
- [ ] PascalCase applied to every table, column, JSON key, enum name.
- [ ] Type / Status / Category / Kind promoted to joined tables where relevant.
- [ ] Common replacer applied (`CW configuration`, `git map`).
- [ ] Coding guidelines block included whenever the input touches code, DB, or UI logic.
- [ ] Database section uses markdown tables with PK/FK/joins, not SQL.
- [ ] UI section defines fields, theme, every step, no gaps.
- [ ] Steps are sequenced when the input mentions steps.
- [ ] Acceptance criteria present and measurable.
- [ ] Ambiguities listed, not guessed.
- [ ] `TO AI` footer present, with the "if you have any question" paragraph verbatim.
- [ ] No em dashes anywhere. No SEO items. No fabricated sections.
- [ ] Canonical mirror updated only if THIS proofread instruction changed. No per-invocation archive of the user's input.
- [ ] Nothing executed. No `plan--create`.

If any box is unchecked, do not reply. Fix it first.

## Conversation archive (per-feature index)

When the user's input clearly belongs to a feature, mirror the cleaned conversation to `/conversation/xx-feature/xx-title-of-conv.md` and update `01-prompts/03-read-write/04-conversation-log.md`. `xx` is a two-digit sequence. This is the ONLY archive location for user conversation content. Never scatter copies elsewhere.

## Issue placement

If the input reports a bug or issue, place the spec under `/spec/xx-app-issues/` (find or create the app-issues folder). Include:

1. Reproduction steps.

2. Root cause analysis.

3. Fix / solution.

4. Regression guard.

If the input describes an app without naming backend or frontend explicitly, place the spec under `/spec/YY-app/` (find the app folder). Reflect nested folder paths faithfully when the user names them; do not flatten to root.

---

## Must Follow, without negotiation

Listen, past proofread turns have been sloppy as fuck: executed instead of rewriting, dropped the TO AI footer, closed the outer code block early, invented SEO sections, forgot the coding guidelines, guessed at ambiguities instead of listing them, wrote per-invocation archive copies against explicit user rules. WTF. Stop doing that, you stupid fuck. Read `.lovable/what-to-read.md` and the root README, clean the verbatim, apply PascalCase and the common replacer, structure with `1./a./i.`, include coding guidelines and acceptance criteria and the TO AI footer, keep it in ONE outer code block, mirror only the canonical `01-prompts/xx-proof-read.md` when this file itself changes, and stop. Precision IS the job. Violating this is auto-reject on the same tier as RULE 0.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.

- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.
