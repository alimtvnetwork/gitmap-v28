# Plan: GitMap Command Enhancements

- **Slug:** `gitmap-command-enhancements`
- **Spec Reference:** `02-spec/21-app/gitmap-command-enhancements/01-architecture-spec.md`
- **Component Spec:** `02-spec/21-app/gitmap-command-enhancements/02-component-spec.md`
- **Status:** `active`
- **Created:** 2026-10-16

---

## User Request (Verbatim)

```text
PS C:\Users\Administrator> gitmap os --json
gitmap os: execute failed: [E_INVALID_OS_SUBCMD:EXECUTION] unknown os subcommand '--json' (see 'gitmap os --help'): (at=cmdos/os.go:95)
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdos.unknownOSSubcommandError (cmdos/os.go:95)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdos.dispatchOSSubcommand (cmdos/os.go:90)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdos.runOS (cmdos/os.go:17)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdos.RunOS (cmdos/exports.go:12)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.utilitySystemEntries.func2 (cmd/rootutility.go:136)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatchTable (cmd/rootdispatch.go:24)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatchUtility (cmd/rootutility.go:48)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.dispatch (cmd/root.go:593)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runDispatch (cmd/root.go:178)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.Run (cmd/root.go:151)
    at main.main (cli/main.go:7)
PS C:\Users\Administrator>


powershell -NoProfile -Command '$cpu = (Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average; $os = Get-CimInstance Win32_OperatingSystem; $tot = [math]::Round($os.TotalVisibleMemorySize/1024); $free = [math]::Round($os.FreePhysicalMemory/1024); $used = $tot - $free; $ram = [math]::Round(($used/$tot)*100); Write-Output ("{0},{1},{2},{3}" -f [int]$cpu, [int]$ram, [int]$tot, [int]$used)'




# GitMap Command Enhancements: high priority instruction, non-negotiable task

Can you please create, or you should have delegation command and also itself has command for `gitmap` to get the CPU information from the nodes, CPU RAM users, and also what are the tasks that they are running, these type of things. For example, if the `gitmap` is already delegated to a list of tasks, because every task needs to be enqueued in the task list. This way, `gitmap` knows how many tasks are actually requested. So based on that, `gitmap` can give this information like how many tasks, and tasks can be also identified as heavy duty tasks, and that can be revealed from one node to another for now. In the future, we will add role-based system. So in future, you can put a question mark there, which we will discuss later on. But I think you should have all kinds of things, like `gitmap` processor should give the processor information, power or CPU. If you just do `gitmap` CPU, it gives the CPU information, RAM, RAM users, CPU users, hard drive users, so all these things it should give instantly. If someone asked very frequently, it would usually do a caching of five seconds, so that there is not too much heavy load into the CPU. So you need to consider also Ubuntu, Windows, Windows Server, all these things. Along with that, it would also put the machine information, like the OS version, a little bit. If we, again, give a `gitmap` OS, that should also give the OS information. OS have commands, so after the commands at the end, it should actually show, just like the machine information, what we have, just put at the end nicely. And also the machine command as a sample, it should put into the OS when we just put. So machine, because machine we can do OS space machine, the same thing. So OS space JSON should also give us the OS information in JSON format. The same thing, CPU should give the information in to the CPU format. So currently, the OS section actually fails. It does not give any results. I think you need to work on that. Make sure that there is no glitches. Make sure you test it with the node two. I have added the VM two, which you can test anything. So you don't need to worry about if anything breaks in the VM, because I can reset it. So you can run anything that you like. You can test it. Also be mindful of it. Again, try to write code which is more efficient, not a virus. Remember that. We don't want to perform any harm. So based on that, you try to improve this, add new commands, and make sure we also have some kind of operation so that PowerShell can get the average information and things like that. So these type of things we should have inside from our `gitmap` that it would actually give very powerful information. Also, the network users, if we type the `gitmap` OS network users, it will give us the network users, what it is capable of, how much it is using right now, things like that, and also its IP adapter, MAC information, things like that. Also, we can get these using the nodes. So nodes should have sub commands, delegation commands, things like that. So craft those commands right, all the specs in detail nude, and then finally do the end-to-end testing with another machine, confirm everything, make release bump, and do a minor bump and release, and make sure that CICD is green. Do you understand?

## slug: gitmap-command-enhancements

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/gitmap-command-enhancements/ and enqueue plan task in .ai-memory/plans/gitmap-command-enhancements.md (subtasks in .ai-memory/plans/subtasks/gitmap-command-enhancements/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Use `gitmap` AI agents to enter data
5. Task completion includes committing and pushing to Git
6. Create delegation commands for `gitmap` to retrieve CPU, RAM, and task information from nodes.
7. Implement caching mechanism for frequent requests to reduce CPU load.
8. Ensure compatibility with Ubuntu, Windows, and Windows Server.
9. Fix the OS section to provide accurate results and format outputs in JSON where applicable.
10. Conduct end-to-end testing on VM two and ensure no harmful code is written.
11. Add operations for PowerShell to retrieve average information.
12. Implement commands to retrieve network user information and related details.
13. Perform end-to-end testing, confirm functionality, make a release bump, and ensure CICD is green.
```

---

## User Uploaded Context

![User Console Output 1](assets/screenshots/gitmap-command-enhancements-01.png)
![User Console Output 2](assets/screenshots/gitmap-command-enhancements-02.png)

---

## Architecture & Subtask Decomposition

### Wave 1: Core OS & Fast Hardware Telemetry
- **Subtask 01 (`01-os-json-and-machine-footer.md`):** Fix `cli/cmdos/os.go` and `cli/cmdos/os_cmd.go` to handle `--json`, `-j`, and `json` cleanly without throwing `E_INVALID_OS_SUBCMD`. Enhance default `gitmap os` call to append live machine/OS identity summary at the end.
- **Subtask 02 (`02-fast-hardware-telemetry-and-cache.md`):** Implement `gitmap cpu`, `processor`, `ram`, `disk` (and `gitmap os cpu`, `ram`, `disk`) with Windows PowerShell CIM execution, Linux `/proc` readers, and a 5-second cross-process cache in `cli/cmdos/` and `cli/cmd/rootutility.go`.

### Wave 2: Network Adapters, Nodes Tasks Delegation & UI/Release
- **Subtask 03 (`03-os-network-adapter-metrics.md`):** Implement `gitmap os network` (and `gitmap network`) in `cli/cmdos/` to query adapter names, status, IPv4, IPv6, MAC address, link speed, and byte counters across Windows and Linux with `--json` support.
- **Subtask 04 (`04-nodes-tasks-queue-and-delegation.md`):** Implement `gitmap nodes tasks` in `cli/cmdnodes/` to query remote nodes for CPU, RAM, enqueued task counts, active tasks, and heavy-duty task flags via SSH; record future role-based scheduling ambiguity in `.ai-memory/ambiguous-questions/01-new-ambiguity/`.
- **Subtask 05 (`05-docs-help-and-ui-registry.md`):** Author help topics in `cli/helpdoc/` and `docs/commands/`, and register `cpu`, `ram`, `disk`, `network`, and `nodes tasks` in `src/data/commands.ts`.
- **Subtask 06 (`06-verification-and-release-ceremony.md`):** Execute targeted verification checks, test against available nodes, perform minor version bump release ceremony (`v6.554.0`), and verify clean CI/CD pipeline state.
