package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zegl/hackagotchi/cats"

	_ "embed"
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

	if *flagPreexec != "" {
		if err := importSingle(*flagPreexec); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		output()
	}
}

func importSingle(cmd string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find config dir: %w", err)
	}

	hackagotchiDir := path.Join(homeDir, ".config", "hackagotchi")
	if err := os.MkdirAll(hackagotchiDir, 0777); err != nil {
		return fmt.Errorf("failed to create ~/.config/hackagotchi directory: %w", err)
	}

	historyFilePath := path.Join(hackagotchiDir, "history_wal")

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

func output() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	frame int
}

func NewModel() *model {
	return &model{}
}

func (m model) Init() tea.Cmd {
	return m.characterAnimation()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "q":
			return m, tea.Quit
		}
	case characterAnimationMsg:
		m.frame++
		return m, m.characterAnimation()
	}

	return m, nil
}

func (m model) View() string {
	achivements := list.Copy().Width(52).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			listHeader("Latest achivements"),
			listItem("Use Git"),
			listItem("Holy trinity (use npm, yarn, and pnpm in one week)"),
			listDone("Use the --force"),
			listItem("Use Slack"),
			listDone("Use Teams"),
		),
	)

	doc := strings.Builder{}

	var historyStyle = lipgloss.NewStyle().
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

	cols := lipgloss.JoinHorizontal(
		lipgloss.Top,
		historyStyle.Copy().Bold(true).Width(24).Render(cat),
		historyStyle.Copy().Align(lipgloss.Left).PaddingLeft(3).Render(fmt.Sprintf("zegl\nMood: Happy\nLevel: %d", m.frame)),
	)

	var device = lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(orange).Render(cols)

	frame := lipgloss.JoinHorizontal(
		lipgloss.Center,
		device,
		achivements,
	)

	var all = lipgloss.NewStyle().Render(frame)

	doc.WriteString(all)

	return doc.String()
}

func (m model) characterAnimation() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return characterAnimationMsg(t)
	})
}

type characterAnimationMsg time.Time
