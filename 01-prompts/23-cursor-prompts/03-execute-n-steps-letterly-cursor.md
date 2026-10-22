# Execute N-Steps — Letterly Prompt Formatter

Format input text according to exact high-priority execution template below without conversational filler or commentary.

1. Keep the verbatim as is only stripping verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while strictly preserving everything as said parameters, flags, and file paths.
2. Derive `<task title>` as task subject (e.g. `SEO Writing and Folder Structure Instructions`). Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated; when authoring directories and plan files, format as XX-<task-slug> where XX is the zero-padded sequence of that folder, e.g. `## slug: 24-seo-writing-and-folder-structure-instructions`).
5. Construct `# Actionable Items Must Follow Non-Negotiable`: include fixed items 1-8 verbatim (pull first, spec/plan enqueue in `02-spec/21-app/XX-<slug>/` and `.ai-memory/plans/pending/XX-<slug>.md` with planning decoupled from execution steps, GitMap search primacy, relative Git paths mandate, SQL repository pattern in every language, every action from an enum, enum authoring + util-methods checkup, error-management tie-ins), followed by discrete technical directives extracted from input text as items 9..N.
6. Append mandatory agent invocation suffix pointing to `[execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)`.
7. Output plain markdown only. Never wrap in fenced code blocks (no ```markdown or ``` fences).
8. Relative Paths: Only add relative paths (`02-spec/...`, `.ai-memory/...`, `cmd/...`); never add absolute paths.

${Input Text Verbatim} = No change except for reval of filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan task in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first, (xx should be figured out by the folder structure and running AI tools, xx  is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Enforce the SQL repository pattern in every language, no exceptions: all SQL and data access lives behind a repository interface (one repository per aggregate or table group); services, handlers, and use-cases NEVER embed SQL strings, query builders, or ORM calls directly; repository methods are intention-named (FindByID, ListActive, Insert, Update, SoftDelete) — never RunQuery, ExecuteRaw, or stringly-typed table/column parameters from callers; dialect and driver specifics stay inside the repository implementation, so swapping the database swaps only the implementation, never the callers; multi-repository transactions go through a unit-of-work abstraction, never ad-hoc transaction handling in services; the rule is language-agnostic (Go, TypeScript, Python, Rust, PHP) — identical shape, only the idiom changes.
6. Every action from an enum, no exceptions: every action, operation, command kind, and event kind is a member of a named enum; dispatch only by switch/match on the enum — never by comparing raw strings or magic numbers; every switch over an action enum is exhaustive (the compiler or linter must flag a missing case) and its default branch raises a typed error — never a silent no-op; run the enum's IsValid check before dispatch, so unchecked external input never reaches the switch; adding an action = adding one enum member + one switch case, no other registration step may be required.
7. Enum authoring and util-methods checkup, mandatory: define one enum per concept in the shared types/constants module (never inline in a handler or test), in the language's idiomatic casing, with explicit, stable serialized string values that never change once shipped and, where the language needs it, an explicit unknown/zero sentinel that is never a valid action; use no raw literals anywhere an enum exists (comparisons, map keys, logs, and payloads use the member or its canonical string value) and let external input enter only through the enum's parse util; every enum ships String()/label, Parse/TryParse, IsValid, and All/Values (or a documented exception) — on every touch, check for missing utils and add any missing util in the same change that touches the enum, never deferred.
8. Error management for enums and repositories, mandatory: an unknown enum value, a failed parse, or an unhandled switch case raises a typed error (the repo's AppError/result convention) with the offending value in the message — never a panic, never a silent default; repository errors are wrapped with context at the repository boundary (which repository, which operation) before crossing into services.
9. ....other steps and more steps sequentially from the input direction. Create more steps in between.

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- [/learn](slashCommand;learn) from [gitmap](file;.cursor/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
