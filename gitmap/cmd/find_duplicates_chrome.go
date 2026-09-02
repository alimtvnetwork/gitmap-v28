package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type chromeProfileDetail struct {
	ID         int64
	Name       string
	SourcePath string
}

func runFindDuplicatesChrome() error {
	mainDB, err := store.OpenDefault()
	if err != nil {
		fmt.Println("  " + constants.ColorYellow + "Gitmap database not available for Chrome profile check." + constants.ColorReset)
		return nil
	}
	defer mainDB.Close()

	profiles := queryChromeProfileDetails(mainDB)
	if len(profiles) == 0 {
		fmt.Println("  " + constants.ColorDim + "No Chrome profiles tracked in database." + constants.ColorReset)
		return nil
	}
	dupGroups := groupChromeDuplicates(profiles)
	if len(dupGroups) == 0 {
		fmt.Printf("  %s✓ Chrome: No duplicate profiles found. Total active: %d%s\n\n",
			constants.ColorGreen, len(profiles), constants.ColorReset)
		return nil
	}
	printChromeDupFindings(dupGroups)
	printChromeRemediations(dupGroups)
	return nil
}

func queryChromeProfileDetails(db *store.DB) []chromeProfileDetail {
	var out []chromeProfileDetail
	rows, err := db.Conn().Query("SELECT ChromeProfileId, Name, SourcePath FROM ChromeProfile")
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var d chromeProfileDetail
		if scanErr := rows.Scan(&d.ID, &d.Name, &d.SourcePath); scanErr == nil {
			out = append(out, d)
		}
	}
	return out
}

func groupChromeDuplicates(profiles []chromeProfileDetail) map[string][]chromeProfileDetail {
	groups := make(map[string][]chromeProfileDetail)
	for _, p := range profiles {
		key := strings.ToLower(filepath.Clean(p.SourcePath))
		if key == "" || key == "." {
			key = strings.ToLower(p.Name)
		}
		groups[key] = append(groups[key], p)
	}
	dupGroups := make(map[string][]chromeProfileDetail)
	for k, list := range groups {
		if len(list) > 1 {
			dupGroups[k] = list
		}
	}
	return dupGroups
}

func printChromeDupFindings(dupGroups map[string][]chromeProfileDetail) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Chrome Duplicate Profiles ──" + constants.ColorReset)
	totalDups := 0
	for _, list := range dupGroups {
		totalDups += len(list) - 1
	}
	fmt.Printf("  Found %s%d%s duplicate profile group(s) (%s%d%s duplicate entries total):\n\n",
		constants.ColorWhite, len(dupGroups), constants.ColorReset,
		constants.ColorYellow, totalDups, constants.ColorReset)

	for key, list := range dupGroups {
		fmt.Printf("  Key: %s%s%s (%d entries)\n", constants.ColorCyan, key, constants.ColorReset, len(list))
		fmt.Printf("    %-6s %-24s %s\n", "ID", "NAME", "SOURCE PATH")
		fmt.Println("    " + strings.Repeat("─", 68))
		for _, p := range list {
			fmt.Printf("    %-6d %-24s %s\n", p.ID, truncateStr(p.Name, 23), truncateStr(p.SourcePath, 34))
		}
		fmt.Println()
	}
}

func printChromeRemediations(dupGroups map[string][]chromeProfileDetail) {
	var sampleName string
	for _, list := range dupGroups {
		if len(list) > 1 {
			sampleName = list[1].Name
			break
		}
	}
	fmt.Println("  " + constants.ColorCyan + "Remediation & Fix Commands for Chrome Profiles:" + constants.ColorReset)
	fmt.Println("  " + strings.Repeat("─", 68))
	fmt.Printf("  ● Fix Single (Delete specific duplicate profile):\n")
	fmt.Printf("    %sgitmap chrome-profile delete \"%s\"%s\n", constants.ColorGreen, sampleName, constants.ColorReset)
	fmt.Printf("    %sgitmap chrome-profile clear --except \"%s\"%s\n\n", constants.ColorGreen, sampleName, constants.ColorReset)
	fmt.Printf("  ● Fix All Together (Deduplicate & merge profiles across database):\n")
	fmt.Printf("    %sgitmap chrome-profile optimize-projects%s   (alias: %sgitmap chrome-profile --repeat-fix%s)\n\n",
		constants.ColorGreen, constants.ColorReset, constants.ColorWhite, constants.ColorReset)
}
