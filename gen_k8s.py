import os
import json

tasks = [
    ("001", "Scaffold gitmap kube command", "Create gitmap/cmd/kube_cmd.go and register it to root"),
    ("002", "Scaffold gitmap kube install", "Create gitmap/cmd/kube_install_cmd.go for the install step"),
    ("003", "Scaffold gitmap kube init-master", "Create gitmap/cmd/kube_init_cmd.go for master init"),
    ("004", "Implement remote script streaming", "Add logic to read 02-ubuntu-prereq and 03-kube-install and stream via SSH"),
    ("005", "Implement token extraction", "Parse stdout of init-master to grab kubeadm join token and hash"),
    ("006", "Implement join-worker", "Create gitmap/cmd/kube_join_cmd.go to execute join command remotely"),
    ("007", "Create SQLite migration", "Add kube_cluster table migration to gitmap/store/migrations_kube.go"),
    ("008", "Implement DB repo", "Add functions to insert and retrieve tokens in gitmap/store/kube_repo.go"),
    ("009", "Implement rollout config parser", "Create gitmap/cmd/kube_rollout_cmd.go to read config.json"),
    ("010", "Implement rollout orchestration", "Add the go-routine orchestration to run all steps automatically")
]

os.makedirs(".lovable/plans/subtasks/kube-integration", exist_ok=True)
for tid, title, desc in tasks:
    with open(f".lovable/plans/subtasks/kube-integration/{tid}-task.md", "w", encoding="utf-8") as f:
        f.write(f"# Task {tid}: {title}\n\n## Description\n{desc}\n\n## Verification\nRun `go test ./...` and `go build ./...`.\n")
