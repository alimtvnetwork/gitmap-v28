# Subtask 04: Minor Release v6.238.0 Orchestration

## Objective
Normalize line endings, format code across modified Go files, stage changes, and orchestrate the automated minor release `v6.238.0` with release notes detailing Plan 170 (Cluster Runner & Ubuntu Provisioning) and Plan 171 (Kubernetes Cluster Lifecycle, CRI-O, Helm & NFS Storage Provisioning).

## Actions
1. Run formatters and linters:
   - `python 03-ai-scripts/26-go-code-formatter.py`
   - `python 03-ai-scripts/04-newline-fixer.py --fix`
   - `python 03-ai-scripts/10-encoding-normalizer.py --fix`
2. Move Plan 171 from pending to completed:
   - Move `.ai-memory/plans/pending/171-kubernetes-cluster-lifecycle-helm-and-nfs.md` to `.ai-memory/plans/completed/171-kubernetes-cluster-lifecycle-helm-and-nfs.md`.
3. Commit all staged changes:
   - `git add -A`
   - `git commit -m "feat(cluster): Kubernetes cluster lifecycle, CRI-O runtime, Helm and NFS storage suite"`
4. Run minor release orchestrator:
   - `python 03-ai-scripts/29-release-orchestrator.py --tier minor --skip-tests`

## Constraints
- TOTAL BAN on running `go test`, `go build`, or CI runner.
