// Package cmd — chrome_which.go: identify which Chrome profile is
// currently active by reading lockfile mtimes + Local State LastUsed.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func runChromeWhich(_ []string) {
	root := chromeUserDataDir()
	statePath := filepath.Join(root, "Local State")
	raw, err := os.ReadFile(statePath) //nolint:gosec
	if err != nil {
		fmt.Fprintf(os.Stderr, "chrome which: ERROR read Local State: %v\n", err)
		os.Exit(1)
	}
	var doc struct {
		Profile struct {
			LastUsed       string   `json:"last_used"`
			LastActiveProf []string `json:"last_active_profiles"`
			InfoCache      map[string]struct {
				Name string `json:"name"`
			} `json:"info_cache"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		fmt.Fprintf(os.Stderr, "chrome which: ERROR parse Local State: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[1;96mChrome User Data\033[0m %s\n", root)
	if doc.Profile.LastUsed != "" {
		name := doc.Profile.InfoCache[doc.Profile.LastUsed].Name
		fmt.Printf("\033[1;92mlast_used\033[0m       %s  \033[2;37m(display: %q)\033[0m\n", doc.Profile.LastUsed, name)
	}
	if len(doc.Profile.LastActiveProf) > 0 {
		fmt.Printf("\033[1;94mlast_active\033[0m\n")
		for _, d := range doc.Profile.LastActiveProf {
			fmt.Printf("  - %s  \033[2;37m(display: %q)\033[0m\n", d, doc.Profile.InfoCache[d].Name)
		}
	}
}
