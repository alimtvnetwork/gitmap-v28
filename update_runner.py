import os
import re

with open("03-ai-scripts/06-cicd-local-runner.py", "r", encoding="utf-8") as f:
    content = f.read()

content = content.replace(
    'def run_job(name: str, cmd: list[str], timeout_sec: int) -> JobResult:',
    'def run_job(name: str, cmd: list[str], timeout_sec: int, env: dict[str, str] = None) -> JobResult:'
)

content = content.replace(
    'timeout=timeout_sec,\n        )',
    'timeout=timeout_sec,\n            env=env,\n        )'
)

content = content.replace(
    'executor.submit(run_job, name, cmd, args.timeout): name',
    'executor.submit(run_job, name, cmd.get("cmd") if isinstance(cmd, dict) else cmd, args.timeout, {**os.environ, **cmd.get("env")} if isinstance(cmd, dict) and "env" in cmd else None): name'
)

# And now I'll inject the missing jobs directly by text replacement.
new_jobs1 = '''
            "Constants Collision Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdConstantsAreUnique", "-count=1"],
            "Helptext Parity Check": ["go", "test", "-C", "gitmap", "./helptext/...", "-count=1"],
            "Lint Script Unit Tests": [sys.executable, ".github/scripts/tests/test_ci_scripts.py"],
            "golangci-lint (strict)": ["golangci-lint", "run", "--issues-exit-code=1", "--timeout=10m", "-c", ".golangci.yml", "--path-prefix", "gitmap", "./..."],
            "Cross-OS Vet (Windows)": {"cmd": ["go", "vet", "-C", "gitmap", "./..."], "env": {"GOOS": "windows", "GOARCH": "amd64"}},
            "Cross-OS Vet (Darwin)": {"cmd": ["go", "vet", "-C", "gitmap", "./..."], "env": {"GOOS": "darwin", "GOARCH": "amd64"}},
'''
content = content.replace(
    '''"Constants Collision Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdConstantsAreUnique", "-count=1"],\n            "Helptext Parity Check": ["go", "test", "-C", "gitmap", "./helptext/...", "-count=1"],''',
    new_jobs1.strip()
)

new_jobs2 = '''
            "Web App Build": ["npm", "run", "build"],
            "GoReleaser Snapshot Build": ["goreleaser", "release", "--snapshot", "--clean"],
'''
content = content.replace(
    '"Web App Build": ["npm", "run", "build"],',
    new_jobs2.strip()
)

new_jobs3 = '''
            "E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],
            "Installer Smoke (source)": [sys.executable, ".github/scripts/smoke-installer.py", "source"],
            "Installer Smoke (release)": [sys.executable, ".github/scripts/smoke-installer.py", "release"],
            "Go Test Race (Hot Packages)": ["go", "test", "-C", "gitmap", "-race", "-count=1", "-timeout=15m", "./cmd/...", "./cloneconcurrency/...", "./visibility/...", "./store/...", "./uipref/..."],
'''
content = content.replace(
    '"E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],',
    new_jobs3.strip()
)

with open("03-ai-scripts/06-cicd-local-runner.py", "w", encoding="utf-8") as f:
    f.write(content)
