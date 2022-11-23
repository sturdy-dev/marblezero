package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sturdy-dev/marblezero/achievements"
)

type showAllAchievementsModel struct {
	events []achievements.HistoryEvent
	pages  [][]achievements.Achievement
	page   int
}

func NewShowAllAchievementsModel(events []achievements.HistoryEvent) tea.Model {

	var pages [][]achievements.Achievement
	var page []achievements.Achievement

	const perPage = 8

	for _, a := range achievements.Achievements {
		page = append(page, a)
		if len(page) == perPage {
			pages = append(pages, page)
			page = []achievements.Achievement{}
		}
	}
	if len(page) > 0 {
		pages = append(pages, page)
	}

	return &showAllAchievementsModel{
		events: events,
		pages:  pages,
	}
}

func (m *showAllAchievementsModel) Init() tea.Cmd {
	return nil
}

func (m *showAllAchievementsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *showAllAchievementsModel) View() string {
	var showAchievements []string = []string{
		listHeader("All Achievements"),
	}

	for _, a := range m.pages[m.page] {
		if ok, _ := a.Func(m.events); ok {
			showAchievements = append(showAchievements, listDone(a.Name))
		} else {
			showAchievements = append(showAchievements, listItem(a.Name))
		}
	}

	// align unstructions with bottom
	if len(m.pages[m.page]) < 8 {
		showAchievements = append(showAchievements, strings.Repeat("\n", 7-len(m.pages[m.page])))
	}

	showAchievements = append(showAchievements, lipgloss.NewStyle().Foreground(subtle).Render(fmt.Sprintf("Page %d/%d (n/p/q)", m.page+1, len(m.pages))))

	return deviceRightStyle.Copy().PaddingLeft(1).Render(lipgloss.JoinVertical(lipgloss.Left, showAchievements...))
}
