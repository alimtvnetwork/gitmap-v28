//go:build windows

package cmdos

import (
	"golang.org/x/sys/windows/registry"
)

const (
	dataCollectionKey = `SOFTWARE\Policies\Microsoft\Windows\DataCollection`
	systemPoliciesKey = `SOFTWARE\Policies\Microsoft\Windows\System`
	searchPoliciesKey = `Software\Microsoft\Windows\CurrentVersion\Search`
)

func positivePolicyValue(isPrivacyEnhanced bool) uint32 {
	if isPrivacyEnhanced {
		return 0
	}
	return 1
}

func applyTelemetry(isPrivacyEnhanced bool) error {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, dataCollectionKey, registry.SET_VALUE)
	if err != nil {
		return nil
	}
	defer k.Close()

	_ = k.SetDWordValue("AllowTelemetry", positivePolicyValue(isPrivacyEnhanced))
	return nil
}

func applyActivityFeed(isPrivacyEnhanced bool) error {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, systemPoliciesKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	val := positivePolicyValue(isPrivacyEnhanced)
	_ = k.SetDWordValue("PublishUserActivities", val)
	_ = k.SetDWordValue("UploadUserActivities", val)
	return nil
}

func applyBingSearch(isPrivacyEnhanced bool) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, searchPoliciesKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	val := positivePolicyValue(isPrivacyEnhanced)
	_ = k.SetDWordValue("BingSearchEnabled", val)
	_ = k.SetDWordValue("CortanaConsent", val)
	return nil
}
