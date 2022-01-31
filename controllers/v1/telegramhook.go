package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramWebhookController struct{}

func (t TelegramWebhookController) Post(c *gin.Context) {
	//cfg := config.GetConfig()

	zap.S().Info("Got a telegram webhook!")

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		// always 200 so tg doesn't start exponential backoff
		c.JSON(http.StatusOK, gin.H{"error binding telegram JSON payload": err.Error()})
		return
	}

	zap.S().Info("Update ID: ", update.UpdateID)

	if update.ChannelPost != nil {
		// channel post
	} else if update.Message != nil {
		// DM
		if update.Message.Entities != nil && update.Message.Entities[0].Type == "bot_command" {
			// it's a command
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
