package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TelegramWebhookController struct{}

func (t TelegramWebhookController) Post(c *gin.Context) {
	zap.S().Debug("Got a telegram webhook!")
}
