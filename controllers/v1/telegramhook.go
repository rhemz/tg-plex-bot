package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rhemz/tg-plex-bot/config"
	"github.com/rhemz/tg-plex-bot/util"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//var (
//	cfg = config.GetConfig()
//)

type TelegramWebhookController struct{}

func (t TelegramWebhookController) Post(c *gin.Context) {
	//

	zap.S().Info("Got a telegram webhook!")

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		// always 200 so tg doesn't start exponential backoff
		c.JSON(http.StatusOK, gin.H{"error binding telegram JSON payload": err.Error()})
		return
	}

	zap.S().Info("Update ID: ", update.UpdateID)
	err := handleUpdateEvent(update, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error handling request"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func handleUpdateEvent(update tgbotapi.Update, c *gin.Context) error {
	cfg := config.GetConfig()
	authUsers := cfg.GetStringSlice("telegram.allowedUsers")

	if update.ChannelPost != nil {
		// channel post
	} else if update.Message != nil {
		// DM
		if !(update.Message.IsCommand() && util.Contains(authUsers, update.Message.From.UserName)) {
			// not a command from an authorized user
			zap.S().Info("Discarding DM from @", update.Message.From.UserName, ": ", update.Message.Text)
			return nil
		}

		// it's a command from an authorized user
		zap.S().Info("Got command from authorized user")

		switch update.Message.Command() {
		case "disable":
			disableMessages(update.Message.CommandArguments(), update)
		case "transcode":
			// TODO: implement transcode toggle
		}
	}

	return nil
}

func disableMessages(args string, update tgbotapi.Update) {
	if len(args) == 0 {
		_ = util.SendMessageResponse("This command requires a duration argument", update)
		return
	}
	duration, err := time.ParseDuration(args)
	if err != nil {
		zap.S().Error("Could not parse duration from: ", args)
		_ = util.SendMessageResponse(fmt.Sprintf("Could not parse duration: %s", args), update)
		return
	}

	zap.S().Info("Disabling posting updates for ", duration)
	_ = util.SendMessageResponse(fmt.Sprintf("Disabling for %s", duration), update)
	config.TelegramMessagesDisabled = true
	time.AfterFunc(duration, func() {
		msg := fmt.Sprintf("Re-enabling telegram messaging!  Sending was disabled for %s", duration)
		zap.S().Info(msg)
		_ = util.SendMessageResponse(msg, update)
		config.TelegramMessagesDisabled = false
	})
}
