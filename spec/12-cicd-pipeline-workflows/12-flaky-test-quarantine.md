# Flaky Test Quarantine Protocol

## 1. Definition of a Flaky Test

A test is considered flaky if it both passes and fails periodically without any changes to the code under test. Flaky tests erode developer trust and slow down the CI/CD pipeline.

## 2. Quarantine Process

- **Detection:** CI pipelines must be configured to detect tests that fail and then pass on a retry.
- **Isolation:** Flaky tests must be immediately marked with a `@quarantine` or `@skip` tag (depending on the framework). They MUST NOT block the main CI/CD pipeline.
- **Ticket Creation:** An automated or manual Jira/GitHub issue must be created immediately to investigate the root cause of the flakiness.
- **Resolution:** A quarantined test can only be reintroduced to the main suite after the root cause is identified, fixed, and the test passes 100 consecutive times in a burn-in environment.
