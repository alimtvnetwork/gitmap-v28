

gitmap pull all (pa)
gitmap pull all ssh (pas)


gitmap fix ignore all  [-y]# will fix on all repos to fix gitignroe issues
gitmap fix-ignore-all (fia)  [-y] # will fix on all repos to fix gitignroe issues
gitmap fix ignores all ssh [-y] # will fix on all repos to fix gitignroe issues
gitmap fix-ignores-all-ssh (fias)  [-y] # will fix on all repos to fix gitignroe issues for all nodes, current node will run as is, follow the gitmap pas formula for this and write in the spec as a term for Gitmap PAS formula so that can be reffered by to any AI properly
gitmap commit-push-all-repos (cpar) [-y]# commit all repos if there is any pending
gitmap commit-push-all-repos (cpar) --review(r) # show the pending commits for review first if approved then commits and push
gitmap commit-push-all-repos (cpar) --review --commit-only(co) # show the pending commits for review first if approved then commits and push
gitmap see commit pending
gitmap see git-ignore/ignore/ig issues
gitmap see errors # same as gitmap errors
gitmap see history # same gitmap history

gitmap see errors ssh (ses) # same as gitmap errors, ssh follow gitmap PAS
gitmap pull-all-ssh(pas) # do pull all right now ly/optimize
gitmap ignore/ig add/scan/scan-ssh(ss)/remove/edit/action/ls/help/ui/app/add-group/remove-group/rm-grp/set-default-group/add-grp-to-default(agtd) <name of the group>/apply ./connect-group-with-repo (cgwp) /export/import
gitmap ignore connect-group-with-repo (cgwp) <group-name> <path1>,<alias of the repo> --add-with-default(awd)
gitmap repo-manage ui

gitmap cache create .
gitmap cache  create <relpath1>,<abs path2>
gitmap cache  create "a.json", "b.json"

gitmap cache ls/add/create/remove/rm/help
gitmap cache search "text search" "*.md" [--lines 10] [--limit 20]
gitmap cache search "text search" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache search "text search" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache search-multi "text search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache search-multi-grep "regex search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache recache/reconcile/sync
gitmap history ssh
gitmap nodes histories/history




Important instructions must follow


Okay, I think we need to improve the Git Map a lot, actually. So, first of all, if you look into this, these are the common problems that you have inside Git Map, that whenever we are doing the pull first, it tries to check that the ignore file is there. Right. So gitignore, it's mentioning the file. That is quite bad because the first thing it should do is pull. That's its job. That's what we requested for. Okay? So once it completes the pull, it shows every report. And also at the same time, when it is pulling or something like this, it could actually send parallel, let's say, async request to check if the ignore files. Why I'm saying ignore files, it's not one, but it can be many in the future that we can add it to the system or a text file inside the Git Map, there should be settings that we can modify the gitignores that would be automatically ignored. Also, there could be apply. We could apply certain type of groups of ignore. So we can create. So that's why I have created a section called gitignore add. Add, which by default, is going to be added to the default group. Okay? So we can have add group command inside Git Map ignore. So group is we can actually use a group. So when we create a group, we can actually apply that group on certain folder, certain repo. Okay? We can tie a group with some settings that connect group with repo, where we could actually provide, which is like CGWR, connect group with repo. So I give a short form all the time. If the short form already exists, then do not reuse the same one, just suggest a different one or add a different one, which should be absolutely fine. So if we use Git Map ignore with connect with a different repo, we can actually provide a path. Okay, path one, we can provide the alias of the repo. Okay? So all kinds of things we can actually connect. So if we connect additional git group, okay. If we connect any group name and then... Okay. So all these things, I need you to provide and put into the help text, terminal help text, also in the UI help text, also in the docs very properly. Okay? So what the ignore would do is by default it would have a default group. So if you want to add something, that would basically going to be added to the gitignore, and that could be a text file that we could also modify and check what are the default ignores that we want to add there. Now, the thing is that default ignore should also have the .gitmap folder. There is a .gitmap folder that we create inside the repository. Inside this, sometimes we have the backup folder. So that folder also be default gitignore. Remember to add it as a default one. Okay. Now, coming to the point, if we use the add, we can actually add any file name or a name or a folder. It does not have to be a specific folder, but it can be a relative to that repo. So we can just give a string and that would be added automatically into that text file and also the SQLite DB. Now, it will also check the duplicates. If it is already added, then it said, "Already added. Don't need to add." Okay? Now that's the first thing. Second is that we should see LS. LS will tell us how many groups are there and each group have how many default gitignores. Okay. And is this group added with the default group as well? That means default will run, then this group will also run automatically on every repo, or some group is connected to a specific repo only. And these all can be import/export using JSON as well. So that means we need to have export/import settings for the JSON file, so we can import that as well. 




And also try to understand the formatting of the JSON import and export. I think you already know how to deal with it. So you understand the formatting. Okay. And also which format should also tell us, like is it the format for the gitignore. Okay. Now coming to the point, we have inside the gitignore, we have duplicates sometime. Okay? And that is a serious offense So the way that Git Map pull all. So I think I have to start with the Git Map pull all or the EA short form will work is that it is going to, how can I say? Pull all is going to pull first, that's the first thing, the habitual behavior, it will start. And also at the same time, for the other repos, it will just for every two, three repos, it will just create one async request. So it would have a full fledge of the repo information. So based on that, how many repos are there, it would create a worker type group, something. So each worker, it would not spawn a process, just the async work. So you should have a util method that actually helps you do that, drying the code. Remember, drying the code all the time. So we might use everywhere, this idea. So what it should do is that, let's say we have 50 projects. So let's say five projects will be handled by each worker. So let's say we spawn five async processes, and each one will do, again, can complete five repository searching to find this, let's say, ignore file has a issue, and that would be summarized and asked at the end when the summary is already done. So it will start first, and then at the end it will join and wait for the system to complete the ignore summary, and it will tell us a few things. First, if a repository is gitignore needs to be optimized, that means it has duplicate files, okay, that needed to be optimized. That's first thing. Second is that a repository might have a gitignore folder or file, and that file is already there. Now, if that is one of the cases, then it would highly say that, at the prompt at the end, all the summaries it will do at once in a very nice, colorful way, if it finds anything. If not, then just ends it. Now, if it finds these issues, then it is going to summarize, each repos have what type of issues. Summarize in a nicer way, cache and sub points. And then it will give two options. One, user can fix all in one shot. Like just prompting, yes, Y. Or user wants, probably first asks, their first question wants to resolve all at once, or user wants to do each one item as a single prompt. So it will be asked with single items at a time, so user can make a decision what to do with them. So when it optimize the gitignore file, remember to not touch any other aspects. Only it will going to ignore if the things are, I mean, optimized, if the things are duplicated. Remember that. Exact same duplicate. Okay. And if the file is already ignored, it does not contain in the git commit. You don't have to check all the commits. You just have to check that if the current file is added to the current commit state or not. If it is there, I mean, added to the git. You check in terms of git. If it is there, then yeah, it's a red flag, it needs to be removed because it's in the gitignore. Only then you will raise flag. But it's not like you're going to raise flag every time there is a ignore file there in the repo. No, you don't. You just close it. Everything is done. This is how it needs to be done. That's the Git Map pull all how it's going to work. But also at the same time, we want to have some automation and better approaches. For example, we want to have Git Map pull all SSH command. That would be PAS altogether. And when a machine runs PAS, by default in the current machine, it's going to run the pull all. And all those nodes, which is except for the current one, if we are running from this, so it would understand this. It would run the command with SSH nodes, but using Git Map, but also at the same time, it's going to run a little bit different. It's going to have a less concurrent request, that's the first thing. Second, the system, it needs to have a different type of code. Because the concurrent request would be less, it would be creating two workers only, that workers is going to probably handle two async operations only, and this is how it's going to complete. So each one of the nodes. Remember to do that. Now, if, think about this, if the processor is in very high, let's say, in pressure, there is just one worker, two async hence. This is how the code needs to be. It's very delicate running in the SSH. Now, at the end, when the SSH is done, it actually communicates using JSON. And we can do the similar process if we actually expose API endpoint, and which we will discuss later on. 



API endpoint I wanted to discuss because we want to Use that Gitmap as an MCP server for the AI in the future. Put a question mark and ambiguity question in your ambiguity folder so we can discuss this later. Put it as a pending task. We will discuss it later on, how it's going to work. Okay. Now, the pull all SSH, it needs to be very much efficient. Okay? So these, let's say number of nodes that we have, the first thing it will do is enqueue the task on those nodes. And remember, Gitmap, every action needs to go through the task servers. So it needs to be added to the task, and then it will start the task, and when completed, it will also mark in the task that it is completed so that we can see it in the history, we can see it in the completion failure, and every error that happens, it needs to be saved inside the Gitmap errors DB. Okay? So then we can see every error using errors. We can also do a limit of the errors, how many we wanted to see, things like that. Okay. These are all right. Now, so this is what you learn. SSH is a different thing. You need to be very efficient, different type of code. You can reuse some of the code which is common, try to make the code dry as much as possible so that the code is not repeated. But also some of the functionalities, like how it starts from the branching out, you don't write too many if else as you write smaller struct functions so that you can delegate the task. Remember, system design is the king. If you have a bad type of system design, the code will be too much. You will have more work to do. Okay? So the things that you do, you will also write in the spec mode, okay, so that you don't make these type of crazy mistakes. Okay, so whatever we have discussed about the Gitmap pull all SSH, which is PAS, Gitmap PAS, I want you to recognize this as a Gitmap PAS formula, which I can refer to you in future, anytime I need a similar situation. Okay? So now we are going to discuss about the Gitmap fix ignore all. Fix ignores all, or fix ignore all. It could be S or without S. Both should work. Okay? Remember that. It could be fix ignores or ignore, both cases, all. And also you have SSH version of it. The SSH version will, again, same way it will run as Gitmap PS. Okay? Now coming to the next part, Gitmap commit push all repositories. So that is basically going to try to commit all these repositories, okay, using the Gitmap feature command, okay, and try to showcase all as a summary like this is what it has done. But also if the user wants to have a review, they could just use a review flag. I can have a review or R flag that will actually summarize all of the repos that has changes, and as a sub-node, it will show as folder tree, like where the changes are. Again, if they wanted to show all first, and then it will propose two different options. Either one, they commit everything as a feature using the Gitmap, okay? That's one thing. Again, all the commits, all the push, everything needs to be in the Gitmap tasks, okay, so that we can see the history. Okay? That is very important. Okay. Now coming to the point. Okay. I think you understand the fix all repo, how it's going to do. It's going to just perform the scan first on all the repos and try to see if there is a ignore inconsistency or duplicates, or the ignore file is already there. Then it will suggest the solution all together as a summary. We could do, let's say, hyphen Y at the beginning to avoid the prompts. Otherwise, it can show us a prompt like we want to resolve it all together, or we want to resolve it single by single. So let's say we pick one or two option. So if you pick one, then it would say how we want to resolve it. Are all these bugs, we want to commit as a bug or feature, how we want to do it. So we pick it and then done all of these commits in one. Now if we don't do that, then it would showcase again, each one of them as a sub-node, sub-item, as a lot of options that what we could do. 




We could see the files, what files we want to see, what file we want to discard, what file we want to commit. Is it, let's say, bug fix or features? All kinds of options will be there, so how we want to do it. Each one of those would be treated as a session. We resolve each one of the repo and then proceed to the next one. Do you understand? And each one of the tracking and things we are doing, it needs to be in the tasks, enqueued, and then when done, task will be marked as done. So remember that you need to have a global system where you could actually deal with these tasks. So your code needs to be that much elegant, dry in terms of that. Okay. And similarly, the SSH, or if we provide the hyphen Y at the beginning, then everything will be confirmed, and it will proceed. And, okay. Now, the SSH would be same as the gitmap pas command, how it's going to work. So you follow the same technique. Now, commit push. Now, if we go to the commit push all repos, this is going to be the same as the ignore, somewhat like this. So it's going to go through all these repos that we have, and then it can ask, these are the changes we have. This is how you want to commit those as feature, bug, or things like that all together or single one, create the session. Same will go for the fix ignores as well. I think you understood the point. If any issues, we can discuss later. So we will have another command segment. Also, it means we have the help, C help, things like that. So C will help us go through the things. For example, we can say gitmap c commit pending. So that would show us all the repos that has changes tab which is not committed yet. Okay, C gitignore issues. Same one, we just can see the ignore issues. So same thing, just can be said as just ignore issues or ig. Okay? So gitignore can also be written as ig, short form. Okay? That would have a lot of sub commands. We will come to this. Git map C commit pending. It would showcase as similar as the commit push all repos, but in this case, it will just showcase. Yeah, it will also suggest like would you like to commit as a prompt, okay? Yeah, you can just delegate the same code, I think similar one. Git map C gitignore. It's just going to show which repos has issues and things like that. Something similar to the fix ignore. Okay? It will also prompt the user if you want to fix or not. Okay. And then repo issues. Repo issues means... Okay, repo issues, I forgot how I thought of it. Yeah, repo issues, I think we have to remove for now because I cannot think of it. C errors, git map space C in order to just do the git map errors, but we could also do git map errors and C errors in both of these. We could do it using SSH as well. That means we can have S-E-S, something like this. And it is again going to follow the same process at git map EAS. Okay? All yes, we have discussed. Now we have gitmap ignore. So gitmap ignore will have lots of sub commands, as I have discussed before. We can add items directly to the default group. We can do scan, scan SSH, remove items, remove groups, actually, also. So add group, remove group, or RM GRP also. So you know how to set the group to this default, or you pick a group to be added as a default ignore group. We can do that. Or also add groups to the, add group to the default group. Add GRP to defaults. GTE. Yeah, and we can just pass path. Sorry. Name of the group. That group will execute after the group, default group. Okay? Gitmap ignore apply means we scan, okay, and we can apply inside a folder or repo. If you're inside the folder, then it would apply inside that folder. Okay, with dot or without dot. Okay? Automatically, it will only apply to that folder. So if we do not have any folder, then it will try to just do the gitmap fix ignore all, something like this. But applying all these rules. Okay? So if we are on a specific rule, first the default ones will be applied, then the rule specific to that repo connected will be applied. Okay, I hope you understand, but if you have any issues, we can discuss it. There is another command like repo manage that would actually showcase all the repos we have in the git. I mean, hub, git lab in the future, we will discuss this. And then what we want to clone, how many repos we already have, I mean, already in the system, where those are. We want to run something on there. So all kinds of things needs to be in terms of the UI level, okay? Then new command, we have gitmap create repo cache. So we will create a different type of split DB close to the root DB we have close to the CLI, right? So there would be cache repo folder. 



Inside this, we will have the repo that we have targeted to do this creating cache. It will create a slug for that repo, and that on that slug, it would create a SQL.db file where it would have-- So SQL.db will have the folder structure and the root level files. Okay? All the files information that it, the repo has and the root level files. And then every folder that it contains, for this, it will create a separate SQLite DB with the folder name as a slug name. Okay? And the folder exact relative path from the repo, it would be mentioned inside that SQL DB that we just created for that cache DB. Okay? Now, this is a split DB concept, so you need to look into the split DB concept to understand this and follow through the split DB concept according to others and the spec. Okay? If you have any questions, we can discuss this again. Now, the cache, when it's started to create, remember we always ignore some of the parts. For example, the .git folder, VS Code folder, or let's say, the IDE folders we ignore. Node modules we ignore. Okay? Some of the commons ones would be automatically ignored. Okay? And that what we should have in the gitignore, GitHub ignore as well. The GitHub ignore would be powerful one that we also use here in the create cache first. Okay? Also in the create cache SQL DB first, the root DB for this repo, we should have the repo URL, most of the information so that we can relate to and query to where it is happening. And also this task also needs to be created in the tasks, and we should be able to see as a history. We can do undo, redo as well. Remember that. So when the cache is started, it would actually spawn agent based on the folder we have. Okay? Now, if a folder has many nested folder, in future, if the workers are resolved, then some of the workers will take the files from this as a split responsibility and start working on it. So in the SQL.db, the main DB that we create, that should contain all the files path, absolute path, relative path, and also the path where that repo actually exist. Okay? Relative path and so on. So usually absolute path is not there, not needed. I think relative path and one settings or one place that would have the where the repo starts, so that it can combine all the time and find it. And it does not have to find it automatically. There would be a table view that would combine this and create a View that actually contains the full absolute path automatically so that we reduce the variable. Remember that. Now, this can be changed also. View can be there where we could think of the path just changed, and it would act like it, things like this. Okay. Now, in the sql.db, which is the root DB for this split DB, it would have the all files information, all folder information. And based on all files, it should not have the hash, but it should have the last modified time. Okay? And that is very crucial because anytime we want to reread or understand that there is a change, we could actually check if the file is outdated. If the file is only outdated, then we take the file again, and the file will be kept only if the root files will be kept inside the sql.db file. So we need to have a file type database, I mean, table, and there should be one too many joins, so that the table is normalized, or database is normalized. So these files which are kept in there, only the files at the root level should be kept in there, and every folder that contains that would have its own DB. And not the subfolders, so it would have subfolder, subfolder. Not the subfolders will be created as a DB, remember that. So similar technique that we have to apply, okay? And you have to confirm that you understood this. Now, when the create cache starts, it actually needs to have the worker process so that it have more hands, more async operation it could do in a store. So usually there is a default threshold. So any file that is more than 200 to 300 KB, 200 KB, it would not store to the database. Okay? And also no binaries, no image, things like that. It can contain the image path, but no image file. No image file, no big type of JSON, it should not have, okay? It can have the JSON file path, but also at the same time, it would have a full overview saying that we don't keep the large JSON files, okay?




 The files which are large, it would also be in the table like we do not keep. Okay? Is keep. There should be a is keep flag, which actually tells us we have it or not in the cache. And also we can query any file if this is already in the cache, and we should reuse it or things like that. Now, the purpose of this cache is to search effective filing or effective search during the search operation, and also the list of files, and also the AUM search. So check the algorithm, how it is done right now, and can we improve some of this using this cache? And also, since it's a cache, it needs to be updated time to time. So we need to think about that as well. So we are still not going fully depending on this cache, but it's an idea. So from this cache, we should be able to see the file histories and things like that. Cache LS should tell us how many cache repos that we have. I think, yeah, the best way is to just git map cache create. Cache create. I think that it would be the best idea because that would have a command segment, things like that. LS add. Add or create would behave like the same. Remove and RM will behave like same. We should have a help command that actually explains how does this work and how efficient that is. So we should be able to search using the cache. We should be able to search. So we can do search, direct search in text or file name. So we should have, wait. Yeah, search would be text search. We can also provide formats, star MD files. We should only search on the MD files. We could do on in between contains, something like this. We can provide multiple MD files that we want to search. Okay. And this all needs to be in the help. Without the help, no one can use it. Also, you need to add it to the docs, okay? You cannot miss anything. So search would do it like this. Search will also do recaching if the file is outdated. If the file is outdated, we will not trust that file. And during the search in a repo, if we find out there is a new file, and that search will automatically go ahead with, let's say another async operation. That would just create another process for the git map. Automatically, it will just create the process. And that process will only check, and we don't need to get the answer here because that can automatically write to the cache DB. If we start working with a repo, then it would automatically start the cache reconciling. So that means it will check the file last modified time. If there is a mismatch, then it will try to update from the file system to the DB automatically, because the SQL.db knows where the file is, in which DB and where. So based on that, it can quickly modify that file, update the file. And we can do multi search as well. Multi search, this means we provide multiple search. And we can do star, we can do star in the middle, stars with we don't care. We can do all kinds of formatting. So you could provide that samples. We can also have search multi or search multi-grep. So we can actually do a regex search. That regex search will be performed in terms of the SQLite, not in our regex. Okay, that would be very much faster. Again, we would reconcile the file data, and we could also do a limit, like how many data we wanted to see, how many lines. So usually when we find the text, it should only display the amount from above and below. Usually, the lines would be 10, actually, 10 lines, but it can be modified with the user's preference. Okay. And also we can limit the search. Usually, the limit search would be, if we are searching for the text, then the limit would be, yes, 20 by default. So it will just showcase the file names, that where it found, where it is, also the line number, and a little bit of the context where that is the 10 lines. This is how it needs to be displayed. But also at the same time, all these flags can be applied to every one of them. Okay? It's very, very powerful. Also, at the same time, we should be able to, cache reconcile or recache. Recache and reconcile is the same thing. Is going to understand the cache or the sync. I mean, all of these will behave the same way. Okay, I think I explained all of these factors. Now, it comes to the factor called the Gitmap. SSH, Gitmap history SSH. Again, it's going to follow the similar technique as the Gitmap PAS, but also at the same time, the Gitmap nodes history, the same thing as the Gitmap history SSH. Same way Gitmap PAS runs, it will follow the same technique. Okay, so there's a lot of things I mentioned. So first thing is that you list out the task, and then you write this exact verbatim to the spec first, and you create multiple tasks. Move these multiple tasks to the planning mode or planning file, how we do it, and then you instruct the agent to do it. Is it clear? Do you understand all these things? Okay. At the end, I want you to release, of course. Okay? So you release at the end, and then you check the Gitmap PE in the repository until it becomes green. At the end, not now. Okay? When everything is done. In between, you don't check anything. So before you make that release, you just check one time the build. Usually, the build is, that's a ban. But in this case, we will do it, but only one case. And also, if we find an issue in the Gitmap pulling, then we do not do this Git solution. We try to provide a solution from the Gitmap to fix it. Try to understand this. And there should be one command that would fix all this. I want this to be fixed as well. I have added a screenshot that is actually bugging me. I think you should be respectful and try to fix this issue as well. Does this make sense? Okay, let's finalize then.






# [V6] Parent Task N-Step Continuous Loop & Mandatory Multi-Agent Subagent Orchestration — Workflow (must follow)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /execute-parent-task-with-n-steps-v6 <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand:goal) Autonomously orchestrate and execute the parent task end-to-end: FIRST showcase and list out the given task in visible chat during Turn 1, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit that holds strictly this task's files.

[/learn](slashCommand:learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ABSOLUTE, NON-NEGOTIABLE MUST:** Spawning subagents via the `invoke_subagent` tool (`A = 2`, `H = 2`) is an **ABSOLUTE, NON-NEGOTIABLE MUST** in both **Phase 1** (parallel codebase discovery reading and modular spec authoring) and **Phase 2** (parallel subtask code execution with `TypeName: "self"`).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing all discovery reads or all subtask code modifications by itself without invoking `invoke_subagent`. Failing to call `invoke_subagent` when `A >= 2` is a critical protocol violation on the same tier as Rule 0.
- **Phase 1 Mandatory Subagent Dispatch:** Immediately after establishing the Confirmed Task Breakdown (Phase 1A) and the single-agent unified blueprint overview (`01-overview.md` or parent plan skeleton), the lead agent MUST call `invoke_subagent` to spawn `A = 2` subagents in parallel for codebase discovery/reading or modular spec sections and yield the turn to await `<SYSTEM_MESSAGE>`.
- **Phase 2 Mandatory Subagent Dispatch (`TypeName: "self"`):** Once subtasks are generated in `.ai-memory/plans/subtasks/xx-<slug>/`, the lead agent MUST call `invoke_subagent` with `TypeName: "self"` to dispatch `A = 2` worker subagents (`H = 2` disjoint subtasks per worker) and yield the turn to await `<SYSTEM_MESSAGE>`.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable.
2. **Showcase Given Task First (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, you MUST output the confirmed task breakdown directly in visible chat. Never execute tools silently without displaying the task breakdown to the user first!
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Screenshots & Media:** Decode base64/screenshots immediately into `assets/screenshots/<slug>-<NN>.png`. Reference via relative markdown links (`![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
5. **Discrete Deliverables Extraction:** Break down whatever user requirements were given into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).
6. **Mandatory Same-Turn Tool Chaining (TOTAL BAN ON TURNING OFF):** Emit the breakdown in chat with clean vertical formatting, and in the **EXACT SAME TURN**, invoke your first tool call (e.g. `write_to_file` to initialize ledger/spec, or run preflight). NEVER emit text alone (which ends the turn prematurely), and never ask "Should I proceed?".

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to Preflight & Phase 1B Spec Generation (Active Tool Call Running Below).
```

---

## 1. Precedence Hierarchy & Scope (Highest First)

1. **User Instructions & Preamble:** Directives and parameters ABOVE this prompt outrank everything below.
2. **Platform Limits:** Native tools, Artifact Review Policy, permission prompts, hooks. Never claim to override them.
3. **Repo Rules:** `AGENTS.md`, `.ai-memory/strictly-avoid.md`, and `coding-guidelines.md`.
4. **This Prompt.**

If sources conflict, follow stricter one and record under `Conflicts:` in ledger.

### Scope Control Rules
- **Turn 1 Task Showcase:** In your very first response turn, you MUST showcase and list out the given task in visible chat. Running tools silently without presenting the task breakdown is strictly banned.
- **Read budget:** Read only requested paths, search hits, and Step 0 context.
- **History read-only:** Never edit past events, changelogs, completed plans, release notes, or `06-old-prompts/` and `19-old-execute-prompts/`.
- **Out-of-scope:** Log under `Follow-ups:` in plan; never fix in this run.
- **Minimal diff:** Change only lines required; never reflow unaffected lines.
- **Indexes:** Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts. Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md` every run. Register recent completed tasks before push.
- **Mandatory Multi-Agent Partitioning:** Even for tasks touching few files, work MUST be partitioned across A workers (e.g. Worker 01 implements changes, Worker 02 implements verification/linters/companion tests). Solo execution is strictly banned.
- **Finish early:** When all Task-IDs are `DONE`, proceed directly to consolidation.
- **Zero releases:** Never bump versions or edit changelogs unless requested (R10).

---

## 2. Core Operational Rules (Cite by ID)

- **R1 Zero Builds or Test Suites (TOTAL BAN).** NEVER run `go build`, `npm run build`, `vite build`, `go test ./...`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. CI verifies builds and suites. Routine turns must never waste time on heavy compilation/tests. Only explicit user command lifts this.
- **R2 Targeted Checks Only.** Run only fast, file-scoped checks on specifically modified files (see Section 10). A check scanning 0 files is a **FAIL**.
- **R3 Evidence or It Did Not Happen.** Every `DONE`, `PASS`, or "verified" claim MUST cite a concrete file path, git diffstat, or command exit code (`exit 0`). Vague assurances are auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths.** Verify commands with a harmless call (`gitmap lf readme.md`), not `--help`. Use documented fallbacks and log in ledger.
- **R5 Mandatory Subagents (`invoke_subagent`).** Spawning subagents via `invoke_subagent` (`A = 2`, `H = 2`) is an **ABSOLUTE MUST** (`research` for discovery in Phase 1, `self` for edits in Phase 2). The lead agent is STRICTLY FORBIDDEN from executing all reads or edits solo. Solo execution without calling `invoke_subagent` is an auto-reject failure on the same tier as Rule 0.
- **R6 One Owner Per File (Disjoint Bounding Boxes).** Within every worker wave, each file has exactly one owner. Shared indexes (`.ai-memory/plans/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/what-to-read.md`, `02-spec/21-app/readme.md`, directory `readme.md`) belong exclusively to lead.
- **R7 Git Safety & Isolation.** Subagents never run git commands or alter git state. Nobody runs `git reset --hard`, `git checkout --`, `git clean`, `git stash`, or force pushes.
- **R8/R9 Atomic Commit & Push via GitMap.** The run ends with one GitMap call: `gitmap cpf "<summary>"` (features) or `gitmap cpb "<summary>"` (fixes). GitMap stages, commits, and pushes. Never commit file-by-file. Before GitMap, all push gates must pass (targeted checks, secrets gate, and `.gitignore` hygiene; untrack any ignored files: `git rm --cached`). Push rejected: `git pull --rebase`, re-run command. Miss after push: allow one follow-up `gitmap cpb "<summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits. Workers never run git commands or GitMap commit tools; only lead does.
- **R10 Zero Unauthorized Releases.** Never bump versions, edit `version.json`, update changelogs, or trigger release scripts unless user explicitly requested release.
- **R11 Strict Relative Git Paths & Lowercase Hygiene.** Strict ban on absolute paths (`C:\...`, `/home/...`) and `file:///` URIs. Paths relative from git root. New filenames and specs strictly lowercase.
- **R12 No Polling / Immediate Turn Yielding.** Print progress line (`Dispatched Worker 01 .. Worker <A> (wave k / WAVES); waiting for their results.`) and **STOP CALLING TOOLS**. Never poll in loop. Check `manage_subagents` once if wave runs long.
- **R13 Two-Strike Retry Cap & Anti-Looping.** Tool failing twice: worker replies `STATUS: BLOCKED` with exact error and stops. Lead takes over and logs `LEAD_FALLBACK: <reason>`. Subtask failing two remediation rounds is marked `FAILED` with RCA (Section 12).
- **R14 100% Ambiguity & Decision Boundaries.** Non-blocking: choose conservative option, log in ledger `Assumptions:`, proceed. Blocking: `ask_question` once, log in `.ai-memory/ambiguous-questions/01-new-ambiguity/`, continue unblocked tasks.
- **R15 Zero Generated Artifacts Committed.** Never commit build caches, logs, temp scripts, or newly generated code (Hard Rule 1) unless repository already tracked them.
- **R16 Zero Secrets in Standard Repos.** Never write credentials, tokens, passwords, or `.env` contents into tracked files, commits, ledger, plans, specs, or prompts (`AGENTS.md` section 9).
  - *Secrets Gate (lead, before GitMap call):*
    1. Check changed/new files via `git status --porcelain`.
    2. Run `python linter-scripts/check-forbidden-strings.py`.
    3. Search files via `git grep -nE` for private keys (`BEGIN [A-Z ]*PRIVATE KEY`), AWS (`AKIA[0-9A-Z]{16}`), GitHub (`gh[pousr]_[A-Za-z0-9]{36}`), OpenAI (`sk-[A-Za-z0-9]{20,}`), Slack (`xox[baprs]-[A-Za-z0-9-]{10,}`), or secret/token assignments.
    4. On hit: if `repo-secrets` exists in default work directory, store via `gitmap rs text "<value>" --slug <slug>` (or `gitmap rs file <path>`) and replace with env var/placeholder; if not, remove value and `ask_question` once. Log `SECRET_OFFLOADED: <file>:<line>` without value.
    5. Never print secrets in chat/logs; refer to file:line only. Workers finding a secret report `BLOCKED: secret at <file>:<line>`. Never put repository URLs or absolute paths into secrets instructions (`AGENTS.md` section 9).

---

## 3. GitMap High-Speed Command Primacy (Run Everything Faster)

GitMap is your **PRIMARY** acceleration engine:

| Operation | Command | Alias | Purpose |
| :

---

## 5. Step 0: Preflight, Platform Handshake & Ledger Creation (Phase 1 Budget)

1. **Platform Handshake:** Confirm tools (`invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, `write_to_file`, `replace_file_content`, `run_command`). If `task_boundary` exists: set `PLANNING` (Phase 1), `EXECUTION` (Phase 2), `VERIFICATION` (Phase 3).
2. **Commands & Directory:** Confirm `gitmap --version` and `python --version` exit 0. Verify GitMap with harmless call (`gitmap lf readme.md`), not `--help`. `run_command` uses `Cwd` in workspace root, paths relative. Never cd to other drives or tool folders.
3. **Working Tree Cleanliness:** Run `git status --porcelain`. Record modified files in ledger; never touch them. Confirm root `readme.md` is lowercase. Read `.ai-memory/what-to-read.md`, `strictly-avoid.md`, `coding-guidelines.md`.
4. **Resume Procedure (Check Before Creating):**
   - Match `.ai-memory/temp-agents/*/ledger.md` on `Request slug:` and `Request first line:`. No match: start fresh.
   - Status `COMPLETE`: verify commit in `git log`, report "Already complete: <sha>", stop.
   - Status `ACTIVE`: log `RESUMED_FROM: step x, phase p, wave k`; if `Pushed: yes`, proceed to final report; if workers in flight, re-dispatch; if subtask has diff, verify and mark `DONE` or re-dispatch; subtasks marked `DONE` are never redone. Continue from `Next action:`. Never run `git reset`, `git stash`, `git clean`, or `git checkout --`.
5. **Ledger Creation:** If not resuming, create `.ai-memory/temp-agents/NN-<slug>/ledger.md`:

```markdown
# Ledger: NN-<slug>
Request slug: <slug>
Request first line: <verbatim first line>
Status: ACTIVE
Phase: 1    Wave: 0 / WAVES    Step: 1 / N
Last completed action: Phase 1A Capture & Task Breakdown
Next action: Phase 1B Spec & Plan
Workers in flight: none
Commits: none    Pushed: no
Branch: <branch> | Tree at start: clean (or dirty with <paths>)
Tools: invoke_subagent=yes send_message=yes ask_question=yes gitmap=yes
| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|

---

## 6. Phase 1B: Spec, Plan & Lean Subtasks (Steps 1 .. PHASE_1_BUDGET)

1. **Single-Agent Blueprint:** Lead orchestrator alone authors initial spec overview and planning skeleton.
2. **Mandatory Parallel Discovery Subagents (A workers):** You MUST dispatch A research subagents (`Research 01 .. Research <A>`) via `invoke_subagent` on disjoint folders to map symbols and dependencies using GitMap commands (`gitmap f`, `gitmap lf`, `gitmap search`, `gitmap cat`) and yield the turn. Do not perform all discovery solo in the main agent.
   - *Research Contract:* Reply with one line per hit formatted as `path:line: text`, then `SUMMARY: <one line>`, then stop.
3. **Canonical Spec Authoring (`02-spec/21-app/`):** Single-domain (<= 150 lines): Write `02-spec/21-app/NN-<slug>.md`. Multi-domain: Write `02-spec/21-app/NN-<slug>/` (`01-overview.md` .. `04-verification-gates.md`). Register in `02-spec/21-app/readme.md`.
4. **Execution Plan:** Write `.ai-memory/plans/pending/NN-<slug>.md` linking to spec and mapping Task-IDs to subtasks. Register in `.ai-memory/plans/readme.md`.
5. **Root Task JSON Manifest & Subtask Files:**
   - Maintain root JSON `.ai-memory/plans/subtasks/NN-<slug>/task.json` describing task and indexing subtasks:
     `{"taskSlug":"NN-<slug>","status":"IN_PROGRESS","subtasks":[{"id":"Task-01","file":"01-<name>.json","owner":"Worker 01","status":"PENDING"}]}`
   - Each subagent creates and updates its assigned subtask in structured JSON format (`.ai-memory/plans/subtasks/NN-<slug>/01-<name>.json`), referenced directly from root JSON manifest.
   - Lean subtask companion: `.ai-memory/plans/subtasks/NN-<slug>/01-<name>.md` (Traceability ID, Spec Reference, Owned Files, Action, Acceptance Criteria, Targeted Verification).
6. **Parallel Subtask Decomposition:** Decompose deliverables across A workers so each worker receives disjoint target files. Solo execution without calling `invoke_subagent` is strictly banned.
7. **Readiness Gate:** Complete Phase 1 planning within `PHASE_1_BUDGET` steps, then proceed **UNCONDITIONALLY** into Phase 2.

---

## 7. Phase 2: Mandatory Worker Waves & Coding Guidelines Enforcement (Steps (PHASE_1_BUDGET + 1) .. N)

> [!CRITICAL]
> **MANDATORY `invoke_subagent` DISPATCH (ZERO SOLO EXECUTION):**
> You MUST spawn A workers (`TypeName: "self"`, up to H subtasks per worker) in parallel via `invoke_subagent`. Executing all subtasks solo in main agent without calling `invoke_subagent` is an immediate auto-reject failure on the same tier as Rule 0. Work MUST be partitioned across A workers.

### 7.1 Dispatch Payload (`invoke_subagent`)

The `invoke_subagent` payload holds A entries (`Worker 01 .. Worker <A>`):

```json
{
  "Subagents": [
    {
      "TypeName": "self",
      "Role": "Worker 01: [Assigned Feature/Module A]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "<Worker Brief Below>"
    },
    {
      "TypeName": "self",
      "Role": "Worker 02: [Assigned Feature/Module B]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "<Worker Brief Below>"
    }
  ]
}
```

### 7.2 Self-Contained Worker Brief (Eliminate Context Blindness)

Subagents spawn with clean context. The prompt envelope MUST inject complete instructions:

```text
You are Worker <NN> for task NN-<slug>. You have no prior chat context; this brief is your complete specification.

### Boundaries:
- Read any file in the workspace; edit only your Owned Files: <relative paths>.
- After C tool calls, stop and report what you have.
- A tool failing twice: reply "STATUS: BLOCKED" with exact error and stop. Never guess paths and never troubleshoot machine.
- Workers that find a secret stop and report "BLOCKED: secret at <file>:<line>". They do not handle it themselves.
- Adhere to R1, R2, and R11 by ID.

### Assigned Subtasks (up to H subtasks):
- Subtask 1: .ai-memory/plans/subtasks/NN-<slug>/01-<name>.md
- Subtask 2: .ai-memory/plans/subtasks/NN-<slug>/02-<name>.md (if assigned)

### 100% Non-Negotiable Coding Guidelines (AUTO-REJECT ON VIOLATION):
1. Positive booleans ONLY: use `is` and `has` prefixes exclusively. NEVER evaluate explicit `== true`. NEVER combine positive and negative checks in the same condition (`if isA && !isB` is BANNED).
2. Go Structured Errors: return `*appfault.AppError`, never bare `error`.
3. Function Sizing: <= 8 lines preferred, hard cap 15 lines. Extract domain structs and raw generics to `types.go`.
4. Strict Relative Git Paths: zero absolute filesystem paths and zero `file:///` URIs.
5. Repo Secrets: if any credentials or private tokens are needed, store them in the `repo-secrets` folder in the default work directory (via `gitmap rs`). Never commit secrets.
6. Zero Builds or Tests: NEVER run `go build`, `npm run build`, `go test`, or `pytest`.
7. Targeted Verification: Run only fast file-scoped linters (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`). A check scanning 0 files is a FAIL.

### Output Contract:
Write your subtask output to .ai-memory/plans/subtasks/NN-<slug>/01-<name>.json and reply with this JSON block, once per subtask, then stop:
{
  "task": "Task-01",
  "status": "DONE",
  "filesChanged": ["<path1>", "<path2>"],
  "checks": "<command> -> exit <code>, <files scanned>",
  "acceptance": { "ac1": "PASS <evidence>" },
  "assumptions": [],
  "blockers": []
}
```

### 7.3 Turn-Yielding & Verification Protocol

1. **Yield:** Print progress line (`Dispatched Worker 01 .. Worker <A> (wave k / WAVES); waiting for their results.`) and **STOP CALLING TOOLS**.
2. **Verify Worker Reports Independently:** Confirm `git diff --stat -- <owned files>` matches `filesChanged`, no files outside owned files modified, re-run targeted checks for `exit 0` on non-zero files.
3. **Reject Violations:** Send failures via `send_message`. On `BLOCKED`, lead does work and logs `LEAD_FALLBACK: <reason>`. After two failed rounds, mark `FAILED`, write RCA, continue (R13).
4. **Update Ledger:** Record status, evidence, changed paths in `ledger.md` via `replace_file_content`.
5. **Loop:** Dispatch subsequent waves until all subtasks are `DONE` or `FAILED`.

---

## 8. Phase 3: Consolidation, Evidence Verification & Atomic GitMap Push

1. **Consolidate Subtasks:** Merge completed subtasks into `.ai-memory/plans/completed/NN-<slug>.md`, logging real steps from ledger; link to canonical spec. Delete `.ai-memory/plans/subtasks/NN-<slug>/` and pending plan. Canonical spec in `02-spec/21-app/` stays permanently.
2. **Update Registers:** Update `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md`. Update `01-prompts/readme.md` and `.ai-memory/prompts.md` only when adding/modifying prompts. Register recent completed tasks before push gate.
3. **Push Gate Verification:** Before GitMap call, verify: (a) targeted checks exit 0 (>0 files), (b) secrets gate clean, (c) `.gitignore` covers caches, build outputs, logs, reports, `.env*` (untrack any tracked ignored files via `git rm --cached <file>` or `git rm -r --cached <dir>`). On failure: abort GitMap call; mark task `FAILED` with RCA.
4. **Atomic Commit & Push:** Call `gitmap cpf "<summary>"` (features) or `gitmap cpb "<summary>"` (fixes). Push rejected: `git pull --rebase` and re-run. Miss after push: allow one follow-up `gitmap cpb "<summary>"`, logged as `FOLLOW_UP_PUSH: <sha>`. Never amend pushed commits.

---

## 9. Final Report Format (Strict Vertical Lines)

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]` — [diff/check evidence]
- ❌ **Task-02: [Descriptive Task Title]** — `[Failed]` — [RCA link]

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Steps Used

- Step x / N (Phase 1: y / PHASE_1_BUDGET, Phase 2: z / PHASE_2_BUDGET), Wave k / WAVES

### Implementation Confidence Score

- Confidence: [passed checks / total checks]
- Rationale: [Verified evidence across all criteria, passing targeted linters, zero regressions]

Independent check: run /verify-parent-task-run <slug>

### 🤖 Independent AI Verification & Audit Prompt

(Emit self-contained audit prompt linking spec, plan, and modified files)
```

---

## 10. Targeted Verification Checks (R2)

Confirm scripts exist via harmless workspace call before invoking (R4). Run on changed files/folders only:

- **Coding Guidelines & Boolean Linter:** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`
- **Relative Path Linter:** `python linter-scripts/check-relative-paths.py`
- **Prompts & Spec Index Linter:** `python linter-scripts/check-prompts-loaded.py`
- **Markdown Link & Doc Path Linter:** `python 03-ai-scripts/22-doc-path-linter.py <folder>`
- **Sequence Integrity Linter:** `python linter-scripts/check-sequence-integrity.py`
- **Forbidden Strings Check:** `python linter-scripts/check-forbidden-strings.py`

---

## 11. Discovery Toolchain (GitMap Primary)

Use Section 3 GitMap table as primary. Python fallbacks (`03-ai-scripts/11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) apply only if GitMap is unavailable.

---

## 12. Issue Destination & RCA Routing

- **CI/CD & Workflow Failures:** `.ai-memory/cicd-issues/NN-<slug>.md`, indexed in `.ai-memory/cicd-index.md`.
- **Application Bugs:** `02-spec/22-app-issues/NN-<slug>.md` with 4-part RCA (Reproduction, Cause, Fix, Prevention), indexed in `02-spec/22-app-issues/readme.md`.
- **Failed Subtasks (R13):** Log RCA in `.ai-memory/memory/issues/` and link from `ledger.md`.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run the targeted checks (builds and full unit tests stay in CI per R1), and execute the final atomic GitMap commit and push before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
