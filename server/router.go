package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhemz/tg-plex-bot/controllers/v1"
)


func NewRouter() *gin.Engine {
	router := gin.New()
	// router.Use(gin.Logger())
	
	// health check
	router.GET("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "ok")
    })

    // plex webhook
	v1API := router.Group("v1")
	{
		plexHookGroup := v1API.Group("plexhook")
		{
			hook := new(v1.PlexWebhookController)
			plexHookGroup.POST("", hook.Post)
		}
	}

	return router
}
