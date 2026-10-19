package cmdnodes

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/pterm/pterm"
)

// RemoteRepoManifestItem models a remote repository discovered during the discovery handshake.
type RemoteRepoManifestItem struct {
	RepoName  string `json:"repoName"`
	Slug      string `json:"slug"`
	RemoteURL string `json:"remoteUrl"`
	Path      string `json:"path"`
}

type rawRemoteManifestItem struct {
	RepoName     string `json:"repoName"`
	RepoNameAlt  string `json:"repo_name"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	RemoteURL    string `json:"remoteUrl"`
	RemoteURLAlt string `json:"remote_url"`
	URL          string `json:"url"`
	Path         string `json:"path"`
	PathAlt      string `json:"absolute_path"`
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func resolveManifestItemIdentifier(item RemoteRepoManifestItem) string {
	if item.RemoteURL != "" {
		return item.RemoteURL
	}
	if item.Slug != "" {
		return item.Slug
	}
	return item.RepoName
}

func parseRemoteManifest(raw string) []RemoteRepoManifestItem {
	var rawItems []rawRemoteManifestItem
	clean := extractJSONPayload(raw)
	if err := json.Unmarshal([]byte(clean), &rawItems); err == nil && len(rawItems) > 0 {
		var items []RemoteRepoManifestItem
		for _, ri := range rawItems {
			name := firstNonEmpty(ri.RepoName, ri.RepoNameAlt, ri.Name)
			slug := ri.Slug
			if slug == "" {
				slug = name
			}
			if name == "" {
				name = slug
			}
			url := firstNonEmpty(ri.RemoteURL, ri.RemoteURLAlt, ri.URL)
			path := firstNonEmpty(ri.Path, ri.PathAlt)

			if name != "" || slug != "" || url != "" {
				items = append(items, RemoteRepoManifestItem{
					RepoName:  name,
					Slug:      slug,
					RemoteURL: url,
					Path:      path,
				})
			}
		}
		if len(items) > 0 {
			return items
		}
	}

	var genericList []map[string]interface{}
	if err := json.Unmarshal([]byte(clean), &genericList); err != nil {
		var env map[string]interface{}
		if errEnv := json.Unmarshal([]byte(clean), &env); errEnv == nil {
			if dataSlice, ok := env["data"].([]interface{}); ok {
				for _, el := range dataSlice {
					if m, ok := el.(map[string]interface{}); ok {
						genericList = append(genericList, m)
					}
				}
			} else if itemSlice, ok := env["items"].([]interface{}); ok {
				for _, el := range itemSlice {
					if m, ok := el.(map[string]interface{}); ok {
						genericList = append(genericList, m)
					}
				}
			}
		}
	}

	var items []RemoteRepoManifestItem
	for _, m := range genericList {
		item := RemoteRepoManifestItem{}
		if v, ok := m["repoName"].(string); ok {
			item.RepoName = v
		} else if v, ok := m["repo_name"].(string); ok {
			item.RepoName = v
		} else if v, ok := m["name"].(string); ok {
			item.RepoName = v
		}
		if v, ok := m["slug"].(string); ok {
			item.Slug = v
		}
		if v, ok := m["remoteUrl"].(string); ok {
			item.RemoteURL = v
		} else if v, ok := m["remote_url"].(string); ok {
			item.RemoteURL = v
		} else if v, ok := m["url"].(string); ok {
			item.RemoteURL = v
		}
		if v, ok := m["path"].(string); ok {
			item.Path = v
		} else if v, ok := m["absolute_path"].(string); ok {
			item.Path = v
		}
		if item.Slug == "" {
			item.Slug = item.RepoName
		}
		if item.RepoName == "" {
			item.RepoName = item.Slug
		}
		if item.RepoName != "" || item.Slug != "" || item.RemoteURL != "" {
			items = append(items, item)
		}
	}

	return items
}

func appendEmptyOnlineNodeResult(results []NodeSummaryResult, conn db.SSHConnection, isJSON bool) []NodeSummaryResult {
	if !isJSON {
		return results
	}
	return append(results, NodeSummaryResult{
		NodeAlias:       conn.Alias,
		Host:            conn.IPAddress,
		IsOnline:        true,
		IsLocalMaster:   false,
		UniqueRepoCount: 0,
		Account:         conn.Username,
		Instances:       1,
	})
}

func printNodeExecError(alias string, err error, isJSON bool) {
	if !isJSON {
		pterm.Error.Printf("Failed to execute on node %s: %v\n", alias, err)
	}
}

func parsePayloadOrRaw(out string) interface{} {
	cleanRem := extractJSONPayload(out)
	var remPayload interface{}
	if err := json.Unmarshal([]byte(cleanRem), &remPayload); err != nil {
		return out
	}
	return remPayload
}
