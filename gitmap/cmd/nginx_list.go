package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func fetchDBSites() []store.SiteRecord {
	db, err := store.OpenSitesSplitDB()
	if err != nil {
		return nil
	}

	defer db.Close()

	sites, err := db.ListSites()
	if err != nil {
		return nil
	}

	return sites
}

func printDBSiteRow(s store.SiteRecord) {
	status := "disabled"
	if s.IsActive {
		status = "enabled"
	}

	fmt.Printf("  %-25s %-12s %-6d %-10s %s\n", s.Domain, s.SiteType, s.ListenPort, status, s.DocumentRoot)
}

func renderDBSitesTable(sites []store.SiteRecord) {
	fmt.Printf("\n  %-25s %-12s %-6s %-10s %s\n", "DOMAIN", "TYPE", "PORT", "STATUS", "DOCUMENT ROOT")
	fmt.Printf("  %-25s %-12s %-6s %-10s %s\n", "------", "----", "----", "------", "-------------")
	for _, s := range sites {
		printDBSiteRow(s)
	}

	fmt.Printf("\n  Total: %d sites tracked in sites.db\n\n", len(sites))
}

// runNginxList lists virtual hosts from sites.db, falling back to disk if empty.
func runNginxList(args []string) error {
	dbSites := fetchDBSites()
	if len(dbSites) > 0 {
		renderDBSitesTable(dbSites)

		return nil
	}

	return runVHostList(args)
}
