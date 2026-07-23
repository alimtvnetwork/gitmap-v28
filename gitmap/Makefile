.PHONY: build run test update no-pull no-deploy clean help generate cross cross-startup fmt

SHELL := /bin/bash

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

# `generate` runs every //go:generate directive in the module. Today the
# only generator is completion/internal/gencommands (writes
# completion/allcommands_generated.go), but new generators added later
# will be picked up automatically. Made a prerequisite of build/run/test
# so the generated file can never silently drift from the Cmd*
# constants — exactly the failure class the CI generate-check guards
# against. Using a phony target instead of a file target keeps it
# dependency-free and always runs (cheap: ~100ms locally).
generate: ## Regenerate completion + any other go:generate outputs
	@go generate ./...

build: generate ## Build and deploy (full pipeline)
	@bash ../run.sh

run: generate ## Build, deploy, and run with optional ARGS (make run ARGS="list --json")
	@bash ../run.sh -r $(ARGS)

lint: ## Check for duplicate type declarations across package files
	@bash scripts/check-duplicate-types.sh .

fmt: ## Auto-format every .go file under gitmap/ with gofmt
	@bash ../scripts/format-go.sh

test: generate ## Run tests with report
	@bash ../run.sh -t

update: ## Pull, build, deploy, and sync PATH binary
	@bash ../run.sh --update

no-pull: generate ## Build without git pull
	@bash ../run.sh --no-pull

no-deploy: generate ## Build without deploying
	@bash ../run.sh --no-deploy

clean: ## Remove build artifacts
	@rm -rf bin/
	@echo "Build artifacts removed."

# Cross-platform compilation matrix. Mirrors the targets in
# .github/workflows/cross-platform.yml + startup-build-tags.yml so a
# single local command catches build-tag leaks (e.g. the historic
# startup/win*.go regression where an untagged dispatcher referenced
# a Windows-only symbol with no `_other.go` stub).
#
# Uses `go build` (not `go vet`) so every referenced symbol must
# resolve on the target OS — vet skips files excluded by build tags
# and would silently miss the leak class. CGO_ENABLED=0 matches the
# release pipeline's static-link contract.
#
# Matrix: 6 targets (windows/linux/darwin × amd64/arm64), same as
# spec/01-app/42-cross-platform.md. Compile-only test pass uses
# `go test -c` so cross-OS validation never tries to execute a foreign
# test binary while still catching symbol leaks in *_test.go files.
#
# Runs sequentially (not parallel) so failures attribute cleanly to
# the offending GOOS/GOARCH pair without interleaved output.
CROSS_TARGETS := \
	windows/amd64 windows/arm64 \
	linux/amd64   linux/arm64 \
	darwin/amd64  darwin/arm64

cross: generate ## Cross-compile all 6 GOOS/GOARCH targets (build + test-compile)
	@set -e; \
	failed=""; \
	for target in $(CROSS_TARGETS); do \
		os=$${target%/*}; arch=$${target#*/}; \
		printf "\033[36m==> %s/%s\033[0m\n" "$$os" "$$arch"; \
		if ! CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build ./...; then \
			failed="$$failed $$os/$$arch(build)"; continue; \
		fi; \
		if ! CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go test -count=0 ./...; then \
			failed="$$failed $$os/$$arch(test-compile)"; continue; \
		fi; \
		printf "    \033[32mok\033[0m\n"; \
	done; \
	if [ -n "$$failed" ]; then \
		printf "\n\033[31m✗ cross-compile failed:\033[0m%s\n" "$$failed" >&2; \
		exit 1; \
	fi; \
	printf "\n\033[32m✓ all 6 GOOS/GOARCH targets compiled cleanly\033[0m\n"

cross-startup: ## Cross-compile only the startup package (fast build-tag guard)
	@set -e; \
	for os in linux darwin windows; do \
		printf "\033[36m==> startup / GOOS=%s\033[0m\n" "$$os"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=amd64 go build ./startup/...; \
		for pkg in $$(CGO_ENABLED=0 GOOS=$$os GOARCH=amd64 go list ./startup/...); do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=amd64 go test -c -o /dev/null "$$pkg"; \
		done; \
	done; \
	printf "\n\033[32m✓ startup package builds cleanly on all 3 OSes\033[0m\n"
