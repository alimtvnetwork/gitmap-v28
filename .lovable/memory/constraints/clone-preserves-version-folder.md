---
name: Base clone never strips -vN folder names
description: gitmap clone (single or multi URL) must keep the repo name verbatim, including trailing -vN; only clone-next flattens versions
type: constraint
---

# Base `clone` must not rewrite folder names

`gitmap clone <url> [url...]` must create the destination folder with
the repo name exactly as it appears in the URL, **including any
trailing `-vN` suffix**.

- `https://github.com/owner/codex-june-6-v2.git` clones into
  `codex-june-6-v2/`, never `codex-june-6/`.
- An explicit folder argument (single-URL form) wins verbatim.
- Single-URL and multi-URL paths must agree. Both go through
  `resolveCloneFolder` in `gitmap/cmd/clonemulti.go`.

**Why:** users expect the folder to match what they typed. The old
multi-URL behavior called `clonenext.ParseRepoName` and flattened
versioned names, silently colliding distinct `-vN` repos into one
folder (fixed in v6.83.0).

**How to apply:** never reintroduce `clonenext.ParseRepoName` (or any
`-vN` regex) into the `clone` path. Version flattening and version
bumping belong to `gitmap clone-next` / `cn` only. Pinned by
`TestResolveCloneFolderPreservesVersionSuffix` and
`TestRepoNameFromURLKeepsVersionSuffix` in
`gitmap/cmd/clonemulti_folder_test.go`.
