#!/usr/bin/env bash
# Smoke test: `gitmap history-pin` collapses every historical version
# of a file to its current content. Builds a temp git repo where X
# has 3 distinct historical states, runs `gitmap hpin X --no-push`,
# and asserts every historical revision of X now hashes to one value.
#
# Usage: smoke-history-pin.sh <path-to-gitmap-binary>
#
# Spec: spec/04-generic-cli/16-history-rewrite.md

set -euo pipefail

GITMAP_BIN="${1:?path to gitmap binary required}"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

git init --bare "$WORK/origin.git" >/dev/null
git clone "$WORK/origin.git" "$WORK/work" >/dev/null
cd "$WORK/work"

# 3 commits, X has a distinct content in each.
for n in 1 2 3; do
  echo "version $n contents of X" > X
  git add X
  git -c user.name=ci -c user.email=ci@x commit -m "X v$n" >/dev/null
done
git push origin main >/dev/null

# Sanity: 3 distinct hashes pre-rewrite.
PRE="$(for sha in $(git log --all --pretty=format:%H -- X); do git show "$sha:X" | sha256sum; done | sort -u | wc -l)"
test "$PRE" = "3"

SANDBOX_LINE="$("$GITMAP_BIN" history-pin X --no-push --keep-sandbox 2>&1 | tee /dev/stderr | grep 'sandbox kept at' || true)"
SANDBOX="$(echo "$SANDBOX_LINE" | sed -E 's/.*sandbox kept at //')"

if [ -z "$SANDBOX" ] || [ ! -d "$SANDBOX" ]; then
  echo "FAIL: could not parse sandbox path from output" >&2
  exit 1
fi

# Every historical revision of X in the sandbox must hash identically.
UNIQ="$(for sha in $(git -C "$SANDBOX" log --all --pretty=format:%H -- X); do
  git -C "$SANDBOX" show "$sha:X" | sha256sum
done | sort -u | wc -l)"

if [ "$UNIQ" != "1" ]; then
  echo "FAIL: history-pin left $UNIQ distinct content hashes (expected 1)" >&2
  exit 1
fi

echo "PASS: history-pin collapsed X to a single content hash across all history"

# ── Scenario B: multi-path pin + --message scoping ──────────────────
# Add a second tracked file Y with 2 distinct historical states, plus
# an unrelated file Z that pin must NOT touch.
for n in 1 2; do
  echo "Y v$n" > Y
  echo "Z v$n" > Z
  git add Y Z
  git -c user.name=ci -c user.email=ci@x commit -m "Y/Z v$n" >/dev/null
done
# A final commit that touches only Z (must be untouched by pin/--message)
echo "Z final" > Z
git add Z
git -c user.name=ci -c user.email=ci@x commit -m "Z only" >/dev/null
git push origin main >/dev/null

rm -rf "$SANDBOX"
SANDBOX_LINE="$("$GITMAP_BIN" history-pin X Y --no-push --keep-sandbox \
  --message "pinned by ci" 2>&1 | tee /dev/stderr | grep 'sandbox kept at' || true)"
SANDBOX="$(echo "$SANDBOX_LINE" | sed -E 's/.*sandbox kept at //')"

for P in X Y; do
  U="$(for sha in $(git -C "$SANDBOX" log --all --pretty=format:%H -- "$P"); do
    git -C "$SANDBOX" show "$sha:$P" | sha256sum
  done | sort -u | wc -l)"
  if [ "$U" != "1" ]; then
    echo "FAIL: multi-path pin: $P has $U distinct hashes (expected 1)" >&2
    exit 1
  fi
done

# Z must NOT be collapsed — it had 3 distinct states and we never asked
# for it.
ZUNIQ="$(for sha in $(git -C "$SANDBOX" log --all --pretty=format:%H -- Z); do
  git -C "$SANDBOX" show "$sha:Z" | sha256sum
done | sort -u | wc -l)"
if [ "$ZUNIQ" -lt 2 ]; then
  echo "FAIL: history-pin leaked into unrelated path Z (collapsed to $ZUNIQ)" >&2
  exit 1
fi

# --message must NOT touch the "Z only" commit.
ZONLY_KEPT="$(git -C "$SANDBOX" log --all --pretty=format:'%s' | grep -c '^Z only$' || true)"
if [ "$ZONLY_KEPT" != "1" ]; then
  echo "FAIL: --message leaked into untouched 'Z only' commit" >&2
  exit 1
fi

echo "PASS: multi-path pin (X, Y) + --message scoping; Z untouched"