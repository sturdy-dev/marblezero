package achivements

import "time"

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
