package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func fileOrDirExists(name string) bool {
	_, err := os.Stat(name)

	return err == nil
}

func resolveCurrentWorkingDir() string {
	pwd, err := os.Getwd()
	if err == nil {
		return pwd
	}

	return ""
}

func detectProjectWorkingDir() string {
	hasProject := fileOrDirExists("public") || fileOrDirExists("wp-config.php")
	if hasProject {
		return resolveCurrentWorkingDir()
	}

	return ""
}

func detectDefaultDocumentRoot(domain string) string {
	pwd := detectProjectWorkingDir()
	if pwd != "" {
		return pwd
	}
	if runtime.GOOS == "windows" {
		return filepath.Join("C:\\www", domain)
	}

	return filepath.Join("/var/www", domain)
}

func detectDefaultSiteType(rawType string) string {
	if rawType != "" && rawType != "php" {
		return rawType
	}
	if _, err := os.Stat("artisan"); err == nil {
		return "laravel"
	}
	if _, err := os.Stat("wp-config.php"); err == nil {
		return "wordpress"
	}

	return "php"
}

func persistSiteRecord(cfg VHostConfig, targetPath string, isDryRun bool) error {
	if isDryRun {
		return nil
	}
	db, err := store.OpenSitesSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec := store.SiteRecord{
		Domain:          cfg.Domain,
		SiteType:        string(cfg.SiteType),
		DocumentRoot:    cfg.DocumentRoot,
		NginxConfigPath: targetPath,
		ListenPort:      cfg.Port,
		IsSslEnabled:    cfg.IsSslEnabled,
		IsActive:        true,
	}

	return db.UpsertSite(rec)
}
