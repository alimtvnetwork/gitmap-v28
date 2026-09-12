package cmdvhost

import (
	"strings"
	"testing"
)

func TestParseVHostSiteTypeValid(t *testing.T) {
	cases := map[string]VHostSiteType{
		"wordpress": VHostSiteTypeWordpress,
		"wp":        VHostSiteTypeWordpress,
		"laravel":   VHostSiteTypeLaravel,
		"art":       VHostSiteTypeLaravel,
		"artisan":   VHostSiteTypeLaravel,
		"php":       VHostSiteTypePHP,
		"static":    VHostSiteTypeStatic,
		"html":      VHostSiteTypeStatic,
	}

	for input, want := range cases {
		got, err := ParseVHostSiteType(input)
		if err != nil {
			t.Fatalf("ParseVHostSiteType(%q) unexpected err: %v", input, err)
		}

		if got != want {
			t.Errorf("ParseVHostSiteType(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseVHostSiteTypeInvalid(t *testing.T) {
	_, err := ParseVHostSiteType("unknown_framework")
	if err == nil {
		t.Fatal("expected error for unknown site type")
	}
}

func TestValidateVHostConfigRequiredFields(t *testing.T) {
	cfg := VHostConfig{}
	err := ValidateVHostConfig(cfg)
	if err == nil {
		t.Fatal("expected validation error for empty config")
	}
}

func TestValidateVHostConfigMissingRoot(t *testing.T) {
	cfg := VHostConfig{Domain: "example.com", SiteType: VHostSiteTypeWordpress}
	err := ValidateVHostConfig(cfg)
	if err == nil {
		t.Fatal("expected validation error for missing document root")
	}
}

func TestRenderVHostWordPress(t *testing.T) {
	cfg := VHostConfig{
		SiteType:     VHostSiteTypeWordpress,
		Domain:       "wp.local",
		DocumentRoot: "/var/www/wordpress",
	}

	out, err := RenderVHostConfig(cfg)
	if err != nil {
		t.Fatalf("RenderVHostConfig failed: %v", err)
	}

	assertContains(t, out, "server_name wp.local;")
	assertContains(t, out, "try_files $uri $uri/ /index.php?$args;")
	assertContains(t, out, "location = /wp-config.php")
	assertContains(t, out, "location = /xmlrpc.php")
	assertContains(t, out, "location ~* /wp-content/uploads/.*\\.php$")
	assertContains(t, out, "# >>> gitmap:vhost/wordpress/wp.local >>>")
	assertContains(t, out, "# <<< gitmap:vhost/wordpress/wp.local <<<")
}

func TestRenderVHostLaravel(t *testing.T) {
	cfg := VHostConfig{
		SiteType:     VHostSiteTypeLaravel,
		Domain:       "laravel.local",
		DocumentRoot: "/var/www/laravel",
	}

	out, err := RenderVHostConfig(cfg)
	if err != nil {
		t.Fatalf("RenderVHostConfig failed: %v", err)
	}

	assertContains(t, out, "root /var/www/laravel/public;")
	assertContains(t, out, "try_files $uri $uri/ /index.php?$query_string;")
	assertContains(t, out, "location ~ /\\.env")
	assertContains(t, out, "# >>> gitmap:vhost/laravel/laravel.local >>>")
}

func TestRenderVHostPHP(t *testing.T) {
	cfg := VHostConfig{
		SiteType:     VHostSiteTypePHP,
		Domain:       "php.local",
		DocumentRoot: "/var/www/phpapp",
	}

	out, err := RenderVHostConfig(cfg)
	if err != nil {
		t.Fatalf("RenderVHostConfig failed: %v", err)
	}

	assertContains(t, out, "fastcgi_pass unix:/run/php/php-fpm.sock;")
	assertContains(t, out, "# >>> gitmap:vhost/php/php.local >>>")
}

func TestRenderVHostStatic(t *testing.T) {
	cfg := VHostConfig{
		SiteType:     VHostSiteTypeStatic,
		Domain:       "static.local",
		DocumentRoot: "/var/www/static",
	}

	out, err := RenderVHostConfig(cfg)
	if err != nil {
		t.Fatalf("RenderVHostConfig failed: %v", err)
	}

	assertContains(t, out, "try_files $uri $uri/ =404;")
	assertNotContains(t, out, "fastcgi_pass")
}

func assertContains(t *testing.T, content, substr string) {
	t.Helper()
	hasMatch := strings.Contains(content, substr)
	if hasMatch {
		return
	}

	t.Errorf("expected content to contain %q", substr)
}

func assertNotContains(t *testing.T, content, substr string) {
	t.Helper()
	hasSub := strings.Contains(content, substr)
	if hasSub {
		t.Errorf("expected content NOT to contain %q", substr)
	}
}
