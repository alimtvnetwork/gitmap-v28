#!/usr/bin/env bash
# Short interactive installer for gitmap on Linux / macOS.
#
# DUAL-MODE EXECUTION (spec/01-app/108-install-quick-auto-source.md):
#
#   1. Eval-mode (RECOMMENDED — auto-activates PATH in current shell):
#        eval "$(curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh)"
#
#      Because `eval` runs the script directly inside the user's interactive
#      shell instead of a child bash process, the trailing `source <profile>`
#      mutates the *current* PATH. No "open a new terminal" step needed.
#
#   2. Pipe-mode (LEGACY — child process, prints reload hint):
#        curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh | bash
#
#      Works exactly as before. The script detects it's running in a child
#      bash and prints a loud "Run this now: source <profile>" banner at the
#      end instead of trying to source (which would be a no-op).
#
#   3. Local-file mode:
#        ./install-quick.sh
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

REPO="alimtvnetwork/gitmap-v27"
INSTALLER_URL="https://raw.githubusercontent.com/${REPO}/main/gitmap/scripts/install.sh"
DEFAULT_DIR="${HOME}/.local/bin"

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
        # PROBE_CEILING retained for backward compat (legacy fail-fast upper
        # bound). The canonical knob per spec/07-generic-release/09 §6 is
        # --discovery-window <K> (default 20, cap 20, or 50 if GITHUB_TOKEN
        # is set).
        local PROBE_CEILING=30
        local DISCOVERY_WINDOW=20

        while [ $# -gt 0 ]; do
            case "$1" in
                --dir)               INSTALL_DIR="$2"; shift 2 ;;
                --version)           VERSION="$2";     shift 2 ;;
                --no-discovery)      NO_DISCOVERY=1;   shift ;;
                --probe-ceiling)     PROBE_CEILING="$2"; shift 2 ;;
                --discovery-window)  DISCOVERY_WINDOW="$2"; shift 2 ;;
                -h|--help)
                    sed -n '2,40p' "${BASH_SOURCE[0]:-$0}" 2>/dev/null || \
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
            [ -n "$INSTALL_DIR" ]   && pass_args+=(--dir "$INSTALL_DIR")
            [ -n "$VERSION" ]       && pass_args+=(--version "$VERSION")
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
            printf '\n'
            printf '  \033[36mgitmap quick installer\033[0m\n'
            printf '  \033[90m---------------------\033[0m\n'
            printf '  Choose install folder. Press Enter to accept the default.\n'
            printf '  \033[90mDefault: %s\033[0m\n' "${DEFAULT_DIR}"
            printf '  Install path: '

            if [ -r /dev/tty ]; then
                IFS= read -r answer < /dev/tty || answer=""
            else
                IFS= read -r answer || answer=""
            fi

            if [ -z "${answer}" ]; then
                echo "${DEFAULT_DIR}"
            else
                echo "${answer}"
            fi
        }

        if [ -z "${INSTALL_DIR}" ]; then
            INSTALL_DIR="$(prompt_dir)"
        fi

        printf '\n  \033[32mInstalling gitmap to: %s\033[0m\n\n' "${INSTALL_DIR}"

        save_deploy_path() {
            local dir="$1"
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

        save_deploy_path "${INSTALL_DIR}" || printf '  \033[33m[WARN] Could not save powershell.json\033[0m\n'

        ARGS=(--dir "${INSTALL_DIR}")
        if [ -n "${VERSION}" ]; then
            ARGS+=(--version "${VERSION}")
        fi

        # Run the canonical installer in a child bash. We CAPTURE its stderr
        # so we can extract the PATH_RELOAD hint and re-run the source
        # ourselves in eval-mode — install.sh prints something like:
        #     source ~/.zshrc   # in zsh
        local installer_log
        installer_log="$(mktemp)"
        if curl -fsSL "${INSTALLER_URL}" | bash -s -- "${ARGS[@]}" 2> >(tee "${installer_log}" >&2); then
            local hint profile
            hint="$(grep -oE '(source |\. )[^ ]+' "${installer_log}" | head -n1 || true)"
            profile="$(printf '%s' "${hint}" | awk '{print $NF}')"
            rm -f "${installer_log}"
            # Persist for the outer driver (subshell can't export upward).
            mkdir -p "${INSTALL_DIR}" 2>/dev/null || true
            printf '%s\n' "${profile:-}" > "${INSTALL_DIR}/.gitmap-last-profile" 2>/dev/null || true
            printf '%s\n' "${INSTALL_DIR}" > "${HOME}/.gitmap-last-install-dir" 2>/dev/null || true
            return 0
        else
            local rc=$?
            rm -f "${installer_log}"
            return $rc
        fi
    )
}

# ── Outer driver: runs in the user's shell when eval'd ─────────────────

__gitmap_quick_install_main "$@"
__gitmap_rc=$?

if [ "${__gitmap_rc}" -ne 0 ]; then
    printf '\n  \033[31m[ERROR]\033[0m install failed (exit %d)\n' "${__gitmap_rc}" >&2
    return "${__gitmap_rc}" 2>/dev/null || exit "${__gitmap_rc}"
fi

# Recover the install dir + reload-target written by the subshell.
__gitmap_install_dir="${HOME}/.local/bin"
if [ -r "${HOME}/.gitmap-last-install-dir" ]; then
    __gitmap_install_dir="$(cat "${HOME}/.gitmap-last-install-dir" 2>/dev/null || true)"
    [ -z "${__gitmap_install_dir}" ] && __gitmap_install_dir="${HOME}/.local/bin"
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

unset __gitmap_quick_install_main __gitmap_detect_eval_mode \
      __gitmap_rc __gitmap_install_dir __gitmap_profile __GITMAP_EVAL_MODE 2>/dev/null || true
