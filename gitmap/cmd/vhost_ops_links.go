package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func parseMarkerSiteType(content string) VHostSiteType {
	const markerPrefix = "# >>> gitmap:vhost/"
	idx := strings.Index(content, markerPrefix)
	isMissing := idx < 0
	if isMissing {
		return VHostSiteTypeStatic
	}

	sub := content[idx+len(markerPrefix):]
	slashIdx := strings.Index(sub, "/")
	isNoSlash := slashIdx < 0
	if isNoSlash {
		return VHostSiteTypeStatic
	}

	rawType := sub[:slashIdx]
	siteType, err := ParseVHostSiteType(rawType)
	hasError := err != nil
	if hasError {
		return VHostSiteTypeStatic
	}

	return siteType
}

func resolveVHostSiteType(path string) VHostSiteType {
	bytes, err := os.ReadFile(path)
	hasError := err != nil
	if hasError {
		return VHostSiteTypeStatic
	}

	return parseMarkerSiteType(string(bytes))
}

func isVHostSymlinkEnabled(enabledDir, name string) bool {
	enabledPath := filepath.Join(enabledDir, name)
	_, err := os.Lstat(enabledPath)
	isEnabled := err == nil
	if isEnabled {
		return true
	}

	return false
}

func readVHostInfo(availDir, enabledDir, name string) VHostInfo {
	cfgPath := filepath.Join(availDir, name)
	siteType := resolveVHostSiteType(cfgPath)
	isEnabled := isVHostSymlinkEnabled(enabledDir, name)
	domain := strings.TrimSuffix(name, ".conf")

	return VHostInfo{
		Domain:     domain,
		SiteType:   siteType,
		IsEnabled:  isEnabled,
		ConfigPath: cfgPath,
	}
}

func findAvailableVHostFile(availDir, domain string) (string, *apperror.AppError) {
	path := filepath.Join(availDir, domain)
	_, err := os.Stat(path)
	isFound := err == nil
	if isFound {
		return path, nil
	}

	confPath := path + ".conf"
	_, confErr := os.Stat(confPath)
	isConfFound := confErr == nil
	if isConfFound {
		return confPath, nil
	}

	return "", apperror.NewValidationError("site config not found for domain: " + domain)
}

func findEnabledVHostPath(enabledDir, domain string) string {
	path := filepath.Join(enabledDir, domain)
	_, err := os.Lstat(path)
	isFound := err == nil
	if isFound {
		return path
	}

	confPath := path + ".conf"
	_, confErr := os.Lstat(confPath)
	isConfFound := confErr == nil
	if isConfFound {
		return confPath
	}

	return ""
}

func createVHostSymlink(source, target string) *apperror.AppError {
	dir := filepath.Dir(target)
	dirErr := os.MkdirAll(dir, 0755)
	if dirErr != nil {
		return apperror.WrapSimple(dirErr, "os.MkdirAll")
	}

	linkErr := os.Symlink(source, target)
	if linkErr != nil {
		return apperror.WrapSimple(linkErr, "os.Symlink")
	}

	return nil
}

// EnableVHost enables an Nginx site by symlinking it into sites-enabled.
func EnableVHost(domain string, opts VHostOptions) *apperror.AppError {
	applied := ApplyVHostOptionDefaults(opts)
	srcPath, srcErr := findAvailableVHostFile(applied.SitesAvailableDir, domain)
	if srcErr != nil {
		return srcErr
	}

	dstPath := filepath.Join(applied.SitesEnabledDir, filepath.Base(srcPath))
	_, statErr := os.Lstat(dstPath)
	isAlreadyEnabled := statErr == nil
	if isAlreadyEnabled {
		return nil
	}

	return createVHostSymlink(srcPath, dstPath)
}

// DisableVHost disables an Nginx site by removing its symlink from sites-enabled.
func DisableVHost(domain string, opts VHostOptions) *apperror.AppError {
	applied := ApplyVHostOptionDefaults(opts)
	target := findEnabledVHostPath(applied.SitesEnabledDir, domain)
	isEmpty := target == ""
	if isEmpty {
		return nil
	}

	remErr := os.Remove(target)
	if remErr != nil {
		return apperror.WrapSimple(remErr, "os.Remove")
	}

	return nil
}
