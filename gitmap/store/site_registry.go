package store

import (
	"database/sql"
	"errors"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// SiteRecord represents a registered Nginx virtual host site in SQLite.
type SiteRecord struct {
	SiteRegistryId  int64  `json:"siteRegistryId"`
	Domain          string `json:"domain"`
	SiteType        string `json:"siteType"`
	DocumentRoot    string `json:"documentRoot"`
	NginxConfigPath string `json:"nginxConfigPath"`
	PhpVersion      string `json:"phpVersion,omitempty"`
	PhpSocketPath   string `json:"phpSocketPath,omitempty"`
	ListenPort      int    `json:"listenPort"`
	IsSslEnabled    bool   `json:"isSslEnabled"`
	IsActive        bool   `json:"isActive"`
	Description     string `json:"description,omitempty"`
	Notes           string `json:"notes,omitempty"`
	Comments        string `json:"comments,omitempty"`
	CreatedAt       string `json:"createdAt,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
}

const sqlUpsertSite = `INSERT INTO SiteRegistry (
    Domain, SiteType, DocumentRoot, NginxConfigPath,
    PhpVersion, PhpSocketPath, ListenPort, IsSslEnabled,
    IsActive, Description, Notes, Comments, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(Domain) DO UPDATE SET
    SiteType=excluded.SiteType,
    DocumentRoot=excluded.DocumentRoot,
    NginxConfigPath=excluded.NginxConfigPath,
    PhpVersion=excluded.PhpVersion,
    PhpSocketPath=excluded.PhpSocketPath,
    ListenPort=excluded.ListenPort,
    IsSslEnabled=excluded.IsSslEnabled,
    IsActive=excluded.IsActive,
    Description=excluded.Description,
    Notes=excluded.Notes,
    Comments=excluded.Comments,
    UpdatedAt=CURRENT_TIMESTAMP;`

// UpsertSite inserts or updates a virtual host site entry.
func (s *SitesSplitDB) UpsertSite(site SiteRecord) error {
	isSsl := 0
	if site.IsSslEnabled {
		isSsl = 1
	}
	isActive := 1
	if !site.IsActive && site.SiteRegistryId > 0 {
		isActive = 0
	}
	args := []any{
		site.Domain, site.SiteType, site.DocumentRoot, site.NginxConfigPath,
		site.PhpVersion, site.PhpSocketPath, site.ListenPort, isSsl,
		isActive, site.Description, site.Notes, site.Comments,
	}
	_, err := s.conn.Exec(sqlUpsertSite, args...)
	if err != nil {
		return apperror.WrapSimple(err, "sites_split.upsertSite")
	}

	return nil
}

func scanSiteRecord(scanner interface{ Scan(...any) error }) (*SiteRecord, error) {
	var r SiteRecord
	var isSsl, isActive int
	err := scanner.Scan(
		&r.SiteRegistryId, &r.Domain, &r.SiteType, &r.DocumentRoot,
		&r.NginxConfigPath, &r.PhpVersion, &r.PhpSocketPath, &r.ListenPort,
		&isSsl, &isActive, &r.Description, &r.Notes, &r.Comments,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.IsSslEnabled = isSsl == 1
	r.IsActive = isActive == 1

	return &r, nil
}

const sqlSelectSiteByDomain = `SELECT
    SiteRegistryId, Domain, SiteType, DocumentRoot, NginxConfigPath,
    COALESCE(PhpVersion, ''), COALESCE(PhpSocketPath, ''), ListenPort,
    IsSslEnabled, IsActive, COALESCE(Description, ''), COALESCE(Notes, ''),
    COALESCE(Comments, ''), CreatedAt, UpdatedAt
FROM SiteRegistry WHERE Domain = ?;`

// GetSite retrieves a virtual host site record by domain name.
func (s *SitesSplitDB) GetSite(domain string) (*SiteRecord, error) {
	row := s.conn.QueryRow(sqlSelectSiteByDomain, domain)
	rec, err := scanSiteRecord(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "sites_split.getSite")
	}

	return rec, nil
}

const sqlSelectAllSites = `SELECT
    SiteRegistryId, Domain, SiteType, DocumentRoot, NginxConfigPath,
    COALESCE(PhpVersion, ''), COALESCE(PhpSocketPath, ''), ListenPort,
    IsSslEnabled, IsActive, COALESCE(Description, ''), COALESCE(Notes, ''),
    COALESCE(Comments, ''), CreatedAt, UpdatedAt
FROM SiteRegistry ORDER BY Domain ASC;`

// ListSites retrieves all registered virtual host sites.
func (s *SitesSplitDB) ListSites() ([]SiteRecord, error) {
	rows, err := s.conn.Query(sqlSelectAllSites)
	if err != nil {
		return nil, apperror.WrapSimple(err, "sites_split.listSites")
	}
	defer rows.Close()

	return iterateSiteRows(rows)
}

func iterateSiteRows(rows *sql.Rows) ([]SiteRecord, error) {
	var results []SiteRecord
	for rows.Next() {
		rec, err := scanSiteRecord(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "sites_split.scanRow")
		}
		results = append(results, *rec)
	}

	return results, nil
}

const sqlDeleteSite = `DELETE FROM SiteRegistry WHERE Domain = ?;`

// DeleteSite removes a virtual host site record by domain name.
func (s *SitesSplitDB) DeleteSite(domain string) error {
	_, err := s.conn.Exec(sqlDeleteSite, domain)
	if err != nil {
		return apperror.WrapSimple(err, "sites_split.deleteSite")
	}

	return nil
}
