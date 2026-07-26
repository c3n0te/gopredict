package main

import (
	"fmt"
	"gopredict/api"
	"log/slog"
	"strings"
	"sync"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/akhenakh/sgp4"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type State int

const (
	ListView State = iota
	PassView
)

type styles struct {
	app           lipgloss.Style
	box           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

func newStyles(darkBG bool) styles {
	lightDark := lipgloss.LightDark(darkBG)

	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),

		box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 3).
			MarginBottom(1).
			MarginTop(1).
			Padding(1, 2),

		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),

		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04B575"), lipgloss.Color("#04B575"))),
	}
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	updateSats       key.Binding
	changeState      key.Binding
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		updateSats: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update tles"),
		),
		changeState: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select sat"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	db            *sqlx.DB
	state         State
	stations      []api.Station
	styles        styles
	darkBG        bool
	width, height int
	once          *sync.Once
	list          list.Model
	table         table.Model
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *model) updateListProperties() {
	// Update list size.
	h, v := m.styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.styles = newStyles(m.darkBG)
	m.list.Styles.Title = m.styles.title
}
func (m model) newPassTable(passes []sgp4.PassDetails) table.Model {
	columns := []table.Column{
		{Title: "MaxElev", Width: 15},
		{Title: "AOS", Width: 35},
		{Title: "LOS", Width: 35},
		{Title: "Duration", Width: 15},
	}

	rows := []table.Row{}
	for _, pass := range passes {
		row := table.Row{
			fmt.Sprintf("%f", pass.MaxElevation),
			pass.AOS.String(),
			pass.LOS.String(),
			pass.Duration.String(),
		}

		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(100),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)
	return t
}

func (m model) newList(tles []api.TLE) {
	numItems := len(m.list.Items())
	for i := range numItems {
		m.list.RemoveItem(i)
	}

	for i, tle := range tles {
		tleItem := item{
			title:       tle.SatName,
			description: "",
		}

		m.list.InsertItem(i, tleItem)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.updateSats):
			tleStr, err := UpdateSats()
			if err != nil || tleStr == "" {
				return m, nil
			}

			tles := ParseTLEs(tleStr)
			UpsertTLEs(m.db, tles)
			go m.newList(tles)
			return m, nil

		case key.Matches(msg, m.keys.changeState):
			switch m.state {
			case ListView:
				m.state = PassView
			case PassView:
				m.state = ListView
			}

			satname := m.list.SelectedItem().FilterValue()
			tle, err := ReadTLEBySatName(m.db, satname)
			if err != nil {
				slog.Error("Error reading TLE by satname: ", "error", err)
			}

			slog.Info(fmt.Sprintf("Retrieved tle from db: %v", tle))
			tleElems := []string{tle.SatName, tle.Line1, tle.Line2}
			tleStr := strings.Join(tleElems, "\n")
			tleSgp4, err := sgp4.ParseTLE(tleStr)
			startTime := time.Now().UTC()
			stopTime := startTime.Add(24 * time.Hour) // Predict for the next 24 hours
			stepSeconds := 30                         // Propagation step in seconds

			for _, stn := range m.stations {
				passes, err := tleSgp4.GeneratePasses(
					stn.Latitude,
					stn.Longitude,
					stn.Altitude,
					startTime,
					stopTime,
					stepSeconds,
				)

				if err != nil {
					slog.Error("Error generating passes: ", "error", err)
				}

				t := m.newPassTable(passes)
				m.table = t
			}

			return m, nil

		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		}
	}

	// This will also call our delegate's update function.
	switch m.state {
	case PassView:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)

	case ListView:
		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	v := tea.NewView("")

	switch m.state {
	case ListView:
		v = tea.NewView(m.styles.app.Render(m.list.View()))
		v.AltScreen = true

	case PassView:
		v = tea.NewView(m.styles.box.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
	}

	return v
}

func initialModel(db *sqlx.DB) model {
	// Initialize the model and list.
	m := model{}
	m.db = db
	m.styles = newStyles(false) // default to dark background styles
	m.state = ListView

	stns, err := ReadStations(m.db)
	if err != nil {
		return m
	}

	if len(stns) == 0 {
		stns, err = ParseStationFile()
		if err != nil {
			return m
		}

		err = UpsertStations(m.db, stns)
		if err != nil {
			return m
		}
	}

	m.stations = stns

	// Make initial list of items.
	items := []list.Item{}
	tles, err := ReadTLEs(m.db)
	if err != nil {
		return m
	}

	if len(tles) == 0 {
		tleStr, err := UpdateSats()
		if err != nil || tleStr == "" {
			return m
		}

		tles = ParseTLEs(tleStr)
		UpsertTLEs(m.db, tles)
	}

	for _, tle := range tles {
		tleItem := item{
			title:       tle.SatName,
			description: "",
		}

		items = append(items, tleItem)
	}

	t := m.newPassTable([]sgp4.PassDetails{})
	m.table = t

	// Setup list.
	delegateKeys := newDelegateKeyMap()
	listKeys := newListKeyMap()
	delegate := newItemDelegate(delegateKeys, &m.styles)
	satList := list.New(items, delegate, 0, 0)
	satList.Title = "GoPredict Satellite Tracking App"
	satList.Styles.Title = m.styles.title
	satList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.updateSats,
			listKeys.changeState,
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	m.list = satList
	m.keys = listKeys
	m.delegateKeys = delegateKeys

	return m
}
