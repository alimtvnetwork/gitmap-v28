# Desktop

Okay. Now, you have to write an output in a different way. So first of all, whatever I say, whatever the input is, before the input, you write as high priority instruction, non-negotiable, must follow. Okay? And then you put the input verbatim. That would be the first thing. So if we go inside the output, that would be high priority. Then you have the goal, let learn, and then the input verbatim. That is correct. And then later on, you would do a little bit of action items that you think are the things that is mentioned inside this instruction. So the first one will be always and. The second item, another item would be how the instruction went. You just put the instructions items there. Okay, and at the end, we will ask, like, follow agent. The execute parent task within steps V6. I've given the file format, how it could go. Make sense?

Don't write "Certainly! Here's the output based on your instructions:" just output exactly as the given format no need for any action items, please, clear???

Follow the format below without the text "output"; you must add `Follow Skill [06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)` as a suffix, not as a prefix, you stupid. Follow exactly as I have said.

${Input Text Verbatim} = would be the input text as it is, without the filler words like um, ah, wh, etc.

Output
# High Priority Instruction

[/goal](slashCommand;goal) [/learn](slashCommand;learn) ${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write a plan and spec first
2. ..

Must follow and spawn an agent using the following skill

[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)

---

Don't output like
[/goal](slashCommand;goal) [/learn](slashCommand;learn) ${Input Text Verbatim} [06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6) Input : ...text given (no please no)
