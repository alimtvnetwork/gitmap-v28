package cmdsetup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateWpConfigContentDefaults(t *testing.T) {
	opts := defaultWordPressOptions()
	opts.DBName = "test_wp_db"
	opts.DBUser = "wp_admin"
	opts.DBPassword = "secret_password"
	opts.DBHost = "10.0.0.1"
	opts.TablePrefix = "wp_custom_"
	opts.IsDebug = true

	out, err := GenerateWpConfigContent("", opts)
	if err != nil {
		t.Fatalf("GenerateWpConfigContent failed: %v", err)
	}

	assertContains(t, out, "define( 'DB_NAME', 'test_wp_db' );")
	assertContains(t, out, "define( 'DB_USER', 'wp_admin' );")
	assertContains(t, out, "define( 'DB_PASSWORD', 'secret_password' );")
	assertContains(t, out, "define( 'DB_HOST', '10.0.0.1' );")
	assertContains(t, out, "$table_prefix = 'wp_custom_';")
	assertContains(t, out, "define( 'WP_DEBUG', true );")
	assertContains(t, out, "// >>> gitmap:wp-salts >>>")
	assertContains(t, out, "// <<< gitmap:wp-salts <<<")
	assertContains(t, out, "define('AUTH_KEY', '")
}

func TestGenerateWpConfigContentSamplePreserved(t *testing.T) {
	sample := `<?php
// Custom comment
define( 'DB_NAME', 'database_name_here' );
define( 'DB_USER', 'username_here' );
define( 'DB_PASSWORD', 'password_here' );
define( 'DB_HOST', 'localhost' );
$table_prefix = 'wp_';
define( 'WP_DEBUG', false );
/* That's all, stop editing! */
`
	opts := defaultWordPressOptions()
	opts.DBName = "prod_db"
	out, err := GenerateWpConfigContent(sample, opts)
	if err != nil {
		t.Fatalf("GenerateWpConfigContent with sample failed: %v", err)
	}

	assertContains(t, out, "// Custom comment")
	assertContains(t, out, "define( 'DB_NAME', 'prod_db' );")
	assertContains(t, out, "// >>> gitmap:wp-salts >>>")
}

func TestParseWordPressFlags(t *testing.T) {
	args := []string{
		"--dir=/tmp/site",
		"--db-name=my_db",
		"--db-user=my_user",
		"--db-pass=my_pass",
		"--db-host=127.0.0.2",
		"--prefix=site_",
		"--debug",
		"--vhost",
		"--domain=site.test",
		"--port=8080",
		"--dry-run",
	}

	opts, err := parseWordPressFlags(args)
	if err != nil {
		t.Fatalf("parseWordPressFlags failed: %v", err)
	}

	if opts.TargetDir != "/tmp/site" {
		t.Errorf("expected TargetDir /tmp/site, got %s", opts.TargetDir)
	}

	if opts.DBName != "my_db" || opts.DBUser != "my_user" {
		t.Errorf("unexpected DB creds: %s, %s", opts.DBName, opts.DBUser)
	}

	if !opts.IsDebug || !opts.IsVHost || !opts.IsDryRun {
		t.Error("expected boolean flags to be true")
	}

	if opts.Domain != "site.test" || opts.Port != 8080 {
		t.Errorf("unexpected vhost opts: %s:%d", opts.Domain, opts.Port)
	}
}

func TestParseWordPressFlagsPositional(t *testing.T) {
	args := []string{"/var/www/mywp", "--db-name=pos_db"}
	opts, err := parseWordPressFlags(args)
	if err != nil {
		t.Fatalf("parseWordPressFlags failed: %v", err)
	}

	if opts.TargetDir != "/var/www/mywp" {
		t.Errorf("expected positional dir /var/www/mywp, got %s", opts.TargetDir)
	}

	if opts.DBName != "pos_db" {
		t.Errorf("expected db pos_db, got %s", opts.DBName)
	}
}

func TestSetupWordPressDryRun(t *testing.T) {
	dir := t.TempDir()
	opts := defaultWordPressOptions()
	opts.TargetDir = dir
	opts.IsDryRun = true

	err := SetupWordPress(opts)
	if err != nil {
		t.Fatalf("SetupWordPress dry-run failed: %v", err)
	}

	wpConfig := filepath.Join(dir, "wp-config.php")
	hasFile := fileExists(wpConfig)
	if hasFile {
		t.Error("wp-config.php should not be written during dry-run")
	}
}

func TestSetupWordPressRealExecution(t *testing.T) {
	dir := t.TempDir()
	opts := defaultWordPressOptions()
	opts.TargetDir = dir
	opts.DBName = "real_wp"
	opts.IsVHost = true
	opts.Domain = "wp.test"

	err := SetupWordPress(opts)
	if err != nil {
		t.Fatalf("SetupWordPress execution failed: %v", err)
	}

	wpConfig := filepath.Join(dir, "wp-config.php")
	content, readErr := os.ReadFile(wpConfig)
	if readErr != nil {
		t.Fatalf("failed reading written wp-config.php: %v", readErr)
	}

	assertContains(t, string(content), "define( 'DB_NAME', 'real_wp' );")

	vhostFile := filepath.Join(dir, "nginx-wp.test.conf")
	vhostContent, vhostErr := os.ReadFile(vhostFile)
	if vhostErr != nil {
		t.Fatalf("failed reading written vhost file: %v", vhostErr)
	}

	assertContains(t, string(vhostContent), "server_name wp.test;")
}

func TestResolveSetupTargetWordPress(t *testing.T) {
	handled, _ := resolveSetupTarget("wp", []string{"--dry-run"})
	if !handled {
		t.Error("expected 'wp' to be handled by resolveSetupTarget")
	}

	handledWp, _ := resolveSetupTarget("wordpress", []string{"--dry-run"})
	if !handledWp {
		t.Error("expected 'wordpress' to be handled by resolveSetupTarget")
	}

	handledUnknown, _ := resolveSetupTarget("unknown_cmd", []string{})
	if handledUnknown {
		t.Error("expected unknown_cmd not to be handled")
	}
}
