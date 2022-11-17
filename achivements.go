package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sturdy-dev/hackagotchi/achivements"
)

type showAllAchivementsModel struct {
	events []achivements.HistoryEvent
	pages  [][]achivements.Achivement
	page   int
}

func NewShowAllAchivementsModel(events []achivements.HistoryEvent) tea.Model {

	var pages [][]achivements.Achivement
	var page []achivements.Achivement

	const perPage = 8

	for _, a := range achivements.Achivements {
		page = append(page, a)
		if len(page) == perPage {
			pages = append(pages, page)
			page = []achivements.Achivement{}
		}
	}
	if len(page) > 0 {
		pages = append(pages, page)
	}

	return &showAllAchivementsModel{
		events: events,
		pages:  pages,
	}
}

func (m *showAllAchivementsModel) Init() tea.Cmd {
	return nil
}

func (m *showAllAchivementsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "right", "tab", "n", "j", "l":
			if m.page < len(m.pages)-1 {
				m.page++
			}
		case "up", "left", "p", "h", "k":
			if m.page > 0 {
				m.page--
			}
		case "home":
			m.page = 0
		case "q", "esc", "enter":
			return m, goToHomeCmd
		}
	}
	return m, nil
}

func (m *showAllAchivementsModel) View() string {
	var showAchivements []string = []string{
		listHeader("All Achivements"),
	}

	for _, a := range m.pages[m.page] {
		if ok, _ := a.Func(m.events); ok {
			showAchivements = append(showAchivements, listDone(a.Name))
		} else {
			showAchivements = append(showAchivements, listItem(a.Name))
		}
	}

	// align unstructions with bottom
	if len(m.pages[m.page]) < 8 {
		showAchivements = append(showAchivements, strings.Repeat("\n", 7-len(m.pages[m.page])))
	}

	showAchivements = append(showAchivements, lipgloss.NewStyle().Foreground(subtle).Render(fmt.Sprintf("Page %d/%d (n/p/q)", m.page+1, len(m.pages))))

	return deviceRightStyle.Copy().PaddingLeft(1).Render(lipgloss.JoinVertical(lipgloss.Left, showAchivements...))
}
