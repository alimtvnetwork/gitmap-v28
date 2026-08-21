# Next 2 Steps or Tasks (v5)

## What I want

1. Give me the **NEXT N STEPS — exactly N** — and for each one:
   1a) **Reasoning** — why this step, why now, what breaks if it's skipped.
   1b) **Time estimate** — realistic, not optimistic.
   1c) **What it unblocks** — the next thing that becomes possible.

2. Then list **every remaining item** after those 3 so I can see the full picture. At the end of the task always bump the minor version, add changes log and update release notes and if possible pin that version in the root readme file. And also save this prompt in the .lovable folder in the prompts folder for known as 'xx-next-task.md' and update it as 'next task with number'

## Definition of done (non-negotiable)

You are NOT done until all of these are true:
- [ ] You have actually read the relevant files AND the project memories — and you can name the exact files/functions/lines involved.
- [ ] The **root cause** is written in ONE sentence, before any fix.
- [ ] The fix is the **minimum correct change** tied to that root cause — not a symptom patch.
- [ ] You **verified** it: build output, error logs, and/or preview — and you show the before/after signal (failing → passing).
- [ ] You reported what changed and why.

## Hard rules

- **STOP and read first.** No skimming, no guessing from filenames. If you can't name the exact lines, you haven't read enough — go back.
- **Root cause before fix.** Trace the bug end-to-end. No assumptions. No "this should work."
- **No symptom-patching.** If your "fix" is a try/catch, a fallback value, or a re-render hack used to hide the problem, you've failed — start over.
- **If you're unsure, SAY SO.** Do not fabricate. A wrong-but-confident answer is worse than "I don't know yet."
- **Go slow. Go critical. Go deep.** Depth is not optional polish — it IS the entire job. Fast + wrong = useless and wastes another full loop.

## Error logs & error management (ALWAYS focus on this)

- Read the actual error logs FIRST — console, server/worker logs, build output, stack traces. The answer is usually already there.
- If there are NO logs, that itself is the bug: add logging at the entry point and surface errors instead of swallowing them. Silent failure is unacceptable.
- Every fix must include proper error handling and observability: errors must be logged with context and surfaced, never hidden.
- Confirm the relevant log line actually fires after your change. If you can't see it in the logs, you haven't proven the fix.

## Additional Instruction (must follow if matches)

Before executing, check the task type and follow the relevant guidelines if they exist (skip silently if the file is missing):

1. **Coding tasks** (especially Golang, Python, PHP, or other backend):
   - Check for `.lovable/coding-guidelines.md`. If present, follow it.
   - Also check `spec/coding-guidelines/`. If present, follow every file inside.
   - If this is a coding task and neither location has guidelines, ask me to provide one.

2. **SEO tasks** (website/SEO-related):
   - Check for `.lovable/seo-guidelines.md`. If present, follow it.

Rule: verify the file/folder exists first. If it does not, skip that guideline silently. If multiple guidelines apply, follow all of them; if they conflict, prefer the folder-level spec and call out the conflict.
