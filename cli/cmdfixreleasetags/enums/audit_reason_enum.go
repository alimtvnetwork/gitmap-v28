package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// AuditReasonType identifies the primary cause for flagging a release tag.
type AuditReasonType string

const (
	ReasonHealthy         AuditReasonType = "HEALTHY"
	ReasonOrphanTag       AuditReasonType = "ORPHAN_TAG"
	ReasonDraftRelease    AuditReasonType = "DRAFT_RELEASE"
	ReasonMissingAssets   AuditReasonType = "MISSING_ASSETS"
	ReasonCorruptedAssets AuditReasonType = "CORRUPTED_ASSETS"
	ReasonCICDFailed      AuditReasonType = "CICD_FAILED"
)

// IsValid reports whether the value matches a known variant.
func (a AuditReasonType) IsValid() bool {
	switch a {
	case ReasonHealthy, ReasonOrphanTag, ReasonDraftRelease, ReasonMissingAssets, ReasonCorruptedAssets, ReasonCICDFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (a AuditReasonType) String() string {
	return string(a)
}

// ParseAuditReasonType parses a raw string into a validated AuditReasonType.
func ParseAuditReasonType(s string) (AuditReasonType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllAuditReasonTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "reasonhealthy":
		return ReasonHealthy, nil
	case "healthy":
		return ReasonHealthy, nil
	case "reasonorphantag":
		return ReasonOrphanTag, nil
	case "orphantag":
		return ReasonOrphanTag, nil
	case "reasondraftrelease":
		return ReasonDraftRelease, nil
	case "draftrelease":
		return ReasonDraftRelease, nil
	case "reasonmissingassets":
		return ReasonMissingAssets, nil
	case "missingassets":
		return ReasonMissingAssets, nil
	case "reasoncorruptedassets":
		return ReasonCorruptedAssets, nil
	case "corruptedassets":
		return ReasonCorruptedAssets, nil
	case "reasoncicdfailed":
		return ReasonCICDFailed, nil
	case "cicdfailed":
		return ReasonCICDFailed, nil
	}
	return "", apperror.NewValidation("ParseAuditReasonType", "E1000", fmt.Sprintf("invalid auditreasontype: %q", s)).WithContext("value", s)
}

// AllAuditReasonTypes returns all valid variants.
func AllAuditReasonTypes() []AuditReasonType {
	return []AuditReasonType{
		ReasonHealthy,
		ReasonOrphanTag,
		ReasonDraftRelease,
		ReasonMissingAssets,
		ReasonCorruptedAssets,
		ReasonCICDFailed,
	}
}

// IsHealthy reports whether the value is ReasonHealthy.
func (a AuditReasonType) IsHealthy() bool {
	return a == ReasonHealthy
}

// IsOrphanTag reports whether the value is ReasonOrphanTag.
func (a AuditReasonType) IsOrphanTag() bool {
	return a == ReasonOrphanTag
}

// IsDraftRelease reports whether the value is ReasonDraftRelease.
func (a AuditReasonType) IsDraftRelease() bool {
	return a == ReasonDraftRelease
}

// IsMissingAssets reports whether the value is ReasonMissingAssets.
func (a AuditReasonType) IsMissingAssets() bool {
	return a == ReasonMissingAssets
}

// IsCorruptedAssets reports whether the value is ReasonCorruptedAssets.
func (a AuditReasonType) IsCorruptedAssets() bool {
	return a == ReasonCorruptedAssets
}

// IsCICDFailed reports whether the value is ReasonCICDFailed.
func (a AuditReasonType) IsCICDFailed() bool {
	return a == ReasonCICDFailed
}

// HasReason reports whether the value is non-empty and valid.
func (a AuditReasonType) HasReason() bool {
	return a.IsValid()
}
