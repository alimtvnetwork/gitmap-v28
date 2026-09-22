package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func handleDNSSet(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("usage: gitmap os dns set <cloudflare|google|quad9|adguard>", "E_MISSING_ARG")
	}

	return handleDNSDirectSet(args[0])
}

func handleDNSDirectSet(providerKey string) error {
	provider, isFound := KnownDNSProviders[strings.ToLower(providerKey)]
	if !isFound {
		msg := fmt.Sprintf("unknown DNS provider %q (valid: cloudflare, google, quad9, adguard)", providerKey)

		return apperror.NewSimple(msg, "E_INVALID_DNS_PROVIDER")
	}

	engine := newPlatformDNSEngine()
	iface, err := engine.GetDefaultInterface()
	if err != nil {
		return err
	}

	if err := engine.SetDNS(iface, provider); err != nil {
		return err
	}

	fmt.Printf("✔ DNS updated to %s (%s, %s) on interface %s\n", provider.Name, provider.Primary, provider.Secondary, iface)

	return nil
}
