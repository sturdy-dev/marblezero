package achivements

import (
	"bufio"
	"encoding/json"
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
	Cmd        string    `json:"cmd"`
	SubCommand string    `json:"subcommand,omitempty"`
	At         time.Time `json:"at"`

	IsForce bool `json:"is_force,omitempty"`
	IsRmRf  bool `json:"is_rm_rf,omitempty"`
}

type AchivementFunc func(events []HistoryEvent) (awarded bool, at *time.Time)

type AchivementFilterFunc func(event HistoryEvent) bool

var (
	trueFunc AchivementFunc = func(events []HistoryEvent) (bool, *time.Time) {
		return true, &time.Time{}
	}

	usedCommandFunc = func(cmd string, times int) AchivementFunc {
		return func(events []HistoryEvent) (bool, *time.Time) {
			var c int
			for _, e := range events {
				if e.Cmd == cmd {
					c++
					if c >= times {
						return true, &e.At
					}
				}
			}
			return false, nil
		}
	}

	usedSubCommandFunc = func(cmd, subcommand string, times int) AchivementFunc {
		return func(events []HistoryEvent) (bool, *time.Time) {
			var c int
			for _, e := range events {
				if e.Cmd == cmd && e.SubCommand == subcommand {
					c++
					if c >= times {
						return true, &e.At
					}
				}
			}
			return false, nil
		}
	}

	withCommand = func(cmd string) AchivementFilterFunc {
		return func(event HistoryEvent) bool {
			return event.Cmd == cmd
		}
	}

	withSubCommand = func(cmd, sub string) AchivementFilterFunc {
		return func(event HistoryEvent) bool {
			return event.Cmd == cmd && event.SubCommand == sub
		}
	}

	withHourRange = func(min, max int) AchivementFilterFunc {
		return func(event HistoryEvent) bool {
			return event.At.Hour() >= min && event.At.Hour() <= max
		}
	}

	filter = func(filters ...AchivementFilterFunc) AchivementFunc {
		return func(events []HistoryEvent) (bool, *time.Time) {
		loopEvents:
			for _, e := range events {
				for _, f := range filters {
					if !f(e) {
						continue loopEvents
					}
				}
				return true, &e.At
			}
			return false, nil
		}
	}

	Achivements = []Achivement{
		{Name: "Name your pet", Func: trueFunc},

		// Deno
		{Name: "node << 2", Description: "Use deno", Func: usedCommandFunc("deno", 1)},

		// Go
		{Name: "Gopher", Description: "Use Go", Func: usedCommandFunc("go", 1)},
		{Name: "Go-go-gadget!", Description: "Use Go 10 times", Func: usedCommandFunc("go", 10)},
		{Name: "I love Rob", Description: "Use Go 1000 times", Func: usedCommandFunc("go", 1000)},

		// Rust
		{Name: "Getting Rusty", Description: "Use Cargo", Func: usedCommandFunc("cargo", 1)},
		{Name: "No bugs to be seen here", Description: "Use Cargo 50 times", Func: usedCommandFunc("cargo", 50)},
		{Name: "Rewrite it in Rust", Description: "Use Cargo 100 times", Func: usedCommandFunc("cargo", 100)},

		// Python
		{Name: "Import from __legacy__", Description: "Use Python2 10 times", Func: usedCommandFunc("python2", 10)},

		// Git
		{Name: "Teamwork makes the dream work", Description: "Use git", Func: usedCommandFunc("git", 1)},
		{Name: "Oncaller", Description: "Make a git commit in the middle of the night", Func: filter(withSubCommand("git", "commit"), withHourRange(2, 5))},
		{Name: "Use the --force", Description: "Use a git command with --force", Func: filter(withCommand("git"), func(e HistoryEvent) bool { return e.IsForce })},

		// Git commit streaks
		{Name: "Contributor", Description: "Make 10 git commits", Func: usedSubCommandFunc("git", "commit", 5)},
		{Name: "Developer", Description: "Make 25 git commits", Func: usedSubCommandFunc("git", "commit", 25)},
		{Name: "Coder", Description: "Make 100 git commits", Func: usedSubCommandFunc("git", "commit", 50)},
		{Name: "10xer", Description: "Make 500 git commits", Func: usedSubCommandFunc("git", "commit", 500)},

		// Editors
		{Name: "How do I exit this thing?", Description: "Edit a file with vim", Func: filter(withCommand("vim"))},
		{Name: "M-x give-me-achivement", Description: "Edit a file with emacs", Func: filter(withCommand("emacs"))},
		{Name: "Keeping it simple", Description: "Edit a file with nano", Func: filter(withCommand("nano"))},

		// Shells
		{Name: "Show 'em whos boss", Description: "Use sudo", Func: filter(withCommand("sudo"))},
		{Name: "Back to the past", Description: "Use sh", Func: filter(withCommand("sh"))},
		{Name: "Gone fishin' 🐟", Description: "Use fish", Func: filter(withCommand("fish"))},

		{Name: "Early bird", Description: "Use a command bewtween 05:00 and 07:00", Func: filter(withHourRange(5, 7))},
		{Name: "I love my cubicle", Description: "Use a command bewtween 07:00 and 17:00", Func: filter(withHourRange(9, 17))},
		{Name: "Night owl", Description: "Use a command bewtween 01:00 and 03:00", Func: filter(withHourRange(1, 3))},

		{Name: "Archivist", Description: "Use fzf to search the archives", Func: usedCommandFunc("fzf", 1)},
		{Name: "No looking back", Description: "Delete a directory with rm -rf", Func: filter(withCommand("rm"), func(e HistoryEvent) bool { return e.IsRmRf })},
	}
)

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
		var e HistoryEvent
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		events = append(events, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan wal: %w", err)
	}

	return events, nil
}
