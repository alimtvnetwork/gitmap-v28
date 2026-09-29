# Subtask 192.2: Version-Specific Update with `--pin` and Fleet `--ssh` Flag

## Spec Reference
- `02-spec/21-app/182-version-pinning-macro-deploy-ui-settings-secret-flags.md` §2.1

## Deliverables
1. **Targeted Update with Version & Pinning**:
   - In `cli/cmdupdate/` and `cli/cmd/rootutility.go`, support `gitmap update -v <version>` and `gitmap update --version <version>`.
   - Support `--pin` (or `-pin`) flag on update for both `gitmap` and `agm`/`agy`.
   - When `--pin` is supplied, write the pinned version to `installation.db`.
2. **Fleet Remote Execution via `--ssh`**:
   - Support `gitmap update -v <version> --pin --ssh` (and `gitmap update --ssh --version <version> [--pin]`).
   - Propagate version and pinning commands across all reachable cluster targets.
