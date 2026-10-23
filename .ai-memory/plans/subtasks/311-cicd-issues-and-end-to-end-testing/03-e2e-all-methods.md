# Subtask 311-03: End-to-End Testing for All Methods

## Objective

Prove every registered gitmap method/command works end-to-end, using the
repo's own e2e and smoke harnesses executed in CI.

## Steps

1. Enumerate the method surface from the command registry
   (cli/constants registry and the helpdoc index) via gitmap aum search /
   gitmap find.
2. Map each method to the harness that exercises it: Installer Smoke
   (source-build, linux + windows), history-rewrite smoke, command
   invocation tests, and any spec-244-style e2e verify programs present in
   the repo.
3. For any method with no harness coverage, add the smallest e2e case that
   proves a real invocation (help surface + one run against a scratch repo).
4. Record the method-coverage list in
   02-spec/21-app/311-cicd-issues-and-end-to-end-testing/02-component-spec.md
   evidence section.

## Coordination

- cli/cmd/cmd_invocation_test.go is currently dirty in the spec-310 stream's
  box. Do not edit it until that stream commits or hands it over; cover other
  harnesses first.

## Exit Criteria

- Every registered method appears in the coverage list with a passing CI
  harness named next to it.
- All smoke/e2e jobs green on the release commit.
