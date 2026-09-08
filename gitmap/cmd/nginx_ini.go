package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func probePHPIniPaths() string {
	out, err := exec.Command("php", "--ini").Output()
	if err != nil {
		return "  PHP CLI: not detected (default fallback: /etc/php/8.x/fpm/conf.d/)"
	}
	lines := strings.Split(string(out), "\n")
	var res []string
	for _, line := range lines {
		if strings.Contains(line, "Configuration File") || strings.Contains(line, "Scan for additional") {
			res = append(res, "  "+strings.TrimSpace(line))
		}
	}

	return strings.Join(res, "\n")
}

func printIniDirectivesWordPress() {
	fmt.Printf("%s--- WordPress Recommended Directives (php.ini) ---%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("; >>> gitmap:php.ini/wordpress >>>")
	fmt.Println("upload_max_filesize = 64M")
	fmt.Println("post_max_size       = 64M")
	fmt.Println("memory_limit        = 256M")
	fmt.Println("max_execution_time  = 300")
	fmt.Println("max_input_vars      = 3000")
	fmt.Println("; <<< gitmap:php.ini/wordpress <<<")
	fmt.Println()
}

func printIniDirectivesLaravel() {
	fmt.Printf("%s--- Laravel Recommended Directives (php.ini) ---%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("; >>> gitmap:php.ini/laravel >>>")
	fmt.Println("upload_max_filesize            = 32M")
	fmt.Println("post_max_size                  = 32M")
	fmt.Println("memory_limit                   = 512M")
	fmt.Println("max_execution_time             = 120")
	fmt.Println("opcache.enable                 = 1")
	fmt.Println("opcache.memory_consumption     = 128")
	fmt.Println("opcache.max_accelerated_files  = 10000")
	fmt.Println("; <<< gitmap:php.ini/laravel <<<")
	fmt.Println()
}

func printIniShowcaseDiff() {
	fmt.Printf("%s--- Showcase: Idempotent Marker-Block Merge ---%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Println("  Existing custom directives in php.ini are 100% preserved.")
	fmt.Println("  Gitmap injects a isolated marker block:")
	fmt.Println("  [Before]:")
	fmt.Println("    engine = On")
	fmt.Println("    short_open_tag = Off")
	fmt.Println("  [After Gitmap Merge]:")
	fmt.Println("    engine = On")
	fmt.Println("    short_open_tag = Off")
	fmt.Println("    ; >>> gitmap:php.ini/domain.com >>>")
	fmt.Println("    upload_max_filesize = 64M")
	fmt.Println("    memory_limit = 256M")
	fmt.Println("    ; <<< gitmap:php.ini/domain.com <<<")
	fmt.Println()
}

func printNginxFastCGIOverrideShowcase() {
	fmt.Printf("%s--- Alternative: Nginx fastcgi_param Per-VHost Override ---%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("  If you prefer not modifying system php.ini, Gitmap injects per-vhost:")
	fmt.Println("  fastcgi_param PHP_VALUE \"upload_max_filesize=64M \\n post_max_size=64M \\n memory_limit=256M\";")
	fmt.Println("  fastcgi_param PHP_ADMIN_VALUE \"max_execution_time=300\";")
	fmt.Println()
}

// runNginxIni displays the PHP/Nginx INI configuration engine and showcase.
func runNginxIni(args []string) error {
	fmt.Printf("\n%s=== Gitmap Nginx & PHP-FPM INI Configuration Engine ===%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Detected PHP Configuration Environment:")
	fmt.Println(probePHPIniPaths())
	fmt.Println()
	printIniDirectivesWordPress()
	printIniDirectivesLaravel()
	printIniShowcaseDiff()
	printNginxFastCGIOverrideShowcase()

	return nil
}
