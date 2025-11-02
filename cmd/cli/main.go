package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mathiasdonoso/harborw/internal/tui"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		log, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			panic(err)
		}
		defer log.Close()
		logger := slog.New((slog.NewTextHandler(log, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelDebug,
		})))
		slog.SetDefault(logger)
	}

	slog.Debug("Init program")

	model, err := tui.NewModel()
	if err != nil {
		panic(err)
	}
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
