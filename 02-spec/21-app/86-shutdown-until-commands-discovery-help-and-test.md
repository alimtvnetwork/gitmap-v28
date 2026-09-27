# Canonical Spec: Shutdown-Until Command Discovery, Root Dispatch & Rich Help Integration

**Spec ID:** SPEC-APP-86
**Status:** Active
**Author:** Antigravity Master Orchestrator
**Date:** 2026-09-28
**Version:** v6.365.0

---

## 1. User Request (Verbatim)

```text
is it done properly???

Okay. Currently, most of the things are working fine. I think it's good, but I think I do have some requests that you could modify. First of all, the gitmap update all. It could be hyphen all or space all or UA. This would do the same thing. But now coming to the point how it's going to do it, the update is a delicate situation. So if we are running from the external machine, we should actually use JSON as a communication so that the help text does not get included in it. And currently, the way that you are doing the formatting or displaying it looks very poor. Okay? So I think this is where we need to work on. We need to make sure the help and things, these are really high quality. Okay? So currently, the output of the part two, how you're doing it, it's quite poor. Okay? And this is where we need to work on. And also make sure that, at the end, you bump the minor version and make a release and check the Gitmap PE. Also now I wanted you to check, so you can even run a test command on, let's say, scripts picture I or something like this. And you check by the running prompts LS that this project is actually running and other projects, but this project needs to be there. Or you could do other stuff like the own project. Let's say the status, alien status sample, just put a hi there and check if this comes up in the running prompt when it is running or put some weight. Okay? Like sleep and then say hi. You can do that so you can trace back. So that's one thing. Another is backing up the running prompts. So each one of the prompts from all the projects, the running prompts would be backed up. Then you can restore the prompt, running prompt, that should inject and run the prompt that we have taken the backup. So make sure that you do the end-to-end testing so that you can confirm it, that it's working very well and there is no confusion. Okay? It's been several times that we are trying that. I hope there should be no issues, right? If possible, you can also test out with some file the Gitmap SC or deploy command, sync left, sync right command. Okay. And also update the help text and also UI help text. These are very important. Do not miss it. Is it clear? Can you please do that? At the end, do a minor bump and a release and do the Gitmap PE to check the errors. Okay.

Also at the same time, we had a command, shut down until. I could not find that command. Check these commands are in the help text. I could not find it. We cannot test it because it would shut down. Okay, so after you do the release, then you can test the shutdown command. Shut down until the prompts are running, something like this. You can test it out for one or two minutes or something like this. I could not find the command, unfortunately. Yeah, I do not. Make sure you work on it, and you make sure that this is finally done as well. Is it understood? Can you please help me with it?
```

---

## 2. Technical Architecture & Requirements

1. **Top-Level & AGY Command Aliases:**
   - Command names:
     - `gitmap shutdown-until-green` (full name)
     - `gitmap shutdown-until` (user-friendly alias)
     - `gitmap sug` (concise 3-letter alias)
     - `gitmap agy shutdown-until-green`
     - `gitmap agy shutdown-until`
     - `gitmap agy sug`
   - Subcommands supported:
     - `ls`: List currently monitored projects and their pipeline status
     - `add-projects <targets...>`: Register projects into watch list
     - `rm <targets...>`: Remove projects from watch list
     - `agy-running-projects`: Automatically register all currently running Antigravity projects
     - `run [-t <duration>] [--dry-run]`: Execute the watch loop until all monitored pipelines are green
     - `help`: Render rich two-column terminal UI help menu
2. **Help Menu Prominence Standards:**
   - Add `shutdown-until` / `sug` to `gitmap agy help` and `gitmap help`.
   - Implement `RenderAgySugHelp()` using `termhelp.HelpMenu` (cyan/yellow boxed two-column layout).
   - Author `cli/helptext/shutdown-until.md` with synopsis, workflow explanation, and examples.
3. **Safety & Dry-Run Invariants:**
   - `--dry-run` flag MUST be supported in `sug run` to allow testing the watch loop without executing actual OS shutdown.
   - When all projects are green under `--dry-run`, output `[DRY-RUN] All projects are green. OS shutdown command would be executed (shutdown /s /t 60)`.
