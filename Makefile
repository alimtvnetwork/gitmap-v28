.PHONY: lint vet test build clean setup vulncheck all release release-dry changelog changelog-check goldens-regen goldens-verify fixtures-bump fixtures-bump-verify

GO       := go
LINT     := golangci-lint
MODULE   := gitmap
BINARY   := gitmap
VERSION  ?= dev
# Build-time identity (v5.60.0+) — stamps source repo metadata into the
# `gitmap binary` footer block so it cannot fall back to the user's CWD.
BUILD_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null)
BUILD_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
BUILD_REPO   ?= $(shell git config --get remote.origin.url 2>/dev/null)
BUILD_DATE   ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  := -s -w \
  -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/constants.Version=$(VERSION)' \
  -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildCommit=$(BUILD_COMMIT)' \
  -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildBranch=$(BUILD_BRANCH)' \
  -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildRepo=$(BUILD_REPO)' \
  -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildDate=$(BUILD_DATE)'

# Force bash for recipes — fixtures-bump targets rely on $${PIPESTATUS[0]}
# which is a bash builtin (dash and POSIX sh do not provide it).
SHELL := /bin/bash

all: lint test build

## Setup — install tools and git hooks
setup:
	@./setup.sh

## Lint — run golangci-lint
lint:
	@cd $(MODULE) && $(LINT) run ./... --timeout=5m

## Vet — run go vet
vet:
	@cd $(MODULE) && $(GO) vet ./...

## Test — run all tests
test:
	@cd $(MODULE) && $(GO) test ./... -v -count=1

## Build — compile for the current platform
build:
	@cd $(MODULE) && CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o ../$(BINARY) .
	@echo "Built $(BINARY) ($(VERSION))"

## Vulncheck — scan for known vulnerabilities
vulncheck:
	@cd $(MODULE) && $(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Release — run full release workflow (usage: make release BUMP=patch)
BUMP ?= patch
release: lint test
	@cd $(MODULE) && $(GO) run . release --bump $(BUMP)

## Release dry-run — preview release without executing
release-dry:
	@cd $(MODULE) && $(GO) run . release --bump $(BUMP) --dry-run

## Clean — remove build artifacts
clean:
	@rm -f $(BINARY)
	@rm -rf $(MODULE)/.gitmap/release-assets
	@echo "Cleaned."

## Changelog — regenerate CHANGELOG.md and src/data/changelog.ts from
## Conventional Commits since the latest annotated git tag.
## Usage:
##   make changelog VERSION=v3.92.0
##   make changelog VERSION=v3.92.0 SINCE=v3.90.0          # partial backfill
##   make changelog RELEASE_TAG=v3.91.0 SINCE=v3.90.0      # rebuild a past release
SINCE       ?=
RELEASE_TAG ?=
changelog:
	@cd scripts/changelog && $(GO) run . -mode=write -version=$(VERSION) -repo=../.. -since=$(SINCE) -release-tag=$(RELEASE_TAG)

## Changelog-check — fail (exit 3) when the on-disk changelogs drift
## from the regenerated output. Wire into CI. Forwards SINCE / RELEASE_TAG
## so partial-update PRs can verify only their slice.
changelog-check:
	@cd scripts/changelog && $(GO) run . -mode=check -version=$(VERSION) -repo=../.. -since=$(SINCE) -release-tag=$(RELEASE_TAG)

## Goldens-regen — regenerate golden fixtures for a specific test pattern.
## REQUIRES RUN=<pattern>. Delegates to `gitmap regoldens`, which is the
## ONLY sanctioned entry point that may unlock the golden-update gate
## (see spec/05-coding-guidelines/21-golden-fixture-regeneration.md §6).
## PKG defaults to ./... but should be narrowed for speed.
## Usage:
##   make goldens-regen RUN=TestStartupListJSONContract
##   make goldens-regen RUN=FindNextJSONContract PKG=./cmd/...
PKG ?= ./...
goldens-regen:
	@if [ -z "$(RUN)" ]; then \
		echo "ERROR: RUN=<test pattern> is required (e.g. make goldens-regen RUN=TestFooContract)"; \
		exit 2; \
	fi
	@echo "▸ Regenerating goldens via gitmap regoldens: pattern=$(RUN) pkg=$(PKG)"
	@cd $(MODULE) && $(GO) run . regoldens --run '$(RUN)' --pkg '$(PKG)'

## Goldens-verify — re-run the same pattern WITHOUT the gating env vars to
## confirm regenerated fixtures pass cleanly. This is the mandatory second
## pass: if it fails, the writers are non-deterministic or drift remains.
## Usage:
##   make goldens-verify RUN=TestStartupListJSONContract
goldens-verify:
	@if [ -z "$(RUN)" ]; then \
		echo "ERROR: RUN=<test pattern> is required (e.g. make goldens-verify RUN=TestFooContract)"; \
		exit 2; \
	fi
	@echo "▸ Verifying goldens (no env gates): pattern=$(RUN) pkg=$(PKG)"
	@cd $(MODULE) && unset GITMAP_UPDATE_GOLDEN GITMAP_ALLOW_GOLDEN_UPDATE && \
		$(GO) test $(PKG) -run '$(RUN)' -count=1 -v

## Fixtures-bump — re-run a test with GITMAP_FIXTURE_AUTOBUMP=1 so any
## test using fixtureversion.MustValidateBodyWithAutobump auto-rewrites
## its stale `// fixture-stamp:` marker in source. Pure no-op for tests
## that are not stale or do not opt in. Always followed up with a
## clean re-run via fixtures-bump-verify so a non-deterministic bump
## cannot land silently.
## REQUIRES RUN=<pattern>. PKG defaults to ./... — narrow it for speed.
## Fails with exit code 3 if RUN matched zero tests — a typo in the
## pattern would otherwise look like a successful no-op bump.
## Usage:
##   make fixtures-bump RUN=TestFixRepoRewriteV9ToV12Fixture PKG=./cmd/...
fixtures-bump:
	@if [ -z "$(RUN)" ]; then \
		echo "ERROR: RUN=<test pattern> is required (e.g. make fixtures-bump RUN=TestFooFixture)"; \
		exit 2; \
	fi
	@echo "▸ Auto-bumping fixture stamps: pattern=$(RUN) pkg=$(PKG)"
	@LOG=$$(mktemp); \
		cd $(MODULE) && GITMAP_FIXTURE_AUTOBUMP=1 \
		$(GO) test $(PKG) -run '$(RUN)' -count=1 -v 2>&1 | tee $$LOG; \
		ran=$$(grep -c '^=== RUN  ' $$LOG || true); \
		rm -f $$LOG; \
		if [ "$$ran" -eq 0 ]; then \
			echo "ERROR: RUN='$(RUN)' matched 0 tests in PKG='$(PKG)' — check the pattern (regex) and package path."; \
			exit 3; \
		fi; \
		echo "▸ Autobump pass executed $$ran test(s)"
	@$(MAKE) fixtures-bump-verify RUN='$(RUN)' PKG='$(PKG)'

## Fixtures-bump-verify — re-run the same pattern WITHOUT the autobump
## env gate to confirm the rewritten fixtures pass cleanly. Mandatory
## second pass (mirrors goldens-verify): if it fails, the rewriter or
## the test logic still disagrees and a human must intervene.
## Also fails with exit code 3 if RUN matched zero tests, so a
## stale or mistyped pattern cannot masquerade as a clean verify.
## Usage:
##   make fixtures-bump-verify RUN=TestFooFixture
fixtures-bump-verify:
	@if [ -z "$(RUN)" ]; then \
		echo "ERROR: RUN=<test pattern> is required"; \
		exit 2; \
	fi
	@echo "▸ Verifying bumped fixtures (no autobump gate): pattern=$(RUN) pkg=$(PKG)"
	@LOG=$$(mktemp); \
		cd $(MODULE) && unset GITMAP_FIXTURE_AUTOBUMP && \
		$(GO) test $(PKG) -run '$(RUN)' -count=1 -v 2>&1 | tee $$LOG; \
		status=$${PIPESTATUS[0]}; \
		ran=$$(grep -c '^=== RUN  ' $$LOG || true); \
		rm -f $$LOG; \
		if [ "$$status" -ne 0 ]; then exit $$status; fi; \
		if [ "$$ran" -eq 0 ]; then \
			echo "ERROR: RUN='$(RUN)' matched 0 tests in PKG='$(PKG)' during verify pass."; \
			exit 3; \
		fi; \
		echo "▸ Verify pass executed $$ran test(s)"
