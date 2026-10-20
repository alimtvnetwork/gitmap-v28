// Package cmdschedule manages scheduled tasks, privacy masking, and mailbox inspection.
package cmdschedule

import (
	"os"
	"strings"
)

// PrivacyMasker governs egress presentation sanitization for email addresses.
type PrivacyMasker struct {
	IsMasked bool
}

// NewPrivacyMasker creates a new masker instance with the specified masking state.
func NewPrivacyMasker(isMasked bool) *PrivacyMasker {
	return &PrivacyMasker{
		IsMasked: isMasked,
	}
}

// Mask sanitizes an email address according to the active masking configuration.
func (p *PrivacyMasker) Mask(email string) string {
	if !p.IsMasked {
		return email
	}
	return MaskEmail(email)
}

// MaskList sanitizes a list of email addresses according to the active masking configuration.
func (p *PrivacyMasker) MaskList(emails []string) []string {
	if !p.IsMasked {
		return emails
	}
	return MaskEmailList(emails)
}

// HasAt verifies that an email string contains exactly one '@' character.
func HasAt(email string) bool {
	return strings.Contains(email, "@") && strings.Count(email, "@") == 1
}

// MaskEmail masks username and domain name while preserving boundary characters.
func MaskEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	if !isValidEmailFormat(trimmed) {
		return "***"
	}

	parts := strings.Split(trimmed, "@")
	localPart := maskSegment(parts[0])
	domainPart := maskDomainLabels(parts[1])

	return localPart + "@" + domainPart
}

// MaskEmailList masks each email in the slice and returns the sanitized collection.
func MaskEmailList(emails []string) []string {
	if len(emails) == 0 {
		return []string{}
	}

	results := make([]string, len(emails))
	for i, e := range emails {
		results[i] = MaskEmail(e)
	}

	return results
}

// MaskString sanitizes sensitive credentials or tokens for safe terminal rendering.
func MaskString(s string) string {
	trimmed := strings.TrimSpace(s)
	n := len(trimmed)
	if n == 0 {
		return ""
	}
	if n <= 4 {
		return "****"
	}
	if n <= 8 {
		return string(trimmed[0]) + "****" + string(trimmed[n-1])
	}

	return trimmed[:2] + "****" + trimmed[n-2:]
}

// ResolveIsMasked determines whether masking is active based on flag and environment.
func ResolveIsMasked(showFullEmail bool) bool {
	if showFullEmail {
		return false
	}
	if os.Getenv("GITMAP_ALLOW_UNMASKED_PII") == "1" {
		return false
	}

	return true
}

func isValidEmailFormat(s string) bool {
	if s == "" || !HasAt(s) {
		return false
	}

	parts := strings.Split(s, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}

	labels := strings.Split(parts[1], ".")
	if labels[len(labels)-1] == "" {
		return false
	}

	return areLabelsNonEmpty(labels[:len(labels)-1])
}

func areLabelsNonEmpty(labels []string) bool {
	for _, l := range labels {
		if len(l) == 0 {
			return false
		}
	}

	return true
}

func maskSegment(s string) string {
	n := len(s)
	if n <= 1 {
		return s + "*"
	}
	if n == 2 {
		return string(s[0]) + "*"
	}
	if n == 3 {
		return string(s[0]) + "*" + string(s[2])
	}

	return string(s[0]) + "***" + string(s[n-1])
}

func maskDomainLabels(domain string) string {
	segments := strings.Split(domain, ".")
	lastIdx := len(segments) - 1
	tld := segments[lastIdx]

	maskedLabels := make([]string, 0, len(segments))
	for i := 0; i < lastIdx; i++ {
		maskedLabels = append(maskedLabels, maskSegment(segments[i]))
	}
	maskedLabels = append(maskedLabels, tld)

	return strings.Join(maskedLabels, ".")
}
