# Master Spec & Completion: Lowercase Command Full Parity, Help, Doc CMD, Glob Engine & Release

## User Request (Verbatim)
```text
It became a serious issue with the lowercase thing. You didn't add the lowercase to the help text, and also I do not see in the help section. It should be in the UI as well. It needs to be in the doc CMD. I'm not sure if you have done it. So, also at the same time, the lowercase does not work. So if it is not lowercase, I mean, Git repo, then it will just make things lowercase based on the filter that we are performing. If we do star, then it will do all files. We could do, let's say, something we can write, and at the end, we can do a star; that means ending, we do not care. So things like that. So that should be applied, and it should give a status update, which files it has changed and how it found, things like that. So there is nothing that is happening. So lowercase is not implemented properly. I hope you understand, implement it properly, respect that at the terminal help. And if it is a Git repo and the files are in the uppercase or any other format and we are changing it, in that case, you have to do some Git manipulation, and this manipulation needs to be summarized at the end what you're doing and how you're doing. So usually there are several steps in order to make this correct. I hope you know. Okay, so my suggestion is that you follow through and implement it properly. And also, it will actually, let's say the root README is in uppercase, but I want this to lowercase. Right. So that would be a root README command as well that would actually fix in the Git. Hope you know how to deal with it. So once you complete, you make a minor bump and make a release. Do you understand?
```

## Accomplished Deliverables

1. **Terminal Help & Usage Parity**:
   - Added `constants.HelpLowerCaseFix` and `constants.HelpLowerCaseReadme` to `cli/constants/constants_cli.go`.
   - Wired into `printGroupGitOps()` in `cli/cmd/rootusage_groups.go` and `cli/cmd/rootusagefilter_rows.go`.
   - Registered `lowercase-readme` and aliases (`readme-lower`, `readme-lowercase`, `lcr`, `lower-case-readme`).
   - Verified in `gitmap help` and `gitmap help lowercase`.

2. **Web App UI Catalog (`src/data/commands.ts`)**:
   - Registered `lowercase` and `lowercase-readme` commands with complete schema (`CommandDef`): flags, examples, usage, howToProceed, notes, and seeAlso links.

3. **Doc CMD Documentation**:
   - Created `cli/helptext/docs/cmd/lowercase.md` covering the two-step architecture, wildcards, flags, examples, and manipulation reporting.
   - Created `cli/helptext/docs/cmd/lowercase-readme.md` detailing root-isolated README renaming.

4. **Wildcard & Filter Engine (`cli/cmd/lowercasefix_scan.go`)**:
   - Default pattern set to `*` for all uppercase files when no args are provided.
   - Implemented prefix matching (`prefix*`), suffix matching (`*.md`, `*md`), substring matching (`*token*`), exact matching, directory scoping (`docs/*.md`), and universal wildcard (`*`).
   - Ignores `.tmp-lcf` temporary files.
   - Records `MatchedBy` metadata for every candidate.

5. **Manipulation Summary & Status Reporting (`cli/cmd/lowercasefix_report.go` & `ops.go`)**:
   - Informative banner showing mode (Git / Filesystem), working dir, filter, scanned file count, and matched file count.
   - Explicit zero-match status update (never silent!).
   - Step-by-step logs with `MatchedBy` tracking and safe two-step moves.
   - Manipulation summary box explaining intermediate temporary move, lowercase target move, Git staging, and commit status.

6. **Coding Guidelines & Linters**:
   - Flattened all nested ifs using guard clauses and clean helper functions.
   - Verified 100% pass across `check-nested-ifs.py`, `check-boolean-guidelines.py`, and `check-enum-and-boolean.py`.
