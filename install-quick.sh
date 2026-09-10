#!/usr/bin/env bash
# Short interactive installer for gitmap on Linux / macOS.
#
# DUAL-MODE EXECUTION (spec/01-app/108-install-quick-auto-source.md):
#
#   1. Eval-mode (RECOMMENDED — auto-activates PATH in current shell):
#        eval "$(curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh)"
#
#      Because `eval` runs the script directly inside the user's interactive
#      shell instead of a child bash process, the trailing `source <profile>`
#      mutates the *current* PATH. No "open a new terminal" step needed.
#
#   2. Pipe-mode (LEGACY — child process, prints reload hint):
#        curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh | bash
#
#      Works exactly as before. The script detects it's running in a child
#      bash and prints a loud "Run this now: source <profile>" banner at the
#      end instead of trying to source (which would be a no-op).
#
#   3. Local-file mode:
#        ./install-quick.sh
#        ./install-quick.sh --interactive
#        ./install-quick.sh --dir /opt/gitmap
#        ./install-quick.sh --no-discovery
#        ./install-quick.sh --probe-ceiling 50
#
# Versioned repo discovery: if the source repo URL ends with -v<N>, this
# script probes for higher-numbered sibling repos (-v<N+1>, -v<N+2>, ...)
# and delegates to the latest available one. See:
#   spec/01-app/95-installer-script-find-latest-repo.md

# NOTE: We deliberately DO NOT `set -euo pipefail` at the top level when the
# script is being eval'd inside the user's interactive shell — those options
# would alter the user's shell options for the rest of the session. We set
# them only inside the main function via a subshell so the user's shell
# state is never modified.

REPO="alimtvnetwork/gitmap-v28"
INSTALLER_URL="https://raw.githubusercontent.com/${REPO}/main/gitmap/scripts/install.sh"
if [ "$(id -u 2>/dev/null || echo 1)" -eq 0 ]; then
    DEFAULT_DIR="/usr/local/bin"
else
    DEFAULT_DIR="${HOME:-~}/.local/bin"
fi

# ── Clean up corrupted installation directories ────────────────────────
# Removes accidental literal `~` directories and folders created by
# bash stdout-pollution bugs (e.g. named after ANSI escape banners).
cleanup_corrupted_install_dirs() {
    local scan_roots=()
    [ -n "${PWD:-}" ] && scan_roots+=("${PWD}")
    [ -n "${HOME:-}" ] && scan_roots+=("${HOME}")
    [ -d "${HOME:-}/.local" ] && scan_roots+=("${HOME}/.local" "${HOME}/.local/bin")
    [ -d "/tmp" ] && scan_roots+=("/tmp")

    # 1. Tier 1: Python 3 Scanner (Highest reliability against newlines/ANSI bytes)
    if command -v python3 >/dev/null 2>&1; then
        python3 - "${scan_roots[@]}" << 'PYEOF' 2>/dev/null || true
import os, sys, shutil

esc = chr(27)
for root in sys.argv[1:]:
    if not os.path.isdir(root):
        continue
    try:
        for name in os.listdir(root):
            p = os.path.join(root, name)
            if not os.path.isdir(p) or os.path.islink(p):
                continue
            if p in ("/", os.path.expanduser("~")):
                continue

            is_bad = False
            if name == "~":
                is_bad = True
            elif esc in name or "\x1b" in name or "\033" in name:
                is_bad = True
            elif "\n" in name or "\r" in name:
                is_bad = True
            elif "quick installer" in name.lower() or ("gitmap" in name.lower() and "installer" in name.lower()):
                is_bad = True
            elif "default:" in name.lower() or "choose install folder" in name.lower():
                is_bad = True

            if is_bad:
                shutil.rmtree(p, ignore_errors=True)
    except Exception:
        pass
PYEOF
        return 0
    fi

    # 2. Tier 2: Portable POSIX find Scanner
    local esc=$'\033' nl=$'\n' cr=$'\r'
    local r bad_dir
    for r in "${scan_roots[@]}"; do
        [ ! -d "$r" ] && continue
        if [ -d "$r/~" ] && [ "$r/~" != "/" ] && [ "$r/~" != "${HOME:-}" ]; then
            rm -rf "$r/~" 2>/dev/null || true
        fi
        find "$r" -mindepth 1 -maxdepth 1 -type d \( \
            -name "*gitmap*installer*" -o \
            -name "*quick installer*" -o \
            -name "*${esc}*" -o \
            -name "*${nl}*" -o \
            -name "*${cr}*" \
        \) 2>/dev/null | while IFS= read -r bad_dir; do
            if [ -n "$bad_dir" ] && [ "$bad_dir" != "/" ] && [ "$bad_dir" != "${HOME:-}" ]; then
                rm -rf "$bad_dir" 2>/dev/null || true
            fi
        done
    done
}

# ── Path sanitizer & tilde expansion ────────────────────────────────────
sanitize_install_dir() {
    local raw="$1"
    if [ -z "${raw}" ]; then
        echo ""
        return 0
    fi

    # Strip ANSI escape sequences: \033[...m or \x1b[...]
    local cleaned
    cleaned="$(printf '%s' "${raw}" | sed -E 's/\x1B\[[0-9;]*[a-zA-Z]//g')"

    # Strip carriage returns (\r)
    cleaned="$(printf '%s' "${cleaned}" | tr -d '\r')"

    # Strip leading/trailing whitespace and surrounding quotes
    cleaned="$(printf '%s' "${cleaned}" | sed -e 's/^[[:space:]"'"'"']*//' -e 's/[[:space:]"'"'"']*$//')"

    # If cleaned contains newlines or prompt text, it is corrupted output -> reject
    case "${cleaned}" in
        *$'\n'*|*$'\r'*|*"installer"*|*"Install path"*|*"Default:"*|*"Choose install folder"*)
            echo ""
            return 0
            ;;
    esac

    # Tilde expansion (bash inside double-quotes does not expand ~)
    if [ "${cleaned}" = "~" ]; then
        cleaned="${HOME:-/tmp}"
    elif [[ "${cleaned}" == "~/"* ]]; then
        cleaned="${HOME:-/tmp}/${cleaned#\~/}"
    fi

    # Strip trailing slash (unless root /)
    if [ "${cleaned}" != "/" ]; then
        cleaned="${cleaned%/}"
    fi

    echo "${cleaned}"
}

# Detect execution mode. EVAL_MODE=1 means the script is running in the
# user's interactive shell (via `eval` or `source`), so we CAN source the
# profile at the end and have it stick. Otherwise we print a manual hint.
__gitmap_detect_eval_mode() {
    # When invoked as `bash script.sh` or `curl | bash`, $0 is "bash" or
    # "/bin/bash" and we are in a child process. When eval'd, $0 is the
    # caller's shell (e.g. "-bash", "-zsh", "/usr/bin/zsh") and BASH_SOURCE
    # is empty (zsh) or equals $0 (bash eval).
    case "${0##*/}" in
        bash|sh|dash|ksh)
            # Could be `bash script.sh` (child) OR `bash -c '...eval...'` —
            # check BASH_SOURCE: in eval'd context BASH_SOURCE[0] is empty.
            if [ -z "${BASH_SOURCE[0]:-}" ]; then
                echo 1
            else
                echo 0
            fi
            ;;
        *)
            # zsh, fish, login shell ("-bash", "-zsh") — eval mode.
            echo 1
            ;;
    esac
}

__GITMAP_EVAL_MODE="$(__gitmap_detect_eval_mode)"

# All the heavy lifting runs inside this function so we can use `local` and
# wrap it in a subshell-safe way without polluting the caller's shell.
__gitmap_quick_install_main() {
    # Subshell isolates `set -e` etc. from the user's interactive shell.
    (
        set -eu
        # Don't fail-fast on pipefail if eval'd — too easy to trip user state.

        local INSTALL_DIR=""
        local VERSION=""
        local NO_DISCOVERY=0
        local INTERACTIVE=0
        # PROBE_CEILING retained for backward compat (legacy fail-fast upper
        # bound). The canonical knob per spec/07-generic-release/09 §6 is
        # --discovery-window <K> (default 20, cap 20, or 50 if GITHUB_TOKEN
        # is set).
        local PROBE_CEILING=30
        local DISCOVERY_WINDOW=20

        while [ $# -gt 0 ]; do
            case "$1" in
                -i|--interactive)    INTERACTIVE=1;    shift ;;
                --dir)               INSTALL_DIR="$2"; shift 2 ;;
                --version)           VERSION="$2";     shift 2 ;;
                --no-discovery)      NO_DISCOVERY=1;   shift ;;
                --probe-ceiling)     PROBE_CEILING="$2"; shift 2 ;;
                --discovery-window)  DISCOVERY_WINDOW="$2"; shift 2 ;;
                -h|--help)
                    sed -n '2,31p' "${BASH_SOURCE[0]:-$0}" 2>/dev/null || \
                        printf '  See https://github.com/%s for usage.\n' "${REPO}"
                    return 0
                    ;;
                *)
                    printf '  Unknown argument: %s\n' "$1" >&2
                    return 1
                    ;;
            esac
        done

        # ── Versioned repo discovery (spec 95) ──────────────────────

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
            local url="$1"
            curl -sfI --max-time 5 "$url" >/dev/null 2>&1
        }

        # Per spec/07-generic-release/09-generic-install-script-behavior.md §4.1:
        # Probe -v<N+1>..-v<N+window> CONCURRENTLY (max 20, or 50 if
        # GITHUB_TOKEN is set). Pick max(M) where HEAD returned 200.
        # Gaps are tolerated (no fail-fast on first MISS).
        resolve_effective_repo() {
            local repo="$1" window="$2"

            if ! parse_repo_suffix "$repo"; then
                printf '  [discovery] no -v<N> suffix on '"'"'%s'"'"'; installing baseline as-is\n' "$repo" >&2
                echo "$repo"
                return 0
            fi

            local owner="$SUFFIX_OWNER" stem="$SUFFIX_STEM" baseline="$SUFFIX_N"

            # Concurrency cap: 20 anonymous, 50 if GITHUB_TOKEN supplied.
            local max_concurrency=20
            if [ -n "${GITHUB_TOKEN:-}" ]; then
                max_concurrency=50
            fi
            if [ "$window" -gt "$max_concurrency" ]; then
                window="$max_concurrency"
            fi

            printf '  [discovery] baseline: %s/%s-v%s\n' "$owner" "$stem" "$baseline" >&2
            printf '  [discovery] window: %d (parallel HEAD, max-hit-wins, gap-tolerant)\n' "$window" >&2

            local tmpdir
            tmpdir="$(mktemp -d 2>/dev/null || mktemp -d -t gmprobe)"
            # shellcheck disable=SC2064
            trap "rm -rf '$tmpdir'" RETURN

            local m url pids=()
            for (( m = baseline + 1; m <= baseline + window; m++ )); do
                url="https://github.com/${owner}/${stem}-v${m}"
                # Each probe runs in background, writes its M to a hit-file
                # on success. 5s connect+total timeout per spec §4.1.
                (
                    if curl -sfI --max-time 5 --connect-timeout 5 "$url" >/dev/null 2>&1; then
                        printf '  [discovery] HEAD %s ... HIT\n' "$url" >&2
                        printf '%d\n' "$m" > "$tmpdir/hit.$m"
                    else
                        printf '  [discovery] HEAD %s ... MISS\n' "$url" >&2
                    fi
                ) &
                pids+=($!)
            done

            # Wait for ALL probes (no early exit — gaps allowed).
            local pid
            for pid in "${pids[@]}"; do
                wait "$pid" 2>/dev/null || true
            done

            # Pick max(M) across hit files.
            local effective="$baseline" hit_m
            for f in "$tmpdir"/hit.*; do
                [ -e "$f" ] || continue
                hit_m="$(cat "$f" 2>/dev/null || echo 0)"
                if [ "$hit_m" -gt "$effective" ]; then
                    effective="$hit_m"
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

        invoke_delegated_installer() {
            local effective_repo="$1"
            local delegated_url="https://raw.githubusercontent.com/${effective_repo}/main/install-quick.sh"

            printf '  [discovery] delegating to %s\n' "$delegated_url" >&2

            local pass_args=()
            [ "$INTERACTIVE" = "1" ] && pass_args+=(--interactive)
            [ -n "$INSTALL_DIR" ]    && pass_args+=(--dir "$INSTALL_DIR")
            [ -n "$VERSION" ]        && pass_args+=(--version "$VERSION")
            pass_args+=(--probe-ceiling "$PROBE_CEILING")
            pass_args+=(--discovery-window "$DISCOVERY_WINDOW")

            export INSTALLER_DELEGATED=1

            local script
            if ! script="$(curl -fsSL --max-time 15 "$delegated_url")"; then
                printf '  [discovery] [WARN] could not fetch delegated installer; falling back to baseline\n' >&2
                unset INSTALLER_DELEGATED
                return 1
            fi

            bash -c "$script" _ "${pass_args[@]}"
            return $?
        }

        if [ "${INSTALLER_DELEGATED:-0}" = "1" ]; then
            printf '  [discovery] INSTALLER_DELEGATED=1; skipping discovery (loop guard)\n' >&2
        elif [ "$NO_DISCOVERY" = "1" ]; then
            printf '  [discovery] --no-discovery set; skipping probe\n' >&2
        elif [ -n "$VERSION" ]; then
            # Strict-tag contract (spec/07-generic-release/09-generic-install-script-behavior.md §3):
            # An explicit --version pins the install to that exact release.
            # MUST NOT probe -v<N+i> sibling repos. MUST NOT call releases/latest.
            # MUST NOT fall back to main on failure. The canonical installer
            # downstream enforces the same contract on the asset download path.
            printf '  [strict] --version %s pinned; skipping repo probe (no fallback)\n' "$VERSION" >&2
        else
            EFFECTIVE_REPO="$(resolve_effective_repo "$REPO" "$DISCOVERY_WINDOW")"
            if [ "$EFFECTIVE_REPO" != "$REPO" ]; then
                if invoke_delegated_installer "$EFFECTIVE_REPO"; then
                    return 0
                fi
                printf '  [discovery] [WARN] delegation failed; falling back to baseline\n' >&2
            fi
        fi

        # ── Baseline install flow ────────────────────────────────────

        prompt_dir() {
            printf '\n' >&2
            printf '  \033[36mgitmap quick installer\033[0m\n' >&2
            printf '  \033[90m---------------------\033[0m\n' >&2
            printf '  Choose install folder. Press Enter to accept the default.\n' >&2
            printf '  \033[90mDefault: %s\033[0m\n' "${DEFAULT_DIR}" >&2
            printf '  Install path: ' >&2

            local answer=""
            if [ -r /dev/tty ]; then
                IFS= read -r answer < /dev/tty || answer=""
            elif [ -t 0 ]; then
                IFS= read -r answer || answer=""
            fi

            answer="$(sanitize_install_dir "${answer}")"
            if [ -z "${answer}" ]; then
                echo "${DEFAULT_DIR}"
            else
                echo "${answer}"
            fi
        }

        # Clean up any legacy corrupted folders right at the start
        cleanup_corrupted_install_dirs

        if [ -n "${INSTALL_DIR}" ]; then
            INSTALL_DIR="$(sanitize_install_dir "${INSTALL_DIR}")"
        fi

        if [ -z "${INSTALL_DIR}" ]; then
            if [ "${INTERACTIVE}" = "1" ] && { [ -t 0 ] || [ -r /dev/tty ]; }; then
                INSTALL_DIR="$(prompt_dir)"
            else
                INSTALL_DIR="${DEFAULT_DIR}"
                printf '  \033[90m[info] Using default install dir: %s (pass -i/--interactive to choose)\033[0m\n' "${INSTALL_DIR}" >&2
            fi
        fi

        [ -z "${INSTALL_DIR}" ] && INSTALL_DIR="${DEFAULT_DIR}"

        printf '\n  \033[32mInstalling gitmap to: %s\033[0m\n\n' "${INSTALL_DIR}"

        save_deploy_path() {
            local dir="$1"
            [ -z "${dir}" ] && return 1
            mkdir -p "${dir}" 2>/dev/null || true
            local cfg="${dir}/powershell.json"
            cat > "${cfg}" <<EOF
{
  "deployPath": "${dir}",
  "buildOutput": "./bin",
  "binaryName": "gitmap",
  "goSource": "./gitmap",
  "copyData": true
}
EOF
            printf '  \033[90mSaved deployPath -> %s\033[0m\n' "${cfg}"
        }

        # Dedicated temporary staging directory for quick installer
        local stage_dir
        stage_dir="$(mktemp -d 2>/dev/null || mktemp -d -t gmquick)"
        # shellcheck disable=SC2064
        trap "rm -rf '${stage_dir}' 2>/dev/null || true" EXIT RETURN INT TERM

        local staged_installer="${stage_dir}/install.sh"
        local installer_log="${stage_dir}/install.log"

        if ! curl -fsSL "${INSTALLER_URL}" -o "${staged_installer}"; then
            printf '  \033[31m[ERROR] Failed to download installer from %s\033[0m\n' "${INSTALLER_URL}" >&2
            rm -rf "${stage_dir}" 2>/dev/null || true
            return 1
        fi
        chmod +x "${staged_installer}" 2>/dev/null || true

        ARGS=(--dir "${INSTALL_DIR}")
        if [ -n "${VERSION}" ]; then
            ARGS+=(--version "${VERSION}")
        fi

        # Run the staged canonical installer in child bash
        if bash "${staged_installer}" "${ARGS[@]}" 2> >(tee "${installer_log}" >&2); then
            local hint profile
            hint="$(grep -oE '(source |\. )[^ ]+' "${installer_log}" | head -n1 || true)"
            profile="$(printf '%s' "${hint}" | awk '{print $NF}')"

            # Clean up temporary staging directory immediately
            rm -rf "${stage_dir}" 2>/dev/null || true

            # Save deployPath config in target directory
            save_deploy_path "${INSTALL_DIR}" || printf '  \033[33m[WARN] Could not save powershell.json\033[0m\n'

            # Persist for the outer driver (subshell can't export upward).
            mkdir -p "${INSTALL_DIR}" 2>/dev/null || true
            printf '%s\n' "${profile:-}" > "${INSTALL_DIR}/.gitmap-last-profile" 2>/dev/null || true
            printf '%s\n' "${INSTALL_DIR}" > "${HOME}/.gitmap-last-install-dir" 2>/dev/null || true
            return 0
        else
            local rc=$?
            rm -rf "${stage_dir}" 2>/dev/null || true
            return $rc
        fi
    )
}

# ── Outer driver: runs in the user's shell when eval'd ─────────────────

__gitmap_quick_install_main "$@"
__gitmap_rc=$?

for __gm_arg in "$@"; do
    if [ "${__gm_arg}" = "-h" ] || [ "${__gm_arg}" = "--help" ]; then
        return 0 2>/dev/null || exit 0
    fi
done

if [ "${__gitmap_rc}" -ne 0 ]; then
    printf '\n  \033[31m[ERROR]\033[0m install failed (exit %d)\n' "${__gitmap_rc}" >&2
    return "${__gitmap_rc}" 2>/dev/null || exit "${__gitmap_rc}"
fi

# Recover the install dir + reload-target written by the subshell.
__gitmap_install_dir="${DEFAULT_DIR}"
if [ -r "${HOME}/.gitmap-last-install-dir" ]; then
    __gitmap_install_dir="$(cat "${HOME}/.gitmap-last-install-dir" 2>/dev/null || true)"
    __gitmap_install_dir="$(sanitize_install_dir "${__gitmap_install_dir}")"
    [ -z "${__gitmap_install_dir}" ] && __gitmap_install_dir="${DEFAULT_DIR}"
fi

__gitmap_profile=""
if [ -r "${__gitmap_install_dir}/.gitmap-last-profile" ]; then
    __gitmap_profile="$(cat "${__gitmap_install_dir}/.gitmap-last-profile" 2>/dev/null || true)"
fi

# Best-effort fallback when install.sh's hint was not parseable: pick the
# profile matching the user's current shell.
if [ -z "${__gitmap_profile}" ]; then
    case "${SHELL##*/}" in
        zsh)  __gitmap_profile="${HOME}/.zshrc" ;;
        bash) __gitmap_profile="${HOME}/.bashrc" ;;
        *)    __gitmap_profile="${HOME}/.profile" ;;
    esac
fi

if [ "${__GITMAP_EVAL_MODE}" = "1" ] && [ -r "${__gitmap_profile}" ]; then
    # We are in the user's live shell — actually source the profile so PATH
    # is updated NOW and `gitmap` is callable on the very next line.
    printf '\n  \033[32mActivating gitmap in this shell:\033[0m source %s\n' "${__gitmap_profile}"
    # shellcheck disable=SC1090
    . "${__gitmap_profile}" || \
        printf '  \033[33m[WARN] source %s returned non-zero; PATH may need a manual reload\033[0m\n' "${__gitmap_profile}" >&2

    if command -v gitmap >/dev/null 2>&1; then
        printf '  \033[32mOK\033[0m  gitmap is on your PATH:  %s\n' "$(command -v gitmap)"
    else
        printf '  \033[33m[WARN]\033[0m gitmap is not yet on PATH. Open a new terminal or run:\n' >&2
        printf '      \033[36msource %s\033[0m\n' "${__gitmap_profile}" >&2
    fi
else
    # Pipe-mode (curl | bash) or no profile found — print a loud manual hint.
    printf '\n  \033[33m> To use gitmap NOW in this shell, run:\033[0m\n'
    printf '      \033[36msource %s\033[0m\n' "${__gitmap_profile}"
    printf '  \033[90m  (or open a new terminal -- POSIX child processes cannot mutate the parent shell)\033[0m\n'
    printf '\n  \033[90m  Tip: next time, install with eval-mode for auto-activation:\033[0m\n'
    printf '      \033[36meval "$(curl -fsSL https://raw.githubusercontent.com/%s/main/install-quick.sh)"\033[0m\n' "${REPO}"
fi

cleanup_corrupted_install_dirs

unset __gitmap_quick_install_main __gitmap_detect_eval_mode \
      __gitmap_rc __gitmap_install_dir __gitmap_profile __GITMAP_EVAL_MODE \
      DEFAULT_DIR sanitize_install_dir cleanup_corrupted_install_dirs 2>/dev/null || true
