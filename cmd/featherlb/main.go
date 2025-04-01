package main

import (
	"featherlb/internal/app/featherlb/log"
	"featherlb/internal/app/featherlb/server"
	"featherlb/internal/pkg/types"
	"flag"
	"log/slog"
)

func main() {
	configPath := flag.String("config", "", "Path to the config file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	if *configPath == "" {
		slog.Error("Config file path is required")
		return
	}

	log.ConfigureLogging(*debug)

	config, err := types.ReadConfigFromFile(*configPath)
	if err != nil {
		slog.Error("Failed to read config file", "error", err)
		return
	}

	slog.Debug("Config loaded", "location", *configPath, "config", config)

	server := server.NewFeatherLBServer()
	server.StartWithConfig(*config)
}
