# Plan 98: SSH Macro/PEA/PEAT Fleet Deployment, Multi-Node Update Telemetry, Remote SSH Clone & Isolated Temporary E2E Validation

> **Plan Status:** Completed  
> **Traceability IDs:** Task-01 .. Task-05  
> **Spec Reference:** [02-spec/21-app/149-ssh-macro-pea-deploy-fleet-update-and-ssh-clone-tempe2e.md](../../../02-spec/21-app/149-ssh-macro-pea-deploy-fleet-update-and-ssh-clone-tempe2e.md)

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 98-01** | `cmdmacro`, `cmd/rootutility.go` | `gitmap macro|peat|pea deploy ssh --except/--excep id,ip,alias` parallel fleet deployment | Completed |
| **Subtask 98-02** | `cmdupdate` | `gitmap update --all` / `gitmap update all` / `gitmap ua` with `--except` & JSON summary table | Completed |
| **Subtask 98-03** | `cmdupdate` | `gitmap update <name> --excep/--except` & `gitmap update ls` JSON software inventory table | Completed |
| **Subtask 98-04** | `cmdssh`, `cmd/roottooling.go` | `gitmap ssh clone` / `ssh-clone` / `ssh-c` with current repo URL & default workdir inference | Completed |
| **Subtask 98-05** | `cli/tests/e2e` | Isolated temporary E2E test suite (`//go:build tempe2e`, `RUN_TEMP_E2E=1`) & release `v6.324.0` | Completed |

---

## 2. Verified Outcomes

1. **Macro, PEAT & PEA Fleet Deployment**: Verified and enhanced `cli/cmdmacro/macro_deploy_ssh.go` and `cli/cmd/rootutility.go` to support `--except`, `--excep`, and `--exclude` filtering across host IDs, aliases, and IPs.
2. **Fleet Update & Inventory (`ua`, `update --all`, `update <name>`, `update ls`)**: Enhanced `cli/cmdupdate/update_fleet.go` and `cli/cmdupdate/update_fleet_ls.go` with `--excep` support and parallel JSON-to-table aggregation.
3. **Remote SSH Clone (`ssh clone`, `ssh-clone`, `ssh-c`)**: Enhanced `cli/cmdssh/ssh_clone_types.go` and `cli/cmdssh/ssh_clone_path.go` with `--excep` support, automatic current repository origin URL detection, and default remote workdir (`~/git/<repo>`) resolution.
4. **Isolated Temporary E2E Suite (`//go:build tempe2e`)**: Created `cli/tests/e2e/ssh_fleet_deploy_update_clone_tempe2e_test.go` and verified both active execution (`RUN_TEMP_E2E=1`) and default runtime skip guard (`RUN_TEMP_E2E` unset).
