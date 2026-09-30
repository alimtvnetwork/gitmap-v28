package cmdos

// buildOSTUITabs constructs all six categorized tabs for the dashboard.
func buildOSTUITabs() []OSTUITab {
	return []OSTUITab{
		buildTweaksTab(),
		buildAutoLoginTab(),
		buildDisplayDMTab(),
		buildDNSNetTab(),
		buildCleanStorageTab(),
		buildUpdateTab(),
	}
}

func countAllCheckedItems(tabs []OSTUITab) int {
	total := 0
	for _, tab := range tabs {
		total += countTabChecked(tab)
	}

	return total
}

func countTabChecked(tab OSTUITab) int {
	checked := 0
	for _, it := range tab.Items {
		if it.IsChecked {
			checked++
		}
	}

	return checked
}
