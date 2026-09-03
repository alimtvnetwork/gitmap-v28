# Acceptance Criteria & Verification Matrix

## Verification Commands

```bash
go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1
go test ./gitmap/helptext/... -count=1
go test -v ./gitmap/cmd -run "TestPickProfile|TestNormalize|TestResolveBackup|TestApplyDefault" -count=1
python .github/scripts/e2e-cli-smoke.py bin/gitmap.exe
```

---

## Given / When / Then Contracts

### Scenario 1: Multi-Account Profile Listing & Sequence Selection

- **Given**: An initialized GitMap installation with discovered GitHub user and organization accounts.
- **When**: Running `gitmap profiles ls`.
- **Then**: An indexed table is printed displaying `#`, `Name`, `Provider`, `Type`, `Default`, `Usage`, and `Last Used`.
- **When**: Running `gitmap profiles set-default 2`.
- **Then**: The second account is persisted as default, and marked with `* (default)` in subsequent listings.

### Scenario 2: Repository Creation with Custom Organization

- **Given**: Authenticated GitHub credentials.
- **When**: Running `gitmap create my-service --org auktvgo --private --dir ./my-service`.
- **Then**: Local folder `./my-service` is initialized with git, `README.md`, and `.gitignore`, initial commit is made, and remote repository `auktvgo/my-service` is created privately and pushed.

### Scenario 3: Cloud Backup Snapshot Creation

- **Given**: Valid default profile and SQLite databases in `data/`.
- **When**: Running `gitmap backup create --note "Pre-update"`.
- **Then**: A new snapshot folder is committed into `gitmap-cloud-backup`, pushed to remote, and recorded in `gitmap backup ls`.
