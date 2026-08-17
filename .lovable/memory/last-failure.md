# Last Failure
The issue is `01-ssh-repo-cloned-as-https.md`. The root cause is that Gitmap only persists HTTPS URLs on clone/scan, so when cloning an SSH repo via `gitmap clone gitmap.json`, it falls back to HTTPS, triggering browser auth.
