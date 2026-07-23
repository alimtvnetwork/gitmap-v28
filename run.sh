#!/usr/bin/env bash
#
# Build, deploy, and run gitmap CLI from the repo root.
#
# Usage:
#   ./run.sh                          # pull, build, deploy
#   ./run.sh --no-pull                # skip git pull
#   ./run.sh --force-pull             # discard local changes + pull (no prompt)
#   ./run.sh --no-deploy              # skip deploy step
#   ./run.sh -r scan                  # build + scan parent folder
#   ./run.sh -r scan ~/repos          # build + scan specific path
#   ./run.sh -r help                  # build + show help
#   ./run.sh -t                       # run all unit tests with reports
#
# Configuration is read from gitmap/powershell.json (same as run.ps1).
# --force-pull automatically discards local changes and removes untracked
# files before pulling. Useful for CI or unattended builds.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"
GITMAP_DIR="$REPO_ROOT/gitmap"

# -- Defaults --------------------------------------------------
NO_PULL=false
NO_DEPLOY=false
FORCE_PULL=false
DEPLOY_PATH=""
UPDATE=false
RUN=false
TEST=false
DEBUG_REPO_DETECT=false
RUN_ARGS=()

# -- Parse arguments -------------------------------------------
while [[ $# -gt 0 ]]; do
    case "$1" in
        --no-pull)             NO_PULL=true; shift ;;
        --force-pull)          FORCE_PULL=true; shift ;;
        --no-deploy)           NO_DEPLOY=true; shift ;;
        --deploy-path)         DEPLOY_PATH="$2"; shift 2 ;;
        --update)              UPDATE=true; shift ;;
        --debug-repo-detect)   DEBUG_REPO_DETECT=true; shift ;;
        -r|--run)              RUN=true; shift
            # Collect remaining args for gitmap
            while [[ $# -gt 0 ]]; do
                RUN_ARGS+=("$1"); shift
            done
            ;;
        -t|--test)             TEST=true; shift ;;
        *)
            echo "Unknown option: $1" >&2; exit 1 ;;
    esac
done

# Honor env-var bridge from `gitmap update --debug-repo-detect`.
if [[ "${GITMAP_DEBUG_REPO_DETECT:-}" == "1" ]]; then
    DEBUG_REPO_DETECT=true
fi

# -- Colors ----------------------------------------------------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
GRAY='\033[0;37m'
NC='\033[0m'

# -- Logging helpers -------------------------------------------
write_step() {
    echo ""
    echo -e "  ${MAGENTA}[$1]${NC} $2"
    echo "  --------------------------------------------------"
}

write_success() { echo -e "  ${GREEN}OK${NC} ${GREEN}$1${NC}"; }
write_info()    { echo -e "  ${CYAN}->${NC} ${GRAY}$1${NC}"; }
write_warn()    { echo -e "  ${YELLOW}!!${NC} ${YELLOW}$1${NC}"; }
write_fail()    { echo -e "  ${RED}XX${NC} ${RED}$1${NC}"; }

# -- Error reporting (JSONL) -----------------------------------
# When run from `gitmap update --report-errors json`, env vars
# GITMAP_REPORT_ERRORS=json and GITMAP_REPORT_ERRORS_FILE=<path>
# are set. Each non-fatal failure appends one JSON object per line.
# Args: stage command exit_code message [extra_json_object]
write_report_error() {
    local stage="$1"
    local command="$2"
    local exit_code="$3"
    local message="$4"
    local extra="${5:-{\}}"
    if [[ "${GITMAP_REPORT_ERRORS:-}" != "json" ]] || [[ -z "${GITMAP_REPORT_ERRORS_FILE:-}" ]]; then
        return 0
    fi
    local ts
    ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    local cwd
    cwd="$(pwd)"
    # Escape backslashes and double quotes for JSON safety.
    local esc_msg="${message//\\/\\\\}"; esc_msg="${esc_msg//\"/\\\"}"
    local esc_cmd="${command//\\/\\\\}"; esc_cmd="${esc_cmd//\"/\\\"}"
    local esc_cwd="${cwd//\\/\\\\}"; esc_cwd="${esc_cwd//\"/\\\"}"
    local line
    line="{\"timestamp\":\"${ts}\",\"stage\":\"${stage}\",\"command\":\"${esc_cmd}\",\"exitCode\":${exit_code},\"cwd\":\"${esc_cwd}\",\"message\":\"${esc_msg}\",\"paths\":${extra},\"os\":\"unix\"}"
    if ! printf '%s\n' "$line" >> "${GITMAP_REPORT_ERRORS_FILE}" 2>/dev/null; then
        write_warn "Could not append to report-errors file: ${GITMAP_REPORT_ERRORS_FILE}"
    fi
}

# -- Repo-detect debug -----------------------------------------
# Active when --debug-repo-detect is passed OR GITMAP_DEBUG_REPO_DETECT=1
# (set by `gitmap update --debug-repo-detect`). Prints structured marker
# checks and decision reasons. Mirrors entries to JSONL report when
# --report-errors json is also active.
is_debug_repo_detect() {
    [[ "${DEBUG_REPO_DETECT:-false}" == "true" ]] || [[ "${GITMAP_DEBUG_REPO_DETECT:-}" == "1" ]]
}

write_repo_detect() {
    local check="$1"
    local result="$2"
    local detail="${3:-}"
    if ! is_debug_repo_detect; then
        return 0
    fi
    if [[ -n "$detail" ]]; then
        printf "  ${CYAN}[DETECT]${NC} %-28s = %s  ${GRAY}(%s)${NC}\n" "$check" "$result" "$detail"
    else
        printf "  ${CYAN}[DETECT]${NC} %-28s = %s\n" "$check" "$result"
    fi
    # Mirror to JSONL report file when active.
    write_report_error "repo-detect" "$check" 0 "$result" \
        "{\"detail\":\"${detail//\"/\\\"}\",\"level\":\"info\"}"
}

write_repo_detect_snippet() {
    local title="$1"
    local path="$2"
    local max_lines="${3:-6}"
    if ! is_debug_repo_detect; then
        return 0
    fi
    echo -e "  ${CYAN}[DETECT]${NC} ${title} :"
    if [[ ! -f "$path" ]]; then
        echo -e "    ${GRAY}(file not found: $path)${NC}"
        return 0
    fi
    head -n "$max_lines" "$path" 2>/dev/null | sed "s/^/    /" | while IFS= read -r line; do
        echo -e "${GRAY}${line}${NC}"
    done
}

# -- Banner ----------------------------------------------------
show_banner() {
    echo ""
    echo -e "  ${CYAN}+======================================+${NC}"
    echo -e "  ${CYAN}|         ${CYAN}gitmap builder${CYAN}              |${NC}"
    echo -e "  ${CYAN}+======================================+${NC}"
    echo ""
}

# -- Load deploy manifest (single source of truth) -------------
# Mirrors run.ps1's Get-DeployManifest. Reads gitmap/constants/deploy-manifest.json
# so APP_SUBDIR / LEGACY_APP_SUBDIRS aren't hardcoded. Renaming the deploy
# folder ONLY requires editing that JSON file.
APP_SUBDIR="gitmap-cli"
LEGACY_APP_SUBDIRS=("gitmap")
load_deploy_manifest() {
    local manifest="$GITMAP_DIR/constants/deploy-manifest.json"
    if [[ ! -f "$manifest" ]]; then
        write_warn "deploy-manifest.json not found at $manifest — using defaults"
        return
    fi
    local app
    app=$(grep -E '"appSubdir"' "$manifest" | head -n1 | sed -E 's/.*"appSubdir"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')
    if [[ -n "$app" ]]; then
        APP_SUBDIR="$app"
    fi
    # legacyAppSubdirs is a JSON array — extract values between [ ... ].
    local legacy_block
    legacy_block=$(awk '/"legacyAppSubdirs"/,/]/' "$manifest")
    if [[ -n "$legacy_block" ]]; then
        local legacy_csv
        legacy_csv=$(echo "$legacy_block" | grep -oE '"[^"]+"' | sed -E 's/^"(.*)"$/\1/' | grep -v '^legacyAppSubdirs$' || true)
        if [[ -n "$legacy_csv" ]]; then
            LEGACY_APP_SUBDIRS=()
            while IFS= read -r line; do
                [[ -n "$line" ]] && LEGACY_APP_SUBDIRS+=("$line")
            done <<< "$legacy_csv"
        fi
    fi
}

# is_known_app_subdir returns 0 if $1 matches APP_SUBDIR or any legacy entry.
is_known_app_subdir() {
    local name="$1"
    [[ "$name" == "$APP_SUBDIR" ]] && return 0
    local legacy
    for legacy in "${LEGACY_APP_SUBDIRS[@]}"; do
        [[ "$name" == "$legacy" ]] && return 0
    done
    return 1
}

# -- Load config -----------------------------------------------
load_config() {
    load_deploy_manifest
    local config_path="$GITMAP_DIR/powershell.json"
    if [[ -f "$config_path" ]]; then
        write_info "Config loaded from powershell.json"
    else
        write_warn "No powershell.json found, using defaults"
    fi
}

get_config_value() {
    local key="$1"
    local default="$2"
    local config_path="$GITMAP_DIR/powershell.json"
    if [[ -f "$config_path" ]] && command -v python3 &>/dev/null; then
        python3 -c "import json,sys; d=json.load(open('$config_path')); print(d.get('$key','$default'))" 2>/dev/null || echo "$default"
    elif [[ -f "$config_path" ]] && command -v jq &>/dev/null; then
        jq -r ".$key // \"$default\"" "$config_path" 2>/dev/null || echo "$default"
    else
        echo "$default"
    fi
}

# -- Platform detection ----------------------------------------
BINARY_NAME="gitmap"
if [[ "$(uname -s)" == *MINGW* ]] || [[ "$(uname -s)" == *MSYS* ]]; then
    BINARY_NAME="gitmap.exe"
fi

BUILD_OUTPUT=$(get_config_value "buildOutput" "./bin")
DEPLOY_TARGET=$(get_config_value "deployPath" "$HOME/bin-run")
COPY_DATA=$(get_config_value "copyData" "true")

# -- Ensure main branch ----------------------------------------
ensure_main_branch() {
    cd "$REPO_ROOT"
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
    if [[ "$current_branch" != "main" ]]; then
        write_warn "Currently on branch '$current_branch', switching to main..."
        if ! git checkout main 2>&1; then
            write_fail "Failed to switch to main branch"
            exit 1
        fi
        write_success "Switched to main branch"
    fi
}

# -- Git pull --------------------------------------------------
invoke_git_pull() {
    write_step "1/4" "Pulling latest changes"
    ensure_main_branch
    cd "$REPO_ROOT"
    local output pull_exit
    set +e
    output=$(git pull 2>&1)
    pull_exit=$?
    set -e

    while IFS= read -r line; do
        [[ -n "$line" ]] && write_info "$line"
    done <<< "$output"

    if [[ $pull_exit -ne 0 ]]; then
        if echo "$output" | grep -qiE "Your local changes|overwritten by merge|not possible because you have unmerged|Please commit your changes or stash them"; then
            if [[ "$FORCE_PULL" == "true" ]]; then
                force_pull_clean
            else
                resolve_pull_conflict
            fi
        else
            write_fail "Git pull failed (exit code $pull_exit)"
            exit 1
        fi
    else
        write_success "Pull complete"
    fi
}

# -- Force pull: discard + clean without prompting -------------
force_pull_clean() {
    write_warn "Force-pull: discarding local changes and removing untracked files..."
    if ! git checkout -- . 2>&1; then
        write_fail "Git checkout failed"
        exit 1
    fi
    write_success "Local changes discarded"

    local clean_output
    clean_output=$(git clean -fd 2>&1) || true
    if [[ -n "$clean_output" ]]; then
        local clean_count
        clean_count=$(echo "$clean_output" | grep -c . || true)
        write_success "Removed $clean_count untracked file(s)"
    fi

    retry_git_pull
}

# -- Resolve pull conflict with local changes ------------------
resolve_pull_conflict() {
    write_warn "Git pull failed due to local changes"
    echo ""
    echo -e "  ${YELLOW}Choose how to proceed:${NC}"
    echo -e "    ${CYAN}[S] Stash changes (save for later, then pull)${NC}"
    echo -e "    ${CYAN}[D] Discard changes (reset working tree, then pull)${NC}"
    echo -e "    ${CYAN}[C] Clean all (discard changes + remove untracked files, then pull)${NC}"
    echo -e "    ${CYAN}[Q] Quit (abort without changes)${NC}"
    echo ""
    read -rp "  Enter choice (S/D/C/Q): " choice

    case "$(echo "$choice" | tr '[:lower:]' '[:upper:]')" in
        S)
            write_info "Stashing local changes..."
            local stash_output
            if stash_output=$(git stash push -m "auto-stash before run.sh pull" 2>&1); then
                write_success "Changes stashed"
                write_info "Run 'git stash pop' later to restore your changes"
                retry_git_pull
            else
                write_fail "Git stash failed"
                echo "$stash_output" >&2
                exit 1
            fi
            ;;
        D)
            write_warn "Discarding all local changes..."
            if git checkout -- . 2>&1; then
                write_success "Local changes discarded"
                retry_git_pull
            else
                write_fail "Git checkout failed"
                exit 1
            fi
            ;;
        C)
            write_warn "Discarding all local changes and removing untracked files..."
            if ! git checkout -- . 2>&1; then
                write_fail "Git checkout failed"
                exit 1
            fi
            write_success "Local changes discarded"

            local clean_output
            clean_output=$(git clean -fd 2>&1) || true
            if [[ -n "$clean_output" ]]; then
                while IFS= read -r line; do
                    [[ -n "$line" ]] && write_info "$line"
                done <<< "$clean_output"
                local clean_count
                clean_count=$(echo "$clean_output" | grep -c . || true)
                write_success "Removed $clean_count untracked file(s)"
            else
                write_info "No untracked files to remove"
            fi

            retry_git_pull
            ;;
        *)
            write_info "Aborted by user"
            exit 0
            ;;
    esac
}

# -- Retry git pull after stash/discard -----------------------
retry_git_pull() {
    write_info "Retrying git pull..."
    local output pull_exit
    set +e
    output=$(git pull 2>&1)
    pull_exit=$?
    set -e

    while IFS= read -r line; do
        [[ -n "$line" ]] && write_info "$line"
    done <<< "$output"

    if [[ $pull_exit -ne 0 ]]; then
        write_fail "Git pull failed again (exit code $pull_exit)"
        exit 1
    fi

    write_success "Pull complete"
}

# -- Resolve dependencies -------------------------------------
resolve_dependencies() {
    write_step "2/4" "Resolving Go dependencies"
    cd "$GITMAP_DIR"
    if ! go mod tidy 2>&1; then
        write_fail "go mod tidy failed"
        exit 1
    fi
    write_success "Dependencies resolved"
}

# -- Pre-build validation --------------------------------------
test_source_files() {
    write_info "Validating source files..."

    local required_files=(
        "main.go"
        "go.mod"
        "cmd/root.go"
        "cmd/scan.go"
        "cmd/clone.go"
        "cmd/update.go"
        "cmd/pull.go"
        "cmd/rescan.go"
        "cmd/desktopsync.go"
        "constants/constants.go"
        "config/config.go"
        "scanner/scanner.go"
        "mapper/mapper.go"
        "model/record.go"
        "formatter/csv.go"
        "formatter/json.go"
        "formatter/terminal.go"
        "formatter/text.go"
        "formatter/structure.go"
        "formatter/clonescript.go"
        "formatter/directclone.go"
        "formatter/desktopscript.go"
        "cloner/cloner.go"
        "cloner/safe_pull.go"
        "gitutil/gitutil.go"
        "desktop/desktop.go"
        "verbose/verbose.go"
        "setup/setup.go"
        "cmd/setup.go"
        "cmd/status.go"
        "cmd/exec.go"
        "cmd/release.go"
        "cmd/releasebranch.go"
        "cmd/releasepending.go"
        "cmd/changelog.go"
        "cmd/doctor.go"
        "release/semver.go"
        "release/metadata.go"
        "release/gitops.go"
        "release/github.go"
        "release/changelog.go"
        "release/workflow.go"
    )

    local missing=()
    for file in "${required_files[@]}"; do
        if [[ ! -f "$GITMAP_DIR/$file" ]]; then
            missing+=("$file")
        fi
    done

    if [[ ${#missing[@]} -gt 0 ]]; then
        write_fail "Missing source files (${#missing[@]}):"
        for f in "${missing[@]}"; do
            echo "  - $f"
        done
        exit 1
    fi

    write_success "All ${#required_files[@]} source files present"
}

# -- Build binary ----------------------------------------------
build_binary() {
    write_step "3/4" "Building $BINARY_NAME"
    test_source_files

    local bin_dir="$REPO_ROOT/$BUILD_OUTPUT"
    local out_path="$bin_dir/$BINARY_NAME"

    mkdir -p "$bin_dir"

    cd "$GITMAP_DIR"
    local abs_repo_root
    abs_repo_root=$(cd "$REPO_ROOT" && pwd)

    # Build-time identity injection — stamps the *source* repo metadata
    # into the binary so the `gitmap binary` footer block never falls
    # back to probing the user's CWD (fixed in v5.60.0).
    local build_commit build_branch build_repo build_date
    build_commit=$(cd "$abs_repo_root" && git rev-parse HEAD 2>/dev/null || echo "")
    build_branch=$(cd "$abs_repo_root" && git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
    build_repo=$(cd "$abs_repo_root" && git config --get remote.origin.url 2>/dev/null || echo "")
    build_date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    local ldflags="-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/constants.RepoPath=$abs_repo_root'"
    ldflags="$ldflags -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildCommit=$build_commit'"
    ldflags="$ldflags -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildBranch=$build_branch'"
    ldflags="$ldflags -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildRepo=$build_repo'"
    ldflags="$ldflags -X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildDate=$build_date'"

    # Pre-build provenance stamp — prints commit SHA, branch, declared
    # version, and a fingerprint of the historically-problematic cmd/
    # files so a stale checkout is obvious in the build log before
    # `go build` runs. Non-fatal: stamp failures never block the build.
    if [ -f "$REPO_ROOT/scripts/build-stamp.sh" ]; then
        bash "$REPO_ROOT/scripts/build-stamp.sh" || true
    fi

    if ! go build -ldflags "$ldflags" -o "$out_path" . 2>&1; then
        write_fail "Go build failed"
        exit 1
    fi

    if [[ "$COPY_DATA" == "true" ]]; then
        copy_data_folder "$bin_dir"
    fi

    local size
    if [[ "$(uname -s)" == "Darwin" ]]; then
        size=$(stat -f%z "$out_path" 2>/dev/null || echo "0")
    else
        size=$(stat -c%s "$out_path" 2>/dev/null || echo "0")
    fi
    local size_mb
    size_mb=$(echo "scale=2; $size / 1048576" | bc 2>/dev/null || echo "?")
    write_success "Binary built (${size_mb} MB) -> $out_path"

    BINARY_PATH="$out_path"
}

# -- Copy data folder -----------------------------------------
copy_data_folder() {
    local bin_dir="$1"
    local data_source="$GITMAP_DIR/data"
    local data_dest="$bin_dir/data"

    if [[ -d "$data_source" ]]; then
        rm -rf "$data_dest"
        cp -r "$data_source" "$data_dest"
        write_info "Copied data folder to bin"
    fi
}

# -- Copy docs-site to deploy directory -----------------------
# Required for `gitmap help-dashboard` (hd) which resolves docs-site/
# relative to the binary directory. Without this, `gitmap hd` fails with:
#   "Docs site directory not found at <deploy>/docs-site"
# -- Copy docs-site to deploy directory -----------------------
# Required for `gitmap help-dashboard` (hd) which resolves docs-site/
# relative to the binary directory. Without this, `gitmap hd` fails with:
#   "Docs site directory not found at <deploy>/docs-site"
#
# Source resolution order (first hit wins):
#   1. <repo>/docs-site/dist/   — legacy layout with a dedicated subdir
#   2. <repo>/dist/             — current layout where the repo root IS the
#                                 Vite docs app (no docs-site/ subdir)
#   3. Auto-build at repo root  — if package.json has a `build` script and
#                                 npm is on PATH, run it and use <repo>/dist/.
#   4. <repo>/docs-site/ source — npm-dev fallback (no prebuilt dist).
#   5. Warn — `gitmap hd` will fail until docs are built.
copy_docs_site() {
    local app_dir="$1"
    local docs_dest="$app_dir/docs-site"
    local legacy_dir="$REPO_ROOT/docs-site"
    local legacy_dist="$legacy_dir/dist"
    local root_dist="$REPO_ROOT/dist"
    local root_pkg="$REPO_ROOT/package.json"
    local gitmap_main="$GITMAP_DIR/main.go"
    local node_modules="$REPO_ROOT/node_modules"

    # Repo-detect diagnostics (active under --debug-repo-detect or env var).
    write_repo_detect "RepoRoot"          "$REPO_ROOT"
    write_repo_detect "GitMapDir"         "$GITMAP_DIR"
    write_repo_detect "gitmap/main.go"    "$([[ -f "$gitmap_main" ]] && echo present || echo missing)" "$gitmap_main"
    write_repo_detect "package.json"      "$([[ -f "$root_pkg" ]]    && echo present || echo missing)" "$root_pkg"
    write_repo_detect "node_modules/"     "$([[ -d "$node_modules" ]] && echo present || echo missing)"
    write_repo_detect "docs-site/dist/"   "$([[ -d "$legacy_dist" ]]  && echo present || echo missing)"
    write_repo_detect "dist/ (root)"      "$([[ -d "$root_dist" ]]    && echo present || echo missing)"
    local npm_path
    npm_path="$(command -v npm 2>/dev/null || true)"
    write_repo_detect "npm on PATH"       "$([[ -n "$npm_path" ]] && echo yes || echo no)" "$npm_path"
    write_repo_detect_snippet "package.json (first 6 lines)" "$root_pkg"

    # 1. Legacy <repo>/docs-site/dist/
    if [[ -d "$legacy_dist" ]]; then
        write_repo_detect "decision" "use-prebuilt-legacy" "$legacy_dist"
        local dist_dest="$docs_dest/dist"
        rm -rf "$dist_dest"
        mkdir -p "$docs_dest"
        cp -r "$legacy_dist" "$dist_dest"
        write_info "Copied docs-site/dist to gitmap app directory"
        return
    fi

    # 2. Current <repo>/dist/ (root-level Vite app)
    if [[ -d "$root_dist" ]]; then
        write_repo_detect "decision" "use-prebuilt-root" "$root_dist"
        local dist_dest="$docs_dest/dist"
        rm -rf "$dist_dest"
        mkdir -p "$docs_dest"
        cp -r "$root_dist" "$dist_dest"
        write_info "Copied root dist/ to gitmap app docs-site/dist"
        return
    fi

    # 3. Auto-build the root Vite app if package.json + npm available
    if [[ -f "$root_pkg" ]] && [[ -n "$npm_path" ]]; then
        local has_build="missing" has_vite="missing"
        grep -q '"build"' "$root_pkg" && has_build="found"
        grep -q '"vite"'  "$root_pkg" && has_vite="found"
        write_repo_detect "package.json:build" "$has_build"
        write_repo_detect "package.json:vite"  "$has_vite"
        if [[ "$has_build" == "found" ]]; then
            write_repo_detect "decision" "auto-build" "npm run build at $REPO_ROOT"
            if [[ ! -d "$node_modules" ]] || [[ ! -x "$node_modules/.bin/vite" ]]; then
                write_info "Installing docs dependencies (npm install) at repo root..."
                local install_exit=0
                (cd "$REPO_ROOT" && npm install --no-audit --no-fund --silent >/dev/null 2>&1) || install_exit=$?
                if [[ $install_exit -ne 0 ]]; then
                    write_warn "npm install failed - skipping docs build"
                    write_report_error "docs-npm-install" \
                        "npm install --no-audit --no-fund --silent" \
                        "$install_exit" \
                        "npm install failed at repo root; docs build skipped" \
                        "{\"repoRoot\":\"${REPO_ROOT//\"/\\\"}\",\"packageJson\":\"${root_pkg//\"/\\\"}\"}"
                    return
                fi
            fi
            write_info "Auto-building docs (npm run build) at repo root..."
            local build_exit=0
            (cd "$REPO_ROOT" && npm run build >/dev/null 2>&1) || build_exit=$?
            if [[ $build_exit -eq 0 ]] && [[ -d "$root_dist" ]]; then
                local dist_dest="$docs_dest/dist"
                rm -rf "$dist_dest"
                mkdir -p "$docs_dest"
                cp -r "$root_dist" "$dist_dest"
                write_info "Built and copied docs to gitmap app docs-site/dist"
                return
            fi
            write_warn "Auto-build failed - 'gitmap hd' will fail"
            write_report_error "docs-npm-build" \
                "npm run build" \
                "$build_exit" \
                "npm run build did not produce dist/ output" \
                "{\"repoRoot\":\"${REPO_ROOT//\"/\\\"}\",\"expectedDist\":\"${root_dist//\"/\\\"}\",\"packageJson\":\"${root_pkg//\"/\\\"}\"}"
            return
        else
            write_repo_detect "decision" "skip-no-build-script" "package.json has no \"build\" entry"
        fi
    else
        local skip_reason="unknown"
        [[ ! -f "$root_pkg" ]] && skip_reason="no package.json"
        [[ -z "$npm_path" ]]   && skip_reason="npm not on PATH"
        write_repo_detect "decision" "skip-not-a-vite-repo" "$skip_reason"
    fi

    # 4. Legacy <repo>/docs-site/ source-only — npm-dev fallback
    if [[ -d "$legacy_dir" ]]; then
        write_repo_detect "decision" "use-legacy-source" "$legacy_dir"
        rm -rf "$docs_dest"
        mkdir -p "$docs_dest"
        (cd "$legacy_dir" && find . -mindepth 1 -maxdepth 1 ! -name 'node_modules' -exec cp -r {} "$docs_dest/" \;)
        write_warn "No prebuilt dist/ found - copied docs-site/ source only (run 'npm run build' for static mode)"
        return
    fi

    # 5. Nothing found
    write_repo_detect "decision" "no-docs-source"
    write_warn "No docs found (checked docs-site/dist, docs-site/, dist/) - 'gitmap hd' will fail"
}

# -- Resolve deploy target -------------------------------------
# Priority: 1) --deploy-path flag  2) globally installed gitmap location  3) powershell.json default
resolve_deploy_target() {
    # 1) Explicit CLI override always wins
    if [[ -n "$DEPLOY_PATH" ]]; then
        write_info "Deploy target: CLI override -> $DEPLOY_PATH"
        echo "$DEPLOY_PATH"

        return
    fi

    # 2) If gitmap is already on PATH, deploy to its parent directory
    local active_cmd
    active_cmd=$(command -v gitmap 2>/dev/null || true)
    if [[ -n "$active_cmd" ]] && [[ -f "$active_cmd" ]]; then
        local resolved_active
        resolved_active=$(readlink -f "$active_cmd" 2>/dev/null || echo "$active_cmd")
        local active_dir
        active_dir=$(dirname "$resolved_active")
        local active_dir_name
        active_dir_name=$(basename "$active_dir")

        # The binary lives in <deploy-target>/$APP_SUBDIR/gitmap (or any
        # legacy folder name from LEGACY_APP_SUBDIRS). Either way the deploy
        # target is the parent of that wrapped folder. Folder names are
        # sourced from gitmap/constants/deploy-manifest.json.
        if is_known_app_subdir "$active_dir_name"; then
            local deploy_target
            deploy_target=$(dirname "$active_dir")
            write_info "Deploy target: detected from PATH -> $deploy_target"
            echo "$deploy_target"

            return
        fi

        # Binary is directly in a folder (not nested under gitmap-cli/)
        local deploy_target
        deploy_target=$(dirname "$active_dir")
        write_info "Deploy target: detected from PATH -> $deploy_target"
        echo "$deploy_target"

        return
    fi

    # 3) Fall back to powershell.json default
    write_info "Deploy target: powershell.json default -> $DEPLOY_TARGET"
    echo "$DEPLOY_TARGET"
}

# -- Repair legacy unwrapped/wrapped layout (DFD-3) ------------
# Migrates two legacy layouts into the canonical <target>/gitmap-cli/:
#   1) Unwrapped (pre-DFD): <target>/<binary> at top level.
#   2) v3.6.0..v3.13.10 wrapped: <target>/gitmap/ folder (renamed to
#      gitmap-cli on Unix in v3.13.11 for parity with run.ps1, which
#      did the same rename in v3.6.0).
# Idempotent — re-running on a correct gitmap-cli/ layout is a no-op.
repair_deploy_layout() {
    local target="$1"
    local app_dir="$target/$APP_SUBDIR"
    local legacy_binary="$target/$BINARY_NAME"
    local wrapped_binary="$app_dir/$BINARY_NAME"

    # --- Migration 2: any legacy app folder -> $APP_SUBDIR ----
    # Run BEFORE the unwrapped check so a stray top-level binary inside
    # an otherwise-correct legacy install gets folded in correctly.
    local legacy
    for legacy in "${LEGACY_APP_SUBDIRS[@]}"; do
        local legacy_app_dir="$target/$legacy"
        [[ "$legacy_app_dir" == "$app_dir" ]] && continue
        [[ ! -d "$legacy_app_dir" ]] && continue
        if [[ -d "$app_dir" ]]; then
            write_warn "Layout: both $legacy/ and $APP_SUBDIR/ exist at $target — leaving legacy $legacy/ for manual review"
        else
            if mv "$legacy_app_dir" "$app_dir" 2>/dev/null; then
                write_info "Layout: migrated legacy $legacy/ -> $APP_SUBDIR/ at $target"
            else
                write_warn "Layout: could not rename $legacy_app_dir -> $app_dir"
            fi
        fi
    done

    # --- Migration 1: legacy unwrapped binary -> $APP_SUBDIR/ --
    if [[ ! -f "$legacy_binary" ]]; then
        write_info "Layout: OK (no legacy binary at $target)"
        return
    fi

    if [[ -f "$wrapped_binary" ]]; then
        if rm -f "$legacy_binary" 2>/dev/null; then
            write_info "Layout: removed leftover legacy binary at $legacy_binary"
        else
            write_warn "Layout: could not remove legacy binary $legacy_binary"
        fi
        return
    fi

    write_info "Layout: migrating legacy unwrapped install -> $app_dir"
    mkdir -p "$app_dir"

    local name src dst
    for name in "$BINARY_NAME" data CHANGELOG.md docs docs-site; do
        src="$target/$name"
        dst="$app_dir/$name"
        [[ ! -e "$src" ]] && continue
        if [[ -e "$dst" ]]; then
            write_info "Layout: $name already inside $APP_SUBDIR/, skipping move"
            continue
        fi
        if mv "$src" "$dst" 2>/dev/null; then
            write_info "Layout: moved $name -> $APP_SUBDIR/$name"
        else
            write_warn "Layout: could not move $name"
        fi
    done
}

# -- Pre-deploy cleanup (DFD-6) --------------------------------
# Removes prior-deploy artifacts before the new binary is copied:
#   *.old, <binary>-update-*, updater-tmp-*, /tmp/<binary>-update-*,
#   *.gitmap-tmp-* swap directories. Logs every removal; never aborts.
invoke_deploy_cleanup() {
    local target="$1"
    local app_dir="$2"
    local stem="${BINARY_NAME%.exe}"
    local removed=0
    local dir pat f

    local scan_dirs=()
    [[ -d "$target" ]]  && scan_dirs+=("$target")
    [[ -d "$app_dir" ]] && [[ "$app_dir" != "$target" ]] && scan_dirs+=("$app_dir")

    local patterns=("*.old" "${stem}-update-*" "updater-tmp-*")
    for dir in "${scan_dirs[@]}"; do
        for pat in "${patterns[@]}"; do
            for f in "$dir"/$pat; do
                [[ ! -e "$f" ]] && continue
                if rm -rf "$f" 2>/dev/null; then
                    write_info "[cleanup] removed $f"
                    removed=$((removed + 1))
                else
                    write_warn "[cleanup] could not remove $f"
                fi
            done
        done
    done

    # Temp-dir scripts: $TMPDIR or /tmp
    local tmp_root="${TMPDIR:-/tmp}"
    if [[ -d "$tmp_root" ]]; then
        for f in "$tmp_root/${stem}-update-"*.sh "$tmp_root/${stem}-update-"*; do
            [[ ! -e "$f" ]] && continue
            if rm -rf "$f" 2>/dev/null; then
                write_info "[cleanup] removed temp script $f"
                removed=$((removed + 1))
            fi
        done
    fi

    # *.gitmap-tmp-* swap dirs left by interrupted clones
    if [[ -d "$target" ]]; then
        for f in "$target"/*.gitmap-tmp-*; do
            [[ ! -d "$f" ]] && continue
            if rm -rf "$f" 2>/dev/null; then
                write_info "[cleanup] removed swap dir $f"
                removed=$((removed + 1))
            else
                write_warn "[cleanup] could not remove $f"
            fi
        done
    fi

    if [[ $removed -gt 0 ]]; then
        write_success "[cleanup] removed $removed artifact(s)"
    else
        write_info "[cleanup] nothing to clean"
    fi
}

# -- Register on user PATH + refresh current session (DFD-4/5) -
# Writes export line into shell rc files using the
# 21-post-install-shell-activation marker block. Updates $PATH in the
# current process and prints the reload one-liner.
register_on_path() {
    local app_dir="$1"
    if [[ ! -d "$app_dir" ]]; then
        write_warn "PATH: skipping (app dir does not exist: $app_dir)"
        return
    fi

    local resolved
    resolved=$(cd "$app_dir" && pwd)

    case ":${PATH}:" in
        *":${resolved}:"*)
            write_info "PATH: already in current session"
            ;;
        *)
            export PATH="$PATH:$resolved"
            write_info "PATH: appended to current session -> $resolved"
            ;;
    esac

    local shell_name
    shell_name=$(basename "${SHELL:-/bin/bash}")

    local profiles=()
    case "$shell_name" in
        zsh)  profiles=("$HOME/.zshrc" "$HOME/.zprofile") ;;
        bash) profiles=("$HOME/.bashrc" "$HOME/.bash_profile") ;;
        fish) profiles=("$HOME/.config/fish/config.fish") ;;
        *)    profiles=("$HOME/.profile") ;;
    esac

    local marker_open="# gitmap shell wrapper v2 - managed by run.sh. Do not edit manually."
    local marker_close="# gitmap shell wrapper v2 end"

    # Single-source-of-truth: ask the freshly-built gitmap binary for the
    # canonical snippet bytes. Falls back to an inline heredoc only if
    # the binary is unreachable (first-run before deploy completed).
    local snippet=""
    local snippet_shell="$shell_name"
    case "$snippet_shell" in bash|zsh|fish) ;; *) snippet_shell="bash" ;; esac
    if [[ -n "${BINARY_PATH:-}" ]] && [[ -x "${BINARY_PATH}" ]]; then
        snippet="$("${BINARY_PATH}" setup print-path-snippet \
            --shell "$snippet_shell" --dir "$resolved" --manager "run.sh" 2>/dev/null || true)"
    fi
    if [[ -z "$snippet" ]]; then
        if [[ "$shell_name" == "fish" ]]; then
            snippet="${marker_open}
set -gx GITMAP_WRAPPER 1
fish_add_path ${resolved}
${marker_close}"
        else
            snippet="${marker_open}
export GITMAP_WRAPPER=1
case \":\${PATH}:\" in *\":${resolved}:\"*) ;; *) export PATH=\"\$PATH:${resolved}\" ;; esac
${marker_close}"
        fi
    fi

    local profile written=0
    for profile in "${profiles[@]}"; do
        [[ -z "$profile" ]] && continue
        mkdir -p "$(dirname "$profile")"
        touch "$profile"
        if grep -qF "$marker_open" "$profile" 2>/dev/null; then
            # Rewrite existing block (sed-based, portable enough for bash/zsh/fish profiles)
            local tmp
            tmp=$(mktemp)
            awk -v open="$marker_open" -v close="$marker_close" -v body="$snippet" '
                $0 == open { skip = 1; print body; next }
                skip && $0 == close { skip = 0; next }
                !skip { print }
            ' "$profile" > "$tmp" && mv "$tmp" "$profile"
            write_info "PATH: refreshed marker block in $profile"
        else
            printf '\n%s\n' "$snippet" >> "$profile"
            write_info "PATH: appended marker block to $profile"
        fi
        written=$((written + 1))
    done

    local primary="${profiles[0]}"
    local reload_cmd="source $primary"
    [[ "$shell_name" != "fish" ]] && reload_cmd=". $primary"
    write_success "PATH: persisted to $written profile(s)"
    write_info "PATH: to activate now in this shell, run: $reload_cmd"
}

# -- Deploy to target directory --------------------------------
deploy_binary() {
    write_step "4/4" "Deploying"

    local target
    target=$(resolve_deploy_target)

    write_info "Target: $target"
    mkdir -p "$target"

    # Migrate legacy unwrapped or older wrapped layouts into the canonical
    # $APP_SUBDIR/ (DFD-3) BEFORE we resolve $app_dir. Folder names come
    # from gitmap/constants/deploy-manifest.json (single source of truth).
    repair_deploy_layout "$target"

    local app_dir="$target/$APP_SUBDIR"
    mkdir -p "$app_dir"

    # Pre-deploy cleanup (DFD-6) — runs BEFORE the new binary is copied.
    invoke_deploy_cleanup "$target" "$app_dir"

    local dest_file="$app_dir/$BINARY_NAME"
    local backup_file="${dest_file}.old"
    local has_backup=false
    local deploy_success=false

    if [[ -f "$dest_file" ]]; then
        # Rename-first strategy: move the existing binary out of the way
        # so the copy target is free (avoids "text file busy" on some systems)
        rm -f "$backup_file" 2>/dev/null || true
        if mv "$dest_file" "$backup_file" 2>/dev/null; then
            has_backup=true
            write_info "Renamed existing binary to ${BINARY_NAME}.old (rename-first)"
        else
            # Fallback: copy-based backup
            if cp "$dest_file" "$backup_file" 2>/dev/null; then
                has_backup=true
                write_info "Backed up existing binary to ${BINARY_NAME}.old"
            else
                write_warn "Could not create backup"
            fi
        fi
    fi

    # Copy new binary — after rename-first, the destination is free
    local max_attempts=5
    local attempt=1
    while [[ $attempt -le $max_attempts ]]; do
        if cp "$BINARY_PATH" "$dest_file" 2>/dev/null; then
            deploy_success=true
            break
        fi
        write_warn "Target still locked; retrying ($attempt/$max_attempts)..."
        sleep 1
        attempt=$((attempt + 1))
    done

    if [[ "$deploy_success" == "false" ]]; then
        if [[ "$has_backup" == "true" ]] && [[ -f "$backup_file" ]] && [[ ! -f "$dest_file" ]]; then
            write_warn "Deploy failed - restoring previous binary from backup"
            mv "$backup_file" "$dest_file" 2>/dev/null && write_success "Rollback complete" || write_fail "Rollback also failed"
        fi
        write_fail "Deploy failed after $max_attempts attempts"
        exit 1
    fi

    if [[ "$has_backup" == "true" ]] && [[ "$deploy_success" == "true" ]]; then
        write_info "Previous binary kept as ${BINARY_NAME}.old (run 'gitmap update-cleanup' to remove)"
    fi

    # Copy data folder to deploy dir
    local bin_dir
    bin_dir=$(dirname "$BINARY_PATH")
    local data_dir="$bin_dir/data"
    local data_dest="$app_dir/data"
    if [[ -d "$data_dir" ]]; then
        rm -rf "$data_dest"
        cp -r "$data_dir" "$data_dest"
        write_info "Copied data folder to gitmap app directory"
    fi

    copy_docs_site "$app_dir"

    write_success "Deployed to $app_dir"

    # Register on user PATH and refresh current session (DFD-4, DFD-5).
    register_on_path "$app_dir"

    DEPLOYED_BINARY_PATH="$dest_file"
}

# -- Run gitmap ------------------------------------------------
invoke_run() {
    echo ""
    write_step "RUN" "Executing gitmap"

    local args=("${RUN_ARGS[@]}")

    # Default to scanning parent folder if no args
    if [[ ${#args[@]} -eq 0 ]]; then
        local parent_dir
        parent_dir=$(dirname "$REPO_ROOT")
        write_info "No args provided, defaulting to: scan $parent_dir"
        args=("scan" "$parent_dir")
    fi

    local arg_string="${args[*]}"
    write_info "Binary: $BINARY_PATH"
    write_info "Runner CWD: $(pwd)"
    write_info "Command: gitmap $arg_string"
    echo "  --------------------------------------------------"
    echo ""

    "$BINARY_PATH" "${args[@]}"
    local exit_code=$?

    echo ""
    if [[ $exit_code -eq 0 ]]; then
        write_success "Run complete"
    else
        write_fail "gitmap exited with code $exit_code"
    fi
}

# -- Run tests -------------------------------------------------
invoke_tests() {
    write_step "TEST" "Running unit tests"

    local report_dir="$GITMAP_DIR/data/unit-test-reports"
    mkdir -p "$report_dir"

    local overall_log="$report_dir/overall.log.txt"
    local failing_log="$report_dir/failingTest.log.txt"

    cd "$GITMAP_DIR"
    write_info "Running: go test ./..."

    local test_output
    test_output=$(go test ./... -v -count=1 2>&1) || true
    local test_exit=$?

    echo "$test_output" > "$overall_log"
    write_info "Overall report: $overall_log"

    # Extract failing tests
    local fail_lines
    fail_lines=$(echo "$test_output" | grep -E "^--- FAIL:|^FAIL\s" || true)

    if [[ -n "$fail_lines" ]]; then
        echo "$fail_lines" > "$failing_log"
        write_fail "Some tests failed. See: $failing_log"
    else
        echo "No failing tests." > "$failing_log"
        write_success "All tests passed"
    fi

    # Print summary counts
    local pass_count fail_count skip_count
    pass_count=$(echo "$test_output" | grep -c "^--- PASS:" || true)
    fail_count=$(echo "$test_output" | grep -c "^--- FAIL:" || true)
    skip_count=$(echo "$test_output" | grep -c "^--- SKIP:" || true)
    write_info "Results: $pass_count passed, $fail_count failed, $skip_count skipped"

    # Show test output
    while IFS= read -r line; do
        if [[ "$line" =~ ^---\ FAIL: ]]; then
            echo -e "  ${RED}$line${NC}"
        elif [[ "$line" =~ ^---\ PASS: ]]; then
            echo -e "  ${GREEN}$line${NC}"
        elif [[ "$line" =~ ^FAIL ]]; then
            echo -e "  ${RED}$line${NC}"
        elif [[ "$line" =~ ^ok\  ]]; then
            echo -e "  ${GREEN}$line${NC}"
        elif [[ -n "$line" ]]; then
            echo -e "  ${GRAY}$line${NC}"
        fi
    done <<< "$test_output"

    if [[ $test_exit -ne 0 ]]; then
        write_fail "Tests failed (exit code $test_exit)"
    fi
}

# -- Main ------------------------------------------------------
BINARY_PATH=""
DEPLOYED_BINARY_PATH=""

show_banner
load_config

if [[ "$TEST" == "true" ]]; then
    write_info "Test mode enabled (-t)"
    resolve_dependencies
    invoke_tests
    echo ""
    write_success "All done!"
    echo ""
    exit 0
fi

if [[ "$UPDATE" == "true" ]]; then
    write_info "Update mode enabled (--update)"
fi

if [[ "$NO_PULL" == "false" ]]; then
    invoke_git_pull
else
    write_info "Skipping git pull (--no-pull)"
fi

resolve_dependencies
build_binary

# Show built version
version_output=$("$BINARY_PATH" version 2>&1 || true)
write_info "Version: $version_output"

if [[ "$NO_DEPLOY" == "false" ]]; then
    deploy_binary

    # Check if PATH points to a different gitmap
    active_cmd=$(command -v gitmap 2>/dev/null || true)
    if [[ -n "$active_cmd" ]] && [[ -n "$DEPLOYED_BINARY_PATH" ]] && [[ -f "$DEPLOYED_BINARY_PATH" ]]; then
        active_resolved=$(readlink -f "$active_cmd" 2>/dev/null || echo "$active_cmd")
        deployed_resolved=$(readlink -f "$DEPLOYED_BINARY_PATH" 2>/dev/null || echo "$DEPLOYED_BINARY_PATH")
        if [[ "$active_resolved" != "$deployed_resolved" ]]; then
            write_warn "PATH points to a different gitmap binary."
            write_info "Active:   $active_resolved"
            write_info "Deployed: $deployed_resolved"
            if cp "$DEPLOYED_BINARY_PATH" "$active_cmd" 2>/dev/null; then
                synced_version=$("$active_cmd" version 2>&1 || true)
                write_success "Synced active PATH binary -> $synced_version"
            else
                write_warn "Could not sync active PATH binary. Run manually:"
                write_info "cp \"$DEPLOYED_BINARY_PATH\" \"$active_cmd\""
            fi
        fi
    fi
else
    write_info "Skipping deploy (--no-deploy)"
fi

# Show latest changelog
changelog_binary="$BINARY_PATH"
active_for_changelog=$(command -v gitmap 2>/dev/null || true)
if [[ -n "$active_for_changelog" ]] && [[ -f "$active_for_changelog" ]]; then
    changelog_binary="$active_for_changelog"
elif [[ -n "$DEPLOYED_BINARY_PATH" ]] && [[ -f "$DEPLOYED_BINARY_PATH" ]]; then
    changelog_binary="$DEPLOYED_BINARY_PATH"
fi

if [[ -f "$changelog_binary" ]]; then
    echo ""
    write_info "Latest changelog:"
    "$changelog_binary" changelog --latest || true

    if [[ "$UPDATE" == "true" ]]; then
        echo ""
        write_info "Running update cleanup"
        "$changelog_binary" update-cleanup || true
    fi
fi

if [[ "$RUN" == "true" ]]; then
    invoke_run
fi

echo ""
write_success "All done!"
echo ""
