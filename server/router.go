package server

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/controllers/v1"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true)) // use zap logger
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))      // log panic to error

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// plex webhook
	v1API := router.Group("v1")
	{
		plexHookGroup := v1API.Group("plexhook")
		{
			hook := new(v1.PlexWebhookController)
			plexHookGroup.POST("", hook.Post)
		}

		telegramGroup := v1API.Group("telegram")
		{
			hook := new(v1.TelegramWebhookController)
			telegramGroup.POST(config.GetConfig().Get("telegram.hookUrl").(string), hook.Post)
		}
	}

	return router
}
