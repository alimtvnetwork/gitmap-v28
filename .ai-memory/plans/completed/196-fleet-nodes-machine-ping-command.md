# Plan 196: Fleet Nodes Machine Ping Command (`gitmap nodes ping` & `gitmap ping`)

## 1. Intent & Context
- The user requested:
  `let's have gitmap nodes ping and gitmap ping`
  `this is basically going to run the machine command and print nicely in the terminal for all nodes, alright?`
- Runs the native host operating system's `ping` command concurrently across registered fleet nodes.
- Augments raw ICMP ping with TCP port reachability to accurately identify nodes with Windows Defender Firewall restrictions dropping ICMP echo requests while allowing SSH.
- Supports target filtering, ad-hoc IP/hostname targets, `--count`, `--timeout`, `--raw`, `--json`, and rich ANSI box table formatting.

## 2. Tasks Executed
1. **Machine Ping Execution & Output Parsing**:
   - Implemented `cli/cmd/nodes_ping_cmd.go`.
   - Added cross-platform machine ping execution (`ping -n` on Windows, `ping -c` on Linux/macOS).
   - Added regex parsers for both Windows and Linux ICMP statistics and RTT metrics.
   - Added concurrent TCP port probing (`net.DialTimeout`).
   - Added status resolution (`● ONLINE`, `● REACHABLE (TCP)`, `○ OFFLINE`, `▲ DEGRADED`).
2. **Top-Level & Subsystem Dispatch**:
   - Wired `gitmap nodes ping [target]` in `cli/cmd/nodes_cmd.go`.
   - Wired `gitmap ping [target]` in `cli/cmd/rootcore.go`.
   - Updated topic catalog in `cli/helptext/catalog.go`.
3. **Unit Tests**:
   - Implemented `cli/cmd/nodes_ping_cmd_test.go` covering option parsing, Windows ping output parsing, Linux ping output parsing, status resolution, ANSI table rendering, JSON serialization, and help text.
   - All 7 tests passed (100% PASS).
4. **Live Verification**:
   - Built and updated `%USERPROFILE%\AppData\Local\gitmap-cli\gitmap.exe`.
   - Tested `gitmap nodes ping`, `gitmap ping`, `gitmap ping w1`, `gitmap ping --raw w1`, `gitmap ping --json w1`, `gitmap ping 127.0.0.1`, and `gitmap nodes ping --help`.
5. **Documentation & Specs**:
   - Authored specification `02-spec/21-app/186-fleet-nodes-machine-ping-command.md`.
   - Updated index in `02-spec/21-app/readme.md`.
