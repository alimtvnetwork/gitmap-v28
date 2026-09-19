package cmdcargo

func extractInstallFlag(args []string) ([]string, bool) {
	var cleaned []string
	hasInstall := false
	for _, a := range args {
		if a == "--install" || a == "-i" {
			hasInstall = true

			continue
		}
		cleaned = append(cleaned, a)
	}

	return cleaned, hasInstall
}
