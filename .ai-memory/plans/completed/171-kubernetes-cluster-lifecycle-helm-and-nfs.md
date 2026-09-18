# Plan 171: Kubernetes Cluster Lifecycle, Helm & NFS Storage Provisioning Suite

## Executive Summary
Integrate the complete Kubernetes cluster installation, CRI-O container runtime configuration, cluster initialization, worker joining, Helm package management, and NFS dynamic storage provisioning from `kubernetes-training/03-kube-Installer` into GitMap CLI (`gitmap cluster k8s ...`). Passwords and tokens remain protected via SQLite database persistence and SSH RSA encryption at rest. Conclude with minor release bump to `v6.238.0`.

## Architecture Guidelines
- TOTAL BAN on running `go test`, `go build`, or `06-cicd-local-runner.py` during execution.
- Function length <= 15 lines (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal *apperror.AppError returns.
- Strict Unix LF line endings.
- Strict disjoint file assignments across subagents.

## Subtasks
1. `01-cluster-k8s-recipe-generators.md`
2. `02-cluster-k8s-command-and-join-pipeline.md`
3. `03-cluster-k8s-help-and-catalog.md`
4. `04-minor-release-v6-238-0.md`
