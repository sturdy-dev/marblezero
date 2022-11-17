package achivements

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/sturdy-dev/hackagotchi/state"
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

	FileExtensions []string `json:"file_extensions,omitempty"`
}

type ConditionFunc func(HistoryEvent) bool

type FilterFunc func([]HistoryEvent) []HistoryEvent

type AchivementFunc func(events []HistoryEvent) (awarded bool, at *time.Time)

const achivementNameMaxLength = 29

var (
	trueFunc AchivementFunc = func(events []HistoryEvent) (bool, *time.Time) {
		return true, &time.Time{}
	}

	withCommand = func(cmd string) ConditionFunc {
		return func(event HistoryEvent) bool {
			return event.Cmd == cmd
		}
	}

	withSubCommand = func(cmd, sub string) ConditionFunc {
		return func(event HistoryEvent) bool {
			return event.Cmd == cmd && event.SubCommand == sub
		}
	}

	withHourRange = func(min, max int) ConditionFunc {
		return func(event HistoryEvent) bool {
			return event.At.Hour() >= min && event.At.Hour() <= max
		}
	}

	withUniqueFileExtsMin = func(min int) ConditionFunc {
		return func(event HistoryEvent) bool {
			m := make(map[string]struct{})
			for _, e := range event.FileExtensions {
				m[e] = struct{}{}
			}
			return len(m) >= min
		}
	}

	withExts = func(exts ...string) ConditionFunc {
		return func(event HistoryEvent) bool {
			m := make(map[string]struct{})
			for _, e := range event.FileExtensions {
				m[e] = struct{}{}
			}
			for _, ext := range exts {
				if _, ok := m[ext]; !ok {
					return false
				}
			}
			return true
		}
	}

	times = func(n int) AchivementFunc {
		return func(events []HistoryEvent) (bool, *time.Time) {
			if len(events) < n {
				return false, nil
			}
			t := events[n]
			return true, &t.At
		}
	}

	and = func(filters ...ConditionFunc) FilterFunc {
		return func(events []HistoryEvent) []HistoryEvent {
			var res []HistoryEvent
		loopEvents:
			for _, e := range events {
				for _, f := range filters {
					if !f(e) {
						continue loopEvents
					}
				}
				res = append(res, e)
			}
			return res
		}
	}

	first = func(filters FilterFunc) AchivementFunc {
		return nth(filters, 0)
	}

	nth = func(filters FilterFunc, n int) AchivementFunc {
		return func(events []HistoryEvent) (bool, *time.Time) {
			filtered := filters(events)
			if len(filtered) <= n {
				return false, nil
			}
			return true, &filtered[n].At
		}
	}

	// Generally, the levels are
	// 1, 50, 250, 1000 times

	Achivements = []Achivement{
		{Name: "Name your pet", Func: trueFunc},

		// Deno
		{Name: "node << 2", Description: "Use deno", Func: first(and(withCommand("deno")))},

		// Go
		{Name: "Gopher", Description: "Use Go", Func: first(and(withCommand("go")))},
		{Name: "Go-go-gadget!", Description: "Use Go 50 times", Func: nth(and(withCommand("go")), 50)},
		{Name: "if err != nil", Description: "Use Go 250 times", Func: nth(and(withCommand("go")), 250)},
		{Name: "I love Rob", Description: "Use Go 1000 times", Func: nth(and(withCommand("go")), 1000)},

		// Rust
		{Name: "Getting Rusty", Description: "Use Cargo", Func: first(and(withCommand("cargo")))}, // TODO: allow cargo OR rustc?
		{Name: "No bugs to be seen here", Description: "Use Cargo 50 times", Func: nth(and(withCommand("cargo")), 50)},
		{Name: "Rewrite it in Rust", Description: "Use Cargo 250 times", Func: nth(and(withCommand("cargo")), 250)},
		// TODO: Rust 1000 times

		// Python
		{Name: "Import from __legacy__", Description: "Use Python2", Func: first(and(withCommand("python2")))},
		{Name: "Early adopter", Description: "Use Python3", Func: first(and(withCommand("python3")))},
		{Name: "Master of indentation", Description: "Use Python 50 times", Func: nth(and(withCommand("python")), 50)},
		{Name: "Parseltongue", Description: "Use Python 250 times", Func: nth(and(withCommand("python")), 250)},
		// TODO: Python 1000 times

		// Git
		{Name: "Teamwork makes the dream work", Description: "Use git", Func: first(and(withCommand("git")))},
		{Name: "Oncaller", Description: "Make a git commit in the middle of the night", Func: first(and(withSubCommand("git", "commit"), withHourRange(2, 5)))},
		{Name: "Use the --force", Description: "Use a git command with --force", Func: first(and(withCommand("git"), func(e HistoryEvent) bool { return e.IsForce }))},

		// Git commit streaks
		{Name: "Contributor", Description: "Make a git commit", Func: first(and(withSubCommand("git", "commit")))},
		{Name: "Developer", Description: "Make 50 git commits", Func: nth(and(withSubCommand("git", "commit")), 50)},
		{Name: "Coder", Description: "Make 250 git commits", Func: nth(and(withSubCommand("git", "commit")), 250)},
		{Name: "10xer", Description: "Make 1000 git commits", Func: nth(and(withSubCommand("git", "commit")), 1000)},

		// Polyglot
		{Name: "Polyglot", Description: "Add 3 files with different extensions to the git staging area", Func: first(and(withSubCommand("git", "add"), withUniqueFileExtsMin(3)))},
		{Name: "International Spy", Description: "Add 3 files with different extensions to the git staging area, 50 times", Func: nth(and(withSubCommand("git", "add"), withUniqueFileExtsMin(3)), 50)},

		// Editors
		{Name: "How do I exit this thing?", Description: "Edit a file with vim", Func: first(and(withCommand("vim")))},
		{Name: "M-x give-me-achivement", Description: "Edit a file with emacs", Func: first(and(withCommand("emacs")))},
		{Name: "Keeping it simple", Description: "Edit a file with nano", Func: first(and(withCommand("nano")))},

		// Shells
		{Name: "Show 'em whos boss", Description: "Use sudo", Func: first(and(withCommand("sudo")))},
		{Name: "Back to the past", Description: "Use sh", Func: first(and(withCommand("sh")))},
		{Name: "Gone fishin' 🐟", Description: "Use fish", Func: first(and(withCommand("fish")))}, // Alternative title: "90s kid"

		{Name: "Early bird", Description: "Use a command bewtween 05:00 and 07:00", Func: first(and(withHourRange(5, 7)))},
		{Name: "I love my cubicle", Description: "Use a command bewtween 07:00 and 17:00", Func: first(and(withHourRange(9, 17)))},
		{Name: "Night owl", Description: "Use a command bewtween 01:00 and 03:00", Func: first(and(withHourRange(1, 3)))},

		{Name: "Archivist", Description: "Use fzf to search the archives", Func: first(and(withCommand("fzf")))},
		{Name: "No backsies", Description: "Delete a directory with rm -rf", Func: first(and(withCommand("rm"), func(e HistoryEvent) bool { return e.IsRmRf }))},
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

func init() {
	for _, a := range Achivements {
		if len(a.Name) > achivementNameMaxLength {
			log.Printf("Achivement name (%s) is too long. Len=%d MaxAllowed=%d", a.Name, len(a.Name), achivementNameMaxLength)
		}
	}
}
