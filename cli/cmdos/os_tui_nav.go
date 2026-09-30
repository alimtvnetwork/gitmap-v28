package cmdos

func handleTabSwitch(m OSTUIModel, delta int) OSTUIModel {
	n := len(m.Tabs)
	if n == 0 {
		return m
	}

	m.ActiveTab = (m.ActiveTab + delta + n) % n
	m.Cursor = 0

	return m
}

func handleJumpTab(m OSTUIModel, tabIdx int) OSTUIModel {
	if tabIdx < 0 || tabIdx >= len(m.Tabs) {
		return m
	}

	m.ActiveTab = tabIdx
	m.Cursor = 0

	return m
}

func handleCursorMove(m OSTUIModel, delta int) OSTUIModel {
	items := m.Tabs[m.ActiveTab].Items
	n := len(items)
	if n == 0 {
		return m
	}

	next := m.Cursor + delta
	if next < 0 {
		next = 0
	}
	if next >= n {
		next = n - 1
	}
	m.Cursor = next

	return m
}

func toggleCurrentItem(m OSTUIModel) OSTUIModel {
	items := m.Tabs[m.ActiveTab].Items
	if m.Cursor < 0 || m.Cursor >= len(items) {
		return m
	}

	items[m.Cursor].IsChecked = !items[m.Cursor].IsChecked
	m.Tabs[m.ActiveTab].Items = items

	return m
}

func setAllItemsInTab(m OSTUIModel, isChecked bool) OSTUIModel {
	items := m.Tabs[m.ActiveTab].Items
	for i := range items {
		items[i].IsChecked = isChecked
	}
	m.Tabs[m.ActiveTab].Items = items

	return m
}
