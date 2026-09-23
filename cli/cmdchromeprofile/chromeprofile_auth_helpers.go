// Package cmd — chromeprofile_auth_helpers.go: extraction, encoding,
// and restoration helpers for Chrome GAIA credentials, Web Data, and Cookies.
package cmdchromeprofile

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type chromeGaiaInfo struct {
	GaiaID        string
	GaiaName      string
	GaiaGivenName string
	Email         string
}

type prefAccountInfoItem struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Gaia      string `json:"gaia"`
	GivenName string `json:"given_name"`
}

func fetchLocalStateGaiaInfo(dirName string) chromeGaiaInfo {
	state := readChromeLocalState()
	if state == nil {
		return chromeGaiaInfo{}
	}

	info, isFound := state.Profile.InfoCache[dirName]
	if !isFound {
		return chromeGaiaInfo{}
	}

	return chromeGaiaInfo{
		GaiaID:        info.GAIAID,
		GaiaName:      info.GAIAName,
		GaiaGivenName: info.GAIAGivenName,
		Email:         info.UserName,
	}
}

func extractGaiaFromPreferences(raw json.RawMessage) chromeGaiaInfo {
	if len(raw) == 0 {
		return chromeGaiaInfo{}
	}

	var pref struct {
		AccountInfo []prefAccountInfoItem `json:"account_info"`
	}
	if err := json.Unmarshal(raw, &pref); err != nil || len(pref.AccountInfo) == 0 {
		return chromeGaiaInfo{}
	}

	return buildGaiaInfoFromPrefAccount(pref.AccountInfo[0])
}

func buildGaiaInfoFromPrefAccount(acc prefAccountInfoItem) chromeGaiaInfo {
	gaiaID := acc.Gaia
	if gaiaID == "" {
		gaiaID = acc.AccountID
	}

	return chromeGaiaInfo{
		GaiaID:        gaiaID,
		GaiaName:      acc.FullName,
		GaiaGivenName: acc.GivenName,
		Email:         acc.Email,
	}
}

func resolveProfileGaiaInfo(dirName string, prefsRaw json.RawMessage, srcProfile string) chromeGaiaInfo {
	info := fetchLocalStateGaiaInfo(dirName)
	prefInfo := extractGaiaFromPreferences(prefsRaw)
	mergeGaiaInfo(&info, prefInfo)

	return info
}

func mergeGaiaInfo(target *chromeGaiaInfo, source chromeGaiaInfo) {
	if target.GaiaID == "" {
		target.GaiaID = source.GaiaID
	}
	if target.GaiaName == "" {
		target.GaiaName = source.GaiaName
	}
	if target.GaiaGivenName == "" {
		target.GaiaGivenName = source.GaiaGivenName
	}
	if target.Email == "" {
		target.Email = source.Email
	}
}

func readOptionalBase64(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil || len(raw) == 0 {
		return ""
	}

	return base64.StdEncoding.EncodeToString(raw)
}

func readProfileCookiesBase64(profileDir string) string {
	netCookies := filepath.Join(profileDir, "Network", "Cookies")
	if b64 := readOptionalBase64(netCookies); b64 != "" {
		return b64
	}

	return readOptionalBase64(filepath.Join(profileDir, "Cookies"))
}

func readProfileWebDataBase64(profileDir string) string {
	return readOptionalBase64(filepath.Join(profileDir, "Web Data"))
}

func restoreProfileCookies(profileDir, cookiesB64 string) error {
	if cookiesB64 == "" {
		return nil
	}

	raw, err := base64.StdEncoding.DecodeString(cookiesB64)
	if err != nil {
		return fmt.Errorf("decode cookies: %w", err)
	}

	netDir := filepath.Join(profileDir, "Network")
	if err := os.MkdirAll(netDir, constants.DirPermission); err != nil {
		return fmt.Errorf("mkdir %s: %w", netDir, err)
	}

	return os.WriteFile(filepath.Join(netDir, "Cookies"), raw, constants.FilePermission)
}

func restoreProfileWebData(profileDir, webDataB64 string) error {
	if webDataB64 == "" {
		return nil
	}

	raw, err := base64.StdEncoding.DecodeString(webDataB64)
	if err != nil {
		return fmt.Errorf("decode web data: %w", err)
	}

	return os.WriteFile(filepath.Join(profileDir, "Web Data"), raw, constants.FilePermission)
}
