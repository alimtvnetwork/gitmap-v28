package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverFastCGIPassCustom(t *testing.T) {
	custom := "unix:/custom/path.sock"
	got := ResolveFastCGIPass(custom)
	if got != custom {
		t.Fatalf("expected %s, got %s", custom, got)
	}
}

func TestDiscoverFastCGIPassPatterns(t *testing.T) {
	dir := t.TempDir()
	sock81 := filepath.Join(dir, "php8.1-fpm.sock")
	sock83 := filepath.Join(dir, "php8.3-fpm.sock")
	writeErr1 := os.WriteFile(sock81, []byte{}, 0600)
	if writeErr1 != nil {
		t.Fatalf("write sock81: %v", writeErr1)
	}

	writeErr2 := os.WriteFile(sock83, []byte{}, 0600)
	if writeErr2 != nil {
		t.Fatalf("write sock83: %v", writeErr2)
	}

	patterns := []string{filepath.Join(dir, "php*-fpm.sock")}
	got := DiscoverFastCGIPassWithPatterns(patterns)
	want := "unix:" + sock83
	if got != want {
		t.Fatalf("expected newest socket %s, got %s", want, got)
	}
}

func TestDiscoverFastCGIPassFallback(t *testing.T) {
	dir := t.TempDir()
	patterns := []string{filepath.Join(dir, "nonexistent*.sock")}
	got := DiscoverFastCGIPassWithPatterns(patterns)
	hasDefault := strings.Contains(got, "php-fpm.sock") || got == "127.0.0.1:9000"
	if !hasDefault {
		t.Fatalf("expected fallback socket or loopback, got %s", got)
	}
}

func createTestVHostOptions(dir string) VHostOptions {
	availDir := filepath.Join(dir, "sites-available")
	enabDir := filepath.Join(dir, "sites-enabled")
	confDDir := filepath.Join(dir, "conf.d")
	_ = os.MkdirAll(availDir, 0755)
	_ = os.MkdirAll(enabDir, 0755)
	_ = os.MkdirAll(confDDir, 0755)

	return VHostOptions{
		SitesAvailableDir: availDir,
		SitesEnabledDir:   enabDir,
		ConfDDir:          confDDir,
	}
}

func TestCreateVHostWordPress(t *testing.T) {
	dir := t.TempDir()
	opts := createTestVHostOptions(dir)
	cfg := VHostConfig{
		SiteType:     VHostSiteTypeWordpress,
		Domain:       "wp.test",
		DocumentRoot: "/var/www/wp",
	}

	path, err := CreateVHost(cfg, opts)
	if err != nil {
		t.Fatalf("CreateVHost failed: %v", err)
	}

	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile failed: %v", readErr)
	}

	body := string(content)
	assertContains(t, body, "server_name wp.test;")
	assertContains(t, body, "# >>> gitmap:vhost/wordpress/wp.test >>>")
	assertContains(t, body, "# <<< gitmap:vhost/wordpress/wp.test <<<")
}

func TestEnableAndDisableVHost(t *testing.T) {
	dir := t.TempDir()
	opts := createTestVHostOptions(dir)
	cfg := VHostConfig{
		SiteType:     VHostSiteTypeLaravel,
		Domain:       "laravel.test",
		DocumentRoot: "/var/www/laravel",
	}

	_, createErr := CreateVHost(cfg, opts)
	if createErr != nil {
		t.Fatalf("CreateVHost failed: %v", createErr)
	}

	enableErr := EnableVHost("laravel.test", opts)
	if enableErr != nil {
		t.Skipf("symlink unavailable on platform: %v", enableErr)
	}

	enableIdempotentErr := EnableVHost("laravel.test", opts)
	if enableIdempotentErr != nil {
		t.Fatalf("second EnableVHost failed: %v", enableIdempotentErr)
	}

	vhosts, listErr := ListVHosts(opts)
	if listErr != nil {
		t.Fatalf("ListVHosts failed: %v", listErr)
	}

	if len(vhosts) != 1 || !vhosts[0].IsEnabled {
		t.Fatalf("expected 1 enabled vhost, got %+v", vhosts)
	}

	disableErr := DisableVHost("laravel.test", opts)
	if disableErr != nil {
		t.Fatalf("DisableVHost failed: %v", disableErr)
	}

	disableIdempotentErr := DisableVHost("laravel.test", opts)
	if disableIdempotentErr != nil {
		t.Fatalf("second DisableVHost failed: %v", disableIdempotentErr)
	}
}

func TestEnableVHostMissing(t *testing.T) {
	dir := t.TempDir()
	opts := createTestVHostOptions(dir)
	err := EnableVHost("notfound.test", opts)
	if err == nil {
		t.Fatal("expected error enabling non-existent vhost")
	}
}

func TestNginxTestAndReloadDryRun(t *testing.T) {
	opts := VHostOptions{IsDryRun: true}
	msg, err := TestNginxConfig(opts)
	if err != nil {
		t.Fatalf("TestNginxConfig dry-run failed: %v", err)
	}

	assertContains(t, msg, "dry-run")
	reloadErr := ReloadNginx(opts)
	if reloadErr != nil {
		t.Fatalf("ReloadNginx dry-run failed: %v", reloadErr)
	}
}

func TestRunVHostDispatchInvalid(t *testing.T) {
	err := runVHost([]string{"unknownsubcommand"})
	if err == nil {
		t.Fatal("expected error for invalid vhost subcommand")
	}
}

func TestRunVHostListEmpty(t *testing.T) {
	err := runVHostList([]string{})
	if err != nil {
		t.Fatalf("runVHostList failed: %v", err)
	}
}
