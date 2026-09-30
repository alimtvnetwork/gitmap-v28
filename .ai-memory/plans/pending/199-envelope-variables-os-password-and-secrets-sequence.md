# Plan 199: JSON Envelope Variables, WorkDirectory Object, OS Password CLI, and Repo-Secrets Sequence

Spec Reference: [02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md](../../../02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md)

## User Request (Verbatim)

```text
also make sure we have samples of

gitmap ssh join {user}@ip <alias> # then provide password if for first time and gitmap stores it using RSA algorithm using salt please
can reuse to login and also needs to prompt user for this as well, if the user don't prompt don't save please



Okay. Can we do a little bit changes here? Inside your, let's say, JSON, if you look into this, there are variables that have been part of the path. For example, the C users, administrator, and then the ID RSA. That is repeated everywhere. So you could just put that to variable and show case also how the variable can be reused. Okay. And, yeah. The OS Windows is there, but that does not mean that it has to be imported to Windows. Remember that. Okay. Remember to mention in the notes as well. Okay? So this is very important. Also here, you have created the work directory, single or flat line arguments. But what I requested is the work directory could be an object inside our JSON, and that object would have other aspects. And one of the aspects is the variable, that we could reuse inside the work directory as well. We wanted to see that the system can actually reuse the variables. So work directory variable could be coming from the work directory section, and also rename this Gitmap SSH node item to just SSH nodes import. Okay? There is no need to mention the Gitmap because it's already inside the Gitmap folder, right? So at the end, I want you to also commit fix to the repo secrets and also test using Gitmap end to end and do a final release, with the minor bump and make sure the CICD is green. Do you understand this? In the other, let's say, items, do you see any other inconsistencies? What do I mean by other items? Other items means the other JSON files into that secret node. And also, please confirm that Gitmap should be able to change OS password, right? We should have the change password option already to where we can give the username and the password, and it will change from the Gitmap if you have the admin rights. But also at the same time, if you're in the current user, if you just use the change password and just provide the password, that would automatically change the password by confirming. That should work on Windows, Linux, Ubuntu, macOS, all the cases. So that also needs to be in the OS help section, with proper examples. Do you understand? Can you please do that for me? Yeah. All the files that we have inside the Gitmap, it should start from zero one, zero two sequence, and no underscore. All the files needs to be lowercase. If we needed to, that's it. Is it clear? And fix all the file internal linkings. For example, the Azure template and the commit pull things, that needs to be updated according to the file and naming change. Okay. Can you please do that in the commit tool complete as well? Make sure that we have concise files, no duplicate files as well. Can you please do that for me? And also show a final report what you have optimized and what you have
```

## Extracted Actionable Task List

- **Task-01**: Multi-pass JSON Envelope variable interpolation & polymorphic `workDirectory` object support (`Subtask 01`).
- **Task-02**: `gitmap os change-password` single-argument current-user password mode (`gitmap os change-password <new-password>`) and dual-argument admin mode (`gitmap os change-password <user> <new-password>`) + `gitmap os --help` examples (`Subtask 01`).
- **Task-03**: `gitmap ssh join {user}@ip <alias>` interactive salted-RSA password prompt samples & CI lint/script regression fixes (`Subtask 02`).
- **Task-04**: Harmonize all JSON files in `D:\work\repo-secrets\01-gitmap` (`01-` through `09-`, including `07-seo-templates.json` and `08-seo-templates-cli.json`) and sibling machine folders (`04-w1-machine` through `07-final-network-machine`) (`Subtasks 03 & 04`).
