package cmdos

import "fmt"

func executeCheckedItems(tabs []OSTUITab, isDryRun bool) []OSTUIActionResult {
	var results []OSTUIActionResult
	for _, tab := range tabs {
		for _, item := range tab.Items {
			if item.IsChecked {
				res := executeSingleTUIItem(item, isDryRun)
				results = append(results, res)
			}
		}
	}

	return results
}

func executeSingleTUIItem(item OSTUIItem, isDryRun bool) OSTUIActionResult {
	if isDryRun {
		return OSTUIActionResult{
			Title:     item.Title,
			IsSuccess: true,
			Message:   fmt.Sprintf("[dry-run] simulated execution for %s", item.ID),
		}
	}

	if !item.IsSupported {
		return OSTUIActionResult{
			Title:     item.Title,
			IsSuccess: false,
			Message:   "action not supported on this platform",
		}
	}

	return dispatchTUIAction(item)
}
