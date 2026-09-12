# Root Cause Analysis: CI Step Timeout Flakiness and Fixture Gofmt Backup Dirtiness

## 1. Why It Happened

Two failures were surfaced in remote GitHub Actions CI workflows (#34675150281 and #34675150171):
1. In macro/macro_test.go, TestExecute_StepTimeout failed on Ubuntu/Linux runners with signal: killed, took: 3.006s when expecting timeout under 2.8s.
2. In 	ests/fixrepo_test/gofmt_e2e_test.go, TestFixRepoGofmtCleanAfterRewrite failed on all platforms (Ubuntu, macOS, Windows) because gofmt -l . detected a dirty file inside .gitmap/backup/.../files/aligned_map.go.

## 2. How It Happened

1. **Macro Timeout**:
   - macro.buildStepCmd launched sh -c "sleep 3" on Unix.
   - When the 1-second step timeout fired, Go signaled the parent sh shell.
   - However, the spawned child sleep 3 remained alive, holding the stdout/stderr pipe write ends open.
   - Because xec.Cmd.WaitDelay was unset, cmd.Run() waited the full duration until the pipes reached EOF (when sleep 3 terminated after 3.003s), causing lapsed > 2.8s to fail.

2. **Fixture Gofmt**:
   - In ixture_helpers_test.go, lignedMapSource defined a string template with trailing empty lines after the closing map brace (}\n\n).
   - When setupFixtureRepo wrote ligned_map.go, the file was unformatted before ix-repo even executed.
   - ix-repo created a pre-rewrite snapshot in .gitmap/backup/, preserving those original unformatted bytes.
   - ix-repo then formatted the working tree copy of ligned_map.go, making it clean.
   - However, unGofmtList called gofmt -l . across the whole fixture root directory, which recursed into .gitmap/backup/ and flagged the preserved backup file as dirty.

## 3. Root Cause

1. uildStepCmd lacked cmd.WaitDelay, allowing orphaned subprocesses to hold stdout/stderr descriptors beyond the context deadline. Furthermore, getSleepCmd did not use xec on Unix, resulting in a child process fork.
2. lignedMapSource had superfluous trailing blank lines, making the initial fixture dirty, and unGofmtList did not exclude internal backup directories (.gitmap/ and .git/) from the working tree formatting check.

## 4. Code Fix

1. In gitmap/macro/execute.go: Added cmd.WaitDelay = 250 * time.Millisecond in uildStepCmd to forcibly terminate lingering I/O handles after context cancellation.
2. In gitmap/macro/macro_test.go: Changed Unix sleep command to xec sleep %d in getSleepCmd so sh replaces itself directly with sleep.
3. In gitmap/tests/fixrepo_test/fixture_helpers_test.go:
   - Stripped trailing blank lines from lignedMapSource.
   - Updated unGofmtList via ilterDirtyGofmtLines to ignore .gitmap/ and .git/ directories.
