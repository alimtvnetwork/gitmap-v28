// Package cmd — chrome_observe.go: inspects active Chrome processes, open pages,
// live tabs (via CDP / session files), and profile status with JSON/YAML support.
package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	_ "modernc.org/sqlite"
)

type chromeTabInfo struct {
	Profile string `json:"profile" yaml:"profile"`
	Title   string `json:"title" yaml:"title"`
	URL     string `json:"url" yaml:"url"`
	Status  string `json:"status" yaml:"status"`
	PID     int    `json:"pid,omitempty" yaml:"pid,omitempty"`
}

type chromeObservationReport struct {
	IsRunning      bool            `json:"isRunning" yaml:"isRunning"`
	ProcessCount   int             `json:"processCount" yaml:"processCount"`
	ActiveProfiles []string        `json:"activeProfiles" yaml:"activeProfiles"`
	Tabs           []chromeTabInfo `json:"tabs" yaml:"tabs"`
	ReportedAt     string          `json:"reportedAt" yaml:"reportedAt"`
}

func runChromeObserve(args []string) error {
	checkHelp(constants.SubCmdChromeObserve, args)
	opts := parseExtensionFilterArgs(args)
	report := collectChromeObservation(opts.Profile, opts.IsAll)
	if opts.Format == constants.OutputJSON {
		return printJSON(report)
	}
	if opts.Format == constants.OutputYAML {
		return printYAML(report)
	}
	return printObservationTable(report)
}

func collectChromeObservation(profFilter string, isAll bool) chromeObservationReport {
	isRunning, _ := isChromeRunning(runtime.GOOS)
	procCount := 0
	if isRunning {
		procCount = 1
	}
	report := chromeObservationReport{
		IsRunning:    isRunning,
		ProcessCount: procCount,
		ReportedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	profiles := resolveTargetProfiles(profFilter, isAll)
	report.ActiveProfiles = profiles
	report.Tabs = collectTabsAcrossProfiles(profiles, report.IsRunning)
	return report
}

func collectTabsAcrossProfiles(profiles []string, isRunning bool) []chromeTabInfo {
	var allTabs []chromeTabInfo
	cdpTabs := fetchCDPTabs()
	if len(cdpTabs) > 0 {
		return cdpTabs
	}
	for _, prof := range profiles {
		tabs := extractProfileSessionTabs(prof, isRunning)
		allTabs = append(allTabs, tabs...)
	}
	return allTabs
}

type cdpTabItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

func fetchCDPTabs() []chromeTabInfo {
	client := http.Client{Timeout: 300 * time.Millisecond}
	resp, err := client.Get("http://127.0.0.1:9222/json/list")
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var rawTabs []cdpTabItem
	_ = json.Unmarshal(body, &rawTabs)
	return mapCDPTabs(rawTabs)
}

func mapCDPTabs(rawTabs []cdpTabItem) []chromeTabInfo {
	var tabs []chromeTabInfo
	for _, t := range rawTabs {
		if t.Type == "page" && t.URL != "" {
			tabs = append(tabs, chromeTabInfo{
				Profile: "Active",
				Title:   t.Title,
				URL:     t.URL,
				Status:  "live (cdp)",
			})
		}
	}
	return tabs
}

func extractProfileSessionTabs(profName string, isRunning bool) []chromeTabInfo {
	srcPath, hasDir := resolveChromeProfileDir(profName)
	if !hasDir {
		return nil
	}
	historyDB := filepath.Join(srcPath, "History")
	if _, err := os.Stat(historyDB); err == nil {
		return readRecentHistoryURLs(historyDB, profName, isRunning)
	}
	return nil
}

func readRecentHistoryURLs(dbPath, profName string, isRunning bool) []chromeTabInfo {
	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=ro")
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT url, title FROM urls ORDER BY last_visit_time DESC LIMIT 10")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanHistoryRows(rows, profName, isRunning)
}

func scanHistoryRows(rows *sql.Rows, profName string, isRunning bool) []chromeTabInfo {
	var tabs []chromeTabInfo
	status := "recent"
	if isRunning {
		status = "active"
	}
	for rows.Next() {
		var u, t string
		if rows.Scan(&u, &t) == nil && u != "" {
			tabs = append(tabs, chromeTabInfo{
				Profile: profName,
				Title:   t,
				URL:     u,
				Status:  status,
			})
		}
	}
	return tabs
}

func printObservationTable(report chromeObservationReport) error {
	runStatus := "\033[1;91m● stopped\033[0m"
	if report.IsRunning {
		runStatus = fmt.Sprintf("\033[1;92m● running\033[0m (%d process(es))", report.ProcessCount)
	}
	fmt.Printf("\n\033[1;96mChrome Browser Status:\033[0m %s\n", runStatus)
	fmt.Printf("Active Profiles: %s\n\n", strings.Join(report.ActiveProfiles, ", "))
	if len(report.Tabs) == 0 {
		fmt.Println("  (no active tabs or recent pages found)")
		return nil
	}
	fmt.Printf("%-14s %-10s %-32s %s\n", "PROFILE", "STATUS", "TITLE", "URL")
	fmt.Println(strings.Repeat("-", 100))
	for _, t := range report.Tabs {
		title := t.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		fmt.Printf("%-14s %-10s %-32s %s\n", t.Profile, t.Status, title, t.URL)
	}
	fmt.Printf("\nTotal: %d page(s)/tab(s) observed\n", len(report.Tabs))
	return nil
}
