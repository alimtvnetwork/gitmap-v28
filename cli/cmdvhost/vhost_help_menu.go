package cmdvhost

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderVHostHelp displays the styled two-column VHost help menu.
func RenderVHostHelp() {
	termout.RenderMenu(buildVHostHelpMenu())
}

func buildVHostHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Nginx Virtual Host Manager (gitmap vhost)",
		UsageLines: []string{
			"gitmap vhost [command] [args]",
			"gitmap vhost create <domain> <root> [flags]",
		},
		Sections: []termout.HelpSection{
			buildVHostSiteSection(),
			buildVHostNginxSection(),
		},
		FooterFlags: buildVHostFooterFlags(),
		Tips: []string{
			"Run 'gitmap vhost create --type=wordpress mysite.local /var/www/mysite' to scaffold WP.",
			"Use 'gitmap vhost test' to verify Nginx syntax before reloading.",
		},
	}
}

func buildVHostFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--type <type>", Description: "Site profile (wordpress, laravel, php, static)"},
		{Command: "--port <n>", Description: "HTTP listen port (default: 80)"},
		{Command: "--ssl", Description: "Enable SSL on port 443 with HTTPS redirect"},
		{Command: "--fastcgi <pass>", Description: "FastCGI socket or host:port upstream"},
		{Command: "--enable", Description: "Link to sites-enabled immediately upon creation"},
		{Command: "--dry-run", Description: "Output configuration to stdout without saving"},
		{Command: "-h, --help", Description: "Show this virtual host help menu"},
	}
}

func buildVHostSiteSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Virtual Host Management",
		Entries: []termout.CommandEntry{
			{Command: "list (ls)", Description: "List configured virtual host sites and active status"},
			{Command: "create <dom> <dir>", Description: "Generate production-grade Nginx configuration"},
			{Command: "enable (en) <domain>", Description: "Enable virtual host via sites-enabled symlink"},
			{Command: "disable (dis) <dom>", Description: "Disable virtual host by removing sites-enabled link"},
		},
	}
}

func buildVHostNginxSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Nginx Operations",
		Entries: []termout.CommandEntry{
			{Command: "test (check, t)", Description: "Validate Nginx configuration syntax (nginx -t)"},
			{Command: "reload (restart, r)", Description: "Gracefully reload Nginx daemon (nginx -s reload)"},
		},
	}
}
