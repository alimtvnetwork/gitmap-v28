package cmdssh

import (
	"fmt"
	"strings"
)

const defaultRouteIP = "192.168.0.1"
const defaultZshTheme = "fletcherm"

const netplanScriptTemplate = `#!/bin/bash
set -e
mkdir -p /etc/netplan
cat <<'EOF' > /etc/netplan/00-installer-config.yaml
network:
  renderer: networkd
  ethernets:
    ens33:
      dhcp4: false
      addresses:
        - %s
      routes:
        - to: default
          via: %s
      nameservers:
        addresses: [8.8.8.8, 1.1.1.1]
  version: 2
EOF
chmod 600 /etc/netplan/00-installer-config.yaml
netplan apply
`

const basePackagesScriptTemplate = `#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y curl wget git zsh net-tools htop build-essential
`

const createUserScriptTemplate = `#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive
id -u %s >/dev/null 2>&1 || useradd -m -s /usr/bin/zsh %s
%secho '%s ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/%s
chmod 0440 /etc/sudoers.d/%s
usermod -aG sudo %s 2>/dev/null || true
HOMEDIR=$(eval echo ~%s)
mkdir -p "$HOMEDIR/.ssh"
chmod 700 "$HOMEDIR/.ssh"
if [ ! -d "$HOMEDIR/.oh-my-zsh" ]; then
  su - %s -c 'RUNZSH=no CHSH=no sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh || wget -qO- https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended' || true
fi
if [ -f "$HOMEDIR/.zshrc" ]; then
  sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="%s"/' "$HOMEDIR/.zshrc"
else
  echo 'ZSH_THEME="%s"' > "$HOMEDIR/.zshrc"
fi
chown -R %s:%s "$HOMEDIR"
`

const updateThemeScriptTemplate = `#!/bin/bash
set -e
TARGET_ZSHRC="${HOME}/.zshrc"
if [ -n "$SUDO_USER" ] && [ "$SUDO_USER" != "root" ]; then
  TARGET_ZSHRC="/home/${SUDO_USER}/.zshrc"
fi
if [ -f "$TARGET_ZSHRC" ] && grep -q '^ZSH_THEME=' "$TARGET_ZSHRC"; then
  sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="%s"/' "$TARGET_ZSHRC"
else
  echo 'ZSH_THEME="%s"' >> "$TARGET_ZSHRC"
fi
if [ -f "${HOME}/.zshrc" ] && [ "${HOME}/.zshrc" != "$TARGET_ZSHRC" ]; then
  if grep -q '^ZSH_THEME=' "${HOME}/.zshrc"; then
    sed -i 's/^ZSH_THEME="[^"]*"/ZSH_THEME="%s"/' "${HOME}/.zshrc"
  else
    echo 'ZSH_THEME="%s"' >> "${HOME}/.zshrc"
  fi
fi
`

const purgeScriptTemplate = `#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive
apt-get autoremove --purge -y && apt-get clean
`

func resolveRouteIP(routeIP string) string {
	hasRoute := routeIP != ""
	if hasRoute {
		return routeIP
	}

	return defaultRouteIP
}

func resolveCIDRIP(ip string) string {
	hasSlash := strings.Contains(ip, "/")
	if hasSlash {
		return ip
	}

	return fmt.Sprintf("%s/24", ip)
}

// GenerateNetplanScript generates bash script to configure netplan static IP.
func GenerateNetplanScript(ip string, routeIP string) string {
	gw := resolveRouteIP(routeIP)
	cidr := resolveCIDRIP(ip)

	return fmt.Sprintf(netplanScriptTemplate, cidr, gw)
}

// GenerateBasePackagesScript generates bash script to install essential tools.
func GenerateBasePackagesScript() string {
	return basePackagesScriptTemplate
}

func resolveThemeName(theme string) string {
	hasTheme := theme != ""
	if hasTheme {
		return theme
	}

	return defaultZshTheme
}

func buildPasswordSnippet(username, password string) string {
	hasPass := password != ""
	if hasPass {
		return fmt.Sprintf("echo '%s:%s' | chpasswd\n", username, password)
	}

	return ""
}

// GenerateCreateUserScript generates bash script to create user with zsh & oh-my-zsh.
func GenerateCreateUserScript(username, password, theme string) string {
	appliedTheme := resolveThemeName(theme)
	passSnippet := buildPasswordSnippet(username, password)

	return fmt.Sprintf(createUserScriptTemplate,
		username, username,
		passSnippet,
		username, username, username, username, username,
		username,
		appliedTheme, appliedTheme,
		username, username)
}

// GenerateThemeScript generates bash script to update ZSH_THEME in ~/.zshrc.
func GenerateThemeScript(theme string) string {
	appliedTheme := resolveThemeName(theme)

	return fmt.Sprintf(updateThemeScriptTemplate, appliedTheme, appliedTheme, appliedTheme, appliedTheme)
}

// GeneratePurgeScript generates bash script to autoremove and clean package cache.
func GeneratePurgeScript() string {
	return purgeScriptTemplate
}
