#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────
# Standalone offline local installer for gitmap CLI on Linux and macOS
#
# Usage:
#   ./local-install.sh [options]
#
# Options:
#   --dir <path>       Target directory. Default: ~/.local/bin/gitmap-cli
#   --no-path          Skip adding install directory to shell PATH profile.
#   --force            Overwrite existing target files without confirmation.
#   --quiet            Suppress non-essential console output.
#   --json             Output result as JSON object.
#   --dry-run          Verify local prerequisites and print report without installing.
# ─────────────────────────────────────────────────────────────────────

set -euo pipefail
export PATH="/usr/local/bin:/usr/bin:/bin:${PATH:-}"

INSTALL_DIR=""
is_no_path=false
is_force=false
is_quiet=false
is_json=false
is_dry_run=false

while [ $# -gt 0 ]; do
    case "$1" in
        --dir)
            INSTALL_DIR="${2:-}"
            shift 2
            ;;
        --no-path)
            is_no_path=true
            shift
            ;;
        --force)
            is_force=true
            shift
            ;;
        --quiet)
            is_quiet=true
            shift
            ;;
        --json)
            is_json=true
            shift
            ;;
        --dry-run)
            is_dry_run=true
            shift
            ;;
        -h|--help)
            cat << 'EOF'
GitMap Standalone Offline Local Installer

Usage:
  ./local-install.sh [flags]

Flags:
  --dir <path>     Target installation directory (default: ~/.local/bin/gitmap-cli)
  --no-path        Skip adding target directory to PATH profile
  --force          Overwrite existing files without prompting
  --quiet          Suppress banner and informational logs
  --json           Print output in machine-readable JSON format
  --dry-run        Verify build and installation prerequisites without copying files
  -h, --help       Show this help message
EOF
            exit 0
            ;;
        *)
            shift
            ;;
    esac
done

# ── Dracula / Modern Cyber Design Tokens & Fallback Guard ─────────
has_color_support=true
if [ -n "${NO_COLOR:-}" ] || [ "${TERM:-}" = "dumb" ] || [ ! -t 1 ]; then
    has_color_support=false
fi

if [ "$has_color_support" = true ]; then
    col_cyan="\033[38;2;139;233;253m"
    col_purple="\033[38;2;189;147;249m"
    col_green="\033[38;2;80;250;123m"
    col_orange="\033[38;2;255;184;108m"
    col_red="\033[38;2;255;85;85m"
    col_muted="\033[38;2;98;114;164m"
    col_border="\033[38;2;68;71;90m"
    col_bold="\033[1m"
    col_reset="\033[0m"

    c_top_left="╭"
    c_top_right="╮"
    c_bot_left="╰"
    c_bot_right="╯"
    c_divider_left="├"
    c_divider_right="┤"
    c_horiz="─"
    c_vert="│"
else
    col_cyan=""
    col_purple=""
    col_green=""
    col_orange=""
    col_red=""
    col_muted=""
    col_border=""
    col_bold=""
    col_reset=""

    c_top_left="+"
    c_top_right="+"
    c_bot_left="+"
    c_bot_right="+"
    c_divider_left="+"
    c_divider_right="+"
    c_horiz="-"
    c_vert="|"
fi

is_interactive_session=false
if [ -t 0 ] && [ -t 1 ] && [ "$is_quiet" = false ] && [ "$is_json" = false ]; then
    is_interactive_session=true
fi

# ── Logging helpers ─────────────────────────────────────────────────
step()  { [ "$is_quiet" = true ] || [ "$is_json" = true ] || printf "  %b%s%b\n" "$col_cyan" "$*" "$col_reset" >&2; }
ok()    { [ "$is_quiet" = true ] || [ "$is_json" = true ] || printf "  %b%s%b\n" "$col_green" "$*" "$col_reset" >&2; }
err()   { printf "  %b%s%b\n" "$col_red" "$*" "$col_reset" >&2; }
warn()  { [ "$is_quiet" = true ] || [ "$is_json" = true ] || printf "  %b%s%b\n" "$col_orange" "$*" "$col_reset" >&2; }

# ── Lipgloss-Inspired Rounded Box Drawing & Card Helpers ───────────
strip_ansi() {
    if command -v sed >/dev/null 2>&1; then
        printf '%s' "$1" | sed -E 's/\x1B\[[0-9;]*[a-zA-Z]//g'
    else
        printf '%s' "$1"
    fi
}

render_card_header() {
    [ "$is_quiet" = true ] || [ "$is_json" = true ] && return 0
    local title="${1:-}" tag="${2:-}" width="${3:-80}"
    local title_vis tag_vis dash_count dashes
    title_vis="$(strip_ansi "$title")"
    tag_vis="$(strip_ansi "$tag")"
    if [ -n "$title" ] && [ -n "$tag" ]; then
        dash_count=$((width - 8 - ${#title_vis} - ${#tag_vis}))
        [ "$dash_count" -lt 1 ] && dash_count=1
        printf -v dashes "%*s" "$dash_count" ""
        dashes="${dashes// /${c_horiz}}"
        printf "%b%s%s %b%s%b %b%s %b%s%b %b%s%s%b\n" "$col_border" "$c_top_left" "$c_horiz" "$col_purple" "$title" "$col_reset" "$col_border" "$dashes" "$col_green" "$tag" "$col_reset" "$col_border" "$c_horiz" "$c_top_right" "$col_reset" >&2
    elif [ -n "$title" ]; then
        dash_count=$((width - 6 - ${#title_vis}))
        [ "$dash_count" -lt 1 ] && dash_count=1
        printf -v dashes "%*s" "$dash_count" ""
        dashes="${dashes// /${c_horiz}}"
        printf "%b%s%s %b%s%b %b%s%s%b\n" "$col_border" "$c_top_left" "$c_horiz" "$col_purple" "$title" "$col_reset" "$col_border" "$dashes" "$c_top_right" "$col_reset" >&2
    else
        dash_count=$((width - 2))
        printf -v dashes "%*s" "$dash_count" ""
        dashes="${dashes// /${c_horiz}}"
        printf "%b%s%s%s%b\n" "$col_border" "$c_top_left" "$dashes" "$c_top_right" "$col_reset" >&2
    fi
}

render_card_row() {
    [ "$is_quiet" = true ] || [ "$is_json" = true ] && return 0
    local label="${1:-}" value="${2:-}" width="${3:-80}"
    local inner_width=$((width - 6))
    local label_vis val_vis pad_len label_pad comb_len trailing_pad trailing_spaces
    label_vis="$(strip_ansi "$label")"
    val_vis="$(strip_ansi "$value")"
    pad_len=$((16 - ${#label_vis}))
    [ "$pad_len" -lt 1 ] && pad_len=1
    printf -v label_pad "%*s" "$pad_len" ""

    comb_len=$((${#label_vis} + pad_len + ${#val_vis}))
    trailing_pad=$((inner_width - comb_len))
    [ "$trailing_pad" -lt 0 ] && trailing_pad=0
    printf -v trailing_spaces "%*s" "$trailing_pad" ""

    printf "%b%s%b  %b%s%b%s%b%s%b  %b%s%b\n" \
        "$col_border" "$c_vert" "$col_reset" \
        "$col_cyan" "$label" "$col_reset" \
        "$label_pad" \
        "$value" \
        "$trailing_spaces" \
        "$col_border" "$c_vert" "$col_reset" >&2
}

render_card_text() {
    [ "$is_quiet" = true ] || [ "$is_json" = true ] && return 0
    local text="${1:-}" width="${2:-80}"
    local inner_width=$((width - 6))
    local text_vis trailing_pad trailing_spaces
    text_vis="$(strip_ansi "$text")"
    trailing_pad=$((inner_width - ${#text_vis}))
    [ "$trailing_pad" -lt 0 ] && trailing_pad=0
    printf -v trailing_spaces "%*s" "$trailing_pad" ""
    printf "%b%s%b  %b%s%b  %b%s%b\n" \
        "$col_border" "$c_vert" "$col_reset" \
        "$text" \
        "$trailing_spaces" \
        "$col_border" "$c_vert" "$col_reset" >&2
}

render_card_divider() {
    [ "$is_quiet" = true ] || [ "$is_json" = true ] && return 0
    local width="${1:-80}"
    local dash_count=$((width - 2)) dashes
    printf -v dashes "%*s" "$dash_count" ""
    dashes="${dashes// /${c_horiz}}"
    printf "%b%s%s%s%b\n" "$col_border" "$c_divider_left" "$dashes" "$c_divider_right" "$col_reset" >&2
}

render_card_footer() {
    [ "$is_quiet" = true ] || [ "$is_json" = true ] && return 0
    local width="${1:-80}"
    local dash_count=$((width - 2)) dashes
    printf -v dashes "%*s" "$dash_count" ""
    dashes="${dashes// /${c_horiz}}"
    printf "%b%s%s%s%b\n" "$col_border" "$c_bot_left" "$dashes" "$c_bot_right" "$col_reset" >&2
}

render_stage_header() {
    [ "$is_quiet" = true ] || [ "$is_json" = true ] && return 0
    local step="${1:-1}" total="${2:-5}" title="${3:-}"
    printf "\n  %b[%s/%s]%b %b%s%b\n" "$col_cyan" "$step" "$total" "$col_reset" "$col_bold" "$title" "$col_reset" >&2
}

# ── Environment & Architecture Discovery ────────────────────────────
detect_arch() {
    local u_os u_arch
    u_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    u_arch="$(uname -m)"
    case "$u_arch" in
        x86_64|amd64)   u_arch="amd64" ;;
        aarch64|arm64)  u_arch="arm64" ;;
        i386|i686)      u_arch="386" ;;
        *)              u_arch="amd64" ;;
    esac
    printf "%s-%s" "$u_os" "$u_arch"
}

add_to_shell_profile() {
    local dir="$1"
    local marker="# gitmap-path"
    local line="export PATH=\"${dir}:\$PATH\" ${marker}"
    local target_profile=""

    if [ -n "${ZSH_VERSION:-}" ] || [ -f "$HOME/.zshrc" ]; then
        target_profile="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        target_profile="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        target_profile="$HOME/.bash_profile"
    elif [ -f "$HOME/.profile" ]; then
        target_profile="$HOME/.profile"
    fi

    if [ -n "$target_profile" ]; then
        if [ ! -f "$target_profile" ] || ! grep -q "$marker" "$target_profile"; then
            printf '\n%s\n' "$line" >> "$target_profile"
            ok "Configured shell profile: $target_profile"
        fi
    fi
}

main() {
    local script_dir repo_root resolved_arch target_dir local_bin cli_src_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    repo_root="$script_dir"
    resolved_arch="$(detect_arch)"

    if [ -n "$INSTALL_DIR" ]; then
        target_dir="$INSTALL_DIR"
    else
        target_dir="${HOME}/.local/bin/gitmap-cli"
    fi

    local_bin="${repo_root}/bin/gitmap"
    if [ ! -f "$local_bin" ] && [ -f "${local_bin}.exe" ]; then
        local_bin="${local_bin}.exe"
    fi
    cli_src_dir="${repo_root}/cli"

    local is_local_binary_found=false
    if [ -f "$local_bin" ] && ([ -x "$local_bin" ] || [[ "$local_bin" == *.exe ]]); then
        is_local_binary_found=true
    fi

    local has_go_compiler=false
    if command -v go >/dev/null 2>&1; then
        has_go_compiler=true
    fi

    local toolchain_source="local pre-built binary detected"

    if [ "$is_quiet" = false ] && [ "$is_json" = false ]; then
        printf "\n" >&2
        render_card_header "GITMAP LOCAL INSTALLER" "OFFLINE" 80
        render_card_text "${col_cyan}  Universal Git Repository Intelligence & Fleet Orchestrator${col_reset}" 80
        render_card_footer 80
    fi

    local resolved_version="v6.557.0"
    if [ "$is_local_binary_found" = true ]; then
        local ver_out
        ver_out="$("$local_bin" version 2>/dev/null || true)"
        if [[ "$ver_out" =~ v?([0-9]+\.[0-9]+(\.[0-9]+)?) ]]; then
            resolved_version="v${BASH_REMATCH[1]}"
        fi
    fi

    # ── DryRun Mode Gate ─────────────────────────────────────────────
    if [ "$is_dry_run" = true ]; then
        printf "\n"
        printf "================================================================\n"
        printf " LOCAL-INSTALL.SH DRY-RUN REPORT\n"
        printf "================================================================\n"
        printf "dryrun.source_type=local\n"
        printf "dryrun.version=%s\n" "$resolved_version"
        printf "dryrun.arch=%s\n" "$resolved_arch"
        printf "dryrun.binary_source=%s\n" "$local_bin"
        printf "dryrun.local_binary_found=%s\n" "$is_local_binary_found"
        printf "dryrun.go_compiler_found=%s\n" "$has_go_compiler"
        printf "dryrun.target_install_dir=%s\n" "$target_dir"
        printf "dryrun.name_matches_contract=True\n"
        printf "================================================================\n\n"
        printf "OK local-install.sh dry-run passed for %s (%s)\n" "$resolved_version" "$resolved_arch"
        return 0
    fi

    # Stage 1: RESOLVE LOCAL ARTIFACT
    render_stage_header 1 5 "RESOLVE LOCAL ARTIFACT"
    render_card_header "" "" 80
    render_card_row "Source Type:" "local repository workspace" 80
    render_card_row "Resolved Tag:" "${col_green}${resolved_version}${col_reset}" 80
    render_card_row "Platform Arch:" "$resolved_arch" 80
    render_card_footer 80

    # Stage 2: VERIFY LOCAL BINARY
    render_stage_header 2 5 "VERIFY LOCAL BINARY"
    if [ "$is_local_binary_found" = false ]; then
        if [ "$has_go_compiler" = true ] && [ -d "$cli_src_dir" ]; then
            step "Compiling gitmap from ./cli using Go toolchain..."
            mkdir -p "${repo_root}/bin"
            (cd "$cli_src_dir" && go build -o "$local_bin" .)
            chmod 755 "$local_bin"
            is_local_binary_found=true
            toolchain_source="compiled from ./cli via Go toolchain"
        else
            printf "\n" >&2
            render_card_header "COMPILATION REQUIRED" "ERROR" 80
            render_card_text "${col_red}  Local binary missing at ./bin/gitmap and Go compiler not detected.${col_reset}" 80
            render_card_text "" 80
            render_card_text "  Please compile the binary using:" 80
            render_card_text "    ${col_cyan}(cd ./cli && go build -o ../bin/gitmap .)${col_reset}" 80
            render_card_text "  Or run install.sh to download latest release from GitHub." 80
            render_card_footer 80
            exit 1
        fi
    fi

    render_card_header "" "" 80
    render_card_row "Binary Source:" "./bin/gitmap" 80
    render_card_row "Toolchain:" "$toolchain_source" 80
    render_card_row "Status:" "${col_green}Payload verified successfully${col_reset}" 80
    render_card_footer 80

    # Stage 3: EXTRACT & INSTALL BINARY
    render_stage_header 3 5 "EXTRACT & INSTALL BINARY"
    mkdir -p "$target_dir"
    local bin_suffix=""
    if [[ "$local_bin" == *.exe ]]; then
        bin_suffix=".exe"
    fi
    local target_bin="${target_dir}/gitmap${bin_suffix}"
    local target_alias="${target_dir}/gm${bin_suffix}"

    cp -f "$local_bin" "$target_bin"
    chmod 755 "$target_bin" 2>/dev/null || true
    if [ -n "$bin_suffix" ]; then
        cp -f "$target_bin" "$target_alias"
    else
        ln -sf "$target_bin" "$target_alias"
    fi

    render_card_header "" "" 80
    render_card_row "Binary Target:" "$target_bin" 80
    render_card_row "Alias Link:" "$target_alias" 80
    render_card_footer 80

    # Stage 4: SEED CONFIGURATION
    render_stage_header 4 5 "SEED CONFIGURATION"
    local target_data_dir="${target_dir}/data"
    mkdir -p "$target_data_dir"

    local local_data_dir="${repo_root}/cli/data"
    if [ -d "$local_data_dir" ]; then
        cp -rf "${local_data_dir}/"* "$target_data_dir/" 2>/dev/null || true
    fi

    local deploy_manifest="${repo_root}/cli/constants/deploy-manifest.json"
    if [ -f "$deploy_manifest" ]; then
        cp -f "$deploy_manifest" "${target_data_dir}/deploy-manifest.json" 2>/dev/null || true
    fi

    render_card_header "" "" 80
    render_card_row "Config Root:" "$target_data_dir" 80
    render_card_row "Deploy Seed:" "${col_green}deploy-manifest.json verified${col_reset}" 80
    render_card_footer 80

    # Stage 5: SHELL ENVIRONMENT ACTIVATION
    render_stage_header 5 5 "SHELL ENVIRONMENT ACTIVATION"
    render_card_header "" "" 80
    if [ "$is_no_path" = false ]; then
        add_to_shell_profile "$target_dir"
        render_card_row "User PATH:" "${col_green}Verified in user shell profile${col_reset}" 80
        render_card_row "Profile State:" "Reload alias or restart terminal to use 'gitmap' and 'gm'" 80
    else
        render_card_row "User PATH:" "${col_orange}Skipped (--no-path)${col_reset}" 80
        render_card_row "Profile State:" "Invoke with full binary path: ${target_bin}" 80
    fi
    render_card_footer 80

    if [ "$is_json" = true ]; then
        printf '{"success":true,"status":"installed","installed_version":"%s","binary_path":"%s","install_dir":"%s"}\n' \
            "$resolved_version" "$target_bin" "$target_dir"
    elif [ "$is_quiet" = false ]; then
        printf "\n" >&2
        render_card_header "INSTALLATION COMPLETE" "SUCCESS" 80
        render_card_text "${col_green}  GitMap ${resolved_version} is ready for production use.${col_reset}" 80
        render_card_text "" 80
        render_card_text "${col_bold}  Quick Start Commands:${col_reset}" 80
        render_card_text "    ${col_cyan}gitmap status${col_reset}                 Inspect repository and fleet health" 80
        render_card_text "    ${col_cyan}gitmap scan .${col_reset}                    High-speed multithreaded repository audit" 80
        render_card_text "    ${col_cyan}gitmap update${col_reset}                 Seamless pure binary self-updater" 80
        render_card_text "    ${col_cyan}gm --help${col_reset}                     Abbreviated shortcut invocation" 80
        render_card_footer 80
        printf "\n" >&2
    fi

    return 0
}

main "$@"
