package ingest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/hackagotchi/achivements"
)

func TestParse(t *testing.T) {
	ts := time.Now()

	cases := []struct {
		cmd      string
		expected achivements.HistoryEvent
	}{
		{
			cmd: "git add foo.go bar.go yes.py lalalalalalalallaa .DS_Store hello. what.s",
			expected: achivements.HistoryEvent{
				Cmd:            "git",
				SubCommand:     "add",
				FileExtensions: []string{"go", "go", "py", "s"},
				At:             ts,
			},
		},
		{
			cmd: "git push --force",
			expected: achivements.HistoryEvent{
				Cmd:        "git",
				SubCommand: "push",
				IsForce:    true,
				At:         ts,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.cmd, func(t *testing.T) {
			event := parse(tc.cmd, ts)
			assert.Equal(t, tc.expected, event)
		})
	}

}
