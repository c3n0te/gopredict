package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jmoiron/sqlx"
)

var ()

type Styles struct {
	Title  lipgloss.Style
	Status lipgloss.Style
	Box    lipgloss.Style
}

func newStyles(darkBG bool) Styles {
	lightDark := lipgloss.LightDark(darkBG)
	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),

		Status: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04B575"), lipgloss.Color("#04B575"))),

		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 3).
			MarginTop(1),
	}
}

type Model struct {
	DB       *sqlx.DB
	styles   Styles
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func InitialModel(db *sqlx.DB) Model {
	return Model{
		DB:       db,
		styles:   newStyles(true),
		choices:  []string{},
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "u":
			tleStr := UpdateSats()
			tles := ParseTLEs(tleStr)
			UpsertTLEs(m.DB, tles)

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", "space":
			_, exists := m.selected[m.cursor]
			if exists {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	s := m.styles.Title.Render("GoPredict Satellite Tracking App") + "\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = m.styles.Status.Render(">")
			s += fmt.Sprintf("%s %s\n", cursor, m.styles.Status.Render(choice))
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}
	s += "\nPress j/k or arrows to move • Space/Enter to select • u to update • q to quit\n"
	return tea.NewView(m.styles.Box.Render(s))
}
