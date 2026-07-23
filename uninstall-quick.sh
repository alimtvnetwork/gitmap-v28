#!/usr/bin/env bash
# One-liner uninstaller for gitmap-v28 on Linux / macOS.
#
# Removes the gitmap-v28 binary, install folder, shell PATH entries, and
# (optionally) the per-user data folder. Works whether gitmap-v28 was installed
# via:
#
#   - install-quick.sh (one-liner)
#   - gitmap-v28/scripts/install.sh (canonical installer)
#   - manual `run.sh` build-and-deploy
#
# Strategy:
#   1. If `gitmap-v28` is on PATH, delegate to `gitmap-v28 self-uninstall -y`.
#      The binary knows how to clean its marker-block PATH entries from
#      ~/.bashrc, ~/.zshrc, ~/.profile, etc.
#   2. If `gitmap-v28` is NOT on PATH (already partially removed, broken
#      install), fall back to a manual sweep:
#        - delete <dir>/gitmap-cli AND legacy <dir>/gitmap-v28
#        - strip the install dir from rc files
#        - prompt before deleting ~/.config/gitmap-v28
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.sh | bash
#   ./uninstall-quick.sh
#   ./uninstall-quick.sh --keep-data
#   ./uninstall-quick.sh --dir /opt/gitmap-v28
#   ./uninstall-quick.sh --yes

set -uo pipefail

INSTALL_DIR=""
KEEP_DATA=0
ASSUME_YES=0

while [ $# -gt 0 ]; do
    case "$1" in
        --dir)        INSTALL_DIR="$2"; shift 2 ;;
        --keep-data)  KEEP_DATA=1;       shift   ;;
        -y|--yes)     ASSUME_YES=1;      shift   ;;
        -h|--help)
            sed -n '2,28p' "$0" 2>/dev/null || true
            exit 0
            ;;
        *)
            printf '  Unknown argument: %s\n' "$1" >&2
            exit 1
            ;;
    esac
done

c_cyan='\033[36m'; c_dgray='\033[90m'; c_green='\033[32m'
c_yellow='\033[33m'; c_red='\033[31m'; c_reset='\033[0m'

step() { printf '  %b%s%b\n' "$c_cyan"   "$1" "$c_reset"; }
info() { printf '    %b%s%b\n' "$c_dgray"  "$1" "$c_reset"; }
ok()   { printf '    %b%s%b\n' "$c_green"  "$1" "$c_reset"; }
warn() { printf '    %b%s%b\n' "$c_yellow" "$1" "$c_reset"; }
err()  { printf '    %b%s%b\n' "$c_red"    "$1" "$c_reset"; }

confirm() {
    if [ "$ASSUME_YES" = "1" ]; then return 0; fi
    local prompt="$1" answer
    printf '\n  %b%s [y/N]: %b' "$c_yellow" "$prompt" "$c_reset"
    if [ -r /dev/tty ]; then
        IFS= read -r answer < /dev/tty || answer=""
    else
        IFS= read -r answer || answer=""
    fi
    case "$answer" in y|Y|yes|YES) return 0 ;; *) return 1 ;; esac
}

# ---------------------------------------------------------------------------
# Step 1 — canonical self-uninstall via the gitmap-v28 binary itself.
# ---------------------------------------------------------------------------

try_self_uninstall() {
    if ! command -v gitmap-v28 >/dev/null 2>&1; then
        info "gitmap-v28 not found on PATH, skipping self-uninstall (will sweep manually)"
        return 1
    fi

    info "Active binary: $(command -v gitmap-v28)"
    info "Delegating to: gitmap-v28 self-uninstall -y"
    printf '\n'
    if gitmap-v28 self-uninstall -y; then
        ok "self-uninstall completed cleanly"
        return 0
    fi

    warn "self-uninstall exited non-zero; falling back to manual sweep"
    return 1
}

# ---------------------------------------------------------------------------
# Step 2 — manual sweep fallback.
# ---------------------------------------------------------------------------

resolve_install_dir() {
    if [ -n "$INSTALL_DIR" ]; then printf '%s\n' "$INSTALL_DIR"; return; fi

    # Active binary -> its parent.
    local active
    active="$(command -v gitmap-v28 2>/dev/null || true)"
    if [ -n "$active" ] && [ -x "$active" ]; then
        local resolved
        resolved="$(readlink -f "$active" 2>/dev/null || printf '%s' "$active")"
        local parent grand
        parent="$(dirname "$resolved")"
        grand="$(dirname "$parent")"

        # If parent is named gitmap-cli/gitmap-v28, return its parent (deploy root).
        case "$(basename "$parent")" in
            gitmap-cli|gitmap-v28) printf '%s\n' "$grand" ;;
            *)                 printf '%s\n' "$parent" ;;
        esac
        return
    fi

    # Common defaults.
    for d in "$HOME/.local/bin" "$HOME/bin" "/opt/gitmap-v28" "/usr/local/bin"; do
        if [ -x "$d/gitmap-v28" ] || [ -x "$d/gitmap-cli/gitmap-v28" ] || [ -x "$d/gitmap-v28/gitmap-v28" ]; then
            printf '%s\n' "$d"
            return
        fi
    done
    printf '\n'
}

remove_install_files() {
    local root="$1"
    if [ -z "$root" ]; then
        warn "could not locate a gitmap-v28 install dir; skipping file removal"
        return
    fi

    for sub in gitmap-cli gitmap-v28; do
        local dir="$root/$sub"
        if [ -d "$dir" ]; then
            if rm -rf "$dir"; then ok "removed $dir"
            else err "could not remove $dir (check permissions)"; fi
        fi
    done

    local flat="$root/gitmap-v28"
    if [ -f "$flat" ]; then
        if rm -f "$flat"; then ok "removed $flat"
        else err "could not remove $flat"; fi
    fi
}

clean_rc_files() {
    local root="$1"
    if [ -z "$root" ]; then return; fi

    local pattern_root
    pattern_root="$(printf '%s' "$root" | sed 's/[][\.*^$()+?{}|\/]/\\&/g')"

    for rc in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile" "$HOME/.bash_profile"; do
        [ -f "$rc" ] || continue
        if grep -qE "(${pattern_root}|gitmap-cli|gitmap-v28/gitmap-v28)" "$rc" 2>/dev/null; then
            local backup="${rc}.gitmap-uninstall.bak"
            cp "$rc" "$backup" 2>/dev/null || true
            # Strip lines that reference the install dir or the gitmap-v28 subfolders.
            sed -i.tmp "/${pattern_root}/d; /gitmap-cli/d; /gitmap-v28\/gitmap-v28/d" "$rc" 2>/dev/null \
                && rm -f "${rc}.tmp" \
                && ok "cleaned PATH entries from $rc (backup: $backup)"
        fi
    done
}

remove_data_folder() {
    local data="${XDG_CONFIG_HOME:-$HOME/.config}/gitmap-v28"
    [ -d "$data" ] || return

    if [ "$KEEP_DATA" = "1" ]; then
        info "keeping data folder: $data"
        return
    fi

    if confirm "Also delete user data at $data?"; then
        if rm -rf "$data"; then ok "removed $data"
        else err "could not remove $data"; fi
    else
        info "kept: $data"
    fi
}

# ---------------------------------------------------------------------------
# Exhaustive PATH sweep — find EVERY `gitmap-v28` still on PATH and remove it.
# Catches stray binaries the canonical self-uninstall and install-dir sweep
# missed (manual copies in ~/bin, /usr/local/bin shims, etc.).
# ---------------------------------------------------------------------------

find_all_gitmap_on_path() {
    # `command -v` returns only the first match; iterate PATH dirs explicitly.
    local IFS=':'
    local seen=""
    for d in $PATH; do
        [ -z "$d" ] && continue
        for name in gitmap-v28 gitmap.exe; do
            local candidate="$d/$name"
            if [ -f "$candidate" ] && [ -x "$candidate" ]; then
                # de-dupe via newline-delimited seen list
                case "$(printf '\n%s\n' "$seen")" in
                    *$'\n'"$candidate"$'\n'*) ;;
                    *) seen="${seen}${candidate}\n"; printf '%s\n' "$candidate" ;;
                esac
            fi
        done
    done
}

remove_stray_binaries() {
    local found
    found="$(find_all_gitmap_on_path)"
    if [ -z "$found" ]; then
        info "no stray gitmap-v28 binaries found on PATH"
        return
    fi

    local count
    count="$(printf '%s\n' "$found" | grep -c .)"
    info "found $count gitmap-v28 binary location(s):"
    printf '%s\n' "$found" | while IFS= read -r b; do info "  - $b"; done

    printf '%s\n' "$found" | while IFS= read -r bin; do
        [ -z "$bin" ] && continue
        if rm -f "$bin" 2>/dev/null; then
            ok "removed $bin"
        else
            # Try with sudo if it's in a system path.
            case "$bin" in
                /usr/*|/opt/*)
                    if sudo rm -f "$bin" 2>/dev/null; then
                        ok "removed $bin (sudo)"
                    else
                        err "could not remove $bin (try: sudo rm -f $bin)"
                    fi
                    ;;
                *)
                    err "could not remove $bin (check permissions)"
                    ;;
            esac
        fi
    done
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

printf '\n  %bgitmap quick uninstaller%b\n' "$c_cyan" "$c_reset"
printf '  %b------------------------%b\n\n' "$c_dgray" "$c_reset"

step "[1/4] Trying canonical self-uninstall"
if ! try_self_uninstall; then
    printf '\n'
    step "[2/4] Manual sweep — locating install dir"
    ROOT="$(resolve_install_dir)"
    if [ -n "$ROOT" ]; then info "Install dir: $ROOT"; else warn "no install dir found"; fi

    printf '\n'
    step "[3/4] Removing install files"
    remove_install_files "$ROOT"

    printf '\n'
    step "[4/4] Cleaning shell rc files"
    clean_rc_files "$ROOT"
fi

printf '\n'
step "Exhaustive PATH sweep — removing any remaining gitmap-v28 binaries"
remove_stray_binaries

printf '\n'
step "User data"
remove_data_folder

printf '\n  %bDone. Open a new shell to refresh PATH.%b\n\n' "$c_green" "$c_reset"
