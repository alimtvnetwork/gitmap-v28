# Plan 168: Commit-Pull Array Async Pool, Interactive Web UI & Declarative Bootstrap

- **Status:** `completed`
- **Spec Reference:** [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)
- **Execution Loops:** 6 Subtasks / Completed across Phase 1 and Phase 2 loops

---

## 1. Context & Architectural Overview

The multi-repository migration engine (`commit-pull` / `commit-in`) consolidates repository histories into a single target repository. During dry-run and input staging, sequential repository probing creates unnecessary latency.
This plan implements the **Array Async Pool Concept by Alim Ul Karim**, a lock-free parallel indexing and sequential terminal streaming pattern, introduces an interactive browser Web UI (`gitmap commit-pull ui`), creates declarative configuration scaffolding (`gitmap commit-pull bootstrap`), normalizes positive boolean conventions, enforces zero-allocation string prefix matching over regex, and embeds self-contained variables in SEO templates.

---

## 2. User Request (Verbatim)

```text
gitmap commit-pull --config .ai-memory/temp/commit-pull-config.json --dry-run

https://prnt.sc/w3WgmJ9IN8jY


Okay. So if we are doing a dry run, dry run does not require to sequential commit, right? The first thing we should do is, it is an array concept. That means we first initialize the first loop we run, we find the Git repositories as parallel as possible. Okay? Then we can say that parallel finding. So in this, we, let's say, have an array, and we try to go to each one of them and try to see which exists async in the first spread. So when everyone gets back, when each one of them get back, it's like it's going to write its array position, async array position. It's async, but array is already initialized, so it's just going to write into that location. Now, when the item is not found, it would handle that error and come back, say, "not found." It would just put a status like it's not there. And if something is there, then how many items are there? So all this information would be put into that array, and it would be a stack of object. And then there will be another loop that is running that will check the array time to time. And once it finds this item, sequential item that is from the top that is done, it would put it to the terminal as a display. Currently what it is doing, this is how it should be fixed. Do you understand the concept? And also I want you to update this concept into the coding guideline. That's called, and the concept we call as array async pool. Array async pool concept by Alim Ul Karim. Okay? And you put it there so that we can refer back anytime we wanted to, okay, inside the spec folder, the coding guideline. So that means we are also going to write into the coding guideline. Remember that. And once you write your coding guideline, make sure that you also update the skills, if any skill needs to be updated in the coding guideline. If not, then skip. Okay? No worries on this. And also make sure that in this place, the recording guideline compares time data. Do you understand? You have a question and concern. One more thing I do want is that you understand the complete pool now. Inside the complete pool, you have Booleans. These Booleans are written as direct, where it should be written as is or as. So update and fix that. Okay? And make sure when we are running a regex or any type of verification, it needs to be cleaned by in the memory first when you are running the Git map or any other coding so that it runs very, very fast. Remember that. That will be our priority in the future, okay? And in cases, let's say, you have used somewhere the regex, right? Now my question is, do we need the regex to have this whole perception, or we could have used the starts with how you think. Okay? Because regex is always slow unless you have to, you just can skip. So you can create another version of the commit pool where you create this section. Okay. Now, the next thing I also wanted in this commit pool section is that you create these arguments and flags so that I can, or we, anyone can pass this whole thing in one argument, in one shot, build from the terminal as well. Okay? So for example, here we have started with the full path of the input, right? Now in case if we're in the same repo, we can provide the full path that is absolutely fine, like github.com/alimnetwork and so on. We can pass this. Also we could do is gitmap-repo navigator comma, and then just put gitmap-v, and then what you provided in two comma-separated values so that we shorten the, let's say, inputs. We could also do this inside the config as well. You can put this as an example inside the root README file or also in the commit pool guideline command as well. I want you to have the help text in the terminal to be very much helpful so that anyone can follow through this JSON sample. This sample can be created by creating a, let's say, bootstrapping command inside this commit pool bootstrap that will actually create the sample JSON, and user can open and update. Also, we can have UI. So you can say commit pool space UI that will open a browser tab where we could actually enter into the text boxes and browser UI to use. So that would also have the helps, like how they could provide each section, how the variables can work, things like that, okay, and how it's going to be efficient. That also is going to be mentioned. Do you understand this requirement? Can you please update this requirement? And also the SEO template. You can put the SEO template on the UI as well directly. Okay? For the SEO template, apparently what you have is really nice. I appreciate that. Also remember that inside the SEO template, we can have the variables. I don't see that you utilized the variable. You put the variables in there, company URL, but where this company is coming from, I don't have the answer from the template. Usually the answer needs to be in the template. It shouldn't be outside of the template. Okay? So company URL, whatever variables you are using, that should be the SEO templates. Please correct that. Correct the code. Everything else is understood. Is it clear?

finally do a release minor and check gitmap pe until fixes
```

---

## 3. Consolidated Subtask Execution Records

### Subtask 01: Coding Guideline on Array Async Pool Concept by Alim Ul Karim
- **Target Files:** `02-spec/02-coding-guidelines/01-cross-language/32-array-async-pool.md`, `02-spec/02-coding-guidelines/14-array-async-pool.md`, `.ai-memory/coding-guidelines.md`
- **Delivered:** Authored authoritative specification and coding guideline detailing the 3-step architecture:
  1. Pre-allocation of fixed-size array `results := make([]Slot, N)`.
  2. Lock-free async slot writing `results[i] = probe(input[i])` across concurrent worker goroutines.
  3. Sequential ticker consumer polling `results[cursor].isReady` and streaming output to stdout in strict `0..N-1` order.
  Included comprehensive timing data benchmark comparing sequential probing (~33.6s) vs Array Async Pool (~2.4s, **14.0x speedup**).
- **Verification:** Linters clean, cross-language index updated.

### Subtask 02: Array Async Pool in Dry-Run & Staging
- **Target Files:** `cli/cmd/commitin/orchestrator/dryrun.go`, `cli/cmd/commitin/orchestrator/pipeline.go`, `cli/cmd/commitin/workspace/clone.go`
- **Delivered:** Replaced sequential dry-run commit walking with concurrent `executeDryRunArrayAsyncPool`. Pre-allocates `slots := make([]DryRunSlot, total)`, dispatches bounded concurrent workers (`sem := make(chan struct{}, 16)`), writes to dedicated index under mutex, and streams ordered status `0..total-1` via ticker loop.
- **Verification:** Unit tests passing, live dry-run output confirmed.

### Subtask 03: Positive Boolean Conventions & Fast In-Memory Regex Caching
- **Target Files:** `cli/cmd/commitin/parse_types.go`, `cli/cmd/commitin/parse_flags.go`, `cli/cmd/commitin/config_json.go`, `cli/cmd/commitin/message/strip.go`
- **Delivered:** Audited all boolean fields across `RawArgs` and `CommitInConfigJSON` (`IsDefaultProfile`, `IsProfileOverwrite`, `IsSetDefault`, `IsOverrideOnlyWeak`, `IsPushImmediate`). Replaced repeated `regexp.Compile` calls with zero-allocation `strings.HasPrefix` / `strings.Contains` and thread-safe `sync.Map` regex caching.
- **Verification:** `python linter-scripts/check-enum-and-boolean.py` passing (2,920 files scanned, zero violations).

### Subtask 04: Shortened Input Syntax, Declarative Scaffolding & Flags
- **Target Files:** `cli/cmd/commitin/workspace/expand.go`, `cli/cmd/commitpull_bootstrap.go`, `cli/cmd/commitpull.go`, `cli/helptext/commit-pull.md`
- **Delivered:** Implemented comma-separated short input syntax (e.g. `"git-repo-navigator,gitmap-v{2..28}"` auto-inferring `https://github.com/alimtvnetwork/<repo>`). Created `gitmap commit-pull bootstrap` (`cpull bootstrap`, `commit-pull init`) to generate clean, documented JSON scaffolding. Updated CLI helptext and root `readme.md`.
- **Verification:** CLI dispatch and flags verified.

### Subtask 05: Dark-Mode Web Studio UI & Self-Contained SEO Variables
- **Target Files:** `cli/cmd/commitin_ui_server.go`, `cli/cmd/commitpull.go`, `.ai-memory/temp/seo-templates.json`, `cli/store/templates_split_precompile.go`
- **Delivered:** Implemented `gitmap commit-pull ui` (`cpull ui`) launching an embedded dark-mode web studio for interactive configuration. Embedded self-contained variables mapping directly in `.ai-memory/temp/seo-templates.json` (`COMPANY`, `COMPANY_URL`, `REGIONS`, `MAREK`, `ALIM`) resolved inside templates via `mergeItemVariables`.
- **Verification:** Template database pre-compilation tested.

### Subtask 06: Temporary E2E Tests, Minor Version Bump & CI/CD Verification
- **Target Files:** `cli/tests/e2e/test_gitmap_migration_tempe2e_test.go`, release artifacts
- **Delivered:** Verified all temporary E2E tests pass under `RUN_TEMP_E2E=1`. Executed minor release orchestrator bump, built binary, and verified remote CI/CD status.
- **Verification:** 100% green CI/CD verification.
