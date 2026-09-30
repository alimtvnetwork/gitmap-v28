# SPEC-APP-90: AGY Add/Add-Read, Running-Projects Tree (`rp prompts ls`), Last-Active-Projects (`lap`), Rerun-With-ID (`rwi` / `rwc`), Machine/Alias Management (`--ssh`), Two-Way Telegram Bot, Email Speed Settings, and Modernized OS/AGY Help

**Spec ID:** `SPEC-APP-90`
**Status:** Active (`v6.369.0`)
**Author:** MD ALIM UL KARIM

---

## 1. User Request (Verbatim)

```text
gitmap agy help

<prefix cli> = gitmap agy

<prefix cli> add
<prefix cli> add path
<prefix cli> add .
<prefix cli> add-read .
<prefix cli> add-read path
<prefix cli> rp ls
<prefix cli> rp prompts ls
<prefix cli> last-active-projects(lap) N [ls/help] [--limit/l Y] --offset/skip Z # list of projects which has communication in last 24 hours , N = 24, Y = 10,

<prefix cli> last-active-projects(lap) N [ls/help] [--limit/l Y] --offset/skip [--json] [--file/f "filepath where the josn will be saved"] Z # list of projects which has communication in last 24 hours , N = 24, Y = 10,

<prefix cli> last-active-projects(lap) N [ls/help] [--page/p 2] --wordcount(wc) T # list of projects T = 200,
<prefix cli> rerun-with-id (rwi) <project id/alais/sequence/path> <convid> "text prompt or file path to file or relative path to a flle" [-p/prompt nameof-the-prompt]
<prefix cli> rerun-with-convid (rwc) <convid> "text prompt or file path to file or relative path to a flle" [-p/prompt nameof-the-prompt]
<prefix cli>

<cli> machine ls/change/set/help
<cli> alias ls/change/set/help


<cli> machine ls/change/set/help --ssh
<cli> alias ls/change/set/help --ssh



gitmap os help

I think there are several issues on the GitHub AGY. So AGY, when we do the help, it shows add project user fine, but you don't have the add project with a single path. You don't have the add project with dot, dot means the current repo would automatically add in. And also you need to have add read. That means that is going to read the current path, add the project, and run the read prompt. So these are the things I want, and you don't have it, and you do not, let's say, fix the OS help command. That is really terrible and poor. And also we should be able to see that the running projects, we can see running project LS should share the running projects. Running projects means active prompts that is running. We could say running prompts, running projects, prompts LS. So that would show the prompts, let's say 200 words, and list out the projects' names and sub item as a prompt, like what is running, so that I can get a common idea. So these things are still missing. I really do not understand. Also, we could see running prompts, or we could say we could see last active projects. What do I mean by that? So usually last active projects is LAP. That is going to show us the list of projects. So we could give a number actually here, N, and LS or no LS. Could do LS, no LS, or help. Okay. So usually if I just put list active projects, so that is going to be not the running projects, but this is going to be list projects which has communication. In the past, let's say, 24 hours. So by default, N could be 24, and this can be changed from the settings default mode, and also user can put their number N. Let's say they can put 12, they can put eight. Okay? And there could be a limit as well. The limit is also optional, or L. So by default, the limit is 10. Default, if we do not put anything, then Y limit is 10, but it can be, yeah, changed as well. So if we do not pass anything, it would mostly limit to 10 projects. So it would mention that there are 10 project and more than that, we are not listing it out. Okay? We are max limiting the 10 projects. So this is important. Also, we could do an offset. Offset here, that means offset. Or it can be skip. Skip is same as offset, and that can be Z. And we can provide the offset number where we want to offset, or also we could see the same thing as paged. So in the page section, it's just going to use the page and limit under item. So you could say page two, so you could see page two or P2, or you could see page two list of active projects. Okay? And each one of the projects in the list of active projects by default would have the prompts title as a sub item. As a tree view, it would show which conversation is running, and inside this conversation, what is the prompt, let's say, last one or two prompts that has run the 200 characters, 200 words of that, but that can also be changed with the WW word count. Okay? Usually it's by default, let's say 200. It is by default 200, but can be changed to something that we want to. Okay. Make sure that this is how the last active projects is created. Okay. Also in this case, we wanted to see the project ID, conversation ID with brackets. Okay. So that I can recall where I wanted to target. Let's say I wanted to target the rerun with ID, and we can pass a project ID, we can pass, or alias actually, for the sequence. So here when the list active projects or prompts, running projects, these would always give the three things, project ID, project alias, project path, sequence number. Sequence number would be saved into the database that will be used for that iteration. Okay? So that would be always generated. That's a list of running projects. So that would have a sequence. That sequence should be saved as a cache for a day to the SQLite database. So if it is past that day, then the function needs to be called to regenerate the sequence again. Try to have this for the prompt. So mention all these things, and then we can just give a prompt ID. Okay. And then we can just do the text prompt or prompt. One more thing. So when we generate the tree, we should have a sequence ID, one ID that you can pick up from the database and would know what is the conversation ID. What do I mean by that? So we will have a rerun with prompt ID, which I could use the ID that you have listed in the list active projects or list RP. So that would have the conversation exact, some sequence of ID, very short form, that you can keep into the database that could be represented by a number or something that you can understand. And using that, you will find that conversation, find that project, everything, and then I could just give any text that I want, and that would actually inject it as a prompt. Okay? So yeah. You can do that, and I hope that you can help me with this. Git map OS help, that needs to be improved, okay? Yeah, so main idea is so that I could know which projects are running. I could switch to fast forward as well. I don't see that this command is added to the help. All this command I'm discussing that also needs to be added to the help. That needs to be, let's say, yeah, added to the help, added to the UI help, added to the commands docs. Do you understand? Can you please help me with this? Are we aligned on this? Also, please check that can we integrate or have a Telegram chatbot connected for the Git map? Can we have this feature? And please add the prompts for the Telegram section where we can add a Telegram bot to the machine where we could use the Git map Telegram And could communicate from the chat, from the messaging, and can communicate back to our system. Do you understand? Can you please do that for me? You can also check the Antigravity Manager for this. We have already done it. Also, for this, we should also add the email setup. So in the settings section, we should have these speed features, actually. So do you understand? Is it clear? Okay, so we can also check the running projects and active projects to a JSON as well. I mean, view it in JSON, also export the JSON as a file as well. So you can do a JSON flag. If you do a JSON flag, it will be outputted as JSON. We could do a F file. Again, these are totally optional. If you do that, then this is the behavior. It will be the file or the F. We can provide a file path. The JSON will be saved. I think you understood. So try to have this feature. Try to have the commands properly integrated. Okay, please help me with this. And also the OS help that needs to be modernized. Okay. And, oh, okay, one or a few things I just remember. So in the CLI, what I want, that is also a must, I think we need to have. So CLI needs to have a way to change the machine name, like the OS machine name. View it, change it, or things like that. That command also be residing as inside OS and also should be available outside as well. That means we could do CLI machine, and then we could do LS change or set or help or do nothing. Do nothing is just help, same thing. Or LS, it will display the machine IP, machine aliasing, machine name, whatever the OS name is there, regardless of the OS, Windows, macOS, Ubuntu, all these cases. Okay? And we can set it from here. So you need to test it, end-to-end test indirectly in this machine by setting some machine name and things like that. Aliasing also we should be able to do. Alias, you could change the current machine aliasing. Alias would be a name for the Git map that how other machines can recognize it. So they could change it by themselves. Also, other machines will know their aliasing. Okay. But by default, if they pick with the IP without the alias, then it would automatically pick this alias. This is where the alias will be helpful. So in machine or alias, in both cases, it would actually show the IP, it would show the alias name. Basically, the identifier for this machine, whatever that is possible. Okay? And we can set it. So when we say set, we just give a name with the format that is applicable. Usually, we give a sample in there in the prompt. Also user can set directly and say Y, and it would be applied automatically. Remember, it should apply Windows, macOS, Linux, especially Ubuntu for now. Do you understand? And also we should be able to do these things for SSH. That means all the machines that is out there, we are going to get its information, whether that is ID, aliasing, and things like that. We can also do set, we can do reverting back. So all kinds of things should be possible with it. So that is a machine identifier, and that should also be available inside the OS command as well. Do you understand? Okay. So apply this, make sure that you make everything correct. So at the end, you make a release bump. Okay? And at the end, you can check git map space PE to check everything is working. Is it clear?
```

---

## 2. Architectural Specifications

### 2.1 `gitmap agy add` & `gitmap agy add-read`
- `gitmap agy add` / `gitmap agy add .`: Resolves the current working directory (`os.Getwd()`) or enclosing Git repository root and adds/registers the project into Antigravity's workspace registry.
- `gitmap agy add <path>`: Resolves `<path>` (relative or absolute) and adds/registers the project in Antigravity.
- `gitmap agy add-read .` / `gitmap agy add-read <path>` (aliases: `ar`, `add-and-read`):
  1. Resolves `.` or `<path>`.
  2. Adds/registers the project in Antigravity.
  3. Immediately loads the canonical **Read Memory** prompt (`01-prompts/read.md`) and injects/starts a Read Memory conversation on the registered project.

### 2.2 24-Hour SQLite Sequence Cache & Tree View (`rp ls`, `rp prompts ls`, `last-active-projects` / `lap`)
- **24-Hour SQLite Sequence Cache (`gitmap-running-prompts.db`):**
  - Table `AgySequenceCache`:
    - `SeqId TEXT PRIMARY KEY` (e.g. `1`, `2`, `P1`, `P2`)
    - `SeqNum INTEGER NOT NULL`
    - `EntryType TEXT NOT NULL` (`project` or `prompt`)
    - `ProjectId TEXT NOT NULL`
    - `ProjectAlias TEXT NOT NULL`
    - `ProjectPath TEXT NOT NULL`
    - `ConversationId TEXT NOT NULL DEFAULT ''`
    - `PromptSnippet TEXT NOT NULL DEFAULT ''`
    - `CreatedAt INTEGER NOT NULL` (Unix timestamp)
    - `ExpiresAt INTEGER NOT NULL` (`CreatedAt + 86400` = 24 hours)
  - If existing cache rows have `ExpiresAt < time.Now().Unix()`, the cache is invalidated and regenerated deterministically (`1..M` for projects, `P1..PK` for conversation prompts).
- **`gitmap agy rp ls` & `gitmap agy rp prompts ls`:**
  - `gitmap agy rp ls`: Displays running projects with columns/fields: `Seq (#1)`, `Project ID`, `Project Alias`, `Project Path`, `Conversation ID [bracketed]`, and `Active Status`.
  - `gitmap agy rp prompts ls`: Displays running projects as a tree view where each project lists its active conversation(s) (`[ProjectID | ConvID | Seq: P1]`) and a sub-item showing the running prompt truncated to `200` words by default (`--wordcount` / `--wc T`).
  - Supports `--json` and `--file` / `-f <filepath>`.
- **`gitmap agy last-active-projects` (`lap`)**:
  - Syntax: `gitmap agy last-active-projects [N] [ls|help] [--limit|-l Y] [--offset|--skip Z] [--page|-p P] [--wordcount|--wc T] [--json] [--file|-f <filepath>]` (also top-level `gitmap lap` / `gitmap last-active-projects`).
  - `N`: Lookback window in hours (`default N = 24`, configurable in settings or passed positionally e.g. `gitmap agy lap 12`).
  - `Y`: Max projects to list (`--limit` / `-l`, `default Y = 10`). When total matching projects `M > Y`, prints explicit notice: `Showing 10 of M active projects (max limit 10 — pass --limit or --page to view more)`.
  - `Z`: Offset/skip (`--offset` or `--skip`, `default Z = 0`). If `--page` / `-p P` (`P >= 1`) is provided, `offset = (P - 1) * Y`.
  - `T`: Word count limit per prompt sub-item (`--wordcount` / `--wc`, `default T = 200`).
  - Tree View: Each project displays `[Seq: #1] [ProjectID: <id>] [Alias: <alias>] [Path: <path>]`, followed by tree branches (`├─ Conv [<convId>] [Seq: P1]`) and the last 1–2 prompts truncated to `T` words.

### 2.3 `gitmap agy rerun-with-id` (`rwi`), `gitmap agy rerun-with-convid` (`rwc`), & `rerun-with-prompt-id` (`rwp`)
- `gitmap agy rerun-with-id` (alias `rwi`, also top-level `gitmap rwi`):
  - Syntax: `gitmap agy rerun-with-id <project-id|alias|sequence|path> <convid> "<text prompt or file path>" [-p|--prompt <name-of-prompt>]`
  - Resolves `<project-id|alias|sequence|path>` (using the 24h SQLite sequence cache if `<sequence>` is numeric like `1` or `#1`) and `<convid>` (or short seq `P1`).
  - Resolves the prompt from literal text, a relative/absolute file path if it exists on disk, or a named prompt template via `-p` / `--prompt`.
  - Enqueues and dispatches the prompt into the resolved project & conversation.
- `gitmap agy rerun-with-convid` (alias `rwc`, also `rerun-with-prompt-id` / `rwp`, and top-level `gitmap rwc`):
  - Syntax: `gitmap agy rerun-with-convid <convid|short-seq-id> "<text prompt or file path>" [-p|--prompt <name-of-prompt>]`
  - Looks up `<convid|short-seq-id>` (e.g. `P1`, `1`, or full/prefix conversation UUID) in the 24h SQLite `AgySequenceCache` and Antigravity conversation dirs to automatically identify both the target project and conversation, then injects/reruns the prompt.

### 2.4 `gitmap machine` & `gitmap alias` (`ls/change/set/revert/help [--ssh]`) + `gitmap os` Integration
- Top-level commands: `gitmap machine [ls|change|set|revert|help] [--ssh] [-y]` and `gitmap alias [ls|change|set|revert|help] [--ssh] [-y]`.
- Also available under `gitmap os`: `gitmap os machine ...` and `gitmap os alias ...`.
- Running with no subcommand or `help` displays the boxed help menu and current machine identity summary.
- `ls`: Shows Machine IP (`Local IPv4`), Machine Alias (defaults to IP if not explicitly set, so other machines always have an identifier), OS Hostname / Machine Name, Previous Name/Alias (for `revert`), and OS Platform (Windows, macOS, Ubuntu/Linux).
- `set <name>` / `change <name> [-y]`:
  - Shows a format sample (`e.g. dev-win-01, build-ubuntu-node2`) and applies the change (automatically when `-y` or non-interactive is passed).
  - Cross-platform OS hostname update + GitMap SQLite identity persistence (`machine.alias`, `machine.name`, `machine.previous_name`, `machine.previous_alias`).
- `revert`: Restores the previous machine name / alias.
- `--ssh`: Executes `ls`, `set`, `change`, or `revert` across all joined SSH fleet machines (`db.ListSSHConnections()`), aggregating their IP, alias, OS hostname, and platform in a unified table or JSON.

### 2.5 Two-Way Telegram Chatbot Integration, Prompts & Email Speed Settings
- `gitmap telegram [setup|start|status|send|poll|help]` (and `gitmap agy telegram ...`):
  - Connects a Telegram bot (`botToken`, `chatId`, auto-discovered from Antigravity Manager config or configured via `gitmap telegram setup --token <TOKEN> --chat <ID>`).
  - Supports two-way command polling/webhook handling (`status`, `rp ls`, `lap`, `rwi <seq> <prompt>`, `asw`, `help`) so users can chat with the machine's GitMap instance from Telegram and receive responses back.
  - Canonical prompt authored at `01-prompts/telegram-bot-setup.md`.
- `gitmap email [setup|status|test|help]` & `gitmap settings` (`gitmap agy settings`):
  - Speed settings for Email (`smtpHost`, `smtpPort`, `senderEmail`, `recipientEmail`, `appPassword`) and Telegram (`botToken`, `chatId`), plus default `lap` hours (`lapDefaultHours = 24`) and account-switch threshold (`15%`).

### 2.6 Modernized `gitmap os help` and `gitmap agy help`
- Boxed cyan/yellow enterprise help menus for `gitmap os help` (`cli/cmdos/`) and `gitmap agy help` (`cli/cmdagy/agy_help.go` & `cli/helptext/agy.md`), documenting every new command and flag with real-world examples.
