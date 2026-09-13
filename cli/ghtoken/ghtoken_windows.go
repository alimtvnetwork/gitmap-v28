//go:build windows

package ghtoken

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

func tokenFromSystemRegistry() (string, SourceType, bool) {
	if tok, isDefined := queryUserRegistryEnv("GH_TOKEN"); isDefined {
		return tok, SourceWinRegistryUser, true
	}

	if tok, isDefined := queryUserRegistryEnv("GITHUB_TOKEN"); isDefined {
		return tok, SourceWinRegistryUser, true
	}

	if tok, isDefined := queryMachineRegistryEnv("GH_TOKEN"); isDefined {
		return tok, SourceWinRegistrySys, true
	}

	if tok, isDefined := queryMachineRegistryEnv("GITHUB_TOKEN"); isDefined {
		return tok, SourceWinRegistrySys, true
	}

	return "", SourceNone, false
}

func queryUserRegistryEnv(valName string) (string, bool) {
	return readRegistryEnvVar(registry.CURRENT_USER, `Environment`, valName)
}

func queryMachineRegistryEnv(valName string) (string, bool) {
	subKey := `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`

	return readRegistryEnvVar(registry.LOCAL_MACHINE, subKey, valName)
}

func readRegistryEnvVar(root registry.Key, subKey, valName string) (string, bool) {
	key, err := registry.OpenKey(root, subKey, registry.QUERY_VALUE)
	if err != nil {
		return "", false
	}

	defer key.Close()

	val, _, err := key.GetStringValue(valName)
	if err != nil {
		return "", false
	}

	trimmed := strings.TrimSpace(val)
	isDefined := len(trimmed) > 0

	return trimmed, isDefined
}
