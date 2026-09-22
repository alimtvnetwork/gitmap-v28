package cmdos

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"

	"golang.org/x/term"
)

func promptAutoLoginConfig() (AutoLoginConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	currentUser, _ := user.Current()
	defaultUser := resolveDefaultUsername(currentUser)
	defaultDomain := resolveDefaultDomain()

	u := promptText(reader, fmt.Sprintf("Username [%s]: ", defaultUser), defaultUser)
	d := promptText(reader, fmt.Sprintf("Domain [%s]: ", defaultDomain), defaultDomain)
	p := promptPassword("Password: ")

	return AutoLoginConfig{
		Username:  u,
		Domain:    d,
		Password:  p,
		IsEnabled: true,
	}, nil
}

func resolveDefaultUsername(u *user.User) string {
	if u != nil && len(u.Username) > 0 {
		parts := strings.Split(u.Username, "\\")
		return parts[len(parts)-1]
	}
	return os.Getenv("USERNAME")
}

func resolveDefaultDomain() string {
	domain := os.Getenv("USERDOMAIN")
	if len(domain) > 0 {
		return domain
	}
	return "."
}

func promptText(r *bufio.Reader, label, fallback string) string {
	fmt.Print(label)
	input, _ := r.ReadString('\n')
	cleaned := strings.TrimSpace(input)
	if len(cleaned) == 0 {
		return fallback
	}
	return cleaned
}

func promptPassword(label string) string {
	fmt.Print(label)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(bytePassword))
}
