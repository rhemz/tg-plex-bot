package server

import (
	"github.com/rhemz/tg-plex-bot/config"
	"go.uber.org/zap"
)

func Start() {
	r := NewRouter()

	bind := config.GetConfig().GetString("server.bindAddr")
	zap.S().Info("Attempting to listen on ", bind)

	// listen and serve
	defer r.Run(bind)
}
