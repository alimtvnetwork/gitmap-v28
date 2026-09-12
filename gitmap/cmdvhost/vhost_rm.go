package cmdvhost

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func parseRmDomain(args []string) (string, bool) {
	for _, a := range args {
		if !isDryRunArg(a) {
			return a, true
		}
	}

	return "", false
}

func deleteSiteFromDB(domain string, isDryRun bool) {
	if isDryRun {
		return
	}

	db, err := store.OpenSitesSplitDB()
	if err != nil {
		return
	}

	defer db.Close()

	_ = db.DeleteSite(domain)
}

func printRmSuccess(domain string, isDryRun bool) {
	dryPrefix := ""
	if isDryRun {
		dryPrefix = "[dry-run] "
	}

	fmt.Printf("%s%s✔ Removed virtual host: %s%s\n", constants.ColorGreen, dryPrefix, domain, constants.ColorReset)
	fmt.Printf("  • Unlinked from sites-enabled and removed from sites-available\n")
	fmt.Printf("  • Deleted record from sites.db\n")
	fmt.Printf("  • Validated Nginx configuration and reloaded\n")
}

func executeNginxRm(domain string, opts VHostOptions) error {
	if err := RemoveVHost(domain, opts); err != nil {
		return err
	}

	deleteSiteFromDB(domain, opts.IsDryRun)
	if err := ReloadNginx(opts); err != nil {
		return err
	}

	printRmSuccess(domain, opts.IsDryRun)

	return nil
}

// runNginxRm removes a virtual host configuration and unregisters it from sites.db.
func runNginxRm(args []string) error {
	domain, hasDomain := parseRmDomain(args)
	if !hasDomain {
		return apperror.NewValidationError("domain name required: gitmap nginx rm <domain>")
	}

	opts := VHostOptions{IsDryRun: parseDryRunOption(args)}

	return executeNginxRm(domain, opts)
}
