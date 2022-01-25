package main

import (
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/server"
	"go.uber.org/zap"
)

func main() {
	// replace the sytem logger
	logger, _ := zap.NewProduction()
	undoLoggerReplace := zap.ReplaceGlobals(logger)
	defer undoLoggerReplace()

	// init config
	config.Init("config")
	if config.GetConfig().Get("telegram.bot_id") == "" || config.GetConfig().Get("telegram.api_token") == "" {
		zap.S().Fatal("Must set bot ID and API token")
	}

	// start doing server-y things
	server.Start()
}
