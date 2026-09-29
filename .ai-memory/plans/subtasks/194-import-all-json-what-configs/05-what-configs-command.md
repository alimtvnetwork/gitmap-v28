# Subtask 05: `gitmap what-configs` (`wc`) Inspection Command
- Created `cli/cmd/what_configs_cmd.go` with multi-format detection (JSON envelopes, YAML, SQLite DB, TOML, ENV, INI).
- Implemented ANSI table report showing file, detected category, recommended command, and system impact.
- Registered command and aliases `what-configs`, `wc`, `what-config`, `whatconfigs` in `cli/cmd/roottooling.go`.
- Added unit tests in `cli/cmd/what_configs_cmd_test.go` and verified PASS.
- Status: Completed.
