# System for Folder-Specific Work Directories and SSH Key Automation

## User Request (Verbatim)

```text
# System for Folder-Specific Work Directories and SSH Key Automation: high priority instruction, non-negotiable task

Can you create a system where folder by folder or also the work group, we can have different work directory, right? So based on that or folder by folder, we can have different types of login or SSH key automatically applied when we do the Git pull, Git push, Git pull, Git everything automatically. So can you customize and code like that and also provide that help for grouping, provide new commands, grouping and help on also inside the LLM train section and also docs in the help text and provide examples in the UI help, HD help. Okay? So provide that and make sure that you do the end-to-end testing as much as possible. You can create separate SSH keys or separate login things that is folder specific. And once it is set up, it does not need to speak again or set again. There should be a system that the Gitmap remembers it through its root Database (Database (Db)), SQLite Database (Database (Db)), and then it can always apply the format and things like that. And also the CPF and CBF, the module name can be put in the brackets and then the colon and then the item should be written. So you can actually enhance the format naming that should be for feature bugs, release, chores, everything can be written in the brackets, which module or whatever that is working on. Can you please do that for me and also update the help text and also coding guideline. Gitmap skills needs to be updated and synced with all the other repositories. Do you understand? Can you please follow through

## task: System for Folder-Specific Work Directories and SSH Key Automation
## slug: system-for-folder-specific-work-directories-and-ssh-key-automation
```

## Architectural Overview

This specification establishes an automated, zero-friction, folder-specific and workgroup-specific working directory and SSH credential system in GitMap:
1. **Work Group & Folder Mapping Engine**: Associating filesystem folders, parents, and work groups with dedicated Git working directories, profiles, and SSH credentials in GitMap's SQLite Split-DB (`gitmap.db` / `installation.db`).
2. **Autonomous Git Operation Dispatch**: Intercepting and decorating `git pull`, `git push`, `git fetch`, and fleet sync operations to dynamically inject the bound SSH private key (`-i <key> -F /dev/null -o IdentitiesOnly=yes`) and author identity based on current repository or folder lineage.
3. **Persistent SQLite Database State**: Zero prompt re-asking; once registered or bound, GitMap persists associations indefinitely across restarts and node reboots.
4. **Enhanced Commit Formatting (`cpf` & `cbf`)**: Enhancing GitMap commit and branch pushers to support bracketed module syntax: `[module]: <summary>` and type-prefixed `type([module]): <summary>` across features, bugs, releases, and chores.
5. **Interactive Grouping Commands & Diagnostics**: Providing dedicated commands for grouping, listing, binding, and diagnosing folder/group work directories.
6. **Documentation, Help Menus & LLM Train Integration**: Synchronizing UI help, HD help docs (`helpdoc/`), and the LLM self-training index (`gitmap llm train`).
