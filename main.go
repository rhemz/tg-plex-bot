package main

import (
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/server"
	"github.com/rhemz/tg-plex-bot/util"
	"go.uber.org/zap"
)

func main() {
	// replace the system logger
	logger, _ := zap.NewProduction()
	undoLoggerReplace := zap.ReplaceGlobals(logger)
	defer undoLoggerReplace()

	// init config
	config.Init("config")
	cfg := config.GetConfig()
	requiredConf := []string{
		cfg.Get("telegram.botId").(string),
		cfg.Get("telegram.apiToken").(string),
		cfg.GetStringSlice("telegram.broadcastChannels")[0],
	}
	if util.Contains(requiredConf, "") {
		zap.S().Fatal("Must set bot ID, API token, and at least 1 broadcast channel ID")
	}

	// init telegram bot api
	err := util.InitTelegramAPI(cfg.GetString("telegram.botId"), cfg.GetString("telegram.apiToken"))
	if err != nil {
		zap.S().Fatal("Error creating telegram api client: ", err)
	}
	zap.S().Info("Authorized telegram bot account: ", util.GetTelegramAPI().Self.UserName)

	// start doing server-y things
	server.Start()
}
