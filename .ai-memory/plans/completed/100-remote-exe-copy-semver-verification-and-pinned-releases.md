# Plan 100: Remote Binary Installation & Copy via Low-Level SSH, Remote Fleet SemVer Verification & Comparison, and Specific Version Pinning/Downgrades for GitMap & Antigravity Manager

> **Status:** `COMPLETED`  
> **Release Version:** `v6.327.0`  
> **Spec Reference:** [`02-spec/21-app/151-remote-exe-copy-semver-verification-and-pinned-releases.md`](../../../02-spec/21-app/151-remote-exe-copy-semver-verification-and-pinned-releases.md)

---

## 1. Completed Subtasks

- [x] **SUBTASK-100-01**: Documented and verified native low-level SSH binary copy and execution architecture via `gitmap ssh cp <src> <node>:<dest>` using pure SSH channel primitives (Base64 memory buffer streaming, in-memory decoding, and atomic disk materialization without SMB, CIFS, NFS, or FTP). Enhanced Unix execution with automated `chmod +x` assignment for binary targets.
- [x] **SUBTASK-100-02**: Implemented remote fleet SemVer inspection and comparison in `cli/cmdssh/ssh_node_version.go` (`gitmap ssh nodes --version`). Evaluates remote node version against local version (`constants.Version`) with colorized markers: `● IDENTICAL (SAME)`, `▲ ABOVE (NEWER)`, `▼ BELOW (OLDER)`, and `○ NOT INSTALLED` / `○ OFFLINE`.
- [x] **SUBTASK-100-03**: Created `cli/cmd/version_tags.go` querying GitHub releases and tags via `/releases` API with fallback to `/tags`, semver integer comparison (`CompareSemver`), and terminal table formatter (`RenderReleaseTagsTable`).
- [x] **SUBTASK-100-04**: Integrated Antigravity Manager release tag listing and specific version targeting (`gitmap agm version ls`, `gitmap agm install <version>`, `gitmap agm update <version>` / `--version <version>`) in `cli/cmdinstall/agm_install.go`, `cli/cmdinstall/agm_update.go`, and `cli/cmdinstall/installagmanager.go`.
- [x] **SUBTASK-100-05**: Integrated GitMap release tag listing and specific version downgrading/upgrading (`gitmap version ls`, `gitmap versions`, `gitmap update <version>` / `gitmap update --version <version>`) in `cli/cmd/rootutility.go`, `cli/cmdupdate/update_installer_cmd.go`, and `cli/cmdupdate/updateremoteinstall.go`.
- [x] **SUBTASK-100-06**: Authored and validated isolated temporary E2E test suite in `cli/tests/e2e/version_ls_and_node_comparison_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`) passing all checks for SemVer comparison, table rendering, pinned version arguments, and low-level SSH write command generation.
