package cmdchromeprofile

import (
	"encoding/json"
	_ "modernc.org/sqlite"
)

func extractEmailFromPreferences(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	var pref struct {
		AccountInfo []struct {
			Email string `json:"email"`
		} `json:"account_info"`
		Signin struct {
			UserName string `json:"user_name"`
		} `json:"signin"`
		Google struct {
			Services struct {
				Username string `json:"username"`
			} `json:"services"`
		} `json:"google"`
	}

	if err := json.Unmarshal(raw, &pref); err != nil {
		return ""
	}

	if len(pref.AccountInfo) > 0 && pref.AccountInfo[0].Email != "" {
		return pref.AccountInfo[0].Email
	}

	if pref.Signin.UserName != "" {
		return pref.Signin.UserName
	}

	return pref.Google.Services.Username
}
