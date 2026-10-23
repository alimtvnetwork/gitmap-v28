package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// UrlCacheEntry models a cached HTTP/HTTPS web resource in the Split-DB.
type UrlCacheEntry struct {
	UrlCacheEntryId int64  `json:"urlCacheEntryId"`
	UrlHash         string `json:"urlHash"`
	RawUrl          string `json:"rawUrl"`
	UrlText         string `json:"urlText,omitempty"`
	HttpMethod      string `json:"httpMethod,omitempty"`
	StatusCode      int    `json:"statusCode,omitempty"`
	ETag            string `json:"etag,omitempty"`
	LastModified    string `json:"lastModified,omitempty"`
	ContentType     string `json:"contentType"`
	ContentLength   int64  `json:"contentLength"`
	PayloadBlob     []byte `json:"payloadBlob,omitempty"`
	CompressionType string `json:"compressionType,omitempty"`
	IsCompressed    bool   `json:"isCompressed"`
	HeadersJson     string `json:"headersJson,omitempty"`
	DurationMs      int64  `json:"durationMs,omitempty"`
	AccessCount     int64  `json:"accessCount"`
	CachedAt        int64  `json:"cachedAt,omitempty"`
	ExpiresAt       int64  `json:"expiresAt"`
	LastAccessedAt  int64  `json:"lastAccessedAt,omitempty"`
	IsActive        bool   `json:"isActive"`
	Description     string `json:"description,omitempty"`
	Notes           string `json:"notes,omitempty"`
	Comments        string `json:"comments,omitempty"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
}

// EffectiveURL returns RawUrl or UrlText whichever is populated.
func (e *UrlCacheEntry) EffectiveURL() string {
	if e.RawUrl != "" {
		return e.RawUrl
	}

	return e.UrlText
}

// UrlAccessLog records an audit log for URL access events.
type UrlAccessLog struct {
	UrlAccessLogId  int64  `json:"urlAccessLogId"`
	UrlCacheEntryId *int64 `json:"urlCacheEntryId,omitempty"`
	UrlHash         string `json:"urlHash"`
	RawUrl          string `json:"rawUrl"`
	IsCacheHit      bool   `json:"isCacheHit"`
	StatusCode      int    `json:"statusCode"`
	BytesServed     int64  `json:"bytesServed"`
	DurationMs      int64  `json:"durationMs"`
	CallerSource    string `json:"callerSource"`
	Notes           string `json:"notes,omitempty"`
	Comments        string `json:"comments,omitempty"`
	CreatedAt       int64  `json:"createdAt"`
}

// Canonical Result wrapper type aliases.
type (
	UrlCacheEntryWrap      = result.Wrap[UrlCacheEntry]
	UrlCacheEntrySliceWrap = result.WrapSlice[UrlCacheEntry]
	UrlAccessLogSliceWrap  = result.WrapSlice[UrlAccessLog]
)

// UrlCacheQueryParam defines parameters for filtering and paginating cached URLs.
type UrlCacheQueryParam struct {
	Limit        int
	Offset       int
	IsActiveOnly bool
}
