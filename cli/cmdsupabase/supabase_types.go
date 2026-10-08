package cmdsupabase

import (
	"time"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusError    = "error"
)

// SupabaseProjectRecord represents an encrypted project entry in the SQLite vault.
type SupabaseProjectRecord struct {
	Alias         string    `json:"alias"`
	ProjectRef    string    `json:"project_ref"`
	ApiUrl        string    `json:"api_url"`
	AnonKeyEnc    string    `json:"anon_key_enc"`
	ServiceKeyEnc string    `json:"service_key_enc"`
	DbUrlEnc      string    `json:"db_url_enc"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SupabaseProjectDecrypted represents in-memory credentials after authorized decryption.
type SupabaseProjectDecrypted struct {
	Alias       string `json:"alias"`
	ProjectRef  string `json:"project_ref"`
	ApiUrl      string `json:"api_url"`
	AnonKey     string `json:"anon_key"`
	ServiceKey  string `json:"service_key"`
	DbUrl       string `json:"db_url"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// SupabaseProjectView represents safe, masked metadata for tabular CLI display.
type SupabaseProjectView struct {
	Alias            string `json:"alias"`
	ProjectRef       string `json:"project_ref"`
	ApiUrl           string `json:"api_url"`
	AnonKeyMasked    string `json:"anon_key_masked"`
	ServiceKeyMasked string `json:"service_key_masked"`
	HasDbUrl         bool   `json:"has_db_url"`
	Status           string `json:"status"`
	UpdatedAt        string `json:"updated_at"`
}

// SupabasePingResult captures connectivity probe telemetry.
type SupabasePingResult struct {
	Alias          string        `json:"alias"`
	ApiUrl         string        `json:"api_url"`
	HttpStatusCode int           `json:"http_status_code"`
	Latency        time.Duration `json:"latency"`
	IsSuccess      bool          `json:"is_success"`
	ErrorMessage   string        `json:"error_message,omitempty"`
}
