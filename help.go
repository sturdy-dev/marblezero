package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type helpModel struct{}

func NewHelpModel() tea.Model {
	return &helpModel{}
}

func (m *helpModel) Init() tea.Cmd {
	return nil
}

func (m *helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "enter":
			return m, goToHomeCmd
		}
	}
	return m, nil
}

func (m *helpModel) View() string {
	var commands []string = []string{
		listHeader("Commands"),
		"a: show achievements",
		"r: rename your pet",
		"q / esc / enter / cmd+c: quit",
		"",
		lipgloss.NewStyle().Foreground(subtle).Render("(press enter to go back)"),
	}

	return deviceRightStyle.Copy().PaddingLeft(1).Render(lipgloss.JoinVertical(lipgloss.Left, commands...))
}
