package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zegl/hackagotchi/achivements"
	"github.com/zegl/hackagotchi/cats"
	"github.com/zegl/hackagotchi/state"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	orange = lipgloss.Color("#f97316")

	list = lipgloss.NewStyle().
		BorderForeground(subtle).
		MarginLeft(4)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}
)

var (
	flagPreexec = flag.String("import-single", "", "Import a single execution. To be used with shell pre/post-exec hooks")
)

func main() {
	flag.Parse()

	storagePath, err := state.NewStoragePath()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config, err := state.LoadConfig(storagePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if *flagPreexec != "" {
		if err := importSingle(storagePath, *flagPreexec); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	events, err := achivements.ParseHistory(storagePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Graphical app
	output(config, events)
}

func importSingle(storagePath state.StoragePath, cmd string) error {
	historyFilePath := path.Join(string(storagePath), "history_wal")

	fp, err := os.OpenFile(historyFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("failed to open wal: %w", err)
	}

	prog := strings.Split(cmd, " ")[0]

	if _, err := fp.WriteString(prog + "\n"); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	if err := fp.Close(); err != nil {
		return fmt.Errorf("failed to close file")
	}

	return nil
}

func output(config *state.Config, events []achivements.HistoryEvent) {
	p := tea.NewProgram(NewModel(config, events))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type Screen int

const (
	HomeScreen Screen = iota
	SetupNameScreen
)

type model struct {
	screen    Screen
	frame     int
	textInput textinput.Model
	config    *state.Config
	events    []achivements.HistoryEvent
}

func NewModel(config *state.Config, events []achivements.HistoryEvent) *model {
	ti := textinput.New()
	ti.Placeholder = "Nala"
	ti.Focus()
	ti.CharLimit = 12
	ti.Width = 12
	ti.BackgroundStyle = lipgloss.NewStyle().Background(orange)
	ti.PlaceholderStyle = lipgloss.NewStyle().Background(orange)
	ti.PromptStyle = lipgloss.NewStyle().Background(orange)
	ti.CursorStyle = lipgloss.NewStyle().Background(orange)
	ti.TextStyle = lipgloss.NewStyle().Background(orange).Foreground(lipgloss.Color("#FAFAFA")).Bold(true)

	screen := HomeScreen
	if config.Name == "" {
		screen = SetupNameScreen
	}

	return &model{
		screen:    screen,
		config:    config,
		textInput: ti,
		events:    events,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.characterAnimation())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.screen == SetupNameScreen {
				m.config.Name = m.textInput.Value()
				if err := m.config.Save(); err != nil {
					log.Println(err)
					return m, tea.Quit
				}
				m.screen = HomeScreen
				m.frame = 0 // reset counter

			} else {
				return m, tea.Quit
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case characterAnimationMsg:
		m.frame++
		return m, m.characterAnimation()
	}

	var cmd tea.Cmd

	if m.screen == SetupNameScreen {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var showAchivements []string = []string{
		listHeader("Achivements"),
	}

	for _, a := range achivements.Achivements {
		if a.Func(m.events) {
			showAchivements = append(showAchivements, listDone(a.Name))
		} else {
			showAchivements = append(showAchivements, listItem(a.Name))
		}
	}

	achivementsFrame := list.Copy().Width(52).Render(lipgloss.JoinVertical(lipgloss.Left, showAchivements...))

	doc := strings.Builder{}

	var inScreenStyle = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(orange).
		Height(11).
		Padding(1, 1).
		AlignVertical(lipgloss.Center)

	cats := []string{
		cats.CatDefault,
		cats.CatDefaultLookRight,
		cats.Cat3,
		cats.CatAmused,
	}
	cat := cats[m.frame%4]

	var deviceRight string
	if m.screen == HomeScreen {
		deviceRight = inScreenStyle.Copy().Align(lipgloss.Left).PaddingLeft(3).Width(32).Render(fmt.Sprintf("%s\nMood: Happy\nLevel: 3", m.config.Name))
	} else if m.screen == SetupNameScreen {
		deviceRight = inScreenStyle.Copy().Align(lipgloss.Left).PaddingLeft(3).Width(32).Render("Hey buddy! What's your name?\n" + m.textInput.View())
	}

	cols := lipgloss.JoinHorizontal(
		lipgloss.Top,
		inScreenStyle.Copy().Bold(true).Width(24).Render(cat),
		deviceRight,
	)

	var device = lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(orange).Render(cols)

	var horizontals []string = []string{device}

	// Show achivements if setup
	if m.screen == HomeScreen {
		horizontals = append(horizontals, achivementsFrame)
	}

	frame := lipgloss.JoinHorizontal(
		lipgloss.Center,
		horizontals...,
	)

	var all = lipgloss.NewStyle().Render(frame)

	doc.WriteString(all + "\n\n")

	return doc.String()
}

func (m model) characterAnimation() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return characterAnimationMsg(t)
	})
}

type characterAnimationMsg time.Time
