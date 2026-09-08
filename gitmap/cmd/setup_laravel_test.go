package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseLaravelFlags(t *testing.T) {
	args := []string{
		"--dir=/tmp/laravel_app",
		"--app-name=CustomApp",
		"--app-env=staging",
		"--app-debug=false",
		"--db-connection=sqlite",
		"--db-name=database.sqlite",
		"--vhost",
		"--domain=laravel.test",
		"--port=8088",
		"--dry-run",
	}
	opts, err := parseLaravelFlags(args)
	if err != nil {
		t.Fatalf("parseLaravelFlags failed: %v", err)
	}
	if opts.TargetDir != "/tmp/laravel_app" {
		t.Errorf("expected TargetDir /tmp/laravel_app, got %s", opts.TargetDir)
	}
	if opts.EnvOpts.AppName != "CustomApp" || opts.EnvOpts.AppEnv != "staging" {
		t.Errorf("unexpected AppName/Env: %s, %s", opts.EnvOpts.AppName, opts.EnvOpts.AppEnv)
	}
	if opts.EnvOpts.DBConnection != "sqlite" || opts.EnvOpts.DBDatabase != "database.sqlite" {
		t.Errorf("unexpected DB opts: %s, %s", opts.EnvOpts.DBConnection, opts.EnvOpts.DBDatabase)
	}
	if !opts.IsVHost || !opts.IsDryRun {
		t.Error("expected boolean flags to be true")
	}
}

func TestParseLaravelFlagsPositional(t *testing.T) {
	args := []string{"/var/www/art_app", "--app-name=PosApp"}
	opts, err := parseLaravelFlags(args)
	if err != nil {
		t.Fatalf("parseLaravelFlags failed: %v", err)
	}
	if opts.TargetDir != "/var/www/art_app" {
		t.Errorf("expected positional dir /var/www/art_app, got %s", opts.TargetDir)
	}
	if opts.EnvOpts.AppName != "PosApp" {
		t.Errorf("expected app name PosApp, got %s", opts.EnvOpts.AppName)
	}
}

func TestSetupLaravelDryRun(t *testing.T) {
	dir := t.TempDir()
	opts := defaultLaravelSetupOptions()
	opts.TargetDir = dir
	opts.IsDryRun = true
	opts.IsStorageLink = false

	err := SetupLaravel(opts)
	if err != nil {
		t.Fatalf("SetupLaravel dry-run failed: %v", err)
	}
	envFile := filepath.Join(dir, ".env")
	hasFile := fileExists(envFile)
	if hasFile {
		t.Error(".env should not be written during dry-run")
	}
}

func TestSetupLaravelRealExecution(t *testing.T) {
	dir := t.TempDir()
	opts := defaultLaravelSetupOptions()
	opts.TargetDir = dir
	opts.EnvOpts.AppName = "ProdLaravel"
	opts.EnvOpts.DBDatabase = "prod_db"
	opts.IsStorageLink = false
	opts.IsVHost = true
	opts.Domain = "laratest.local"

	err := SetupLaravel(opts)
	if err != nil {
		t.Fatalf("SetupLaravel execution failed: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	content, readErr := os.ReadFile(envFile)
	if readErr != nil {
		t.Fatalf("failed reading written .env: %v", readErr)
	}
	str := string(content)
	assertContains(t, str, "APP_NAME=ProdLaravel")
	assertContains(t, str, "DB_DATABASE=prod_db")
	assertContains(t, str, "APP_KEY=base64:")

	vhostFile := filepath.Join(dir, "nginx-laratest.local.conf")
	vhostContent, vhostErr := os.ReadFile(vhostFile)
	if vhostErr != nil {
		t.Fatalf("failed reading written vhost file: %v", vhostErr)
	}
	assertContains(t, string(vhostContent), "server_name laratest.local;")
}

func TestSetupLaravelExistingEnvPreserved(t *testing.T) {
	dir := t.TempDir()
	exampleContent := "APP_NAME=OldApp\nCUSTOM_TOKEN=secret123\n"
	writeErr := os.WriteFile(filepath.Join(dir, ".env.example"), []byte(exampleContent), 0644)
	if writeErr != nil {
		t.Fatalf("failed writing .env.example: %v", writeErr)
	}

	opts := defaultLaravelSetupOptions()
	opts.TargetDir = dir
	opts.EnvOpts.AppName = "NewApp"
	opts.IsStorageLink = false

	err := SetupLaravel(opts)
	if err != nil {
		t.Fatalf("SetupLaravel failed: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	content, readErr := os.ReadFile(envFile)
	if readErr != nil {
		t.Fatalf("failed reading .env: %v", readErr)
	}
	str := string(content)
	assertContains(t, str, "APP_NAME=NewApp")
	assertContains(t, str, "CUSTOM_TOKEN=secret123")
	assertContains(t, str, "APP_KEY=base64:")
}

func TestResolveSetupTargetLaravel(t *testing.T) {
	handled, _ := resolveSetupTarget("laravel", []string{"--dry-run"})
	if !handled {
		t.Error("expected 'laravel' to be handled by resolveSetupTarget")
	}
	handledArt, _ := resolveSetupTarget("art", []string{"--dry-run"})
	if !handledArt {
		t.Error("expected 'art' to be handled by resolveSetupTarget")
	}
	handledArtisan, _ := resolveSetupTarget("artisan", []string{"--dry-run"})
	if !handledArtisan {
		t.Error("expected 'artisan' to be handled by resolveSetupTarget")
	}
}

func TestCreateStorageSymlink(t *testing.T) {
	dir := t.TempDir()
	err := createStorageSymlink(dir)
	if err != nil {
		t.Fatalf("createStorageSymlink failed: %v", err)
	}
	targetPath := filepath.Join(dir, "storage", "app", "public")
	hasTarget := dirExists(targetPath)
	if !hasTarget {
		t.Error("expected storage/app/public directory to be created")
	}
}
