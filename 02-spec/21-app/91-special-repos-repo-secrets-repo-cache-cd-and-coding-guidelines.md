# Spec 91: Special Repositories (`repo-secrets` / `rs` & `repo-cache` / `rc`), One-Time Scan Discovery, Sequenced File/Folder/Text Put & Auto-Push, `gitmap cd rs/rc`, and Coding Guidelines V2 Prompt Integration

Status: active
Version: 6.380.0

## User Request (Verbatim)

Improve this and also release minor and do a gitmap pe check

Okay. In the Git Map system, we can have some repositories. These are fixed. Okay? So one of them would be repo-secrets, another could be repo-storage. And these would be default repositories that would be in the settings. So from the settings, user can change it, the default names. Now, the way that this would work, so default names, and we can actually have default categories in the future as well, some of the default repositories. The way that it would work is that when the scan happens for the first time, the Git Map will try to see in their account if they have a repo-secrets, that could be in the folder or can be inside their, let's say, repository. If they have the same thing, then Git Map will automatically suggest, and if it is not there in the work directory, then Git Map will automatically suggest that, "Do you like to clone the repo-secrets?" And we use the repo-secret for keeping the secret files from the repository to here, so that your public repositories and others does not expose anything public or secrets. So that would be the purpose of it. Also mentioned that this is the default one. User can change it from the settings. Settings command will also be shared, how user can change it. Okay, this is one point. Another would be repo-storage. Repo-storage is a similar one as repo-secrets. The secrets will contain the secure stuff, password type thing, like environment variables, storage, things like that. But repo-storage is any type of weak, let's say, scripts or testing item that could be stored folder by folder as a sequence. Okay? So if those are not there, then also Git Map can recommend that, "Would you like to create those?" Each one, they could ask them a separate question and say, "This should be used for this purpose." Okay? So based on the answer that the user provides, if the answer is already given, then it's not like every scan is going to ask this question. Remember that. So once Git Map finds this, this would be considered as special repositories. This would be saved in a special repository table. Okay? And we can do RS, Git Map CD RS. RS would be repo-secrets, and it would have a shortcut, as I mentioned, LSC. Oh, so yeah. So how we can do the storage, right? So repo hyphen... Yeah, rather than saying storage, let's say cache. Repo-cache. So one would be repo-secrets, that is RS. Repo-cache, which is RC. So we can do the short forms to go to any places and make sure that these actually contains no issue at all. Do you understand this? And

Now, coming to a little bit of detail. So when we have the repo secrets, any environment variable or things, the Git map will automatically keep it there. And also, the coding guideline, I think we need to pull the coding guideline in. And also mention that coding guideline will also use this repo secret to put files, put text in a new file for the repo. So there will be a command to put things into the repo Git automatically. That would be the Git map, repo secrets or RC space, then the file. If we use the file command, we will put the file. If we use the folder command, folder. For which repo, we could mention the repo, or we could just run it inside the repo, then it would find that repo, whatever sequence that is, and it's going to put the file inside that repo so that next time also Git map knows where that is, also others know as well. So according to this, also update to the coding guideline sections, like if there is any secret stuff that needs to be saved into this repo secrets folder based on the repo sequence, XX hyphen repo name, and the sequence. And then, if there is this... If we want to store something like a PowerShell script or something that is temporary, but also we might need, then we can put it to the repo cache folder. That also needs to be observed and mentioned to the guideline, and also do it for the putting, I mean Git map as well. So it should have commands for file, folder. Even the text, we can say text that is saved as a file automatically based on the text given that we're trying to create a slug as an empty file. Or, yeah, it will always keep as an empty file, as a sequence, 0102 sequence automatically. So this will have its own idea, and it will always commit and push. So that is the idea using the Git map. And it should also appear in the help text. Okay, and also the UI help. And, yeah. And also we can use this to feed the LLM model. So the coding guideline prompts, v2 prompts, needs to have this mentioned so that they would understand why each folder is there, repository is there. And if they wanted to put something secret, how they put it, and commit and push it. So that can be added. Okay? And this still needs to be updated in the coding guideline. But working with coding guideline, you must do a Git pull before you change it. Are we aligned on this? Do you understand everything here?

---

## 1. Architecture & Data Contracts

### 1.1 SQLite Split-DB: `gitmap-special-repos.db` (`SpecialRepository` & `SpecialRepoMapping`)
- **`SpecialRepository` Table:**
  - `ShortKey TEXT PRIMARY KEY` (`"rs"` for `repo-secrets`, `"rc"` for `repo-cache`)
  - `DefaultName TEXT NOT NULL` (`"repo-secrets"`, `"repo-cache"`)
  - `ConfiguredName TEXT NOT NULL` (configurable via `gitmap settings --repo-secrets <name> --repo-cache <name>`)
  - `Category TEXT NOT NULL` (`"secrets"`, `"cache"`)
  - `LocalPath TEXT NOT NULL` (e.g. `D:\work\repo-secrets`, `D:\work\repo-cache`)
  - `RemoteURL TEXT NOT NULL DEFAULT ''`
  - `Purpose TEXT NOT NULL`
  - `IsPromptAnswered INTEGER NOT NULL DEFAULT 0` (1 once user answers or repo is detected/created, ensuring `gitmap scan` asks at most ONCE)
  - `UserDecision TEXT NOT NULL DEFAULT ''` (`"accepted"`, `"declined"`, `"detected"`)
  - `UpdatedAt TEXT NOT NULL`
- **`SpecialRepoFolderSeq` Table:**
  - `RepoName TEXT PRIMARY KEY` (e.g. `"gitmap"`, `"coding-guidelines"`)
  - `SeqPrefix TEXT NOT NULL` (e.g. `"01-gitmap"`, `"02-coding-guidelines"`)
  - `SeqNumber INTEGER NOT NULL` (`1`, `2`, ...)

### 1.2 Navigation & CLI Commands
- **Navigation Shortcuts:**
  - `gitmap cd rs` -> navigates directly to `repo-secrets` (`D:\work\repo-secrets`).
  - `gitmap cd rc` -> navigates directly to `repo-cache` (`D:\work\repo-cache`).
- **Special Repo Operations (`gitmap rs` / `gitmap repo-secrets` and `gitmap rc` / `gitmap repo-cache` / `gitmap repo-storage`):**
  - `gitmap rs file <filepath> [--repo <name>]`: copies `<filepath>` into `<repo-secrets>/<XX-repo-name>/<NN-filename>`, then stages, commits (`Secret: add <NN-filename> for <XX-repo-name>`), and pushes.
  - `gitmap rs folder <folderpath> [--repo <name>]`: copies `<folderpath>` into `<repo-secrets>/<XX-repo-name>/<NN-foldername>`, then stages, commits, and pushes.
  - `gitmap rs text "<text-or-slug>" [--slug <slug>] [--ext <ext>] [--repo <name>]`: derives a slug from `<text-or-slug>`, creates `<repo-secrets>/<XX-repo-name>/<NN-slug>.md` (or `.env`/`.txt`), writes content, stages, commits, and pushes.
  - `gitmap rc file <filepath> [--repo <name>]`: copies temporary/reusable script or test artifact into `<repo-cache>/<XX-repo-name>/<NN-filename>`, stages, commits, and pushes.
  - `gitmap rc folder <folderpath> [--repo <name>]`: copies folder into `<repo-cache>/<XX-repo-name>/<NN-foldername>`, stages, commits, and pushes.
  - `gitmap rc text "<text-or-slug>" [--slug <slug>] [--ext <ext>] [--repo <name>]`: creates `<repo-cache>/<XX-repo-name>/<NN-slug>.ps1` (or `.md`/`.txt`), writes content, stages, commits, and pushes.
  - `gitmap rs ls` / `gitmap rc ls [--repo <name>] [--json]`: lists sequenced repo folders (`01-gitmap`, etc.) and their `01-`, `02-` sequenced files.
  - `gitmap rs scan-check` / `gitmap scan`: runs one-time discovery for `repo-secrets` and `repo-cache`, prompting only if `IsPromptAnswered == 0` and never prompting again once answered.
