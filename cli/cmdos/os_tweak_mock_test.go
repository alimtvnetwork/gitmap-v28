package cmdos

type mockTweakEngine struct {
	isClassicContext bool
	isClassicStart   bool
	isUltimatePower  bool
	isHibernate      bool
	isTelemetry      bool
	isActivity       bool
	isBing           bool
}

func (m *mockTweakEngine) SetContextMenu(isClassic bool) error {
	m.isClassicContext = isClassic
	return nil
}

func (m *mockTweakEngine) SetStartMenu(isClassic bool) error {
	m.isClassicStart = isClassic
	return nil
}

func (m *mockTweakEngine) SetPowerScheme(isUltimate bool) error {
	m.isUltimatePower = isUltimate
	return nil
}

func (m *mockTweakEngine) SetHibernate(isEnabled bool) error {
	m.isHibernate = isEnabled
	return nil
}

func (m *mockTweakEngine) SetTelemetry(isDisabled bool) error {
	m.isTelemetry = isDisabled
	return nil
}

func (m *mockTweakEngine) SetActivityFeed(isDisabled bool) error {
	m.isActivity = isDisabled
	return nil
}

func (m *mockTweakEngine) SetBingSearch(isDisabled bool) error {
	m.isBing = isDisabled
	return nil
}

func (m *mockTweakEngine) GetStatus() (TweakStatus, error) {
	return TweakStatus{
		IsClassicContextMenu:   m.isClassicContext,
		IsClassicStartMenu:     m.isClassicStart,
		IsUltimatePower:        m.isUltimatePower,
		IsHibernateEnabled:     m.isHibernate,
		IsTelemetryDisabled:    m.isTelemetry,
		IsActivityFeedDisabled: m.isActivity,
		IsBingSearchDisabled:   m.isBing,
	}, nil
}
