package cmdhooks

// commitMsgGuard is the commit-msg hook installed by `hooks install --all`.
// It enforces the house rule: no colons in the commit subject line.
const commitMsgGuard = `#!/bin/sh
# gitmap commit-msg guard (installed by ` + "`gitmap hooks install --all`" + `)
# House rule: the commit subject line must not contain a colon.
subject=$(head -n1 "$1")
case "$subject" in
*:*)
  echo "commit-msg guard: subject must not contain ':'" >&2
  echo "  got: $subject" >&2
  exit 1
  ;;
esac
exit 0
`

// guardMarker identifies hooks previously installed by this command.
const guardMarker = "gitmap commit-msg guard"
