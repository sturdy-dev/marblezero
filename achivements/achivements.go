package achivements

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/zegl/hackagotchi/state"
)

type Achivement struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AwardedAt   time.Time `json:"awarded_at"`
	Func        AchivementFunc
}

type HistoryEvent struct {
	Cmd     string
	At      time.Time
	IsForce bool
}

type AchivementFunc func(events []HistoryEvent) bool

var (
	trueFunc AchivementFunc = func(events []HistoryEvent) bool {
		return true
	}

	falseFunc AchivementFunc = func(events []HistoryEvent) bool {
		return false
	}

	usedCommandFunc = func(cmd string) AchivementFunc {
		return func(events []HistoryEvent) bool {
			for _, e := range events {
				if e.Cmd == cmd {
					return true
				}
			}
			return false
		}
	}

	Achivements = []Achivement{
		{Name: "Name your pet", Func: trueFunc},

		// Use commands
		{Name: "node << 2", Description: "Use deno", Func: usedCommandFunc("deno")},
		{Name: "Gopher", Description: "Use Go", Func: usedCommandFunc("go")},
		{Name: "Silly Snake 🐍", Description: "Use Python 2", Func: usedCommandFunc("python2")},
		{Name: "Teamwork makes the dream work", Description: "Use Git", Func: usedCommandFunc("git")},
		{Name: "Archivist", Description: "Use fzf", Func: usedCommandFunc("fzf")},

		// TODO
		{Name: "Use the --force", Description: "Use a command with --force", Func: falseFunc},
	}
)

func init() {
	_ = falseFunc
}

func ParseHistory(storagePath state.StoragePath) ([]HistoryEvent, error) {
	file, err := os.Open(path.Join(string(storagePath), "history_wal"))
	if errors.Is(err, os.ErrNotExist) {
		return []HistoryEvent{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read wal: %w", err)
	}
	defer file.Close()

	var events []HistoryEvent

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		events = append(events, HistoryEvent{Cmd: scanner.Text()})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan wal: %w", err)
	}

	return events, nil
}
