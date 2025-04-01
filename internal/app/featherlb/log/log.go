package log

import (
	"log/slog"
	"os"
)

// ConfigureLogging sets up the logging configuration based on the debug flag.
func ConfigureLogging(debug bool) {
	var logger *slog.Logger
	if debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(logger)
}
