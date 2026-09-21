package cmdssh

import "strings"

func hasRawFlag(args []string) bool {
	for _, a := range args {
		isRaw := a == "--raw" || a == "-r"
		if isRaw {
			return true
		}
	}

	return false
}

func maskKeyBlob(blob string) string {
	isLong := len(blob) > 12
	if isLong {
		return blob[:12] + "...[redacted, pass --raw to view]..."
	}

	return "...[redacted, pass --raw to view]..."
}

func formatDisplayPublicKey(pubKey string, isRaw bool) string {
	if isRaw {
		return strings.TrimSpace(pubKey)
	}

	clean := strings.TrimSpace(pubKey)
	parts := strings.Fields(clean)
	hasParts := len(parts) >= 2
	if !hasParts {
		return maskKeyBlob(clean)
	}

	return formatMaskedKeyParts(parts)
}

func formatMaskedKeyParts(parts []string) string {
	keyType := parts[0]
	maskedBlob := maskKeyBlob(parts[1])
	hasComment := len(parts) >= 3
	if hasComment {
		comment := strings.Join(parts[2:], " ")

		return keyType + " " + maskedBlob + " " + comment
	}

	return keyType + " " + maskedBlob
}
