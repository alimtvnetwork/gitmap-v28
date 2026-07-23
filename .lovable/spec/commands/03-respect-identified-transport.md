# Command: Always respect a repo's identified transport (SSH vs HTTPS)

**Scope:** all gitmap operations that read or act on a scanned repo.
**Applies to:** scan, clone, clone-next, pull, safe-pull, probe, clone scripts, terminal report, docs.

## Rule
1. Scan MUST persist BOTH `HTTPSUrl` and `SSHUrl` for every repo whenever derivable.
2. Scan MUST record the **identified transport** of `origin` (`ssh` | `https`) as a first-class field on the repo record.
3. Any subsequent operation that needs a URL MUST select the URL matching the identified transport. Falling back to the other transport is only allowed when the matching URL is empty, and MUST be logged.
4. Terminal scan report MUST display, per repo: identified transport, HTTPS url, SSH url, and the clone command in the identified transport.
5. Generated clone scripts MUST emit commands in the identified transport per repo (mixed scripts are allowed).
