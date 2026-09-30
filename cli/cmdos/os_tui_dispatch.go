package cmdos

import "fmt"

func dispatchTUIAction(item OSTUIItem) OSTUIActionResult {
	switch {
	case isTweakAction(item.ID):
		return dispatchTweakAction(item.ID, item.Title)
	case isAutoLoginAction(item.ID):
		return dispatchAutoLoginAction(item.ID, item.Title)
	case isDisplayDMAction(item.ID):
		return dispatchDisplayDMAction(item.ID, item.Title)
	case isDNSAction(item.ID):
		return dispatchDNSAction(item.ID, item.Title)
	case isCleanAction(item.ID):
		return dispatchCleanAction(item.ID, item.Title)
	case isUpdateAction(item.ID):
		return dispatchUpdateAction(item.ID, item.Title)
	default:
		return OSTUIActionResult{
			Title:     item.Title,
			IsSuccess: false,
			Message:   fmt.Sprintf("unknown action id %s", item.ID),
		}
	}
}

func actionResultFromError(title string, err error, successMsg string) OSTUIActionResult {
	if err != nil {
		return OSTUIActionResult{
			Title:     title,
			IsSuccess: false,
			Message:   err.Error(),
		}
	}

	return OSTUIActionResult{
		Title:     title,
		IsSuccess: true,
		Message:   successMsg,
	}
}
