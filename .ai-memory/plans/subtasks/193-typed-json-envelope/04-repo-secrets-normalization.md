# Subtask 193.4: Repo-Secrets JSON Manifests Normalization
Traceability ID: Task-04
Spec Reference: [02-spec/21-app/183-typed-json-envelope-and-format-inspection.md](../../../02-spec/21-app/183-typed-json-envelope-and-format-inspection.md)
Target Files: ./repo-secrets\01-gitmap\*.json, ./repo-secrets\04-w1-machine\*.json, ./repo-secrets\05-w2-machine\*.json, ./repo-secrets\06-w3-machine\*.json, ./repo-secrets\07-final-network-machine\*.json
Action: Git pull ./repo-secrets. Transform JSON manifests into standard typed envelope with attributes and data. Keep secrets strictly within repo-secrets. Commit and push repo-secrets to origin/main.
Acceptance Criteria:
1. `./repo-secrets\01-gitmap\gitmap-ssh-nodes.json` and others adhere to typed envelope standard.
2. Machine folders (`04-w1-machine`, etc.) adhere to typed envelope standard.
3. No secrets or credentials leaked outside repo-secrets.
4. Committed and pushed to `origin/main`.
Targeted Verification: `git -C ./repo-secrets status`
