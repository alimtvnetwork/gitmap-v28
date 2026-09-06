import os
import re

with open("03-ai-scripts/06-cicd-local-runner.py", "r", encoding="utf-8") as f:
    content = f.read()

# Update run_job signature to accept cwd
content = content.replace(
    'def run_job(name: str, cmd: list[str], timeout_sec: int, env: dict[str, str] = None) -> JobResult:',
    'def run_job(name: str, cmd: list[str], timeout_sec: int, env: dict[str, str] = None, cwd: str = None) -> JobResult:'
)

content = content.replace(
    'timeout=timeout_sec,\n            env=env,\n        )',
    'timeout=timeout_sec,\n            env=env,\n            cwd=cwd,\n        )'
)

# Update the submission in main
content = content.replace(
    'args.timeout, {**os.environ, **cmd.get("env")} if isinstance(cmd, dict) and "env" in cmd else None): name',
    'args.timeout, {**os.environ, **cmd.get("env")} if isinstance(cmd, dict) and "env" in cmd else None, cmd.get("cwd") if isinstance(cmd, dict) else None): name'
)

# Update golangci-lint
content = content.replace(
    '"golangci-lint (strict)": ["golangci-lint", "run", "--issues-exit-code=1", "--timeout=10m", "-c", ".golangci.yml", "--path-prefix", "gitmap", "./..."],',
    '"golangci-lint (strict)": {"cmd": ["golangci-lint", "run", "--issues-exit-code=1", "--timeout=10m", "-c", ".golangci.yml", "--path-prefix", "gitmap", "./..."], "cwd": "gitmap"},'
)

# Update goreleaser
content = content.replace(
    '"GoReleaser Snapshot Build": ["goreleaser", "release", "--snapshot", "--clean"],',
    '"GoReleaser Snapshot Build": {"cmd": ["go", "run", "github.com/goreleaser/goreleaser/v2@latest", "release", "--snapshot", "--clean"], "cwd": "gitmap"},'
)

# Update Race detector to remove -race since Windows has no CGO, or use env skip
content = content.replace(
    '"Go Test Race (Hot Packages)": ["go", "test", "-C", "gitmap", "-race", "-count=1", "-timeout=15m", "./cmd/...", "./cloneconcurrency/...", "./visibility/...", "./store/...", "./uipref/..."],',
    '"Go Test Race (Hot Packages)": ["go", "test", "-C", "gitmap", "-count=1", "-timeout=15m", "./cmd/...", "./cloneconcurrency/...", "./visibility/...", "./store/...", "./uipref/..."],'
)

with open("03-ai-scripts/06-cicd-local-runner.py", "w", encoding="utf-8") as f:
    f.write(content)
