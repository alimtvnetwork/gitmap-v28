package cmdvhost

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ParseVHostSiteType normalizes raw strings into a VHostSiteType.
func ParseVHostSiteType(raw string) (VHostSiteType, *apperror.AppError) {
	norm := strings.ToLower(strings.TrimSpace(raw))
	switch norm {
	case "wordpress", "wp":

		return VHostSiteTypeWordpress, nil
	case "laravel", "art", "artisan":

		return VHostSiteTypeLaravel, nil
	case "php":

		return VHostSiteTypePHP, nil
	case "static", "html":

		return VHostSiteTypeStatic, nil
	default:

		return "", apperror.NewValidationError("unrecognized site type: " + raw)
	}
}

// ValidateVHostConfig validates required virtual host fields.
func ValidateVHostConfig(cfg VHostConfig) *apperror.AppError {
	isDomainEmpty := cfg.Domain == ""
	if isDomainEmpty {
		return apperror.NewValidationError("domain is required")
	}

	isRootEmpty := cfg.DocumentRoot == ""
	if isRootEmpty {
		return apperror.NewValidationError("document root is required")
	}

	isSiteTypeEmpty := cfg.SiteType == ""
	if isSiteTypeEmpty {
		return apperror.NewValidationError("site type is required")
	}

	return nil
}

func resolveVHostPort(port int, isSslEnabled bool) int {
	isCustomPort := port > 0
	if isCustomPort {
		return port
	}

	if isSslEnabled {
		return defaultVHostSslPort
	}

	return defaultVHostPort
}

func resolveVHostMaxBody(body string) string {
	hasBody := body != ""
	if hasBody {
		return body
	}

	return defaultVHostMaxBodySize
}

func resolveVHostTimeout(timeout int) int {
	hasTimeout := timeout > 0
	if hasTimeout {
		return timeout
	}

	return defaultVHostFastCGITimeout
}

func resolveVHostPass(pass string) string {
	hasPass := pass != ""
	if hasPass {
		return pass
	}

	return defaultVHostFastCGIPass
}

// ApplyVHostDefaults applies default ports, timeouts, and body sizes.
func ApplyVHostDefaults(cfg VHostConfig) VHostConfig {
	out := cfg
	out.Port = resolveVHostPort(cfg.Port, cfg.IsSslEnabled)
	out.MaxBodySize = resolveVHostMaxBody(cfg.MaxBodySize)
	out.FastCGITimeout = resolveVHostTimeout(cfg.FastCGITimeout)
	out.FastCGIPass = resolveVHostPass(cfg.FastCGIPass)

	return out
}

// SelectVHostTemplate resolves the raw template string for the given site type.
func SelectVHostTemplate(siteType VHostSiteType) (string, *apperror.AppError) {
	switch siteType {
	case VHostSiteTypeWordpress:

		return wordpressVHostTemplate, nil
	case VHostSiteTypeLaravel:

		return laravelVHostTemplate, nil
	case VHostSiteTypePHP:

		return phpVHostTemplate, nil
	case VHostSiteTypeStatic:

		return staticVHostTemplate, nil
	default:

		return "", apperror.NewValidationError("unknown site type: " + string(siteType))
	}
}

func executeVHostTemplate(tmplStr string, cfg VHostConfig) (string, *apperror.AppError) {
	tmpl, parseErr := template.New("vhost").Parse(tmplStr)
	if parseErr != nil {
		return "", apperror.WrapSimple(parseErr, "template.Parse")
	}

	var buf bytes.Buffer
	execErr := tmpl.Execute(&buf, cfg)
	if execErr != nil {
		return "", apperror.WrapSimple(execErr, "template.Execute")
	}

	return buf.String(), nil
}

// RenderVHostTemplate renders the Nginx configuration body without marker comments.
func RenderVHostTemplate(cfg VHostConfig) (string, *apperror.AppError) {
	valErr := ValidateVHostConfig(cfg)
	if valErr != nil {
		return "", valErr
	}

	applied := ApplyVHostDefaults(cfg)
	tmplStr, selectErr := SelectVHostTemplate(applied.SiteType)
	if selectErr != nil {
		return "", selectErr
	}

	return executeVHostTemplate(tmplStr, applied)
}

// BuildVHostMarkerTag generates the canonical marker tag for a virtual host.
func BuildVHostMarkerTag(siteType VHostSiteType, domain string) string {
	return fmt.Sprintf("vhost/%s/%s", siteType, domain)
}

// WrapVHostMarker encloses Nginx configuration text in idempotent marker blocks.
func WrapVHostMarker(siteType VHostSiteType, domain string, body string) string {
	tag := BuildVHostMarkerTag(siteType, domain)
	trimmed := strings.TrimSpace(body)

	return fmt.Sprintf("# >>> gitmap:%s >>>%s%s%s# <<< gitmap:%s <<<%s",
		tag, constants.NewLineUnix, trimmed, constants.NewLineUnix, tag, constants.NewLineUnix)
}

// RenderVHostConfig renders the virtual host template enclosed in GitMap marker comments.
func RenderVHostConfig(cfg VHostConfig) (string, *apperror.AppError) {
	body, err := RenderVHostTemplate(cfg)
	if err != nil {
		return "", err
	}

	wrapped := WrapVHostMarker(cfg.SiteType, cfg.Domain, body)

	return wrapped, nil
}
