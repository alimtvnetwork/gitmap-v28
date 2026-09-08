package cmd

// VHostSiteType represents the supported site types for virtual hosts.
type VHostSiteType string

// SiteType alias for VHostSiteType.
type SiteType = VHostSiteType

const (
	VHostSiteTypeWordpress VHostSiteType = "wordpress"
	VHostSiteTypeLaravel   VHostSiteType = "laravel"
	VHostSiteTypePHP       VHostSiteType = "php"
	VHostSiteTypeStatic    VHostSiteType = "static"
)

const (
	defaultVHostPort           = 80
	defaultVHostSslPort        = 443
	defaultVHostMaxBodySize    = "64M"
	defaultVHostFastCGITimeout = 60
	defaultVHostFastCGIPass    = "unix:/run/php/php-fpm.sock"
	defaultSitesAvailableDir   = "/etc/nginx/sites-available"
	defaultSitesEnabledDir     = "/etc/nginx/sites-enabled"
	defaultConfDDir            = "/etc/nginx/conf.d"
)

// VHostConfig holds parameters for virtual host template generation.
type VHostConfig struct {
	SiteType       VHostSiteType
	Domain         string
	Aliases        string
	DocumentRoot   string
	Port           int
	FastCGIPass    string
	MaxBodySize    string
	FastCGITimeout int
	IsSslEnabled   bool
}

// VHostOptions holds operational flags and directory overrides for vhost commands.
type VHostOptions struct {
	SitesAvailableDir string
	SitesEnabledDir   string
	ConfDDir          string
	NginxBin          string
	IsEnabled         bool
	IsDryRun          bool
}

// VHostInfo describes a discovered or managed virtual host site.
type VHostInfo struct {
	Domain     string
	SiteType   VHostSiteType
	IsEnabled  bool
	ConfigPath string
	TargetPath string
}
