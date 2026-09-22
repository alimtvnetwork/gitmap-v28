//go:build windows

package cmdos

import "golang.org/x/sys/windows/registry"

func isTelemetryDisabled() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, dataCollectionKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("AllowTelemetry")

	return err == nil && val == 0
}

func isActivityFeedDisabled() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, systemPoliciesKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("PublishUserActivities")

	return err == nil && val == 0
}

func isBingSearchDisabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, searchPoliciesKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("BingSearchEnabled")

	return err == nil && val == 0
}
