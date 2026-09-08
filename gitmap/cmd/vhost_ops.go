package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func applyVHostDirDefaults(opts VHostOptions) VHostOptions {
	res := opts
	isAvailEmpty := res.SitesAvailableDir == ""
	if isAvailEmpty {
		res.SitesAvailableDir = defaultSitesAvailableDir
	}
	isEnabledEmpty := res.SitesEnabledDir == ""
	if isEnabledEmpty {
		res.SitesEnabledDir = defaultSitesEnabledDir
	}

	return res
}

// ApplyVHostOptionDefaults applies standard directory paths and defaults.
func ApplyVHostOptionDefaults(opts VHostOptions) VHostOptions {
	res := applyVHostDirDefaults(opts)
	isConfDEmpty := res.ConfDDir == ""
	if isConfDEmpty {
		res.ConfDDir = defaultConfDDir
	}
	isNginxEmpty := res.NginxBin == ""
	if isNginxEmpty {
		res.NginxBin = "nginx"
	}

	return res
}

func scanVHostEntries(availDir, enabledDir string) []VHostInfo {
	entries, err := os.ReadDir(availDir)
	hasError := err != nil
	if hasError {
		return []VHostInfo{}
	}
	var res []VHostInfo
	for _, entry := range entries {
		info := readVHostInfo(availDir, enabledDir, entry.Name())
		res = append(res, info)
	}

	return res
}

// ListVHosts retrieves virtual hosts from sites-available or conf.d.
func ListVHosts(opts VHostOptions) ([]VHostInfo, *apperror.AppError) {
	applied := ApplyVHostOptionDefaults(opts)
	_, err := os.Stat(applied.SitesAvailableDir)
	hasDir := err == nil
	if hasDir {
		return scanVHostEntries(applied.SitesAvailableDir, applied.SitesEnabledDir), nil
	}
	_, confDErr := os.Stat(applied.ConfDDir)
	hasConfD := confDErr == nil
	if hasConfD {
		return scanVHostEntries(applied.ConfDDir, applied.ConfDDir), nil
	}

	return []VHostInfo{}, nil
}

func prepareVHostConfig(cfg VHostConfig) VHostConfig {
	out := cfg
	isStatic := cfg.SiteType == VHostSiteTypeStatic
	if isStatic {
		return out
	}
	isPassEmpty := cfg.FastCGIPass == ""
	if isPassEmpty {
		out.FastCGIPass = DiscoverFastCGIPass()
	}

	return out
}

func resolveVHostTargetPath(opts VHostOptions, domain string) string {
	_, err := os.Stat(opts.SitesAvailableDir)
	hasAvail := err == nil
	if hasAvail {
		return filepath.Join(opts.SitesAvailableDir, domain)
	}

	return filepath.Join(opts.ConfDDir, domain+".conf")
}

func writeVHostFile(targetPath, content string) *apperror.AppError {
	dir := filepath.Dir(targetPath)
	dirErr := os.MkdirAll(dir, 0755)
	if dirErr != nil {
		return apperror.WrapSimple(dirErr, "os.MkdirAll")
	}
	writeErr := os.WriteFile(targetPath, []byte(content), 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "os.WriteFile")
	}

	return nil
}

func maybeEnableVHost(domain string, opts VHostOptions) *apperror.AppError {
	if opts.IsEnabled {
		return EnableVHost(domain, opts)
	}

	return nil
}

func writeAndEnableVHost(domain, targetPath, rendered string, opts VHostOptions) *apperror.AppError {
	writeErr := writeVHostFile(targetPath, rendered)
	if writeErr != nil {
		return writeErr
	}

	return maybeEnableVHost(domain, opts)
}

// CreateVHost renders and writes an Nginx virtual host configuration.
func CreateVHost(cfg VHostConfig, opts VHostOptions) (string, *apperror.AppError) {
	applied := ApplyVHostOptionDefaults(opts)
	prepCfg := prepareVHostConfig(cfg)
	rendered, renderErr := RenderVHostConfig(prepCfg)
	if renderErr != nil {
		return "", renderErr
	}
	targetPath := resolveVHostTargetPath(applied, prepCfg.Domain)
	if applied.IsDryRun {
		return targetPath, nil
	}
	err := writeAndEnableVHost(prepCfg.Domain, targetPath, rendered, applied)

	return targetPath, err
}

func removeVHostFiles(domain string, opts VHostOptions) {
	_ = DisableVHost(domain, opts)
	availTarget := filepath.Join(opts.SitesAvailableDir, domain)
	_ = os.Remove(availTarget)
	availConfTarget := filepath.Join(opts.SitesAvailableDir, domain+".conf")
	_ = os.Remove(availConfTarget)
	confDTarget := filepath.Join(opts.ConfDDir, domain+".conf")
	_ = os.Remove(confDTarget)
}

// RemoveVHost disables and removes an Nginx virtual host configuration.
func RemoveVHost(domain string, opts VHostOptions) *apperror.AppError {
	applied := ApplyVHostOptionDefaults(opts)
	if applied.IsDryRun {
		return nil
	}
	removeVHostFiles(domain, applied)

	return nil
}

func isNginxInstalled(bin string) bool {
	_, err := exec.LookPath(bin)

	return err == nil
}

// TestNginxConfig executes 'nginx -t' to validate configuration syntax.
func TestNginxConfig(opts VHostOptions) (string, *apperror.AppError) {
	applied := ApplyVHostOptionDefaults(opts)
	if applied.IsDryRun {
		return "nginx: configuration file test is successful (dry-run)", nil
	}
	if !isNginxInstalled(applied.NginxBin) {
		return "nginx: binary not in PATH (syntax test skipped)", nil
	}
	cmd := exec.Command(applied.NginxBin, "-t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), apperror.NewExecutionError("nginx test failed: " + string(out))
	}

	return string(out), nil
}

// ReloadNginx tests and triggers 'nginx -s reload'.
func ReloadNginx(opts VHostOptions) *apperror.AppError {
	applied := ApplyVHostOptionDefaults(opts)
	if applied.IsDryRun || !isNginxInstalled(applied.NginxBin) {
		return nil
	}
	_, testErr := TestNginxConfig(applied)
	if testErr != nil {
		return testErr
	}
	cmd := exec.Command(applied.NginxBin, "-s", "reload")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError("nginx reload failed: " + string(out))
	}

	return nil
}
