# Subtask 01: Root Usage Alignment & Parameter Structs

## Scope
1. Update `gitmap/cmd/rootusage.go`:
   - Encapsulate `renderStandardHelpRow` parameters into `HelpRowParams`.
   - Encapsulate `parseExpandableMarker` return into `ExpandableMarkerResult`.
2. Update `gitmap/cmd/rootusagefooter.go`:
   - Encapsulate `emitIdentityRows` parameters into `IdentityRowParams`.
3. Verify with `go test -v ./...` in `gitmap/cmd`.
