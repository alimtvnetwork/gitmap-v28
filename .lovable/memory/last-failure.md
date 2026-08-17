# Last Failure
The issues are `02-reclone-loses-ssh-transport.md` and `03-no-gitmap-code-command.md`. 
The root cause for 02 is that `cfr` (and similar reclone commands) use a URL picker (`cloner.pickURL` / `clonefixrepo.go` / `clone.go` / `clonenow.go`) that hard-codes HTTPS-first because the `Repo` schema lacks an `IdentifiedTransport` column to persist the original transport used (e.g., SSH). Thus, transport is never remembered by next reclone. 
The root cause for 03 is that the `gitmap code` (and aliases `vcode`, `vscode`) command dispatch entry and its handler (`cmd/code.go`) simply do not exist.
