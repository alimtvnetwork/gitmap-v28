# Step 3: Wire reclone to reuse transport

This step required ensuring that `cfr`, `cfrp`, `clone-now`, `reclone`, and direct `clone` all coerce the URL to the stored transport from the DB (if one exists) and persist the transport into the DB once cloned successfully.
This ensures an SSH-cloned repository is not silently downgraded to HTTPS when re-cloned.

## Changes:

- **`cmd/clonefixrepo.go`**: Was already wired to use `coerceURLToStoredTransport` and `persistRecloneTransport` by a previous task.
- **`cmd/reporeclone.go`** (`gitmap reclone <path>`): Added calls to `coerceURLToStoredTransport` before cloning and `persistRecloneTransport` after a successful clone.
- **`cmd/clone.go`** (`gitmap clone <url>`): Updated `executeDirectClone` to coerce the URL before issuing the `clone` command and print block, and added `persistRecloneTransport` on both replace and no-replace success branches.
- **`cmd/clonemulti.go`** (`gitmap clone <url1> <url2> ...`): Updated `executeDirectCloneOne` to similarly call both helpers.
- **`cmd/clonenow.go` & `clonenow/clonenow.go`** (`gitmap clone-now <file>`): Injected `CoerceURL` and `PersistURL` callbacks into the `clonenow.Plan` struct so that the execution engine inside `clonenow/execute.go` can invoke `cmd`-level helper functions without triggering circular package dependencies.

The `gitmap history` logging requirement is already implemented directly inside `persistRecloneTransport`, ensuring it fires whenever a repo's transport is persisted!

With these additions, all reclone entrypoints properly honor and persist the `IdentifiedTransport`.
