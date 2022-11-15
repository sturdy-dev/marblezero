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

func Single(storagePath state.StoragePath, cmd string) error {
	historyFilePath := path.Join(string(storagePath), "history_wal")

	fp, err := os.OpenFile(historyFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("failed to open wal: %w", err)
	}

	parts := strings.Split(cmd, " ")
	prog := parts[0]

	event := achivements.HistoryEvent{
		Cmd:     prog,
		At:      time.Now(),
		IsForce: strings.Contains(cmd, "--force"),
	}

	// find subcommand (first non flag)
	if prog == "git" ||
		prog == "npm" || prog == "yarn" || prog == "pnpm" || //	JS
		prog == "pip3" || prog == "pip" || // Python
		false {
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "-") {
				continue
			}
			event.SubCommand = p
			break
		}
	}

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
