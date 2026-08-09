# RCA: Git LFS smudge failure — `[404] Object does not exist on the server`

Repo: `git@github.com:alimtvnetwork/lara-licensing-v4.git`
Offending path: `assets/01-licensing.xmind`
LFS OID: `1929d9c77f0c1c44704e0a2d4809f93a5bca60b48397cd6a3f5d3b3ee951bb56`

## The Issue
`git clone` has two phases. Phase 1 fetches objects over the wire (succeeds because the pointer file is just text). Phase 2 writes worktree files and triggers `git-lfs filter-process`. If the LFS server returns a 404 for the OID, the filter exits non-zero and checkout is aborted.

## Machine-readable summary (for an AI agent)

```yaml
error_class: git_lfs_smudge_404
signature:
  - "Smudge error: Error downloading"
  - "Object does not exist on the server: [404]"
  - "external filter 'git-lfs filter-process' failed"
  - "Clone succeeded, but checkout failed"
meaning: git object db is intact; the LFS blob store is missing the referenced OID
never_do:
  - retry the clone unchanged (deterministic failure)
  - delete and re-clone hoping it resolves
  - git lfs pull (will hit the same 404)
first_action: set GIT_LFS_SKIP_SMUDGE=1 to obtain a working tree
diagnose:
  - git lfs ls-files -l
  - git lfs fetch --all
  - git cat-file -p HEAD:<path>
resolve:
  - blob_recoverable: re-add file, then `git lfs push origin --all`
  - blob_lost_and_file_unneeded: git rm --cached <path>; commit; push
  - blob_lost_and_path_needed: remove LFS rule from .gitattributes; commit placeholder
verify: fresh clone into a temp dir with `git lfs install` active must exit 0
prevent: pre-push LFS hook, nightly clone smoke test, `git lfs push --all` on every migration
```

## The Fallback Solution (PowerShell)
```powershell
$env:GIT_LFS_SKIP_SMUDGE=1; git clone git@github.com:alimtvnetwork/lara-licensing-v4.git; cd lara-licensing-v4; git restore --source=HEAD :/; git rm --cached "assets/01-licensing.xmind" -q; Remove-Item "assets/01-licensing.xmind" -Force -ErrorAction SilentlyContinue; git commit -m "chore(lfs): remove pointer for missing LFS object 1929d9c7"; git push; Remove-Item Env:GIT_LFS_SKIP_SMUDGE
```
