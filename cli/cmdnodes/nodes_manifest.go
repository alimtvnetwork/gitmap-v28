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
	clean := extractJSONPayload(raw)
	items := parseTypedManifest(clean)
	if len(items) > 0 {
		return items
	}

	return parseGenericManifest(clean)
}

func parseTypedManifest(clean string) []RemoteRepoManifestItem {
	var rawItems []rawRemoteManifestItem
	if err := json.Unmarshal([]byte(clean), &rawItems); err != nil || len(rawItems) == 0 {
		return nil
	}

	var items []RemoteRepoManifestItem
	for _, ri := range rawItems {
		item := convertRawManifestItem(ri)
		if item.RepoName != "" || item.Slug != "" || item.RemoteURL != "" {
			items = append(items, item)
		}
	}
	return items
}

func convertRawManifestItem(ri rawRemoteManifestItem) RemoteRepoManifestItem {
	name := firstNonEmpty(ri.RepoName, ri.RepoNameAlt, ri.Name)
	slug := ri.Slug
	if slug == "" {
		slug = name
	}
	if name == "" {
		name = slug
	}
	return RemoteRepoManifestItem{
		RepoName:  name,
		Slug:      slug,
		RemoteURL: firstNonEmpty(ri.RemoteURL, ri.RemoteURLAlt, ri.URL),
		Path:      firstNonEmpty(ri.Path, ri.PathAlt),
	}
}

func parseGenericManifest(clean string) []RemoteRepoManifestItem {
	var genericList []map[string]interface{}
	if err := json.Unmarshal([]byte(clean), &genericList); err != nil {
		genericList = extractEnvelopeMapSlice(clean)
	}

	var items []RemoteRepoManifestItem
	for _, m := range genericList {
		item := convertGenericMapItem(m)
		if item.RepoName != "" || item.Slug != "" || item.RemoteURL != "" {
			items = append(items, item)
		}
	}
	return items
}

func extractEnvelopeMapSlice(clean string) []map[string]interface{} {
	var env map[string]interface{}
	if errEnv := json.Unmarshal([]byte(clean), &env); errEnv != nil {
		return nil
	}
	rawSlice := extractRawSliceFromEnv(env)
	var list []map[string]interface{}
	for _, el := range rawSlice {
		m, ok := el.(map[string]interface{})
		if ok {
			list = append(list, m)
		}
	}
	return list
}

func extractRawSliceFromEnv(env map[string]interface{}) []interface{} {
	if dataSlice, ok := env["data"].([]interface{}); ok {
		return dataSlice
	}
	if itemSlice, ok := env["items"].([]interface{}); ok {
		return itemSlice
	}
	return nil
}

func convertGenericMapItem(m map[string]interface{}) RemoteRepoManifestItem {
	item := RemoteRepoManifestItem{}
	item.RepoName = extractFirstString(m, "repoName", "repo_name", "name")
	item.Slug = extractFirstString(m, "slug")
	item.RemoteURL = extractFirstString(m, "remoteUrl", "remote_url", "url")
	item.Path = extractFirstString(m, "path", "absolute_path")
	if item.Slug == "" {
		item.Slug = item.RepoName
	}
	if item.RepoName == "" {
		item.RepoName = item.Slug
	}
	return item
}

func extractFirstString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
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
