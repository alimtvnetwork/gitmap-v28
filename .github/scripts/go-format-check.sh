#!/usr/bin/env bash
# go-format-check.sh — Hard-gate that runs `gofmt -l` and `goimports -l`
# over the entire gitmap/ Go module tree. Fails on any file that would
# be rewritten, surfacing the exact diff so the contributor can fix
# locally with `gofmt -w` / `goimports -w` and recommit.
#
# `goimports` is a superset of gofmt that ALSO normalizes import blocks
# (adds missing imports, removes unused ones, groups std / third-party
# / local imports in canonical order). Local-prefix is read from go.mod
# so internal imports are grouped in their own block.
#
# goimports is pinned to v0.24.0 per the project's static-analysis rule
# (no @latest tool installs). gofmt ships with the Go toolchain.
#
# Run from the gitmap/ working directory.

set -uo pipefail

# ---- gofmt check ---------------------------------------------------------
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "::error::The following .go files are not gofmt-clean:"
  echo "$unformatted" | sed 's/^/  /'
  echo ""
  echo "----- gofmt -d (diff) -----"
  # shellcheck disable=SC2086
  gofmt -d $unformatted
  echo "---------------------------"
  echo ""
  echo "Fix locally with:  cd gitmap && gofmt -w ."
  exit 1
fi
echo "✅ All .go files are gofmt-clean."

# ---- goimports check -----------------------------------------------------
go install golang.org/x/tools/cmd/goimports@v0.24.0
GOIMPORTS="$(go env GOPATH)/bin/goimports"
LOCAL_PREFIX="$(awk '/^module /{print $2; exit}' go.mod)"

unformatted=$("$GOIMPORTS" -l -local "$LOCAL_PREFIX" .)
if [ -n "$unformatted" ]; then
  echo "::error::The following .go files have goimports issues (import grouping or formatting):"
  echo "$unformatted" | sed 's/^/  /'
  echo ""
  echo "----- goimports -d (diff) -----"
  # shellcheck disable=SC2086
  "$GOIMPORTS" -d -local "$LOCAL_PREFIX" $unformatted
  echo "-------------------------------"
  echo ""
  echo "Fix locally with:  cd gitmap && goimports -w -local $LOCAL_PREFIX ."
  exit 1
fi
echo "✅ All .go files pass goimports (local prefix: $LOCAL_PREFIX)."
