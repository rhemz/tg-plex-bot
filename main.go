package main

import (
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/server"
	"github.com/rhemz/tg-plex-bot/util"
	"go.uber.org/zap"
)

func main() {
	// replace the sytem logger
	logger, _ := zap.NewProduction()
	undoLoggerReplace := zap.ReplaceGlobals(logger)
	defer undoLoggerReplace()

	// init config
	config.Init("config")
	requiredConf := []string{
		config.GetConfig().Get("telegram.botId").(string),
		config.GetConfig().Get("telegram.apiToken").(string),
		config.GetConfig().GetStringSlice("telegram.broadcastChannels")[0],
	}

	if util.Contains(requiredConf, "") {
		zap.S().Fatal("Must set bot ID, API token, and at least 1 broadcast channel ID")
	}

	// start doing server-y things
	server.Start()
}
