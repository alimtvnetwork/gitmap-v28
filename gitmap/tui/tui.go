package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

const viewCount = 9

// view indices.
const (
	viewBrowser      = 0
	viewActions      = 1
	viewGroups       = 2
	viewDashboard    = 3
	viewReleases     = 4
	viewTempReleases = 5
	viewZipGroups    = 6
	viewAliases      = 7
	viewLogs         = 8
)

// rootModel is the top-level Bubble Tea model.
type rootModel struct {
	db           *store.DB
	repos        []model.ScanRecord
	groups       []model.Group
	activeTab    int
	width        int
	height       int
	browser      browserModel
	actions      actionsModel
	groupsMgr    groupsModel
	dashboard    dashboardModel
	releases     releasesModel
	tempReleases tempReleasesModel
	zipGroups    zipGroupsModel
	aliases      aliasesModel
	logs         logsModel
	quitting     bool
}

// Run launches the interactive TUI.
func Run(db *store.DB, cfg model.Config) error {
	repos, err := db.ListRepos()
	if err != nil {
		return fmt.Errorf(constants.ErrTUIDBOpen, err)
	}

	groups, _ := db.ListGroups()

	m := newRootModel(db, repos, groups, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()

	return err
}

func newRootModel(db *store.DB, repos []model.ScanRecord, groups []model.Group, cfg model.Config) rootModel {
	return rootModel{
		db:           db,
		repos:        repos,
		groups:       groups,
		activeTab:    viewBrowser,
		browser:      newBrowserModel(repos),
		actions:      newActionsModel(),
		groupsMgr:    newGroupsModel(groups),
		dashboard:    newDashboardModel(repos, cfg.DashboardRefresh),
		releases:     newReleasesModel(db),
		tempReleases: newTempReleasesModel(db),
		zipGroups:    newZipGroupsModel(db),
		aliases:      newAliasesModel(db),
		logs:         newLogsModel(db),
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.dashboard.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case refreshMsg, tickMsg:
		dm, cmd := m.dashboard.Update(msg)
		m.dashboard = dm

		return m, cmd
	case tea.KeyMsg:
		if keys.quit(msg) && !m.browser.searching {
			m.quitting = true

			return m, tea.Quit
		}
		if keys.tab(msg) && !m.browser.searching {
			m.activeTab = (m.activeTab + 1) % viewCount

			return m, nil
		}
	}

	return m.updateActiveView(msg)
}

func (m rootModel) updateActiveView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.activeTab {
	case viewBrowser:
		bm, cmd := m.browser.Update(msg)
		m.browser = bm

		return m, cmd
	case viewActions:
		am, cmd := m.actions.Update(msg, m.browser.selected())
		m.actions = am

		return m, cmd
	case viewGroups:
		gm, cmd := m.groupsMgr.Update(msg)
		m.groupsMgr = gm

		return m, cmd
	case viewDashboard:
		dm, cmd := m.dashboard.Update(msg)
		m.dashboard = dm

		return m, cmd
	case viewReleases:
		rm, cmd := m.releases.Update(msg)
		m.releases = rm

		return m, cmd
	case viewTempReleases:
		tm, cmd := m.tempReleases.Update(msg)
		m.tempReleases = tm

		return m, cmd
	case viewZipGroups:
		zm, cmd := m.zipGroups.Update(msg)
		m.zipGroups = zm

		return m, cmd
	case viewAliases:
		am, cmd := m.aliases.Update(msg)
		m.aliases = am

		return m, cmd
	case viewLogs:
		lm, cmd := m.logs.Update(msg)
		m.logs = lm

		return m, cmd
	}

	return m, nil
}
