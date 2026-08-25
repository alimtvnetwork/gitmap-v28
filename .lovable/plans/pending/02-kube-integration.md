# 02-kube-integration.md

## Goal
Integrate the `scripts/kubernetes/` bash scripts into the `gitmap kube` CLI commands, enabling automated remote cluster provisioning via SSH.

## Constraints
- Max 3 agents per execution block.
- Follow `.lovable/temp/` lockfile rules.
- Test all CLI mappings locally without live SSH API calls.

## Execution Waves

### Wave 1: Core Command Scaffolding
- [ ] 001-task.md - Scaffold `gitmap kube` root command.
- [ ] 002-task.md - Scaffold `gitmap kube install` command.
- [ ] 003-task.md - Scaffold `gitmap kube init-master` command.

### Wave 2: Execution Logic & Script Embedding
- [ ] 004-task.md - Implement remote script streaming for `install`.
- [ ] 005-task.md - Implement token extraction for `init-master`.
- [ ] 006-task.md - Scaffold and implement `gitmap kube join-worker`.

### Wave 3: DB & Orchestration
- [ ] 007-task.md - Create SQLite migration for `kube_cluster`.
- [ ] 008-task.md - Implement DB read/write for tokens.
- [ ] 009-task.md - Implement `gitmap kube rollout` config JSON parser.
- [ ] 010-task.md - Implement rollout orchestration logic.
