STATUS: DONE

# Subtask 2: End-to-End Tests for commit-right

1. Create `gitmap/cmd/committransfer_e2e_test.go` or `gitmap/committransfer/e2e_test.go`.
2. Write a test that initializes two git repositories, creates 5 commits in repo A, and runs `commit-right` to copy them to repo B.
3. Assert that exactly 5 commits were transferred.
4. Clean up cleanly.
