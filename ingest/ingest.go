package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/zegl/hackagotchi/achivements"
	"github.com/zegl/hackagotchi/state"
)

func parse(cmd string, ts time.Time) achivements.HistoryEvent {
	parts := strings.Split(cmd, " ")
	prog := parts[0]
	var subcommand string
	var exts []string

	// find subcommand (first non flag)
	if prog == "git" ||
		prog == "npm" || prog == "yarn" || prog == "pnpm" || //	JS
		prog == "pip3" || prog == "pip" || // Python
		false {
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "-") {
				continue
			}
			subcommand = p
			break
		}
	}

	for _, p := range parts {
		if last := strings.LastIndexByte(p, '.'); last > len(p)-4 && last < len(p)-1 {
			exts = append(exts, p[last+1:])
		}
	}

	return achivements.HistoryEvent{
		Cmd:            prog,
		At:             ts,
		IsForce:        strings.Contains(cmd, "--force"),
		SubCommand:     subcommand,
		FileExtensions: exts,
	}
}

func Single(storagePath state.StoragePath, cmd string) error {
	historyFilePath := path.Join(string(storagePath), "history_wal")

	fp, err := os.OpenFile(historyFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("failed to open wal: %w", err)
	}

	ts := time.Now()

	event := parse(cmd, ts)

	raw, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if _, err := fp.WriteString(string(raw) + "\n"); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	if err := fp.Close(); err != nil {
		return fmt.Errorf("failed to close file")
	}

	return nil
}
