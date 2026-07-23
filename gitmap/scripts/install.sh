#!/usr/bin/env bash
# Re-exec under bash if invoked via sh/dash (which lack pipefail, local, etc.)
#
# IMPORTANT: When this script is piped through `sh` (e.g.
# `curl ... | sh`), the script source arrives on stdin and is consumed
# line-by-line by the POSIX shell. We cannot simply `exec bash -s` in that
# case because bash would then read only whatever bytes the parent shell
# had not yet consumed — producing "command not found" errors for every
# function defined later in the file.
#
# The reliable fix is to materialize the full script to a temp file
# (either by copying ourselves when invoked from disk, or by re-fetching
# from the canonical URL when streamed from a pipe) and then exec bash
# against that file.
if [ -z "${BASH_VERSION:-}" ]; then
    if ! command -v bash >/dev/null 2>&1; then
        printf '\033[31m  Error: bash is required but not found. Install bash first.\033[0m\n' >&2
        exit 1
    fi

    # If $0 points at a real readable file on disk, just re-exec bash on it.
    if [ -f "$0" ] && [ -r "$0" ]; then
        exec bash "$0" "$@"
    fi

    # Otherwise we were piped from stdin (curl | sh). Capture stdin to a
    # temp file first; if stdin is empty/exhausted, fall back to fetching
    # the canonical install.sh from GitHub.
    _gm_tmp="$(mktemp 2>/dev/null || echo "/tmp/gitmap-install.$$.sh")"
    # Try draining whatever is left on stdin into the temp file.
    if [ ! -t 0 ]; then
        cat > "$_gm_tmp" 2>/dev/null || true
    fi

    # If the captured file is missing the shebang or is suspiciously short,
    # the parent `sh` already consumed most of it. Re-fetch from GitHub.
    if [ ! -s "$_gm_tmp" ] || ! head -1 "$_gm_tmp" 2>/dev/null | grep -q '^#!'; then
        _gm_url="https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.sh"
        if command -v curl >/dev/null 2>&1; then
            curl -fsSL "$_gm_url" -o "$_gm_tmp" || {
                printf '\033[31m  Error: failed to re-fetch installer from %s\033[0m\n' "$_gm_url" >&2
                exit 1
            }
        elif command -v wget >/dev/null 2>&1; then
            wget -qO "$_gm_tmp" "$_gm_url" || {
                printf '\033[31m  Error: failed to re-fetch installer from %s\033[0m\n' "$_gm_url" >&2
                exit 1
            }
        else
            printf '\033[31m  Error: neither curl nor wget available to re-fetch installer.\033[0m\n' >&2
            exit 1
        fi
    fi

    chmod +x "$_gm_tmp" 2>/dev/null || true
    # Mark the temp file for cleanup after bash exits (best-effort).
    export GITMAP_INSTALL_SELF_TMP="$_gm_tmp"
    exec bash "$_gm_tmp" "$@"
fi

# If we were re-exec'd from a self-written temp file, schedule cleanup.
if [ -n "${GITMAP_INSTALL_SELF_TMP:-}" ]; then
    trap 'rm -f "$GITMAP_INSTALL_SELF_TMP" 2>/dev/null || true' EXIT
fi
# ─────────────────────────────────────────────────────────────────────
# gitmap installer for Linux and macOS
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.sh | bash
#
# Options:
#   --version <tag>    Install a specific version (e.g. v2.55.0). Default: latest.
#   --dir <path>       Target directory. Default: ~/.local/bin
#   --arch <arch>      Force architecture (amd64, arm64). Default: auto-detect.
#   --no-path          Skip adding install directory to PATH.
#
# Examples:
#   curl -fsSL .../install.sh | bash
#   curl -fsSL .../install.sh | bash -s -- --version v2.55.0
#   ./install.sh --dir /opt/gitmap --arch arm64
# ─────────────────────────────────────────────────────────────────────

set -euo pipefail

REPO="alimtvnetwork/gitmap-v28"
BINARY_NAME="gitmap"
BINARY_ALIAS="gm"
TMP_DIR=""
APP_DIR=""
PATH_SHELL=""
PATH_TARGET=""
PATH_LINE=""
PATH_STATUS=""
PATH_RELOAD=""
# PATH_RELOAD_ALT holds a secondary reload command for the *other* shell
# when both a POSIX profile and the pwsh profile were written (typical
# under --dual-shell or when pwsh is detected alongside zsh/bash). The
# post-install block renders it as an "or in <shell>:" hint so users in
# either shell see the syntactically correct reload command — never a
# `source ~/.zshrc` shown to someone sitting in pwsh, and vice versa.
PATH_RELOAD_ALT=""
PATH_RELOAD_ALT_SHELL=""
# Per-profile audit trail used by --show-path. Populated by add_to_path
# and consumed by print_install_summary; safe to read even when empty.
PATH_PROFILES_WRITTEN=""    # union of added + updated; back-compat
PATH_PROFILES_SKIPPED=""    # mirrors PATH_PROFILES_UNCHANGED for back-compat
PATH_PROFILES_ADDED=""      # profile files where the snippet was newly appended
PATH_PROFILES_UPDATED=""    # profile files whose snippet body changed
PATH_PROFILES_UNCHANGED=""  # profile files where the snippet was already byte-identical
PATH_PWSH_DETECTED="no"

cleanup() {
    if [ -n "${TMP_DIR}" ] && [ -d "${TMP_DIR}" ]; then
        rm -rf "${TMP_DIR}"
    fi
}
trap cleanup EXIT

# ── Logging helpers ─────────────────────────────────────────────────

step()  { printf '  \033[36m%s\033[0m\n' "$*" >&2; }
ok()    { printf '  \033[32m%s\033[0m\n' "$*" >&2; }
err()   { printf '  \033[31m%s\033[0m\n' "$*" >&2; }
warn()  { printf '  \033[33m%s\033[0m\n' "$*" >&2; }

# ── Deploy manifest (single source of truth) ───────────────────────
# Mirrors run.sh's load_deploy_manifest and run.ps1's Get-DeployManifest.
# Fetches gitmap/constants/deploy-manifest.json from the install repo so
# APP_SUBDIR / LEGACY_APP_SUBDIRS aren't hardcoded across scripts and Go.
# Renaming the deploy folder ONLY requires editing that JSON file.
APP_SUBDIR="gitmap-cli"
LEGACY_APP_SUBDIRS=("gitmap")
load_deploy_manifest() {
    local manifest_url="https://raw.githubusercontent.com/${REPO}/main/gitmap/constants/deploy-manifest.json"
    local manifest
    manifest=$(curl -fsSL --max-time 5 "$manifest_url" 2>/dev/null || true)
    if [ -z "$manifest" ]; then
        warn "deploy-manifest.json not reachable - using defaults ($APP_SUBDIR)"
        return
    fi
    local app
    app=$(printf '%s' "$manifest" | grep -E '"appSubdir"' | head -n1 | sed -E 's/.*"appSubdir"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')
    if [ -n "$app" ]; then
        APP_SUBDIR="$app"
    fi
    local legacy_block legacy_csv
    legacy_block=$(printf '%s' "$manifest" | awk '/"legacyAppSubdirs"/,/]/')
    if [ -n "$legacy_block" ]; then
        legacy_csv=$(printf '%s' "$legacy_block" | grep -oE '"[^"]+"' | sed -E 's/^"(.*)"$/\1/' | grep -v '^legacyAppSubdirs$' || true)
        if [ -n "$legacy_csv" ]; then
            LEGACY_APP_SUBDIRS=()
            while IFS= read -r line; do
                [ -n "$line" ] && LEGACY_APP_SUBDIRS+=("$line")
            done <<< "$legacy_csv"
        fi
    fi
}

# ── Versioned repo discovery ────────────────────────────────────────
# spec/01-app/95-installer-script-find-latest-repo.md

# Parses "<owner>/<stem>-v<N>". Sets SUFFIX_OWNER, SUFFIX_STEM, SUFFIX_N.
parse_repo_suffix() {
    local repo="$1"
    if [[ "$repo" =~ ^([^/]+)/(.+)-v([0-9]+)$ ]]; then
        SUFFIX_OWNER="${BASH_REMATCH[1]}"
        SUFFIX_STEM="${BASH_REMATCH[2]}"
        SUFFIX_N="${BASH_REMATCH[3]}"
        return 0
    fi
    return 1
}

repo_exists() {
    curl -sfI --max-time 5 "$1" >/dev/null 2>&1
}

# Echoes the effective "<owner>/<stem>-v<M>" (or original repo when none higher).
resolve_effective_repo() {
    local repo="$1" ceiling="$2"
    if ! parse_repo_suffix "$repo"; then
        printf '  [discovery] no -v<N> suffix on '"'"'%s'"'"'; installing baseline as-is\n' "$repo" >&2
        echo "$repo"
        return 0
    fi

    local owner="$SUFFIX_OWNER" stem="$SUFFIX_STEM" baseline="$SUFFIX_N"
    local effective="$baseline" m url

    printf '  [discovery] baseline: %s/%s-v%s\n' "$owner" "$stem" "$baseline" >&2
    printf '  [discovery] probe ceiling: %s\n' "$ceiling" >&2

    for (( m = baseline + 1; m <= ceiling; m++ )); do
        url="https://github.com/${owner}/${stem}-v${m}"
        if repo_exists "$url"; then
            printf '  [discovery] HEAD %s ... HIT\n' "$url" >&2
            effective=$m
        else
            printf '  [discovery] HEAD %s ... MISS (fail-fast)\n' "$url" >&2
            break
        fi
    done

    if [ "$effective" = "$baseline" ]; then
        printf '  [discovery] no higher version found; using baseline -v%s\n' "$baseline" >&2
        echo "$repo"
    else
        printf '  [discovery] effective: %s/%s-v%s (was -v%s)\n' "$owner" "$stem" "$effective" "$baseline" >&2
        echo "${owner}/${stem}-v${effective}"
    fi
}

# Re-exec the full installer from the effective repo, passing through flags.
invoke_delegated_full_installer() {
    local effective_repo="$1"
    shift
    local delegated_url="https://raw.githubusercontent.com/${effective_repo}/main/gitmap/scripts/install.sh"
    printf '  [discovery] delegating to %s\n' "$delegated_url" >&2

    export INSTALLER_DELEGATED=1

    local script
    if ! script="$(curl -fsSL --max-time 15 "$delegated_url")"; then
        printf '  [discovery] [WARN] could not fetch delegated installer; falling back to baseline\n' >&2
        unset INSTALLER_DELEGATED
        return 1
    fi

    bash -c "$script" _ "$@"
    exit $?
}

# ── Detect OS ───────────────────────────────────────────────────────

detect_os() {
    local uname_out
    uname_out="$(uname -s)"
    case "${uname_out}" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*)
            err "Windows detected. Use the PowerShell installer instead:"
            err "  irm https://raw.githubusercontent.com/${REPO}/main/gitmap/scripts/install.ps1 | iex"
            exit 1
            ;;
        *)
            err "Unsupported OS: ${uname_out}"
            exit 1
            ;;
    esac
}

# ── Detect architecture ────────────────────────────────────────────

detect_arch() {
    local arch_flag="$1"
    if [ -n "${arch_flag}" ]; then
        echo "${arch_flag}"
        return
    fi

    local machine
    machine="$(uname -m)"
    case "${machine}" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)
            err "Unsupported architecture: ${machine}"
            exit 1
            ;;
    esac
}

# ── Resolve version (latest or pinned) ─────────────────────────────

resolve_version() {
    local version="$1"
    if [ -n "${version}" ]; then
        echo "${version}"
        return
    fi

    step "Fetching latest release..."
    local url="https://api.github.com/repos/${REPO}/releases/latest"
    local tag

    if command -v curl >/dev/null 2>&1; then
        tag="$(curl -fsSL "${url}" | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')"
    elif command -v wget >/dev/null 2>&1; then
        tag="$(wget -qO- "${url}" | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')"
    else
        err "Neither curl nor wget found. Cannot fetch latest release."
        exit 1
    fi

    if [ -z "${tag}" ]; then
        err "Failed to determine latest version."
        exit 1
    fi

    echo "${tag}"
}

# ── Download helper ────────────────────────────────────────────────

download() {
    local url="$1" dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "${dest}" "${url}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${dest}" "${url}"
    else
        err "Neither curl nor wget found."
        exit 1
    fi
}

# ── Strict-tag failure (spec/07-generic-release/09 §3) ─────────────
# Print the canonical no-fallback message and exit 1. Called from
# download_asset whenever VERSION was supplied explicitly and the
# requested release asset cannot be downloaded or verified.
strict_fail() {
    local detail="$1"
    err ""
    err "Error: requested release ${VERSION} not found in ${REPO};"
    err "       refusing to fall back per strict-tag contract."
    err "       See spec/07-generic-release/09-generic-install-script-behavior.md §3."
    if [ -n "${detail}" ]; then
        err "       Detail: ${detail}"
    fi
    exit 1
}

# ── Pre-flight asset existence ─────────────────────────────────────

# preflight_asset_exists: HEAD-probe the asset URL. Returns 0 on
# HTTP 200/302, non-zero otherwise. Used to fail fast with a clear
# error before any download/extract work.
preflight_asset_exists() {
    local url="$1"
    local code
    if command -v curl >/dev/null 2>&1; then
        code="$(curl -fsSLI -o /dev/null -w '%{http_code}' \
            --max-time 10 "${url}" 2>/dev/null || echo "000")"
    elif command -v wget >/dev/null 2>&1; then
        code="$(wget --spider --timeout=10 -S "${url}" 2>&1 \
            | awk '/HTTP\// {code=$2} END {print code+0}')"
    else
        # No HEAD tool; let the caller try the real download.
        return 0
    fi
    case "${code}" in
        200|302|301) return 0 ;;
        *) return 1 ;;
    esac
}

# emit_missing_asset_error: print a clearly-formatted error block
# explaining which asset pattern was expected, the full URL probed,
# and the release page to inspect.
emit_missing_asset_error() {
    local version="$1" os="$2" arch="$3" asset_name="$4" asset_url="$5"
    local release_page="https://github.com/${REPO}/releases/tag/${version}"
    err ""
    err "================================================================"
    err " RELEASE ASSET NOT FOUND"
    err "================================================================"
    err " Expected pattern: ${BINARY_NAME}-${version}-${os}-${arch}.tar.gz"
    err " Resolved name:    ${asset_name}"
    err " Probed URL:       ${asset_url}"
    err " Release page:     ${release_page}"
    err ""
    err " The release tag '${version}' exists but does not publish the"
    err " expected platform asset for ${os}/${arch}. This usually means:"
    err "   • the cross-compile job for this platform failed in CI,"
    err "   • the release was cut from a branch missing this target, or"
    err "   • you pinned a pre-release/draft tag with partial assets."
    err ""
    err " Inspect the release page above to see which assets ARE present,"
    err " then either pin a different --version or wait for the release"
    err " pipeline to re-publish."
    err "================================================================"
    err ""
}

# ── Download and verify asset ──────────────────────────────────────

download_asset() {
    local version="$1" os="$2" arch="$3"
    local asset_name="${BINARY_NAME}-${version}-${os}-${arch}.tar.gz"
    local base_url="https://github.com/${REPO}/releases/download/${version}"
    local asset_url="${base_url}/${asset_name}"
    local checksum_url="${base_url}/checksums.txt"

    # Strict mode: VERSION was supplied explicitly via --version.
    # On any failure below we MUST exit 1, not fall through to a .zip
    # probe or anything else that could mask a missing-tag situation.
    local strict=0
    if [ -n "${VERSION}" ]; then
        strict=1
        printf '  [strict] download: %s\n' "${asset_url}" >&2
    fi

    # TMP_DIR is set by the caller (main).

    local archive_path="${TMP_DIR}/${asset_name}"
    local checksum_path="${TMP_DIR}/checksums.txt"

    # Pre-flight asset-existence check. Emits a clearly-formatted
    # error block when the expected release-asset pattern is missing,
    # so users see WHAT was expected, WHERE we looked, and the
    # release page to inspect — not just "Download failed".
    if ! preflight_asset_exists "${asset_url}"; then
        emit_missing_asset_error "${version}" "${os}" "${arch}" \
            "${asset_name}" "${asset_url}"
        if [ "${strict}" = "1" ]; then
            strict_fail "expected asset ${asset_name} not found at ${asset_url}"
        fi
        exit 1
    fi

    step "Downloading ${asset_name} (${version})..."
    if ! download "${asset_url}" "${archive_path}"; then
        if [ "${strict}" = "1" ]; then
            strict_fail "asset ${asset_name} not downloadable from ${base_url}"
        fi
        err "Download failed: ${asset_url}"
        exit 1
    fi
    if ! download "${checksum_url}" "${checksum_path}"; then
        if [ "${strict}" = "1" ]; then
            strict_fail "checksums.txt missing at ${checksum_url}"
        fi
        err "Download failed: ${checksum_url}"
        exit 1
    fi

    # Verify checksum
    step "Verifying checksum..."
    local expected_line
    expected_line="$(grep "${asset_name}" "${checksum_path}" || true)"
    if [ -z "${expected_line}" ]; then
        if [ "${strict}" = "1" ]; then
            # Strict mode: do NOT probe an alternate asset name. The
            # tag's own checksums.txt is the source of truth.
            strict_fail "asset ${asset_name} not listed in checksums.txt for ${version}"
        fi

        # Try .zip variant (some releases may only have zip)
        asset_name="${BINARY_NAME}-${version}-${os}-${arch}.zip"
        asset_url="${base_url}/${asset_name}"
        archive_path="${TMP_DIR}/${asset_name}"

        step "Trying .zip variant..."
        download "${asset_url}" "${archive_path}"
        expected_line="$(grep "${asset_name}" "${checksum_path}" || true)"

        if [ -z "${expected_line}" ]; then
            err "Asset not found in checksums.txt"
            err "Tried: ${BINARY_NAME}-${version}-${os}-${arch}.tar.gz"
            err "Tried: ${asset_name}"
            exit 1
        fi
    fi

    local expected_hash
    expected_hash="$(echo "${expected_line}" | awk '{print $1}')"

    local actual_hash
    if command -v sha256sum >/dev/null 2>&1; then
        actual_hash="$(sha256sum "${archive_path}" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
        actual_hash="$(shasum -a 256 "${archive_path}" | awk '{print $1}')"
    else
        err "No SHA256 tool found (sha256sum or shasum required)."
        exit 1
    fi

    if [ "${actual_hash}" != "${expected_hash}" ]; then
        if [ "${strict}" = "1" ]; then
            strict_fail "checksum mismatch for ${asset_name} (expected ${expected_hash}, got ${actual_hash})"
        fi
        err "Checksum mismatch!"
        err "  Expected: ${expected_hash}"
        err "  Got:      ${actual_hash}"
        exit 1
    fi

    ok "Checksum verified."
    echo "${archive_path}"
}

# ── Layout repair + pre-deploy cleanup (DFD-3, DFD-6) ──────────────
# Migrates two legacy layouts into the canonical <dir>/gitmap-cli/:
#   1) Unwrapped (pre-DFD): <dir>/gitmap binary at the top level.
#   2) v3.6.0..v3.13.10 wrapped: <dir>/gitmap/ folder (renamed to
#      gitmap-cli on Unix in v3.13.11 to match the Windows installer
#      which did the same rename in v3.6.0).
# Idempotent — re-running on a correct gitmap-cli/ layout is a no-op.

repair_layout() {
    local target="$1"
    local app_dir="$target/$APP_SUBDIR"
    local legacy_binary="$target/${BINARY_NAME}"
    local wrapped_binary="$app_dir/${BINARY_NAME}"

    # --- Migration 2: any legacy app folder -> $APP_SUBDIR ----
    # Distinguish folder vs file at $target/<legacy>. A directory means the
    # old wrapped layout; a file means the very-old unwrapped binary.
    local legacy
    for legacy in "${LEGACY_APP_SUBDIRS[@]}"; do
        local legacy_app_dir="$target/$legacy"
        [ "$legacy_app_dir" = "$app_dir" ] && continue
        if [ -d "$legacy_app_dir" ]; then
            if [ -d "$app_dir" ]; then
                warn "Layout: both $legacy/ and ${APP_SUBDIR}/ exist at ${target} — leaving legacy $legacy/ for manual review"
            else
                mv "$legacy_app_dir" "$app_dir" 2>/dev/null && \
                    step "Layout: migrated legacy $legacy/ -> ${APP_SUBDIR}/ at ${target}"
            fi
        fi
    done

    # --- Migration 1: legacy unwrapped binary -> $APP_SUBDIR/ --
    if [ -f "$legacy_binary" ] && [ ! -d "$legacy_binary" ]; then
        if [ -f "$wrapped_binary" ]; then
            rm -f "$legacy_binary" 2>/dev/null && \
                step "Layout: removed leftover legacy binary $legacy_binary"
        else
            step "Layout: migrating legacy unwrapped install -> ${app_dir}"
            mkdir -p "$app_dir"
            local name src dst
            for name in "${BINARY_NAME}" data CHANGELOG.md docs docs-site; do
                src="$target/$name"
                dst="$app_dir/$name"
                [ ! -e "$src" ] && continue
                [ -e "$dst" ] && continue
                mv "$src" "$dst" 2>/dev/null && \
                    step "  moved $name -> ${APP_SUBDIR}/$name"
            done
        fi
    else
        step "Layout: OK"
    fi
}

cleanup_prior_artifacts() {
    local target="$1" app_dir="$2"
    local stem="${BINARY_NAME}"
    local removed=0
    local dir pat f

    for dir in "$target" "$app_dir"; do
        [ ! -d "$dir" ] && continue
        for pat in "*.old" "${stem}-update-*" "updater-tmp-*"; do
            for f in "$dir"/$pat; do
                [ ! -e "$f" ] && continue
                rm -rf "$f" 2>/dev/null && {
                    step "[cleanup] removed $f"
                    removed=$((removed + 1))
                }
            done
        done
    done

    local tmp_root="${TMPDIR:-/tmp}"
    if [ -d "$tmp_root" ]; then
        for f in "$tmp_root/${stem}-update-"*; do
            [ ! -e "$f" ] && continue
            rm -rf "$f" 2>/dev/null && {
                step "[cleanup] removed temp $f"
                removed=$((removed + 1))
            }
        done
    fi

    if [ -d "$target" ]; then
        for f in "$target"/*.gitmap-tmp-*; do
            [ ! -d "$f" ] && continue
            rm -rf "$f" 2>/dev/null && {
                step "[cleanup] removed swap dir $f"
                removed=$((removed + 1))
            }
        done
    fi

    if [ "$removed" -gt 0 ]; then
        ok "[cleanup] removed $removed artifact(s)"
    else
        step "[cleanup] nothing to clean"
    fi
}

# ── Extract and install binary ─────────────────────────────────────

install_binary() {
    local archive_path="$1" install_dir="$2" os="$3" arch="$4" version="$5"

    # DFD-1/DFD-3: nested layout. install_dir is the deploy ROOT (e.g.
    # ~/.local/bin); the actual app folder is ${install_dir}/gitmap-cli.
    # The folder name was renamed from ${BINARY_NAME} ("gitmap") to
    # "gitmap-cli" in v3.13.11 for cross-platform parity with run.ps1.
    repair_layout "${install_dir}"
    local app_dir="${install_dir}/${APP_SUBDIR}"
    cleanup_prior_artifacts "${install_dir}" "${app_dir}"

    step "Installing to ${app_dir}..."
    mkdir -p "${app_dir}"

    local extract_dir="${TMP_DIR}/extract"
    mkdir -p "${extract_dir}"

    case "${archive_path}" in
        *.tar.gz|*.tgz)
            tar -xzf "${archive_path}" -C "${extract_dir}"
            ;;
        *.zip)
            if command -v unzip >/dev/null 2>&1; then
                unzip -qo "${archive_path}" -d "${extract_dir}"
            else
                err "unzip not found. Cannot extract .zip archive."
                exit 1
            fi
            ;;
        *)
            err "Unknown archive format: ${archive_path}"
            exit 1
            ;;
    esac

    local binary_path=""
    local candidate

    candidate="$(find "${extract_dir}" -type f -name "${BINARY_NAME}" | head -1)"
    [ -n "${candidate}" ] && binary_path="${candidate}"

    if [ -z "${binary_path}" ]; then
        candidate="$(find "${extract_dir}" -type f -name "${BINARY_NAME}-${os}-${arch}" | head -1)"
        [ -n "${candidate}" ] && binary_path="${candidate}"
    fi

    if [ -z "${binary_path}" ]; then
        candidate="$(find "${extract_dir}" -type f -regex ".*/${BINARY_NAME}-v[0-9][0-9.]*-${os}-${arch}" | head -1)"
        [ -n "${candidate}" ] && binary_path="${candidate}"
    fi

    if [ -z "${binary_path}" ]; then
        candidate="$(find "${extract_dir}" -type f -executable | head -1)"
        [ -n "${candidate}" ] && binary_path="${candidate}"
    fi

    if [ -z "${binary_path}" ]; then
        err "Archive did not contain a recognizable binary."
        find "${extract_dir}" -type f | while read -r f; do err "  ${f}"; done
        exit 1
    fi

    local target_path="${app_dir}/${BINARY_NAME}"

    if [ -f "${target_path}" ]; then
        mv -f "${target_path}" "${target_path}.old" 2>/dev/null || true
    fi

    mv -f "${binary_path}" "${target_path}"
    chmod +x "${target_path}"

    rm -f "${target_path}.old" 2>/dev/null || true

    if [ ! -f "${target_path}" ]; then
        err "Install failed: ${BINARY_NAME} was not written to ${app_dir}"
        exit 1
    fi

    ok "Installed ${BINARY_NAME} to ${app_dir}"

    # Install short alias `gm` alongside `gitmap` so both commands resolve to
    # the same binary on PATH. Prefer a symlink (atomic, single source of
    # truth); fall back to a copy when the filesystem rejects symlinks
    # (e.g. some FAT/exFAT mounts). Non-fatal: the primary install already
    # succeeded by this point.
    local alias_path="${app_dir}/${BINARY_ALIAS}"
    rm -f "${alias_path}" 2>/dev/null || true
    if ln -s "${BINARY_NAME}" "${alias_path}" 2>/dev/null; then
        ok "Installed alias ${BINARY_ALIAS} -> ${BINARY_NAME} (symlink)"
    elif cp -f "${target_path}" "${alias_path}" 2>/dev/null; then
        chmod +x "${alias_path}" 2>/dev/null || true
        ok "Installed alias ${BINARY_ALIAS} -> ${BINARY_NAME} (copy fallback)"
    else
        err "  ⚠ Could not create '${BINARY_ALIAS}' alias next to ${BINARY_NAME}; only the full '${BINARY_NAME}' command will be available."
    fi

    # Echo the app dir so main() can use it for PATH + summary.
    APP_DIR="${app_dir}"
}

# ── Download and extract docs-site.zip release asset ───────────────
# Required for `gitmap help-dashboard` (hd). Best-effort: skip silently
# if the release does not bundle docs-site.zip (older versions).
install_docs_site() {
    local version="$1" install_dir="$2"
    local asset_name="docs-site.zip"
    local asset_url="https://github.com/${REPO}/releases/download/${version}/${asset_name}"
    local tmp_zip="${TMP_DIR}/${asset_name}"

    step "Downloading docs-site.zip (${version})..."

    if ! download "${asset_url}" "${tmp_zip}" 2>/dev/null; then
        step "  docs-site.zip not available for ${version} - skipping (gitmap hd will fall back to hosted docs)"
        rm -f "${tmp_zip}" 2>/dev/null || true
        return 0
    fi

    # Remove any existing docs-site/ before extracting fresh.
    rm -rf "${install_dir}/docs-site" 2>/dev/null || true

    if ! command -v unzip >/dev/null 2>&1; then
        err "unzip not found - cannot extract docs-site.zip (install unzip and re-run)"
        rm -f "${tmp_zip}" 2>/dev/null || true
        return 0
    fi

    # The zip's internal layout is docs-site/dist/... so it extracts directly.
    if unzip -qo "${tmp_zip}" -d "${install_dir}"; then
        ok "Installed docs-site to ${install_dir}/docs-site"
    else
        err "Failed to extract docs-site.zip"
    fi

    rm -f "${tmp_zip}" 2>/dev/null || true
}

# ── Add to PATH ────────────────────────────────────────────────────

# detect_active_pwsh reports whether the user is running this installer
# from inside a PowerShell (pwsh) session. Checks multiple env signals
# because some shells / sudo wrappers strip PSModulePath but leave
# other pwsh-specific variables intact.
# Spec: spec/02-app-issues/29-macos-pwsh-shell-not-activated-after-install.md
detect_active_pwsh() {
    # PSModulePath: classic pwsh marker, always set in interactive sessions.
    if [ -n "${PSModulePath:-}" ]; then
        return 0
    fi
    # POWERSHELL_DISTRIBUTION_CHANNEL: set by the official pwsh package.
    if [ -n "${POWERSHELL_DISTRIBUTION_CHANNEL:-}" ]; then
        return 0
    fi
    # PSExecutionPolicyPreference: set when pwsh exports its policy.
    if [ -n "${PSExecutionPolicyPreference:-}" ]; then
        return 0
    fi
    # GITMAP_DUAL_SHELL: explicit override from `gitmap self-install --dual-shell`.
    if [ "${GITMAP_DUAL_SHELL:-}" = "1" ]; then
        return 0
    fi

    return 1
}

# pwsh_profile_path echoes the per-user pwsh profile path on Unix
# (Linux / macOS), creating the parent directory if needed.
pwsh_profile_path() {
    local dir="${HOME}/.config/powershell"
    mkdir -p "${dir}" 2>/dev/null || true
    echo "${dir}/Microsoft.PowerShell_profile.ps1"
}

# add_path_to_profile writes a marker-block snippet (per
# spec/04-generic-cli/21-post-install-shell-activation) to a single
# profile file. Idempotent across all three outcomes: appends when
# absent, rewrites in place when present-but-different, no-ops when
# present-and-identical. Third arg is the snippet shell flavour:
# "bash" | "fish" | "pwsh". (Legacy callers passed "false"/"true" —
# both are coerced to bash/fish.)
#
# Exit codes — used by add_to_path to render the per-file status table:
#   0 = added    (new block appended to the file)
#   1 = unchanged (block already present and byte-identical to snippet)
#   2 = updated  (block was present but body changed; rewritten in place)
add_path_to_profile() {
    local dir="$1" profile_file="$2" shell_kind="$3"

    # Back-compat coercion for old boolean callers.
    case "${shell_kind}" in
        true)  shell_kind="fish" ;;
        false|"") shell_kind="bash" ;;
    esac

    local marker_open="# gitmap shell wrapper v2 - managed by gitmap installer. Do not edit manually."
    local marker_close="# gitmap shell wrapper v2 end"

    local snippet
    snippet="$(resolve_path_snippet "${shell_kind}" "${dir}" "${marker_open}" "${marker_close}")"

    mkdir -p "$(dirname "${profile_file}")"
    touch "${profile_file}"

    if grep -qF "${marker_open}" "${profile_file}" 2>/dev/null; then
        rewrite_marker_block "${profile_file}" "${marker_open}" "${marker_close}" "${snippet}"

        return $?
    fi

    printf '\n%s\n' "${snippet}" >> "${profile_file}"

    return 0
}

# resolve_path_snippet returns the canonical PATH snippet body for a
# shell, preferring the freshly-installed gitmap binary's
# `setup print-path-snippet` output and falling back to a hard-coded
# template when the binary isn't yet on disk.
resolve_path_snippet() {
    local shell_kind="$1" dir="$2" open="$3" close="$4"
    local gitmap_bin=""
    if [ -x "${INSTALL_DIR:-}/${APP_SUBDIR}/gitmap" ]; then
        gitmap_bin="${INSTALL_DIR}/${APP_SUBDIR}/gitmap"
    elif [ -x "${INSTALL_DIR:-}/gitmap" ]; then
        # Pre-v3.13.11 fallback: top-level binary (very old unwrapped install).
        gitmap_bin="${INSTALL_DIR}/gitmap"
    elif command -v gitmap >/dev/null 2>&1; then
        gitmap_bin="$(command -v gitmap)"
    fi
    local snippet=""
    if [ -n "${gitmap_bin}" ]; then
        snippet="$("${gitmap_bin}" setup print-path-snippet \
            --shell "${shell_kind}" --dir "${dir}" --manager "installer" 2>/dev/null || true)"
    fi
    if [ -z "${snippet}" ]; then
        snippet="$(fallback_snippet "${shell_kind}" "${dir}" "${open}" "${close}")"
    fi
    printf '%s' "${snippet}"
}

# rewrite_marker_block replaces the existing marker block in profile_file
# with the supplied snippet. Compares byte-for-byte before writing so we
# can return 1 (unchanged) vs 2 (updated) — this is what powers the
# "already present" vs "updated" distinction in the per-file report.
# rewrite_marker_block replaces the existing marker block in profile_file
# with the supplied snippet. Compares byte-for-byte before writing so we
# can return 1 (unchanged) vs 2 (updated) — this is what powers the
# "already present" vs "updated" distinction in the per-file report.
#
# NOTE: awk vars are named open_marker / close_marker because `close`
# is a reserved gawk builtin and `close=...` triggers a fatal error.
# rewrite_marker_block replaces the existing marker block in profile_file
# with the supplied snippet. Compares byte-for-byte before writing so we
# can return 1 (unchanged) vs 2 (updated) — this is what powers the
# "already present" vs "updated" distinction in the per-file report.
#
# NOTE: awk vars are named open_marker / close_marker because `close`
# is a reserved gawk builtin and `close=...` triggers a fatal error.
# Diff uses `diff -q` (POSIX, ubiquitous) instead of `cmp` because some
# minimal containers omit cmp from the busybox build.
# rewrite_marker_block replaces the existing marker block in profile_file
# with the supplied snippet. Compares byte-for-byte before writing so we
# can return 1 (unchanged) vs 2 (updated) — this is what powers the
# "already present" vs "updated" distinction in the per-file report.
#
# NOTE: awk vars are named open_marker / close_marker because `close`
# is a reserved gawk builtin and `close=...` triggers a fatal error.
# Comparison uses files_equal (pure-bash byte read) so we don't depend
# on `cmp` or `diff` — both are missing in some minimal containers.
rewrite_marker_block() {
    local profile_file="$1" open="$2" close="$3" snippet="$4"
    local tmp
    tmp="$(mktemp)"
    awk -v open_marker="${open}" -v close_marker="${close}" -v body="${snippet}" '
        $0 == open_marker { skip = 1; print body; next }
        skip && $0 == close_marker { skip = 0; next }
        !skip { print }
    ' "${profile_file}" > "${tmp}"
    if files_equal "${tmp}" "${profile_file}"; then
        rm -f "${tmp}"

        return 1 # unchanged
    fi
    mv "${tmp}" "${profile_file}"

    return 2 # updated
}

# files_equal returns 0 iff $1 and $2 have identical bytes. Pure bash
# (no cmp/diff dependency). Compares size first as a fast path.
files_equal() {
    local a="$1" b="$2"
    local size_a size_b
    size_a=$(wc -c < "${a}" 2>/dev/null || echo -1)
    size_b=$(wc -c < "${b}" 2>/dev/null || echo -1)
    if [ "${size_a}" != "${size_b}" ]; then
        return 1
    fi
    local hash_a hash_b
    hash_a=$(sha1sum "${a}" 2>/dev/null | awk '{print $1}')
    hash_b=$(sha1sum "${b}" 2>/dev/null | awk '{print $1}')
    [ -n "${hash_a}" ] && [ "${hash_a}" = "${hash_b}" ]
}

# fallback_snippet renders the marker-block snippet body when the
# freshly-installed gitmap binary is unavailable to do it for us.
# Kept tiny and shell-specific so the per-shell flavour stays explicit.
fallback_snippet() {
    local kind="$1" dir="$2" open="$3" close="$4"
    case "${kind}" in
        fish)
            printf '%s\nset -gx GITMAP_WRAPPER 1\nfish_add_path %s\n%s' "${open}" "${dir}" "${close}"
            ;;
        pwsh)
            printf '%s\n$env:GITMAP_WRAPPER = "1"\nif (-not ($env:PATH -split '"'"':'"'"' | Where-Object { $_ -eq '"'"'%s'"'"' })) { $env:PATH = "$env:PATH:%s" }\n%s' "${open}" "${dir}" "${dir}" "${close}"
            ;;
        *)
            printf '%s\nexport GITMAP_WRAPPER=1\ncase ":${PATH}:" in *":%s:"*) ;; *) export PATH="$PATH:%s" ;; esac\n%s' "${open}" "${dir}" "${dir}" "${close}"
            ;;
    esac
}

add_to_path() {
    local dir="$1"
    local has_session_path=false

    case ":${PATH}:" in
        *":${dir}:"*)
            has_session_path=true
            ;;
    esac

    # Detect primary shell
    local shell_name
    shell_name="$(basename "${SHELL:-/bin/bash}")"
    PATH_SHELL="${shell_name}"

    local primary_profile=""
    # Three outcome lists drive both the live per-file table and the
    # final summary. Lists are space-separated paths, populated by
    # record_profile_outcome from add_path_to_profile's exit code.
    local profiles_added=""
    local profiles_updated=""
    local profiles_unchanged=""

    # ── Write to all relevant POSIX/bash/zsh profiles ──────────────
    # This ensures gitmap is available regardless of which shell the user opens.

    step "Writing PATH snippet to shell profiles..."

    # zsh profiles (both, to cover login + interactive shells)
    if should_write_profile zsh && { [ "${shell_name}" = "zsh" ] || [ -f "${HOME}/.zshrc" ] || [ -f "${HOME}/.zprofile" ]; }; then
        # .zshrc — interactive shells (most terminal emulators)
        add_path_to_profile "${dir}" "${HOME}/.zshrc" false
        record_profile_outcome $? "~/.zshrc"
        # .zprofile — login shells (macOS Terminal.app)
        add_path_to_profile "${dir}" "${HOME}/.zprofile" false
        record_profile_outcome $? "~/.zprofile"
    fi

    # bash profiles
    if should_write_profile bash && { [ "${shell_name}" = "bash" ] || [ -f "${HOME}/.bashrc" ] || [ -f "${HOME}/.bash_profile" ]; }; then
        add_path_to_profile "${dir}" "${HOME}/.bashrc" false
        record_profile_outcome $? "~/.bashrc"
        if [ -f "${HOME}/.bash_profile" ]; then
            add_path_to_profile "${dir}" "${HOME}/.bash_profile" false
            record_profile_outcome $? "~/.bash_profile"
        fi
    fi

    # POSIX ~/.profile — catch-all for sh and other POSIX shells.
    # Skipped under single-shell modes AND combo modes (zsh+pwsh, etc.)
    # to honour the "only the listed families" strict contract. Only
    # `auto` (detect everything) and `both` (write everything) include it.
    if [ "${PROFILE_MODE}" = "auto" ] || [ "${PROFILE_MODE}" = "both" ]; then
        add_path_to_profile "${dir}" "${HOME}/.profile" false
        record_profile_outcome $? "~/.profile"
    fi

    # fish (only if fish is installed or is the default shell)
    if should_write_profile fish && { [ "${shell_name}" = "fish" ] || command -v fish >/dev/null 2>&1; }; then
        local fish_config="${HOME}/.config/fish/config.fish"
        add_path_to_profile "${dir}" "${fish_config}" fish
        record_profile_outcome $? "~/.config/fish/config.fish"
    fi

    # PowerShell on Unix — detected when the installer was launched from
    # inside a pwsh session (PSModulePath is set), or when pwsh is on PATH.
    # The --shell-mode both / pwsh-containing combo (DUAL_SHELL=true)
    # forces this branch even when neither detection signal fires.
    # Issue: spec/02-app-issues/29-macos-pwsh-shell-not-activated-after-install.md
    local pwsh_active=false
    if detect_active_pwsh; then
        pwsh_active=true
        PATH_PWSH_DETECTED="yes (env signal)"
    fi
    local pwsh_force=false
    if [ "${DUAL_SHELL:-false}" = true ]; then
        pwsh_force=true
        PATH_PWSH_DETECTED="forced (--shell-mode ${PROFILE_MODE})"
    fi
    if should_write_profile pwsh && { [ "${pwsh_active}" = true ] || [ "${pwsh_force}" = true ] || command -v pwsh >/dev/null 2>&1; }; then
        if [ "${pwsh_active}" = false ] && [ "${pwsh_force}" = false ]; then
            PATH_PWSH_DETECTED="yes (pwsh on PATH)"
        fi
        local pwsh_profile
        pwsh_profile="$(pwsh_profile_path)"
        add_path_to_profile "${dir}" "${pwsh_profile}" pwsh
        record_profile_outcome $? "~/.config/powershell/Microsoft.PowerShell_profile.ps1"
    fi

    # If the user is actively in pwsh, the pwsh profile becomes the primary
    # reload target — overrides $SHELL-based detection (which on macOS still
    # reports zsh even when pwsh is the active interpreter).
    if [ "${pwsh_active}" = true ]; then
        PATH_SHELL="pwsh"
        shell_name="pwsh"
    fi

    # Determine primary profile for reload instruction
    case "${shell_name}" in
        zsh)    primary_profile="${HOME}/.zshrc" ;;
        bash)   primary_profile="${HOME}/.bashrc" ;;
        fish)   primary_profile="${HOME}/.config/fish/config.fish" ;;
        pwsh)   primary_profile="$(pwsh_profile_path)" ;;
        *)      primary_profile="${HOME}/.profile" ;;
    esac

    PATH_TARGET="${primary_profile}"

    case "${shell_name}" in
        fish)
            PATH_LINE="fish_add_path ${dir}"
            PATH_RELOAD="source ${primary_profile}"
            ;;
        pwsh)
            # Dot-source the pwsh profile — `. $PROFILE` is the canonical
            # PowerShell idiom for re-loading the current user's profile
            # in the active session. `source` is a bash builtin and would
            # error inside pwsh (`source: The term 'source' is not recognized`).
            PATH_LINE="\$env:PATH = \"\$env:PATH:${dir}\""
            PATH_RELOAD=". \$PROFILE"
            ;;
        *)
            PATH_LINE="export PATH=\"\$PATH:${dir}\""
            PATH_RELOAD=". ${primary_profile}"
            ;;
    esac

    # Cross-shell reload hint: when the installer touched BOTH the active
    # shell's profile AND the other family's profile (typical under
    # --profile both / --dual-shell, or when pwsh is detected alongside
    # zsh on macOS), expose the *other* shell's reload command as
    # PATH_RELOAD_ALT. Pass the UNION of all three lists — the alt is
    # about which profile *contains* the snippet, not whether we
    # rewrote it on this run.
    resolve_alt_reload "${profiles_added} ${profiles_updated} ${profiles_unchanged}"

    # Per-file change report: one line per profile with its idempotency
    # outcome so the user can audit exactly what changed and what was
    # already in place. Always printed (not gated behind --show-path)
    # because this is the answer to "did the installer touch my dotfiles?"
    print_profile_change_table

    # Persist the per-profile lists so print_install_summary + the audit
    # block can echo them. PATH_PROFILES_WRITTEN is the union of added +
    # updated for backward-compat with --show-path output; the new
    # PATH_PROFILES_UNCHANGED list exposes the no-op count separately.
    PATH_PROFILES_WRITTEN="${profiles_added# }${profiles_updated:+ }${profiles_updated# }"
    PATH_PROFILES_SKIPPED="${profiles_unchanged# }"
    PATH_PROFILES_ADDED="${profiles_added# }"
    PATH_PROFILES_UPDATED="${profiles_updated# }"
    PATH_PROFILES_UNCHANGED="${profiles_unchanged# }"

    # PATH_STATUS feeds the one-line summary in print_install_summary.
    if [ -n "${profiles_added}" ]; then
        PATH_STATUS="added"
    elif [ -n "${profiles_updated}" ]; then
        PATH_STATUS="updated"
    else
        PATH_STATUS="already present"
    fi

    if [ "${has_session_path}" = true ]; then
        return
    fi

    # Update current session (only effective when script is sourced, not piped)
    export PATH="${PATH}:${dir}"
}

# record_profile_outcome appends $2 to the appropriate per-status list
# based on the exit code from add_path_to_profile (passed as $1).
#   0 = added → profiles_added
#   1 = unchanged → profiles_unchanged
#   2 = updated → profiles_updated
# Uses the caller's locals via dynamic scope (bash function-local vars
# from add_to_path are visible here). Kept tiny so the call sites in
# add_to_path stay one line each.
record_profile_outcome() {
    local code="$1" path="$2"
    case "${code}" in
        0) profiles_added="${profiles_added} ${path}" ;;
        2) profiles_updated="${profiles_updated} ${path}" ;;
        *) profiles_unchanged="${profiles_unchanged} ${path}" ;;
    esac
}

# print_profile_change_table renders a per-file status line for every
# profile we touched. Symbols mirror common diff conventions so the
# output reads at a glance:
#   [+] added     — new gitmap block appended to a previously-untouched file
#   [~] updated   — block was present but body changed (e.g. install dir moved)
#   [=] unchanged — block already present and identical (true no-op)
print_profile_change_table() {
    local total
    total=$(count_words "${profiles_added}")
    total=$((total + $(count_words "${profiles_updated}")))
    total=$((total + $(count_words "${profiles_unchanged}")))
    if [ "${total}" -eq 0 ]; then
        warn "No shell profiles were eligible for the PATH snippet."
        return
    fi
    step "PATH snippet status (${total} profile$( [ ${total} -ne 1 ] && echo s)):"
    print_status_lines "+" 32 "${profiles_added}"
    print_status_lines "~" 33 "${profiles_updated}"
    print_status_lines "=" 90 "${profiles_unchanged}"
}

# print_status_lines emits one indented coloured line per path in $3,
# tagged with the symbol $1 (e.g. +/~/=) and ANSI colour code $2.
# Empty lists are silently skipped so a no-op category produces no row.
print_status_lines() {
    local symbol="$1" color="$2" list="$3" path
    for path in ${list}; do
        printf '    \033[%sm[%s]\033[0m %s\n' "${color}" "${symbol}" "${path}" >&2
    done
}

# count_words returns the number of whitespace-separated tokens in $1.
# Wrapper around `set --` so callers don't pollute their argv.
count_words() {
    # shellcheck disable=SC2086
    set -- ${1}
    echo $#
}

# resolve_alt_reload picks a secondary reload command for the *other*
# shell family, so dual-shell installs never show a syntactically wrong
# command (e.g. `source ~/.zshrc` while the user is sitting in pwsh).
#
# Selection rules (PATH_SHELL is the primary, already set above):
#   - PATH_SHELL=pwsh: alt is the first POSIX profile we wrote
#     (.zshrc preferred, then .bashrc, then .profile).
#   - PATH_SHELL=zsh|bash|fish|other: alt is the pwsh dot-source command
#     iff the pwsh profile was written (dual-shell or pwsh on PATH).
#
# Inputs:  $1 = profiles_written list (space-separated, e.g. " ~/.zshrc ~/.bashrc ~/.config/powershell/...")
# Outputs: PATH_RELOAD_ALT, PATH_RELOAD_ALT_SHELL (label like "zsh" or "pwsh")
resolve_alt_reload() {
    local written=" $1 "
    PATH_RELOAD_ALT=""
    PATH_RELOAD_ALT_SHELL=""
    if [ "${PATH_SHELL}" = "pwsh" ]; then
        pick_posix_alt "${written}"
        return
    fi
    case "${written}" in
        *powershell*)
            PATH_RELOAD_ALT=". \$PROFILE"
            PATH_RELOAD_ALT_SHELL="pwsh"
            ;;
    esac
}

# pick_posix_alt scans the written-profiles list for the highest-priority
# POSIX profile and sets PATH_RELOAD_ALT to its `source` command. Split
# out so resolve_alt_reload stays a clear top-level dispatch.
pick_posix_alt() {
    local written="$1"
    case "${written}" in
        *.zshrc*)
            PATH_RELOAD_ALT="source ~/.zshrc"
            PATH_RELOAD_ALT_SHELL="zsh"
            ;;
        *.bashrc*)
            PATH_RELOAD_ALT="source ~/.bashrc"
            PATH_RELOAD_ALT_SHELL="bash"
            ;;
        *.profile*)
            PATH_RELOAD_ALT=". ~/.profile"
            PATH_RELOAD_ALT_SHELL="sh"
            ;;
    esac
}

print_install_summary() {
    local installed_version="$1" bin_path="$2"

    echo ""
    step "Install summary"
    printf '    Version: %s\n' "${installed_version}" >&2
    printf '    Binary: %s\n' "${bin_path}" >&2
    printf '    Install dir: %s\n' "$(dirname "${bin_path}")" >&2
    if [ "${NO_PATH}" = true ]; then
        printf '    PATH target: skipped (--no-path)\n' >&2

        return
    fi
    printf '    Shell: %s\n' "${PATH_SHELL}" >&2
    printf '    PATH target: %s (%s)\n' "${PATH_TARGET}" "${PATH_STATUS}" >&2
    printf '    Reload: %s\n' "${PATH_RELOAD}" >&2
    if [ -n "${PATH_RELOAD_ALT}" ]; then
        printf '    Reload (%s): %s\n' "${PATH_RELOAD_ALT_SHELL}" "${PATH_RELOAD_ALT}" >&2
    fi

    # --show-path expands the summary with the full audit trail so the
    # user can confirm every profile file we touched and why a particular
    # shell was chosen. Triggered by SHOW_PATH=true (set in parse_args).
    if [ "${SHOW_PATH:-false}" = true ]; then
        print_path_audit
    fi
}

# print_path_audit dumps the detected shell, pwsh signal, and the full
# list of profile files written/skipped. Called only under --show-path
# to keep the default summary terse.
# print_path_audit dumps the detected shell, pwsh signal, and the full
# per-status profile lists. Called only under --show-path to keep the
# default summary terse — but the per-file change table (printed
# unconditionally by print_profile_change_table) already shows the
# important info; this audit adds the "why this shell was chosen" detail.
print_path_audit() {
    printf '\n  \033[36m%s\033[0m\n' "PATH audit (--show-path)" >&2
    printf '    Detected SHELL env: %s\n' "${SHELL:-<unset>}" >&2
    printf '    pwsh detected: %s\n' "${PATH_PWSH_DETECTED}" >&2
    printf '    Profile mode: %s\n' "${PROFILE_MODE:-auto}" >&2
    print_audit_list "Profiles added"           "${PATH_PROFILES_ADDED}"
    print_audit_list "Profiles updated"         "${PATH_PROFILES_UPDATED}"
    print_audit_list "Profiles already in sync" "${PATH_PROFILES_UNCHANGED}"
    printf '    PATH line: %s\n' "${PATH_LINE}" >&2
}

# print_audit_list renders one labeled line for a profile-status list.
# Empty lists print "<none>" so users can confirm a category was checked
# but produced no entries (vs being skipped entirely).
print_audit_list() {
    local label="$1" list="$2"
    if [ -n "${list}" ]; then
        printf '    %s: %s\n' "${label}" "${list}" >&2
    else
        printf '    %s: <none>\n' "${label}" >&2
    fi
}

# ── Resolve install directory ──────────────────────────────────────

resolve_install_dir() {
    local dir="$1"
    if [ -n "${dir}" ]; then
        echo "${dir}"
        return
    fi

    # Use ~/.local/bin if it exists or is standard; fallback to /usr/local/bin
    if [ -d "${HOME}/.local/bin" ] || [ -w "${HOME}/.local" ]; then
        echo "${HOME}/.local/bin"
    elif [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
    else
        echo "${HOME}/.local/bin"
    fi
}

# ── Parse arguments ────────────────────────────────────────────────

parse_args() {
    VERSION=""
    INSTALL_DIR=""
    ARCH_FLAG=""
    NO_PATH=false
    NO_DISCOVERY=false
    PROBE_CEILING=30
    DUAL_SHELL=false
    SHOW_PATH=false
    # PROFILE_MODE controls which shell profiles get the PATH snippet.
    # Accepted values:
    #   auto   = run detection (default; current behavior)
    #   both   = write all supported profiles (zsh + bash + .profile + fish + pwsh)
    #   zsh|bash|pwsh|fish = restrict writes to exactly that one shell family
    #   <a>+<b>[+<c>...] = strict combo, only the listed families (e.g. zsh+pwsh)
    # Combos are STRICT — ~/.profile and undeclared families are skipped.
    # The Go caller (gitmap self-install --shell-mode <mode>) already
    # validated this; we re-validate below for direct (curl|bash) users.
    # See spec/02-app-issues/29-macos-pwsh-shell-not-activated-after-install.md
    # for the original motivating use case (pwsh user on macOS).
    PROFILE_MODE="auto"

    while [ $# -gt 0 ]; do
        case "$1" in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --arch)
                ARCH_FLAG="$2"
                shift 2
                ;;
            --no-path)
                NO_PATH=true
                shift
                ;;
            --no-discovery)
                NO_DISCOVERY=true
                shift
                ;;
            --probe-ceiling)
                PROBE_CEILING="$2"
                shift 2
                ;;
            --shell-mode|--profile)
                # --shell-mode is canonical (v3.48.0+); --profile is the
                # v3.46.0 alias. Both accept singletons (auto|both|zsh|
                # bash|pwsh|fish) and `+`-joined combos (e.g. zsh+pwsh).
                # Validation lives in validate_shell_mode so curl|bash
                # users get the same error contract as gitmap self-install
                # users. Combos containing pwsh export GITMAP_DUAL_SHELL=1
                # so detect_active_pwsh fires regardless of $SHELL.
                PROFILE_MODE="$2"
                validate_shell_mode "${PROFILE_MODE}"
                if shell_mode_includes_pwsh "${PROFILE_MODE}"; then
                    DUAL_SHELL=true
                    export GITMAP_DUAL_SHELL=1
                fi
                shift 2
                ;;
            --dual-shell)
                # Hidden alias for --shell-mode both. Kept for backward
                # compatibility with v3.43–v3.45 callers; new code should
                # use --shell-mode both.
                PROFILE_MODE="both"
                DUAL_SHELL=true
                export GITMAP_DUAL_SHELL=1
                shift
                ;;
            --show-path)
                # Expand the install summary with the per-profile audit
                # trail (detected shell, pwsh signal, every profile file
                # touched). Diagnostic flag; no behavior change otherwise.
                SHOW_PATH=true
                shift
                ;;
            --help|-h)
                echo "Usage: install.sh [--version <tag>] [--dir <path>] [--arch <arch>] [--no-path] [--no-discovery] [--probe-ceiling <N>] [--shell-mode <mode>] [--show-path]"
                echo ""
                echo "Options:"
                echo "  --version <tag>        Install a specific version (e.g. v2.55.0)"
                echo "  --dir <path>           Target directory (default: ~/.local/bin)"
                echo "  --arch <arch>          Force architecture: amd64, arm64 (default: auto)"
                echo "  --no-path              Skip adding install directory to PATH"
                echo "  --no-discovery         Skip versioned-repo discovery (install baseline)"
                echo "  --probe-ceiling <N>    Highest -v<N> to probe (default: 30)"
                echo "  --shell-mode <mode>    Which shell profile(s) to write:"
                echo "                           auto|both|zsh|bash|pwsh|fish"
                echo "                           or '+'-joined combo (e.g. zsh+pwsh, bash+fish)"
                echo "  --profile <mode>       Deprecated alias for --shell-mode"
                echo "  --show-path            Print detected shell + every profile file written"
                exit 0
                ;;
            *)
                err "Unknown option: $1"
                err "Run with --help for usage."
                exit 1
                ;;
        esac
    done
}

# validate_shell_mode rejects unknown --shell-mode / --profile values.
# Accepts the six singletons (auto|both|zsh|bash|pwsh|fish) and any
# `+`-joined combo of concrete shell families (zsh, bash, pwsh, fish);
# meta values `auto` and `both` are not valid inside combos.
# Exits 1 with a clear list so curl|bash users get the same error
# contract as `gitmap self-install` users.
validate_shell_mode() {
    local mode="$1"
    case "${mode}" in
        auto|both|zsh|bash|pwsh|fish) return 0 ;;
    esac
    if validate_shell_mode_combo "${mode}"; then
        return 0
    fi
    err "--shell-mode '${mode}' is not valid."
    err "Accepted singletons: auto|both|zsh|bash|pwsh|fish"
    err "Accepted combos:     '+'-joined concrete shells (e.g. zsh+pwsh, bash+fish, zsh+bash+pwsh)"
    exit 1
}

# validate_shell_mode_combo returns 0 iff $1 is a `+`-joined list of two
# or more distinct concrete shell families. Helper for validate_shell_mode.
validate_shell_mode_combo() {
    local mode="$1"
    case "${mode}" in
        *+*) ;;
        *) return 1 ;;
    esac
    local oldifs="${IFS}"
    IFS='+'
    # shellcheck disable=SC2086
    set -- ${mode}
    IFS="${oldifs}"
    if [ $# -lt 2 ]; then
        return 1
    fi
    local seen=""
    local tok
    for tok in "$@"; do
        case "${tok}" in
            zsh|bash|pwsh|fish) ;;
            *) return 1 ;;
        esac
        case " ${seen} " in
            *" ${tok} "*) return 1 ;;
        esac
        seen="${seen} ${tok}"
    done

    return 0
}

# validate_profile_mode is the legacy entry point kept for any external
# scripts that source install.sh and call it directly. Delegates to
# validate_shell_mode.
validate_profile_mode() {
    validate_shell_mode "$1"
}

# shell_mode_includes_pwsh reports whether $1 forces a pwsh profile
# write. Covers `both`, the bare `pwsh` singleton, and any combo whose
# `+`-joined tokens include `pwsh`. Used to decide whether to export
# GITMAP_DUAL_SHELL=1.
shell_mode_includes_pwsh() {
    local mode="$1"
    case "${mode}" in
        both|pwsh) return 0 ;;
        *+pwsh|pwsh+*|*+pwsh+*) return 0 ;;
        *) return 1 ;;
    esac
}

# should_write_profile reports whether the profile family $1 should
# receive the PATH snippet, given the current PROFILE_MODE.
#   auto → caller's existing detection logic decides (return 0 always;
#          the caller's own `if` already filters).
#   both → always yes.
#   <shell> singleton → only if it matches.
#   <a>+<b>[+<c>...] combo → yes iff family appears as a `+`-token.
# Combos are strict: anything not listed (including ~/.profile, which
# the caller checks separately) is excluded.
should_write_profile() {
    local family="$1"
    case "${PROFILE_MODE}" in
        auto|both) return 0 ;;
        "${family}") return 0 ;;
        *+*)
            case "+${PROFILE_MODE}+" in
                *"+${family}+"*) return 0 ;;
                *) return 1 ;;
            esac
            ;;
        *) return 1 ;;
    esac
}

# ── Main ───────────────────────────────────────────────────────────

# verify_installation runs the three post-install checks the user
# asked for: (1) print the installed version by invoking the binary,
# (2) confirm `gitmap` resolves on PATH in the current shell, (3)
# ensure the per-install data folder exists (create on miss). Each
# check prints PASS/WARN with a one-line reason; failures never
# abort because the binary is already on disk and recoverable.
verify_installation() {
    local bin_path="$1" app_dir="$2"
    local data_dir="${app_dir}/data"

    echo ""
    step "Verifying installation"

    # 1. Version check — use the binary directly (does not depend on
    # PATH), so a PATH miss in step 2 still shows a real version.
    if [ -x "${bin_path}" ] && "${bin_path}" version >/dev/null 2>&1; then
        printf '    \033[32mPASS\033[0m  Version: %s\n' \
            "$("${bin_path}" version 2>&1 | head -n 1)" >&2
    else
        printf '    \033[33mWARN\033[0m  Could not run %s version\n' "${bin_path}" >&2
    fi

    # 2. PATH check — `command -v` reflects the *current* shell, which
    # for fresh installs may still be stale until the user reloads
    # their profile. Treat a miss as a warning, not a failure.
    if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
        printf '    \033[32mPASS\033[0m  PATH active: %s resolves to %s\n' \
            "${BINARY_NAME}" "$(command -v "${BINARY_NAME}")" >&2
    elif [ "${NO_PATH}" = true ]; then
        printf '    \033[33mWARN\033[0m  PATH skipped (--no-path); invoke with full path: %s\n' \
            "${bin_path}" >&2
    else
        printf '    \033[33mWARN\033[0m  %s not on PATH yet — reload your shell: %s\n' \
            "${BINARY_NAME}" "${PATH_RELOAD}" >&2
    fi

    # 3. Data folder — ensure it exists so the first `gitmap scan`
    # does not race on directory creation. Creating it here is safe
    # because the install dir is already owned by the current user.
    if [ -d "${data_dir}" ]; then
        printf '    \033[32mPASS\033[0m  Data folder exists: %s\n' "${data_dir}" >&2
    elif mkdir -p "${data_dir}" 2>/dev/null; then
        printf '    \033[32mPASS\033[0m  Data folder created: %s\n' "${data_dir}" >&2
    else
        printf '    \033[33mWARN\033[0m  Could not create data folder: %s\n' "${data_dir}" >&2
    fi
}

main() {
    echo ""
    echo "  gitmap installer"
    printf '  \033[90mgithub.com/%s\033[0m\n' "${REPO}"
    echo ""

    parse_args "$@"
    load_deploy_manifest

    # Versioned repo discovery: re-exec from the latest -v<M> sibling repo.
    if [ "${INSTALLER_DELEGATED:-0}" = "1" ]; then
        printf '  [discovery] INSTALLER_DELEGATED=1; skipping discovery (loop guard)\n' >&2
    elif [ "${NO_DISCOVERY}" = "true" ]; then
        printf '  [discovery] --no-discovery set; skipping probe\n' >&2
    elif [ -n "${VERSION}" ]; then
        # Pinned-version contract (spec/07-generic-release/08-pinned-version-install-snippet.md):
        # When --version is supplied, install EXACTLY that version from the embedded REPO.
        # Skip versioned-repo discovery so a snippet copied from a v3.x release page
        # never silently jumps to the v4 repo's latest tag.
        printf '  [discovery] --version %s pinned; skipping repo probe (exact-version install)\n' "${VERSION}" >&2
    else
        local effective_repo
        effective_repo="$(resolve_effective_repo "${REPO}" "${PROBE_CEILING}")"
        if [ "${effective_repo}" != "${REPO}" ]; then
            invoke_delegated_full_installer "${effective_repo}" "$@" || true
        fi
    fi

    local os arch version install_dir archive_path

    os="$(detect_os)"
    arch="$(detect_arch "${ARCH_FLAG}")"
    version="$(resolve_version "${VERSION}")"
    install_dir="$(resolve_install_dir "${INSTALL_DIR}")"

    # Create TMP_DIR in parent scope so install_binary and cleanup can access it.
    TMP_DIR="$(mktemp -d)"
    archive_path="$(download_asset "${version}" "${os}" "${arch}")"

    APP_DIR=""
    install_binary "${archive_path}" "${install_dir}" "${os}" "${arch}" "${version}"

    # Bundle the docs site so `gitmap help-dashboard` works after install.
    install_docs_site "${version}" "${APP_DIR}"

    if [ "${NO_PATH}" = false ]; then
        add_to_path "${APP_DIR}"
    fi

    # Verify the binary works.
    #
    # Discard stderr (`2>/dev/null`) and parse only stdout. First-run side
    # effects like SeedDownloaderConfig print informational lines to stderr
    # (e.g. "◦ Downloader config already customized …") that would otherwise
    # be concatenated into ${installed_version} and shown to the user. We
    # then filter to the canonical `gitmap vX.Y.Z` line via awk (no pipe to
    # `head` to avoid SIGPIPE under `set -o pipefail`).
    local bin_path="${APP_DIR}/${BINARY_NAME}"
    local installed_version="${version}"
    if [ -f "${bin_path}" ]; then
        echo ""
        local version_output
        if version_output="$("${bin_path}" version 2>/dev/null)"; then
            local version_line
            version_line="$(printf '%s\n' "${version_output}" | awk '/^gitmap v[0-9]/{print; exit}')"
            if [ -n "${version_line}" ]; then
                installed_version="${version_line}"
                ok "${version_line}"
            else
                ok "gitmap ${version}"
            fi
        else
            err "Binary found but failed to run."
        fi
    else
        err "Binary not found at ${bin_path}"
    fi

    print_install_summary "${installed_version}" "${bin_path}"
    if [ "${NO_PATH}" = false ]; then
        echo ""
        printf '  \033[32mOK\033[0m  To start using gitmap \033[1mright now\033[0m, run:\n' >&2
        echo "" >&2
        # Label the primary command with its shell so the user can see at
        # a glance whether it matches the shell they're sitting in. The
        # PATH_SHELL value is already pwsh-overridden when detect_active_pwsh
        # fired, so this label is the source of truth for the active shell.
        printf '      \033[36m%s\033[0m   \033[90m# in %s\033[0m\n' "${PATH_RELOAD}" "${PATH_SHELL}" >&2
        # When dual-shell wrote both a POSIX profile AND the pwsh profile,
        # show the alternate command so users in the *other* shell aren't
        # left with a syntactically wrong hint (e.g. `source ~/.zshrc`
        # pasted into a pwsh prompt).
        if [ -n "${PATH_RELOAD_ALT}" ]; then
            printf '      \033[36m%s\033[0m   \033[90m# in %s\033[0m\n' "${PATH_RELOAD_ALT}" "${PATH_RELOAD_ALT_SHELL}" >&2
        fi
        echo "" >&2
        printf '     Or open a new terminal window.\n' >&2
        echo "" >&2
        printf '  \033[90mInstalled to: %s\033[0m\n' "${bin_path}" >&2
        printf '  \033[90mApp folder on PATH: %s\033[0m\n' "${APP_DIR}" >&2
    fi

    verify_installation "${bin_path}" "${APP_DIR}"

    # v5.18.0+: auto-run `gitmap setup` so the shell wrapper (gcd /
    # gitmap function) + completions are installed without a second
    # command. Setup is idempotent (marker `# gitmap shell wrapper v2`).
    # Non-fatal: install itself already succeeded.
    if [ -x "${bin_path}" ]; then
        echo ""
        printf '  \033[36m→ Running '\''gitmap setup'\'' to install shell wrapper + completions...\033[0m\n' >&2
        "${bin_path}" setup || printf '  \033[33m(setup auto-run skipped — run '\''gitmap setup'\'' manually)\033[0m\n' >&2
    fi

    echo ""
    ok "Done! Run 'gitmap --help' to get started."
    echo ""
}

main "$@"

