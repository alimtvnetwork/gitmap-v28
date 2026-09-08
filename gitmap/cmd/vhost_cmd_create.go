package cmd

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type createFlagsHolder struct {
	rawType   string
	domain    string
	root      string
	port      int
	isSsl     bool
	fastcgi   string
	aliases   string
	isEnabled bool
	isDryRun  bool
	availDir  string
	enabDir   string
}

func setupCreateFlagSet(fs *flag.FlagSet, h *createFlagsHolder) {
	fs.StringVar(&h.rawType, "type", "php", "Site type (wordpress, laravel, php, static)")
	fs.StringVar(&h.domain, "domain", "", "Virtual host domain name")
	fs.StringVar(&h.root, "root", "", "Document root directory")
	fs.IntVar(&h.port, "port", 80, "Listen port")
	fs.BoolVar(&h.isSsl, "ssl", false, "Enable SSL port 443")
	fs.StringVar(&h.fastcgi, "fastcgi", "", "FastCGI pass upstream or socket")
	fs.StringVar(&h.fastcgi, "php-sock", "", "FastCGI pass upstream alias")
	fs.StringVar(&h.aliases, "aliases", "", "Server aliases")
	fs.BoolVar(&h.isEnabled, "enable", false, "Enable virtual host immediately")
	fs.BoolVar(&h.isDryRun, "dry-run", false, "Simulate execution without modifying files")
	fs.StringVar(&h.availDir, "sites-available", "", "Custom sites-available path")
	fs.StringVar(&h.enabDir, "sites-enabled", "", "Custom sites-enabled path")
}

func resolvePositionalDomain(posArgs []string, current string) string {
	hasCurrent := current != ""
	if hasCurrent {
		return current
	}
	hasPos := len(posArgs) > 0
	if hasPos {
		return posArgs[0]
	}

	return ""
}

func resolvePositionalRoot(posArgs []string, current string) string {
	hasCurrent := current != ""
	if hasCurrent {
		return current
	}
	hasPos := len(posArgs) > 1
	if hasPos {
		return posArgs[1]
	}

	return ""
}

func populateCreatePositional(fs *flag.FlagSet, h *createFlagsHolder) {
	posArgs := fs.Args()
	h.domain = resolvePositionalDomain(posArgs, h.domain)
	h.root = resolvePositionalRoot(posArgs, h.root)
}

func buildCreateConfig(h createFlagsHolder) (VHostConfig, VHostOptions, *apperror.AppError) {
	st, parseErr := ParseVHostSiteType(h.rawType)
	if parseErr != nil {
		return VHostConfig{}, VHostOptions{}, parseErr
	}
	cfg := VHostConfig{
		SiteType:     st,
		Domain:       h.domain,
		Aliases:      h.aliases,
		DocumentRoot: h.root,
		Port:         h.port,
		FastCGIPass:  h.fastcgi,
		IsSslEnabled: h.isSsl,
	}
	opts := VHostOptions{
		SitesAvailableDir: h.availDir,
		SitesEnabledDir:   h.enabDir,
		IsEnabled:         h.isEnabled,
		IsDryRun:          h.isDryRun,
	}

	return cfg, opts, nil
}

func runVHostCreate(args []string) error {
	var holder createFlagsHolder
	fs := flag.NewFlagSet("vhost create", flag.ContinueOnError)
	setupCreateFlagSet(fs, &holder)
	err := fs.Parse(args)
	if err == flag.ErrHelp {
		cliexit.Exit(0)
	}
	if err != nil {
		return apperror.WrapSimple(err, "flag.Parse")
	}
	populateCreatePositional(fs, &holder)
	cfg, opts, buildErr := buildCreateConfig(holder)
	if buildErr != nil {
		return buildErr
	}
	targetPath, createErr := CreateVHost(cfg, opts)
	if createErr != nil {
		return createErr
	}
	fmt.Printf("%sVirtual host created for %s at %s%s\n", constants.ColorGreen, cfg.Domain, targetPath, constants.ColorReset)

	return nil
}
