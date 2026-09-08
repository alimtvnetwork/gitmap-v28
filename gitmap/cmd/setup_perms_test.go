package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePermsAppType(t *testing.T) {
	if ParsePermsAppType("wordpress") != PermsAppTypeWordPress {
		t.Error("expected wordpress type")
	}
	if ParsePermsAppType("wp") != PermsAppTypeWordPress {
		t.Error("expected wp type")
	}
	if ParsePermsAppType("laravel") != PermsAppTypeLaravel {
		t.Error("expected laravel type")
	}
	if ParsePermsAppType("artisan") != PermsAppTypeLaravel {
		t.Error("expected artisan type")
	}
	if ParsePermsAppType("other") != PermsAppTypeGeneric {
		t.Error("expected generic type")
	}
}

func TestParsePermsFlags(t *testing.T) {
	args := []string{"--dir=/var/www", "--app=wordpress", "--owner=www-data:www-data", "--fix", "--dry-run"}
	opts, err := parsePermsFlags(args)
	if err != nil {
		t.Fatalf("parsePermsFlags failed: %v", err)
	}
	if opts.TargetDir != "/var/www" {
		t.Errorf("expected TargetDir /var/www, got %s", opts.TargetDir)
	}
	if opts.AppType != PermsAppTypeWordPress {
		t.Errorf("expected AppType wordpress, got %s", opts.AppType)
	}
	if opts.Owner != "www-data:www-data" {
		t.Errorf("expected owner www-data:www-data, got %s", opts.Owner)
	}
	if !opts.IsFix || !opts.IsDryRun {
		t.Error("expected IsFix and IsDryRun to be true")
	}
}

func TestParsePermsFlagsPositional(t *testing.T) {
	args := []string{"/srv/www/site", "--fix"}
	opts, err := parsePermsFlags(args)
	if err != nil {
		t.Fatalf("parsePermsFlags failed: %v", err)
	}
	if opts.TargetDir != "/srv/www/site" {
		t.Errorf("expected positional dir /srv/www/site, got %s", opts.TargetDir)
	}
	if !opts.IsFix {
		t.Error("expected IsFix to be true")
	}
}

func TestIsSensitiveFile(t *testing.T) {
	if !isSensitiveFile("wp-config.php") {
		t.Error("expected wp-config.php to be sensitive")
	}
	if !isSensitiveFile(".env") {
		t.Error("expected .env to be sensitive")
	}
	if !isSensitiveFile("cert.key") {
		t.Error("expected cert.key to be sensitive")
	}
	if isSensitiveFile("index.php") {
		t.Error("index.php should not be sensitive")
	}
}

func TestIsWritablePathWordPress(t *testing.T) {
	if !isWritablePath("wp-content/uploads/2026/file.jpg", PermsAppTypeWordPress) {
		t.Error("expected wp-content/uploads to be writable")
	}
	if isWritablePath("wp-includes/version.php", PermsAppTypeWordPress) {
		t.Error("wp-includes should not be writable")
	}
}

func TestIsWritablePathLaravel(t *testing.T) {
	if !isWritablePath("storage/logs/laravel.log", PermsAppTypeLaravel) {
		t.Error("expected storage to be writable")
	}
	if !isWritablePath("bootstrap/cache/services.php", PermsAppTypeLaravel) {
		t.Error("expected bootstrap/cache to be writable")
	}
	if isWritablePath("app/Models/User.php", PermsAppTypeLaravel) {
		t.Error("app/Models should not be writable")
	}
}

func TestResolveExpectedUnixMode(t *testing.T) {
	dirMode := resolveExpectedUnixMode("app", true, PermsAppTypeGeneric)
	if dirMode != 0755 {
		t.Errorf("expected dir mode 0755, got %04o", dirMode)
	}
	writeDirMode := resolveExpectedUnixMode("storage", true, PermsAppTypeLaravel)
	if writeDirMode != 0775 {
		t.Errorf("expected writable dir mode 0775, got %04o", writeDirMode)
	}
	credMode := resolveExpectedUnixMode(".env", false, PermsAppTypeLaravel)
	if credMode != 0600 {
		t.Errorf("expected cred mode 0600, got %04o", credMode)
	}
	fileMode := resolveExpectedUnixMode("index.php", false, PermsAppTypeGeneric)
	if fileMode != 0644 {
		t.Errorf("expected file mode 0644, got %04o", fileMode)
	}
}

func TestApplyUnixPermissionsAudit(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	writeErr := os.WriteFile(envPath, []byte("APP_KEY="), 0644)
	if writeErr != nil {
		t.Fatalf("failed writing test .env: %v", writeErr)
	}

	opts := PermsOptions{
		TargetDir: dir,
		AppType:   PermsAppTypeLaravel,
		IsFix:     true,
		IsDryRun:  false,
	}
	report, err := ApplyUnixPermissions(opts)
	if err != nil {
		t.Fatalf("ApplyUnixPermissions failed: %v", err)
	}
	if report.ScannedCount < 2 {
		t.Errorf("expected at least 2 scanned items, got %d", report.ScannedCount)
	}
}

func TestApplyWindowsPermissionsMock(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	writeErr := os.WriteFile(envFile, []byte("KEY=test"), 0644)
	if writeErr != nil {
		t.Fatalf("failed creating .env: %v", writeErr)
	}

	executed := make([]string, 0)
	origRunner := commandRunner
	commandRunner = func(name string, args ...string) error {
		executed = append(executed, name)

		return nil
	}
	defer func() { commandRunner = origRunner }()

	opts := PermsOptions{TargetDir: dir, IsFix: true, IsDryRun: false}
	report, err := ApplyWindowsPermissions(opts)
	if err != nil {
		t.Fatalf("ApplyWindowsPermissions failed: %v", err)
	}
	if report.FixedCount == 0 {
		t.Error("expected fixed count > 0")
	}
	if len(executed) < 2 {
		t.Errorf("expected at least 2 commands executed, got %d", len(executed))
	}
}

func TestResolveSetupTargetPerms(t *testing.T) {
	handled, _ := resolveSetupTarget("perms", []string{"--dry-run"})
	if !handled {
		t.Error("expected 'perms' to be handled by resolveSetupTarget")
	}
	handledFull, _ := resolveSetupTarget("permissions", []string{"--dry-run"})
	if !handledFull {
		t.Error("expected 'permissions' to be handled by resolveSetupTarget")
	}
}

func TestApplyAppPermissionsWrappers(t *testing.T) {
	dir := t.TempDir()
	origRunner := commandRunner
	commandRunner = func(name string, args ...string) error { return nil }
	defer func() { commandRunner = origRunner }()

	wpErr := ApplyWordPressPermissions(dir)
	if wpErr != nil {
		t.Fatalf("ApplyWordPressPermissions failed: %v", wpErr)
	}
	larErr := ApplyLaravelPermissions(dir)
	if larErr != nil {
		t.Fatalf("ApplyLaravelPermissions failed: %v", larErr)
	}
}
