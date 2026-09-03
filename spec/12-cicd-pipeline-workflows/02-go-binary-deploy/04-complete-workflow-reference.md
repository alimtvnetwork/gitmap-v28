# Go Binary — Complete Workflow Reference

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

Provide a single, end-to-end, copy-paste-ready GitHub Actions workflow that assembles all patterns from the modular spec files into one cohesive release pipeline. This file exists so an AI or engineer can generate a working `.github/workflows/release.yml` without manually stitching fragments from 7+ spec files.

**This is the integration document** — it shows how every piece connects.

---

## Complete Release Workflow

```yaml
name: Release

on:
  push:
    branches:
      - "release/**"
    tags:
      - "v*"

concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false

permissions:
  contents: write

env:
  FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
  BINARY_NAME: "<binary>"
  UPDATER_NAME: "<binary>-updater"
  MAIN_MODULE_DIR: "<binary>"
  UPDATER_MODULE_DIR: "<binary>-updater"

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      # ──────────────────────────────────────────────
      # Stage 1: Checkout
      # ──────────────────────────────────────────────
      - name: Checkout repository
        uses: actions/checkout@v6
        with:
          fetch-depth: 0

      # ──────────────────────────────────────────────
      # Stage 2: Setup Go toolchain
      # ──────────────────────────────────────────────
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version-file: ${{ env.MAIN_MODULE_DIR }}/go.mod
          cache-dependency-path: ${{ env.MAIN_MODULE_DIR }}/go.sum

      # ──────────────────────────────────────────────
      # Stage 3: Resolve version from Git ref
      # ──────────────────────────────────────────────
      - name: Resolve version
        id: version
        run: |
          if [[ "$GITHUB_REF" == refs/tags/* ]]; then
            VERSION="${GITHUB_REF_NAME}"
          elif [[ "$GITHUB_REF" == refs/heads/release/* ]]; then
            VERSION="${GITHUB_REF_NAME#release/}"
          else
            echo "::error::Unexpected ref: $GITHUB_REF"
            exit 1
          fi
          echo "version=$VERSION" >> "$GITHUB_OUTPUT"
          echo "clean_version=${VERSION#v}" >> "$GITHUB_OUTPUT"
          echo "Resolved version: $VERSION"

      # ──────────────────────────────────────────────
      # Stage 4a: Generate Windows resources (icon + version info)
      # ──────────────────────────────────────────────
      - name: Generate Windows resources (main binary)
        working-directory: ${{ env.MAIN_MODULE_DIR }}
        run: |
          go install github.com/tc-hib/go-winres@v0.3.3
          CLEAN_VERSION="${{ steps.version.outputs.clean_version }}"
          go-winres make --product-version "$CLEAN_VERSION" --file-version "$CLEAN_VERSION"
          echo "Generated .syso files:"
          ls -la rsrc_windows_*.syso

      - name: Generate Windows resources (updater)
        working-directory: ${{ env.UPDATER_MODULE_DIR }}
        run: |
          CLEAN_VERSION="${{ steps.version.outputs.clean_version }}"
          go-winres make --product-version "$CLEAN_VERSION" --file-version "$CLEAN_VERSION"
          ls -la rsrc_windows_*.syso

      # ──────────────────────────────────────────────
      # Stage 4b: Cross-compile ALL binaries (6 targets × 2 binaries)
      # ──────────────────────────────────────────────
      - name: Build main binary (all targets)
        working-directory: ${{ env.MAIN_MODULE_DIR }}
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          COMMIT="${GITHUB_SHA:0:10}"
          BUILD_DATE="$(date -u '+%Y-%m-%d')"
          LDFLAGS="-s -w \
            -X '<module>/version.Version=$VERSION' \
            -X '<module>/version.Commit=$COMMIT' \
            -X '<module>/version.BuildDate=$BUILD_DATE' \
            -X '<module>/version.RepoPath=<owner>/<repo>'"

          mkdir -p dist

          TARGETS="windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"
          for target in $TARGETS; do
            os="${target%/*}"
            arch="${target#*/}"
            ext=""
            [[ "$os" == "windows" ]] && ext=".exe"

            output="dist/${BINARY_NAME}-${os}-${arch}${ext}"
            echo "Building: $output"
            CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
              go build -ldflags "$LDFLAGS" -o "$output" .
            echo "  Size: $(stat -c%s "$output" 2>/dev/null || stat -f%z "$output") bytes"
          done

      - name: Build updater binary (all targets)
        working-directory: ${{ env.UPDATER_MODULE_DIR }}
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          COMMIT="${GITHUB_SHA:0:10}"
          LDFLAGS="-s -w \
            -X '<module>/version.Version=$VERSION' \
            -X '<module>/version.Commit=$COMMIT'"

          TARGETS="windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"
          for target in $TARGETS; do
            os="${target%/*}"
            arch="${target#*/}"
            ext=""
            [[ "$os" == "windows" ]] && ext=".exe"

            # NOTE: Output goes to the MAIN module's dist/ folder
            output="../${MAIN_MODULE_DIR}/dist/${UPDATER_NAME}-${os}-${arch}${ext}"
            echo "Building: $output"
            CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
              go build -ldflags "$LDFLAGS" -o "$output" .
          done

      # ──────────────────────────────────────────────
      # Stage 4c: Build docs-site (if applicable)
      # ──────────────────────────────────────────────
      - name: Setup Node.js (for docs-site)
        uses: actions/setup-node@v4
        with:
          node-version: "20"

      - name: Build docs-site
        working-directory: docs-site
        run: |
          npm ci
          npm run build

      - name: Package docs-site
        run: |
          cd docs-site/dist
          zip -r "../../${MAIN_MODULE_DIR}/dist/docs-site.zip" .

      # ──────────────────────────────────────────────
      # Stage 5: Code signing (conditional)
      # See: 02-spec/12-cicd-pipeline-workflows/10-code-signing.md
      # ──────────────────────────────────────────────
      - name: Sign Windows binaries
        if: vars.SIGNPATH_SIGNING_ENABLED == 'true'
        uses: signpath/github-action-submit-signing-request@v1
        with:
          api-token: ${{ secrets.SIGNPATH_API_TOKEN }}
          organization-id: ${{ secrets.SIGNPATH_ORGANIZATION_ID }}
          project-slug: ${{ secrets.SIGNPATH_PROJECT_SLUG }}
          signing-policy-slug: ${{ secrets.SIGNPATH_SIGNING_POLICY_SLUG }}
          artifact-configuration-slug: "exe"
          input-artifact-path: "${{ env.MAIN_MODULE_DIR }}/dist"
          output-artifact-path: "${{ env.MAIN_MODULE_DIR }}/dist"
          wait-for-completion: true
          wait-for-completion-timeout-in-seconds: 600
          parameters: |
            include: "*.exe"

      # ──────────────────────────────────────────────
      # Stage 6: Compress binaries into platform archives
      # ──────────────────────────────────────────────
      - name: Compress binaries
        working-directory: ${{ env.MAIN_MODULE_DIR }}/dist
        run: |
          for f in *-*; do
            # Skip non-binary files
            [[ "$f" == *.zip || "$f" == *.tar.gz || "$f" == *.txt || "$f" == *.ps1 || "$f" == *.sh || "$f" == *.md ]] && continue
            [ -f "$f" ] || continue

            if [[ "$f" == *.exe ]]; then
              name="${f%.exe}"
              zip "${name}.zip" "$f"
            else
              tar czf "${f}.tar.gz" "$f"
            fi
            rm "$f"
          done

      # ──────────────────────────────────────────────
      # Stage 7: Generate checksums
      # ──────────────────────────────────────────────
      - name: Generate checksums
        working-directory: ${{ env.MAIN_MODULE_DIR }}/dist
        run: |
          sha256sum *.zip *.tar.gz docs-site.zip 2>/dev/null > checksums.txt || true
          echo "Checksums:"
          cat checksums.txt

      # ──────────────────────────────────────────────
      # Stage 8: Generate version-pinned install scripts
      # See: 02-spec/12-cicd-pipeline-workflows/07-install-script-generation.md
      # ──────────────────────────────────────────────
      - name: Generate install scripts
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          REPO="${{ github.repository }}"

          cp scripts/install.ps1 ${MAIN_MODULE_DIR}/dist/install.ps1
          cp scripts/install.sh  ${MAIN_MODULE_DIR}/dist/install.sh

          sed -i "s|VERSION_PLACEHOLDER|$VERSION|g" ${MAIN_MODULE_DIR}/dist/install.ps1
          sed -i "s|REPO_PLACEHOLDER|$REPO|g"       ${MAIN_MODULE_DIR}/dist/install.ps1
          sed -i "s|VERSION_PLACEHOLDER|$VERSION|g" ${MAIN_MODULE_DIR}/dist/install.sh
          sed -i "s|REPO_PLACEHOLDER|$REPO|g"       ${MAIN_MODULE_DIR}/dist/install.sh

          chmod +x ${MAIN_MODULE_DIR}/dist/install.sh

      # ──────────────────────────────────────────────
      # Stage 9: Extract changelog
      # See: 02-spec/12-cicd-pipeline-workflows/14-release-body-and-changelog.md
      # ──────────────────────────────────────────────
      - name: Extract changelog
        id: changelog
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          CHANGELOG_ENTRY=$(awk -v ver="$VERSION" '
            /^## / {
              if (found) exit
              if (index($0, ver)) found=1
            }
            found { print }
          ' changelog.md 2>/dev/null || echo "")

          if [ -z "$CHANGELOG_ENTRY" ]; then
            CHANGELOG_ENTRY="No changelog entry found for $VERSION."
          fi

          echo "$CHANGELOG_ENTRY" > /tmp/changelog-entry.md

      # ──────────────────────────────────────────────
      # Stage 10: Assemble release body
      # ──────────────────────────────────────────────
      - name: Generate release body
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          COMMIT="${GITHUB_SHA:0:10}"
          BRANCH="${GITHUB_REF_NAME}"
          BUILD_DATE="$(date -u '+%Y-%m-%d')"
          GO_VERSION="$(go version | awk '{print $3}')"
          REPO="${{ github.repository }}"

          cat > /tmp/release-body.md << 'HEADER'
          $(cat /tmp/changelog-entry.md)

          ---

          ## Release Info

          | Field | Value |
          |-------|-------|
          HEADER

          cat >> /tmp/release-body.md << EOF
          | Version | $VERSION |
          | Commit | $COMMIT |
          | Branch | $BRANCH |
          | Build Date | $BUILD_DATE |
          | Go Version | $GO_VERSION |

          ## Checksums (SHA-256)

          \`\`\`
          $(cat ${MAIN_MODULE_DIR}/dist/checksums.txt)
          \`\`\`

          ## Quick Install

          **Windows (PowerShell)**
          \`\`\`powershell
          irm https://github.com/$REPO/releases/download/$VERSION/install.ps1 | iex
          \`\`\`

          **Linux / macOS**
          \`\`\`bash
          curl -fsSL https://github.com/$REPO/releases/download/$VERSION/install.sh | bash
          \`\`\`
          EOF

      # ──────────────────────────────────────────────
      # Stage 11: Create GitHub Release
      # ──────────────────────────────────────────────
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: ${{ steps.version.outputs.version }}
          name: "${{ env.BINARY_NAME }} ${{ steps.version.outputs.version }}"
          body_path: /tmp/release-body.md
          files: ${{ env.MAIN_MODULE_DIR }}/dist/*
          draft: false
          prerelease: ${{ contains(steps.version.outputs.version, '-') }}
          make_latest: ${{ !contains(steps.version.outputs.version, '-') }}
```

---

## Stage-by-Stage Cross-References

| Stage | Spec Reference |
|-------|---------------|
| Checkout & Setup | [Shared Conventions](../03-shared-conventions.md) |
| Version Resolution | [Shared Conventions — Version Resolution](../03-shared-conventions.md#version-resolution) |
| Windows Resources | [Binary Icon Branding](../16-binary-icon-branding.md) |
| Cross-Compilation | [Cross-Compilation](../../14-update/17-cross-compilation.md) |
| Multi-Module Build | [Release Pipeline — Multiple Binaries](../05-release-pipeline.md#multiple-binaries) |
| Docs-Site Bundling | [Docs-Site Bundling](../05-release-pipeline.md#docs-site-bundling) |
| Code Signing | [Code Signing](../10-code-signing.md) |
| Compression | [Release Assets](../../14-update/14-release-assets.md) |
| Checksums | [Checksums & Verification](../../14-update/15-checksums-verification.md) |
| Install Scripts | [Install Script Generation](../07-install-script-generation.md) |
| Changelog | [Release Body and Changelog](../14-release-body-and-changelog.md) |
| GitHub Release | [GitHub Release Standard](../04-github-release-standard.md) |

---

## Multi-Module Directory Layout

This workflow assumes the following project structure for multi-binary projects:

```
<repo-root>/
├── <binary>/                  # Main binary Go module
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── winres.json
│   ├── assets/icon.png
│   └── dist/                  # ALL build outputs go here
│       ├── <binary>-linux-amd64.tar.gz
│       ├── <binary>-updater-linux-amd64.tar.gz
│       ├── docs-site.zip
│       ├── checksums.txt
│       ├── install.ps1
│       └── install.sh
├── <binary>-updater/          # Updater binary Go module
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── winres.json
├── docs-site/                 # Documentation site (Node.js)
│   ├── package.json
│   └── dist/
├── scripts/
│   ├── install.ps1            # Template with placeholders
│   └── install.sh             # Template with placeholders
└── changelog.md
```

Key rules:

- The updater binary outputs to `../<binary>/dist/` (the main module's dist folder)
- `docs-site.zip` is built from a Node.js project and placed in the same `dist/`
- All assets converge into a **single `dist/` directory** for unified packaging and release

---

## Error Recovery

| Failure Point | Behavior | Fallback |
|---------------|----------|----------|
| Version resolution fails | Pipeline exits immediately | Fix the Git ref or branch name |
| One cross-compile target fails | Pipeline fails — no partial releases | Fix the build error and re-tag |
| Code signing times out (600s) | Signing step fails; pipeline stops | Disable signing via `SIGNPATH_SIGNING_ENABLED=false` and release unsigned |
| Changelog extraction finds nothing | Uses fallback text: "No changelog entry found for $VERSION." | Add changelog entry and re-release |
| Docs-site build fails | Pipeline fails | Fix docs build or remove docs-site step |
| GitHub Release creation fails | Pipeline fails | Check permissions and re-run |

**Policy**: Release builds are all-or-nothing. A partial release (e.g., 5 of 6 targets) is never published. If any stage fails, the entire pipeline fails and must be re-run after fixing the issue.

---

## Constraints

- This is a **single-job workflow** — all stages run sequentially on one runner
- Binaries are built **exactly once** — no stage triggers a rebuild
- The updater builds into the main module's `dist/` directory
- Release pipelines **never cancel** — `cancel-in-progress: false`
- All tools are pinned to exact versions
- `working-directory` is used instead of `cd`

---

## Cross-References

- [CI Pipeline](./02-ci-pipeline.md) — Validation pipeline that precedes releases
- [Release Pipeline](../05-release-pipeline.md) — Modular spec this workflow implements
- [Shared Conventions](../03-shared-conventions.md) — Platform, triggers, permissions

---

*Complete workflow reference — v3.2.0 — 2026-04-13*
