package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sturdy-dev/marblezero/achivements"
	"github.com/sturdy-dev/marblezero/cats"
	"github.com/sturdy-dev/marblezero/ingest"
	"github.com/sturdy-dev/marblezero/shells"
	"github.com/sturdy-dev/marblezero/state"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	defaultText = lipgloss.AdaptiveColor{Light: "#262626", Dark: "#FAFAFA"}

	orange = lipgloss.Color("#f97316")
	yellow = lipgloss.Color("#f9cf16")

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			Background(orange).
			Width(31).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			Background(orange).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Background(orange).
			Width(29).
			Render(s)
	}

	speechBubble = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())

	inScreenStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(orange)

	catStyle         = inScreenStyle.Copy().Bold(true).Height(11).Width(24)              // left side of screen
	deviceRightStyle = inScreenStyle.Copy().Height(11).Width(32).Foreground(defaultText) // right side of screen
)

var (
	flagPreexec        = flag.String("import-single", "", "Import a single execution. To be used with shell pre/post-exec hooks")
	flagFish           = flag.Bool("fish", false, "Print shell integration for the fish shell")
	flagZsh            = flag.Bool("zsh", false, "Print shell integration for the zsh shell")
	flagDebugColorMode = flag.Bool("debug-colors", false, "Debug layout")
)

func main() {
	flag.Parse()

	if *flagFish {
		fmt.Println(shells.Fish)
		return
	} else if *flagZsh {
		fmt.Println(shells.Zsh)
		return
	}

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
		if err := ingest.Single(storagePath, *flagPreexec); err != nil {
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

func output(config *state.Config, events []achivements.HistoryEvent) {

	// Set debug colors
	if *flagDebugColorMode {
		catStyle.Background(lipgloss.Color("#4d7c0f"))
		deviceRightStyle.Background(lipgloss.Color("#b91c1c"))
	}

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
	ListAllAchivementsScreen
	HelpScreen
)

type model struct {
	screen    Screen
	frame     int
	textInput textinput.Model
	config    *state.Config

	events               []achivements.HistoryEvent
	completedAchivements []achivements.Achivement

	rightScreenModel tea.Model
}

func NewModel(config *state.Config, events []achivements.HistoryEvent) *model {
	ti := textinput.New()
	ti.Placeholder = "Marble"
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

	// Calculate awarded achivements
	var completedAchivements []achivements.Achivement
	for _, a := range achivements.Achivements {
		if ok, at := a.Func(events); ok {
			a.AwardedAt = *at
			completedAchivements = append(completedAchivements, a)
		}
	}
	sort.Slice(completedAchivements, func(a, b int) bool {
		return completedAchivements[a].AwardedAt.After(completedAchivements[b].AwardedAt)
	})

	return &model{
		screen:               screen,
		config:               config,
		textInput:            ti,
		events:               events,
		completedAchivements: completedAchivements,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.characterAnimation())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	preScreen := m.screen

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.screen == SetupNameScreen {
				newName := strings.TrimSpace(m.textInput.Value())
				if len(newName) > 0 {
					m.config.Name = newName
					if err := m.config.Save(); err != nil {
						log.Println(err)
						return m, tea.Quit
					}
					m.screen = HomeScreen
					m.frame = 0 // reset counter
				}
			} else if m.screen == ListAllAchivementsScreen {
				m.screen = HomeScreen // go back
			} else if m.screen == HelpScreen {
				m.screen = HomeScreen // go back
			} else {
				return m, tea.Quit
			}

		// Rename
		case "r":
			if m.screen == HomeScreen {
				m.screen = SetupNameScreen
			}

		// List all achivements
		case "a":
			if m.screen == HomeScreen {
				m.screen = ListAllAchivementsScreen
				m.rightScreenModel = NewShowAllAchivementsModel(m.events)
			}

		// show help
		case "?", "h":
			if m.screen == HomeScreen {
				m.screen = HelpScreen
				m.rightScreenModel = NewHelpModel()
			}

		// Quit program if on home, else go back
		case "q", "esc":
			if m.screen == HomeScreen {
				return m, tea.Quit
			} else if m.screen == ListAllAchivementsScreen {
				m.screen = HomeScreen // go back
			} else if m.screen == HelpScreen {
				m.screen = HomeScreen // go back
			}
		// Quit program
		case "ctrl+c":
			return m, tea.Quit
		}

	case characterAnimationMsg:
		m.frame++
		return m, m.characterAnimation()

	case goToHomeMsg:
		m.screen = HomeScreen
		m.rightScreenModel = nil
	}

	if m.rightScreenModel != nil {
		m.rightScreenModel.Update(msg)
	}

	var cmd tea.Cmd
	if preScreen == SetupNameScreen && m.screen == SetupNameScreen {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	doc := strings.Builder{}

	var cat string

	switch m.screen {
	case SetupNameScreen:
		cat = cats.CatCurious
	default:
		cats := []string{
			cats.CatNormalStraight,
			cats.CatNormalStraightRaisedTail,
			cats.CatNormalStraight,
			cats.CatNormalRight,
			cats.CatNormalStraight,
			cats.CatAmused,
			cats.CatNormalStraight,
			cats.CatNormalStraightFoldedLeftEar,
			cats.CatNormalStraight,
		}
		cat = cats[m.frame%len(cats)]
	}

	level := len(m.completedAchivements)/3 + 1
	xp := len(m.completedAchivements) * 13

	var deviceRight string
	switch m.screen {
	case HomeScreen:
		charStats := fmt.Sprintf("%s\nMood: Happy\nLevel: %d (%d XP)", m.config.Name, level, xp)

		latestAchivementHeader := inScreenStyle.Copy().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(yellow).
			Foreground(yellow).
			MarginTop(2).
			Width(29).
			Render("Latest Achivement")

		a := m.completedAchivements[0]

		achivementName := inScreenStyle.Copy().Bold(true).Render(a.Name)
		// achivementXP := inScreenStyle.Copy().Render(fmt.Sprintf(" (%d XP)", 25))
		achivementDescription := inScreenStyle.Copy().Foreground(subtle).Render(a.Description)

		latestAchivement := inScreenStyle.Copy().Width(29).Render(fmt.Sprintf("%s\n%s\n(%d XP)", achivementName, achivementDescription, 12))

		deviceRight = deviceRightStyle.Copy().PaddingLeft(3).Render(
			lipgloss.JoinVertical(lipgloss.Left, charStats, latestAchivementHeader, latestAchivement),
		)

	case SetupNameScreen:
		bubble := inScreenStyle.Copy().Padding(0).Height(0).Render(lipgloss.JoinHorizontal(lipgloss.Bottom, "<\n", speechBubble.Render("Meow! Meow!\nWhat's my name?")))
		deviceRight = deviceRightStyle.PaddingLeft(3).Render(lipgloss.JoinVertical(lipgloss.Left, "\n\n", bubble, m.textInput.View()))

	case ListAllAchivementsScreen, HelpScreen:
		deviceRight = m.rightScreenModel.View()
	}

	cols := lipgloss.JoinHorizontal(
		lipgloss.Center,
		catStyle.Render(cat),
		deviceRight,
	)

	var device = lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(yellow).
		BorderBackground(orange).
		Render(cols)

	var horizontals []string = []string{device}

	frame := lipgloss.JoinHorizontal(
		lipgloss.Top,
		horizontals...,
	)

	var all = lipgloss.NewStyle().Render(frame)

	doc.WriteString(all + "\n")

	return doc.String()
}

func (m model) characterAnimation() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return characterAnimationMsg(t)
	})
}

type characterAnimationMsg time.Time

type goToHomeMsg struct{}

func goToHomeCmd() tea.Msg {
	return goToHomeMsg{}
}
