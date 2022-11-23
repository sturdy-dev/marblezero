package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sturdy-dev/marblezero/achievements"
	"github.com/sturdy-dev/marblezero/state"
)

func parse(cmd string, ts time.Time) achievements.HistoryEvent {
	parts := strings.Split(cmd, " ")
	prog := parts[0]

	// find subcommand (first non flag)
	var subcommand string
	switch prog {
	case "git",
		// JS fanboys
		"npm", "yarn", "pnpm",
		// Python
		"pip3", "pip":
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "-") {
				continue
			}
			subcommand = p
			break
		}
	}

	var exts []string
	for _, p := range parts {
		if last := strings.LastIndexByte(p, '.'); last > len(p)-4 && last < len(p)-1 {
			exts = append(exts, p[last+1:])
		}
	}

	var flags []string
	switch prog {
	case "xcode-select", "rm", "git":
		for _, p := range parts {
			if strings.HasPrefix(p, "-") {
				flags = append(flags, p)
			}
		}
	}

	return achievements.HistoryEvent{
		Cmd:            prog,
		At:             ts,
		SubCommand:     subcommand,
		FileExtensions: exts,
		Flags:          flags,

		// Deprecated
		IsForce: strings.Contains(cmd, "--force"),
		IsRmRf:  strings.Contains(cmd, "-rf") || strings.Contains(cmd, "-fr"),
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
