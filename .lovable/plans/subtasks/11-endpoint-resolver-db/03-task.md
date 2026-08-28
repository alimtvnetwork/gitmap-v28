# Subtask 3: Integration

1. Open `gitmap/cmd/committransfer.go`.
2. Locate `resolveCommitEndpoints`.
3. Wrap `leftRaw` and `rightRaw` with `resolveEndpointString`:
   ```go
   resolvedLeft := resolveEndpointString(leftRaw)
   resolvedRight := resolveEndpointString(rightRaw)
   left, err := movemerge.ResolveEndpoint(resolvedLeft, true, mmOpts)
   // ...
   ```
4. Verify tests pass.
