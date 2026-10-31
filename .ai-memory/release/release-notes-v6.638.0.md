# Release Notes — GitMap v6.638.0

## Summary
GitMap v6.638.0 enhances the native `gitmap cat` inspection command with multi-token OR query matching (`-q "token1|token2|token3"`), file metadata inspection (`--info`, `-i`), and adds on-the-fly interactive teaching diagnostics to `gitmap py` when ad-hoc Python file reading (`open(...)`) is attempted.

## Key Changes
1. **Multi-Token Pipe Filtering in `gitmap cat`**:
   - `gitmap cat <filepath> -q "token1|token2|token3" -l <limit>` enables matching lines containing any of the given tokens (case-insensitive OR matching).
   - Fully replaces complex Python one-liners (`[print(l) for l in open(...)...]`).

2. **File Metadata & Statistics (`--info` / `-i` / `--stats`)**:
   - `gitmap cat <filepath> --info` reports exact line counts, byte sizes, and formats cleanly without dumping full file content.

3. **On-The-Fly Teaching in `gitmap py`**:
   - `gitmap py -c "open(...)"` intercepts ad-hoc file reading anti-patterns and surfaces immediate, non-blocking guidance recommending `gitmap cat` or native tools.

4. **Help & Documentation**:
   - Updated `cli/helpdoc/cat.md` with multi-token filtering and `--info` flags and usage examples.
